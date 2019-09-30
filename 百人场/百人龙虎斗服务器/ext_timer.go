// 定时器回调函数，控制百人场状态
package main

import (
	// "encoding/json"
	"fmt"
	"logs"
	"time"

	"bl.com/util"
)

// 洗牌
func (this *ExtDesk) TimerShuffle(d interface{}) {
	logs.Debug("洗牌中")
	this.Lock()
	defer this.Unlock()

	this.ExtDeskInit()

	this.LeftCount = int32(this.CardMgr.GetLeftCardCount())
	this.RightCount = int32(this.CardMgr.GetSendCardCount())
	sd := GGameShuffleNotify{
		Id:         MSG_GAME_INFO_SHUFFLE_NOTIFY,
		Timer:      int32(gameConfig.Timer.ShuffleNum) * 1000,
		LeftCount:  this.LeftCount,
		RightCount: this.RightCount,
		GameCount:  this.Count,
	}

	this.BroadcastAll(MSG_GAME_INFO_SHUFFLE_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_SHUFFLE_NOTIFY
	this.AddTimer(gameConfig.Timer.Shuffle, gameConfig.Timer.ShuffleNum, this.TimerReady, nil)
}

// 准备
func (this *ExtDesk) TimerReady(d interface{}) {
	logs.Debug("准备中")
	this.Lock()
	defer this.Unlock()

	// 局数+1
	this.Count++

	this.HandleExit()
	this.HandleUndo()

	// 刷新座位
	this.UpdatePlayer()

	this.GameId = GetJuHao() // util.BuildGameId(GCONFIG.GameType)

	players := this.SeatMgr.GetUserList(len(this.Players))

	for _, v := range players {
		// 用户没有下注，增加未下注次数
		if v.(*ExtPlayer).GetMsgId() == 0 {
			fmt.Println("发现玩家本局未下注:", v.(*ExtPlayer).Nick, v.(*ExtPlayer).Robot)
			v.(*ExtPlayer).AddUndoTimes()
		}
		v.(*ExtPlayer).ResetAreaList()
	}

	this.LeftCount = int32(this.CardMgr.GetLeftCardCount())
	this.RightCount = int32(this.CardMgr.GetSendCardCount())
	sd := GGameReadyNotify{
		Id:         MSG_GAME_INFO_READY_NOTIFY,
		Timer:      int32(gameConfig.Timer.ReadyNum) * 1000,
		LeftCount:  this.LeftCount,
		RightCount: this.RightCount,
		GameCount:  this.Count,
		GameId:     this.GameId,
	}

	allUser := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range allUser {
		sd.SeatList = this.GetSeatInfo(v.(*ExtPlayer))
		v.(*ExtPlayer).SendNativeMsg(MSG_GAME_INFO_READY_NOTIFY, sd)
	}

	this.GameState = MSG_GAME_INFO_READY_NOTIFY
	this.AddTimer(gameConfig.Timer.Ready, gameConfig.Timer.ReadyNum, this.TimerSendCard, nil)
}

// 发牌
func (this *ExtDesk) TimerSendCard(d interface{}) {
	logs.Debug("发牌中")
	this.Lock()
	defer this.Unlock()

	this.OpenCard = this.CardMgr.SendOneCard()

	sd := GGameSendCardNotify{
		Id:       MSG_GAME_INFO_SEND_NOTIFY,
		Timer:    int32(gameConfig.Timer.SendCardNum) * 1000,
		OpenCard: this.OpenCard,
	}

	this.BroadcastAll(MSG_GAME_INFO_SEND_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_SEND_NOTIFY
	this.AddTimer(gameConfig.Timer.SendCard, gameConfig.Timer.SendCardNum, this.TimerBet, nil)
}

// 下注
func (this *ExtDesk) TimerBet(d interface{}) {
	logs.Debug("下注中")
	this.Lock()
	defer this.Unlock()

	if this.Count == int32(gameConfig.DeskInfo.BetLimit)+1 {
		this.BetArea[INDEX_DRAGONSPADE-1] = false
		this.BetArea[INDEX_DRAGONPLUM-1] = false
		this.BetArea[INDEX_DRAGONRED-1] = false
		this.BetArea[INDEX_DRAGONBLOCK-1] = false

		this.BetArea[INDEX_TIGERSPADE-1] = false
		this.BetArea[INDEX_TIGERPLUM-1] = false
		this.BetArea[INDEX_TIGERRED-1] = false
		this.BetArea[INDEX_TIGERBLOCK-1] = false
	}

	sd := GGameBetNotify{
		Id:      MSG_GAME_INFO_BET_NOTIFY,
		Timer:   int32(gameConfig.Timer.BetNum) * 1000,
		BetArea: this.BetArea,
	}

	this.BroadcastAll(MSG_GAME_INFO_BET_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_BET_NOTIFY
	this.AddTimer(gameConfig.Timer.Bet, gameConfig.Timer.BetNum, this.TimerStopBet, nil)

	// 定时广播
	this.AddTimer(gameConfig.Timer.NewBet, gameConfig.Timer.NewBetNum, this.HandleTimeOutBet, nil)
}

// 停止下注
func (this *ExtDesk) TimerStopBet(d interface{}) {
	logs.Debug("停止下注")
	this.Lock()
	defer this.Unlock()

	tAreaCoins := this.GetAreaCoinsList()

	// 消息公共部分
	sd := GGameStopBetNotify{
		Id:         MSG_GAME_INFO_STOP_BET_NOTIFY,
		Timer:      int32(gameConfig.Timer.StopBetNum) * 1000,
		TAreaCoins: tAreaCoins,
	}

	// 座位玩家下注信息
	// sd.SeatBetList = this.SeatMgr.GetSeatNewBetList()

	// 新下注，排除座位玩家
	// otherNewBetList := this.SeatMgr.GetOtherNewBetList()

	// 每个人的下注都不一样  需要单独处理
	players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range players {
		p := v.(*ExtPlayer)
		// newBetList := p.GetNewBetList()
		// p.ColAreaCoins()
		// TPAreaCoins := p.GetTotBetList()

		sd.PAreaCoins = p.GetNTAreaCoinsList()

		sd.OtherBetList = this.SeatMgr.GetOtherNewBetList2(p.Uid)
		// if this.SeatMgr.IsOnSeat(p) {
		// 	sd.OtherBetList = otherNewBetList
		// } else {
		// 	sd.OtherBetList = util.LessInt64List(otherNewBetList, newBetList)
		// }

		p.SendNativeMsg(MSG_GAME_INFO_STOP_BET_NOTIFY, sd)
	}

	for _, user := range this.Players {
		user.ColAreaCoins()
	}

	this.GameState = MSG_GAME_INFO_STOP_BET_NOTIFY
	this.AddTimer(gameConfig.Timer.StopBet, gameConfig.Timer.StopBetNum, this.TimerOpen, nil)

	//发牌
	this.DragonCard = this.CardMgr.SendOneCard()
	this.TigerCard = this.CardMgr.SendOneCard()

	if CD-CalPkAll(StartControlTime, time.Now().Unix()) < 0 && this.getPlayerBet() > 0 && GetCostType() == 1 {
		logs.Debug("进入风控")
		this.ControllerWinOrLose(true)
	} else if this.getPlayerBet() <= 0 {
		logs.Debug("胜率75")
		ra := RandInt64(100)
		if ra < 75 {
			this.ControllerWinOrLose(false)
		}
	}
	// 计算牌
	// this.BuildCards()
}
func (this *ExtDesk) ControllerWinOrLose(win bool) {
	chart := this.CardMgr.CompareCard(this.DragonCard, this.TigerCard)
	getCoins, _, _ := this.GameBeforeEnd(chart)
	if win {
		logs.Debug("需要赢")
		if getCoins < 0 {
			this.TigerCard, this.DragonCard = this.DragonCard, this.TigerCard
		}
	} else {
		logs.Debug("需要输")
		if getCoins > 0 {
			this.TigerCard, this.DragonCard = this.DragonCard, this.TigerCard
		}
	}
}

//获取真实玩家下注
func (this *ExtDesk) getPlayerBet() (allbet int64) {
	for _, v := range this.Players {
		if !v.Robot {
			for _, v1 := range v.BetList {
				allbet += v1
			}
		}
	}
	return
}

// 开牌
func (this *ExtDesk) TimerOpen(d interface{}) {
	logs.Debug("开牌")
	this.Lock()
	defer this.Unlock()

	sd := GGameOpenNotify{
		Id:         MSG_GAME_INFO_OPEN_NOTIFY,
		Timer:      int32(gameConfig.Timer.OpenNum) * 1000,
		DragonCard: this.DragonCard,
		TigerCard:  this.TigerCard,
	}

	this.BroadcastAll(MSG_GAME_INFO_OPEN_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_OPEN_NOTIFY
	this.AddTimer(gameConfig.Timer.Open, gameConfig.Timer.OpenNum, this.TimerAward, nil)
}

func (this *ExtDesk) BuildValid(winArea []bool, areaList []int64) int64 {
	// 同时押龙、虎，区域不计算有效打码
	if areaList[INDEX_DRAGON-1] > 0 && areaList[INDEX_TIGER-1] > 0 {
		areaList[INDEX_TIGER-1] = 0
		areaList[INDEX_DRAGON-1] = 0
	}

	// 同时押输、赢，区域不计算有效打码
	if areaList[INDEX_BANKERLOSE-1] > 0 && areaList[INDEX_BANKERWIN-1] > 0 {
		areaList[INDEX_BANKERWIN-1] = 0
		areaList[INDEX_BANKERLOSE-1] = 0
	}

	// 龙/虎花色下注不小于3，龙/虎区域不计算有效打码
	tigerCount := 0
	dragonCount := 0
	for i := 0; i < 4; i++ {
		if areaList[INDEX_TIGERSPADE+i-1] > 0 {
			tigerCount++
		}
		if areaList[INDEX_DRAGONSPADE+i-1] > 0 {
			dragonCount++
		}
	}
	if tigerCount > 2 {
		for i := 0; i < 4; i++ {
			areaList[INDEX_TIGERSPADE+i-1] = 0
		}
	}
	if dragonCount > 2 {
		for i := 0; i < 4; i++ {
			areaList[INDEX_DRAGONSPADE+i-1] = 0
		}
	}

	// 退还不计算打码
	if winArea[INDEX_BANKERLOSE-1] && winArea[INDEX_BANKERWIN-1] {
		areaList[INDEX_BANKERWIN-1] = 0
		areaList[INDEX_BANKERLOSE-1] = 0
	}

	var ret int64
	for _, v := range areaList {
		ret += v
	}

	return ret
}

// 派奖
func (this *ExtDesk) TimerAward(d interface{}) {
	logs.Debug("派奖中")
	this.Lock()
	defer this.Unlock()

	sd := GGameAwardNotify{
		Id:    MSG_GAME_INFO_AWARD_NOTIFY,
		Timer: int32(gameConfig.Timer.AwardNum) * 1000,
	}
	endTime := time.Now().Format("2006-01-02 15:04:05")
	var roomName string
	if GCONFIG.RoomType == 1 {
		roomName = "荣耀厅"
	} else if GCONFIG.RoomType == 2 {
		roomName = "王牌厅"
	} else if GCONFIG.RoomType == 3 {
		roomName = "战神厅"
	} else {
		roomName = "体验厅"
	}

	// 添加走势
	chart := this.CardMgr.CompareCard(this.DragonCard, this.TigerCard)
	this.RunChart = append(this.RunChart, int32(chart))
	//添加大厅走势
	// js, err := json.Marshal(this.RunChart)
	// if err != nil {
	// 	logs.Debug("添加大厅走势，生成json失败", err)
	// }
	// this.DeskMgr.SetZouShi(&GameZouShi{
	// 	SerId: int32(G_DbGetGameServerData.Sid),
	// 	GameInfo: GameTypeDetail{GameType: int32(GCONFIG.GameType),
	// 		RoomType:  int32(GCONFIG.RoomType),
	// 		GradeType: int32(GCONFIG.GradeType),
	// 	},
	// 	ZouShi:    string(js),
	// 	PlayerNum: int32(len(this.Players)),
	// 	UpdateT:   time.Now().Unix(),
	// })
	////////
	length := len(this.RunChart)
	if length > gameConfig.DeskInfo.RunChartCount {
		sd.RunChart = this.RunChart[length-gameConfig.DeskInfo.RunChartCount:]
	} else {
		sd.RunChart = this.RunChart
	}
	// 游戏结算，返回区域输赢情况
	_, WinArea, TWinArea := this.GameEnd(chart)
	for i, v := range this.BetArea {
		WinArea[i] = WinArea[i] && v
	}

	sd.WinArea = WinArea[:]
	sd.TWinArea = TWinArea[:]
	this.WinArea = WinArea[:]
	winarea := []int{}
	for i, v1 := range this.WinArea {
		if v1 {
			winarea = append(winarea, i)
		}
	}
	seatList := this.SeatMgr.GetSeatList()
	for _, seat := range seatList {
		win := seat.(*ExtPlayer).GetWinCoins()
		sd.SeatWinCoins = append(sd.SeatWinCoins, win)
		winArea := seat.(*ExtPlayer).GetWinList()
		TWinArea = util.LessInt64List(TWinArea, winArea)
	}

	sd.OtherWinArea = TWinArea

	//发送结算消息给数据库, 简单记录
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.GameId,
		Mini:        false,
		SetLeave:    1, //是否设置离开，0离开，1不离开
	}

	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.GameId,
		OpenCard:    this.OpenCard,
		DragonCard:  this.DragonCard,
		TigerCard:   this.TigerCard,
	}

	players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range players {
		//玩家游戏记录
		re := RecordData{
			MatchNum: this.JuHao,
			RoomName: roomName,
			EndTime:  endTime,
			Date:     time.Now().Unix(),
			LongCard: this.DragonCard,
			HuCard:   this.TigerCard,
		}
		re.WinArea = winarea
		p := v.(*ExtPlayer)
		sd.PWin = p.GetWinCoins()

		sd.PWinArea = p.GetWinList()

		areaList := p.GetTotBetList()
		valid := this.BuildValid(WinArea, areaList)

		betCoins := p.GetTotAreaCoins()

		// 结算
		sd.PrizeCoins = sd.PWin - betCoins
		logs.Debug("真是输赢:", sd.PrizeCoins)

		p.Award()
		sd.PCoins = p.GetCoins()

		if betCoins > 0 {
			dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
				UserId:      p.GetUid(),
				UserAccount: p.Account,
				BetCoins:    betCoins,
				ValidBet:    valid,
				PrizeCoins:  sd.PWin - betCoins,
				Robot:       p.Robot,
			})
			if GetCostType() == 1 && !p.Robot {
				p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
			}
			dbreq.UserCoin = []GGameEndInfo{}

			if p.Robot {
				continue
			}
			rddata := GGameRecordInfo{
				UserId:      p.GetUid(),
				UserAccount: p.Account,
				BetCoins:    betCoins,
				BetArea:     p.GetTotBetList(),
				PrizeCoins:  sd.PWin - betCoins,
				CoinsAfter:  sd.PCoins,
				Robot:       p.Robot,
			}
			rddata.CoinsBefore = rddata.CoinsAfter - rddata.PrizeCoins
			rdreq.UserRecord = append(rdreq.UserRecord, rddata)
			if GetCostType() == 1 {
				p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
			}
			rdreq.UserRecord = []GGameRecordInfo{}
		}
		for i, v1 := range p.GetNTAreaCoinsList() {
			re.AllBet += v1
			fmt.Println("区域：", i, "压住:", v1)
			if v1 > 0 {
				re.BetArea = append(re.BetArea, BetInfo{
					AreaIndex: i,
					BetCoins:  v1,
				})
			}
		}
		re.WinOrLost = sd.PrizeCoins
		if re.AllBet > 0 {
			logs.Debug("该玩家下注，将其记录存入！！")
			G_AllRecord.AddRecord(p.Uid, &re)
		} else {
			logs.Debug("该玩家没有下注，不存本局记录!!")
		}
		p.SendNativeMsg(MSG_GAME_INFO_AWARD_NOTIFY, sd)
	}

	this.GameState = MSG_GAME_INFO_AWARD_NOTIFY

	this.ResetAreaCoins()
	logs.Debug("Count:", this.CardMgr.GetLeftCardCount())
	logs.Debug("当前库存!!::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::lkfjdsklfjsadflskdafjs:", CD)
	// 剩余张数 < 3 结束，否则继续准备
	if this.CardMgr.GetLeftCardCount() < 3 {
		this.TimerOver(nil)
		return
	}

	this.AddTimer(gameConfig.Timer.Award, gameConfig.Timer.AwardNum, this.TimerReady, nil)
}

// 结束
func (this *ExtDesk) TimerOver(d interface{}) {
	logs.Debug("结束")
	this.AddTimer(gameConfig.Timer.Over, gameConfig.Timer.OverNum, this.TimerShuffle, nil)
}

// 结算
func (this *ExtDesk) GameEnd(chart byte) (int64, []bool, []int64) {
	logs.Debug("结算中")
	var winArea [13]bool
	var tWinArea [13]int64
	var double [13]float64
	for i := range winArea {
		winArea[i] = false
		double[i] = 0
		tWinArea[i] = 0
	}

	// 计算闲输赢
	var loseCoins int64 = 0
	var index int

	// 龙虎和区域中奖
	switch chart {
	case DRAGON:
		index = INDEX_DRAGON - 1
	case TIGER:
		index = INDEX_TIGER - 1
	case DRAW:
		index = INDEX_DRAW - 1
	default:
		index = INDEX_ERROR - 1
	}
	coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	loseCoins += int64(coins)
	double[index] = gameConfig.Double[index]
	winArea[index] = true

	// 花色区域中奖
	dColor := GetCardColor(this.DragonCard)
	switch dColor {
	case CARD_COLOR_Fang:
		index = INDEX_DRAGONSPADE - 1
	case CARD_COLOR_Mei:
		index = INDEX_DRAGONPLUM - 1
	case CARD_COLOR_Hong:
		index = INDEX_DRAGONRED - 1
	case CARD_COLOR_Hei:
		index = INDEX_DRAGONBLOCK - 1
	default:
		index = INDEX_ERROR - 1
	}
	coins = float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	loseCoins += int64(coins)
	double[index] = gameConfig.Double[index]
	winArea[index] = true

	tColor := GetCardColor(this.TigerCard)
	switch tColor {
	case CARD_COLOR_Fang:
		index = INDEX_TIGERSPADE - 1
	case CARD_COLOR_Mei:
		index = INDEX_TIGERPLUM - 1
	case CARD_COLOR_Hong:
		index = INDEX_TIGERRED - 1
	case CARD_COLOR_Hei:
		index = INDEX_TIGERBLOCK - 1
	default:
		index = INDEX_ERROR - 1
	}
	coins = float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	loseCoins += int64(coins)
	double[index] = gameConfig.Double[index]
	winArea[index] = true

	dMidT := this.GetAreaCoin(INDEX_DRAGON-1) - this.GetAreaCoin(INDEX_TIGER-1)
	if dMidT == 0 {
		// 龙虎押注一样，上庄区域退还
		index = INDEX_BANKERWIN - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true

		index = INDEX_BANKERLOSE - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true
	}

	switch chart {
	case DRAGON:
		if dMidT < 0 {
			index = INDEX_BANKERWIN - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else if dMidT > 0 {
			index = INDEX_BANKERLOSE - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		}
	case TIGER:
		if dMidT < 0 {
			index = INDEX_BANKERLOSE - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else if dMidT > 0 {
			index = INDEX_BANKERWIN - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		}
	case DRAW:
		// 返还上庄
		index = INDEX_BANKERWIN - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true

		index = INDEX_BANKERLOSE - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true
	}

	// 计算玩家输赢
	for i, v := range winArea {
		if !v {
			continue
		}

		coins := int64(float64(this.GetUserAreaCoin(i)) * double[i])
		tWinArea[i] = coins
	}

	players := this.SeatMgr.GetUserList(len(this.Players))

	//以下为属性修改
	for _, v := range players {
		p := v.(*ExtPlayer)
		p.BuildWinList(double[:])
		p.AddWinList()
		p.AddBetList()
	}

	totCoins := this.GetUserAreaCoins()
	//修改库存
	AddLocalStock(totCoins - loseCoins)
	AddCD(totCoins - loseCoins)

	// if gameConfig.DeskInfo.Win == 0 {
	// 	return totCoins - loseCoins, winArea[:], tWinArea[:]
	// }

	// // 有设置盈利率，需要计算盈利率
	// if this.GetUserAreaCoins() > 0 {
	// 	this.tCount++
	// 	this.totCoins += float64(totCoins)
	// 	this.wCoins += float64(totCoins - loseCoins)
	// 	// logs.Debug("测试次数，总下注，总赢取：", this.tCount, this.totCoins, this.wCoins)
	// }

	// rate := this.wCoins / this.totCoins * 100
	// if int(rate) < gameConfig.DeskInfo.Win+5 && int(rate) > gameConfig.DeskInfo.Win-5 {
	// 	this.wCoins = 0
	// 	this.totCoins = 0
	// 	this.tCount = 0
	// }

	return totCoins - loseCoins, winArea[:], tWinArea[:]
}

// 计算开奖结果
func (this *ExtDesk) GameBeforeEnd(chart byte) (int64, []bool, []int64) {
	logs.Debug("结算中")
	var winArea [13]bool
	var tWinArea [13]int64
	var double [13]float64
	for i := range winArea {
		winArea[i] = false
		double[i] = 0
		tWinArea[i] = 0
	}

	// 计算闲输赢
	var loseCoins int64 = 0
	var index int

	// 龙虎和区域中奖
	switch chart {
	case DRAGON:
		index = INDEX_DRAGON - 1
	case TIGER:
		index = INDEX_TIGER - 1
	case DRAW:
		index = INDEX_DRAW - 1
	default:
		index = INDEX_ERROR - 1
	}
	coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	loseCoins += int64(coins)
	double[index] = gameConfig.Double[index]
	winArea[index] = true

	// 花色区域中奖
	dColor := GetCardColor(this.DragonCard)
	switch dColor {
	case CARD_COLOR_Fang:
		index = INDEX_DRAGONSPADE - 1
	case CARD_COLOR_Mei:
		index = INDEX_DRAGONPLUM - 1
	case CARD_COLOR_Hong:
		index = INDEX_DRAGONRED - 1
	case CARD_COLOR_Hei:
		index = INDEX_DRAGONBLOCK - 1
	default:
		index = INDEX_ERROR - 1
	}
	coins = float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	loseCoins += int64(coins)
	double[index] = gameConfig.Double[index]
	winArea[index] = true

	tColor := GetCardColor(this.TigerCard)
	switch tColor {
	case CARD_COLOR_Fang:
		index = INDEX_TIGERSPADE - 1
	case CARD_COLOR_Mei:
		index = INDEX_TIGERPLUM - 1
	case CARD_COLOR_Hong:
		index = INDEX_TIGERRED - 1
	case CARD_COLOR_Hei:
		index = INDEX_TIGERBLOCK - 1
	default:
		index = INDEX_ERROR - 1
	}
	coins = float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	loseCoins += int64(coins)
	double[index] = gameConfig.Double[index]
	winArea[index] = true

	dMidT := this.GetAreaCoin(INDEX_DRAGON-1) - this.GetAreaCoin(INDEX_TIGER-1)
	if dMidT == 0 {
		// 龙虎押注一样，上庄区域退还
		index = INDEX_BANKERWIN - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true

		index = INDEX_BANKERLOSE - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true
	}

	switch chart {
	case DRAGON:
		if dMidT < 0 {
			index = INDEX_BANKERWIN - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else if dMidT > 0 {
			index = INDEX_BANKERLOSE - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		}
	case TIGER:
		if dMidT < 0 {
			index = INDEX_BANKERLOSE - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else if dMidT > 0 {
			index = INDEX_BANKERWIN - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		}
	case DRAW:
		// 返还上庄
		index = INDEX_BANKERWIN - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true

		index = INDEX_BANKERLOSE - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true
	}

	// 计算玩家输赢
	for i, v := range winArea {
		if !v {
			continue
		}
		coins := int64(float64(this.GetUserAreaCoin(i)) * double[i])
		tWinArea[i] = coins
	}
	totCoins := this.GetUserAreaCoins()
	return totCoins - loseCoins, winArea[:], tWinArea[:]
}

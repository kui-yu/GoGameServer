package main

import (
	"fmt"
	"logs"
	"time"

	"bl.com/util"
)

//洗牌
func (this *ExtDesk) TimerShuffle(d interface{}) {
	logs.Debug("**********洗牌中**********")
	this.Lock()
	defer this.Unlock()

	this.ExtDeskInit()
	sd := GGameShuffleNotify{
		Id:        MSG_GAME_INFO_SHUFFLE_NOTIFY,
		Timer:     int32(gameConfig.Timer.ShuffleNum) * 1000,
		GameCount: this.Count,
	}
	this.BroadcastAll(MSG_GAME_INFO_SHUFFLE_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_SHUFFLE_NOTIFY
	this.AddTimer(gameConfig.Timer.Shuffle, gameConfig.Timer.ShuffleNum, this.TimerReady, nil)
}

//准备
func (this *ExtDesk) TimerReady(d interface{}) {
	logs.Debug("**********准备中**********")
	//	游戏局数加一
	this.Count++
	//处理玩家退出
	this.HandleExit()
	//处理某些玩家多局未下注
	this.HandleUndo()
	//	刷新座位
	this.UpdatePlayer()
	//	获取游戏ID
	this.GameId = GetJuHao() //----->util.BuildGameId(GCONFIG.GameType)
	Players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range Players {
		v.(*ExtPlayer).ResetAreaList() //重置用户下注列表
		if v.(*ExtPlayer).GetMsgId() == 0 {
			v.(*ExtPlayer).AddUndoTimes()
		}
	}
	sd := GGameReadyNotify{
		Id:          MSG_GAME_INFO_READY_NOTIFY,
		Timer:       int32(gameConfig.Timer.ReadyNum) * 1000,
		GameCount:   this.Count,
		GameId:      this.GameId, // 局号
		Limitcoinid: this.Limitcoinid,
	}
	allUser := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range allUser {
		//向全部玩家发送游戏准备的阶段号
		sd.SeatList = this.GetSeatInfo(v.(*ExtPlayer))
		v.(*ExtPlayer).SendNativeMsg(MSG_GAME_INFO_READY_NOTIFY, sd)
	}
	this.GameState = MSG_GAME_INFO_READY_NOTIFY
	this.AddTimer(gameConfig.Timer.Ready, gameConfig.Timer.ReadyNum, this.TimerSendCard, nil)
}

//发牌
func (this *ExtDesk) TimerSendCard(d interface{}) {
	logs.Debug("**********发牌中**********")
	this.Lock()
	defer this.Unlock()
	this.CardList = []int32{}
	showCard := this.CardMgr.HandCardInfo(6) //获取已经发好的牌
	for _, v := range showCard {
		this.CardList = append(this.CardList, int32(v))
	}
	this.RedCard = this.CardList[:3]
	this.BlackCard = this.CardList[3:]
	Rcard, Rcolor := SortHandCard(this.RedCard)   //将红方的牌进行排序并分出花色，数值
	Bcard, Bcolor := SortHandCard(this.BlackCard) //将黑方的牌进行排序并分出花色，数值
	//	判断牌的类型以及等级
	RGrade, _ := GetCardType(Rcard, Rcolor)
	BGrade, _ := GetCardType(Bcard, Bcolor)

	sd := GGameSendCardNotify{
		Id:             MSG_GAME_INFO_SEND_NOTIFY,
		Timer:          int32(gameConfig.Timer.SendCardNum) * 1000,
		RedCard:        Rcard,  //红方牌
		RedCardColor:   Rcolor, //红方牌花色
		RGrade:         RGrade, //红方牌等级
		BlackCard:      Bcard,  //黑方牌
		BlackCardColor: Bcolor, //黑方牌花色
		BGrade:         BGrade, //黑方牌等级
	}
	fmt.Println("发牌结构体:", sd)
	this.BroadcastAll(MSG_GAME_INFO_SEND_NOTIFY, sd)
	this.GameState = MSG_GAME_INFO_SEND_NOTIFY
	this.AddTimer(gameConfig.Timer.SendCard, gameConfig.Timer.SendCardNum, this.TimerBet, nil)
}

//下注
func (this *ExtDesk) TimerBet(d interface{}) {
	logs.Debug("**********下注中**********")
	this.Lock()
	defer this.Unlock()
	// logs.Debug("可下注区域：", this.BetArea)
	sd := GGameBetNotify{
		Id:    MSG_GAME_INFO_BET_NOTIFY,
		Timer: int32(gameConfig.Timer.BetNum) * 1000,
	}
	this.BroadcastAllSpec(MSG_GAME_INFO_BET_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_BET_NOTIFY
	this.AddTimer(gameConfig.Timer.Bet, gameConfig.Timer.BetNum, this.TimerStopBet, nil)
	this.AddTimer(gameConfig.Timer.NewBet, gameConfig.Timer.NewBetNum, this.HandleTimeOutBet, nil)

}

//停止下注
func (this *ExtDesk) TimerStopBet(d interface{}) {
	logs.Debug("**********停止下注**********")
	this.Lock()
	defer this.Unlock()

	tAreaCoins := this.GetAreaCoinsList()

	//消息公共部分
	sd := GGameStopBetNotify{
		Id:         MSG_GAME_INFO_STOP_BET_NOTIFY,
		Timer:      int32(gameConfig.Timer.StopBetNum) * 1000,
		TAreaCoins: tAreaCoins,
	}
	//座位玩家下注信息
	sd.SeatBetList = this.SeatMgr.GetSeatNewBetList()

	//新下注，排除座位玩家
	OtherNewBetList := this.SeatMgr.GetOtherNewBetList()

	//每个人的下注都不一样 需要简单处理
	playersList := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range playersList {
		p := v.(*ExtPlayer)
		newBetList := p.GetNewBetList()
		p.ColAreaCoins()
		TpAreaCoins := p.GetTotBetList()
		sd.PAreaCoins = TpAreaCoins

		if this.SeatMgr.IsOnSeat(p) {
			sd.OtherBetList = OtherNewBetList
		} else {
			sd.OtherBetList = util.LessInt64List(OtherNewBetList, newBetList)
		}
		p.SendNativeMsg(MSG_GAME_INFO_STOP_BET_NOTIFY, sd)
	}
	this.GameState = MSG_GAME_INFO_STOP_BET_NOTIFY
	this.AddTimer(gameConfig.Timer.StopBet, gameConfig.Timer.StopBetNum, this.TimerOpen, nil)
	/*-----------------------风控--------------------------*/

}

//开牌
func (this *ExtDesk) TimerOpen(d interface{}) {
	logs.Debug("**********开牌中**********")
	this.Lock()
	defer this.Unlock()
	Rcard, Rcolor := SortHandCard(this.RedCard)   //将红方的牌进行排序并分出花色，数值
	Bcard, Bcolor := SortHandCard(this.BlackCard) //将黑方的牌进行排序并分出花色，数值
	//	判断牌的类型以及等级
	RGrade, _ := GetCardType(Rcard, Rcolor)
	BGrade, _ := GetCardType(Bcard, Bcolor)
	sd := GGameOpenNotify{
		Id:        MSG_GAME_INFO_OPEN_NOTIFY,
		Timer:     int32(gameConfig.Timer.OpenNum) * 1000,
		RedCard:   this.RedCard,
		RedType:   RGrade,
		BlackCard: this.BlackCard,
		BlackType: BGrade,
	}
	logs.Debug("开牌阶段消息:%v", sd)
	this.BroadcastAll(MSG_GAME_INFO_OPEN_NOTIFY, sd)
	this.GameState = MSG_GAME_INFO_OPEN_NOTIFY
	this.AddTimer(gameConfig.Timer.Open, gameConfig.Timer.OpenNum, this.TimerAward, nil)

}

//创建有效投注
func (this *ExtDesk) BuildVaild(winArea []bool, areaList []int64) {
	logs.Debug("建立有效投注")
}

//派奖
func (this *ExtDesk) TimerAward(d interface{}) {
	logs.Debug("**********派奖中**********")
	this.Lock()
	defer this.Unlock()

	sd := GGameAwardNotify{
		Id:    MSG_GAME_INFO_AWARD_NOTIFY,
		Timer: int32(gameConfig.Timer.AwardNum) * 1000,
	}
	endTime := time.Now().Format("2006-01-02 15:04:05")
	var roomName string
	if GCONFIG.GradeType == 1 {
		roomName = "荣耀厅"
	} else if GCONFIG.GradeType == 2 {
		roomName = "王牌厅"
	} else if GCONFIG.GradeType == 3 {
		roomName = "战神厅"
	} else {
		roomName = "体验厅"
	}
	RD := RecordData{
		RoomName: roomName, //房间名称
		EndTime:  endTime,  //结束时间
		Date:     time.Now().Unix(),
	}
	Rcard, Rcolor := SortHandCard(this.RedCard)   //将红方的牌进行排序并分出花色，数值
	Bcard, Bcolor := SortHandCard(this.BlackCard) //将黑方的牌进行排序并分出花色，数值
	Rtype, _ := GetCardType(Rcard, Rcolor)
	Btype, _ := GetCardType(Bcard, Bcolor)
	sd.Rtype = Rtype
	sd.Btype = Btype
	chart, typeChart := this.CardMgr.CompareCard(Rcard, Rcolor, Bcard, Bcolor)

	sd.PairValue = this.CardMgr.CardValue
	// logs.Debug("*************赢方为对子的值:%v", sd.PairValue)
	this.RunChart = append([]int32{chart}, this.RunChart...)
	this.CardTypeChart = append([]int{typeChart}, this.CardTypeChart...)
	Rlength := len(this.RunChart)
	Clength := len(this.CardTypeChart)
	//判断输赢走势长度是否大于规定次数
	if Rlength > gameConfig.DeskInfo.RunChartCount {
		sd.RunChart = this.RunChart[:gameConfig.DeskInfo.RunChartCount]
	} else {
		sd.RunChart = this.RunChart
	}
	//判断牌型走势长度是否大于规定次数
	if Clength > gameConfig.DeskInfo.CardTypeChartCount {
		sd.CardTypeChart = this.CardTypeChart[:gameConfig.DeskInfo.CardTypeChartCount]
	} else {
		sd.CardTypeChart = this.CardTypeChart
	}
	//游戏结算，返回区域输赢情况
	WinArea, TWinArea := this.GameEnd(chart, typeChart)

	sd.WinArea = WinArea[:]
	sd.TWinArea = TWinArea[:]
	this.WinArea = WinArea[:]
	seatList := this.SeatMgr.GetSeatList()

	for _, seat := range seatList {
		win := seat.(*ExtPlayer).GetWinCoins()
		sd.SeatWinCoins = append(sd.SeatWinCoins, win)
		winArea := seat.(*ExtPlayer).GetWinList()
		TWinArea = util.LessInt64List(TWinArea, winArea)
	}
	betlist := G_DbGetGameServerData.GameConfig.TenChips
	if GetCostType() == 1 { //不是体验场
		//发送结算消息给数据库，简单记录
		GGE := GGameEnd{
			Id:          MSG_GAME_END_NOTIFY,
			GameId:      GCONFIG.GameType,
			GradeId:     GCONFIG.GradeType,
			RoomId:      GCONFIG.RoomType,
			GameRoundNo: this.GameId,
			Mini:        false,
			SetLeave:    1, //是否设置离开，0离开，1不离开
		}

		//发送消息给大厅去记录游戏记录
		GGR := GGameRecord{
			Id:          MSG_GAME_END_RECORD,
			GameId:      GCONFIG.GameType,
			GradeId:     GCONFIG.GradeType,
			RoomId:      GCONFIG.RoomType,
			GradeNumber: 1,
			GameRoundNo: this.GameId,
			RedCard:     this.RedCard,
			BlackCard:   this.BlackCard,
		}
		//游戏记录牌型
		Wintype1 := Wintype{} //游戏记录
		switch this.WinId {   //红黑
		case 0:
			Wintype1.WinID = "红"
		case 1:
			Wintype1.WinID = "黑"
		}
		switch typeChart { //牌型
		case 1:
			Wintype1.Type = "单张"
		case 2:
			Wintype1.Type = "对子"
		case 3:
			Wintype1.Type = "顺子"
		case 4:
			Wintype1.Type = "同花"
		case 5:
			Wintype1.Type = "同花顺"
		case 6:
			Wintype1.Type = "豹子"
		}
		RD.TypeList = append(RD.TypeList, Wintype1)

		players := this.SeatMgr.GetUserList(len(this.Players))
		for _, v := range players {
			var allbet1 int64
			var vaild int64
			var area = []string{"红", "黑", "幸运一击"}
			var Barea []string
			p := v.(*ExtPlayer)
			sd.PWin = p.GetWinCoins()
			sd.PWinArea = p.GetWinList()
			areaList := p.GetTotBetList()

			//判断这局用户是否下注了
			if len(p.betArealist) == 0 {
				//如果没有下注，则返回此局赢方ID
				sd.WinAreaId = this.WinId
			} else {
				//如果下注，则返回这局下注在红方或黑方的ID
				sd.WinAreaId = this.betId
			}
			for _, v := range areaList {
				vaild += v
			}
			betCoins := p.GetTotAreaCoins()
			// 结算
			sd.PrizeCoins = sd.PWin - betCoins
			p.Award()
			sd.PCoins = p.GetCoins()
			//判定下注金币列表的限制
			var as int32 = -1
			for i := len(betlist) - 1; i >= 0; i-- {
				if sd.PCoins >= betlist[i] {
					as = int32(i)
					break
				}
			}
			this.Limitcoinid = as
			sd.Limitcoinid = as

			if betCoins > 0 {
				GGE.UserCoin = append(GGE.UserCoin, GGameEndInfo{
					UserId:      p.GetUid(),
					UserAccount: p.Account,
					BetCoins:    betCoins,
					ValidBet:    vaild,
					PrizeCoins:  sd.PWin - betCoins,
					Robot:       p.Robot,
				})
				fmt.Println("发送数据库数据:%v", GGE)
				p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &GGE)
				GGE.UserCoin = []GGameEndInfo{}

				if p.Robot {
					continue
				}
				GGRInfo := GGameRecordInfo{
					UserId:      p.GetUid(),
					UserAccount: p.Account,
					BetCoins:    betCoins,
					BetArea:     p.GetTotBetList(),
					PrizeCoins:  sd.PWin - betCoins,
					CoinsAfter:  sd.PCoins,
					Robot:       p.Robot,
				}
				GGRInfo.CoinsBefore = GGRInfo.CoinsAfter - GGRInfo.PrizeCoins
				GGR.UserRecord = append(GGR.UserRecord, GGRInfo)
				fmt.Println("发送消息给大厅去记录游戏记录:%v", GGR)
				p.SendNativeMsgForce(MSG_GAME_END_RECORD, &GGR)
				GGR.UserRecord = []GGameRecordInfo{}
			}

			//游戏记录
			for _, bet := range p.GetTotBetList() {
				allbet1 += bet
			}
			RD.AllBet = allbet1
			// logs.Debug("总下注：", p.Nick, RD.AllBet)
			RD.WinOrLost = sd.PWin - betCoins
			for i, _ := range p.GetTotBetList() {
				if p.GetTotBetList()[i] > 0 {
					Barea = append(Barea, area[i])
				}
			}
			RD.BetArea = Barea
			if allbet1 > 0 {
				fmt.Println("该玩家有投注，所以存储游戏记录")
				G_AllRecord.AddRecord(p.Uid, &RD)
			} else {
				fmt.Println("该玩家没有投注，不存游戏记录")
			}
			//发送派奖记录
			logs.Debug("派奖记录：", sd)
			p.SendNativeMsg(MSG_GAME_INFO_AWARD_NOTIFY, sd)
			p.betArealist = []int{}
		}
	} else { //如果是体验场，则发送派奖信息
		logs.Debug("这是体验场")
		Wintype1 := Wintype{} //游戏记录
		switch this.WinId {   //红黑
		case 0:
			Wintype1.WinID = "红"
		case 1:
			Wintype1.WinID = "黑"
		}
		switch typeChart { //牌型
		case 1:
			Wintype1.Type = "单张"
		case 2:
			Wintype1.Type = "对子"
		case 3:
			Wintype1.Type = "顺子"
		case 4:
			Wintype1.Type = "同花"
		case 5:
			Wintype1.Type = "同花顺"
		case 6:
			Wintype1.Type = "豹子"
		}
		RD.TypeList = append(RD.TypeList, Wintype1)

		player := this.SeatMgr.GetUserList(len(this.Players))
		for _, v := range player {
			var allbet2 int64
			var area = []string{"红", "黑", "幸运一击"}
			var Barea []string
			p := v.(*ExtPlayer)
			sd.PWin = p.GetWinCoins()
			sd.PWinArea = p.GetWinList()

			//判断这局用户是否下注了
			if len(p.betArealist) == 0 {
				//如果没有下注，则返回此局赢方ID
				sd.WinAreaId = this.WinId
			} else {
				//如果下注，则返回这局下注在红方或黑方的ID
				sd.WinAreaId = this.betId
			}

			betCoins := p.GetTotAreaCoins()
			// 结算
			sd.PrizeCoins = sd.PWin - betCoins
			p.Award()
			sd.PCoins = p.GetCoins()
			//判定玩家下注金币列表的限制
			var as int32 = -1
			for i := len(betlist) - 1; i >= 0; i-- {
				if sd.PCoins >= betlist[i] {
					as = int32(i)
					break
				}
			}
			this.Limitcoinid = as
			sd.Limitcoinid = as
			for _, bet := range p.GetTotBetList() {
				allbet2 += bet
			}
			RD.AllBet = allbet2
			RD.WinOrLost = sd.PWin - betCoins
			for i, _ := range p.GetTotBetList() {
				if p.GetTotBetList()[i] > 0 {
					Barea = append(Barea, area[i])
				}
			}
			RD.BetArea = Barea
			logs.Debug("游戏记录：%v", RD)
			if allbet2 > 0 {
				fmt.Println("该玩家有投注，所以存储游戏记录")
				G_AllRecord.AddRecord(p.Uid, &RD)
			} else {
				fmt.Println("该玩家没有投注，不存游戏记录")
			}
			p.SendNativeMsg(MSG_GAME_INFO_AWARD_NOTIFY, sd)
			p.betArealist = []int{}
			logs.Debug("玩家昵称：%v  ;游戏派奖记录：%v", p.Nick, sd)
		}
	}
	this.GameState = MSG_GAME_INFO_AWARD_NOTIFY
	this.ResetAreaCoins()
	this.AddTimer(gameConfig.Timer.Award, gameConfig.Timer.AwardNum, this.TimerOver, nil)
	for _, v := range this.Players {
		v.Rest()
	}
}

//此局游戏结束，将跳转洗牌
func (this *ExtDesk) TimerOver(d interface{}) {
	logs.Debug("游戏结束")
	this.AddTimer(gameConfig.Timer.Over, gameConfig.Timer.OverNum, this.TimerShuffle, nil)
}

//结算
func (this *ExtDesk) GameEnd(chart int32, typeChart int) ([]bool, []int64) {
	logs.Debug("********游戏结算中********")
	var winArea [3]bool
	var tWinArea [3]int64
	var double [3]float64
	for i := range winArea {
		winArea[i] = false
		double[i] = 0
		tWinArea[i] = 0
	}
	//计算输赢
	var loseCoins int64 = 0
	var index int
	var coins float64
	//红黑区域中奖
	switch chart {
	case RED:
		index = INDEX_RED - 1
		winArea[0] = true
	case BLACK:
		index = INDEX_BLACK - 1
		winArea[1] = true
	default:
		index = INDEX_ERROR - 1
	}
	this.WinId = int32(index)
	coins = float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index] * 2
	loseCoins += int64(coins)
	double[index] = gameConfig.Double[index]

	//根据牌类型计算输赢金额
	switch typeChart {
	case CARD_BOOM:
		double[2] = gameConfig.Double[2]
		winArea[2] = true
	case CARD_FLUSH:
		double[2] = gameConfig.Double[3]
		winArea[2] = true
	case CARD_TONGHUA:
		double[2] = gameConfig.Double[4]
		winArea[2] = true
	case CARD_SHUNZI:
		double[2] = gameConfig.Double[5]
		winArea[2] = true
	case CARD_PAIR:
		if this.CardMgr.CardValue >= 9 {
			double[2] = gameConfig.Double[6]
			winArea[2] = true
		} else {
			double[2] = 0
			winArea[2] = false
		}
	default:
		double[2] = 0
		winArea[2] = false
	}
	if this.betId == INDEX_LUCKYBLOW-1 {
		switch typeChart {
		case CARD_BOOM:
			double[2] = gameConfig.Double[2]
			winArea[2] = true
		case CARD_FLUSH:
			double[2] = gameConfig.Double[3]
			winArea[2] = true
		case CARD_TONGHUA:
			double[2] = gameConfig.Double[4]
			winArea[2] = true
		case CARD_SHUNZI:
			double[2] = gameConfig.Double[5]
			winArea[2] = true
		case CARD_PAIR:
			// logs.Debug("对子的值！：", this.CardMgr.CardValue)
			if this.CardMgr.CardValue >= 9 {
				double[2] = gameConfig.Double[6]
				winArea[2] = true
			} else {
				double[2] = 0
				winArea[2] = false
			}
		default:
			double[2] = 0
			winArea[2] = false
		}
	}
	coins = float64(this.GetUserAreaCoin(2)) * gameConfig.Double[2]
	loseCoins += int64(coins)
	//计算玩家输赢fallthrough
	for i, v := range winArea {
		if !v {
			continue
		}
		coins := int64(float64(this.GetUserAreaCoin(i)) * double[i])

		tWinArea[i] = coins
	}
	players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range players {
		p := v.(*ExtPlayer)
		p.BuildWinList(double[:])
		p.AddWinList()
		p.AddBetList()
	}
	torCoins := this.GetAreaCoins()
	AddLocalStock(torCoins - loseCoins)
	return winArea[:], tWinArea[:]
}

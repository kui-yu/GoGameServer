// 定时器回调函数，控制百人场状态
package main

import (
	"logs"
	"time"

	"bl.com/util"
)

// 洗牌
func (this *ExtDesk) TimerShuffle(d interface{}) {
	this.Lock()
	defer this.Unlock()

	this.ExtDeskInit()

	sd := GGameShuffleNotify{
		Id:         MSG_GAME_INFO_SHUFFLE_NOTIFY,
		Timer:      int32(gameConfig.Timer.ShuffleNum) * 1000,
		LeftCount:  int32(len(this.CardMgr.MVSourceCard)),
		RightCount: int32(len(this.CardMgr.OutCards)),
		GameCount:  this.Count,
	}

	this.BroadcastAll(MSG_GAME_INFO_SHUFFLE_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_SHUFFLE_NOTIFY
	this.AddTimer(gameConfig.Timer.Shuffle, gameConfig.Timer.ShuffleNum, this.TimerReady, nil)
}

// 准备
func (this *ExtDesk) TimerReady(d interface{}) {
	this.Lock()
	defer this.Unlock()

	// 局数+1
	this.Count++

	this.HandleExit()
	this.HandleUndo()

	this.IdleCard = []int32{}
	this.BankerCard = []int32{}
	this.IdleDians = []int32{}
	this.BankerDians = []int32{}

	// 刷新座位
	this.UpdatePlayer()

	this.GameId = GetJuHao() // util.BuildGameId(GCONFIG.GameType)
	players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range players {
		v.(*ExtPlayer).ResetAreaList()
		// 用户没有下注，增加未下注次数
		if v.(*ExtPlayer).GetMsgId() == 0 {
			v.(*ExtPlayer).AddUndoTimes()
		}
		if v.(*ExtPlayer).LiXian { //判断玩家是否离线
			// this.SeatMgr.DelPlayer(v)
			// this.LeaveByForce(v)
			logs.Debug("---------准备阶段查看离线玩家昵称：%v", v.(*ExtPlayer).Nick)
		}
	}
	sd := GGameReadyNotify{
		Id:         MSG_GAME_INFO_READY_NOTIFY,
		Timer:      int32(gameConfig.Timer.ReadyNum) * 1000,
		LeftCount:  int32(len(this.CardMgr.MVSourceCard)),
		RightCount: int32(len(this.CardMgr.OutCards)),
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
	this.Lock()
	defer this.Unlock()

	sd := GGameSendCardNotify{
		Id:    MSG_GAME_INFO_SEND_NOTIFY,
		Timer: int32(gameConfig.Timer.SendCardNum) * 1000,
	}

	this.BroadcastAll(MSG_GAME_INFO_SEND_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_SEND_NOTIFY
	this.AddTimer(gameConfig.Timer.SendCard, gameConfig.Timer.SendCardNum, this.TimerBet, nil)
}

// 下注
func (this *ExtDesk) TimerBet(d interface{}) {
	this.Lock()
	defer this.Unlock()

	if this.Count == int32(gameConfig.DeskInfo.BetLimit)+1 {
		this.BetArea[INDEX_BIG-1] = false
		this.BetArea[INDEX_SMALL-1] = false
	}

	sd := GGameBetNotify{
		Id:      MSG_GAME_INFO_BET_NOTIFY,
		Timer:   int32(gameConfig.Timer.BetNum) * 1000,
		BetArea: this.BetArea,
	}
	for _, v := range this.Players {
		if v.LiXian { //判断玩家是否离线
			logs.Debug("---------下注阶段查看离线玩家昵称：%v", v.Nick)
		}
	}
	this.BroadcastAll(MSG_GAME_INFO_BET_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_BET_NOTIFY
	this.AddTimer(gameConfig.Timer.Bet, gameConfig.Timer.BetNum, this.TimerStopBet, nil)

	// 定时广播
	this.AddTimer(gameConfig.Timer.NewBet, gameConfig.Timer.NewBetNum, this.HandleTimeOutBet, nil)
}

// 停止下注
func (this *ExtDesk) TimerStopBet(d interface{}) {
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

		// sd.PAreaCoins = TPAreaCoins
		sd.PAreaCoins = p.GetNTAreaCoinsList()
		sd.OtherBetList = this.SeatMgr.GetOtherNewBetList2(p.Uid)

		// if this.SeatMgr.IsOnSeat(p) {
		// 	sd.OtherBetList = otherNewBetList
		// } else {
		// 	sd.OtherBetList = util.LessInt64List(otherNewBetList, newBetList)
		// }
		if p.LiXian { //判断玩家是否离线
			logs.Debug("**********停止下注阶段查看离线玩家昵称：%v", p.Nick)
			// this.SeatMgr.DelPlayer(v)
			// this.LeaveByForce(v)
		}
		p.SendNativeMsg(MSG_GAME_INFO_STOP_BET_NOTIFY, sd)
	}

	for _, user := range this.Players {
		user.ColAreaCoins()
	}

	this.GameState = MSG_GAME_INFO_STOP_BET_NOTIFY
	this.AddTimer(gameConfig.Timer.StopBet, gameConfig.Timer.StopBetNum, this.TimerOpen, nil)

	// 制作发牌
	this.BuildCards()
}

// 开牌
func (this *ExtDesk) TimerOpen(d interface{}) {
	this.Lock()
	defer this.Unlock()

	timer := gameConfig.Timer.OpenNum
	timer += (len(this.BankerCard) + len(this.IdleCard) - 4) * gameConfig.Timer.AddNum
	AllStageTime[5] = timer //赋值阶段时间
	sd := GGameOpenNotify{
		Id:          MSG_GAME_INFO_OPEN_NOTIFY,
		Timer:       int32(timer) * 1000,
		IdleCard:    this.IdleCard,
		BankerCard:  this.BankerCard,
		IdleDians:   this.IdleDians,
		BankerDians: this.BankerDians,
	}

	this.BroadcastAll(MSG_GAME_INFO_OPEN_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_OPEN_NOTIFY
	this.OpenCardTime = timer //开牌阶段时间记录
	this.AddTimer(gameConfig.Timer.Open, timer, this.TimerAward, nil)
}

func (this *ExtDesk) BuildValid(winArea []bool, areaList []int64) int64 {
	// 同时押庄、闲，区域不计算有效打码
	if areaList[INDEX_BANKER-1] > 0 && areaList[INDEX_IDLE-1] > 0 {
		areaList[INDEX_IDLE-1] = -1
		areaList[INDEX_BANKER-1] = -1
	}

	// 同时押输、赢，区域不计算有效打码
	if areaList[INDEX_BANKERLOSE-1] > 0 && areaList[INDEX_BANKERWIN-1] > 0 {
		areaList[INDEX_BANKERWIN-1] = -1
		areaList[INDEX_BANKERLOSE-1] = -1
	}

	// 同时押大、小，区域不计算有效打码
	if areaList[INDEX_SMALL-1] > 0 && areaList[INDEX_BIG-1] > 0 {
		areaList[INDEX_BIG-1] = -1
		areaList[INDEX_SMALL-1] = -1
	}

	// 退还不计算打码
	// if winArea[INDEX_BANKER-1] && winArea[INDEX_IDLE-1] {
	// 	areaList[INDEX_IDLE-1] = -1
	// 	areaList[INDEX_BANKER-1] = -1
	// }
	// if winArea[INDEX_BANKERLOSE-1] && winArea[INDEX_BANKERWIN-1] {
	// 	areaList[INDEX_BANKERWIN-1] = -1
	// 	areaList[INDEX_BANKERLOSE-1] = -1
	// }

	var ret int64
	for _, v := range areaList {
		ret += v
	}

	return ret
}

// 派奖
func (this *ExtDesk) TimerAward(d interface{}) {
	this.OpenCardTime = 0 //开牌阶段时间记录清空
	this.Lock()
	defer this.Unlock()
	sd := GGameAwardNotify{
		Id:    MSG_GAME_INFO_AWARD_NOTIFY,
		Timer: int32(gameConfig.Timer.AwardNum) * 1000,
	}

	// 添加走势
	var double byte
	if GetCardValue(this.IdleCard[0]) == GetCardValue(this.IdleCard[1]) {
		double += IDLE << 4
	}

	if GetCardValue(this.BankerCard[0]) == GetCardValue(this.BankerCard[1]) {
		double += BANKER << 4
	}

	chart := this.CardMgr.CompareCard(this.IdleCard, this.BankerCard)
	this.RunChart = append(this.RunChart, int32(double+chart))

	//更新走势
	sd.RunChart = this.RunChart

	// 游戏结算，返回区域输赢情况
	_, WinArea, TWinArea := this.GameEnd(chart, double)

	for i, v := range this.BetArea {
		WinArea[i] = WinArea[i] && v
	}

	sd.WinArea = WinArea[:]
	sd.TWinArea = TWinArea[:]
	this.WinArea = WinArea[:]

	// 添加走势类型次数
	sd.TypeTimes = this.TypeTimes

	seatList := this.SeatMgr.GetSeatList()
	for _, seat := range seatList {
		win := seat.(*ExtPlayer).GetWinCoins()
		sd.SeatWinCoins = append(sd.SeatWinCoins, win)
		winArea := seat.(*ExtPlayer).GetWinList()
		TWinArea = util.LessInt64List(TWinArea, winArea)
	}

	sd.OtherWinArea = TWinArea
	if GetCostType() == 1 { //如果不是体验场
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
			BankerCard:  this.BankerCard,
			IdleCard:    this.IdleCard,
		}
		players := this.SeatMgr.GetUserList(len(this.Players))
		for _, v := range players {
			p := v.(*ExtPlayer)
			sd.PWin = p.GetWinCoins()
			sd.PWinArea = p.GetWinList()

			areaList := p.GetTotBetList()
			valid := this.BuildValid(WinArea, areaList)
			betCoins := p.GetTotAreaCoins()
			// 结算
			sd.PrizeCoins = sd.PWin - betCoins
			p.Award()
			sd.PCoins = p.GetCoins()
			//
			if betCoins > 0 {
				//游戏记录
				rd := RecordData{}
				rd.Grade = GCONFIG.GradeType
				rd.Date = time.Now().Unix()
				rd.YaZhu = int(betCoins)
				rd.YingKui = int(sd.PrizeCoins)
				for _, ae := range areaList {
					rd.YzQuYu = append(rd.YzQuYu, int(ae))
				}
				rd.XianPai = append([]int32{}, this.IdleCard...)
				rd.ZhuangPai = append([]int32{}, this.BankerCard...)
				if sd.WinArea[INDEX_BIG-1] {
					rd.DaXiao = 1
				} else if sd.WinArea[INDEX_SMALL-1] {
					rd.DaXiao = 2
				}
				G_AllRecord.AddRecord(p.Uid, &rd)
				//
				dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
					UserId:      p.GetUid(),
					UserAccount: p.Account,
					BetCoins:    betCoins,
					ValidBet:    valid,
					PrizeCoins:  sd.PWin - betCoins,
					Robot:       p.Robot,
				})
				if !p.Robot {
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
				p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
				rdreq.UserRecord = []GGameRecordInfo{}
			}
			if GetCostType() == 1 && !p.Robot {
				AddCD(-sd.PrizeCoins)
				AddLocalStock(-sd.PrizeCoins)
				//fmt.Println("添加库存:", -sd.PrizeCoins)
			}
			p.SendNativeMsg(MSG_GAME_INFO_AWARD_NOTIFY, sd)
		}
	} else {
		//如果是体验场，只发送派奖信息
		player := this.SeatMgr.GetUserList(len(this.Players))
		for _, v := range player {
			p := v.(*ExtPlayer)
			sd.PWin = p.GetWinCoins()
			sd.PWinArea = p.GetWinList()
			betCoins := p.GetTotAreaCoins()
			// 结算
			sd.PrizeCoins = sd.PWin - betCoins
			p.Award()
			sd.PCoins = p.GetCoins()
			p.SendNativeMsg(MSG_GAME_INFO_AWARD_NOTIFY, sd)
		}
	}
	this.GameState = MSG_GAME_INFO_AWARD_NOTIFY

	this.ResetAreaCoins()
	//重置玩家下注
	for _, v := range this.Players {
		v.ResetAreaList()
		if v.LiXian { //判断玩家是否离线
			logs.Debug("++++++++派奖阶段查看离线玩家昵称：%v", v.Nick)
			this.SeatMgr.DelPlayer(v)
			this.LeaveByForce(v)
		}
	}
	// 剩余张数 < 6 结束，否则继续准备
	if len(this.CardMgr.MVSourceCard) < 6 {
		this.TimerOver(nil)
		return
	}
	this.AddTimer(gameConfig.Timer.Award, gameConfig.Timer.AwardNum, this.TimerReady, nil)
}

// 结束
func (this *ExtDesk) TimerOver(d interface{}) {
	this.AddTimer(gameConfig.Timer.Over, gameConfig.Timer.OverNum, this.TimerShuffle, nil)
}

// 结算
func (this *ExtDesk) GameEnd(chart, pair byte) (int64, []bool, []int64) {
	//
	var winArea [9]bool
	var tWinArea [9]int64
	var double [9]float64
	for i := range winArea {
		winArea[i] = false
		double[i] = 0
		tWinArea[i] = 0
	}

	// 计算闲输赢
	var loseCoins int64 = 0
	var index int
	bankerIndex := INDEX_BANKER - 1
	idleIndex := INDEX_IDLE - 1
	bankerCoins := this.GetAreaCoin(bankerIndex)
	idleCoins := this.GetAreaCoin(idleIndex)

	switch chart {
	case IDLE:
		coins := float64(idleCoins) * gameConfig.Double[idleIndex]
		loseCoins += int64(coins)
		double[idleIndex] = gameConfig.Double[idleIndex]
		winArea[idleIndex] = true

		if idleCoins < bankerCoins {
			index = INDEX_BANKERWIN - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else if idleCoins > bankerCoins {
			index = INDEX_BANKERLOSE - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else {
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
		this.TypeTimes[1]++
	case BANKER:
		coins := float64(bankerCoins) * gameConfig.Double[bankerIndex]
		loseCoins += int64(coins)
		double[bankerIndex] = gameConfig.Double[bankerIndex]
		winArea[bankerIndex] = true

		if bankerCoins < idleCoins {
			index = INDEX_BANKERWIN - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else if bankerCoins > idleCoins {
			index = INDEX_BANKERLOSE - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else {
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
		this.TypeTimes[0]++
	case DRAW:
		index = INDEX_DRAW - 1
		coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
		loseCoins += int64(coins)
		double[index] = gameConfig.Double[index]
		winArea[index] = true

		// 返还上庄
		index = INDEX_BANKERWIN - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true
		index = INDEX_BANKERLOSE - 1
		loseCoins += this.GetUserAreaCoin(index)
		double[index] = 1
		winArea[index] = true
		// 返回庄、闲
		loseCoins += idleCoins
		double[idleIndex] = 1
		winArea[idleIndex] = true
		loseCoins += bankerCoins
		double[bankerIndex] = 1
		winArea[bankerIndex] = true

		this.TypeTimes[2]++
	}

	if len(this.IdleCard)+len(this.BankerCard) == 4 {
		index = INDEX_SMALL - 1
		coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
		loseCoins += int64(coins)
		double[index] = gameConfig.Double[index]
		winArea[index] = true
	} else {
		index = INDEX_BIG - 1
		coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
		loseCoins += int64(coins)
		double[index] = gameConfig.Double[index]
		winArea[index] = true
	}

	if pair&(IDLE<<4) > 0 {
		index = INDEX_IDLEPAIR - 1
		coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
		loseCoins += int64(coins)
		double[index] = gameConfig.Double[index]
		winArea[index] = true
		this.TypeTimes[4]++
	}

	if pair&(BANKER<<4) > 0 {
		index = INDEX_BANKERPAIR - 1
		coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
		loseCoins += int64(coins)
		double[index] = gameConfig.Double[index]
		winArea[index] = true
		this.TypeTimes[3]++
	}

	this.TypeTimes[5]++
	// 计算玩家输赢
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

	totCoins := this.GetUserAreaCoins()
	// if GetCostType() == 1 {
	// 	AddLocalStock(totCoins - loseCoins)
	// }
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

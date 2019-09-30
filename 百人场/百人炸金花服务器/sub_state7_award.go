package main

import (
	"fmt"
	"logs"

	"time"

	"bl.com/util"
)

// 派奖
func (this *ExtDesk) TimerAward(d interface{}) {
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
	//牌型应该为东南西北，不是东西南北
	arr := append([]int32{}, this.TypeList...)
	arr[1], arr[2] = arr[2], arr[1]
	RD := RecordData{
		RoomName:    roomName,         //房间名称
		TBankerCard: this.TBankerType, //庄牌型
		TypeList:    arr,              //闲（东南西北）牌型
		EndTime:     endTime,          //结束时间
		Date:        time.Now().Unix(),
	}
	// 添加走势
	// 庄
	_, bType := this.CardMgr.GetMaxCards(this.BankerCard)
	this.BankRunChart = append(this.BankRunChart, bType)
	if len(this.BankRunChart) > gameConfig.DeskInfo.RunChartCount {
		this.BankRunChart = this.BankRunChart[1:]
	}

	// 闲
	for i := 0; i < 4; i++ {
		win := this.CardMgr.CompareCard(this.MIdleCard[i], this.MBankerCard)
		this.RunChart[i] = append(this.RunChart[i], win)
	}

	for i := 0; i < 4; i++ {
		length := len(this.RunChart[i])

		if length > gameConfig.DeskInfo.RunChartCount {
			this.RunChart[i] = this.RunChart[i][1:]
		}
	}
	sd.BankRunChart = this.BankRunChart

	// 游戏结算，返回区域输赢情况
	_, WinArea, TWinArea := this.GameEnd(this.TypeList)
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
			var allbet1 int64
			p := v.(*ExtPlayer)
			sd.PWin = p.GetWinCoins()
			sd.PWinArea = p.GetWinList()

			valid := p.GetTotAreaCoins()

			betCoins := p.GetTotAreaCoins()
			// 结算
			sd.PrizeCoins = sd.PWin - betCoins
			p.Award()
			sd.PCoins = p.GetCoins()

			if betCoins > 0 && !p.Robot {
				dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
					UserId:      p.GetUid(),
					UserAccount: p.Account,
					BetCoins:    betCoins,
					ValidBet:    valid,
					PrizeCoins:  sd.PWin - betCoins,
					Robot:       p.Robot,
				})
				p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
				dbreq.UserCoin = []GGameEndInfo{}

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

			for _, bet := range p.GetTotBetList() {
				allbet1 += bet
				RD.AllBet = allbet1
			}
			logs.Debug("总下注：", p.Nick, RD.AllBet)
			RD.WinOrLost = sd.PWin - betCoins
			RD.BetArea = p.GetTotBetList()
			logs.Debug("游戏记录：", RD)
			if allbet1 > 0 {
				logs.Debug("该玩家有投注，所以存储游戏记录")
				G_AllRecord.AddRecord(p.Uid, &RD)
			} else {
				logs.Debug("该玩家没有投注，不存游戏记录")
			}
			p.SendNativeMsg(MSG_GAME_INFO_AWARD_NOTIFY, sd)
		}
	} else { //如果是体验场，只发送派奖信息
		logs.Debug("体验场")
		player := this.SeatMgr.GetUserList(len(this.Players))
		for _, v := range player {
			var allbet2 int64
			p := v.(*ExtPlayer)
			sd.PWin = p.GetWinCoins()
			sd.PWinArea = p.GetWinList()
			betCoins := p.GetTotAreaCoins()
			// 结算
			sd.PrizeCoins = sd.PWin - betCoins
			p.Award()
			sd.PCoins = p.GetCoins()
			for _, bet := range p.GetTotBetList() {
				allbet2 += bet
			}
			logs.Debug("体验场总下注：", p.Nick, allbet2)
			RD.AllBet = allbet2
			RD.WinOrLost = sd.PWin - betCoins
			RD.BetArea = p.GetTotBetList()
			logs.Debug("游戏记录：", RD)
			if allbet2 > 0 {
				logs.Debug("该玩家有投注，所以存储游戏记录")
				G_AllRecord.AddRecord(p.Uid, &RD)
			} else {
				logs.Debug("该玩家没有投注，不存游戏记录")
			}
			sd.PCoins = p.GetCoins()
			p.SendNativeMsg(MSG_GAME_INFO_AWARD_NOTIFY, sd)
		}
	}

	this.GameState = MSG_GAME_INFO_AWARD_NOTIFY
	fmt.Printf("派奖阶段=>牌1:%v,牌2:%v\n", this.BankerCard, this.IdleCard)
	this.ResetAreaCoins()
	this.TimerOver(nil)
}

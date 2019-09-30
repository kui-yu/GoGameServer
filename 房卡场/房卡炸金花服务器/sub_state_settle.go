package main

import (
	"logs"
	"time"
)

func (this *ExtDesk) GameStateSettle() {
	logs.Debug("结算阶段")
	this.BroadStageTime(STAGE_SETTLE_TIME)
	this.runTimer(STAGE_SETTLE_TIME, this.HandleGameSettleCoin)
}

//阶段-结算
func (this *ExtDesk) HandleGameSettleCoin(d interface{}) {
	// logs.Debug("阶段-结算")
	//发送结算信息
	this.ClearTimer()

	var sum int64
	sum = 0
	info := GSSettlePlayInfo{
		Id: MSG_GAME_INFO_SETTLE,
	}
	if len(this.SettleContest) > 0 {
		//金币不足比牌-结算
		info.Count = len(this.SettleContest)
		// fmt.Println("结算 out of range", len(this.SettleContest))
		info.PContest = make([]Contest, 0, info.Count)
		for i := 0; i < info.Count; i++ {
			info.PContest = append(info.PContest, Contest{Person_1: this.SettleContest[i].Person_1, Person_2: this.SettleContest[i].Person_2, Winner: this.SettleContest[i].Winner})
		}
	}

	if this.Round >= GameRound && GameRound != -1 {
		// logs.Debug("强制比牌")
		var contestPlayer *ExtPlayer
		var allPlayer []*ExtPlayer
		for _, v := range this.Players {
			if v.CardType != 2 {
				allPlayer = append(allPlayer, v)
			}
		}
		contestPlayer = allPlayer[0]

		for i := 1; i < len(allPlayer); i++ { //
			result := GetResult(contestPlayer.HandCards, contestPlayer.HandColor, allPlayer[i].HandCards, allPlayer[i].HandColor)
			msg := Contest{Person_1: contestPlayer.ChairId, Person_2: allPlayer[i].ChairId}
			var b bool = true
			if contestPlayer.CardType == 2 || allPlayer[i].CardType == 2 {
				b = false
			}
			if result == 0 { //输家牌型变换
				allPlayer[i].CardType = 2
				msg.Winner = contestPlayer.ChairId
			} else {
				contestPlayer.CardType = 2
				msg.Winner = allPlayer[i].ChairId
				contestPlayer = allPlayer[i]
			}
			if b {
				info.PContest = append(info.PContest, msg)
			}
		}
	}

	info.Count = len(info.PContest)
	for _, v := range this.Players {
		if v.CardType != 2 {
			info.SCard = append(info.SCard, SettleCard{ChairId: v.ChairId, Identity: 0, HandCard: v.OldHandCard, Lv: v.CardLv})
		} else {
			if !v.IsGU {
				info.SCard = append(info.SCard, SettleCard{ChairId: v.ChairId, Identity: 1, HandCard: v.OldHandCard, Lv: v.CardLv})
			}
		}
	}

	for _, v := range this.Players { //
		obj := GSRecordInfo{} //记录战绩
		obj.WinDate = time.Now().Format("2006-01-02 15:04:05")
		if v.CardType != 2 {
			for i := 0; i < len(this.CoinList); i++ {
				sum += this.CoinList[i]
			}
			// allsum := sum
			for i := 0; i < len(v.PayCoin); i++ {
				sum -= v.PayCoin[i]
			}
			//v.RateCoins = float64(allsum*this.Bscore) * this.Rate
			//v.WinCoins = (sum * this.Bscore) - int64(v.RateCoins)
			v.WinCoins = (sum * this.Bscore)
			v.Coins += v.WinCoins
			info.CList = append(info.CList, CoinList{ChairId: v.ChairId, WinCoins: v.WinCoins, Coins: v.Coins})
			obj.WinCoins += v.WinCoins //每局战绩记录
			v.TotalCoins += v.WinCoins //记录总结算金币
			// fmt.Println("win:", sum, "cardLv:", v.CardLv, v.Coins)
			if v.Robot {
				AddLocalStock(v.WinCoins)
			}
		} else {
			var dsum int64
			dsum = 0
			for i := 0; i < len(v.PayCoin); i++ {
				dsum -= v.PayCoin[i]
			}
			v.WinCoins = dsum * this.Bscore
			v.Coins = v.Coins + dsum*this.Bscore
			info.CList = append(info.CList, CoinList{ChairId: v.ChairId, WinCoins: v.WinCoins, Coins: v.Coins})
			obj.WinCoins += v.WinCoins //每局战绩记录
			v.TotalCoins += v.WinCoins //记录总结算金币
			if v.Robot {
				AddLocalStock(v.WinCoins)
			}
			// fmt.Println("lost:", dsum, "cardLv:", v.CardLv, v.Coins)
		}
		v.RecordInfos = append(v.RecordInfos, obj) //添加战绩记录
		//当局结算发送游戏记录
		v.SendNativeMsg(MSG_GAME_INFO_RECORD_INFO_REPLY, &GSRecordInfos{
			Id:    MSG_GAME_INFO_RECORD_INFO_REPLY,
			Infos: v.RecordInfos,
		})
	}
	this.BroadcastAll(MSG_GAME_INFO_SETTLE, &info)
	// this.nextStage(STAGE_SETTLE)
	this.PutSqlData()
	//this.SendMsgSql()

	//判断房卡轮数
	this.GameRound++
	if this.GameRound-1 >= this.TableConfig.TotalRound { //轮数到了，游戏结束
		//系统定义结束
		this.GameState = GAME_STATUS_END
		this.BroadStageTime(TIMER_OVER_NUM)

		// for _, v := range this.Players {
		// 	if v.IsLeave == 1 {
		// 		continue
		// 	}
		// 	this.DeskMgr.LeaveDo(v.Uid)
		// }
		//开始归还桌子定时器
		//this.AddTimer(TIMER_OVER, TIMER_OVER_NUM, this.TimerOver, nil)
		this.TimerOver("")
	} else { //轮数没到，继续下一局
		this.InitGame()
		//重置玩家信息
		for _, v := range this.Players {
			this.ResetPlayer(v)
		}
		this.CoinList = []int64{}
		this.Round = 0
		//this.nextStage(GAME_STATUS_START)
		for _, v := range this.Players {
			v.IsReady = 0
		}
		this.GameState = STAGE_INIT
		this.BroadStageTime(0)
		this.nextStage(STAGE_INIT)
	}

}

func (this *ExtDesk) SendMsgSql() {
	//处理玩家数据
	if GetCostType() == 1 {
		this.PutSqlData()
	}

	//系统定义结束
	this.GameState = GAME_STATUS_END
	this.BroadStageTime(TIMER_OVER_NUM)
	if GetCostType() != 1 {
		for _, p := range this.Players {
			p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:      MSG_GAME_LEAVE_REPLY,
				Result:  0,
				Cid:     p.ChairId,
				Uid:     p.Uid,
				Robot:   p.Robot,
				NoToCli: true,
			})
		}
	}

	for _, v := range this.Players {
		if v.IsLeave == 1 {
			continue
		}
		this.DeskMgr.LeaveDo(v.Uid)
	}

	//this.Players = []*ExtPlayer{}
	//开始归还桌子定时器
	this.AddTimer(TIMER_OVER, TIMER_OVER_NUM, this.TimerOver, nil)
}

func (this *ExtDesk) TimerOver(d interface{}) {
	this.BroadcastAll(MSG_GAME_INFO_STAGE, &GSStageInfo{
		Id:        MSG_GAME_INFO_STAGE,
		Stage:     GAME_STATUS_END,
		StageTime: 0,
	})
	//发送总结算
	obj := new(GSLumpSum)
	obj.Id = MSG_GAME_INFO_LUMPSUM_READY
	//this.BroadcastAll(MSG_GAME_INFO_LUMPSUM_READY, obj)
	for _, v := range this.Players {
		for _, j := range this.Players {
			info := GSPlayerSum{
				ChairId: j.ChairId,
				Coin:    j.TotalCoins,
			}
			obj.Info = append(obj.Info, info)
		}
		v.SendNativeMsg(MSG_GAME_INFO_LUMPSUM_READY, obj)
		obj.Info = []GSPlayerSum{}
	}
	for _, p := range this.Players {
		p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:      MSG_GAME_LEAVE_REPLY,
			Result:  0,
			Cid:     p.ChairId,
			Uid:     p.Uid,
			Token:   p.Token,
			NoToCli: true,
		})
	}
	this.GameState = GAME_STATUS_FREE
	this.JuHao = ""
	//重置桌面属性
	this.CoinList = []int64{}
	this.Round = 0
	this.SettleContest = []PlayerContest{}
	this.ChairList = []int32{}

	// this.Players = []*ExtPlayer{}
	// this.CardMgr.InitCards()
	//房卡添加
	this.TableConfig = GATableConfig{}
	this.DisPlayer = []int32{}
	this.ClearTimer()
	//玩家离开
	this.GameOverLeave()
	this.DeskMgr.BackDesk(this)
}

//房主离开，所有人也离开
func (this *ExtDesk) HouseOwnerLeave() {
	logs.Debug("房间结束")
	//游戏结束
	// logs.Debug("游戏结束", this.Players)
	this.ClearTimer()
	this.GameState = GAME_STATUS_END
	this.BroadStageTime(0)

	//玩家离开
	for _, p := range this.Players {
		p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:      MSG_GAME_LEAVE_REPLY,
			Result:  0,
			Cid:     p.ChairId,
			Uid:     p.Uid,
			Token:   p.Token,
			NoToCli: true,
		})
	}
	this.GameOverLeave()
	//归还桌子
	this.GameState = GAME_STATUS_FREE
	this.ResetTable()
	this.DeskMgr.BackDesk(this)
}

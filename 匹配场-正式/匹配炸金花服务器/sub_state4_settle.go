package main

// import (
// 	"logs"
// )

func (this *ExtDesk) GameStateSettle() {
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

	if len(this.SettleContest) > 0 { //金币不足比牌-结算
		info.Count = len(this.SettleContest)
		// fmt.Println("结算 out of range", len(this.SettleContest))
		info.PContest = make([]Contest, 0, info.Count)
		for i := 0; i < info.Count; i++ {
			info.PContest = append(info.PContest, Contest{Person_1: this.SettleContest[i].Person_1, Person_2: this.SettleContest[i].Person_2, Winner: this.SettleContest[i].Winner})
		}
	}

	if this.Round >= GameRound {
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

			if result == 0 { //输家牌型变换
				allPlayer[i].CardType = 2
				msg.Winner = contestPlayer.ChairId
			} else {
				contestPlayer.CardType = 2
				msg.Winner = allPlayer[i].ChairId
				contestPlayer = allPlayer[i]
			}
			info.PContest = append(info.PContest, msg)
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
		if v.CardType != 2 {
			for i := 0; i < len(this.CoinList); i++ {
				sum += this.CoinList[i]
			}
			allsum := sum
			for i := 0; i < len(v.PayCoin); i++ {
				sum -= v.PayCoin[i]
			}
			// logs.Debug("抽水", this.Rate)
			v.RateCoins = float64(allsum*this.Bscore) * this.Rate
			v.WinCoins = (sum * this.Bscore) - int64(v.RateCoins)
			v.Coins += v.WinCoins
			info.CList = append(info.CList, CoinList{ChairId: v.ChairId, WinCoins: v.WinCoins, Coins: v.Coins})
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
			if v.Robot {
				AddLocalStock(v.WinCoins)
			}
			// fmt.Println("lost:", dsum, "cardLv:", v.CardLv, v.Coins)
		}
	}

	this.BroadcastAll(MSG_GAME_INFO_SETTLE, &info)
	// this.nextStage(STAGE_SETTLE)
	this.SendMsgSql()
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

	this.Players = []*ExtPlayer{}
	//开始归还桌子定时器
	this.AddTimer(TIMER_OVER, TIMER_OVER_NUM, this.TimerOver, nil)
}

func (this *ExtDesk) TimerOver(d interface{}) {
	// fmt.Println("run this")
	this.GameOver()
}

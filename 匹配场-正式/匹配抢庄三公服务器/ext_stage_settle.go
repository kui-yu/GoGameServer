package main

import (
	"fmt"
)

func (this *ExtDesk) StageSettle(d interface{}) {
	this.GameState = STAGE_GAME_SETTLE
	this.BroadStageTime(gameConfig.Stage_Game_Settle_Timer)
	//结算信息
	settle := GSettleInfo{
		Id: MSG_GAME_INFO_SETTLE_REPLY,
	}
	banker := this.Players[this.DeskBankerInfos.BankerId] //庄家
	bankerRes := PlayerResult{
		Uid: banker.Uid,
	}
	idleArr := append([]*ExtPlayer{}, this.Players[0:this.DeskBankerInfos.BankerId]...)
	idleArr = append(idleArr, this.Players[this.DeskBankerInfos.BankerId+1:]...)
	trend := Trend{}
	for _, v := range idleArr {
		playerRes := PlayerResult{Uid: v.Uid} //闲家输赢结构
		if v.HandCards.CardType > banker.HandCards.CardType {
			r := float64(v.HandCards.Multiple*this.BScore*v.BankerInfos.Multiple) * (1 - G_DbGetGameServerData.Rate) //扣除费率
			v.WinCoins = int64(r)                                                                                    //数据库记录
			playerRes.WinCoins += int64(r)                                                                           //闲家赢得金币
			v.Coins += playerRes.WinCoins                                                                            //更新玩家金币
			playerRes.Coins = v.Coins                                                                                //更新后的闲家金币
			bankerRes.WinCoins -= v.HandCards.Multiple * this.BScore * v.BankerInfos.Multiple                        //庄家输的金币
			//bankerRes.Coins -= v.HandCards.Multiple * this.BScore * v.BankerInfos.Multiple                           //更新后的庄家金币
			settle.IdleResult = append(settle.IdleResult, playerRes)
			//添加走势
			trend.WinCoins = v.WinCoins
			//v.AddTrend(trend)
			this.PutSqlData(v)
			continue
		} else if v.HandCards.CardType == banker.HandCards.CardType {
			if GetGongCardCount(v.HandCards.CardValue) > GetGongCardCount(banker.HandCards.CardValue) {
				r := float64(v.HandCards.Multiple*this.BScore*v.BankerInfos.Multiple) * (1 - G_DbGetGameServerData.Rate)
				v.WinCoins = int64(r)
				playerRes.WinCoins += int64(r)
				v.Coins += playerRes.WinCoins
				playerRes.Coins = v.Coins
				bankerRes.WinCoins -= v.HandCards.Multiple * this.BScore * v.BankerInfos.Multiple
				//bankerRes.Coins -= playerRes.WinCoins
				settle.IdleResult = append(settle.IdleResult, playerRes)
				//添加走势
				trend.WinCoins = v.WinCoins
				//v.AddTrend(trend)
				this.PutSqlData(v)
				continue
			} else if GetGongCardCount(v.HandCards.CardValue) == GetGongCardCount(banker.HandCards.CardValue) {
				if GetCradValue(v.HandCards.CardValue[0]) > GetCradValue(banker.HandCards.CardValue[0]) {
					r := float64(v.HandCards.Multiple*this.BScore*v.BankerInfos.Multiple) * (1 - G_DbGetGameServerData.Rate)
					v.WinCoins = int64(r)
					playerRes.WinCoins += int64(r)
					v.Coins += playerRes.WinCoins
					playerRes.Coins = v.Coins
					bankerRes.WinCoins -= v.HandCards.Multiple * this.BScore * v.BankerInfos.Multiple
					//bankerRes.Coins -= playerRes.WinCoins
					settle.IdleResult = append(settle.IdleResult, playerRes)
					//添加走势
					trend.WinCoins = v.WinCoins
					//v.AddTrend(trend)
					this.PutSqlData(v)
					continue
				} else if GetCradValue(v.HandCards.CardValue[0]) == GetCradValue(banker.HandCards.CardValue[0]) {
					if GetCardColor(v.HandCards.CardValue[0]) > GetCardColor(banker.HandCards.CardValue[0]) {
						r := float64(v.HandCards.Multiple*this.BScore*v.BankerInfos.Multiple) * (1 - G_DbGetGameServerData.Rate)
						v.WinCoins = int64(r)
						playerRes.WinCoins += int64(r)
						v.Coins += playerRes.WinCoins
						playerRes.Coins = v.Coins
						bankerRes.WinCoins -= v.HandCards.Multiple * this.BScore * v.BankerInfos.Multiple
						//bankerRes.Coins -= playerRes.WinCoins
						settle.IdleResult = append(settle.IdleResult, playerRes)
						//添加走势
						trend.WinCoins = v.WinCoins
						//v.AddTrend(trend)
						this.PutSqlData(v)
						continue
					} else {
						playerRes.WinCoins -= banker.HandCards.Multiple * this.BScore * v.BankerInfos.Multiple //输的金币
						v.WinCoins = playerRes.WinCoins                                                        //数据库记录
						v.Coins += playerRes.WinCoins                                                          //更新金币
						playerRes.Coins = v.Coins                                                              //更新后的闲家金币
						//r := float64(-playerRes.WinCoins) * (1 - G_DbGetGameServerData.Rate)                   //扣除费率
						bankerRes.WinCoins += -playerRes.WinCoins //int64(r)                                                         //庄家输的金币
						//bankerRes.Coins += bankerRes.WinCoins                                                  //更新后的庄家金币
						settle.IdleResult = append(settle.IdleResult, playerRes) //添加闲家结算记录
						//添加走势
						trend.WinCoins = v.WinCoins
						//v.AddTrend(trend)
						this.PutSqlData(v)
						continue
					}
				}
			}

		}
		playerRes.WinCoins -= banker.HandCards.Multiple * this.BScore * v.BankerInfos.Multiple //输的金币
		v.WinCoins = playerRes.WinCoins                                                        //数据库记录
		v.Coins += playerRes.WinCoins                                                          //更新金币
		playerRes.Coins = v.Coins                                                              //更新后的闲家金币
		//r := float64(-playerRes.WinCoins) * (1 - G_DbGetGameServerData.Rate)                   //扣除费率
		bankerRes.WinCoins += -playerRes.WinCoins //int64(r)                                                         //庄家输的金币
		//bankerRes.Coins += bankerRes.WinCoins                                                  //更新后的庄家金币
		settle.IdleResult = append(settle.IdleResult, playerRes) //添加闲家结算记录
		//添加走势
		trend.WinCoins = v.WinCoins
		//v.AddTrend(trend)
		this.PutSqlData(v) //请求后台，记录数据
	}
	if bankerRes.WinCoins > 0 {
		bankerRes.WinCoins = int64(float64(bankerRes.WinCoins) * float64(1-G_DbGetGameServerData.Rate))
	}
	banker.WinCoins = bankerRes.WinCoins //数据库记录
	banker.Coins += bankerRes.WinCoins
	bankerRes.Coins = banker.Coins
	settle.BankerResult = bankerRes //添加庄家结算记录
	//添加走势
	trend.WinCoins = banker.WinCoins
	trend.Player = 1
	//banker.AddTrend(trend)
	this.PutSqlData(banker)
	this.SettleInfo = settle                              //赋值结算消息，用于断线重连
	this.BroadcastAll(MSG_GAME_INFO_SETTLE_REPLY, settle) //发送结算记录
	this.runTimer(gameConfig.Stage_Game_Settle_Timer, this.GameEnd)
}

//数据库和大厅记录
func (this *ExtDesk) PutSqlData(p *ExtPlayer) {
	if GetCostType() != 2 {
		//数据库记录
		gameAndData := GGameEnd{
			Id:          MSG_GAME_END_NOTIFY,
			GameId:      GCONFIG.GameType,
			GradeId:     GCONFIG.GradeType,
			RoomId:      GCONFIG.RoomType,
			GameRoundNo: this.JuHao,
			Mini:        false,
			SetLeave:    1,
		}
		gameAndData.UserCoin = append(gameAndData.UserCoin, GGameEndInfo{
			UserId:      p.Uid,
			UserAccount: p.Account,
			BetCoins:    p.BankerInfos.Multiple,
			ValidBet:    p.BankerInfos.Multiple,
			PrizeCoins:  p.WinCoins,
			Robot:       p.Robot,
			WaterRate:   G_DbGetGameServerData.Rate,
		})
		p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, gameAndData)
		//大厅记录
		cards := make([][]int, 0)
		for k, v := range this.Players {
			if k == int(this.DeskBankerInfos.BankerId) {
				cards = append([][]int{v.HandCards.CardValue}, cards...)
			} else {
				cards = append(cards, v.HandCards.CardValue)
			}
		}
		gameHallData := GGameRecord{
			Id:             MSG_GAME_END_RECORD,
			GameId:         GCONFIG.GameType,
			GradeId:        GCONFIG.GradeType,
			RoomId:         GCONFIG.RoomType,
			GradeNumber:    1,
			GameRoundNo:    this.JuHao,
			PlayerCard:     p.HandCards.CardValue,
			SettlementCard: cards,
		}
		d := GGameRecordInfo{
			UserId:         p.Uid,
			UserAccount:    p.Account,
			Robot:          p.Robot,
			CoinsBefore:    p.Coins - p.WinCoins,
			BetCoins:       this.BScore * p.BankerInfos.Multiple,
			Coins:          p.Coins,
			CoinsAfter:     p.Coins - p.WinCoins,
			BankerMultiple: 1,
			BetMultiple:    p.BankerInfos.Multiple,
			Banker:         p.BankerInfos.IsBanker,
			BaseScore:      this.BScore,
			PrizeCoins:     p.WinCoins,
		}
		gameHallData.UserRecord = append(gameHallData.UserRecord, d)
		fmt.Println("游戏记录:", gameHallData)
		p.SendNativeMsgForce(MSG_GAME_END_RECORD, gameHallData)
	}
}

//结束
func (this *ExtDesk) GameEnd(d interface{}) {
	this.initDesk() //初始化桌子
	for _, v := range this.Players {
		v.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, GLeaveReply{
			Id:      MSG_GAME_LEAVE_REPLY,
			Result:  0,
			Cid:     v.ChairId,
			Uid:     v.Uid,
			Token:   v.Token,
			NoToCli: true,
		})
	}
	this.GameOverLeave()
	this.DeskMgr.BackDesk(this) //归还桌子
}

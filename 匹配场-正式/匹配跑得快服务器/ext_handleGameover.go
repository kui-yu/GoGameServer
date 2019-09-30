package main

import (
	"fmt"
	"logs"
)

func (this *ExtDesk) GameOver(p *ExtPlayer) {
	logs.Debug("进入结算阶段")
	this.TList = []*Timer{}
	this.GameState = GAME_STATUS_BALANCE
	this.BroadStageTime(0)
	var quanguan int32 = -1
	var quanguan2 int32 = -1
	var baopeiCoins int = 0

	//判断全关玩家
	for _, v := range this.Players {
		if len(v.HandCards) == 16 {
			v.IsQuanGUan = true
		}
	}
	// 遍历玩家集合,判断是否存在包赔玩家,全关玩家
	for _, v := range this.Players {
		if v.IsQuanGUan {
			if quanguan == -1 {
				quanguan = v.ChairId
			} else {
				quanguan2 = v.ChairId
			}
		}
	}
	// 保存总记录，用于大厅显示游戏记录
	gameEnd := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Mini:        false,
	}
	//保存游戏详情
	gameRecord := GameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		GradeNumber: 1,
	}
	for _, v := range this.Players {
		//计算玩家总输赢
		if v.Uid != p.Uid {
			if quanguan != -1 {
				if quanguan == v.ChairId || quanguan2 == v.ChairId {
					v.GetCoins -= len(v.HandCards) * this.Bscore * 2
				} else {
					v.GetCoins -= len(v.HandCards) * this.Bscore
				}
			} else {
				v.GetCoins -= len(v.HandCards) * this.Bscore
			}
		}
	}
	for _, v1 := range this.Players {
		if v1.Uid == p.Uid {
			for _, v := range this.Players {
				if v.Uid != v1.Uid {
					v1.GetCoins -= v.GetCoins
				}
			}
		}
	}
	//寻找不用包赔的id
	noBaopeis := []int32{}
	for _, v := range this.Players {
		if v.Uid != p.Uid {
			if !v.IsBaoPei {
				noBaopeis = append(noBaopeis, v.ChairId)
			}
		}
	}
	if len(noBaopeis) < 2 {
		logs.Debug("发现包赔玩家")
		//修正总结算
		for _, v := range this.Players {
			if v.Uid != p.Uid {
				if v.IsBaoPei {
					baopeiCoins = this.Players[noBaopeis[0]].GetCoins
					v.GetCoins += baopeiCoins
					this.Players[noBaopeis[0]].GetCoins = 0
				}
			}
		}
	}

	for _, v := range this.Players {
		this.BalanceResult[v.ChairId] = v.GetCoins
	}

	for _, v := range this.Players {
		logs.Debug("接下来判断各个玩家的炸弹情况:")
		fmt.Println("玩家:", v.Nick, "的炸弹数量是:", v.Booms, "结算金额是:", v.BoomBalance)
	}
	// //防止以小博大机制,玩家最多只能够赢取本身所含金币的最大值，如果玩家没有那么多钱，则金币清0
	// //判断正常结算
	// for _, v := range this.Players {
	// 	if p.Uid == v.Uid {
	// 		if v.Coins < int64(v.GetCoins) {
	// 			//如果玩家金币小于赢取金币，则默认赢取自身金币，并且将其余金币根据比例退还
	// 			bl := float32(v.Coins) / float32(v.GetCoins)
	// 			for _, v1 := range this.Players {
	// 				if v1.Uid != v.Uid {
	// 					v1.GetCoins = int(bl * float32(v1.GetCoins))
	// 				}
	// 			}
	// 			v.GetCoins = int(v.Coins)
	// 			break
	// 		}
	// 	}
	// }
	// for _, v := range this.Players {
	// 	//如果输家的金币不足以正常支付
	// 	if v.Coins < int64(-v.GetCoins) {
	// 		p.GetCoins += v.GetCoins
	// 		p.GetCoins += int(v.Coins)
	// 		v.GetCoins = int(v.Coins)
	// 	}
	// 	v.Coins = 0
	// }
	p.WinForMap[(p.ChairId+1)%int32(len(this.Players))] += (-this.Players[(p.ChairId+1)%int32(len(this.Players))].GetCoins)
	p.WinForMap[(p.ChairId+2)%int32(len(this.Players))] += (-this.Players[(p.ChairId+2)%int32(len(this.Players))].GetCoins)
	for _, v := range this.Players {
		if v.Uid != p.Uid {
			v.LoseForMap[p.ChairId] += -(v.GetCoins)
		}
	}
	for _, v := range this.Players {
		fmt.Println("玩家", v.Nick, "输给", this.Players[(v.ChairId+1)%int32(len(this.Players))].Nick, "的金币是", v.LoseForMap[(v.ChairId+1)%int32(len(this.Players))])
		fmt.Println("玩家", v.Nick, "输给", this.Players[(v.ChairId+2)%int32(len(this.Players))].Nick, "的金币是", v.LoseForMap[(v.ChairId+2)%int32(len(this.Players))])
		fmt.Println("玩家", v.Nick, "赢取", this.Players[(v.ChairId+1)%int32(len(this.Players))].Nick, "的金币是", v.WinForMap[(v.ChairId+1)%int32(len(this.Players))])
		fmt.Println("玩家", v.Nick, "赢取", this.Players[(v.ChairId+2)%int32(len(this.Players))].Nick, "的金币是", v.WinForMap[(v.ChairId+2)%int32(len(this.Players))])
	}
	for _, v := range this.Players {
		//炸弹另外结算
		v.GetCoins += v.BoomBalance
	}
	//玩家最多只能赢取自身金币数值一样的金币
	for _, v := range this.Players {
		if v.GetCoins > 0 && v.GetCoins > int(v.Coins) {
			logs.Debug("发现玩家金币不足，无法按照计划赚取金币:", v.Nick)
			bl := float32(v.Coins) / float32(v.GetCoins)
			if v.WinForMap[(v.ChairId+1)%int32(len(this.Players))] != 0 {
				v.WinForMap[(v.ChairId+1)%int32(len(this.Players))] = int(float32(v.WinForMap[(v.ChairId+1)%int32(len(this.Players))]) * bl)
			}
			if v.WinForMap[(v.ChairId+2)%int32(len(this.Players))] != 0 {
				v.WinForMap[(v.ChairId+2)%int32(len(this.Players))] = int(float32(v.WinForMap[(v.ChairId+2)%int32(len(this.Players))]) * bl)
			}
			for _, v1 := range this.Players {
				if v1.Uid != v.Uid {
					if v1.LoseForMap[v.ChairId] != 0 {
						v1.LoseForMap[v.ChairId] = int(float32(v1.LoseForMap[v.ChairId]) * bl)
					}
				}
			}
		}
	}
	fmt.Println("================================================")
	for _, v := range this.Players {
		fmt.Println("玩家", v.Nick, "输给", this.Players[(v.ChairId+1)%int32(len(this.Players))].Nick, "的金币是", v.LoseForMap[(v.ChairId+1)%int32(len(this.Players))])
		fmt.Println("玩家", v.Nick, "输给", this.Players[(v.ChairId+2)%int32(len(this.Players))].Nick, "的金币是", v.LoseForMap[(v.ChairId+2)%int32(len(this.Players))])
		fmt.Println("玩家", v.Nick, "赢取", this.Players[(v.ChairId+1)%int32(len(this.Players))].Nick, "的金币是", v.WinForMap[(v.ChairId+1)%int32(len(this.Players))])
		fmt.Println("玩家", v.Nick, "赢取", this.Players[(v.ChairId+2)%int32(len(this.Players))].Nick, "的金币是", v.WinForMap[(v.ChairId+2)%int32(len(this.Players))])
	}
	//根据玩家的两个集合判断最终的输赢情况
	for _, v := range this.Players {
		add := v.WinForMap[(v.ChairId+1)%int32(len(this.Players))] + v.WinForMap[(v.ChairId+2)%int32(len(this.Players))]
		jian := v.LoseForMap[(v.ChairId+1)%int32(len(this.Players))] + v.LoseForMap[(v.ChairId+2)%int32(len(this.Players))]
		v.GetCoins = add - jian
	}
	//玩家最多只能输自身金币数值一样的金币
	for _, v := range this.Players {
		if v.GetCoins < 0 && (-v.GetCoins) > int(v.Coins) {
			logs.Debug("发现玩家金币不足，无法按照计划扣除金币：", v.Nick)
			bl := float32(v.Coins) / float32(-v.GetCoins)
			if v.LoseForMap[(v.ChairId+1)%int32(len(this.Players))] != 0 {
				v.LoseForMap[(v.ChairId+1)%int32(len(this.Players))] = int(float32(v.LoseForMap[(v.ChairId+1)%int32(len(this.Players))]) * bl)
			}
			if v.LoseForMap[(v.ChairId+2)%int32(len(this.Players))] != 0 {
				v.LoseForMap[(v.ChairId+2)%int32(len(this.Players))] = int(float32(v.LoseForMap[(v.ChairId+2)%int32(len(this.Players))]) * bl)
			}
			for _, v1 := range this.Players {
				if v1.Uid != v.Uid {
					if v1.WinForMap[v.ChairId] != 0 {
						v1.WinForMap[v.ChairId] = int(float32(v1.WinForMap[v.ChairId]) * bl)
					}
				}
			}
		}
	}

	logs.Debug("==================================================")
	for _, v := range this.Players {
		fmt.Println("玩家", v.Nick, "输给", this.Players[(v.ChairId+1)%int32(len(this.Players))].Nick, "的金币是", v.LoseForMap[(v.ChairId+1)%int32(len(this.Players))])
		fmt.Println("玩家", v.Nick, "输给", this.Players[(v.ChairId+2)%int32(len(this.Players))].Nick, "的金币是", v.LoseForMap[(v.ChairId+2)%int32(len(this.Players))])
		fmt.Println("玩家", v.Nick, "赢取", this.Players[(v.ChairId+1)%int32(len(this.Players))].Nick, "的金币是", v.WinForMap[(v.ChairId+1)%int32(len(this.Players))])
		fmt.Println("玩家", v.Nick, "赢取", this.Players[(v.ChairId+2)%int32(len(this.Players))].Nick, "的金币是", v.WinForMap[(v.ChairId+2)%int32(len(this.Players))])
	}
	//根据玩家的两个集合判断最终的输赢情况
	for _, v := range this.Players {
		add := v.WinForMap[(v.ChairId+1)%int32(len(this.Players))] + v.WinForMap[(v.ChairId+2)%int32(len(this.Players))]
		jian := v.LoseForMap[(v.ChairId+1)%int32(len(this.Players))] + v.LoseForMap[(v.ChairId+2)%int32(len(this.Players))]
		v.GetCoins = add - jian
	}
	//储存到数据库
	for _, v := range this.Players {
		fmt.Println("结算完:", v.GetCoins, "名字为:", v.Nick)
		if v.GetCoins > 0 {
			v.WaterProft = float64(v.GetCoins) * G_DbGetGameServerData.Rate
			v.GetCoins -= int(v.WaterProft)
			logs.Debug("抽水率：", G_DbGetGameServerData.Rate)
			logs.Debug("抽了", v.WaterProft)
			v.Coins += int64(v.GetCoins)
		}

		gameEndInfo := GGameEndInfo{
			UserId:      v.Uid,
			UserAccount: v.Account,
			PrizeCoins:  int64(v.GetCoins),
			Robot:       v.Robot,
			WaterProfit: v.WaterProft,
			WaterRate:   G_DbGetGameServerData.Rate,
		}
		gameEnd.UserCoin = append(gameEnd.UserCoin, gameEndInfo)
		// 如果不是体验场，则发送给数据库
		if GetCostType() == 1 {
			v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, gameEnd)
		}
		if !v.Robot {
			gameRecordInfo := GameRecordInfo{
				UserId:          v.Uid,
				UserAccount:     v.Account,
				Robot:           v.Robot,
				CoinsBefore:     v.Coins - int64(v.GetCoins),
				PrizeCoins:      int64(v.GetCoins),
				CoinsAfter:      v.Coins,
				BaseScore:       G_DbGetGameServerData.Bscore,
				SurPlusCardsNum: len(v.HandCards),
				CoverBombNum:    v.BeBooms,
				BombNum:         v.Booms,
			}
			if v.ChairId == int32(baopeiCoins) {
				gameRecordInfo.CompensateNum = int64(baopeiCoins)
			} else {
				gameRecordInfo.CompensateNum = 0
			}
			gameRecord.UserRecord = append(gameRecord.UserRecord, gameRecordInfo)
			//如果不是体验场则发送游戏记录
			if GetCostType() == 1 {
				v.SendNativeMsgForce(MSG_GAME_END_RECORD, gameRecord)
			}
		}
		//清空
		gameEnd.UserCoin = []GGameEndInfo{}
		gameRecord.UserRecord = []GameRecordInfo{}
	}
	//发送结算信息
	balance := BalanceToClient{
		Id: MSG_GAME_INFO_BALANCE_BRO,
	}
	for _, v := range this.Players {
		PlayerMsg := PlayerMsgToBa{
			Booms:      v.Booms,
			Balance:    v.GetCoins,
			IsQuanGuan: v.IsQuanGUan,
			BaoPei:     v.IsBaoPei,
			Coins:      v.Coins,
		}
		v.HandCards = Sort(v.HandCards)
		for _, v1 := range v.HandCards {
			PlayerMsg.Handcards = append(PlayerMsg.Handcards, int(v1))
		}
		balance.PlayerMsgToBa = append(balance.PlayerMsgToBa, PlayerMsg)
	}
	this.BroadcastAll(MSG_GAME_INFO_BALANCE_BRO, balance)
	if GetCostType() == 2 {
		for _, v := range this.Players {
			//玩家离开
			v.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:      MSG_GAME_LEAVE_REPLY,
				Result:  0,
				Cid:     v.ChairId,
				Uid:     v.Uid,
				Token:   v.Token,
				Robot:   v.Robot,
				NoToCli: true,
			})
		}
	}
	this.Rest()
	this.GameOverLeave()
	this.GameState = GAME_STATUS_FREE
	this.JuHao = ""
	this.DeskMgr.BackDesk(this)
}

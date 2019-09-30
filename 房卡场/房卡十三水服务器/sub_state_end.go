package main

import (
	"logs"
	"time"
)

func (this *ExtDesk) GameStateSettle() {
	this.BroadStageTime(STAGE_SETTLE_TIME)
	//进入倒计时
	this.runTimer(STAGE_SETTLE_TIME, this.GameStateSettleEnd)
}

//阶段-结算
func (this *ExtDesk) GameStateSettleEnd(d interface{}) {

	//结算比大小
	for i := 0; i < len(this.Players)-1; i++ {
		for j := i + 1; j < len(this.Players); j++ {
			var specialScore1 int
			var specialScore2 int
			if this.Players[i].IsPlay == 2 {
				specialScore1 = getSpecialPoint(this.Players[i].SpecialType)
				this.Players[i].SpecialScore = specialScore1
			}
			if this.Players[j].IsPlay == 2 {
				specialScore2 = getSpecialPoint(this.Players[j].SpecialType)
				this.Players[j].SpecialScore = specialScore2
			}
			//比牌
			if this.Players[i].IsPlay == 2 && this.Players[j].IsPlay == 2 {

				//特殊牌型
				this.Players[i].WinCoinList[3] += specialScore1
				this.Players[j].WinCoinList[3] -= specialScore1
				//特殊牌型
				this.Players[j].WinCoinList[3] += specialScore2
				this.Players[i].WinCoinList[3] -= specialScore2

			} else {
				if this.Players[i].SpecialType > 0 && this.Players[i].IsPlay == 2 {
					//特殊牌型
					this.Players[i].WinCoinList[3] += specialScore1
					this.Players[j].WinCoinList[3] -= specialScore1
				} else {
					if this.Players[j].SpecialType > 0 && this.Players[j].IsPlay == 2 {
						//特殊牌型
						this.Players[j].WinCoinList[3] += specialScore2
						this.Players[i].WinCoinList[3] -= specialScore2
					} else {
						//普通比牌
						comparePoker(this.Players[i], this.Players[j])
					}
				}
			}
		}
	}

	//总得分
	for _, v := range this.Players {
		// logs.Debug("全垒打", len(v.ShootPlayers), len(this.Players))
		if len(v.ShootPlayers) == len(this.Players)-1 && len(this.Players) > 2 {
			// logs.Debug("------------------------全垒打-------------------------------")
			//全垒打
			v.IsAllWin = true
			//底分
			winCoin := v.WinCoinList[3] / (len(this.Players) - 1)

			v.WinCoinList[3] *= 2

			for _, v2 := range this.Players {
				if v != v2 {
					v2.WinCoinList[3] -= winCoin
				}
			}
		}

	}

	//土豪专场
	if this.TableConfig.GameModule == 2 {
		var tableMultiple int64 //桌面池子总倍数
		var tableMoney int64    //桌面池子总金额
		for _, v := range this.Players {
			winCoins := int64(v.WinCoinList[3] * this.Bscore)
			if winCoins < 0 {
				winCoins = -winCoins
				if winCoins > v.Coins {
					tableMoney += v.Coins
					v.WinCoins = -v.Coins
				} else {
					tableMoney += winCoins
					v.WinCoins = -winCoins
				}
			} else {
				tableMultiple += int64(v.WinCoinList[3])
			}
		}

		//平摊池子
		for _, v := range this.Players {
			if v.WinCoinList[3] > 0 {
				v.WinCoins = int64(int64(v.WinCoinList[3]) * tableMoney / tableMultiple)
				logs.Debug("输赢得分", v.WinCoins)
				//手续费
				v.RateCoins = float64(v.WinCoins) * this.Rate
				v.WinCoins = v.WinCoins - int64(v.RateCoins)
			}
		}
	} else {
		for _, v := range this.Players {
			v.WinCoins = int64(v.WinCoinList[3] * this.Bscore)
		}
	}

	//添加游戏记录
	for _, v := range this.Players {
		recordInfo := GSRecordInfo{
			WinCoins: v.WinCoins,
			WinDate:  time.Now().Format("2006-01-02 15:04:05"),
		}
		v.RecordInfos = append(v.RecordInfos, recordInfo)
	}

	//发送结算消息
	var settleInfos GSSettleInfos
	settleInfos.AllWinChairId = -1
	for _, v := range this.Players {
		var coins int64
		//土豪专场
		if this.TableConfig.GameModule == 2 {
			v.Coins += v.WinCoins
			coins = v.Coins
		} else {
			v.TotalCoins += v.WinCoins
			coins = v.TotalCoins
		}
		if v.IsAllWin {
			settleInfos.AllWinChairId = v.ChairId
		}

		info := GSettlePlayerInfo{
			Uid:          v.Uid,
			ChairId:      v.ChairId,
			WinCoinList:  v.WinCoinList,
			WinCoins:     v.WinCoins,
			ShootList:    v.ShootPlayers,
			Coins:        coins,
			SpecialScore: v.SpecialScore,
			NormalScores: v.NormalScores,
			ShootScoress: v.ShootScoress,
		}
		//牌型
		if v.SpecialType > 0 && v.IsPlay == 2 {
			//特殊牌型
			info.PlayCards = v.SpecialCards
			info.SpecialType = v.SpecialType
		} else {
			//正常牌型
			info.PlayCards = v.PlayCards
			info.SpecialType = 0
			info.NormalTypes = []int{v.PlayTypes[0], v.PlayTypes[1], v.PlayTypes[2]}
		}
		settleInfos.PlayInfo = append(settleInfos.PlayInfo, info)
	}

	settleInfos.Id = MSG_GAME_INFO_SETTLE_INFO_REPLY
	// logs.Debug("ces", settleInfos)
	var isLeave int32 = 0
	if this.TableConfig.TotalRound > this.Round {
		isLeave = 1
	}
	//数据交互
	this.PutSqlData(isLeave)
	//发送结算信息
	this.BroadcastAll(MSG_GAME_INFO_SETTLE_INFO_REPLY, &settleInfos)

	//游戏记录
	for _, v := range this.Players {
		recordInfo := GSRecordInfos{
			Id:    MSG_GAME_INFO_RECORD_INFO_REPLY,
			Infos: v.RecordInfos,
		}
		v.SendNativeMsg(MSG_GAME_INFO_RECORD_INFO_REPLY, &recordInfo)
	}

	if this.TableConfig.TotalRound <= this.Round {
		//结束
		this.TimerOver()
	} else {
		this.nextStage(STAGE_INIT)
	}
}

func (this *ExtDesk) TimerOver() {
	logs.Debug("房间结束")
	//总结算消息
	if this.Round > 0 {
		var rsInfoEnd GSSettleInfoEnd
		for _, v := range this.Players {
			var coins int64 = 0
			for _, c := range v.RecordInfos {
				coins += c.WinCoins
			}
			infoEnd := GSSettlePlayInfoEnd{
				ChairId:  v.ChairId,
				WinCoins: coins,
			}
			rsInfoEnd.PlayInfos = append(rsInfoEnd.PlayInfos, infoEnd)
		}
		rsInfoEnd.Id = MSG_GAME_INFO_SETTLE_INFO_END_REPLY
		this.BroadcastAll(MSG_GAME_INFO_SETTLE_INFO_END_REPLY, &rsInfoEnd)
	}
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

//正常牌型 position:0 头墩, 1 中墩, 2 底墩
func getNormalPoint(cardType int, position int) int {
	if cardType == NORMAL_FIVE_KIND {
		if position == 1 {
			return 20
		} else {
			return 10
		}
	} else if cardType == NORMAL_COLOR_STRAIGHT {
		if position == 1 {
			return 10
		} else {
			return 5
		}
	} else if cardType == NORMAL_FOUR_KIND {
		if position == 1 {
			return 8
		} else {
			return 4
		}
	} else if cardType == NORMAL_GOURD {
		if position == 1 {
			return 2
		}
	} else if cardType == NORMAL_THREE_KIND {
		if position == 0 {
			return 3
		}
	}
	return 1
}

//特殊牌型分数
func getSpecialPoint(cardType int) int {
	if cardType == SPECIAL_COLOR_DRAGON {
		return 108
	} else if cardType == SPECIAL_DRAGON {
		return 26
	} else if cardType == SPECIAL_THREE_SAME_COLOR_STRAIGHT {
		return 18
	} else if cardType == SPECIAL_THREE_FOUR_KIND {
		return 16
	} else if cardType == SPECIAL_FOUR_THREE_KIND {
		return 8
	} else if cardType == SPECIAL_SIX_PAIR {
		return 4
	} else if cardType == SPECIAL_THREE_STRAIGHT {
		return 4
	} else if cardType == SPECIAL_THREE_SAME_COLOR {
		return 4
	} else if cardType == SPECIAL_FAIL {
		return -3
	}
	return 0
}

//比牌
func comparePoker(player1 *ExtPlayer, player2 *ExtPlayer) {
	logs.Debug("比牌")
	var shootCnt1 int
	var shootCnt2 int
	winCoinList1 := []int{0, 0, 0}
	winCoinList2 := []int{0, 0, 0}
	//比较三墩
	for count := 0; count < len(player1.PlayTypes); count++ {
		var ct1, ct2 GCardsType
		if count == 0 {
			//头墩
			ct1 = GCardsType{
				Type:  player1.PlayTypes[0],
				Cards: ListGet(player1.PlayCards, 0, 3),
			}
			ct2 = GCardsType{
				Type:  player2.PlayTypes[0],
				Cards: ListGet(player2.PlayCards, 0, 3),
			}
		} else if count == 1 {
			//中墩
			ct1 = GCardsType{
				Type:  player1.PlayTypes[1],
				Cards: ListGet(player1.PlayCards, 3, 5),
			}
			ct2 = GCardsType{
				Type:  player2.PlayTypes[1],
				Cards: ListGet(player2.PlayCards, 3, 5),
			}
		} else if count == 2 {
			//底墩
			ct1 = GCardsType{
				Type:  player1.PlayTypes[2],
				Cards: ListGet(player1.PlayCards, 8, 5),
			}
			ct2 = GCardsType{
				Type:  player2.PlayTypes[2],
				Cards: ListGet(player2.PlayCards, 8, 5),
			}
		}
		//两人都不倒水
		if player1.SpecialType == SPECIAL_FAIL && player2.SpecialType == SPECIAL_FAIL {
			continue
		} else {
			logs.Debug("player1.SpecialType", player1.SpecialType)
			logs.Debug("player2.SpecialType", player2.SpecialType)
			if player2.SpecialType == SPECIAL_FAIL {
				shootCnt1++
				winCoinList1[count] += getNormalPoint(player1.PlayTypes[count], count)
				winCoinList2[count] -= getNormalPoint(player1.PlayTypes[count], count)
				continue
			} else if player1.SpecialType == SPECIAL_FAIL {
				shootCnt2++
				winCoinList1[count] -= getNormalPoint(player2.PlayTypes[count], count)
				winCoinList2[count] += getNormalPoint(player2.PlayTypes[count], count)
				continue
			}
			//比大小
			rs := MCardsCompare(ct1, ct2)
			if rs == 1 {
				shootCnt1++
				winCoinList1[count] += getNormalPoint(player1.PlayTypes[count], count)
				winCoinList2[count] -= getNormalPoint(player1.PlayTypes[count], count)
			} else if rs == 2 {
				shootCnt2++
				winCoinList1[count] -= getNormalPoint(player2.PlayTypes[count], count)
				winCoinList2[count] += getNormalPoint(player2.PlayTypes[count], count)
			}
		}
	}

	//更新结果
	for i := 0; i < 3; i++ {
		player1.WinCoinList[i] += winCoinList1[i]
		player2.WinCoinList[i] += winCoinList2[i]
	}

	winCoin1 := winCoinList1[0] + winCoinList1[1] + winCoinList1[2]
	winCoin2 := winCoinList2[0] + winCoinList2[1] + winCoinList2[2]
	//普通得分
	getNormalScores(player1, winCoinList1)
	getNormalScores(player2, winCoinList2)
	//打枪
	if shootCnt1 >= 3 {
		//打枪玩家 动画
		player1.ShootPlayers = append(player1.ShootPlayers, player2.ChairId)
		player1.ShootScoress = append(player1.ShootScoress, winCoinList1)
		//总得分
		player1.WinCoinList[3] = player1.WinCoinList[3] + winCoin1*2
		player2.WinCoinList[3] = player2.WinCoinList[3] + winCoin2*2
	} else if shootCnt2 >= 3 {
		//打枪玩家 动画
		player2.ShootPlayers = append(player2.ShootPlayers, player1.ChairId)
		player2.ShootScoress = append(player2.ShootScoress, winCoinList2)
		//总得分
		player1.WinCoinList[3] = player1.WinCoinList[3] + winCoin1*2
		player2.WinCoinList[3] = player2.WinCoinList[3] + winCoin2*2
	} else {
		//总得分
		player1.WinCoinList[3] = player1.WinCoinList[3] + winCoin1
		player2.WinCoinList[3] = player2.WinCoinList[3] + winCoin2
	}
}

//玩家结果普通得分
func getNormalScores(p *ExtPlayer, normalScores []int) {
	p.NormalScores[0] += normalScores[0]
	p.NormalScores[1] += normalScores[1]
	p.NormalScores[2] += normalScores[2]
}

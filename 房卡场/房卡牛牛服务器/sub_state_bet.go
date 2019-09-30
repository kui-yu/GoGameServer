package main

func (this *ExtDesk) GameStateBet() {

	this.BroadStageTime(STAGE_CALL_SCORE_TIME)
	//判断玩家可下倍数
	var betList []int
	if this.TableConfig.GameType == 1 {
		//抢庄
		if this.Players[this.Banker].CallMultiple == 200 {
			betList = []int{1, 3, 6, 10, 13}
		} else if this.Players[this.Banker].CallMultiple == 100 {
			betList = []int{1, 2, 3, 6}
		} else {
			betList = []int{1, 2, 3}
		}

	} else {
		//通比
		betList = []int{1, 2, 3, 4, 5}
	}
	for _, v := range this.Players {
		var playerBet []int
		for _, bet := range betList {
			if v.Coins >= int64(bet*this.TableConfig.BaseScore*5) && this.TableConfig.GameModule == 2 {
				playerBet = append(playerBet, bet)
			} else {
				playerBet = append(playerBet, bet)
			}
		}
		v.SendNativeMsg(MSG_GAME_INFO_BET_LIST, &GCallListMsg{
			Id:         MSG_GAME_INFO_BET_LIST,
			BetListCnt: len(playerBet),
			BetList:    playerBet,
		})
		v.PlayerBets = playerBet
	}

	//进入倒计时
	this.runTimer(STAGE_CALL_SCORE_TIME, this.GameStateBetEnd)
}

//阶段-叫分结束
func (this *ExtDesk) GameStateBetEnd(d interface{}) {
	// logs.Debug("叫分阶段结束", len(this.Players))

	for _, v := range this.Players {
		if v.BetMultiple == 0 && v.ChairId != this.Banker {
			v.BetMultiple = 1
			//回复
			re := GSCallMsg{
				Id:       MSG_GAME_INFO_CALL_REPLY,
				ChairId:  v.ChairId,
				Multiple: v.BetMultiple,
			}
			//广播
			this.BroadcastAll(MSG_GAME_INFO_CALL_REPLY, &re)
		}
	}

	//进入发牌
	this.nextStage(STAGE_DEAL)
}

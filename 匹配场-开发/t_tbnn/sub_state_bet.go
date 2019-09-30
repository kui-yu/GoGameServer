package main

func (this *ExtDesk) GameStateBet() {

	this.BroadStageTime(TIME_STAGE_CALL_NUM)

	bets := []int{1, 2, 3, 4, 5}
	for _, v := range this.Players {
		var playerBet []int
		if GetCostType() == 1 { //如果不是体验场再进行玩家金币下注能力判断
			for _, bet := range bets {
				if v.Coins >= int64(this.Bscore*5*bet) {
					playerBet = append(playerBet, bet)
				}
			}
		} else { //如果是体验场就不限制
			playerBet = bets
		}

		v.SendNativeMsg(MSG_GAME_INFO_BET_LIST, &GCallListMsg{
			Id:         MSG_GAME_INFO_BET_LIST,
			BetListCnt: len(playerBet),
			BetList:    playerBet,
		})
		v.PlayerBets = playerBet
	}

	//进入倒计时
	this.runTimer(TIME_STAGE_CALL_NUM, this.GameStateBetEnd)
}

//叫分阶段-结束
func (this *ExtDesk) GameStateBetEnd(d interface{}) {
	// logs.Debug("叫分阶段结束", len(this.Players))

	for _, v := range this.Players {
		if v.BetMultiple == 0 {
			v.BetMultiple = 1
			//回复
			re := GCallMsgReply{
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

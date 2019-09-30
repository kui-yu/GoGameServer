package main

func (this *ExtDesk) GameStateBet() {

	this.BroadStageTime(TIME_STAGE_BET_NUM)
	//进入倒计时
	this.runTimer(TIME_STAGE_BET_NUM, this.GameStateBetEnd)
}

//叫分阶段-结束
func (this *ExtDesk) GameStateBetEnd(d interface{}) {

	for _, v := range this.Players {
		if v.BetMultiple == 0 && v.ChairId != this.Banker {
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

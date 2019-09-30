package main

func (this *ExtDesk) GameStateDismiss() {
	this.BroadStageTime(STAGE_DISMISS_TIME)
	//进入倒计时
	this.runTimer(STAGE_DISMISS_TIME, this.GameStateDismissEnd)
}

//阶段-解散
func (this *ExtDesk) GameStateDismissEnd(d interface{}) {
	
	this.BroadcastAll(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
		Id:        MSG_GAME_INFO_DISMISS_REPLY,
		DisPlayer: this.DisPlayer,
		IsDismiss: 3,
	})
	//同意人数已满，解散
	//房间结束
	this.TimerOver()
}

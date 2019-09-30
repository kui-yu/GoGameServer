package main

func (this *ExtDesk) GameStateRestart() {
	//结束阶段-返回开始阶段
	this.BroadStageTime(STAGE_RESTART_TIME)
	this.runTimer(STAGE_RESTART_TIME, this.GameStateRestartEnd)
}

//阶段-重新开始
func (this *ExtDesk) GameStateRestartEnd(d interface{}) {
	//进入游戏开始阶段
	this.nextStage(GAME_STATUS_START)
}

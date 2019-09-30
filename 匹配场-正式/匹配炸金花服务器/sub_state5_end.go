package main

// import "logs"

func (this *ExtDesk) GameStateEnd() {
	// this.BroadStageTime(STAGE_END_TIME)
	this.runTimer(STAGE_END_TIME, this.HandleGameEnd)
}

func (this *ExtDesk) HandleGameEnd(d interface{}) {
	// logs.Debug("游戏结束")
	if this.GameState != GAME_STATUS_END {
		return
	}
	this.nextStage(STAGE_SETTLE)
}

package main

func (this *ExtDesk) GameStateStart() {
	this.BroadStageTime(STAGE_START_TIME)
	//开始发牌-广播玩家
	this.DealPoker()
	//进入倒计时
	this.runTimer(STAGE_START_TIME, this.GameStateStartEnd)
}

//阶段-开始
func (this *ExtDesk) GameStateStartEnd(d interface{}) {
	//进入玩牌阶段
	this.nextStage(STAGE_PLAY)
}

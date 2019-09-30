package main

//开始阶段
func (this *ExtDesk) GameStart(d interface{}) {
	this.Stage = STAGE_GAME_START
	this.BroadStageTime(gameConfigInfo.Start_Timer)
	this.runTimer(gameConfigInfo.Start_Timer, this.GameStageBet)
}

//下注阶段
func (this *ExtDesk) GameStageBet(d interface{}) {
	this.Stage = STAGE_GAME_BET
	//群发阶段消息=>下注阶段
	this.BroadStageTime(gameConfigInfo.Bet_Timer)
	//下一阶段=>开牌阶段
	this.runTimer(gameConfigInfo.Bet_Timer, this.GameStageStopBet)
}

//停止下注阶段
func (this *ExtDesk) GameStageStopBet(d interface{}) {
	this.Stage = STAGE_GAME_STOP_BET
	//群发阶段消息=>开牌阶段
	this.BroadStageTime(gameConfigInfo.Stop_Bet_Timer)
	this.runTimer(gameConfigInfo.Stop_Bet_Timer, this.GameStageSendCard)
}

//发牌阶段
func (this *ExtDesk) GameStageSendCard(d interface{}) {
	this.Stage = STAGE_GAME_SEND_CARD
	this.BroadStageTime(gameConfigInfo.Send_Card_Timer)
	this.runTimer(gameConfigInfo.Send_Card_Timer, this.GameStageOpenCard)
}

package main

//开始阶段
func (this *ExtDesk) GameStart(i interface{}) {
	//初始化
	this.DeskBankerInfos.BankerId = -1
	this.GameState = STAGE_START
	this.BroadStageTime(gameConfig.Stage_Start_Timer)
	this.DeskCards = DisturbCards(InitCards()) //初始化卡组并打乱
	this.runTimer(gameConfig.Stage_Start_Timer, this.ShuffleCards)
}

//洗牌阶段
func (this *ExtDesk) ShuffleCards(d interface{}) {
	this.GameState = STAGE_SHUFFLE_CARDS
	this.BroadStageTime(gameConfig.Stage_Shuffle_Cards_Timer)
	this.runTimer(gameConfig.Stage_Shuffle_Cards_Timer, this.StageSendCards)
}

//发牌阶段
func (this *ExtDesk) StageSendCards(d interface{}) {
	this.GameState = STAGE_SEND_CARDS
	this.BroadStageTime(gameConfig.Stage_Send_Cards_Timer)
	this.SendCards() //发牌
	this.runTimer(gameConfig.Stage_Send_Cards_Timer, this.StageChoiceBanker)
}

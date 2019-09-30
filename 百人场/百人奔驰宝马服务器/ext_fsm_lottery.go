package main

type FSMLottery struct {
	Mark int

	EDesk       *ExtDesk
	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMLottery) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMLottery) Run() {
	DebugLog("游戏状态-开奖")

	this.EndDateTime = GetTimeMS() + int64(gameConfig.StateInfo.RunCarLogoTime)

	this.addListen() // 添加监听
	this.EDesk.GameState = GAME_STATUS_LOTTERY
	this.EDesk.SendGameState(GAME_STATUS_LOTTERY, int64(gameConfig.StateInfo.RunCarLogoTime))

	this.EDesk.AddTimer(GAME_STATUS_LOTTERY, gameConfig.StateInfo.RunCarLogoTime/1000, this.TimerCall, nil)

	this.RunLogo()
}

func (this *FSMLottery) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_BALANCE)
}

func (this *FSMLottery) GetMark() int {
	return this.Mark
}

func (this *FSMLottery) Leave() {
	this.removeListen()
}

func (this *FSMLottery) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FSMLottery) addListen() {}

func (this *FSMLottery) removeListen() {}

func (this *FSMLottery) RunLogo() {
	result := this.EDesk.allotCard()
	this.EDesk.GameResult = result
	info := GNLottery{
		Id:  MSG_GAME_INFO_NLOTTERY,
		Car: result,
	}
	DebugLog("开奖结果：", result)
	this.EDesk.BroadcastAll(MSG_GAME_INFO_NLOTTERY, &info)
}

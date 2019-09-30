package main

type FSMReady struct {
	Mark int

	EDesk       *ExtDesk
	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMReady) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMReady) Run() {
	DebugLog("游戏状态-准备")

	this.EndDateTime = GetTimeMS() + int64(gameConfig.StateInfo.RunCarLogoTime)

	this.addListen() // 添加监听
	this.EDesk.GameState = GAME_STATUS_READY
	this.EDesk.SendGameState(GAME_STATUS_READY, int64(gameConfig.StateInfo.ReadyTime))
	this.EDesk.ResetExtDesk()
	this.EDesk.AddTimer(GAME_STATUS_READY, gameConfig.StateInfo.ReadyTime/1000, this.TimerCall, nil)

	this.clear()
}

func (this *FSMReady) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_DOWNBET)
}

func (this *FSMReady) GetMark() int {
	return this.Mark
}

func (this *FSMReady) Leave() {
	this.removeListen()
}

func (this *FSMReady) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FSMReady) addListen() {}

func (this *FSMReady) removeListen() {}

func (this *FSMReady) clear() {
	for _, v := range this.EDesk.Players {
		v.PAreaCoins = []int64{0, 0, 0, 0, 0, 0, 0, 0}
	}
}

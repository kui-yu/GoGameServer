package main

// 下注
func (this *ExtDesk) TimerBet(d interface{}) {
	this.Lock()
	defer this.Unlock()

	sd := GGameBetNotify{
		Id:    MSG_GAME_INFO_BET_NOTIFY,
		Timer: int32(gameConfig.Timer.BetNum) * 1000,
	}

	this.BroadcastAll(MSG_GAME_INFO_BET_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_BET_NOTIFY
	this.AddTimer(gameConfig.Timer.Bet, gameConfig.Timer.BetNum, this.TimerStopBet, nil)

	// 定时广播
	this.AddTimer(gameConfig.Timer.NewBet, gameConfig.Timer.NewBetNum, this.HandleTimeOutBet, nil)
}

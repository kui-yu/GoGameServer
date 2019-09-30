package main

// 洗牌
func (this *ExtDesk) TimerShuffle(d interface{}) {
	this.Lock()
	defer this.Unlock()

	this.ExtDeskInit()

	sd := GGameShuffleNotify{
		Id:    MSG_GAME_INFO_SHUFFLE_NOTIFY,
		Timer: int32(gameConfig.Timer.ShuffleNum) * 1000,
	}

	this.BroadcastAll(MSG_GAME_INFO_SHUFFLE_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_SHUFFLE_NOTIFY
	this.AddTimer(gameConfig.Timer.Shuffle, gameConfig.Timer.ShuffleNum, this.TimerReady, nil)
}

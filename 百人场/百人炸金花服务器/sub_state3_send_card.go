package main

// 发牌
func (this *ExtDesk) TimerSendCard(d interface{}) {
	this.Lock()
	defer this.Unlock()

	// 东、西、南、北  发牌两张
	for j := 0; j < 4; j++ {
		this.IdleCard[j] = append(this.IdleCard[j], this.CardMgr.SendCard(2)...)
	}
	this.BankerCard = append(this.BankerCard, this.CardMgr.SendCard(2)...)

	sd := GGameSendCardNotify{
		Id:         MSG_GAME_INFO_SEND_NOTIFY,
		Timer:      int32(gameConfig.Timer.SendCardNum) * 1000,
		BankerCard: this.BankerCard,
		IdleCard:   this.IdleCard,
	}

	this.BroadcastAll(MSG_GAME_INFO_SEND_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_SEND_NOTIFY
	this.AddTimer(gameConfig.Timer.SendCard, gameConfig.Timer.SendCardNum, this.TimerBet, nil)
}

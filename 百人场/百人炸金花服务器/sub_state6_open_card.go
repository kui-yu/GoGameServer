package main

import (
	"fmt"
)

// 开牌
func (this *ExtDesk) TimerOpen(d interface{}) {
	this.Lock()
	defer this.Unlock()

	sd := GGameOpenNotify{
		Id:          MSG_GAME_INFO_OPEN_NOTIFY,
		Timer:       int32(gameConfig.Timer.OpenNum) * 1000,
		BankerCard:  this.BankerCard,
		MBankerCard: this.MBankerCard,
		TBankerCard: this.TBankerType,
		IdleCard:    this.IdleCard,
		MIdleCard:   this.MIdleCard,
		TypeList:    this.TypeList,
	}

	this.BroadcastAll(MSG_GAME_INFO_OPEN_NOTIFY, sd)
	fmt.Printf("开牌阶段=>牌1:%v,牌2:%v\n", this.BankerCard, this.IdleCard)
	this.GameState = MSG_GAME_INFO_OPEN_NOTIFY
	this.AddTimer(gameConfig.Timer.Open, gameConfig.Timer.OpenNum, this.TimerAward, nil)
}

package main

import (
	"logs"

	"bl.com/paigow"
)

// 开牌
func (this *ExtDesk) TimerOpen(d interface{}) {
	this.Lock()
	defer this.Unlock()
	for _, v := range this.Players {
		if v.LiXian {
			logs.Debug("--------开牌阶段有用户离线：%v", v.Nick)
		}
	}
	this.BankerCard = paigow.Sort(this.BankerCard)
	for i := 0; i < 3; i++ {
		this.IdleCard[i] = paigow.Sort(this.IdleCard[i])
	}

	bType := paigow.GetCardsType(this.BankerCard)
	this.TypeList = append(this.TypeList, bType)
	for i := 0; i < 3; i++ {
		pType := paigow.GetCardsType(this.IdleCard[i])
		this.TypeList = append(this.TypeList, pType)
	}

	sd := GGameOpenNotify{
		Id:         MSG_GAME_INFO_OPEN_NOTIFY,
		Timer:      int32(gameConfig.Timer.OpenNum) * 1000,
		BankerCard: this.BankerCard,
		IdleCard:   this.IdleCard,
		TypeList:   this.TypeList,
	}

	this.BroadcastAll(MSG_GAME_INFO_OPEN_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_OPEN_NOTIFY
	this.AddTimer(gameConfig.Timer.Open, gameConfig.Timer.OpenNum, this.TimerAward, nil)
}

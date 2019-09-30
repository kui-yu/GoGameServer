package main

import (
	"logs"
)

// 发牌
func (this *ExtDesk) TimerSendCard(d interface{}) {
	this.Lock()
	defer this.Unlock()
	for _, v := range this.Players {
		if v.LiXian {
			logs.Debug("++++++发牌阶段有用户离线：%v", v.Nick)
		}
	}
	sd := GGameSendCardNotify{
		Id:    MSG_GAME_INFO_SEND_NOTIFY,
		Timer: int32(gameConfig.Timer.SendCardNum) * 1000,
	}

	this.BroadcastAll(MSG_GAME_INFO_SEND_NOTIFY, sd)

	this.GameState = MSG_GAME_INFO_SEND_NOTIFY
	this.AddTimer(gameConfig.Timer.SendCard, gameConfig.Timer.SendCardNum, this.TimerBet, nil)
}

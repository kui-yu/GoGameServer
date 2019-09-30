package main

import (
	"encoding/json"
	// "logs"
	"math/rand"
	"time"
)

//抢地主 叫分阶段
func (this *ExtDesk) TimerCall(nextplayer *ExtPlayer) {

	if nextplayer.Robot {
		//机器人叫分,随机秒数
		rand.Seed(time.Now().UnixNano())
		timerCallNum := rand.Intn(3) + 2

		this.AddTimer(TIMER_CALL, timerCallNum, this.robotCallBank, nextplayer)
	} else {
		//等待真实玩家叫分
		this.AddTimer(TIMER_CALL, TIMER_CALL_NUM, this.TimerCallEnd, nil)
	}
}

//玩家抢地主，叫分阶段结束
func (this *ExtDesk) TimerCallEnd(d interface{}) {
	data := GCallMsg{
		Coins: -1,
	}
	p := this.Players[this.CurCid]
	dv, _ := json.Marshal(data)
	this.HandleGameCall(p, &DkInMsg{
		Uid:  p.Uid,
		Data: string(dv),
	})
}

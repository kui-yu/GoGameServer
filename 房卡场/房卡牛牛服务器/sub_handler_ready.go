package main

import (
	"encoding/json"
)

//玩家准备
func (this *ExtDesk) HandleReady(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != GAME_STATUS_FREE && this.GameState != STAGE_INIT {
		return
	}
	data := GAPlayerReady{}
	json.Unmarshal([]byte(d.Data), &data)

	p.IsReady = data.IsReady

	this.BroadcastAll(MSG_GAME_INFO_READY_REPLY, &GSPlayerReady{
		Id:      MSG_GAME_INFO_READY_REPLY,
		ChairId: p.ChairId,
		IsReady: p.IsReady,
	})

	flag := true
	for _, v := range this.Players {
		if v.IsReady == 0 {
			flag = false
		}
	}

	if flag && len(this.Players) >= this.TableConfig.PlayerNumber {
		//进入游戏开始
		this.nextStage(GAME_STATUS_START)
	}
}

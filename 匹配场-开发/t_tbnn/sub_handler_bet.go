package main

import (
	"encoding/json"
	// "logs"
)

//玩家叫分
func (this *ExtDesk) HandleGameCall(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != GAME_STATUS_START {
		return
	}
	//已叫分
	if p.BetMultiple > 0 {
		return
	}

	data := GCallMsg{}
	err := json.Unmarshal([]byte(d.Data), &data)
	if err != nil {
		return
	}

	//下注倍数{1,2,3,4,5}
	multFlag := true
	for _, bet := range p.PlayerBets {
		if bet == data.Multiple {
			multFlag = false
			break
		}
	}
	if multFlag {
		return
	}

	//倍数
	p.BetMultiple = data.Multiple
	//回复
	re := GCallMsgReply{
		Id:       MSG_GAME_INFO_CALL_REPLY,
		ChairId:  p.ChairId,
		Multiple: data.Multiple,
	}
	//标志
	flag := true
	//广播
	this.BroadcastAll(MSG_GAME_INFO_CALL_REPLY, &re)

	for _, v := range this.Players {
		if v.BetMultiple == 0 {
			flag = false
		}
	}

	if flag {
		// logs.Debug("都叫分")
		this.nextStage(STAGE_DEAL)
	}
}

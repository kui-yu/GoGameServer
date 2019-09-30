package main

import (
	"encoding/json"
)

//玩家叫分
func (this *ExtDesk) HandleGameBet(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != STAGE_BET {
		return
	}
	//庄家不叫分
	if p.ChairId == this.Banker {
		return
	}
	//玩家已叫分
	if p.BetMultiple > 0 || this.Banker == p.ChairId {
		return
	}

	data := GCallMsg{}
	err := json.Unmarshal([]byte(d.Data), &data)
	if err != nil {
		return
	}

	betFlag := true
	for _, bet := range p.PlayerBets {
		if bet == data.Multiple {
			betFlag = false
			break
		}
	}
	if betFlag {
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
		if v.BetMultiple == 0 && v.ChairId != this.Banker {
			flag = false
		}
	}

	if flag {
		//进入发牌
		this.nextStage(STAGE_DEAL)
	}

	return
}

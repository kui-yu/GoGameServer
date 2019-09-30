package main

import (
	"encoding/json"
)

//玩家下注
func (this *ExtDesk) HandleGameBet(p *ExtPlayer, d *DkInMsg) {
	//不是下注阶段
	if this.GameState != STAGE_BET {
		return
	}
	//玩家已下注
	if p.PlayMultiple > -1 || this.Banker == p.ChairId {
		return
	}
	//解析指令
	data := GAPlayerPlayInfo{}
	err := json.Unmarshal([]byte(d.Data), &data)
	//解析失败
	if err != nil {
		return
	}

	betFlag := true
	for _, bet := range p.PlayerBets {
		if bet == data.PlayMultiple {
			betFlag = false
			break
		}
	}
	if betFlag {
		return
	}

	p.PlayMultiple = data.PlayMultiple

	//下注返回
	info := GSPlayerPlayInfo{
		Id:           MSG_GAME_INFO_PLAY_INFO_REPLY,
		ChairId:      p.ChairId,
		PlayMultiple: p.PlayMultiple,
	}
	this.BroadcastAll(MSG_GAME_INFO_PLAY_INFO_REPLY, &info)

	flag := true
	for _, v := range this.Players {
		if v.PlayMultiple == -1 && v.ChairId != this.Banker {
			flag = false
		}
	}

	if flag {
		//发牌游戏开始
		this.nextStage(STAGE_DEAL)
	}
}

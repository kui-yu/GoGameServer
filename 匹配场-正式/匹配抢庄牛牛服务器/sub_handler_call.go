package main

import (
	"encoding/json"
)

//玩家抢庄
func (this *ExtDesk) HandleGameCall(p *ExtPlayer, d *DkInMsg) {

	if this.GameState != GAME_STATUS_START {
		return
	}

	//已抢庄
	if p.CallBankFlag {
		return
	}
	//解析
	data := GCallMsg{}
	err := json.Unmarshal([]byte(d.Data), &data)
	//解析失败
	if err != nil {
		return
	}

	//{不抢,50,100,200}
	callFlag := true
	for _, bet := range p.PlayerCalls {
		if bet == data.Multiple {
			callFlag = false
			break
		}
	}
	if callFlag {
		return
	}

	//抢庄倍数
	p.CallMultiple = data.Multiple
	p.CallBankFlag = true

	//回复
	re := GCallMsgReply{
		Id:       MSG_GAME_INFO_CALL_BANKER_NOTIFY,
		ChairId:  p.ChairId,
		Multiple: data.Multiple,
	}

	//广播
	this.BroadcastAll(MSG_GAME_INFO_CALL_BANKER_NOTIFY, &re)

	//全部抢完庄
	var count int
	for _, v := range this.Players {
		if v.CallBankFlag {
			count++
		}
	}
	//广播庄家
	if count >= len(this.Players) {
		//选庄
		this.ChooseBanker()
		//下个阶段
		this.nextStage(STAGE_BET)
	}
	return
}

package main

import (
	"encoding/json"
)

//玩家抢庄
func (this *ExtDesk) HandleCallBank(p *ExtPlayer, d *DkInMsg) {

	if this.GameState != STAGE_CALL_BANKER {
		//返回Err
		return
	}

	data := GACallMsg{}
	err := json.Unmarshal([]byte(d.Data), &data)
	if err != nil {
		return
	}
	// logs.Debug("抢庄", d)
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

	//回复
	re := GSCallMsg{
		Id:       MSG_GAME_INFO_CALL_BANKER_NOTIFY,
		ChairId:  p.ChairId,
		Multiple: data.Multiple,
	}

	//广播
	this.BroadcastAll(MSG_GAME_INFO_CALL_BANKER_NOTIFY, &re)

	//全部抢完庄
	var count int
	for _, v := range this.Players {
		if v.CallMultiple >= 0 {
			count++
		}
	}
	//广播庄家
	if count >= len(this.Players) {

		//选庄
		this.ChooseBanker()

		//下个阶段
		this.nextStage(STAGE_CALL_SCORE)
	}
	return
}

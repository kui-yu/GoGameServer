package main

import (
	"encoding/json"
)

func (this *ExtDesk) ShowCards(p *ExtPlayer, d *DkInMsg) {
	req := GShowCardsREQ{}
	_ = json.Unmarshal([]byte(d.Data), req)
	res := GShowCardsRES{
		Id:  MSG_GAME_INFO_SHOW_CARDS_REPLY,
		Uid: p.Uid,
	}
	if this.GameState != STAGE_OPEN_CARDS {
		res.Err = "不是开牌阶段"
		res.Result = 1
		p.SendNativeMsg(MSG_GAME_INFO_SHOW_CARDS_REPLY, res)
		return
	}
	if p.IsOpenCards {
		res.Err = "已经亮牌了，不要重复操作"
		res.Result = 2
		p.SendNativeMsg(MSG_GAME_INFO_SHOW_CARDS_REPLY, res)
		return
	}
	//正常亮牌
	p.IsOpenCards = true
	res.Cards = p.HandCards
	this.BroadcastAll(MSG_GAME_INFO_SHOW_CARDS_REPLY, res)
	count := 0
	for _, v := range this.Players {
		if v.IsOpenCards {
			count++
		}
	}
	//所有玩家都亮牌了，进入结算
	if count == len(this.Players) {
		this.DelTimer(3)
		this.StageSettle("")
	}
}

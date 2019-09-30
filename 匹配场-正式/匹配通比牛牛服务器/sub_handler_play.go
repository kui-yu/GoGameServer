package main

//玩牌
func (this *ExtDesk) HandlePlayCard(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != STAGE_PLAY {
		return
	}
	if p.IsLook {
		return
	}
	p.IsLook = true

	niuHand := GHandNiuReply{
		Id:       MSG_GAME_INFO_PLAY_REPLY,
		ChairId:  p.ChairId,
		NiuPoint: p.NiuPoint,
		NiuCards: p.NiuCards,
	}

	this.BroadcastAll(MSG_GAME_INFO_PLAY_REPLY, niuHand)

	flag := true
	for _, v := range this.Players {
		if !v.IsLook {
			flag = false
		}
	}
	if flag {
		//结算
		this.nextStage(STAGE_SETTLE)
	}
}

package main

func (this *ExtDesk) HandleTips(p *ExtPlayer, d *DkInMsg) {

	var tips [][]byte

	if this.MaxChuPai != nil {
		tips = this.CalcTips(this.MaxChuPai.Cards, p.HandCard)
	} else {
		tips = this.CalcFirstTips(p.HandCard)
	}

	p.SendNativeMsg(MSG_GAME_INFO_TIPS_REPLY, &GGameTips{
		Id:   MSG_GAME_INFO_TIPS_REPLY,
		Tips: tips,
	})

}

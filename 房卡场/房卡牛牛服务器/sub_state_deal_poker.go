package main

func (this *ExtDesk) GameStateDeal() {
	this.GameDeal()
	this.BroadStageTime(STAGE_DEAL_POKER_TIME)
	//进入倒计时
	this.runTimer(STAGE_DEAL_POKER_TIME, this.GameStateDealEnd)
}

//阶段-发牌发牌结束
func (this *ExtDesk) GameStateDealEnd(d interface{}) {
	//进入玩牌
	this.nextStage(STAGE_PLAY)
}

//发牌
func (this *ExtDesk) GameDeal() {
	this.CardMgr.Shuffle()
	for _, v := range this.Players {
		v.HandCard = this.CardMgr.SendHandCard(5)
		v.NiuPoint, v.NiuCards = GetResult(v.HandCard)
		v.NiuMultiple = GetNiuMultiple(v.NiuPoint)
		handcard := GSSendHandCards{
			Id:       MSG_GAME_INFO_DEAL_REPLY,
			ChairId:  v.ChairId,
			NiuPoint: v.NiuPoint,
			NiuCards: v.HandCard,
		}

		v.SendNativeMsg(MSG_GAME_INFO_DEAL_REPLY, &handcard)
	}
}

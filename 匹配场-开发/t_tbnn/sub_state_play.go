package main

func (this *ExtDesk) GameStatePlay() {
	this.BroadStageTime(TIME_STAGE_PLAYCARD_NUM)
	//进入倒计时
	this.runTimer(TIME_STAGE_PLAYCARD_NUM, this.GameStatePlayEnd)
}

//玩牌阶段-结束
func (this *ExtDesk) GameStatePlayEnd(d interface{}) {
	// logs.Debug("玩牌")

	for _, v := range this.Players {
		if !v.IsLook {
			//到时间，未点击
			niuHand := GHandNiuReply{
				Id:       MSG_GAME_INFO_PLAY_REPLY,
				ChairId:  v.ChairId,
				NiuPoint: v.NiuPoint,
				NiuCards: v.NiuCards,
			}

			this.BroadcastAll(MSG_GAME_INFO_PLAY_REPLY, niuHand)
			v.IsLook = true
		}
	}

	//结算
	this.nextStage(STAGE_SETTLE)

}

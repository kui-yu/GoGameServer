package main

func (this *ExtDesk) GameStatePlay() {
	this.BroadStageTime(STAGE_PLAY_TIME)
	//进入倒计时
	this.runTimer(STAGE_PLAY_TIME, this.GameStatePlayEnd)
}

//阶段-玩牌
func (this *ExtDesk) GameStatePlayEnd(d interface{}) {
	for _, v := range this.Players {
		if v.IsPlay == 0 {
			v.IsPlay = 1
			v.PlayTypes, v.PlayCards = RecommendPoker(v.HandCards, NORMAL_FIVE_KIND)

			result := GSPlayInfo{
				Id:      MSG_GAME_INFO_PLAY_REPLY,
				ChairId: v.ChairId,
			}
			this.BroadcastAll(MSG_GAME_INFO_PLAY_REPLY, &result)
		}
	}
	//进入玩牌阶段
	this.nextStage(STAGE_SETTLE)
}

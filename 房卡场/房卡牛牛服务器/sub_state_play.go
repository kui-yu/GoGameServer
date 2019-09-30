package main

func (this *ExtDesk) GameStatePlay() {
	this.BroadStageTime(STAGE_PLAYCARD_TIME)
	//进入倒计时
	this.runTimer(STAGE_PLAYCARD_TIME, this.GameStatePlayEnd)
}

//阶段-玩牌结束
func (this *ExtDesk) GameStatePlayEnd(d interface{}) {
	// logs.Debug("玩牌")

	for _, v := range this.Players {
		if !v.IsLook {
			//到时间，未点击
			finish := GSPlayCard{
				Id:      MSG_GAME_INFO_PLAY_REPLY,
				ChairId: v.ChairId,
			}
			this.BroadcastAll(MSG_GAME_INFO_PLAY_REPLY, finish)
			v.IsLook = true
		}
	}

	//结算
	this.nextStage(STAGE_SETTLE)
}

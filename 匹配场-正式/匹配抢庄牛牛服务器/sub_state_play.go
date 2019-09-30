package main

func (this *ExtDesk) GameStatePlay() {

	this.BroadStageTime(TIME_STAGE_PLAYCARD_NUM)
	//进入倒计时
	this.runTimer(TIME_STAGE_PLAYCARD_NUM, this.GameStatePlayEnd)
}

//玩牌阶段-结束
func (this *ExtDesk) GameStatePlayEnd(d interface{}) {

	for _, v := range this.Players {
		if !v.IsLook {
			//到时间，未点击
			finish := GPlayCard{
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

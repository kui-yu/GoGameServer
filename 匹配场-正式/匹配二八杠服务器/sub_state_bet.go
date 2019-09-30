package main

func (this *ExtDesk) GameStateBet() {
	//下注阶段
	this.BroadStageTime(STAGE_PLAY_TIME)
	//进入倒计时
	this.runTimer(STAGE_PLAY_TIME, this.GameStateBetEnd)
}

//阶段-下注
func (this *ExtDesk) GameStateBetEnd(d interface{}) {
	//未下注
	for _, v := range this.Players {
		if v.PlayMultiple == -1 && v.ChairId != this.Banker {
			//默认1倍
			v.PlayMultiple = 1
			//下注返回
			info := GSPlayerPlayInfo{
				Id:           MSG_GAME_INFO_PLAY_INFO_REPLY,
				ChairId:      v.ChairId,
				PlayMultiple: v.PlayMultiple,
			}
			this.BroadcastAll(MSG_GAME_INFO_PLAY_INFO_REPLY, &info)
		}
	}
	//发牌游戏开始
	this.nextStage(STAGE_DEAL)
}

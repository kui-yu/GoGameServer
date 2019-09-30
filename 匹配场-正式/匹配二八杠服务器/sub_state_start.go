package main

func (this *ExtDesk) GameStateStart() {
	//开始阶段
	this.Round++
	//重置玩家信息
	for _, v := range this.Players {
		this.ResetPlayer(v)
	}
	this.BroadStageTime(STAGE_START_TIME)
	//返回开始与下注信息
	info := GSStartInfo{
		Id:    MSG_GAME_INFO_START_INFO_REPLY,
		Round: this.Round,
	}
	this.BroadcastAll(MSG_GAME_INFO_START_INFO_REPLY, &info)
	//进入倒计时
	this.runTimer(STAGE_START_TIME, this.GameStateStartEnd)
}

//阶段-开始
func (this *ExtDesk) GameStateStartEnd(d interface{}) {
	// 进入抢庄阶段
	this.nextStage(STAGE_CALL)
}

package main

func (this *ExtDesk) GameStateStart() {
	//开始阶段
	this.Round++
	this.BroadcastAll(MSG_GAME_INFO_START_INFO, &GSStartInfo{
		Id:    MSG_GAME_INFO_START_INFO,
		Round: this.Round,
	})
	//第一局扣房费
	if this.Round == 1 {
		minGoal := this.getPayMoney()
		this.GToHAddRoomCard(this.FkOwner, int64(-minGoal))
	}
	//游戏开始
	this.BroadStageTime(STAGE_START_TIME)
	this.runTimer(STAGE_START_TIME, this.GameStateStartEnd)
}

func (this *ExtDesk) getPayMoney() int64 {
	var roundValue int64
	if this.TableConfig.TotalRound <= 5 {
		roundValue = 1
	} else if this.TableConfig.TotalRound <= 10 {
		roundValue = 2
	} else if this.TableConfig.TotalRound <= 15 {
		roundValue = 3
	} else {
		roundValue = 4
	}
	var feeValue int64
	var money int64 = 1
	feeValue = roundValue * money
	return feeValue
}

//阶段-游戏开始
func (this *ExtDesk) GameStateStartEnd(d interface{}) {
	if this.TableConfig.GameType == 1 {
		//抢庄牛牛
		this.nextStage(STAGE_CALL_BANKER)
	} else {
		//通比牛牛
		this.nextStage(STAGE_CALL_SCORE)
	}
}

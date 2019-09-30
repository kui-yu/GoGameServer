package main

// import (
// 	"logs"
// )

func (this *ExtDesk) GameStateStart() {
	//开始阶段
	this.Round++
	this.BroadcastAll(MSG_GAME_INFO_START_INFO, &GSStartInfo{
		Id:    MSG_GAME_INFO_START_INFO,
		Round: this.Round,
	})
	//第一局扣房费
	if this.Round == 1 {
		minGoal := this.TableConfig.TotalRound / 5
		this.GToHAddRoomCard(this.FkOwner, int64(-minGoal))
		// feeValue := this.getPayMoney()
		// if this.TableConfig.PayType == 1 {
		// 	this.GToHAddCoin(this.FkOwner, -feeValue)
		// } else {
		// 	for _, v := range this.Players {
		// 		this.GToHAddCoin(v.Uid, -feeValue)
		// 	}
		// }
	}

	this.BroadStageTime(STAGE_START_TIME)
	//开始发牌-广播玩家
	this.DealPoker()
	//进入倒计时
	this.runTimer(STAGE_START_TIME, this.GameStateStartEnd)
}

func (this *ExtDesk) getPayMoney() int64 {
	var roundValue int
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
	if this.TableConfig.PayType == 1 {
		if this.TableConfig.GameModule == 2 {
			//房主支付
			feeValue = int64(5000 * roundValue * this.TableConfig.PlayerNumber)
		} else {
			//房主支付
			feeValue = int64(float64(10000*roundValue*this.TableConfig.PlayerNumber*88/100) + 0.5)
		}
	} else {
		if this.TableConfig.GameModule == 2 {
			//AA支付
			feeValue = int64(5000 * roundValue)
		} else {
			//AA支付
			feeValue = int64(10000 * roundValue)
		}

	}
	return feeValue
}

//阶段-开始
func (this *ExtDesk) GameStateStartEnd(d interface{}) {
	//进入玩牌阶段
	this.nextStage(STAGE_PLAY)
}

package main

func (this *ExtDesk) GameStateCall() {
	//抢庄阶段
	this.BroadStageTime(STAGE_CALL_TIME)
	if GetCostType() == 1 { //如果是体验场就不做筹码判断
		//抢庄筹码计算
		for _, v := range this.Players {
			var callList []int
			if v.Coins >= int64(22*this.Bscore*3) {
				callList = []int{0, 1, 2, 3}
			} else if v.Coins >= int64(22*this.Bscore*2) {
				callList = []int{0, 1, 2}
			} else if v.Coins >= int64(22*this.Bscore*1) {
				callList = []int{0, 1}
			} else {
				callList = []int{0}
			}

			v.PlayerCalls = callList
			var info GSCallList
			info.Id = MSG_GAME_INFO_CALL_LIST
			info.CallListCnt = len(callList)
			info.CallList = callList
			v.SendNativeMsg(MSG_GAME_INFO_CALL_LIST, &info)
		}
	} else {
		for _, v := range this.Players {
			v.PlayerCalls = []int{0, 1, 2, 3}
			v.SendNativeMsg(MSG_GAME_INFO_CALL_LIST, &GSCallList{
				Id:          MSG_GAME_INFO_CALL_LIST,
				CallListCnt: 4,
				CallList:    []int{0, 1, 2, 3},
			})
		}

	}

	//进入倒计时
	this.runTimer(STAGE_CALL_TIME, this.GameStateCallEnd)
}

//阶段-抢庄
func (this *ExtDesk) GameStateCallEnd(d interface{}) {

	for _, v := range this.Players {
		if v.CallMultiple == -1 { //如果等于负1说明玩家未操作，通知客户端玩家不抢
			//叫庄返回
			info := GSPlayerCallInfo{
				Id:           MSG_GAME_INFO_CALL_INFO_REPLY,
				ChairId:      v.ChairId,
				CallMultiple: 0,
			}
			this.BroadcastAll(MSG_GAME_INFO_CALL_INFO_REPLY, &info)
		}
	}

	this.ChooseBanker()

	//抢庄时间到,进入下注阶段
	this.nextStage(STAGE_BET)
}

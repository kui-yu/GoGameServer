package main

func (this *ExtDesk) GameStateCall() {

	this.BroadStageTime(STAGE_CALL_BANKER_TIME)

	callList := []int{0, 50, 100, 200}
	for _, v := range this.Players {
		var playerCalls []int
		for _, call := range callList {
			if v.Coins > int64(this.TableConfig.BaseScore*call) && this.TableConfig.GameModule == 2 {
				playerCalls = append(playerCalls, call)
			} else {
				playerCalls = append(playerCalls, call)
			}
		}
		v.SendNativeMsg(MSG_GAME_INFO_CALL_LIST, &GSCallList{
			Id:          MSG_GAME_INFO_CALL_LIST,
			CallListCnt: len(playerCalls),
			CallList:    playerCalls,
		})
		v.PlayerCalls = playerCalls
	}

	//进入倒计时
	this.runTimer(STAGE_CALL_BANKER_TIME, this.GameStateCallEnd)
}

//阶段-抢庄结束
func (this *ExtDesk) GameStateCallEnd(d interface{}) {

	//判断是否都没有抢庄
	for _, v := range this.Players {

		//没有抢庄
		if v.CallMultiple == -1 {
			v.CallMultiple = 0
			re := GSCallMsg{
				Id:       MSG_GAME_INFO_CALL_BANKER_NOTIFY,
				ChairId:  v.ChairId,
				Multiple: v.CallMultiple,
			}
			//广播
			this.BroadcastAll(MSG_GAME_INFO_CALL_BANKER_NOTIFY, &re)
		}
	}

	//都没人叫庄，广播随机庄家
	this.ChooseBanker()
	//跳到下一个阶段
	this.nextStage(STAGE_CALL_SCORE)
}

//指定庄家
func (this *ExtDesk) ChooseBanker() {

	var count int
	for _, v := range this.Players {
		if v.CallMultiple == 0 {
			count++
		}
	}

	var callMax int
	//取最大倍数玩家
	for i := 1; i < len(this.Players); i++ {
		if this.Players[callMax].CallMultiple < this.Players[i].CallMultiple {
			callMax = i
		}
	}

	callMaxList := []int{callMax}
	//判断有没有一样大的玩家
	for i := callMax + 1; i < len(this.Players); i++ {
		if this.Players[callMax].CallMultiple == this.Players[i].CallMultiple {
			//倍数一样
			callMaxList = append(callMaxList, i)
		}
	}

	//选出庄家
	if len(callMaxList) > 1 {
		// 2.叫分一样，随机当庄
		//随机
		tempList := ListShuffle(callMaxList)
		//取第一个
		this.Banker = this.Players[tempList[0]].ChairId
	} else {
		this.Banker = this.Players[callMax].ChairId
	}

	//没人叫庄，默认
	if count >= len(this.Players) {
		this.Players[this.Banker].CallMultiple = 50
		// logs.Debug("庄家倍数", this.Players[this.Banker].CallBank)
	}
	// logs.Debug("庄家位置", this.Banker)
	//计算倍数

	banker := GSCallBank{
		Id:              MSG_GAME_INFO_CHOOSE_BANKER,
		Banker:          this.Banker,
		BankerList:      callMaxList,
		BankerMultiples: this.Players[this.Banker].CallMultiple,
	}
	//广播
	this.BroadcastAll(MSG_GAME_INFO_CHOOSE_BANKER, &banker)
}

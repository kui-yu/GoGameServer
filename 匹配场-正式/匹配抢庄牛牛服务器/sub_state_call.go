package main

func (this *ExtDesk) GameStateCall() {
	//跳入抢庄阶段
	this.BroadStageTime(TIME_STAGE_CALL_NUM)

	callList := []int{0, 50, 100, 200}
	for _, v := range this.Players {
		var playerCalls []int
		if GetCostType() == 1 { //如果不是体验场再判断持有金币的下注能力
			for _, call := range callList {
				if v.Coins > int64(this.Bscore*call) {
					playerCalls = append(playerCalls, call)
				}
			}
		} else { //体验场则不判断
			playerCalls = callList
		}
		v.SendNativeMsg(MSG_GAME_INFO_CALL_LIST, &GSCallList{
			Id:          MSG_GAME_INFO_CALL_LIST,
			CallListCnt: len(playerCalls),
			CallList:    playerCalls,
		})
		v.PlayerCalls = playerCalls
	}

	//进入倒计时
	this.runTimer(TIME_STAGE_CALL_NUM, this.GameStateCallEnd)
}

//抢庄时间到
func (this *ExtDesk) GameStateCallEnd(d interface{}) {

	//判断是否都没有抢庄
	for _, v := range this.Players {

		//没有抢庄
		if !v.CallBankFlag {
			re := GCallMsgReply{
				Id:       MSG_GAME_INFO_CALL_BANKER_NOTIFY,
				ChairId:  v.ChairId,
				Multiple: v.CallMultiple,
			}
			//广播
			this.BroadcastAll(MSG_GAME_INFO_CALL_BANKER_NOTIFY, &re)
		}
	}

	//选庄
	this.ChooseBanker()
	//下个阶段
	this.nextStage(STAGE_BET)
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
		//转化
		tempList := ListIntToInt32(callMaxList)
		//随机
		tempList = ListShuffle(tempList)
		//取第一个
		this.Banker = this.Players[tempList[0]].ChairId
	} else {
		this.Banker = this.Players[callMax].ChairId
	}

	//没人叫庄，默认
	if count >= len(this.Players) {
		this.Players[this.Banker].CallMultiple = 50
	}

	//判断玩家可下倍数
	var betList []int
	if this.Players[this.Banker].CallMultiple == 200 {
		betList = []int{1, 3, 6, 10}
	} else if this.Players[this.Banker].CallMultiple == 100 {
		betList = []int{1, 2, 3, 5}
	} else {
		betList = []int{1, 2, 3}
	}

	for _, v := range this.Players {
		var playerBet []int
		if GetCostType() == 1 { //如果不是体验场再玩家进行下注倍数能力判断
			for _, bet := range betList {
				if v.Coins >= int64(bet*this.Bscore*5) {
					playerBet = append(playerBet, bet)
				}
			}
		} else { //如果是体验场则不做限制
			playerBet = betList
		}

		v.PlayerBets = playerBet
		banker := GCallBankReply{
			Id:              MSG_GAME_INFO_CHOOSE_BANKER,
			Banker:          this.Banker,
			BankerList:      callMaxList,
			BankerMultiples: this.Players[this.Banker].CallMultiple,
			BetListCnt:      len(v.PlayerBets),
			BetList:         v.PlayerBets,
		}
		//广播
		v.SendNativeMsg(MSG_GAME_INFO_CHOOSE_BANKER, &banker)
	}
}

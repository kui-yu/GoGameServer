package main

import (
	"encoding/json"
)

//玩家叫庄
func (this *ExtDesk) HandleGameCall(p *ExtPlayer, d *DkInMsg) {
	//不是叫庄阶段
	if this.GameState != STAGE_CALL {
		return
	}

	//已经叫庄过
	if p.CallMultiple > -1 {
		return
	}
	//解析指令
	data := GAPlayerCallInfo{}
	err := json.Unmarshal([]byte(d.Data), &data)
	if err != nil {
		return
	}

	callFlag := true
	//抢庄限制 {不抢,1,2,3}
	for _, call := range p.PlayerCalls {
		if call == data.CallMultiple {
			callFlag = false
			break
		}
	}
	//没有抢庄条件
	if callFlag {
		return
	}

	p.CallMultiple = data.CallMultiple
	//叫庄返回
	info := GSPlayerCallInfo{
		Id:           MSG_GAME_INFO_CALL_INFO_REPLY,
		ChairId:      p.ChairId,
		CallMultiple: p.CallMultiple,
	}
	this.BroadcastAll(MSG_GAME_INFO_CALL_INFO_REPLY, &info)

	flag := true
	for _, v := range this.Players {
		if v.CallMultiple == -1 {
			flag = false
		}
	}

	if flag {
		//指定庄家
		this.ChooseBanker()
		//进入下一个阶段
		this.nextStage(STAGE_BET)
	}

}

//指定庄家
func (this *ExtDesk) ChooseBanker() {

	var noCall []int
	for _, v := range this.Players {
		if v.CallMultiple <= 0 {
			//叫庄返回
			ListAdd(&noCall, int(v.ChairId))
		}
	}

	//都不抢庄
	if len(noCall) == len(this.Players) {
		for _, v := range this.Players {
			if v.CallMultiple <= 0 {
				v.CallMultiple = 1
			}
		}
	}

	//取最大倍数玩家
	var callMax int
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
		//随机一个庄家
		tempList := ListShuffle(callMaxList)
		//取第一个
		this.Banker = this.Players[tempList[0]].ChairId
	} else {
		this.Banker = this.Players[callMax].ChairId
	}

	//谁是庄家
	bankerInfo := GSPlayerCallBank{
		Id:              MSG_GAME_INFO_BANKER_REPLY,
		Banker:          this.Banker,
		BankerList:      callMaxList,
		BankerMultiples: this.Players[this.Banker].CallMultiple,
	}

	var betList = []int{1, 6, 12, 18, 22}
	//如果不是体验场则判断玩家下注范围
	if GetCostType() == 1 {
		for _, v := range this.Players {
			var playerBet []int
			for _, bet := range betList {
				maxBet := int64(bet * this.Bscore * this.Players[this.Banker].CallMultiple)
				if v.Coins >= maxBet && this.Players[this.Banker].Coins >= maxBet*3 {
					playerBet = append(playerBet, bet)
				}
			}
			if len(playerBet) == 5 {
				playerBet = []int{1, 6, 18, 22}
			}
			v.PlayerBets = playerBet
			bankerInfo.BetListCnt = len(playerBet)
			bankerInfo.BetList = playerBet
			v.SendNativeMsg(MSG_GAME_INFO_BANKER_REPLY, &bankerInfo)
		}
	} else {
		for _, v := range this.Players {
			v.PlayerBets = []int{1, 6, 18, 22}
			bankerInfo.BetListCnt = 4
			bankerInfo.BetList = []int{1, 6, 18, 22}
			v.SendNativeMsg(MSG_GAME_INFO_BANKER_REPLY, &bankerInfo)
		}
	}
}

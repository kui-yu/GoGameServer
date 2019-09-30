package main

import (
	"logs"
	"math/rand"
)

type GControlValues struct {
	num      int          //步数
	handcard []byte       //手牌
	rsValues []GTypeValue //结果
}

//游戏开始，进入发牌阶段
func (this *ExtDesk) TimerSendCard(d interface{}) {

	this.CurCid = int32(rand.Intn(10) % len(this.Players)) //随机出一个id 作为叫分人
	//洗牌
	this.CardMgr.Shuffle()

	// 发牌 17 张
	var pokerValues []GControlValues
	for i := 0; i < len(this.Players); i++ {
		hd := this.CardMgr.SendHandCard(17)
		Sort(hd)
		hdValues := R_GetBestCalc(hd)
		pokerValues = append(pokerValues, GControlValues{
			num:      hdValues.num,
			handcard: hd,
			rsValues: hdValues.rsValues,
		})
	}
	pokerValues = this.PokerControl(pokerValues)
	winPlayers := this.ControlResult()

	for i := 0; i < len(winPlayers); i++ {
		v := this.Players[winPlayers[i]]
		v.SetHandCard(pokerValues[i].handcard)
		logs.Debug("座位：", v.ChairId, "手牌", v.HandCard)
		logs.Debug("小牌", R_GetValues(v.HandCard))
	}

	// 底牌
	this.DiPai = []byte{}
	this.DiPai = this.CardMgr.SendHandCard(3)

	// this.Players[0].SetHandCard(Sort(StrSplitToCards("66 2 1 61 44 12 11 58 55 39 22 37 5 20 4 19 3")))
	// this.Players[1].SetHandCard(Sort(StrSplitToCards("34 45 13 28 59 27 57 41 25 9 56 24 38 6 53 21 35")))
	// this.Players[2].SetHandCard(Sort(StrSplitToCards("65 18 49 33 17 29 60 43 26 10 8 23 7 54 52 36 51")))

	// this.DiPai = StrSplitToCards("42 40 50")

	logs.Debug("底牌", this.DiPai)
	logs.Debug("小牌", R_GetValues(this.DiPai))
	//发送牌消息
	sd := GGameSendCardNotify{
		Id:  MSG_GAME_INFO_SEND_NOTIFY,
		Cid: this.CurCid,
	}
	for _, v := range this.Players {
		sd.HandsCards = v.HandCard
		v.SendNativeMsg(MSG_GAME_INFO_SEND_NOTIFY, &sd)
	}
	//发牌后进入叫分阶段，开启叫分阶段的定时器
	this.AddTimer(TIMER_START, TIMER_START_NUM, this.TimerDealPoker, nil)
}

//等待发牌动画结束，进入叫分抢庄阶段
func (this *ExtDesk) TimerDealPoker(d interface{}) {
	//进入抢庄
	this.GameState = GAME_STATUS_CALL
	this.BroadStageTime(TIMER_CALL_NUM)
	//判断叫分是否是机器人
	nextplayer := this.Players[this.CurCid]
	if nextplayer.Robot {
		this.TimerCall(nextplayer)
	} else {
		//发牌后进入叫分阶段，开启叫分阶段的定时器
		this.AddTimer(TIMER_CALL, TIMER_CALL_NUM, this.TimerCallEnd, nil)
	}
}

//排序手牌
func (this *ExtDesk) PokerControl(pokerValues []GControlValues) []GControlValues {

	for i := 0; i < len(pokerValues)-1; i++ {
		minValues := pokerValues[i]
		for j := i + 1; j < len(pokerValues); j++ {
			if pokerValues[i].num > pokerValues[j].num {
				minValues = pokerValues[j]
				pokerValues[j] = pokerValues[i]
				pokerValues[i] = minValues
			}
		}
	}
	return pokerValues
}

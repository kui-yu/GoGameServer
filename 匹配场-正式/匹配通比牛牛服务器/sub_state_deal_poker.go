package main

import (
	"logs"
)

func (this *ExtDesk) GameStateDeal() {

	this.CardMgr.Shuffle()

	var cardTypes []GCardType //手牌数组
	for i := 0; i < len(this.Players); i++ {
		handCard := this.CardMgr.SendHandCard(5)
		niuPoint, niuCards := GetResult(handCard)
		cardType := GCardType{
			HandCard: handCard,
			NiuPoint: niuPoint,
			NiuCards: niuCards,
		}
		cardTypes = append(cardTypes, cardType)
	}
	//排序
	cardTypes = this.ControlPoker(cardTypes)
	//生成赢家座位
	winChairs := this.ControlResult()
	logs.Debug("winChair:::", winChairs)
	for i := 0; i < len(winChairs); i++ {
		this.Players[winChairs[i]].HandCard = cardTypes[i].HandCard
		this.Players[winChairs[i]].NiuPoint = cardTypes[i].NiuPoint
		this.Players[winChairs[i]].NiuCards = cardTypes[i].NiuCards
		this.Players[winChairs[i]].NiuMultiple = GetNiuMultiple(cardTypes[i].NiuPoint)

		handcard := GHandNiuReply{
			Id:       MSG_GAME_INFO_DEAL_REPLY,
			ChairId:  this.Players[winChairs[i]].ChairId,
			NiuPoint: this.Players[winChairs[i]].NiuPoint,
			NiuCards: this.Players[winChairs[i]].HandCard,
		}
		this.Players[winChairs[i]].SendNativeMsg(MSG_GAME_INFO_DEAL_REPLY, &handcard)
	}

	//进入发牌动画
	this.BroadStageTime(TIME_STAGE_START_NUM)
	//进入倒计时
	this.runTimer(TIME_STAGE_START_NUM, this.GameStartEnd)
}

//发牌动画-结束
func (this *ExtDesk) GameStartEnd(d interface{}) {
	// logs.Debug("发牌动画")
	//进入玩牌
	this.nextStage(STAGE_PLAY)
}

//获取机器人
func (this *ExtDesk) GetRobot() []*ExtPlayer {
	result := []*ExtPlayer{}
	for _, v := range this.Players {
		if v.Robot {
			result = append(result, v)
		}
	}
	return result
}

//返回牌组，排序的牌组（从大到小）
func (this *ExtDesk) ControlPoker(cardTypes []GCardType) []GCardType {
	// logs.Debug("排序前", cardTypes)
	for i := 0; i < len(cardTypes); i++ {
		maxCardTypes := cardTypes[i]
		for j := i + 1; j < len(cardTypes); j++ {
			rs := SoloResult(cardTypes[i], cardTypes[j])
			if rs == 2 {
				maxCardTypes = cardTypes[j]
				cardTypes[j] = cardTypes[i]
				cardTypes[i] = maxCardTypes
			}
		}
	}
	// logs.Debug("排序后", cardTypes)
	return cardTypes
}

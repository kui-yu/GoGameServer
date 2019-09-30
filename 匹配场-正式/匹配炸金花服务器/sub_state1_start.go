package main

import (
	"logs"
)

//开始阶段
func (this *ExtDesk) GameStateStart() {

	this.BroadStageTime(STAGE_START_TIME)
	//重置玩家信息
	for _, v := range this.Players {
		this.ResetPlayer(v)
	}

	//最小下注额
	this.MinCoin = 1

	//发牌
	//按人数发牌 人数*3张
	handcard := this.ControlPoker()
	winChairs := this.ControlResult()
	for i := 0; i < len(winChairs); i++ {
		v := this.Players[winChairs[i]]
		if v.Robot {
			this.RobotWinner = true
		}

		cardNum := i + 1
		v.OldHandCard = Sort(handcard[cardNum])
		card, color := SortHandCard(handcard[cardNum])
		v.HandCards = Sort(card)
		v.HandColor = color
		// logs.Debug("排序：", v.HandCards)
		//获取卡牌类型
		lv, _ := GetCardType(v.HandCards, v.HandColor)
		v.CardLv = lv
		//桌子下注筹码列表
		this.CoinList = append(this.CoinList, 1)
		//玩家下注筹码列表
		v.PayCoin = append(v.PayCoin, 1)
	}
	//金币下注
	this.CoinPush()
	this.runTimer(STAGE_START_TIME, this.HandlePayCardCoin)
}

//阶段-发牌和底注
func (this *ExtDesk) HandlePayCardCoin(d interface{}) {
	// logs.Debug("阶段-发牌和底注")
	this.nextStage(STAGE_PLAY_OPERATION)
}

func (this *ExtDesk) ControlPoker() [][]int {

	var count int
	var handCards [][]int
	for {
		count++
		//洗牌
		this.CardMgr.Shuffle()
		//按人数发牌 人数*3张
		handCards = this.CardMgr.HandCardInfo(30)

		logs.Debug("排序前", handCards)
		for i := 0; i < len(handCards)-1; i++ {
			maxCard := handCards[i]
			for j := i + 1; j < len(handCards); j++ {
				maxCard, maxColor := SortHandCard(maxCard)
				card, color := SortHandCard(handCards[j])
				if GetResult(maxCard, maxColor, card, color) == 1 {
					maxCard = handCards[j]
					handCards[j] = handCards[i]
					handCards[i] = maxCard
				}
			}
		}
		logs.Debug("排序后", handCards)
		threeCard, threeColor := SortHandCard(handCards[3])
		types, _ := GetCardType(threeCard, threeColor)
		if types > CARD_SINGLE || count > 500 {
			break
		}
	}

	return handCards
}

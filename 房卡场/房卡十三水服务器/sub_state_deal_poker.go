package main

// import (
// 	// "logs"
// 	"math/rand"
// 	"time"
// )

//发牌
func (this *ExtDesk) DealPoker() {
	//洗牌
	this.CardMgr.Shuffle()

	var handCards [][]int
	for i := 0; i < len(this.Players); i++ {
		//发牌
		handCard := this.CardMgr.SendHandCard(13)
		handCards = append(handCards, handCard)
	}
	handCards = this.ControlPoker(handCards)
	winChairs := this.ControlResult()
	//发送手牌消息
	for i := 0; i < len(winChairs); i++ {
		chairId := winChairs[i]
		this.Players[chairId].HandCards = handCards[i]
		//计算特殊牌
		specialType, specialCards := MathSpecial(this.Players[chairId].HandCards)
		if specialType > 0 {
			this.Players[chairId].SpecialType = specialType
			this.Players[chairId].SpecialCards = specialCards
		}

		result := GSHandInfo{
			Id:           MSG_GAME_INFO_HANDINFO_REPLY,
			ChairId:      this.Players[chairId].ChairId,
			HandCards:    this.Players[chairId].HandCards,
			SpecialType:  this.Players[chairId].SpecialType,
			SpecialCards: this.Players[chairId].SpecialCards,
		}
		this.Players[chairId].SendNativeMsg(MSG_GAME_INFO_HANDINFO_REPLY, &result)
	}
}

//手牌大小
func (this *ExtDesk) ControlPoker(handCards [][]int) [][]int {
	// logs.Debug("排序前", handCards)

	var specialCards [][]int
	var normalCards [][]int

	for i := 0; i < len(handCards); i++ {
		specialType, _ := MathSpecial(handCards[i])
		if specialType > 0 {
			specialCards = append(specialCards, handCards[i])
		} else {
			normalCards = append(normalCards, handCards[i])
		}
	}

	var plays []GRecommendPoker
	for i := 0; i < len(normalCards); i++ {
		types, cards := RecommendPoker(normalCards[i], NORMAL_FIVE_KIND)
		info := GRecommendPoker{
			Types: types,
			Cards: cards,
		}
		plays = append(plays, info)
	}
	for i := 0; i < len(normalCards)-1; i++ {
		for j := i + 1; j < len(normalCards); j++ {
			maxPlays := normalCards[i]
			var compareCount int
			//比较三墩
			for count := 0; count < len(plays[i].Types); count++ {
				var ct1, ct2 GCardsType
				if count == 0 {
					//头墩
					ct1 = GCardsType{
						Type:  plays[i].Types[0],
						Cards: ListGet(plays[i].Cards, 0, 3),
					}
					ct2 = GCardsType{
						Type:  plays[j].Types[0],
						Cards: ListGet(plays[j].Cards, 0, 3),
					}
				} else if count == 1 {
					//中墩
					ct1 = GCardsType{
						Type:  plays[i].Types[1],
						Cards: ListGet(plays[i].Cards, 3, 5),
					}
					ct2 = GCardsType{
						Type:  plays[j].Types[1],
						Cards: ListGet(plays[j].Cards, 3, 5),
					}
				} else if count == 2 {
					//底墩
					ct1 = GCardsType{
						Type:  plays[i].Types[2],
						Cards: ListGet(plays[i].Cards, 8, 5),
					}
					ct2 = GCardsType{
						Type:  plays[j].Types[2],
						Cards: ListGet(plays[j].Cards, 8, 5),
					}
				}

				rs := MCardsCompare(ct1, ct2)
				if rs == 2 {
					compareCount += getNormalPoint(ct2.Type, count)
				}
			}
			if compareCount >= 2 {
				maxPlays = normalCards[j]
				normalCards[j] = normalCards[i]
				normalCards[i] = maxPlays
			}
		}
	}

	var rsCards [][]int
	if len(specialCards) > 0 {
		rsCards = append(rsCards, specialCards...)
	}
	rsCards = append(rsCards, normalCards...)

	// logs.Debug("排序后", rsCards)
	return rsCards
}

func (this *ExtDesk) ControlResult() []int32 {

	//玩家随机
	var winChairs []int32

	for _, v := range this.Players {
		winChairs = append(winChairs, v.ChairId)
	}
	winChairs = ListShuffleByInt32(winChairs)

	// logs.Debug("玩家排序", winChairs)
	return winChairs
}

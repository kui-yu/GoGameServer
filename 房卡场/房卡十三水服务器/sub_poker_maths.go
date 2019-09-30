package main

// import (
// 	"logs"
// )

//判断特殊牌型
func MathSpecial(cards []int) (int, []int) {
	//青龙 108分
	if CheckColor(cards) && CheckDragon(cards) {
		return SPECIAL_COLOR_DRAGON, ListSortAsc(cards)
	}
	//一条龙 26分
	if CheckDragon(cards) {
		return SPECIAL_DRAGON, ListSortAsc(cards)
	}
	//三同花顺 18分
	threeCount, threeLists := GetStraightLists(cards)
	if threeCount >= 2 {
		var specialList []int
		var threePoker []int
		var twoPoker []int
		if len(threeLists) >= 2 {
			threePoker = ListDelGet(threeLists[0], 0, 5)
			twoPoker = ListDelGet(threeLists[1], 0, 5)
		} else {
			threePoker = ListDelGet(threeLists[0], 0, 5)
			twoPoker = ListDelGet(threeLists[0], 5, 5)
		}
		//判断余数是否是顺子
		var tempList []int
		tempList = append(tempList, threePoker...)
		tempList = append(tempList, twoPoker...)
		lastPoker := ListDelList(cards, tempList)
		if CheckStraight(ListSortDesc(lastPoker)) {
			if CheckColor(threePoker) &&
				CheckColor(twoPoker) &&
				CheckColor(lastPoker) {
				specialList = append(specialList, lastPoker...)
				specialList = append(specialList, twoPoker...)
				specialList = append(specialList, threePoker...)
				return SPECIAL_THREE_SAME_COLOR_STRAIGHT, specialList
			}
		}

	}
	//三分天下 16分
	fourKindNum, _ := GetKindLists(cards, 4)
	if fourKindNum >= 3 {
		return SPECIAL_THREE_FOUR_KIND, ListSortAsc(cards)
	}
	//四套三条 8分
	threeKindNum, _ := GetKindLists(cards, 3)
	if threeKindNum >= 4 {
		return SPECIAL_FOUR_THREE_KIND, ListSortAsc(cards)
	}
	//六队半 4分
	pairNum, _ := GetKindLists(cards, 2)
	if pairNum >= 6 {
		return SPECIAL_SIX_PAIR, ListSortAsc(cards)
	}

	//三顺子 4分
	if threeCount >= 2 {
		var specialList []int
		var threePoker []int
		var twoPoker []int
		if len(threeLists) >= 2 {
			threePoker = ListDelGet(threeLists[0], 0, 5)
			twoPoker = ListDelGet(threeLists[1], 0, 5)
		} else {
			threePoker = ListDelGet(threeLists[0], 0, 5)
			twoPoker = ListDelGet(threeLists[0], 5, 5)
		}
		//判断余数是否是顺子
		var tempList []int
		tempList = append(tempList, threePoker...)
		tempList = append(tempList, twoPoker...)
		lastPoker := ListDelList(cards, tempList)
		if CheckStraight(ListSortDesc(lastPoker)) {
			specialList = append(specialList, lastPoker...)
			specialList = append(specialList, twoPoker...)
			specialList = append(specialList, threePoker...)
			return SPECIAL_THREE_STRAIGHT, specialList
		}

	}
	//三同花 4分
	threeCount, threeLists = GetColorLists(cards)
	if threeCount >= 2 {
		var specialList []int
		var threePoker []int
		var twoPoker []int
		if len(threeLists) >= 2 {
			threePoker = ListDelGet(threeLists[0], 0, 5)
			twoPoker = ListDelGet(threeLists[1], 0, 5)
		} else {
			threePoker = ListDelGet(threeLists[0], 0, 5)
			twoPoker = ListDelGet(threeLists[0], 5, 5)
		}
		//判断余数是否是同花
		var tempList []int
		tempList = append(tempList, threePoker...)
		tempList = append(tempList, twoPoker...)
		lastPoker := ListDelList(cards, tempList)

		if CheckColor(ListSortDesc(lastPoker)) {
			specialList = append(specialList, lastPoker...)
			specialList = append(specialList, twoPoker...)
			specialList = append(specialList, threePoker...)
			return SPECIAL_THREE_SAME_COLOR, specialList
		}
	}

	return 0, []int{}
}

//获取同花
func GetColorLists(cards []int) (int, [][]int) {
	var fangList []int
	var meiList []int
	var hongList []int
	var heiList []int
	for _, v := range cards {
		if GetCardColor(v) == CARD_COLOR_Fang {
			ListAdd(&fangList, v)
		} else if GetCardColor(v) == CARD_COLOR_Mei {
			ListAdd(&meiList, v)
		} else if GetCardColor(v) == CARD_COLOR_Hong {
			ListAdd(&hongList, v)
		} else if GetCardColor(v) == CARD_COLOR_Hei {
			ListAdd(&heiList, v)
		}
	}

	var colorLists [][]int
	var colorCount int

	if len(heiList) >= 5 {
		colorLists = append(colorLists, heiList)
		colorCount++
		if len(heiList) >= 10 {
			colorCount++
		}
	}
	if len(hongList) >= 5 {
		colorLists = append(colorLists, hongList)
		colorCount++
		if len(hongList) >= 10 {
			colorCount++
		}
	}
	if len(meiList) >= 5 {
		colorLists = append(colorLists, meiList)
		colorCount++
		if len(meiList) >= 10 {
			colorCount++
		}
	}
	if len(fangList) >= 5 {
		colorLists = append(colorLists, fangList)
		colorCount++
		if len(fangList) >= 10 {
			colorCount++
		}
	}

	return colorCount, colorLists
}

//获取顺子
func GetStraightLists(cards []int) (int, [][]int) {

	cards = ListSortDesc(cards)

	var straightLists [][]int
	var straightCount int

	var straightList []int
	for i := 0; i < len(cards)-1; i++ {
		if GetCardValue(cards[i]) == GetCardValue(cards[i+1]) {
			continue
		}
		if GetCardValue(cards[i]) == GetCardValue(cards[i+1])+1 {
			if len(straightList) == 0 {
				ListAdd(&straightList, cards[i])
			}
			ListAdd(&straightList, cards[i+1])
		} else {
			if len(straightList) >= 5 {
				straightCount++
				if len(straightList) >= 10 {
					straightCount++
				}
				straightLists = append(straightLists, straightList)
			} else if len(straightList) >= 4 {
				//或者5，4，3，2
				rs := false
				tempCards2 := []int{2, 3, 4, 5}
				count := 0
				for i := 0; i < len(tempCards2); i++ {
					for j := 0; j < len(straightList); j++ {
						if tempCards2[i] == GetCardValue(straightList[j]) {
							count++
							break
						}
					}
				}
				if count == 4 {
					rs = true
				}
				//A牌
				aList := []int{}
				for a := 0; a < len(cards); a++ {
					if GetCardValue(cards[a]) == 14 {
						aList = append(aList, cards[a])
					}
				}

				aUseList := []int{}
				for as := 0; as < len(straightLists); as++ {
					for a := 0; a < len(straightLists[as]); a++ {
						if GetCardValue(straightLists[as][a]) == 14 {
							aUseList = append(aUseList, straightLists[as][a])
						}
					}
				}
				tempCards := ListDelList(aList, aUseList)
				if rs && len(tempCards) > 0 {
					straightList = append(straightList, tempCards[0])
					straightLists = append(straightLists, ListSortDesc(straightList))
					straightCount++
				}
			}
			straightList = []int{}
		}
	}

	if len(straightList) >= 5 {
		straightCount++
		if len(straightList) >= 10 {
			straightCount++
		}
		straightLists = append(straightLists, straightList)
		straightList = []int{}
	}
	//或者5，4，3，2
	if len(straightList) >= 4 {
		rs := false
		tempCards2 := []int{2, 3, 4, 5}
		count := 0
		for i := 0; i < len(tempCards2); i++ {
			for j := 0; j < len(straightList); j++ {
				if tempCards2[i] == GetCardValue(straightList[j]) {
					count++
					break
				}
			}
		}
		if count == 4 {
			rs = true
		}
		//A牌
		aList := []int{}
		for a := 0; a < len(cards); a++ {
			if GetCardValue(cards[a]) == 14 {
				aList = append(aList, cards[a])
			}
		}

		aUseList := []int{}
		for as := 0; as < len(straightLists); as++ {
			for a := 0; a < len(straightLists[as]); a++ {
				if GetCardValue(straightLists[as][a]) == 14 {
					aUseList = append(aUseList, straightLists[as][a])
				}
			}
		}
		tempCards := ListDelList(aList, aUseList)
		if rs && len(tempCards) > 0 {
			straightList = append(straightList, tempCards[0])
			straightLists = append(straightLists, ListSortDesc(straightList))
			straightCount++
		}
		straightList = []int{}
	}

	return straightCount, straightLists
}

//获取条子
func GetKindLists(cards []int, kindNum int) (int, [][]int) {
	var kindCount int
	var kindLists [][]int

	var kindList []int
	for i := 0; i < len(cards)-1; i++ {
		if GetCardValue(cards[i]) == GetCardValue(cards[i+1]) {
			if len(kindList) == 0 {
				ListAdd(&kindList, cards[i])
			}
			ListAdd(&kindList, cards[i+1])
		} else {
			if len(kindList) >= kindNum {
				kindCount++
				if len(kindList) >= kindNum*2 {
					kindCount++
				}
				kindLists = append(kindLists, kindList)
			}
			kindList = []int{}
		}
	}

	if len(kindList) >= kindNum {
		kindCount++
		if len(kindList) >= kindNum*2 {
			kindCount++
		}
		kindLists = append(kindLists, kindList)
	}

	return kindCount, kindLists
}

//获取扑克类型
func GetPokerType(cards []int, maxType int) (int, []int) {
	//五同
	if maxType >= NORMAL_FIVE_KIND {
		fiveKindCount, fiveKindLists := GetKindLists(cards, 5)
		if fiveKindCount > 0 {
			return NORMAL_FIVE_KIND, fiveKindLists[0]
		}
	}
	//同花顺
	if maxType >= NORMAL_COLOR_STRAIGHT {
		colorCount, colorLists := GetColorLists(cards)
		if colorCount > 0 {
			//判断是否顺子
			for i := 0; i < len(colorLists); i++ {
				straightCount, straightLists := GetStraightLists(colorLists[i])
				if straightCount > 0 {
					//判断是否同花
					for i := 0; i < len(straightLists); i++ {
						for j := 0; j < len(straightLists[i])-4; j++ {
							straightList := ListGet(straightLists[i], j, 5)
							return NORMAL_COLOR_STRAIGHT, straightList
						}
					}
				}
			}
		}
	}
	//铁支
	if maxType >= NORMAL_FOUR_KIND {
		fourKindCount, fourKindLists := GetKindLists(cards, 4)
		if fourKindCount > 0 {
			return NORMAL_FOUR_KIND, fourKindLists[0]
		}
	}

	//葫芦
	if maxType >= NORMAL_GOURD {
		threeCount, threeKindLists := GetKindLists(cards, 3)
		if threeCount > 0 {
			gourdPoker := threeKindLists[0]

			//先移除三条
			otherCards := append([]int{}, cards...)
			for i := 0; i < threeCount; i++ {
				otherCards = ListDelList(otherCards, threeKindLists[i])
			}
			otherCards = ListSortAsc(otherCards)
			pairCount, pairLists := GetKindLists(otherCards, 2)
			if pairCount > 0 {
				gourdPoker = append(gourdPoker, pairLists[0]...)
				return NORMAL_GOURD, gourdPoker
			} else {
				otherCards = ListDelList(cards, gourdPoker)
				otherCards = ListSortAsc(otherCards)
				pairCount, pairLists = GetKindLists(otherCards, 2)
				if pairCount > 0 {
					gourdPoker = append(gourdPoker, ListGet(pairLists[0], 0, 2)...)
					return NORMAL_GOURD, gourdPoker
				}
			}
		}
	}
	//同花
	if maxType >= NORMAL_SAME_COLOR {
		colorCount, colorLists := GetColorLists(cards)
		if colorCount > 0 {
			for i := 0; i < len(colorLists); i++ {
				if len(colorLists[i]) > 5 {
					pairCount, pairLists := GetKindLists(colorLists[i], 2)
					if pairCount > 0 {
						var colorList []int
						for j := 0; j < len(pairLists); j++ {
							colorList = append(colorList, pairLists[j]...)
						}
						lastPokerNum := 5 - len(colorList)
						if lastPokerNum <= 0 {
							return NORMAL_SAME_COLOR, ListGet(colorList, 0, 5)
						}
						//余下同花
						lastPoker := append([]int{}, ListDelList(colorLists[i], colorList)...)
						//对子同花排前
						colorList = append(colorList, lastPoker...)
						return NORMAL_SAME_COLOR, ListGet(colorList, 0, 5)
					}
				}
			}
			//可以递归求算最优牌型
			return NORMAL_SAME_COLOR, ListGet(colorLists[0], 0, 5)
		}
	}
	//顺子
	if maxType >= NORMAL_STRAIGHT {
		straightCount, straightLists := GetStraightLists(cards)
		if straightCount > 0 {
			//可以递归求算最优牌型
			return NORMAL_STRAIGHT, ListGet(straightLists[0], 0, 5)
		}
	}
	//三条
	if maxType >= NORMAL_THREE_KIND {
		threeCount, threeKindLists := GetKindLists(cards, 3)
		if threeCount > 0 {
			return NORMAL_THREE_KIND, threeKindLists[0]
		}
	}

	//两队/对子
	if maxType >= NORMAL_PAIR {
		pairCount, pairLists := GetKindLists(cards, 2)
		if pairCount > 2 {
			var twoPairs []int
			twoPairs = append(twoPairs, pairLists[0]...)
			twoPairs = append(twoPairs, pairLists[len(pairLists)-1]...)
			return NORMAL_TWO_PAIR, twoPairs
		}
		if pairCount > 1 {
			var twoPairs []int
			twoPairs = append(twoPairs, pairLists[0]...)
			twoPairs = append(twoPairs, pairLists[1]...)
			return NORMAL_TWO_PAIR, twoPairs
		}
		if pairCount > 0 {
			return NORMAL_PAIR, pairLists[0]
		}
	}

	//乌龙
	if len(cards) > 5 {
		tempCard := []int{cards[0]}
		tempCard = append(tempCard, ListGet(cards, len(cards)-4, 4)...)
		return NORMAL_ONE, tempCard
	}
	return NORMAL_ONE, cards
}

//推荐牌型
func RecommendPoker(cards []int, maxType int) ([]int, []int) {
	//第三墩
	threeType, threeCards := GetPokerType(cards, maxType)
	// logs.Debug("第三墩", threeType, threeCards)
	//第二墩
	secondLastPoker := ListDelList(cards, threeCards)
	secondType, secondCards := GetPokerType(secondLastPoker, threeType)
	// logs.Debug("第二墩", secondType, secondCards)
	//第一墩
	firstLastPoker := ListDelList(secondLastPoker, secondCards)
	firstType, firstCards := GetPokerType(firstLastPoker, NORMAL_THREE_KIND)
	if firstType == NORMAL_TWO_PAIR {
		firstType = NORMAL_PAIR
	}
	// logs.Debug("第一墩", firstType, firstCards)
	firstCards = ListGet(firstCards, 0, 3)
	// logs.Debug("第一墩", firstType, firstCards)
	//
	lastPoker := append([]int{}, ListDelList(firstLastPoker, firstCards)...)
	// logs.Debug("剩余", lastPoker)

	//补足手牌
	threeLength := len(threeCards)
	if threeLength < 5 && len(lastPoker) >= (5-threeLength) && (5-threeLength) > 0 {
		threeCards = append(threeCards, lastPoker[0:(5-threeLength)]...)
		lastPoker = append([]int{}, ListDelMore(lastPoker, 0, (5-threeLength))...)
		// logs.Debug("剩余手牌3", threeCards, lastPoker)
	}
	if threeType != ValidateFive(threeCards) {
		return nil, nil
	}

	secondLength := len(secondCards)
	if secondLength < 5 && len(lastPoker) >= (5-secondLength) && (5-secondLength) > 0 {
		secondCards = append(secondCards, lastPoker[0:(5-secondLength)]...)
		lastPoker = append([]int{}, ListDelMore(lastPoker, 0, (5-secondLength))...)
		// logs.Debug("剩余手牌2", secondCards, lastPoker)
	}
	if secondType != ValidateFive(secondCards) {
		return nil, nil
	}

	firstLength := len(firstCards)
	if firstLength < 3 && len(lastPoker) >= (3-firstLength) && (3-firstLength) > 0 {
		firstCards = append(firstCards, lastPoker[0:3-firstLength]...)
		lastPoker = append([]int{}, ListDelMore(lastPoker, 0, (3-firstLength))...)
		// logs.Debug("剩余手牌1", firstCards, lastPoker)
	}
	if firstType != ValidateThree(firstCards) {
		return nil, nil
	}

	//结果牌型数组
	var handCardType []int
	var handCards []int
	handCards = append(handCards, firstCards...)
	handCardType = append(handCardType, firstType)

	//中墩
	ct2 := GCardsType{
		Type:  secondType,
		Cards: secondCards,
	}
	//底墩
	ct3 := GCardsType{
		Type:  threeType,
		Cards: threeCards,
	}
	rs2 := MCardsCompare(ct2, ct3)
	if rs2 == 1 {
		//第三墩
		handCards = append(handCards, threeCards...)
		handCardType = append(handCardType, threeType)
		//第二墩
		handCards = append(handCards, secondCards...)
		handCardType = append(handCardType, secondType)

	} else {
		//第二墩
		handCards = append(handCards, secondCards...)
		handCardType = append(handCardType, secondType)
		//第三墩
		handCards = append(handCards, threeCards...)
		handCardType = append(handCardType, threeType)
	}
	return handCardType, handCards
}

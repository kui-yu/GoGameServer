package main

import (
	"logs"
)

//判断正常牌型
func ValidateFive(cards []int) int {

	if len(cards) != 5 {
		logs.Error("判断数组长度为5异常")
		return -1
	}

	//五同
	if GetCardValue(cards[0]) == GetCardValue(cards[4]) {
		return NORMAL_FIVE_KIND
	}
	//同花顺
	if CheckColor(cards) && CheckStraight(cards) {
		return NORMAL_COLOR_STRAIGHT
	}
	//铁支
	for i := 0; i < 2; i++ {
		fourKindCount := 0
		for j := 0; j < len(cards); j++ {
			if GetCardValue(cards[i]) == GetCardValue(cards[j]) {
				fourKindCount++
			}
			if fourKindCount >= 4 {
				return NORMAL_FOUR_KIND
			}
		}
	}
	//葫芦
	if GetCardValue(cards[0]) == GetCardValue(cards[2]) &&
		GetCardValue(cards[3]) == GetCardValue(cards[4]) {
		return NORMAL_GOURD
	}
	if GetCardValue(cards[0]) == GetCardValue(cards[1]) &&
		GetCardValue(cards[2]) == GetCardValue(cards[4]) {
		return NORMAL_GOURD
	}
	//同花
	if CheckColor(cards) {
		return NORMAL_SAME_COLOR
	}
	//顺子
	if CheckStraight(cards) {
		return NORMAL_STRAIGHT
	}
	//三条
	for i := 0; i < 3; i++ {
		threeKindCount := 0
		for j := 0; j < len(cards); j++ {
			if GetCardValue(cards[i]) == GetCardValue(cards[j]) {
				threeKindCount++
			}
			if threeKindCount >= 3 {
				return NORMAL_THREE_KIND
			}
		}
	}
	//两对
	if GetCardValue(cards[0]) == GetCardValue(cards[1]) &&
		GetCardValue(cards[2]) == GetCardValue(cards[3]) {
		return NORMAL_TWO_PAIR
	}
	if GetCardValue(cards[0]) == GetCardValue(cards[1]) &&
		GetCardValue(cards[3]) == GetCardValue(cards[4]) {
		return NORMAL_TWO_PAIR
	}
	if GetCardValue(cards[1]) == GetCardValue(cards[2]) &&
		GetCardValue(cards[3]) == GetCardValue(cards[4]) {
		return NORMAL_TWO_PAIR
	}
	//对子
	for i := 0; i < 4; i++ {
		pairCount := 0
		for j := 0; j < len(cards); j++ {
			if GetCardValue(cards[i]) == GetCardValue(cards[j]) {
				pairCount++
			}
			if pairCount >= 2 {
				return NORMAL_PAIR
			}
		}
	}
	//乌龙
	return NORMAL_ONE
}

//判断头墩牌型
func ValidateThree(cards []int) int {

	if len(cards) != 3 {
		logs.Error("判断数组长度为3异常")
		return -1
	}

	//三条
	threeKindCount := 0
	for j := 0; j < len(cards); j++ {
		if GetCardValue(cards[0]) == GetCardValue(cards[j]) {
			threeKindCount++
		}
		if threeKindCount >= 3 {
			return NORMAL_THREE_KIND
		}
	}

	//对子
	for i := 0; i < 2; i++ {
		pairCount := 0
		for j := 0; j < len(cards); j++ {
			if GetCardValue(cards[i]) == GetCardValue(cards[j]) {
				pairCount++
			}
			if pairCount >= 2 {
				return NORMAL_PAIR
			}
		}
	}
	//乌龙
	return NORMAL_ONE
}

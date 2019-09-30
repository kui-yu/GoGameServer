package main

import (
	"math/rand"
	"time"
)

//牌类型
const (
	CARDTYPE_SANGONG = 10 //3//三公
	CARDTYPE_ZHADAN  = 11 //4//炸弹
	CARDTYPE_BAOJIU  = 12 //5//爆玖
)

//走势记录
type Trend struct {
	Player   int   //1为庄家，0为闲家
	WinCoins int64 //输赢金币
}

//牌结构
type Card struct {
	CardValue []int //牌
	CardType  int   //10=三公，11=炸弹，12=爆玖,0到9对应点数
	Multiple  int64 //倍数
}

//初始化卡组
func InitCards() []int {
	cards := make([]int, 0)
	for i := 1; i <= 13; i++ {
		cards = append(cards, i+0x00, i+0x10, i+0x20, i+0x30)
	}
	return cards
}

//打乱卡组
func DisturbCards(cards []int) []int {
	newcards := make([]int, len(cards))
	rand.Seed(time.Now().UnixNano())
	randArr := rand.Perm(len(cards))
	for k, v := range randArr {
		newcards[k] = cards[v]
	}
	return newcards
}

//发牌
func (this *ExtDesk) SendCards() {
	allHandCards := this.GetHandCards(3) //所有的玩家牌
	for k, v := range this.Players {
		v.HandCards = allHandCards[k] //发牌
	}
}

//获取指定人数的手牌,num等于手牌数
func (this *ExtDesk) GetHandCards(num int) []Card {
	cards := make([]Card, 0)
	count := 0
	for i := 0; i < len(this.Players); i++ {
		cardsValue := append([]int{}, this.DeskCards[count:count+num]...) //卡值
		CardSort(cardsValue)                                              //排序
		cardsType := GetCardType(cardsValue)                              //卡类型
		cards = append(cards, Card{
			CardValue: cardsValue,
			CardType:  cardsType,
			Multiple:  GetCardMultiple(cardsType),
		})
		count += num
	}
	return cards
}

//获取牌的数值
func GetCradValue(card int) int {
	return card % 16
}

//获取牌的花色
func GetCardColor(card int) int {
	return card / 16
}

//获取手牌类型
func GetCardType(cards []int) int {
	if IsBaoJiu(cards) {
		return CARDTYPE_BAOJIU
	} else if IsZhaDan(cards) {
		return CARDTYPE_ZHADAN
	} else if IsSanGong(cards) {
		return CARDTYPE_SANGONG
	} else {
		return GetDianCount(cards)
	}
}

//判断爆玖
func IsBaoJiu(cards []int) bool {
	if GetCradValue(cards[0]) == GetCradValue(cards[2]) &&
		GetCradValue(cards[1]) == GetCradValue(cards[2]) &&
		GetCradValue(cards[0]) == 3 {
		return true
	} else {
		return false
	}
}

//判断炸弹
func IsZhaDan(cards []int) bool {
	if GetCradValue(cards[0]) == GetCradValue(cards[2]) &&
		GetCradValue(cards[1]) == GetCradValue(cards[2]) &&
		GetCradValue(cards[0]) != 3 {
		return true
	} else {
		return false
	}
}

//判断三公
func IsSanGong(cards []int) bool {
	if GetGongCardCount(cards) != 3 {
		return false
	}
	if !(GetCradValue(cards[0]) == GetCradValue(cards[2]) && GetCradValue(cards[1]) == GetCradValue(cards[2])) {
		return true
	} else {
		return false
	}
}

//获取公牌数量
func GetGongCardCount(cards []int) int {
	count := 0
	for _, v := range cards {
		if GetCradValue(v) >= 11 && GetCradValue(v) <= 13 {
			count++
		}
	}
	return count
}

//获取点数牌点数
func GetDianCount(cards []int) int {
	fcards := append([]int{}, cards...)
	count := 0
	for k, v := range fcards {
		if GetCradValue(v) > count && GetCradValue(v) <= 10 {
			count = v
		}
		if GetCradValue(v) >= 11 {
			fcards[k] = 10
		}
	}
	return (GetCradValue(fcards[0]) + GetCradValue(fcards[1]) + GetCradValue(fcards[2])) % 10
}

//获取牌型倍数
func GetCardMultiple(count int) int64 {
	switch {
	case count == 10:
		return 3
	case count == 11:
		return 4
	case count == 12:
		return 5
	case count >= 7:
		return 2
	default:
		return 1
	}
}

//根据数值和花色排序
func CardSort(cards []int) {
	for i := 0; i < len(cards)-1; i++ {
		for j := 0; j < len(cards)-1-i; j++ {
			if GetCradValue(cards[j]) < GetCradValue(cards[j+1]) {
				cards[j], cards[j+1] = cards[j+1], cards[j]
			} else if GetCradValue(cards[j]) == GetCradValue(cards[j+1]) {
				if GetCardColor(cards[j]) < GetCardColor(cards[j+1]) {
					cards[j], cards[j+1] = cards[j+1], cards[j]
				}
			}
		}
	}
}

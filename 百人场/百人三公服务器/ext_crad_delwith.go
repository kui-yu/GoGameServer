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
	CardType int //牌型
	Player   int //1为庄家，0为闲家
}

//牌结构
type Card struct {
	CardValue []int
	CardType  int
	Multiple  int64
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
	count := 0
	handCards := make([]Card, 0)
	for k := 0; k < len(this.HandCards); k++ {
		cards := append(this.HandCards[k].CardValue, this.DeskCards[count:count+3]...)
		CardSort(cards)
		handCards = append(handCards, Card{
			CardValue: cards,                               //牌值
			CardType:  GetCardType(cards),                  //牌型
			Multiple:  GetCardMultiple(GetCardType(cards)), //牌对应的倍数
		})
		count += 3
	}
	this.DeskCards = this.DeskCards[len(this.HandCards)*3:] //发完牌桌面牌要删除对应的牌
	//风控
	var isPlayerBet bool
	for _, v := range this.Players { //确定是否有真实玩家下注
		if !v.Robot && v.IsBet {
			isPlayerBet = true
			break
		}
	}
	if !isPlayerBet && this.IsHaveBet() && BankerLose() { //如果没真实玩家下注，有一点几率控制庄家输
		this.GetWinResult(handCards, 1)
		return
	}
	curCd := CalPkAll(StartControlTime, time.Now().Unix())
	if CD-curCd >= 0 || !isPlayerBet || GetCostType() == 2 { //不进风控
		for k := range this.HandCards {
			this.HandCards[k] = handCards[k]
		}
	} else { //进入风控
		this.GetWinResult(handCards, 0) //换牌
	}
}

//判断是否有下注
func (this *ExtDesk) IsHaveBet() bool {
	var count int64
	for _, v := range this.PlaceBet {
		count += v
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

//百分75概率
func BankerLose() bool {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(100) + 1
	if r <= 75 {
		return true
	} else {
		return false
	}
}

//统计
//递归计算能赢的牌
func (this *ExtDesk) GetWinResult(cards []Card, c int) {
	var winCoins int64 = 0
	for k := range this.HandCards {
		this.HandCards[k] = cards[k]
	}
	zCards := this.HandCards[0]
	roundResult := make([]int, 4)
	for k, v := range this.HandCards[1:] {
		if zCards.CardType > v.CardType {
			roundResult[k] = -1
			continue
		} else if zCards.CardType == v.CardType {
			if GetCradValue(zCards.CardValue[0]) > GetCradValue(v.CardValue[0]) {
				roundResult[k] = -1
				continue
			} else if GetCradValue(zCards.CardValue[0]) == GetCradValue(v.CardValue[0]) {
				if GetCardColr(zCards.CardValue[0]) > GetCardColr(v.CardValue[0]) {
					roundResult[k] = -1
					continue
				} else {
					roundResult[k] = 1
					continue
				}
			}
		}
		roundResult[k] = 1
	}
	for _, v := range this.Players {
		if !v.IsBet {
			continue
		}
		for l, k := range roundResult {
			if v.PlaceBet[l] <= 0 {
				continue
			}
			if k < 0 {
				if v.IsDouble {
					winCoins += zCards.Multiple * v.PlaceBet[l]
				} else {
					winCoins += v.PlaceBet[l]
				}
			} else {
				if v.IsDouble {
					winCoins -= zCards.Multiple * v.PlaceBet[l]
				} else {
					winCoins -= v.PlaceBet[l]
				}
			}
		}
	}
	if c == 0 {
		if winCoins > 0 {
			return
		} else {
			c := append([]Card{}, cards[len(cards)-1])
			c = append(c, cards[:len(cards)-1]...)
			this.GetWinResult(c, 0)
		}
	} else if c == 1 {
		if winCoins < 0 {
			return
		} else {
			c := append([]Card{}, cards[len(cards)-1])
			c = append(c, cards[:len(cards)-1]...)
			this.GetWinResult(c, 1)
		}
	}
}

//牌排序
func PlayersCardsSort(cards []Card) {
	for i := 0; i < len(cards)-1; i++ {
		for j := 0; j < len(cards)-1-i; j++ {
			if cards[j].CardType < cards[j+1].CardType {
				cards[j], cards[j+1] = cards[j+1], cards[j]
			} else if cards[j].CardType == cards[j+1].CardType {
				if GetCradValue(cards[j].CardValue[0]) < GetCradValue(cards[j+1].CardValue[0]) {
					cards[j], cards[j+1] = cards[j+1], cards[j]
				} else if GetCradValue(cards[j].CardValue[0]) == GetCradValue(cards[j+1].CardValue[0]) {
					if GetCardColr(cards[j].CardValue[0]) < GetCardColr(cards[j+1].CardValue[0]) {
						cards[j], cards[j+1] = cards[j+1], cards[j]
					}
				}
			}
		}
	}
}

//获取牌的数值
func GetCradValue(card int) int {
	return card % 16
}

//获取牌的花色
func GetCardColr(card int) int {
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

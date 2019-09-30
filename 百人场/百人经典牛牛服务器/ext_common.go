package main

import (
	"crypto/rand"

	// "fmt"
	"math/big"
)

// 牌的类型
const (
	Card_Fang = 0x00
	Card_Mei  = 0x10
	Card_Hong = 0x20
	Card_Hei  = 0x30
)

type CardGroupType int

const (
	_ CardGroupType = iota
	CardGroupType_Cattle_1
	CardGroupType_Cattle_2
	CardGroupType_Cattle_3
	CardGroupType_Cattle_4
	CardGroupType_Cattle_5
	CardGroupType_Cattle_6
	CardGroupType_Cattle_7     // 2倍
	CardGroupType_Cattle_8     // 3倍
	CardGroupType_Cattle_9     // 3倍
	CardGroupType_Cattle_C     // 4倍
	CardGroupType_Cattle_BOMB  // 炸弹 5倍
	CardGroupType_Cattle_WUHUA // 五花牛 不包括10 6倍
	CardGroupType_None
	CardGroupType_NotCattle
)

var BetDoubleMap map[CardGroupType]int32

func init() {
	BetDoubleMap = make(map[CardGroupType]int32)
	BetDoubleMap[CardGroupType_NotCattle] = 1
	BetDoubleMap[CardGroupType_Cattle_1] = 1
	BetDoubleMap[CardGroupType_Cattle_2] = 1
	BetDoubleMap[CardGroupType_Cattle_3] = 1
	BetDoubleMap[CardGroupType_Cattle_4] = 1
	BetDoubleMap[CardGroupType_Cattle_5] = 1
	BetDoubleMap[CardGroupType_Cattle_6] = 1
	BetDoubleMap[CardGroupType_Cattle_7] = 2
	BetDoubleMap[CardGroupType_Cattle_8] = 3
	BetDoubleMap[CardGroupType_Cattle_9] = 3
	BetDoubleMap[CardGroupType_Cattle_C] = 4
	BetDoubleMap[CardGroupType_Cattle_BOMB] = 5
	BetDoubleMap[CardGroupType_Cattle_WUHUA] = 5 //五花牛，炸弹为5倍
}

// 牌值获取函数
func GetCardColor(card byte) byte {
	return (card & 0xF0) >> 4
}

// 洗牌
func (this *ExtDesk) ShuffleCard() {
	var cards []uint8 = []uint8{}
	for i := 0; i < 4; i++ {
		for j := 1; j < 14; j++ {
			cards = append(cards, uint8((i<<4)|j))
		}
	}

	clen := 52

	for i := 0; i < clen; i++ {
		r, _ := GetRandomNum(i, clen)
		cards[i], cards[r] = cards[r], cards[i]
	}

	// //自定义牌组
	// var cards []uint8 = []uint8{}
	// cards = append(cards, []uint8{0x21 + 11, 0xa, 0x11 + 0x09, 0x31 + 5, 0x31 + 4}...)
	// cards = append(cards, []uint8{0x31 + 12, 0x31 + 10, 0x01 + 7, 0x11 + 2, 0x31}...)
	// cards = append(cards, []uint8{0x11 + 8, 0x21 + 8, 0x11 + 1, 0x21 + 5, 0x31 + 2}...)
	// cards = append(cards, []uint8{0x21 + 4, 0x01 + 3, 0x11, 0x11 + 12, 0x11 + 11}...)
	// cards = append(cards, []uint8{0x21 + 12, 0x21 + 10, 0x01 + 5, 0x11 + 4, 0x01 + 2}...)
	this.DownCards = cards
}

// 发牌
func (this *ExtDesk) SendCard(num int) []int {
	//添加要发的手牌
	list := append([]uint8{}, this.DownCards[:num]...)
	//删除要发的手牌元素
	this.DownCards = append([]uint8{}, this.DownCards[num:]...)

	var sourceList []int
	for _, card := range list {
		sourceList = append(sourceList, int(card))
	}
	//返回
	return sourceList
}

//回收手牌
func (this *ExtDesk) RecoverCard(cards []int) {
	for _, card := range cards {
		this.DownCards = append(this.DownCards, uint8(card))
	}
}

//打乱牌组
func (this *ExtDesk) UpsetCard() {

	var sourceList []uint8

	MVCard := append([]uint8{}, this.DownCards...)

	// 随机打乱牌型
	for i := 0; i < len(this.DownCards); i++ {
		//打乱
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(MVCard))))
		//添加
		sourceList = append(sourceList, MVCard[randIndex.Int64()])
		// 移除已经添加的牌
		MVCard = append(MVCard[:int(randIndex.Int64())], MVCard[int(randIndex.Int64())+1:]...)
	}

	this.DownCards = sourceList
}

// 算出牌的类型
func CalcCards(cards []int) (CardGroupType, uint8, []int) {
	cards = OrderCard(cards)
	var nums [5]int
	var maxCard int = 0
	groupType := CardGroupType_NotCattle

	for i, card := range cards {
		if (card & 0xF) > (maxCard & 0xF) {
			maxCard = card
		}
		if ((card & 0xF) == (maxCard & 0xF)) && ((card >> 4) > (maxCard >> 4)) {
			maxCard = card
		}

		if (card & 0xF) > 9 {
			nums[i] = 10
		} else {
			nums[i] = card & 0xF
		}
	}

	// 是否五花牛 不包括10, 判断排序后最后一张是否大于10
	if cards[4]&0xF > 10 {
		return CardGroupType_Cattle_WUHUA, uint8(maxCard), cards
	}
	// 是否炸弹
	var cardKeyMap map[int]int = make(map[int]int)
	for _, card := range cards {
		if _, ok := cardKeyMap[card&0xF]; ok {
			cardKeyMap[card&0xF] += 1
		} else {
			cardKeyMap[card&0xF] = 1
		}
	}
	for _, num := range cardKeyMap {
		if num == 4 {
			return CardGroupType_Cattle_BOMB, uint8(maxCard), cards
		}
	}

	// 找出牛
	isFindNiu := false
	for i := 0; i < 5; i++ {
		if isFindNiu {
			break
		}
		for j := i + 1; j < 5; j++ {
			if isFindNiu {
				break
			}
			for k := j + 1; k < 5; k++ {
				n3 := nums[i] + nums[j] + nums[k]

				// 找出牛
				if n3%10 == 0 {
					sumnum := 0
					cards[i], cards[0] = cards[0], cards[i]
					cards[j], cards[1] = cards[1], cards[j]
					cards[k], cards[2] = cards[2], cards[k]

					if (cards[3] & 0xF) < (cards[4] & 0xF) {
						cards[3], cards[4] = cards[4], cards[3]
					}

					for w := 0; w < 5; w++ {
						if w != i && w != j && w != k {
							sumnum += int(nums[w])
						}
					}

					groupType = CardGroupType(sumnum % 10)
					if groupType == 0 {
						groupType = CardGroupType_Cattle_C
					}

					isFindNiu = true
					break
				}

			}
		}
	}

	return groupType, uint8(maxCard), cards
}

// 牌排序
func OrderCard(cards []int) []int {
	tempCards := cards
	clen := len(tempCards)
	for i := 0; i < clen; i++ {
		for j := i + 1; j < clen; j++ {
			if tempCards[i]&0xF < tempCards[j]&0xF {
				tempCards[i], tempCards[j] = tempCards[j], tempCards[i]
			}
		}
	}
	return tempCards
}

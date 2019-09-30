package main

import (
	"math/rand"
	"time"
)

const (
	CARD_COLOR = 0xF0 //花色掩码
	CARD_VALUE = 0x0F //数值掩码
)

const (
	//方块
	Card_Fang_1 = iota + 0x01
	Card_Fang_2
	Card_Fang_3
	Card_Fang_4
	Card_Fang_5
	Card_Fang_6
	Card_Fang_7
	Card_Fang_8
	Card_Fang_9
	Card_Fang_10
	Card_Fang_J
	Card_Fang_Q
	Card_Fang_K
)

const (
	//梅花
	Card_Mei_1 = iota + 0x11
	Card_Mei_2
	Card_Mei_3
	Card_Mei_4
	Card_Mei_5
	Card_Mei_6
	Card_Mei_7
	Card_Mei_8
	Card_Mei_9
	Card_Mei_10
	Card_Mei_J
	Card_Mei_Q
	Card_Mei_K
)

const (
	//红桃
	Card_Hong_1 = iota + 0x21
	Card_Hong_2
	Card_Hong_3
	Card_Hong_4
	Card_Hong_5
	Card_Hong_6
	Card_Hong_7
	Card_Hong_8
	Card_Hong_9
	Card_Hong_10
	Card_Hong_J
	Card_Hong_Q
	Card_Hong_K
)

const (
	//黑桃
	Card_Hei_1 = iota + 0x31
	Card_Hei_2
	Card_Hei_3
	Card_Hei_4
	Card_Hei_5
	Card_Hei_6
	Card_Hei_7
	Card_Hei_8
	Card_Hei_9
	Card_Hei_10
	Card_Hei_J
	Card_Hei_Q
	Card_Hei_K
)

const (
	CARD_COLOR_Fang = iota
	CARD_COLOR_Mei
	CARD_COLOR_Hong
	CARD_COLOR_Hei
)

const (
	NIU_FIVE_SMALL = 14 //五小牛
	NIU_BOOM       = 13 //四炸
	NIU_FIVE_COLOR = 12 //五花牛
	NIU_FORE_COLOR = 11 //四花牛
	NIU_NIU        = 10 //牛牛
	NIU_NINE       = 9  //牛9
	NIU_EIGHT      = 8  //牛8
	NIU_SEVEN      = 7  //牛7
	NIU_SIX        = 6  //牛6
	NIU_FIVE       = 5  //牛5
	NIU_FORE       = 4  //牛4
	NIU_THREE      = 3  //牛3
	NIU_TWO        = 2  //牛2
	NIU_ONE        = 1  //牛1
	NIU_ZERO       = 0  //无牛
)

/////////////////////////////////////////////////////////
//卡牌管理器，负责做牌
type MgrCard struct {
	MVCard       []int32
	MVSourceCard []int32
}

//初始化
func (this *MgrCard) InitCards() {
	this.MVCard = []int32{}
	this.MVSourceCard = []int32{}
}

//赋值
func (this *MgrCard) InitNormalCards() {
	begaincard := []int32{Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1}
	for _, v := range begaincard {
		for j := int32(0); j < 13; j++ {
			this.MVCard = append(this.MVCard, v+j)
		}
	}
}

//洗牌
func (this *MgrCard) Shuffle() {
	this.MVSourceCard = append([]int32{}, this.MVCard...)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(this.MVCard))
	for i, randIndex := range perm {
		this.MVSourceCard[i] = this.MVCard[randIndex]
	}
}

//发手牌
func (this *MgrCard) SendHandCard(num int) []int32 {

	list := append([]int32{}, this.MVSourceCard[0:num]...)
	//删除第0个元素开始
	i := 0
	this.MVSourceCard = append(this.MVSourceCard[:i], this.MVSourceCard[i+num:]...)
	//返回
	return list
}

// 获取卡牌花色
func GetCardColor(card int32) int32 {
	return (card & CARD_COLOR) >> 4
}

// 获取卡牌值
func GetCardValue(card int32) int32 {
	return (card & CARD_VALUE)
}

//转化值
func GetLogicValue(card int32) int32 {
	d := GetCardValue(card)

	if d > 9 {
		return 10
	}
	return d
}

//获取卡牌值
func GetSmallValue(cards []int32) []int32 {
	var smallCards []int32
	for i := 0; i < len(cards); i++ {
		vi := GetCardValue(cards[i])
		smallCards = append(smallCards, vi)
	}
	return smallCards
}

//获取卡牌牛牛值
func GetSmallNiuValue(cards []int32) []int32 {
	var smallCards []int32
	for i := 0; i < len(cards); i++ {
		vi := GetLogicValue(cards[i])
		smallCards = append(smallCards, vi)
	}
	return smallCards
}

//排序
func Sort(cards []int32) []int32 {
	cs := append([]int32{}, cards...)
	for i := 0; i < len(cs)-1; i++ {
		for j := i + 1; j < len(cs); j++ {
			vi := GetCardValue(cs[i])
			vj := GetCardValue(cs[j])
			if vi < vj || ((vi == vj) && (GetCardColor(cs[i]) > GetCardColor(cs[j]))) {
				vt := cs[i]
				cs[i] = cs[j]
				cs[j] = vt
			}
		}
	}
	return cs
}

//算五小牛
func MathFiveSmall(cards []int32) int32 {
	var smallCards = GetSmallValue(cards)

	var sum int32

	for i := 0; i < len(smallCards); i++ {
		if smallCards[i] > 4 {
			return NIU_ZERO
		}
		sum += smallCards[i]
	}
	if sum < 10 {
		return NIU_FIVE_SMALL
	}
	return NIU_ZERO
}

//算四炸
func MathBoom(cards []int32) (int32, []int32) {

	niuCards := append([]int32{}, cards...)

	smallCards := GetSmallValue(cards)

	var count1 int32
	var count2 int32

	for i := 0; i < len(smallCards); i++ {
		if smallCards[0] == smallCards[i] {
			count1++
		}
		if smallCards[1] == smallCards[i] {
			count2++
		}
	}
	if count1 == 4 {
		for i := 0; i < len(smallCards); i++ {
			if smallCards[0] != smallCards[i] {
				tempCard := niuCards[i]
				//移除
				niuCards = append(niuCards[:i], niuCards[i+1:]...)
				//后面追加
				niuCards = append(niuCards, tempCard)
				break
			}
		}
		return NIU_BOOM, niuCards
	}
	if count2 == 4 {
		for i := 0; i < len(smallCards); i++ {
			if smallCards[1] != smallCards[i] {
				tempCard := niuCards[i]
				//移除
				niuCards = append(niuCards[:i], niuCards[i+1:]...)
				//后面追加
				niuCards = append(niuCards, tempCard)
				break
			}
		}
		return NIU_BOOM, niuCards
	}

	return NIU_ZERO, niuCards
}

//算五花牛/四花牛
func MathColor(cards []int32) (int32, []int32) {

	var smallCards = GetSmallValue(cards)

	var count int32
	for i := 0; i < len(smallCards); i++ {
		if smallCards[i] > 10 {
			count++
		}
	}
	//五花牛
	if count == 5 {
		return NIU_FIVE_COLOR, cards
	}
	//四花牛
	// if count == 4 {
	// 	return NIU_FORE_COLOR, cards
	// }
	return NIU_ZERO, cards
}

//算牛牛
func MathTen(cards []int32) (int32, []int32) {
	var smallCards = GetSmallNiuValue(cards)

	//牛点
	var niuPoint int32

	for i := 0; i < len(smallCards); i++ {
		for j := i + 1; j < len(smallCards); j++ {
			for k := j + 1; k < len(smallCards); k++ {
				//有牛
				if (smallCards[i]+smallCards[j]+smallCards[k])%10 == 0 {
					var fourKey int
					for l := 0; l < len(smallCards); l++ {
						if l != i && l != j && l != k {
							fourKey = l
							break
						}
					}
					var fiveKey int
					for l := 0; l < len(smallCards); l++ {
						if l != i && l != j && l != k && l != fourKey {
							fiveKey = l
							break
						}
					}
					//123
					niuPoint = (smallCards[fourKey] + smallCards[fiveKey]) % 10
					if niuPoint == 0 {
						niuPoint = NIU_NIU
					}
					temp := []int32{cards[i], cards[j], cards[k], cards[fourKey], cards[fiveKey]}
					return niuPoint, temp
				}
			}
		}
	}
	return niuPoint, cards
}

//获取结果
func GetResult(cards []int32) (int32, []int32) {

	var sortCards []int32 = Sort(cards)
	// logs.Debug("收到手牌", sortCards)

	//五小牛判断
	// fiveSmall := MathFiveSmall(sortCards)
	// if fiveSmall > 0 {
	// 	return fiveSmall, sortCards
	// }

	//炸弹判断
	boom, boomCards := MathBoom(sortCards)
	if boom > 0 {
		return boom, boomCards
	}

	//花牛
	color, colorCards := MathColor(sortCards)
	if color > 0 {
		return color, colorCards
	}

	//判断牛牛
	return MathTen(sortCards)
}

//比牌
func SoloResult(player1 GCardType, player2 GCardType) int32 {
	if player1.NiuPoint > player2.NiuPoint {
		return 1
	} else if player1.NiuPoint < player2.NiuPoint {
		return 2
	} else {
		//比大小
		if player1.NiuPoint == NIU_BOOM {
			//炸弹比较
			if player1.NiuCards[0] > player2.NiuCards[0] {
				return 1
			} else {
				return 2
			}
		} else {
			var sortCards1 []int32 = Sort(player1.HandCard)
			var sortCards2 []int32 = Sort(player2.HandCard)
			//比最大
			var card1 = GetCardValue(sortCards1[0])
			var card2 = GetCardValue(sortCards2[0])
			if card1 > card2 {
				return 1
			} else if card1 < card2 {
				return 2
			} else {
				//比花色
				if GetCardColor(sortCards1[0]) > GetCardColor(sortCards2[0]) {
					return 1
				} else {
					return 2
				}
			}
		}
	}
}

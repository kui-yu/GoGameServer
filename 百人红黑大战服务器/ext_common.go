package main

import (
	// "logs"
	"math/rand"
	"time"
)

//扑克牌花色牌值常量
const (
	//方
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
	//梅
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
	//红
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
	//黑
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

//花色权重常量
const (
	CARD_COLOR_Fang = iota
	CARD_COLOR_Mei
	CARD_COLOR_Hong
	CARD_COLOR_Hei
)
const (
	CARD_BOOM    = 6 //豹子
	CARD_FLUSH   = 5 //同花顺
	CARD_TONGHUA = 4 //同花
	CARD_SHUNZI  = 3 //顺子
	CARD_PAIR    = 2 //对子
	CARD_SINGLE  = 1 //单张
)
const (
	RED = 1 + iota
	BLACK
)

//花色获取函数
func GetCardColor(card int32) int32 {
	return (card & 0xF0)
}

//牌值获取函数
func GetCardValue(card int32) int32 {
	return (card & 0x0F)
}

//逻辑牌值获取函数
func GetLogicValue(card int32) int32 {
	d := GetCardValue(card)
	return d
}

//卡牌管理器，负责做牌
type MgrCard struct {
	MVCard       []byte //初始牌
	MVSourceCard []byte //打乱以后的牌
	CardValue    int32  //对子牌型的值
}

func (this *MgrCard) InitCards() {
	this.MVCard = []byte{}
	this.MVSourceCard = []byte{}
}

//初始牌
func (this *MgrCard) InitNormalCards() {
	begaincard := []byte{Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1}
	for _, v := range begaincard {
		for j := byte(0); j < 13; j++ {
			this.MVCard = append(this.MVCard, v+j)
		}
	}
}

//洗牌
func (this *MgrCard) Shuffle() {
	this.MVSourceCard = append([]byte{}, this.MVCard...)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(this.MVCard))
	for i, randIndex := range perm {
		this.MVSourceCard[i] = this.MVCard[randIndex]
	}
}

//发牌
func (this *MgrCard) HandCardInfo(num int) []byte {
	list := this.MVSourceCard[:num]
	this.MVSourceCard = ListDelMore(this.MVSourceCard, 0, num)
	var cardList []byte
	for i := 0; i < num/3; i++ {
		for j := 0; j < 3; j++ {
			cardList = append(cardList, list[i+j*num/3])
		}
	}
	return cardList
}

func ListDelMore(list []byte, from int, num int) []byte {
	sourceList := append([]byte{}, list...)
	sourceList = append(sourceList[:from], sourceList[from+num:]...)
	return sourceList
}

// 整理手牌
func SortHandCard(card1 []int32) ([]int32, []int32) {
	var card []int32
	card = append(card, card1[:]...)
	color := make([]int32, 0, 16)
	for i := 0; i < len(card); i++ {
		color = append(color, GetCardColor(card[i])/16)
		card[i] = GetCardValue(card[i])
	}

	for i := 0; i < len(card); i++ {
		if card[i] == 1 {
			card[i] = 14
		}
	}
	//对扑克牌进行排序
	for n := len(card); n > 0; n-- {
		for i := 0; i < n-1; i++ {
			if card[i] > card[i+1] {
				m := card[i]
				card[i] = card[i+1]
				card[i+1] = m
				c := color[i]
				color[i] = color[i+1]
				color[i+1] = c
			} else if card[i] == card[i+1] {
				if color[i] < color[i+1] {
					c := color[i]
					color[i] = color[i+1]
					color[i+1] = c
					m := card[i]
					card[i] = card[i+1]
					card[i+1] = m
				}
			}
		}
	}
	return card, color
}

//判断牌型是否为炸弹
func CheckBoom(card []int32) (bool, int32, int) { //返回为炸弹值和牌类型
	for i := 0; i < len(card); i++ {
		if card[0] != card[i] {
			return false, 0, 0
		}
	}
	return true, card[0], CARD_BOOM
}

//判断牌型是否为同花顺
func CheckFlush(card []int32, color []int32) (bool, int32, int) {
	for i := 0; i < len(color)-1; i++ { //判断是不是同花
		// for j := 0; j < len(color)-1; j++ {
		if color[i] != color[i+1] {
			return false, 0, 0
		}
		// }
	}
	Check := false
	for i := 0; i < len(card); i++ { //如果牌中有A 则特殊处理
		if card[i] == 14 {
			Check = true
			break
		}
	}
	if Check { //特殊处理A
		check2 := false
		for i := 1; i < 3; i++ {
			if card[i]-card[i-1] != 1 {
				check2 = true
				break
			}
			if i == 2 {
				return true, card[0], CARD_FLUSH
			}
		}
		if check2 {
			list := make([]int32, 0, 3)
			list = append(list, 1)
			list = append(list, card[0], card[1])
			for i := 1; i < 3; i++ {
				if list[i]-list[i-1] != 1 {
					return false, 0, 0
				}
			}
			return true, 1, CARD_FLUSH
		}
	}
	for i := 1; i < 3; i++ {
		if card[i]-card[i-1] != 1 {
			return false, 0, 0
		}
	}
	return true, card[0], CARD_FLUSH
}

//同花
func CheckTonghua(color []int32) (bool, int) {
	for i := 0; i < len(color); i++ {
		if color[0] != color[i] {
			return false, 0
		}
	}
	return true, CARD_TONGHUA
}

//顺子
func CheckShunzi(card []int32) (bool, int32, int) {
	ch := false
	for i := 0; i < 3; i++ {
		if card[i] == 14 {
			ch = true
			break
		}
	}
	if ch { //特殊处理A
		ch2 := false
		for i := 1; i < 3; i++ {
			if card[i]-card[i-1] != 1 {
				ch2 = true
				break
			}
			if i == 2 {
				return true, card[0], CARD_SHUNZI
			}
		}
		if ch2 {
			list := make([]int32, 0, 3)
			list = append(list, 1)
			list = append(list, card[0], card[1])
			for i := 1; i < 3; i++ {
				if list[i]-list[i-1] != 1 {
					return false, 0, 0
				}
			}
			return true, 1, CARD_SHUNZI
		}
	}

	for i := 1; i < 3; i++ {
		if card[i]-card[i-1] != 1 {
			return false, 0, 0
		}
	}
	return true, card[0], CARD_SHUNZI
}

//对子
func CheckPair(card []int32) (bool, int32, int) { //返回对子的值，牌类型
	mid := card[0]

	for i := 1; i < 3; i++ {
		if mid == card[i] {
			return true, mid, CARD_PAIR
		}
		if mid != card[i] {
			mid = card[i]
		}
	}
	return false, 0, 0
}

//特殊牌
func Check235(card []int32) int {
	list := []int32{2, 3, 5}
	for i := 0; i < 3; i++ {
		if card[i] != list[i] {
			return 0
		}
	}
	return 1
}

//获取牌的类型
func GetCardType(card []int32, color []int32) (int, int32) { //牌型，值
	if c, n, ctype := CheckBoom(card); c {
		return ctype, n
	} else if c, n, ctype := CheckFlush(card, color); c {
		return ctype, n
	} else if c, ctype := CheckTonghua(color); c {
		return ctype, 0
	} else if c, n, ctype := CheckShunzi(card); c {
		return ctype, n
	} else if c, n, ctype := CheckPair(card); c {
		return ctype, n
	} else {
		return CARD_SINGLE, 0
	}
}

//比牌操作
func (this *MgrCard) CompareCard(card1 []int32, color1 []int32, card2 []int32, color2 []int32) (int32, int) {
	this.CardValue = 0
	ctype1, n1 := GetCardType(card1, color1) //红方牌的类型
	ctype2, n2 := GetCardType(card2, color2) //黑方牌的类型
	if ctype1 == 6 || ctype2 == 6 {          //特殊牌处理
		special1 := Check235(card1)
		special2 := Check235(card2)
		if special1 == 1 {
			return RED, ctype1
		} else if special2 == 1 {
			return BLACK, ctype2
		}
	}
	if ctype1 > ctype2 {
		if ctype1 == 2 {
			this.CardValue = n1
		}
		return RED, ctype1
	} else if ctype2 > ctype1 {
		if ctype2 == 2 {
			this.CardValue = n2
		}
		return BLACK, ctype2
	} else {
		if ctype1 != 1 && ctype1 != 4 {
			if n1 > n2 {
				if ctype1 == 2 {
					this.CardValue = n1
				}
				return RED, ctype1
			} else if n2 > n1 {
				if ctype2 == 2 {
					this.CardValue = n2
				}
				return BLACK, ctype2
			} else {
				if ctype1 == 2 {
					for i := 0; i < 3; i++ {
						if card1[i] > card2[i] {
							this.CardValue = n1
							return RED, ctype1
						} else if card1[i] < card2[i] {
							this.CardValue = n2
							return BLACK, ctype2
						}
					}
				}
				for i := 2; i > 0; i-- {
					if color1[i] > color2[i] {
						return RED, ctype1
					} else {
						return BLACK, ctype2
					}
				}
			}
		} else {
			for i := 2; i >= 0; i-- {
				// fmt.Println(card1[i], card2[i])
				if card1[i] > card2[i] {
					// fmt.Println("红方获胜", card1[i])
					return RED, ctype1
				} else if card2[i] > card1[i] {
					// fmt.Println("黑方获胜", card2[i])
					return BLACK, ctype2
				}
			}
			for i := 2; i > 0; i-- {
				if color1[i] > color2[i] {
					return RED, ctype1
				} else {
					return BLACK, ctype2
				}
			}
		}

	}
	return BLACK, ctype2
}

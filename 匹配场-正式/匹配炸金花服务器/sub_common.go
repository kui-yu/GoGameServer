package main

//import "fmt"

const (
	CARD_COLOR = 0xF0 //花色掩码
	CARD_VALUE = 0x0F //数值掩码
)

const (
	CARD_TAO  = 0x31
	CARD_XIN  = 0x21
	CARD_HUA  = 0x11
	CARD_FANG = 0x01
)

const (
	CARD_BOOM    = 6 //炸弹
	CARD_FLUSH   = 5 //同花顺
	CARD_TONGHUA = 4 //同花
	CARD_SHUNZI  = 3 //顺子
	CARD_PAIR    = 2 //对子
	CARD_SINGLE  = 1 //单张
)

//卡牌管理器，负责做牌，数据结构
type MgrCard struct {
	MVSourceCard []int
}

//洗牌
func (this *MgrCard) Shuffle() {
	this.MVSourceCard = []int{}
	cardNew := []int{CARD_TAO, CARD_XIN, CARD_HUA, CARD_FANG}
	for i := 0; i < 13; i++ {
		for _, v := range cardNew {
			this.MVSourceCard = append(this.MVSourceCard, v+i)
		}
	}
	this.MVSourceCard = ListShuffle(this.MVSourceCard)
	// fmt.Println(len(this.MVCard), "牌长度", this.MVCard)
}

//发牌
func (this *MgrCard) HandCardInfo(num int) [][]int {
	list := this.MVSourceCard[:num]
	//删除第0个元素开始----ak
	this.MVSourceCard = ListDelMore(this.MVSourceCard, 0, num)

	// fmt.Println(list)
	var cardList [][]int
	for i := 0; i < num/3; i++ {
		playerlist := make([]int, 0, 3)
		for j := 0; j < 3; j++ {
			playerlist = append(playerlist, list[i+j*num/3])
		}
		cardList = append(cardList, playerlist)
	}
	// fmt.Println(cardList)
	return cardList
}

// 获取卡牌值
func GetCardValue(card int) int {
	return (card & CARD_VALUE)
}

// 获取卡牌花色
func GetCardColor(card int) int {
	return (card & CARD_COLOR)
}

//整理手牌
func SortHandCard(card1 []int) ([]int, []int) {
	var card []int
	card = append(card, card1[:]...)
	color := make([]int, 0, 16)
	for i := 0; i < len(card); i++ {
		color = append(color, GetCardColor(card[i])/16)
		card[i] = GetCardValue(card[i])
	}

	for i := 0; i < len(card); i++ {
		if card[i] == 1 {
			card[i] = 14
		}
	}

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

//炸弹
func CheckBoom(card []int) (bool, int, int) { //炸弹值,牌类型
	for i := 0; i < 3; i++ {
		if card[i] != card[0] {
			return false, 0, 0
		}
	}
	return true, card[0], CARD_BOOM
}

//同花顺
// func CheckFlush(card []int, color []int) (bool, int, int) { //最小值，牌类型
// 	for i := 0; i < len(color); i++ {
// 		if color[0] != color[i] {
// 			return false, 0, 0
// 		}
// 	}

// 	check := false
// 	for i := 0; i < len(color); i++ {
// 		if card[i] == 14 {
// 			check = true
// 			break
// 		}
// 	}

// 	if check { //特殊处理A
// 		check2 := false
// 		for i := 1; i < 3; i++ {
// 			if card[i]-card[i-1] != 1 {
// 				check2 = true
// 				break
// 			}
// 			if i == 2 {
// 				return true, card[0], CARD_FLUSH
// 			}
// 		}
// 		if check2 {
// 			list := make([]int, 0, 3)
// 			list = append(list, 1)
// 			list = append(list, card[0], card[1])
// 			for i := 1; i < 3; i++ {
// 				if list[i]-list[i-1] != 1 {
// 					return false, 0, 0
// 				}
// 			}
// 			return true, 1, CARD_FLUSH
// 		}
// 	}

// 	for i := 1; i < 3; i++ {
// 		if card[i]-card[i-1] != 1 {
// 			return false, 0, 0
// 		}
// 	}
// 	return true, card[0], CARD_FLUSH
// }
//改
func CheckFlush(card []int, color []int) (bool, int, int) {
	b, _, _ := CheckShunzi(card)
	if b && color[0] == color[1] && color[1] == color[2] {
		return true, card[0], CARD_FLUSH
	} else {
		return false, 0, 0
	}
}

//金花
func CheckTonghua(color []int) (bool, int) { //牌类型
	for i := 1; i < 3; i++ {
		if color[0] != color[i] {
			return false, 0
		}
	}
	return true, CARD_TONGHUA
}

//顺子
func CheckShunzi(card []int) (bool, int, int) { //最小值，牌类型
	check := false
	for i := 0; i < 3; i++ {
		if card[i] == 14 {
			check = true
			break
		}
	}

	if check { //特殊处理A
		//
		if card[0]-card[1] == 11 && card[1]-card[2] == 1 ||
			card[0]-card[1] == 1 && card[1]-card[2] == 1 {
			return true, card[0], CARD_SHUNZI
		}
		//
		// check2 := false
		// for i := 1; i < 3; i++ {
		// 	if card[i]-card[i-1] != -1 {
		// 		check2 = true
		// 		break
		// 	}
		// 	if i == 2 {
		// 		return true, card[0], CARD_SHUNZI
		// 	}
		// }
		// if check2 {
		// 	list := make([]int, 0, 3)
		// 	list = append(list, 1)
		// 	list = append(list, card[0], card[1])
		// 	for i := 1; i < 3; i++ {
		// 		if list[i]-list[i-1] != -1 {
		// 			return false, 0, 0
		// 		}
		// 	}
		// 	return true, 1, CARD_SHUNZI
		// }
	}

	for i := 1; i < 3; i++ {
		if card[i]-card[i-1] != -1 {
			return false, 0, 0
		}
	}
	return true, card[0], CARD_SHUNZI
}

//对子
func CheckPair(card []int) (bool, int, int) { //对子值，牌类型
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

//特殊牌235
func Check235(card []int) int {
	list := []int{2, 3, 5}
	for i := 0; i < 3; i++ {
		if card[i] != list[i] {
			return 0
		}
	}
	return 1
}

//比牌
//返回获胜牌位置
func GetResult(card1 []int, color1 []int, card2 []int, color2 []int) int {
	ctype1, n1 := GetCardType(card1, color1)
	ctype2, n2 := GetCardType(card2, color2)
	//fmt.Println(ctype1, ctype2, n1, n2, "输出牌类型和最小牌")
	if ctype1 == 6 || ctype2 == 6 { //特殊牌处理
		special1 := Check235(card1)
		special2 := Check235(card2)
		if special1 == 1 {
			return 0
		} else if special2 == 1 {
			return 1
		}
	}
	if ctype1 > ctype2 {
		return 0
	} else if ctype2 > ctype1 {
		return 1
	} else {
		if ctype1 != 1 && ctype1 != 4 {
			if n1 > n2 {
				return 0
			} else if n2 > n1 {
				return 1
			} else {
				if ctype1 == 2 {
					for i := 0; i < 3; i++ {
						if card1[i] > card2[i] {
							return 0
						} else if card1[i] < card2[i] {
							return 1
						}
					}
				}
				for i := 2; i > 0; i-- {
					if color1[i] > color2[i] {
						return 0
					} else {
						return 1
					}
				}
			}
		} else {
			// fmt.Println(card1, card2)
			// for i := 2; i >= 0; i-- {
			// 	// fmt.Println(card1[i], card2[i])
			// 	if card1[i] > card2[i] {
			// 		// fmt.Println("玩家0获胜", card1[i])
			// 		return 0
			// 	} else if card2[i] > card1[i] {
			// 		// fmt.Println("玩家1获胜", card2[i])
			// 		return 1
			// 	}
			// }
			//改,ayu
			for i := 0; i < 3; i++ {
				if card1[i] > card2[i] {
					// fmt.Println("玩家0获胜", card1[i])
					return 0
				} else if card2[i] > card1[i] {
					// fmt.Println("玩家1获胜", card2[i])
					return 1
				}
			}
			for i := 2; i > 0; i-- {
				if color1[i] > color2[i] {
					return 0
				} else {
					return 1
				}
			}
		}

	}
	return 0
}

//获取牌类型
func GetCardType(card []int, color []int) (int, int) { //牌型，值
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

func GetLogicValue(card int) int {
	d := GetCardValue(card)
	if card == 0x41 {
		return 16
	}
	if card == 0x42 {
		return 17
	}
	if d < 2 {
		return d + 13
	}
	return d
}

//排序
func Sort(cs []int) []int {
	for i := 0; i < len(cs)-1; i++ {
		for j := i + 1; j < len(cs); j++ {
			vi := GetLogicValue(cs[i])
			vj := GetLogicValue(cs[j])
			if vi < vj || ((vi == vj) && (GetCardColor(cs[i]) < GetCardColor(cs[j]))) {
				vt := cs[i]
				cs[i] = cs[j]
				cs[j] = vt
			}
		}
	}
	return append([]int{}, cs...)
}

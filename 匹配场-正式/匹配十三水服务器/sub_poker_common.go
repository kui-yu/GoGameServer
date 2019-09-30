package main

const (
	CARD_COLOR = 0xF0 //花色掩码
	CARD_VALUE = 0x0F //数值掩码
)

//方块
const (
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

//梅花
const (
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

//红桃
const (
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

//黑桃
const (
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

//花色值
const (
	CARD_COLOR_Fang = iota
	CARD_COLOR_Mei
	CARD_COLOR_Hong
	CARD_COLOR_Hei
)

//正常牌型
const (
	NORMAL_ONE            = 1  //乌龙
	NORMAL_PAIR           = 2  //对子
	NORMAL_TWO_PAIR       = 3  //两对
	NORMAL_THREE_KIND     = 4  //三条
	NORMAL_STRAIGHT       = 5  //顺子
	NORMAL_SAME_COLOR     = 6  //同花
	NORMAL_GOURD          = 7  //葫芦
	NORMAL_FOUR_KIND      = 8  //铁支
	NORMAL_COLOR_STRAIGHT = 9  //同花顺
	NORMAL_FIVE_KIND      = 10 //五同
)

//特殊牌型
const (
	SPECIAL_THREE_SAME_COLOR          = 11 //三同花 4分
	SPECIAL_THREE_STRAIGHT            = 12 //三顺子 4分
	SPECIAL_SIX_PAIR                  = 13 //六队半 4分
	SPECIAL_FOUR_THREE_KIND           = 14 //四套三条 8分
	SPECIAL_THREE_FOUR_KIND           = 15 //三分天下 16分
	SPECIAL_THREE_SAME_COLOR_STRAIGHT = 16 //三同花顺 18分
	SPECIAL_DRAGON                    = 17 //一条龙 26分
	SPECIAL_COLOR_DRAGON              = 18 //青龙 108分
	SPECIAL_FAIL                      = 30 //倒水 相公
)

/////////////////////////////////////////////////////////
//卡牌管理器，负责做牌
type MgrCard struct {
	MVCard       []int
	MVSourceCard []int
}

//初始化
func (this *MgrCard) InitCards() {
	this.MVCard = []int{}
	this.MVSourceCard = []int{}
}

//赋值
func (this *MgrCard) InitNormalCards() {
	begaincard := []int{Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1, Card_Hei_1}
	for _, v := range begaincard {
		for j := int(0); j < 13; j++ {
			this.MVCard = append(this.MVCard, v+j)
		}
	}
}

//洗牌
func (this *MgrCard) Shuffle() {
	this.MVSourceCard = ListShuffle(this.MVCard)

}

//发手牌
func (this *MgrCard) SendHandCard(num int) []int {
	list := append([]int{}, this.MVSourceCard[0:num]...)
	//删除第0个元素开始
	this.MVSourceCard = ListDelMore(this.MVSourceCard, 0, num)
	//返回
	return ListSortDesc(list)
}

//获取卡牌值数组
func GetSmallCards(cards []int) []int {
	var smallCards []int
	for i := 0; i < len(cards); i++ {
		vi := GetCardValue(cards[i])
		smallCards = append(smallCards, vi)
	}
	return smallCards
}

// 获取卡牌花色
func GetCardColor(card int) int {
	return (card & CARD_COLOR) >> 4
}

// 获取卡牌值
func GetCardValue(card int) int {
	tempCard := (card & CARD_VALUE)
	if tempCard == 1 {
		return 14
	}
	return tempCard
}

//判断是否同花
func CheckColor(cards []int) bool {

	var flag bool = true
	for i := 0; i < len(cards)-1; i++ {
		if GetCardColor(cards[i]) != GetCardColor(cards[i+1]) {
			flag = false
		}
	}
	return flag
}

func CheckDragon(cards []int) bool {
	tempCards := ListSortDesc(cards)

	rs := true
	//2-A
	for i := 0; i < len(tempCards)-1; i++ {
		if GetCardValue(tempCards[i]) != GetCardValue(tempCards[i+1])+1 {
			rs = false
			break
		}
	}

	return rs
}

//判断是否是顺子
func CheckStraight(cards []int) bool {
	tempCards := ListSortDesc(cards)

	rs := true
	//2-A
	for i := 0; i < len(tempCards)-1; i++ {
		if GetCardValue(tempCards[i]) != GetCardValue(tempCards[i+1])+1 {
			rs = false
			break
		}
	}

	if rs {
		return rs
	}

	//特殊顺子 1，2，3，4，5
	tempCards2 := []int{14, 2, 3, 4, 5}
	count := 0
	for i := 0; i < len(tempCards2); i++ {
		for j := 0; j < len(tempCards); j++ {
			if tempCards2[i] == GetCardValue(tempCards[j]) {
				count++
				break
			}
		}
	}
	if count == 5 {
		rs = true
	}
	return rs
}

//计算比大小 1 前面大 2 后面大
func MCardsCompare(p1 GCardsType, p2 GCardsType) int {
	//判断类型大小
	if p1.Type > p2.Type {
		return 1
	} else if p1.Type < p2.Type {
		return 2
	} else {
		p1.Cards = ListSortDesc(p1.Cards)
		p2.Cards = ListSortDesc(p2.Cards)
		//类型相等
		if p1.Type == NORMAL_PAIR || p1.Type == NORMAL_TWO_PAIR {
			//对子,两对
			_, cardsLists1 := GetKindLists(p1.Cards, 2)
			_, cardsLists2 := GetKindLists(p2.Cards, 2)
			for i := 0; i < len(cardsLists1); i++ {
				if GetCardValue(cardsLists1[i][0]) > GetCardValue(cardsLists2[i][0]) {
					return 1
				} else if GetCardValue(cardsLists1[i][0]) < GetCardValue(cardsLists2[i][0]) {
					return 2
				}
			}
		} else if p1.Type == NORMAL_FOUR_KIND {
			//铁支
			_, cardsLists1 := GetKindLists(p1.Cards, 4)
			_, cardsLists2 := GetKindLists(p2.Cards, 4)
			for i := 0; i < len(cardsLists1); i++ {
				if GetCardValue(cardsLists1[i][0]) > GetCardValue(cardsLists2[i][0]) {
					return 1
				} else if GetCardValue(cardsLists1[i][0]) < GetCardValue(cardsLists2[i][0]) {
					return 2
				}
			}
		} else if p1.Type == NORMAL_THREE_KIND || p1.Type == NORMAL_GOURD {
			//三条，葫芦
			_, cardsLists1 := GetKindLists(p1.Cards, 3)
			_, cardsLists2 := GetKindLists(p2.Cards, 3)
			for i := 0; i < len(cardsLists1); i++ {
				if GetCardValue(cardsLists1[i][0]) > GetCardValue(cardsLists2[i][0]) {
					return 1
				} else if GetCardValue(cardsLists1[i][0]) < GetCardValue(cardsLists2[i][0]) {
					return 2
				}
			}
		} else if p1.Type == NORMAL_COLOR_STRAIGHT || p1.Type == NORMAL_STRAIGHT {
			//同花顺，顺子
			tempPoker1 := GetSmallCards(p1.Cards)
			tempPoker2 := GetSmallCards(p2.Cards)
			for i := 0; i < len(tempPoker1); i++ {
				if tempPoker1[i] > tempPoker2[i] {
					return 1
				} else if tempPoker1[i] < tempPoker2[i] {
					return 2
				}
			}
		} else {
			//五同，同花，乌龙
			var count1, count2 int
			var cardsLists1, cardsLists2 [][]int
			//判断是否含有3条
			count1, cardsLists1 = GetKindLists(p1.Cards, 3)
			count2, cardsLists2 = GetKindLists(p2.Cards, 3)
			if count1 > count2 {
				return 1
			} else if count1 < count2 {
				return 2
			} else {
				if count1 > 0 && count1 == count2 {
					for i := 0; i < len(cardsLists1); i++ {
						if GetCardValue(cardsLists1[i][0]) > GetCardValue(cardsLists2[i][0]) {
							return 1
						} else if GetCardValue(cardsLists1[i][0]) < GetCardValue(cardsLists2[i][0]) {
							return 2
						}
					}
				} else {
					//判断是否含有对子
					count1, cardsLists1 = GetKindLists(p1.Cards, 2)
					count2, cardsLists2 = GetKindLists(p2.Cards, 2)
					if count1 > count2 {
						return 1
					} else if count1 < count2 {
						return 2
					} else {
						if count1 > 0 && count1 == count2 {
							for i := 0; i < len(cardsLists1); i++ {
								if GetCardValue(cardsLists1[i][0]) > GetCardValue(cardsLists2[i][0]) {
									return 1
								} else if GetCardValue(cardsLists1[i][0]) < GetCardValue(cardsLists2[i][0]) {
									return 2
								}
							}
						}
					}
				}
			}
		}
	}

	//判断单支最大
	for i := 0; i < len(p1.Cards); i++ {
		if GetCardValue(p1.Cards[i]) > GetCardValue(p2.Cards[i]) {
			return 1
		} else if GetCardValue(p1.Cards[i]) < GetCardValue(p2.Cards[i]) {
			return 2
		}
	}
	return 0
}

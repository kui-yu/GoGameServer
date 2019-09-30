package main

const (
	//筒子
	Card_Tong_1 = iota + 1
	Card_Tong_2
	Card_Tong_3
	Card_Tong_4
	Card_Tong_5
	Card_Tong_6
	Card_Tong_7
	Card_Tong_8
	Card_Tong_9
)

const (
	//白板
	Card_White = 10
)

//卡牌管理器，负责做牌
type MgrCard struct {
	MVCard       []int
	MVSourceCard []int
}

//初始化置空卡牌
func (this *MgrCard) InitCards() {
	this.MVCard = []int{}
	this.MVSourceCard = []int{}
}

//赋值
func (this *MgrCard) InitNormalCards() {
	for i := 1; i <= 10; i++ {
		for j := 0; j < 4; j++ {
			this.MVCard = append(this.MVCard, i)
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
	return sortCard(ListSortDesc(list))
}

//判断对子
func checkPairs(cards []int) int {
	if cards[0] == cards[1] {
		return 1
	}
	return 0
}

//排序
func sortCard(cards []int) []int {
	sourceList := append([]int{}, cards...)

	if cards[0] == 10 {
		sourceList[0] = cards[1]
		sourceList[1] = cards[0]
	}
	return sourceList
}

//比牌
func compareCards(cards1 []int, cards2 []int) int {
	//先判断对子
	pair1 := checkPairs(cards1)
	pair2 := checkPairs(cards2)
	if pair1 > pair2 {
		return 1
	} else if pair1 < pair2 {
		return 2
	} else if pair1 == pair2 && pair1 == 1 {
		if cards1[0] > cards2[0] {
			return 1
		} else if cards1[0] < cards2[0] {
			return 2
		} else {
			return 0
		}
	} else {
		//在比较点数
		point1 := (cards1[0] + cards1[1]) % 10
		point2 := (cards2[0] + cards2[1]) % 10
		//二八杠
		if point1 == 0 && point2 > 0 && cards1[0] == 8 {
			return 1
		} else if point1 > 0 && point2 == 0 && cards2[0] == 8 {
			return 2
		} else if point1 == 0 && point2 == 0 && cards1[0] == 8 && cards2[0] != 8 {
			return 1
		} else if point1 == 0 && point2 == 0 && cards1[0] != 8 && cards2[0] == 8 {
			return 2
		} else {

			if point1 > point2 {
				return 1
			} else if point1 < point2 {
				return 2
			} else {
				cards1 = ListSortDesc(cards1)
				cards2 = ListSortDesc(cards2)
				//点数一样
				if cards1[0] > cards2[0] {
					return 1
				} else if cards1[0] < cards2[0] {
					return 2
				}
			}

		}
	}

	return 0
}

package main

import (
	"crypto/rand"
	"fmt"

	"math/big"
	"time"
)

// 玩家历史
type BetHistory struct {
	IsVictory bool  // 是否胜利
	DownBet   int64 // 下注的金额
}

// 座位信息
type SeatInfo struct {
	Id              int             // 座位id
	UserId          int64           // 玩家id
	DownBetValue    int64           // 总下注金额
	UserBetValue    int64           // 用户下注额
	SeatDownCount   int             // 坐下的次数
	MinDownBetCount int             // 最低下注次数
	TrendHistory    []CardGroupType // 走势
}

/////////////////////////////////////////////////////////
// 获得当前时间的毫毛
func GetTimeMS() int64 {
	return time.Now().UnixNano() / 1e6
}

// 随机数生成器
func GetRandomNum(min, max int) (int, error) {
	maxBigInt := big.NewInt(int64(max - min))
	i, err := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < 0 {
		return 0, err
	}
	return int(i.Int64()) + min, err
}

// 洗牌
func (this *ExtDesk) ShuffleCard() {
	var cards []uint8 = []uint8{}
	for i := 1; i < 5; i++ {
		for j := 1; j < 14; j++ {
			cards = append(cards, uint8((i<<4)|j))
		}
	}

	clen := 52
	if gameConfig.GameLimtInfo.ExistsMaxminking {
		cards = append(cards, uint8(Card_King|14))
		cards = append(cards, uint8(Card_King|15))
		clen += 2
	}

	for i := 0; i < clen; i++ {
		r, _ := GetRandomNum(i, clen)
		cards[i], cards[r] = cards[r], cards[i]
	}

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

func DebugLog(format string, args ...interface{}) {
	return
	fmt.Println(fmt.Sprintf(format, args))
}

func TestLog(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args))
}

func FormatDeskId(deskId int, grade int) string {
	first := ""
	switch grade {
	case 1:
		first = "C"
	case 2:
		first = "Z"
	case 3:
		first = "G"
	}

	return fmt.Sprintf("%s%04d", first, deskId+1)
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

// 牌排序
func OrderCardU8(cards []uint8) []uint8 {
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

// 名字***
func MarkNickName(name string) string {
	if len(name) > 4 {
		return "***" + name[len(name)-4:]
	} else {
		return "***" + name
	}
}

// 生成庄家牌
func MarkBankerCard(cards []uint8, oneCard uint8, carType CardGroupType) []int {
	clen := len(cards)
	for i := 0; i < clen; i++ {
		r, _ := GetRandomNum(i, clen)
		cards[i], cards[r] = cards[r], cards[i]
	}

	bankerCards := []int{}
	isFind := false
	for i := 1; i < clen; i++ {
		if isFind {
			break
		}
		for j := i + 1; j < clen; j++ {
			if isFind {
				break
			}
			for k := j + 1; k < clen; k++ {
				if isFind {
					break
				}
				for l := k + 1; l < clen; l++ {
					bankerCards = []int{
						int(oneCard),
						int(cards[i]),
						int(cards[j]),
						int(cards[k]),
						int(cards[l]),
					}

					ctype, _, _ := CalcCards(bankerCards)
					if ctype == carType {
						isFind = true
						break
					}
				}
			}
		}
	}

	if isFind {
		return bankerCards
	} else {
		for i := 0; i < 4; i++ {
			bankerCards[i] = int(cards[i])
		}
		bankerCards[4] = int(oneCard)
		return bankerCards
	}
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Abs(x int) int {
	if x > 0 {
		return x
	}
	return x * -1
}

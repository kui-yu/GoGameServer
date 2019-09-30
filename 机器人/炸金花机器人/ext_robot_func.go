package main

const (
	CARD_COLOR = 0xF0 //花色掩码
	CARD_VALUE = 0x0F //数值掩码
)

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

func DelRepeat(ary []int32) []int32 {
	var rary []int32
	if len(ary) == 0 {
		return nil
	}
	rary = append(rary, ary[0])
	for i := 1; i < len(ary); i++ {
		for j := 0; j < len(rary); j++ {
			if ary[i] == rary[j] {
				break
			}
			if j == len(rary)-1 {
				rary = append(rary, ary[i])
			}
		}
	}
	return rary
}

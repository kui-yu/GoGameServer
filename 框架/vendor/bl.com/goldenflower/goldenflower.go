package goldenflower

const (
	GFError   byte = iota // 类型错误
	GFSingle              // 单牌类型
	GFDouble              // 对子类型
	GFShunZi              // 顺子类型
	GFJinHua              // 金花类型
	GFShunJin             // 顺金类型
	GFBaoZi               // 豹子类型
	GFSpecial             // 特殊类型235
)

const (
	Card_A byte = iota + 1
	Card_2
	Card_3
	Card_4
	Card_5
	Card_6
	Card_7
	Card_8
	Card_9
	Card_10
	Card_J
	Card_Q
	Card_K
)

// 花色获取函数
func GetCardColor(card byte) byte {
	return (card & 0xF0) >> 4
}

// 牌值获取函数
func GetCardValue(card byte) byte {
	return (card & 0x0F)
}

// 逻辑牌值获取函数
func GetLogicValue(card byte) byte {
	d := GetCardValue(card)
	if d <= 1 {
		return d + 13
	}
	return d
}

// 牌型排序
func Sort(v []byte) []byte {
	cs := v
	length := len(cs)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			vi := GetLogicValue(cs[i])
			vj := GetLogicValue(cs[j])
			if vi < vj || ((vi == vj) && (GetCardColor(cs[i]) < GetCardColor(cs[j]))) {
				cs[i], cs[j] = cs[j], cs[i]
			}
		}
	}
	return append([]byte{}, cs...)
}

// 牌型获取
func GetCardsType(cards []byte, spe bool) byte {
	if len(cards) != 3 {
		return GFError
	}

	// 牌排序
	cards = Sort(cards)

	suit := GetCardColor(cards[0])

	c1 := GetCardValue(cards[0])
	c2 := GetCardValue(cards[1])
	c3 := GetCardValue(cards[2])

	sameColor := true // 同色
	lineCard := true  // 顺子

	// 牌型分析
	for i := 1; i < len(cards); i++ {
		// 数据分析：金花、顺子
		if GetCardColor(cards[i]) != suit {
			sameColor = false
		}

		if c1-GetCardValue(cards[i]) != byte(i) {
			lineCard = false
		}

		// 判断结束
		if sameColor == false && lineCard == false {
			break
		}
	}

	// 特殊牌型 A32
	if !lineCard {
		if c1 == Card_A && c2 == Card_3 && c3 == Card_2 {
			lineCard = true
		} else if c1 == Card_A && c2 == Card_K && c3 == Card_Q {
			lineCard = true
		}
	}

	// 顺金
	if sameColor && lineCard {
		return GFShunJin
	}

	// 顺子
	if lineCard {
		return GFShunZi
	}

	// 金花
	if sameColor {
		return GFJinHua
	}

	// 豹子
	if c1 == c2 && c2 == c3 {
		return GFBaoZi
	}

	// 对子
	if c1 == c2 || c2 == c3 {
		return GFDouble
	}

	// 特殊牌型 235
	if spe && c3 == 2 && c2 == 3 && c1 == 5 {
		return GFSpecial
	}

	return GFSingle
}

// 比牌
func CompareCard(f []byte, n []byte, spe bool) bool {
	if len(f) != 3 || len(n) != 3 {
		return false
	}

	fV := GetCardsType(f, spe)
	nV := GetCardsType(n, spe)

	for {
		// 特殊牌型235 > 豹子
		if spe && fV+nV == GFSpecial+GFBaoZi {
			break
		}

		// 特殊牌型235 还原为单牌
		if spe && nV == GFSpecial {
			nV = GFSingle
		}
		if spe && fV == GFSpecial {
			fV = GFSingle
		}

		if fV != nV {
			break
		}

		// 同类型判断
		switch fV {
		case GFBaoZi, GFSingle, GFJinHua: // 豹子\单牌\金花
			// 比牌值
			var isBreak bool = false
			for i := 0; i < 3; i++ {
				fV = GetLogicValue(f[i])
				nV = GetLogicValue(n[i])
				if fV != nV {
					isBreak = true
					break
				}
			}
			if isBreak {
				break
			}

			// 比花色
			fV = f[0]
			nV = n[0]

			break
		case GFShunJin, GFShunZi: // 432 > A32
			fV = GetLogicValue(f[0])
			nV = GetLogicValue(n[0])

			if fV == 14 && GetCardValue(f[2]) == 2 {
				fV = 1
			}

			if nV == 14 && GetCardValue(n[2]) == 2 {
				nV = 1
			}
			if fV != nV {
				break
			}

			// 顺子大牌一致，比花色
			fV = f[0]
			nV = n[0]
			break
		case GFDouble: // 对子
			fV = GetLogicValue(f[1])
			nV = GetLogicValue(n[1])

			if fV != nV {
				break
			}

			// 对子一样，取单牌大小
			if GetLogicValue(f[0]) == fV {
				fV = GetLogicValue(f[2])
			} else {
				fV = GetLogicValue(f[0])
			}

			if GetLogicValue(n[0]) == nV {
				nV = GetLogicValue(n[2])
			} else {
				nV = GetLogicValue(n[0])
			}

			if fV != nV {
				break
			}

			// 对子一样,单牌一样 取对子牌花色
			if GetLogicValue(f[0]) == fV {
				fV = f[0]
			} else {
				fV = f[1]
			}

			if GetLogicValue(n[0]) == nV {
				nV = n[0]
			} else {
				nV = n[1]
			}
			break
		}

		break
	}

	return fV > nV
}

/**
* 公共方法
**/
package main

func FindIdxFromInts(s []int, val int) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}

	return -1
}

// 牌排序
func OrderCard(cards []int, orderAsc bool) []int {
	tempCards := cards
	clen := len(tempCards)
	for i := 0; i < clen; i++ {
		for j := i + 1; j < clen; j++ {
			a, b := tempCards[i], tempCards[j]
			if orderAsc {
				if a&0xF > b&0xF {
					tempCards[i], tempCards[j] = b, a
				} else if (a&0xF == b&0xF) && ((a >> 4) > (b >> 4)) {
					tempCards[i], tempCards[j] = b, a
				}
			} else {
				if a&0xF < b&0xF {
					tempCards[i], tempCards[j] = b, a
				} else if (a&0xF == b&0xF) && ((a >> 4) < (b >> 4)) {
					tempCards[i], tempCards[j] = b, a
				}
			}
		}
	}
	return tempCards
}

func CalcCard(cards []int) GCardGroupInfo {
	groupInfo := GCardGroupInfo{}
	// 类型对应cardVal
	typeAtVals := make(map[int][]int)
	// 值对应cardType
	valAtTypes := make(map[int][]int)

	//倒序排序
	ocards := OrderCard(cards, false)

	for _, card := range ocards {
		ctype := card >> 4
		cnum := card & 0xF

		tval, tok := typeAtVals[ctype]
		if tok {
			typeAtVals[ctype] = append(tval, cnum)
		} else {
			typeAtVals[ctype] = []int{cnum}
		}

		vval, vok := valAtTypes[cnum]
		if vok {
			valAtTypes[cnum] = append(vval, ctype)
		} else {
			valAtTypes[cnum] = []int{ctype}
		}
	}

	// 皇家同花顺
	{
		hjCards := []int{
			Card_Hei | 13,
			Card_Hei | 12,
			Card_Hei | 11,
			Card_Hei | 10,
			Card_Hei | 1,
		}

		isExists := true
		for _, v := range hjCards {
			if FindIdxFromInts(ocards, v) == -1 {
				isExists = false
				break
			}
		}

		if isExists == true {
			groupInfo.GroupType = CardGroupRoyalFlush
			groupInfo.Cards = hjCards
			return groupInfo
		}
	}

	//同花顺
	{
		isExists := false
		for k, vals := range typeAtVals {
			if len(vals) < 5 {
				continue
			}
			isExists = false
			oneIdx, curVal := 0, 0
			for i, val := range vals {
				if i == 0 {
					oneIdx, curVal = 0, val
				}
				if curVal-1 != val {
					oneIdx, curVal = i, val
				}

				if i-oneIdx == 4 {
					isExists = true
					break
				}
			}
			if isExists == true {
				groupInfo.GroupType = CardGroupSFlush
				for i := oneIdx; i < oneIdx+5; i++ {
					groupInfo.Cards = append(groupInfo.Cards, k<<4|vals[i])
				}

				return groupInfo
			}
		}
	}

	// 四条
	{
		scard := -1
		for v, types := range valAtTypes {
			if len(types) == 4 {
				scard = v
				break
			}
		}

		if scard != -1 {
			groupInfo.GroupType = CardGroupFourT
			groupInfo.Cards = []int{
				Card_Hei | scard,
				Card_Hong | scard,
				Card_Mei | scard,
				Card_Fang | scard,
			}
			for _, card := range ocards {
				if card&0xF != scard {
					groupInfo.Cards = append(groupInfo.Cards, card)
					return groupInfo
				}
			}
		}
	}

	//三张加一对
	{
		isExists, val := false, -1
		for v, types := range valAtTypes {
			if len(types) == 3 {
				isExists, val = true, v
				break
			}
		}

		if isExists {
			for v, types := range valAtTypes {
				if v != val && len(types) > 1 {
					groupInfo.GroupType = CardGroupFullhouse
					st, _ := valAtTypes[val]
					for _, t := range st {
						groupInfo.Cards = append(groupInfo.Cards, t<<4|val)
					}
					groupInfo.Cards = append(groupInfo.Cards, types[0]<<4|v)
					groupInfo.Cards = append(groupInfo.Cards, types[1]<<4|v)
					return groupInfo
				}
			}
		}
	}

	//同花
	{
		for t, vals := range typeAtVals {
			if len(vals) < 5 {
				continue
			}

			groupInfo.GroupType = CardGroupFlush

			for i := 0; i < 5; i++ {
				groupInfo.Cards = append(groupInfo.Cards, t<<4|vals[i])
			}
			return groupInfo
		}
	}

	//顺子
	{
		ns := []int{}
		for _, v := range ocards {
			if len(ns) == 0 {
				ns = append(ns, v)
			}
			lc := ns[len(ns)-1]

			if (v & 0xF) == (lc & 0xF) {
				continue
			} else if (v & 0xF) == ((lc & 0xF) - 1) {
				ns = append(ns, v)
			} else {
				ns = []int{v}
			}
		}

		if len(ns) > 4 {
			groupInfo.GroupType = CardGroupStraight
			groupInfo.Cards = ns[:5]
		}
	}

	//三张
	{
		isExists, val := false, -1
		for v, types := range valAtTypes {
			if len(types) == 3 {
				isExists, val = true, v
				break
			}
		}

		if isExists {
			groupInfo.GroupType = CardGroupFullhouse
			st, _ := valAtTypes[val]
			for _, t := range st {
				groupInfo.Cards = append(groupInfo.Cards, t<<4|val)
			}

			for _, c := range ocards {
				if (c & 0xF) != val {
					groupInfo.Cards = append(groupInfo.Cards, c)
					if len(groupInfo.Cards) == 5 {
						return groupInfo
					}
				}
			}
		}
	}

	// 两对 | 一对 | 高牌
	{
		nc := []int{}
		for t, vals := range typeAtVals {
			if len(vals) == 1 {
				continue
			}

			nc = append(nc, t<<4|vals[0])
			nc = append(nc, t<<4|vals[1])
			if len(nc) == 4 {
				break
			}
		}

		if len(nc) == 2 {
			groupInfo.GroupType = CardGroupOnePair
		} else if len(nc) == 4 {
			groupInfo.GroupType = CardGroupTwoPair
		} else {
			groupInfo.GroupType = CardGroupHighCard
		}
		for _, v := range ocards {
			if FindIdxFromInts(nc, v) == -1 {
				nc = append(nc, v)
				groupInfo.Cards = nc
				if len(nc) == 5 {
					return groupInfo
				}
			}
		}
	}

	return groupInfo
}

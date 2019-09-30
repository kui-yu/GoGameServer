/**
* 公共函数
**/

package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"logs"
	"math/big"
	"strconv"
	"time"
)

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
func ShuffleCard(isAddKing bool) []int {
	var cards []int = []int{}
	for i := 1; i < 5; i++ {
		for j := 1; j < 14; j++ {
			cards = append(cards, (i<<4)|j)
		}
	}

	if isAddKing {
		cards = append(cards, Card_King|14)
		cards = append(cards, Card_King|15)
	}
	clen := len(cards)
	for i := 0; i < clen; i++ {
		r, _ := GetRandomNum(i, clen)
		cards[i], cards[r] = cards[r], cards[i]
	}
	return cards
}

func DebugLog(format string, args ...interface{}) {
	logs.Debug(format, args)
}

func ErrorLog(format string, args ...interface{}) {
	logs.Error(format, args)
	logs.Debug(format, args)
}

func TestLog(format string, args ...interface{}) {
	logs.Error(format, args)
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

// 名字***
func MarkNickName(name string) string {
	if len(name) > 4 {
		return "***" + name[len(name)-4:]
	} else {
		return "***" + name
	}
}

func FindIndexFromInt64(arr []int64, val int64) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

//查找牌
func FindCardsFromNum(arr []int, num int) []int {
	r := []int{}
	for _, v := range arr {
		if v&0xF == num {
			r = append(r, v)
		}
	}
	return r
}

func ReplaceMapField(obj interface{}, paths []string, callback func(val interface{}) interface{}) {
	isLast := len(paths) == 1
	switch t := obj.(type) {
	case map[string]interface{}:
		if isLast {
			t[paths[0]] = callback(t[paths[0]])
		} else {
			if paths[0] == "*" {
				for _, nobj := range t {
					ReplaceMapField(nobj, paths[1:], callback)
				}
			} else {
				ReplaceMapField(t[paths[0]], paths[1:], callback)
			}
		}

	case []interface{}:
		if isLast {
			idx, _ := strconv.Atoi(paths[0])
			t[idx] = callback(t[idx])
		} else {
			if paths[0] == "*" {
				for _, nobj := range t {
					ReplaceMapField(nobj, paths[1:], callback)
				}
			} else {
				idx, _ := strconv.Atoi(paths[0])
				ReplaceMapField(t[idx], paths[1:], callback)
			}
		}
	default:

	}
}

func ConvertObjToMap(obj interface{}) map[string]interface{} {
	mdata := map[string]interface{}{}
	jdata, _ := json.Marshal(obj)
	json.Unmarshal(jdata, &mdata)
	return mdata
}

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
	typeKeys := []int{}
	// 值对应cardType
	valAtTypes := make(map[int][]int)
	valKeys := []int{}

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
			typeKeys = append(typeKeys, ctype)
		}

		vval, vok := valAtTypes[cnum]
		if vok {
			valAtTypes[cnum] = append(vval, ctype)
		} else {
			valAtTypes[cnum] = []int{ctype}
			valKeys = append(valKeys, cnum)
		}
	}

	// 皇家同花顺
	{
		hjCards := []int{
			Card_Hei | 1,
			Card_Hei | 13,
			Card_Hei | 12,
			Card_Hei | 11,
			Card_Hei | 10,
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
		//考虑AKQJT
		ns := []int{}
		dsCards := []int{}

		for k, vals := range typeAtVals {
			if len(vals) < 5 {
				continue
			}
			ns = []int{}
			dsCards = []int{
				1,
				13,
				12,
				11,
				10,
			}

			for {
				num := dsCards[0]
				if FindIdxFromInts(vals, num) == -1 {
					break
				} else {
					dsCards = dsCards[1:]
					ns = append(ns, (k<<4)|num)
				}

				if len(dsCards) == 0 {
					break
				}
			}
			if len(ns) == 5 {
				break
			}
		}

		if len(ns) == 5 {
			groupInfo.GroupType = CardGroupSFlush
			groupInfo.Cards = ns
			return groupInfo
		}

		isExists := false

		for k, vals := range typeAtVals {
			if len(vals) < 5 {
				continue
			}
			isExists = false

			oneIdx, upVal := 0, 0
			for i, val := range vals {
				if i == 0 || (upVal-1 != val) {
					oneIdx, upVal = i, -1
				}

				if i-oneIdx == 4 {
					isExists = true
					break
				}
				upVal = val

			}
			if isExists == true {
				groupInfo.GroupType = CardGroupSFlush
				for i := oneIdx; i < oneIdx+5; i++ {
					groupInfo.Cards = append(groupInfo.Cards, (k<<4)|vals[i])
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
						groupInfo.Cards = append(groupInfo.Cards, (t<<4)|val)
					}
					groupInfo.Cards = append(groupInfo.Cards, (types[0])<<4|v)
					groupInfo.Cards = append(groupInfo.Cards, (types[1])<<4|v)
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

			alen := 5
			if vals[len(vals)-1] == 1 {
				alen = 4
				groupInfo.Cards = append(groupInfo.Cards, (t<<4)|1)
			}

			for i := 0; i < alen; i++ {
				groupInfo.Cards = append(groupInfo.Cards, (t<<4)|vals[i])
			}
			return groupInfo
		}
	}

	//顺子
	{
		ns := []int{}

		//考虑AKQJT
		dsCards := []int{
			1,
			13,
			12,
			11,
			10,
		}
		for _, num := range dsCards {
			isExists := false
			for _, v := range ocards {
				if v&0xF == num {
					isExists = true
					ns = append(ns, v)
					break
				}
			}
			if !isExists {
				break
			}
		}
		if len(ns) == 5 {
			groupInfo.GroupType = CardGroupStraight
			groupInfo.Cards = ns
			return groupInfo
		}

		ns = []int{}
		// 其他顺子
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
			return groupInfo
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
			groupInfo.GroupType = CardGroupThreeT
			st, _ := valAtTypes[val]
			for _, t := range st {
				groupInfo.Cards = append(groupInfo.Cards, (t<<4)|val)
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
		// 是否一对A
		for _, val := range valKeys {
			types := valAtTypes[val]
			if len(types) == 1 {
				continue
			}
			if val == 1 {
				nc = append(nc, (types[0]<<4)|val)
				nc = append(nc, (types[1]<<4)|val)
			}
		}

		for _, val := range valKeys {
			types := valAtTypes[val]
			if len(types) == 1 || val == 1 {
				continue
			}
			nc = append(nc, (types[0]<<4)|val)
			nc = append(nc, (types[1]<<4)|val)
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

		// 加入A
		for _, v := range ocards {
			if FindIdxFromInts(nc, v) == -1 {
				if v&0xF == 1 {
					nc = append(nc, v)
					groupInfo.Cards = nc
				}

				if len(nc) == 5 {
					return groupInfo
				}
			}
		}

		for _, v := range ocards {
			if FindIdxFromInts(nc, v) == -1 && v&0xF != 1 {
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

// a<b return -1 a==b return 0 a>b return 1
func CompareCardGroup(a GCardGroupInfo, b GCardGroupInfo) int {
	if a.GroupType > b.GroupType {
		return 1
	} else if a.GroupType < b.GroupType {
		return -1
	} else {
		for i := 0; i < 5; i++ {
			c1, c2 := a.Cards[i]&0xF, b.Cards[i]&0xF

			// 处理A
			if c1&0xF == 1 && c2&0xF != 1 {
				return 1
			}
			if c2&0xF == 1 && c1&0xF != 1 {
				return -1
			}

			// 处理普通
			if c1 == c2 {
				continue
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		}
	}
	return 0
}

func TestFormatCard(cards []int) string {
	str := "["
	for _, card := range cards {
		str += fmt.Sprintf("%d,", card&0xF)
	}
	str += "]"
	return str
}

// 牌型和牌型比较 测试数据
func Test() {
	cards := [][]int{}
	//0皇家同花顺
	cards = append(cards, []int{4<<4 | 1, 4<<4 | 13, 4<<4 | 12, 4<<4 | 11, 4<<4 | 10, 1<<4 | 4, 2<<4 | 12})
	//1最大同花顺
	cards = append(cards, []int{2<<4 | 1, 2<<4 | 13, 2<<4 | 12, 2<<4 | 11, 2<<4 | 10, 2<<4 | 4, 3<<4 | 12})
	//2同花顺12345
	cards = append(cards, []int{2<<4 | 1, 2<<4 | 2, 2<<4 | 3, 2<<4 | 4, 2<<4 | 5, 1<<4 | 4, 1<<4 | 12})
	//3同花顺23456
	cards = append(cards, []int{2<<4 | 2, 2<<4 | 3, 2<<4 | 4, 2<<4 | 5, 2<<4 | 6, 1<<4 | 4, 1<<4 | 12})
	//4四条2
	cards = append(cards, []int{1<<4 | 2, 2<<4 | 2, 3<<4 | 2, 4<<4 | 2, 2<<4 | 1, 1<<4 | 4, 1<<4 | 12})
	//5四条A
	cards = append(cards, []int{1<<4 | 1, 2<<4 | 2, 3<<4 | 1, 4<<4 | 1, 2<<4 | 1, 1<<4 | 4, 1<<4 | 12})
	//6三条A+一对2
	cards = append(cards, []int{1<<4 | 1, 2<<4 | 1, 3<<4 | 1, 4<<4 | 2, 2<<4 | 2, 1<<4 | 8, 1<<4 | 12})
	//7三条2+一对3
	cards = append(cards, []int{1<<4 | 2, 2<<4 | 2, 3<<4 | 2, 4<<4 | 3, 2<<4 | 3, 1<<4 | 8, 1<<4 | 12})
	//8三条2+一对1
	cards = append(cards, []int{1<<4 | 2, 2<<4 | 2, 3<<4 | 2, 4<<4 | 1, 2<<4 | 1, 1<<4 | 8, 1<<4 | 12})
	//9同花最大K
	cards = append(cards, []int{1<<4 | 13, 1<<4 | 11, 1<<4 | 9, 1<<4 | 8, 1<<4 | 6, 2<<4 | 8, 1<<4 | 12})
	//10同花最大A
	cards = append(cards, []int{1<<4 | 1, 1<<4 | 11, 1<<4 | 9, 1<<4 | 8, 1<<4 | 6, 2<<4 | 8, 1<<4 | 12})
	//11顺子12345
	cards = append(cards, []int{3<<4 | 1, 1<<4 | 2, 2<<4 | 3, 2<<4 | 4, 2<<4 | 5, 1<<4 | 4, 1<<4 | 12})
	//12顺子23456
	cards = append(cards, []int{2<<4 | 2, 2<<4 | 3, 1<<4 | 4, 2<<4 | 5, 2<<4 | 6, 1<<4 | 4, 1<<4 | 12})
	//13顺子AKQJT
	cards = append(cards, []int{2<<4 | 1, 3<<4 | 13, 1<<4 | 12, 2<<4 | 11, 2<<4 | 10, 1<<4 | 4, 1<<4 | 12})
	//14三条A+最大K
	cards = append(cards, []int{1<<4 | 1, 2<<4 | 1, 3<<4 | 1, 4<<4 | 2, 2<<4 | 4, 1<<4 | 8, 1<<4 | 13})
	//15三条2+最大A
	cards = append(cards, []int{1<<4 | 2, 2<<4 | 2, 3<<4 | 2, 4<<4 | 1, 2<<4 | 4, 1<<4 | 8, 1<<4 | 13})
	//16三条2+最大10
	cards = append(cards, []int{1<<4 | 2, 2<<4 | 2, 3<<4 | 2, 4<<4 | 10, 2<<4 | 6, 1<<4 | 8, 1<<4 | 13})
	//17两对1,2 + K
	cards = append(cards, []int{1<<4 | 1, 2<<4 | 1, 3<<4 | 2, 4<<4 | 2, 2<<4 | 6, 1<<4 | 8, 1<<4 | 13})
	//18两对2,3 + A
	cards = append(cards, []int{1<<4 | 2, 2<<4 | 2, 3<<4 | 3, 4<<4 | 3, 2<<4 | 1, 1<<4 | 8, 1<<4 | 13})
	//19两对2,3 + K
	cards = append(cards, []int{1<<4 | 2, 2<<4 | 2, 3<<4 | 3, 4<<4 | 3, 2<<4 | 4, 1<<4 | 8, 1<<4 | 13})
	//20一对1 + k
	cards = append(cards, []int{1<<4 | 1, 2<<4 | 1, 3<<4 | 3, 4<<4 | 6, 2<<4 | 9, 1<<4 | 8, 1<<4 | 13})
	//21一对2 + k
	cards = append(cards, []int{1<<4 | 2, 2<<4 | 2, 3<<4 | 3, 4<<4 | 6, 2<<4 | 9, 1<<4 | 8, 1<<4 | 13})
	//22一对2 + 1
	cards = append(cards, []int{1<<4 | 2, 2<<4 | 2, 3<<4 | 3, 4<<4 | 6, 2<<4 | 9, 1<<4 | 8, 1<<4 | 1})
	//23高牌1
	cards = append(cards, []int{1<<4 | 1, 2<<4 | 2, 3<<4 | 4, 4<<4 | 5, 2<<4 | 6, 1<<4 | 8, 1<<4 | 11})
	//24高牌k
	cards = append(cards, []int{1<<4 | 13, 2<<4 | 2, 3<<4 | 4, 4<<4 | 5, 2<<4 | 6, 1<<4 | 8, 1<<4 | 11})

	for _, cs := range cards {
		rlen := len(cs)
		for i := 0; i < rlen; i++ {
			ridx, _ := GetRandomNum(0, rlen)
			cs[i], cs[ridx] = cs[ridx], cs[i]
		}
	}
	DebugLog("皇家同花顺<=>最大同花顺")
	ca := CalcCard(cards[0])
	cb := CalcCard(cards[1])
	ta := TestFormatCard(ca.Cards)
	tb := TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("同花顺12345<=>同花顺23456")
	ca = CalcCard(cards[2])
	cb = CalcCard(cards[3])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("同花顺12345<=>最大同花顺")
	ca = CalcCard(cards[2])
	cb = CalcCard(cards[1])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("四条2<=>四条A")
	ca = CalcCard(cards[4])
	cb = CalcCard(cards[5])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("6三条A+一对2<=>7三条2+一对3")
	ca = CalcCard(cards[6])
	cb = CalcCard(cards[7])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("8三条2+一对1<=>7三条2+一对3")
	ca = CalcCard(cards[8])
	cb = CalcCard(cards[7])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("9同花最大K<=>10同花最大A")
	ca = CalcCard(cards[9])
	cb = CalcCard(cards[10])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("11顺子12345<=>12顺子23456")
	ca = CalcCard(cards[11])
	cb = CalcCard(cards[12])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("11顺子12345<=>13顺子AKQJT")
	ca = CalcCard(cards[11])
	cb = CalcCard(cards[13])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("14三条A+最大K<=>15三条2+最大A")
	ca = CalcCard(cards[14])
	cb = CalcCard(cards[15])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("16三条2+最大10<=>15三条2+最大A")
	ca = CalcCard(cards[16])
	cb = CalcCard(cards[15])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("17两对1,2 + K<=>18两对2,3 + A")
	ca = CalcCard(cards[17])
	cb = CalcCard(cards[18])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("19两对2,3 + K<=>18两对2,3 + A")
	ca = CalcCard(cards[19])
	cb = CalcCard(cards[18])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("20一对1 + k<=>21一对2 + k")
	ca = CalcCard(cards[20])
	cb = CalcCard(cards[21])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("22一对2 + 1<=>21一对2 + k")
	ca = CalcCard(cards[22])
	cb = CalcCard(cards[21])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))

	DebugLog("23高牌1<=>24高牌k")
	ca = CalcCard(cards[23])
	cb = CalcCard(cards[24])
	ta = TestFormatCard(ca.Cards)
	tb = TestFormatCard(cb.Cards)
	DebugLog("比较结果", ca.GroupType, cb.GroupType, ta, tb, CompareCardGroup(ca, cb))
}

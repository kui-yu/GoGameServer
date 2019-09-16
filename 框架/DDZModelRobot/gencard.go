package main

import (
	// "logs"
	"math/rand"
)

var G_OutMap map[int32]string

func init() {
	G_OutMap = make(map[int32]string)
	G_OutMap[CT_SINGLE] = "单牌"
	G_OutMap[CT_DOUBLE] = "对子"
	G_OutMap[CT_SINGLE_CONNECT] = "单顺"
	G_OutMap[CT_DOUBLE_CONNECT] = "双顺"
	G_OutMap[CT_THREE] = "三代"
	G_OutMap[CT_THREE_LINE_TAKE_ONE] = "三代一"
	G_OutMap[CT_THREE_LINE_TAKE_TWO] = "三代二"
	G_OutMap[CT_FOUR_LINE_TAKE_ONE] = "四代单"
	G_OutMap[CT_FOUR_LINE_TAKE_TWO] = "四代双"
	G_OutMap[CT_AIRCRAFT] = "飞机"
	G_OutMap[CT_AIRCRAFT_ONE] = "飞机代单"
	G_OutMap[CT_AIRCRAFT_TWO] = "飞机代对"
	G_OutMap[CT_BOMB_FOUR] = "炸弹"
	G_OutMap[CT_TWOKING] = "王炸"
}

func Sort(cs []byte) []byte {
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
	return append([]byte{}, cs...)
}

func VecDelMulti(d []byte, cs []byte) ([]byte, bool) {
	for _, v := range cs {
		isDel := false
		for i := len(d) - 1; i >= 0; i-- {
			if d[i] == v {
				d = append(d[:i], d[i+1:]...)
				isDel = true
				break
			}
		}
		//
		if !isDel {
			return d, isDel
		}
	}
	//
	return d, true
}

func GetCardColor(card byte) byte {
	return (card & CARD_COLOR) >> 4
}

func GetCardValue(card byte) byte {
	return (card & CARD_VALUE)
}

func GetLogicValue(card byte) byte {
	d := GetCardValue(card)
	if card == 0x41 {
		return 16
	}
	if card == 0x42 {
		return 17
	}
	if d <= 2 {
		return d + 13
	}
	return d
}

//
func GenOutCard(cards []byte) {

}

func Gen1OutCard(cards []byte, out *GGameOutCard) bool {
	out.Max = cards[0]
	out.Cards = append([]byte{}, cards...)
	out.Type = CT_SINGLE
	return true
}

func Gen2OutCard(cards []byte, out *GGameOutCard) bool {
	c := Sort(cards)
	if c[0] == 0x42 && c[1] == 0x41 {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_TWOKING
	} else {
		if GetLogicValue(c[0]) != GetLogicValue(c[1]) {
			return false
		}
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_DOUBLE
	}
	return true
}

func Gen3OutCard(cards []byte, out *GGameOutCard) bool {
	c := Sort(cards)
	if GetLogicValue(c[0]) != GetLogicValue(c[2]) {
		return false
	}
	out.Max = c[0]
	out.Cards = append([]byte{}, c...)
	out.Type = CT_THREE
	return true
}

//炸，三带一
func Gen4OutCard(cards []byte, out *GGameOutCard) bool {
	c := Sort(cards)
	if GetLogicValue(c[0]) == GetLogicValue(c[3]) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_BOMB_FOUR
		return true
	}
	//
	if GetLogicValue(c[0]) == GetLogicValue(c[2]) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_THREE_LINE_TAKE_ONE
		return true
	} else if GetLogicValue(c[1]) == GetLogicValue(c[3]) {
		out.Max = c[1]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_THREE_LINE_TAKE_ONE
		return true
	}
	return false
}

//三带二，单顺
func Gen5OutCard(cards []byte, out *GGameOutCard) bool {
	c := Sort(cards)
	if (GetLogicValue(c[0]) == GetLogicValue(c[2])) &&
		(GetLogicValue(c[3]) == GetLogicValue(c[4])) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_THREE_LINE_TAKE_TWO
		return true
	}
	if (GetLogicValue(c[0]) == GetLogicValue(c[1])) &&
		(GetLogicValue(c[2]) == GetLogicValue(c[4])) {
		out.Max = c[2]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_THREE_LINE_TAKE_TWO
		return true
	}

	if GenDanShun(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_SINGLE_CONNECT
		return true
	}

	return false
}

//双顺，单顺，飞机，四带单
func Gen6OutCard(cards []byte, out *GGameOutCard) bool {
	c := Sort(cards)
	if GenDanShun(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_SINGLE_CONNECT
		return true
	}
	if GenShuangShun(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_DOUBLE_CONNECT
		return true
	}
	if GenFeiJi(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_AIRCRAFT
		return true
	}
	if (GetLogicValue(c[0]) == GetLogicValue(c[3])) &&
		(GetLogicValue(c[4]) != GetLogicValue(c[5])) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_FOUR_LINE_TAKE_ONE
		return true
	}
	if (GetLogicValue(c[1]) == GetLogicValue(c[4])) &&
		(GetLogicValue(c[0]) != GetLogicValue(c[5])) {
		out.Max = c[1]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_FOUR_LINE_TAKE_ONE
		return true
	}
	if (GetLogicValue(c[2]) == GetLogicValue(c[5])) &&
		(GetLogicValue(c[0]) != GetLogicValue(c[1])) {
		out.Max = c[2]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_FOUR_LINE_TAKE_ONE
		return true
	}
	return false
}

//单顺
func Gen7OutCard(cards []byte, out *GGameOutCard) bool {
	c := Sort(cards)
	if GenDanShun(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_SINGLE_CONNECT
		return true
	}

	return false
}

//单顺，双顺，飞机带单,四带对
func Gen8OutCard(cards []byte, out *GGameOutCard) bool {
	c := Sort(cards)
	if GenDanShun(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_SINGLE_CONNECT
		return true
	}
	if GenShuangShun(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_DOUBLE_CONNECT
		return true
	}
	if GenFeiJiDaiDan(c, out) {
		out.Cards = append([]byte{}, c...)
		out.Type = CT_AIRCRAFT_ONE
		return true
	}

	if (GetLogicValue(c[0]) == GetLogicValue(c[3])) &&
		(GetLogicValue(c[4]) == GetLogicValue(c[5])) &&
		(GetLogicValue(c[6]) == GetLogicValue(c[7])) &&
		(GetLogicValue(c[4]) != GetLogicValue(c[6])) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_FOUR_LINE_TAKE_TWO
		return true
	}

	if (GetLogicValue(c[2]) == GetLogicValue(c[5])) &&
		(GetLogicValue(c[0]) == GetLogicValue(c[1])) &&
		(GetLogicValue(c[6]) == GetLogicValue(c[7])) &&
		(GetLogicValue(c[0]) != GetLogicValue(c[6])) {
		out.Max = c[2]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_FOUR_LINE_TAKE_TWO
		return true
	}

	if (GetLogicValue(c[4]) == GetLogicValue(c[7])) &&
		(GetLogicValue(c[0]) == GetLogicValue(c[1])) &&
		(GetLogicValue(c[2]) == GetLogicValue(c[3])) &&
		(GetLogicValue(c[0]) != GetLogicValue(c[2])) {
		out.Max = c[4]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_FOUR_LINE_TAKE_TWO
		return true
	}

	return false
}

func GenDefault(cards []byte, out *GGameOutCard) bool {
	c := Sort(cards)
	cl := len(c)
	//
	if cl <= 12 && GenDanShun(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_SINGLE_CONNECT
		return true
	}

	if cl%2 == 0 && GenShuangShun(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_DOUBLE_CONNECT
		return true
	}

	if cl%3 == 0 && GenFeiJi(c, out) {
		out.Max = c[0]
		out.Cards = append([]byte{}, c...)
		out.Type = CT_AIRCRAFT
		return true
	}

	if cl%4 == 0 && GenFeiJiDaiDan(c, out) {
		out.Cards = append([]byte{}, c...)
		out.Type = CT_AIRCRAFT_ONE
		return true
	}

	if cl%5 == 0 && GenFeiJiDaiDui(c, out) {
		out.Cards = append([]byte{}, c...)
		out.Type = CT_AIRCRAFT_TWO
		return true
	}
	return false
}

func GenDanShun(cards []byte, out *GGameOutCard) bool {
	if len(cards) == 0 {
		return false
	}
	if GetLogicValue(cards[0]) >= 15 {
		return false
	}

	for i := 0; i < len(cards)-1; i++ {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+1])+1 {
			return false
		}
	}
	return true
}

func GenShuangShun(cards []byte, out *GGameOutCard) bool {
	if len(cards) == 0 {
		return false
	}
	if GetLogicValue(cards[0]) >= 15 {
		return false
	}

	for i := 0; i < len(cards)-2; i++ {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+2])+1 {
			return false
		}
	}
	return true
}

func GenFeiJi(cards []byte, out *GGameOutCard) bool {
	for i := 0; i < len(cards)-3; i++ {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+3])+1 {
			return false
		}
	}
	return true
}

func GenFeiJiDaiDan(cards []byte, out *GGameOutCard) bool {
	fnum := len(cards) / 4
	for _, v := range cards {
		dais := []byte{}
		feijis := []byte{}
		max := GetLogicValue(v)
		for j := 0; j < len(cards); j++ {
			if (GetLogicValue(cards[j]) <= max-byte(fnum)) ||
				(GetLogicValue(cards[j]) > max) {
				dais = append(dais, cards[j])
			} else {
				feijis = append(feijis, cards[j])
			}
		}
		if len(dais) != fnum {
			continue
		}
		if !GenFeiJi(feijis, out) {
			continue
		}
		out.Max = v
		return true
	}
	return false
}

func GenFeiJiDaiDui(cards []byte, out *GGameOutCard) bool {
	fnum := len(cards) / 5
	for _, v := range cards {
		dais := []byte{}
		feijis := []byte{}
		max := GetLogicValue(v)
		for j := 0; j < len(cards); j++ {
			if (GetLogicValue(cards[j]) <= max-byte(fnum)) ||
				(GetLogicValue(cards[j]) > max) {
				dais = append(dais, cards[j])
			} else {
				feijis = append(feijis, cards[j])
			}
		}
		if len(dais) != fnum*2 {
			continue
		}
		if !GenFeiJi(feijis, out) {
			continue
		}
		//检查带的是不是对
		meet := true
		for l := 0; l < len(dais)-1; l += 2 {
			if GetLogicValue(dais[l]) != GetLogicValue(dais[l+1]) {
				meet = false
				break
			}
		}
		if !meet {
			continue
		}
		out.Max = v
		return true
	}
	return false
}

/////////////////////////////////////////////////////////////////
//最大的人出牌，计算出的牌
var G_Type []int32 = []int32{CT_SINGLE,
	CT_DOUBLE,
	CT_SINGLE_CONNECT,
	CT_DOUBLE_CONNECT,
	CT_THREE,
	CT_THREE_LINE_TAKE_ONE,
	CT_THREE_LINE_TAKE_TWO,
	CT_FOUR_LINE_TAKE_ONE,
	CT_FOUR_LINE_TAKE_TWO,
	CT_AIRCRAFT,
	CT_AIRCRAFT_ONE,
	CT_AIRCRAFT_TWO,
	CT_BOMB_FOUR,
	CT_TWOKING,
}

type OutType struct {
	Num    int
	BStart int
	BEnd   int
}

//出牌的张数，可以自己修改
var G_OutType []OutType = []OutType{{1, 0, 20},
	{2, 20, 30},
	{5, 30, 200},
	{6, 200, 400},
	{3, 400, 500},
	{4, 500, 600},
	{5, 600, 700},
	{6, 700, 800},
	{8, 800, 830},
	{6, 830, 860},
	{8, 860, 880},
	{10, 880, 900},
	{4, 900, 980},
	{2, 980, 1000},
}

func GetOutCard(cards []byte, out *GGameOutCard) bool {
	if len(cards) == 0 {
		return false
	}
	rid := rand.Intn(1000)
	outtype := 0
	outnum := 0
	for i, v := range G_OutType {
		if rid >= v.BStart && rid < v.BEnd {
			outtype = i + 1
			outnum = v.Num
			break
		}
	}
	//
	lenc := len(cards)
	if lenc < outnum {
		if lenc >= 4 {
			outtype = CT_THREE_LINE_TAKE_ONE
			outnum = 4
		} else if lenc >= 3 {
			outtype = CT_THREE
			outnum = 3
		} else if lenc >= 2 {
			outtype = CT_DOUBLE
			outnum = 2
		} else {
			outtype = CT_SINGLE
			outnum = 1
		}

	}
	//
	re := false
	for i := lenc; i >= lenc-(lenc-outnum); i-- {
		if DoGenCard(cards[i-outnum:i], int32(outnum), int32(outtype), out) {
			re = true
			break
		}
	}
	if !re && lenc >= 4 {
		outtype = CT_THREE_LINE_TAKE_ONE
		outnum = 4
		for i := lenc; i >= lenc-(lenc-outnum); i-- {
			if DoGenCard(cards[i-outnum:i], int32(outnum), int32(outtype), out) {
				re = true
				break
			}
		}
	}
	if !re && lenc >= 2 {
		outtype = CT_DOUBLE
		outnum = 2
		for i := lenc; i >= lenc-(lenc-outnum); i-- {
			if DoGenCard(cards[i-outnum:i], int32(outnum), int32(outtype), out) {
				re = true
				break
			}
		}
	}
	if !re {
		re = DoGenCard(cards[lenc-1:lenc], 1, 1, out)
		if !re {
			return false
		}
	}
	return true
}

func DoGenCard(cards []byte, num, style int32, out *GGameOutCard) bool {
	re := false
	switch num {
	case 1:
		re = Gen1OutCard(cards, out)
	case 2:
		re = Gen2OutCard(cards, out)
	case 3:
		re = Gen3OutCard(cards, out)
	case 4:
		re = Gen4OutCard(cards, out)
	case 5:
		re = Gen5OutCard(cards, out)
	case 6:
		re = Gen6OutCard(cards, out)
	case 7:
		re = Gen7OutCard(cards, out)
	case 8:
		re = Gen8OutCard(cards, out)
	default:
		re = GenDefault(cards, out)
	}
	return re
}

func GetSecondOutCard(cards []byte, out *GGameOutCard, other *GGameOutCardReply) bool {
	outtype := other.Type
	outnum := len(other.Cards)
	lenc := len(cards)
	if lenc < outnum {
		//过
		return false
	}
	//
	re := false
	for i := lenc; i >= lenc-(lenc-outnum); i-- {
		gcards := cards[i-outnum : i]
		if len(gcards) == 0 {
			return false
		}
		if DoGenCard(gcards, int32(outnum), outtype, out) {
			if GetLogicValue(out.Max) > GetLogicValue(other.Max) && out.Type == outtype {
				re = true
				break
			}
		}
	}
	//
	//
	return re
}

//////////////////////////////////////////////////////////////////

package main

import (
	// "logs"
	"math/rand"
	"time"
)

const (
	CARD_COLOR   = 0xF0 //花色掩码
	CARD_VALUE   = 0x0F //数值掩码
	Card_Invalid = 0x00
	Card_Rear    = 0xFF
)

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

const (
	// 王
	Card_King_1 = iota + 0x41
	Card_King_2
)

const (
	CARD_COLOR_Fang = iota
	CARD_COLOR_Mei
	CARD_COLOR_Hong
	CARD_COLOR_Hei
	CARD_COLOR_King
	CARD_COLOR_Invalid
)

//扑克类型
const (
	CT_ERROR               = 0  //错误类型
	CT_SINGLE              = 1  //单牌类型
	CT_DOUBLE              = 2  //对子类型
	CT_SINGLE_CONNECT      = 3  //单龙
	CT_DOUBLE_CONNECT      = 4  //双龙
	CT_THREE               = 5  //三张
	CT_THREE_LINE_TAKE_ONE = 6  //三带一单
	CT_THREE_LINE_TAKE_TWO = 7  //三带一对
	CT_FOUR_LINE_TAKE_ONE  = 8  //四带两单
	CT_FOUR_LINE_TAKE_TWO  = 9  //四带两对
	CT_AIRCRAFT            = 10 //飞机
	CT_AIRCRAFT_ONE        = 11 //飞机带单
	CT_AIRCRAFT_TWO        = 12 //飞机带对
	CT_BOMB_FOUR           = 13 //炸弹
	CT_TWOKING             = 14 //对王类型
)

/////////////////////////////////////////////////////////
//卡牌管理器，负责做牌
type MgrCard struct {
	MVCard       []byte
	MVSourceCard []byte
	MSendId      int
}

func (this *MgrCard) InitCards() {
	this.MVCard = []byte{}
	this.MVSourceCard = []byte{}
	this.MSendId = 0
}

func (this *MgrCard) InitNormalCards() {
	begaincard := []byte{Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1}
	for _, v := range begaincard {
		for j := byte(0); j < 13; j++ {
			this.MVCard = append(this.MVCard, v+j)
		}
	}

	// 添加大小王
	this.MVCard = append(this.MVCard, Card_King_1, Card_King_2)
}

//发手牌
func (this *MgrCard) SendHandCard(num int) []byte {
	this.MSendId += num
	return append([]byte{}, this.MVSourceCard[this.MSendId-num:this.MSendId]...)
}

//剩余牌数，超过返回0
func (this *MgrCard) GetLeftCardCount() int {
	if this.MSendId > int(len(this.MVSourceCard)) {
		return 0
	}
	return len(this.MVSourceCard) - this.MSendId
}

func (this *MgrCard) GetSendCardCount() int {
	return this.MSendId
}

//洗牌
func (this *MgrCard) Shuffle() {

	this.MSendId = 0
	this.MVSourceCard = append([]byte{}, this.MVCard...)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(this.MVCard))
	for i, randIndex := range perm {
		this.MVSourceCard[i] = this.MVCard[randIndex]
	}
}

/////////////////////////////////////////////////////////
// 牌值获取函数
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

func CheckDan(out []byte, max byte) bool {
	if len(out) != 1 {
		return false
	}
	return out[0] == max
}

func CheckDui(out []byte, max byte) bool {
	if len(out) != 2 {
		return false
	}
	if GetLogicValue(out[0]) != GetLogicValue(out[1]) {
		return false
	}
	return out[0] == max
}

func CheckDanShun(out []byte, max byte) bool {
	if len(out) < 5 {
		return false
	}
	if GetLogicValue(out[0]) >= 15 {
		return false
	}
	for i := 0; i < len(out)-1; i++ {
		if GetLogicValue(out[i]) != GetLogicValue(out[i+1])+1 {
			return false
		}
	}
	return out[0] == max
}

func CheckShuangShun(out []byte, max byte) bool {
	if len(out) < 6 || len(out)%2 != 0 {
		return false
	}
	if GetLogicValue(out[0]) >= 15 {
		return false
	}
	for i := 0; i < len(out)-2; i++ {
		if GetLogicValue(out[i]) != GetLogicValue(out[i+2])+1 {
			return false
		}
	}
	return out[0] == max
}

func CheckSanZhang(out []byte, max byte) bool {
	if len(out) != 3 {
		return false
	}
	if GetLogicValue(out[0]) != GetLogicValue(out[2]) {
		return false
	}
	return out[0] == max
}

func CheckSanDaiYi(out []byte, max byte) bool {
	if len(out) != 4 {
		return false
	}
	admax := byte(0)
	if GetLogicValue(out[0]) == GetLogicValue(out[2]) && GetLogicValue(out[1]) != GetLogicValue(out[3]) {
		admax = out[0]
	} else if GetLogicValue(out[0]) != GetLogicValue(out[2]) && GetLogicValue(out[1]) == GetLogicValue(out[3]) {
		admax = out[1]
	}
	if admax == 0 {
		return false
	}
	return admax == max
}

func CheckSanDaiEr(out []byte, max byte) bool {
	if len(out) != 5 {
		return false
	}
	admax := byte(0)
	if GetLogicValue(out[0]) == GetLogicValue(out[2]) && GetLogicValue(out[3]) == GetLogicValue(out[4]) {
		admax = out[0]
	} else if GetLogicValue(out[0]) == GetLogicValue(out[1]) && GetLogicValue(out[2]) == GetLogicValue(out[4]) {
		admax = out[2]
	}
	if admax == 0 {
		return false
	}
	return admax == max
}

func CheckSiDaiDan(out []byte, max byte) bool {
	if len(out) != 6 {
		return false
	}
	admax := byte(0)
	if GetLogicValue(out[0]) == GetLogicValue(out[3]) {
		admax = out[0]
	} else if GetLogicValue(out[1]) == GetLogicValue(out[4]) {
		admax = out[1]
	} else if GetLogicValue(out[2]) == GetLogicValue(out[5]) {
		admax = out[2]
		if out[0] == 0x42 && out[1] == 0x41 {
			return false
		}
	}
	if admax == 0 {
		return false
	}
	return admax == max
}

func CheckSiDaiDui(out []byte, max byte) bool {
	if len(out) != 8 {
		return false
	}
	admax := byte(0)
	if GetLogicValue(out[0]) == GetLogicValue(out[3]) &&
		GetLogicValue(out[4]) == GetLogicValue(out[5]) &&
		GetLogicValue(out[6]) == GetLogicValue(out[7]) {
		admax = out[0]
	} else if GetLogicValue(out[4]) == GetLogicValue(out[7]) &&
		GetLogicValue(out[0]) == GetLogicValue(out[1]) &&
		GetLogicValue(out[2]) == GetLogicValue(out[3]) {
		admax = out[4]
	} else if GetLogicValue(out[2]) == GetLogicValue(out[5]) &&
		GetLogicValue(out[0]) == GetLogicValue(out[1]) &&
		GetLogicValue(out[6]) == GetLogicValue(out[7]) {
		admax = out[2]
	}
	if admax == 0 {
		return false
	}
	return admax == max
}

func CheckFeiJi(out []byte, max byte) bool {
	if len(out) < 6 || len(out)%3 != 0 {
		return false
	}
	if GetLogicValue(out[0]) >= 15 {
		return false
	}
	for i := 0; i < len(out)-3; i++ {
		if GetLogicValue(out[i]) != GetLogicValue(out[i+3])+1 {
			return false
		}
	}
	return out[0] == max
}

func CheckFeiJiDaiDan(out []byte, max byte) bool {
	feijinum := byte(0)
	if len(out) == 8 {
		feijinum = 2
	} else if len(out) == 12 {
		feijinum = 3
	} else if len(out) == 16 {
		feijinum = 4
	} else if len(out) == 20 {
		feijinum = 5
	} else {
		return false
	}

	vDaiCards := []byte{}
	vFeiCards := []byte{}
	admax := GetLogicValue(max)
	for i := 0; i < len(out); i++ {
		if GetLogicValue(out[i]) <= admax-feijinum || GetLogicValue(out[i]) > admax {
			vDaiCards = append(vDaiCards, out[i])
		} else {
			vFeiCards = append(vFeiCards, out[i])
		}
	}
	if len(vDaiCards) != int(feijinum) {
		return false
	}
	return CheckFeiJi(vFeiCards, max)
}

func CheckFeiJiDaiDui(out []byte, max byte) bool {
	feijinum := byte(0)
	if len(out) == 10 {
		feijinum = 2
	} else if len(out) == 15 {
		feijinum = 3
	} else if len(out) == 20 {
		feijinum = 4
	} else {
		return false
	}

	vDaiCards := []byte{}
	vFeiCards := []byte{}
	admax := GetLogicValue(max)
	for i := 0; i < len(out); i++ {
		if GetLogicValue(out[i]) <= admax-feijinum || GetLogicValue(out[i]) > admax {
			vDaiCards = append(vDaiCards, out[i])
		} else {
			vFeiCards = append(vFeiCards, out[i])
		}
	}
	if len(vDaiCards) != int(feijinum*2) {
		return false
	}
	//判断是否带对
	for i := 0; i < len(vDaiCards)-1; i += 2 {
		if GetLogicValue(vDaiCards[i]) != GetLogicValue(vDaiCards[i+1]) {
			return false
		}
	}
	return CheckFeiJi(vFeiCards, max)
}

func CheckZhadan(out []byte, max byte) bool {
	if len(out) != 4 {
		return false
	}
	//
	if GetLogicValue(out[0]) != GetLogicValue(out[3]) {
		return false
	}
	//
	return out[0] == max
}

func CheckWangZha(out []byte, max byte) bool {
	if len(out) != 2 {
		return false
	}
	if out[0] != 0x42 || out[1] != 0x41 {
		return false
	}
	//
	return out[0] == max
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

///////
func Gen1OutCard(cards []byte, out *GOutCard) bool {
	out.Max = cards[0]
	out.Cards = append([]byte{}, cards...)
	out.Type = CT_SINGLE
	return true
}

func Gen2OutCard(cards []byte, out *GOutCard) bool {
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

func Gen3OutCard(cards []byte, out *GOutCard) bool {
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
func Gen4OutCard(cards []byte, out *GOutCard) bool {
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
func Gen5OutCard(cards []byte, out *GOutCard) bool {
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
func Gen6OutCard(cards []byte, out *GOutCard) bool {
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
func Gen7OutCard(cards []byte, out *GOutCard) bool {
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
func Gen8OutCard(cards []byte, out *GOutCard) bool {
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

func GenDefault(cards []byte, out *GOutCard) bool {
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

func GenDanShun(cards []byte, out *GOutCard) bool {
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

func GenShuangShun(cards []byte, out *GOutCard) bool {
	if GetLogicValue(cards[0]) >= 15 {
		return false
	}

	if GetLogicValue(cards[0]) != GetLogicValue(cards[1]) {
		return false
	}

	for i := 0; i < len(cards)-2; i++ {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+2])+1 {
			return false
		}
	}
	return true
}

func GenFeiJi(cards []byte, out *GOutCard) bool {
	if GetLogicValue(cards[0]) != GetLogicValue(cards[1]) ||
		GetLogicValue(cards[1]) != GetLogicValue(cards[2]) {
		return false
	}
	if GetLogicValue(cards[0]) >= 15 {
		return false
	}
	for i := 0; i < len(cards)-3; i++ {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+3])+1 {
			return false
		}
	}
	return true
}

func GenFeiJiDaiDan(cards []byte, out *GOutCard) bool {
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

		ok := true
		for di := 1; di < len(dais); di++ {
			if GetLogicValue(dais[di]) == GetLogicValue(dais[di-1]) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		out.Max = v
		return true
	}
	return false
}

func GenFeiJiDaiDui(cards []byte, out *GOutCard) bool {
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

func DoGenCard(cards []byte, out *GOutCard) bool {
	re := false
	num := len(cards)
	if num == 0 {
		return false
	}
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

func DoGenCard2(cards []byte, num, style int32, out *GOutCard) bool {
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

func GenFirstOutCard(cards []byte, out *GOutCard) bool {
	if len(cards) == 0 {
		return false
	}
	//
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
		if DoGenCard2(cards[i-outnum:i], int32(outnum), int32(outtype), out) {
			re = true
			break
		}
	}
	if !re && lenc >= 4 {
		outtype = CT_THREE_LINE_TAKE_ONE
		outnum = 4
		for i := lenc; i >= lenc-(lenc-outnum); i-- {
			if DoGenCard2(cards[i-outnum:i], int32(outnum), int32(outtype), out) {
				re = true
				break
			}
		}
	}
	if !re && lenc >= 2 {
		outtype = CT_DOUBLE
		outnum = 2
		for i := lenc; i >= lenc-(lenc-outnum); i-- {
			if DoGenCard2(cards[i-outnum:i], int32(outnum), int32(outtype), out) {
				re = true
				break
			}
		}
	}
	if !re {
		re = DoGenCard2(cards[lenc-1:lenc], 1, 1, out)
		if !re {
			return false
		}
	}
	return true
}

func GenSecondOutCard(cards []byte, out *GOutCard, other *GOutCard) bool {
	outtype := other.Type
	outnum := len(other.Cards)
	lenc := len(cards)
	if lenc < outnum {
		//过
		return false
	}
	//
	likeout := rand.Intn(100)
	if likeout > 50 {
		return false
	}
	//
	re := false
	for i := lenc; i >= lenc-(lenc-outnum); i-- {
		gcards := cards[i-outnum : i]
		if len(gcards) == 0 {
			return false
		}
		if DoGenCard2(gcards, int32(outnum), outtype, out) {
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

//
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

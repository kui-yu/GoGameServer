package main

import (
	"logs"
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
// // 王
// Card_King_1 = iota + 0x41
// Card_King_2
)

const (
	CARD_COLOR_Fang = iota
	CARD_COLOR_Mei
	CARD_COLOR_Hong
	CARD_COLOR_Hei
	// CARD_COLOR_King
	CARD_COLOR_Invalid
)

//扑克类型
const (
	CT_ERROR               = iota //错误类型
	CT_SINGLE                     //单牌类型
	CT_DOUBLE                     //对子类型
	CT_SINGLE_CONNECT             //单龙
	CT_DOUBLE_CONNECT             //双龙
	CT_THREE                      //三张
	CT_THREE_LINE_TAKE_ONE        //三带一单
	CT_THREE_LINE_TAKE_TWO        //三带一对
	CT_THREE_DARGON
	CT_FOUR_LINE_TAKE_THREE //四带三
	CT_AIRCRAFT             //飞机(3-2)
	CT_BOMB_FOUR            //炸弹
)

type LastOutCards struct {
	Max   byte
	Type  int
	Cards []byte
	Cid   int32
}

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
			if v+j != 0x31 && v+j != 0x02 && v+j != 0x12 && v+j != 0x22 {
				this.MVCard = append(this.MVCard, v+j)
			}
		}
	}
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

//获取牌值
func GetLogicValue(card byte) byte {
	d := GetCardValue(card)
	if d <= 2 {
		return d + 13
	}
	return d
}

//将牌从大到小排序
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

//检测时候手牌中是否存在某些牌，如果有，就返回true  并且返回去除这些牌之后的手牌
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

//判断单牌类型
func Check_CT_SINGE(cards []byte) (bool, byte) {
	if len(cards) != 1 {
		return false, 0
	}
	return true, cards[0]
}

//判断对子类型
func Check_CT_DOUBLE(cards []byte) (bool, byte) {
	if len(cards) != 2 {
		return false, 0

	}
	pd := true
	if GetLogicValue(cards[0]) != GetLogicValue(cards[1]) {
		pd = false
	}
	return pd, cards[0]
}

//判断单龙
func Check_CT_SINGLE_CONNECT(cards []byte) (bool, byte) {
	pd := true
	if len(cards) < 5 || GetLogicValue(cards[0]) == GetLogicValue(Card_Hei_2) {
		pd = false
		return pd, cards[0]
	}

	if GetLogicValue(cards[0]) < 15 {
		for i := 0; i < len(cards)-1; i++ {
			if GetLogicValue(cards[i]) != GetLogicValue(cards[i+1])+1 {
				pd = false
			}
		}
	} else {
		pd = false
	}
	return pd, cards[0]
}

//判断双龙
func Check_CT_DOUBLE_CONNECT(cards []byte) (bool, byte) {
	pd := true
	if len(cards) < 4 || len(cards)%2 != 0 || GetLogicValue(cards[0]) == GetLogicValue(Card_Hei_2) {
		pd = false
		return pd, cards[0]
	}
	for i := 0; i < len(cards)-3; i += 2 {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+2])+1 {
			pd = false
		}
	}
	for i := 1; i < len(cards)-2; i += 2 {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+1])+1 {
			pd = false
		}
	}
	return pd, cards[0]
}

//判断三顺
func Check_CT_THREE_DARGON(cards []byte) (bool, byte) {
	pd := true
	if len(cards) < 6 || len(cards)%3 != 0 || GetLogicValue(cards[0]) == GetLogicValue(Card_Hei_2) {
		pd = false
		return pd, cards[0]
	}
	for i := 0; i < len(cards)-5; i += 3 {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+4])+1 {
			pd = false
		}
	}
	for i := 1; i < len(cards)-4; i += 3 {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+3])+1 {
			pd = false
		}
	}
	for i := 1; i < len(cards)-3; i += 3 {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+2])+1 {
			pd = false
		}
	}
	return pd, cards[0]
}

//三张
func Check_CT_THREE(cards []byte, handcardsnum int) (bool, byte) {
	pd := true
	if len(cards) != 3 {
		pd = false
		return pd, cards[0]
	}
	if handcardsnum != len(cards) {
		pd = false
	} else {
		if GetLogicValue(cards[0]) != GetLogicValue(cards[len(cards)-1]) {
			pd = false
		}
	}
	return pd, cards[0]
}

//小郑
//查找卡牌中是否存在n张重复的牌(返回找到与否和找到第一个重复的n张牌的首张下标)如n为3，有4张重复的返回false
func Find_repetitCard(c []byte, n int) (bool, int) {
	if n > len(c) {
		return false, 0
	}
	for i := 0; i <= len(c)-n; i++ {
		if GetLogicValue(c[i]) == GetLogicValue(c[i+n-1]) {
			if i+n < len(c) && GetLogicValue(c[i]) == GetLogicValue(c[i+n]) {
				return false, 0
			}
			return true, i
		}
	}
	return false, 0
}

//三代一
func Check_CT_THREE_LINE_TAKE_ONE(cards []byte, handcardsnum int) (bool, byte) {
	if len(cards) != handcardsnum || len(cards) != 4 || handcardsnum != 4 {
		return false, 0
	}

	i, index := Find_repetitCard(cards, 3)
	if i == true {
		return true, cards[index]
	}
	return false, 0
}

//三代二
func Check_CT_THREE_LINE_TAKE_TWO(c []byte) (bool, byte) {
	if len(c) != 5 {
		return false, 0
	}
	if i, index := Find_repetitCard(c, 3); i {
		return true, c[index]
	}
	return false, 0
}

//四代三
func Check_CT_FOUR_LINE_TAKE_THREE(c []byte) (bool, byte) {
	if len(c) != 7 {
		return false, 0
	}

	if i, index := Find_repetitCard(c, 4); i {
		return true, c[index]
	}
	return false, 0
}

func Check_CT_AIRCRAFT(cards []byte) (bool, byte) {
	num := len(cards)
	if num%5 != 0 || num == 0 {
		return false, 0
	}
	//不允许飞机里面有带炸弹
	if c4, _ := Find_repetitCard(cards, 4); c4 {
		return false, 0
	}
	for index := 0; index < num; {
		j, index1 := Find_repetitCard(cards[index:], 3)
		if j == false {
			return false, 0
		}
		//判断后面是否真的有num/5个连飞机
		if index+index1+num/5*3 > len(cards) {
			break
		}
		//判断是否真的有num/5个飞机
		if x, _ := Check_CT_THREE_DARGON(cards[index+index1 : index+index1+num/5*3]); x {
			n := append([]byte{}, cards[index+index1:index+index1+num/5*3]...)
			n = append(n, cards[:index+index1]...)
			n = append(n, cards[index+index1+num/5*3:]...)
			return true, n[0]
		}
		index += index1 + 3*num/5
	}
	return false, 0
}

//炸弹
func Check_CT_BOMB_FOUR(cards []byte) (bool, byte) {
	if len(cards) != 4 {
		return false, 0
	}
	if i, index := Find_repetitCard(cards, 4); i {
		return true, cards[index]
	}
	return false, 0
}

//检测牌型总方法
func Check_Chard(cards []byte, handcardsnum int) (canout bool, lastoutcards LastOutCards) {
	if len(cards) == 0 {
		return false, LastOutCards{}
	}
	cards = Sort(cards)
	lastoutcards.Cards = cards
	if pd, max := Check_CT_SINGE(cards); pd {
		lastoutcards.Max = max
		lastoutcards.Type = CT_SINGLE
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_DOUBLE(cards); pd {
		lastoutcards.Max = max
		lastoutcards.Type = CT_DOUBLE
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_SINGLE_CONNECT(cards); pd {
		logs.Debug("检测到该牌为顺子")
		lastoutcards.Max = max
		lastoutcards.Type = CT_SINGLE_CONNECT
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_DOUBLE_CONNECT(cards); pd {
		logs.Debug("检测到该牌为双龙")
		lastoutcards.Max = max
		lastoutcards.Type = CT_DOUBLE_CONNECT
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_THREE(cards, handcardsnum); pd {
		lastoutcards.Max = max
		lastoutcards.Type = CT_THREE
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_THREE_DARGON(cards); pd {
		lastoutcards.Max = max
		lastoutcards.Type = CT_THREE_DARGON
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_THREE_LINE_TAKE_ONE(cards, handcardsnum); pd {
		lastoutcards.Max = max
		lastoutcards.Type = CT_THREE_LINE_TAKE_ONE
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_THREE_LINE_TAKE_TWO(cards); pd {
		lastoutcards.Max = max
		lastoutcards.Type = CT_THREE_LINE_TAKE_TWO
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_FOUR_LINE_TAKE_THREE(cards); pd {
		lastoutcards.Max = max
		lastoutcards.Type = CT_FOUR_LINE_TAKE_THREE
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_AIRCRAFT(cards); pd {
		logs.Debug("检测到改牌为飞机")
		lastoutcards.Max = max
		lastoutcards.Type = CT_AIRCRAFT
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}
	if pd, max := Check_CT_BOMB_FOUR(cards); pd {
		lastoutcards.Max = max
		lastoutcards.Type = CT_BOMB_FOUR
		lastoutcards.Cards = SortToSee(lastoutcards)
		return true, lastoutcards
	}

	return false, LastOutCards{}
}

/////////////////////////////////////////
//将玩家出的牌，进行排序
func SortToSee(lastout LastOutCards) []byte {
	//根据出牌类型 整理牌组用于显示桌面上
	switch lastout.Type {
	case CT_SINGLE:
		return Sort(lastout.Cards)
		break
	case CT_DOUBLE:
		return Sort(lastout.Cards)
		break
	case CT_SINGLE_CONNECT:
		return Sort(lastout.Cards)
		break
	case CT_DOUBLE_CONNECT:
		return Sort(lastout.Cards)
		break
	case CT_THREE_DARGON:
		return Sort(lastout.Cards)
		break
	case CT_THREE:
		return Sort(lastout.Cards)
		break
	case CT_THREE_LINE_TAKE_ONE:
		return SortSanDaiAndFeiJi(lastout)
		break
	case CT_THREE_LINE_TAKE_TWO:
		return SortSanDaiAndFeiJi(lastout)
		break
	case CT_FOUR_LINE_TAKE_THREE:
		return SortSiDaiSan(lastout)
		break
	case CT_AIRCRAFT:
		return SortSanDaiAndFeiJi(lastout)
		break
	case CT_BOMB_FOUR:
		return Sort(lastout.Cards)
		break
	default:
		logs.Debug("在进行玩家出牌 牌组整理的时候，玩家类型出现了错误")
		break
	}
	return nil
}

//排序3带1,3带2 飞机
func SortSanDaiAndFeiJi(lastout LastOutCards) []byte {
	lastout.Cards = Sort(lastout.Cards)
	//根据玩家出牌长度判断飞机个数
	feijinum := len(lastout.Cards) / 5
	//在根据最大值推算飞机
	nums := []byte{}
	if feijinum < 1 {
		nums = append(nums, GetLogicValue(lastout.Max))
	} else {
		for i := 0; i < feijinum; i++ {
			nums = append(nums, GetLogicValue(lastout.Max)-byte(i))
		}
	}

	result, quchu := ChouNumToCards(lastout.Cards, nums)
	result = Sort(result)
	quchu = append(quchu, result...)
	return quchu
}

//排序4带3
func SortSiDaiSan(lastout LastOutCards) []byte {
	lastout.Cards = Sort(lastout.Cards)
	nums := []byte{}
	nums = append(nums, lastout.Max)
	result, quchu := ChouNumToCards(lastout.Cards, nums)
	result = Sort(result)
	quchu = append(quchu, result...)
	return quchu
}

//将值为Num的牌从牌组中抽出
func ChouNumToCards(cards []byte, nums []byte) ([]byte, []byte) {
	cards = Sort(cards)
	quchu := []byte{}
	for _, num := range nums {
		for i := len(cards) - 1; i >= 0; i-- {
			if GetLogicValue(cards[i]) == num {
				quchu = append(quchu, cards[i])
				cards = append(cards[:i], cards[i+1:]...)
			}
		}
	}
	return cards, quchu
}

//从手牌中返回 num张一样的扑克牌数组
func FindCardsByNum(num int, handCards1 []byte, notoc byte) [][]byte {
	handCards := append([]byte{}, handCards1...)
	result := [][]byte{}
	item := []byte{}
	if num > len(handCards) {
		return result
	}
	handCards = Sort(handCards)
	for i := 0; i <= len(handCards)-num; i++ {
		if GetLogicValue(handCards[i]) == notoc {
			continue
		}
		if GetLogicValue(handCards[i]) == GetLogicValue(handCards[i+num-1]) {
			if i+num < len(handCards) && GetLogicValue(handCards[i]) != GetLogicValue(handCards[i+num]) || i+num >= len(handCards) {
				if i != 0 && GetLogicValue(handCards[i]) != GetLogicValue(handCards[i-1]) || i == 0 {
					for j := 0; j < num; j++ {
						item = append(item, handCards[i+j])
					}
					result = append(result, item)
					fresult := FindCardsByNum(num, handCards[i+num-1:], GetLogicValue(handCards[i]))
					if len(fresult) > 0 {
						result = append(result, fresult...)
					}
					return result
				}
			}
		}
	}
	return result
}

//将牌组删减为num相同数量的牌 例如 num=2 则将牌组 [1,2,3,4,4,4,55,6666] 变成 [4,4,5,5,6,6]
func CutCardsByNum(num int, handCards []byte) []byte {
	handCards = Sort(handCards)
	result := []byte{}
	item := []byte{}
	var noneed byte
	for i := 0; i < len(handCards); i++ {
		if len(item) == 0 {
			if noneed != GetLogicValue(handCards[i]) {
				item = append(item, handCards[i])
			}
		} else {
			if GetLogicValue(item[0]) == GetLogicValue(handCards[i]) && noneed != GetLogicValue(handCards[i]) {
				item = append(item, handCards[i])
			} else if GetLogicValue(item[0]) != GetLogicValue(handCards[i]) {
				item = []byte{}
				item = append(item, handCards[i])
			}
		}
		if len(item) == num {
			result = append(result, item...)
			noneed = GetLogicValue(item[0])
			item = []byte{}
		}
	}
	return result
}

//返回剪切之后的Num数量牌组集合
func ChangeCutCardsByNum(num int, handCards []byte) [][]byte {
	result := [][]byte{}
	handCards = Sort(handCards)
	//调用剪切方法
	res := CutCardsByNum(num, handCards)
	re := []byte{}
	count := 0
	for i := 0; i < len(res); i++ {
		if count == num {
			count = 1
			result = append(result, re)
			re = []byte{}
			re = append(re, res[i])
		} else {
			re = append(re, res[i])
			count++
		}
	}
	if len(re) != 0 {
		result = append(result, re)
	}
	return result
}

////////////////////////////////////
//从手牌中查找单牌
func FindCardsDan(handCards []byte) [][]byte {
	return ChangeCutCardsByNum(1, handCards)
}

//从手牌中查找对子
func FindCardsShuang(handCards []byte) [][]byte {
	return ChangeCutCardsByNum(2, handCards)
}

//从手牌中查找顺子
func FindCardsSunzi(num int, handCards []byte) [][]byte {
	//num: 查找num连顺
	//handCards: 手牌
	//去除手牌中多余的牌,将其变为单排
	handCards = Sort(handCards)
	result := [][]byte{}
	item := []byte{}
	handCards = CutCardsByNum(1, handCards)
	if len(handCards) < num {
		return result
	}
	for i := 0; i <= len(handCards)-num; i++ {
		if GetLogicValue(handCards[i]) == GetLogicValue(handCards[i+num-1])+byte(num-1) {
			for j := 0; j < num; j++ {
				item = append(item, handCards[i+j])
			}
			result = append(result, item)
			fresult := FindCardsSunzi(num, handCards[i+1:])
			if len(fresult) > 0 {
				result = append(result, fresult...)
			}
			//顺子中不能存在2，所以我们需要将带有2的牌组删除
			if len(result) > 0 {
				for i := len(result) - 1; i >= 0; i-- {
					//循环每一个牌组 判断是否含有2
					pd := false
					for _, v := range result[i] {
						if GetLogicValue(v) == GetLogicValue(Card_Hei_2) {
							pd = true
							break
						}
					}
					if pd {
						result = append(result[:i], result[i+1:]...)
					}
				}
			}
			return result
		}
	}
	return result
}

//从手牌中查找双顺
func FindCardsShuangShun(num int, handCards []byte) [][]byte {
	handCards = Sort(handCards)
	result := [][]byte{}
	item := []byte{}
	handCards = CutCardsByNum(2, handCards)
	if len(handCards) < num*2 {
		return result
	}
	for i := 0; i <= len(handCards)-num*2; i++ {
		if GetLogicValue(handCards[i]) == GetLogicValue(handCards[i+num*2-1])+byte(num-1) {
			for j := 0; j < num*2; j++ {
				item = append(item, handCards[i+j])
			}
			result = append(result, item)
			fresult := FindCardsShuangShun(num, handCards[i+2:])
			if len(fresult) > 0 {
				result = append(result, fresult...)
			}
			//顺子中不能存在2，所以我们需要将带有2的牌组删除
			if len(result) > 0 {
				for i := len(result) - 1; i >= 0; i-- {
					//循环每一个牌组 判断是否含有2
					pd := false
					for _, v := range result[i] {
						if GetLogicValue(v) == GetLogicValue(Card_Hei_2) {
							pd = true
							break
						}
					}
					if pd {
						result = append(result[:i], result[i+1:]...)
					}
				}
			}

			return result
		}
	}
	return result
}

//从手牌中查找三顺
func FindCardsSanShun(num int, handCards []byte) [][]byte {
	handCards = Sort(handCards)
	result := [][]byte{}
	item := []byte{}
	handCards = CutCardsByNum(3, handCards)
	if len(handCards) < num*3 {
		return result
	}
	for i := 0; i <= len(handCards)-num*3; i++ {
		if GetLogicValue(handCards[i]) == GetLogicValue(handCards[i+num*3-1])+byte(num-1) {
			for j := 0; j < num*3; j++ {
				item = append(item, handCards[i+j])
			}
			result = append(result, item)
			fresult := FindCardsSanShun(num, handCards[i+3:])
			if len(fresult) > 0 {
				result = append(result, fresult...)
			}
			//顺子中不能存在2，所以我们需要将带有2的牌组删除
			if len(result) > 0 {
				for i := len(result) - 1; i >= 0; i-- {
					//循环每一个牌组 判断是否含有2
					pd := false
					for _, v := range result[i] {
						if GetLogicValue(v) == GetLogicValue(Card_Hei_2) {
							pd = true
							break
						}
					}
					if pd {
						result = append(result[:i], result[i+1:]...)
					}
				}
			}

			return result
		}
	}
	return result
}

//从手牌中查找到飞机 num:几个飞机
func FindCardsFeiJi(num int, handCards []byte) [][]byte {
	handCards = Sort(handCards)
	result := [][]byte{}
	threeAry := FindCardsSanShun(num, handCards)
	handc := []byte{}
	if len(threeAry) > 0 {
		//如果有对应的三顺，则需要判断该三顺所对应需要携带的翅膀
		needTake := num * 2
		//循环得到的3顺集合
		for _, v := range threeAry {
			handc = []byte{}
			handc = append(handc, handCards...) //疑问不知为何handcards会改变值，所以只能使用hanc来代替
			hc, ok := VecDelMulti(handc, v)
			if ok {
				//去除可能会和飞机组成炸弹的带牌、
				for i := len(hc) - 1; i >= 0; i-- {
					for _, v := range threeAry {
						for _, v1 := range v {
							if GetLogicValue(v1) == GetLogicValue(hc[i]) {
								hc = append([]byte{}, hc[0:len(hc)-1]...)
								break
							}
						}
					}
				}
				//查找出炸弹 （飞机无法带炸弹）
				boomAry := FindCardsByNum(4, hc, 0)
				for _, v1 := range boomAry {
					hc, _ = VecDelMulti(hc, v1)
				}
				if len(hc) >= needTake {
					v := append(v, hc[len(hc)-needTake:]...)
					result = append(result, v)
				} else {
					//只要不带整个炸弹，需要几张炸弹牌凑翅磅也是可以行的通的
					needcards := needTake - len(hc)
					if !(needcards > len(boomAry)*4-len(boomAry)) {
						res := []byte{}
						for i := len(boomAry) - 1; i >= 0; i-- {
							conti := false
							for j := 0; j < len(boomAry[i]); j++ {
								if needcards == 0 {
									res = append(res, hc...)
									res = Sort(res)
									v = append(v, res...)
									break
								} else {
									if j == 3 {
										conti = true
										break
									}
									res = append(res, boomAry[i][j])
									needcards--
								}
							}
							if conti {
								continue
							} else {
								break
							}
						}
						result = append(result, v)
					}
				}
			}
		}
	}

	for i := len(result) - 1; i >= 0; i-- {
		res := FindCardsByNum(4, result[i], 0)
		if len(res) > 0 {
			if i == len(result)-1 {
				result = append([][]byte{}, result[:len(result)-1]...)
			} else {
				result = append(result[:i], result[i+1:]...)
			}
		}
	}
	return result
}

//从手牌中判断是否存在3带二(顺便判断3带一，3张 是否可以出牌)
func FindCardsSanDaiEr(handCards []byte) [][]byte {
	result := [][]byte{}
	handCards = Sort(handCards)
	//查找出所有3个相同牌的牌组
	threeAry := ChangeCutCardsByNum(3, handCards)

	for _, v := range threeAry {
		hanc := append([]byte{}, handCards...)
		//需要去除3张，判断是否有可以带的牌
		hanc, ok := VecDelMulti(hanc, v)
		if ok {
			pd := true
			switch len(hanc) {
			case 1:
				v = append(v, hanc[0])
				if GetLogicValue(hanc[0]) == GetLogicValue(v[0]) {
					pd = false
				}
				break
			case 0:
				break
			default:
				for _, v1 := range hanc {
					if GetLogicValue(v1) == GetLogicValue(v[0]) {
						pd = false
					}
				}
				v = append(v, hanc[len(hanc)-2:]...)
				break
			}
			if pd {
				result = append(result, v)
			}
		}
	}
	return result
}

//从手牌中查找炸弹
func FindCardsBoom(hancards []byte) [][]byte {
	return ChangeCutCardsByNum(4, hancards)
}

//从手牌中查找四带三
func FindCardsBoomDaiSan(hancards []byte) [][]byte {
	result := [][]byte{}
	//整理手牌
	hancards = Sort(hancards)
	//取出4个牌组
	fourArry := ChangeCutCardsByNum(4, hancards)
	//判单剩余手牌是否>=需要带的牌数量
	for _, v := range fourArry {
		hanc := append([]byte{}, hancards...)
		hc, ok := VecDelMulti(hanc, v)
		hc = Sort(hc)
		if ok {
			if len(hc) >= 3 {
				v := append(v, hc[len(hc)-3:]...)
				result = append(result, v)
			}
		}
	}
	return result
}

////////////////////////////////////
//根据玩家出牌 判断自己手牌中是否存在能够出的牌
func FindType(handcards1 []byte, outcards LastOutCards) ([][]byte, [][]byte) {
	handcards := append([]byte{}, handcards1...)
	result := [][]byte{}
	//从大到小进行排列
	//通过玩家牌型判断自己能够出的牌型
	//查找炸弹
	boom := FindCardsBoom(handcards)
	if outcards.Max != 0 {
		re := [][]byte{}
		num := len(outcards.Cards)
		switch outcards.Type {
		case CT_SINGLE:
			re = FindCardsDan(handcards)
			break
		case CT_DOUBLE:
			re = FindCardsShuang(handcards)
			break
		case CT_SINGLE_CONNECT:
			re = FindCardsSunzi(num, handcards)
			break
		case CT_DOUBLE_CONNECT:
			re = FindCardsShuangShun(num/2, handcards)
			break
		case CT_THREE_DARGON:
			re = FindCardsSanShun(num/3, handcards)
			break
		case CT_FOUR_LINE_TAKE_THREE:
			re = FindCardsBoomDaiSan(handcards)
			break
		case CT_THREE_LINE_TAKE_TWO:
			re = FindCardsSanDaiEr(handcards)
		case CT_AIRCRAFT:
			re = FindCardsFeiJi(num/5, handcards)
			break
		case CT_BOMB_FOUR:
			break
		default:
			logs.Debug("根据上一家出牌类型，判断手牌中可出的牌失败")
			break
		}
		result = append(result, re...)
	} else {
		result = append(result, FindCardsDan(handcards)...)
		result = append(result, FindCardsShuang(handcards)...)
		result = append(result, FindCardsBoomDaiSan(handcards)...)
		result = append(result, FindCardsSanDaiEr(handcards)...)
	}
	return boom, result
}

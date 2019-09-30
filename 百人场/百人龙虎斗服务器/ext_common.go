package main

import (
	"bl.com/util"
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
	CARD_COLOR_Fang = iota
	CARD_COLOR_Mei
	CARD_COLOR_Hong
	CARD_COLOR_Hei
)

const (
	DRAGON byte = 1 + iota
	TIGER
	DRAW
)

// 花色获取函数
func GetCardColor(card int32) byte {
	return (byte(card) & 0xF0) >> 4
}

// 牌值获取函数
func GetCardValue(card int32) byte {
	return (byte(card) & 0x0F)
}

// 逻辑牌值获取函数
func GetLogicValue(card int32) byte {
	d := GetCardValue(card)

	return d
}

/////////////////////////////////////////////////////////
//卡牌管理器，负责做牌
type MgrCard struct {
	MVCard       []byte
	MVSourceCard []byte

	MSendId int // 已经切掉的牌
	MUsedId int // 已经使用过的牌
}

func (this *MgrCard) InitCards() {
	this.MVCard = []byte{}
	this.MVSourceCard = []byte{}
	this.MSendId = 0
	this.MUsedId = 0
}

func (this *MgrCard) InitNormalCards() {
	// 8 副牌
	begaincard := []byte{Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1,
		Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1,
		Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1,
		Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1,
		Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1,
		Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1,
		Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1,
		Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1}

	for _, v := range begaincard {
		for j := byte(0); j < 13; j++ {
			this.MVCard = append(this.MVCard, v+j)
		}
	}
}

// 获取一张牌
func (this *MgrCard) SendOneCard() int32 {
	this.MVSourceCard[this.MUsedId], this.MVSourceCard[this.MSendId] = this.MVSourceCard[this.MSendId], this.MVSourceCard[this.MUsedId]
	ret := int32(this.MVSourceCard[this.MUsedId])

	this.MSendId += 1
	this.MUsedId += 1

	return ret
}

func (this *MgrCard) SendCard(id int) int32 {
	this.MVSourceCard[this.MUsedId], this.MVSourceCard[id] = this.MVSourceCard[id], this.MVSourceCard[this.MUsedId]

	ret := int32(this.MVSourceCard[this.MUsedId])

	this.MUsedId += 1

	return ret
}

// 剩余牌数，超过返回0
func (this *MgrCard) GetLeftCardCount() int {
	if this.MSendId > len(this.MVSourceCard) {
		return 0
	}
	return len(this.MVSourceCard) - this.MSendId
}

func (this *MgrCard) GetSendCardCount() int {
	return this.MSendId
}

func (this *MgrCard) SetSendCardCount(count int) {
	this.MSendId = count
}

//洗牌
func (this *MgrCard) Shuffle() {
	this.MSendId = 0

	this.MVSourceCard = util.Shuffle(this.MVCard)

	last, _ := util.GetRandomNum(0, 3)     // 剩余牌张数，0-2张
	count, _ := util.GetRandomNum(80, 101) // 保留局数，80-100局

	this.MSendId = len(this.MVSourceCard) - last - 3*count

	this.MUsedId = 0
}

// 比牌
func (this *MgrCard) CompareCard(d, t int32) byte {
	f := GetCardValue(d)
	n := GetCardValue(t)

	if f > n {
		return DRAGON
	}

	if f < n {
		return TIGER
	}

	f = GetCardColor(d)
	n = GetCardColor(t)

	if f > n {
		return DRAGON
	}

	if f < n {
		return TIGER
	}

	return DRAW
}

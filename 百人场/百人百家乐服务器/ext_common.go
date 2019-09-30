package main

import (
	"math/rand"
	"time"

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
	IDLE byte = 1 + iota
	BANKER
	DRAW
)

// 牌值获取函数
func GetCardValue(card int32) byte {
	return (byte(card) & 0x0F)
}

// 逻辑牌值获取函数
func GetLogicValue(card int32) byte {
	d := GetCardValue(card)
	if d > 10 {
		d = 10
	}
	return d
}

/////////////////////////////////////////////////////////
//卡牌管理器，负责做牌
type MgrCard struct {
	MVCard       []byte
	MVSourceCard []byte
	OutCards     []byte
}

func (this *MgrCard) InitCards() {
	this.MVCard = []byte{}
	this.MVSourceCard = []byte{}
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
	this.OutCards = append(this.OutCards, this.MVSourceCard[0])
	this.MVSourceCard = this.MVSourceCard[1:]
	return int32(this.OutCards[len(this.OutCards)-1])
}

//洗牌
func (this *MgrCard) Shuffle() {
	this.MVSourceCard = util.Shuffle(this.MVCard)
	this.OutCards = make([]byte, 0)
}

// 比牌
func (this *MgrCard) CompareCard(idea, banker []int32) byte {
	var f byte
	var n byte
	for i := 0; i < len(idea); i++ {
		f = (f + GetLogicValue(idea[i])) % 10
	}

	for i := 0; i < len(banker); i++ {
		n = (n + GetLogicValue(banker[i])) % 10
	}

	if f > n {
		return IDLE
	}

	if f < n {
		return BANKER
	}

	return DRAW
}

//打乱牌
func (this *MgrCard) DisturbCards() {
	newCards := make([]byte, 0)
	rand.Seed(time.Now().UnixNano())
	rArr := rand.Perm(len(this.MVSourceCard))
	for _, v := range rArr {
		newCards = append(newCards, this.MVSourceCard[v])
	}
	this.MVSourceCard = newCards
}

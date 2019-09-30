package main

import (
	"crypto/rand"
	"math/big"

	"bl.com/goldenflower"
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

type MgrCard struct {
	MVCard       []byte
	MVSourceCard []byte

	GFError   byte // 类型错误
	GFSingle  byte // 单牌类型
	GFDouble  byte // 对子类型
	GFShunZi  byte // 顺子类型
	GFJinHua  byte // 金花类型
	GFShunJin byte // 顺金类型
	GFBaoZi   byte // 豹子类型
	GFSpecial byte // 特殊类型235
}

func (this *MgrCard) InitCards() {
	this.MVCard = []byte{}
	this.MVSourceCard = []byte{}
	// this.MSendId = 0

	this.GFError = goldenflower.GFError
	this.GFSingle = goldenflower.GFSingle
	this.GFDouble = goldenflower.GFDouble
	this.GFShunZi = goldenflower.GFShunZi
	this.GFJinHua = goldenflower.GFJinHua
	this.GFShunJin = goldenflower.GFShunJin
	this.GFBaoZi = goldenflower.GFBaoZi
}

func (this *MgrCard) InitNormalCards() {
	begaincard := []byte{Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1}
	for _, v := range begaincard {
		for j := byte(0); j < 13; j++ {
			this.MVCard = append(this.MVCard, v+j)
		}
	}
}

// 发牌
func (this *MgrCard) SendCard(num int) []int32 {
	//添加要发的手牌
	list := append([]byte{}, this.MVSourceCard[:num]...)
	//删除要发的手牌元素
	this.MVSourceCard = append([]byte{}, this.MVSourceCard[num:]...)

	var sourceList []int32
	for _, card := range list {
		sourceList = append(sourceList, int32(card))
	}
	//返回
	return sourceList
}

//洗牌
func (this *MgrCard) Shuffle() {
	this.MVSourceCard = util.Shuffle(this.MVCard)
}

//回收手牌
func (this *MgrCard) RecoverCard(cards []int32) {
	for _, card := range cards {
		this.MVSourceCard = append(this.MVSourceCard, uint8(card))
	}
}

//打乱牌组
func (this *MgrCard) UpsetCard() {

	var sourceList []uint8

	MVCard := append([]uint8{}, this.MVSourceCard...)

	// 随机打乱牌型
	for i := 0; i < len(this.MVSourceCard); i++ {
		//打乱
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(MVCard))))
		//添加
		sourceList = append(sourceList, MVCard[randIndex.Int64()])
		// 移除已经添加的牌
		MVCard = append(MVCard[:int(randIndex.Int64())], MVCard[int(randIndex.Int64())+1:]...)
	}

	this.MVSourceCard = sourceList
}

// 排序组合，获取最大牌
func (this *MgrCard) GetMaxCards(cards []int32) ([]int32, int32) {
	max := []byte{byte(cards[0]), byte(cards[1]), byte(cards[2])}
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			for k := j + 1; k < 5; k++ {
				c := []byte{byte(cards[i]), byte(cards[j]), byte(cards[k])}
				if goldenflower.CompareCard(c, max, false) {
					max = c[:]
				}
			}
		}
	}

	return []int32{int32(max[0]), int32(max[1]), int32(max[2])}, int32(goldenflower.GetCardsType(max, false))
}

// 牌型比较
func (this *MgrCard) CompareCard(f, n []int32) bool {
	if len(f) != len(n) || len(f) != 3 {
		return false
	}

	first := []byte{}
	next := []byte{}
	for i := 0; i < 3; i++ {
		first = append(first, byte(f[i]))
		next = append(next, byte(n[i]))
	}

	return goldenflower.CompareCard(first, next, false)
}

func (this *MgrCard) GetCardsType(c []int32) byte {
	if len(c) != 3 {
		return goldenflower.GFError
	}

	cards := []byte{byte(c[0]), byte(c[1]), byte(c[2])}

	return goldenflower.GetCardsType(cards, false)
}

func (this *MgrCard) GetLogicValue(v int32) byte {
	return goldenflower.GetCardValue(byte(v))
}

func (this *MgrCard) Sort(c []int32) []int32 {
	s := []byte{}
	for _, v := range c {
		s = append(s, byte(v))
	}

	s = goldenflower.Sort(s)

	ret := []int32{}
	for _, v := range s {
		ret = append(ret, int32(v))
	}

	return ret
}

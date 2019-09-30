package main

import (
	"crypto/rand"
	"math/big"

	"bl.com/paigow"
	"bl.com/util"
)

type MgrCard struct {
	MVCard       []byte
	MVSourceCard []byte

	//MSendId int // 已经切掉的牌
}

func (this *MgrCard) InitCards() {
	this.MVCard = []byte{}
	this.MVSourceCard = []byte{}
	// this.MSendId = 0
}

func (this *MgrCard) InitNormalCards() {
	this.MVCard = paigow.GetInitNormalCards()
}

//洗牌
func (this *MgrCard) Shuffle() {
	this.MVSourceCard = util.Shuffle(this.MVCard)
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

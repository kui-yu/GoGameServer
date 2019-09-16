package MaJiangTool

import (
	"logs"
	"math/rand"
	"time"
)

type CardManager108 struct {
	SendGangCardCount int //发送杠牌张数，即牌尾减少了多少张牌
	SendCardCount     int //已经发牌张数（前向发牌）
	AllCardCount      int //游戏中麻将牌的总个数
	HandCardCount     int //游戏中手牌个数包括当前牌
	BaoPos            int
	AllCardSource     []byte //牌池，用于洗牌用
	AllCard           []byte //牌池，用于发牌
}

func (this *CardManager108) Initialize() {
	for i := byte(0); i < 9; i++ {
		for j := byte(0); j < 4; j++ {
			this.AllCardSource = append(this.AllCardSource, Card_Wan_1+i)
		}
	}
	for i := byte(0); i < 9; i++ {
		for j := byte(0); j < 4; j++ {
			this.AllCardSource = append(this.AllCardSource, Card_Tiao_1+i)
		}
	}
	for i := byte(0); i < 9; i++ {
		for j := byte(0); j < 4; j++ {
			this.AllCardSource = append(this.AllCardSource, Card_Bing_1+i)
		}
	}
}

func (this *CardManager108) Shuffle() {
	this.AllCard = append([]byte{}, this.AllCardSource...)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(this.AllCardSource))
	for i, randIndex := range perm {
		this.AllCard[i] = this.AllCardSource[randIndex]
	}
	// //做牌
	// this.AllCard = append([]byte{}, Card_Wan_1, Card_Wan_1, Card_Wan_1)
	// this.AllCard = append(this.AllCard, Card_Wan_1, Card_Tiao_2, Card_Tiao_3, Card_Tiao_4, Card_Tiao_5)
	// this.AllCard = append(this.AllCard, Card_Tiao_5, Card_Tiao_6, Card_Tiao_7, Card_Tiao_8, Card_Tiao_9)
	// //
	// this.AllCard = append(this.AllCard, Card_Wan_1, Card_Wan_2, Card_Wan_3)
	// this.AllCard = append(this.AllCard, Card_Tiao_1, Card_Tiao_2, Card_Tiao_3, Card_Tiao_4, Card_Tiao_5)
	// this.AllCard = append(this.AllCard, Card_Tiao_5, Card_Tiao_6, Card_Tiao_7, Card_Tiao_8, Card_Tiao_9)
	// //
	// this.AllCard = append(this.AllCard, Card_Wan_1, Card_Wan_2, Card_Wan_3)
	// this.AllCard = append(this.AllCard, Card_Tiao_1, Card_Tiao_2, Card_Tiao_3, Card_Tiao_4, Card_Tiao_5)
	// this.AllCard = append(this.AllCard, Card_Tiao_5, Card_Tiao_6, Card_Tiao_7, Card_Tiao_8, Card_Tiao_9)
	// //
	// this.AllCard = append(this.AllCard, Card_Wan_1, Card_Wan_2, Card_Wan_3)
	// this.AllCard = append(this.AllCard, Card_Tiao_1, Card_Tiao_2, Card_Tiao_3, Card_Tiao_4, Card_Tiao_5)
	// this.AllCard = append(this.AllCard, Card_Tiao_5, Card_Tiao_6, Card_Tiao_7, Card_Tiao_8, Card_Tiao_9)
	// //
	// this.AllCard = append(this.AllCard, Card_Bing_1, Card_Wan_1, Card_Wan_1, Card_Wan_1, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	// this.AllCard = append(this.AllCard, Card_Tiao_3, Card_Tiao_3, Card_Tiao_4, Card_Tiao_4, Card_Tiao_4)
	this.SendCardCount = 0
	this.SendGangCardCount = 0
	this.BaoPos = Card_Rear
}

func (this *CardManager108) SendStartCard(shouPai *[]byte) {
	*shouPai = append(*shouPai, this.AllCard[this.SendCardCount:this.SendCardCount+13]...)
	this.SendCardCount += 13
	Sort(*shouPai)
}

func (this *CardManager108) SendCard(gang bool) byte {
	c := byte(Card_Invalid)
	if this.SendCardCount+this.SendGangCardCount == 108 {
		return c
	}
	if !gang {
		c = this.AllCard[this.SendCardCount]
		this.SendCardCount++
		return c
	} else {
		logs.Debug("牌数：", len(this.AllCard), this.SendGangCardCount)
		c = this.AllCard[108-this.SendGangCardCount-1]
		this.SendGangCardCount++
		return c
	}
}

func (this *CardManager108) GetLeftCardCount() int {
	if this.SendCardCount+this.SendGangCardCount > 108 {
		return Card_Invalid
	} else {
		return 108 - (this.SendCardCount + this.SendGangCardCount)
	}
}

package main

import (
	. "MaJiangTool"
)

// //初始化展示：其他都是具体项目具体定义
// type ExtPlayer struct {
// 	Player
// }

//有效动作
type ValidAction struct {
	Style   int
	Card    byte
	HuStyle []int //如果是胡，胡的类型，如果不是就忽略
}

type ExtPlayer struct {
	Player
	//以下是用户自己定义的变量
	HandCard     []byte         // 手牌
	HuanPais     []byte         //换牌
	QueColor     byte           //定缺的颜色
	FuZis        []FuZi         //已有的附子
	OutCards     []byte         //出牌牌池
	ValidActions []*ValidAction //有效动作
	HuType       []int          //胡牌的类型
	LunHuEd      bool           //本轮是否胡过
	AcEvent      *ActionEvent   // 未验证动作
	GiveUp       bool           //是否放弃
	TuoGuan      bool           //是否托管
	RateCoins    float64        //手续费
	Double       int
}

func (this *ExtPlayer) ContainValidAction(style int) bool {
	for _, v := range this.ValidActions {
		if v.Style == style {
			return true
		}
	}
	return false
}

func (this *ExtPlayer) ContainValidActionMul(styles []int) bool {
	for _, s := range styles {
		for _, v := range this.ValidActions {
			if s == v.Style {
				return true
			}
		}
	}
	return false
}

func (this *ExtPlayer) AddVaildAcion(style int, card byte) {
	this.ValidActions = append(this.ValidActions, &ValidAction{
		Style: style,
		Card:  card,
	})
}

func (this *ExtPlayer) GetHuTypes() []int {
	for _, v := range this.ValidActions {
		if v.Style == ActionType_Hu {
			return v.HuStyle
		}
	}
	return []int{}
}

func (this *ExtPlayer) GetLastCard() byte {
	return this.HandCard[len(this.HandCard)-1]
}

// 设置手牌,手牌由大到小排序
func (this *ExtPlayer) SetHandCard(c []byte) {
	this.HandCard = c
}

func (this *ExtPlayer) HaveQueColor() bool {
	for _, v := range this.HandCard {
		if GetCardColor(v) == this.QueColor {
			return true
		}
	}
	return false
}

func (this *ExtPlayer) GetOperateAction() int {
	if this.AcEvent != nil {
		return this.AcEvent.Fu.WeaveKind
	} else {
		return ActionType_None
	}
}

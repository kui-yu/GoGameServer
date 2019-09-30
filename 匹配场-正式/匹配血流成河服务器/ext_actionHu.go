package main

import (
	. "MaJiangTool"
	// "logs"
)

type ExtActionHuPai struct {
	BaseAction
	Hu ActionHu
	//
	HuValid ValidAction
}

func (this *ExtActionHuPai) Init() {
	this.Style = ActionType_Hu
}

func (this *ExtActionHuPai) GetResult(ChairId int, SelfFuZi []FuZi, ShouPai []byte, LastEvent EventIe) bool {
	if !this.CheckCondition(ChairId, SelfFuZi, ShouPai, LastEvent) {
		return false
	}
	//向手牌中添加胡的牌
	vHand := append([]byte{}, ShouPai...)
	if LastEvent.GetStyle() == EventType_OutCard {
		e := LastEvent.(*OutCardEvent)
		vHand = append(vHand, e.Card)
	} else if LastEvent.GetStyle() == EventType_Action {
		e := LastEvent.(*ActionEvent)
		vHand = append(vHand, e.Fu.OperateCard)
	}
	//是否是清一色
	if CheckCardColorCount(SelfFuZi, ShouPai, nil) <= 1 {
		this.AddHuType(HuType_QingYiSe)
	}
	//
	duicards := append([]byte{}, vHand...)
	Sort(duicards)
	if Hu7Dui(ChairId, SelfFuZi, &duicards, LastEvent, this.Hui) {
		//，此时duicards不一定包含所有牌，有的会牌没有配对是不加入的
		//还有赖子就已经全部配对完了,表示肯定有4张一样的
		if len(duicards) != 14 {
			this.AddHuType(HuType_7Dui_HaoHua)
			return true
		}
		//还原会牌
		for i, v := range duicards {
			if v&Hui_Mask > 0 {
				duicards[i] = v & (^byte(Hui_Mask))
			}
		}
		//
		for _, c := range duicards {
			if CountSameCard(duicards, c) == 4 {
				this.AddHuType(HuType_7Dui_HaoHua)
				return true
			}
		}
		//
		this.AddHuType(HuType_7Dui)
		return true
	}
	//普通胡牌
	lstHu := [][]byte{}
	if this.Hu.GetResult(ChairId, SelfFuZi, ShouPai, LastEvent, &lstHu) {
		for _, lst := range lstHu {
			if CheckPiaoHu(SelfFuZi, lst) {
				this.AddHuType(HuType_PiaoHu)
				if CheckAllJiang(SelfFuZi, vHand, this.Hui) {
					this.AddHuType(HuType_JiangHu)
				}
				if len(lst) == 2 { //金钩钓
					if CheckShiBaLuoHan(SelfFuZi, vHand) {
						this.AddHuType(HuType_ShiBaLuoHan)
					} else {
						this.AddHuType(HuType_JinGouDiao)
					}
				}
			} else {
				if CheckYaoJiuHu(SelfFuZi, lst) {
					this.AddHuType(HuType_YaoJiuHu)
				} else {
					this.AddHuType(HuType_PingHu)
				}
			}
		}
		return true
	}
	return false
}

func (this *ExtActionHuPai) AddHuType(t int) {
	for _, v := range this.HuValid.HuStyle {
		if v == t {
			return
		}
	}
	this.HuValid.HuStyle = append(this.HuValid.HuStyle, t)
}

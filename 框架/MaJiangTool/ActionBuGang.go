package MaJiangTool

import (
	"logs"
)

type ActionPuBuGang struct {
	BaseAction
}

func (this *ActionPuBuGang) Init() {
	this.InitData(nil, ActionType_Gang_PuBuGang)
	this.Supper = this
}

func (this *ActionPuBuGang) GetResult(ChairId int, SelfFuZi []FuZi, ShouPai []byte, LastEvent EventIe, Out interface{}) bool {
	//检查触发条件
	if !this.CheckCondition(ChairId, SelfFuZi, ShouPai, LastEvent) {
		return false
	}
	//可补杠列表
	bucards := []byte{}
	for _, v := range SelfFuZi {
		if v.WeaveKind == ActionType_Peng {
			for _, sit := range ShouPai {
				if sit == v.CardData[0] {
					bucards = append(bucards, sit)
					break
				}
			}
		}
	}
	//
	if len(bucards) == 0 {
		return false
	}
	//
	if Out != nil {
		result, ok := Out.(*[]FuZi)
		for _, it := range bucards {
			if ok {
				gang := FuZi{}
				gang.WeaveKind = ActionType_Gang_PuBuGang
				gang.OperateCard = it
				gang.CardData = append(gang.CardData, it)
				gang.CardData = append(gang.CardData, it)
				gang.CardData = append(gang.CardData, it)
				gang.CardData = append(gang.CardData, it)
				*result = append(*result, gang)
			}
		}
	}
	//
	return true
}

func (this *ActionPuBuGang) ReNew(SelfFuZi *[]FuZi, ShouPai *[]byte, Event *ActionEvent, DoCheck bool) bool {
	fz := Event.GetActionFuZi()
	if ActionType_Gang_PuBuGang != fz.WeaveKind || fz.OperateCard == byte(Card_Invalid) {
		return false
	}
	for _, v := range *SelfFuZi {
		if ActionType_Gang_Ming == v.WeaveKind && fz.OperateCard == v.CardData[0] {
			return false
		}
	}
	//从手牌中删除操作牌，如果失败则操作失败
	Event.DelCard = []byte{fz.OperateCard}
	// dc := Event.GetDelCard()
	// *dc = append(*dc, fz.OperateCard)
	if !DeleteCard(ShouPai, Event.DelCard) {
		return false
	}
	//更新玩家附子列表和动作事件
	for i, fv := range *SelfFuZi {
		logs.Debug(".......明杠", ActionType_Peng, fv.WeaveKind, fz.OperateCard, fv.CardData[0])
		if (ActionType_Peng == fv.WeaveKind) && fz.OperateCard == fv.CardData[0] {
			(*SelfFuZi)[i].WeaveKind = ActionType_Gang_PuBuGang
			(*SelfFuZi)[i].CardData = append((*SelfFuZi)[i].CardData, fz.OperateCard)
			Event.GetActionFuZi().CardData = append(Event.GetActionFuZi().CardData, fz.OperateCard)
			(*SelfFuZi)[i].ProvideUser = Card_Rear
			Event.GetActionFuZi().ProvideUser = Card_Rear
			break
		}
	}
	return true
}

func (this *ActionPuBuGang) RollBack(ChairId int, SelfFuZi *[]FuZi, ShouPai []byte, Event *ActionEvent, DoCheck bool) bool {
	if ActionType_Gang_PuBuGang != Event.GetActionFuZi().WeaveKind ||
		Event.GetActionFuZi().OperateCard == Card_Invalid {
		return false
	}

	var minggang *FuZi
	for i, v := range *SelfFuZi {
		if ActionType_Gang_Ming == v.WeaveKind && Event.GetActionFuZi().OperateCard == v.OperateCard {
			minggang = &((*SelfFuZi)[i])
			break
		}
	}
	if minggang == nil {
		return false
	}
	minggang.WeaveKind = ActionType_Peng
	minggang.CardData = minggang.CardData[:3]
	return true
}

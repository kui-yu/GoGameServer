package MaJiangTool

import (
	"logs"
)

type ActionMingGang struct {
	BaseAction
}

func (this *ActionMingGang) Init(Hui HuiIe) {
	this.InitData(Hui, ActionType_Gang_Ming)
	this.Supper = this
}

func (this *ActionMingGang) GetResult(ChairId int, SelfFuZi []FuZi, ShouPai []byte, LastEvent EventIe, Out interface{}) bool {
	//检查触发条件
	if !this.CheckCondition(ChairId, SelfFuZi, ShouPai, LastEvent) {
		return false
	}
	//
	opcard := byte(Card_Invalid)
	if LastEvent.GetStyle() == EventType_OutCard {
		ev := LastEvent.(*OutCardEvent)
		opcard = ev.GetCard()
	} else if LastEvent.GetStyle() == EventType_Action {
		ev := LastEvent.(*ActionEvent)
		opcard = ev.GetActionFuZi().OperateCard
	}

	if opcard == Card_Invalid {
		return false
	}
	//
	handCard := append([]byte{}, ShouPai...)
	SignHuiToCardVector(this.Hui, handCard)
	//
	cnt := 0
	for _, hc := range handCard {
		if hc == opcard {
			cnt++
		}
	}
	if cnt < 3 {
		return false
	}

	if Out == nil {
		return true
	}
	//
	outPut := Out.(*[]FuZi)
	addItem := FuZi{}
	addItem.WeaveKind = ActionType_Gang_Ming
	for i := 0; i < 4; i++ {
		addItem.CardData = append(addItem.CardData, opcard)
	}
	addItem.OperateCard = opcard
	logs.Debug("明杠触发：", addItem)
	*outPut = append(*outPut, addItem)
	return true
}

func (this *ActionMingGang) ReNew(SelfFuZi *[]FuZi, ShouPai *[]byte, Event *ActionEvent, DoCheck bool) bool {
	reNewFuZi := Event.GetActionFuZi()
	if ActionType_Gang_Ming != reNewFuZi.WeaveKind {
		return false
	}
	gangCard := Event.GetActionFuZi().OperateCard
	for i := 0; i < 4; i++ {
		reNewFuZi.CardData = append(reNewFuZi.CardData, gangCard)

	}
	//
	Event.DelCard = []byte{gangCard, gangCard, gangCard}
	if !DeleteCard(ShouPai, Event.DelCard) {
		return false
	}
	*SelfFuZi = append(*SelfFuZi, *reNewFuZi)
	return true
}

func (this *ActionMingGang) RollBack(ChairId int, SelfFuZi *[]FuZi, ShouPai []byte, Event *ActionEvent, DoCheck bool) bool {
	return true
}

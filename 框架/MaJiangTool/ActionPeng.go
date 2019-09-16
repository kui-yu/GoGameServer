package MaJiangTool

// import (
// 	"logs"
// )

type ActionPeng struct {
	BaseAction
}

func (this *ActionPeng) Init(Hui HuiIe) {
	this.InitData(Hui, ActionType_Peng)
	this.Supper = this
}

func (this *ActionPeng) GetResult(ChairId int, SelfFuZi []FuZi, ShouPai []byte, LastEvent EventIe, Out interface{}) bool {
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

	if opcard == Card_Invalid || (this.Hui != nil && this.Hui.IsHui(opcard)) {
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
	if cnt < 2 {
		return false
	}

	if Out == nil {
		return true
	}
	//
	outPut := Out.(*[]FuZi)
	addItem := FuZi{}
	addItem.WeaveKind = ActionType_Peng
	for i := 0; i < 3; i++ {
		addItem.CardData = append(addItem.CardData, opcard)
	}
	addItem.OperateCard = opcard
	*outPut = append(*outPut, addItem)
	return true
}

func (this *ActionPeng) ReNew(SelfFuZi *[]FuZi, ShouPai *[]byte, Event *ActionEvent, DoCheck bool) bool {
	reNewFuZi := Event.GetActionFuZi()
	if ActionType_Peng != reNewFuZi.WeaveKind {
		return false
	}
	gangCard := Event.GetActionFuZi().OperateCard
	reNewFuZi.CardData = []byte{gangCard, gangCard, gangCard}
	//
	Event.DelCard = []byte{gangCard, gangCard}
	if !DeleteCard(ShouPai, Event.DelCard) {
		return false
	}
	*SelfFuZi = append(*SelfFuZi, *reNewFuZi)
	return true
}

func (this *ActionPeng) RollBack(ChairId int, SelfFuZi *[]FuZi, ShouPai []byte, Event *ActionEvent, DoCheck bool) bool {
	return true
}

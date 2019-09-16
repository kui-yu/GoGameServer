package MaJiangTool

type ActionAnGang struct {
	BaseAction
}

func (this *ActionAnGang) Init(Hui HuiIe) {
	this.InitData(Hui, ActionType_Gang_An)
	this.Supper = this
}

func (this *ActionAnGang) GetResult(ChairId int, SelfFuZi []FuZi, ShouPai []byte, LastEvent EventIe, Out interface{}) bool {
	//检查触发条件
	if !this.CheckCondition(ChairId, SelfFuZi, ShouPai, LastEvent) {
		return false
	}
	//把手牌遍历出数量
	cardCount := [34]int{}
	for _, v := range ShouPai {
		if this.Hui != nil && this.Hui.IsHui(v) {
			continue
		}
		cardCount[SwitchToCardIndex(v)]++
	}

	anCards := []byte{}
	for i, v := range cardCount {
		if v == 4 {
			anCards = append(anCards, SwitchToCardData(i))
		}
	}

	if len(anCards) == 0 {
		return false
	}

	if Out == nil {
		return true
	}

	outPut := Out.(*[]FuZi)
	for _, ac := range anCards {
		addItem := FuZi{}
		addItem.WeaveKind = ActionType_Gang_An
		addItem.CardData = []byte{ac, ac, ac, ac}
		addItem.OperateCard = ac
		*outPut = append(*outPut, addItem)
	}
	return true
}

func (this *ActionAnGang) ReNew(SelfFuZi *[]FuZi, ShouPai *[]byte, Event *ActionEvent, DoCheck bool) bool {

	reNewFuZi := Event.GetActionFuZi()
	if ActionType_Gang_An != reNewFuZi.WeaveKind {
		return false
	}
	gangCard := Event.GetActionFuZi().OperateCard
	for i := 0; i < 4; i++ {
		reNewFuZi.CardData = append(reNewFuZi.CardData, gangCard)

	}
	//
	Event.DelCard = []byte{gangCard, gangCard, gangCard, gangCard}
	if !DeleteCard(ShouPai, Event.DelCard) {
		return false
	}
	*SelfFuZi = append(*SelfFuZi, *reNewFuZi)
	return true
}

func (this *ActionAnGang) RollBack(ChairId int, SelfFuZi *[]FuZi, ShouPai []byte, Event *ActionEvent, DoCheck bool) bool {
	return true
}

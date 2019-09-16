package MaJiangTool

const (
	EventType_SendCard = iota
	EventType_OutCard
	EventType_Action
	EventType_BaoPai
)

type BaseEvent struct {
	Style   int //事件类型
	ChairId int //座位号
}

func (this *BaseEvent) GetStyle() int {
	return this.Style
}

func (this *BaseEvent) GetChairId() int {
	return this.ChairId
}

//发牌事件
type SendCardEvent struct {
	BaseEvent
	Card   byte //牌的牌值
	BuGang bool //补杠发牌
}

func (this *SendCardEvent) GetCard() byte {
	return this.Card
}

//出牌事件
type OutCardEvent struct {
	BaseEvent
	Card     byte //牌的牌值
	OverTime bool
}

func (this *OutCardEvent) GetCard() byte {
	return this.Card
}

func (this *OutCardEvent) IsOverTime() bool {
	return this.OverTime
}

//
type ActionEvent struct {
	BaseEvent
	Fu      FuZi
	DelCard []byte
}

func (this *ActionEvent) Init(ChairId int) {
	this.Style = EventType_Action
	this.ChairId = ChairId
}

func (this *ActionEvent) GetActionFuZi() *FuZi {
	return &this.Fu
}

func (this *ActionEvent) GetCard() byte {
	return this.Fu.OperateCard
}

func (this *ActionEvent) GetDelCard() *[]byte {
	return &this.DelCard
}

/////////////////////////////////////////////
func CreateSendCardEvent(chairId int, sendCard byte, buGang bool) *SendCardEvent {
	ne := &SendCardEvent{
		Card:   sendCard,
		BuGang: buGang,
	}
	ne.ChairId = chairId
	ne.Style = EventType_SendCard
	return ne
}

func CreateOutCardEvent(chairId int, outCard byte, overTime bool) *OutCardEvent {
	ne := &OutCardEvent{
		Card:     outCard,
		OverTime: overTime,
	}
	ne.ChairId = chairId
	ne.Style = EventType_OutCard
	return ne
}

func CreateActionEvent(chairId int) *ActionEvent {
	ne := &ActionEvent{}
	ne.ChairId = chairId
	ne.Style = EventType_Action
	return ne
}

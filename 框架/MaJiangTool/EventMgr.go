package MaJiangTool

type EventMgr struct {
	Events []EventIe
}

func (this *EventMgr) Reset() {
	this.Events = []EventIe{}
}

func (this *EventMgr) AddEvent(e EventIe) {
	this.Events = append(this.Events, e)
}

//返回最后一个事件
func (this *EventMgr) GetLastEvent() EventIe {
	if len(this.Events) == 0 {
		return nil
	}
	return this.Events[len(this.Events)-1]
}

//返回倒数第几个事件
func (this *EventMgr) GetBackEvent(id int) EventIe {
	if id == 0 || id > len(this.Events) {
		return nil
	}
	return this.Events[len(this.Events)-id]
}

func (this *EventMgr) GetLastOutCardEvent() EventIe {
	for _, v := range this.Events {
		if v.GetStyle() == EventType_OutCard {
			return v
		}
	}
	return nil
}

func (this *EventMgr) GetLastEventChairId() int {
	if len(this.Events) == 0 {
		return -1
	}
	return this.Events[len(this.Events)-1].GetChairId()
}

func (this *EventMgr) GetEventCount() int {
	return len(this.Events)
}

//碰先不碰后规则
//玩家这轮摸牌之后到结束，可以碰的牌是否出现过
func (this *EventMgr) IsFirstOutThisCardByLun(chairId int, outCard byte) bool {
	cnt := 0
	for _, v := range this.Events {
		if v.GetChairId() == chairId && v.GetStyle() == EventType_SendCard {
			break
		}
		//
		if EventType_OutCard == v.GetStyle() && v.GetCard() == outCard {
			cnt++
		}
	}
	if cnt > 1 {
		return false
	}
	return true
}

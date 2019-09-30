package main

import (
	// "github.com/astaxie/beego/logs"
	. "MaJiangTool"
	"encoding/json"
	"logs"
)

func (this *ExtDesk) HandleGameOutCard(p *ExtPlayer, d *DkInMsg) {
	if int(p.ChairId) != this.CurCid || this.GameState != GAME_STATE_PLAY {
		logs.Debug("出牌1", p.ChairId, this.CurCid)
		return
	}
	//
	if p.GiveUp {
		return
	}
	//
	data := GOutCard{}
	json.Unmarshal([]byte(d.Data), &data)
	outcard := byte(data.Card)
	//检查手牌中是否有这这张牌
	exist := false
	for _, c := range p.HandCard {
		if c == outcard {
			exist = true
			break
		}
	}
	if !exist {
		return
	}
	//如果已经胡牌了，出的不是最后一张
	if len(p.HuType) > 0 && outcard != p.HandCard[len(p.HandCard)-1] {
		return
	}
	//关闭定时器
	this.ClearTimer()
	//更新玩家数据(出牌牌池，删除手牌), 根据定缺排序
	p.OutCards = append(p.OutCards, outcard)
	DeleteCard(&p.HandCard, []byte{outcard})
	SortByQue(&p.HandCard, p.QueColor)
	//创建玩家出牌事件
	outEvent := CreateOutCardEvent(int(p.ChairId), outcard, false)
	this.EventManager.AddEvent(outEvent)
	//
	this.UpdateValidAction_AfterOutCardEvent(outEvent)
	//是否触发其他玩家操作
	this.NextProcessAfterSuccessEvent(outEvent)
}

//进入出牌阶段
func (this *ExtDesk) TimerOutCard(d interface{}) {
	this.GameState = GAME_STATE_PLAY
	this.BroadStageTime(TIMER_OUTCARD_NUM)
	//更新当前玩家动作，发送给玩家
	card := this.Players[this.CurCid].GetLastCard()
	sevent := CreateSendCardEvent(this.CurCid, card, false)
	this.EventManager.AddEvent(sevent)
	this.UpdateValidAction_AfterSendCardEvent(sevent)
	//开启定缺托管定时器
	this.NextProcessAfterSuccessEvent(sevent)

}

func (this *ExtDesk) GenSendAction(cid int) []HaveAction {
	hav := []HaveAction{}
	user := this.Players[cid]
	for _, v := range user.ValidActions {
		hav = append(hav, HaveAction{
			Style:   v.Style,
			Card:    int(v.Card),
			HuTypes: append([]int{}, v.HuStyle...),
		})
	}
	return hav
}

/////////////////////////////////////////////////////
func (this *ExtDesk) UpdateValidAction_AfterSendCardEvent(event *SendCardEvent) {
	//清空所有玩家身上的动作
	this.ClearActionReply()
	//
	user := this.Players[event.ChairId]
	//
	if user.GiveUp {
		return
	}
	//
	for _, ac := range this.ActionsAfterSendCard {
		fus := []FuZi{}
		if ac.GetResult(event.ChairId, user.FuZis, user.HandCard, event, &fus) {
			for _, f := range fus {
				if GetCardColor(f.OperateCard) == byte(user.QueColor) {
					continue
				} else {
					user.AddVaildAcion(f.WeaveKind, f.OperateCard)
				}
			}
		}
	}
	//
	if !user.HaveQueColor() {
		//更新触发胡牌动作标志位
		if this.HuPai.GetResult(event.ChairId, user.FuZis, user.HandCard, event) {
			this.HuPai.HuValid.Style = ActionType_Hu
			this.HuPai.HuValid.Card = event.GetCard()
			va := this.HuPai.HuValid
			user.ValidActions = append(user.ValidActions, &va)
			this.HuPai.HuValid = ValidAction{}
		}
	}

	//发送玩家有的动作
	acs := this.GenSendAction(int(user.ChairId))
	if len(acs) != 0 {
		user.SendNativeMsg(MSG_GAME_INFO_HAVEACTION_NOTIFY, &GHaveActionNotify{
			Id:   MSG_GAME_INFO_HAVEACTION_NOTIFY,
			Data: acs,
		})
	}
}

func (this *ExtDesk) UpdateValidAction_AfterOutCardEvent(event *OutCardEvent) {
	this.ClearActionReply()
	//
	for i := 1; i < len(this.Players); i++ {
		cid := (event.ChairId + i) % len(this.Players)
		p := this.Players[cid]
		//放弃或者是缺的颜色就不做处理
		if p.GiveUp || GetCardColor(event.Card) == byte(p.QueColor) {
			continue
		}
		//更新触发动作标志位
		for _, ac := range this.ActionsAfterOutCard {
			fus := []FuZi{}
			if ac.GetResult(cid, p.FuZis, p.HandCard, event, &fus) {
				for _, f := range fus {
					if GetCardColor(f.OperateCard) == byte(p.QueColor) {
						continue
					} else {
						p.AddVaildAcion(f.WeaveKind, f.OperateCard)
					}
				}
			}
		}
		if !p.HaveQueColor() {
			if this.HuPai.GetResult(int(p.ChairId), p.FuZis, p.HandCard, event) {
				this.HuPai.HuValid.Style = ActionType_Hu
				this.HuPai.HuValid.Card = event.GetCard()
				va := this.HuPai.HuValid
				p.ValidActions = append(p.ValidActions, &va)
				this.HuPai.HuValid = ValidAction{}
			}
		}

		//广播触发动作
		acs := this.GenSendAction(int(p.ChairId))
		if len(acs) != 0 {
			p.SendNativeMsg(MSG_GAME_INFO_HAVEACTION_NOTIFY, &GHaveActionNotify{
				Id:   MSG_GAME_INFO_HAVEACTION_NOTIFY,
				Data: acs,
			})
		}
	}
}

//抢杠胡
func (this *ExtDesk) UpdateValidAction_AfterActionEvent_Other(event *ActionEvent) {
	this.ClearActionReply()
	//
	if event.GetActionFuZi().WeaveKind != ActionType_Gang_PuBuGang {
		return
	}
	//
	for i := 1; i < len(this.Players); i++ {
		cid := (event.ChairId + i) % len(this.Players)
		p := this.Players[cid]
		//放弃或者是缺的颜色就不做处理
		if p.GiveUp || GetCardColor(event.Fu.OperateCard) == byte(p.QueColor) {
			continue
		}
		for _, ac := range this.ActionsAfterActionOther {
			fus := []FuZi{}
			if ac.GetResult(cid, p.FuZis, p.HandCard, event, &fus) {
				for _, f := range fus {
					if GetCardColor(f.OperateCard) == p.QueColor {
						continue
					} else {
						p.AddVaildAcion(f.WeaveKind, f.OperateCard)
					}
				}
			}
		}
		//
		if !p.HaveQueColor() {
			//
			if this.HuPai.GetResult(int(p.ChairId), p.FuZis, p.HandCard, event) {
				this.HuPai.HuValid.Style = ActionType_Hu
				this.HuPai.HuValid.Card = event.GetCard()
				va := this.HuPai.HuValid
				p.ValidActions = append(p.ValidActions, &va)
				this.HuPai.HuValid = ValidAction{}
			}
		}

		//广播触发动作
		acs := this.GenSendAction(int(p.ChairId))
		if len(acs) != 0 {
			p.SendNativeMsg(MSG_GAME_INFO_HAVEACTION_NOTIFY, &GHaveActionNotify{
				Id:   MSG_GAME_INFO_HAVEACTION_NOTIFY,
				Data: acs,
			})
		}
	}
}

func (this *ExtDesk) UpdateValidAction_AfterActionEvent_Self(event *ActionEvent) {
	//清空所有玩家身上的动作
	this.ClearActionReply()
	//
	user := this.Players[event.ChairId]
	//
	if user.GiveUp {
		return
	}
	//
	for _, ac := range this.ActionsAfterActionSelf {
		fus := []FuZi{}
		if ac.GetResult(event.ChairId, user.FuZis, user.HandCard, event, &fus) {
			for _, f := range fus {
				if GetCardColor(f.OperateCard) == user.QueColor {
					continue
				} else {
					user.AddVaildAcion(f.WeaveKind, f.OperateCard)
				}
			}
		}
	}
	//
	//发送玩家有的动作
	acs := this.GenSendAction(int(user.ChairId))
	if len(acs) != 0 {
		user.SendNativeMsg(MSG_GAME_INFO_HAVEACTION_NOTIFY, &GHaveActionNotify{
			Id:   MSG_GAME_INFO_HAVEACTION_NOTIFY,
			Data: acs,
		})
	}
}

func (this *ExtDesk) ClearActionReply() {
	for _, v := range this.Players {
		v.ValidActions = []*ValidAction{}
		v.AcEvent = nil
		v.LunHuEd = false
	}
}

func (this *ExtDesk) GetHaveActionCount() int {
	cnt := 0
	for _, v := range this.Players {
		if len(v.ValidActions) > 0 {
			cnt++
		}
	}
	return cnt
}

func (this *ExtDesk) NextProcessAfterSuccessEvent(event EventIe) {
	user := this.Players[event.GetChairId()]
	//发牌成功后直接开启出牌定时器
	this.ClearTimer()
	if event.GetStyle() == EventType_SendCard {
		this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCardDo, nil)
	} else if event.GetStyle() == EventType_OutCard {
		//出牌成功后，有动作等待动作，没动作发牌
		logs.Debug("出牌,动作,牌", user.ChairId, this.GetHaveActionCount(), event.GetCard())
		if this.GetHaveActionCount() > 0 {
			this.AddTimer(TIMER_ACTION, TIMER_ACTION_NUM, this.TimerActionReplyDo, nil)
		} else {
			this.CurCid = (this.CurCid + 1) % len(this.Players)
			this.SendCard(this.CurCid, false)
		}
	} else if event.GetStyle() == EventType_Action {
		//动作成功后，如果是杠发牌给自己，如果是碰则出牌或者继续动作,如果还有动作则等待
		logs.Debug("玩家状态：", user.ChairId, user.HandCard, user.FuZis)
		this.CurCid = int(user.ChairId)
		if this.GetHaveActionCount() == 0 {
			actionevent := event.(*ActionEvent)
			if HandCardSize(user.HandCard, user.FuZis)%3 == 1 {
				this.SendCard(event.GetChairId(), this.NeedBuGang(actionevent))
			} else {
				this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCardDo, nil)
			}
		} else {
			this.AddTimer(TIMER_ACTION, TIMER_ACTION_NUM, this.TimerActionReplyDo, nil)
		}
	}
	//打印当前玩家的手牌和附子

}

func (this *ExtDesk) NeedBuGang(event *ActionEvent) bool {
	if event.Fu.WeaveKind == ActionType_Gang_Ming ||
		event.Fu.WeaveKind == ActionType_Gang_An ||
		event.Fu.WeaveKind == ActionType_Gang_PuBuGang {
		return true
	}
	return false
}

func (this *ExtDesk) SendCard(chairId int, gang bool) {
	//游戏是否结束
	if this.CardMgr.GetLeftCardCount() == 0 {
		//游戏结束处理
		this.GameOver()
		return
	}
	user := this.Players[chairId]
	nextId := -1
	for i := 0; i < len(this.Players); i++ {
		user = this.Players[(chairId+i)%len(this.Players)]
		if user.GiveUp {
			continue
		}
		nextId = int(chairId + i)
		break
	}
	if nextId == -1 {
		return
	}
	this.CurCid = nextId
	//清空定时器
	//发牌
	nc := this.CardMgr.SendCard(gang)
	logs.Debug("发牌给玩家:", this.CurCid, nc)
	user.HandCard = append(user.HandCard, nc)
	//创建系统发牌事件
	sendEvent := CreateSendCardEvent(int(user.ChairId), nc, gang)
	this.EventManager.AddEvent(sendEvent)
	//广播发牌
	sd := &GSendCardNofify{
		Id:   MSG_GAME_INFO_SENDCARD_NOTIFY,
		Cid:  user.ChairId,
		Card: int(nc),
		Gang: gang,
	}
	for _, v := range this.Players {
		if v.ChairId != user.ChairId {
			sd.Card = -1
		} else {
			sd.Card = int(nc)
		}
		v.SendNativeMsg(MSG_GAME_INFO_SENDCARD_NOTIFY, sd)
	}
	//更新发牌玩家动作标志位
	this.UpdateValidAction_AfterSendCardEvent(sendEvent)
	//
	this.NextProcessAfterSuccessEvent(sendEvent)
}

func (this *ExtDesk) TimerOutCardDo(d interface{}) {
	user := this.Players[this.CurCid]
	out := user.HandCard[len(user.HandCard)-1]
	user.HandCard = user.HandCard[:len(user.HandCard)-1]
	req := GOutCard{
		Id:   MSG_GAME_INFO_OUTCARD,
		Card: int(out),
	}
	js, _ := json.Marshal(&req)
	InMsg := DkInMsg{
		Id:   MSG_GAME_INFO_OUTCARD,
		Data: string(js),
	}
	this.HandleGameOutCard(user, &InMsg)
}

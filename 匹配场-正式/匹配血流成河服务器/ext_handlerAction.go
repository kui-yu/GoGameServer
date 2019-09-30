package main

import (
	// "github.com/astaxie/beego/logs"
	. "MaJiangTool"
	"encoding/json"
	"logs"
)

func (this *ExtDesk) HandleAction(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("玩家操作动作1：", p.ChairId)
	if this.GameState != GAME_STATE_PLAY {
		return
	}
	//
	if p.GiveUp || len(p.ValidActions) == 0 {
		return
	}
	//
	data := GAction{}
	json.Unmarshal([]byte(d.Data), &data)
	card := byte(data.Card)
	//没有这个动作
	logs.Debug("玩家操作动作2：", p.ChairId, data.Style)
	if data.Style != ActionType_FangQi && !p.ContainValidAction(data.Style) {
		return
	}
	//已经有动作过
	if data.Style != ActionType_FangQi && p.GetOperateAction() != ActionType_None {
		return
	}
	//
	if data.Style != ActionType_FangQi {
		//操作牌非法，碰完之后还可以杠，所以chairid不能为自己才能判断
		if (this.EventManager.GetLastEvent().GetChairId() != int(p.ChairId)) &&
			(card != this.EventManager.GetLastEvent().GetCard()) {
			return
		}
	} else {
		//轮到自己出牌不能放弃
		if HandCardSize(p.HandCard, p.FuZis)%3 == 2 {
			return
		}
	}
	//创建玩家动作事件
	event := CreateActionEvent(int(p.ChairId))
	event.Fu.ProvideUser = this.EventManager.GetLastEvent().GetChairId()
	event.Fu.WeaveKind = data.Style
	event.Fu.OperateCard = card
	//缓存下动作
	p.AcEvent = event
	//
	logs.Debug("玩家操作动作3：")
	this.ProcessPriorityAction(this.EventManager.GetLastEvent(), int(p.ChairId))
}

//确定动作最高优先级
func (this *ExtDesk) ProcessPriorityAction(event EventIe, chairId int) {
	needWaitOther := false //是否需要等待其他玩家动作（胡>杠>碰>吃）,有多个胡也要等待
	priorityChairId := -1  //最高动作优先级的玩家
	//胡牌优先级最高
	user := this.Players[chairId]
	//
	if user.AcEvent.Fu.WeaveKind == ActionType_Hu {
		//清空自己的动作,设置胡牌类型，然后判断还有没有其他胡，有则等待，没有就发送中途结算，然后再发牌
		user.HuType = user.GetHuTypes()
		user.ValidActions = []*ValidAction{}
		//清空其他玩家不是胡的动作
		for i := 0; i < len(this.Players); i++ {
			if !this.Players[i].ContainValidAction(ActionType_Hu) {
				this.Players[i].ValidActions = []*ValidAction{}
			}
		}
		user.LunHuEd = true
		//如果还有胡就等待
		if this.GetPendingActionCount() > 0 {
			//发送动作信息
			this.BroadcastAll(MSG_GAME_INFO_ACTION_NOTIFY, &GActionDoNotify{
				Id:         MSG_GAME_INFO_ACTION_NOTIFY,
				Cid:        chairId,
				ActionType: ActionType_Hu,
				Cards:      []int{},
			})
			return
		}
		//如果胡都操作了。那么发送结算信息，然后继续发牌
		logs.Debug("胡牌,手牌，附子1:", user.HandCard, user.FuZis)
		this.GameHu()
		//如果上一个动作是普补杠，则还要回滚
		if event.GetStyle() == EventType_Action {
			ev := event.(*ActionEvent)
			this.RollBack(ev)
		} else if event.GetStyle() == EventType_SendCard {
			//删除发的牌
			DeleteCard(&user.HandCard, []byte{event.GetCard()})
		}
		//发牌
		this.ClearActionReply()
		this.SendCard((chairId+1)%len(this.Players), false)
		return
	}
	////////////
	//以下是不是胡的操作
	//
	//1.胡牌优先级,如果还有胡，等待
	priorityChairId = this.CheckPriority([]int{ActionType_Hu}, event.GetChairId(), &needWaitOther)
	if needWaitOther {
		return
	}
	//碰，明杠(没有胡的情况),是否需要等待其他人
	logs.Debug("玩家操作动作4：")
	if priorityChairId == -1 {
		priorityChairId = this.CheckPriority([]int{ActionType_Peng, ActionType_Gang_Ming}, event.GetChairId(), &needWaitOther)
		if needWaitOther {
			return
		}
	}

	//有效动作
	logs.Debug("玩家操作动作5：")
	if priorityChairId == -1 {
		//从上个动作的下家开始遍历
		for i := 1; i <= len(this.Players); i++ {
			cid := (event.GetChairId() + i) % len(this.Players)
			p := this.Players[cid]
			if p.GetOperateAction() != ActionType_None &&
				p.GetOperateAction() != ActionType_FangQi {
				priorityChairId = cid
				break
			}
		}
	}
	logs.Debug("玩家操作动作6：", priorityChairId)
	//确定最高优先级后
	if priorityChairId != -1 {
		proPlayer := this.Players[priorityChairId]
		proAction := this.Players[priorityChairId].AcEvent
		logs.Debug("生效动作", chairId, proAction.Fu.WeaveKind)
		if !this.ReNew(event, chairId, proAction) {
			logs.Debug("玩家操作动作7：")
			user.AcEvent.Fu.WeaveKind = ActionType_FangQi
			this.ProcessPriorityAction(event, chairId)
			return
		}
		//生效成功后
		if event.GetStyle() == EventType_OutCard {
			u := this.Players[event.GetChairId()]
			u.OutCards = u.OutCards[:len(u.OutCards)-1]
		} else if event.GetStyle() == EventType_Action {
			if event.GetChairId() != proAction.GetChairId() {
				ev := event.(*ActionEvent)
				this.RollBack(ev)
			}
			//判断是否抢杠胡
			if proAction.GetStyle() == ActionType_Hu &&
				proAction.GetChairId() != event.GetChairId() &&
				event.GetStyle() == ActionType_Gang_PuBuGang {
				//广播更新玩家附子
			}
		}
		this.EventManager.AddEvent(proAction)
		this.Players[priorityChairId].AcEvent = nil
		//注意，这边可能有多个胡
		if proAction.GetStyle() == ActionType_Hu {
			//如果胡都操作了。那么发送结算信息，然后继续发牌
			logs.Debug("胡牌,手牌，附子1:", proPlayer.HandCard, proPlayer.FuZis)
			this.GameHu()
			//
			if event.GetStyle() == EventType_SendCard {
				DeleteCard(&proPlayer.HandCard, []byte{proAction.GetCard()})
			}
			//发牌
			this.ClearActionReply()
			this.SendCard((chairId+1)%len(this.Players), false)
			return
		} else {
			//如果是普补杠，更新以下是否其他人可以抢杠胡
			this.UpdateValidAction_AfterActionEvent_Other(proAction)
			//自己动作完成后，比如碰后玩，自身还有暗杠就更新动作
			if this.GetHaveActionCount() == 0 {
				if HandCardSize(proPlayer.HandCard, proPlayer.FuZis)%3 != 1 {
					this.UpdateValidAction_AfterActionEvent_Self(proAction)
				}
			}
			//发送动作信息
			scards := []int{}
			for _, sc := range proAction.DelCard {
				scards = append(scards, int(sc))
			}
			logs.Debug("通知玩家操作哪个动作：", priorityChairId, proAction.Fu.WeaveKind)
			this.BroadcastAll(MSG_GAME_INFO_ACTION_NOTIFY, &GActionDoNotify{
				Id:         MSG_GAME_INFO_ACTION_NOTIFY,
				Cid:        priorityChairId,
				ActionType: proAction.Fu.WeaveKind,
				Cards:      scards,
			})
			this.NextProcessAfterSuccessEvent(proAction)
		}
	} else { //放弃动作，其他人也没有动作的时候
		if this.GetPendingActionCount() != 0 {
			return
		}
		this.ClearActionReply()
		this.ClearTimer()
		this.NextProcessAfterSuccessEvent(event)
	}
}

//回滚
func (this *ExtDesk) RollBack(event *ActionEvent) bool {
	p := this.Players[event.ChairId]
	for _, f := range this.ActionsAfterSendCard {
		if f.RollBack(event.ChairId, &p.FuZis, p.HandCard, event, false) {
			return true
		}
	}
	for _, f := range this.ActionsAfterActionSelf {
		if f.RollBack(event.ChairId, &p.FuZis, p.HandCard, event, false) {
			return true
		}
	}
	return false
}

//有动作未操作的玩家数量
func (this *ExtDesk) GetPendingActionCount() int {
	cnt := 0
	for _, p := range this.Players {
		if len(p.ValidActions) > 0 && p.GetOperateAction() == ActionType_None {
			cnt++
		}
	}
	return cnt
}

func (this *ExtDesk) ReNew(lastEvent EventIe, chairId int, event *ActionEvent) bool {
	user := this.Players[chairId]
	if lastEvent.GetStyle() == EventType_SendCard {
		for _, v := range this.ActionsAfterSendCard {
			if v.ReNew(&user.FuZis, &user.HandCard, event, true) {
				return true
			}
		}
	} else if lastEvent.GetStyle() == EventType_OutCard {
		for _, v := range this.ActionsAfterOutCard {
			if v.ReNew(&user.FuZis, &user.HandCard, event, true) {
				return true
			}
		}
	} else if lastEvent.GetStyle() == EventType_Action {
		if lastEvent.GetChairId() == event.GetChairId() {
			for _, v := range this.ActionsAfterActionSelf {
				if v.ReNew(&user.FuZis, &user.HandCard, event, true) {
					return true
				}
			}
		} else {
			for _, v := range this.ActionsAfterActionOther {
				if v.ReNew(&user.FuZis, &user.HandCard, event, true) {
					return true
				}
			}
		}
	}
	//如果上面都没满足，那就只能是胡
	if event.GetStyle() != ActionType_Hu {
		return false
	}
	return false
	// if this.HuPai.ReNew(&user.FuZis, &user.HandCard, event, true) {
	// 	return true
	// }
	// return false
}

//确定哪个玩家身上有最高优先级动作
func (this *ExtDesk) CheckPriority(checkType []int, startChairId int, needWait *bool) int {
	for i := 1; i <= len(this.Players); i++ {
		id := (startChairId + i) % len(this.Players)
		user := this.Players[id]
		if user.ContainValidActionMul(checkType) {
			if user.GetOperateAction() == ActionType_None {
				*needWait = true
				return -1
			} else if IntContain(checkType, user.GetOperateAction()) {
				return id
			}
		}
	}
	return -1
}

func (this *ExtDesk) TimerActionReplyDo(d interface{}) {
	if this.GetHaveActionCount() == 0 {
		return
	}
	//
	for _, v := range this.Players {
		if len(v.ValidActions) > 0 {
			msg := GAction{
				Id:    MSG_GAME_INFO_ACTION,
				Style: ActionType_FangQi,
				Card:  Card_Invalid,
			}
			js, _ := json.Marshal(&msg)
			InMsg := DkInMsg{
				Id:   MSG_GAME_INFO_ACTION,
				Data: string(js),
			}
			//执行动作
			this.HandleAction(v, &InMsg)
		}
	}
}

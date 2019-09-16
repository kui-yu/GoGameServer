package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) HandleGameCall(p *ExtPlayer, d *DkInMsg) {
	//判断是否叫分阶段
	if this.GameState != GAME_STATUS_CALL {
		logs.Error("叫分游戏状态错误:", this.GameState, GAME_STATUS_CALL)
		return
	}

	// logs.Debug("玩家叫分")
	data := GCallMsg{}
	json.Unmarshal([]byte(d.Data), &data)

	//判断是否玩家叫分
	if p.ChairId != this.CurCid || data.Coins == 0 {
		logs.Error("不是当前叫分用户或叫分为0:", p.ChairId, this.CurCid, data.Coins)
		return
	}
	//叫分不能小于最大叫分
	if data.Coins > 0 && data.Coins <= this.CallFen {
		logs.Error("叫分小于当前叫分:", data.Coins, this.CallFen)
		return
	}
	p.CFen = data.Coins
	if p.CFen > 0 {
		this.CallFen = p.CFen
		this.Banker = p.ChairId
	}
	this.DelTimer(TIMER_CALL)
	//是否还有下一个
	isover := false
	if p.CFen < 3 {
		nextid := (p.ChairId + 1) % int32(len(this.Players))
		nextPlayer := this.Players[nextid]
		if nextPlayer.CFen != 0 {
			isover = true
		} else {
			this.CurCid = nextid
			if p.TuoGuan {
				this.AddTimer(TIMER_CALL, 1, this.TimerCall, nil)
			} else {
				this.AddTimer(TIMER_CALL, TIMER_CALL_NUM, this.TimerCall, nil)
			}

		}
	} else {
		isover = true
	}

	//
	re := GCallMsgReply{
		Id:    MSG_GAME_INFO_CALL_REPLY,
		Cid:   p.ChairId,
		Coins: data.Coins,
	}
	if isover {
		re.End = 1
	} else {
		re.End = 0
	}
	this.BroadcastAll(MSG_GAME_INFO_CALL_REPLY, &re)
	if !isover {
		return
	}
	//判断所有人是不是都不叫。如果都不叫就重新进去发牌阶段
	allnocall := true
	for _, v := range this.Players {
		if v.CFen != -1 {
			allnocall = false
			break
		}
	}
	if allnocall {
		for _, v := range this.Players {
			v.CFen = 0
		}
		this.GameState = GAME_STATUS_START
		this.CallTimes++
		if this.CallTimes >= 3 {
			this.GameOverByNoCall(p)
		} else {
			this.BroadStageTime(TIMER_SENDCARD_NUM)
			this.AddTimer(TIMER_SENDCARD, TIMER_SENDCARD_NUM, this.TimerSendCard, nil)
		}
		return
	}
	//设置倍数
	for _, v := range this.Players {
		if v.ChairId == this.Banker {
			v.Double = this.CallFen * 2
		} else {
			v.Double = this.CallFen
		}
	}
	//
	banker := this.Players[this.Banker]
	banker.HandCard = append(banker.HandCard, this.DiPai...)
	this.DiPaiDoulbe = int32(this.CalDiPaiDouble(this.DiPai))
	//都叫分过了。定庄
	this.GameState = GAME_STATUS_PLAY
	//出牌阶段消息
	this.BroadStageTime(TIMER_OUTCARD_NUM)

	this.CurCid = banker.ChairId
	notify := GBankerNotify{
		Id:     MSG_GAME_INFO_BANKER_NOTIFY,
		Banker: banker.ChairId,
		DiPai:  this.DiPai,
		Double: this.DiPaiDoulbe,
	}
	this.BroadcastAll(MSG_GAME_INFO_BANKER_NOTIFY, &notify)

	//添加定时器，进入出牌阶段
	nextplayer := this.Players[this.CurCid]
	if nextplayer.TuoGuan {
		this.AddTimer(TIMER_OUTCARD, 1, this.TimerOutCard, nil)
	} else {
		this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCard, nil)
	}
}

func (this *ExtDesk) TimerCall(d interface{}) {
	data := GCallMsg{
		Coins: -1,
	}
	p := this.Players[this.CurCid]
	dv, _ := json.Marshal(data)
	this.HandleGameCall(p, &DkInMsg{
		Uid:  p.Uid,
		Data: string(dv),
	})
}

//计算底牌倍数，暂时不用
func (this *ExtDesk) CalDiPaiDouble(cards []byte) int {
	vdipai := Sort(cards)
	if vdipai[0] == 0x42 && vdipai[1] == 0x41 {
		return 4
	} else if vdipai[0] == 0x42 {
		return 2
	} else if vdipai[0] == 0x41 {
		return 2
	} else if GetLogicValue(vdipai[0]) == GetLogicValue(vdipai[1]) &&
		GetLogicValue(vdipai[0]) == GetLogicValue(vdipai[2]) {
		return 4
	} else if GetCardColor(vdipai[0]) == GetCardColor(vdipai[1]) &&
		GetCardColor(vdipai[0]) == GetCardColor(vdipai[2]) {
		if this.IsShunZiType(vdipai) {
			return 4
		} else {
			return 3
		}
	} else if this.IsShunZiType(vdipai) {
		return 3
	}
	return 1
}

func (this *ExtDesk) IsShunZiType(cards []byte) bool {
	for i := 0; i < len(cards)-1; i++ {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+1])+1 {
			return false
		}
	}
	return true
}

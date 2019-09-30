package main

import (
	. "MaJiangTool"
	"encoding/json"
	"math/rand"
)

// "encoding/json"
// "logs"
//	"sort"

func (this *ExtDesk) HandleDingQue(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != GAME_STATE_DINGQUE {
		return
	}
	//
	req := GDingQue{}
	json.Unmarshal([]byte(d.Data), &req)
	if req.Color < 0 || req.Color > 3 {
		return
	}
	//
	if p.QueColor != 255 { //已经操作定缺
		return
	}
	p.QueColor = byte(req.Color)
	//广播哪个玩家操作了定缺
	this.BroadcastAll(MSG_GAME_INFO_DINGQUE_NOTIFY, &GDingQueNotify{
		Id:  MSG_GAME_INFO_DINGQUE_NOTIFY,
		Cid: int(p.ChairId),
	})
	//判断所有玩家是否都操作结束了
	for _, v := range this.Players {
		if v.QueColor == 255 {
			return
		}
	}
	//以下是操作结束
	//广播进入出牌阶段
	this.TimerOutCard(nil)
}

func (this *ExtDesk) TimerDingQue(d interface{}) {
	this.GameState = GAME_STATE_DINGQUE
	this.BroadStageTime(TIMER_DINGQUE_NUM)
	//开启定缺托管定时器
	this.AddTimer(TIMER_DINGQUE, TIMER_DINGQUE_NUM, this.TimerDingQueDo, nil)
}

func (this *ExtDesk) TimerDingQueDo(d interface{}) {
	for _, v := range this.Players {
		if v.QueColor == 255 {
			v.QueColor = byte(CARD_COLOR_Wan + rand.Intn(10)%3)
			msg := GDingQue{
				Id:    MSG_GAME_INFO_DINGQUE,
				Color: int(v.QueColor),
			}
			js, _ := json.Marshal(&msg)
			InMsg := DkInMsg{
				Id:   MSG_GAME_INFO_DINGQUE,
				Data: string(js),
			}
			this.HandleDingQue(v, &InMsg)
		}
	}
}

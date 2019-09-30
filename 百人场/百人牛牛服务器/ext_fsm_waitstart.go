/**
* 等待开始状态
**/
package main

import (
	"encoding/json"
	"logs"
)

type FSMWaitStart struct {
	UpMark int
	Mark   int
	EDesk  *ExtDesk

	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMWaitStart) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMWaitStart) GetMark() int {
	return this.Mark
}

func (this *FSMWaitStart) Run(upMark int) {
	DebugLog("游戏状态：等待开始")
	logs.Debug("游戏状态：等待开始")
	this.UpMark = upMark

	timeId := gameConfig.GameStatusTimer.WaitstartId
	timeMs := int64(gameConfig.GameStatusTimer.WaitstartMS)

	this.EndDateTime = GetTimeMS() + timeMs

	this.addListener()                          // 添加监听
	this.EDesk.SendGameState(this.Mark, timeMs) // 发送桌子状态

	this.EDesk.JuHao = GetJuHao()

	this.EDesk.SendNotice(MSG_GAME_NDESKCHANGE, &struct {
		Id    int
		JuHao string
	}{
		Id:    MSG_GAME_NDESKCHANGE,
		JuHao: this.EDesk.JuHao,
	}, true, nil)

	this.EDesk.AddTimer(timeId, int(timeMs)/1000, this.TimerCall, nil)
}

func (this *FSMWaitStart) Leave() {
	this.removeListener()
}

func (this *FSMWaitStart) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_SEATBET)
}

func (this *FSMWaitStart) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()

	return remainTimeMS
}

// 添加网络监听
func (this *FSMWaitStart) addListener() {
	this.EDesk.Handle[MSG_GAME_QSEATDOWN] = this.recvRSeatDown
	this.EDesk.Handle[MSG_GAME_QSEATUP] = this.recvRSeatUp
	this.EDesk.Handle[MSG_GAME_QDOWNBET] = this.recvRSeatBet
}

// 接收到玩家下注
func (this *FSMWaitStart) recvRSeatBet(p *ExtPlayer, d *DkInMsg) {
	req := GClientQDownBet{}
	json.Unmarshal([]byte(d.Data), &req)
	this.EDesk.UserDownBet(p, req.SeatIdx, req.CoinIdx, true)
}

// 删除网络监听
func (this *FSMWaitStart) removeListener() {
	delete(this.EDesk.Handle, MSG_GAME_QSEATDOWN)
	delete(this.EDesk.Handle, MSG_GAME_QSEATUP)
}

func (this *FSMWaitStart) onUserOnline(p *ExtPlayer) {

}

func (this *FSMWaitStart) onUserOffline(p *ExtPlayer) {

}

// 接收到玩家请求坐下
func (this *FSMWaitStart) recvRSeatDown(p *ExtPlayer, d *DkInMsg) {
	req := GClientQSeatDown{}
	json.Unmarshal([]byte(d.Data), &req)

	this.EDesk.UserSeatDown(p, req.SeatIdx)
}

// 接收到玩家请求站立
func (this *FSMWaitStart) recvRSeatUp(p *ExtPlayer, d *DkInMsg) {
	this.EDesk.UserSeatUp(p, true, false)
}

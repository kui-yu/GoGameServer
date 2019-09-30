/**
* 等待玩家下注
**/
package main

import (
	"encoding/json"
	"logs"
)

type FSMDownBets struct {
	UpMark int
	Mark   int
	EDesk  *ExtDesk

	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMDownBets) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMDownBets) GetMark() int {
	return this.Mark
}

func (this *FSMDownBets) Run(upMark int) {
	DebugLog("游戏状态：所有玩家下注")
	logs.Debug("游戏状态：所有玩家下注")

	this.UpMark = upMark

	timeId := gameConfig.GameStatusTimer.DownBetsId
	timeMs := int64(gameConfig.GameStatusTimer.DownBetsMS)

	this.EndDateTime = GetTimeMS() + timeMs

	this.addListener()                          // 添加监听
	this.EDesk.SendGameState(this.Mark, timeMs) // 发送桌子状态

	this.EDesk.AddTimer(timeId, int(timeMs)/1000, this.TimerCall, nil)
}

func (this *FSMDownBets) Leave() {
	this.removeListener()
}

func (this *FSMDownBets) TimerCall(d interface{}) {
	// 跳转状态
	this.EDesk.RunFSM(GAME_STATUS_OPENCARD)
}

func (this *FSMDownBets) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()

	return remainTimeMS
}

// 添加网络监听
func (this *FSMDownBets) addListener() {
	// 玩家请求下注
	this.EDesk.Handle[MSG_GAME_QDOWNBET] = this.recvRSeatBet
}

// 删除网络监听
func (this *FSMDownBets) removeListener() {
	delete(this.EDesk.Handle, MSG_GAME_QDOWNBET)
}

func (this *FSMDownBets) onUserOnline(p *ExtPlayer) {

}

func (this *FSMDownBets) onUserOffline(p *ExtPlayer) {

}

// 接收到玩家下注
func (this *FSMDownBets) recvRSeatBet(p *ExtPlayer, d *DkInMsg) {
	req := GClientQDownBet{}
	json.Unmarshal([]byte(d.Data), &req)
	this.EDesk.UserDownBet(p, req.SeatIdx, req.CoinIdx, true)
}

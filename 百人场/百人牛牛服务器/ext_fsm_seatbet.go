/**
* 坐下玩家下注和抢座状态
**/
package main

import (
	"encoding/json"
	"logs"
)

type FSMSeatBet struct {
	UpMark int
	Mark   int
	EDesk  *ExtDesk

	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMSeatBet) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMSeatBet) GetMark() int {
	return this.Mark
}

func (this *FSMSeatBet) Run(upMark int) {
	DebugLog("游戏状态：座位玩家下注")
	logs.Debug("游戏状态：座位玩家下注")
	this.UpMark = upMark

	timeId := gameConfig.GameStatusTimer.RobSeatId
	timeMs := int64(gameConfig.GameStatusTimer.RobSeatMS)

	this.EndDateTime = GetTimeMS() + timeMs

	this.addListener()                          // 添加监听
	this.EDesk.SendGameState(this.Mark, timeMs) // 发送桌子状态

	// 发送强制下注
	for _, seat := range this.EDesk.Seats {
		if seat.UserId != 0 {
			this.forceDownBet(seat.Id, seat.UserId)
		}
	}

	this.EDesk.AddTimer(timeId, int(timeMs)/1000, this.TimerCall, nil)
}

func (this *FSMSeatBet) Leave() {
	this.removeListener()
}

func (this *FSMSeatBet) TimerCall(d interface{}) {
	// 跳转状态
	this.EDesk.RunFSM(GAME_STATUS_FACARD)
}

func (this *FSMSeatBet) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()

	return remainTimeMS
}

// 添加网络监听
func (this *FSMSeatBet) addListener() {
	// 玩家请求坐下
	this.EDesk.Handle[MSG_GAME_QSEATDOWN] = this.recvRSeatDown
	this.EDesk.Handle[MSG_GAME_QSEATUP] = this.recvRSeatUp

	// 玩家请求下注
	this.EDesk.Handle[MSG_GAME_QDOWNBET] = this.recvRSeatBet
}

// 删除网络监听
func (this *FSMSeatBet) removeListener() {
	delete(this.EDesk.Handle, MSG_GAME_QSEATDOWN)
	delete(this.EDesk.Handle, MSG_GAME_QSEATUP)
	delete(this.EDesk.Handle, MSG_GAME_QDOWNBET)
}

func (this *FSMSeatBet) onUserOnline(p *ExtPlayer) {

}

func (this *FSMSeatBet) onUserOffline(p *ExtPlayer) {

}

func (this *FSMSeatBet) forceDownBet(seatId int, userId int64) {
	// 强制坐下的玩家下注
	p := this.EDesk.GetPlayer(userId)
	if p != nil {
		this.EDesk.UserDownBet(p, seatId, -1, false)
	}
}

// 接收到玩家请求坐下
func (this *FSMSeatBet) recvRSeatDown(p *ExtPlayer, d *DkInMsg) {
	req := GClientQSeatDown{}
	json.Unmarshal([]byte(d.Data), &req)

	succ := this.EDesk.UserSeatDown(p, req.SeatIdx)

	if succ {
		this.forceDownBet(req.SeatIdx, p.Uid)
	}
}

// 接收到玩家请求站立， 当前状态不允许玩家站起
func (this *FSMSeatBet) recvRSeatUp(p *ExtPlayer, d *DkInMsg) {
	// this.EDesk.UserSeatUp(p, true)
}

// 接收到玩家下注
func (this *FSMSeatBet) recvRSeatBet(p *ExtPlayer, d *DkInMsg) {
	req := GClientQDownBet{}
	json.Unmarshal([]byte(d.Data), &req)

	this.EDesk.UserDownBet(p, req.SeatIdx, req.CoinIdx, true)
}

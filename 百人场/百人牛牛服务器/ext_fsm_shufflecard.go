/**
* 洗牌状态
**/
package main

import (
	"encoding/json"
	"logs"
)

type FSMShuffleCard struct {
	UpMark int
	Mark   int
	EDesk  *ExtDesk

	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMShuffleCard) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMShuffleCard) GetMark() int {
	return this.Mark
}

func (this *FSMShuffleCard) Run(upMark int) {
	DebugLog("游戏状态：洗牌")
	logs.Debug("游戏状态：洗牌")

	this.UpMark = upMark

	timeId := gameConfig.GameStatusTimer.ShufflecardId
	timeMs := int64(gameConfig.GameStatusTimer.ShufflecardMS)

	this.EndDateTime = GetTimeMS() + timeMs

	this.addListener()                          // 添加监听
	this.EDesk.SendGameState(this.Mark, timeMs) // 发送桌子状态

	this.EDesk.AddTimer(timeId, int(timeMs)/1000, this.TimerCall, nil)
}

func (this *FSMShuffleCard) Leave() {
	this.removeListener()
}

func (this *FSMShuffleCard) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_WAITSTART)
}

func (this *FSMShuffleCard) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()

	return remainTimeMS
}

// 添加网络监听
func (this *FSMShuffleCard) addListener() {
	this.EDesk.Handle[MSG_GAME_QDOWNBET] = this.recvRSeatBet
}

// 接收到玩家下注
func (this *FSMShuffleCard) recvRSeatBet(p *ExtPlayer, d *DkInMsg) {
	req := GClientQDownBet{}
	json.Unmarshal([]byte(d.Data), &req)
	this.EDesk.UserDownBet(p, req.SeatIdx, req.CoinIdx, true)
}

// 删除网络监听
func (this *FSMShuffleCard) removeListener() {
}

func (this *FSMShuffleCard) onUserOnline(p *ExtPlayer) {

}

func (this *FSMShuffleCard) onUserOffline(p *ExtPlayer) {

}

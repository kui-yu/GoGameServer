package main

import (
	"logs"

	"bl.com/util"
)

// 准备
func (this *ExtDesk) TimerReady(d interface{}) {
	this.Lock()
	defer this.Unlock()

	this.HandleExit()
	this.HandleUndo()

	// 刷新座位
	this.UpdatePlayer()

	this.GameId = GetJuHao() // util.BuildGameId(GCONFIG.GameType)

	players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range players {
		v.(*ExtPlayer).ResetAreaList()
		// 用户没有下注，增加未下注次数
		if v.(*ExtPlayer).GetMsgId() == 0 {
			v.(*ExtPlayer).AddUndoTimes()
		}
		if v.(*ExtPlayer).LiXian {
			logs.Debug("^^^^^^^^^^^^^准备阶段有用户离线：%v", v.(*ExtPlayer).Nick)
			this.SeatMgr.DelPlayer(v.(*ExtPlayer))
			this.LeaveByForce(v.(*ExtPlayer))
		}
	}

	dices1, _ := util.GetRandomNum(1, 6)
	dices2, _ := util.GetRandomNum(1, 6)
	this.dices1 = dices1
	this.dices2 = dices2
	sd := GGameReadyNotify{
		Id:     MSG_GAME_INFO_READY_NOTIFY,
		Timer:  int32(gameConfig.Timer.ReadyNum) * 1000,
		GameId: this.GameId,
		Dices:  []int{this.dices1, this.dices2},
	}

	allUser := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range allUser {
		sd.SeatList = this.GetSeatInfo(v.(*ExtPlayer))
		v.(*ExtPlayer).SendNativeMsg(MSG_GAME_INFO_READY_NOTIFY, sd)
	}

	this.GameState = MSG_GAME_INFO_READY_NOTIFY
	this.AddTimer(gameConfig.Timer.Ready, gameConfig.Timer.ReadyNum, this.TimerSendCard, nil)
}

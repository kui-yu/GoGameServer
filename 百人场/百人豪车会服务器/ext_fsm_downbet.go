package main

import (
	"encoding/json"
	//"fmt"
)

type FSMDownBet struct {
	Mark int

	EndDateTime int64 // 当前状态的结束时间

	EDesk *ExtDesk
}

func (this *FSMDownBet) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMDownBet) Run() {
	DebugLog("游戏状态-下注")
	this.EndDateTime = GetTimeMS() + int64(gameConfig.StateInfo.DownBetTime)

	this.EDesk.JuHao = GetJuHao()
	DebugLog("局号改变通知：", this.EDesk.JuHao)
	this.EDesk.SendNotice(MSG_GAME_INFO_NDESKCHANGE, &struct {
		Id    int
		JuHao string
	}{
		Id:    MSG_GAME_INFO_NDESKCHANGE,
		JuHao: this.EDesk.JuHao,
	}, true, nil)
	// this.EDesk.ChipList = G_DbGetGameServerData.GameConfig.TenChips
	// this.EDesk.MaxBet = G_DbGetGameServerData.GameConfig.LimitRedMax
	// fmt.Println("限红：", this.EDesk.MaxBet)
	// fmt.Println("筹码列表", this.EDesk.ChipList)
	this.addListen() // 添加监听
	this.EDesk.GameState = GAME_STATUS_DOWNBET
	this.EDesk.AddTimer(gameConfig.StateInfo.BroMsg, 0, this.EDesk.BroBetMsg, nil)
	this.EDesk.SendGameState(GAME_STATUS_DOWNBET, int64(gameConfig.StateInfo.DownBetTime))
	this.EDesk.AddTimer(GAME_STATUS_DOWNBET, gameConfig.StateInfo.DownBetTime/1000, this.TimerCall, nil)
}

func (this *FSMDownBet) GetMark() int {
	return this.Mark
}

func (this *FSMDownBet) Leave() {
	this.removeListen()
}

func (this *FSMDownBet) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FSMDownBet) addListen() {
	this.EDesk.Handle[MSG_GAME_INFO_QDOWNBET] = this.HandleBet
}

func (this *FSMDownBet) removeListen() {
	delete(this.EDesk.fsms, MSG_GAME_INFO_QDOWNBET)
}

func (this *FSMDownBet) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_LOTTERY)
}

func (this *FSMDownBet) HandleBet(p *ExtPlayer, d *DkInMsg) {
	data := GADownBet{}
	json.Unmarshal([]byte(d.Data), &data)
	//fmt.Println("data:", data)
	this.EDesk.UserDownBet(p, data.BetsIdx, data.CoinIdx)
}

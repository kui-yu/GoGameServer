/**
* 发给玩家的两张牌
**/

package main

type FsmHoleCards struct {
	Mark        int
	EDesk       *ExtDesk
	EndDateTime int64
}

func (this *FsmHoleCards) InitFsm(mark int, extDesk *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDesk
}

func (this *FsmHoleCards) GetMark() int {
	return this.Mark
}

func (this *FsmHoleCards) Run(upMark int, args ...interface{}) {
	DebugLog("状态 发给玩家的两张牌")
	timerMs := gameConfig.GameTimer.HoleCards
	this.EndDateTime = GetTimeMS() + int64(timerMs)

	this.EDesk.SendDeskStatus(this.Mark, timerMs)
	this.EDesk.AddUniueTimer(TimerId, int(timerMs/1000), this.TimerCall, nil)

	// 状态逻辑
	this.EDesk.IsOperateOpen = false
	this.EDesk.CurrStage = this.GetMark()
	cards := ShuffleCard(false)
	this.EDesk.DownCards = cards

	cardIdx := 0
	for _, p := range this.EDesk.Players {
		if p.State != UserStateGameIn {
			continue
		}

		holecardInfo := struct {
			Cards []int
		}{
			Cards: cards[cardIdx : cardIdx+2],
		}
		p.Cards = holecardInfo.Cards
		p.SendNetMessage(MSG_GAME_NGameHoleCards, holecardInfo)
		cardIdx += 2
	}
	this.EDesk.CurrDownCardIdx = cardIdx
}

func (this *FsmHoleCards) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GameStatusUserOperate)
}

func (this *FsmHoleCards) Leave() {

}

func (this *FsmHoleCards) Reset() {
}

func (this *FsmHoleCards) GetRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FsmHoleCards) OnUserOnline(p *ExtPlayer) {
}

func (this *FsmHoleCards) OnUserOffline(p *ExtPlayer) {
}

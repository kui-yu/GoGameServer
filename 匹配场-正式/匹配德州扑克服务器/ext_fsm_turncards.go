/**
* 第四张公共牌
**/

package main

type FsmTurnCards struct {
	Mark        int
	EDesk       *ExtDesk
	EndDateTime int64
}

func (this *FsmTurnCards) InitFsm(mark int, extDesk *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDesk
}

func (this *FsmTurnCards) GetMark() int {
	return this.Mark
}

func (this *FsmTurnCards) Run(upMark int, args ...interface{}) {
	DebugLog("状态 第四张公共牌")
	timerMs := gameConfig.GameTimer.PublicCards
	this.EndDateTime = GetTimeMS() + int64(timerMs)

	this.EDesk.SendDeskStatus(this.Mark, timerMs)
	this.EDesk.AddUniueTimer(TimerId, int(timerMs/1000), this.TimerCall, nil)

	this.EDesk.IsOperateOpen = true
	this.EDesk.CurrStage = this.GetMark()

	card := this.EDesk.DownCards[this.EDesk.CurrDownCardIdx]
	this.EDesk.CurrDownCardIdx++
	this.EDesk.PublicCards = append(this.EDesk.PublicCards, card)

	this.EDesk.SendNetMessage(MSG_GAME_NGamePublicCards, struct {
		Type  int
		Cards []int
	}{
		Type:  2,
		Cards: []int{card},
	})
}

func (this *FsmTurnCards) TimerCall(d interface{}) {
	this.EDesk.IsOperateOpen = true
	if this.EDesk.IsExistOperateFsm {
		this.EDesk.SetFirstOperateUser()
		this.EDesk.RunFSM(GameStatusUserOperate)
	} else {
		this.EDesk.RunFSM(this.EDesk.GetNextStageMark())
	}
}

func (this *FsmTurnCards) Leave() {

}

func (this *FsmTurnCards) Reset() {
}

func (this *FsmTurnCards) GetRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FsmTurnCards) OnUserOnline(p *ExtPlayer) {
}

func (this *FsmTurnCards) OnUserOffline(p *ExtPlayer) {
}

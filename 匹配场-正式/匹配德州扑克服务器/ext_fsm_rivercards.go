/**
* 第五张公共牌
**/

package main

type FsmRiverCards struct {
	Mark        int
	EDesk       *ExtDesk
	EndDateTime int64
}

func (this *FsmRiverCards) InitFsm(mark int, extDesk *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDesk
}

func (this *FsmRiverCards) GetMark() int {
	return this.Mark
}

func (this *FsmRiverCards) Run(upMark int, args ...interface{}) {
	DebugLog("状态 第五张公共牌")
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
		Type:  3,
		Cards: []int{card},
	})
}

func (this *FsmRiverCards) TimerCall(d interface{}) {
	this.EDesk.IsOperateOpen = true
	if this.EDesk.IsExistOperateFsm {
		this.EDesk.SetFirstOperateUser()
		this.EDesk.RunFSM(GameStatusUserOperate)
	} else {
		this.EDesk.RunFSM(this.EDesk.GetNextStageMark())
	}
}

func (this *FsmRiverCards) Leave() {

}

func (this *FsmRiverCards) Reset() {
}

func (this *FsmRiverCards) GetRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FsmRiverCards) OnUserOnline(p *ExtPlayer) {
}

func (this *FsmRiverCards) OnUserOffline(p *ExtPlayer) {
}

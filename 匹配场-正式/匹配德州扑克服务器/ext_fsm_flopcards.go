/**
* 头三张公共牌
**/

package main

type FsmFlopCards struct {
	Mark        int
	EDesk       *ExtDesk
	EndDateTime int64
}

func (this *FsmFlopCards) InitFsm(mark int, extDesk *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDesk
}

func (this *FsmFlopCards) GetMark() int {
	return this.Mark
}

func (this *FsmFlopCards) Run(upMark int, args ...interface{}) {
	DebugLog("状态 头三张公共牌")
	timerMs := gameConfig.GameTimer.PublicCards
	this.EndDateTime = GetTimeMS() + int64(timerMs)

	this.EDesk.SendDeskStatus(this.Mark, timerMs)
	this.EDesk.AddUniueTimer(TimerId, int(timerMs/1000), this.TimerCall, nil)

	this.EDesk.IsOperateOpen = true
	this.EDesk.CurrStage = this.GetMark()

	cardIdx := this.EDesk.CurrDownCardIdx
	this.EDesk.CurrDownCardIdx += 3

	cards := this.EDesk.DownCards[cardIdx : cardIdx+3]
	this.EDesk.PublicCards = append(this.EDesk.PublicCards, cards...)
	// 发牌
	this.EDesk.SendNetMessage(MSG_GAME_NGamePublicCards, struct {
		Type  int
		Cards []int
	}{
		Type:  1,
		Cards: cards,
	})

}

func (this *FsmFlopCards) TimerCall(d interface{}) {
	this.EDesk.IsOperateOpen = true
	if this.EDesk.IsExistOperateFsm {
		DebugLog(">>>111")
		this.EDesk.SetFirstOperateUser()
		DebugLog(">>>333")
		this.EDesk.RunFSM(GameStatusUserOperate)
		DebugLog(">>>444")
	} else {
		DebugLog(">>>222")
		this.EDesk.RunFSM(this.EDesk.GetNextStageMark())
	}
}

func (this *FsmFlopCards) Leave() {

}

func (this *FsmFlopCards) Reset() {
}

func (this *FsmFlopCards) GetRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FsmFlopCards) OnUserOnline(p *ExtPlayer) {
}

func (this *FsmFlopCards) OnUserOffline(p *ExtPlayer) {
}

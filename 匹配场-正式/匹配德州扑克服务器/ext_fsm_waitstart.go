/**
* 等待开始游戏
**/

package main

type FsmWaitStart struct {
	mark    int
	extDesk *ExtDesk
}

func (this *FsmWaitStart) InitFsm(mark int, extDesk *ExtDesk) {
	this.mark = mark
	this.extDesk = extDesk
}

func (this *FsmWaitStart) GetMark() int {
	return this.mark
}

func (this *FsmWaitStart) Run(upMark int, args ...interface{}) {
}

func (this *FsmWaitStart) Leave() {

}

func (this *FsmWaitStart) Reset() {
}

func (this *FsmWaitStart) GetRestTime() int64 {
	return 0
}

func (this *FsmWaitStart) OnUserOnline(p *ExtPlayer) {
}

func (this *FsmWaitStart) OnUserOffline(p *ExtPlayer) {
}

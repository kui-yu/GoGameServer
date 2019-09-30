/**
* 等待开始状态
**/
package main

type FsmWaitStart struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FsmWaitStart) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FsmWaitStart) GetMark() int {
	return this.Mark
}

func (this *FsmWaitStart) Run(upMark int, overtime int) {
	DebugLog("进入游戏状态：等待开始")
	this.UpMark = upMark

	this.addListener() // 添加监听
}

func (this *FsmWaitStart) Leave() {
	this.removeListener()
}

func (this *FsmWaitStart) onEvent(interface{}) {

}

// 添加网络监听
func (this *FsmWaitStart) addListener() {

}

// 删除网络监听
func (this *FsmWaitStart) removeListener() {

}

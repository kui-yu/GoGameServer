/**
* 随机庄家，下盲注
**/
package main

type FsmRandBank struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FsmRandBank) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FsmRandBank) GetMark() int {
	return this.Mark
}

func (this *FsmRandBank) Run(upMark int, overtime int) {
	DebugLog("进入游戏状态：随机庄家，下盲注")
	this.UpMark = upMark

	this.RC.GameIn = true
	this.addListener() // 添加监听
}

func (this *FsmRandBank) Leave() {
	this.removeListener()
}

func (this *FsmRandBank) onEvent(interface{}) {

}

// 添加网络监听
func (this *FsmRandBank) addListener() {

}

// 删除网络监听
func (this *FsmRandBank) removeListener() {

}

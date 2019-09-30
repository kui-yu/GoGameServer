/**
* 开牌状态
**/
package main

type FSMOpenCard struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FSMOpenCard) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FSMOpenCard) GetMark() int {
	return this.Mark
}

func (this *FSMOpenCard) Run(upMark int) {
	DebugLog("进入游戏状态：开牌")
	this.UpMark = upMark

	this.addListener() // 添加监听

}

func (this *FSMOpenCard) Leave() {
	this.removeListener()
}

func (this *FSMOpenCard) onEvent(interface{}) {

}

// 添加网络监听
func (this *FSMOpenCard) addListener() {
}

// 删除网络监听
func (this *FSMOpenCard) removeListener() {
}

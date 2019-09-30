/**
* 发牌状态
**/
package main

type FSMFaCard struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FSMFaCard) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FSMFaCard) GetMark() int {
	return this.Mark
}

func (this *FSMFaCard) Run(upMark int) {
	DebugLog("进入游戏状态：发牌")

	this.UpMark = upMark

	this.addListener() // 添加监听
}

func (this *FSMFaCard) Leave() {
	this.removeListener()
}

func (this *FSMFaCard) onEvent(interface{}) {

}

// 添加网络监听
func (this *FSMFaCard) addListener() {
}

// 删除网络监听
func (this *FSMFaCard) removeListener() {
}

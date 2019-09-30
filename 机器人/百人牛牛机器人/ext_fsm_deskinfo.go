/**
* 请求桌子信息
**/
package main

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

type FSMGameDesk struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FSMGameDesk) InitFSM(mark int, client *ExtRobotClient) {
	this.Mark = mark
	this.RC = client
}

func (this *FSMGameDesk) GetMark() int {
	return this.Mark
}

func (this *FSMGameDesk) Run(upMark int) {
	DebugLog("状态：请求桌子")

	this.UpMark = upMark

	this.addListener()

	this.sendGameDesk()
}

func (this *FSMGameDesk) Leave() {
	this.removeListener()
}

func (this *FSMGameDesk) onEvent(interface{}) {

}

// 添加网络监听
func (this *FSMGameDesk) addListener() {
	this.RC.Handle[MSG_GAME_RDESKINFO] = this.onGameDesk
}

// 删除网络监听
func (this *FSMGameDesk) removeListener() {
	delete(this.RC.Handle, MSG_GAME_RDESKINFO)
}

// 请求桌子消息
func (this *FSMGameDesk) sendGameDesk() {
	this.RC.AddMsgNative(MSG_GAME_QDESKINFO, struct {
		Id int32 //协议号
	}{
		Id: MSG_GAME_QDESKINFO,
	})
}

// 接收到桌子消息
func (this *FSMGameDesk) onGameDesk(d string) {
	DebugLog("桌子消息", d)
	if gjson.Get(d, "Result").Int() != 0 {
		ErrorLog("请求桌子信息失败 [%s]", gjson.Get(d, "Err").String())
		controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this.RC)
		return
	}

	json.Unmarshal([]byte(d), &this.RC.DeskInfo)

	this.RC.Coin = this.RC.DeskInfo.MyUserCoin

	// 进入网络驱动状态机模式
	this.RC.GameIn = true
	this.RC.AddGameHandler()

	DebugLog("桌子状态：%d 持续时间：%d", this.RC.DeskInfo.GameStatus, this.RC.DeskInfo.GameStatusDuration)
}

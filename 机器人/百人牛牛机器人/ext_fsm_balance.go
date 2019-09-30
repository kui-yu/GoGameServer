/**
* 结算状态
**/
package main

import (
	"github.com/tidwall/gjson"
)

type FSMBalance struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FSMBalance) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FSMBalance) GetMark() int {
	return this.Mark
}

func (this *FSMBalance) Run(upMark int) {
	DebugLog("进入游戏状态：结算")

	this.UpMark = upMark

	this.addListener() // 添加监听
}

func (this *FSMBalance) Leave() {
	this.removeListener()
}

func (this *FSMBalance) onEvent(interface{}) {

}

// 添加网络监听
func (this *FSMBalance) addListener() {
	this.RC.Handle[MSG_GAME_BALANCE] = this.onGameDesk
}

// 删除网络监听
func (this *FSMBalance) removeListener() {
	delete(this.RC.Handle, MSG_GAME_BALANCE)
}

// 结算
func (this *FSMBalance) onGameDesk(d string) {
	DebugLog("结算数据", d)

	this.RC.DeskInfo.MyUserCoin = gjson.Get(d, "MyCoin").Int()
	this.RC.Coin = this.RC.DeskInfo.MyUserCoin

	seatIdx := -1
	seats := this.RC.DeskInfo.Seats
	for k, v := range seats {
		if v.UserId == this.RC.UserInfo.Uid {
			seatIdx = k
			break
		}
	}

	// 发送当前局结速
	if seatIdx == -1 || ((seats[seatIdx].SeatDownCount + 1) >= this.RC.DeskInfo.SeatUpTotalCount) {
		controller.sendEvent(EVENT_CONT_ROBOTIDLE, this.RC)
	}
}

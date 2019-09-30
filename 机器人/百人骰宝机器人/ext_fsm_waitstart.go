/**
* 等待开始状态
**/
package main

type FSMWaitStart struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FSMWaitStart) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FSMWaitStart) GetMark() int {
	return this.Mark
}

func (this *FSMWaitStart) Run(upMark int) {
	DebugLog("进入游戏状态：等待开始")
	this.UpMark = upMark

	this.addListener() // 添加监听

	this.checkAndSendSeatUp() // 检查桌位信息
}

func (this *FSMWaitStart) Leave() {
	this.removeListener()
}

func (this *FSMWaitStart) onEvent(interface{}) {

}

// 添加网络监听
func (this *FSMWaitStart) addListener() {

}

// 删除网络监听
func (this *FSMWaitStart) removeListener() {

}

// 检查并发送机器人站起
func (this *FSMWaitStart) checkAndSendSeatUp() {
	DebugLog("========= 检查并发送机器人站起 ", this.RC.UserInfo.Uid)

	seatIdx := -1
	seats := this.RC.DeskInfo.Seats
	for k, v := range seats {
		if v.UserId == this.RC.UserInfo.Uid {
			seatIdx = k
			break
		}
	}

	if seatIdx != -1 {
		seatInfo := seats[seatIdx]
		seatInfo.SeatDownCount++

		this.RC.DeskInfo.Seats[seatIdx] = seatInfo

		if seatInfo.SeatDownCount >= this.RC.DeskInfo.SeatUpTotalCount {
			DebugLog("发送机器人站起")
			this.RC.AddMsgNative(MSG_GAME_QSEATUP, struct {
				Id int
			}{
				Id: MSG_GAME_QSEATUP,
			})
		}
	}
}

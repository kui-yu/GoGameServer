/**
* 坐下玩家下注和抢座状态
**/
package main

type FSMSeatBet struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient

	TimerIds []int
}

func (this *FSMSeatBet) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc

	this.RC.EventHandle[EVENT_ROBOT_SEATDOWN] = this.onSetSeatDown
}

func (this *FSMSeatBet) GetMark() int {
	return this.Mark
}

func (this *FSMSeatBet) Run(upMark int) {
	DebugLog("进入游戏状态：座位玩家下注")
	this.UpMark = upMark

	this.addListener() // 添加监听

	if this.RC.IsSeatDown() {
		var time = this.RC.DeskInfo.GameStatusDuration / 1000
		if time < 1 {
			return
		}

		downbetnum := gameConfig.getGameConfigInt("downbetnum")

		for i := 0; i < downbetnum; i++ {
			// 开始下注
			t, _ := GetRandomNum(1, int(time))
			timeId := this.RC.TimeTicker.AddTimer(t, this.downBet, nil)
			this.TimerIds = append(this.TimerIds, timeId)
		}
	}

}

func (this *FSMSeatBet) Leave() {
	this.removeListener()
}

func (this *FSMSeatBet) onEvent(interface{}) {

}

// 添加网络监听
func (this *FSMSeatBet) addListener() {
	this.RC.Handle[MSG_GAME_RSEATDOWN] = this.onRSeatDown // 玩家请求坐下的回复
	this.RC.Handle[MSG_GAME_RDOWNBET] = this.onDownBetR
}

// 删除网络监听
func (this *FSMSeatBet) removeListener() {
	delete(this.RC.Handle, MSG_GAME_RSEATDOWN)
	delete(this.RC.Handle, MSG_GAME_RDOWNBET)
}

// 坐下回复
func (this *FSMSeatBet) onRSeatDown(str string) {
	DebugLog("玩家请求坐下的回复", str)
}

// 机器人坐下
func (this *FSMSeatBet) onSetSeatDown(d interface{}) {
	ms := d.(int64)
	DebugLog("机器人坐下 time %d", ms)

	minsecondseat := gameConfig.getGameConfigInt("minsecondseat")
	t, _ := GetRandomNum(int(minsecondseat), int(ms)/1000)
	this.RC.TimeTicker.AddTimer(t, func(timerId int, d interface{}) {
		seatLen := 0
		for _, seat := range this.RC.DeskInfo.Seats {
			if seat.UserId != 0 {
				seatLen++
			}
		}

		if (4 - seatLen) > gameConfig.getGameConfigInt("minemptyseat") {
			seatIdx := 0
			isExists := false
			for i := 0; i < 4; i++ {
				isExists = false
				for _, v := range this.RC.DeskInfo.Seats {
					if v.Id == i && v.UserId != 0 {
						isExists = true
						break
					}
				}

				if isExists == false {
					seatIdx = i
					break
				}
			}

			DebugLog("发送机器人坐下消息 %d，%d", MSG_GAME_QSEATDOWN, seatIdx)
			this.RC.AddMsgNative(MSG_GAME_QSEATDOWN, struct {
				Id      int32 //协议号
				SeatIdx int   // 座位索引
			}{
				Id:      MSG_GAME_QSEATDOWN,
				SeatIdx: seatIdx,
			})
		}
	}, nil)
}

// 下注
func (this *FSMSeatBet) downBet(id int, d interface{}) {
	if this.RC.DeskInfo.GameStatus != this.Mark {
		return
	}

	this.TimerIds = DeleteIntArrayFromValue(this.TimerIds, id)

	levelLen := len(this.RC.DeskInfo.BetLevels)
	DebugLog("发送下注levelLen:%d", levelLen)

	coinIdx, _ := GetRandomNum(0, levelLen)
	seatIdx, _ := GetRandomNum(0, 4)

	this.RC.AddMsgNative(MSG_GAME_QDOWNBET, struct {
		Id      int32
		SeatIdx int
		CoinIdx int
	}{
		Id:      MSG_GAME_QDOWNBET,
		SeatIdx: seatIdx,
		CoinIdx: coinIdx,
	})
}

func (this *FSMSeatBet) onDownBetR(str string) {
	DebugLog("接收到下注回复", str)
}

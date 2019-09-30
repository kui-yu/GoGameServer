/**
* 等待玩家下注
**/
package main

type FSMDownBets struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient

	TimerIds []int
}

func (this *FSMDownBets) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FSMDownBets) GetMark() int {
	return this.Mark
}

func (this *FSMDownBets) Run(upMark int) {
	DebugLog("进入游戏状态：所有玩家下注")
	this.UpMark = upMark

	this.addListener()

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

func (this *FSMDownBets) Leave() {
	this.removeListener()
	for _, timeId := range this.TimerIds {
		this.RC.TimeTicker.DelTimer(timeId)
	}
	this.TimerIds = []int{}
}

func (this *FSMDownBets) onEvent(interface{}) {

}

// 添加网络监听
func (this *FSMDownBets) addListener() {
	this.RC.Handle[MSG_GAME_RDOWNBET] = this.onDownBetR
}

// 删除网络监听
func (this *FSMDownBets) removeListener() {
	delete(this.RC.Handle, MSG_GAME_RDOWNBET)
}

// 下注
func (this *FSMDownBets) downBet(id int, d interface{}) {
	if this.RC.DeskInfo.GameStatus != this.Mark {
		return
	}

	this.TimerIds = DeleteIntArrayFromValue(this.TimerIds, id)
	levelLen := len(this.RC.DeskInfo.BetLevels)
	DebugLog("发送下注 levelLen:%d", levelLen)

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

func (this *FSMDownBets) onDownBetR(str string) {
	DebugLog("接收到下注回复", str)
}

package main

import (
	"time"

	"github.com/tidwall/gjson"
)

type ExtRobotClient struct {
	RobotClient

	DeskInfo   GClientDeskInfo
	fsms       map[int]FSMBase // fsm状态机集合
	currentFSM FSMBase         // 当前状态机
	upFSM      FSMBase         // 上一个状态机
}

//=======================对外方法=====================
// 机器人开始
func (this *ExtRobotClient) Start() {
	this.ChildPtr = this
	this.IsStop = false
	this.GameIn = false
	this.IsConnect = false

	this.Coin = int64(gameConfig.getGameConfigCoin("initcoin"))

	this.fsms = make(map[int]FSMBase)
	this.SendMsg = make(chan *SendMsg, 1000)
	this.RecvMsg = make(chan *RecvMsg, 1000)
	this.Handle = make(map[uint32]func(string))
	this.EventMsg = make(chan *EventMsg, 1000)
	this.EventHandle = make(map[int32]func(interface{}))
	this.TimeTicker = new(TimeTicker)
	this.SysTicker = nil

	DebugLog("机器人开始了")
	this.addEventHandler()
	this.addFSM(GAME_STATUS_GAME_DESK, new(FSMGameDesk))
	this.addFSM(GAME_STATUS_WAITSTART, new(FSMWaitStart))
	this.addFSM(GAME_STATUS_SEATBET, new(FSMSeatBet))
	this.addFSM(GAME_STATUS_FACARD, new(FSMFaCard))
	this.addFSM(GAME_STATUS_DOWNBTES, new(FSMDownBets))
	this.addFSM(GAME_STATUS_OPENCARD, new(FSMOpenCard))
	this.addFSM(GAME_STATUS_BALANCE, new(FSMBalance))
	this.CheckOnLineId = this.TimeTicker.AddTimer(30, func(int, interface{}) {
		this.RestChekIsNo()
	}, nil)
	go func() {
		this.AddHandler()
		this.requstAccountInfo(this.Coin)
		this.OnMessageHandle()
	}()
}
func (this *ExtRobotClient) RestChekIsNo() {
	time.Sleep(time.Second * 30)
	controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this)
}
func (this *ExtRobotClient) GetFSM(mark int) FSMBase {
	if mark != 0 {
		return this.fsms[mark]
	}
	return this.currentFSM
}

// 是否有这个状态
func (this *ExtRobotClient) exitstFSM(mark int) bool {
	_, ok := this.fsms[mark]
	return ok
}

func (this *ExtRobotClient) RunFSM(mark int) {
	if this.IsRobotStop() {
		return
	}

	if _, ok := this.fsms[mark]; ok == false {
		return
	}

	var upMark int = 0
	if this.currentFSM != nil {
		this.upFSM = this.currentFSM
		this.upFSM.Leave()
		upMark = this.upFSM.GetMark()
	}

	this.currentFSM = this.GetFSM(mark)
	this.currentFSM.Run(upMark)
}

func (this *ExtRobotClient) addFSM(mark int, fsm FSMBase) {
	fsm.InitFSM(mark, this)
	this.fsms[mark] = fsm
}

func (this *ExtRobotClient) addEventHandler() {
	this.EventHandle[EVENT_CONNECT_SUCCESS] = this.onConnected
}

func (this *ExtRobotClient) AddHandler() {
	this.Handle[MSG_HALL_ROBOT_LOGIN_REPLY] = this.onLogin
	this.Handle[MSG_HALL_JOIN_GAME_REPLY] = this.onJoinGame
	this.Handle[MSG_GAME_AUTO_REPLY] = this.onGameAuto
}

// 添加网络消息监听
func (this *ExtRobotClient) AddGameHandler() {

	this.Handle[MSG_GAME_NSTATUS_CHANGE] = this.onGameState
	this.Handle[MSG_GAME_NDOWNBET] = this.onDownBetN
	this.Handle[MSG_GAME_NSEATDOWN] = this.onNSeatDown // 座位信息改变通知
	this.Handle[MSG_GAME_RMANYUSER] = this.onManyUser
}

func (this *ExtRobotClient) RemoveAllHandler() {
	delete(this.Handle, MSG_HALL_ROBOT_LOGIN_REPLY)
	delete(this.Handle, MSG_HALL_JOIN_GAME_REPLY)
	delete(this.Handle, MSG_GAME_AUTO_REPLY)

	delete(this.Handle, MSG_GAME_NSTATUS_CHANGE)
	delete(this.Handle, MSG_GAME_NDOWNBET)
	delete(this.Handle, MSG_GAME_NSEATDOWN)
	delete(this.Handle, MSG_GAME_RMANYUSER)
}

// 游戏状态改变
func (this *ExtRobotClient) onGameState(str string) {
	if this.IsRobotStop() {
		return
	}

	DebugLog("游戏状态改变", str)
	status := int(gjson.Get(str, "GameStatus").Int())

	// 当前状态机是否存在
	if this.exitstFSM(status) == false {
		return
	}

	this.DeskInfo.GameStatus = status
	this.DeskInfo.GameStatusDuration = gjson.Get(str, "GameStatusDuration").Int()

	this.RunFSM(this.DeskInfo.GameStatus)

	// 如果当前状态是抢座状态
	if controller.getRobotClients()[0] == this && this.DeskInfo.GameStatus == GAME_STATUS_SEATBET {
		this.onAllocSeat()
	}
}

// 下注通知
func (this *ExtRobotClient) onDownBetN(str string) {
	DebugLog("接收到下注", str)

	uid := gjson.Get(str, "Uid").Int()
	if this.UserInfo.Uid == uid {
		coin := gjson.Get(str, "Coin").Int()
		DebugLog("用户自己下注扣除金币", coin)
	}
}

func (this *ExtRobotClient) onManyUser(json string) {
	TestLog("请求更多玩家返回", json)
}

func (this *ExtRobotClient) onAllocSeat() {
	// 判断是否需要机器人开始抢座
	seatLen := 0
	for _, seat := range this.DeskInfo.Seats {
		if seat.UserId != 0 {
			seatLen++
		}
	}

	emptySeatLen := 4 - seatLen
	minemptyseat := gameConfig.getGameConfigInt("minemptyseat")

	if emptySeatLen <= minemptyseat {
		return
	}

	num := emptySeatLen - minemptyseat

	clients := controller.getRobotClients()
	// 删除坐下的玩家
	var useclients []*ExtRobotClient
	for _, v := range clients {
		var exists = false
		for _, s := range v.DeskInfo.Seats {
			if v.UserInfo.Uid == s.UserId {
				exists = true
				break
			}
		}

		if exists == false {
			useclients = append(useclients, v)
		}
	}

	clen := len(useclients)

	if clen <= num {
		for i := 0; i < clen; i++ {
			useclients[i].AddEventNative(EVENT_ROBOT_SEATDOWN, this.DeskInfo.GameStatusDuration)
		}
	} else {
		var oldIdxs map[int]bool = make(map[int]bool)

		for {

			idx, _ := GetRandomNum(0, clen)
			if _, ok := oldIdxs[idx]; ok == false {
				oldIdxs[idx] = true
				useclients[idx].AddEventNative(EVENT_ROBOT_SEATDOWN, this.DeskInfo.GameStatusDuration)
				num--
				if num == 0 {
					break
				}
			}
		}
	}
}

// 坐下通知
func (this *ExtRobotClient) onNSeatDown(str string) {
	DebugLog("座位信息改变通知", str)

	idx := gjson.Get(str, "SeatId").Int()
	stype := gjson.Get(str, "Type").Int() //0添加 1修改 2删除

	if stype == 0 {
		this.DeskInfo.Seats = append(this.DeskInfo.Seats, GClientSeatInfo{
			Id:            int(idx),
			UserId:        gjson.Get(str, "NewUserId").Int(),
			Name:          gjson.Get(str, "NewUserName").String(),
			Avatar:        gjson.Get(str, "NewUserAvatar").String(),
			SeatDownCount: 1,
			DownBetTotal:  0,
		})
	} else if stype == 1 {
		for i, seat := range this.DeskInfo.Seats {
			if seat.Id == int(idx) {
				seat.UserId = gjson.Get(str, "NewUserId").Int()
				seat.Name = gjson.Get(str, "NewUserName").String()
				seat.Avatar = gjson.Get(str, "NewUserAvatar").String()
				this.DeskInfo.Seats[i] = seat
				break
			}
		}
	} else if stype == 2 {
		for i, seat := range this.DeskInfo.Seats {
			if seat.Id == int(idx) {
				this.DeskInfo.Seats = append(this.DeskInfo.Seats[:i], this.DeskInfo.Seats[i+1:]...)
				break
			}
		}
	}
}

// 判断是否在座位上
func (this *ExtRobotClient) IsSeatDown() bool {
	seats := this.DeskInfo.Seats

	isSeatDown := false
	for _, v := range seats {
		if v.UserId == this.UserInfo.Uid {
			isSeatDown = true
			break
		}
	}

	return isSeatDown
}

// 判断是否在座位上
func (this *ExtRobotClient) GetMySeatDown() *GClientSeatInfo {
	seats := this.DeskInfo.Seats

	for _, v := range seats {
		if v.UserId == this.UserInfo.Uid {
			return &v
		}
	}

	return nil
}

func (this *ExtRobotClient) Stop() bool {
	if this.IsRobotStop() {
		return false
	}

	this.RWMutex.Lock()
	this.IsStop = true
	defer this.RWMutex.Unlock()

	this.BaseStop()

	fsm := this.GetFSM(0)
	if fsm != nil {
		fsm.Leave()
	}
	this.RemoveAllHandler()

	return true
}

//匹配结果
func (this *ExtRobotClient) onGameAuto(d string) {
	result := gjson.Get(d, "Result").Int()
	if result != 0 {
		ErrorLog("请求匹配房间失败归还机器人重新请求账号 [%s]", gjson.Get(d, "Err").String())
		controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this)

		return
	}

	this.RunFSM(GAME_STATUS_GAME_DESK)
}

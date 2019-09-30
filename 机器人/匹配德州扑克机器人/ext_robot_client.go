package main

import (
	"encoding/json"
	"time"

	"github.com/tidwall/gjson"
)

type ExtRobotClient struct {
	RobotClient

	Sid        int   //座位号
	CarryCoin  int64 //携带金币
	SmallBlind int64 //小盲
	Cards      []int //牌集合

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
	coin, _ := GetRandomNum64(gameConfig.getGameConfigCoin("mincoin"), gameConfig.getGameConfigCoin("maxcoin"))

	this.Coin = coin

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

	this.addFSM(GameStatusWaitStart, new(FsmWaitStart))
	this.addFSM(GameStatusRandBank, new(FsmRandBank))
	this.addFSM(GameStatusHoleCards, new(FsmHoleCards))
	this.addFSM(GameStatusFlopCards, new(FsmFlopCards))
	this.addFSM(GameStatusTurnCards, new(FsmTurnCards))
	this.addFSM(GameStatusRiverCards, new(FsmRiverCards))
	this.addFSM(GameStatusUserOperate, new(FsmUserOperate))
	this.addFSM(GameStatusResults, new(FsmResults))

	go func() {
		this.AddHandler()
		this.requstAccountInfo(this.Coin)
		this.OnMessageHandle()
	}()
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

func (this *ExtRobotClient) RunFSM(mark int, overtime int) {
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
	this.currentFSM.Run(upMark, overtime)
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
	this.Handle[MSG_GAME_NGameStatus] = this.onGameState
	this.Handle[MSG_GAME_NGameInfo] = this.onGameDeskInfo
	this.Handle[MSG_GAME_NGameReconnectInfo] = this.onGameReconnectInfo
	this.Handle[MSG_GAME_NGameJackpotChange] = this.onGameJackpotChange
	this.Handle[MSG_GAME_NGameSetCoin] = this.onGameSetCoin
	this.Handle[MSG_GAME_RGameSetCoin] = this.onRepGameSetCoin
}

func (this *ExtRobotClient) RemoveAllHandler() {
	delete(this.Handle, MSG_HALL_ROBOT_LOGIN_REPLY)
	delete(this.Handle, MSG_HALL_JOIN_GAME_REPLY)
	delete(this.Handle, MSG_GAME_AUTO_REPLY)

	delete(this.Handle, MSG_GAME_NGameStatus)
	delete(this.Handle, MSG_GAME_NGameInfo)
	delete(this.Handle, MSG_GAME_NGameReconnectInfo)
	delete(this.Handle, MSG_GAME_NGameJackpotChange)
	delete(this.Handle, MSG_GAME_NGameSetCoin)
	delete(this.Handle, MSG_GAME_RGameSetCoin)
}

// 游戏状态改变
func (this *ExtRobotClient) onGameState(str string) {
	if this.IsRobotStop() {
		return
	}

	DebugLog("游戏状态改变", str)
	status := int(gjson.Get(str, "GameStatus").Int())
	overtime := int(gjson.Get(str, "OverTime").Int())

	// 当前状态机是否存在
	if this.exitstFSM(status) == false {
		return
	}

	this.RunFSM(status, overtime)
}

func (this *ExtRobotClient) Reset() {
	this.GameIn = false
	this.Cards = []int{}

	//通知可以离开
	controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
}

// 桌子消息
func (this *ExtRobotClient) onGameDeskInfo(str string) {
	if this.IsRobotStop() {
		return
	}

	this.SmallBlind = gjson.Get(str, "DownBlind").Int()

	for _, itemInfo := range gjson.Get(str, "Users").Array() {
		if itemInfo.Get("Uid").Int() == this.UserInfo.Uid {
			this.Sid = int(itemInfo.Get("Sid").Int())
			this.CarryCoin = itemInfo.Get("Coin").Int()
			break
		}
	}

	DebugLog("接收到桌子消息", str)
	DebugLog("用户金币", this.CarryCoin)
}

func (this *ExtRobotClient) onGameReconnectInfo(str string) {
	if this.IsRobotStop() {
		return
	}

	for _, itemInfo := range gjson.Get(str, "Users").Array() {
		if itemInfo.Get("Uid").Int() == this.UserInfo.Uid {
			this.Sid = int(itemInfo.Get("Sid").Int())
			this.CarryCoin = itemInfo.Get("Coin").Int()
			break
		}
	}

	this.SmallBlind = gjson.Get(str, "DownBlind").Int()
	status := int(gjson.Get(str, "GameStatus").Int())
	overtime := int(gjson.Get(str, "OverTime").Int())

	DebugLog("接收到断线重连信息", str)
	DebugLog("用户金币", this.CarryCoin)

	// 当前状态机是否存在
	if this.exitstFSM(status) == false {
		return
	}

	this.RunFSM(status, overtime)
}

func (this *ExtRobotClient) onGameJackpotChange(str string) {
	data := struct {
		Sid       int
		ChangeVal int64
		Value     int64
	}{}

	json.Unmarshal([]byte(str), &data)

	if data.Sid == this.Sid {
		this.Coin -= data.ChangeVal
		this.CarryCoin -= data.ChangeVal
	}
}

func (this *ExtRobotClient) onGameSetCoin(str string) {
	controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
}

func (this *ExtRobotClient) onRepGameSetCoin(str string) {
	data := struct {
		Coin int64
		Err  string
	}{}

	json.Unmarshal([]byte(str), &data)
	if len(data.Err) != 0 {
		ErrorLog("设置金币错误 %s", data.Err)
		return
	}

	this.CarryCoin = data.Coin
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

	if result == 11 {
		DebugLog("玩家已在游戏中，换机器人")
		controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
		return
	} else if result == 9 {
		DebugLog("匹配失败金币不足，换机器人")
		controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
		return
	}

	if result != 0 {
		time.Sleep(time.Second * time.Duration(5))
		this.sendGameAuto()
	} else {
		this.AddGameHandler()
	}
}

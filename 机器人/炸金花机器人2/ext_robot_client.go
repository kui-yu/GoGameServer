package main

import (
	"math/rand"
	"time"

	"github.com/tidwall/gjson"
)

type ExtRobotClient struct {
	RobotClient

	HandCard   []byte
	SeatId     int32
	CurId      int32
	MTimerNum  int
	MGameState int32
	MJiaoFen   int
	MMaxOut    *GGameOutCardReply
	CallTimes  int
	CallEd     bool
}

//=======================对外方法=====================
// 机器人开始
func (this *ExtRobotClient) Start() {
	this.ChildPtr = this
	this.IsStop = false
	this.GameIn = false
	this.IsConnect = false

	//初始化金币
	time.Sleep(time.Second * time.Duration(1))
	minCoin := gameConfig.getGameConfigCoin("initCoinMin")
	maxCoin := gameConfig.getGameConfigCoin("initCoinMax")

	rand.Seed(time.Now().UnixNano())
	this.Coin = (int64(rand.Int63n(maxCoin-minCoin)+minCoin) / 1000) * 1000

	this.SendMsg = make(chan *SendMsg, 1000)
	this.RecvMsg = make(chan *RecvMsg, 1000)
	this.Handle = make(map[uint32]func(string))
	this.EventMsg = make(chan *EventMsg, 1000)
	this.EventHandle = make(map[int32]func(interface{}))
	this.TimeTicker = new(TimeTicker)
	this.SysTicker = nil

	DebugLog("机器人开始了")
	this.addEventHandler()
	this.AddHandler()
	go func() {
		//todo:0005此处添加获取机器人账号
		this.requstAccountInfo(this.Coin)
		this.OnMessageHandle()
	}()
}

func (this *ExtRobotClient) addEventHandler() {
	this.EventHandle[EVENT_CONNECT_SUCCESS] = this.onConnected
}

// 添加网络消息监听
func (this *ExtRobotClient) AddHandler() {
	//todo:0003此处添加游戏中的消息处理函数
	this.Handle[MSG_HALL_ROBOT_LOGIN_REPLY] = this.onLogin
	this.Handle[MSG_HALL_JOIN_GAME_REPLY] = this.onJoinGame
	this.Handle[MSG_GAME_AUTO_REPLY] = this.onGameAuto
	//机器人消息处理
	this.Handle[MSG_GAME_INFO_AUTO_REPLY] = this.HandleGameAutoReply
	this.Handle[MSG_GAME_INFO_END_NOTIFY] = this.GameEndNotify
}

func (this *ExtRobotClient) Stop() bool {
	if this.IsRobotStop() {
		return false
	}

	this.RWMutex.Lock()
	this.IsStop = true
	defer this.RWMutex.Unlock()

	this.BaseStop()
	//todo:0002此处添加离开处理

	return true
}

//匹配结果
func (this *ExtRobotClient) onGameAuto(d string) {
	// DebugLog("匹配房间消息结果", d)
	result := gjson.Get(d, "Result").Int()
	if result == 1 {
		// 服务器繁忙，请稍后重试
		time.Sleep(time.Second * time.Duration(10))
		this.sendGameAuto()
		return
	} else if result == 11 {
		// 玩家已在游戏中
		controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
		return
	} else if result == 9 {
		// 匹配失败，金币不足
		controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
		return
	}

	if result != 0 {
		time.Sleep(time.Second * time.Duration(5))
		this.sendGameAuto()
		return
	}
}

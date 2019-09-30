package main

import (
	"logs"
	"math/rand"
	"time"

	"github.com/tidwall/gjson"
)

type ExtRobotClient struct {
	RobotClient

	Status       int
	LastOutCards LastOutCards
	NextPlayer   int32
	IsDan        int32 //报单玩家下标
	SeatId       int32
	HandCard     []byte
}

func (this *ExtRobotClient) Rest() {
	this.NextPlayer = -1
	this.IsDan = -1
	this.HandCard = []byte{}
	this.LastOutCards = LastOutCards{}
}

//=======================对外方法=====================
// 机器人开始
func (this *ExtRobotClient) Start() {
	this.ChildPtr = this
	this.IsStop = false
	this.GameIn = false
	this.IsConnect = false
	//初始化金币
	// time.Sleep(time.Second * time.Duration(1))
	// minCoin := int64(gameConfig.getGameConfigInt("shiftRobotCoinMin")) * 100
	// maxCoin := int64(gameConfig.getGameConfigInt("shiftRobotCoinMax")) * 100
	rand.Seed(time.Now().UnixNano())
	// this.Coin = int64(rand.Int63n(maxCoin-minCoin)+minCoin) * 100
	this.Coin = int64(gameConfig.getGameConfigCoin("initcoin"))

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
	//百人经典牛牛定义
	//接收匹配更多信息返回
	this.Handle[MSG_GAME_INFO_AUTO_REPALY] = this.HandleAutoInfo
	//接收状态信息
	this.Handle[MSG_GAME_INFO_STAGE] = this.HandleStage
	//接受发牌信息
	this.Handle[MSG_GAME_INFO_SENDCARD] = this.HandleSendCards
	//接收玩家出牌广播
	this.Handle[MSG_GAME_INFO_OUTCARD_BRO] = this.HandleOutCards
	//接收玩家过广播
	this.Handle[MSG_GAME_INFO_PASS_BRO] = this.HandlePass
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
		logs.Error("匹配失败，服务器繁忙。。。。。10s")
		time.Sleep(time.Second * time.Duration(10))
		this.sendGameAuto()
		return
	} else if result == 11 {
		// 玩家已在游戏中
		logs.Error("匹配失败，玩家正在游戏")
		time.Sleep(time.Second * time.Duration(5))
		controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
		return
	} else if result == 9 {
		// 匹配失败，金币不足
		logs.Error("匹配失败，金币不足")
		time.Sleep(time.Second * time.Duration(5))
		controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
		return
	}
	if result != 0 {
		logs.Error("匹配失败，5秒后重新匹配")
		time.Sleep(time.Second * time.Duration(5))
		this.sendGameAuto()
		return
	}
}

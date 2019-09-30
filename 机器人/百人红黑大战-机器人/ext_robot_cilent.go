package main

import (
	"logs"
	"math/rand"
	"time"

	"github.com/tidwall/gjson"
)

type ExtRobotClient struct {
	RobotClient

	BetList  []int64 // 下注筹码列表
	AreaList []int
	PlayList []int
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
	this.CheckOnLineId = this.TimeTicker.AddTimer(30, func(int, interface{}) {
		this.RestChekIsNo()
	}, nil)
	go func() {
		//todo:0005此处添加获取机器人账号
		this.requstAccountInfo(this.Coin)
		this.OnMessageHandle()
	}()
}
func (this *ExtRobotClient) RestChekIsNo() {
	time.Sleep(time.Second * 30)
	controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this)
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
	//百人红黑大战
	this.Handle[MSG_GAME_INFO_AUTO_REPLY] = this.AutoReplyFinal
	this.Handle[MSG_GAME_INFO_RECONNECT_REPLY] = this.RoomInfo
	this.Handle[MSG_GAME_INFO_BET_NOTIFY] = this.BetTime

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
	if result != 0 {
		ErrorLog("请求匹配房间失败归还机器人重新请求账号 [%s]", gjson.Get(d, "Err").String())
		controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this)
		time.Sleep(time.Second * 2)
		return
	}
	this.AddMsgNative(MSG_GAME_INFO_INTO, struct{ Id int }{
		Id: MSG_GAME_INFO_INTO,
	})
	logs.Debug("匹配成功")
}

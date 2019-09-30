package main

import (
	"math/rand"
	"time"

	// "sync"
	"github.com/tidwall/gjson"
)

type ExtRobotClient struct {
	RobotClient

	ChairId    int32   //座位号
	CallPlayer int32   //叫牌玩家
	MinCoin    int64   //最小下注
	Round      int     //当前轮数
	HandCard   []int   //手牌
	CardLv     int     //手牌等级
	Alive      []int32 //存活玩家
	Leaves     []int32 //离开玩家
	SignRole   []int   //标记角色 0机器人，1玩家
	CardType   int     // 0未看牌 1已看牌
	PayCoin    int64   //记录下注筹码数
	AllCoin    int64   //桌面筹码

	MaxIsRobot     int         //0玩家 1机器人
	MaxHandCard    []int       //最大的手牌
	MaxCardLv      int         //牌型
	MaxChairId     int32       //座位号
	WinnerRole     int         // 0玩家，1机器人
	PlayerHandCard []PHandCard //真实玩家手牌

	IsChange int //0无 1有
	Special  int //0关，1开
	IsDie    int //0 false, 1true

	RobotChange  int //换牌概率
	IsCardChange int // 0.未换牌 1.已换牌
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

	//ext_conf
	// this.RobotChange = GExtConfig.RobotChange[gameConfig.GCEnterGame.Gradetype-1]
	// if this.RobotChange == 0 {
	// 	this.RobotChange = 100
	// }
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
	this.Handle[MSG_GAME_INFO_AUTO_REPLY] = this.HandleAuto
	this.Handle[MSG_GAME_INFO_STAGE] = this.HandlerGameStatus
	this.Handle[MSG_GAME_INFO_CALLPLAYER_REPLY] = this.HandleCallPlayer
	this.Handle[MSG_GAME_SETTLE] = this.HandlerSettle
	this.Handle[MSG_GAME_GIVE_UP] = this.HandleGiveUp
	this.Handle[MSG_GAME_LOOK_CARD] = this.HandleLookCard
	this.Handle[MSG_GAME_INFO_PLAY_INFO_REPLY] = this.HandleContestLoser
	this.Handle[MSG_GAME_INFO_MAX] = this.HandleGetMax
	this.Handle[MSG_GAME_INFO_CHANGE_CARD] = this.HandleChangeCard
	this.Handle[MSG_GAME_COIN] = this.HandleGetCoinMsg
	this.Handle[MSG_GAME_INFO_LEAVE_REPLY] = this.HandleLeave
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
	DebugLog("匹配房间消息结果", d)
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
		DebugLog("匹配结果失败")
		time.Sleep(time.Second * time.Duration(5))
		this.sendGameAuto()
		return
	}
}

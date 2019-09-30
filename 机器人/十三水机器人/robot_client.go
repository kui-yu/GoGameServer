package main

import (
	"encoding/json"
	"fmt"
	"logs"
	"sync"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/tidwall/gjson"
)

type RobotClient struct {
	IsStop    bool // 是否已停止
	IsConnect bool // 是否连接了网络

	HttpUrl    string
	HttpOrigin string // 网络连接服务器 https://www.abqpht.com

	HallToken string
	UserInfo  UserInfo
	Coin      int64 // 用户金币，用于归还机器人时调用

	SendMsg chan *SendMsg
	RecvMsg chan *RecvMsg
	Ws      *websocket.Conn
	Handle  map[uint32]func(string)

	EventMsg    chan *EventMsg
	EventHandle map[int32]func(interface{})

	TimeTicker *TimeTicker
	SysTicker  *time.Ticker
	RWMutex    sync.RWMutex // 读取锁
	//
	GameIn bool //是否在游戏中
	//

	HeartTimerId int // 心跳定时器id
	ChildPtr     interface{}
}

func (this *RobotClient) IsRobotStop() bool {
	isStop := false
	this.RWMutex.Lock()
	isStop = this.IsStop
	this.RWMutex.Unlock()
	return isStop
}

func (this *RobotClient) AddMsgNative(id uint32, d interface{}) {
	this.RWMutex.Lock()
	if this.IsStop {
		this.RWMutex.Unlock()
		return
	}
	omsg := SendMsg{
		Id: id,
	}
	if d != nil {
		omsg.Data, _ = json.Marshal(d)
	}

	select {
	case this.SendMsg <- &omsg:

	default:
		ErrorLog("丢失信息")
	}

	this.RWMutex.Unlock()
}

func (this *RobotClient) AddEventNative(id uint32, d interface{}) {
	omsg := EventMsg{
		Id:   int32(id),
		Data: d,
	}

	if this.IsRobotStop() {
		return
	}

	select {
	case this.EventMsg <- &omsg:

	default:
		ErrorLog("丢失信息")
	}
}

func (this *RobotClient) OpenHeart() {
	this.TimeTicker.AddTimer(10, func(id int, d interface{}) {
		this.AddMsgNative(MSG_HALL_HEART, struct {
			Id int
		}{
			Id: MSG_HALL_HEART,
		})

		this.OpenHeart()
	}, nil)
}

func (this *RobotClient) CloseHeart() {
	if this.HeartTimerId != 0 {
		this.TimeTicker.DelTimer(this.HeartTimerId)
		this.HeartTimerId = 0
	}
}

//=======================网络消息=====================
// 接收消息
func (this *RobotClient) RecvMessageThread() {
	go func() {
		for {
			// 获取数据
			var strRecv string
			err := websocket.Message.Receive(this.Ws, &strRecv)
			if err != nil {
				//断线 跳出循环，协程关闭
				DebugLog("接收消息错误，重新请求账号")
				this.IsConnect = false
				break
			}

			// 数据解密
			strRecv, err = Dncrypt(strRecv)
			id := uint32(gjson.Get(strRecv, "Id").Int())

			if this.IsRobotStop() {
				return
			}

			this.RWMutex.Lock()
			if this.IsStop {
				this.RWMutex.Unlock()
				return
			}
			select {
			case this.RecvMsg <- &RecvMsg{
				Id:   id,
				Data: strRecv,
			}:

			default:
			}
			this.RWMutex.Unlock()
		}

		if this.IsRobotStop() == false {
			time.Sleep(time.Second * 2)
			controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this.ChildPtr)
		}
	}()
}

// 服务端发送消息到客户端
func (this *RobotClient) SendMessageThread() {
	go func() {
		for s := range this.SendMsg {
			if this.IsRobotStop() {
				break
			}
			str, _ := Encrypt(s.Data)

			if _, err := this.Ws.Write([]byte(str)); err != nil {
				ErrorLog("发送消息错误", err, string(s.Data))
				break
			}
		}

		this.IsConnect = false

		if this.IsRobotStop() == false {
			time.Sleep(time.Second * 2)
			controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this.ChildPtr)
		}

	}()

}

// 发送消息到处理层
func (this *RobotClient) OnMessageHandle() {
	this.SysTicker = time.NewTicker(time.Second)
	for {
		if this.IsRobotStop() {
			return
		}
		select {
		case s := <-this.RecvMsg:
			if this.IsRobotStop() {
				return
			}
			h, ok := this.Handle[s.Id]
			if ok {
				h(s.Data)
			}
		case e := <-this.EventMsg:
			if this.IsRobotStop() {
				return
			}
			eh, ok := this.EventHandle[e.Id]
			if ok {
				eh(e.Data)
			}

		case <-this.SysTicker.C:
			if this.IsRobotStop() {
				return
			}
			this.TimeTicker.DoTimer()
		}
	}
}

// 连接网络
func (this *RobotClient) Connect() bool {
	for {
		if this.IsRobotStop() {
			return false
		}

		this.IsConnect = false
		DebugLog("userInfo", this.UserInfo)

		var err error
		this.Ws, err = websocket.Dial(this.HttpUrl, "", this.HttpOrigin)

		if err != nil {
			this.IsConnect = false
			ErrorLog("连接网络失败，重新请求账号", this.HttpUrl, this.HttpOrigin)
			DebugLog("错误信息", err)
			time.Sleep(time.Second * 2)
			controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this.ChildPtr)
			return false
		}

		this.IsConnect = true
		DebugLog("游戏服务器连接成功 %s", this.HttpUrl, this.HttpOrigin)

		this.RecvMessageThread()
		this.SendMessageThread()

		this.AddEventNative(EVENT_CONNECT_SUCCESS, nil)

		return true
	}
}

// 停止
func (this *RobotClient) BaseStop() {
	if this.SysTicker != nil {
		this.SysTicker.Stop()
	}

	this.TimeTicker.ClearTimer()

	defer close(this.EventMsg)
	defer close(this.RecvMsg)
	defer close(this.SendMsg)
	if this.IsConnect {
		err := this.Ws.Close()

		if err != nil {
			ErrorLog("关闭连接错误", err)
		}
	}
}

//获取机器人账号，然后连接大厅
func (this *RobotClient) requstAccountInfo(initCoin int64) {
	url := GRobotServer.getRobotUrl()
	token := GRobotServer.getRobotToken()
	hallId := GRobotServer.getRobotHallId()

	gameId := gameConfig.GCEnterGame.Gametype
	gradeId := gameConfig.GCEnterGame.Gradetype

	DebugLog("获取机器人参数 url:%s token:%s hallId:%d gameId:%d gradeId:%d\n", url, token, hallId, gameId, gradeId)

	args := fmt.Sprintf("?gameId=%d&gradeId=%d&hallId=%d&coin=%d", gameId, gradeId, hallId, initCoin)
	data, err := SendRequest(url+args, nil, "GET", token)

	if err == nil {
		DebugLog("得到机器人账号信息", data)

		if gjson.Get(data, "code").Int() != 200 {
			ErrorLog("获得机器人错误信息:%s\n", gjson.Get(data, "msg").String())

			this.TimeTicker.AddTimer(1, func(id int, d interface{}) {
				this.requstAccountInfo(initCoin)
			}, nil)
			return
		}

		d := gjson.Get(data, "data")
		hallIp := d.Get("hallIp").String()
		hallPort := d.Get("port").Int()
		this.HallToken = d.Get("token").String()

		if gameConfig.GCGameServer.IsCustomHallConnect {
			this.HttpOrigin = fmt.Sprintf("http://%s:%d", gameConfig.GCGameServer.HallIp, gameConfig.GCGameServer.HallPort)
			this.HttpUrl = fmt.Sprintf("ws://%s:%d/hall", gameConfig.GCGameServer.HallIp, gameConfig.GCGameServer.HallPort)
		} else if hallIp[0:4] == "hall" {
			this.HttpOrigin = fmt.Sprintf("https://%s:%d", gameConfig.GCGameServer.WebsocketUrl, 443)
			this.HttpUrl = fmt.Sprintf("wss://%s:%d/%s", gameConfig.GCGameServer.WebsocketUrl, 443, hallIp)
		} else {
			this.HttpOrigin = fmt.Sprintf("http://%s:%d", hallIp, hallPort)
			this.HttpUrl = fmt.Sprintf("ws://%s:%d/hall", hallIp, hallPort)
		}

	} else {
		ErrorLog("得到机器人账号信息错误", err)

		this.TimeTicker.AddTimer(1, func(id int, d interface{}) {
			this.requstAccountInfo(initCoin)
		}, nil)
		return
	}
	//此处获取机器人账号成功，开始连接大厅
	this.Connect()
}

/////////////////////////////////////////////////////////////////////////
//发送登陆大厅的请求
func (this *RobotClient) sendLogin() {
	DebugLog("发送大厅登录")

	this.AddMsgNative(MSG_HALL_ROBOT_LOGIN, struct {
		Id  int32  //协议号
		Gid string //web 登录后获取的 token
	}{
		Id:  MSG_HALL_ROBOT_LOGIN,
		Gid: this.HallToken,
	})
}

///////////////////////////////////////////////////////////
// 通知连接成功
func (this *ExtRobotClient) onConnected(d interface{}) {
	// 开启心跳
	this.OpenHeart()
	//登陆游戏大厅
	this.sendLogin()
}

//登陆大厅应答
func (this *ExtRobotClient) onLogin(d string) {
	DebugLog("接收到消息", d)
	logs.Debug("我的金币:", this.Coin)
	// 登录错误，重新走登录流程
	if gjson.Get(d, "Result").Int() != 0 {
		ErrorLog("【登录错误】 token:%s err:%s", this.HallToken, d)
		time.Sleep(time.Second * 2)
		controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this)
		return
	}
	DebugLog("【登录成功】 token:", this.HallToken)

	this.UserInfo.GameServerId = 0
	this.UserInfo.Uid = gjson.Get(d, "Uid").Int()

	// 是否走重新连接游戏，还是进入游戏
	this.sendReEnterGame()
}

func (this *ExtRobotClient) sendReEnterGame() {
	this.AddMsgNative(MSG_HALL_JOIN_GMAE, struct {
		Id        int32 //协议号
		GameType  int32 //游戏类型
		RoomType  int32 //房间类型
		GradeType int32 //场次类型
	}{
		Id:        MSG_HALL_JOIN_GMAE,
		GameType:  int32(gameConfig.GCEnterGame.Gametype),
		RoomType:  int32(gameConfig.GCEnterGame.Roomtype),
		GradeType: int32(gameConfig.GCEnterGame.Gradetype),
	})
}

// 接收到进入游戏消息
func (this *ExtRobotClient) onJoinGame(d string) {
	if gjson.Get(d, "Result").Int() != 0 {
		ErrorLog("进入游戏失败 [%s]", gjson.Get(d, "Err").String(), this.Coin)
		time.Sleep(time.Second * 2)
		controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this)
		return
	}

	DebugLog("进入游戏成功")
	//发送匹配
	this.sendGameAuto()
}

//发送匹配请求
func (this *ExtRobotClient) sendGameAuto() {
	this.AddMsgNative(MSG_GAME_AUTO, struct {
		Id        int32 //协议号
		IsRobot   bool  // 是否机器人
		GameType  int32 //游戏类型
		RoomType  int32 //房间类型
		GradeType int32 //场次类型
	}{
		Id:        MSG_GAME_AUTO,
		IsRobot:   true,
		GameType:  int32(gameConfig.GCEnterGame.Gametype),
		RoomType:  int32(gameConfig.GCEnterGame.Roomtype),
		GradeType: int32(gameConfig.GCEnterGame.Gradetype),
	})
}

// 连接机器人服务器
package main

import (
	"encoding/json"
	"fmt"

	// "strconv"
	"time"

	"github.com/tidwall/gjson"

	"code.google.com/p/go.net/websocket"
)

var GRobotServer RobotServer = RobotServer{
	SendMsg:       make(chan *SendMsg, 1000),
	RecvMsg:       make(chan *RecvMsg, 1000),
	isInitSuccess: false,
	Handle:        make(map[uint32]func(string)),
}

type RobotServer struct {
	SendMsg chan *SendMsg
	RecvMsg chan *RecvMsg

	Handle map[uint32]func(string)
	Ws     *websocket.Conn
	Info   RepBgInfo

	isInitSuccess bool
}

func (this *RobotServer) getRobotUrl() string {
	return this.Info.BgRobotGetUrl
}

func (this *RobotServer) getRobotRestUrl() string {
	return this.Info.BgRobotRestUrl
}

func (this *RobotServer) getRobotRestAllUrl() string {
	return this.Info.BgRobotRestAllUrl
}

func (this *RobotServer) getRobotAddCoinUrl() string {
	return this.Info.BgRobotAddCoinUrl
}

func (this *RobotServer) getRobotToken() string {
	return this.Info.BgRobotToken
}

func (this *RobotServer) getRobotHallId() int {
	return this.Info.BgRobotHallId
}

func (this *RobotServer) AddMsgNative(id uint32, d interface{}) {
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
}

// 接收消息
func (this *RobotServer) RecvMessageThread() {
	go func() {
		for {
			// 获取数据
			var strRecv string
			err := websocket.Message.Receive(this.Ws, &strRecv)
			if err != nil {
				//断线 跳出循环，协程关闭
				break
			}

			// 数据解密
			strRecv, err = Dncrypt(strRecv)
			DebugLog("接收数据：", strRecv)
			id := uint32(gjson.Get(strRecv, "Id").Int())

			this.RecvMsg <- &RecvMsg{
				Id:   id,
				Data: strRecv,
			}
		}

		this.Connect()
	}()
}

// 服务端发送消息到客户端
func (this *RobotServer) SendMessageThread() {
	go func() {
		for s := range this.SendMsg {
			str, _ := Encrypt(s.Data)
			if _, err := this.Ws.Write([]byte(str)); err != nil {
				ErrorLog("发送消息错误", err)
				continue
			}
		}
		this.Ws.Close()
		close(this.SendMsg)

		this.Connect()
	}()

}

// 发送消息到处理层
func (this *RobotServer) OnMessageHandle() {
	go func() {
		for s := range this.RecvMsg {
			h, ok := this.Handle[s.Id]
			if ok {
				h(s.Data)
			}
		}
	}()
}

func (this *RobotServer) addHandler() {
	//机器人上线应答
	this.Handle[MSG_GAME_ROBOT_ONLINE] = this.onRecvRobotOnline
}

//注册机器人的配置到机器人管理中心,供后台获取显示和修改控制
func (this *RobotServer) sendOnlineData() {
	url := fmt.Sprintf("http://%s:%d", gameConfig.OpenWebInterface.Ip, gameConfig.OpenWebInterface.Port)
	info := &ResGameInfo{
		Id:                MSG_GAME_ROBOT_ONLINE,
		GroupId:           uint32(gameConfig.GCEnterGame.Roomtype),
		GameId:            uint32(gameConfig.GCEnterGame.Gametype),
		GradeId:           uint32(gameConfig.GCEnterGame.Gradetype),
		Name:              gameConfig.GCEnterGame.Name,
		GetRobotConfigUrl: url + "/" + gameConfig.OpenWebInterface.GetRobotConfigUrl,
		PutRobotConfigUrl: url + "/" + gameConfig.OpenWebInterface.PutRobotConfigUrl,
		CheckRobotUrl:     url + "/" + gameConfig.OpenWebInterface.CheckRobotUrl,
		OfflineRobotUrl:   url + "/" + gameConfig.OpenWebInterface.OfflineRobotUrl,
		Forceoffroboturl:  url + "/" + gameConfig.OpenWebInterface.Forceoffroboturl,
		Forceonroboturl:   url + "/" + gameConfig.OpenWebInterface.Forceonroboturl,
		RobotCount:        uint32(gameConfig.getGameConfigInt("num")),
	}

	this.AddMsgNative(MSG_GAME_ROBOT_ONLINE, info)
}

func (this *RobotServer) ResetAllRobot() {
	allurl := GRobotServer.getRobotRestAllUrl()
	DebugLog("归还所有机器人", allurl)
	str, err := SendRequest(allurl, nil, "GET", this.getRobotToken())
	if err != nil {
		ErrorLog("归还所有机器人错误", err)
	} else {
		DebugLog("归还所有机器人结果", str)
	}
}

// 接收到配置信息
//此处添加机器人
func (this *RobotServer) onRecvRobotOnline(data string) {
	DebugLog("接收到消息 ", data)

	json.Unmarshal([]byte(data), &this.Info)

	// 归还所有机器人
	// this.ResetAllRobot()

	controller.sendEvent(EVENT_CONT_ADDROBOT, gameConfig.getGameConfigInt("num"))
}

//连接机器人管理中心服务器
//地址端口在配置文件里面配置
func (this *RobotServer) Connect() bool {
	for {
		hallip := gameConfig.GCRobotManager.Ip + ":" + gameConfig.GCRobotManager.Port
		origin := "http://" + hallip
		url := "ws://" + hallip + "/robotserver"

		var err error
		GRobotServer.Ws, err = websocket.Dial(url, "", origin)
		//连接失败，3秒后重新连接
		if err != nil {
			ErrorLog("机器人服务器连接失败3秒后重连 ", err)
			time.Sleep(3 * time.Second)
			return this.Connect()
		}

		DebugLog("机器人服务器连接成功 %s", url)
		//开启与机器人中心服务器的发送与接收线程
		this.RecvMessageThread()
		this.SendMessageThread()
		// 发送机器人上线给后台
		this.sendOnlineData()
		return true
	}
}

//初始化结构体，协议，还有和机器人中心服务器的连接
func (this *RobotServer) Init() {
	this.addHandler()
	this.OnMessageHandle()
	this.Connect()
}

func StartRobotServer() {
	GRobotServer.Init()
}

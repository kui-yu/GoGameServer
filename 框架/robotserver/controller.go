package main

import (
	"encoding/json"
	"fmt"
	"logs"

	//	"strconv"

	"code.google.com/p/go.net/websocket"
	"github.com/tidwall/gjson"
)

type RobotClient struct {
	SendMsg chan *SendMsg
	RecvMsg chan *RecvMsg
	Conn    *websocket.Conn
	Handle  map[uint32]func(string)

	Closed bool

	IsLogin           bool
	GroupId           uint32
	GameId            uint32
	GradeId           uint32
	Name              string
	GetRobotConfigUrl string
	PutRobotConfigUrl string
	CheckRobotUrl     string
	OfflineRobotUrl   string
	Forceoffroboturl  string
	Forceonroboturl   string
	RobotCount        uint32
}

var robotClients []*RobotClient

func (this *RobotClient) AddMsgNative(id uint32, d interface{}) {
	fmt.Println("发送消息AddMsgNative")
	omsg := SendMsg{
		Id: id,
	}
	if d != nil {
		omsg.Data, _ = json.Marshal(d)
	}

	select {
	case this.SendMsg <- &omsg:
	default:
		logs.Debug("丢失信息")
	}
}

// 从客户端接收消息
func (this *RobotClient) RecvMessageThread() {
	for {
		// 获取数据
		var strRecv string
		err := websocket.Message.Receive(this.Conn, &strRecv)
		if err != nil {
			//断线 跳出循环，协程关闭
			fmt.Println("关闭连接")
			RemoveClient(this)
			return
		}

		// 数据解密
		strRecv, err = Dncrypt(strRecv)
		logs.Debug("接收数据：", strRecv)
		id := uint32(gjson.Get(strRecv, "Id").Int())

		this.RecvMsg <- &RecvMsg{
			Id:   id,
			Data: strRecv,
		}
	}
}

// 服务端发送消息到客户端
func (this *RobotClient) SendMessageThread() {
	go func() {
		for s := range this.SendMsg {
			str, _ := Encrypt(s.Data)
			fmt.Println("发送消息到机器人")
			if _, err := this.Conn.Write([]byte(str)); err != nil {
				logs.Debug("发送消息错误", err)
				RemoveClient(this)
				return
			}
		}
	}()
}

// 服务端发送消息到客户端
func (this *RobotClient) OnMessageHandle() {
	go func() {
		for s := range this.RecvMsg {
			h, ok := this.Handle[s.Id]
			if ok {
				h(s.Data)
			}
		}
	}()
}

func (this *RobotClient) Stop() {
	this.removeHandler()

	this.Conn.Close()

	close(this.RecvMsg)
	close(this.SendMsg)
}

// socket 连接触发，产生两个协程  RecvMessageThread 与 SendMessageThread
func DoHandler(ws *websocket.Conn) {
	client := &RobotClient{Conn: ws,
		SendMsg: make(chan *SendMsg, 1000),
		RecvMsg: make(chan *RecvMsg, 1000),
		Handle:  make(map[uint32]func(string)),
		IsLogin: false,
	}

	robotClients = append(robotClients, client)

	client.addHandler()
	client.OnMessageHandle()
	client.SendMessageThread()
	client.RecvMessageThread()
}

func FindClient(groupId uint32, gameId uint32, gradeId uint32) *RobotClient {
	for _, v := range robotClients {
		if v.GroupId == groupId && v.GameId == gameId && v.GradeId == gradeId {
			return v
		}
	}
	return nil
}

func RemoveClient(robot *RobotClient) {
	for i, client := range robotClients {
		if client == robot {
			client.Stop()
			robotClients = append(robotClients[:i], robotClients[i+1:]...)
			break
		}
	}
}

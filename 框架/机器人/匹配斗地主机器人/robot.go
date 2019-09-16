package main

import (
	"encoding/json"
	"fmt"

	//	"time"

	"logs"
	"strconv"
	"sync"

	"code.google.com/p/go.net/websocket"
	"github.com/tidwall/gjson"
)

type DoMsg struct {
	Handle int32
	Msg    []byte
}

type Roboter struct {
	firstCall int32 //第一个叫地主的人
	Account   string
	Uid       int64
	Name      string
	Coins     int32
	Sex       int32
	HandCard  []byte
	SeatId    int32
	CurId     int32
	Handle    map[int32]func([]byte)
	Ws        *websocket.Conn
	XuHao     int64
	ShowId    int64

	MMsgChan chan *DoMsg

	MTimerChan chan bool
	MTimerNum  int
	MGameState int32
	MJiaoFen   int

	// Cards map[int]CardGroupInfo
	Lk      sync.Mutex // 互斥锁
	ChairId int32      // 座位号
}

func (this *Roboter) SInit() {
	this.MMsgChan = make(chan *DoMsg)
	this.Handle = make(map[int32]func([]byte))
	this.MTimerChan = make(chan bool)
	this.MTimerNum = 0
	this.MGameState = 0
	this.MJiaoFen = 0
	this.InitHandle()
}

func (this *Roboter) Start(account string, xuhao, uid int64) {
	//游客登录，获取uid和token（guid）
	var err error
	this.Account = account
	this.XuHao = xuhao
	this.ShowId = uid
	// 登录大厅
	hallip := GCONFIG.Ip + ":" + strconv.Itoa(GCONFIG.Port)
	origin := "http://" + hallip
	url := "ws://" + hallip + "/hall"
	this.Ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		logs.Debug("websocket连接失败 ", err)
		return
	}

	//登录
	var reqlogin HMsgHallLogin
	reqlogin.Account = this.Account
	reqlogin.Gid = ""
	reqlogin.Id = MSG_HALL_LOGIN // 登录Id号

	jsv, _ := json.Marshal(reqlogin)
	strRecv, _ := Encrypt(jsv)
	// logs.Debug("发送登录消息：", string(jsv))
	if _, err := this.Ws.Write([]byte(strRecv)); err != nil {
		logs.Debug("ws.Write ", err)
		return
	}

	// 循环接收数据
	for {
		var recvdata []byte
		err := websocket.Message.Receive(this.Ws, &recvdata)
		if err != nil {
			logs.Debug("读取数据失败", err)
			break
		}
		str, _ := Dncrypt(string(recvdata))

		id := int32(gjson.Get(str, "Id").Int())

		msg := DoMsg{
			Handle: id,
		}
		msg.Msg = []byte(str)
		if this.MMsgChan != nil {
			this.MMsgChan <- &msg
		} else {
			logs.Debug("发送消息给管理器失败", msg)
		}
	}
	this.Ws.Close()
}

//消息转发
func (this *Roboter) HandleMsg() {
	for {
		select {
		case msg := <-this.MMsgChan:
			funchandle, ok := this.Handle[msg.Handle]
			if ok {
				funchandle(msg.Msg)
			}
		case <-this.MTimerChan:
			if this.MTimerNum != 0 {
				this.MTimerNum--
				if this.MTimerNum == 0 {
					//this.Do()
				}
			}
		}
	}
}

//发送消息给服务端
func (this *Roboter) SendToServer(v interface{}) {
	msg, _ := json.Marshal(v)
	strRecv, _ := Encrypt(msg)
	if _, err := this.Ws.Write([]byte(strRecv)); err != nil {
		fmt.Println("发送消息给客户端失败 ", err)
	}
}

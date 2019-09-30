package main

import (
	"encoding/json"
	"fmt"

	//	"time"

	"strconv"

	"logs"

	"code.google.com/p/go.net/websocket"
	"github.com/tidwall/gjson"
)

type DoMsg struct {
	Handle int32
	Msg    []byte
}

type Roboter struct {
	Account  string
	Uid      int64
	Name     string
	Coins    int32
	Sex      int32
	HandCard []int
	SeatId   int32
	CurId    int32
	Handle   map[int32]func([]byte)
	Ws       *websocket.Conn
	XuHao    int64
	ShowId   int64

	MMsgChan   chan *DoMsg
	MTimerChan chan bool
	MTimerNum  int
	MGameState int32
	MJiaoFen   int
	Banker     int
	TList      []*Timer
	QueColor   int
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
		// logs.Debug("接收的数据:", str)
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
			this.DoTimer()
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

type Timer struct {
	Id int
	H  func(interface{})
	T  int //定时时间
	D  interface{}
}

//定时器
func (this *Roboter) DoTimer() {
	if len(this.TList) == 0 {
		return
	}
	nlist := []*Timer{}
	olist := []*Timer{}
	for _, v := range this.TList {
		v.T--
		if v.T <= 0 {
			olist = append(olist, v)
		} else {
			nlist = append(nlist, v)
		}
	}
	this.TList = nlist
	for _, v := range olist {
		v.H(v.D)
	}
}

func (this *Roboter) AddTimer(id int, t int, h func(interface{}), d interface{}) {
	this.TList = append(this.TList, &Timer{
		Id: id,
		H:  h,
		T:  t,
		D:  d,
	})
}

//同一id的定时器只能存在一个
func (this *Roboter) AddUniueTimer(id int, t int, h func(interface{}), d interface{}) {
	for i := len(this.TList) - 1; i >= 0; i-- {
		if this.TList[i].Id == id {
			this.TList = append(this.TList[:i], this.TList[i+1:]...)
		}
	}
	this.AddTimer(id, t, h, d)
}

func (this *Roboter) DelTimer(id int) {
	for i, v := range this.TList {
		if v.Id == id {
			this.TList = append(this.TList[:i], this.TList[i+1:]...)
			break
		}
	}
}

func (this *Roboter) GetTimerNum(id int) int {
	for _, v := range this.TList {
		if v.Id == id {
			return v.T
		}
	}
	return 0
}

//清空定时器
func (this *Roboter) ClearTimer() {
	this.TList = []*Timer{}
}

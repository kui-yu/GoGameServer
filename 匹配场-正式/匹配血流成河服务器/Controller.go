package main

import (
	"encoding/json"
	"logs"

	"time"

	//	"strconv"
	"sync"

	"code.google.com/p/go.net/websocket"
	"github.com/tidwall/gjson"
)

type Contoller struct {
	Uid         int64 //玩家id
	Account     string
	Nick        string
	Coin        int64
	Sex         int32
	Lv          int32
	GameId      int32
	ServerId    int32
	Token       string
	Head        string
	Robot       bool //是否是机器人,true是 ， false否
	HierarchyId int
	RoomCard    int64
	//
	OutMsg chan *OutMsg //桌子-》controller
	Conn   *websocket.Conn
	Closed bool
	Lk     sync.Mutex
}

func (this *Contoller) AddMsgNative(id int32, d interface{}, force bool) {
	//
	omsg := OutMsg{
		Id: id,
	}
	if d != nil {
		omsg.Data, _ = json.Marshal(d)
	}

	this.Lk.Lock()
	if this.Closed {
		this.Lk.Unlock()
		return
	}
	if id == MSG_HALL_LEAVE_REPLY {
		this.Closed = true
	}
	if force {
		this.OutMsg <- &omsg
	} else {
		select {
		case this.OutMsg <- &omsg:
		default:
			logs.Error("missing data:", id, d)
		}
	}
	this.Lk.Unlock()
}

func (this *Contoller) AddMsg(id int32, d []byte, force bool) {
	this.Lk.Lock()
	if this.Closed {
		this.Lk.Unlock()
		return
	}
	if id == MSG_HALL_LEAVE_REPLY {
		this.Closed = true
	}
	omsg := OutMsg{
		Id:   id,
		Data: d,
	}
	if force {
		this.OutMsg <- &omsg
	} else {
		select {
		case this.OutMsg <- &omsg:
		default:
			logs.Debug("丢失信息")
		}
	}
	this.Lk.Unlock()
}

// 从客户端接收消息
func (this *Contoller) PullFromClient() {
	for {
		// 获取数据
		var strRecv string
		err := websocket.Message.Receive(this.Conn, &strRecv)
		if err != nil {
			//断线
			if this.Uid == 0 {
				this.AddMsgNative(MSG_HALL_LEAVE_REPLY, nil, true)
			} else {
				GHall.AddMsg(&InMsg{
					Id:  MSG_HALL_LEAVE,
					Uid: this.Uid,
					Col: nil,
				})
			}
			//跳出循环，协程关闭
			break
		}

		// 数据解密
		strRecv, err = Dncrypt(strRecv)
		// logs.Debug("接收数据：", strRecv)
		id := int32(gjson.Get(strRecv, "Id").Int())
		if id == MSG_HALL_LOGIN || id == MSG_HALL_ROBOT_LOGIN {
			if this.Uid != 0 {
				logs.Error("player already login", strRecv)
				continue
			}
			GHall.AddMsg(&InMsg{
				Id:   id,
				Uid:  this.Uid,
				Data: strRecv,
				Col:  this,
			})
		} else {
			if this.Uid == 0 {
				logs.Error("player no login", strRecv)
				continue
			} else {
				GHall.AddMsg(&InMsg{
					Id:   id,
					Uid:  this.Uid,
					Data: strRecv,
					Col:  this,
				})
			}
		}

	}
}

// 服务端发送消息到客户端
func (this *Contoller) PushToClient() {
	go func() {
		for s := range this.OutMsg {
			if s.Id == MSG_HALL_LEAVE_REPLY { //断线必须关掉此协程
				break
			}
			str, _ := Encrypt(s.Data)
			if _, err := this.Conn.Write([]byte(str)); err != nil {
				logs.Debug("发送消息错误", err)
				time.Sleep(time.Second)
				continue
			}
		}
		this.Conn.Close()
		close(this.OutMsg)
	}()
}

// socket 连接触发，产生两个协程  PushToClient 与 PullFromClient
func DoHandler(ws *websocket.Conn) {
	Connect := &Contoller{Conn: ws,
		OutMsg: make(chan *OutMsg, 1000),
	}

	Connect.PushToClient()
	Connect.PullFromClient()
}

package main

import (
	"crypto/rand"
	"encoding/json"
	"logs"
	"math/big"
)

///////////////////////////////////////////////////////
//玩家数据结构
type Player struct {
	Account     string   // 账号
	Uid         int64    // 用户ID
	Nick        string   // 个性签名
	Sex         int32    // 性别
	Head        string   // 头像
	Lv          int32    // 等级
	ChairId     int32    // 座位号
	Ready       bool     // 准备状态
	Dk          *ExtDesk // 桌子号
	Coins       int64    // 金币数
	LiXian      bool
	Token       string
	Robot       bool
	HierarchyId int32
	RoomCard    int64
	Mgr         *DeskMgr // 桌子管理器
}

func (this *Player) SetMgr(mg *DeskMgr) {
	this.Mgr = mg
}

func (this *Player) SetDesk(dk *ExtDesk) {
	this.Dk = dk
}

func (this *Player) SendNativeMsg(id int, data interface{}) {
	if this.Mgr != nil && !this.LiXian {
		this.Mgr.SendNativeMsg(id, this.Uid, data)
	}
}

func (this *Player) SendNativeMsgForce(id int, data interface{}) {
	if this.Mgr != nil {
		this.Mgr.SendNativeMsg(id, this.Uid, data)
	}
}

///////////////////////////////////////////////////////
// 真随机
func RandInt64(max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < 0 {
		return RandInt64(max)
	}
	return i.Int64()
}

///////////////////////////////////////////////////////
type Timer struct {
	Id int
	H  func(interface{})
	T  int //定时时间
	D  interface{}
}

//定时器
func (this *Desk) DoTimer() {
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

func (this *Desk) AddTimer(id int, t int, h func(interface{}), d interface{}) {
	this.TList = append(this.TList, &Timer{
		Id: id,
		H:  h,
		T:  t,
		D:  d,
	})
}

//同一id的定时器只能存在一个
func (this *Desk) AddUniueTimer(id int, t int, h func(interface{}), d interface{}) {
	for i := len(this.TList) - 1; i >= 0; i-- {
		if this.TList[i].Id == id {
			this.TList = append(this.TList[:i], this.TList[i+1:]...)
		}
	}
	this.AddTimer(id, t, h, d)
}

func (this *Desk) DelTimer(id int) {
	for i, v := range this.TList {
		if v.Id == id {
			this.TList = append(this.TList[:i], this.TList[i+1:]...)
			break
		}
	}
}

func (this *Desk) GetTimerNum(id int) int {
	for _, v := range this.TList {
		if v.Id == id {
			return v.T
		}
	}
	return 0
}

//清空定时器
func (this *Desk) ClearTimer() {
	this.TList = []*Timer{}
}

///////////////////////////////////////////////////////

type Desk struct {
	Handle    map[int32]func(*ExtPlayer, *DkInMsg) // 注册函数
	Id        int                                  // 桌号
	Players   []*ExtPlayer                         // 所有玩家信息和状态
	DeskMgr   *DeskMgr                             // 游戏管理器
	TList     []*Timer                             // 定时器列表
	GameState int                                  // 游戏状态
	MsgId     int32                                //消息自增id，防止中途丢包
	JuHao     string                               //局号
	//
	FkNo    string //房间号
	FkInfo  string //房间规则信息
	FkOwner int64  //房主uid
	FkTime  int64  //创建时间
	//
	RobotIn int //机器人进入开关,0:禁止进入，1禁止进入但开启打开定时器，2允许进入
}

// 初始化结构体

func (this *Desk) InitData() {
	this.Handle = make(map[int32]func(*ExtPlayer, *DkInMsg))
}

func (this *Desk) SetDeskMgr(p *DeskMgr) {
	this.DeskMgr = p
}

//初始化数值
func (this *Desk) DoHandle(p *ExtPlayer, msg *DkInMsg) {
	h, ok := this.Handle[msg.Id]
	if ok {
		h(p, msg)
	}
}

/////////////////////////////////////////////////////
//公共接口
func (this *Desk) GetPlayer(uid int64) *ExtPlayer {
	for _, v := range this.Players {
		if v.Uid == uid {
			return v
		}
	}
	return nil
}

//添加玩家
//返回具体人数
func (this *Desk) AddPlayer(p *ExtPlayer) int {
	if len(this.Players) >= GCONFIG.PlayerNum {
		return -1
	}
	//设置chairid
	doinsert := false
	for i, v := range this.Players {
		if i != int(v.ChairId) {
			doinsert = true
			p.ChairId = int32(i)
			nl := append([]*ExtPlayer{}, this.Players[:i]...)
			nl = append(nl, p)
			nl = append(nl, this.Players[i:]...)
			this.Players = nl
			break
		}
	}
	if !doinsert {
		p.ChairId = int32(len(this.Players))
		this.Players = append(this.Players, p)
	}
	//
	return len(this.Players)
}

func (this *Desk) DelPlayer(uid int64) *ExtPlayer {
	for i, v := range this.Players {
		if v.Uid == uid {
			this.Players = append(this.Players[:i], this.Players[i+1:]...)
			return v
		}
	}
	return nil
}

//
func (this *Desk) BroadcastAll(id int, d interface{}) {
	for _, v := range this.Players {
		v.SendNativeMsg(id, d)
	}
}

func (this *Desk) BroadcastOther(p *ExtPlayer, id int, d interface{}) {
	for _, v := range this.Players {
		if v.Uid == p.Uid {
			continue
		}
		v.SendNativeMsg(id, d)
	}
}

func (this *Desk) BroadcastSelf(p *ExtPlayer, id int, d interface{}) {
	p.SendNativeMsg(id, d)
}

/////////////////////////////////////////////////////
func (this *Desk) GetZouShi(fromId int32) {

}

func (this *Desk) GetServerId() int32 {
	return 1
}

//消息协议函数
func (this *Desk) Leave(p *ExtPlayer) bool {
	if this.GameState == GAME_STATUS_FREE {
		this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Result: 0,
			Token:  p.Token,
			Robot:  p.Robot,
		})
		this.DelPlayer(p.Uid)
		this.DeskMgr.LeaveDo(p.Uid)
	} else if this.GameState == GAME_STATUS_END {
		return true
	} else {
		// p.LiXian = true
		// this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
		// 	Id:     MSG_GAME_LEAVE_REPLY,
		// 	Result: 1,
		// 	Cid:    p.ChairId,
		// 	Uid:    p.Uid,
		// 	Err:    "玩家正在游戏中，不能离开",
		// })
		p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Result: 1,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Err:    "玩家正在游戏中，不能离开",
			Robot:  p.Robot,
		})
		return false
	}
	return true
}

func (this *Desk) GameOverLeave() {
	for _, v := range this.Players {
		logs.Debug("*******踢掉玩家信息：%v", v.Nick)
		if v.Uid == 1 {
			logs.Debug("玩家离开", v.Uid, this.GameState)
		}
		this.DeskMgr.LeaveDo(v.Uid)
	}
	this.Players = []*ExtPlayer{}
}

//系统更新金币
//通知玩家更新金币
func (this *Desk) UpdatePlayerInfo(p *ExtPlayer, d *DkInMsg) {
	req := GUpdatePlayerInfo{}
	json.Unmarshal([]byte(d.Data), &req)
	p.Coins += req.PlayerInfo.Coins
	p.Head = req.PlayerInfo.Head
	p.HierarchyId = req.PlayerInfo.HierarchyId
	p.RoomCard += req.PlayerInfo.RoomCard
	// push 有通知
	logs.Debug("充值或者修改头像推送成功")
	p.SendNativeMsg(MSG_GAME_UPDATEPLAYER_NOTIFY, &GUpdatePlayerNotify{
		Id:       MSG_GAME_UPDATEPLAYER_NOTIFY,
		Coin:     req.PlayerInfo.Coins,
		Head:     req.PlayerInfo.Head,
		RoomCard: req.PlayerInfo.RoomCard,
	})
}

func (this *Desk) LeaveByForce(p *ExtPlayer) {
	p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
		Id:     MSG_GAME_LEAVE_REPLY,
		Result: 0,
		Cid:    p.ChairId,
		Uid:    p.Uid,
		Token:  p.Token,
		Robot:  p.Robot,
	})
	this.DelPlayer(p.Uid)
	this.DeskMgr.LeaveDo(p.Uid)
}

func (this *Desk) GToHAddCoin(uid, coin int64) {
	if this.DeskMgr != nil {
		this.DeskMgr.SendNativeMsg(MSG_GAME_ADDCOIN_NOTIFY, uid, &GToHAddCoin{
			Id:   MSG_GAME_ADDCOIN_NOTIFY,
			Uid:  uid,
			Coin: coin,
		})
	}

}

//消耗房卡,card要传负数
func (this *Desk) GToHAddRoomCard(uid, card int64) {
	if this.DeskMgr != nil {
		this.DeskMgr.SendNativeMsg(MSG_GAME_ADDROOMCARD_NOTIFY, uid, &GToHAddRoomCard{
			Id:       MSG_GAME_ADDROOMCARD_NOTIFY,
			Uid:      uid,
			RoomCard: card,
		})
	}
}

////////////////////////////////////////////
//机器人允许进入定时器
const (
	TIMER_ENABLE_ROBOTIN = 90000
)

func (this *Desk) AddRobotInTimer(t int) {
	if this.RobotIn != 0 {
		return
	}
	this.RobotIn = 1
	this.AddTimer(TIMER_ENABLE_ROBOTIN, t, this.TimerRobotInTimer, nil)
}

//开启机器人允许进入标志
func (this *Desk) TimerRobotInTimer(d interface{}) {
	this.RobotIn = 2
}

////////////////////////////////////////////
//房卡场自动回收没有使用的房间
const (
	TIMER_FKBACKDESK     = 100000
	TIMER_FKBACKDESK_NUM = 600
)

func (this *ExtDesk) AddBackDeskTimer() {
	this.AddTimer(TIMER_FKBACKDESK, TIMER_FKBACKDESK_NUM, this.FkBackDeskTimer, nil)
}

func (this *ExtDesk) FkBackDeskTimer(d interface{}) {
	this.DeskMgr.BackDesk(this)
}

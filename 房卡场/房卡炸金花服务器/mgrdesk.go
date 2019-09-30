package main

import (
	"encoding/json"
	"logs"
	"os"
	"strconv"
	"time"
	//	"github.com/tidwall/gjson"
)

////////////////////////////////////
//读取游戏配置，失败就关掉服务器，提示
var G_DbGetGameServerData DbGetGameServerData

var G_DbRoleRateListData []DbHierarchyGameRate

func init() {
	GetGameInfoFromDb()
	GetRoleRateList()
}

////////////////////////////////////
var GDeskMgr DeskMgr

type DeskMgr struct {
	CInMsg     *chan *DkInMsg
	COutMsg    *chan *DkOutMsg
	Handle     map[int32]func(*DkInMsg)
	Free       []*ExtDesk
	All        []*ExtDesk
	MapPlayers map[int64]*ExtPlayer
	//
	MFkDesk map[string]*ExtDesk
	//
	ChangeStock int64 //一段时间内变化的库存值
	//
	DeskAllotModel int //0不参与分配，1中途加入离开模式(德州扑克)
}

func (this *DeskMgr) InitData() {
	this.Handle = make(map[int32]func(*DkInMsg))
	this.MapPlayers = make(map[int64]*ExtPlayer)
	this.MFkDesk = make(map[string]*ExtDesk)
	//
	this.InitHandle()
	for i := 0; i < GCONFIG.DeskNum; i++ {
		d := new(ExtDesk)
		//
		d.Id = i
		this.All = append(this.All, d)
		this.Free = append(this.Free, d)
		d.InitData()
		d.SetDeskMgr(&GDeskMgr)
		d.InitExtData()
	}
}

func (this *DeskMgr) SetInChan(c *chan *DkInMsg) {
	this.CInMsg = c
}

func (this *DeskMgr) SetOutChan(c *chan *DkOutMsg) {
	this.COutMsg = c
}

func (this *DeskMgr) SendNativeMsg(id int, uid int64, d interface{}) {
	msg := DkOutMsg{
		Id:  int32(id),
		Uid: uid,
	}
	if d != nil {
		msg.Data, _ = json.Marshal(d)
	}
	if this.COutMsg != nil {
		*this.COutMsg <- &msg
	} else {
		logs.Error("desk mgr send miss", msg)
	}
}

func (this *DeskMgr) AddMsg(id int32, uid int64, d string) {
	inmsg := DkInMsg{
		Id:  id,
		Uid: uid,
	}
	inmsg.Data = d
	*this.CInMsg <- &inmsg
}

func (this *DeskMgr) AddNativeMsg(id int32, uid int64, d interface{}) {
	v, _ := json.Marshal(d)
	this.AddMsg(id, uid, string(v))
}

// 归还桌子指针，调用此函数前必须充值desk所有数据
func (this *DeskMgr) BackDesk(p *ExtDesk) {
	this.Free = append(this.Free, p)
	//
	if p.FkNo != "" {
		delete(this.MFkDesk, p.FkNo)
	}
}

func (this *DeskMgr) Run() {
	t := time.NewTicker(time.Second)
	for {
		select {
		case v := <-*this.CInMsg:
			h, ok := this.Handle[v.Id]
			if ok {
				h(v)
			} else {
				p, ok := this.MapPlayers[v.Uid]
				if ok {
					dh, ok := p.Dk.Handle[v.Id]
					if ok {
						dh(p, v)
					}
				}
			}
		case <-t.C:
			for _, d := range this.All {
				d.DoTimer()
			}
		}
	}
}

//////////////////////////////////////////////////////////
//以下为消息处理函数
func (this *DeskMgr) InitHandle() {
	// 玩家退出游戏
	this.Handle[MSG_GAME_LEAVE] = this.Leave
	// 玩家进入匹配模式
	this.Handle[MSG_GAME_AUTO] = this.GameAuto
	//断线重连
	this.Handle[MSG_GAME_RECONNECT] = this.ReConnect
	//
	this.Handle[MSG_GAME_UPDATE_PLAYER_INFO] = this.UpdatePlayerInfo
	//
	this.Handle[MSG_GAME_FK_CREATEDESK] = this.CreateFkDesk
	//
	this.Handle[MSG_GAME_FK_JOIN] = this.FkJoin
}

func (this *DeskMgr) GameAuto(msg *DkInMsg) {
	//
	jmsg := GAutoGame{}
	err := json.Unmarshal([]byte(msg.Data), &jmsg)
	if err != nil {
		logs.Error("DeskMgr GameAuto err:", err)
		return
	}
	//用户是否存在
	p, ok := this.MapPlayers[jmsg.Uid]
	if ok {
		logs.Error("DeskMgr player exist:", jmsg.Uid)
		p.SendNativeMsg(MSG_GAME_AUTO_REPLY, GAutoGameReply{
			Id:     MSG_GAME_AUTO_REPLY,
			Err:    "玩家已在游戏中",
			Result: 11,
		})
		return
	}
	//
	if jmsg.Coin < int64(G_DbGetGameServerData.Restrict) {
		this.SendNativeMsg(MSG_GAME_AUTO_REPLY, msg.Uid, &GAutoGameReply{
			Id:     MSG_GAME_AUTO_REPLY,
			Err:    "匹配失败，金币不足",
			Result: 9,
		})
		return
	}
	//
	p = &ExtPlayer{}
	//游戏注册用户
	p.Account = jmsg.Account
	p.Uid = jmsg.Uid
	p.Nick = jmsg.Nick
	p.Sex = jmsg.Sex
	p.Head = jmsg.Head
	p.Lv = jmsg.Lv
	p.SetMgr(this)
	p.Coins = jmsg.Coin
	p.Token = jmsg.Token
	p.Robot = jmsg.Robot
	p.HierarchyId = jmsg.HierarchyId
	this.MapPlayers[p.Uid] = p

	if GCONFIG.RoomType == 3 { //百人场
		dk := this.Free[0]
		pnum := dk.AddPlayer(p)
		if pnum > 0 {
			p.SetDesk(dk)
			dk.DoHandle(p, msg)
			return
		}
	} else if this.DeskAllotModel == 1 {
		for _, v := range this.Free {
			pnum := v.AddPlayer(p)
			if pnum > 0 {
				p.SetDesk(v)
				v.DoHandle(p, msg)
				return
			}
		}
	} else {
		newfree := []*ExtDesk{}
		for i, v := range this.Free {
			if v.GameState != GAME_STATUS_FREE {
				continue
			}
			pnum := v.AddPlayer(p)
			if pnum > 0 {
				p.SetDesk(v)
				v.DoHandle(p, msg)
				if pnum < GCONFIG.PlayerNum {
					newfree = append(newfree, v)
				}
				this.Free = append(newfree, this.Free[i+1:]...)
				return
			}
			newfree = append(newfree, v)
		}
		this.Free = newfree
	}

	if p.Dk == nil {
		delete(this.MapPlayers, p.Uid)
	}
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, GAutoGameReply{
		Id:     MSG_GAME_AUTO_REPLY,
		Err:    "服务器繁忙，请稍后重试",
		Result: 1, // 没有多余的桌子
	})
}

func (this *DeskMgr) SetDeskAllotModel(mo int) {
	this.DeskAllotModel = mo
}

func (this *DeskMgr) FkJoin(msg *DkInMsg) {
	req := GFkJoinToGame{}
	err := json.Unmarshal([]byte(msg.Data), &req)
	if err != nil {
		logs.Error("加入房卡场参数传错", err, msg)
		this.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, msg.Uid, &GFkJoinReply{
			Id:     MSG_GAME_FK_JOIN_REPLY,
			Err:    "参数传错",
			Result: 1,
		})
		return
	}
	//用户是否存在
	p, ok := this.MapPlayers[req.Uid]
	if ok {
		p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, GFkJoinReply{
			Id:     MSG_GAME_FK_JOIN_REPLY,
			Err:    "玩家已存在",
			Result: 2,
		})
		return
	}
	//查找房间号是否存在
	dk, ok := this.MFkDesk[req.FkNo]
	if !ok {
		this.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, req.Uid, &GFkJoinReply{
			Id:     MSG_GAME_FK_JOIN_REPLY,
			Err:    "房间号不存在",
			Result: 3,
		})
		return
	}
	//注册用户
	p = &ExtPlayer{}
	p.Account = req.Account
	p.Uid = req.Uid
	p.Nick = req.Nick
	p.Sex = req.Sex
	p.Head = req.Head
	p.Lv = req.Lv
	p.SetMgr(this)
	p.Coins = req.Coin
	p.Token = req.Token
	p.Robot = req.Robot
	p.RoomCard = req.RoomCard
	p.SetDesk(dk)
	this.MapPlayers[p.Uid] = p
	//调用桌子的这个协议
	dk.DoHandle(p, msg)
}

func (this *DeskMgr) ReConnect(msg *DkInMsg) {
	// re := GReConnect{}
	// json.Unmarshal([]byte(msg.Data), &re)
	//
	p, ok := this.MapPlayers[msg.Uid]
	if !ok {
		this.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, msg.Uid, &GReConnectReply{
			Id:     MSG_GAME_RECONNECT_REPLY,
			Result: 2,
			Err:    "游戏中没有该玩家",
		})
		return
	}
	//
	if p.Dk == nil {
		this.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, msg.Uid, &GReConnectReply{
			Id:     MSG_GAME_RECONNECT_REPLY,
			Result: 3,
			Err:    "该玩家没有在游戏中，可直接离开",
		})
		return
	}
	//
	p.LiXian = false
	p.Dk.HandleReconnect(p, msg)
}

func (this *DeskMgr) Leave(msg *DkInMsg) {
	p, ok := this.MapPlayers[msg.Uid]
	if !ok {
		logs.Debug("用户不在匹配和游戏中")
		return
	}
	if p.Dk != nil {
		p.Dk.Leave(p)
	} else {
		this.LeaveDo(p.Uid)
		// this.SendNativeMsg(MSG_GAME_LEAVE_REPLY, p.Uid, &GLeaveReply{
		// 	Id:     MSG_GAME_LEAVE_REPLY,
		// 	Cid:    -1,
		// 	Result: 0,
		// })
	}

}

func (this *DeskMgr) UpdatePlayerInfo(msg *DkInMsg) {
	p, ok := this.MapPlayers[msg.Uid]
	if ok {
		if p.Dk != nil {
			p.Dk.UpdatePlayerInfo(p, msg)
		}
	}
}

func (this *DeskMgr) CreateFkDesk(msg *DkInMsg) {
	//没有空余的桌子
	if len(this.Free) == 0 {
		this.SendNativeMsg(MSG_HALL_CREATE_FKROOM_REPLY, msg.Uid, &GFkCreateDeskReply{
			Id:     MSG_HALL_CREATE_FKROOM_REPLY,
			Result: 1,
			Err:    "没有空余的桌子，请稍后重试",
		})
		return
	}
	//
	req := GFkCreateDesk{}
	err := json.Unmarshal([]byte(msg.Data), &req)
	if err != nil {
		this.SendNativeMsg(MSG_HALL_CREATE_FKROOM_REPLY, msg.Uid, &GFkCreateDeskReply{
			Id:     MSG_HALL_CREATE_FKROOM_REPLY,
			Result: 2,
			Err:    "创建房卡桌子，参数传错",
		})
		return
	}
	//判断房间号是否存在
	_, ok := this.MFkDesk[req.FkNo]
	if ok {
		this.SendNativeMsg(MSG_HALL_CREATE_FKROOM_REPLY, msg.Uid, &GFkCreateDeskReply{
			Id:     MSG_HALL_CREATE_FKROOM_REPLY,
			Result: 3,
			Err:    "房间号已存在",
		})
		return
	}
	//赋值房卡信息
	nk := this.Free[0]
	nk.FkNo = req.FkNo
	nk.FkOwner = msg.Uid
	nk.FkInfo = req.FkInfo
	this.MFkDesk[req.FkNo] = nk
	this.Free = this.Free[1:]
	//开启桌子回收定时器，10分钟
	nk.AddBackDeskTimer()
	//
	this.SendNativeMsg(MSG_HALL_CREATE_FKROOM_REPLY, msg.Uid, &GFkCreateDeskReply{
		Id:        MSG_HALL_CREATE_FKROOM_REPLY,
		Result:    0,
		FkNo:      req.FkNo,
		GameType:  int32(GCONFIG.GameType),
		RoomType:  int32(GCONFIG.RoomType),
		GradeType: int32(GCONFIG.GradeType),
	})
}

func (this *DeskMgr) LeaveDo(uid int64) {
	// logs.Debug("删除管理器里面的玩家")
	delete(this.MapPlayers, uid)
}

func GetGameInfoFromDb() {
	if DB_SAVE_MODEL == 1 {
		getip := GCONFIG.WebGameIp + "/web/gameInfo/gameinfo?"
		getip = getip + "gradeId=" + strconv.Itoa(GCONFIG.GradeType)
		getip = getip + "&gameId=" + strconv.Itoa(GCONFIG.GameType)
		getip = getip + "&roomId=" + strconv.Itoa(GCONFIG.RoomType)
		token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJoYWxscmVxdWVzdHNtYWxsIiwiaWF0IjoxNTUwODE1Mzc4fQ.mNr0cT3i-F-E5twxX3RSbtFd7tgfhFRr2kbZY8o3RLQ"
		re, err := SendRequest(getip, nil, "GET", token)
		if err != nil {
			logs.Debug("获取游戏服务器信息失败", err)
			os.Exit(5)
		}
		rsp := DbGetGameServerRsp{}
		json.Unmarshal([]byte(re), &rsp)
		if rsp.Code != 200 || len(rsp.Data.Game) == 0 {
			logs.Debug("获取游戏服务器信息失败", rsp)
			os.Exit(5)
		}
		//
		G_DbGetGameServerData = rsp.Data.Game[0]
		logs.Debug("获取游戏服务器信息成功", G_DbGetGameServerData)
	} else {
		G_DbGetGameServerData = DbGetGameServerData{
			Bscore:   100,
			MaxTimes: 8,
		}
	}

}

func GetRoleRateList() {
	if DB_SAVE_MODEL == 1 {
		getip := GCONFIG.WebDbIp + "/V1.0/App/hierarchy/rateList"
		token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJoYWxscmVxdWVzdGhvbWUiLCJpYXQiOjE1NTA4MTYxOTR9.m8zSOaTzkXKTKoUuEISF-fKryM2KDKPO-_YwEXizN54"
		re, err := SendRequest(getip, nil, "GET", token)
		if err != nil {
			logs.Debug("获取玩家胜率列表失败", err)
			os.Exit(5)
		}
		rsp := DbGetRoleRateListRsp{}
		json.Unmarshal([]byte(re), &rsp)
		if rsp.Code != 200 {
			logs.Debug("获取游戏服务器信息失败", rsp)
			os.Exit(5)
		}
		G_DbRoleRateListData = rsp.Data
		logs.Debug("获取玩家层级胜率成功", G_DbRoleRateListData)
	}
}

func GetJuHao() string {
	nowt := time.Now()
	return nowt.Format("20060102150405")
}

/////////////////////////////////////////////////////
//系统游戏配置接口
//获取当前库存值
func GetLocalStock() int64 {
	return G_DbGetGameServerData.GameConfig.CurrentStock
}

//修改当前库存值
func AddLocalStock(a int64) {
	G_DbGetGameServerData.GameConfig.CurrentStock += a
	GDeskMgr.ChangeStock += a
}

//根据玩家层级获取概率，不存在返回-1
func GetRateByHierarchyId(id int32) float32 {
	for _, v := range G_DbRoleRateListData {
		if v.HierarchyId == int(id) {
			return v.Rate
		}
	}
	return -1
}

//根据当前库存获取概率，不存在返回-1
func GetRateByInterval() float32 {
	for _, v := range G_DbGetGameServerData.GameConfig.IntervalGameRate {
		if G_DbGetGameServerData.GameConfig.CurrentStock >= v.IntervalStart &&
			G_DbGetGameServerData.GameConfig.CurrentStock <= v.IntervalEnd {
			return v.Rate
		}
	}
	return -1
}

//消耗类型，1金币，2代币（不能操作金币）
func GetCostType() int {
	if GCONFIG.GradeType == 6 {
		return 2
	} else {
		return 1
	}
}

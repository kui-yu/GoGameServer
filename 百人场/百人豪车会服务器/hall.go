package main

import (
	"encoding/json"
	"fmt"
	"logs"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	DB_SAVE_MODEL = 1 //1数据库模式，2内存模式
)

var GHall Hall

type Hall struct {
	CInMsg    chan *InMsg
	CDkInMsg  chan *DkInMsg  //hall->deskmgr
	CDkOutMsg chan *DkOutMsg //deskmgr->hall
	CDbInMsg  chan *DbInMsg
	CDbOutMsg chan *DbOutMsg
	Handle    map[int32]func(*InMsg)
	Players   map[int64]*Contoller
	Lk        sync.RWMutex
	MGameInfo map[GameTypeDetail]*DbGetGameServerData
	//数据库
	DbData  map[string]*DbPlayerData
	DbData2 map[int64]*DbPlayerData
	DbLk    sync.RWMutex
	GUid    int64
	//
	MRoomNo map[string]*GameTypeDetail
}

// type DbPlayer struct {
// 	Uid     int64 //玩家id
// 	Account string
// 	Nick    string
// 	Coin    int32
// 	Sex     int32
// 	GameId  int32
// }

func init() {
	GHall.CInMsg = make(chan *InMsg, 1000)
	GHall.CDkInMsg = make(chan *DkInMsg, 1000)
	GHall.CDkOutMsg = make(chan *DkOutMsg, 1000)
	GHall.CDbInMsg = make(chan *DbInMsg, 1000)
	GHall.CDbOutMsg = make(chan *DbOutMsg, 1000)
	GHall.Handle = make(map[int32]func(*InMsg))
	GHall.Players = make(map[int64]*Contoller)
	GHall.DbData = make(map[string]*DbPlayerData)
	GHall.DbData2 = make(map[int64]*DbPlayerData)
	GHall.MGameInfo = make(map[GameTypeDetail]*DbGetGameServerData)
	GHall.MRoomNo = make(map[string]*GameTypeDetail)
	GHall.InitHandle()
	//
	GDeskMgr.SetInChan(&GHall.CDkInMsg)
	GDeskMgr.SetOutChan(&GHall.CDkOutMsg)
	GDeskMgr.InitData()
	//
	// G_DbProxy.SetInChan(&GHall.CDbInMsg)
	// G_DbProxy.SetOutChan(&GHall.CDbOutMsg)
	// G_DbProxy.InitData()
	//
	go GHall.Run()
	go GHall.Run2()
	// go GHall.Run3()
	logs.Debug("大厅初始化结束")
	go GDeskMgr.Run()
	logs.Debug("游服初始化结束")
	// go G_DbProxy.Run()
	// logs.Debug("数据库代理初始化结束")
}

func (this *Hall) AddMsg(msg *InMsg) {
	this.CInMsg <- msg
}

func (this *Hall) Run() {
	for v := range this.CInMsg {
		if v.Id >= 400000 && v.Id < 490000 {
			if v.Id == MSG_GAME_AUTO {
				re := GAutoGameToHall{
					IsRobot: false,
				}
				json.Unmarshal([]byte(v.Data), &re)
				//
				if v.Col.GameId != 0 {
					v.Col.AddMsgNative(MSG_GAME_AUTO_REPLY, GAutoGameReply{
						Id:     MSG_GAME_AUTO_REPLY,
						Result: 2,
						Err:    "玩家已加入游戏",
					}, false)
					continue
				}
				nid := MarshalGameAllType(re.RoomType, re.GradeType, re.GameType)
				if !v.Col.Robot {
					this.SetGameIdByUid(v.Col.Uid, nid, 1)
				}

				v.Col.GameId = nid
				v.Col.ServerId = 1
				takeCoin := v.Col.Coin
				if GetCostType() == 2 {
					takeCoin = 200000
				}
				GDeskMgr.AddNativeMsg(v.Id, v.Uid, &GAutoGame{
					Id:          v.Id,
					Account:     v.Col.Account,
					Uid:         v.Col.Uid,
					Nick:        v.Col.Nick,
					Sex:         v.Col.Sex,
					Head:        v.Col.Head,
					Lv:          1,
					Coin:        takeCoin,
					HierarchyId: int32(v.Col.HierarchyId),
					// Token:   v.Col.Token,
					Robot: v.Col.Robot,
				})
			} else if v.Id == MSG_GAME_FK_JOIN {
				re := GFkJoinToHall{}
				json.Unmarshal([]byte(v.Data), &re)
				//
				if v.Col.GameId != 0 {
					v.Col.AddMsgNative(MSG_GAME_FK_JOIN_REPLY, GFkJoinReply{
						Id:     MSG_GAME_FK_JOIN_REPLY,
						Result: 1,
						Err:    "房卡玩家已加入游戏",
					}, false)
					continue
				}
				//
				nid := MarshalGameAllType(re.RoomType, re.GradeType, re.GameType)
				v.Col.GameId = nid
				//
				GDeskMgr.AddNativeMsg(v.Id, v.Uid, &GFkJoinToGame{
					Id:      v.Id,
					Account: v.Col.Account,
					Uid:     v.Col.Uid,
					Nick:    v.Col.Nick,
					Sex:     v.Col.Sex,
					Head:    v.Col.Head,
					Lv:      1,
					Coin:    v.Col.Coin,
					// Token:   v.Col.Token,
					Robot:    v.Col.Robot,
					RoomCard: v.Col.RoomCard,
					FkNo:     re.FkNo,
				})
			} else {
				GDeskMgr.AddMsg(v.Id, v.Uid, v.Data)
			}
		} else {
			h, ok := this.Handle[v.Id]
			if ok {
				h(v)
			}
		}
	}
}

func (this *Hall) Run2() {
	for m := range this.CDkOutMsg {
		this.Lk.RLock()
		p, ok := this.Players[m.Uid]
		this.Lk.RUnlock()
		if ok {
			if m.Id == MSG_GAME_AUTO_REPLY {
				re := GAutoGameReply{}
				json.Unmarshal(m.Data, &re)
				if re.Result != 0 {
					p.GameId = 0
					p.ServerId = 0
					if !p.Robot {
						this.SetGameIdByUid(p.Uid, 0, 0)
					}
				}
			} else if m.Id == MSG_GAME_RECONNECT_REPLY {
				re := GReConnectReply{}
				json.Unmarshal(m.Data, &re)
				if re.Result != 0 {
					p.GameId = 0
					p.ServerId = 0
					if !p.Robot {
						this.SetGameIdByUid(p.Uid, 0, 0)
					}
				}
			} else if m.Id == MSG_GAME_FK_JOIN_REPLY {
				re := GFkJoinReply{}
				json.Unmarshal(m.Data, &re)
				if re.Result != 0 {
					p.GameId = 0
					p.ServerId = 0
					this.SetGameIdByUid(p.Uid, 0, 0)
				} else {
					this.SetGameIdByUid(p.Uid, p.GameId, p.ServerId)
				}
			}
			if m.Id == MSG_GAME_LEAVE_REPLY {
				re := GLeaveReply{}
				json.Unmarshal(m.Data, &re)
				if !re.NoToCli {
					p.AddMsg(m.Id, m.Data, false)
				}
			} else {
				p.AddMsg(m.Id, m.Data, false)
			}

		}
		if m.Id == MSG_GAME_END_NOTIFY {
			re := GGameEnd{}
			json.Unmarshal(m.Data, &re)
			robot := true
			for _, v := range re.UserCoin {
				if !v.Robot {
					this.SetCoin(v.UserAccount, v.UserId, v.PrizeCoins)
				}
				robot = v.Robot
				if re.SetLeave == 1 {
					continue
				}
				this.Lk.RLock()
				pr, ok := this.Players[v.UserId]
				this.Lk.RUnlock()
				if ok {
					pr.Coin += int64(v.PrizeCoins)
					pr.GameId = 0
					pr.ServerId = 0
				}
				// logs.Debug("............", v.UserId, v.Token)
				if !v.Robot {
					this.SetGameIdByUid(v.UserId, 0, 0)
				}
			}
			//
			if robot {
				this.PushRobotRecord(&re)
			} else {
				this.PushRecord(&re)
			}
		} else if m.Id == MSG_GAME_END_RECORD {
			this.PushOtherRecord(string(m.Data))
		} else if m.Id == MSG_GAME_LEAVE_REPLY {
			re := GLeaveReply{}
			json.Unmarshal(m.Data, &re)
			if re.Result == 0 {
				if re.Uid == m.Uid {
					if ok {
						p.GameId = 0
						p.ServerId = 0
					}
					//等机器人uid和玩家uid分离后再添加处理
					if !re.Robot {
						this.SetGameIdByUid(m.Uid, 0, 0)
					}
				}
			}
		} else if m.Id == MSG_GAME_ADDCOIN_NOTIFY {
			re := GToHAddCoin{}
			json.Unmarshal(m.Data, &re)
			//
			this.SetCoin("", re.Uid, re.Coin)
			//
			if ok {
				p.Coin += re.Coin
				p.AddMsg(m.Id, m.Data, false)
			}
		} else if m.Id == MSG_GAME_ADDROOMCARD_NOTIFY {
			re := GToHAddRoomCard{}
			json.Unmarshal(m.Data, &re)
			//
			this.SetRoomCard("", re.Uid, re.RoomCard)
			//
			if ok {
				p.RoomCard += re.RoomCard
				p.AddMsg(m.Id, m.Data, false)
			}
		} else if m.Id == MSG_GAME_GETZOUSHI_REPLY {
			re := GGetZouShiReply{}
			json.Unmarshal([]byte(m.Data), &re)
			//
			GMgrGameZouShi.AddZouShi(&re.Data)
		}
	}
}

///////////////////////////////////////
//消息注册函数
func (this *Hall) InitHandle() {
	this.Handle[MSG_HALL_LOGIN] = this.Login
	this.Handle[MSG_HALL_LEAVE] = this.Leave
	this.Handle[MSG_HALL_JOIN_GMAE] = this.JoinGame
	this.Handle[MSG_HALL_GETNEWESTCOIN] = this.GetNewestCoin
	this.Handle[MSG_HALL_ROBOT_LOGIN] = this.RobotLogin
	this.Handle[MSG_HALL_CREATE_FKROOM] = this.CreateFkRoom
	this.Handle[MSG_HALL_JOIN_FKROOM] = this.JoinFkRoom
	this.Handle[MSG_HALL_HEART] = this.Heart
	this.Handle[MSG_HALL_GETZOUSHI] = this.GetZouShi
	this.Handle[MSG_HALL_GETZOUSHI_SINGLE] = this.GetZouShiSingle
}

func (this *Hall) Login(msg *InMsg) {
	jmsg := HMsgHallLogin{}
	err := json.Unmarshal([]byte(msg.Data), &jmsg)
	if err != nil {
		//断掉长连接
		logs.Error("login error:", err, msg)
		msg.Col.AddMsgNative(MSG_HALL_LOGIN_REPLY, &HMsgHallLoginReply{
			Id:     MSG_HALL_LOGIN_REPLY,
			Result: 1,
			Err:    "参数传错",
		}, false)
		msg.Col.AddMsgNative(MSG_HALL_LEAVE_REPLY, nil, true)
		return
	}
	//根据账号获取用户信息,不存在则存入
	player, err := this.GetPlayer(jmsg.Account, jmsg.Gid)
	if err != nil {
		if DB_SAVE_MODEL == 1 {
			msg.Col.AddMsgNative(MSG_HALL_LOGIN_REPLY, &HMsgHallLoginReply{
				Id:     MSG_HALL_LOGIN_REPLY,
				Result: 2,
				Err:    "玩家不存在或者数据库错误",
			}, false)
			msg.Col.AddMsgNative(MSG_HALL_LEAVE_REPLY, nil, true)
			// msg.Col.AddMsgNative(MSG_HALL_LEAVE_REPLY, &GLeaveReply{
			// 	Id:     MSG_HALL_LEAVE_REPLY,
			// 	Result: 10,
			// 	Err:    "玩家不存在或者数据库错误",
			// }, true)
			return
		} else {
			this.GUid++
			nuid := this.GUid
			player = &DbPlayerData{
				Uid:          nuid,
				Account:      jmsg.Account,
				Coin:         10000,
				GameInfoId:   0,
				GameServerId: 0,
			}
			this.SavePlayer(jmsg.Account, player)
		}
	}
	//
	this.Lk.RLock()
	coler, ok := this.Players[player.Uid]
	this.Lk.RUnlock()
	if ok {
		if coler.GameId != 0 {
			GDeskMgr.AddMsg(MSG_GAME_DISCONNECT, player.Uid, "")
		}
		this.Lk.Lock()
		delete(this.Players, player.Uid)
		this.Lk.Unlock()
		coler.AddMsgNative(MSG_HALL_OTHERLOGIN_NOTIFY, &HMsgHallQiangDengNotify{
			Id: MSG_HALL_OTHERLOGIN_NOTIFY,
		}, false)
		coler.AddMsgNative(MSG_HALL_LEAVE_REPLY, nil, true)
	}

	//
	msg.Col.Token = jmsg.Gid
	msg.Col.Account = player.Account
	msg.Col.Uid = player.Uid
	msg.Col.GameId = int32(player.GameInfoId)
	msg.Col.ServerId = int32(player.GameServerId)
	msg.Col.Coin = player.Coin
	msg.Col.Nick = player.Account
	msg.Col.Head = player.Portrait
	msg.Col.Robot = false
	msg.Col.HierarchyId = player.HierarchyId
	msg.Col.RoomCard = player.RoomCard
	if player.Sex {
		msg.Col.Sex = 1
	} else {
		msg.Col.Sex = 0
	}
	this.Lk.Lock()
	this.Players[player.Uid] = msg.Col
	this.Lk.Unlock()
	//判断gameid是否为0，不是则断线重连
	msg.Col.AddMsgNative(MSG_HALL_LOGIN_REPLY, &HMsgHallLoginReply{
		Id:              MSG_HALL_LOGIN_REPLY,
		Result:          0,
		Account:         player.Account,
		Uid:             player.Uid,
		Nick:            player.Account,
		Sex:             msg.Col.Sex,
		Head:            player.Portrait,
		Coin:            player.Coin,
		GameSerId:       int32(player.GameInfoId),
		AliPayAccount:   player.AliPayAccount,
		BankCard:        player.BankCard,
		Commission:      player.Commission,
		Money:           player.Money,
		UnReadNum:       player.UnreadNum,
		FrozenEnable:    player.FrozenEnable,
		ForbiddenEnable: player.ForbiddenEnable,
		UserName:        player.UserName,
		BindPassword:    player.BindPassword,
		RoomCard:        player.RoomCard,
		NickName:        player.NickName,
		Phone:           player.Phone,
		InvitationCode:  player.InvitationCode,
	}, false)
}

func (this *Hall) RobotLogin(msg *InMsg) {
	jmsg := HMsgHallRobotLogin{}
	err := json.Unmarshal([]byte(msg.Data), &jmsg)
	if err != nil {
		//断掉长连接
		logs.Error("robot login error:", err, msg)
		msg.Col.AddMsgNative(MSG_HALL_ROBOT_LOGIN_REPLY, &HMsgHallRobotLoginReply{
			Id:     MSG_HALL_ROBOT_LOGIN_REPLY,
			Result: 1,
			Err:    "参数传错",
		}, false)
		msg.Col.AddMsgNative(MSG_HALL_LEAVE_REPLY, nil, true)
		return
	}
	robot, err := this.GetRobot(jmsg.Gid)
	if err != nil {
		msg.Col.AddMsgNative(MSG_HALL_ROBOT_LOGIN_REPLY, &HMsgHallRobotLoginReply{
			Id:     MSG_HALL_ROBOT_LOGIN_REPLY,
			Result: 2,
			Err:    "机器人不存在或者数据库错误",
		}, false)
		msg.Col.AddMsgNative(MSG_HALL_LEAVE_REPLY, nil, true)
		return
	}
	//
	this.Lk.RLock()
	_, ok := this.Players[robot.Uid]
	this.Lk.RUnlock()
	if ok {
		msg.Col.AddMsgNative(MSG_HALL_ROBOT_LOGIN_REPLY, &HMsgHallRobotLoginReply{
			Id:     MSG_HALL_ROBOT_LOGIN_REPLY,
			Result: 3,
			Err:    "机器人已经登录过",
		}, false)
		msg.Col.AddMsgNative(MSG_HALL_LEAVE_REPLY, nil, true)
		return
	}
	//
	msg.Col.Token = jmsg.Gid
	msg.Col.Account = robot.Name
	msg.Col.Uid = robot.Uid
	msg.Col.Coin = robot.Coin
	msg.Col.Nick = robot.Name
	msg.Col.Head = robot.Head
	msg.Col.Robot = true
	if robot.Sex {
		msg.Col.Sex = 1
	} else {
		msg.Col.Sex = 0
	}
	this.Lk.Lock()
	this.Players[robot.Uid] = msg.Col
	this.Lk.Unlock()
	//
	msg.Col.AddMsgNative(MSG_HALL_ROBOT_LOGIN_REPLY, &HMsgHallRobotLoginReply{
		Id:      MSG_HALL_ROBOT_LOGIN_REPLY,
		Result:  0,
		Uid:     robot.Uid,
		Coin:    robot.Coin,
		Account: robot.Name,
		Head:    robot.Head,
		Sex:     robot.Sex,
	}, false)
	//
}

func (this *Hall) JoinGame(msg *InMsg) {
	req := HMsgHallJoinGame{}
	err := json.Unmarshal([]byte(msg.Data), &req)
	if err != nil {
		logs.Error("JoinGame error:", err, msg)
		msg.Col.AddMsgNative(MSG_HALL_JOIN_GAME_REPLY, &HMsgHallJoinGameReply{
			Id:     MSG_HALL_JOIN_GAME_REPLY,
			Result: 1,
			Err:    "参数传错",
		}, false)
		return
	}
	//判断游服是否支持
	//判断金币是否足够
	if DB_SAVE_MODEL == 1 {
		gifo, ok := this.MGameInfo[GameTypeDetail{
			GameType:  req.GameType,
			RoomType:  req.RoomType,
			GradeType: req.GradeType,
		}]
		if !ok {
			gifo, err = GetGameDetailByWeb(req.RoomType, req.GradeType, req.GameType)
			if err != nil {
				msg.Col.AddMsgNative(MSG_HALL_JOIN_GAME_REPLY, &HMsgHallJoinGameReply{
					Id:     MSG_HALL_JOIN_GAME_REPLY,
					Result: 2,
					Err:    err.Error(),
				}, false)
				return
			}
		}
		this.MGameInfo[GameTypeDetail{
			GameType:  req.GameType,
			RoomType:  req.RoomType,
			GradeType: req.GradeType,
		}] = gifo
		//
		p, ok := this.Players[msg.Uid]
		if !ok {
			msg.Col.AddMsgNative(MSG_HALL_JOIN_GAME_REPLY, &HMsgHallJoinGameReply{
				Id:     MSG_HALL_JOIN_GAME_REPLY,
				Result: 3,
				Err:    "用户不存在",
			}, false)
			return
		}
		if gifo.Restrict > int(p.Coin) {
			msg.Col.AddMsgNative(MSG_HALL_JOIN_GAME_REPLY, &HMsgHallJoinGameReply{
				Id:     MSG_HALL_JOIN_GAME_REPLY,
				Result: 4,
				Err:    "金币不足,最低：" + strconv.Itoa(int(gifo.Restrict)),
			}, false)
			return
		}
	}
	//
	msg.Col.AddMsgNative(MSG_HALL_JOIN_GAME_REPLY, &HMsgHallJoinGameReply{
		Id:     MSG_HALL_JOIN_GAME_REPLY,
		Result: 0,
	}, false)
}

func (this *Hall) Heart(msg *InMsg) {
	msg.Col.AddMsgNative(MSG_HALL_HEART_REPLY, &HMsgHallHeartReply{
		Id: MSG_HALL_HEART_REPLY,
	}, false)
}

func (this *Hall) GetZouShi(msg *InMsg) {
	jmsg := HMsgHallGetZouShi{}
	err := json.Unmarshal([]byte(msg.Data), &jmsg)
	if err != nil {
		msg.Col.AddMsgNative(MSG_HALL_GETZOUSHI_REPLY, &HMsgHallGetZouShiReply{
			Id:     MSG_HALL_GETZOUSHI_REPLY,
			Result: 1,
			Err:    "参数传错",
		}, false)
		return
	}
	//
	rsp := HMsgHallGetZouShiReply{
		Id:     MSG_HALL_GETZOUSHI_REPLY,
		Result: 0,
	}
	rsp.ZouShi = GMgrGameZouShi.GetZouShi(&GameTypeDetail{
		GameType:  jmsg.GameType,
		RoomType:  jmsg.RoomType,
		GradeType: jmsg.GradeType,
	})
	//
	msg.Col.AddMsgNative(MSG_HALL_GETZOUSHI_REPLY, &rsp, false)
}

func (this *Hall) GetZouShiSingle(msg *InMsg) {
	jmsg := HMsgHallGetZouShiSingle{}
	err := json.Unmarshal([]byte(msg.Data), &jmsg)
	if err != nil {
		msg.Col.AddMsgNative(MSG_HALL_GETZOUSHI_SINGLE_REPLY, &HMsgHallGetZouShiSingleReply{
			Id:     MSG_HALL_GETZOUSHI_SINGLE_REPLY,
			Result: 1,
			Err:    "参数传错",
		}, false)
		return
	}

	zoushi := GMgrGameZouShi.GetZouShiSingle(&GameTypeDetail{
		GameType:  jmsg.GameType,
		RoomType:  jmsg.RoomType,
		GradeType: jmsg.GradeType,
	}, jmsg.SerId)

	if zoushi == nil {
		msg.Col.AddMsgNative(MSG_HALL_GETZOUSHI_SINGLE_REPLY, &HMsgHallGetZouShiSingleReply{
			Id:     MSG_HALL_GETZOUSHI_SINGLE_REPLY,
			Result: 2,
			Err:    "游戏服不存在",
		}, false)
		return
	}
	msg.Col.AddMsgNative(MSG_HALL_GETZOUSHI_SINGLE_REPLY, &HMsgHallGetZouShiSingleReply{
		Id:     MSG_HALL_GETZOUSHI_SINGLE_REPLY,
		Result: 0,
		ZouShi: *zoushi,
	}, false)

}

func (this *Hall) GetNewestCoin(msg *InMsg) {
	coin := int64(0)
	this.Lk.RLock()
	p, ok := this.Players[msg.Uid]
	if ok {
		coin = p.Coin
	}
	this.Lk.RUnlock()
	if !ok {
		msg.Col.AddMsgNative(MSG_HALL_GETNEWESTCOIN_REPLY, HMsgHallGetNewestCoinReply{
			Id:     MSG_HALL_GETNEWESTCOIN_REPLY,
			Result: 1,
			Err:    "更新金币的时候，用户不在线",
		}, false)
		return
	}
	msg.Col.AddMsgNative(MSG_HALL_GETNEWESTCOIN_REPLY, HMsgHallGetNewestCoinReply{
		Id:     MSG_HALL_GETNEWESTCOIN_REPLY,
		Result: 0,
		Coin:   coin,
	}, false)
}

//创建房卡房间
func (this *Hall) CreateFkRoom(msg *InMsg) {
	req := HMsgHallCreateFkRoom{}
	err := json.Unmarshal([]byte(msg.Data), &req)
	if err != nil {
		logs.Error("创建房间失败，参数解析错误", err, msg)
		msg.Col.AddMsgNative(MSG_HALL_CREATE_FKROOM_REPLY, &HMsgHallCreateFkRoomReply{
			Id:     MSG_HALL_CREATE_FKROOM_REPLY,
			Result: 1,
			Err:    "参数传错",
		}, false)
		return
	}
	//获取唯一房间号
	fkno := this.CreateRoomNo(&GameTypeDetail{
		GameType:  int32(req.GameType),
		RoomType:  int32(req.RoomType),
		GradeType: int32(req.GradeType),
	})
	//
	msg.Col.ServerId = 1
	//去游戏服注册房间信息
	GDeskMgr.AddNativeMsg(MSG_GAME_FK_CREATEDESK, msg.Uid, &GFkCreateDesk{
		Id:     MSG_GAME_FK_CREATEDESK,
		FkNo:   fkno,
		FkInfo: req.FkInfo,
	})
}

func (this *Hall) JoinFkRoom(msg *InMsg) {
	req := HMsgHallJoinFkRoom{}
	err := json.Unmarshal([]byte(msg.Data), &req)
	if err != nil {
		logs.Error("加入房卡房间参数错误", err, msg)
		msg.Col.AddMsgNative(MSG_HALL_JOIN_FKROOM_REPLY, &HMsgHallJoinFkRoomReply{
			Id:     MSG_HALL_JOIN_FKROOM_REPLY,
			Result: 1,
			Err:    "参数传错",
		}, false)
		return
	}
	//查询房间是否存在
	msg.Col.ServerId = 1
	//
	gt := this.GetFkRoom(req.FkNo)
	if gt == nil {
		logs.Error("无此房卡号", req.FkNo)
		msg.Col.AddMsgNative(MSG_HALL_JOIN_FKROOM_REPLY, &HMsgHallJoinFkRoomReply{
			Id:     MSG_HALL_JOIN_FKROOM_REPLY,
			Result: 2,
			Err:    "无此房卡号",
		}, false)
		return
	}
	//
	msg.Col.AddMsgNative(MSG_HALL_JOIN_FKROOM_REPLY, &HMsgHallJoinFkRoomReply{
		Id:        MSG_HALL_JOIN_FKROOM_REPLY,
		Result:    0,
		FkNo:      req.FkNo,
		GameType:  gt.GameType,
		RoomType:  gt.RoomType,
		GradeType: gt.GradeType,
	}, false)
}

////////////////////////////////////////////////////////////////////////////////////
//数据库接口
func (this *Hall) GetPlayer(act string, token string) (*DbPlayerData, error) {
	//接入数据库模式
	if DB_SAVE_MODEL == 1 {
		getip := GCONFIG.WebDbIp + "/V1.0/App/UserInfo"
		re, err := SendRequest(getip, nil, "GET", token)
		if err != nil {
			return nil, err
		}
		rsp := DbGetPlayerRsp{}
		err = json.Unmarshal([]byte(re), &rsp)
		if err != nil {
			return nil, err
		}
		if rsp.Code != 200 {
			return nil, fmt.Errorf("http错误:", rsp.Code, rsp.Msg)
		}
		return &rsp.Data, nil
	} else {
		//接入内存模式
		this.DbLk.RLock()
		p, ok := this.DbData[act]
		if !ok {
			this.DbLk.RUnlock()
			return nil, fmt.Errorf("player no exist")
		}
		this.DbLk.RUnlock()
		return p, nil
	}
}

func (this *Hall) GetRobot(token string) (*DbRobotData, error) {
	getip := GCONFIG.WebRobotIp + "/V1.0/robot/robotInfo"

	backstageToken := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJ0b2tlbiIsImlhdCI6MTU1Mjg5ODIyMn0.uSfTj7U9XcEAHddfjv742Ti44ZthUhRTWk4_ke1s9R8"
	url := fmt.Sprintf("%s?token=%s", getip, token)
	re, err := SendRequest(url, nil, "GET", backstageToken)

	if err != nil {
		return nil, err
	}
	rsp := DbGetRobotRsp{}
	err = json.Unmarshal([]byte(re), &rsp)
	if err != nil {
		return nil, err
	}
	if rsp.Code != 200 {
		return nil, fmt.Errorf("http错误:", rsp.Code, rsp.Msg)
	}
	return &rsp.Data, nil
}

func (this *Hall) SavePlayer(act string, p *DbPlayerData) {
	this.DbLk.Lock()
	this.DbData[act] = p
	this.DbData2[p.Uid] = p
	this.DbLk.Unlock()
}

func (this *Hall) SetPlayerGameInfo(act string, gameId int32, serverId int32) {
	this.DbLk.Lock()
	p, ok := this.DbData[act]
	if ok {
		p.GameInfoId = int(gameId)
		p.GameServerId = int(serverId)
	}
	this.DbLk.Unlock()
}

func (this *Hall) SetGameIdByUid(uid int64, id int32, serverId int32) error {
	if DB_SAVE_MODEL == 1 {
		getip := GCONFIG.WebDbIp + "/V1.0/App/EnterGame"
		req := struct {
			Uid          int64 `json:"userId"`
			GameInfoId   int   `json:"gameInfoId"`
			GameServerId int   `json:"gameServerId"`
		}{
			Uid:          uid,
			GameInfoId:   int(id),
			GameServerId: int(serverId),
		}
		token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJoYWxscmVxdWVzdGhvbWUiLCJpYXQiOjE1NTA4MTYxOTR9.m8zSOaTzkXKTKoUuEISF-fKryM2KDKPO-_YwEXizN54"
		re, err := SendRequest(getip, &req, "POST", token)
		if err != nil {
			logs.Error("保存游戏信息失败", token)
			return err
		}
		rsp := struct {
			Msg  string `json:"msg"`
			Code int    `json:"code"`
		}{}
		json.Unmarshal([]byte(re), &rsp)
		if rsp.Code != 200 {
			logs.Error("保存游戏信息失败", rsp)
			return fmt.Errorf("保存游戏信息时错误", rsp)
		}
		return nil
	} else {
		this.DbLk.Lock()
		p, ok := this.DbData2[uid]
		if ok {
			p.GameInfoId = int(id)
			p.GameServerId = int(serverId)
		}
		this.DbLk.Unlock()
		return nil
	}
}

func (this *Hall) SetGameIdByRobotUid(uid int64, id int32, serverId int32) error {
	if DB_SAVE_MODEL == 1 {
		getip := GCONFIG.WebRobotIp + "/V1.0/robot/EnterGame"
		req := struct {
			GameInfoId   int `json:"gameInfoId"`
			GameServerId int `json:"gameServerId"`
		}{
			GameInfoId:   int(id),
			GameServerId: int(serverId),
		}
		token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJoYWxscmVxdWVzdGhvbWUiLCJpYXQiOjE1NTA4MTYxOTR9.m8zSOaTzkXKTKoUuEISF-fKryM2KDKPO-_YwEXizN54"
		re, err := SendRequest(getip, &req, "POST", token)
		if err != nil {
			logs.Error("保存Robot游戏信息失败", token)
			return err
		}
		rsp := struct {
			Msg  string `json:"msg"`
			Code int    `json:"code"`
		}{}
		json.Unmarshal([]byte(re), &rsp)
		if rsp.Code != 200 {
			logs.Error("保存Robot游戏信息失败", rsp, re)
			return fmt.Errorf("保存Robot游戏信息时错误", rsp)
		}
		return nil
	} else {
		this.DbLk.Lock()
		p, ok := this.DbData2[uid]
		if ok {
			p.GameInfoId = int(id)
			p.GameServerId = int(serverId)
		}
		this.DbLk.Unlock()
		return nil
	}
}

func (this *Hall) PushRecord(d interface{}) error {
	if DB_SAVE_MODEL == 1 {
		logs.Debug("发送保存记录")
		token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJoYWxscmVxdWVzdGhvbWUiLCJpYXQiOjE1NTA4MTYxOTR9.m8zSOaTzkXKTKoUuEISF-fKryM2KDKPO-_YwEXizN54"
		getip := GCONFIG.WebDbIp + "/V1.0/Game/GameRecord"
		re, err := SendRequest(getip, d, "POST", token)
		if err != nil {
			logs.Error("保存游戏记录失败")
			return err
		}
		rsp := struct {
			Msg  string `json:"msg"`
			Code int    `json:"code"`
		}{}
		json.Unmarshal([]byte(re), &rsp)
		if rsp.Code != 200 {
			logs.Error("保存游戏记录错误", rsp)
			return fmt.Errorf("保存游戏记录错误", rsp)
		}
		return nil
	}
	return nil
}

func (this *Hall) PushRobotRecord(d interface{}) error {
	if DB_SAVE_MODEL == 1 {
		logs.Debug("发送机器人保存记录")
		token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJ0b2tlbiIsImlhdCI6MTU1Mjg5ODIyMn0.uSfTj7U9XcEAHddfjv742Ti44ZthUhRTWk4_ke1s9R8"
		getip := GCONFIG.WebRobotIp + "/V1.0/robotGame/GameRecord"
		re, err := SendRequest(getip, d, "POST", token)
		if err != nil {
			logs.Error("机器人保存记录失败")
			return err
		}
		rsp := struct {
			Msg  string `json:"msg"`
			Code int    `json:"code"`
		}{}
		json.Unmarshal([]byte(re), &rsp)
		if rsp.Code != 200 {
			logs.Error("机器人保存记录失败", rsp)
			return fmt.Errorf("机器人保存记录失败", rsp)
		}
		return nil
	}
	return nil
}
func (this *Hall) PushOtherRecord(d string) error {
	if DB_SAVE_MODEL == 1 {
		logs.Debug("发送斗地主的保存记录")
		getip := GCONFIG.WebDbIp + "/V1.0/Game/GameRecordDetails"
		token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJoYWxscmVxdWVzdGhvbWUiLCJpYXQiOjE1NTA4MTYxOTR9.m8zSOaTzkXKTKoUuEISF-fKryM2KDKPO-_YwEXizN54"
		re, err := SendRequestByString(getip, d, "POST", token)
		if err != nil {
			logs.Error("保存游戏详情失败")
			return err
		}
		rsp := struct {
			Msg  string `json:"msg"`
			Code int    `json:"code"`
		}{}
		json.Unmarshal([]byte(re), &rsp)
		if rsp.Code != 200 {
			logs.Error("保存游戏详情错误", rsp)
			return fmt.Errorf("保存游戏详情错误", rsp)
		}
		return nil
	}
	return nil
}

func (this *Hall) SetCoin(act string, uid, coin int64) error {
	if DB_SAVE_MODEL == 1 {
		getip := GCONFIG.WebDbIp + "/V1.0/App/UserCoin"
		req := struct {
			Uid  int64 `json:"userId"`
			Coin int64 `json:"coin"`
		}{
			Uid:  uid,
			Coin: coin,
		}
		logs.Debug("发送修改金币")
		token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJoYWxscmVxdWVzdGhvbWUiLCJpYXQiOjE1NTA4MTYxOTR9.m8zSOaTzkXKTKoUuEISF-fKryM2KDKPO-_YwEXizN54"
		re, err := SendRequest(getip, &req, "POST", token)
		if err != nil {
			logs.Error("保存游戏信息失败")
			return err
		}
		rsp := struct {
			Msg  string `json:"msg"`
			Code int    `json:"code"`
		}{}
		json.Unmarshal([]byte(re), &rsp)
		if rsp.Code != 200 {
			logs.Error("保存游戏信息失败", rsp)
			return fmt.Errorf("保存游戏信息时错误", rsp)
		}
		return nil
	} else {
		this.DbLk.Lock()
		p, ok := this.DbData[act]
		if ok {
			p.Coin += coin
		}
		this.DbLk.Unlock()
		return nil
	}
}

func (this *Hall) SetRoomCard(act string, uid, card int64) error {
	logs.Debug("修改房卡", uid, card)
	if DB_SAVE_MODEL == 1 {
		getip := GCONFIG.WebDbIp + "/V1.0/App/UseRoomCard"
		req := struct {
			Uid      int64 `json:"userId"`
			RoomCard int64 `json:"roomCard"`
		}{
			Uid:      uid,
			RoomCard: card,
		}
		logs.Debug("发送修改房卡")
		token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJoYWxscmVxdWVzdGhvbWUiLCJpYXQiOjE1NTA4MTYxOTR9.m8zSOaTzkXKTKoUuEISF-fKryM2KDKPO-_YwEXizN54"
		re, err := SendRequest(getip, &req, "POST", token)
		if err != nil {
			logs.Error("发送修改房卡失败")
			return err
		}
		rsp := struct {
			Msg  string `json:"msg"`
			Code int    `json:"code"`
		}{}
		json.Unmarshal([]byte(re), &rsp)
		if rsp.Code != 200 {
			logs.Error("发送修改房卡失败", rsp)
			return fmt.Errorf("发送修改房卡失败", rsp)
		}
		return nil
	} else {
		this.DbLk.Lock()
		p, ok := this.DbData[act]
		if ok {
			p.RoomCard += card
		}
		this.DbLk.Unlock()
		return nil
	}
}

func (this *Hall) Leave(msg *InMsg) {
	// GDeskMgr.AddMsg(MSG_GAME_LEAVE, msg.Uid, "")
	this.Lk.RLock()
	p, ok := this.Players[msg.Uid]
	this.Lk.RUnlock()
	if ok {
		if p.GameId != 0 {
			GDeskMgr.AddMsg(MSG_GAME_DISCONNECT, msg.Uid, "")
		}
		this.Lk.Lock()
		delete(this.Players, msg.Uid)
		this.Lk.Unlock()
		p.AddMsgNative(MSG_HALL_LEAVE_REPLY, nil, true)
	}
}

func GetGameDetailByWeb(roomtype, gradetype, gametype int32) (*DbGetGameServerData, error) {
	getip := GCONFIG.WebGameIp + "/web/gameInfo/gameinfo?"
	getip = getip + "gradeId=" + strconv.Itoa(int(gradetype))
	getip = getip + "&gameId=" + strconv.Itoa(int(gametype))
	getip = getip + "&roomId=" + strconv.Itoa(int(roomtype))
	token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJoYWxscmVxdWVzdHNtYWxsIiwiaWF0IjoxNTUwODE1Mzc4fQ.mNr0cT3i-F-E5twxX3RSbtFd7tgfhFRr2kbZY8o3RLQ"
	re, err := SendRequest(getip, nil, "GET", token)
	if err != nil {
		return nil, err
	}
	rsp := DbGetGameServerRsp{}
	json.Unmarshal([]byte(re), &rsp)
	if rsp.Code != 200 || len(rsp.Data.Game) == 0 {
		return nil, fmt.Errorf("数据库无此数据")
	}
	//
	resp := rsp.Data.Game[0]
	return &resp, nil
}

//获取房间号
func (this *Hall) CreateRoomNo(gt *GameTypeDetail) string {
	for {
		rno := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
		_, ok := this.MRoomNo[rno]
		if !ok {
			this.MRoomNo[rno] = gt
			return rno
		}
	}
}

func (this *Hall) GetFkRoom(no string) *GameTypeDetail {
	v, ok := this.MRoomNo[no]
	if ok {
		return v
	}
	return nil
}

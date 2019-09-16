package main

import (
	"encoding/json"
	"fmt"
	"logs"
	"strconv"
	"time"
)

func (this *Roboter) InitHandle() {
	this.Handle[MSG_HALL_LOGIN_REPLY] = this.HandleHallLoginReply
	//
	this.Handle[MSG_GAME_INFO_AUTO_REPLY] = this.HandleGameAutoReply
	this.Handle[MSG_GAME_AUTO_REPLY] = this.HandleAutoReply
	this.Handle[MSG_GAME_INFO_START] = this.HandleGameStart
	this.Handle[MSG_GAME_INFO_SEND_NOTIFY] = this.HandleGameSend
	this.Handle[MSG_GAME_INFO_CALL_REPLY] = this.handleCallReply
	this.Handle[MSG_GAME_INFO_BANKER_NOTIFY] = this.HandleBankerNotify
	this.Handle[MSG_GAME_INFO_OUTCARD_REPLY] = this.OutCardReply
	this.Handle[MSG_GAME_INFO_PASS_REPLY] = this.PassReply
	this.Handle[MSG_GAME_INFO_END_NOTIFY] = this.GameEndNotify
	this.Handle[MSG_GAME_INFO_RECONNECT_REPLY] = this.GameReconnect
}

func (this *Roboter) HandleHallLoginReply(data []byte) {
	d := HMsgHallLoginReply{}
	json.Unmarshal(data, &d)
	if d.Result != 0 {
		time.Sleep(time.Second * 3)
		this.SendLogin()
		return
	}
	this.Coins = d.Coin
	this.Uid = d.Uid
	if this.XuHao == this.ShowId {
		logs.Debug("登录大厅收到的消息", d)
	}
	//
	if d.GameSerId == 0 {
		sd := GAutoGame2{
			Id:        MSG_GAME_AUTO,
			GameType:  1,
			RoomType:  1,
			GradeType: 1,
		}
		this.SendToServer(sd)
	} else {
		sd := GReconnect{
			Id: MSG_GAME_RECONNECT,
		}
		this.SendToServer(sd)
	}
}

func (this *Roboter) HandleGameAutoReply(data []byte) {
	d := GInfoAutoGameReply{}
	json.Unmarshal(data, &d)
	for _, v := range d.Seat {
		if v.Uid == this.Uid {
			this.SeatId = v.Cid
		}
	}
	if this.XuHao == this.ShowId {
		logs.Debug("自动匹配收到的消息", d)
	}
}

func (this *Roboter) HandleAutoReply(data []byte) {
	if this.XuHao == this.ShowId {
		logs.Debug("自动匹配收到的消息结果", string(data))
	}
	re := GAutoGameReply{}
	json.Unmarshal(data, &re)
	if re.Result != 0 {
		time.Sleep(time.Second * 3)
		sd := GAutoGame2{
			Id:        MSG_GAME_AUTO,
			GameType:  1,
			RoomType:  1,
			GradeType: 1,
		}
		this.SendToServer(sd)
	}

}

func (this *Roboter) HandleGameStart(data []byte) {
	if this.XuHao == this.ShowId {
		logs.Debug("游戏开始：", string(data))
		this.MMaxOut = nil
	}
}

func (this *Roboter) HandleGameSend(data []byte) {
	if this.XuHao == this.ShowId {
		logs.Debug("游戏发牌开始", string(data))
	}
	//
	d := GGameSendCardNotify{}
	json.Unmarshal(data, &d)
	this.HandCard = Sort(d.HandsCards)
	//开始叫分
	this.CallTimes = 0
	this.CallEd = false
	if this.SeatId == d.Cid {
		time.Sleep(time.Second * 1)
		sd := GCallMsg{
			Id:    MSG_GAME_INFO_CALL,
			Coins: -1,
		}
		this.SendToServer(sd)
	}
}

func (this *Roboter) handleCallReply(data []byte) {
	re := GCallMsgReply{}
	json.Unmarshal(data, &re)
	if this.XuHao == this.ShowId {
		logs.Debug("有人叫分", re, this.SeatId)
	}
	this.CallTimes++
	if this.CallTimes >= 3 || re.Coins >= 3 {
		return
	}
	//
	if (re.Cid+1)%3 == this.SeatId {
		time.Sleep(time.Second * 1)
		sd := GCallMsg{
			Id:    MSG_GAME_INFO_CALL,
			Coins: re.Coins + 2,
		}
		this.SendToServer(sd)
	}
}

func (this *Roboter) HandleBankerNotify(data []byte) {
	if this.XuHao == this.ShowId {
		logs.Debug("定庄通知")
	}
	re := GBankerNotify{}
	json.Unmarshal(data, &re)
	//
	if this.SeatId == re.Banker {
		this.HandCard = append(this.HandCard, re.DiPai...)
		//出牌
		out := GGameOutCard{
			Id: MSG_GAME_INFO_OUTCARD,
		}
		ok := GetOutCard(this.HandCard, &out)
		if ok {
			this.HandCard, _ = VecDelMulti(this.HandCard, out.Cards)
			rout := GRealOutCard{
				Id:   MSG_GAME_INFO_OUTCARD,
				Type: out.Type,
			}
			for _, v := range out.Cards {
				rout.Cards = append(rout.Cards, int32(v))
			}
			this.SendToServer(rout)
		}
	}
}

func (this *Roboter) OutCardReply(data []byte) {
	re := GGameOutCardReply{}
	json.Unmarshal(data, &re)
	if this.XuHao == this.ShowId {
		v, _ := G_OutMap[re.Type]
		info := "出牌玩家：" + strconv.Itoa(int(re.Cid)) + "---牌：" + v + strconv.Itoa(int(GetLogicValue(re.Max)))
		info = info + fmt.Sprintf("---牌数据：", re.Cards)
		logs.Debug(info)
	}
	this.MMaxOut = &re
	//轮到自己出牌
	if (re.Cid+1)%3 == this.SeatId {
		time.Sleep(time.Second)
		//出牌
		out := GGameOutCard{
			Id: MSG_GAME_INFO_OUTCARD,
		}
		if GetSecondOutCard(this.HandCard, &out, this.MMaxOut) {
			this.HandCard, _ = VecDelMulti(this.HandCard, out.Cards)
			rout := GRealOutCard{
				Id:   MSG_GAME_INFO_OUTCARD,
				Type: out.Type,
			}
			for _, v := range out.Cards {
				rout.Cards = append(rout.Cards, int32(v))
			}
			// logs.Debug("出牌", rout)
			this.SendToServer(rout)
		} else {
			this.Pass()
		}
	}
}

func (this *Roboter) GameEndNotify(data []byte) {
	if this.XuHao == this.ShowId {
		logs.Debug("游戏结束", string(data))
	}
	//初始化数据
	this.HandCard = []byte{}
	this.SeatId = -2
	this.CurId = -2
	this.CallTimes = 0
	this.CallEd = false
	//
	time.Sleep(time.Second * 3)
	if this.XuHao == this.ShowId {
		logs.Debug("重新匹配")
	}
	sd := GAutoGame2{
		Id:        MSG_GAME_AUTO,
		GameType:  1,
		RoomType:  1,
		GradeType: 1,
	}
	this.SendToServer(sd)
}

func (this *Roboter) PassReply(data []byte) {
	re := GGamePassReply{}
	json.Unmarshal(data, &re)

	if this.XuHao == this.ShowId {
		logs.Debug("有人过", string(data))
	}
	//轮到自己出牌
	if (re.Cid+1)%3 == this.SeatId {
		time.Sleep(time.Second)
		if this.MMaxOut == nil || this.SeatId == this.MMaxOut.Cid {
			//出牌
			out := GGameOutCard{
				Id: MSG_GAME_INFO_OUTCARD,
			}
			ok := GetOutCard(this.HandCard, &out)
			if ok {
				this.HandCard, _ = VecDelMulti(this.HandCard, out.Cards)
				rout := GRealOutCard{
					Id:   MSG_GAME_INFO_OUTCARD,
					Type: out.Type,
				}
				for _, v := range out.Cards {
					rout.Cards = append(rout.Cards, int32(v))
				}
				this.SendToServer(rout)
			}
		} else {
			//出牌
			out := GGameOutCard{
				Id: MSG_GAME_INFO_OUTCARD,
			}
			if GetSecondOutCard(this.HandCard, &out, this.MMaxOut) {
				this.HandCard, _ = VecDelMulti(this.HandCard, out.Cards)
				rout := GRealOutCard{
					Id:   MSG_GAME_INFO_OUTCARD,
					Type: out.Type,
				}
				for _, v := range out.Cards {
					rout.Cards = append(rout.Cards, int32(v))
				}
				this.SendToServer(rout)
			} else {
				this.Pass()
			}
		}
	}
}

func (this *Roboter) Pass() {
	sd := GGamePass{
		Id: MSG_GAME_INFO_PASS,
	}
	this.SendToServer(sd)
}

func (this *Roboter) SendLogin() {
	var reqlogin HMsgHallLogin
	reqlogin.Account = this.Account
	reqlogin.Gid = ""
	reqlogin.Id = MSG_HALL_LOGIN // 登录Id号
	this.SendToServer(&reqlogin)
}

func (this *Roboter) GameReconnect(data []byte) {
	re := GInfoReConnectReply{}
	json.Unmarshal(data, &re)
	// logs.Debug("断线重连", string(data))
	if this.XuHao == this.ShowId {
		logs.Debug("断线重连", string(data))
	}

	if re.GameState == GAME_STATUS_FREE || re.GameState == GAME_STATUS_END {
		return
	}
	this.MGameState = re.GameState
	this.SeatId = re.Cid
	this.CurId = re.CurCid
	this.HandCard = re.Cards
	if len(re.OutEd) == 0 {
		this.MMaxOut = nil
	} else {
		this.MMaxOut = &GGameOutCardReply{
			Cid:   re.OutEd[len(re.OutEd)-1].Cid,
			Type:  re.OutEd[len(re.OutEd)-1].Type,
			Cards: re.OutEd[len(re.OutEd)-1].Cards,
			Max:   re.OutEd[len(re.OutEd)-1].Max,
		}
	}

	//
	if re.GameState == GAME_STATUS_CALL {
		if re.CurCid == this.SeatId {
			time.Sleep(time.Second * 1)
			sd := GCallMsg{
				Id:    MSG_GAME_INFO_CALL,
				Coins: 3,
			}
			this.SendToServer(sd)
		}
	} else if re.GameState == GAME_STATUS_PLAY {
		if re.CurCid == this.SeatId {
			if this.MMaxOut == nil {
				//出牌
				out := GGameOutCard{
					Id: MSG_GAME_INFO_OUTCARD,
				}
				ok := GetOutCard(this.HandCard, &out)
				if ok {
					this.HandCard, _ = VecDelMulti(this.HandCard, out.Cards)
					rout := GRealOutCard{
						Id:   MSG_GAME_INFO_OUTCARD,
						Type: out.Type,
					}
					for _, v := range out.Cards {
						rout.Cards = append(rout.Cards, int32(v))
					}
					this.SendToServer(rout)
				}
			} else {
				//出牌
				out := GGameOutCard{
					Id: MSG_GAME_INFO_OUTCARD,
				}
				if GetSecondOutCard(this.HandCard, &out, this.MMaxOut) {
					this.HandCard, _ = VecDelMulti(this.HandCard, out.Cards)
					rout := GRealOutCard{
						Id:   MSG_GAME_INFO_OUTCARD,
						Type: out.Type,
					}
					for _, v := range out.Cards {
						rout.Cards = append(rout.Cards, int32(v))
					}
					this.SendToServer(rout)
				} else {
					this.Pass()
				}
			}
		}
	}
}

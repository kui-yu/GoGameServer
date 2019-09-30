package main

import (
	"encoding/json"
	"math/rand"
	"time"
	// "github.com/tidwall/gjson"
)

//匹配成功发回的座位信息
func (this *ExtRobotClient) HandleGameAutoReply(d string) {
	rsp := GInfoAutoGameReply{}
	json.Unmarshal([]byte(d), &rsp)
	for _, v := range rsp.Seat {
		if v.Uid == this.UserInfo.Uid {
			this.SeatId = v.Cid
		}
	}
	//已经在游戏中
	this.GameIn = true
}

//收到发牌信息
func (this *ExtRobotClient) HandleGameSend(data string) {
	d := GGameSendCardNotify{}
	json.Unmarshal([]byte(data), &d)
	this.HandCard = Sort(d.HandsCards)
	//开始叫分
	this.CallTimes = 0
	this.CallEd = false
	if this.SeatId == d.Cid {
		stime := rand.Intn(2) + 3
		time.Sleep(time.Second * time.Duration(stime))
		sd := GCallMsg{
			Id:    MSG_GAME_INFO_CALL,
			Coins: -1,
		}
		this.AddMsgNative(MSG_GAME_INFO_CALL, &sd)
	}
}

//叫分
func (this *ExtRobotClient) handleCallReply(data string) {
	re := GCallMsgReply{}
	json.Unmarshal([]byte(data), &re)
	this.CallTimes++
	if this.CallTimes >= 3 || re.Coins >= 3 {
		return
	}
	//
	if (re.Cid+1)%3 == this.SeatId {
		stime := rand.Intn(4) + 1
		time.Sleep(time.Second * time.Duration(stime))
		sd := GCallMsg{
			Id:    MSG_GAME_INFO_CALL,
			Coins: re.Coins + 2,
		}
		this.AddMsgNative(MSG_GAME_INFO_CALL, &sd)
	}
}

//定庄通知
func (this *ExtRobotClient) HandleBankerNotify(data string) {
	re := GBankerNotify{}
	json.Unmarshal([]byte(data), &re)
	//
	if this.SeatId == re.Banker {
		this.HandCard = append(this.HandCard, re.DiPai...)
		//出牌
		stime := rand.Intn(4) + 1
		time.Sleep(time.Second * time.Duration(stime))
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
			this.AddMsgNative(MSG_GAME_INFO_OUTCARD, &rout)
		}
	}
}

func (this *ExtRobotClient) OutCardReply(data string) {
	re := GGameOutCardReply{}
	json.Unmarshal([]byte(data), &re)
	this.MMaxOut = &re
	//轮到自己出牌
	if (re.Cid+1)%3 == this.SeatId {
		stime := rand.Intn(4) + 1
		time.Sleep(time.Second * time.Duration(stime))
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
			this.AddMsgNative(MSG_GAME_INFO_OUTCARD, &rout)
		} else {
			this.Pass()
		}
	}
}

func (this *ExtRobotClient) Pass() {
	sd := GGamePass{
		Id: MSG_GAME_INFO_PASS,
	}
	this.AddMsgNative(MSG_GAME_INFO_PASS, &sd)
}

func (this *ExtRobotClient) PassReply(data string) {
	re := GGamePassReply{}
	json.Unmarshal([]byte(data), &re)
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
				this.AddMsgNative(MSG_GAME_INFO_OUTCARD, &rout)
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
				this.AddMsgNative(MSG_GAME_INFO_OUTCARD, &rout)
			} else {
				this.Pass()
			}
		}
	}
}

func (this *ExtRobotClient) GameEndNotify(data string) {
	//
	rsp := GInfoGameEnd{}
	json.Unmarshal([]byte(data), &rsp)
	//
	this.GameIn = false

	for i, v := range rsp.Coins {
		if int32(i) == this.SeatId {
			this.Coin = v
			break
		}
	}
	//
	controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
}

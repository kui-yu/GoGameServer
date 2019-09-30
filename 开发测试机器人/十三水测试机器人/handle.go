package main

import (
	"encoding/json"
	// "logs"
	// "math/rand"
	// "strconv"
	// "time"
)

func (this *Roboter) InitHandle() {
	//大厅登陆
	this.Handle[MSG_HALL_LOGIN_REPLY] = this.HandleHallLoginReply
	//匹配应答
	this.Handle[MSG_GAME_INFO_AUTO_REPLY] = this.HandleGameAutoReply
	//重新开始
	this.Handle[MSG_GAME_INFO_SETTLE_INFO_REPLY] = this.GameEndNotify
	//
	this.Handle[MSG_GAME_INFO_STAGE] = this.HandleGameStage
}

//状态消息
func (this *Roboter) HandleGameStage(data []byte) {
	d := GStageInfo{}
	json.Unmarshal(data, &d)
	// logs.Debug("阶段", d, this.SeatId)

	if d.Stage == STAGE_PLAY {

		callMsg := GAPlayInfo{
			Id:       MSG_GAME_INFO_PLAY,
			PlayType: 1,
		}
		// logs.Debug("叫分", randNum)
		this.SendToServer(callMsg)
	}
}

//大厅登陆
func (this *Roboter) HandleHallLoginReply(data []byte) {
	d := HMsgHallLoginReply{}
	json.Unmarshal(data, &d)
	// logs.Debug("登录大厅收到的消息", d)
	this.Coins = d.Coin
	this.Uid = d.Uid

	if d.GameId == 0 {
		sd := GAutoGame2{
			Id: MSG_GAME_AUTO,
		}
		this.SendToServer(sd)
	} else {
		sd := GReconnect{
			Id: MSG_GAME_RECONNECT,
		}
		this.SendToServer(sd)
	}
}

//匹配应答
func (this *Roboter) HandleGameAutoReply(data []byte) {
	d := GInfoAutoGameReply{}
	json.Unmarshal(data, &d)
	for _, v := range d.Seat {
		if v.Uid == this.Uid {
			this.SeatId = v.Cid
		}
	}
	if this.XuHao == this.ShowId {
		// logs.Debug("自动匹配收到的消息", d)
	}
}

//游戏结束
func (this *Roboter) GameEndNotify(data []byte) {
	d := GSSettleInfos{}
	json.Unmarshal(data, &d)
	// logs.Debug("收到结束的消息", d)

	//初始化数据
	this.HandCard = []byte{}
	this.SeatId = -2
	this.CurId = -2
	//
	// time.Sleep(time.Second * 3)

	sd := GAutoGame{
		Id:      MSG_GAME_AUTO,
		Account: this.Account,
		Uid:     this.Uid,
		Nick:    this.Name,
		Sex:     this.Sex,
		Head:    "",
		Lv:      1,
		Coin:    this.Coins,
	}
	this.SendToServer(sd)
}

package main

import (
	// "encoding/json"
	// . "MaJiangTool"
	// "logs"
	"math/rand"
)

func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {
	// //
	// if p.ChairId == 1 {
	// 	testcard := []byte{1, 1, 1, 2, 2, 2, 3, 3, 5, 5, 6, 6, 7, 7}
	// 	se := CreateSendCardEvent(1, 3, false)
	// 	this.HuPai.GetResult(1, []FuZi{}, testcard, se)
	// }
	//
	var robotCnt int = 0
	var playerCnt int = 0
	//机器人不能超过4个
	for _, v := range this.Players {
		if v.Robot {
			robotCnt++
		} else {
			playerCnt++
		}
	}
	if robotCnt >= GCONFIG.PlayerNum || playerCnt == 0 {
		//发送匹配成功
		p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
			Id:     MSG_GAME_AUTO_REPLY,
			Result: 13,
		})
		//踢出
		this.LeaveByForce(p)
		return
	}
	//发送匹配成功
	p.QueColor = 255
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:     MSG_GAME_AUTO_REPLY,
		Result: 0,
	})
	//群发用户信息
	for _, v := range this.Players {
		gameReply := GInfoAutoGameReply{
			Id: MSG_GAME_INFO_AUTO_REPLY,
		}
		for _, p := range this.Players {
			seat := GSeatInfo{
				Uid:  p.Uid,
				Nick: p.Nick,
				Cid:  p.ChairId,
				Sex:  p.Sex,
				Head: p.Head,
				Lv:   p.Lv,
				Coin: p.Coins,
			}
			if p.Uid != v.Uid {
				seat.Nick = "***" + seat.Nick[len(seat.Nick)-4:]
			}
			gameReply.Seat = append(gameReply.Seat, seat)
		}
		v.SendNativeMsg(MSG_GAME_INFO_AUTO_REPLY, &gameReply)
	}
	//发送房间信息
	this.MaxDouble = G_DbGetGameServerData.MaxTimes
	this.Bscore = G_DbGetGameServerData.Bscore
	if this.JuHao == "" {
		this.JuHao = GetJuHao()
		// juhao, _ := uuid.NewV4()
		// this.JuHao = juhao.String()
	}

	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &GGameInfoNotify{
		Id:     MSG_GAME_INFO_ROOM_NOTIFY,
		BScore: this.Bscore,
		MaxBei: this.MaxDouble,
		JuHao:  this.JuHao,
	})
	//判断人员是否已满，开启游戏
	if len(this.Players) >= GCONFIG.PlayerNum {
		//阶段消息。游戏开始
		this.GameState = GAME_STATUS_START
		this.BroadStageTime(TIMER_START_NUM)
		//几秒后进入发牌阶段
		this.AddTimer(TIMER_START, TIMER_START_NUM, this.TimerSendCard, nil)
	}
}

func (this *ExtDesk) TimerSendCard(d interface{}) {
	//定庄
	this.Banker = rand.Intn(10) % len(this.Players)
	this.Banker = 0
	this.CurCid = this.Banker
	//
	this.GameState = GAME_STATE_SENDCARD
	this.BroadStageTime(TIMER_SENDCARD_NUM) //广播这个阶段的时间
	//洗牌
	this.CardMgr.Shuffle()
	// 发牌 17 张
	for _, v := range this.Players {
		v.HandCard = []byte{}
		this.CardMgr.SendStartCard(&v.HandCard)
	}
	//庄家多发一张
	bankerCard := this.CardMgr.SendCard(false)
	this.Players[this.Banker].HandCard = append(this.Players[this.Banker].HandCard, bankerCard)
	//发送牌消息
	sd := GGameSendCardNotify{
		Id:     MSG_GAME_INFO_SEND_NOTIFY,
		Banker: this.Banker,
	}
	for _, v := range this.Players {
		sd.HandsCards = []int{}
		for _, c := range v.HandCard {
			sd.HandsCards = append(sd.HandsCards, int(c))
		}
		v.SendNativeMsg(MSG_GAME_INFO_SEND_NOTIFY, &sd)
	}

	//发牌后进入进入换牌阶段
	this.AddTimer(TIMER_SENDCARD, TIMER_SENDCARD_NUM, this.TimerHuanPai, nil)
}

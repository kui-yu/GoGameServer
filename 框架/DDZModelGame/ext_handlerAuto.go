package main

import (
	// "logs"
	"math/rand"
)

func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {
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
	if robotCnt >= 4 || playerCnt == 0 {
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
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:     MSG_GAME_AUTO_REPLY,
		Result: 0,
	})
	//群发用户信息
	// gameReply := GInfoAutoGameReply{
	// 	Id: MSG_GAME_INFO_AUTO_REPLY,
	// }
	// for _, v := range this.Players {
	// 	seat := GSeatInfo{
	// 		Uid:  v.Uid,
	// 		Nick: v.Nick,
	// 		Cid:  v.ChairId,
	// 		Sex:  v.Sex,
	// 		Head: v.Head,
	// 		Lv:   v.Lv,
	// 		Coin: v.Coins,
	// 	}
	// 	gameReply.Seat = append(gameReply.Seat, seat)
	// }
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
	// this.BroadcastAll(MSG_GAME_INFO_AUTO_REPLY, &gameReply)
	//发送房间信息
	this.MaxDouble = int32(G_DbGetGameServerData.MaxTimes)
	this.Bscore = int32(G_DbGetGameServerData.Bscore)
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

		this.GameState = GAME_STATUS_START
		this.BroadStageTime(TIMER_START_NUM)
		//发送游戏开始通知
		sd := GGameStartNotify{
			Id: MSG_GAME_START,
		}
		this.BroadcastAll(MSG_GAME_START, &sd)

		//直接进入发牌
		this.TimerSendCard(nil)
	}
}

func (this *ExtDesk) TimerSendCard(d interface{}) {

	this.CurCid = int32(rand.Intn(10) % len(this.Players))
	//洗牌
	this.CardMgr.Shuffle()
	// 发牌 17 张
	for _, v := range this.Players {
		v.HandCard = []byte{}
		hd := this.CardMgr.SendHandCard(17)
		v.SetHandCard(Sort(hd))
	}
	// 底牌
	this.DiPai = []byte{}
	this.DiPai = this.CardMgr.SendHandCard(3)
	//发送牌消息
	sd := GGameSendCardNotify{
		Id:  MSG_GAME_INFO_SEND_NOTIFY,
		Cid: this.CurCid,
	}
	for _, v := range this.Players {
		sd.HandsCards = v.HandCard
		v.Mgr.SendNativeMsg(MSG_GAME_INFO_SEND_NOTIFY, v.Uid, &sd)
	}
	//发牌后进入叫分阶段，开启叫分阶段的定时器
	this.AddTimer(TIMER_START, TIMER_START_NUM, this.TimerDealPoker, nil)
}

//发牌动画
func (this *ExtDesk) TimerDealPoker(d interface{}) {
	//进入抢庄
	this.GameState = GAME_STATUS_CALL
	this.BroadStageTime(TIMER_CALL_NUM)
	//发牌后进入叫分阶段，开启叫分阶段的定时器
	this.AddTimer(TIMER_CALL, TIMER_CALL_NUM, this.TimerCall, nil)
}

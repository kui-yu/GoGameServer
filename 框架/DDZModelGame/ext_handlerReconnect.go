package main

// import (
// 	"logs"
// )

func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
		p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
			Id:     MSG_GAME_RECONNECT_REPLY,
			Result: 1,
			Err:    "本桌子没有正在的游戏",
		})
		return
	}

	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:     MSG_GAME_RECONNECT_REPLY,
		Result: 0,
	})

	result := GInfoReConnectReply{
		Id: MSG_GAME_INFO_RECONNECT_REPLY,
	}
	result.GameState = int32(this.GameState)
	if this.GameState == GAME_STATUS_CALL {
		result.GameStateTime = int32(TIMER_CALL_NUM)
	} else if this.GameState == GAME_STATUS_PLAY {
		result.GameStateTime = int32(TIMER_OUTCARD_NUM)
	}

	result.Cid = p.ChairId
	for _, v := range this.Players {
		seat := GSeatInfo{
			Uid:  v.Uid,
			Nick: v.Nick,
			Cid:  v.ChairId,
			Sex:  v.Sex,
			Head: v.Head,
			Lv:   v.Lv,
			Coin: v.Coins,
		}
		if len(seat.Nick) > 4 && p.Uid != seat.Uid {
			seat.Nick = "***" + seat.Nick[len(seat.Nick)-4:]
		}
		result.Seats = append(result.Seats, seat)
	}
	//叫分阶段
	result.Cards = append([]byte{}, p.HandCard...)
	for _, v := range this.Players {
		result.CardNum = append(result.CardNum, int32(len(v.HandCard)))
		result.CallFens = append(result.CallFens, v.CFen)
		if v.TuoGuan {
			result.TuoGuans = append(result.TuoGuans, 1)
		} else {
			result.TuoGuans = append(result.TuoGuans, 0)
		}
		if v.LiXian {
			result.LiXians = append(result.LiXians, 1)
		} else {
			result.LiXians = append(result.LiXians, 0)
		}

	}
	result.CurCid = this.CurCid
	result.BScore = this.Bscore
	result.MaxBei = this.MaxDouble
	result.JuHao = this.JuHao
	//
	result.Banker = this.Banker
	result.LastCall = this.CallFen
	result.DiPai = append(result.DiPai, this.DiPai...)
	result.Double = p.Double
	for i := len(this.RdChuPai) - 1; i >= 0; i-- {
		result.OutEd = append(result.OutEd, *this.RdChuPai[i])
		if len(result.OutEd) >= 2 {
			break
		}
	}
	//
	if this.GameState == GAME_STATUS_CALL {
		result.TimerNum = int32(this.GetTimerNum(TIMER_CALL))
	} else if this.GameState == GAME_STATUS_PLAY {
		result.TimerNum = int32(this.GetTimerNum(TIMER_OUTCARD))
	}
	//
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT_REPLY, &result)
	p.LiXian = false
	//
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 1,
	})
}

func (this *ExtDesk) HandleDisConnect(p *ExtPlayer, d *DkInMsg) {
	//广播给其他人，掉线
	if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
		this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Result: 0,
			Token:  p.Token,
		})
		this.DelPlayer(p.Uid)
		this.DeskMgr.LeaveDo(p.Uid)
	} else {
		p.LiXian = true
		this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
			Id:    MSG_GAME_ONLINE_NOTIFY,
			Cid:   p.ChairId,
			State: 2,
		})
	}
}

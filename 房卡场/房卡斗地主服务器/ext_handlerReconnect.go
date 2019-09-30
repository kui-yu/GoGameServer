package main

import (
	"fmt"
	"logs"
)

// import (
// 	"logs"
// )

func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
		p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
			Id:     MSG_GAME_RECONNECT_REPLY,
			Result: 1,
			Err:    "本桌子没有正在的游戏",
		})
		return
	}

	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:       MSG_GAME_RECONNECT_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})

	result := GInfoReConnectReply{
		Id: MSG_GAME_INFO_RECONNECT_REPLY,
	}
	result.GameState = int32(this.GameState)
	if this.GameState == GAME_STATUS_CALL {
		if this.TableConfig.CallType == 1 {
			result.GameStateTime = int32(TIMER_CALL_NUM)
		} else {
			result.GameStateTime = int32(TIMER_GETGMS_NUM)
		}
	} else if this.GameState == GAME_STATUS_PLAY {
		result.GameStateTime = int32(TIMER_OUTCARD_NUM)
	}
	result.CallOrGet = this.CallOrGet
	fmt.Println("下一次叫地主 还是抢地主", result.CallOrGet)
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
		if p.isReady == 0 {
			seat.Ready = false
		} else {
			// fmt.Println("玩家:", v.Nick, "已经准备就绪！！！！！！！！！！！")
			seat.Ready = true
		}
		if len(seat.Nick) > 4 && p.Uid != seat.Uid {
			seat.Nick = "***" + seat.Nick[len(seat.Nick)-4:]
		}
		result.Seats = append(result.Seats, seat)
	}
	//叫分阶段
	for _, v := range p.HandCard {
		result.Cards = append(result.Cards, int(v))
	}
	for _, v := range this.Players {
		fmt.Println("断线重连的时候Calltype", this.TableConfig.CallType)
		result.CardNum = append(result.CardNum, int32(len(v.HandCard)))
		// if v.Uid != p.Uid {
		if this.TableConfig.CallType == 2 {
			logs.Debug("进入抢地主模式的赋值")
			fmt.Println("玩家", v.Nick, "叫分:", v.GetMSG)
			result.CallFens = append(result.CallFens, v.GetMSG)
		} else {
			fmt.Println("进入的是叫分制？")
			result.CallFens = append(result.CallFens, v.CFen)
		}
		// }
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
	result.Round = this.Round
	//
	result.Banker = this.Banker
	result.LastCall = this.CallFen
	var agreetCha []int32
	for _, v := range this.Players {
		if v.IsDimiss == 1 {
			agreetCha = append(agreetCha, v.ChairId)
		}
	}
	result.DisPlayer = agreetCha
	for _, v := range this.DiPai {
		result.DiPai = append(result.DiPai, int(v))
	}
	if this.GameState == GAME_STATUS_PLAY {
		result.LastCall = this.Double
	}
	for i := len(this.RdChuPai) - 1; i >= 0; i-- {
		o1 := GOutCard1{}
		o1.Max = this.RdChuPai[i].Max
		o1.Cid = this.RdChuPai[i].Cid
		o1.Type = this.RdChuPai[i].Type
		for _, v := range this.RdChuPai[i].Cards {
			o1.Cards = append(o1.Cards, int(v))
		}
		result.OutEd = append(result.OutEd, o1)
		if len(result.OutEd) >= 2 {
			break
		}
	}
	//
	if this.GameState == GAME_STATUS_CALL {
		if this.TableConfig.CallType == 1 {
			result.TimerNum = int32(this.GetTimerNum(TIMER_CALL))
		} else {
			result.TimerNum = int32(this.GetTimerNum(TIMER_GETGMS))
		}
	} else if this.GameState == GAME_STATUS_PLAY {
		result.TimerNum = int32(this.GetTimerNum(TIMER_OUTCARD))
	} else if this.GameState == GAME_STATUS_BreakRoomVote {
		result.TimerNum = int32(this.GetTimerNum(TIMER_BREAKROOM))
	}
	//
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT_REPLY, &result)
	gtablejinfoReplay := GTableInfoReply{
		Id:      MSG_GAME_INFO_ROOM_NOTIFY,
		TableId: this.FkNo,
		Config:  this.TableConfig,
	}
	fmt.Println(result.TimerNum, "time---------------------------------------------")
	fmt.Println("当前倍数：：：：：：：：：：：：", result.LastCall)

	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &gtablejinfoReplay)
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
	logs.Debug("接收到断线")
	// if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
	// 	logs.Debug("删除玩家")
	// 	for _, v := range this.Players {
	// 		v.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
	// 			Id:     MSG_GAME_LEAVE_REPLY,
	// 			Cid:    p.ChairId,
	// 			Uid:    p.Uid,
	// 			Result: 0,
	// 			Token:  p.Token,
	// 		})
	// 	}
	// 	this.LeaveByForce(p)
	// } else {
	// 	logs.Debug("不删除玩家")
	// 	p.LiXian = true
	// 	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
	// 		Id:    MSG_GAME_ONLINE_NOTIFY,
	// 		Cid:   p.ChairId,
	// 		State: 2,
	// 	})
	// }
	p.LiXian = true
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 2,
	})
	if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
		if this.FkOwner == p.Uid {
			this.ClearTimer()
			this.GameState = GAME_STATUS_END
			this.BroadStageTime(0)
			//玩家离开
			for _, p := range this.Players {
				p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
					Id:     MSG_GAME_LEAVE_REPLY,
					Result: 0,
					Cid:    p.ChairId,
					Uid:    p.Uid,
					Token:  p.Token,
				})
			}
			this.GameOverLeave()
			//归还桌子
			this.GameState = GAME_STATUS_FREE
			this.ResetTable()
			this.DeskMgr.BackDesk(this)
		} else {
			for _, v := range this.Players {
				v.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
					Id:     MSG_GAME_LEAVE_REPLY,
					Result: 0,
					Cid:    p.ChairId,
					Uid:    p.Uid,
					Token:  p.Token,
					Robot:  p.Robot,
				})
			}
			this.DelPlayer(p.Uid)
			this.DeskMgr.LeaveDo(p.Uid)
		}
	}

}

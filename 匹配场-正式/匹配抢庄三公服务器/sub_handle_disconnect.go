package main

func (this *ExtDesk) HandleDisConnect(p *ExtPlayer, d *DkInMsg) {
	//广播给其他人，掉线
	if this.GameState == GAME_STATUS_FREE {
		this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Result: 0,
			Token:  p.Token,
		})
		this.DelPlayer(p.Uid)
		this.DeskMgr.LeaveDo(p.Uid)
		if len(this.Players) <= 1 {
			this.DelTimer(7) //删除匹配计时器
		}
	} else {
		p.LiXian = true
		this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, Ext_GOnLineNotify{
			Id:    MSG_GAME_ONLINE_NOTIFY,
			Uid:   int32(p.Uid),
			State: 2,
		})
	}
}

//请求离开房间
func (this *ExtDesk) Leave(p *ExtPlayer) bool {
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
		if len(this.Players) <= 1 {
			this.DelTimer(7) //删除匹配计时器
		}
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

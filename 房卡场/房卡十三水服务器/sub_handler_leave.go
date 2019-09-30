package main

import (
	"logs"
)

//玩家离开
func (this *ExtDesk) HandleLeave(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("玩家离开")
	//广播给其他人，掉线
	if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
		if p.ChairId == 0 {
			this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:     MSG_GAME_LEAVE_REPLY,
				Cid:    p.ChairId,
				Uid:    p.Uid,
				Result: 1,
				Token:  p.Token,
				Err:    "房主解散该房间",
			})
			logs.Debug("房主解散该房间")
			this.TimerOver()
		} else {
			this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:     MSG_GAME_LEAVE_REPLY,
				Cid:    p.ChairId,
				Uid:    p.Uid,
				Result: 0,
				Token:  p.Token,
			})
			this.DelPlayer(p.Uid)
			this.DeskMgr.LeaveDo(p.Uid)
		}
	} else {
		this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Result: 2,
			Token:  p.Token,
			Err:    "游戏已开始，请发起解散",
		})
	}
}

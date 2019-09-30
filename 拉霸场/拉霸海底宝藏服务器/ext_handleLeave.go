package main

func (this *ExtDesk) HandleLeave(p *ExtPlayer, d *DkInMsg) {
	//玩家离开
	this.ClearTimer()
	this.GameState = GAME_STATUS_END

	this.BroadcastAll(MSG_GAME_INFO_LEAVE_REPLY, GGameExitNotify{
		Id: MSG_GAME_INFO_LEAVE_REPLY,
	})

	//玩家强制离开
	this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
		Id:     MSG_GAME_LEAVE_REPLY,
		Cid:    p.ChairId,
		Uid:    p.Uid,
		Result: 0,
		Token:  p.Token,
	})

	this.GameOverLeave()
	//归还桌子
	this.GameState = GAME_STATUS_FREE
	this.ResetTable()
	this.DeskMgr.BackDesk(this)
}

func (this *ExtDesk) Leave(p *ExtPlayer) bool {
	//玩家离开
	this.ClearTimer()
	this.GameState = GAME_STATUS_END

	//玩家强制离开
	this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
		Id:     MSG_GAME_LEAVE_REPLY,
		Cid:    p.ChairId,
		Uid:    p.Uid,
		Result: 0,
		Token:  p.Token,
	})

	this.GameOverLeave()
	//归还桌子
	this.GameState = GAME_STATUS_FREE
	this.ResetTable()
	this.DeskMgr.BackDesk(this)

	return true
}

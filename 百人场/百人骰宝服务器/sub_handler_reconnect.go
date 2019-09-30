package main

//掉线重连
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	obj := new(LostConnection)
	arr := p.MaxBet(0)
	p.MaxBetArr = arr
	obj.Id = MSG_GAME_RECONNECT_REPLY
	obj.Big = this.Big
	obj.Small = this.Small
	obj.Odd = this.Odd
	obj.Even = this.Even
	obj.DeskMoney = this.DeskMoney
	obj.History = this.History
	obj.Stage = this.Stage
	obj.MaxBet = arr
	obj.Time = this.TList[0].T
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, obj)
}

//掉线信息
func (this *ExtDesk) HandleDisConnect(p *ExtPlayer, d *DkInMsg) {
	p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
		Id:  MSG_GAME_LEAVE_REPLY,
		Uid: p.Uid,
	})
}

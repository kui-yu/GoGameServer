package main

//请求返回大厅
func (this *ExtDesk) HandleGameBack(p *ExtPlayer, d *DkInMsg) {
	if p.IsBet {
		p.SendNativeMsg(MSG_GAME_INFO_THREE_THMES, &struct {
			Id     int
			ErrStr string
			Code   int
		}{
			Id:     MSG_GAME_INFO_THREE_THMES,
			ErrStr: "已投注不能返回大厅",
			Code:   1,
		})
		return
	} else {
		this.LeaveByForce(p)
		this.IsUpdate(p)
	}
}

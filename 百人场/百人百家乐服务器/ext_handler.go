package main

// 走势图
func (this *ExtDesk) HandleGameRunChart(p *ExtPlayer, d *DkInMsg) {
	this.RLock()
	defer this.RUnlock()

	runChart := GRunChartReply{
		Id:        MSG_GAME_INFO_RUN_CHART_REPLY,
		PRunChart: this.RunChart,
	}

	p.SendNativeMsg(MSG_GAME_INFO_RUN_CHART_REPLY, runChart)
}

// 玩家列表
func (this *ExtDesk) HandleGameUserList(p *ExtPlayer, d *DkInMsg) {
	this.RLock()
	defer this.RUnlock()

	userInfo := this.GetUserList(p)
	msg := GUserInfoReply{
		Id:       MSG_GAME_INFO_USER_LIST_REPLY,
		UserInfo: userInfo,
	}

	p.SendNativeMsg(MSG_GAME_INFO_USER_LIST_REPLY, msg)
}

package main

// 走势图
func (this *ExtDesk) HandleGameRunChart(p *ExtPlayer, d *DkInMsg) {
	this.RLock()
	defer this.RUnlock()
	//改=>发给客户端不是东南西北而是东西南北
	zoushi := append([][]bool{}, this.RunChart...)
	zoushi[1], zoushi[2] = zoushi[2], zoushi[1]
	runChart := GRunChartReply{
		Id:        MSG_GAME_INFO_RUN_CHART_REPLY,
		ARunChart: zoushi,
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

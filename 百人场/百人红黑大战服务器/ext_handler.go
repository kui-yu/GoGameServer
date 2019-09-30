package main

// "logs"

//走势图以及牌型记录
func (this *ExtDesk) HandleGameRunChart(p *ExtPlayer, d *DkInMsg) {
	this.RLock()
	defer this.RUnlock()
	var (
		ChartRedCount   int //走势中红方赢的局数
		ChartBlackCount int //走势中黑方赢的局数
		Rcount          int //近20局红方次数
		Bcount          int //近20局黑方次数
		Chart           []int32
	)
	Rlength := len(this.RunChart)
	Clength := len(this.CardTypeChart)
	RChart := GRunChartReply{
		Id: MSG_GAME_INFO_RUN_CHART_REPLY,
	}

	//判断输赢走势长度是否大于规定次数
	if Rlength > int(144) {
		this.RunChart = this.RunChart[:139]
		for i := 139 - 1; i >= 0; i-- {
			Chart = append(Chart, this.RunChart[i])
		}
		RChart.PRunchart = Chart
		RChart.ChartCount = len(RChart.PRunchart)
		for i, _ := range RChart.PRunchart {
			if RChart.PRunchart[i] == RED {
				ChartRedCount += 1
			} else {
				ChartBlackCount += 1
			}
		}
	} else {
		if Rlength == 0 || Rlength == 1 {
			Chart = this.RunChart
		} else {
			for i := Rlength - 1; i >= 0; i-- {
				Chart = append(Chart, this.RunChart[i])
			}
		}
		RChart.PRunchart = Chart
		RChart.ChartCount = Rlength
		for i, _ := range this.RunChart {
			if RChart.PRunchart[i] == RED {
				ChartRedCount += 1
			} else {
				ChartBlackCount += 1
			}

		}
	}
	RChart.ChartRedCount = ChartRedCount
	RChart.ChartBlackCount = ChartBlackCount

	//判断近20局输赢走势是否大于规定次数
	if Rlength > int(20) {
		RChart.RunChartTwenty = this.RunChart[:int(20)]
	} else {
		RChart.RunChartTwenty = this.RunChart
	}
	for i, _ := range RChart.RunChartTwenty {
		if RChart.RunChartTwenty[i] == RED {
			Rcount += 1
		} else {
			Bcount += 1
		}
	}

	//计算近20局红黑方的出现的次数
	RChart.Rcount = Rcount
	RChart.Bcount = Bcount
	//判断牌型走势长度是否大于规定次数
	if Clength > int(18) {
		this.CardTypeChart = this.CardTypeChart[:17]
		RChart.CardTypeChart = this.CardTypeChart
	} else {
		RChart.CardTypeChart = this.CardTypeChart
	}
	// logs.Debug("走势结构体：", RChart)
	p.SendNativeMsg(MSG_GAME_INFO_RUN_CHART_REPLY, RChart)
}

//玩家列表处理
func (this *ExtDesk) HandleGameUserList(p *ExtPlayer, d *DkInMsg) {
	this.RLock()
	defer this.RUnlock()

	userInfo := this.GetUserList(p)
	msg := GUserInfoReply{
		Id:        MSG_GAME_INFO_USER_LIST_REPLY,
		UserInfo:  userInfo,
		UserCount: this.UserCount,
	}
	p.SendNativeMsg(MSG_GAME_INFO_USER_LIST_REPLY, msg)
}

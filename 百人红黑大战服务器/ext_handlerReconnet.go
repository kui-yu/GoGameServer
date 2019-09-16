package main

// "logs"
// "logs"
// "time"

// "logs"

//处理房间信息
func (this *ExtDesk) HandleRoomInfo(p *ExtPlayer) {
	this.Lock()
	defer this.Unlock()

	result := GInfoReConnectReply{
		Id: MSG_GAME_INFO_RECONNECT_REPLY,

		GameCount: this.Count,
	}

	// 游戏状态
	result.GameState = int32(this.GameState)

	// 房号
	result.RoomId = this.RoomId
	// 局号
	result.GameId = this.GameId
	// 限红
	result.GameLimit = G_DbGetGameServerData.GameConfig.LimitRedMax
	// 下注金币限制
	result.BetList = G_DbGetGameServerData.GameConfig.TenChips
	// 可下注区域
	result.BetArea = p.PBetArea

	// 区域总下注
	result.TAreaCoins = this.GetAreaCoinsList()

	// 玩家下注
	p.ColAreaCoins()
	result.PAreaCoins = p.GetTotBetList()
	// 玩家金币
	result.PCoins = p.GetCoins()
	// 座位玩家
	result.SeatList = this.GetSeatInfo(p)
	// 打开的牌
	result.CardList = this.CardList
	//限制金币池ID
	var a int = -1
	for i := len(result.BetList) - 1; i >= 0; i-- {
		if p.Coins >= result.BetList[i] {
			a = i
			break
		}
	}
	result.LimitCoinId = int32(a)
	timerNum := 0
	// 当前状态时间
	switch this.GameState {
	case MSG_GAME_INFO_SHUFFLE_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Shuffle)
	case MSG_GAME_INFO_READY_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Ready)
	case MSG_GAME_INFO_SEND_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.SendCard)
	case MSG_GAME_INFO_BET_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Bet)
	case MSG_GAME_INFO_STOP_BET_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.StopBet)
	case MSG_GAME_INFO_OPEN_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Open)
	case MSG_GAME_INFO_AWARD_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Award)
	}
	result.Timer = int32(timerNum) * 1000

	// 红黑方牌
	if this.GameState >= MSG_GAME_INFO_OPEN_NOTIFY {
		result.RedCard = this.RedCard
		result.BlackCard = this.BlackCard
		result.WinArea = this.WinArea
		Rcard, Rcolor := SortHandCard(this.RedCard)   //将红方的牌进行排序并分出花色，数值
		Bcard, Bcolor := SortHandCard(this.BlackCard) //将黑方的牌进行排序并分出花色，数值
		//	判断牌的类型以及等级
		RGrade, _ := GetCardType(Rcard, Rcolor)
		BGrade, _ := GetCardType(Bcard, Bcolor)
		result.Btype = BGrade
		result.Rtype = RGrade
	}

	// 添加走势
	length := len(this.RunChart)
	Clength := len(this.CardTypeChart)
	if length > gameConfig.DeskInfo.RunChartCount {
		result.RunChart = this.RunChart[:gameConfig.DeskInfo.RunChartCount]
	} else {
		result.RunChart = this.RunChart
	}
	//判断牌型走势长度是否大于规定次数
	if Clength > gameConfig.DeskInfo.CardTypeChartCount {
		result.CardTypeChart = this.CardTypeChart[:gameConfig.DeskInfo.CardTypeChartCount]
	} else {
		result.CardTypeChart = this.CardTypeChart
	}
	p.LiXian = false
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT_REPLY, &result)
	// logs.Debug("断线重连信息:%v", result)
}

//处理断线重连
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
		Id:       MSG_GAME_RECONNECT_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})

	this.HandleRoomInfo(p)
}

//处理用户掉线
func (this *ExtDesk) HandleDisConnet(p *ExtPlayer, d *DkInMsg) {
	if p.GetTotAreaCoins() == 0 {
		this.SeatMgr.DelPlayer(p)
		this.LeaveByForce(p)
	} else {
		p.LiXian = true //方便踢出用户
	}
}

//处理用户退出房间
func (this *ExtDesk) HandleGameExit(p *ExtPlayer, d *DkInMsg) {
	//用户掉线处理
	if p.GetTotAreaCoins() == 0 {
		p.SendNativeMsg(MSG_GAME_INFO_EXIT_REPLY, GGameExitReply{
			Id:     MSG_GAME_INFO_EXIT_REPLY,
			Result: 0,
		})
		this.SeatMgr.DelPlayer(p)
		this.LeaveByForce(p)
	} else {
		p.SendNativeMsg(MSG_GAME_INFO_EXIT_REPLY, GGameExitReply{
			Id:     MSG_GAME_INFO_EXIT_REPLY,
			Result: 1,
		})
	}
}

/*用户离开*/
func (this *ExtDesk) Leave(p *ExtPlayer) bool {
	if p.GetTotAreaCoins() == 0 {
		this.SeatMgr.DelPlayer(p)
		this.LeaveByForce(p)
	} else {
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

//用户踢出房间
func (this *ExtDesk) HandleExit() {
	for _, v := range this.Players {
		limit := false
		if v.GetCoins() > int64(G_DbGetGameServerData.LimitHigh) {
			//如果用户的金币大于允许进入此房间级别的最大金额时，则限制玩家进入
			limit = true
			v.SendNativeMsg(MSG_GAME_INFO_EXIT_LIMIT_HIGHT, &GLeaveReply{
				Id: MSG_GAME_INFO_EXIT_LIMIT_HIGHT,
			})
		} /* else if v.GetCoins() <= int64(G_DbGetGameServerData.LimitLower) {
			//如果用户的金币小于允许进入此房间级别的最小金额时，则限制玩家进入
			limit = true
			v.SendNativeMsg(MSG_GAME_INFO_EXIT_LIMIT_LOW, &GLeaveReply{
				Id: MSG_GAME_INFO_EXIT_LIMIT_LOW,
			})
		}*/

		if v.LiXian || limit {
			//如果玩家确实离线或者处于上面两个条件中，则强制删除玩家
			this.SeatMgr.DelPlayer(v)
			this.LeaveByForce(v)
		}
	}
}

//用户长时间未下注，则发出警告和踢出
func (this *ExtDesk) HandleUndo() {
	var times int32
	for _, v := range this.Players {
		times = v.GetUndoTimes() //获取用户未下注的次数
		// logs.Debug("未下注次数！！！！！！", v.Nick, times)
		if times >= gameConfig.Undo.Exit {
			//	如果次数大于游戏强制用户退出的次数时，则会发出强制用户离开通知以及强制用户退出
			v.SendNativeMsg(MSG_GAME_INFO_EXIT_NOTIFY, &GLeaveReply{
				Id: MSG_GAME_INFO_EXIT_NOTIFY,
			})
			this.SeatMgr.DelPlayer(v)
			this.LeaveByForce(v)
			continue
		} else if times == gameConfig.Undo.Warning {
			// logs.Debug("达到警告次数！！！！！！！！！！！！！")
			//如果次数达到警告次数时，会发出强制用户离开通知
			v.SendNativeMsg(MSG_GAME_INFO_UNDO_NOTIFY, &GLeaveReply{
				Id: MSG_GAME_INFO_UNDO_NOTIFY,
			})
			// logs.Debug("警告已发出！！！！！！！！！！！！")
		}
	}
}

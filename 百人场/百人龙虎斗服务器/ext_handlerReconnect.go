package main

import (
	"logs"
)

func (this *ExtDesk) HandleRoomInfo(p *ExtPlayer) {
	this.Lock()
	defer this.Unlock()

	result := GInfoReConnectReply{
		Id:         MSG_GAME_INFO_RECONNECT_REPLY,
		LeftCount:  this.LeftCount,
		RightCount: this.RightCount,
		GameCount:  this.Count,
	}

	// 游戏状态
	result.GameState = int32(this.GameState)

	// 房号
	result.RoomId = this.RoomId
	// 局号
	result.GameId = this.GameId
	// 限红
	result.GameLimit = []int64{this.GameLimit.Low, this.GameLimit.High}
	// 下注金币限制
	result.BetList = this.BetList
	// 可下注区域
	result.BetArea = this.BetArea

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
	result.OpenCard = this.OpenCard

	timerNum := 0
	// 当前状态时间
	switch this.GameState {
	case MSG_GAME_INFO_SHUFFLE_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Shuffle) * 1000
	case MSG_GAME_INFO_READY_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Ready) * 1000
	case MSG_GAME_INFO_SEND_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.SendCard) * 1000
	case MSG_GAME_INFO_BET_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Bet) * 1000
	case MSG_GAME_INFO_STOP_BET_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.StopBet) * 1000
	case MSG_GAME_INFO_OPEN_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Open) * 1000
	case MSG_GAME_INFO_AWARD_NOTIFY:
		timerNum = this.GetTimerNum(gameConfig.Timer.Award) * 1000
	}
	result.Timer = int32(timerNum)

	// 龙、虎牌
	if this.GameState >= MSG_GAME_INFO_OPEN_NOTIFY {
		result.DragonCard = this.DragonCard
		result.TigerCard = this.TigerCard
		result.WinArea = this.WinArea
	}

	// 添加走势
	length := len(this.RunChart)
	if length > gameConfig.DeskInfo.RunChartCount {
		result.RunChart = this.RunChart[length-gameConfig.DeskInfo.RunChartCount:]
	} else {
		result.RunChart = this.RunChart
	}

	p.LiXian = false
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT_REPLY, &result)
}

func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d2 *DkInMsg) {
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

// 用户退出房间
func (this *ExtDesk) HandleGameExit(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("接收到用户退出房间:", p.Nick, p.Robot)
	// 用户掉线处理
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

func (this *ExtDesk) Leave(p *ExtPlayer) bool {
	// 用户掉线处理
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

// 用户掉线，处理与退出房间一致
func (this *ExtDesk) HandleDisConnect(p *ExtPlayer, d *DkInMsg) {
	// 用户掉线处理
	if p.GetTotAreaCoins() == 0 {
		this.SeatMgr.DelPlayer(p)
		this.LeaveByForce(p)
	} else {
		p.LiXian = true // 方便结算剔除用户
	}
}

// 用户踢出房间
func (this *ExtDesk) HandleExit() {
	for _, v := range this.Players {
		limit := false
		logs.Debug("用户金币", v.GetCoins())
		if v.GetCoins() >= int64(G_DbGetGameServerData.LimitHigh) {

			limit = true
			v.SendNativeMsg(MSG_GAME_INFO_EXIT_LIMIT_HIGHT, &GLeaveReply{
				Id: MSG_GAME_INFO_EXIT_LIMIT_HIGHT,
			})
		}
		// else if v.GetCoins() < int64(G_DbGetGameServerData.Restrict) {
		// 	limit = true
		// 	v.SendNativeMsg(MSG_GAME_INFO_EXIT_LIMIT_LOW, &GLeaveReply{
		// 		Id: MSG_GAME_INFO_EXIT_LIMIT_LOW,
		// 	})
		// }
		if v.LiXian || limit {
			this.SeatMgr.DelPlayer(v)
			this.LeaveByForce(v)
		}
	}
}

// 用户踢出和踢出警告
func (this *ExtDesk) HandleUndo() {
	var times int32
	for _, v := range this.Players {
		times = v.GetUndoTimes()
		if times >= gameConfig.Undo.Exit {
			logs.Debug("发现用户没有下注::", v.Nick, v.Robot, v.GetUndoTimes())
			v.SendNativeMsg(MSG_GAME_INFO_EXIT_NOTIFY, &GLeaveReply{
				Id: MSG_GAME_INFO_EXIT_NOTIFY,
			})
			this.SeatMgr.DelPlayer(v)
			this.LeaveByForce(v)
			continue
		} else if times == gameConfig.Undo.Warning {
			logs.Debug("警告用户没有下注::", v.Nick, v.Robot, v.GetUndoTimes())
			v.SendNativeMsg(MSG_GAME_INFO_UNDO_NOTIFY, &GLeaveReply{
				Id: MSG_GAME_INFO_UNDO_NOTIFY,
			})
		}
	}
}

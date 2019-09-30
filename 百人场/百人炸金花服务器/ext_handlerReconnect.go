package main

import (
	"fmt"
)

func (this *ExtDesk) HandleRoomInfo(p *ExtPlayer) {
	this.Lock()
	defer this.Unlock()

	result := GInfoReConnectReply{
		Id: MSG_GAME_INFO_RECONNECT_REPLY,
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

	// 区域总下注
	result.TAreaCoins = this.GetAreaCoinsList()

	// 玩家下注
	p.ColAreaCoins()
	result.PAreaCoins = p.GetTotBetList()
	// 玩家金币
	result.PCoins = p.GetCoins()
	// 座位玩家
	result.SeatList = this.GetSeatInfo(p)
	// 骰子
	result.Dices = []int{this.dices1, this.dices2}

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

	// 庄、闲牌
	if this.GameState >= MSG_GAME_INFO_OPEN_NOTIFY {
		fmt.Println("if1")
		result.BankerCard = this.BankerCard
		result.IdleCard = this.IdleCard
		result.TypeList = this.TypeList
		result.MBankerCard = this.MBankerCard
		result.MIdleCard = this.MIdleCard
		result.WinArea = this.WinArea
		result.TBankerCard = this.TBankerType
	} else if this.GameState >= MSG_GAME_INFO_SEND_NOTIFY {
		fmt.Println("if2")
		// 下注阶段，庄、闲开两张牌
		result.BankerCard = this.BankerCard[0:2]
		for i := 0; i < 4; i++ {
			result.IdleCard = append(result.IdleCard, this.IdleCard[i][0:2])
		}
	}

	result.BankRunChart = this.BankRunChart

	p.LiXian = false
	fmt.Printf("庄牌:%v,闲牌:%v,阶段:%v\n", result.BankerCard, result.IdleCard, result.GameState)
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
		Result:   0,
		CostType: GetCostType(),
	})

	this.HandleRoomInfo(p)
}

// 用户退出房间
func (this *ExtDesk) HandleGameExit(p *ExtPlayer, d *DkInMsg) {
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
		if v.GetCoins() >= int64(G_DbGetGameServerData.LimitHigh) && GetCostType() == 1 { //如果是体验场才判断金币过高踢出
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
			v.SendNativeMsg(MSG_GAME_INFO_EXIT_NOTIFY, &GLeaveReply{
				Id: MSG_GAME_INFO_EXIT_NOTIFY,
			})
			this.SeatMgr.DelPlayer(v)
			this.LeaveByForce(v)
			continue
		} else if times == gameConfig.Undo.Warning {
			v.SendNativeMsg(MSG_GAME_INFO_UNDO_NOTIFY, &GLeaveReply{
				Id: MSG_GAME_INFO_UNDO_NOTIFY,
			})
		}
	}
}

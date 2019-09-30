package main

import (
	"sync"
)

type ExtPlayer struct {
	Player

	//以下是用户自己定义的变量
	sync.RWMutex // 锁、防止多线程读写

	Online       bool            // 是否在线
	DownBets     map[uint8]int64 // 下注的金额集合
	BalaDownBets map[uint8]int64 // 结算的金额集合
	UnbetsCount  int32           // 未参与下注局数
	BetHistorys  []BetHistory    // 输赢统计

	HVictoryCount int32 // 历史输赢的局数
	HDownBetTotal int64 // 历史下注的统计
}

// 初始化用户信息(用户刚进入时调用)
func (this *ExtPlayer) Init() {
	this.Online = false
	this.UnbetsCount = 0
	this.DownBets = make(map[uint8]int64)
	this.BalaDownBets = make(map[uint8]int64)
	this.HVictoryCount = 0
	this.HDownBetTotal = 0
}

// 获得玩家的下注总额
func (this *ExtPlayer) getDownBetTotal() int64 {
	var total int64 = 0
	for _, value := range this.DownBets {
		total += value
	}

	return total
}

// 添加到下注历史
func (this *ExtPlayer) addBetHistory(isVictory bool, downBet int64) {
	this.BetHistorys = append(this.BetHistorys, BetHistory{
		IsVictory: isVictory,
		DownBet:   downBet,
	})

	this.HDownBetTotal += downBet
	if isVictory {
		this.HVictoryCount += 1
	}

	userListRecordCount := gameConfig.GameLimtInfo.UserListRecordCount

	if len(this.BetHistorys) > userListRecordCount {
		this.HDownBetTotal -= this.BetHistorys[0].DownBet
		if this.BetHistorys[0].IsVictory {
			this.HVictoryCount -= 1
		}

		this.BetHistorys = this.BetHistorys[1:]
	}
}

package main

import (
	"sync"
)

type ExtPlayer struct {
	Player

	//以下是用户自己定义的变量
	sync.RWMutex // 锁、防止多线程读写

	Online       bool          // 是否在线
	DownBets     map[int]int64 // 下注的金额集合
	PAreaCoins   []int64       //自己区域总下注
	BalaDownBets map[int]int64 // 结算的金额集合
	WinCoins     int64         // 结算金币
	UnbetsCount  int32         // 未参与下注局数
	BetHistorys  []int         // 输赢统计   0,win 1,lose
	// HGameRecord  []GPlayerRocord //玩家游戏记录
	TotalBet     int64   //用户总下注
	ElseWinAndOr []int64 //其他玩家输赢
	// OtherBet     []int64 //其他玩家下注
	// ChairWinOrLost []int64 //椅子玩家输赢
	Match int //局数
}

// 初始化用户信息(用户刚进入时调用)
func (this *ExtPlayer) Init() {
	this.UnbetsCount = 0
	this.DownBets = make(map[int]int64)

	this.BalaDownBets = make(map[int]int64)
	for i := 0; i < 8; i++ {
		this.BalaDownBets[i] = 0
	}
	this.WinCoins = 0
	this.PAreaCoins = []int64{0, 0, 0, 0, 0, 0, 0, 0}
	this.ElseWinAndOr = []int64{0, 0, 0, 0, 0, 0, 0, 0}
	// this.OtherBet = []int64{0, 0, 0, 0, 0, 0, 0, 0}
	// this.ChairWinOrLost = []int64{0, 0, 0, 0, 0, 0}
}

func (this *ExtPlayer) ResetExtPlayer() {
	this.DownBets = make(map[int]int64)
	for i := 0; i < gameConfig.LimitInfo.BetCount; i++ {
		this.DownBets[i] = 0
	}
	this.BalaDownBets = make(map[int]int64)
	this.WinCoins = 0
	this.ElseWinAndOr = []int64{0, 0, 0, 0, 0, 0, 0, 0}
	// this.ChairWinOrLost = []int64{0, 0, 0, 0, 0, 0}
}

// 获得玩家的下注总额
func (this *ExtPlayer) getDownBetTotal() int64 {
	var total int64 = 0
	for _, value := range this.DownBets {
		total += value
	}
	return total
}

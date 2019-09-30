package main

import (
	"sync"

	"bl.com/util"
)

type ExtPlayer struct {
	Player

	//以下是用户自己定义的变量
	sync.RWMutex // 锁、防止多线程读写

	msgId            int32         // 下注消息号，防止重复处理下注
	arealistCoins    util.AreaList // 总下注情况
	newAreaListCoins util.AreaList // 新下注情况

	winCoins  int64 // 赢取金币
	undoTimes int32 // 未参与下注次数

	WList   []bool
	BetList []int64

	WinCount int64
	TotBet   int64

	WinList [8]int64
	IsBet   bool
}

func (this *ExtPlayer) GetUid() int64 {
	return this.Uid
}

// 初始化用户信息(用户刚进入时调用)
func (this *ExtPlayer) Init() {
	this.msgId = 0
	this.arealistCoins.Init(8)
	this.newAreaListCoins.Init(8)

	this.undoTimes = 0
	this.WinList = [8]int64{}
}

func (this *ExtPlayer) GetMsgId() int32 {
	return this.msgId
}

func (this *ExtPlayer) SetMsgId(id int32) {
	this.msgId = id
}

func (this *ExtPlayer) AddCoins(coin int64) {
	this.Lock()
	defer this.Unlock()

	this.Coins += coin
}

func (this *ExtPlayer) GetCoins() int64 {
	this.RLock()
	defer this.RUnlock()

	return this.Coins
}

func (this *ExtPlayer) AddUndoTimes() {
	this.RLock()
	defer this.RUnlock()

	this.undoTimes++
}

func (this *ExtPlayer) GetUndoTimes() int32 {
	this.RLock()
	defer this.RUnlock()

	return this.undoTimes
}

// 清空下注
func (this *ExtPlayer) ResetAreaList() {
	this.msgId = 0
	this.arealistCoins.Init(8)
	this.newAreaListCoins.Init(8)

	this.winCoins = 0
}

// 获取总下注
func (this *ExtPlayer) GetTotAreaCoins() int64 {
	this.RLock()
	defer this.RUnlock()

	areaCoins := this.arealistCoins.GetTotValue()

	return areaCoins
}

// 获取下注列表
func (this *ExtPlayer) GetNTAreaCoinsList() []int64 {
	ret := this.arealistCoins.GetValueList()
	newList := this.newAreaListCoins.GetValueList()

	for i, v := range newList {
		ret[i] += v
	}

	return ret
}
func (this *ExtPlayer) GetTotBetList() []int64 {
	ret := this.arealistCoins.GetValueList()
	return ret
}

func (this *ExtPlayer) GetNewBetList() []int64 {
	ret := this.newAreaListCoins.GetValueList()
	return ret
}

// 获取下注金币
func (this *ExtPlayer) GetTotAreaCoin(area int) int64 {
	ret := this.arealistCoins.GetValue(area)
	return ret
}
func (this *ExtPlayer) GetNewAreaCoin(area int) int64 {
	ret := this.newAreaListCoins.GetValue(area)
	return ret
}

// 添加下注
func (this *ExtPlayer) AddTotAreaCoins(area int, coins int64) bool {
	ret := this.arealistCoins.AddValue(area, coins)

	return ret
}
func (this *ExtPlayer) AddNewAreaCoins(area int, coins int64) bool {
	ret := this.newAreaListCoins.AddValue(area, coins)
	if ret {
		this.Coins -= coins
		this.undoTimes = 0
	}
	return ret
}

// 新下注添加到总下注
func (this *ExtPlayer) ColAreaCoins() {
	this.Lock()

	length := this.newAreaListCoins.GetLength()
	for i := 0; i < length; i++ {
		this.AddTotAreaCoins(i, this.newAreaListCoins.GetValue(i))
		this.newAreaListCoins.SetValue(i, 0)
	}

	this.Unlock()
}

// 添加赢取金币
func (this *ExtPlayer) AddWin(win int64) {
	this.Lock()

	this.winCoins += win

	this.Unlock()
}

// 获取赢取金币
func (this *ExtPlayer) GetWinCoins() int64 {
	this.Lock()
	defer this.Unlock()

	return this.winCoins
}

// 构建赢取列表
func (this *ExtPlayer) BuildWinList(double []float64) {
	this.Lock()
	defer this.Unlock()

	for i, v := range double {
		coins := this.GetTotAreaCoin(i)

		coins = int64(float64(coins) * v)

		this.winCoins += coins
		this.WinList[i] = coins
		if coins == 0 {
			this.WinList[i] = 0 - this.GetTotAreaCoin(i)
		}
	}
}

// 获取赢取列表
func (this *ExtPlayer) GetWinList() []int64 {
	this.Lock()
	defer this.Unlock()

	return this.WinList[:]
}

// 结算
func (this *ExtPlayer) Award() {
	this.Lock()

	this.Coins += this.winCoins
	this.winCoins = 0

	this.Unlock()
}

// 玩家记录
func (this *ExtPlayer) AddWinList() {
	// 保留最近20条输赢记录
	if len(this.WList) == gameConfig.DeskInfo.ListCount {
		this.WList = this.WList[1:]
	}

	this.WList = append(this.WList, this.winCoins > 0)
}

func (this *ExtPlayer) GetWinCount() int32 {
	var count int32 = 0
	for _, v := range this.WList {
		if v {
			count++
		}
	}

	return count
}

func (this *ExtPlayer) AddBetList() {
	bet := this.arealistCoins.GetTotValue()

	// 保留最近20条下注记录
	if len(this.BetList) == gameConfig.DeskInfo.ListCount {
		this.BetList = this.BetList[1:]
	}

	this.BetList = append(this.BetList, bet)
}

func (this *ExtPlayer) GetBetCoins() int64 {
	var totBet int64 = 0
	for _, v := range this.BetList {
		totBet += v
	}

	return totBet
}

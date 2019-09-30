package main

type ExtPlayer struct {
	Player
	GameBoundCount int32 // Bound 出现的次数
	GameWin        int32 // 次元门获取
}

func (this *ExtPlayer) ResetExtPlayer() {
	this.GameBoundCount = 0
	this.GameWin = 0
}

// 添加Bound
func (this *ExtPlayer) AddGameBound(count int32) int32 {
	this.GameBoundCount = this.GameBoundCount + count

	return this.GameBoundCount
}

// 重置Bound次数
func (this *ExtPlayer) ResetGameBound() {
	this.GameBoundCount = 0
}

// 获取Bound次数
func (this *ExtPlayer) GetGameBound() int32 {
	return this.GameBoundCount
}

// 添加次元门奖励
func (this *ExtPlayer) AddGameWin(win int32) int32 {
	this.GameWin = this.GameWin + win
	return this.GameWin
}

// 重置次元门奖励
func (this *ExtPlayer) ResetGameWin() {
	this.GameWin = 0
}

// 次元门奖励翻倍
func (this *ExtPlayer) GameWinDouble(double int32) int32 {
	this.GameWin = this.GameWin * double

	return this.GameWin
}

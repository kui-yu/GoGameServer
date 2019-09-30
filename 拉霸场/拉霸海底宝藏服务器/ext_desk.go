package main

import (
	"fmt"
)

type ExtDesk struct {
	Desk

	imgMgr       MgrImg // 图标管理器
	BonusCount   int64  // 累计Bonus图标数
	BoxCount     int64  // 奖金游戏宝箱数
	MinBonus     int64  // 一局出现Bonus图标数 触发奖金游戏
	IsBonusGame  bool   // 是否已触发奖金游戏
	IsBonusStart bool   // 是否已经开启奖金游戏

	NormalCoins  []int64 // 宝箱内容
	NormalWeight []int32 // 内容权值

	Lines     int64
	Times     int64
	TotLines  int64
	LineCount int //记录总线束

	GameId string

	BoxCoins []int64 // 宝箱金币
}

//重置桌子
func (this *ExtDesk) ResetTable() {
	this.BonusCount = 0
	this.BoxCount = int64(len(gameConfig.DeskInfo.DimensionDoor[0].Coins))
	this.MinBonus = int64(gameConfig.DeskInfo.OpenDimension)
	this.Lines = 0
	this.Times = 0
	this.TotLines = 0
	this.IsBonusGame = false
	this.IsBonusStart = false
	this.imgMgr.Init()
}

func (this *ExtDesk) InitExtData() {
	//
	this.Handle = make(map[int32]func(*ExtPlayer, *DkInMsg))
	this.BonusCount = 0
	this.BoxCount = int64(len(gameConfig.DeskInfo.DimensionDoor[0].Coins))
	this.MinBonus = int64(gameConfig.DeskInfo.OpenDimension)
	this.IsBonusGame = false
	this.IsBonusStart = false
	this.imgMgr.Init()

	// 游戏匹配成功
	this.Handle[MSG_GAME_AUTO] = this.HandleGameAuto
	// 游戏开始
	this.Handle[MSG_GAME_INFO_PLAY] = this.HandleGameInfoPlay
	// 奖金游戏开始
	this.Handle[MSG_GAME_INFO_BONUS_START] = this.BonusStart
	// 奖金游戏
	this.Handle[MSG_GAME_INFO_BONUS] = this.HandleGameBonus
	// 游戏退出
	this.Handle[MSG_GAME_INFO_LEAVE] = this.HandleLeave
	//断线消息
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleLeave
	//断线重连 400010
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
}

func (this *ExtDesk) InitBonus() {
	// 累计图标清空
	this.BonusCount = 0

	this.BoxCoins = []int64{}

	this.NormalCoins = []int64{}
	if this.Lines > int64(len(gameConfig.DeskInfo.DimensionDoor)) {
		this.Lines = int64(len(gameConfig.DeskInfo.DimensionDoor))
	} else if this.Lines < 1 {
		this.Lines = 1
	}
	result := this.LineCount / gameConfig.DeskInfo.OpenDimension
	var index int
	if result > 9 {
		index = 9
	} else {
		index = result
	}
	for _, v := range gameConfig.DeskInfo.DimensionDoor[index-1].Coins {
		this.NormalCoins = append(this.NormalCoins, int64(v))
		this.BoxCoins = append(this.BoxCoins, 0)
	}
	fmt.Println(gameConfig.DeskInfo.DimensionDoor[index-1].Coins)
	this.NormalWeight = []int32{}
	for _, v := range gameConfig.DeskInfo.DimensionWeight {
		this.NormalWeight = append(this.NormalWeight, int32(v))
	}

	this.IsBonusGame = true
}

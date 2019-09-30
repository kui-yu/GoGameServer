package main

//import (
//	"sort"
//)

type ExtPlayer struct {
	Player
	//以下是用户自己定义的变量
	HandCard    []int32 //手牌
	NiuPoint    int32   //牛点
	NiuCards    []int32 //牛牌
	NiuMultiple int     //牛点倍数
	BetMultiple int     //倍数
	IsLook      bool    //是否看牌
	WinMultiple int     //输赢总倍数
	WinCoins    int64   //金币数
	WinPlayer   []int32 //赢的玩家
	RateCoin    float64 //抽水金币
	PlayerBets  []int
}

type GCardType struct {
	HandCard []int32 //手牌
	NiuPoint int32   //牛点
	NiuCards []int32 //牛牌
}

//重置玩家信息
func (this *ExtDesk) ResetPlayer(p *ExtPlayer) {
	p.HandCard = []int32{}
	p.NiuPoint = 0
	p.NiuCards = []int32{}
	p.NiuMultiple = 0
	p.BetMultiple = 0
	p.IsLook = false
	p.WinMultiple = 0
	p.WinCoins = 0
	p.WinPlayer = []int32{}
	p.RateCoin = 0
	p.PlayerBets = []int{}
}

//注册玩家消息
func (this *ExtDesk) InitExtData() {

	//内容初始化
	this.InitGame()

	//玩家匹配 400001
	this.Handle[MSG_GAME_AUTO] = this.HandleGameAuto
	//玩家叫分 410001
	this.Handle[MSG_GAME_INFO_CALL] = this.HandleGameCall
	//玩家开牌 410005
	this.Handle[MSG_GAME_INFO_PLAY] = this.HandlePlayCard
	//断线消息
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
	//重连
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
}

//广播阶段
func (this *ExtPlayer) BroadStageTime(gameState int32, time int32) {
	stage := GStageInfo{
		Id:        MSG_GAME_INFO_STAGE,
		Stage:     gameState,
		StageTime: time,
	}
	this.SendNativeMsg(MSG_GAME_INFO_STAGE, &stage)
}

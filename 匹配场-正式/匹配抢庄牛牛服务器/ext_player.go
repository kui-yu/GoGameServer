package main

//import (
//	"sort"
//)

type ExtPlayer struct {
	Player
	//以下是用户自己定义的变量
	HandCard     []int32 //手牌
	NiuPoint     int32   //牛点
	NiuCards     []int32 //牛牌
	NiuMultiple  int     //牛点倍数
	BetMultiple  int     //倍数
	IsLook       bool    //是否看牌
	WinMultiple  int     //输赢总倍数
	WinCoins     int64   //金币数
	WinPlayer    []int32 //赢的玩家
	CallMultiple int     //叫庄
	CallBankFlag bool    //是否叫庄
	RateCoins    float64 //手续费
	PlayerBets   []int
	PlayerCalls  []int
}

type GCardType struct {
	HandCard []int32 //手牌
	NiuPoint int32   //牛点
	NiuCards []int32 //牛牌
}

//注册玩家消息
func (this *ExtDesk) InitExtData() {

	//内容初始化
	this.InitGame()

	//玩家匹配 400001
	this.Handle[MSG_GAME_AUTO] = this.HandleGameAuto
	//玩家叫分 410001
	this.Handle[MSG_GAME_INFO_CALL] = this.HandleGameBet
	//玩家开牌 410005
	this.Handle[MSG_GAME_INFO_PLAY] = this.HandlePlayCard
	//玩家抢庄 410007
	this.Handle[MSG_GAME_INFO_CALL_BANKER] = this.HandleGameCall
	//断线消息 400013断线
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
	//重连 400010断线重连
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

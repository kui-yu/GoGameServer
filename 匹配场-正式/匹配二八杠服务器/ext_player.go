package main

//import (
//	"sort"
//)

type ExtPlayer struct {
	Player
	//以下是用户自己定义的变量
	HandCards    []int          //初始手牌
	WinMultiple  int            //赢得倍数
	WinCoins     int64          //输赢总得分
	WinList      []int32        //赢的玩家
	CallMultiple int            //叫庄倍数
	PlayMultiple int            //下注倍数
	RecordInfos  []GSRecordInfo //游戏记录
	RateCoins    float64        //手续费
	PlayerBets   []int
	PlayerCalls  []int
}

//注册玩家消息
func (this *ExtDesk) InitExtData() {

	//初始化
	this.InitGame()
	//玩家匹配 400001
	this.Handle[MSG_GAME_AUTO] = this.HandleGameAuto
	//玩家叫庄 410004
	this.Handle[MSG_GAME_INFO_CALL_INFO] = this.HandleGameCall
	//玩家下注 410005
	this.Handle[MSG_GAME_INFO_PLAY_INFO] = this.HandleGameBet
	//断线重连 400010
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
	//断线消息 400013断线
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
	//个人战绩查询
	this.Handle[MSG_GAME_INFO_RECORD_INFO] = this.HandleRecord
}

//重置玩家信息
func (this *ExtDesk) ResetPlayer(p *ExtPlayer) {
	p.HandCards = []int{}
	p.WinMultiple = 0
	p.WinCoins = 0
	p.WinList = []int32{}
	p.CallMultiple = -1
	p.PlayMultiple = -1
	this.Banker = -1
	p.RateCoins = 0
	p.PlayerBets = []int{}
	p.PlayerCalls = []int{}
}

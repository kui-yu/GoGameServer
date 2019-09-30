package main

//import (
//	"sort"
//)

type ExtPlayer struct {
	Player
	//以下是用户自己定义的变量
	IsReady      int            //是否准备
	HandCard     []int          //手牌
	NiuPoint     int            //牛点
	NiuCards     []int          //牛牌
	NiuMultiple  int            //牛点倍数
	BetMultiple  int            //倍数
	IsLook       bool           //是否看牌
	WinMultiple  int            //输赢总倍数
	WinCoins     int64          //金币数
	WinPlayer    []int32        //赢的玩家
	CallMultiple int            //叫庄
	RecordInfos  []GSRecordInfo //游戏记录
	RateCoins    float64        //手续费用
	TotalCoins   int64          //总金币

	PlayerBets  []int
	PlayerCalls []int
}

//注册玩家消息
func (this *ExtDesk) InitExtData() {

	//牌内容初始化
	this.CardMgr.InitCards()
	this.CardMgr.InitNormalCards()

	//玩家匹配 400019
	this.Handle[MSG_GAME_FK_JOIN] = this.HandleGameAuto
	//玩家叫分 410001
	this.Handle[MSG_GAME_INFO_CALL] = this.HandleGameCall
	//玩家开牌 410005
	this.Handle[MSG_GAME_INFO_PLAY] = this.HandlePlayCard
	//玩家抢庄 410007
	this.Handle[MSG_GAME_INFO_CALL_BANKER] = this.HandleCallBank
	//断线消息 400013断线
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
	//重连 400010断线重连
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
	//玩家准备
	this.Handle[MSG_GAME_INFO_READY] = this.HandleReady
	//玩家离开
	this.Handle[MSG_GAME_INFO_LEAVE] = this.HandleLeave
	//玩家解散
	this.Handle[MSG_GAME_INFO_DISMISS] = this.HandleDisMiss
	//个人战绩查询
	this.Handle[MSG_GAME_INFO_RECORD_INFO] = this.HandleRecord
}

//重置玩家信息
func (this *ExtDesk) ResetPlayer(p *ExtPlayer) {
	p.IsReady = 0
	p.HandCard = []int{}
	p.NiuPoint = -1
	p.NiuCards = []int{}
	p.BetMultiple = 0
	p.IsLook = false
	p.WinCoins = 0
	p.WinPlayer = []int32{}
	p.CallMultiple = -1
	p.PlayerBets = []int{}
	p.PlayerCalls = []int{}
	p.NiuMultiple = 0
	p.WinMultiple = 0
	p.RateCoins = 0
}

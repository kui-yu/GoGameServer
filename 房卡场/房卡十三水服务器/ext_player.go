package main

//import (
//	"sort"
//)

type ExtPlayer struct {
	Player
	//以下是用户自己定义的变量
	HandCards    []int          //初始手牌
	SpecialType  int            //特殊牌型值
	SpecialCards []int          //特殊牌型
	SpecialScore int            //特殊得分
	PlayTypes    []int          //玩家牌型值集合[头墩，中墩，底墩]
	PlayCards    []int          //玩家摆牌
	IsPlay       int            //是否摆完牌 0未摆牌 1摆普通牌 2摆特殊牌
	WinCoins     int64          //输赢总得分
	WinCoinList  []int          //输赢列表[头墩，中墩，底墩，总得分]
	NormalScores []int          //普通得分
	ShootPlayers []int32        //打枪
	ShootScoress [][]int        //打枪分数{[头墩得分，中墩得分，底分得分]，[头墩得分，中墩得分，底分得分]}
	IsAllWin     bool           //全垒打
	RateCoins    float64        //手续费用
	IsReady      int            //是否准备 0.未准备；1.准备
	RecordInfos  []GSRecordInfo //游戏记录
	TotalCoins   int64          //总金币
}

//注册玩家消息
func (this *ExtDesk) InitExtData() {
	//玩家匹配 400019
	this.Handle[MSG_GAME_FK_JOIN] = this.HandleGameAuto
	//断线重连 400010
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
	//玩家摆牌
	this.Handle[MSG_GAME_INFO_PLAY] = this.HandlePlay
	//断线消息 400013断线
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
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
	p.HandCards = []int{}             //初始手牌
	p.SpecialType = 0                 //特殊牌型值
	p.SpecialCards = []int{}          //特殊牌型
	p.SpecialScore = 0                //特殊牌型得分
	p.PlayTypes = []int{0, 0, 0}      //玩家牌型值集合[头墩，中墩，底墩]
	p.PlayCards = []int{}             //玩家摆牌
	p.IsPlay = 0                      //是否摆完牌 0未摆牌 1摆普通牌 2摆特殊牌
	p.WinCoins = 0                    //输赢总得分
	p.WinCoinList = []int{0, 0, 0, 0} //输赢列表[头墩，中墩，底墩，总得分]
	p.NormalScores = []int{0, 0, 0}   //普通得分
	p.ShootPlayers = []int32{}        //打枪
	p.ShootScoress = [][]int{}        //打枪分数{[头墩得分，中墩得分，底分得分]，[头墩得分，中墩得分，底分得分]}
	p.IsAllWin = false                //全垒打
	p.IsReady = 0
	p.RateCoins = 0
}

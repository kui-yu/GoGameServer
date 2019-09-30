package main

import "time"

type ExtPlayer struct {
	Player
	IsDouble   bool    //是否翻倍
	IsBet      bool    //是否有下注
	PlaceBet   []int64 //区域下注，0到3分别是黑红梅方
	AllBet     int64   //总下注总额
	BetArrAble int     //可下注筹码
	NotBet     int     //累积没下注的局数
	//RoundSettleInfo []RoundSettleInfo //每局输赢信息记录
	SecondBet []int64 //每秒群发玩家下注
	//玩家列表字段
	Round              int   //累积下注局数
	AccumulateBet      int64 //累积下注
	AccumulateCoins    int64 //累积输赢
	AccumulateCoinsArr []int //前20局输赢
}

//初始化玩家
func (this *ExtPlayer) init() {
	this.IsBet = false
	this.PlaceBet = make([]int64, 4)
	this.AllBet = 0
	this.GetBetArr(false, 0)
	this.SecondBet = []int64{-1, -1, -1, -1} //初始化群发下注，-1是没有下注
}

//金币不足变回平倍
func (this *ExtPlayer) IsChangeDouble() bool {
	if this.BetArrAble == -1 {
		this.GetBetArr(false, 0)
		return true
	} else {
		return false
	}
}

//添加输赢记录，保持最新的10条
func (this *ExtPlayer) AddRoundSettle(data RoundSettleInfo) {
	//瀑布流
	var info = &RecordData{
		GradeType: data.GradeType,
		AllBet:    data.AllBet,
		WinCoins:  data.WinCoins,
		BetArea:   data.BetArea,
		CardType:  data.CardType,
		Time:      time.Now().Format("2006-01-02 15:04:05"),
		Date:      time.Now().Unix(),
	}
	G_AllRecord.AddRecord(this.Uid, info) //添加
	//如果没有记录直接添加，返回
	// if len(this.RoundSettleInfo) == 0 {
	// 	this.RoundSettleInfo = append(this.RoundSettleInfo, data)
	// 	return
	// }
	// //有记录判断是否大于10条
	// if len(this.RoundSettleInfo) < 10 {
	// 	this.RoundSettleInfo = append([]RoundSettleInfo{data}, this.RoundSettleInfo...)
	// } else {
	// 	this.RoundSettleInfo = append([]RoundSettleInfo{data}, this.RoundSettleInfo[0:9]...)
	// }
}

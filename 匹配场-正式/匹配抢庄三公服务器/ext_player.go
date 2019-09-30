package main

type ExtPlayer struct {
	Player
	HandCards   Card  //手牌
	IsOpenCards bool  //是否亮牌了
	WinCoins    int64 //输赢的钱
	//GameTrend   []Trend    //玩家庄家输赢走势
	BankerInfos BankerInfo //抢庄信息
}

// func (this *ExtPlayer) AddTrend(t Trend) {
// 	if len(this.GameTrend) < 15 {
// 		this.GameTrend = append(this.GameTrend, t)
// 	} else {
// 		this.GameTrend = append(this.GameTrend[1:], t)
// 	}
// }

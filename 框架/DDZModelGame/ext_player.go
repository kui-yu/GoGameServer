package main

//import (
//	"sort"
//)

// //初始化展示：其他都是具体项目具体定义
// type ExtPlayer struct {
// 	Player
// }

type ExtPlayer struct {
	Player
	//以下是用户自己定义的变量
	HandCard  []byte // 手牌
	CFen      int32  //叫分
	Pass      bool   //是否过牌
	TuoGuan   bool   //是否托管
	Outed     []*GOutCard
	Double    int32
	RateCoins float64 //手续费
}

// 设置手牌,手牌由大到小排序
func (this *ExtPlayer) SetHandCard(c []byte) {
	this.HandCard = c
}

package main

//import (
//	"sort"
//)

// //初始化展示：其他都是具体项目具体定义
// type ExtPlayer struct {
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
	isReady   int            //是否准备 0，未准备  1，准备
	GetMSG    int32          //是否抢（叫） 过地主 0,初始值 没做任何操作 1，叫地主 2，抢地主，3，不叫，4，不抢
	GGameLogs GGameLogReplay //玩家游戏记录
	LogCoins  int64          //玩家输赢总和
	IsDimiss  int32          //是否同意解散房间  -1，表示未操作    0， 不同意   1,同意
	JiFen     int64
}

// 设置手牌,手牌由大到小排序
func (this *ExtPlayer) SetHandCard(c []byte) {
	this.HandCard = c
	this.Double = 1
	this.IsDimiss = -1
}

func (this *ExtPlayer) ResetPlayer() {
	this.isReady = 0
	this.HandCard = []byte{}
	this.CFen = 0
	this.Pass = false
	this.TuoGuan = false
	this.Outed = []*GOutCard{}
	this.Double = 1
	this.GetMSG = 0
	this.IsDimiss = -1
}

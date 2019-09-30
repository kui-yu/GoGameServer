package main

import (
	"logs"
	"sync"
)

type ExtPlayer struct {
	sync.RWMutex // 锁、防止多线程读写
	Player
	DownBet            []int64      //玩家下注区域集合
	NoToBet            int          //未参与下注数
	OtherDownBet       []int64      //其他玩家下注
	BalanceResult      []int64      //玩家各个区域结算
	OtherBalanceResult []int64      //其他玩家各个区域结算
	History            []BetHistory //历史输赢
	IsOnChair          bool         //玩家是否在座位上
	OldtherDownBet     []int64      //旧的其他玩家下注
	GetCoins           int64        //玩家输赢
}

func (this *ExtPlayer) InitExtData() {
	for i := 0; i < gameConfig.LimitInfo.BetCount; i++ {
		this.DownBet = append(this.DownBet, 0)
		this.BalanceResult = append(this.BalanceResult, 0)
		this.OtherDownBet = append(this.OtherDownBet, 0)
		this.OtherBalanceResult = append(this.OtherBalanceResult, 0)
		this.OldtherDownBet = append(this.OldtherDownBet, 0)
	}
}

func (this *ExtPlayer) Rest() {
	for i := 0; i < gameConfig.LimitInfo.BetCount; i++ {
		this.DownBet[i] = 0
		this.BalanceResult[i] = 0
		this.OtherDownBet[i] = 0
		this.OtherBalanceResult[i] = 0
		this.OldtherDownBet[i] = 0
		this.GetCoins = 0
	}
}

//返回 座位玩家信息(与ExtDesk中的ChairList不同的是，这里的作为玩家除了自己之外Id是有隐藏的)
func (this *ExtPlayer) getChairList() []PlayerInfoByChair {
	//遍历桌子的 座位信息(即展示用户)
	charlist := []PlayerInfoByChair{}
	for _, v := range this.Dk.ChairList {
		charlist = append(charlist, v)
	}
	for i, v := range charlist {
		if v.Uid != this.Uid {
			charlist[i].Nick = this.Dk.ChangeNick(charlist[i].Nick)
		}
	}
	logs.Debug("原本桌子座位信息集合:", this.Dk.ChairList)
	logs.Debug("发送给", this.Nick, "经过隐藏之后的集合:", charlist)
	return charlist
}

//返回 更多玩家信息
func (this *ExtPlayer) getMorePlayer() []PlayerMsgByMore {
	var info []PlayerMsgByMore
	for _, v := range this.Dk.Players {
		var betall int64 //记录中总下注
		var wincount int //记录中赢得次数
		var nick string  //玩家名称
		if v.Uid != this.Uid {
			nick = this.Dk.ChangeNick(v.Nick)
		} else {
			nick = v.Nick
		}
		for _, v1 := range v.History {
			betall += v1.DownBet
			if v1.IsVictory {
				wincount += 1
			}
		}
		playerinfo := PlayerMsgByMore{}
		playerinfo.MatchCount = 20
		playerinfo.BetAll = betall
		playerinfo.WinCount = wincount
		playerinfo.Coins = v.Coins
		playerinfo.Head = v.Head
		playerinfo.Nick = nick
		info = append(info, playerinfo)
	}
	return info
}

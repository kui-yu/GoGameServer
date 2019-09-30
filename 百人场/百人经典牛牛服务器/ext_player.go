package main

import (
	"sync"
)

type ExtPlayer struct {
	sync.RWMutex // 锁、防止多线程读写
	Player
	DownBet            []int64      //玩家下注区域集合
	NoToBet            int          //未参与下注数
	OtherDownBet       []int64      //其他玩家下注
	BalanceResult      []int64      //玩家各个区域结算
	BalanceResultTosee []int64      //用来显示的玩家结算
	OtherBalanceResult []int64      //其他玩家输赢
	History            []BetHistory //历史输赢
	IsOnChair          bool         //玩家是否在座位上
	OldtherDownBet     []int64      //旧的其他玩家下注
	BetCount           int          //玩家下注次数
}

func (this *ExtPlayer) InitExtData() {
	//初始化玩家下注区域集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.DownBet = append(this.DownBet, 0)
	}
	//初始化其他玩家下注区域集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.OtherDownBet = append(this.OtherDownBet, 0)
	}
	//初始化自己玩家结算集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.BalanceResult = append(this.BalanceResult, 0)
	}
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.BalanceResultTosee = append(this.BalanceResultTosee, 0)
	}
	//初始化其他玩家结算集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.OtherBalanceResult = append(this.OtherBalanceResult, 0)
	}
	//初始化旧玩家下注集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.OldtherDownBet = append(this.OldtherDownBet, 0)
	}
	this.BetCount = 1
}

func (this *ExtPlayer) Rest() {
	this.DownBet = []int64{}
	//初始化玩家下注区域集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.DownBet = append(this.DownBet, 0)
	}
	this.OtherDownBet = []int64{}
	//初始化其他玩家下注区域集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.OtherDownBet = append(this.OtherDownBet, 0)
	}
	this.BalanceResult = []int64{}
	//初始化自己玩家结算集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.BalanceResult = append(this.BalanceResult, 0)
	}
	this.OtherBalanceResult = []int64{}
	//初始化其他玩家结算集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.OtherBalanceResult = append(this.OtherBalanceResult, 0)
	}
	this.OldtherDownBet = []int64{}
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.OldtherDownBet = append(this.OldtherDownBet, 0)
	}
	this.BalanceResultTosee = []int64{}
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.BalanceResultTosee = append(this.BalanceResultTosee, 0)
	}
	this.BetCount = 1
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
		if len(info) == gameConfig.LimitInfo.Userlist_count {
			break
		}
	}
	return info
}

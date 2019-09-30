package main

import (
	"strings"
)

type ExtPlayer struct {
	Player

	Sid          int     //座位id
	IsBank       bool    //是否是庄家
	IsFold       bool    //是否弃牌
	AllInStage   int     //第几轮全下,默认0xFF
	Cards        []int   //玩家手牌
	State        int     //用户状态
	CarryCoin    int64   //携带金币
	DownCoins    []int64 //下注的金币
	StageOperate int     //当前操作
}

func (this *ExtPlayer) Init() {
	this.Sid = -1
	this.CarryCoin = 0
	this.LiXian = false
	this.StageOperate = 0

	if this.Coins < 200000 {
		this.Coins = 200000
	}

	if this.Coins > gameConfig.DefSettCoin {
		this.CarryCoin = gameConfig.DefSettCoin
	} else {
		this.CarryCoin = this.Coins
	}

	this.State = UserStateWaitStart
	this.Reset()
}

func (this *ExtPlayer) Reset() {
	this.StageOperate = 0
	this.IsBank = false
	this.IsFold = false
	this.AllInStage = 0xFF
	this.Cards = []int{}
	this.DownCoins = []int64{}
}

func (this *ExtPlayer) SendNetMessage(cmd int, data interface{}, ouserNamePath ...interface{}) {

	// 判断data是否包含cmd,不包含添加cmd
	mdata := ConvertObjToMap(data)
	if v, ok := mdata["Id"]; !ok || v == 0 {
		mdata["Id"] = cmd
	}

	// 名字字段专门处理***
	if len(ouserNamePath) != 0 {
		userNamePath := ouserNamePath[0].(string)
		paths := strings.Split(userNamePath, "/")

		ReplaceMapField(mdata, paths, func(val interface{}) interface{} {
			oldname := val.(string)
			if oldname == this.Account {
				return oldname
			} else {
				return MarkNickName(oldname)
			}
		})
		this.SendNativeMsg(cmd, mdata)
	} else {
		this.SendNativeMsg(cmd, mdata)
	}
}

func (this *ExtPlayer) GetDownBet() int64 {
	var totalCoin int64 = 0
	for _, coin := range this.DownCoins {
		totalCoin += int64(coin)
	}
	return totalCoin
}

func (this *ExtPlayer) AddDownBet(stageIdx int, val int64) {
	this.CarryCoin -= val
	this.Coins -= val

	if len(this.DownCoins) == stageIdx {
		this.DownCoins = append(this.DownCoins, val)
	} else {
		this.DownCoins[stageIdx] += val
	}
}

func (this *ExtPlayer) getLastDownBet() int64 {
	dlen := len(this.DownCoins)
	return this.DownCoins[dlen-1]
}

package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) HandleDownBet(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != GAME_STATUS_DOWNBET {
		logs.Error("目前不是下注状态")
		p.SendNativeMsg(MSG_GAME_INFO_DOWNBET_REPLAY, DownBetReplay{
			Id:     MSG_GAME_INFO_DOWNBET_REPLAY,
			Result: 1,
			ErrStr: "目前不是下注状态！！！",
			Coins:  p.Coins,
		})
		return
	}
	data := DownBet{}
	err := json.Unmarshal([]byte(d.Data), &data)
	if err != nil {
		logs.Error("玩家下注处理----json转换出错:", err)
	}
	//判断筹码Id 和区域Id是否正确
	if data.AreaIndex < 0 || data.AreaIndex >= gameConfig.LimitInfo.BetCount {
		logs.Error("玩家下注处理----下注区域信息错误")
		p.SendNativeMsg(MSG_GAME_INFO_DOWNBET_REPLAY, DownBetReplay{
			Id:     MSG_GAME_INFO_DOWNBET_REPLAY,
			Result: 2,
			ErrStr: "下注区域信息错误！！！",
			Coins:  p.Coins,
		})
		return
	}
	if data.ChipIndex < 0 || data.ChipIndex >= len(G_DbGetGameServerData.GameConfig.TenChips) {
		logs.Error("玩家下注处理----筹码Id错误")
		p.SendNativeMsg(MSG_GAME_INFO_DOWNBET_REPLAY, DownBetReplay{
			Id:     MSG_GAME_INFO_DOWNBET_REPLAY,
			Result: 3,
			ErrStr: "筹码信息错误！！！",
			Coins:  p.Coins,
		})
		return
	}
	downCoins := G_DbGetGameServerData.GameConfig.TenChips[data.ChipIndex] //获取玩家下注金币
	//判断玩家下注数量是否操作玩家本身金币的4/1
	var allBet int64 //目前玩家所有下注数量
	for _, v := range p.DownBet {
		allBet += v
	}
	//判断限红
	if p.DownBet[data.AreaIndex]+downCoins > G_DbGetGameServerData.GameConfig.LimitRedMax {
		logs.Error("下注区域超过限红")
		p.SendNativeMsg(MSG_GAME_INFO_DOWNBET_REPLAY, DownBetReplay{
			Id:     MSG_GAME_INFO_DOWNBET_REPLAY,
			Result: 5,
			ErrStr: "下注区域超过限红！！！",
			Coins:  p.Coins,
		})
		return
	}
	//判断该玩家是否为庄
	if this.Zhuang.Uid == p.Uid {
		logs.Debug("该玩家为庄，无法下注", p.Nick)
		p.SendNativeMsg(MSG_GAME_INFO_DOWNBET_REPLAY, DownBetReplay{
			Id:     MSG_GAME_INFO_DOWNBET_REPLAY,
			Result: 6,
			ErrStr: "庄家无法下注！！！",
			Coins:  p.Coins,
		})
	}
	//将玩家下注保存
	this.DownBet[data.AreaIndex] += downCoins
	if !p.Robot {
		this.ZhenShiDownBet[data.AreaIndex] += downCoins
	}
	p.DownBet[data.AreaIndex] += downCoins
	//返回下注成功。
	logs.Debug("下注成功！")
	p.NoToBet = 0
	p.Coins -= downCoins
	p.SendNativeMsg(MSG_GAME_INFO_DOWNBET_REPLAY, DownBetReplay{
		Id:           MSG_GAME_INFO_DOWNBET_REPLAY,
		Result:       0,
		CoinsIndex:   data.ChipIndex,
		AreaIndex:    data.AreaIndex,
		SelfAllCoins: p.DownBet[data.AreaIndex],
		AllCoins:     this.DownBet[data.AreaIndex],
		BetAbleIndex: this.CanUseChip(p),
		Coins:        p.Coins,
	})
	//将除了该玩家的其他玩家的  OtherDownBet属性更新
	for _, v := range this.Players {
		//更新OldOtherdownbet
		if this.NeedUpdata {
			var oldother []int64
			for _, v1 := range v.OtherDownBet {
				oldother = append(oldother, v1)
			}
			v.OldtherDownBet = oldother
		}
		if p.Uid != v.Uid {
			v.OtherDownBet[data.AreaIndex] += downCoins
		}
	}
	this.NeedBro = true
	this.NeedUpdata = false
}

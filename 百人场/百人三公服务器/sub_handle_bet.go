package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) PlayerBet(p *ExtPlayer, m *DkInMsg) {
	req := &GABetInfo{}
	err := json.Unmarshal([]byte(m.Data), req)
	if err != nil {
		logs.Error("转化json失败 => PlayerBet()", err)
	}
	//不是下注阶段
	if this.Stage != STAGE_GAME_BET {
		return
	}
	//判断下注索引和区域是否正确
	if req.Coin < 0 || req.Coin > 9 || req.Place < 0 || req.Place > 3 {
		return
	}
	if req.Coin > len(G_DbGetGameServerData.GameConfig.TenChips)-1 {
		return
	}
	resplayer := &GSPlayerBet{
		Id: MSG_GAME_INFO_PLAYER_BET_REPLAY,
	}
	if G_DbGetGameServerData.GameConfig.TenChips[req.Coin] > p.Coins {
		resplayer.Code = 1
		resplayer.ErrStr = "金币不足，无法下注"
		p.SendNativeMsg(MSG_GAME_INFO_PLAYER_BET_REPLAY, resplayer)
		return
	}
	if p.IsDouble && (p.AllBet+G_DbGetGameServerData.GameConfig.TenChips[req.Coin])*4 > p.Coins-G_DbGetGameServerData.GameConfig.TenChips[req.Coin] {
		resplayer.Code = 2
		resplayer.ErrStr = "翻倍模式金币不足，请充值"
		p.SendNativeMsg(MSG_GAME_INFO_PLAYER_BET_REPLAY, resplayer)
		return
	}
	if p.PlaceBet[req.Place]+G_DbGetGameServerData.GameConfig.TenChips[req.Coin] > G_DbGetGameServerData.GameConfig.LimitRedMax {
		resplayer.Code = 3
		resplayer.ErrStr = "超出区域限红"
		p.SendNativeMsg(MSG_GAME_INFO_PLAYER_BET_REPLAY, resplayer)
		return
	}
	this.PlaceBet[req.Place] += G_DbGetGameServerData.GameConfig.TenChips[req.Coin] //桌子区域下注金币增加
	p.PlaceBet[req.Place] += G_DbGetGameServerData.GameConfig.TenChips[req.Coin]    //自己区域下注金币增加
	p.AllBet += G_DbGetGameServerData.GameConfig.TenChips[req.Coin]                 //总下注增加
	p.Coins -= G_DbGetGameServerData.GameConfig.TenChips[req.Coin]                  //自己金币减少
	p.IsBet = true                                                                  //有下注
	p.NotBet = 0                                                                    //没下注统计次数变0
	p.AccumulateBet += G_DbGetGameServerData.GameConfig.TenChips[req.Coin]          //累积下注金额
	//赋值发送下注给自己
	p.GetBetArr(p.IsDouble, p.AllBet)
	resplayer.BetArrAble = p.BetArrAble
	resplayer.Coins = p.Coins
	resplayer.CoinIndex = req.Coin
	resplayer.PlaceIndex = req.Place
	resplayer.MeAllCoin = p.PlaceBet[req.Place]
	resplayer.AllCoin = this.PlaceBet[req.Place]
	p.SendNativeMsg(MSG_GAME_INFO_PLAYER_BET_REPLAY, resplayer)
	//每秒群发下注信息给客户端
	p.SecondBet[req.Place] += G_DbGetGameServerData.GameConfig.TenChips[req.Coin]
	var b = true
	//避免重复添加计时器
	for _, v := range this.TList {
		if v.Id == int(p.Uid) {
			b = false
		}
	}
	//不存在再添加
	if b {
		this.AddTimer(int(p.Uid), 1, SentPerSecond, SecondToCli{
			Desk:    this,
			Player:  p,
			Uid:     p.Uid,
			AllCoin: this.PlaceBet,
			Place:   p.SecondBet,
		})
	}
}

//群发下注信息给其他客户端
func SentPerSecond(data interface{}) {
	d := data.(SecondToCli)
	for _, v := range d.Desk.Players {
		if v.Uid == d.Uid {
			continue
		}
		v.SendNativeMsg(MSG_GAME_INFO_PLAYER_BET_MASS, &GSPlayerBetMass{
			Id:      MSG_GAME_INFO_PLAYER_BET_MASS,
			AllCoin: d.AllCoin,
			Place:   d.Place,
		})
		d.Player.SecondBet = make([]int64, 4) //重置
	}
}

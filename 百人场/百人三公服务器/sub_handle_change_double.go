package main

import (
	"encoding/json"
)

//翻倍平倍改变
func (this *ExtDesk) ChangeDouble(p *ExtPlayer, m *DkInMsg) {
	var dobuleMode = new(GAChangedouble)
	_ = json.Unmarshal([]byte(m.Data), dobuleMode)
	res := &GSChangedouble{
		Id: MSG_GAME_INFO_CHANGE_DOUBLE_REPLAY,
	}
	//不是选倍数阶段
	if this.Stage != STAGE_GAME_DOUBLE {
		res.Code = 1
		res.ErrStr = "不是选倍数阶段"
		p.SendNativeMsg(MSG_GAME_INFO_CHANGE_DOUBLE_REPLAY, res)
		return
	}
	if dobuleMode.DoubleMode == 0 { //改为平倍
		p.IsDouble = false
		res.DoubleMode = 0
		p.GetBetArr(false, 0) //获取可下注筹码
		res.BetArrAble = p.BetArrAble
		p.SendNativeMsg(MSG_GAME_INFO_CHANGE_DOUBLE_REPLAY, res)
	} else if dobuleMode.DoubleMode == 1 { //改为翻倍
		p.IsDouble = true
		res.DoubleMode = 1
		p.GetBetArr(true, 0)
		res.BetArrAble = p.BetArrAble
		p.SendNativeMsg(MSG_GAME_INFO_CHANGE_DOUBLE_REPLAY, res)
	}
}

//获取最新下注筹码区间
func (this *ExtPlayer) GetBetArr(isDouble bool, allBet int64) []int64 {
	betArr := append([]int64{}, G_DbGetGameServerData.GameConfig.TenChips...)
	var betArrAble []int64
	if isDouble { //翻倍
		for i := len(G_DbGetGameServerData.GameConfig.TenChips) - 1; i >= 0; i-- {
			if (allBet+G_DbGetGameServerData.GameConfig.TenChips[i])*4 <= this.Coins-G_DbGetGameServerData.GameConfig.TenChips[i] {
				betArrAble = betArr[:i+1]
				break
			}
		}
	} else { //平倍
		for i := len(G_DbGetGameServerData.GameConfig.TenChips) - 1; i >= 0; i-- {
			if G_DbGetGameServerData.GameConfig.TenChips[i] <= this.Coins {
				betArrAble = betArr[:i+1]
				break
			}
		}
	}
	this.BetArrAble = len(betArrAble) - 1
	return betArrAble //返回下注金币切片的最大索引
}

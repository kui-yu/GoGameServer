package main

import (
	//"fmt"
	"logs"
	"time"
)

// 停止下注
func (this *ExtDesk) TimerStopBet(d interface{}) {
	this.Lock()
	defer this.Unlock()
	tAreaCoins := this.GetAreaCoinsList()

	// 消息公共部分
	sd := GGameStopBetNotify{
		Id:         MSG_GAME_INFO_STOP_BET_NOTIFY,
		Timer:      int32(gameConfig.Timer.StopBetNum) * 1000,
		TAreaCoins: tAreaCoins,
	}
	// 每个人的下注都不一样  需要单独处理
	players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range players {
		p := v.(*ExtPlayer)
		sd.PAreaCoins = p.GetNTAreaCoinsList()
		sd.OtherBetList = this.SeatMgr.GetOtherNewBetList2(p.Uid)
		p.SendNativeMsg(MSG_GAME_INFO_STOP_BET_NOTIFY, sd)
	}

	for _, user := range this.Players {
		user.ColAreaCoins()
	}

	this.GameState = MSG_GAME_INFO_STOP_BET_NOTIFY
	this.AddTimer(gameConfig.Timer.StopBet, gameConfig.Timer.StopBetNum, this.TimerOpen, nil)

	logs.Debug("-------------------当前库存的值", CD, "目标库存", CalPkAll(StartControlTime, time.Now().Unix()))

	//真实玩家未下注
	if this.GetUserAreaCoins() == 0 || GetCostType() == 2 { //控制75%的胜率
		w := RandInt64(4)
		if w >= 1 { //123庄输
			logs.Debug("-------------------控制75%的胜率控制输")
			this.allotCard(false)
		} else { //0赢
			logs.Debug("-------------------控制75%的胜率控制赢")
			this.allotCard(true)
		}
		return
	} else if CD-CalPkAll(StartControlTime, time.Now().Unix()) >= 0 {
		//如果是体验场也不进入库存、层级判断
		// 东、西、南、北  发剩余牌3张
		logs.Debug("-------------------随机发牌")
		for i := 0; i < 4; i++ {
			this.IdleCard[i] = append(this.IdleCard[i], this.CardMgr.SendCard(3)...)
			MCard, TCard := this.CardMgr.GetMaxCards(this.IdleCard[i])
			this.MIdleCard = append(this.MIdleCard, MCard)
			this.TypeList = append(this.TypeList, TCard)
		}
		// 庄家牌
		this.BankerCard = append(this.BankerCard, this.CardMgr.SendCard(3)...)
		this.MBankerCard, this.TBankerType = this.CardMgr.GetMaxCards(this.BankerCard)
		return
	}
	//小郑风控
	logs.Debug("-------------------未达到目标库存控制赢")
	this.allotCard(true)
	return
}

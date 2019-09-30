// 定时器回调函数，控制百人场状态
package main

import (
	"bl.com/util"
)

// 结束
func (this *ExtDesk) TimerOver(d interface{}) {
	this.AddTimer(gameConfig.Timer.Over, gameConfig.Timer.OverNum, this.TimerShuffle, nil)
}

// 结算
func (this *ExtDesk) GameEnd(typeList []int32) (int64, []bool, []int64) {
	var winArea [9]bool
	var tWinArea [9]int64
	var double [9]float64
	for i := range winArea {
		winArea[i] = false
		double[i] = 0
		tWinArea[i] = 0
	}

	// 计算输赢
	var loseCoins int64 = 0

	// 东南西北输赢情况
	for i := 0; i < 4; i++ {
		winArea[i] = this.CardMgr.CompareCard(this.MIdleCard[i], this.MBankerCard)
		if winArea[i] {
			double[i] = gameConfig.Double[i]
		}
	}

	// 押庄输赢情况
	btype := this.CardMgr.GetCardsType(this.MBankerCard)
	switch btype {
	case this.CardMgr.GFDouble:
		if this.CardMgr.GetLogicValue(this.MBankerCard[1]) < util.Card_8 && this.CardMgr.GetLogicValue(this.MBankerCard[1]) != util.Card_A {
			break
		}

		winArea[INDEX_BANKER_DOUBLE-1] = true
	case this.CardMgr.GFShunZi:
		winArea[INDEX_BANKER_SHUNZI-1] = true
	case this.CardMgr.GFJinHua:
		if this.CardMgr.GetLogicValue(this.MBankerCard[0]) < util.Card_10 && this.CardMgr.GetLogicValue(this.MBankerCard[0]) != util.Card_A {
			break
		}
		winArea[INDEX_BANKER_JINHUA-1] = true
	case this.CardMgr.GFShunJin:
		if this.CardMgr.GetLogicValue(this.MBankerCard[0]) < util.Card_8 && this.CardMgr.GetLogicValue(this.MBankerCard[0]) != util.Card_A {
			break
		}
		winArea[INDEX_BANKER_SHUNJIN-1] = true
	case this.CardMgr.GFBaoZi:
		if this.CardMgr.GetLogicValue(this.MBankerCard[0]) < util.Card_8 && this.CardMgr.GetLogicValue(this.MBankerCard[0]) != util.Card_A {
			break
		}
		winArea[INDEX_BANKER_BAOZI-1] = true
	}
	//控制盈利制
	zong := this.IsWin(this.BankerCard, this.IdleCard)
	AddCD(zong)
	// 计算输赢
	for i, v := range winArea {
		if !v {
			continue
		}

		double[i] = gameConfig.Double[i]
		coins := int64(float64(this.GetUserAreaCoin(i)) * double[i])
		loseCoins += coins
		tWinArea[i] = coins
	}

	players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range players {
		p := v.(*ExtPlayer)
		p.BuildWinList(double[:])
		p.AddWinList()
		p.AddBetList()
	}
	totCoins := this.GetUserAreaCoins()
	if GetCostType() == 1 { //如果不是体验场再修改库存
		stock := totCoins - loseCoins
		AddLocalStock(stock)
		if stock > 0 {
			this.totalStock += stock
		}
	}
	// if gameConfig.DeskInfo.Win == 0 {
	// 	return totCoins - loseCoins, winArea[:], tWinArea[:]
	// }

	// // 有设置盈利率，需要计算盈利率
	// if this.GetUserAreaCoins() > 0 {
	// 	this.tCount++
	// 	this.totCoins += float64(totCoins)
	// 	this.wCoins += float64(totCoins - loseCoins)
	// 	logs.Debug("局号，测试次数，总下注，总赢取, 输赢率：", this.GameId, this.tCount, this.totCoins/100, this.wCoins/100, this.wCoins/this.totCoins*100)
	// }

	// rate := this.wCoins / this.totCoins * 100
	// if int(rate) < gameConfig.DeskInfo.Win+5 && int(rate) > gameConfig.DeskInfo.Win-5 || this.tCount >= 100 {
	// 	this.wCoins = 0
	// 	this.totCoins = 0
	// 	this.tCount = 0
	// }

	return totCoins - loseCoins, winArea[:], tWinArea[:]
}

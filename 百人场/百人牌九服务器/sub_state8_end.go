package main

import (
	"logs"

	"bl.com/paigow"
)

// 结束
func (this *ExtDesk) TimerOver(d interface{}) {
	for _, v := range this.Players {
		v.IsBet = false
		if v.LiXian {
			logs.Debug("*******有用户离线：%v", v.Nick)
			this.SeatMgr.DelPlayer(v)
			this.LeaveByForce(v)
		}
	}
	this.AddTimer(gameConfig.Timer.Over, gameConfig.Timer.OverNum, this.TimerShuffle, nil)
}

// 结算
func (this *ExtDesk) GameEnd(typeList []int32) (int64, []bool, []int64) {
	var winArea [8]bool
	var tWinArea [8]int64
	var double [8]float64
	for i := range winArea {
		winArea[i] = false
		double[i] = 0
		tWinArea[i] = 0
	}

	// 区域输赢情况
	if typeList[0] == paigow.TYPE_ZHIZUN {
		winArea[INDEX_BANKER_ZHIZUN-1] = true
	} else if typeList[0] == paigow.TYPE_TIAN {
		winArea[INDEX_BANKER_TIAN-1] = true
	}
	// 天
	isWin := paigow.CompareCard(this.IdleCard[0], this.BankerCard)
	winArea[INDEX_TIAN_WIN-1] = isWin
	winArea[INDEX_TIAN_LOSS-1] = !isWin
	// 地
	isWin = paigow.CompareCard(this.IdleCard[1], this.BankerCard)
	winArea[INDEX_DI_WIN-1] = isWin
	winArea[INDEX_DI_LOSS-1] = !isWin
	// 人
	isWin = paigow.CompareCard(this.IdleCard[2], this.BankerCard)
	winArea[INDEX_REN_WIN-1] = isWin
	winArea[INDEX_REN_LOSS-1] = !isWin

	// 计算输赢
	var loseCoins int64 = 0
	for i, v := range winArea {
		var d float64
		switch i + 1 {
		case INDEX_TIAN_WIN:
			fallthrough
		case INDEX_TIAN_LOSS:
			fallthrough
		case INDEX_DI_WIN:
			fallthrough
		case INDEX_DI_LOSS:
			fallthrough
		case INDEX_REN_WIN:
			fallthrough
		case INDEX_REN_LOSS:
			fallthrough
		case INDEX_BANKER_TIAN:
			fallthrough
		case INDEX_BANKER_ZHIZUN:
			d = gameConfig.Double[i]
		}

		if !v || double[i] != 0 {
			continue
		}

		loseCoins += int64(float64(this.GetUserAreaCoin(i)) * d)
		double[i] = d
	}

	// 计算玩家输赢
	for i, v := range winArea {
		if !v {
			continue
		}

		coins := int64(float64(this.GetUserAreaCoin(i)) * double[i])

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
	if GetCostType() == 1 {
		stock := totCoins - loseCoins
		//AddLocalStock(stock)
		if stock > 0 {
			this.totalStock += stock
		}
	}
	return totCoins - loseCoins, winArea[:], tWinArea[:]
}

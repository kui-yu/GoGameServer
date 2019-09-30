package main

import (
	"logs"
)

func (this *ExtDesk) GameStateSettle() {
	this.BroadStageTime(TIME_STAGE_ZERO_NUM)
	//进入倒计时
	this.runTimer(TIME_STAGE_ZERO_NUM, this.GameStateSettleEnd)
}

//结算阶段
func (this *ExtDesk) GameStateSettleEnd(d interface{}) {
	this.ClearTimer()
	var tempWin int = 0
	//先比牌
	for i := 1; i < len(this.Players); i++ {
		p1 := GCardType{
			NiuPoint: this.Players[tempWin].NiuPoint,
			NiuCards: this.Players[tempWin].NiuCards,
			HandCard: this.Players[tempWin].HandCard,
		}
		p2 := GCardType{
			NiuPoint: this.Players[i].NiuPoint,
			NiuCards: this.Players[i].NiuCards,
			HandCard: this.Players[i].HandCard,
		}
		rs := SoloResult(p1, p2)
		if rs > 1 {
			tempWin = i
		}
	}
	logs.Debug("tempWin", tempWin)
	//算分
	for i := 0; i < len(this.Players); i++ {
		if i != tempWin {
			winMultiple := this.Players[tempWin].NiuMultiple * this.Players[tempWin].BetMultiple * this.Players[i].BetMultiple
			this.Players[tempWin].WinMultiple += winMultiple
			this.Players[i].WinMultiple -= winMultiple
		}
	}

	var tableMultiple int64 //桌面池子总倍数
	var tableMoney int64    //桌面池子总金额
	for _, v := range this.Players {
		winCoins := int64(v.WinMultiple * this.Bscore)
		if winCoins < 0 {
			winCoins = -winCoins
			if winCoins > v.Coins && GetCostType() == 1 { //如果不体验场再进入判断玩家金币不够输的情况
				tableMoney += v.Coins
				v.WinCoins = -v.Coins
			} else {
				tableMoney += winCoins
				v.WinCoins = -winCoins
			}
		} else {
			tableMultiple += int64(v.WinMultiple)
		}
	}

	//平摊池子
	for _, v := range this.Players {

		if v.WinMultiple > 0 {
			v.WinCoins = int64(int64(v.WinMultiple) * tableMoney / tableMultiple)
			//手续费
			v.RateCoin = float64(v.WinCoins) * this.Rate
			v.WinCoins = v.WinCoins - int64(v.RateCoin)
		}
	}

	//广播结算
	var winInfos GWinInfosReply
	for _, v := range this.Players {

		v.Coins += int64(v.WinCoins)

		if v.Robot && GetCostType() == 1 { //如果不是体验场在进行层级概率更新
			//当前库存
			AddLocalStock(v.WinCoins)
			AddCD(v.WinCoins)
			logs.Debug("当前库存:", CD)
		}

		info := GWinInfo{
			Uid:     v.Uid,
			ChairId: v.ChairId,
			WinCoin: v.WinCoins,
			Coins:   v.Coins,
		}
		winInfos.Infos = append(winInfos.Infos, info)
	}
	winInfos.Id = MSG_GAME_INFO_SETTLE
	winInfos.InfoCount = int32(len(this.Players))
	winInfos.WinChairId = this.Players[tempWin].ChairId
	if GetCostType() == 1 { //如果不是体验场再进行数据库和大厅消息的更新
		//数据交互
		this.PutSqlData()
	}

	//发送结算信息
	this.BroadcastAll(MSG_GAME_INFO_SETTLE, winInfos)

	this.GameState = GAME_STATUS_END
	this.BroadStageTime(TIMER_OVER_NUM)
	//发送游戏结束消息
	for _, p := range this.Players {
		p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:      MSG_GAME_LEAVE_REPLY,
			Cid:     p.ChairId,
			Uid:     p.Uid,
			Result:  0,
			Token:   p.Token,
			NoToCli: true,
		})
	}

	//开始归还桌子定时器
	this.GameOver()
}

func GetNiuMultiple(niuPoint int32) int {
	if niuPoint > 10 {
		return 5
	} else if niuPoint == 10 {
		return 4
	} else if niuPoint > 7 {
		return 3
	} else if niuPoint == 7 {
		return 2
	} else {
		return 1
	}
}

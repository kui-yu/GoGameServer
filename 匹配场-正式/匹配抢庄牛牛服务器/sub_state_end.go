package main

func (this *ExtDesk) GameStateSettle() {

	this.BroadStageTime(TIME_STAGE_ZERO_NUM)
	//进入倒计时
	this.runTimer(TIME_STAGE_ZERO_NUM, this.GameStateSettleEnd)
}

//结算阶段
func (this *ExtDesk) GameStateSettleEnd(d interface{}) {

	//先比牌
	for i := 0; i < len(this.Players); i++ {
		if this.Banker != this.Players[i].ChairId {
			p1 := GCardType{
				NiuPoint: this.Players[this.Banker].NiuPoint,
				NiuCards: this.Players[this.Banker].NiuCards,
				HandCard: this.Players[this.Banker].HandCard,
			}
			p2 := GCardType{
				NiuPoint: this.Players[i].NiuPoint,
				NiuCards: this.Players[i].NiuCards,
				HandCard: this.Players[i].HandCard,
			}
			rs := SoloResult(p1, p2)
			if rs > 1 {
				//玩家赢
				winMultiple := this.Players[i].NiuMultiple * this.Players[i].BetMultiple
				this.Players[i].WinMultiple += winMultiple
				this.Players[this.Banker].WinMultiple -= winMultiple
			} else {
				//庄家赢
				winMultiple := this.Players[this.Banker].NiuMultiple * this.Players[i].BetMultiple
				this.Players[this.Banker].WinMultiple += winMultiple
				this.Players[i].WinMultiple -= winMultiple
			}
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
			v.RateCoins = float64(v.WinCoins) * this.Rate
			v.WinCoins = v.WinCoins - int64(v.RateCoins)
		}
	}

	//广播结算
	var winInfos GWinInfosReply
	for _, v := range this.Players {

		v.Coins += int64(v.WinCoins)

		if v.Robot && GetCostType() == 1 {
			//当前库存
			AddLocalStock(v.WinCoins)
			AddCD(v.WinCoins)
		}
		info := GWinInfo{
			Uid:      v.Uid,
			ChairId:  v.ChairId,
			WinCoin:  v.WinCoins,
			Coins:    v.Coins,
			NiuCards: v.NiuCards,
			NiuPoint: v.NiuPoint,
		}
		winInfos.Infos = append(winInfos.Infos, info)
	}
	winInfos.Id = MSG_GAME_INFO_SETTLE
	winInfos.InfoCount = int32(len(this.Players))
	//数据交互
	if GetCostType() == 1 { //如果不是体验场再进行数据库更新和大厅消息记录
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
	this.ClearTimer()
	//
	this.GameOverLeave()
	//开始归还桌子定时器
	this.AddTimer(TIMER_OVER, TIMER_OVER_NUM, this.TimerOver, nil)
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

func (this *ExtDesk) TimerOver(d interface{}) {
	this.GameState = GAME_STATUS_FREE
	this.JuHao = ""
	this.DeskMgr.BackDesk(this)
}

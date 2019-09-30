package main

import (
	"time"
)

func (this *ExtDesk) GameStateSettle() {
	//结算阶段
	this.BroadStageTime(STAGE_SETTLE_TIME)

	//结算结果处理
	for _, v := range this.Players {
		if v.ChairId != this.Banker {

			winMultiple := v.PlayMultiple * this.Players[this.Banker].CallMultiple
			winCoins := int64(this.Bscore * winMultiple)
			rs := compareCards(this.Players[this.Banker].HandCards, v.HandCards)
			if rs == 2 {
				//闲家赢
				v.WinMultiple = winMultiple
				v.WinCoins += winCoins
				this.Players[this.Banker].WinCoins -= winCoins
				this.Players[this.Banker].WinMultiple -= winMultiple
			} else {
				//庄家赢
				this.Players[this.Banker].WinList = append(this.Players[this.Banker].WinList, v.ChairId)
				v.WinCoins -= winCoins
				v.WinMultiple = winMultiple
				this.Players[this.Banker].WinCoins += winCoins
				this.Players[this.Banker].WinMultiple += winMultiple
			}
		}
	}

	//发送结算信息
	this.RsInfo = GSSettleInfo{}
	for _, v := range this.Players {
		//游戏消费税
		if v.WinCoins > 0 {
			v.RateCoins = float64(v.WinCoins) * this.Rate
			v.WinCoins = v.WinCoins - int64(v.RateCoins)
		}
		//添加游戏记录
		recordInfo := GSRecordInfo{
			WinCoins: v.WinCoins,
			WinDate:  time.Now().Format("2006-01-02 15:04:05"),
		}
		v.RecordInfos = append(v.RecordInfos, recordInfo)
		// logs.Debug("输赢", recordInfo)
		//结算信息
		v.Coins += v.WinCoins
		//当前库存
		if GetCostType() == 1 { //如果不是体验场并且是机器人再更新库存概率
			if v.Robot {
				AddLocalStock(v.WinCoins)
			}
		}

		info := GSSettlePlayInfo{
			ChairId:         v.ChairId,
			HandCard:        v.HandCards,
			WinCoins:        v.WinCoins,
			Coins:           v.Coins,
			WinList:         v.WinList,
			BankerMultiples: this.Players[this.Banker].CallMultiple,
		}
		if v.ChairId != this.Banker {
			info.PlayerMultiples = v.PlayMultiple
		} else {
			info.PlayerMultiples = v.CallMultiple
		}

		this.RsInfo.PlayInfos = append(this.RsInfo.PlayInfos, info)
	}
	this.RsInfo.Round = this.Round
	this.RsInfo.Id = MSG_GAME_INFO_SETTLE_INFO_REPLY
	this.RsInfo.PutInfos = this.PutInfos

	//是否设置离开，0离开，1不离开
	var isLeave int32 = 0
	if this.Round < this.TotalRound {
		isLeave = 1
	}

	if GetCostType() == 1 { //如果不是体验场的话再进行金币限制和数据库更新
		for _, v := range this.Players {
			if v.Coins < int64(this.Bscore*len(this.Players)) {
				isLeave = 0
				break
			}
		}
		//数据交互
		this.PutSqlData(isLeave)
	}

	// logs.Debug("结算信息", this.RsInfo)

	for _, v := range this.Players {
		recordInfo := GSRecordInfos{
			Id:    MSG_GAME_INFO_RECORD_INFO_REPLY,
			Infos: v.RecordInfos,
		}
		v.SendNativeMsg(MSG_GAME_INFO_RECORD_INFO_REPLY, &recordInfo)
	}
	//通知记录之后
	this.BroadcastAll(MSG_GAME_INFO_SETTLE_INFO_REPLY, &this.RsInfo)

	if isLeave == 0 {
		this.ClearTimer()
		//总结算
		var rsInfoEnd GSSettleInfoEnd
		for _, v := range this.Players {
			var coins int64 = 0
			for _, c := range v.RecordInfos {
				coins += c.WinCoins
			}
			infoEnd := GSSettlePlayInfoEnd{
				Uid:      v.Uid,
				ChairId:  v.ChairId,
				WinCoins: coins,
				Coins:    v.Coins,
			}
			rsInfoEnd.PlayInfos = append(rsInfoEnd.PlayInfos, infoEnd)
		}
		rsInfoEnd.Id = MSG_GAME_INFO_SETTLE_INFO_END_REPLY
		// logs.Debug("总结算", rsInfoEnd)
		//通知记录之后
		this.BroadcastAll(MSG_GAME_INFO_SETTLE_INFO_END_REPLY, &rsInfoEnd)

		//系统定义结束
		this.SysTableEnd()
	} else {
		//等待重新开始
		this.runTimer(STAGE_SETTLE_TIME, this.GameStateSettleEnd)
	}
}

//阶段-结算跳转重新开始
func (this *ExtDesk) GameStateSettleEnd(d interface{}) {
	this.nextStage(STAGE_RESTART)
}

//系统定义结束
func (this *ExtDesk) SysTableEnd() {
	this.GameState = GAME_STATUS_END
	this.BroadStageTime(TIMER_OVER_NUM)
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
	this.GameOverLeave()

	//桌子初始化
	this.GameState = GAME_STATUS_FREE
	this.ResetTable()
	this.DeskMgr.BackDesk(this)
}

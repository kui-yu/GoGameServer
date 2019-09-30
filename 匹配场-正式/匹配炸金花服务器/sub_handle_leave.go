package main

import "logs"

// 底层统一离开
func (this *ExtDesk) Leave(p *ExtPlayer) {
	logs.Debug("run this leave")
	if this.GameState == GAME_STATUS_FREE {
		info := GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Result: 0,
			Cid:    p.ChairId,
			Uid:    p.Uid,
		}
		this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &info)
		logs.Debug("离开消息发送", info)
		this.DeskMgr.LeaveDo(p.Uid)
		logs.Debug("从桌子管理器离开")
		// logs.Debug("run this")
		for id, v := range this.Players { //清理桌子玩家信息
			if v == p {
				this.Players = append(this.Players[:id], this.Players[id+1:]...)
			}
		}
		return
	}
	if p.CardType != 2 {
		info := GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Result: 1,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Err:    "已在游戏中，退出失败。",
		}
		p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &info)
		return
	}

	info := GLeaveReply{
		Id:     MSG_GAME_LEAVE_REPLY,
		Result: 0,
		Cid:    p.ChairId,
		Uid:    p.Uid,
	}
	this.PlayerLeave(p)

	p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &info)
	p.LiXian = true
	this.DeskMgr.LeaveDo(p.Uid)
}

func (this *ExtDesk) HandleIsLeave(p *ExtPlayer, d *DkInMsg) {
	if this.GameState == GAME_STATUS_END || p.CardType != 2 || this.GameState == STAGE_SETTLE {
		return
	}

	info := GSPlayerLeave{}
	info.ChairId = p.ChairId
	info.Id = MSG_GAME_INFO_LEAVE_REPLY
	this.PlayerLeave(p)
	info.LeaveType = 2

	p.SendNativeMsg(MSG_GAME_INFO_LEAVE_REPLY, &info)
	p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
		Id:      MSG_GAME_LEAVE_REPLY,
		Result:  0,
		Cid:     p.ChairId,
		Uid:     p.Uid,
		Err:     "玩家正在游戏中，不能离开",
		Robot:   p.Robot,
		NoToCli: true,
	})
	p.LiXian = true
	this.DeskMgr.LeaveDo(p.Uid)
}

func (this *ExtDesk) PlayerLeave(p *ExtPlayer) {
	p.IsLeave = 1
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Mini:        false,
		SetLeave:    0,
	}

	valid := int64(0) //下注*低分
	for i := 0; i < len(p.PayCoin); i++ {
		valid += p.PayCoin[i]
	}
	valid = valid * this.Bscore
	dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
		UserId:      p.Uid,
		UserAccount: p.Account,
		BetCoins:    int64(this.Bscore),
		ValidBet:    valid,  //下注*低分
		PrizeCoins:  -valid, //输赢金币
		Robot:       p.Robot,
		WaterProfit: p.RateCoins,
		WaterRate:   this.Rate,
	})
	if GetCostType() == 1 {
		p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
	}

	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Round:       this.Round,
	}

	sum := int64(0)
	for i := 0; i < len(p.PayCoin); i++ {
		sum += p.PayCoin[i]
	}
	rddata := GGameRecordInfo{
		UserId:      p.Uid,
		UserAccount: p.Account,
		Robot:       p.Robot,
		CoinsBefore: p.Coins + sum*this.Bscore,
		BetCoins:    sum * this.Bscore, //下注金币
		Coins:       -valid,
		CoinsAfter:  p.Coins,
		Cards:       p.OldHandCard,
		Multiple:    1,
		Score:       this.Bscore,
	}
	// logs.Debug("中途离开wincoins", rddata)
	rdreq.UserRecord = append(rdreq.UserRecord, rddata)
	if GetCostType() == 1 {
		p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
	}
}

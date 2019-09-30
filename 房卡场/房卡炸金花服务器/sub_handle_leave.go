package main

import (
	"encoding/json"
	"logs"
)

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
	logs.Debug("玩家离开")
	//广播给其他人，掉线
	if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
		if p.ChairId == 0 {
			this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:     MSG_GAME_LEAVE_REPLY,
				Cid:    p.ChairId,
				Uid:    p.Uid,
				Result: 1,
				Token:  p.Token,
				Err:    "房主解散该房间",
			})
			logs.Debug("房主解散该房间")
			//this.TimerOver()
			this.HouseOwnerLeave()
		} else {
			this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:     MSG_GAME_LEAVE_REPLY,
				Cid:    p.ChairId,
				Uid:    p.Uid,
				Result: 0,
				Token:  p.Token,
			})
			this.DelPlayer(p.Uid)
			this.DeskMgr.LeaveDo(p.Uid)
		}
	} else {
		this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Result: 2,
			Token:  p.Token,
			Err:    "游戏已开始，请发起解散",
		})
	}
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
		GameModule:  this.TableConfig.GameModule,
		PayType:     this.TableConfig.PayType,
		GameRoundNo: this.JuHao,
	}

	sum := int64(0)
	for i := 0; i < len(p.PayCoin); i++ {
		sum += p.PayCoin[i]
	}
	rddata := GGameRecordInfo{
		UserId:        p.Uid,
		UserAccount:   p.Account,
		Robot:         p.Robot,
		CoinsBefore:   p.Coins - p.WinCoins,
		BetCoins:      sum * this.Bscore, //下注金币
		Coins:         p.WinCoins,
		CoinsAfter:    p.Coins,
		Cards:         p.OldHandCard,
		BetMultiple:   0,
		BrandMultiple: 0,
		Multiple:      1,
		Score:         this.Bscore,
	}
	// logs.Debug("中途离开wincoins", rddata)
	rdreq.UserRecord = append(rdreq.UserRecord, rddata)
	if GetCostType() == 1 {
		p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
	}
}

//房卡
//玩家解散
func (this *ExtDesk) HandleDisMiss(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("玩家解散")
	data := GADismiss{}
	json.Unmarshal([]byte(d.Data), &data)

	if data.IsDismiss == 1 {
		//同意解散
		if this.GameState != STAGE_DISMISS && this.GameState == STAGE_INIT {
			//跳到游戏解散阶段
			this.nextStage(STAGE_DISMISS)
		}
		if this.GameState != STAGE_DISMISS {
			p.SendNativeMsg(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
				Id:        MSG_GAME_INFO_DISMISS_REPLY,
				DisPlayer: this.DisPlayer,
				IsDismiss: 2,
				Message:   "请在准备阶段发起解散",
			})
			return
		}
		this.DisPlayer = append(this.DisPlayer, p.ChairId)

	} else {
		this.ClearTimer() /////
		//解散清空
		this.DisPlayer = []int32{}
		//回到游戏准备阶段
		this.GameState = STAGE_INIT
		this.BroadStageTime(0)

		this.BroadcastAll(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
			Id:        MSG_GAME_INFO_DISMISS_REPLY,
			DisPlayer: this.DisPlayer,
			IsDismiss: 0,
			Message:   p.Nick + "不同意解散",
		})
		return
	}
	//

	if len(this.DisPlayer) == len(this.Players) {
		this.BroadcastAll(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
			Id:        MSG_GAME_INFO_DISMISS_REPLY,
			DisPlayer: this.DisPlayer,
			IsDismiss: 3,
		})
		////
		this.GameState = GAME_STATUS_END
		this.BroadStageTime(TIMER_OVER_NUM)
		if GetCostType() != 1 {
			for _, p := range this.Players {
				p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
					Id:      MSG_GAME_LEAVE_REPLY,
					Result:  0,
					Cid:     p.ChairId,
					Uid:     p.Uid,
					Robot:   p.Robot,
					NoToCli: true,
				})
			}
		}
		//同意人数已满，解散
		this.TimerOver("")
	} else {
		//部分同意
		this.BroadcastAll(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
			Id:        MSG_GAME_INFO_DISMISS_REPLY,
			DisPlayer: this.DisPlayer,
			IsDismiss: 1,
		})
	}
}

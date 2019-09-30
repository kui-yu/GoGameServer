package main

import (
	"fmt"
)

type ExtDesk struct {
	Desk
	CardMgr MgrCard // 扑克牌牌管理
	//======配置信息======
	TableConfig GATableConfig
	Rate        float64
	//=====自定义信息=====
	Round     int     //回合
	Banker    int32   //庄家
	DisPlayer []int32 //解散玩家
}

//重置桌子
func (this *ExtDesk) ResetTable() {
	this.JuHao = ""
	this.TableConfig = GATableConfig{}
	this.Rate = 0
	this.Round = 0
	this.Banker = -1
	this.DisPlayer = []int32{}
}

//阶段执行
func (this *ExtDesk) nextStage(stage int) {
	this.GameState = stage
	//清空所有定时器
	this.ClearTimer()

	if this.GameState == STAGE_INIT {
		this.GameStateInit()
	} else if this.GameState == GAME_STATUS_START {
		//进入游戏开始
		this.GameStateStart()
	} else if this.GameState == STAGE_CALL_BANKER {
		//进入玩家叫庄
		this.GameStateCall()
	} else if this.GameState == STAGE_CALL_SCORE {
		//进入玩家叫分
		this.GameStateBet()
	} else if this.GameState == STAGE_DEAL {
		//进入发牌
		this.GameStateDeal()
	} else if this.GameState == STAGE_PLAY {
		//进入玩牌
		this.GameStatePlay()
	} else if this.GameState == STAGE_SETTLE {
		//进入结算
		this.GameStateSettle()
	} else if this.GameState == STAGE_DISMISS {
		//进入解散阶段
		this.GameStateDismiss()
	}
}

//广播阶段
func (this *ExtDesk) BroadStageTime(time int32) {
	stage := GSStageInfo{
		Id:        MSG_GAME_INFO_STAGE,
		Stage:     int32(this.GameState),
		StageTime: time,
	}
	this.BroadcastAll(MSG_GAME_INFO_STAGE, &stage)
}

//自封装定时器
func (this *ExtDesk) runTimer(t int, h func(interface{})) {
	//定时器ID，定时器时间，可执行函数，可执行参数
	this.AddTimer(10, t, h, nil)
}

//玩家掉线广播
func (this *ExtDesk) HandleDisConnect(p *ExtPlayer, d *DkInMsg) {
	p.LiXian = true
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 2,
	})
	// fmt.Println("桌子状态:", this.GameState)
	if this.GameState == GAME_STATUS_FREE {
		fmt.Println("玩家掉线:", p.Nick)
		if this.FkOwner == p.Uid {
			this.ClearTimer()
			this.GameState = GAME_STATUS_END
			this.BroadStageTime(0)
			//玩家离开
			for _, p := range this.Players {
				p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
					Id:     MSG_GAME_LEAVE_REPLY,
					Result: 0,
					Cid:    p.ChairId,
					Uid:    p.Uid,
					Token:  p.Token,
				})
			}
			this.GameOverLeave()
			//归还桌子
			this.GameState = GAME_STATUS_FREE
			this.ResetTable()
			this.DeskMgr.BackDesk(this)
		} else {
			for _, v := range this.Players {
				v.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
					Id:     MSG_GAME_LEAVE_REPLY,
					Result: 0,
					Cid:    p.ChairId,
					Uid:    p.Uid,
					Token:  p.Token,
					Robot:  p.Robot,
				})
			}
			this.DelPlayer(p.Uid)
			this.DeskMgr.LeaveDo(p.Uid)
		}
	}

}

//数据通信
func (this *ExtDesk) PutSqlData(isLeave int32) {
	//发送结算消息给数据库
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Mini:        false,
		SetLeave:    isLeave,
		Round:       this.Round,
		NoSaveCoin:  1,
		RoomNo:      this.FkNo,
	}
	if this.TableConfig.GameModule == 2 {
		for _, v := range this.Players {
			valid := v.WinCoins
			if valid < 0 {
				valid = -valid
			} else {
				valid = int64(this.TableConfig.BaseScore)
			}
			dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
				UserId:      v.Uid,
				UserAccount: v.Account,
				BetCoins:    0,
				ValidBet:    valid,
				PrizeCoins:  int64(v.WinCoins),
				Robot:       v.Robot,
				WaterProfit: v.RateCoins,
				WaterRate:   this.Rate,
			})
			v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
			dbreq.UserCoin = []GGameEndInfo{}
		}
	} else {
		for _, v := range this.Players {
			// dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
			// 	UserId:      v.Uid,
			// 	UserAccount: v.Account,
			// 	BetCoins:    0,
			// 	ValidBet:    0,
			// 	PrizeCoins:  0,
			// 	Robot:       v.Robot,
			// 	WaterProfit: 0,
			// 	WaterRate:   0,
			// })
			valid := v.WinCoins
			if valid < 0 {
				valid = -valid
			} else {
				valid = int64(this.TableConfig.BaseScore)
			}
			dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
				UserId:      v.Uid,
				UserAccount: v.Account,
				BetCoins:    0,
				ValidBet:    valid,
				PrizeCoins:  int64(v.WinCoins),
				Robot:       v.Robot,
				WaterProfit: v.RateCoins,
				WaterRate:   this.Rate,
			})
			v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
			dbreq.UserCoin = []GGameEndInfo{}
		}
	}
	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameType:    this.TableConfig.GameType,
		GameRoundNo: this.JuHao,
		Round:       this.Round,
		RoomNo:      this.FkNo,
	}
	for _, v := range this.Players {
		if v.Robot {
			continue
		}
		isBanker := 0
		betMultiple := v.BetMultiple
		if this.Banker == v.ChairId {
			betMultiple = v.CallMultiple
			isBanker = 1
		}

		multiple := v.WinMultiple
		if multiple < 0 {
			multiple = -multiple
		}

		rddata := GGameRecordInfo{
			UserId:        v.Uid,
			UserAccount:   v.Account,
			Robot:         v.Robot,
			CoinsBefore:   v.Coins - v.WinCoins,
			BetCoins:      int64(v.BetMultiple * this.TableConfig.BaseScore),
			Coins:         v.WinCoins,
			CoinsAfter:    v.Coins,
			Cards:         v.NiuCards,
			BrandMultiply: v.NiuMultiple,
			BetMultiple:   betMultiple,
			Multiple:      multiple,
			Banker:        int32(isBanker),
			Score:         this.TableConfig.BaseScore,
		}

		rdreq.UserRecord = append(rdreq.UserRecord, rddata)
		v.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
		rdreq.UserRecord = []GGameRecordInfo{}
	}
}

//如果房主退出执行的另一个Leave
func (this *ExtDesk) Leave(p *ExtPlayer) bool {
	//如果是房主退出,通知所以客户端解散
	if p.ChairId == 0 {
		for _, v := range this.Players {
			v.SendNativeMsg(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
				Id: MSG_GAME_INFO_DISMISS_REPLY,
				//DisPlayer: []int32{0, 1},
				IsDismiss: 3,
			})
		}
	}

	if this.GameState == GAME_STATUS_FREE {
		this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Result: 0,
			Token:  p.Token,
			Robot:  p.Robot,
		})
		this.DelPlayer(p.Uid)
		this.DeskMgr.LeaveDo(p.Uid)
	} else if this.GameState == GAME_STATUS_END {
		return true
	} else {
		// p.LiXian = true
		// this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
		// 	Id:     MSG_GAME_LEAVE_REPLY,
		// 	Result: 1,
		// 	Cid:    p.ChairId,
		// 	Uid:    p.Uid,
		// 	Err:    "玩家正在游戏中，不能离开",
		// })
		p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Result: 1,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Err:    "玩家正在游戏中，不能离开",
			Robot:  p.Robot,
		})
		return false
	}
	return true
}

package main

import (
	. "MaJiangTool"
	"logs"
	// "math"
)

//
func (this *ExtDesk) GameHu() {
	for _, v := range this.Players {
		if v.LunHuEd {
			if CheckMenQing(v.FuZis) {
				v.HuType = append(v.HuType, HuType_MenQing)
			}
			if CheckZhongZhang(v.FuZis, v.HandCard, nil) {
				v.HuType = append(v.HuType, HuType_ZhongZhang)
			}
			//庄检查天胡，其他检查地和
			if int(v.ChairId) == this.Banker {
				if CheckTianHu(&this.EventManager) {
					v.HuType = append(v.HuType, HuType_TianHu)
				}
			} else {
				if CheckDiHu(&this.EventManager) {
					v.HuType = append(v.HuType, HuType_DiHu)
				}
			}
			//
			if CheckGangKai(&this.EventManager, int(v.ChairId)) {
				v.HuType = append(v.HuType, HuType_GangKai)
			}
			if CheckGangHouPao(&this.EventManager, int(v.ChairId)) {
				v.HuType = append(v.HuType, HuType_GangHouPao)
			}
			if CheckQiangGangHu(&this.EventManager, int(v.ChairId)) {
				v.HuType = append(v.HuType, HuType_QiangGangHu)
			}
			if CheckZiMo(&this.EventManager, int(v.ChairId)) {
				v.HuType = append(v.HuType, HuType_ZiMo)
			}
			logs.Debug("玩家胡牌：", v.HuType)
		}
	}

}

func (this *ExtDesk) GameOver() {
	//广播给客户端，游戏结束
	logs.Debug("游戏结束----------------------------------------------------")
	this.BroadcastAll(MSG_GAME_INFO_GAMEOVER, &struct {
		Id int
	}{
		Id: MSG_GAME_INFO_GAMEOVER,
	})
	//离开游戏，归还桌子
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
	this.GameOverLeave()
	//开始归还桌子定时器
	this.AddTimer(TIMER_OVER, TIMER_OVER_NUM, this.TimerOver, nil)
}

func (this *ExtDesk) TimerOver(d interface{}) {
	this.GameState = GAME_STATUS_FREE
	this.JuHao = ""
	this.ClearTimer()
	this.EventManager.Events = []EventIe{}
	this.DeskMgr.BackDesk(this)
}

// func (this *ExtDesk) GameOver(player *ExtPlayer) {
// 	this.TList = []*Timer{}
// 	this.GameState = GAME_STATUS_END
// 	this.BroadStageTime(0)
// 	//
// 	overdata := GInfoGameEnd{
// 		Id: MSG_GAME_INFO_END_NOTIFY,
// 	}
// 	if int(player.ChairId) == this.Banker {
// 		overdata.EndType = 0 //地主赢了
// 	} else {
// 		overdata.EndType = 1 //农民赢了
// 	}
// 	overdata.ChunTian = this.IsChunTian()
// 	if overdata.ChunTian == 1 {
// 		for _, v := range this.Players {
// 			v.Double *= 2
// 			if v.ChairId == this.Banker && v.Double > this.MaxDouble*2 {
// 				v.Double = this.MaxDouble * 2
// 			} else if v.ChairId != this.Banker && v.Double > this.MaxDouble {
// 				v.Double = this.MaxDouble
// 			}
// 		}
// 	}
// 	for _, v := range this.Players {
// 		// overdata.Accouts = append(overdata.Accouts, v.Account)
// 		overdata.Double = append(overdata.Double, v.Double)
// 	}
// 	//计算得分
// 	for _, v := range this.Players {
// 		if overdata.EndType == 0 {
// 			if v.ChairId == player.ChairId {
// 				overdata.Scores = append(overdata.Scores, int64(v.Double*this.Bscore))
// 			} else {
// 				overdata.Scores = append(overdata.Scores, int64(-v.Double*this.Bscore))
// 			}
// 		} else {
// 			if v.ChairId == this.Banker {
// 				overdata.Scores = append(overdata.Scores, int64(-v.Double*this.Bscore))
// 			} else {
// 				overdata.Scores = append(overdata.Scores, int64(v.Double*this.Bscore))
// 			}
// 		}
// 	}
// 	//扣除玩家金币
// 	for _, v := range this.Players {
// 		//游戏消费税
// 		if overdata.Scores[v.ChairId] > 0 {
// 			v.RateCoins = float64(overdata.Scores[v.ChairId]) * G_DbGetGameServerData.Rate
// 			overdata.Scores[v.ChairId] = overdata.Scores[v.ChairId] - int64(v.RateCoins)
// 		}
// 		v.Coins += overdata.Scores[v.ChairId]
// 		overdata.Coins = append(overdata.Coins, v.Coins)
// 	}
// 	//剩余手牌 --lgh
// 	var playInfos []GInfoGameEndPlayInfo
// 	for _, p := range this.Players {
// 		if len(p.HandCard) > 0 {

// 			playInfo := GInfoGameEndPlayInfo{
// 				Cid:       p.ChairId,
// 				HandCards: p.HandCard,
// 			}
// 			playInfos = append(playInfos, playInfo)
// 		}
// 	}
// 	overdata.PlayInfos = playInfos
// 	for _, p1 := range this.Players {
// 		overdata.Accouts = []string{}
// 		for _, p2 := range this.Players {
// 			if p1.Uid != p2.Uid && len(p2.Nick) >= 4 {
// 				overdata.Accouts = append(overdata.Accouts, "***"+p2.Nick[len(p2.Nick)-4:])
// 			} else {
// 				overdata.Accouts = append(overdata.Accouts, p2.Nick)
// 			}
// 		}
// 		//
// 		p1.SendNativeMsg(MSG_GAME_INFO_END_NOTIFY, &overdata)
// 	}
// 	//
// 	// this.BroadcastAll(MSG_GAME_INFO_END_NOTIFY, &overdata)
// 	//发送结算消息给数据库
// 	dbreq := GGameEnd{
// 		Id:          MSG_GAME_END_NOTIFY,
// 		GameId:      GCONFIG.GameType,
// 		GradeId:     GCONFIG.GradeType,
// 		RoomId:      GCONFIG.RoomType,
// 		GameRoundNo: this.JuHao,
// 		Mini:        false,
// 		ActiveUid:   player.Uid,
// 	}
// 	for _, v := range this.Players {
// 		valid := overdata.Scores[v.ChairId]
// 		if valid < 0 {
// 			valid = -valid
// 		} else {
// 			valid = int64(this.Bscore)
// 		}
// 		dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
// 			UserId:      v.Uid,
// 			UserAccount: v.Account,
// 			BetCoins:    int64(this.Bscore),
// 			ValidBet:    valid,
// 			PrizeCoins:  overdata.Scores[v.ChairId],
// 			Robot:       v.Robot,
// 			WaterProfit: v.RateCoins,
// 			WaterRate:   G_DbGetGameServerData.Rate,
// 		})
// 		v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
// 		dbreq.UserCoin = []GGameEndInfo{}
// 	}
// 	// player.SendNativeMsg(MSG_GAME_END_NOTIFY, &dbreq)
// 	//发送消息给大厅去记录游戏记录
// 	rdreq := GGameRecord{
// 		Id:          MSG_GAME_END_RECORD,
// 		GameId:      GCONFIG.GameType,
// 		GradeId:     GCONFIG.GradeType,
// 		RoomId:      GCONFIG.RoomType,
// 		GameRoundNo: this.JuHao,
// 	}
// 	for _, v := range this.Players {
// 		if v.Robot {
// 			continue
// 		}
// 		rddata := GGameRecordInfo{
// 			UserId:      v.Uid,
// 			UserAccount: v.Account,
// 			Coins:       overdata.Scores[v.ChairId],
// 			Score:       int(this.Bscore),
// 			Multiple:    int(v.Double),
// 			CoinsBefore: v.Coins - overdata.Scores[v.ChairId],
// 			CoinsAfter:  v.Coins,
// 			Robot:       v.Robot,
// 		}
// 		if v.ChairId == this.Banker {
// 			rddata.Landlord = true
// 		} else {
// 			rddata.Landlord = false
// 		}
// 		rdreq.UserRecord = append(rdreq.UserRecord, rddata)
// 		v.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
// 		rdreq.UserRecord = []GGameRecordInfo{}
// 	}
// 	// player.SendNativeMsg(MSG_GAME_END_RECORD, &rdreq)
// 	//
// 	this.GameOverLeave()
// 	//开始归还桌子定时器
// 	this.AddTimer(TIMER_OVER, TIMER_OVER_NUM, this.TimerOver, nil)
// }

// func (this *ExtDesk) GameOverByNoCall(player *ExtPlayer) {
// 	this.TList = []*Timer{}
// 	this.GameState = GAME_STATUS_END
// 	this.CallFen = 0
// 	this.MaxChuPai = nil
// 	this.RdChuPai = []*GOutCard{}
// 	this.CallTimes = 0
// 	//
// 	overdata := GInfoGameEnd{
// 		Id: MSG_GAME_INFO_END_NOTIFY,
// 	}
// 	//扣除玩家金币
// 	for _, v := range this.Players {
// 		overdata.Coins = append(overdata.Coins, v.Coins)
// 		overdata.Accouts = append(overdata.Accouts, v.Account)
// 		overdata.Double = append(overdata.Double, 0)
// 		overdata.Scores = append(overdata.Scores, 0)
// 	}

// 	//剩余手牌 --lgh
// 	var playInfos []GInfoGameEndPlayInfo
// 	for _, p := range this.Players {
// 		if len(p.HandCard) > 0 {

// 			playInfo := GInfoGameEndPlayInfo{
// 				Cid:       p.ChairId,
// 				HandCards: p.HandCard,
// 			}
// 			playInfos = append(playInfos, playInfo)
// 		}
// 	}
// 	overdata.PlayInfos = playInfos
// 	//
// 	for _, p1 := range this.Players {
// 		overdata.Accouts = []string{}
// 		for _, p2 := range this.Players {
// 			if p1.Uid != p2.Uid && len(p2.Nick) >= 4 {
// 				overdata.Accouts = append(overdata.Accouts, "***"+p2.Nick[len(p2.Nick)-4:])
// 			} else {
// 				overdata.Accouts = append(overdata.Accouts, p2.Nick)
// 			}
// 		}
// 		//
// 		p1.SendNativeMsg(MSG_GAME_INFO_END_NOTIFY, &overdata)
// 		// p1.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
// 		// 	Id:     MSG_GAME_LEAVE_REPLY,
// 		// 	Result: 0,
// 		// 	Cid:    p1.ChairId,
// 		// 	Uid:    p1.Uid,
// 		// 	Token:  p1.Token,
// 		// 	Robot:  p1.Robot,
// 		// })
// 	}
// 	//
// 	dbreq := GGameEnd{
// 		Id:          MSG_GAME_END_NOTIFY,
// 		GameId:      GCONFIG.GameType,
// 		GradeId:     GCONFIG.GradeType,
// 		RoomId:      GCONFIG.RoomType,
// 		GameRoundNo: this.JuHao,
// 		Mini:        false,
// 		ActiveUid:   player.Uid,
// 	}
// 	for _, v := range this.Players {
// 		valid := overdata.Scores[v.ChairId]
// 		if valid < 0 {
// 			valid = -valid
// 		} else {
// 			valid = int64(this.Bscore)
// 		}
// 		dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
// 			UserId:      v.Uid,
// 			UserAccount: v.Account,
// 			BetCoins:    int64(this.Bscore),
// 			ValidBet:    valid,
// 			PrizeCoins:  overdata.Scores[v.ChairId],
// 			Robot:       v.Robot,
// 			WaterProfit: v.RateCoins,
// 			WaterRate:   G_DbGetGameServerData.Rate,
// 		})
// 		v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
// 		dbreq.UserCoin = []GGameEndInfo{}
// 	}
// 	// this.BroadcastAll(MSG_GAME_INFO_END_NOTIFY, &overdata)
// 	//发送结算消息给数据库
// 	// dbreq := GGameEnd{
// 	// 	Id:          MSG_GAME_END_NOTIFY,
// 	// 	GameId:      GCONFIG.GameType,
// 	// 	GradeId:     GCONFIG.GradeType,
// 	// 	RoomId:      GCONFIG.RoomType,
// 	// 	GameRoundNo: this.JuHao,
// 	// 	Mini:        false,
// 	// }
// 	// for _, v := range this.Players {
// 	// 	valid := overdata.Scores[v.ChairId]
// 	// 	if valid < 0 {
// 	// 		valid = -valid
// 	// 	}

// 	// 	dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
// 	// 		UserId:      v.Uid,
// 	// 		UserAccount: v.Account,
// 	// 		BetCoins:    0,
// 	// 		ValidBet:    valid,
// 	// 		PrizeCoins:  overdata.Scores[v.ChairId],
// 	// 		Robot:       v.Robot,
// 	// 	})
// 	// 	v.SendNativeMsg(MSG_GAME_END_NOTIFY, &dbreq)
// 	// 	dbreq.UserCoin = []GGameEndInfo{}
// 	// }
// 	// // player.SendNativeMsg(MSG_GAME_END_NOTIFY, &dbreq)
// 	// //发送消息给大厅去记录游戏记录
// 	// rdreq := GGameRecord{
// 	// 	Id:          MSG_GAME_END_RECORD,
// 	// 	GameId:      GCONFIG.GameType,
// 	// 	GradeId:     GCONFIG.GradeType,
// 	// 	RoomId:      GCONFIG.RoomType,
// 	// 	GameRoundNo: this.JuHao,
// 	// }
// 	// for _, v := range this.Players {
// 	// 	rddata := GGameRecordInfo{
// 	// 		UserId:      v.Uid,
// 	// 		UserAccount: v.Account,
// 	// 		Coins:       overdata.Scores[v.ChairId],
// 	// 		Score:       int(this.Bscore),
// 	// 		Multiple:    int(v.Double),
// 	// 		Robot:       v.Robot,
// 	// 	}
// 	// 	if v.ChairId == this.Banker {
// 	// 		rddata.Landlord = true
// 	// 	} else {
// 	// 		rddata.Landlord = false
// 	// 	}
// 	// 	rdreq.UserRecord = append(rdreq.UserRecord, rddata)
// 	// }
// 	// player.SendNativeMsg(MSG_GAME_END_RECORD, &rdreq)
// 	//
// 	this.GameOverLeave()
// 	//开始归还桌子定时器
// 	this.AddTimer(TIMER_OVER, TIMER_OVER_NUM, this.TimerOver, nil)
// }

// func (this *ExtDesk) TimerOver(d interface{}) {
// 	this.GameState = GAME_STATUS_FREE
// 	this.JuHao = ""
// 	this.DeskMgr.BackDesk(this)
// }

// func (this *ExtDesk) IsChunTian() int32 {
// 	//农民春天
// 	if len(this.Players[this.Banker].Outed) == 1 {
// 		return 1
// 	}
// 	//
// 	for _, v := range this.Players {
// 		if v.ChairId == this.Banker {
// 			continue
// 		}
// 		if len(v.Outed) != 0 {
// 			return 0
// 		}
// 	}
// 	return 1
// }

package main

import (
	"fmt"
	"logs"
	"time"
)

// import (
// 	"logs"
// 	// "math"
// )

func (this *ExtDesk) GameOver(player *ExtPlayer) {
	BoomCout = 0
	for _, v := range this.Players {
		fmt.Println("玩家:", v.Nick, "的倍数是:", v.Double)
	}
	// for _, v := range this.Players {
	// 	fmt.Println("添加记录之前", v.Nick, "的记录：", v.GGameLogs)
	// }
	logs.Debug("go bynormal")
	this.TList = []*Timer{}
	// this.GameState = GAME_STATUS_END
	// this.BroadStageTime(0)
	this.CallFen = 0
	this.MaxChuPai = nil
	this.RdChuPai = []*GOutCard{}
	this.CallTimes = 0
	//
	logs.Debug("游戏结束")
	overdata := GInfoGameEnd{
		Id: MSG_GAME_INFO_END_NOTIFY,
	}
	if player.ChairId == this.Banker {
		overdata.EndType = 0 //地主赢了
	} else {
		overdata.EndType = 1 //农民赢了
	}
	overdata.ChunTian = this.IsChunTian()
	fmt.Println("是否春天！！！！！！！！！！！？：：：", overdata.ChunTian)
	if overdata.ChunTian == 1 {
		for _, v := range this.Players {
			v.Double *= 2
			fmt.Println("进入春天短暂翻倍,Name:", v.Nick, "DOUBLE:", v.Double)
			if v.ChairId == this.Banker && v.Double > this.MaxDouble*2 {
				v.Double = this.MaxDouble * 2
			} else if v.ChairId != this.Banker && v.Double > this.MaxDouble {
				v.Double = this.MaxDouble
			}
		}
	}
	fmt.Println("封顶倍率：", this.MaxDouble)
	if this.DiPaiDoulbe > 1 {
		for _, v := range this.Players {
			v.Double *= this.DiPaiDoulbe
			if v.ChairId == this.Banker && v.Double > this.MaxDouble*2 {
				v.Double = this.MaxDouble * 2
			} else if v.ChairId != this.Banker && v.Double > this.MaxDouble {
				v.Double = this.MaxDouble
			}
		}
	}
	for _, v := range this.Players {
		// overdata.Accouts = append(overdata.Accouts, v.Account)
		overdata.Double = append(overdata.Double, v.Double)
	}
	//计算得分
	for _, v := range this.Players {

		if overdata.EndType == 0 { //如果地主赢了
			if v.ChairId == player.ChairId {
				overdata.Scores = append(overdata.Scores, int64(v.Double*this.Bscore))
			} else {
				overdata.Scores = append(overdata.Scores, int64(-v.Double*this.Bscore))
			}
		} else {
			if v.ChairId == this.Banker {
				overdata.Scores = append(overdata.Scores, int64(-v.Double*this.Bscore))
			} else {
				overdata.Scores = append(overdata.Scores, int64(v.Double*this.Bscore))
			}
		}

	}

	//扣除玩家金币
	for _, v := range this.Players {
		v.Coins += overdata.Scores[v.ChairId]
		v.LogCoins += overdata.Scores[v.ChairId]
		overdata.Coins = append(overdata.Coins, v.Coins)
		v.GGameLogs.Id = MSG_GAME_INFO_QPLAYERLOGS_REPLY
		v.GGameLogs.GGameLogs = append(v.GGameLogs.GGameLogs, GGameLog{
			EndTime: time.Now().Format("2006-01-02 15:04:05"),
			Coins:   overdata.Scores[v.ChairId],
		})
		fmt.Println("现在是", v.Nick, "在存游戏记录：", v.GGameLogs)
	}
	//剩余手牌 --lgh
	var playInfos []GInfoGameEndPlayInfo
	for _, p := range this.Players {
		if len(p.HandCard) > 0 {
			playInfo := GInfoGameEndPlayInfo{
				Cid: p.ChairId,
			}
			hdc := Sort(p.HandCard)
			for _, v := range hdc {
				playInfo.HandCards = append(playInfo.HandCards, int(v))
			}
			playInfos = append(playInfos, playInfo)
		}
	}
	overdata.PlayInfos = playInfos
	for _, p1 := range this.Players {
		overdata.Accouts = []string{}
		for _, p2 := range this.Players {
			if p1.Uid != p2.Uid && len(p2.Nick) >= 4 {
				overdata.Accouts = append(overdata.Accouts, "***"+p2.Nick[len(p2.Nick)-4:])
			} else {
				overdata.Accouts = append(overdata.Accouts, p2.Nick)
			}
		}
		//
		p1.SendNativeMsg(MSG_GAME_INFO_END_NOTIFY, &overdata)
		fmt.Println("正常结算倍数！！！！！！！！！！！！！！！:", overdata.Double)
	}
	//
	// this.BroadcastAll(MSG_GAME_INFO_END_NOTIFY, &overdata)
	//发送结算消息给数据库
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Mini:        false,
		SetLeave:    1,
		NoSaveCoin:  1,
		RoomNo:      this.FkNo,
		Round:       this.Round,
	}
	for _, v := range this.Players {
		valid := overdata.Scores[v.ChairId]
		if valid < 0 {
			valid = -valid
		} else {
			valid = int64(this.Bscore)
		}
		dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
			UserId:      v.Uid,
			UserAccount: v.Account,
			BetCoins:    int64(this.Bscore),
			ValidBet:    valid,
			PrizeCoins:  overdata.Scores[v.ChairId] * 100,
			Robot:       v.Robot,
			WaterRate:   G_DbGetGameServerData.Rate,
		})
		v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
		dbreq.UserCoin = []GGameEndInfo{}
	}
	// player.SendNativeMsg(MSG_GAME_END_NOTIFY, &dbreq)
	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameModule:  1,
		PayType:     1,
		GameType:    1,
		GameRoundNo: this.JuHao,
		RoomNo:      this.FkNo,
		Round:       this.Round,
	}
	for _, v := range this.Players {
		if v.Robot {
			continue
		}
		rddata := GGameRecordInfo{
			UserId:      v.Uid,
			UserAccount: v.Account,
			Coins:       overdata.Scores[v.ChairId],
			Score:       int(this.Bscore),
			Multiple:    int(v.Double),
			CoinsBefore: v.Coins - overdata.Scores[v.ChairId],
			CoinsAfter:  v.Coins,
			Robot:       v.Robot,
		}
		if v.ChairId == this.Banker {
			rddata.Landlord = true
		} else {
			rddata.Landlord = false
		}
		rdreq.UserRecord = append(rdreq.UserRecord, rddata)
		v.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
		rdreq.UserRecord = []GGameRecordInfo{}
	}
	// player.SendNativeMsg(MSG_GAME_END_RECORD, &rdreq)

	//发送游戏记录
	for _, v := range this.Players {
		fmt.Println(v.Nick, "的游戏记录：", v.GGameLogs)
		v.SendNativeMsg(MSG_GAME_INFO_QPLAYERLOGS_REPLY, v.GGameLogs)
	}
	if this.Round >= this.TableConfig.MatchCount {
		this.allB()
	} else {
		//进行下一把
		fmt.Println("还剩余", this.TableConfig.MatchCount-this.Round, "局游戏")
		for _, v := range this.Players { //初始化所有玩家
			v.ResetPlayer()
		}
		//清楚桌子数据
		this.ResetExtDest()
		//再次进入准备状态
		this.GameState = GAME_STATUS_READ
		this.BroadStageTime(TIMER_READ_NUM)
	}
}

func (this *ExtDesk) allB() {
	logs.Debug("游戏结束", this.Players)
	this.ClearTimer()

	this.GameState = GAME_STATUS_END
	this.BroadStageTime(0)
	var allc []int64 = []int64{0, 0, 0}
	//玩家总结算
	for i, v := range this.Players {
		allc[i] += v.LogCoins
	}
	this.BroadcastAll(MSG_GAME_INFO_ALLBALANCE, AllBalance{
		Id:       MSG_GAME_INFO_ALLBALANCE,
		AllConis: allc,
	})

	//玩家离开
	for _, p := range this.Players {
		p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:      MSG_GAME_LEAVE_REPLY,
			Result:  0,
			Cid:     p.ChairId,
			Uid:     p.Uid,
			Token:   p.Token,
			NoToCli: true,
		})
	}
	this.GameOverLeave()
	//归还桌子
	this.GameState = GAME_STATUS_FREE
	this.ResetTable()
	this.DeskMgr.BackDesk(this)
}
func (this *ExtDesk) GameOverByNoCall(player *ExtPlayer) {
	logs.Debug("go bynocall")
	this.TList = []*Timer{}
	// this.GameState = GAME_STATUS_END
	this.CallFen = 0
	this.MaxChuPai = nil
	this.RdChuPai = []*GOutCard{}
	this.CallTimes = 0
	//
	overdata := GInfoGameEnd{
		Id: MSG_GAME_INFO_END_NOTIFY,
	}
	//扣除玩家金币
	for _, v := range this.Players {
		overdata.Coins = append(overdata.Coins, v.Coins)
		overdata.Accouts = append(overdata.Accouts, v.Account)
		overdata.Double = append(overdata.Double, 0)
		overdata.Scores = append(overdata.Scores, 0)
		v.GGameLogs.Id = MSG_GAME_INFO_QPLAYERLOGS_REPLY
		v.GGameLogs.GGameLogs = append(v.GGameLogs.GGameLogs, GGameLog{
			EndTime: time.Now().Format("2006-01-02 15:04:05"),
			Coins:   0,
		})
	}
	//剩余手牌 --lgh
	var playInfos []GInfoGameEndPlayInfo
	for _, p := range this.Players {
		if len(p.HandCard) > 0 {

			playInfo := GInfoGameEndPlayInfo{
				Cid: p.ChairId,
			}
			hdc := Sort(p.HandCard)
			for _, v := range hdc {
				playInfo.HandCards = append(playInfo.HandCards, int(v))
			}
			playInfos = append(playInfos, playInfo)
		}
	}
	overdata.PlayInfos = playInfos
	//
	for _, p1 := range this.Players {
		overdata.Accouts = []string{}
		for _, p2 := range this.Players {
			if p1.Uid != p2.Uid && len(p2.Nick) >= 4 {
				overdata.Accouts = append(overdata.Accouts, "***"+p2.Nick[len(p2.Nick)-4:])
			} else {
				overdata.Accouts = append(overdata.Accouts, p2.Nick)
			}
		}
		p1.SendNativeMsg(MSG_GAME_INFO_END_NOTIFY, &overdata)
		fmt.Println("不正常结算倍数！！！！！！！！！！！！！！！:", overdata.Double)
	}
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Mini:        false,
		SetLeave:    1,
		RoomNo:      this.FkNo,
		Round:       this.Round,
	}
	for _, v := range this.Players {
		valid := overdata.Scores[v.ChairId]
		if valid < 0 {
			valid = -valid
		} else {
			valid = int64(this.Bscore)
		}
		dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
			UserId:      v.Uid,
			UserAccount: v.Account,
			BetCoins:    int64(this.Bscore),
			ValidBet:    valid,
			PrizeCoins:  overdata.Scores[v.ChairId],
			Robot:       v.Robot,
		})
		dbreq.UserCoin = []GGameEndInfo{}
	}
	//发送游戏记录	v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
	for _, v := range this.Players {
		fmt.Println(v.Nick, "的游戏记录：", v.GGameLogs)
		v.SendNativeMsg(MSG_GAME_INFO_QPLAYERLOGS_REPLY, v.GGameLogs)
	}
	if this.Round >= this.TableConfig.MatchCount {
		this.allB()
	} else {
		//进行下一把
		fmt.Println("还剩余", this.TableConfig.MatchCount-this.Round, "局游戏")
		for _, v := range this.Players { //初始化所有玩家
			v.ResetPlayer()
		}
		//清楚桌子数据
		this.ResetExtDest()
		//再次进入准备状态
		this.GameState = GAME_STATUS_READ
		this.BroadStageTime(TIMER_READ_NUM)
	}
}

func (this *ExtDesk) TimerOver(d interface{}) {
	this.GameState = GAME_STATUS_FREE
	this.JuHao = ""
	this.DeskMgr.BackDesk(this)
}

func (this *ExtDesk) IsChunTian() int32 {
	//农民春天
	if len(this.Players[this.Banker].Outed) == 1 {
		return 1
	}
	//
	for _, v := range this.Players {
		if v.ChairId == this.Banker {
			continue
		}
		if len(v.Outed) != 0 {
			return 0
		}
	}
	return 1
}

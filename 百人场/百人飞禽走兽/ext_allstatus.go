package main

import (
	"time"
	// "encoding/json"
	"logs"
)

func (this *ExtDesk) status_ready(d interface{}) {
	logs.Debug("准备阶段")
	this.GameState = GAME_STATUS_READY
	this.BroStatusTime(GAME_STATUS_READY)
	//准备阶段多用于客户端与服务端清除数据,改变局号等等
	this.JuHao = GetJuHao()
	if this.Zhuang.Uid != 0 {
		//判断是否进入风控
		this.GameResult = this.randresult()
		if CD-CalPkAll(StartControlTime, time.Now().Unix()) < 0 {
			logs.Debug("进入风控,当前有玩家是庄，通知机器人开奖结果")
			for _, v := range this.Players {
				if v.Robot {
					v.SendNativeMsg(MSG_GAME_INFO_TOROBOTRESULT, ToRobot{
						Id:    MSG_GAME_INFO_TOROBOTRESULT,
						Index: this.GameResult,
					})
				}
			}
		}
	}
	//清除数据
	for _, v := range this.Players {
		v.Rest()
	}
	this.BroadcastAll(MSG_GAME_INFO_JUHAOCHANGE, struct {
		Id    int
		JuHao string
	}{
		Id:    MSG_GAME_INFO_JUHAOCHANGE,
		JuHao: this.JuHao,
	})
	this.AddTimer(GAME_STATUS_READY, GAME_STATUS_READY_TIME/1000, this.status_startdownbet, nil)
}

//开始下注阶段
func (this *ExtDesk) status_startdownbet(d interface{}) {
	logs.Debug("开始下注阶段")
	this.GameState = GAME_STATUS_STARTBET
	this.BroStatusTime(GAME_STATUS_STARTBET_TIME)
	//洗牌
	this.AddTimer(999, 1, this.BroDownBetInfo, nil) //开始每隔一秒更新一次下注情况
	this.AddTimer(GAME_STATUS_STARTBET, GAME_STATUS_STARTBET_TIME/1000, this.status_downbet, nil)
}

//下注阶段
func (this *ExtDesk) status_downbet(d interface{}) {
	logs.Debug("下注阶段")
	this.GameState = GAME_STATUS_DOWNBET
	this.BroStatusTime(GAME_STATUS_DOWNBET_TIME)
	//下注
	this.AddTimer(GAME_STATUS_DOWNBET, GAME_STATUS_DOWNBET_TIME/1000, this.status_stopdownbet, nil)
}

//停止下注阶段
func (this *ExtDesk) status_stopdownbet(d interface{}) {
	logs.Debug("下注结束")
	this.GameState = GAME_STATUS_ENDBET
	this.BroStatusTime(GAME_STATUS_ENDBET_TIME)
	//下注
	this.AddTimer(GAME_STATUS_ENDBET, GAME_STATUS_ENDBET_TIME/1000, this.status_lottery, nil)
}

//开牌阶段
func (this *ExtDesk) status_lottery(d interface{}) {
	logs.Debug("开奖阶段")
	this.GameState = GAME_STATUS_LOTTERY
	this.BroStatusTime(GAME_STATUS_LOTTERY_TIME)
	//获取真实玩家下注信息
	bet := this.getPlayerBet()
	if this.Zhuang.Uid == 0 {
		//判断是否进入风控(无玩家上庄时下注阶段)
		if CD-CalPkAll(StartControlTime, time.Now().Unix()) < 0 && this.Zhuang.Uid == 0 && bet > 0 {
			logs.Debug("进入风控")
			this.GameResult = this.ControlWinOrLose(true)
		} else if bet <= 0 {
			logs.Debug("不进入风控  75")
			ra := RandInt64(100)
			if ra < 75 {
				this.GameResult = this.ControlWinOrLose(false)
			} else {
				this.GameResult = this.randresult()
			}
		} else {
			this.GameResult = this.randresult()
		}
	}
	this.addHistory()
	//发送开奖通知
	this.BroadcastAll(MSG_GAME_INFO_LOTTERY_BRO, &struct {
		Id     int
		Result int
	}{
		Id:     MSG_GAME_INFO_LOTTERY_BRO,
		Result: this.GameResult,
	})
	this.AddTimer(GAME_STATUS_LOTTERY, GAME_STATUS_LOTTERY_TIME/1000, this.status_balance, nil)
}

//结算阶段
func (this *ExtDesk) status_balance(d interface{}) {
	logs.Debug("结算阶段")
	this.GameState = GAME_STATUS_BALANCE
	this.BroStatusTime(GAME_STATUS_BALANCE_TIME)
	//结算阶段 用于用户金币结算，用户游戏记录存储
	if this.Zhuang.Uid != 0 {
		logs.Debug("发现有玩家上庄情况产生")
		this.BalanceWhenZhuang()
	} else {
		logs.Debug("本局未发现玩家上庄使用普通结算")
		this.BalanceWhenNoZhuang()
	}
	//获取庄家真实输赢
	this.getDeskGetcoins()
	if GetCostType() == 1 {
		AddCD(this.ZhenShisGetcoins)
		AddLocalStock(this.ZhenShisGetcoins)
	}
	this.needDowhenBalance()
	//增加在庄局数
	if this.Zhuang.Uid != 0 {
		this.ZhuangMatchCount += 1
	}
	this.AddTimer(GAME_STATUS_LOTTERY, GAME_STATUS_LOTTERY_TIME/1000, this.status_ready, nil)
}

func (this *ExtDesk) needDowhenBalance() {
	nobetwarning := gameConfig.LimitInfo.NobetWarning
	nobetremove := gameConfig.LimitInfo.NobetRemove
	for _, v := range this.Players {
		//未下注提示或退出
		if v.NoToBet >= nobetwarning {
			if v.NoToBet >= nobetremove {
				v.SendNativeMsg(MSG_GAME_INFO_WARNING, MsgWarning{
					Id:     MSG_GAME_INFO_WARNING,
					Result: 2,
				})
				v.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
					Id:      MSG_GAME_LEAVE_REPLY,
					Result:  0,
					Cid:     v.ChairId,
					Uid:     v.Uid,
					Token:   v.Token,
					Robot:   v.Robot,
					NoToCli: true,
				})
				this.UpChair(v)
				this.BroChairChange()
				this.DelPlayer(v.Uid)
				this.DeskMgr.LeaveDo(v.Uid)
			} else if v.NoToBet == gameConfig.LimitInfo.NobetWarning {
				v.SendNativeMsg(MSG_GAME_INFO_WARNING, MsgWarning{
					Id:     MSG_GAME_INFO_WARNING,
					Result: 1,
				})
			}
		}
		if v.LiXian {
			this.UpChair(v)
			this.BroChairChange()
			this.LeaveByForce(v)
		}

	}
	//每次结算之后，更新队列中的值
	for _, v := range this.Players {
		for _, v1 := range this.WaitZhuang {
			if v.Uid == v1.Uid {
				if v.Coins < gameConfig.LimitInfo.UpZhuangNeed {
					logs.Debug("发现等待队列中的玩家金币变得不足 ，将其移除等待上庄队列")
					v.SendNativeMsg(MSG_GAME_INFO_WARNING, MsgWarning{
						Id:     MSG_GAME_INFO_WARNING,
						Result: 6,
					})
					this.RemoveWait(v)
				} else {
					v1.MyUserCoin = v.Coins
				}
			}
		}
	}
	this.ChangeZhuang()
	if this.Zhuang.Uid == 0 {
		return
	}
	//查找庄
	var zhuang *ExtPlayer
	for _, v := range this.Players {
		if v.Uid == this.Zhuang.Uid {
			zhuang = v
		}
	}
	//查找所有离线玩家，如果离线玩家中存在庄，便下庄，如果 存在队列中，则移除队列
	lixianList := []*ExtPlayer{}
	for _, v := range this.Players {
		if v.LiXian {
			lixianList = append(lixianList, v)
		}
	}
	for _, v := range lixianList {
		if v.Uid == zhuang.Uid {
			this.DownZhuang(v)
		}
	}
	for _, v := range lixianList {
		for _, v1 := range this.WaitZhuang {
			if v.Uid == v1.Uid {
				this.RemoveWait(v)
			}
		}
	}
	//判断当前庄金币是否足够
	if zhuang.Uid < gameConfig.LimitInfo.UpZhuangNeed {
		logs.Debug("当前庄家金币不够！！！")
		zhuang.SendNativeMsg(MSG_GAME_INFO_WARNING, MsgWarning{
			Id:     MSG_GAME_INFO_WARNING,
			Result: 4,
		})
		this.DownZhuang(zhuang)
	}
	//判断玩家当庄局数是否到达限制
	if this.ZhuangMatchCount > gameConfig.LimitInfo.ChangeZhuangCount {
		logs.Debug("当前庄家局数已经满了，换装")
		zhuang.SendNativeMsg(MSG_GAME_INFO_WARNING, MsgWarning{
			Id:     MSG_GAME_INFO_WARNING,
			Result: 5,
		})
		this.DownZhuang(zhuang)
	}
	//查找金币少于入场限制的玩家
	for _, v := range this.Players {
		if v.Coins < int64(G_DbGetGameServerData.LimitLower) {
			logs.Debug("发现玩家小于入场限制 将其踢出", v.Nick)
			this.LeaveByForce(v)
		}
	}
}

//查找真实玩家下注
func (this *ExtDesk) getPlayerBet() (bet int64) {
	for _, v := range this.Players {
		if !v.Robot {
			for _, v1 := range v.DownBet {
				bet += v1
			}
		}
	}
	return
}

//获取我们真实输赢
func (this *ExtDesk) getDeskGetcoins() (deskget int64) {
	for i, v := range this.ZhenShiDownBet {
		if i == this.GameResult {
			this.ZhenShisGetcoins -= v*(int64(LotteryDouble[i])*10)/10 - v
		} else {
			this.ZhenShisGetcoins += v
		}
	}
	return
}

//有庄家结算细节
func (this *ExtDesk) BalanceWhenZhuang() {
	//发送结算消息给数据库,简单记录
	dbre := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Mini:        false,
		SetLeave:    1,
	}
	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
	}
	var zhuang *ExtPlayer
	//找出庄家
	for _, v := range this.Players {
		if v.Uid == this.Zhuang.Uid {
			zhuang = v
		}
	}
	var zhuangGetCoins int64
	var choushui int64
	if zhuangGetCoins > 0 {
		choushui = zhuangGetCoins * (int64(G_DbGetGameServerData.Rate) * 100) / 100
	}
	zhuangGetCoins = zhuangGetCoins - choushui
	for _, v := range this.Players {
		if v.Uid != this.Zhuang.Uid {
			var betCoins int64 //获取玩家下注金币
			for i, v1 := range v.DownBet {
				betCoins += v1 //玩家下注金币
				if i == this.GameResult {
					v.BalanceResult[i] += v1 * (int64(LotteryDouble[i]) * 10) / 10
					v.GetCoins += v1 * (int64(LotteryDouble[i]) * 10) / 10
					this.ZhuangBalanceResult[i] -= v.GetCoins - v1
					zhuangGetCoins -= v.GetCoins - v1
					v.OtherBalanceResult[i] += v.OtherDownBet[i] * (int64(LotteryDouble[i]) * 10) / 10
				} else {
					v.GetCoins -= v1
					v.BalanceResult[i] = -v1
					this.ZhuangBalanceResult[i] = +v1
					zhuangGetCoins += v1
					v.OtherBalanceResult[i] -= v.OtherDownBet[i] * (int64(LotteryDouble[i]) * 10) / 10
				}
			}
			if betCoins > 0 && !v.Robot {
				logs.Debug("该玩家本局有下注，所以向数据库发送游戏记录!")
				gameEndInfo := GGameEndInfo{
					UserId:      v.Uid,
					UserAccount: v.Nick,
					BetCoins:    betCoins,
					ValidBet:    betCoins,
					PrizeCoins:  v.GetCoins,
					Robot:       v.Robot,
					WaterRate:   G_DbGetGameServerData.Rate,
				}
				dbre.UserCoin = append(dbre.UserCoin, gameEndInfo)
				v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, dbre)
				recordInfo := GGameRecordInfo{
					UserId:      v.Uid,
					UserAccount: v.Nick,
					Robot:       v.Robot,
					BetCoins:    betCoins,
					BetArea:     v.DownBet,
					PrizeCoins:  v.GetCoins,
					CoinsBefore: v.Coins + betCoins,
					CoinsAfter:  v.Coins + v.GetCoins + betCoins,
					Banker:      false,
				}
				rdreq.UserRecord = append(rdreq.UserRecord, recordInfo)
				v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, rdreq)
				dbre.UserCoin = []GGameEndInfo{} //清空
				rdreq.UserRecord = []GGameRecordInfo{}
			}
			if v.GetCoins > 0 {
				v.Coins += v.GetCoins
			}
		}
	}
	//单独结算坐庄玩家的数据
	gameEndInfo := GGameEndInfo{
		UserId:      zhuang.Uid,
		UserAccount: zhuang.Nick,
		BetCoins:    0,
		ValidBet:    0,
		PrizeCoins:  zhuangGetCoins,
		Robot:       zhuang.Robot,
		WaterRate:   G_DbGetGameServerData.Rate,
		WaterProfit: float64(choushui),
	}
	dbre.UserCoin = append(dbre.UserCoin, gameEndInfo)
	zhuang.SendNativeMsgForce(MSG_GAME_END_NOTIFY, dbre)
	recordInfo := GGameRecordInfo{
		UserId:      zhuang.Uid,
		UserAccount: zhuang.Nick,
		Robot:       zhuang.Robot,
		BetCoins:    0,
		BetArea:     zhuang.DownBet,
		PrizeCoins:  zhuangGetCoins,
		CoinsBefore: zhuang.Coins,
		CoinsAfter:  zhuang.Coins + zhuangGetCoins,
		Banker:      true,
	}
	zhuang.Coins += zhuangGetCoins
	rdreq.UserRecord = append(rdreq.UserRecord, recordInfo)
	zhuang.SendNativeMsgForce(MSG_GAME_END_NOTIFY, rdreq)
	//发送结算给客户端
	for _, v := range this.Players {
		if v.Uid != this.Zhuang.Uid {
			toc := ToClientBalance{
				Id:              MSG_GAME_INFO_BALANCE,
				MyCoins:         v.Coins,
				MyResult:        v.BalanceResult,
				AllResult:       this.BalanceResult,
				OtherResult:     v.OtherBalanceResult,
				MyGetCoins:      v.GetCoins,
				History:         this.GameResultHistory,
				BetAbleIndex:    this.CanUseChip(v),
				AreaWinDouble:   int(LotteryDouble[this.GameResult]),
				WinArea:         this.GameResult,
				IsZhuang:        false,
				ZhuangWinOrLose: zhuangGetCoins,
				ZhuangCoins:     zhuang.Coins,
				IsHasZhuang:     true,
			}
			v.SendNativeMsg(MSG_GAME_INFO_BALANCE, toc)
		} else {
			toc := ToClientBalance{
				Id:              MSG_GAME_INFO_BALANCE,
				MyCoins:         v.Coins,
				MyResult:        this.ZhuangBalanceResult,
				AllResult:       this.BalanceResult,
				OtherResult:     v.OtherBalanceResult,
				MyGetCoins:      zhuangGetCoins,
				History:         this.GameResultHistory,
				BetAbleIndex:    this.CanUseChip(v),
				AreaWinDouble:   int(LotteryDouble[this.GameResult]),
				WinArea:         this.GameResult,
				IsZhuang:        true,
				ZhuangWinOrLose: zhuangGetCoins,
				ZhuangCoins:     zhuang.Coins,
				IsHasZhuang:     true,
			}
			v.SendNativeMsg(MSG_GAME_INFO_BALANCE, toc)
		}
	}
}

//无庄家结算细节
func (this *ExtDesk) BalanceWhenNoZhuang() {
	//发送结算消息给数据库,简单记录
	dbre := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Mini:        false,
		SetLeave:    1,
	}
	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
	}

	for _, v := range this.Players {
		var betCoins int64 //获取玩家下注金币
		for i, v1 := range v.DownBet {
			betCoins += v1 //玩家下注金币
			if i == this.GameResult {
				v.BalanceResult[i] += v1 * (int64(LotteryDouble[i]) * 10) / 10
				v.GetCoins += v1 * (int64(LotteryDouble[i]) * 10) / 10
				v.OtherBalanceResult[i] += v.OtherDownBet[i] * (int64(LotteryDouble[i]) * 10) / 10
			} else {
				v.GetCoins -= v1
				v.BalanceResult[i] = -v1
				v.OtherBalanceResult[i] -= v.OtherDownBet[i] * (int64(LotteryDouble[i]) * 10) / 10
			}
		}
		if betCoins > 0 && !v.Robot {
			logs.Debug("该玩家本局有下注，所以向数据库发送游戏记录!")
			gameEndInfo := GGameEndInfo{
				UserId:      v.Uid,
				UserAccount: v.Nick,
				BetCoins:    betCoins,
				ValidBet:    betCoins,
				PrizeCoins:  v.GetCoins,
				Robot:       v.Robot,
				WaterRate:   G_DbGetGameServerData.Rate,
			}
			dbre.UserCoin = append(dbre.UserCoin, gameEndInfo)
			v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, dbre)
			recordInfo := GGameRecordInfo{
				UserId:      v.Uid,
				UserAccount: v.Nick,
				Robot:       v.Robot,
				BetCoins:    betCoins,
				BetArea:     v.DownBet,
				PrizeCoins:  v.GetCoins,
				CoinsBefore: v.Coins + betCoins,
				CoinsAfter:  v.Coins + v.GetCoins + betCoins,
				Banker:      false,
			}
			rdreq.UserRecord = append(rdreq.UserRecord, recordInfo)
			v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, rdreq)
			dbre.UserCoin = []GGameEndInfo{} //清空
			rdreq.UserRecord = []GGameRecordInfo{}
		}
		if v.GetCoins > 0 {
			v.Coins += v.GetCoins
		}
	}

	//发送结算给客户端
	for _, v := range this.Players {
		toc := ToClientBalance{
			Id:            MSG_GAME_INFO_BALANCE,
			MyCoins:       v.Coins,
			MyResult:      v.BalanceResult,
			AllResult:     this.BalanceResult,
			OtherResult:   v.OtherBalanceResult,
			MyGetCoins:    v.GetCoins,
			History:       this.GameResultHistory,
			BetAbleIndex:  this.CanUseChip(v),
			AreaWinDouble: this.GameResult,
			IsHasZhuang:   false,
		}
		v.SendNativeMsg(MSG_GAME_INFO_BALANCE, toc)
	}
}

//更新游戏记录
func (this *ExtDesk) addHistory() {
	if len(this.GameResultHistory) >= gameConfig.LimitInfo.HistoryNum {
		this.GameResultHistory = this.GameResultHistory[1:]
	}
	this.GameResultHistory = append(this.GameResultHistory, this.GameResult)
}

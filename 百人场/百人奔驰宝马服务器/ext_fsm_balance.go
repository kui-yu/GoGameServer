package main

type FSMSettle struct {
	Mark int

	EDesk       *ExtDesk
	EndDataTime int64 // 当前状态的结束时间
}

func (this *FSMSettle) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMSettle) Run() {
	DebugLog("游戏状态-结算")

	this.EndDataTime = GetTimeMS() + int64(gameConfig.StateInfo.BalanceTime)

	this.addListen() // 添加监听
	this.EDesk.GameState = GAME_STATUS_BALANCE
	this.EDesk.SendGameState(GAME_STATUS_BALANCE, int64(gameConfig.StateInfo.BalanceTime))

	this.EDesk.AddTimer(GAME_STATUS_BALANCE, gameConfig.StateInfo.BalanceTime/1000, this.TimerCall, nil)

	this.sendRecord()
	this.balanceAll()
}

func (this *FSMSettle) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_DOWNBET)
}

func (this *FSMSettle) GetMark() int {
	return this.Mark
}
func (this *FSMSettle) Leave() {
	this.removeListen()
}

func (this *FSMSettle) getRestTime() int64 {
	remainTimeMS := this.EndDataTime - GetTimeMS()
	return remainTimeMS
}

func (this *FSMSettle) addListen() {}

func (this *FSMSettle) removeListen() {}

func (this *FSMSettle) balanceAll() {
	//发送结算消息给数据库, 简单记录
	var activeUid int64 = 0
	if len(this.EDesk.Players) != 0 {
		for _, v := range this.EDesk.Players {
			if !v.Robot {
				activeUid = v.Uid
			}
		}
	}
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.EDesk.JuHao,
		Mini:        false,
		SetLeave:    1,
		ActiveUid:   activeUid,
	}

	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.EDesk.JuHao,
		BankerCard:  this.EDesk.GameResult,
	}

	var BankResult int64 = 0
	for id, v := range this.EDesk.Bets {
		BankResult += v.UserBetValue
		if id == this.EDesk.GameResult {
			BankResult -= (v.UserBetValue * int64(CarTypeMultiple[id]*10)) / 10
		}
	}

	for _, v := range this.EDesk.Players {
		var betCoins int64 = 0
		var waterRate float64 = 0
		for id, value := range v.DownBets {
			betCoins += value
			if id != this.EDesk.GameResult {
				v.BalaDownBets[id] = -value
			} else {
				v.BalaDownBets[id] = (value * int64(CarTypeMultiple[id]*10)) / 10
			}
		}

		coin, ok := v.BalaDownBets[this.EDesk.GameResult]
		if ok {
			v.WinCoins = coin - betCoins
		} else {
			v.WinCoins = -betCoins
		}

		if v.BalaDownBets[this.EDesk.GameResult] > betCoins {
			waterRate = float64(v.BalaDownBets[this.EDesk.GameResult]-betCoins) * this.EDesk.Rate
		}

		if betCoins != 0 {
			v.UnbetsCount = 0
			uc := GGameEndInfo{
				UserId:      v.Uid,
				UserAccount: v.Account,
				BetCoins:    betCoins,
				ValidBet:    betCoins,
				PrizeCoins:  v.WinCoins,
				Robot:       v.Robot,
				WaterProfit: 0,
				WaterRate:   waterRate,
			}

			dbreq.UserCoin = append(dbreq.UserCoin, uc)
			v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
			dbreq.UserCoin = []GGameEndInfo{}
		} else {
			v.UnbetsCount++
		}

		var rdDownBets [8]int32
		for id, value := range v.DownBets {
			rdDownBets[id] = int32(value)
		}

		if !v.Robot && betCoins != 0 {
			rddata := GGameRecordInfo{
				UserId:      v.Uid,
				UserAccount: v.Account,
				BetCoins:    betCoins,             // 下注的金币
				BetArea:     rdDownBets,           // 区域下注情况
				PrizeCoins:  v.WinCoins,           // 赢取的金币
				CoinsAfter:  v.Coins,              // 结束后金币
				CoinsBefore: v.Coins + v.WinCoins, // 下注前金币
				Robot:       v.Robot,
			}
			rdreq.UserRecord = append(rdreq.UserRecord, rddata)

			// 发送记录到存储
			v.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
			rdreq.UserRecord = []GGameRecordInfo{}
		}

	}

	//修改库存值
	AddLocalStock(BankResult)

	for _, v := range this.EDesk.Players {
		info := GNBalance{
			Id:     MSG_GAME_INFO_NBALANCE,
			MyCoin: v.Coins,
		}

		info.Results = make(map[int]GBetBalance)

		for id, value := range v.BalaDownBets {
			msg := GBetBalance{
				Bottom:   this.EDesk.Bets[id].DownBetValue,
				MyResult: value,
				MyBottom: v.DownBets[id],
			}
			if value > 0 {
				msg.Result = 1
			}
			info.Results[id] = msg
		}

		v.SendNativeMsg(MSG_GAME_INFO_NBALANCE, &info)
	}

	this.EDesk.RemoveAllOfflineAndExistSeat()
	this.EDesk.ResetExtDesk()
}

func (this *FSMSettle) sendRecord() {
	if len(this.EDesk.LogoRecord) < gameConfig.LimitInfo.LogoLimit {
		this.EDesk.LogoRecord = append(this.EDesk.LogoRecord, this.EDesk.GameResult)
	} else {
		this.EDesk.LogoRecord = append(this.EDesk.LogoRecord, this.EDesk.GameResult)
		this.EDesk.LogoRecord = this.EDesk.LogoRecord[1:]
	}

	info := GNRecord{
		Id:              MSG_GAME_INFO_NRECORD,
		Record:          this.EDesk.LogoRecord,
		OnlinePlayerNum: len(this.EDesk.Players),
	}
	DebugLog("开奖记录通知：", this.EDesk.LogoRecord)
	this.EDesk.BroadcastAll(MSG_GAME_INFO_NRECORD, &info)
}

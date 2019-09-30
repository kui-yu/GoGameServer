package main

type ExtDesk struct {
	Desk
	CardMgr   MgrCard // 扑克牌牌管理
	Bscore    int
	Rate      float64
	MaxRobot  int //最大机器人数
	RobotRate int //机器人概率
}

//阶段执行
func (this *ExtDesk) nextStage(stage int) {

	this.GameState = stage
	//清空所有定时器
	this.ClearTimer()

	if this.GameState == GAME_STATUS_START {
		this.GameStateBet()
	} else if this.GameState == STAGE_DEAL {
		this.GameStateDeal()
	} else if this.GameState == STAGE_PLAY {
		this.GameStatePlay()
	} else if this.GameState == STAGE_SETTLE {
		this.GameStateSettle()
	}

}

//初始化
func (this *ExtDesk) InitGame() {
	//牌内容初始化
	this.CardMgr.InitCards()
	this.CardMgr.InitNormalCards()
	//最大机器人数
	this.MaxRobot = GExtRobot.MaxRobot
	this.RobotRate = GExtRobot.RobotRate[GCONFIG.GradeType-1]
}

//自封装定时器
func (this *ExtDesk) runTimer(t int, h func(interface{})) {
	//定时器ID，定时器时间，可执行函数，可执行参数
	this.AddTimer(10, t, h, nil)
}

//广播阶段
func (this *ExtDesk) BroadStageTime(time int32) {
	stage := GStageInfo{
		Id:        MSG_GAME_INFO_STAGE,
		Stage:     int32(this.GameState),
		StageTime: time,
	}
	this.BroadcastAll(MSG_GAME_INFO_STAGE, &stage)
}

//玩家掉线广播
func (this *ExtDesk) HandleDisConnect(p *ExtPlayer, d *DkInMsg) {
	//广播给其他人，掉线
	if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
		this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Result: 0,
			Token:  p.Token,
		})
		this.DelPlayer(p.Uid)
		this.DeskMgr.LeaveDo(p.Uid)
	} else {
		p.LiXian = true
		this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
			Id:    MSG_GAME_ONLINE_NOTIFY,
			Cid:   p.ChairId,
			State: 2,
		})
	}
}

//数据通信
func (this *ExtDesk) PutSqlData() {
	//发送结算消息给数据库
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Mini:        false,
	}

	for _, v := range this.Players {
		//有效投注
		valid := v.WinCoins
		if valid < 0 {
			valid = -valid
		} else {
			valid = int64(v.BetMultiple * this.Bscore)
		}
		dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
			UserId:      v.Uid,
			UserAccount: v.Account,
			BetCoins:    int64(this.Bscore),
			ValidBet:    valid,
			PrizeCoins:  v.WinCoins,
			Robot:       v.Robot,
			WaterProfit: v.RateCoin,
			WaterRate:   this.Rate,
		})
		v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
		dbreq.UserCoin = []GGameEndInfo{}
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
		if v.Robot {
			continue
		}

		winMultiple := int(v.WinMultiple)
		if winMultiple < 0 {
			winMultiple = -winMultiple
		}

		rddata := GGameRecordInfo{
			UserId:        v.Uid,
			UserAccount:   v.Account,
			Robot:         v.Robot,
			CoinsBefore:   v.Coins - v.WinCoins,
			BetCoins:      int64(v.BetMultiple * this.Bscore),
			Coins:         v.WinCoins,
			CoinsAfter:    v.Coins,
			Cards:         ListInt32ToInt(v.NiuCards),
			BrandMultiply: v.NiuMultiple,
			BetMultiple:   v.BetMultiple,
			Multiple:      winMultiple,
			Score:         int(this.Bscore),
		}

		rdreq.UserRecord = append(rdreq.UserRecord, rddata)
		v.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
		rdreq.UserRecord = []GGameRecordInfo{}
	}
}

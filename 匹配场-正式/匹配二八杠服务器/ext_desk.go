package main

import (
	"fmt"
	"logs"
)

type ExtDesk struct {
	Desk
	//以下是用户自己定义的变量
	Bscore     int          //底分
	CardMgr    MgrCard      // 扑克牌牌管理
	Banker     int32        //庄家
	Round      int          //回合
	PutInfos   []G_PutInfo  //已出牌数量
	TotalRound int          //总回合数
	RsInfo     GSSettleInfo //当局结算信息
	Rate       float64      //费率
	MaxRobot   int          //最大机器人数
	RobotRate  int          //机器人概率
	WinList    []int
}

//出牌统计
type G_PutInfo struct {
	Value  int //牌值
	Number int //已出牌数量
}

//清空桌子
func (this *ExtDesk) ResetTable() {
	this.Bscore = 0
	this.Banker = -1
	this.Round = 0
	this.PutInfos = []G_PutInfo{}
	this.JuHao = ""
}

//初始化
func (this *ExtDesk) InitGame() {
	//总回合数
	this.TotalRound = 5
	//牌内容初始化
	this.CardMgr.InitCards()
	this.CardMgr.InitNormalCards()
	//最大机器人数
	this.MaxRobot = GExtRobot.MaxRobot
	this.RobotRate = GExtRobot.RobotRate[GCONFIG.GradeType-1]
}

//添加统计
func (this *ExtDesk) AddPutInfos(cards []int) {
	for _, c := range cards {
		for index, p := range this.PutInfos {
			if p.Value == c {
				this.PutInfos[index].Number++
				break
			}
		}
	}
}

//阶段执行
func (this *ExtDesk) nextStage(stage int) {

	this.GameState = stage
	//清空所有定时器
	this.ClearTimer()
	logs.Debug("阶段：", this.GameState)
	fmt.Println("")
	if this.GameState == GAME_STATUS_START {
		this.GameStateStart()
	} else if this.GameState == STAGE_CALL {
		this.GameStateCall()
	} else if this.GameState == STAGE_BET {
		this.GameStateBet()
	} else if this.GameState == STAGE_DEAL {
		this.GameStateDeal()
	} else if this.GameState == STAGE_SETTLE {
		this.GameStateSettle()
	} else if this.GameState == STAGE_RESTART {
		if this.Round > this.TotalRound {
			this.SysTableEnd()
		} else {
			this.GameStateRestart()
		}
	}
}

//广播阶段
func (this *ExtDesk) BroadStageTime(time int) {
	stage := GSStageInfo{
		Id:        MSG_GAME_INFO_STAGE,
		Stage:     this.GameState,
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
	}

	for _, v := range this.Players {
		valid := v.WinCoins
		if valid < 0 {
			valid = -valid
		} else {
			if v.ChairId == this.Banker {
				valid = int64(this.Bscore)
			} else {
				valid = int64(this.Bscore * v.PlayMultiple)
			}
		}
		dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
			UserId:      v.Uid,
			UserAccount: v.Account,
			BetCoins:    int64(this.Bscore),
			ValidBet:    valid,
			PrizeCoins:  int64(v.WinCoins),
			Robot:       v.Robot,
			WaterProfit: v.RateCoins,
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
		Round:       this.Round,
	}
	for _, v := range this.Players {
		if v.Robot {
			continue
		}

		playMultiple := v.PlayMultiple
		if v.ChairId == this.Banker {
			playMultiple = v.CallMultiple
		}

		winMultiple := v.WinMultiple
		if winMultiple < 0 {
			winMultiple = -winMultiple
		}

		isBanker := 0
		if v.ChairId == this.Banker {
			isBanker = 1
		}
		rddata := GGameRecordInfo{
			UserId:        v.Uid,
			UserAccount:   v.Account,
			Robot:         v.Robot,
			CoinsBefore:   v.Coins - v.WinCoins,
			BetCoins:      int64(playMultiple * this.Bscore),
			Coins:         v.WinCoins,
			CoinsAfter:    v.Coins,
			Cards:         v.HandCards,
			BrandMultiply: 1,
			BetMultiple:   playMultiple,
			Multiple:      winMultiple,
			Banker:        isBanker,
			Score:         this.Bscore,
		}

		rdreq.UserRecord = append(rdreq.UserRecord, rddata)
		v.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
		rdreq.UserRecord = []GGameRecordInfo{}
	}
}

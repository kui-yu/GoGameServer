package main

import (
	"encoding/json"
	// "logs"
)

type ExtDesk struct {
	Desk
	//以下是用户自己定义的变量
	Bscore        int64           // 底分
	CardMgr       MgrCard         // 扑克牌牌管理
	CallPlayer    int32           // 叫牌玩家
	Round         int             // 记录下注轮数
	Pround        int             // 计算轮数的玩家位置
	CoinList      []int64         // 全部玩家下注记录
	ChairList     []int32         // 玩家座位表
	MinCoin       int64           // 最小跟注
	SettleContest []PlayerContest //特殊情况金币不足比牌记录
	Rate          float64         //概率

	NotAllowRobotInRoom bool //机器人进入房间状态
	MaxRobot            int  //机器人数
	RobotRate           int  //概率
	RobotWinner         bool //是否机器人赢
	// RsInfo  GSSettleInfo //当局结算信息
}

func (this *ExtDesk) InitAttribute() {
	this.RobotRate = GExtRobot.RobotRate[GCONFIG.GradeType-1]
}

func (this *ExtDesk) GameOver() {
	// fmt.Println("run this")
	this.GameState = GAME_STATUS_FREE
	this.JuHao = ""
	//重置桌面属性
	this.CoinList = []int64{}
	this.Round = 0
	this.SettleContest = []PlayerContest{}
	this.ChairList = []int32{}
	this.DeskMgr.BackDesk(this)
	this.Players = []*ExtPlayer{}
	// this.CardMgr.InitCards()
	this.RobotWinner = false
	this.MaxRobot = 0
}

//阶段执行
func (this *ExtDesk) nextStage(stage int) {
	filter := false
	if stage == this.GameState {
		filter = true
	}

	this.GameState = stage
	//清空所有定时器
	this.ClearTimer()

	if this.GameState == GAME_STATUS_START {
		this.GameStateStart()
	} else if this.GameState == STAGE_SOLO { //比牌阶段
		this.GameStateContest()
	} else if this.GameState == STAGE_PLAY_OPERATION { //操作阶段
		this.GameStateOperation(filter)
	} else if this.GameState == STAGE_SETTLE { //结算阶段
		this.GameStateSettle()
	} else if this.GameState == GAME_STATUS_END { //游戏结束
		this.GameStateEnd()
	}
}

//叫牌玩家通知//想搞这个的可以参考房卡炸金花
func (this *ExtDesk) MsgCallPlayer() {
	count := 0
	for _, v := range this.Players {
		if v.CardType != 2 {
			count++
		}
	}
	if count == 1 { //进入结算阶段
		this.nextStage(GAME_STATUS_END)
		return
	}

	num := 0
	for k, v := range this.Players {
		if v.ChairId == this.CallPlayer {
			num = k
			break
		}
	}

	cnum := 0
	for i := num; i < len(this.Players); i++ {
		if i+1 == this.Pround || (this.Pround == 0 && i == len(this.Players)-1) { //轮数
			this.Round += 1
		}
		if i == len(this.Players)-1 { //下一个叫牌玩家
			this.CallPlayer = this.ChairList[0]
			cnum = 0
		} else {
			this.CallPlayer = this.ChairList[i+1]
			cnum = i + 1
		}
		Check := false
		for _, v := range this.Players {
			if v.ChairId == this.CallPlayer && v.CardType == 2 {
				Check = true
			}
		}
		if Check {
			if i != len(this.Players)-1 {
				continue
			} else {
				i = -1
				continue
			}
		}
		break
	}

	if this.Round == GameRound { //轮数记录
		this.nextStage(GAME_STATUS_END)
		return
	}

	// fmt.Println(this.Round)
	info := GSPlayerCallPlayer{
		Id:         MSG_GAME_INFO_CALLPLAYER_REPLY,
		Player:     this.CallPlayer,
		Round:      this.Round,
		CoinEnough: IsCoinEnough(this.Players[cnum].Coins, this.Players[cnum].PayCoin, this.Bscore, this.MinCoin, this.Players[cnum].CardType),
		MinCoin:    this.MinCoin,
	}

	if this.Players[cnum].AutoFollowUp == 1 { //自动跟注
		this.BroadcastAll(MSG_GAME_INFO_CALLPLAYER_REPLY, &info)
		this.AddTimer(2, 1, this.AutoFollowIn, nil)
	} else {
		this.BroadcastAll(MSG_GAME_INFO_CALLPLAYER_REPLY, &info)
		this.nextStage(STAGE_PLAY_OPERATION)
		//机器人跳转
		for _, v := range this.Players {
			if v.ChairId == this.CallPlayer {
				if v.Robot {
					//随机秒数
					//跳转机器人操作
					this.AddTimer(11, 3, this.RobotPlay, v)
				}
				break
			}
		}
	}
}

//自动跟注
func (this *ExtDesk) AutoFollowIn(d interface{}) {
	for _, v := range this.Players {
		if this.CallPlayer == v.ChairId {
			this.GetGamePlay(4, v)
			break
		}
	}
}

//下注操作方法调用
func (this *ExtDesk) GetGamePlay(Op int, p *ExtPlayer) {
	info := GAPlayerOperation{
		Id:        MSG_GAME_INFO_PLAY_INFO,
		Operation: Op,
	}
	fd, _ := json.Marshal(info)
	msg := DkInMsg{
		Data: string(fd),
	}

	this.HandleGamePlay(p, &msg)
}

//下注金币推送
func (this *ExtDesk) CoinPush() {
	info := GSCoinMsg{
		Id: MSG_GAME_INFO_COIN,
	}
	var sumCoin [][]int64
	sumCoin = append(sumCoin, this.CoinList)
	for _, v := range this.Players {
		sumCoin = append(sumCoin, v.PayCoin)
	}

	for i := 0; i < len(sumCoin); i++ {
		sum := int64(0)
		for j := 0; j < len(sumCoin[i]); j++ {
			sum += sumCoin[i][j]
		}
		if i == 0 {
			info.AllCoin = sum
		} else {
			info.PCoin = append(info.PCoin, sum)
		}
	}
	info.Round = this.Round
	this.BroadcastAll(MSG_GAME_INFO_COIN, &info)
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

//广播  排除操作者
func (this *ExtDesk) BroadExceptOpPlay(p *ExtPlayer, id int, msg interface{}) {
	for _, v := range this.Players {
		if v == p {
			continue
		}
		v.SendNativeMsg(id, &msg)
	}
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
		this.BroadcastAll(MSG_GAME_ONLINE_NOTIFY, &GOnLineNotify{
			Id:    MSG_GAME_ONLINE_NOTIFY,
			Cid:   p.ChairId,
			State: 2, //1上线，2掉线
		})
	}
}

//重写   离开的玩家不发送消息
func (this *ExtDesk) BroadcastAll(id int, d interface{}) {
	for _, v := range this.Players {
		if v.IsLeave == 1 {
			continue
		}
		v.SendNativeMsg(id, d)
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
		if v.IsLeave == 1 {
			continue
		}
		valid := int64(0) //下注*低分
		for i := 0; i < len(v.PayCoin); i++ {
			valid += v.PayCoin[i]
		}
		valid = valid * this.Bscore
		dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
			UserId:      v.Uid,
			UserAccount: v.Account,
			BetCoins:    int64(this.Bscore),
			ValidBet:    valid,             //下注*低分
			PrizeCoins:  int64(v.WinCoins), //输赢金币
			Robot:       v.Robot,
			WaterProfit: v.RateCoins,
			WaterRate:   this.Rate,
		})
		v.SendNativeMsg(MSG_GAME_END_NOTIFY, &dbreq)
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
		if v.IsLeave == 1 {
			continue
		}
		sum := int64(0)
		for i := 0; i < len(v.PayCoin); i++ {
			sum += v.PayCoin[i]
		}
		// logs.Debug("结算wincoins", v.WinCoins)
		rddata := GGameRecordInfo{
			UserId:      v.Uid,
			UserAccount: v.Account,
			Robot:       v.Robot,
			CoinsBefore: v.Coins - v.WinCoins,
			BetCoins:    sum * this.Bscore, //下注金币
			Coins:       v.WinCoins,
			CoinsAfter:  v.Coins,
			Cards:       v.OldHandCard,
			Multiple:    1,
			Score:       this.Bscore,
		}

		rdreq.UserRecord = append(rdreq.UserRecord, rddata)
		v.SendNativeMsg(MSG_GAME_END_RECORD, &rdreq)
		rdreq.UserRecord = []GGameRecordInfo{}
	}
	// fmt.Println("not msg")
}

package main

import (
	"encoding/json"
)

type ExtDesk struct {
	Desk

	Mark              string          //房间标志
	DownCards         []int           //所有的牌
	CurrDownCardIdx   int             //当前发牌的位置
	PublicCards       []int           //公共牌
	SeatBank          int             //庄家座位号
	SeatOperateId     int             //在哪个座位玩家操作
	IsOperateOpen     bool            //是否有开牌操作
	JackpotVal        int64           //奖池总额
	CurrStage         int             //当前阶段
	fsms              map[int]FSMBase //状态机集合
	currFSM           FSMBase         //当前状态机
	upFSM             FSMBase         //上一个状态机
	IsExistOperateFsm bool            //是否存在操作状态
}

func (this *ExtDesk) InitExtData() {
	this.DeskMgr.SetDeskAllotModel(1)
	//玩家匹配 400001
	this.Handle[MSG_GAME_AUTO] = this.HandleGameAuto
	//断线重连 400010
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
	//断线消息 400013断线
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect

	//玩家设置金币
	this.Handle[MSG_GAME_QGameSetCoin] = this.HandleSettCoin

	this.fsms = make(map[int]FSMBase)
	this.addFSM(GameStatusWaitStart, new(FsmWaitStart))
	this.addFSM(GameStatusRandBank, new(FsmRandBank))
	this.addFSM(GameStatusHoleCards, new(FsmHoleCards))
	this.addFSM(GameStatusFlopCards, new(FsmFlopCards))
	this.addFSM(GameStatusTurnCards, new(FsmTurnCards))
	this.addFSM(GameStatusRiverCards, new(FsmRiverCards))
	this.addFSM(GameStatusUserOperate, new(FsmUserOperate))
	this.addFSM(GameStatusResults, new(FsmResults))

	this.SeatBank = 0xFF
	this.Mark = FormatDeskId(this.Id, GCONFIG.GradeType)

	this.ResetDeskInfo()

	this.RunFSM(GameStatusWaitStart)
}

// ==============底层消息==========
func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {
	DebugLog("HandleGameAuto", p.Uid, p.Account)
	if len(this.Players) > GCONFIG.PlayerNum {
		p.SendNetMessage(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
			Id:       MSG_GAME_AUTO_REPLY,
			CostType: GetCostType(),
			Result:   1,
			Err:      "桌子已满，加入失败",
		})
		return
	}

	// 判断机器人数量
	robotNum := 0
	for _, v := range this.Players {
		if v.Robot {
			robotNum += 1
		}
	}
	if robotNum > gameConfig.RobotNum {
		//发送匹配成功
		p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
			Id:       MSG_GAME_AUTO_REPLY,
			CostType: GetCostType(),
			Result:   13,
			Err:      "桌子已满，加入失败",
		})
		//踢出
		this.LeaveByForce(p)
		return
	}

	p.Init()
	p.Sid = this.GetEmptySeat()

	p.SendNetMessage(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})

	// 群发用户信息
	this.SendNetMessage(MSG_GAME_NGameUserChange, struct {
		ChangeType int
		UserInfo   GCUserInfo
	}{
		ChangeType: 1,
		UserInfo: GCUserInfo{
			Uid:      p.Uid,
			NickName: p.Account,
			Avatar:   p.Head,
			Sid:      p.Sid,
			Coin:     p.Coins,
			Online:   !p.LiXian,
			State:    p.State,
		},
	}, p.Uid, "UserInfo/NickName")
	// 发送房间信息
	this.SendDeskInfo(p)

	this.CheckGameStateAndStart()

	this.GetFSM(0).OnUserOnline(p)
}

func (this *ExtDesk) SendUserOnlineState(p *ExtPlayer) {
	this.SendNetMessage(MSG_GAME_ONLINE_NOTIFY, struct {
		Uid    int64
		Online bool
	}{
		Uid:    p.Uid,
		Online: !p.LiXian,
	}, p.Uid)
}

func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	DebugLog("HandleReconnect", p.Uid, p.Account)

	p.SendNetMessage(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:       MSG_GAME_RECONNECT_REPLY,
		Result:   0,
		CostType: GetCostType(),
	})

	p.LiXian = false
	this.SendUserOnlineState(p)

	//发送房间信息
	this.SendDeskInfo(p)

	this.GetFSM(0).OnUserOnline(p)
}

func (this *ExtDesk) HandleDisConnect(p *ExtPlayer, d *DkInMsg) {
	p.LiXian = true
	if p.State != UserStateGameIn {
		DebugLog("DisConnect 用户离开", p.Uid)
		p.LiXian = false
		this.SendNetMessage(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    int32(p.Sid),
			Uid:    p.Uid,
			Result: 0,
			Token:  p.Token,
		})
		this.DelPlayer(p.Uid)
		this.DeskMgr.LeaveDo(p.Uid)
	} else {
		DebugLog("DisConnect 用户离线", p.Uid)
		this.SendUserOnlineState(p)
	}
}

// 重写底层用户离开
func (this *ExtDesk) Leave(p *ExtPlayer) bool {
	if p.State != UserStateGameIn {
		DebugLog("Leave 用户离线")
		this.SendNetMessage(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    int32(p.Sid),
			Uid:    p.Uid,
			Result: 0,
			Token:  p.Token,
		})
		p.LiXian = true
		this.DelPlayer(p.Uid)
		this.DeskMgr.LeaveDo(p.Uid)
	} else {
		this.SendNetMessage(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    int32(p.Sid),
			Uid:    p.Uid,
			Result: 1,
			Token:  p.Token,
			Err:    "正在游戏中，离开失败",
		})
	}
	return true
}

func (this *ExtDesk) SendNetMessage(cmd int, data interface{}, args ...interface{}) {
	// 需要排除的集合
	excludes := []int64{}
	if len(args) > 0 {
		switch v := args[0].(type) {
		case int64:
			excludes = append(excludes, v)
		case []int64:
			excludes = append(excludes, v...)
		}
	}

	for _, p := range this.Players {
		if FindIndexFromInt64(excludes, p.Uid) == -1 {
			if len(args) > 1 {
				p.SendNetMessage(cmd, data, args[1:]...)
			} else {
				p.SendNetMessage(cmd, data)
			}
		}
	}
}

func (this *ExtDesk) addFSM(mark int, fsm FSMBase) {
	fsm.InitFsm(mark, this)
	this.fsms[mark] = fsm
}

func (this *ExtDesk) GetFSM(mark int) FSMBase {
	if mark != 0 {
		return this.fsms[mark]
	}
	return this.currFSM
}

func (this *ExtDesk) RunFSM(mark int, args ...interface{}) {
	var upMark int = 0
	if this.currFSM != nil {
		this.upFSM = this.currFSM
		this.upFSM.Leave()
		upMark = this.upFSM.GetMark()
	}

	this.currFSM = this.GetFSM(mark)
	this.currFSM.Run(upMark, args)
}

//===================自定义消息===============

func (this *ExtDesk) CheckGameStateAndStart() {
	// 检查游戏状态，是否要开始游戏
	if this.GetFSM(0).GetMark() == GameStatusWaitStart {
		waitStartNum := 0
		for _, player := range this.Players {
			if player.State == UserStateWaitStart {
				waitStartNum += 1
			}
		}

		if waitStartNum > 1 {
			this.RunFSM(GameStatusRandBank)
		}
	}
}

// 接收到玩家设置金币
func (this *ExtDesk) HandleSettCoin(p *ExtPlayer, d *DkInMsg) {
	data := struct {
		Coin int64
	}{}

	json.Unmarshal([]byte(d.Data), &data)

	rep := struct {
		Coin int64
		Err  string
	}{}

	if data.Coin > p.Coins {
		rep.Coin = 0
		rep.Err = "金币不足，设置失败"
		p.SendNetMessage(MSG_GAME_RGameSetCoin, rep)
		return
	}
	if data.Coin <= gameConfig.BigBlindCoin {
		rep.Coin = 0
		rep.Err = "金币不足，设置失败"
		p.SendNetMessage(MSG_GAME_RGameSetCoin, rep)
		return
	}

	p.CarryCoin = p.Coins
	p.State = UserStateWaitStart

	//回复更新成功
	rep.Coin = data.Coin
	p.SendNetMessage(MSG_GAME_RGameSetCoin, rep)

	//通知所有玩家当前玩家更新
	ndata := struct {
		Uid   int64
		Sid   int
		State int
		Coin  int64
	}{
		Uid:   p.Uid,
		Sid:   p.Sid,
		State: p.State,
		Coin:  p.CarryCoin,
	}

	this.SendNetMessage(MSG_GAME_NGameUserChange, ndata)

	// 检查游戏状态，是否要开始游戏
	this.CheckGameStateAndStart()
}

func (this *ExtDesk) ResetDeskInfo() {
	this.JuHao = GetJuHao()

	this.DownCards = []int{}
	this.CurrDownCardIdx = 0
	this.PublicCards = []int{}
	this.SeatOperateId = 0xFF
	this.IsOperateOpen = false
	this.JackpotVal = 0
	this.CurrStage = 0xFF
	this.IsExistOperateFsm = false

	// 重置所有状态机
	for _, v := range this.fsms {
		v.Reset()
	}

	// 重置所有玩家
	for _, p := range this.Players {
		p.Reset()
	}
}

//发送桌子信息
func (this *ExtDesk) SendDeskInfo(player *ExtPlayer) {
	DebugLog("发送房间消息")

	gameInfo := struct {
		JuHao       string
		GameStatus  int
		OverTime    int64
		DownBlind   int64
		JackpotVal  int64
		SeatBank    int
		PublicCards []int
		Users       []interface{}
	}{
		JuHao:       this.JuHao,
		GameStatus:  this.GetFSM(0).GetMark(),
		OverTime:    this.GetFSM(0).GetRestTime(),
		DownBlind:   gameConfig.SmallBlindCoin,
		JackpotVal:  this.JackpotVal,
		SeatBank:    this.SeatBank,
		PublicCards: this.PublicCards,
	}

	for _, p := range this.Players {
		userInfo := GCUserInfo{
			Uid:          p.Uid,
			NickName:     p.Account,
			Avatar:       p.Head,
			Sid:          p.Sid,
			State:        p.State,
			Coin:         p.CarryCoin,
			IsBank:       p.IsBank,
			IsFold:       p.IsFold,
			IsAllIn:      p.AllInStage != 0xFF,
			DownCoins:    p.GetDownBet(),
			Online:       !p.LiXian,
			CurrStageOpt: p.StageOperate,
		}

		if this.GetFSM(0).GetMark() == GameStatusResults || p.Uid == player.Uid {
			userInfo.Cards = p.Cards
		}
		gameInfo.Users = append(gameInfo.Users, userInfo)
	}

	player.SendNetMessage(MSG_GAME_NGameReconnectInfo, gameInfo, "Users/*/NickName")
}

//发送设置携带筹码
func (this *ExtDesk) SendUserSettCoin(p *ExtPlayer) {
	maxCoin := p.Coins
	p.State = UserStateSettCoin
	if maxCoin > gameConfig.UserSettCoinMax {
		maxCoin = gameConfig.UserSettCoinMax
	}

	p.SendNetMessage(MSG_GAME_NGameSetCoin, struct {
		MinCoin int64
		MaxCoin int64
	}{
		MinCoin: gameConfig.UserSettCoinMin,
		MaxCoin: maxCoin,
	})
}

//获取空位子
func (this *ExtDesk) GetEmptySeat() int {
	isExist := false
	start := 0
	num := GCONFIG.PlayerNum

	for ; start < num; start++ {
		isExist = false
		for _, p := range this.Players {
			if p.Sid == start {
				isExist = true
				break
			}
		}

		if isExist == false {
			break
		}
	}

	return start
}

//发送桌子状态
func (this *ExtDesk) SendDeskStatus(status int, ms int) {
	this.SendNetMessage(MSG_GAME_NGameStatus, GCGameStatusInfo{
		GameStatus: status,
		OverTime:   ms,
	})
}

func (this *ExtDesk) SendUpdateRoomNo() {
	this.SendNetMessage(MSG_GAME_NDeskUpdate, struct {
		JuHao string
	}{
		JuHao: this.JuHao,
	})
}

func (this *ExtDesk) GetPlayerFromSid(sid int) *ExtPlayer {
	var p *ExtPlayer = nil
	for _, v := range this.Players {
		if v.Sid == sid {
			p = v
			break
		}
	}

	return p
}

//查找同一状态的上一个玩家
func (this *ExtDesk) GetUpPlayer(sid int, s ...interface{}) *ExtPlayer {
	if sid == 0xFF {
		sid = this.SeatOperateId
	}

	state := UserStateGameIn
	sp := this.GetPlayerFromSid(sid)
	if sp != nil {
		state = sp.State
	}

	if len(s) != 0 {
		state = s[0].(int)
	}

	var upP *ExtPlayer = nil
	for _, v := range this.Players {
		if v.Sid < sid && v.State == state && (upP == nil || v.Sid > upP.Sid) {
			upP = v
		}
	}

	if upP != nil {
		return upP
	}
	// 找到最大的
	for _, v := range this.Players {
		if v.Sid > sid && v.State == state && (upP == nil || v.Sid > upP.Sid) {
			upP = v
		}
	}
	if upP != nil {
		return upP
	}

	return sp
}

//查找同一状态的下一个玩家
func (this *ExtDesk) GetNextPlayer(sid int, s ...interface{}) *ExtPlayer {
	if sid == 0xFF {
		sid = this.SeatOperateId
	}

	state := UserStateGameIn
	sp := this.GetPlayerFromSid(sid)
	if sp != nil {
		state = sp.State
	}

	if len(s) != 0 {
		state = s[0].(int)
	}

	var nextP *ExtPlayer = nil
	for _, v := range this.Players {
		if v.Sid > sid && v.State == state && (nextP == nil || v.Sid < nextP.Sid) {
			nextP = v
		}
	}

	if nextP != nil {
		return nextP
	}
	// 找到最小的
	for _, v := range this.Players {
		if v.Sid < sid && v.State == state && (nextP == nil || v.Sid < nextP.Sid) {
			nextP = v
		}
	}
	if nextP != nil {
		return nextP
	}

	return sp
}

//获得玩家当前状态下注额
func (this *ExtDesk) GetPlayerStageDownBet(p *ExtPlayer) int64 {
	stageIdx := this.GetStageIdx(0xFF)
	if len(p.DownCoins) > stageIdx {
		return p.DownCoins[stageIdx]
	}
	return 0
}

//添加玩家下注,并通知所有玩家
func (this *ExtDesk) AddPlayerDownBet(p *ExtPlayer, val int64) {
	if val != 0 {
		this.SendNetMessage(MSG_GAME_NGameJackpotChange, struct {
			Sid       int   //用户id
			ChangeVal int64 //改变的值
			Value     int64 //奖池新金额
		}{
			Sid:       p.Sid,
			ChangeVal: val,
			Value:     this.JackpotVal + val,
		})
		this.JackpotVal += val
	}
	p.AddDownBet(this.GetStageIdx(0xFF), val)
}

func (this *ExtDesk) GetStageIdx(stage int) int {
	if stage == 0xFF {
		stage = this.CurrStage
	}
	for i, v := range StageDefines {
		if v == stage {
			return i
		}
	}
	return 0
}

func (this *ExtDesk) GetNextStageMark() int {
	return StageDefines[this.GetStageIdx(0xFF)+1]
}

func (this *ExtDesk) GetStageMaxBet() int64 {
	DebugLog("StageDownBetLimit stageIdx:", this.GetStageIdx(0xFF), gameConfig.StageDownBetLimit)

	return gameConfig.StageDownBetLimit[this.GetStageIdx(0xFF)]
}

func (this *ExtDesk) GetStageBetTotal(sid int) int64 {
	var totalCoin int64 = 0

	p := this.GetPlayerFromSid(sid)

	sIdx := this.GetStageIdx(p.AllInStage)
	sVal := p.DownCoins[sIdx]

	for i := 0; i <= sIdx; i++ {
		for _, p := range this.Players {
			if len(p.DownCoins) > i {
				if i < sIdx {
					totalCoin += p.DownCoins[i]
				} else {
					if p.DownCoins[i] < sVal {
						totalCoin += p.DownCoins[i]
					} else {
						totalCoin += sVal
					}
				}
			}
		}
	}

	return totalCoin
}

// 获得某个状态
func (this *ExtDesk) GetUserStateTotal(state int) int {
	num := 0
	for _, p := range this.Players {
		if p.State == state {
			num += 1
		}
	}

	return num
}

// 设置第一个操作的玩家
func (this *ExtDesk) SetFirstOperateUser() {
	sid := this.SeatBank
	for {
		nextP := this.GetNextPlayer(sid)
		sid = nextP.Sid
		if nextP.State != UserStateGameIn {
			continue
		}

		if nextP.IsFold == true || nextP.AllInStage != 0xFF {
			continue
		}
		break
	}

	this.SeatOperateId = sid
}

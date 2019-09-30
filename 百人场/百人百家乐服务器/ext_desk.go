package main

import (
	"encoding/json"
	"fmt"
	"logs"
	"math/rand"
	"sync"
	"time"

	"bl.com/seatlist"
	"bl.com/util"
)

type ExtDesk struct {
	Desk
	OpenCardTime int //开牌阶段时间，用于大厅记录
	sync.RWMutex
	CardMgr MgrCard // 扑克牌牌管理

	SeatMgr seatlist.MgrSeat // 座位管理

	arealistCoins     util.AreaList // 总下注情况
	userArealistCoins util.AreaList // 用户下注情况

	Count       int32   // 当前局数
	IdleCard    []int32 // 闲
	BankerCard  []int32 // 庄
	IdleDians   []int32 // 闲点
	BankerDians []int32 // 庄点
	WinArea     []bool  // 赢取区域
	RoomId      string  // 房号
	GameId      string  // 局号
	GameLimit   Limit   // 限红
	BetList     []int64 // 下注金币
	BetArea     []bool  // 可下注金币
	RunChart    []int32 // 走势
	TypeTimes   []int32 // 走势各类型次数

	NewBet bool

	gameUserListLK sync.RWMutex // 玩家列表读写锁
	Seat           []GSInfo     // 座位信息
	wCoins         float64      // 总赢取
	tCount         float64      // 总局数
	totCoins       float64      // 总投注
}

func (this *ExtDesk) InitExtData() {
	//牌内容初始化
	this.CardMgr.InitCards()
	this.CardMgr.InitNormalCards()
	//

	this.Handle[MSG_GAME_AUTO] = this.HandleGameAuto
	this.Handle[MSG_GAME_INFO_BET] = this.HandleGameBet
	this.Handle[MSG_GAME_INFO_RUN_CHART] = this.HandleGameRunChart
	this.Handle[MSG_GAME_INFO_USER_LIST] = this.HandleGameUserList
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
	this.Handle[MSG_GAME_INFO_INTO] = this.HandleGameAutoFinal
	this.Handle[MSG_GAME_INFO_EXIT] = this.HandleGameExit
	this.Handle[MSG_GAME_INFO_GET_RECORD] = this.HandleGetRecord

	this.RoomId = util.BuildRoomId(GCONFIG.GradeType+1, this.Id+1)
	this.GameLimit = gameConfig.LimitInfo.Limit[GCONFIG.GradeType-1]
	this.BetList = gameConfig.LimitInfo.BetCoins[GCONFIG.GradeType-1].Bet[:]
	logs.Debug("筹码范围：", this.GameLimit, this.BetList)

	// 设置座位数
	this.SeatMgr.SetSeatNum(gameConfig.DeskInfo.SeatCount)

	this.NewBet = false
	this.ResetAreaCoins()
	// 开始洗牌
	this.TimerShuffle(nil)
}

//
func (this *ExtDesk) BroadcastAll(id int, d interface{}) {
	allUser := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range allUser {
		v.(*ExtPlayer).SendNativeMsg(id, d)
	}
}

func (this *ExtDesk) BroadcastOther(p *ExtPlayer, id int, d interface{}) {
	allUser := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range allUser {
		if v.(*ExtPlayer).Uid == p.Uid {
			continue
		}
		v.(*ExtPlayer).SendNativeMsg(id, d)
	}
}

// 重置下注列表
func (this *ExtDesk) ResetAreaCoins() {
	this.arealistCoins.Init(9)
	this.userArealistCoins.Init(9)
}

// 洗牌初始化信息
func (this *ExtDesk) ExtDeskInit() {
	this.Count = 0
	this.CardMgr.Shuffle()
	this.RunChart = []int32{}
	this.BetArea = []bool{true, true, true, true, true, true, true, true, true}
	// 闲、庄、和、庄对、闲对、总次数
	this.TypeTimes = []int32{0, 0, 0, 0, 0, 0}
}

// 获取下注列表
func (this *ExtDesk) GetAreaCoinsList() []int64 {
	ret := this.arealistCoins.GetValueList()
	return ret
}

// 获取下注金币
func (this *ExtDesk) GetAreaCoins() int64 {
	ret := this.arealistCoins.GetTotValue()
	return ret
}
func (this *ExtDesk) GetAreaCoin(area int) int64 {
	ret := this.arealistCoins.GetValue(area)
	return ret
}
func (this *ExtDesk) GetUserAreaCoin(area int) int64 {
	ret := this.userArealistCoins.GetValue(area)
	return ret
}

// 添加下注
func (this *ExtDesk) AddAreaCoins(area int, coins int64) bool {
	this.NewBet = true
	ret := this.arealistCoins.AddValue(area, coins)

	return ret
}
func (this *ExtDesk) AddUserAreaCoins(area int, coins int64) bool {
	this.NewBet = true
	ret := this.userArealistCoins.AddValue(area, coins)

	return ret
}

func (this *ExtDesk) Replenish() {
	dian := int32(GetLogicValue(this.IdleCard[0])+GetLogicValue(this.IdleCard[1])) % 10
	this.IdleDians = append(this.IdleDians, dian)

	bankerDian := int32(GetLogicValue(this.BankerCard[0])+GetLogicValue(this.BankerCard[1])) % 10
	this.BankerDians = append(this.BankerDians, bankerDian)
	if dian >= 8 || bankerDian >= 8 {
		return
	}

	// 闲小于等于5 补牌
	if dian <= 5 {
		this.IdleCard = append(this.IdleCard, this.CardMgr.SendOneCard())
		dian = (dian + int32(GetLogicValue(this.IdleCard[2]))) % 10
		this.IdleDians = append(this.IdleDians, dian)
	}

	// 庄补牌
	var idleCard int32 = -1
	if 3 == len(this.IdleCard) {
		idleCard = int32(GetLogicValue(this.IdleCard[2])) % 10
	}
	isNeddAppend := false
	switch bankerDian {
	case 0, 1, 2:
		isNeddAppend = true
	case 3:
		if idleCard != 8 {
			isNeddAppend = true
		}
	case 4:
		if (idleCard > 1 && idleCard < 8) || idleCard == -1 {
			isNeddAppend = true
		}
	case 5:
		if (idleCard > 3 && idleCard < 8) || idleCard == -1 {
			isNeddAppend = true
		}
	case 6:
		if (idleCard > 5 && idleCard < 8) || idleCard == -1 {
			isNeddAppend = true
		}
	}

	if isNeddAppend {
		this.BankerCard = append(this.BankerCard, this.CardMgr.SendOneCard())
		bankerDian = (bankerDian + int32(GetLogicValue(this.BankerCard[2]))) % 10
		this.BankerDians = append(this.BankerDians, bankerDian)
	}
}
func (this *ExtDesk) ResetDian() {
	this.BankerDians = []int32{}
	this.IdleDians = []int32{}
}

// 获取制作牌
func (this *ExtDesk) BuildCards() {
	var idleCard []int32
	var bankerCard []int32

	for i := 0; i < 2; i++ {
		idleCard = append(idleCard, this.CardMgr.SendOneCard())
		bankerCard = append(bankerCard, this.CardMgr.SendOneCard())
	}
	//没有玩家下注，纯随机
	this.IdleCard = idleCard
	this.BankerCard = bankerCard
	this.Replenish()
	//如果没真实玩家下注，有一点几率控制庄家输
	if !this.IsValidBet() && this.IsValidBetByRobot() && BankerLose() && GetCostType() == 1 {
		fmt.Println("进入庄家百分75概率")
		this.GetWinOrLoseResult(this.BankerCard, this.IdleCard, 0, 6)
		return
	}
	//风控
	curCd := CalPkAll(StartControlTime, time.Now().Unix())
	if CD-curCd < 0 && GetCostType() == 1 && this.IsValidBet() { //进入风控换牌
		this.GetWinOrLoseResult(this.BankerCard, this.IdleCard, 1, 6)
	}
}

//判断是否有有效下注
func (this *ExtDesk) IsValidBet() bool {
	for _, v := range this.Players {
		if !v.Robot && v.IsBet {
			return true
		}
	}
	return false
}

//判断是否只有机器人下注
func (this *ExtDesk) IsValidBetByRobot() bool {
	for _, v := range this.Players {
		if v.Robot && v.IsBet {
			return true
		}
	}
	return false
}

//百分75概率
func BankerLose() bool {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(100) + 1
	if r <= 75 {
		return true
	} else {
		return false
	}
}

//递归换牌
func (this *ExtDesk) GetWinOrLoseResult(zCard []int32, xCard []int32, c int, flag int) int64 {
	this.BankerCard = zCard
	this.IdleCard = xCard
	if flag == 1 {
		return 0
	}
	if c == 0 {
		if res := this.testResult(); res > 0 {
			return res
		} else {
			for _, v := range zCard {
				this.CardMgr.MVSourceCard = append(this.CardMgr.MVSourceCard, byte(v))
				this.CardMgr.OutCards = this.CardMgr.OutCards[:len(this.CardMgr.OutCards)-1]
			}
			for _, l := range xCard {
				this.CardMgr.MVSourceCard = append(this.CardMgr.MVSourceCard, byte(l))
				this.CardMgr.OutCards = this.CardMgr.OutCards[:len(this.CardMgr.OutCards)-1]
			}
			this.CardMgr.DisturbCards() //打乱牌
			this.ResetDian()            //初始化点
			//发牌
			idleCard := []int32{}
			bankerCard := []int32{}
			for i := 0; i < 2; i++ {
				idleCard = append(idleCard, this.CardMgr.SendOneCard())
				bankerCard = append(bankerCard, this.CardMgr.SendOneCard())
			}
			this.IdleCard = idleCard
			this.BankerCard = bankerCard
			this.Replenish()
			flag--
			this.GetWinOrLoseResult(this.BankerCard, this.IdleCard, 0, flag)
		}
	} else if c == 1 {
		if res := this.testResult(); res < 0 {
			return res
		} else {
			for _, v := range zCard {
				this.CardMgr.MVSourceCard = append(this.CardMgr.MVSourceCard, byte(v))
				this.CardMgr.OutCards = this.CardMgr.OutCards[:len(this.CardMgr.OutCards)-1]
			}
			for _, l := range xCard {
				this.CardMgr.MVSourceCard = append(this.CardMgr.MVSourceCard, byte(l))
				this.CardMgr.OutCards = this.CardMgr.OutCards[:len(this.CardMgr.OutCards)-1]
			}
			this.CardMgr.DisturbCards() //打乱牌
			this.ResetDian()            //初始化点
			//发牌
			idleCard := []int32{}
			bankerCard := []int32{}
			for i := 0; i < 2; i++ {
				idleCard = append(idleCard, this.CardMgr.SendOneCard())
				bankerCard = append(bankerCard, this.CardMgr.SendOneCard())
			}
			this.IdleCard = idleCard
			this.BankerCard = bankerCard
			this.Replenish()
			flag--
			this.GetWinOrLoseResult(this.BankerCard, this.IdleCard, 1, flag)
		}
	}
	return 0
}
func (this *ExtDesk) testResult() int64 {
	var winCoins int64
	areaRes := this.getAreaResult()
	for _, j := range this.Players {
		if j.Robot || !j.IsBet { //机器人或者没下注跳过
			continue
		}
		for k, v := range areaRes {
			switch k {
			case 0:
				if j.GetTotAreaCoin(k) > 0 {
					if v {
						winCoins += j.GetTotAreaCoin(k)
					} else {
						winCoins -= j.GetTotAreaCoin(k)
					}
				}
			case 1:
				if j.GetTotAreaCoin(k) > 0 {
					if v {
						winCoins += int64(float64(j.GetTotAreaCoin(k)) * 0.95)
					} else {
						winCoins -= j.GetTotAreaCoin(k)
					}
				}
			case 2:
				if j.GetTotAreaCoin(k) > 0 {
					if v {
						winCoins += (j.GetTotAreaCoin(k) * 8)
					} else {
						winCoins -= j.GetTotAreaCoin(k)
					}
				}
			case 3:
				if j.GetTotAreaCoin(k) > 0 {
					if v {
						winCoins += int64(float64(j.GetTotAreaCoin(k)) * 1.5)
					} else {
						winCoins -= j.GetTotAreaCoin(k)
					}
				}
			case 4:
				if j.GetTotAreaCoin(k) > 0 {
					if v {
						winCoins += int64(float64(j.GetTotAreaCoin(k)) * 0.5)
					} else {
						winCoins -= j.GetTotAreaCoin(k)
					}
				}
			case 5:
				if j.GetTotAreaCoin(k) > 0 {
					if v {
						winCoins += (j.GetTotAreaCoin(k) * 11)
					} else {
						winCoins -= j.GetTotAreaCoin(k)
					}
				}
			case 6:
				if j.GetTotAreaCoin(k) > 0 {
					if v {
						winCoins += (j.GetTotAreaCoin(k) * 11)
					} else {
						winCoins -= j.GetTotAreaCoin(k)
					}
				}
			case 7:
				if j.GetTotAreaCoin(k) > 0 {
					if v && !areaRes[8] {
						winCoins += int64(float64(j.GetTotAreaCoin(k)) * 0.97)
					} else if !v && areaRes[8] {
						winCoins -= j.GetTotAreaCoin(k)
					}
				}
			case 8:
				if j.GetTotAreaCoin(k) > 0 {
					if v && !areaRes[7] {
						winCoins += int64(float64(j.GetTotAreaCoin(k)) * 0.97)
					} else if !v && areaRes[7] {
						winCoins -= j.GetTotAreaCoin(k)
					}
				}
			}
		}
	}
	return winCoins
}
func (this *ExtDesk) getAreaResult() []bool {
	var pair byte
	if GetCardValue(this.IdleCard[0]) == GetCardValue(this.IdleCard[1]) {
		pair += IDLE << 4
	}
	if GetCardValue(this.BankerCard[0]) == GetCardValue(this.BankerCard[1]) {
		pair += BANKER << 4
	}
	chart := this.CardMgr.CompareCard(this.IdleCard, this.BankerCard)
	//this.RunChart = append(this.RunChart, int32(pair+chart))
	//
	var winArea = make([]bool, 9)
	var double [9]float64
	for i := range winArea {
		winArea[i] = false
		double[i] = 0
	}

	// 计算闲输赢
	var index int
	bankerIndex := INDEX_BANKER - 1
	idleIndex := INDEX_IDLE - 1
	bankerCoins := this.GetAreaCoin(bankerIndex)
	idleCoins := this.GetAreaCoin(idleIndex)

	switch chart {
	case IDLE:
		double[idleIndex] = gameConfig.Double[idleIndex]
		winArea[idleIndex] = true
		if idleCoins < bankerCoins {
			index = INDEX_BANKERWIN - 1
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else if idleCoins > bankerCoins {
			index = INDEX_BANKERLOSE - 1
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else {
			// 返还上庄
			index = INDEX_BANKERWIN - 1
			double[index] = 1
			winArea[index] = true
			index = INDEX_BANKERLOSE - 1
			double[index] = 1
			winArea[index] = true
		}
	case BANKER:
		double[bankerIndex] = gameConfig.Double[bankerIndex]
		winArea[bankerIndex] = true

		if bankerCoins < idleCoins {
			index = INDEX_BANKERWIN - 1
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else if bankerCoins > idleCoins {
			index = INDEX_BANKERLOSE - 1
			double[index] = gameConfig.Double[index]
			winArea[index] = true
		} else {
			// 返还上庄
			index = INDEX_BANKERWIN - 1
			double[index] = 1
			winArea[index] = true
			index = INDEX_BANKERLOSE - 1
			double[index] = 1
			winArea[index] = true
		}
	case DRAW:
		index = INDEX_DRAW - 1
		double[index] = gameConfig.Double[index]
		winArea[index] = true
		// 返还上庄
		index = INDEX_BANKERWIN - 1
		double[index] = 1
		winArea[index] = true
		index = INDEX_BANKERLOSE - 1
		double[index] = 1
		winArea[index] = true
		// 返回庄、闲
		double[idleIndex] = 1
		winArea[idleIndex] = true
		double[bankerIndex] = 1
		winArea[bankerIndex] = true
	}
	if len(this.IdleCard)+len(this.BankerCard) == 4 {
		index = INDEX_SMALL - 1
		double[index] = gameConfig.Double[index]
		winArea[index] = true
	} else {
		index = INDEX_BIG - 1
		double[index] = gameConfig.Double[index]
		winArea[index] = true
	}
	if pair&(IDLE<<4) > 0 {
		index = INDEX_IDLEPAIR - 1
		double[index] = gameConfig.Double[index]
		winArea[index] = true
	}
	if pair&(BANKER<<4) > 0 {
		index = INDEX_BANKERPAIR - 1
		double[index] = gameConfig.Double[index]
		winArea[index] = true
	}
	return winArea
}

// 获取盈利率
func (this *ExtDesk) GetWinRate() float64 {
	if this.totCoins == 0 {
		return float64(gameConfig.DeskInfo.Win)
	}

	oldRate := this.wCoins / this.totCoins * 100
	return float64(gameConfig.DeskInfo.Win)*2 - oldRate
}

// 获取最接近的值
func (this *ExtDesk) GetNealIndex(scale float64, winCard []GWinCard) int {
	if winCard == nil {
		return -1
	}

	for i, v := range winCard {
		if v.WinScale < scale {
			if i == 0 {
				return 0
			}

			if this.tCount < 100 || scale > float64(gameConfig.DeskInfo.Win) {
				return i - 1
			}

			return i
		}
	}

	return len(winCard) - 1
}

// 获取座位名单
func (this *ExtDesk) GetSeatInfo(ep *ExtPlayer) []GSInfo {
	SeatList := this.SeatMgr.GetSeatList()
	ret := []GSInfo{}
	for _, v := range SeatList {
		p := v.(*ExtPlayer)
		seat := GSInfo{
			Nick: p.Nick,
			Head: p.Head,
		}

		if len(p.Nick) > 4 && p.Uid != ep.Uid {
			seat.Nick = "***" + p.Nick[len(p.Nick)-4:]
		}
		ret = append(ret, seat)
		// var flag bool
		// for _, v := range ret {
		// 	if v.Nick == seat.Nick {
		// 		flag = true
		// 	}
		// }
		// if !flag {
		// 	ret = append(ret, seat)
		// }
	}
	return ret
}

func (this *ExtDesk) UpdatePlayer() {
	this.SeatMgr.OrderByBetCoins()
	this.SeatMgr.UpdateSeatList()

	this.gameUserListLK.Lock()
	this.Seat = []GSInfo{}
	SeatList := this.SeatMgr.GetSeatList()
	for _, v := range SeatList {
		p := v.(*ExtPlayer)
		seat := GSInfo{
			Nick: p.Nick,
			Head: p.Head,
		}
		this.Seat = append(this.Seat, seat)
	}
	this.gameUserListLK.Unlock()
}

func (this *ExtDesk) GetUserList(ep *ExtPlayer) []GUserInfo {
	this.gameUserListLK.RLock()
	defer this.gameUserListLK.RUnlock()

	ret := []GUserInfo{}
	userList := this.SeatMgr.GetUserList(gameConfig.DeskInfo.ListCount)
	for _, v := range userList {
		p := v.(*ExtPlayer)
		userInfo := GUserInfo{
			Uid:      p.Uid,
			Nick:     p.Nick,
			Head:     p.Head,
			TotBet:   p.GetBetCoins(),
			WinCount: p.GetWinCount(),
			Coins:    p.Coins,
		}

		if len(userInfo.Nick) > 4 && p.Uid != ep.Uid {
			userInfo.Nick = "***" + userInfo.Nick[len(userInfo.Nick)-4:]
		}

		ret = append(ret, userInfo)
	}

	return ret
}

// func (this *ExtDesk) UpdatePlayerInfo(p *ExtPlayer, d *DkInMsg) {
// 	msg := GUpdatePlayerInfo{}
// 	json.Unmarshal([]byte(d.Data), &msg)
// 	p.AddCoins(msg.PlayerInfo.Coins)
// 	//if msg.PlayerInfo.Account == "zheng001" {
// 	logs.Debug("updata-head: %v", msg.PlayerInfo.Head)
// 	logs.Debug("updata-head-len: %v", len(msg.PlayerInfo.Head))
// 	//}

// 	p.Head = msg.PlayerInfo.Head

// 	// 通知玩家，大厅部分有通知充值
// 	// p.SendNativeMsg(MSG_HALL_PUSH_CHANGECOIN, &PMsgToClientChangeCoin{
// 	// 	Id:   MSG_HALL_PUSH_CHANGECOIN,
// 	// 	Coin: p.Coins,
// 	// })
// }
func (this *ExtDesk) GetZouShi(id int32) {
	//
	rsp := GGetZouShiReply{
		Id: MSG_GAME_GETZOUSHI_REPLY,
	}
	t1, t2 := this.GetStageTime()
	trend := Trends{RunChart: this.RunChart, TypeTimes: this.TypeTimes, Time: t1}
	b, _ := json.Marshal(trend)
	rsp.Data.Data.ZouShi = string(b)
	rsp.Data.Data.GameStatus = int32(this.GameState)
	rsp.Data.Data.PlayerNum = int32(len(this.DeskMgr.MapPlayers))
	rsp.Data.GradeNumber = int32(GCONFIG.GradeNumber)
	rsp.Data.Data.MaxBet = this.GameLimit.High
	rsp.Data.Data.LowBet = this.GameLimit.Low
	//stageTime, tId := GetAllTime(AllStageTime)
	rsp.Data.Data.StageTime = GetAllTime(AllStageTime)
	rsp.Data.UpdateT = int64(this.GetTimerNum(t2))
	rsp.Data.SerId = this.GetServerId()
	rsp.Data.GameInfo.GameType = int32(GCONFIG.GameType)
	rsp.Data.GameInfo.RoomType = int32(GCONFIG.RoomType)
	rsp.Data.GameInfo.GradeType = int32(GCONFIG.GradeType)
	//
	this.DeskMgr.SendNativeMsgNoPlayer(MSG_GAME_GETZOUSHI_REPLY, 0, id, &rsp)
}

//获取阶段时间
func (this *ExtDesk) GetStageTime() (int, int) {
	switch this.GameState {
	case MSG_GAME_INFO_SHUFFLE_NOTIFY:
		return this.GetTimerNum(1), 1
	case MSG_GAME_INFO_READY_NOTIFY:
		return GetAllTime(AllStageTime[:1]) + this.GetTimerNum(2), 2
	case MSG_GAME_INFO_SEND_NOTIFY:
		return GetAllTime(AllStageTime[:2]) + this.GetTimerNum(3), 3
	case MSG_GAME_INFO_BET_NOTIFY:
		return GetAllTime(AllStageTime[:3]) + this.GetTimerNum(4), 4
	case MSG_GAME_INFO_STOP_BET_NOTIFY:
		return GetAllTime(AllStageTime[:4]) + this.GetTimerNum(5), 5
	case MSG_GAME_INFO_OPEN_NOTIFY:
		return GetAllTime(AllStageTime[:5]) + this.GetTimerNum(6), 6
	case MSG_GAME_INFO_AWARD_NOTIFY:
		return GetAllTime(AllStageTime[:6]) + this.GetTimerNum(7), 7
	default:
		return 0, 0
	}
}
func (this *ExtDesk) GetUserAreaCoins() int64 {
	ret := this.userArealistCoins.GetTotValue()
	return ret
}

package main

import (
	"sync"

	// "encoding/json"

	"bl.com/seatlist"
	"bl.com/util"
)

type ExtDesk struct {
	Desk

	sync.RWMutex
	CardMgr MgrCard // 扑克牌牌管理

	SeatMgr seatlist.MgrSeat // 座位管理

	arealistCoins     util.AreaList // 总下注情况
	userArealistCoins util.AreaList // 用户下注情况

	Count      int32   // 当前局数
	OpenCard   int32   // 展示出来的牌
	DragonCard int32   // 龙
	TigerCard  int32   // 虎
	WinArea    []bool  // 赢取区域
	RoomId     string  // 房号
	GameId     string  // 局号
	GameLimit  Limit   // 限红
	BetList    []int64 // 下注金币
	BetArea    []bool  // 可下注金币
	RunChart   []int32 // 走势

	NewBet bool

	gameUserListLK sync.RWMutex // 玩家列表读写锁
	Seat           []GSInfo     // 座位信息

	LeftCount  int32 // 剩余牌
	RightCount int32 // 废牌

	wCoins   float64 // 总赢取
	tCount   float64 // 总局数
	totCoins float64 // 总投注
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
	this.arealistCoins.Init(13)
	this.userArealistCoins.Init(13)
}

// 洗牌初始化信息
func (this *ExtDesk) ExtDeskInit() {
	this.Count = 0
	this.CardMgr.Shuffle()
	this.RunChart = []int32{}
	this.BetArea = []bool{true, true, true, true, true, true, true, true, true, true, true, true, true}
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
func (this *ExtDesk) GetUserAreaCoins() int64 {
	ret := this.userArealistCoins.GetTotValue()
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

func (this *ExtDesk) Predict(dragonCard int32, tigerCard int32) float64 {
	chart := this.CardMgr.CompareCard(dragonCard, tigerCard)

	// 计算闲输赢
	var loseCoins int64 = 0
	var index int

	// 龙虎和区域中奖
	switch chart {
	case DRAGON:
		index = INDEX_DRAGON - 1
	case TIGER:
		index = INDEX_TIGER - 1
	case DRAW:
		index = INDEX_DRAW - 1
	default:
		index = INDEX_ERROR - 1
	}
	coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	loseCoins += int64(coins)

	// 花色区域中奖
	dColor := GetCardColor(dragonCard)
	switch dColor {
	case CARD_COLOR_Fang:
		index = INDEX_DRAGONSPADE - 1
	case CARD_COLOR_Mei:
		index = INDEX_DRAGONPLUM - 1
	case CARD_COLOR_Hong:
		index = INDEX_DRAGONRED - 1
	case CARD_COLOR_Hei:
		index = INDEX_DRAGONBLOCK - 1
	default:
		index = INDEX_ERROR - 1
	}
	coins = float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	loseCoins += int64(coins)

	tColor := GetCardColor(tigerCard)
	switch tColor {
	case CARD_COLOR_Fang:
		index = INDEX_TIGERSPADE - 1
	case CARD_COLOR_Mei:
		index = INDEX_TIGERPLUM - 1
	case CARD_COLOR_Hong:
		index = INDEX_TIGERRED - 1
	case CARD_COLOR_Hei:
		index = INDEX_TIGERBLOCK - 1
	default:
		index = INDEX_ERROR - 1
	}
	coins = float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	loseCoins += int64(coins)

	dMidT := this.GetAreaCoin(INDEX_DRAGON-1) - this.GetAreaCoin(INDEX_TIGER-1)
	if dMidT == 0 {
		// 龙虎押注一样，上庄区域退还
		index = INDEX_BANKERWIN - 1
		loseCoins += this.GetUserAreaCoin(index)

		index = INDEX_BANKERLOSE - 1
		loseCoins += this.GetUserAreaCoin(index)
	}

	switch chart {
	case DRAGON:
		if dMidT < 0 {
			index = INDEX_BANKERWIN - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
		} else if dMidT > 0 {
			index = INDEX_BANKERLOSE - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
		}
	case TIGER:
		if dMidT < 0 {
			index = INDEX_BANKERLOSE - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
		} else if dMidT > 0 {
			index = INDEX_BANKERWIN - 1
			coins := float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
			loseCoins += int64(coins)
		}
	case DRAW:
		// 返还上庄
		index = INDEX_BANKERWIN - 1
		loseCoins += this.GetUserAreaCoin(index)

		index = INDEX_BANKERLOSE - 1
		loseCoins += this.GetUserAreaCoin(index)
	}

	return float64(loseCoins)
}

func (this *ExtDesk) IsRight(card1 int32, card2 int32) bool {
	wCoins := float64(this.GetUserAreaCoins())
	loseCoins1 := this.Predict(card1, card2)
	loseCoins2 := this.Predict(card2, card1)

	rate1 := (this.wCoins + wCoins - loseCoins1) / (this.totCoins + wCoins) * 100
	rate2 := (this.wCoins + wCoins - loseCoins2) / (this.totCoins + wCoins) * 100

	isChange := false
	if rate1 < rate2 {
		rate1, rate2 = rate2, rate1
		isChange = true
	}

	// 小赢取，高于概率设置。取小
	// 否则、取大
	if int(rate2) > gameConfig.DeskInfo.Win-5 {
		return isChange
	}

	// 否则取大
	return !isChange
}

// 获取制作牌
func (this *ExtDesk) BuildCards() {
	rand, _ := util.GetRandomNum(0, 20)
	rand = 100 - (10 - rand)
	if this.GetUserAreaCoins() == 0 || gameConfig.DeskInfo.Win <= 0 || GetLocalStock() > G_DbGetGameServerData.GameConfig.GoalStock*int64(rand) {
		// 用户未下注，正常流程
		this.DragonCard = this.CardMgr.SendOneCard()
		this.TigerCard = this.CardMgr.SendOneCard()
		return
	}

	card1 := this.CardMgr.SendOneCard()
	card2 := this.CardMgr.SendOneCard()

	if !this.IsRight(card1, card2) {
		card1, card2 = card2, card1
	}
	this.DragonCard = card1
	this.TigerCard = card2
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
// 	p.Head = msg.PlayerInfo.Head

// 	// 通知玩家，大厅部分有通知充值
// 	// p.SendNativeMsg(MSG_HALL_PUSH_CHANGECOIN, &PMsgToClientChangeCoin{
// 	// 	Id:   MSG_HALL_PUSH_CHANGECOIN,
// 	// 	Coin: p.Coins,
// 	// })
// }

// //定时器
// func (this *ExtDesk) DoTimer() {
// 	if len(this.TList) == 0 {
// 		return
// 	}
// 	nlist := []*Timer{}
// 	olist := []*Timer{}
// 	for _, v := range this.TList {
// 		v.T -= 100
// 		if v.T <= 0 {
// 			olist = append(olist, v)
// 		} else {
// 			nlist = append(nlist, v)
// 		}
// 	}
// 	this.TList = nlist
// 	for _, v := range olist {
// 		v.H(v.D)
// 	}
// }

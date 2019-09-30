package main

import (
	"sync"

	"bl.com/util"

	"bl.com/paigow"
	"bl.com/seatlist"
)

type ExtDesk struct {
	Desk

	sync.RWMutex
	CardMgr MgrCard // 扑克牌牌管理

	SeatMgr seatlist.MgrSeat // 座位管理

	arealistCoins     util.AreaList // 总下注情况
	userArealistCoins util.AreaList // 用户下注情况

	RoomId    string  // 房号
	GameId    string  // 局号
	GameLimit Limit   // 限红
	BetList   []int64 // 下注金币

	BankerCard []int32   // 庄牌
	IdleCard   [][]int32 // 闲牌(天、地、人)
	RunChart   [][]int32 // 走势(庄、天、地、人)
	TypeList   []int32   // 牌型(庄、天、地、人)
	WinArea    []bool    // 赢取区域

	// 两个骰子
	dices1 int
	dices2 int

	NewBet bool

	gameUserListLK sync.RWMutex // 玩家列表读写锁
	Seat           []GSInfo     // 座位信息

	wCoins   float64 // 总赢取
	tCount   float64 // 总局数
	totCoins float64 // 总投注

	totalStock int64 //总累计库存
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

	this.RunChart = [][]int32{}
	this.RunChart = append(this.RunChart, []int32{}, []int32{}, []int32{}, []int32{})

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
	this.arealistCoins.Init(8)
	this.userArealistCoins.Init(8)
}

// 洗牌初始化信息
func (this *ExtDesk) ExtDeskInit() {
	this.GameId = util.BuildGameId(GCONFIG.GameType)
	this.CardMgr.Shuffle()

	this.BankerCard = []int32{} // 庄牌
	this.IdleCard = [][]int32{} // 闲牌(天、地、人)
	this.TypeList = []int32{}
	this.IdleCard = append(this.IdleCard, []int32{}, []int32{}, []int32{})
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

func (this *ExtDesk) tryWinRate(cards []int32, pCards [][]int32, wCoins float64) float64 {
	var typeList []int32
	bType := paigow.GetCardsType(cards)
	typeList = append(typeList, bType)
	for i := 0; i < 3; i++ {
		pType := paigow.GetCardsType(pCards[i])
		typeList = append(typeList, pType)
	}

	// 计算输赢
	var loseCoins float64
	var index int

	// 区域输赢情况
	if typeList[0] == paigow.TYPE_ZHIZUN {
		index = INDEX_BANKER_ZHIZUN - 1
		loseCoins += float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	} else if typeList[0] == paigow.TYPE_TIAN {
		index = INDEX_BANKER_TIAN - 1
		loseCoins += float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]
	}

	// 天
	isWin := paigow.CompareCard(pCards[0], cards)
	if isWin {
		index = INDEX_TIAN_WIN - 1
	} else {
		index = INDEX_TIAN_LOSS - 1
	}
	loseCoins += float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]

	// 地
	isWin = paigow.CompareCard(pCards[1], cards)
	if isWin {
		index = INDEX_DI_WIN - 1
	} else {
		index = INDEX_DI_LOSS - 1
	}
	loseCoins += float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]

	// 人
	isWin = paigow.CompareCard(pCards[2], cards)
	if isWin {
		index = INDEX_REN_WIN - 1
	} else {
		index = INDEX_REN_LOSS - 1
	}
	loseCoins += float64(this.GetUserAreaCoin(index)) * gameConfig.Double[index]

	return (this.wCoins + wCoins - loseCoins) / (this.totCoins + wCoins) * 100
}

// 获取制作牌
// func (this *ExtDesk) BuildBCards() {
// 	rand, _ := util.GetRandomNum(0, 20)
// 	rand = 100 - (10 - rand)
// 	if this.GetUserAreaCoins() == 0 || gameConfig.DeskInfo.Win <= 0 || GetLocalStock() > G_DbGetGameServerData.GameConfig.GoalStock*int64(rand) {
// 		// 用户未下注，正常流程
// 		this.BankerCard = this.CardMgr.SendCard(2)
// 		this.BankerCard = paigow.Sort(this.BankerCard)
// 		return
// 	}

// 	var idleType []int32
// 	for i := 0; i < 3; i++ {
// 		pType := paigow.GetCardsType(this.IdleCard[i])
// 		idleType = append(idleType, pType)
// 	}

// 	var winList []GWinCard
// 	start := this.CardMgr.MSendId
// 	AllCards := this.CardMgr.MVSourceCard[start:]

// 	winCoin := float64(this.GetUserAreaCoins())
// 	// 用户已下注，进行管控
// 	for i := 0; i < len(AllCards)-1; i++ {
// 		for j := i + 1; j < len(AllCards); j++ {
// 			if i == j {
// 				continue
// 			}
// 			rate := this.tryWinRate([]int32{int32(AllCards[i]), int32(AllCards[j])}, this.IdleCard, winCoin)
// 			winCard := GWinCard{
// 				Index:    append([]int{}, i, j),
// 				WinScale: rate,
// 			}
// 			winList = append(winList, winCard)
// 		}
// 	}

// 	// 按照概率降序
// 	for i := 0; i < len(winList)-1; i++ {
// 		for j := i + 1; j < len(winList); j++ {
// 			if winList[i].WinScale < winList[j].WinScale {
// 				winList[i], winList[j] = winList[j], winList[i]
// 			}
// 		}
// 	}

// 	// 通过map主键唯一的特性过滤重复元素
// 	result := []GWinCard{}
// 	tempMap := map[float64]byte{} // 存放不重复主键
// 	for _, e := range winList {
// 		l := len(tempMap)
// 		tempMap[e.WinScale] = 0
// 		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
// 			result = append(result, e)
// 		}
// 	}

// 	// 查找概率最接近的
// 	index := this.GetNealIndex(result)
// 	if index < 0 {
// 		// 查找错误，，正常流程
// 		for i := 0; i < 2; i++ {
// 			this.BankerCard = append(this.BankerCard, this.CardMgr.SendOneCard())
// 		}
// 		this.BankerCard = paigow.Sort(this.BankerCard)
// 	} else {
// 		for _, v := range result[index].Index {
// 			this.BankerCard = append(this.BankerCard, this.CardMgr.SendCard(v))
// 		}
// 		this.BankerCard = paigow.Sort(this.BankerCard)
// 	}

// 	return
// }

// 获取最接近的值
func (this *ExtDesk) GetNealIndex(winCard []GWinCard) int {
	if winCard == nil {
		return -1
	}

	for i, v := range winCard {
		if v.WinScale < float64(gameConfig.DeskInfo.Win) {
			if i == 0 {
				return 0
			}

			if this.tCount < 100 {
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

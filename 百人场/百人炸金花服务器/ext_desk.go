package main

import (
	"logs"
	"sync"

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

	RoomId    string  // 房号
	GameId    string  // 局号
	GameLimit Limit   // 限红
	BetList   []int64 // 下注金币

	BankerCard   []int32   // 庄牌
	MBankerCard  []int32   // 庄最大牌
	TBankerType  int32     // 庄牌型
	IdleCard     [][]int32 // 闲牌(东西南北)
	MIdleCard    [][]int32 // 闲最大牌
	BankRunChart []int32   // 庄走势
	RunChart     [][]bool  // 走势(东西南北)
	TypeList     []int32   // 牌型
	WinArea      []bool    // 赢取区域

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
	//fmt.Println(gameConfig.LimitInfo.Limit)
	this.RoomId = util.BuildRoomId(GCONFIG.GradeType+1, this.Id+1)
	this.GameLimit = gameConfig.LimitInfo.Limit[GCONFIG.GradeType-1]
	this.BetList = gameConfig.LimitInfo.BetCoins[GCONFIG.GradeType-1].Bet[:]

	// 设置座位数
	this.SeatMgr.SetSeatNum(gameConfig.DeskInfo.SeatCount)

	this.NewBet = false

	this.BankRunChart = []int32{}
	this.RunChart = [][]bool{}
	this.RunChart = append(this.RunChart, []bool{}, []bool{}, []bool{}, []bool{})

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
	this.CardMgr.Shuffle()

	this.BankerCard = []int32{}  // 庄牌
	this.MBankerCard = []int32{} // 庄最大牌
	this.TBankerType = 0         // 庄牌型

	this.IdleCard = [][]int32{}  // 闲牌(东西南北)
	this.MIdleCard = [][]int32{} // 闲最大牌
	this.TypeList = []int32{}
	this.IdleCard = append(this.IdleCard, []int32{}, []int32{}, []int32{}, []int32{})

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

// 统计桌面金币，并分配牌，实现输赢控制器
func (this *ExtDesk) allotCard(wantwin bool) {
	logs.Debug("-----------------进入风控，控制输赢：", wantwin)
	//从剩余的牌中取出5组3张牌
	var IdleCard [][]int32
	for i := 0; i < 5; i++ {
		IdleCard = append(IdleCard, this.CardMgr.SendCard(3))
	}
	//取出各家手牌中的2张牌
	fachuCard := append([][]int32{}, this.IdleCard...)
	fachuCard = append(fachuCard, this.BankerCard)
	//组牌
	for k := 0; k < 5; k++ {
		var zucards [][]int32
		for j := 0; j < 5; j++ {
			a := append([]int32{}, fachuCard[j]...)
			a = append(a, IdleCard[j]...)
			zucards = append(zucards, a)
		}
		if wantwin == true { //控制赢
			if coin := this.IsWin(zucards[4], zucards[:4]); coin >= 0 {
				for i := 0; i < 4; i++ {
					this.IdleCard[i] = append(this.IdleCard[i], IdleCard[i]...)
					MCard, TCard := this.CardMgr.GetMaxCards(this.IdleCard[i])
					this.MIdleCard = append(this.MIdleCard, MCard)
					this.TypeList = append(this.TypeList, TCard)
				}
				// 庄家牌
				this.BankerCard = append(this.BankerCard, IdleCard[4]...)
				this.MBankerCard, this.TBankerType = this.CardMgr.GetMaxCards(this.BankerCard)
				logs.Debug("控制庄家赢：", coin)
				return
			}
		} else { //控制输
			if coin := this.IsWin(zucards[4], zucards[:4]); coin < 0 {
				for i := 0; i < 4; i++ {
					this.IdleCard[i] = append(this.IdleCard[i], IdleCard[i]...)
					MCard, TCard := this.CardMgr.GetMaxCards(this.IdleCard[i])
					this.MIdleCard = append(this.MIdleCard, MCard)
					this.TypeList = append(this.TypeList, TCard)
				}
				// 庄家牌
				this.BankerCard = append(this.BankerCard, IdleCard[4]...)
				this.MBankerCard, this.TBankerType = this.CardMgr.GetMaxCards(this.BankerCard)
				logs.Debug("控制庄家输：", coin)
				return
			}
		}

		//换牌
		IdleCard = append(IdleCard[1:], IdleCard[:1]...)
	}
	//找不到的wantwin牌组就随机赋值
	for i := 0; i < 4; i++ {
		this.IdleCard[i] = append(this.IdleCard[i], IdleCard[i]...)
		MCard, TCard := this.CardMgr.GetMaxCards(this.IdleCard[i])
		this.MIdleCard = append(this.MIdleCard, MCard)
		this.TypeList = append(this.TypeList, TCard)
	}
	// 庄家牌
	this.BankerCard = append(this.BankerCard, IdleCard[4]...)
	this.MBankerCard, this.TBankerType = this.CardMgr.GetMaxCards(this.BankerCard)
	return
}

//计算庄家输赢，返回庄家赢的金额（可为负数）
func (this *ExtDesk) IsWin(cards []int32, pCards [][]int32) int64 {
	//玩家赢的金额
	bankerAllWin := float64(this.GetUserAreaCoins())
	//计算所有区域的庄家输赢情况
	for i := 0; i < 4; i++ {
		if this.CardMgr.CompareCard(pCards[i][:], cards) {
			bankerAllWin -= float64(this.GetUserAreaCoin(i)) * gameConfig.Double[i]
		}
	}
	//判断庄家下注区输赢
	C, btype := this.CardMgr.GetMaxCards(cards)
	switch byte(btype) {
	case this.CardMgr.GFDouble:
		if this.CardMgr.GetLogicValue(C[1]) < util.Card_8 && this.CardMgr.GetLogicValue(C[1]) != util.Card_A {
			break
		}
		i := INDEX_BANKER_DOUBLE - 1
		bankerAllWin -= float64(this.GetUserAreaCoin(i)) * gameConfig.Double[i]
	case this.CardMgr.GFShunZi:
		i := INDEX_BANKER_SHUNZI - 1
		bankerAllWin -= float64(this.GetUserAreaCoin(i)) * gameConfig.Double[i]
	case this.CardMgr.GFJinHua:
		if this.CardMgr.GetLogicValue(C[0]) < util.Card_10 && this.CardMgr.GetLogicValue(C[0]) != util.Card_A {
			break
		}
		i := INDEX_BANKER_JINHUA - 1
		bankerAllWin -= float64(this.GetUserAreaCoin(i)) * gameConfig.Double[i]
	case this.CardMgr.GFShunJin:
		if this.CardMgr.GetLogicValue(C[0]) < util.Card_8 && this.CardMgr.GetLogicValue(C[0]) != util.Card_A {
			break
		}
		i := INDEX_BANKER_SHUNJIN - 1
		bankerAllWin -= float64(this.GetUserAreaCoin(i)) * gameConfig.Double[i]
	case this.CardMgr.GFBaoZi:
		if this.CardMgr.GetLogicValue(C[0]) < util.Card_8 && this.CardMgr.GetLogicValue(C[0]) != util.Card_A {
			break
		}
		i := INDEX_BANKER_BAOZI - 1
		bankerAllWin -= float64(this.GetUserAreaCoin(i)) * gameConfig.Double[i]
	}
	return int64(bankerAllWin)

}

func (this *ExtDesk) GetBankCard(maxCard []int32, allCards []byte) {

	AllCards := append([]byte{}, allCards...)

	for i := 0; i < len(AllCards)-2; i++ {
		for j := i + 1; j < len(AllCards)-1; j++ {
			for k := j + 1; k < len(AllCards); k++ {
				cards := append(this.BankerCard, int32(AllCards[i]), int32(AllCards[j]), int32(AllCards[k]))
				bankerCard, bankerType := this.CardMgr.GetMaxCards(cards)
				if this.CardMgr.CompareCard(bankerCard, maxCard) {
					cards = this.CardMgr.Sort(cards)
					logs.Debug("庄最大牌", cards)
					this.BankerCard = cards
					this.MBankerCard, this.TBankerType = bankerCard, bankerType
					return
				}
			}
		}
	}

}

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

	return 0
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

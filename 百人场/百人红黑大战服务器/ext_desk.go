package main

import (
	"logs"
	// "encoding/json"
	"sync"

	"bl.com/seatlist"
	"bl.com/util"
)

type ExtDesk struct {
	Desk
	sync.RWMutex
	CardMgr           MgrCard          //扑克牌牌管理
	SeatMgr           seatlist.MgrSeat //座位管理
	arealistCoins     util.AreaList    //总下注情况
	userArealistCoins util.AreaList    //用户下注情况

	Count          int32   //当前局数
	CardList       []int32 //展示出来的牌
	RedCard        []int32 //红方牌
	BlackCard      []int32 //黑方牌
	WinArea        []bool  //赢取区域
	RoomId         string  //房号
	GameId         string  //局号
	GameLimit      int64   //限红
	BetList        []int64 //下注金币
	BetArea        []bool  //可下注金币区域
	RunChart       []int32 //输赢走势
	CardTypeChart  []int   //牌型记录走势
	betId          int32   //获取红黑方下注的ID
	WinId          int32   //获取红黑方赢取的ID
	NewBet         bool
	UserCount      int          //玩家人数
	gameUserListLK sync.RWMutex //玩家列表读写锁
	Seat           []GSInfo     //座位信息

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
	//游戏各个阶段
	this.Handle[MSG_GAME_AUTO] = this.HandleGameAuto               //自由匹配
	this.Handle[MSG_GAME_INFO_BET] = this.HandleGameBet            //处理游戏下注
	this.Handle[MSG_GAME_INFO_RUN_CHART] = this.HandleGameRunChart //处理游戏走势图
	this.Handle[MSG_GAME_INFO_USER_LIST] = this.HandleGameUserList //处理玩家列表
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect         //处理断线重连
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnet        //处理用户掉线
	this.Handle[MSG_GAME_INFO_INTO] = this.HandleGameAutoFinal     //自由匹配后续处理
	this.Handle[MSG_GAME_INFO_EXIT] = this.HandleGameExit          //玩家退出游戏
	this.Handle[MSG_GAME_INFO_GET_RECORD] = this.HandleGetRecord   //游戏记录

	//创建房间号
	this.RoomId = util.BuildRoomId(GCONFIG.GradeType+1, this.Id+1)
	//游戏限红
	this.GameLimit = G_DbGetGameServerData.GameConfig.LimitRedMax
	//	下注限制
	this.BetList = G_DbGetGameServerData.GameConfig.TenChips
	//	设置座位数
	this.SeatMgr.SetSeatNum(gameConfig.DeskInfo.SeatCount)

	this.NewBet = false
	this.ResetAreaCoins()  //重置下注列表
	this.TimerShuffle(nil) //开始洗牌
}

//广播
func (this *ExtDesk) BroadcastAll(id int, d interface{}) {
	allUser := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range allUser {
		v.(*ExtPlayer).SendNativeMsg(id, d)
	}
}

//广播-特殊
func (this *ExtDesk) BroadcastAllSpec(id int, d interface{}) {
	o := d.(GGameBetNotify)
	allUser := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range allUser {
		o.BetArea = v.(*ExtPlayer).PBetArea
		v.(*ExtPlayer).SendNativeMsg(id, o)
	}
}

//重置下注列表
func (this *ExtDesk) ResetAreaCoins() {
	this.arealistCoins.Init(3)
	this.userArealistCoins.Init(3)
	// this.Limitcoinid = int32(9)
}

//洗牌初始化信息
func (this *ExtDesk) ExtDeskInit() {
	this.Count = 0
	this.CardMgr.Shuffle()
	this.BetArea = []bool{true, true, true}

}

//获得下注列表
func (this *ExtDesk) GetAreaCoinsList() []int64 {
	ret := this.arealistCoins.GetValueList()
	return ret
}

//获取下注金币
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

//添加下注
func (this *ExtDesk) AddAreaCoins(area int, coins int64) bool {
	//添加区域下注
	this.NewBet = true
	ret := this.arealistCoins.AddValue(area, coins)
	return ret
}
func (this *ExtDesk) AddUserAreaCoins(area int, coins int64) bool {
	/*添加玩家下注区域*/
	this.NewBet = true
	ret := this.userArealistCoins.AddValue(area, coins)
	return ret
}

//获取座位名单
func (this *ExtDesk) GetSeatInfo(ep *ExtPlayer) []GSInfo {
	SeatList := this.SeatMgr.GetSeatList()
	ret := []GSInfo{}
	for _, v := range SeatList {
		p := v.(*ExtPlayer)
		seat := GSInfo{
			Nick:  p.Nick,
			Head:  p.Head,
			Coins: p.Coins,
		}
		if len(p.Nick) > 4 && p.Uid != ep.Uid {
			seat.Nick = "***" + p.Nick[len(p.Nick)-4:]
		}
		ret = append(ret, seat)
	}
	return ret
}

//更新游戏玩家
func (this *ExtDesk) UpdatePlayer() {
	this.SeatMgr.OrderByBetCoins()
	this.SeatMgr.UpdateSeatList()

	this.gameUserListLK.Lock()
	this.Seat = []GSInfo{}                 //获取座位信息
	SeatList := this.SeatMgr.GetSeatList() //获取座位列表
	for _, v := range SeatList {
		p := v.(*ExtPlayer)
		seat := GSInfo{
			Nick:  p.Nick,
			Head:  p.Head,
			Coins: p.Coins,
		}
		this.Seat = append(this.Seat, seat)
	}
	this.gameUserListLK.Unlock()
}

//获取用户列表
func (this *ExtDesk) GetUserList(ep *ExtPlayer) []GUserInfo {
	this.gameUserListLK.RLock()
	defer this.gameUserListLK.RUnlock()

	ret := []GUserInfo{}
	userList := this.SeatMgr.GetUserList(gameConfig.DeskInfo.ListCount)
	this.UserCount = len(userList)
	for _, v := range userList {
		p := v.(*ExtPlayer)
		userInfo := GUserInfo{
			Uid:       p.Uid,
			Nick:      p.Nick,
			Head:      p.Head,
			TotBet:    p.GetBetCoins(),
			WinCount:  p.GetWinCount(),
			Coins:     p.Coins,
			UserCount: len(userList),
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

// 通知玩家，大厅部分有通知充值
// p.SendNativeMsg(MSG_HALL_PUSH_CHANGECOIN, &PMsgToClientChangeCoin{
// 	Id:   MSG_HALL_PUSH_CHANGECOIN,
// 	Coin: p.Coins,
// })
// }
func (this *ExtDesk) BuildCard(i int) ([]int32, []int32) {
	card1, card2 := this.Licensing()
	chart, typechart := this.changeCard(card1, card2)
	logs.Debug("走势：%v;牌型走势:%v", chart, typechart)
	var redCard []int32
	var blackCard []int32
	switch i {
	case 0:
		if chart == 1 {
			redCard = card2
			blackCard = card1
		} else {
			redCard = card1
			blackCard = card2
		}
	case 1:
		if chart == 2 {
			redCard = card2
			blackCard = card1
		} else {
			redCard = card1
			blackCard = card2
		}
	case 2:
		if typechart >= 2 {
			for {
				card1, card2 := this.Licensing()
				Rcard, Rcolor := SortHandCard(card1) //将红方的牌进行排序并分出花色，数值
				Bcard, Bcolor := SortHandCard(card2) //将黑方的牌进行排序并分出花色，数值
				RGrade, _ := GetCardType(Rcard, Rcolor)
				BGrade, _ := GetCardType(Bcard, Bcolor)
				if RGrade == 1 && BGrade == 1 {
					redCard = card1
					blackCard = card2
					break
				}
			}
		} else {
			redCard = card1
			blackCard = card2
		}
	default:
		return this.RedCard, this.BlackCard
	}
	return redCard, blackCard
}

func (this *ExtDesk) changeCard(card1 []int32, card2 []int32) (int32, int) {
	Rcard, Rcolor := SortHandCard(card1) //将红方的牌进行排序并分出花色，数值
	Bcard, Bcolor := SortHandCard(card2) //将黑方的牌进行排序并分出花色，数值
	chart, typechart := this.CardMgr.CompareCard(Rcard, Rcolor, Bcard, Bcolor)
	return chart, typechart
}

func (this *ExtDesk) Licensing() ([]int32, []int32) {
	this.CardMgr.Shuffle()
	cards := this.CardMgr.HandCardInfo(6)
	var cardList []int32
	for _, v := range cards {
		cardList = append(cardList, int32(v))
	}
	card1 := cardList[:3]
	card2 := cardList[3:]
	return card1, card2
}

func (this *ExtDesk) Predict(chart int32, typeChart int) ([]bool, []int64) {
	logs.Debug("********游戏预测结算********")
	var winArea [3]bool
	var tWinArea [3]int64
	var double [3]float64
	for i := range winArea {
		winArea[i] = false
		double[i] = 0
		tWinArea[i] = 0
	}
	//计算输赢
	var index int
	//红黑区域中奖
	switch chart {
	case RED:
		index = INDEX_RED - 1
		winArea[0] = true
	case BLACK:
		index = INDEX_BLACK - 1
		winArea[1] = true
	default:
		index = INDEX_ERROR - 1
	}
	this.WinId = int32(index)
	double[index] = gameConfig.Double[index]

	//根据牌类型计算输赢金额
	switch typeChart {
	case CARD_BOOM:
		double[2] = gameConfig.Double[2]
		winArea[2] = true
	case CARD_FLUSH:
		double[2] = gameConfig.Double[3]
		winArea[2] = true
	case CARD_TONGHUA:
		double[2] = gameConfig.Double[4]
		winArea[2] = true
	case CARD_SHUNZI:
		double[2] = gameConfig.Double[5]
		winArea[2] = true
	case CARD_PAIR:
		if this.CardMgr.CardValue >= 9 {
			double[2] = gameConfig.Double[6]
			winArea[2] = true
		} else {
			double[2] = 0
			winArea[2] = false
		}
	default:
		double[2] = 0
		winArea[2] = false
	}
	if this.betId == INDEX_LUCKYBLOW-1 {
		switch typeChart {
		case CARD_BOOM:
			double[2] = gameConfig.Double[2]
			winArea[2] = true
		case CARD_FLUSH:
			double[2] = gameConfig.Double[3]
			winArea[2] = true
		case CARD_TONGHUA:
			double[2] = gameConfig.Double[4]
			winArea[2] = true
		case CARD_SHUNZI:
			double[2] = gameConfig.Double[5]
			winArea[2] = true
		case CARD_PAIR:
			if this.CardMgr.CardValue >= 9 {
				double[2] = gameConfig.Double[6]
				winArea[2] = true
			} else {
				double[2] = 0
				winArea[2] = false
			}
		default:
			double[2] = 0
			winArea[2] = false
		}
	}
	for i, v := range winArea {
		if !v {
			continue
		}
		coins := int64(float64(this.GetUserAreaCoin(i)) * double[i])

		tWinArea[i] = coins
	}
	players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range players {
		p := v.(*ExtPlayer)
		p.BuildWinList(double[:])
	}
	return winArea[:], tWinArea[:]
}

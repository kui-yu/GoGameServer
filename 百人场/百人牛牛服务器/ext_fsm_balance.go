/**
* 结算状态
**/
package main

import (
	"encoding/json"
	// "encoding/json"
	// "fmt"
	"logs"
	"time"
)

type FSMBalance struct {
	UpMark int
	Mark   int
	EDesk  *ExtDesk

	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMBalance) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMBalance) GetMark() int {
	return this.Mark
}

func (this *FSMBalance) Run(upMark int) {
	DebugLog("游戏状态：结算")
	logs.Debug("游戏状态：结算")

	balanceId := gameConfig.GameStatusTimer.BalanceId
	balanceMs := int64(gameConfig.GameStatusTimer.BalanceMS)

	this.UpMark = upMark
	this.EndDateTime = GetTimeMS() + balanceMs

	this.addListener()                             // 添加监听
	this.EDesk.SendGameState(this.Mark, balanceMs) // 发送桌子状态

	this.EDesk.AddTimer(balanceId, int(balanceMs/1000), this.TimerCall, nil)

	this.balance() // 结算
	this.EDesk.RefreshManyUsers()
	this.EDesk.RefreshRank()
}

func (this *FSMBalance) Leave() {
	this.removeListener()

	// this.EDesk.DownCardIdx = 0
	this.EDesk.DownCards = []uint8{}
	this.EDesk.CardGroupArray = make(map[int]CardGroupInfo)

	// 清除座位下注
	for k, seat := range this.EDesk.Seats {
		seat.DownBetValue = 0
		seat.UserBetValue = 0
		this.EDesk.Seats[k] = seat
	}

	plen := len(this.EDesk.Players)
	for i := 0; i < plen; i++ {
		p := this.EDesk.Players[i]

		p.DownBets = make(map[uint8]int64)
		p.BalaDownBets = make(map[uint8]int64)
	}
}

func (this *FSMBalance) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_SHUFFLECARD)
}

func (this *FSMBalance) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()

	return remainTimeMS
}

// 添加网络监听
func (this *FSMBalance) addListener() {
	this.EDesk.Handle[MSG_GAME_QDOWNBET] = this.recvRSeatBet
}

// 接收到玩家下注
func (this *FSMBalance) recvRSeatBet(p *ExtPlayer, d *DkInMsg) {
	req := GClientQDownBet{}
	json.Unmarshal([]byte(d.Data), &req)
	this.EDesk.UserDownBet(p, req.SeatIdx, req.CoinIdx, true)
}

// 删除网络监听
func (this *FSMBalance) removeListener() {
}

func (this *FSMBalance) onUserOnline(p *ExtPlayer) {

}

func (this *FSMBalance) onUserOffline(p *ExtPlayer) {

}

// 结算
func (this *FSMBalance) balance() {
	logs.Debug("正在结算中！！！！！！！！！！！！！！")

	//发送结算消息给数据库, 简单记录
	// var activeUid int64 = 0
	if len(this.EDesk.Players) != 0 {
		for _, v := range this.EDesk.Players {
			if !v.Robot {
				// activeUid = v.Uid
			}
		}
	}
	//发送到大厅走势
	GCHH := GClientHallHistory{
		HomeName: GCONFIG.GradeType,
		HomeOdds: 1,
		// Trend:    this.EDesk.Gchh.Trend,
	}
	GCHH.LimitRed = gameConfig.GameLimtInfo.AreaMaxCoin
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.EDesk.JuHao,
		Mini:        false,
		SetLeave:    1,
		// ActiveUid:   activeUid,
	}

	BankerCard := []int32{}
	// 座位数量
	seatCount := gameConfig.GameLimtInfo.SeatCount
	// 座位历史数量
	runchartCount := gameConfig.GameLimtInfo.RunchartCount //走势条数

	for _, v := range this.EDesk.CardGroupArray[seatCount].Cards {
		BankerCard = append(BankerCard, int32(v))
	}
	IdleCard := [][]int32{}
	for i := 0; i < seatCount; i++ {
		cards := []int32{}
		for _, v := range this.EDesk.CardGroupArray[i].Cards {
			cards = append(cards, int32(v))
		}
		IdleCard = append(IdleCard, cards)
	}

	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.EDesk.JuHao,
		BankerCard:  BankerCard,
		IdleCard:    IdleCard,
	}

	// 计算玩家和庄家的输赢
	zcardInfo := this.EDesk.CardGroupArray[seatCount] //庄家的牌
	zCardType := zcardInfo.CardGroupType              //庄家的牌类型
	zCardMax := zcardInfo.MaxCard                     //庄家的最大牌
	//添加游戏记录:
	endTime := time.Now().Format("2006-01-02 15:04:05")
	var roomName string
	if GCONFIG.GradeType == 1 {
		roomName = "荣耀厅"
	} else if GCONFIG.GradeType == 2 {
		roomName = "王牌厅"
	} else if GCONFIG.GradeType == 3 {
		roomName = "战神厅"
	} else {
		roomName = "体验厅"
	}

	//获取玩家牌型
	var results map[int]GClientSeatBalance = make(map[int]GClientSeatBalance) //创建每一个位置的结算集合

	// 修改当前座位信息
	for id, seatInfo := range this.EDesk.Seats {
		seatCardType := this.EDesk.CardGroupArray[id].CardGroupType
		seatInfo.TrendHistory = append(seatInfo.TrendHistory, seatCardType)

		if seatInfo.UserId != 0 { //如果该座位上有存在玩家，那么该作为的玩家坐下次数+1
			seatInfo.SeatDownCount++
		}

		if len(seatInfo.TrendHistory) > runchartCount { //如果该作为的历史数量 已经大于走势条数限制，则剔除最靠前的历史走势
			seatInfo.TrendHistory = seatInfo.TrendHistory[1:]
		}

		this.EDesk.Seats[id] = seatInfo //重新更新桌子座位信息Seats

		results[id] = GClientSeatBalance{
			Bottom:   seatInfo.DownBetValue, //赋值用户下注金额
			Result:   0,
			MyBottom: 0,
			MyResult: 0,
		}
	}
	xianCardTs := []int{}
	for i := 0; i < gameConfig.GameLimtInfo.SeatCount; i++ {
		xianCardTs = append(xianCardTs, -1)
	}
	for i, v := range this.EDesk.CardGroupArray {
		if i != 4 {
			xianCardTs[i] = int(v.CardGroupType)
		}
	}
	for _, p := range this.EDesk.Players {
		//玩家游戏记录
		re := RecordData{
			MatchNum:       this.EDesk.JuHao,
			RoomName:       roomName,
			EndTime:        endTime,
			Date:           time.Now().Unix(),
			ZhuangCardType: int(zCardType),
		}
		re.XianCardType = []int{}
		for i := 0; i < gameConfig.GameLimtInfo.SeatCount; i++ {
			re.XianCardType = append(re.XianCardType, -1)
		}
		for i, v := range this.EDesk.CardGroupArray {
			if i != 4 {
				re.XianCardType[i] = int(v.CardGroupType)
			}
		}
		pdownBetTotal := p.getDownBetTotal()   //获取玩家下注总额
		coinsBefore := p.Coins + pdownBetTotal // 下注前金币

		if pdownBetTotal == 0 { //如果该玩家本局未下注，那么该玩家未下注数量+1
			p.UnbetsCount += 1
		} else { //否则重置未下注局数
			p.UnbetsCount = 0
		}

		var rdDownBets [4]int32 = [4]int32{}

		_, isSeatDown := this.EDesk.findUserSeatDown(p.Uid) //获取玩家是否在座位上
		var coinChange int64 = 0
		// logs.Debug("庄家牌组:", zhuanhuan(OrderCard(zcardInfo.Cards)))
		// logs.Debug("庄家最大牌:", zCardMax&0xF)
		for k, result := range results {
			//logs.Debug("现在是座位:", k, "的比拼！！！！！！！！！！！！！！！！！")
			if value, ok := p.DownBets[uint8(k)]; ok {
				result.MyBottom = int64(value)
				//cards := this.EDesk.CardGroupArray[k].Cards
				groupType := this.EDesk.CardGroupArray[k].CardGroupType
				maxCard := this.EDesk.CardGroupArray[k].MaxCard
				//logs.Debug("座位牌组:", zhuanhuan(OrderCard(cards)))
				//logs.Debug("座位最大牌:", maxCard&0xF)
				if zCardType == CardGroupType_NotCattle { // 庄家没有牛
					// logs.Debug("庄家没有牛")
					if groupType == CardGroupType_NotCattle { // 玩家没有牛
						// logs.Debug("玩家也没有牛")
						if isSeatDown { // 在座位上
							// logs.Debug("玩家在座位上..")
							GCHH.LimitRed = gameConfig.GameLimtInfo.AreaMaxCoinDownSeat
							if (zCardMax&0xF) > (maxCard&0xF) || ((zCardMax&0xF) == (maxCard&0xF) && (zCardMax>>4) > (maxCard>>4)) {
								// logs.Debug("玩家比牌比输了！！！！")
								result.MyResult = -result.MyBottom * int64(BetDoubleMap[zCardType])
							} else {
								// logs.Debug("玩家比牌比赢了！！！！")
								result.MyResult = result.MyBottom * int64(BetDoubleMap[groupType])
							}
							//logs.Debug("结算的结果是：", result.MyResult)
						} else { // 不在座位上
							// logs.Debug("玩家不在座位上..")
							GCHH.LimitRed = gameConfig.GameLimtInfo.AreaMaxCoin
							if ((maxCard & 0xF) < 11) || (zCardMax&0xF) > (maxCard&0xF) || ((zCardMax&0xF) == (maxCard&0xF) && (zCardMax>>4) > (maxCard>>4)) {
								// logs.Debug("玩家比牌比输了！！！！")
								result.MyResult = -result.MyBottom * int64(BetDoubleMap[zCardType])
							} else {
								// logs.Debug("玩家比牌比赢了！！！！")
								result.MyResult = result.MyBottom * int64(BetDoubleMap[groupType])
							}
							//logs.Debug("结算的结果是：", result.MyResult)
						}
					} else { // 玩家有牛
						// logs.Debug("玩家有牛")
						// logs.Debug("玩家比牌比赢了！！！！")
						result.MyResult = result.MyBottom * int64(BetDoubleMap[groupType])
						//logs.Debug("结算的结果是：", result.MyResult)
					}
				} else { // 庄家有牛
					// logs.Debug("庄家有牛")
					if groupType == CardGroupType_NotCattle { // 玩家没有牛
						// logs.Debug("玩家没牛")
						// logs.Debug("玩家比牌比输了！！！！")
						result.MyResult = -result.MyBottom * int64(BetDoubleMap[zCardType])
					} else { // 玩家有牛
						// logs.Debug("玩家有牛")
						if (zCardType > groupType) ||
							(zCardType == groupType &&
								((zCardMax&0xF) > (maxCard&0xF) ||
									(((zCardMax & 0xF) == (maxCard & 0xF)) &&
										((zCardMax >> 4) > (maxCard >> 4))))) {
							// logs.Debug("玩家比牌比输了！！！！")
							result.MyResult = -result.MyBottom * int64(BetDoubleMap[zCardType])
						} else {
							// logs.Debug("玩家比牌比赢了！！！！")
							result.MyResult = result.MyBottom * int64(BetDoubleMap[groupType])
						}
					}
				}
				result.Result += result.MyResult
			} else {
				result.MyBottom = 0
				result.MyResult = 0
			}
			p.BalaDownBets[uint8(k)] = result.MyResult
			results[k] = result
			rdDownBets[k] = int32(result.MyBottom)
			coinChange += (result.MyBottom + result.MyResult)
		}
		p.Coins += coinChange // 输赢的金币
		for i := 0; i < gameConfig.GameLimtInfo.SeatCount; i++ {
			re.BetArea = append(re.BetArea, 0)
		}
		for i, v := range p.DownBets {
			re.AllBet += v
			if v > 0 {
				re.BetArea[i] = v
			}
		}
		re.WinOrLost = coinChange - int64(pdownBetTotal)

		// logs.Debug("查看各方赢取数量：", GCHH.WinCount)
		// logs.Debug("查看发送大厅走势:", GCHH)
		// zs, err := json.Marshal(GCHH)
		// if err != nil {
		// 	logs.Error("更新走势时json转换出错")
		// }
		// this.EDesk.DeskMgr.SetZouShi(&GameZouShi{
		// 	SerId: int32(G_DbGetGameServerData.Sid),
		// 	GameInfo: GameTypeDetail{GameType: int32(GCONFIG.GameType),
		// 		RoomType:  int32(GCONFIG.RoomType),
		// 		GradeType: int32(GCONFIG.GradeType),
		// 	},
		// 	ZouShi:    string(zs),
		// 	PlayerNum: int32(len(this.EDesk.Players)),
		// 	UpdateT:   time.Now().Unix(),
		// })
		if re.AllBet > 0 {
			//logs.Debug("该玩家有下注,将其存起来")
			G_AllRecord.AddRecord(p.Uid, &re)
		} else {
			//logs.Debug("该玩家本局没有下注，所以不保存")
		}
		if pdownBetTotal != 0 {
			p.addBetHistory(coinChange > 0, pdownBetTotal)
		}
		// 发送记录到数据库
		if GetCostType() == 1 {
			if pdownBetTotal != 0 {
				validBet := pdownBetTotal // 有效下注，
				if coinChange < 0 {
					validBet = -coinChange
				}

				dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
					UserId:      p.Uid,
					UserAccount: p.Account,
					BetCoins:    int64(pdownBetTotal),
					ValidBet:    validBet,
					PrizeCoins:  coinChange - int64(pdownBetTotal),
					Robot:       p.Robot,
				})
				if !p.Robot {
					p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
				}
				dbreq.UserCoin = []GGameEndInfo{}
			}
		}
		if GetCostType() == 1 {
			if !p.Robot && pdownBetTotal != 0 {
				rddata := GGameRecordInfo{
					UserId:      p.Uid,
					UserAccount: p.Account,
					BetCoins:    pdownBetTotal, // 下注的金币
					BetArea:     rdDownBets,    // 区域下注情况
					PrizeCoins:  coinChange,    // 赢取的金币
					CoinsAfter:  p.Coins,       // 结束后金币
					CoinsBefore: coinsBefore,   // 下注前金币
					Robot:       p.Robot,
				}
				rdreq.UserRecord = append(rdreq.UserRecord, rddata)

				// 发送记录到存储
				p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
				rdreq.UserRecord = []GGameRecordInfo{}
			}
		}
		logs.Debug("玩家金币：%v", p.Coins)
	}
	// logs.Debug("走到这！！！！！！！！！！！！！！！！！！！！！")
	// logs.Debug("庄家牌类型：%v", int(zCardType))
	// logs.Debug("闲家牌类型：%v", xianCardTs)
	if int(zCardType) > 0 && int(zCardType) <= 12 {
		for i, _ := range xianCardTs {
			if xianCardTs[i] > int(zCardType) && xianCardTs[i] < 14 {
				if len(this.EDesk.Gchh.Trend) == 0 || len(this.EDesk.Gchh.Trend) == i {
					this.EDesk.Gchh.Trend = append(this.EDesk.Gchh.Trend, []int{1})
					// logs.Debug("走势#############################:%v", this.EDesk.Gchh.Trend)
				} else {
					this.EDesk.Gchh.Trend[i] = append(this.EDesk.Gchh.Trend[i], 1)
				}
			} else {
				if len(this.EDesk.Gchh.Trend) == 0 || len(this.EDesk.Gchh.Trend) == i {
					this.EDesk.Gchh.Trend = append(this.EDesk.Gchh.Trend, []int{0})
					// logs.Debug("走势#############################:%v", this.EDesk.Gchh.Trend)
				} else {
					this.EDesk.Gchh.Trend[i] = append(this.EDesk.Gchh.Trend[i], 0)
				}
			}
		}
	} else if int(zCardType) > 12 {
		for i, _ := range xianCardTs {
			if xianCardTs[i] == 14 {
				if len(this.EDesk.Gchh.Trend) == 0 || len(this.EDesk.Gchh.Trend) == i {
					this.EDesk.Gchh.Trend = append(this.EDesk.Gchh.Trend, []int{0})
					// logs.Debug("走势#############################:%v", this.EDesk.Gchh.Trend)
				} else {
					this.EDesk.Gchh.Trend[i] = append(this.EDesk.Gchh.Trend[i], 0)
				}
			} else {
				if len(this.EDesk.Gchh.Trend) == 0 || len(this.EDesk.Gchh.Trend) == i {
					this.EDesk.Gchh.Trend = append(this.EDesk.Gchh.Trend, []int{1})
					// logs.Debug("走势#############################:%v", this.EDesk.Gchh.Trend)
				} else {
					this.EDesk.Gchh.Trend[i] = append(this.EDesk.Gchh.Trend[i], 1)
				}
			}
		}
	}
	// logs.Debug("~~~~~~~~~~~~~~~~~~走势数组", this.EDesk.Gchh.Trend)
	GCHH.WinCount = []int{0, 0, 0, 0}
	for i, _ := range this.EDesk.Gchh.Trend {
		if len(this.EDesk.Gchh.Trend[i]) > 20 {
			this.EDesk.Gchh.Trend[i] = this.EDesk.Gchh.Trend[i][1:]
		}
		for j, _ := range this.EDesk.Gchh.Trend[i] {
			if this.EDesk.Gchh.Trend[i][j] > 0 {
				GCHH.WinCount[i] += 1
			}
		}
	}
	GCHH.Trend = this.EDesk.Gchh.Trend
	// logs.Debug("GCHH:%v", GCHH)
	this.EDesk.Gchh = GCHH

	var bankerResult int64 = 0 // 庄家输赢结果

	for _, p := range this.EDesk.Players {
		var sendresults map[int]GClientSeatBalance = make(map[int]GClientSeatBalance)
		for id, seatInfo := range this.EDesk.Seats {
			var myButtom int64 = 0
			var myResult int64 = 0

			if value, ok := p.DownBets[uint8(id)]; ok {
				myButtom = value
			}
			if value, ok := p.BalaDownBets[uint8(id)]; ok {
				myResult = value
			}

			sendresults[id] = GClientSeatBalance{
				Bottom:   seatInfo.DownBetValue,
				Result:   results[id].Result,
				MyBottom: myButtom,
				MyResult: myResult,
			}

			if !p.Robot {
				//（此原始的庄家结算方法似乎有问题，应改为 bankerResult -= myResult ，有待考证）
				//bankerResult -= (myResult - myButtom) //原始方法
				bankerResult -= myResult //新方法
			}
		}
		// 发送结算给用户
		p.SendNativeMsg(MSG_GAME_BALANCE, &GClientBalance{
			Id:      MSG_GAME_BALANCE,
			Results: sendresults,
			MyCoin:  p.Coins,
		})
	}

	// 修改库存值
	if GetCostType() == 1 {
		AddLocalStock(bankerResult)
		AddCD(bankerResult)
	}

	this.EDesk.RemoveAllOfflineAndExistSeat()

}

func zhuanhuan(card []int) []int {
	res := []int{}
	for _, v := range card {
		res = append(res, v&0xF)
	}
	return res
}

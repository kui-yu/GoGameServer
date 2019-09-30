package main

import (
	"logs"
	"math/rand"
	"time"

	"bl.com/paigow"
	"bl.com/util"
)

// 派奖
func (this *ExtDesk) TimerAward(d interface{}) {
	this.Lock()
	defer this.Unlock()
	sd := GGameAwardNotify{
		Id:    MSG_GAME_INFO_AWARD_NOTIFY,
		Timer: int32(gameConfig.Timer.AwardNum) * 1000,
	}

	// 添加走势
	// 庄
	bType := paigow.GetCardsType(this.BankerCard) + 1
	this.RunChart[0] = append(this.RunChart[0], bType)

	// 闲
	for i := 0; i < 3; i++ {
		chart := paigow.CompareCard(this.BankerCard, this.IdleCard[i])
		if chart {
			this.RunChart[i+1] = append(this.RunChart[i+1], 0)
		} else {
			this.RunChart[i+1] = append(this.RunChart[i+1], 1)
		}
	}

	for i := 0; i < 4; i++ {
		length := len(this.RunChart[i])

		if length > gameConfig.DeskInfo.RunChartCount {
			this.RunChart[i] = this.RunChart[i][1:]
		}
	}
	sd.BankRunChart = this.RunChart[0]
	// 游戏结算，返回区域输赢情况
	_, WinArea, TWinArea := this.GameEnd(this.TypeList)
	sd.WinArea = WinArea[:]
	sd.TWinArea = TWinArea[:]
	this.WinArea = WinArea[:]
	//玩家游戏记录
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
	seatList := this.SeatMgr.GetSeatList()
	for _, seat := range seatList {
		win := seat.(*ExtPlayer).GetWinCoins()
		sd.SeatWinCoins = append(sd.SeatWinCoins, win)
		winArea := seat.(*ExtPlayer).GetWinList()
		TWinArea = util.LessInt64List(TWinArea, winArea)
	}

	sd.OtherWinArea = TWinArea

	if GetCostType() == 1 { //如果不是体验场
		//发送结算消息给数据库, 简单记录
		dbreq := GGameEnd{
			Id:          MSG_GAME_END_NOTIFY,
			GameId:      GCONFIG.GameType,
			GradeId:     GCONFIG.GradeType,
			RoomId:      GCONFIG.RoomType,
			GameRoundNo: this.GameId,
			Mini:        false,
			SetLeave:    1, //是否设置离开，0离开，1不离开
		}

		//发送消息给大厅去记录游戏记录
		rdreq := GGameRecord{
			Id:          MSG_GAME_END_RECORD,
			GameId:      GCONFIG.GameType,
			GradeId:     GCONFIG.GradeType,
			RoomId:      GCONFIG.RoomType,
			GameRoundNo: this.GameId,
			BankerCard:  this.BankerCard,
			IdleCard:    this.IdleCard,
		}

		players := this.SeatMgr.GetUserList(len(this.Players))
		for _, v := range players {
			re := RecordData{
				MatchNum:   this.JuHao,
				RoomName:   roomName,
				EndTime:    endTime,
				Date:       time.Now().Unix(),
				TypeList:   this.TypeList,
				IdleCard:   this.IdleCard,
				BankerCard: this.BankerCard,
			}

			//获取开奖区域
			for i, v := range this.WinArea {
				if v {
					re.WinArea = append(re.WinArea, i)
				}
			}
			p := v.(*ExtPlayer)
			sd.PWin = p.GetWinCoins()
			sd.PWinArea = p.GetWinList()

			areaList := p.GetTotBetList()
			valid := this.BuildValid(WinArea, areaList)

			betCoins := p.GetTotAreaCoins()

			// 结算
			sd.PrizeCoins = sd.PWin - betCoins
			p.Award()
			sd.PCoins = p.GetCoins()

			if betCoins > 0 {
				dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
					UserId:      p.GetUid(),
					UserAccount: p.Account,
					BetCoins:    betCoins,
					ValidBet:    valid,
					PrizeCoins:  sd.PWin - betCoins,
					Robot:       p.Robot,
				})
				if !p.Robot {
					p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
				}
				dbreq.UserCoin = []GGameEndInfo{}

				if p.Robot {
					continue
				}
				rddata := GGameRecordInfo{
					UserId:      p.GetUid(),
					UserAccount: p.Account,
					BetCoins:    betCoins,
					BetArea:     p.GetTotBetList(),
					PrizeCoins:  sd.PWin - betCoins,
					CoinsAfter:  sd.PCoins,
					Robot:       p.Robot,
				}
				rddata.CoinsBefore = rddata.CoinsAfter - rddata.PrizeCoins
				rdreq.UserRecord = append(rdreq.UserRecord, rddata)
				p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
				rdreq.UserRecord = []GGameRecordInfo{}
			}
			for i, v := range p.GetNTAreaCoinsList() {
				re.AllBet += v
				if v > 0 {
					re.BetArea = append(re.BetArea, AreaAndCoins{
						AreaIndx: i,
						Coins:    v,
					})
				}
			}
			re.WinOrLost = sd.PrizeCoins
			if re.AllBet > 0 {
				//logs.Debug("该玩家本局有下注，所以保存该局游戏记录")
				G_AllRecord.AddRecord(p.Uid, &re)
			} else {
				//logs.Debug("改玩家没有下注，所以不保存该局游戏记录")
			}
			if GetCostType() == 1 && !p.Robot {
				AddCD(-sd.PrizeCoins)
				AddLocalStock(-sd.PrizeCoins)
				//fmt.Println("添加库存:", -sd.PrizeCoins)
			}
			p.SendNativeMsg(MSG_GAME_INFO_AWARD_NOTIFY, sd)
		}
	} else {
		//如果是体验场，只发送派奖信息
		player := this.SeatMgr.GetUserList(len(this.Players))
		for _, v := range player {
			re := RecordData{
				MatchNum:   this.JuHao,
				RoomName:   roomName,
				EndTime:    endTime,
				Date:       time.Now().Unix(),
				TypeList:   this.TypeList,
				IdleCard:   this.IdleCard,
				BankerCard: this.BankerCard,
			}

			//获取开奖区域
			for i, v := range this.WinArea {
				if v {
					re.WinArea = append(re.WinArea, i)
				}
			}
			p := v.(*ExtPlayer)
			sd.PWin = p.GetWinCoins()
			sd.PWinArea = p.GetWinList()
			betCoins := p.GetTotAreaCoins()
			// 结算
			sd.PrizeCoins = sd.PWin - betCoins
			p.Award()
			sd.PCoins = p.GetCoins()
			for i, v := range p.GetNTAreaCoinsList() {
				re.AllBet += v
				if v > 0 {
					re.BetArea = append(re.BetArea, AreaAndCoins{
						AreaIndx: i,
						Coins:    v,
					})
				}
			}
			if re.AllBet > 0 {
				logs.Debug("该玩家本局有下注，所以保存该局游戏记录")
				G_AllRecord.AddRecord(p.Uid, &re)
			} else {
				logs.Debug("改玩家没有下注，所以不保存该局游戏记录")
			}
			p.SendNativeMsg(MSG_GAME_INFO_AWARD_NOTIFY, sd)
		}
	}
	this.GameState = MSG_GAME_INFO_AWARD_NOTIFY

	this.ResetAreaCoins()
	this.TimerOver(nil)
}

func (this *ExtDesk) BuildValid(winArea []bool, areaList []int64) int64 {
	// 同时押输、赢，区域不计算有效打码
	if areaList[INDEX_TIAN_WIN-1] > 0 && areaList[INDEX_TIAN_LOSS-1] > 0 {
		areaList[INDEX_TIAN_WIN-1] = 0
		areaList[INDEX_TIAN_LOSS-1] = 0
	}
	if areaList[INDEX_DI_WIN-1] > 0 && areaList[INDEX_DI_LOSS-1] > 0 {
		areaList[INDEX_DI_WIN-1] = 0
		areaList[INDEX_DI_LOSS-1] = 0
	}
	if areaList[INDEX_REN_WIN-1] > 0 && areaList[INDEX_REN_LOSS-1] > 0 {
		areaList[INDEX_REN_WIN-1] = 0
		areaList[INDEX_REN_LOSS-1] = 0
	}

	var ret int64
	for _, v := range areaList {
		ret += v
	}

	return ret
}
func (this *ExtDesk) testResult() int64 {
	placeResult := make([]int, 4) //区域输赢结果
	//庄牌排序
	bankerCard := this.BankerCard
	bankerCard = paigow.Sort(this.BankerCard)
	//闲牌排序
	idleCard := this.IdleCard
	for i := 0; i < 3; i++ {
		idleCard[i] = paigow.Sort(idleCard[i])
	}
	placeResult[len(placeResult)-1] = int(paigow.GetCardsType(bankerCard) + 1) //庄
	// 闲
	for i := 0; i < 3; i++ {
		chart := paigow.CompareCard(bankerCard, idleCard[i])
		if chart {
			placeResult[i] = -1
		} else {
			placeResult[i] = 1
		}
	}
	var winCoins int64
	for _, v := range this.Players {
		if v.Robot {
			continue
		}
		for i := 0; i <= 7; i++ {
			switch i {
			case 0, 1:
				if v.GetTotAreaCoin(i) > 0 {
					if i == 0 { //下注赢
						if placeResult[0] == 1 { //赢了
							winCoins += int64(float64(v.GetTotAreaCoin(i)) * 0.99)
						} else { //输了
							winCoins -= v.GetTotAreaCoin(i)
						}
					} else { //下注输
						if placeResult[0] == -1 { //输了
							winCoins += int64(float64(v.GetTotAreaCoin(i)) * 0.93)
						} else { //输了
							winCoins -= v.GetTotAreaCoin(i)
						}
					}
				}
			case 2, 3:
				if v.GetTotAreaCoin(i) > 0 {
					if i == 2 { //下注赢
						if placeResult[1] == 1 { //赢了
							winCoins += int64(float64(v.GetTotAreaCoin(i)) * 0.99)
						} else { //输了
							winCoins -= v.GetTotAreaCoin(i)
						}
					} else { //下注输
						if placeResult[1] == -1 { //输了
							winCoins += int64(float64(v.GetTotAreaCoin(i)) * 0.93)
						} else { //输了
							winCoins -= v.GetTotAreaCoin(i)
						}
					}
				}
			case 4, 5:
				if v.GetTotAreaCoin(i) > 0 {
					if i == 4 { //下注赢
						if placeResult[2] == 1 { //赢了
							winCoins += int64(float64(v.GetTotAreaCoin(i)) * 0.99)
						} else { //输了
							winCoins -= v.GetTotAreaCoin(i)
						}
					} else { //下注输
						if placeResult[2] == -1 { //输了
							winCoins += int64(float64(v.GetTotAreaCoin(i)) * 0.93)
						} else { //输了
							winCoins -= v.GetTotAreaCoin(i)
						}
					}
				}
			case 6: //天拖
				if v.GetTotAreaCoin(i) > 0 {
					if placeResult[3] == 2 {
						winCoins += v.GetTotAreaCoin(i) * 79
					} else {
						winCoins -= v.GetTotAreaCoin(i)
					}
				}
			case 7:
				if v.GetTotAreaCoin(i) > 0 {
					if placeResult[3] == 1 {
						winCoins += v.GetTotAreaCoin(i) * 179
					} else {
						winCoins -= v.GetTotAreaCoin(i)
					}
				}
			}
		}
	}
	//fmt.Println("输赢金币:", winCoins)
	return winCoins
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

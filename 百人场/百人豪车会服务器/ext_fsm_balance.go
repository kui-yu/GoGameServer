package main

import (
	//"fmt"
	"logs"
	"time"
)

type FSMSettle struct {
	Mark int

	EDesk       *ExtDesk
	EndDataTime int64 // 当前状态的结束时间
}

func (this *FSMSettle) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMSettle) Run() {
	// DebugLog("游戏状态-结算")
	broCount = 1
	this.EndDataTime = GetTimeMS() + int64(gameConfig.StateInfo.BalanceTime)

	this.addListen() // 添加监听
	this.EDesk.GameState = GAME_STATUS_BALANCE
	this.EDesk.SendGameState(GAME_STATUS_BALANCE, int64(gameConfig.StateInfo.BalanceTime))

	this.EDesk.AddTimer(GAME_STATUS_BALANCE, gameConfig.StateInfo.BalanceTime/1000, this.TimerCall, nil)

	this.sendRecord()
	this.balanceAll()

}

func (this *FSMSettle) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_READY)
}

func (this *FSMSettle) GetMark() int {
	return this.Mark
}
func (this *FSMSettle) Leave() {
	this.removeListen()
}

func (this *FSMSettle) getRestTime() int64 {
	remainTimeMS := this.EndDataTime - GetTimeMS()
	return remainTimeMS
}

func (this *FSMSettle) addListen() {}

func (this *FSMSettle) removeListen() {}

func (this *FSMSettle) balanceAll() {

	//添加风控的当前库存
	var allcoin int64
	allAreacoin := make([]int64, 8, 8)
	for i, b := range this.EDesk.Bets {
		allcoin += b.UserBetValue
		allAreacoin[i] += b.UserBetValue
	}
	wincoins := allcoin - int64(float32(allAreacoin[this.EDesk.GameResult])*CarTypeMultiple[this.EDesk.GameResult])
	AddCD(wincoins)
	logs.Debug("开奖区域：", this.EDesk.GameResult, "赢的金额为：", wincoins)
	//发送结算消息给数据库, 简单记录
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.EDesk.JuHao,
		Mini:        false,
		SetLeave:    1,
	}
	//fmt.Println("GGameEnd:", dbreq)
	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.EDesk.JuHao,
		GradeNumber: "001",
		LotteryCard: this.EDesk.GameResult,
	}
	var BankResult int64 = 0
	for id, v := range this.EDesk.Bets {
		BankResult += v.UserBetValue
		if id == this.EDesk.GameResult {
			BankResult -= (v.UserBetValue * int64((CarTypeMultiple[id])*10)) / 10
		}
	}
	for _, v := range this.EDesk.Players {
		var betCoins int64 = 0
		var waterRate float64 = 0
		for id, value := range v.DownBets {
			betCoins += value
			if id != this.EDesk.GameResult {
				v.BalaDownBets[id] = -value
			} else {
				v.BalaDownBets[id] = (value * int64((CarTypeMultiple[id])*10)) / 10
			}
		}

		coin, ok := v.BalaDownBets[this.EDesk.GameResult]
		if ok {
			v.WinCoins = coin - betCoins
		} else {
			v.WinCoins = -betCoins
		}

		if v.BalaDownBets[this.EDesk.GameResult] > betCoins {
			waterRate = float64(v.BalaDownBets[this.EDesk.GameResult]-betCoins) * this.EDesk.Rate
		}
		if !v.Robot && betCoins != 0 {
			v.UnbetsCount = 0
			uc := GGameEndInfo{
				UserId:      v.Uid,
				UserAccount: v.Account,
				BetCoins:    betCoins,
				ValidBet:    betCoins,
				PrizeCoins:  v.WinCoins,
				Robot:       v.Robot,
				WaterProfit: 0,
				WaterRate:   waterRate,
			}

			dbreq.UserCoin = append(dbreq.UserCoin, uc)
			if GetCostType() == 1 {
				v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
			}
			dbreq.UserCoin = []GGameEndInfo{}
		} else if betCoins == 0 {
			v.UnbetsCount++
		}
		var rdDownBets [8]int32
		for id, value := range v.DownBets {
			rdDownBets[id] = int32(value)
		}

		if !v.Robot && betCoins != 0 {
			rddata := GGameRecordInfo{
				UserId:      v.Uid,
				UserAccount: v.Account,
				BetCoins:    betCoins,             // 下注的金币
				BetArea:     rdDownBets,           // 区域下注情况
				PrizeCoins:  v.WinCoins,           // 赢取的金币
				CoinsAfter:  v.Coins,              // 结束后金币
				CoinsBefore: v.Coins + v.WinCoins, // 下注前金币
				Robot:       v.Robot,
			}
			rdreq.UserRecord = append(rdreq.UserRecord, rddata)
			if GetCostType() == 1 {
				// 发送记录到存储
				v.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
			}
			rdreq.UserRecord = []GGameRecordInfo{}
		}

	}

	//修改库存值
	if GetCostType() == 1 {
		AddLocalStock(BankResult)
	}
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

	ss := this.EDesk.GameResult
	for _, v := range this.EDesk.Players {
		//玩家游戏记录
		re := RecordData{
			MatchNum: this.EDesk.JuHao,
			RoomName: roomName,
			EndTime:  endTime,
			Date:     time.Now().Unix(),
		}
		re.Car = numBecomString([]int{ss})[0]
		var ab int64
		for _, v1 := range v.DownBets {
			ab += v1
		}
		info := GNBalance{
			Id:      MSG_GAME_INFO_NBALANCE,
			Head:    v.Head,
			BetAll:  ab,
			Nick:    v.Nick,
			CarName: numBecomString([]int{this.EDesk.GameResult})[0],
		}
		info.Results = make(map[int]GBetBalance)
		// fmt.Println("baladownbet:", v.BalaDownBets)
		// fmt.Println("Deskbet:", this.EDesk.Bets)
		var yin int64
		for id, value := range v.BalaDownBets {
			msg := GBetBalance{
				Bottom:   this.EDesk.Bets[id].DownBetValue,
				MyResult: value,
				MyBottom: v.DownBets[id],
			}
			// info.WinOrLoseCoins += vaue

			if value > 0 {
				yin += value
			}
			// for _, v := range this.EDesk.Players {
			// fmt.Println(v.Nick, "的结算集合：", v.BalaDownBets)
			// }
			if id == this.EDesk.GameResult {
				// fmt.Println("当前位置：", id, "将位置 标记为中奖区域")
				msg.Win = true
				for _, p1 := range this.EDesk.Players {
					for i, value1 := range p1.BalaDownBets {
						if v.Uid != p1.Uid {
							v.ElseWinAndOr[i] += value1
						}
					}
				}
				// for i, p1 := range this.EDesk.getChairPlayer(v) {
				// 	for _, value1 := range p1.BalaDownBets {
				// 		v.ChairWinOrLost[i] += value1
				// 	}
				// }
			}
			info.ElseWinAndLose = v.ElseWinAndOr
			// info.SeatWinCoins = v.ChairWinOrLost
			if value > 0 {
				msg.Result = 1
				// v.BetHistorys = append(v.BetHistorys, 0)
			} else {
				// v.BetHistorys = append(v.BetHistorys, 1)
			}
			// fmt.Println("msg", msg)
			info.Results[id] = msg
		}
		// fmt.Println("没加上赢钱：", v.Coins)
		// fmt.Println("赢的钱:", yin)
		// fmt.Println("需要扳回一局的钱：", jia)
		// info.MyCoin = v.Coins + jia + yin
		info.MyCoin = v.Coins + yin
		// fmt.Println("金币比比iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii:", info.MyCoin)
		v.Coins = info.MyCoin
		info.CanUserChip = this.EDesk.CanUseChip(v)
		info.WinOrLoseCoins = v.WinCoins
		v.SendNativeMsg(MSG_GAME_INFO_NBALANCE, &info)
		//logs.Debug("结算的玩家输赢~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~：", info.WinOrLoseCoins)
		// fmt.Println("结算！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！--------------------------------：", info.WinOrLoseCoins)
		var pdd bool = false
		for _, v11 := range v.PAreaCoins {
			if v11 != 0 {
				pdd = true
			}
		}
		if pdd {
			if info.WinOrLoseCoins > 0 {
				// fmt.Println("对了！！！！！！！！@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
				v.BetHistorys = append(v.BetHistorys, 0)
			} else {
				v.BetHistorys = append(v.BetHistorys, 1)
			}

		}
		// fmt.Println("结算完成之后", v.Nick, "的金币为:", v.Coins)
		re.WinOrLost = v.WinCoins
		//赋值游戏记录的值
		var allB int64
		for i, v1 := range v.PAreaCoins {
			allB += v1
			if v1 > 0 {
				re.BetArea = append(re.BetArea, AreaAndBet{
					BetCoins: v1,
					BetArea:  i,
				})
			}
		}
		re.AllBet = allB
		if allB > 0 {
			G_AllRecord.AddRecord(v.Uid, &re)
		}
		// //清除其他玩家下注
		// for _, v := range this.EDesk.Players {
		// 	v.OtherBet = []int64{0, 0, 0, 0, 0, 0, 0, 0}
		// }
		var pd22 bool = false
		for _, a1 := range v.PAreaCoins {
			if a1 > 0 {
				pd22 = true
			}
		}
		if pd22 {
			// fmt.Println("有存玩家!!!!!!!!!!!!!!!!!!!!!!!!!!!：", v.Nick, "的数据")
			this.EDesk.BetAgain[v.Uid] = v.PAreaCoins
		}
		// //玩家下注局数
		// for _, v1 := range v.PAreaCoins {
		// 	if v1 > 0 {
		// 		v.Match++
		// 		continue
		// 	}
		// }
	}
	oldOther = make(map[int64][]int64)
	// for _, v := range this.EDesk.Players {
	// fmt.Println("金币公布!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!:")
	// fmt.Println("玩家：", v.Nick, "的金币是：", v.Coins)
	// }
	this.EDesk.RemoveAllOfflineAndExistSeat()
}

func (this *FSMSettle) sendRecord() {
	if len(this.EDesk.LogoRecord) < gameConfig.LimitInfo.LogoLimit {
		this.EDesk.LogoRecord = append([]int{this.EDesk.GameResult}, this.EDesk.LogoRecord...)
	} else {
		this.EDesk.LogoRecord = append([]int{this.EDesk.GameResult}, this.EDesk.LogoRecord...)
		this.EDesk.LogoRecord = this.EDesk.LogoRecord[:gameConfig.LimitInfo.LogoLimit]
	}

	info := GNRecord{
		Id:              MSG_GAME_INFO_NRECORD,
		Record:          this.EDesk.LogoRecord,
		OnlinePlayerNum: len(this.EDesk.Players),
	}
	DebugLog("开奖记录通知：", this.EDesk.LogoRecord)
	this.EDesk.BroadcastAll(MSG_GAME_INFO_NRECORD, &info)
}

func numBecomString(cars []int) []string {
	var arr []string
	for _, v := range cars {
		switch v {
		case 0:
			arr = append(arr, "大法拉利")
			break
		case 1:
			arr = append(arr, "大玛莎拉蒂")
			break
		case 2:
			arr = append(arr, "大保时捷")
			break
		case 3:
			arr = append(arr, "大奔驰")
			break
		case 4:
			arr = append(arr, "小法拉利")
			break
		case 5:
			arr = append(arr, "小玛莎拉蒂")
			break
		case 6:
			arr = append(arr, "小保时捷")
			break
		case 7:
			arr = append(arr, "小奔驰")
			break
		default:
			arr = append(arr, "无")
			break
		}
	}
	return arr
}

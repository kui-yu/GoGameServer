package main

import (
	// "encoding/json"
	"logs"
	"time"
)

// func (this *ExtDesk) status_ready(d interface{}) {
// 	logs.Debug("准备阶段")
// 	this.GameState = this.StatusAndTimes.ReadId
// 	this.BroStatusTime(this.StatusAndTimes.ReadMs)
// 	//准备阶段多用于客户端与服务端清除数据,回收手牌,改变局号等等
// 	this.JuHao = GetJuHao()
// 	this.BroadcastAll(MSG_GAME_INFO_JUHAOCHANGE, JuHaoChange{
// 		Id:    MSG_GAME_INFO_JUHAOCHANGE,
// 		JuHao: this.JuHao,
// 	})
// 	this.AddTimer(this.StatusAndTimes.ReadId, this.StatusAndTimes.ReadMs/1000, this.status_shufflecard, nil)
// }

//洗牌阶段
func (this *ExtDesk) status_shuffle(d interface{}) {
	logs.Debug("洗牌阶段")
	this.GameState = this.StatusAndTimes.ShuffleId
	this.BroStatusTime(this.StatusAndTimes.ShuffleMs)
	//下注
	this.AddTimer(this.StatusAndTimes.ShuffleId, this.StatusAndTimes.ShuffleMs, this.status_startdownbet, nil)
}

//开始下注阶段
func (this *ExtDesk) status_startdownbet(d interface{}) {
	logs.Debug("开始下注阶段")
	this.GameState = this.StatusAndTimes.StartdownbetsId
	this.BroStatusTime(this.StatusAndTimes.StartdownbetsMs)
	//清除数据
	for _, v := range this.Players {
		v.Rest()
	}
	this.Rest()
	//开始下注阶段用于客户端与服务端清除数据,回收手牌,改变局号等等
	this.JuHao = GetJuHao()
	//发送句号改变通知
	this.BroadcastAll(MSG_GAME_INFO_JUHAOCHANGE, &struct {
		Id    int
		JuHao string
	}{
		Id:    MSG_GAME_INFO_JUHAOCHANGE,
		JuHao: this.JuHao,
	})
	//洗牌
	this.ShuffleCard()
	this.AddTimer(999, 1, this.BroDownBetInfo, nil) //开始每隔一秒更新一次下注情况
	this.AddTimer(this.StatusAndTimes.StartdownbetsId, this.StatusAndTimes.StartdownbetsMs, this.status_downbet, nil)
}

//下注阶段
func (this *ExtDesk) status_downbet(d interface{}) {
	logs.Debug("下注阶段")
	this.GameState = this.StatusAndTimes.DownBetsId
	this.BroStatusTime(this.StatusAndTimes.DownBetsMS)
	//下注
	this.AddTimer(this.StatusAndTimes.DownBetsId, this.StatusAndTimes.DownBetsMS, this.status_stopdownbet, nil)
}

//停止下注阶段
func (this *ExtDesk) status_stopdownbet(d interface{}) {
	logs.Debug("下注结束")
	this.GameState = this.StatusAndTimes.StopdownbetsId
	this.BroStatusTime(this.StatusAndTimes.StopdownbetsMs)
	//下注
	this.AddTimer(this.StatusAndTimes.StopdownbetsId, this.StatusAndTimes.StopdownbetsMs, this.status_facard, nil)
}

//发牌阶段
func (this *ExtDesk) status_facard(d interface{}) {
	logs.Debug("发牌阶段")
	this.GameState = this.StatusAndTimes.FaCardId
	this.BroStatusTime(this.StatusAndTimes.FaCardMS)
	//发牌

	this.AddTimer(this.StatusAndTimes.FaCardId, this.StatusAndTimes.FaCardMS, this.status_opencard, nil)
}

//开牌阶段
func (this *ExtDesk) status_opencard(d interface{}) {
	logs.Debug("开牌阶段")
	this.GameState = this.StatusAndTimes.OpenCardId
	this.BroStatusTime(this.StatusAndTimes.OpenCardMS)
	//开牌
	this.allotCard()
	// for i, v := range this.CardGroupArray {
	// 	logs.Debug("", i, "区域手牌情况：", v)
	// }
	this.BroadcastAll(MSG_GAME_INFO_FACARD_BRO, &struct {
		Id    int
		Cards map[int]CardGroupInfo
	}{
		Id:    MSG_GAME_INFO_FACARD_BRO,
		Cards: this.CardGroupArray,
	})
	this.AddTimer(this.StatusAndTimes.OpenCardId, this.StatusAndTimes.OpenCardMS, this.status_balance, nil)
}

//结算阶段
func (this *ExtDesk) status_balance(d interface{}) {
	logs.Debug("结算阶段")
	this.GameState = this.StatusAndTimes.BalanceId
	this.BroStatusTime(this.StatusAndTimes.BalanceMS)
	//结算阶段 用于用户金币结算，用户游戏记录储存
	//	发送结算消息给数据库, 简单记录
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		Mini:        false,
		SetLeave:    1,
	}
	BankerCard := this.CardGroupArray[gameConfig.LimitInfo.BetAreaCount].Cards
	IdleCard := [][]int{}
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		IdleCard = append(IdleCard, this.CardGroupArray[i].Cards)
	}
	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.JuHao,
		BankerCard:  BankerCard,
		IdleCard:    IdleCard,
	}
	cga := map[int]CardGroupInfo{}
	for i, v := range this.CardGroupArray {
		cga[i] = v
	}
	areaResult, multipleResult := GetResult(cga)
	logs.Debug("执行神秘方法之后的 >>>GR:", this.CardGroupArray)
	this.AreaWinDouble = multipleResult
	//添加游戏走势
	this.GameZs.Zoushi[0] = append(this.GameZs.Zoushi[0], Pzshi{
		Mark:   -1,
		ZouShi: typetostring(int(this.CardGroupArray[gameConfig.LimitInfo.BetAreaCount].CardGroupType)),
	})
	//闲家
	for i, v := range areaResult {
		var str string
		if v == 0 {
			str = typetostring(int(this.CardGroupArray[gameConfig.LimitInfo.BetAreaCount].CardGroupType))
		} else {
			str = typetostring(int(this.CardGroupArray[i].CardGroupType))
		}
		this.GameZs.Zoushi[i+1] = append(this.GameZs.Zoushi[i+1], Pzshi{
			Mark:   v,
			ZouShi: str,
		})
	}
	//庄家
	//限制走势长度
	if len(this.GameZs.Zoushi[1]) > gameConfig.LimitInfo.Qzhoushicount {
		for i, _ := range this.GameZs.Zoushi {
			if i != 0 {
				this.GameZs.Zoushi[i] = this.GameZs.Zoushi[i][6:]
			}
		}
	}
	if len(this.GameZs.Zoushi[0]) > gameConfig.LimitInfo.Zzhoushicount {
		this.GameZs.Zoushi[0] = this.GameZs.Zoushi[0][1:]
	}

	//更新走势
	// zs, err := json.Marshal(this.GameZs)
	// if err != nil {
	// 	logs.Error("更新走势时json转换出错")
	// }
	// this.DeskMgr.SetZouShi(&GameZouShi{
	// 	SerId: int32(G_DbGetGameServerData.Sid),
	// 	GameInfo: GameTypeDetail{GameType: int32(GCONFIG.GameType),
	// 		RoomType:  int32(GCONFIG.RoomType),
	// 		GradeType: int32(GCONFIG.GradeType),
	// 	},
	// 	ZouShi:    string(zs),
	// 	PlayerNum: int32(len(this.Players)),
	// 	UpdateT:   time.Now().Unix(),
	// })

	//结算
	var bankerResult int64 = 0 // 庄家输赢结果
	var b int64 = 0            //区域总输赢结果
	areare := []int64{}        //每个区域总输赢结果
	for i, v := range areaResult {
		if v == 1 {
			this.WinArea[i] = 1
		} else {
			this.WinArea[i] = -1
		}
		if v == 0 {
			//庄家赢
			bankerResult += this.DownBetZhenshi[i] * int64(multipleResult[i])
			b += this.DownBet[i] * int64(multipleResult[i])
			areare = append(areare, b)
		} else {
			//庄家输
			bankerResult -= this.DownBetZhenshi[i] * int64(multipleResult[i])
			b -= this.DownBet[i] * int64(multipleResult[i])
			areare = append(areare, b)
		}
	}
	logs.Debug("区域总输赢结果:::", b)
	for _, v := range this.Players {
		var allBet int64 //玩家总投注
		for _, v1 := range v.DownBet {
			allBet += v1
		}
		for i, v1 := range v.DownBet {
			if areaResult[i] == 0 {
				//玩家输
				v.BalanceResult[i] -= v1 * int64(multipleResult[i]-1)
				v.BalanceResultTosee[i] -= v1 * int64(multipleResult[i])
			} else {
				//玩家赢
				v.BalanceResult[i] += v1 * int64(multipleResult[i]+1)
				v.BalanceResultTosee[i] += v1 * int64(multipleResult[i])
			}
			this.BalanceResult[i] += v.BalanceResult[i]
		}
		var getCoins int64
		// var aftergetCoins int64
		for _, v1 := range v.BalanceResult {
			// if v1 > 0 {
			// 	aftergetCoins += v1
			// } else {
			// 	aftergetCoins += v1 + v.DownBet[i]
			// }
			// if v1 > 0 {
			// 	getCoins += v1 - v.DownBet[i]
			// } else {
			getCoins += v1
			// }
		}
		var getCoins2 int64
		for _, v1 := range v.BalanceResultTosee {
			getCoins2 += v1
		}
		v.Coins += getCoins
		//根据实际获取金币判断输赢
		//判断本局是否下注
		if allBet > 0 {
			//玩家本局有下注
			beth := BetHistory{}
			if getCoins > 0 {
				beth.IsVictory = true
			}
			beth.DownBet = allBet
			v.History = append(v.History, beth)
			if len(v.History) >= gameConfig.LimitInfo.Userlist_record_count {
				v.History = v.History[1:]
			}
		} else {
			v.NoToBet += 1
		}

		if allBet > 0 {

			dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
				UserId:      v.Uid,
				UserAccount: v.Account,
				BetCoins:    allBet,
				ValidBet:    allBet,
				PrizeCoins:  getCoins2,
				Robot:       v.Robot,
			})

			if GetCostType() == 1 && !v.Robot {
				v.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)

			}
			dbreq.UserCoin = []GGameEndInfo{}
		}
		if !v.Robot && allBet > 0 {
			rddata := GGameRecordInfo{
				UserId:      v.Uid,
				UserAccount: v.Account,
				BetCoins:    allBet,             // 下注的金币
				BetArea:     v.DownBet,          // 区域下注情况
				PrizeCoins:  getCoins2,          // 赢取的金币
				CoinsAfter:  v.Coins,            // 结束后金币
				CoinsBefore: v.Coins - getCoins, // 下注前金币
				Robot:       v.Robot,
			}
			rdreq.UserRecord = append(rdreq.UserRecord, rddata)
			if GetCostType() == 1 {
				v.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
			}
			rdreq.UserRecord = []GGameRecordInfo{}
		}
	}
	//修改库存值
	if GetCostType() != 2 {
		logs.Debug("庄家输赢：：：：：：：：：：：：：", bankerResult)
		AddCD(bankerResult)
		AddLocalStock(bankerResult)
	}
	//添加游戏记录
	roomName := this.GetRoomName()
	endTime := time.Now().Format("2006-01-02 15:04:05")
	for _, v := range this.Players {
		//玩家游戏记录
		re := RecordData{
			MatchNum:  this.JuHao,
			RoomName:  roomName,
			EndTime:   endTime,
			Date:      time.Now().Unix(),
			ZCardType: int(this.CardGroupArray[gameConfig.LimitInfo.BetAreaCount].CardGroupType),
		}
		for i, v1 := range v.DownBet {
			re.AllBet += v1
			if v1 > 0 {
				re.BetArea = append(re.BetArea, BetArea{
					AreaIndex: i,
					BetCoins:  v1,
					CardType:  int(this.CardGroupArray[i].CardGroupType),
				})
			}
		}
		for _, v1 := range v.BalanceResultTosee {
			// if v1 > 0 {
			// re.WinOrLost += v1 - v.DownBet[i]
			// } else {
			re.WinOrLost += v1
			// }
		}
		if re.AllBet > 0 {
			logs.Debug("玩家有游戏，将其记录存入")
			G_AllRecord.AddRecord(v.Uid, &re)
		} else {
			logs.Debug("玩家没在游戏，不存记录")
		}

	}
	//向玩家发送结算
	for _, v := range this.Players {
		var mygetcoins int64 //玩家自身盈利
		for _, v1 := range v.BalanceResultTosee {
			mygetcoins += v1
		}
		v.DownBet = []int64{}
		//初始化玩家下注区域集合
		for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
			v.DownBet = append(v.DownBet, 0)
		}
		v.SendNativeMsg(MSG_GAME_INFO_BALANCE, ToClientBalance{
			Id:            MSG_GAME_INFO_BALANCE,
			MyCoins:       v.Coins,
			MyResult:      v.BalanceResultTosee,
			AllResult:     this.BalanceResult,
			MyGetCoins:    mygetcoins,
			Zoushi:        this.GameZs.Zoushi,
			BetAbleIndex:  this.CanUseChip(v),
			AreaWinDouble: this.AreaWinDouble,
			WinArea:       this.WinArea,
		})
	}
	logs.Debug("当前库存:::::::::::::::::::::::::::::::::", CD)
	logs.Debug("结算完成之后的 >>>GR:", this.CardGroupArray)
	this.removePlayerByneed()
	this.AddTimer(this.StatusAndTimes.BalanceId, this.StatusAndTimes.BalanceMS, this.status_shuffle, nil)
}
func typetostring(cardtype int) string {
	var result string
	switch cardtype {
	case 1:
		result = "一"
		break
	case 2:
		result = "二"
		break
	case 3:
		result = "三"
		break
	case 4:
		result = "四"
		break
	case 5:
		result = "五"
		break
	case 6:
		result = "六"
		break
	case 7:
		result = "七"
		break
	case 8:
		result = "八"
		break
	case 9:
		result = "九"
		break
	case 10:
		result = "牛"
		break
	case 11:
		result = "炸"
		break
	case 12:
		result = "花"
		break
	default:
		result = "无"
		break
	}
	return result
}
func (this *ExtDesk) removePlayerByneed() {
	nobetwarning := gameConfig.LimitInfo.NobetWarning
	nobetremove := gameConfig.LimitInfo.NobetRemove
	for _, v := range this.Players {
		//未下注提示或退出
		if v.NoToBet >= nobetwarning {
			if v.NoToBet >= nobetremove {
				v.SendNativeMsg(MSG_GAME_INFO_WARNING, MsgWarning{
					Id:     MSG_GAME_INFO_WARNING,
					Result: 2,
				})
				v.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
					Id:      MSG_GAME_LEAVE_REPLY,
					Result:  0,
					Cid:     v.ChairId,
					Uid:     v.Uid,
					Token:   v.Token,
					Robot:   v.Robot,
					NoToCli: true,
				})
				this.UpChair(v)
				this.BroChairChange()
				this.DelPlayer(v.Uid)
				this.DeskMgr.LeaveDo(v.Uid)
			} else if v.NoToBet == gameConfig.LimitInfo.NobetWarning {
				v.SendNativeMsg(MSG_GAME_INFO_WARNING, MsgWarning{
					Id:     MSG_GAME_INFO_WARNING,
					Result: 1,
				})
			}
		}
		if v.LiXian {
			this.UpChair(v)
			this.BroChairChange()
			this.LeaveByForce(v)
		}
	}
	for _, v := range this.Players {
		//查找玩家
		if v.LiXian {
			if v.LiXian {
				this.UpChair(v)
				this.BroChairChange()
				this.LeaveByForce(v)
			}
		}
	}
}

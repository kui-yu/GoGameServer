/**
* 等待开始游戏
**/

package main

import (
	"encoding/json"
	"fmt"
)

type Results struct {
	Uid           int64
	Sid           int
	CarryCoin     int64
	Coin          int64
	Value         int64
	WaterProfit   float64
	Cards         []int
	HandCards     []int
	CardGroupType int
}

type UserCardGroupInfo struct {
	Uid      int64
	CardInfo GCardGroupInfo
}

type GameResult struct {
	JackpotVal int64
	Results    []Results
}

type FsmResults struct {
	Mark        int
	EDesk       *ExtDesk
	EndDateTime int64
}

func (this *FsmResults) InitFsm(mark int, extDesk *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDesk
}

func (this *FsmResults) GetMark() int {
	return this.Mark
}

func (this *FsmResults) Run(upMark int, args ...interface{}) {
	DebugLog("状态 结算")
	timerMs := gameConfig.GameTimer.GameOver

	// 根据人数+2秒
	userWaitTime := 2000
	for _, p := range this.EDesk.Players {
		if p.State == UserStateGameIn {
			userWaitTime += 2000
			DebugLog("人数：+1")
		} else {
			DebugLog("玩家其他状态", p.State)
		}
	}

	DebugLog("倒计时 ", userWaitTime)

	if userWaitTime > timerMs {
		timerMs = userWaitTime
	}

	this.EndDateTime = GetTimeMS() + int64(timerMs)

	this.EDesk.SendDeskStatus(this.Mark, timerMs)
	this.EDesk.AddUniueTimer(TimerId, int(timerMs/1000), this.TimerCall, nil)

	userCardInfos := []UserCardGroupInfo{}
	for _, p := range this.EDesk.Players {
		if p.State != UserStateGameIn {
			continue
		}
		if !p.IsFold {
			cards := []int{}
			cards = append(cards, p.Cards...)
			cards = append(cards, this.EDesk.PublicCards...)

			userCardInfo := UserCardGroupInfo{}
			userCardInfo.Uid = p.Uid

			cstr := ""
			for _, c := range cards {
				cstr = fmt.Sprintf("%s[%d,%d]", cstr, c>>4, c&0xF)
			}
			DebugLog("计算牌型 原数组", cstr, cards)
			userCardInfo.CardInfo = CalcCard(cards)
			DebugLog("计算牌型结果", userCardInfo.CardInfo)
			userCardInfos = append(userCardInfos, userCardInfo)
		}
	}

	equalUidArrs := [][]int64{}
	findEqualUid := func(uid int64) int {
		for i, vs := range equalUidArrs {
			for _, v := range vs {
				if v == uid {
					return i
				}
			}
		}
		return -1
	}

	ulen := len(userCardInfos)
	for i := 0; i < ulen; i++ {
		for j := i + 1; j < ulen; j++ {
			compRet := CompareCardGroup(userCardInfos[i].CardInfo, userCardInfos[j].CardInfo)
			if compRet == 0 { //相等
				eId := findEqualUid(userCardInfos[i].Uid)
				if eId != -1 {
					equalUidArrs[eId] = append(equalUidArrs[eId], userCardInfos[j].Uid)
				}
			} else if compRet < 0 {
				userCardInfos[i], userCardInfos[j] = userCardInfos[j], userCardInfos[i]
			}
		}
	}

	gameResult := GameResult{
		JackpotVal: this.EDesk.JackpotVal,
	}

	//剩余筹码
	var remainBet int64 = this.EDesk.JackpotVal

	for {
		uInfo := userCardInfos[0]
		p := this.EDesk.GetPlayer(uInfo.Uid)
		isExit := false //是否结束
		eId := findEqualUid(p.Uid)
		if eId == -1 { // 其他人没有跟玩家一样的牌
			userCardInfos = userCardInfos[1:]
			winBet := remainBet
			if p.AllInStage != 0xFF {
				winBet = this.EDesk.GetStageBetTotal(p.Sid)
			} else {
				isExit = true
			}
			remainBet -= winBet
			result := this.genResult(p, uInfo.CardInfo, winBet)
			gameResult.Results = append(gameResult.Results, result)

		} else { // 有人跟他一样的牌，平分筹码
			uids := equalUidArrs[eId]
			ulen := len(uids)
			type UBetInfo struct {
				Uid int64
				Bet int64
			}
			betInfos := []UBetInfo{}

			for _, uid := range uids {
				for i, uc := range userCardInfos {
					if uc.Uid == uid {
						userCardInfos = append(userCardInfos[:i], userCardInfos[i+1:]...)
						break
					}
				}

				p := this.EDesk.GetPlayer(uid)
				winBet := this.EDesk.GetStageBetTotal(p.Sid)
				betInfos = append(betInfos, UBetInfo{
					Uid: p.Uid,
					Bet: winBet,
				})

				if p.AllInStage == 0xFF {
					isExit = true
				}
			}
			//筹码排序
			for i := 0; i < ulen; i++ {
				for j := i + 1; j < ulen; j++ {
					if betInfos[i].Bet > betInfos[j].Bet {
						betInfos[i], betInfos[j] = betInfos[j], betInfos[i]
					}
				}
			}

			lastCount := 0
			//allin的分筹码
			for i := 0; i < ulen; i++ {
				betInfo := betInfos[i]
				p := this.EDesk.GetPlayer(betInfo.Uid)
				if p.AllInStage != 0xFF {
					winBet := betInfo.Bet / int64(ulen)
					remainBet -= winBet
					result := this.genResult(p, uInfo.CardInfo, winBet)
					gameResult.Results = append(gameResult.Results, result)
				} else {
					lastCount += 1
				}
			}
			//坚持到最后的分筹码
			for i := 0; i < ulen; i++ {
				betInfo := betInfos[i]
				p := this.EDesk.GetPlayer(betInfo.Uid)
				if p.AllInStage == 0xFF {
					winBet := remainBet / int64(lastCount)
					remainBet -= winBet
					result := this.genResult(p, uInfo.CardInfo, winBet)
					gameResult.Results = append(gameResult.Results, result)
				}
			}
		}
		if isExit || remainBet <= 0 || len(userCardInfos) == 0 {
			break
		}
	}

	// 统计输钱的玩家
	for _, p := range this.EDesk.Players {
		if p.State != UserStateGameIn {
			continue
		}

		isExist := false
		for _, result := range gameResult.Results {
			if result.Uid == p.Uid {
				isExist = true
				break
			}
		}
		if isExist {
			continue
		}

		cards := []int{}
		cards = append(cards, p.Cards...)
		cards = append(cards, this.EDesk.PublicCards...)

		cstr := ""
		for _, c := range cards {
			cstr = fmt.Sprintf("%s[%d,%d]", cstr, c>>4, c&0xF)
		}
		DebugLog("计算牌型 原数组", cstr, cards)
		cardInfo := CalcCard(cards)
		DebugLog("计算牌型结果", cardInfo)

		result := this.genResult(p, cardInfo, p.GetDownBet()*-1)
		gameResult.Results = append(gameResult.Results, result)
	}

	DebugLog("结算数据", gameResult.Results)

	// 是否存在玩家
	isExistPlayer := false
	for _, v := range this.EDesk.Players {
		if !v.Robot {
			isExistPlayer = true
		}
	}
	//发送结算消息给数据库
	if GetCostType() == 1 && isExistPlayer {
		this.saveToDb(gameResult)
	}
	//发送消息给大厅去记录游戏记录
	if GetCostType() == 1 && isExistPlayer {
		this.saveToHall(gameResult)
	}

	resultJson, _ := json.Marshal(gameResult)
	jb, _ := json.Marshal(resultJson)
	DebugLog("结算数据", string(jb))

	this.EDesk.SendNetMessage(MSG_GAME_NGameResult, gameResult)

	// 重置数据
	index := 0
	for {
		p := this.EDesk.Players[index]
		if p.State == UserStateGameIn {
			p.State = UserStateWaitStart
		}

		if p.LiXian {
			DebugLog("踢出离线玩家 %d", p.Uid, p.Account)
			this.EDesk.HandleDisConnect(p, nil)
		} else if p.CarryCoin < gameConfig.BigBlindCoin {
			if p.Robot {
				this.EDesk.HandleDisConnect(p, nil)
			} else {
				this.EDesk.SendUserSettCoin(p)
				index += 1
			}
		} else {
			index += 1
		}

		if index >= len(this.EDesk.Players) {
			break
		}
	}
}

func (this *FsmResults) saveToDb(gameResult GameResult) {
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.EDesk.JuHao,
		Mini:        false,
		SetLeave:    1,
	}

	for _, p := range this.EDesk.Players {
		if p.State != UserStateGameIn {
			continue
		}
		var prizeCoins int64 = 0
		var WaterProfit float64 = 0
		for _, r := range gameResult.Results {
			if r.Uid == p.Uid {
				prizeCoins, WaterProfit = r.Value, r.WaterProfit
			}
		}

		dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
			UserId:      p.Uid,
			UserAccount: p.Account,
			BetCoins:    p.GetDownBet(),
			ValidBet:    p.GetDownBet(),
			PrizeCoins:  prizeCoins,
			Robot:       p.Robot,
			WaterProfit: WaterProfit,
			WaterRate:   G_DbGetGameServerData.Rate,
		})
		p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
		dbreq.UserCoin = []GGameEndInfo{}
	}
}

func (this *FsmResults) saveToHall(gameResult GameResult) {
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.EDesk.JuHao,
	}

	for _, p := range this.EDesk.Players {
		if p.State != UserStateGameIn {
			continue
		}
		if p.Robot {
			continue
		}
		var prizeCoins int64 = 0
		var waterProfit float64 = 0
		for _, r := range gameResult.Results {
			if r.Uid == p.Uid {
				prizeCoins, waterProfit = r.Value, r.WaterProfit
			}
		}

		rddata := GGameRecordInfo{
			UserId:      p.Uid,
			UserAccount: p.Account,
			Coins:       prizeCoins,
			Score:       int(gameConfig.SmallBlindCoin),
			Multiple:    1,
			CoinsBefore: p.Coins - prizeCoins + int64(waterProfit),
			CoinsAfter:  p.Coins,
			Robot:       p.Robot,
		}
		if p.IsBank {
			rddata.Landlord = true
		} else {
			rddata.Landlord = false
		}
		rdreq.UserRecord = append(rdreq.UserRecord, rddata)
		p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
		rdreq.UserRecord = []GGameRecordInfo{}
	}
}

func (this *FsmResults) genResult(p *ExtPlayer, cardInfo GCardGroupInfo, result int64) Results {
	var value int64 = result
	var waterProfit float64 = 0
	if value > 0 {
		waterProfit = float64(value) * G_DbGetGameServerData.Rate
		value = value - int64(waterProfit)
	}

	if value > 0 {
		p.Coins += value
		p.CarryCoin += value
	}

	resultInfo := Results{
		Uid:           p.Uid,
		Sid:           p.Sid,
		Value:         value,
		CarryCoin:     p.CarryCoin,
		Coin:          p.Coins,
		WaterProfit:   waterProfit,
		Cards:         cardInfo.Cards,
		HandCards:     p.Cards,
		CardGroupType: cardInfo.GroupType,
	}

	return resultInfo
}

func (this *FsmResults) TimerCall(d interface{}) {
	DebugLog("开始下一局游戏")
	this.EDesk.ResetDeskInfo()
	this.EDesk.SendUpdateRoomNo()

	// 检查等待开始玩家的数量是否足够
	waitNum := this.EDesk.GetUserStateTotal(UserStateWaitStart)
	if waitNum > 1 {
		this.EDesk.RunFSM(GameStatusRandBank)
	} else {
		this.EDesk.RunFSM(GameStatusWaitStart)
	}
}

func (this *FsmResults) Leave() {

}

func (this *FsmResults) Reset() {
}

func (this *FsmResults) GetRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FsmResults) OnUserOnline(p *ExtPlayer) {
}

func (this *FsmResults) OnUserOffline(p *ExtPlayer) {
}

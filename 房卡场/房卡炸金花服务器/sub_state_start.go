package main

import (
	"logs"
	"math/rand"
	"time"
)

func (this *ExtDesk) GameStateStart() {
	logs.Debug("开始阶段...")
	//开始阶段
	this.BroadStageTime(STAGE_START_TIME)
	this.BroadcastAll(MSG_GAME_INFO_START_INFO_REPLY, &GSGameStart{
		Id:    MSG_GAME_INFO_START_INFO_REPLY,
		Round: this.GameRound,
	})
	//第一局扣房费
	if this.GameRound == 1 {
		minGoal := this.getPayMoney()
		this.GToHAddRoomCard(this.FkOwner, int64(-minGoal))
	}
	//重置玩家信息
	for _, v := range this.Players {
		this.ResetPlayer(v)
	}

	// 发牌和底注
	this.DealPoker()
	this.CoinPush()
	this.runTimer(STAGE_START_TIME, this.HandlePayCardCoin)
}
func (this *ExtDesk) getPayMoney() int64 {
	var roundValue int64
	if this.TableConfig.TotalRound <= 5 {
		roundValue = 1
	} else if this.TableConfig.TotalRound <= 10 {
		roundValue = 2
	} else if this.TableConfig.TotalRound <= 15 {
		roundValue = 3
	} else {
		roundValue = 4
	}
	var feeValue int64
	var money int64 = 1
	feeValue = roundValue * money
	return feeValue
}

//发牌
func (this *ExtDesk) DealPoker() {
	//洗牌
	this.CardMgr.Shuffle()
	handcard := this.CardMgr.HandCardInfo(len(this.Desk.Players) * 3)
	this.MinCoin = 1

	var newCard [][]int
	for k, v := range this.Players {
		v.OldHandCard = handcard[k]

		card, color := SortHandCard(handcard[k])
		v.HandCards = card
		v.HandColor = color
		// logs.Debug("排序：", v.HandCards)
		lv, _ := GetCardType(v.HandCards, v.HandColor)
		v.CardLv = lv
		this.CoinList = append(this.CoinList, 1)
		v.PayCoin = append(v.PayCoin, 1)
		if !v.Robot && v.CardLv == 1 {
			this.PlayerCard(v)
		}
		newCard = append(newCard, v.OldHandCard)
	}

	this.GetMaxHandCard(newCard)

	this.ProbabilityWinnerRole()
	this.SendMaxCardPlayer()
	// fmt.Println(this.CoinList, "下注列表")
}

//阶段-发牌和底注
func (this *ExtDesk) HandlePayCardCoin(d interface{}) {
	// logs.Debug("阶段-发牌和底注")
	this.nextStage(STAGE_PLAY_OPERATION)
}

func (this *ExtDesk) SendMaxCardPlayer() {
	if len(this.MaxCard) != 2 {
		return
	}
	info := GSMaxCard{
		Id: MSG_GAME_INFO_MAX,
	}
	for k, v := range this.Players {
		if !v.Robot {
			info.PlayerHandCard = append(info.PlayerHandCard, PHandCard{HandCards: v.HandCards, CardLv: v.CardLv, ChairId: v.ChairId})
		}
		if k == this.MaxCard[1] {
			info.CardLv = v.CardLv
			info.HandCard = v.HandCards
			info.ChairId = v.ChairId
		}
	}

	info.WinnerRole = this.WinnerRole
	info.IsRobot = this.MaxCard[0]
	for _, v := range this.Players {
		if v.Robot {
			v.SendNativeMsg(MSG_GAME_INFO_MAX, &info)
		}
	}
}

func (this *ExtDesk) PlayerCard(v *ExtPlayer) {
	rand.Seed(time.Now().UnixNano())
	// logs.Debug("change")
	var newcard []int
	if rand.Intn(100) < 50 {
		newcard = this.ChangeCard(1, v.HandCards, 2)
	}
	if rand.Intn(100) < 10 {
		newcard = this.ChangeCard(1, v.HandCards, 3)
	}

	if newcard == nil || len(newcard) != 3 || len(newcard) == 0 {
		logs.Debug("fail", newcard, rand.Intn(100))
		return
	}
	// logs.Debug("玩家换牌")
	// logs.Debug("长度", len(this.CardMgr.MVSourceCard), this.CardMgr.MVSourceCard)
	this.CardMgr.MVSourceCard = append(this.CardMgr.MVSourceCard, v.OldHandCard[:]...) // 还牌
	this.CardMgr.MVSourceCard = DuplicateRemoval(this.CardMgr.MVSourceCard)
	// logs.Debug("长度", len(this.CardMgr.MVSourceCard), this.CardMgr.MVSourceCard)
	v.OldHandCard = newcard
	v.HandCards, v.HandColor = SortHandCard(newcard)
	v.CardLv, _ = GetCardType(v.HandCards, v.HandColor)
}

func (this *ExtDesk) ProbabilityWinnerRole() {
	for _, v := range this.Players {
		if !v.Robot {
			rand.Seed(time.Now().UnixNano())
			if v.CardLv == 1 && v.HandCards[2] <= 10 && rand.Intn(1000) <= 5 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 1 && v.HandCards[2] <= 14 && rand.Intn(100) <= 10 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 2 && v.HandCards[1] <= 10 && rand.Intn(100) <= 15 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 2 && rand.Intn(100) <= 20 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 3 && v.HandCards[2] <= 10 && rand.Intn(100) <= 30 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 3 && rand.Intn(100) <= 40 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 4 && v.HandCards[2] <= 10 && rand.Intn(100) <= 50 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 4 && rand.Intn(100) <= 60 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 5 && v.HandCards[2] <= 10 && rand.Intn(100) <= 70 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 5 && rand.Intn(100) <= 70 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 6 && v.HandCards[2] <= 10 && rand.Intn(100) <= 90 {
				this.WinnerRole = 0
				return
			} else if v.CardLv == 6 && rand.Intn(100) <= 95 {
				this.WinnerRole = 0
				return
			} else {
				this.WinnerRole = 1
			}
		}
	}
}

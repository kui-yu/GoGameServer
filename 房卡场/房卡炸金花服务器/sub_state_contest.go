package main

import "logs"

func (this *ExtDesk) GameStateContest() {
	logs.Debug("比牌阶段")
	this.BroadStageTime(STAGE_CONTEST_TIME)
	this.CoinEnd()
	this.runTimer(STAGE_CONTEST_TIME, this.HandleGameContest)
}

//金币不足比牌
func (this *ExtDesk) CoinEnd() {
	if this.GameState == STAGE_SETTLE || this.GameState == GAME_STATUS_END {
		return
	}
	var p *ExtPlayer
	sumCoin := int64(0)
	for _, v := range this.Players {
		if this.CallPlayer == v.ChairId {
			p = v
			for i := 0; i < len(v.PayCoin); i++ {
				sumCoin += v.PayCoin[i]
			}
			break
		}
	}

	// fmt.Println("比牌")
	info := GSPlayerContest{
		Id: MSG_GAME_INFO_CONTEST,
	}

	followIn := (p.Coins - sumCoin*this.Bscore) / this.Bscore
	this.CoinList = append(this.CoinList, followIn) //金币记录处理
	p.PayCoin = append(p.PayCoin, followIn)
	info.PlayCoin = followIn

	var allPlayer []*ExtPlayer
	outKey := false
	for i := 0; i < len(this.ChairList); i++ { //取出比牌对象顺序
		if p.ChairId == this.ChairList[i] && outKey {
			break
		}
		if outKey && this.Players[i].CardType != 2 {
			allPlayer = append(allPlayer, this.Players[i])
		}
		if p.ChairId == this.ChairList[i] {
			outKey = true
			if i == len(this.ChairList)-1 {
				i = -1
			}
			continue
		} else if outKey && i == len(this.ChairList)-1 {
			i = -1
		}
	}

	for _, v := range allPlayer {
		pc := PlayerContest{}
		pc.Person_1 = p.ChairId
		result := GetResult(p.HandCards, p.HandColor, v.HandCards, v.HandColor)
		pc.Person_2 = v.ChairId
		if result == 0 { //给输家看自己手牌，继续比牌
			pc.Winner = p.ChairId
			v.CardType = 2
			pc.LoserCard = v.OldHandCard
			pc.CardLv = v.CardLv
			this.SettleContest = append(this.SettleContest, pc)
		} else { //该玩家输牌，给该玩家看牌
			pc.Winner = v.ChairId
			p.CardType = 2
			pc.LoserCard = p.OldHandCard
			pc.CardLv = p.CardLv
			this.SettleContest = append(this.SettleContest, pc)
			break
		}
	}

	info.Count = len(this.SettleContest)
	info.PContest = append(info.PContest, this.SettleContest[:]...)

	count := 0
	for _, v := range this.Players {
		if v.CardType != 2 {
			count++
		}
	}

	//下注金币总数推送
	this.CoinPush()
	if count >= 2 {
		this.BroadcastAll(MSG_GAME_INFO_CONTEST, &info)
	} else {
		this.nextStage(GAME_STATUS_END)
	}
}

//阶段-比牌
func (this *ExtDesk) HandleGameContest(d interface{}) {
	// logs.Debug("阶段-比牌")
	pcount := 0
	for _, v := range this.Players {
		if v.CardType != 2 {
			pcount++
		}
	}
	if pcount < 2 {
		return
	}
	count := len(this.SettleContest) * 3
	this.runTimer(count, this.ContestWaitTime)
}

func (this *ExtDesk) ContestWaitTime(d interface{}) {
	this.SettleContest = []PlayerContest{}
	this.nextStage(STAGE_PLAY_OPERATION)
	this.MsgCallPlayer()
}

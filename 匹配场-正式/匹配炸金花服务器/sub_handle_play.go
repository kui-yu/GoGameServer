package main

import (
	"encoding/json"
	"logs"
)

//玩家下注,加注或者跟注比牌
func (this *ExtDesk) HandleGamePlay(p *ExtPlayer, d *DkInMsg) {
	if this.Round == GameRound || this.GameState == STAGE_SETTLE || this.GameState == GAME_STATUS_END {
		return
	}
	if p.ChairId != this.CallPlayer || p.CardType == 2 {
		return
	}
	data := GAPlayerOperation{
		ChairId: -1,
	}
	json.Unmarshal([]byte(d.Data), &data)
	if data.ChairId == -1 {
		return
	}

	if data.PlayCoin < this.MinCoin && data.Operation == 3 {
		return
	}

	if data.Operation == 2 && this.Round < 2 {
		data.Operation = 4
	}

	sumCoin := int64(0)
	for i := 0; i < len(p.PayCoin); i++ {
		sumCoin += p.PayCoin[i]
	}

	if (p.Coins <= this.Bscore*(sumCoin+this.MinCoin) && p.CardType == 0) || (p.Coins <= this.Bscore*(sumCoin+this.MinCoin*2) && p.CardType == 1) { //金币不足，强制比牌 金币不足判断重写
		logs.Debug("金币不足进入比牌", data)
		this.nextStage(STAGE_SOLO)
		return
	} else if data.Operation == 4 { //跟注
		info := GSPlayerPayCoin{
			Id:        MSG_GAME_INFO_PLAY_INFO_REPLY,
			PChairId:  p.ChairId,
			Operation: 4,
		}
		if p.CardType == 0 {
			this.CoinList = append(this.CoinList, this.MinCoin)
			p.PayCoin = append(p.PayCoin, this.MinCoin)
			info.PlayCoin = this.MinCoin
			this.BroadcastAll(MSG_GAME_INFO_PLAY_INFO_REPLY, &info)
		} else if p.CardType == 1 {
			this.CoinList = append(this.CoinList, this.MinCoin*2)
			p.PayCoin = append(p.PayCoin, this.MinCoin*2)
			info.PlayCoin = this.MinCoin * 2
			this.BroadcastAll(MSG_GAME_INFO_PLAY_INFO_REPLY, &info)
		}

	} else if data.Operation == 3 { //加注
		if data.PlayCoin%2 != 0 && p.CardType == 1 {
			data.PlayCoin = data.PlayCoin + 1
		}
		info := GSPlayerPayCoin{
			Id:        MSG_GAME_INFO_PLAY_INFO_REPLY,
			PChairId:  p.ChairId,
			PlayCoin:  data.PlayCoin,
			Operation: 3,
		}
		if p.CardType == 0 {
			this.MinCoin = data.PlayCoin
		} else if p.CardType == 1 {
			this.MinCoin = data.PlayCoin / 2
		}
		this.CoinList = append(this.CoinList, data.PlayCoin)
		p.PayCoin = append(p.PayCoin, data.PlayCoin)
		this.BroadcastAll(MSG_GAME_INFO_PLAY_INFO_REPLY, &info)

	} else if data.Operation == 2 && this.Round >= 2 { //比牌
		// fmt.Println("比牌")
		info := GSPlayerPayCoin{
			Id:        MSG_GAME_INFO_PLAY_INFO_REPLY,
			PChairId:  p.ChairId,
			ChairId:   []int32{p.ChairId},
			Operation: 2,
		}

		if p.CardType == 0 {
			this.CoinList = append(this.CoinList, this.MinCoin)
			p.PayCoin = append(p.PayCoin, this.MinCoin)
			info.PlayCoin = this.MinCoin
		} else {
			this.CoinList = append(this.CoinList, this.MinCoin*2)
			p.PayCoin = append(p.PayCoin, this.MinCoin*2)
			info.PlayCoin = this.MinCoin * 2
		}

		var LookObj *ExtPlayer
		for _, v := range this.Players {
			if data.ChairId == v.ChairId {
				if v.CardType == 2 {
					return
				}
				result := GetResult(p.HandCards, p.HandColor, v.HandCards, v.HandColor)
				info.ChairId = append(info.ChairId, v.ChairId)

				playinfo := GSCardInfo{
					Id: MSG_GAME_INFO_LOOK_CARD,
				}
				if result == 0 { //输家看牌
					info.Winner = p.ChairId
					playinfo.HandCards = v.OldHandCard
					playinfo.Lv = v.CardLv
					playinfo.ChairId = v.ChairId
					playinfo.Model = 1
					v.CardType = 2
					// fmt.Println(p.OldHandCard, v.OldHandCard, "查看手牌")
					LookObj = v
				} else {
					info.Winner = v.ChairId
					playinfo.HandCards = p.OldHandCard
					playinfo.Lv = p.CardLv
					playinfo.ChairId = p.ChairId
					p.CardType = 2
					playinfo.Model = 1
					// fmt.Println(p.OldHandCard, v.OldHandCard, "查看手牌")
					LookObj = p
				}

				alive := 0
				for _, v := range this.Players {
					if v.CardType != 2 {
						alive++
					}
				}
				if alive >= 2 {
					LookObj.SendNativeMsg(MSG_GAME_INFO_LOOK_CARD, &playinfo)
					this.BroadcastAll(MSG_GAME_INFO_PLAY_INFO_REPLY, &info)
				} else {
					// fmt.Println("操作 out of range")
					this.SettleContest = make([]PlayerContest, 0, 1)
					this.SettleContest = append(this.SettleContest, PlayerContest{Person_1: p.ChairId, Person_2: v.ChairId, Winner: info.Winner})
				}
				break
			}
		}
	}

	//下注金币总数推送
	this.CoinPush()
	//指针后移
	this.MsgCallPlayer()
}

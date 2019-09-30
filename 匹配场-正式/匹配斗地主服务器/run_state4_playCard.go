package main

import (
	"encoding/json"
	// "logs"
	"math/rand"
	"time"
)

//玩牌阶段 指定玩家出牌
func (this *ExtDesk) TimerPlayCard() {

	// logs.Debug("玩家", this.CurCid)
	//添加定时器，进入出牌阶段
	nextplayer := this.Players[this.CurCid]

	if nextplayer.TuoGuan {
		//托管玩家
		this.AddTimer(TIMER_OUTCARD, 1, this.TuoGuanOut, nil)
	} else {
		if nextplayer.Robot {
			//机器人随机出牌
			rand.Seed(time.Now().UnixNano())
			timerOutCardNum := rand.Intn(3) + 2

			//机器人出牌
			this.AddTimer(TIMER_OUTCARD, timerOutCardNum, this.robotPlayCard, nextplayer)
		} else {
			//等待真实玩家出牌
			this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCard, nil)
		}
	}
}
func (this *ExtDesk) TuoGuanOut(d interface{}) {
	// logs.Debug(".超时出牌定时器触发")
	p := this.Players[this.CurCid]
	p.Pass = false
	cards := p.HandCard
	//排序
	cards = Sort(cards)

	if this.MaxChuPai == nil {

		req := GGameOutCard{
			Id: MSG_GAME_INFO_OUTCARD,
		}

		// 智能出牌
		var foeCards []byte
		if p.ChairId == this.Players[this.Banker].ChairId {
			//上一个玩家
			lastCid := p.ChairId - 1
			if lastCid < 0 {
				lastCid = 2
			}
			//下一个玩家
			nextCid := p.ChairId + 1
			if nextCid > 2 {
				nextCid = 0
			}
			if len(this.Players[lastCid].HandCard) < len(this.Players[nextCid].HandCard) {
				foeCards = this.Players[lastCid].HandCard
			} else {
				foeCards = this.Players[nextCid].HandCard
			}
		} else {
			foeCards = this.Players[this.Banker].HandCard
		}

		playCards := R_OTOffensive(p.HandCard, foeCards)
		for _, v := range playCards {
			req.Cards = append(req.Cards, int32(v))
		}
		// //出最小的牌
		// req.Cards = append(req.Cards, int32(cards[len(cards)-1]))

		dv, _ := json.Marshal(req)
		this.HandleGameOutCard(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(dv),
		})

		return
	} else {
		//托管出牌
		//=============================================
		// playCards := this.TimerOutPlayCard(this.MaxChuPai, cards)
		_, playCards := R_DefPosition1(this.MaxChuPai, p.HandCard)
		if len(playCards) > 0 {
			req := GGameOutCard{
				Id: MSG_GAME_INFO_OUTCARD,
			}
			for _, v := range playCards {
				req.Cards = append(req.Cards, int32(v))
			}
			dv, _ := json.Marshal(req)
			this.HandleGameOutCard(p, &DkInMsg{
				Uid:  p.Uid,
				Data: string(dv),
			})
		} else {
			//不是第一个出牌就过
			this.HandlePass(p, &DkInMsg{
				Uid: p.Uid,
			})
			p.Pass = true
		}

	}
}

//玩牌阶段 超时出牌
func (this *ExtDesk) TimerOutCard(d interface{}) {
	// logs.Debug(".超时出牌定时器触发")

	p := this.Players[this.CurCid]
	p.Pass = false
	if !p.TuoGuan {
		p.TuoGuan = true
		this.BroadcastAll(MSG_GAME_INFO_TUOGUAN_REPLY, &GTuoGuanReply{
			Id:     MSG_GAME_INFO_TUOGUAN_REPLY,
			Ctl:    1,
			Result: 0,
			Cid:    p.ChairId,
		})
	}

	cards := p.HandCard
	//排序
	cards = Sort(cards)

	if this.MaxChuPai == nil {

		req := GGameOutCard{
			Id: MSG_GAME_INFO_OUTCARD,
		}

		// 智能出牌
		var foeCards []byte
		if p.ChairId == this.Players[this.Banker].ChairId {
			//上一个玩家
			lastCid := p.ChairId - 1
			if lastCid < 0 {
				lastCid = 2
			}
			//下一个玩家
			nextCid := p.ChairId + 1
			if nextCid > 2 {
				nextCid = 0
			}
			if len(this.Players[lastCid].HandCard) < len(this.Players[nextCid].HandCard) {
				foeCards = this.Players[lastCid].HandCard
			} else {
				foeCards = this.Players[nextCid].HandCard
			}
		} else {
			foeCards = this.Players[this.Banker].HandCard
		}

		playCards := R_OTOffensive(p.HandCard, foeCards)
		for _, v := range playCards {
			req.Cards = append(req.Cards, int32(v))
		}
		// //出最小的牌
		// req.Cards = append(req.Cards, int32(cards[len(cards)-1]))

		dv, _ := json.Marshal(req)
		this.HandleGameOutCard(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(dv),
		})

		return
	} else {
		//托管出牌
		//=============================================
		// playCards := this.TimerOutPlayCard(this.MaxChuPai, cards)
		_, playCards := R_DefPosition1(this.MaxChuPai, p.HandCard)
		if len(playCards) > 0 {
			req := GGameOutCard{
				Id: MSG_GAME_INFO_OUTCARD,
			}
			for _, v := range playCards {
				req.Cards = append(req.Cards, int32(v))
			}
			dv, _ := json.Marshal(req)
			this.HandleGameOutCard(p, &DkInMsg{
				Uid:  p.Uid,
				Data: string(dv),
			})
		} else {
			//不是第一个出牌就过
			this.HandlePass(p, &DkInMsg{
				Uid: p.Uid,
			})
			p.Pass = true
		}

	}
}

package main

import (
	// "github.com/astaxie/beego/logs"
	"encoding/json"
	"logs"
	//	"sort"
)

//玩家出牌
func (this *ExtDesk) HandleGameOutCard(p *ExtPlayer, d *DkInMsg) {
	// logs.Debug("玩家出牌", p.Account, string(d.Data))

	if this.GameState != GAME_STATUS_PLAY {
		logs.Error("出牌状态错误:", this.GameState, GAME_STATUS_PLAY)
		return
	}
	if this.CurCid != p.ChairId {
		logs.Error("不是当前玩家出牌:", this.CurCid, p.ChairId)
		return
	}
	//////
	rdata := GGameOutCard{}
	json.Unmarshal([]byte(d.Data), &rdata)

	//
	req := GOutCard{
		Cid: p.ChairId,
	}
	for _, v := range rdata.Cards {
		req.Cards = append(req.Cards, byte(v))
	}
	Sort(req.Cards)
	if !this.CheckStyle(req.Type, req.Cards, &req) {
		logs.Error("出牌类型检测错误:", req)
		return
	}
	vhand := append([]byte{}, p.HandCard...)
	if this.MaxChuPai != nil {
		if !this.CampareChuPai(&req) {
			logs.Error("出牌不大于最大牌:", req, this.MaxChuPai)
			return
		}
	} /*else {
		//清空出牌记录
		this.RdChuPai = []*GOutCard{}
	}*/
	nh, ok := VecDelMulti(vhand, req.Cards)
	if !ok {
		logs.Error("没有这些手牌:", req)
		return
	}
	p.HandCard = nh
	this.MaxChuPai = &req
	p.Outed = append(p.Outed, &req)
	this.RdChuPai = append(this.RdChuPai, &req)
	//
	outdouble := 1
	if req.Type == CT_TWOKING || req.Type == CT_BOMB_FOUR {
		outdouble = 2
		// banker := this.Players[this.Banker]
		for _, v := range this.Players {
			v.Double *= 2
			if v.ChairId == this.Banker && v.Double > this.MaxDouble*2 {
				v.Double = this.MaxDouble * 2
			} else if v.ChairId != this.Banker && v.Double > this.MaxDouble {
				v.Double = this.MaxDouble
			}
		}
	}
	this.CurCid = (this.CurCid + 1) % int32(len(this.Players))

	//出牌回复
	rsp := GGameOutCardReply{
		Id:     MSG_GAME_INFO_OUTCARD_REPLY,
		Cid:    p.ChairId,
		Type:   req.Type,
		Cards:  req.Cards,
		Max:    req.Max,
		Double: int32(outdouble),
	}
	this.BroadcastAll(MSG_GAME_INFO_OUTCARD_REPLY, &rsp)

	this.DelTimer(TIMER_OUTCARD)

	//判断是否结束了
	if len(p.HandCard) == 0 {
		// logs.Debug("游戏结束")
		this.GameOver(p)
	} else {
		nextplayer := this.Players[this.CurCid]
		if nextplayer.TuoGuan {
			this.AddTimer(TIMER_OUTCARD, 1, this.TuoGuanOut, nil)
		} else {

			this.TimerPlayCard()
			// this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCard, nil)
		}
	}
	p.Pass = false
}

//检查牌型
func (this *ExtDesk) CheckStyle(style int32, out []byte, outc *GOutCard) bool {
	re := DoGenCard(out, outc)
	if !re {
		return re
	}
	if this.MaxChuPai != nil {
		if this.MaxChuPai.Type != outc.Type && outc.Type < CT_BOMB_FOUR {
			return false
		}
		if this.MaxChuPai.Type != outc.Type && outc.Type >= CT_BOMB_FOUR && outc.Type < this.MaxChuPai.Type {
			return false
		}
	}
	return true
}

//比较出牌
func (this *ExtDesk) CampareChuPai(chupai *GOutCard) bool {
	if this.MaxChuPai.Type == CT_TWOKING {
		return true
	}
	if this.MaxChuPai.Type != chupai.Type {
		if chupai.Type < CT_BOMB_FOUR {
			return false
		} else if this.MaxChuPai.Type > chupai.Type {
			return false
		}
	} else {
		if GetLogicValue(this.MaxChuPai.Max) >= GetLogicValue(chupai.Max) {
			return false
		}
	}
	return true
}

package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) handleOutCardSelect_lz(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != GAME_STATUS_PLAY {
		logs.Error("游戏阶段错误———选择牌型")
		return
	}
	if this.CurCid != p.ChairId {
		logs.Error("不是当前玩家出牌—————选择牌型")
		return
	}
	outcard := CanOutType{}
	err := json.Unmarshal([]byte(d.Data), &outcard)
	if err != nil {
		logs.Error("转换消息 类型错误————选择牌型", err)
	}
	logs.Debug("outcar————选择牌型:", outcard)
	vhand := p.HandCard
	nh, ok1 := VecDelMulti(vhand, outcard.Cards)
	if !ok1 {
		logs.Error("没有这些手牌!")
		return
	}
	p.HandCard = nh
	this.MaxChuPai = &GOutCard{
		Cid:   p.ChairId,
		Max:   outcard.Max,
		Type:  int32(outcard.Max),
		Cards: outcard.Cards,
	}
	p.Outed = append(p.Outed, this.MaxChuPai)
	this.RdChuPai = append(this.RdChuPai, this.MaxChuPai)

	outdouble := 1
	//判断是否炸弹，如果是翻倍
	if this.MaxChuPai.Type == CT_TWOKING || this.MaxChuPai.Type == CT_BOMB_FOUR {
		var bomcout int
		if this.TableConfig.Boom == -1 {
			bomcout = 999
		} else {
			bomcout = this.TableConfig.Boom
		}
		if BoomCout < bomcout { //如果小于规定的炸弹 才会翻倍
			BoomCout++
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
	}
	rsp := GGameOutCardReply{
		Id:     MSG_GAME_INFO_OUTCARD_REPLY,
		Ishas:  true,
		Cid:    p.ChairId,
		Type:   this.MaxChuPai.Type,
		Max:    this.MaxChuPai.Max,
		Double: int32(outdouble),
		// LzAndBecome: outcard.LzBecome,
		// Ptcon:       outcard.Ptcon,
		NextCid: (this.CurCid + 1) % int32(len(this.Players)),
	}
	this.CurCid = (this.CurCid + 1) % int32(len(this.Players))

	var lzcon []byte
	lzType := 0x50
	var cd []byte
	for _, v := range outcard.Ptcon {
		cd = append(cd, v)
	}
	for _, v := range outcard.LzBecome {
		cd = append(cd, byte(lzType)+v.Become)
		lzcon = append(lzcon, v.Lz)
	}
	for _, v := range cd {
		rsp.CardsLz = append(rsp.CardsLz, int(v))
	}
	for index := range rsp.CardsLz {
		rsp.CardsLz[index] = 0
	}
	Sort(cd)
	for _, v := range cd {
		rsp.Cards = append(rsp.Cards, int(v))
	}
	var NeedIndex []int
	for i, v := range cd {
		if GetCardColor(v) > 4 {
			NeedIndex = append(NeedIndex, i)
		}
	}
	for i, v := range NeedIndex {
		rsp.CardsLz[v] = int(lzcon[i])
	}
	this.BroadcastAll(MSG_GAME_INFO_OUTCARD_REPLY, &rsp)
	this.DelTimer(TIMER_OUTCARD)
	//判断是否结束了
	if len(p.HandCard) == 0 {
		// logs.Debug("游戏结束")
		this.GameState = GAME_STATUS_BALANCE
		this.BroadStageTime(0)
		this.GameOver(p)
		//重新设置桌子数据，等所有玩家离开，归还桌子
		//
	} else {
		nextplayer := this.Players[this.CurCid]
		if nextplayer.TuoGuan {
			this.AddTimer(TIMER_OUTCARD, 1, this.TimerOutCard, nil)
		} else {
			this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCard, nil)
		}
	}
}

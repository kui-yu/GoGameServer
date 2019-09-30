package main

import (
	"encoding/json"
	"fmt"
	"logs"
)

//处理玩家出牌
func (this *ExtDesk) HandleOutCard(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("接收到玩家出牌请求", p.Nick)
	if this.GameState != GAME_STATUS_OUTCARD {
		logs.Error("不是出牌阶段", this.GameState)
		p.SendNativeMsg(MSG_GAME_INFO_OUTCARD_REPLY, OutCardsReply{
			Id:     MSG_GAME_INFO_OUTCARD_REPLY,
			Result: 1,
			Err:    "不是出牌阶段",
		})
		return
	}
	if this.CurCid != p.ChairId {
		logs.Error("不是轮到该玩家出牌", this.CurCid, p.ChairId)
		p.SendNativeMsg(MSG_GAME_INFO_OUTCARD_REPLY, OutCardsReply{
			Id:     MSG_GAME_INFO_OUTCARD_REPLY,
			Result: 2,
			Err:    "不是轮到您出牌",
		})
		return
	}
	p1 := &ExtPlayer{}
	p1.HandCards = append(p1.HandCards, p.HandCards...)
	//获取玩家出牌牌组
	data := OutCard{}
	json.Unmarshal([]byte(d.Data), &data)
	outc := []byte{}
	//将cards转型为byte存入
	for _, v := range data.Cards {
		outc = append(outc, byte(v))
	}
	outc = Sort(outc)
	//检测出牌是否合法
	ok, lastout := Check_Chard(outc, len(p.HandCards))
	if !ok {
		logs.Error("该玩家出牌类型不合法！！", lastout.Cards)
		p.SendNativeMsg(MSG_GAME_INFO_OUTCARD_REPLY, OutCardsReply{
			Id:     MSG_GAME_INFO_OUTCARD_REPLY,
			Result: 3,
			Err:    "您的出牌类型不合法",
		})
		return
	}
	lastout.Cid = p.ChairId
	if !this.comp(len(outc), lastout.Type, lastout.Max, len(p.HandCards)) {
		logs.Error("您的出牌未能大过上家", lastout.Cards, this.LastOutCards)
		p.SendNativeMsg(MSG_GAME_INFO_OUTCARD_REPLY, OutCardsReply{
			Id:     MSG_GAME_INFO_OUTCARD_REPLY,
			Result: 4,
			Err:    "您的出牌未大过上家",
		})
		return
	}
	hcard := append([]byte{}, p.HandCards...)
	hc, ok2 := VecDelMulti(hcard, lastout.Cards)
	if !ok2 {
		logs.Error("没有这些手牌", lastout.Cards)
		p.SendNativeMsg(MSG_GAME_INFO_OUTCARD_REPLY, OutCardsReply{
			Id:     MSG_GAME_INFO_OUTCARD_REPLY,
			Result: 5,
			Err:    "没有这些手牌!",
		})
		return
	}
	this.DelTimer(GAME_STATUS_OUTCARD)
	//设置下一出牌玩家
	this.CurCid = (this.CurCid + 1) % int32(len(this.Players))
	nextPlayer := this.Players[this.CurCid]
	_, _, result := this.CanOutCards(p1)
	//判断包赔，如果本次出牌类型为单数,并且下家为
	if lastout.Type == CT_SINGLE && nextPlayer.IsDan {
		// 判断是否包赔  未完成
		if GetLogicValue(result[len(result)-1][0]) > GetLogicValue(lastout.Cards[0]) {
			logs.Debug("玩家并没有出最大的单牌，如果放走下家，则包赔")
			p.IsBaoPei = true
		} else {
			logs.Debug("玩家出的最大单牌，即使下家出完牌，也不需要包赔")
		}
	}
	p.HandCards = hc
	this.LastOutCards = lastout
	this.OutCards[p.ChairId] = lastout
	//获取下一个玩家的提示

	//判断玩家手牌 是否只剩一张，是的话需要保单
	isDan := false
	if len(p.HandCards) == 1 {
		isDan = true
		p.IsDan = true
	}
	//广播玩家出牌
	outcardsbr := OutCardsBro{
		Id:      MSG_GAME_INFO_OUTCARD_BRO,
		Cid:     p.ChairId,
		Type:    this.LastOutCards.Type,
		Max:     this.LastOutCards.Max,
		NextCid: this.CurCid,
		IsDan:   isDan,
	}
	for _, v1 := range this.LastOutCards.Cards {
		outcardsbr.Cards = append(outcardsbr.Cards, int(v1))
	}
	ts := [][]int{}
	for _, v := range this.Players {
		if v.ChairId == this.CurCid {
			tishi, _, _ := this.CanOutCards(v)
			for _, v := range tishi {
				list := []int{}
				for _, v1 := range v {
					list = append(list, int(v1))
				}
				outcardsbr.Hint = append(outcardsbr.Hint, list)
			}
			ts = append(ts, outcardsbr.Hint...)
		}
		v.SendNativeMsg(MSG_GAME_INFO_OUTCARD_BRO, outcardsbr)
		outcardsbr.Hint = [][]int{}
	}

	if len(p.HandCards) == 0 {
		this.GameOver(p)
	} else {
		if len(ts) <= 0 {
			this.AddTimer(GAME_STATUS_OUTCARD, 2, this.TimerPass, nil)
		} else {
			if nextPlayer.TuoGuan {
				this.AddTimer(GAME_STATUS_OUTCARD, 1, this.TuoGuanOut, nil)
			} else {
				this.AddTimer(GAME_STATUS_OUTCARD, GAME_STATUS_OUTCARD_TIME, this.TimerOutCard, nil)
			}
		}
	}
	p.Pass = false
}

func (this *ExtDesk) TuoGuanOut(d interface{}) {
	p := this.Players[this.CurCid]
	p.Pass = false
	cards := append([]byte{}, p.HandCards...)
	//排序
	cards = Sort(cards)
	out := this.TuoGuanOutCards(p)
	out = Sort(out)
	for _, v := range out {
		fmt.Print(v)
	}
	fmt.Println("")
	req := OutCard{
		Id: MSG_GAME_INFO_OUTCARD,
	}
	if len(out) <= 0 {
		this.AddTimer(GAME_STATUS_OUTCARD, 2, this.TimerPass, nil)
	} else {
		for _, v := range out {
			req.Cards = append(req.Cards, int32(v))
		}
		data, err := json.Marshal(req)
		if err != nil {
			logs.Error("出牌超时方法中json转换W错误")
			return
		}
		this.HandleOutCard(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(data),
		})
	}

}

//出牌阶段超时处理
func (this *ExtDesk) TimerOutCard(d interface{}) {
	p := this.Players[this.CurCid]
	p.Pass = false
	if !p.TuoGuan {
		//将玩家设置为托管
		p.TuoGuan = true
		this.BroadcastAll(MSG_GAME_INFO_TUOGUAN_BRO, TuoGuanReply{
			Id:  MSG_GAME_INFO_TUOGUAN_BRO,
			Ctl: 1,
			Cid: p.ChairId,
		})
	}
	cards := append([]byte{}, p.HandCards...)
	//排序
	cards = Sort(cards)
	out := this.TuoGuanOutCards(p)
	out = Sort(out)
	for _, v := range out {
		fmt.Print(v)
	}
	fmt.Println("")
	req := OutCard{
		Id: MSG_GAME_INFO_OUTCARD,
	}
	if len(out) <= 0 {
		this.AddTimer(GAME_STATUS_OUTCARD, 2, this.TimerPass, nil)
	} else {
		for _, v := range out {
			req.Cards = append(req.Cards, int32(v))
		}
		data, err := json.Marshal(req)
		if err != nil {
			logs.Error("出牌超时方法中json转换W错误")
			return
		}
		this.HandleOutCard(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(data),
		})
	}
}

//玩家没有大过上家的牌，自动帮玩家过牌
func (this *ExtDesk) TimerPass(d interface{}) {
	nextPlayer := this.Players[this.CurCid]
	this.HandlePass(nextPlayer, &DkInMsg{
		Uid: nextPlayer.Uid,
	})
}

//玩家出牌与上家出牌对比
func (this *ExtDesk) comp(cardslen int, cardtype int, max byte, handlen int) bool {
	//是否存在上一家出牌
	if this.LastOutCards.Max == 0 {
		//不存在上一家出牌，所以直接出牌成功
		return true
	}
	//判断类型是否一样，或者是炸弹
	if cardtype == this.LastOutCards.Type {
		if cardslen == len(this.LastOutCards.Cards) && GetLogicValue(max) > GetLogicValue(this.LastOutCards.Max) {
			return true
		}
	}
	//如果玩家出牌炸弹
	if cardtype == CT_BOMB_FOUR {
		return true
	}

	if handlen == 3 && cardtype == CT_THREE || handlen == 4 && cardtype == CT_THREE_LINE_TAKE_ONE {
		if GetLogicValue(max) >= GetLogicValue(this.LastOutCards.Max) {
			return true
		}
	}
	logs.Debug("出牌时对比不过")
	return false
}

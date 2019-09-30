package main

import (
	"fmt"
	// "github.com/astaxie/beego/logs"
	"encoding/json"
	"logs"
	//	"sort"
)

var BoomCout int //用来计算炸的次数
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
	fmt.Println("req.Cards", req.Cards)
	fmt.Println("req.Max", req.Max)
	fmt.Println("req.Type", req.Type)
	var SanDaiEr bool = true
	var SiDaiEr bool = true
	if this.TableConfig.CanSelect[0] == 1 {
		//不可三带二
		SanDaiEr = false
	}
	if this.TableConfig.CanSelect[1] == 1 {
		//不可四带二对
		SanDaiEr = false
	}

	if !this.CheckStyle(req.Type, req.Cards, &req, SanDaiEr, SiDaiEr) {
		logs.Error("出牌类型检测错误:", req)
		return
	}
	vhand := append([]byte{}, p.HandCard...)
	if this.MaxChuPai != nil {
		if !this.CampareChuPai(&req) {
			logs.Error("出牌不大于最大牌:", req, this.MaxChuPai)
			return
		}
	}
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
	//判断是否炸弹，如果是翻倍
	if req.Type == CT_TWOKING || req.Type == CT_BOMB_FOUR {
		logs.Debug("发现是炸弹")
		var bomcout int
		if this.TableConfig.Boom == -1 {
			bomcout = 999
		} else {
			bomcout = this.TableConfig.Boom
		}
		if BoomCout < bomcout { //如果小于规定的炸弹 才会翻倍
			logs.Debug("进来了！")
			BoomCout++
			outdouble = 2
			// banker := this.Players[this.Banker]
			this.Double *= 2
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
		Id:      MSG_GAME_INFO_OUTCARD_REPLY,
		Cid:     p.ChairId,
		Type:    req.Type,
		Max:     req.Max,
		Double:  int32(outdouble),
		NextCid: (this.CurCid + 1) % int32(len(this.Players)),
	}
	fmt.Println("出牌倍数：", rsp.Double)
	this.CurCid = (this.CurCid + 1) % int32(len(this.Players))
	for _, v := range req.Cards {
		rsp.Cards = append(rsp.Cards, int(v))
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
			this.AddTimer(TIMER_OUTCARD, 1, this.TuoGuanOut, nil)
		} else {
			this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCard, nil)
		}
	}

}
func (this *ExtDesk) TuoGuanOut(d interface{}) {
	// logs.Debug(".超时出牌定时器触发")
	logs.Debug("进入出牌阶段")
	p := this.Players[this.CurCid]
	if this.MaxChuPai == nil { //如果托管玩家为第一个出牌的人,默认只先出第一张牌
		req := GGameOutCard{
			Id: MSG_GAME_INFO_OUTCARD,
		}
		req.Cards = append(req.Cards, int32(p.HandCard[len(p.HandCard)-1]))
		dv, _ := json.Marshal(req)
		this.HandleGameOutCard(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(dv),
		})
		return
	} else { //如果托管玩家上面已经有人出牌
		this.HandlePass(p, &DkInMsg{
			Uid: p.Uid,
		})
	}
	if this.TableConfig.GameType == 4 {
		logs.Debug("癞子模式出牌超时")
	} else if this.TableConfig.GameType == 1 {
		logs.Debug("普通模式出牌超时")
	}
}

func (this *ExtDesk) TimerOutCard(d interface{}) {
	// logs.Debug(".超时出牌定时器触发")
	logs.Debug("进入出牌阶段")
	p := this.Players[this.CurCid]
	if !p.TuoGuan {
		p.TuoGuan = true
		this.BroadcastAll(MSG_GAME_INFO_TUOGUAN_REPLY, &GTuoGuanReply{
			Id:     MSG_GAME_INFO_TUOGUAN_REPLY,
			Ctl:    1,
			Result: 0,
			Cid:    p.ChairId,
		})
	}

	if this.MaxChuPai == nil { //如果托管玩家为第一个出牌的人,默认只先出第一张牌
		req := GGameOutCard{
			Id: MSG_GAME_INFO_OUTCARD,
		}
		req.Cards = append(req.Cards, int32(p.HandCard[len(p.HandCard)-1]))
		dv, _ := json.Marshal(req)
		this.HandleGameOutCard(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(dv),
		})
		return
	} else { //如果托管玩家上面已经有人出牌
		this.HandlePass(p, &DkInMsg{
			Uid: p.Uid,
		})
	}
	if this.TableConfig.GameType == 4 {
		logs.Debug("癞子模式出牌超时")
	} else if this.TableConfig.GameType == 1 {
		logs.Debug("普通模式出牌超时")
	}
}

func (this *ExtDesk) CheckStyle(style int32, out []byte, outc *GOutCard, SanDaiEr bool, SiDaiEr bool) bool {
	re := DoGenCard(out, outc, SanDaiEr, SiDaiEr)
	logs.Debug("re:", re)
	if !re {
		return re
	}
	if this.MaxChuPai != nil {
		logs.Debug("上一家出牌不等于空")
		if this.MaxChuPai.Type != outc.Type && outc.Type < CT_BOMB_FOUR {
			return false
		}
		if this.MaxChuPai.Type != outc.Type && outc.Type >= CT_BOMB_FOUR && outc.Type < this.MaxChuPai.Type {
			return false
		}
	}
	return true
}

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

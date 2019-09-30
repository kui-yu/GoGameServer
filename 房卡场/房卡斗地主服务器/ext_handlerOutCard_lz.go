package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) handleGameOutCard_lz(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("癞子出牌阶段！")
	if this.GameState != GAME_STATUS_PLAY {
		logs.Error("阶段错误", this.GameState, GAME_STATUS_PLAY)
		return
	}

	if p.ChairId != this.CurCid {
		logs.Error("不是轮到该玩家出牌", p.ChairId, this.CurCid)
		return
	}

	poutcards := GGameOutCard{}
	json.Unmarshal([]byte(d.Data), &poutcards)
	var pcards []byte
	for _, v := range poutcards.Cards {
		pcards = append(pcards, byte(v))
	}
	canout, ok := this.checkOut([]byte(pcards))
	if !ok {
		logs.Error("牌型检测错误！")
		return
	}

	vhand := p.HandCard
	nh, ok1 := VecDelMulti(vhand, pcards)
	if !ok1 {
		logs.Error("没有这些手牌!")
		return
	}
	if canout.Ishas {
		//判断是否只有一种可能
		if len(canout.Canout) == 1 {
			p.HandCard = nh
			this.MaxChuPai = &GOutCard{
				Cid:   p.ChairId,
				Max:   canout.Canout[0].Max,
				Type:  int32(canout.Canout[0].Max),
				Cards: canout.Canout[0].Cards,
			}
			p.Outed = append(p.Outed, this.MaxChuPai)
			this.RdChuPai = append(this.RdChuPai, this.MaxChuPai)
			//
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
						this.Double *= 2
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
				// LzAndBecome: canout.Canout[0].LzBecome,
				// Ptcon:       canout.Canout[0].Ptcon,
				NextCid: (this.CurCid + 1) % int32(len(this.Players)),
			}
			this.CurCid = (this.CurCid + 1) % int32(len(this.Players))

			var lzcon []byte
			lzType := 0x50
			var cd []byte
			for _, v := range canout.Canout[0].Ptcon {
				cd = append(cd, v)
			}
			for _, v := range canout.Canout[0].LzBecome {
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
		} else {
			//发送 选择信息
			p.SendNativeMsg(MSG_GAME_INFO_OUTCARD_LZ_SELECT, GGameSelectOut{
				Id:     MSG_GAME_INFO_OUTCARD_LZ_SELECT,
				Canout: canout.Canout,
			})
		}
	} else {
		p.HandCard = nh
		this.MaxChuPai = &GOutCard{
			Cid:   p.ChairId,
			Max:   canout.GoutCard.Max,
			Type:  canout.GoutCard.Type,
			Cards: canout.GoutCard.Cards,
		}
		p.Outed = append(p.Outed, this.MaxChuPai)
		this.RdChuPai = append(this.RdChuPai, this.MaxChuPai)
		//
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
				logs.Debug("可以翻倍")
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
			Id:      MSG_GAME_INFO_OUTCARD_REPLY,
			Ishas:   false,
			Cid:     p.ChairId,
			Type:    this.MaxChuPai.Type,
			Max:     this.MaxChuPai.Max,
			Double:  int32(outdouble),
			NextCid: (this.CurCid + 1) % int32(len(this.Players)),
		}
		this.CurCid = (this.CurCid + 1) % int32(len(this.Players))
		for _, v := range this.MaxChuPai.Cards {
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
				this.AddTimer(TIMER_OUTCARD, 1, this.TimerOutCard, nil)
			} else {
				this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCard, nil)
			}
		}
	}
}

func (this *ExtDesk) checkOut(cards []byte) (CanoutCards, bool) {
	canoutC := CanoutCards{}
	canOut := []CanOutType{}
	gout := GOutCard{}
	cards = Sort(cards)
	isHas, lzcon, ptcon := IsHasLz(cards, this.Lz_CardMgr.Lz_Lz)
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
	if isHas {
		canoutC.Ishas = true
		if !IsAllLz(cards, this.Lz_CardMgr.Lz_Lz) {
			var sandaidui bool = true
			var sidaidui bool = true
			for i, v := range this.TableConfig.CanSelect {
				if i == 0 {
					if v == 1 {
						sandaidui = false
					}
				}
				if i == 1 {
					if v == 1 {
						sandaidui = false
					}
				}
			}
			Lz_OutCard_Has(cards, ptcon, lzcon, &canOut, sandaidui, sidaidui)
			canoutC.Canout = canOut
		} else {
			if len(cards) >= 4 {
				co := CanOutType{
					CT:    CT_BOME_FOUR_UP,
					Cards: cards,
					Max:   GetLogicValue(lzcon[0]),
					Ptcon: ptcon,
				}
				for _, v := range lzcon {
					co.LzBecome = append(co.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: lzcon[0],
					})
				}
				canOut = append(canOut, co)
			} else {

				if this.CheckStyle(0, cards, &gout, SanDaiEr, SiDaiEr) {
					co := CanOutType{
						CT:    int(gout.Type),
						Ptcon: ptcon,
						Cards: cards,
						Max:   gout.Max,
					}
					for _, v := range lzcon {
						co.LzBecome = append(co.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: 1,
						})
					}
					canOut = append(canOut, co)
				}
			}
			canoutC.Canout = canOut
		}
		if len(canoutC.Canout) == 0 {
			return canoutC, false
		}
	} else {
		canoutC.Ishas = false
		if this.CheckStyle(0, cards, &gout, SanDaiEr, SiDaiEr) {
			canoutC.GoutCard = gout
		} else {
			logs.Error("纯普通牌类型判断错误")
			return canoutC, false
		}
	}

	if this.MaxChuPai != nil {
		if canoutC.Ishas {
			newcanout := []CanOutType{}
			//寻找符合规则的可出牌
			for i := 0; i < len(canoutC.Canout); i++ {
				if this.MaxChuPai.Type < CT_BOMB_FOUR_SOFT {
					if canoutC.Canout[i].CT >= CT_BOMB_FOUR_SOFT {
						newcanout = append(newcanout, canoutC.Canout[i])
					}
					if this.MaxChuPai.Type == int32(canoutC.Canout[i].CT) {
						if len(this.MaxChuPai.Cards) == len(canoutC.Canout[i].Cards) {
							if this.MaxChuPai.Max < canoutC.Canout[i].Max {
								newcanout = append(newcanout, canoutC.Canout[i])
							}
						}
					}
				} else {
					if this.MaxChuPai.Type != CT_BOME_FOUR_UP { //如果上次出牌 不是硬炸弹
						if this.MaxChuPai.Type < int32(canoutC.Canout[i].CT) {
							newcanout = append(newcanout, canoutC.Canout[i])
						}
						if this.MaxChuPai.Type == int32(canoutC.Canout[i].CT) {
							if this.MaxChuPai.Max < canoutC.Canout[i].Max {
								newcanout = append(newcanout, canoutC.Canout[i])
							}
						}
					} else {
						if len(this.MaxChuPai.Cards) <= len(canoutC.Canout[i].Cards) {
							newcanout = append(newcanout, canoutC.Canout[i])
						}
					}
				}
			}
			canoutC.Canout = newcanout
			if len(canoutC.Canout) != 0 {
				return canoutC, true
			} else {
				return canoutC, false
			}

		} else {
			this.CampareChuPai(&(canoutC.GoutCard))
			if &(canoutC.GoutCard) != nil {
				return canoutC, true
			} else {
				return canoutC, false
			}
		}
	} else {
		return canoutC, true
	}
}

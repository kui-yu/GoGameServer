package main

import (
	"encoding/json"
	"logs"
)

//机器人出牌
func (this *ExtDesk) robotPlayCard(d interface{}) {

	if this.GameState != GAME_STATUS_PLAY {
		return
	}

	p := d.(*ExtPlayer)

	//上一个玩家
	lastCid := p.ChairId - 1
	if lastCid < 0 {
		lastCid = 2
	}

	// 地主
	if p.ChairId == this.Players[this.Banker].ChairId {
		this.robotLandlord(p)
		return
	}
	// 农民1 农民2
	if lastCid == this.Players[this.Banker].ChairId {
		this.robotFarmer1(p)
		return
	} else {
		this.robotFarmer2(p)
		return
	}
}

//地主出牌
func (this *ExtDesk) robotLandlord(d interface{}) {

	p := d.(*ExtPlayer)

	logs.Debug("地主出牌")
	//出牌
	var playCards []byte

	//上一个玩家
	lastCid := p.ChairId - 1
	if lastCid < 0 {
		lastCid = 2
	}
	lastCardsNum := len(this.Players[lastCid].HandCard)

	//下一个玩家
	nextCid := p.ChairId + 1
	if nextCid > 2 {
		nextCid = 0
	}
	nextCardsNum := len(this.Players[nextCid].HandCard)

	//机器人是庄家
	var foeCards []byte

	if lastCardsNum > nextCardsNum {
		foeCards = this.Players[nextCid].HandCard
	} else {
		foeCards = this.Players[lastCid].HandCard
	}

	if this.MaxChuPai == nil {
		var ifLastCards []byte
		//先手最后两步牌
		playTypeValues := R_GetLastBestCalc(p.HandCard)
		if playTypeValues.num <= 2 {
			//判断剩余最后2手牌，大于等于对手
			for _, vartype := range playTypeValues.rsValues {
				wantPlays := vartype.cardValue
				for _, wantPlay := range wantPlays {
					nextTips := this.CalcTips(wantPlay, this.Players[nextCid].HandCard)
					lastTips := this.CalcTips(wantPlay, this.Players[lastCid].HandCard)
					if len(nextTips) == 0 && len(lastTips) == 0 {
						ifLastCards = wantPlay
						break
					}
				}
				if len(ifLastCards) > 0 {
					break
				}
			}
		}
		if len(ifLastCards) > 0 {
			playCards = ifLastCards
		} else {
			//先手牌
			playCards = R_OTOffensive(p.HandCard, foeCards)
		}
	} else {

		//判断是否是最后一手牌，能出直接结束
		var ifLastCards []byte
		if this.MaxChuPai != nil {

			_, ifLastCards = R_DefPosition1(this.MaxChuPai, p.HandCard)
			//出牌
			if len(ifLastCards) > 0 {
				logs.Debug("出牌", ifLastCards)
				tempCards := p.HandCard
				//剩下的手牌
				lastTempCards := ListDelListByByte(tempCards, ifLastCards)
				if len(lastTempCards) > 1 {
					playTypeValues := R_GetLastBestCalc(lastTempCards)
					logs.Debug("lastTypeValuesNum", playTypeValues.num)
					//步数刚好小等于2步
					if playTypeValues.num <= 2 {
						var wantPlayCount int
						//判断剩余最后2手牌，大于等于对手
						for _, vartype := range playTypeValues.rsValues {
							wantPlays := vartype.cardValue
							for _, wantPlay := range wantPlays {
								nextTips := this.CalcTips(wantPlay, this.Players[nextCid].HandCard)
								lastTips := this.CalcTips(wantPlay, this.Players[lastCid].HandCard)
								if len(nextTips) == 0 && len(lastTips) == 0 {
									wantPlayCount++
								}
							}
						}
						logs.Debug("wantPlayCount", wantPlayCount)
						if wantPlayCount == 0 && playTypeValues.num == 2 {
							ifLastCards = []byte{}
						}
					} else {
						ifLastCards = []byte{}
					}

				}
			}
		}
		var hasBomb bool
		if len(ifLastCards) > 0 {
			playCards = ifLastCards
		} else {
			//后手牌
			var playType byte
			//后手牌1
			playType, playCards = R_DefPosition1(this.MaxChuPai, p.HandCard)
			if playType == CT_BOMB_FOUR || playType == CT_TWOKING {
				//如果出炸弹
				if len(ifLastCards) == 0 {
					playCards = []byte{}
					hasBomb = true
				}
			}
		}
		//如果不出牌，走提示出牌
		if len(playCards) == 0 && !hasBomb {
			tips := this.CalcTips(this.MaxChuPai.Cards, p.HandCard)
			if len(tips) > 1 {
				playCards = tips[0]
			}
		}
	}

	//返回结果
	if len(playCards) > 0 {
		//检测出牌错误
		if !this.CheckOutCard(playCards, p) {

			if this.MaxChuPai != nil {
				//机器人不要
				this.HandlePass(p, &DkInMsg{
					Uid: p.Uid,
				})
			} else {
				playCards = []byte{p.HandCard[0]}
			}
		}

		//机器人要打的牌
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
		logs.Debug("机器人打的牌", playCards)
	} else {
		//机器人不要
		this.HandlePass(p, &DkInMsg{
			Uid: p.Uid,
		})
	}
}

//农民1出牌
func (this *ExtDesk) robotFarmer1(d interface{}) {

	p := d.(*ExtPlayer)

	logs.Debug("农民1出牌")
	//出牌
	var playCards []byte

	//农民2
	nextCid := p.ChairId + 1
	if nextCid > 2 {
		nextCid = 0
	}
	nextCardsNum := len(this.Players[nextCid].HandCard)

	//对手牌组
	foeCards := this.Players[this.Banker].HandCard
	foeTypeValues := R_GetLastBestCalc(foeCards)

	//判断是否是最后一手牌，能出直接结束
	var ifLastCards []byte
	if this.MaxChuPai != nil {

		_, ifLastCards = R_DefPosition1(this.MaxChuPai, p.HandCard)
		//出牌
		if len(ifLastCards) > 0 {
			logs.Debug("出牌", ifLastCards)
			tempCards := p.HandCard
			//剩下的手牌
			lastTempCards := ListDelListByByte(tempCards, ifLastCards)
			if len(lastTempCards) > 1 {
				playTypeValues := R_GetLastBestCalc(lastTempCards)
				logs.Debug("lastTypeValuesNum", playTypeValues.num)
				//步数刚好小等于2步
				if playTypeValues.num <= 2 {
					var wantPlayCount int
					//判断剩余最后2手牌，大于等于对手
					for _, vartype := range playTypeValues.rsValues {
						wantPlays := vartype.cardValue
						for _, wantPlay := range wantPlays {
							tips := this.CalcTips(wantPlay, this.Players[this.Banker].HandCard)
							if len(tips) == 0 {
								wantPlayCount++
							}
						}
					}
					logs.Debug("wantPlayCount", wantPlayCount)
					if wantPlayCount == 0 && playTypeValues.num == 2 {
						ifLastCards = []byte{}
					}
				} else {
					ifLastCards = []byte{}
				}

			}
		}
	}

	if len(ifLastCards) != 0 {
		playCards = ifLastCards
	} else {

		//正常打牌
		if this.MaxChuPai == nil {
			//先手牌
			if nextCardsNum == 1 {
				//下家剩余1张牌
				playCards = append([]byte{p.HandCard[len(p.HandCard)-1]})
			} else {
				playCards = R_OTOffensive(p.HandCard, foeCards)
			}
		} else {
			if this.Players[this.Banker].Pass {
				playCards = []byte{}
			} else {
				var hasBomb bool
				if foeTypeValues.num <= 3 && len(foeCards) < 3 {
					var playType byte
					//后手牌2
					playType, playCards = R_DefPosition2(this.MaxChuPai, p.HandCard)
					//判断王炸或炸弹，对手剩1步
					if (playType == CT_BOMB_FOUR || playType == CT_TWOKING) && foeTypeValues.num == 1 {
						tempCards := p.HandCard
						//剩下的手牌
						lastTempCards := ListDelListByByte(tempCards, playCards)
						if len(lastTempCards) > 0 {
							playTypeValues := R_GetLastBestCalc(lastTempCards)
							if playTypeValues.num < 3 {
								//对手剩1步，自己剩两步，如果能比地主大，则炸
								if playTypeValues.num == 2 {
									var wantPlayCount int
									//判断剩余最后2手牌，大于等于对手
									for _, vartype := range playTypeValues.rsValues {
										wantPlays := vartype.cardValue
										for _, wantPlay := range wantPlays {
											tips := this.CalcTips(wantPlay, this.Players[this.Banker].HandCard)
											if len(tips) == 0 {
												wantPlayCount++
											}
										}
									}
									if wantPlayCount == 0 {
										playCards = []byte{}
									}
								}
							} else {
								playCards = []byte{}
							}
						}
					}
				} else {
					var playType byte

					//后手牌1
					playType, playCards = R_DefPosition1(this.MaxChuPai, p.HandCard)
					if playType == CT_BOMB_FOUR || playType == CT_TWOKING {
						//如果出炸弹
						if len(ifLastCards) == 0 {
							playCards = []byte{}
							hasBomb = true
						}
					}
				}
				//如果出不了牌，获取提示出牌
				if len(playCards) == 0 && !hasBomb {
					tips := this.CalcTips(this.MaxChuPai.Cards, p.HandCard)
					if len(tips) > 1 {
						playCards = tips[0]
					}
				}
			}
		}
	}
	//返回结果
	if len(playCards) > 0 {
		//检测出牌错误
		if !this.CheckOutCard(playCards, p) {

			if this.MaxChuPai != nil {
				//机器人不要
				this.HandlePass(p, &DkInMsg{
					Uid: p.Uid,
				})
			} else {
				playCards = []byte{p.HandCard[0]}
			}
		}

		//机器人要打的牌
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
		logs.Debug("机器人打的牌", playCards)
	} else {
		//机器人不要
		this.HandlePass(p, &DkInMsg{
			Uid: p.Uid,
		})
	}
}

//农民2出牌
func (this *ExtDesk) robotFarmer2(d interface{}) {

	p := d.(*ExtPlayer)
	logs.Debug("农民2出牌")
	//出牌
	var playCards []byte

	//上一个玩家
	lastCid := p.ChairId - 1
	if lastCid < 0 {
		lastCid = 2
	}

	var badGuy bool
	//上个玩家不是地主
	if len(p.HandCard) > 2 {
		//上一家不是地主 && 上一家有出牌
		if !this.Players[lastCid].Pass {
			badGuy = true
		}
	}
	//对手牌组
	foeCards := this.Players[this.Banker].HandCard
	foeTypeValues := R_GetLastBestCalc(foeCards)
	//上一家牌组
	lastCards := this.Players[lastCid].HandCard
	lastTypeValues := R_GetLastBestCalc(lastCards)

	//判断是否是最后一手牌，能出直接结束
	var ifLastCards []byte
	if this.MaxChuPai != nil {

		_, ifLastCards = R_DefPosition1(this.MaxChuPai, p.HandCard)
		//出牌
		if len(ifLastCards) > 0 {
			logs.Debug("出牌", ifLastCards)
			tempCards := p.HandCard
			//剩下的手牌
			lastTempCards := ListDelListByByte(tempCards, ifLastCards)
			if len(lastTempCards) > 1 {
				playTypeValues := R_GetLastBestCalc(lastTempCards)
				logs.Debug("lastTypeValuesNum", playTypeValues.num)
				//步数刚好小等于2步
				if playTypeValues.num <= 2 {
					var wantPlayCount int
					//判断剩余最后2手牌，大于等于对手
					for _, vartype := range playTypeValues.rsValues {
						wantPlays := vartype.cardValue
						for _, wantPlay := range wantPlays {
							tips := this.CalcTips(wantPlay, this.Players[this.Banker].HandCard)
							if len(tips) == 0 {
								wantPlayCount++
							}
						}
					}
					logs.Debug("wantPlayCount", wantPlayCount)
					if wantPlayCount == 0 && playTypeValues.num == 2 {
						ifLastCards = []byte{}
					}
				} else {
					ifLastCards = []byte{}
				}

			}
		}
	}

	if len(ifLastCards) != 0 {
		playCards = ifLastCards
	} else {

		if badGuy {

			var passIng bool

			if !this.Players[lastCid].Pass && lastTypeValues.num == 1 {
				//上一家出牌 && 剩最后一步
				passIng = true
			} else if !this.Players[lastCid].Pass && !this.Players[this.Banker].Pass && foeTypeValues.num < 3 && len(foeCards) < 5 {
				//上一家出牌 && 地主也有出牌 && 地主剩两手牌
				tips := this.CalcTips(this.MaxChuPai.Cards, this.Players[this.Banker].HandCard)
				//如果 上一家出的牌比地主的牌大
				if len(tips) == 0 {
					passIng = true
				}
			}

			//判断 是否 过牌
			if passIng {
				playCards = []byte{}
			} else {
				playCards = R_Trip(this.MaxChuPai, p.HandCard, foeCards, this.Players[lastCid].HandCard)
				logs.Debug(" 顶牌2", playCards)

			}

		} else {
			//正常打牌
			var hasBomb bool
			if this.MaxChuPai == nil {
				//先手牌
				playCards = R_OTOffensive(p.HandCard, foeCards)
			} else {
				if foeTypeValues.num <= 3 && len(foeCards) < 5 {
					// logs.Debug("后手牌2")
					//后手牌
					_, playCards = R_DefPosition2(this.MaxChuPai, p.HandCard)
				} else {
					var playType byte
					//后手牌1
					playType, playCards = R_DefPosition1(this.MaxChuPai, p.HandCard)
					if playType == CT_BOMB_FOUR || playType == CT_TWOKING {
						//如果出炸弹
						if len(ifLastCards) == 0 {
							hasBomb = true
							playCards = []byte{}
						}
					}
				}
				if len(playCards) == 0 && !hasBomb {
					tips := this.CalcTips(this.MaxChuPai.Cards, p.HandCard)
					if len(tips) > 1 {
						playCards = tips[0]
					}
				}
			}
		}
	}

	//返回结果
	if len(playCards) > 0 {
		//检测出牌错误
		if !this.CheckOutCard(playCards, p) {

			if this.MaxChuPai != nil {
				//机器人不要
				this.HandlePass(p, &DkInMsg{
					Uid: p.Uid,
				})
			} else {
				playCards = []byte{p.HandCard[0]}
			}
		}

		//机器人要打的牌
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
		logs.Debug("机器人打的牌", playCards)
	} else {
		//机器人不要
		this.HandlePass(p, &DkInMsg{
			Uid: p.Uid,
		})
	}
}

//检验机器出牌
func (this *ExtDesk) CheckOutCard(cards []byte, p *ExtPlayer) bool {

	req := GOutCard{
		Cid: p.ChairId,
	}
	req.Cards = cards

	Sort(req.Cards)
	if !this.CheckStyle(req.Type, req.Cards, &req) {
		logs.Error("机器人出牌类型检测错误:", req)
		return false
	}
	return true
}

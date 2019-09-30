package main

import (
	"logs"
)

func (this *ExtDesk) CalcFirstTips(myCards []byte) [][]byte {

	var rsCards [][]byte
	Sort(myCards)

	//获取单牌
	var singleCards []byte
	typeValues := R_GetCalcTypes1(myCards)
	for _, v := range typeValues[0].rsValues {
		if v.cardType == CT_SINGLE {
			singleCards = v.cardValue[0]
			break
		}
	}

	var pairsCards [][]byte
	//获取对子
	typeValues2 := R_GetCalcTypes1(myCards)
	for _, v := range typeValues2[0].rsValues {
		if v.cardType == CT_DOUBLE {
			pairsCards = v.cardValue
			break
		}
	}

	//单牌所有数组
	var allSingleCards [][]byte
	if len(singleCards) > 0 {
		for i := len(singleCards) - 1; i >= 0; i-- {
			allSingleCards = append(allSingleCards, []byte{singleCards[i]})
		}
	}
	logs.Debug("单牌组", allSingleCards)

	//对子所有数组
	var allPairsCards [][]byte
	if len(pairsCards) > 0 {
		for i := len(pairsCards) - 1; i >= 0; i-- {
			allPairsCards = append(allPairsCards, pairsCards[i])
		}
	}
	logs.Debug("对子组", allPairsCards)

	if len(allSingleCards) > 0 && len(allPairsCards) > 0 {

		al := len(allSingleCards)
		bl := len(allPairsCards)
		cl := al + bl
		rsCards = make([][]byte, cl)

		ai := 0
		bi := 0
		ci := 0

		for ai < al && bi < bl {
			if GetLogicValue(allSingleCards[ai][0]) < GetLogicValue(allPairsCards[bi][0]) {
				rsCards[ci] = allSingleCards[ai]
				ci++
				ai++
			} else {
				rsCards[ci] = allPairsCards[bi]
				ci++
				bi++
			}
		}

		for ai < al {
			rsCards[ci] = allSingleCards[ai]
			ci++
			ai++
		}

		for bi < bl {
			rsCards[ci] = allPairsCards[bi]
			ci++
			bi++
		}

	} else if len(allSingleCards) > 0 {
		rsCards = allSingleCards
	} else if len(allPairsCards) > 0 {
		rsCards = allPairsCards
	}

	logs.Debug("提示1", rsCards)
	return rsCards
}

//智能提示
func (this *ExtDesk) CalcTips(playCards []byte, myCards []byte) [][]byte {
	//出牌类型
	req := GOutCard{}
	req.Cards = playCards
	Sort(req.Cards)
	this.CheckStyle(req.Type, req.Cards, &req)

	var rs [][]byte
	//计算
	if req.Type == CT_SINGLE_CONNECT {
		//顺子 √
		for {
			stra := CalcTipsOneStraight(req.Cards, myCards)
			if len(stra) == 0 {
				break
			}
			rs = append(rs, stra)
			req.Cards = stra
		}

	} else if req.Type == CT_DOUBLE_CONNECT {
		//连对 √
		for {
			stra := CalcTipsPairStraight(req.Cards, myCards)
			if len(stra) == 0 {
				break
			}
			rs = append(rs, stra)
			req.Cards = stra
		}

	} else if req.Type == CT_AIRCRAFT || req.Type == CT_AIRCRAFT_ONE || req.Type == CT_AIRCRAFT_TWO {

		lenPlane := len(req.Cards)
		//飞机，飞机带一，飞机带队
		maxTypeValues := R_GetBestCalc(req.Cards)
		for _, v := range maxTypeValues.rsValues {
			if v.cardType == CT_AIRCRAFT {
				req.Cards = v.cardValue[0]
				break
			}
		}
		//翅膀长度
		lenWing := lenPlane - len(req.Cards)

		//飞机
		for {
			stra := CalcTipsPlane(req.Cards, myCards)
			if len(stra) == 0 {
				break
			}
			rs = append(rs, stra)
			req.Cards = stra
		}
		if len(rs) > 0 {
			if req.Type == CT_AIRCRAFT_ONE {
				lastCards := myCards
				for _, rgCards := range rs {
					lastCards = ListDelListByByte(lastCards, rgCards)
				}

				if len(lastCards) >= lenWing {
					var singleCards []byte
					//获取单牌
					typeValues := R_GetBestCalc(lastCards)
					for _, v := range typeValues.rsValues {
						if v.cardType == CT_SINGLE {
							singleCards = v.cardValue[0]
							break
						}
					}

					var wings []byte
					if len(singleCards) < lenWing {
						//单牌不足，余牌来凑
						lastCards = ListDelListByByte(lastCards, singleCards)
						for i := 0; i < (lenWing - len(singleCards)); i++ {
							singleCards = append(singleCards, lastCards[len(lastCards)-(i+1)])
						}
					}
					//组装翅膀
					for i := len(singleCards) - 1; i >= 0; i-- {
						wings = append(wings, singleCards[i])
						if len(wings) == lenWing {
							break
						}
					}

					var rsFor [][]byte
					for _, rgCards := range rs {
						rgCards = append(rgCards, wings...)
						rsFor = append(rsFor, rgCards)
					}
					rs = rsFor
				}
			}

			if req.Type == CT_AIRCRAFT_TWO {
				lastCards := myCards
				for _, rgCards := range rs {
					lastCards = ListDelListByByte(lastCards, rgCards)
				}
				if len(lastCards) >= lenWing {
					var pairs [][]byte
					//获取对子
					typeValues := R_GetBestCalc(lastCards)
					for _, v := range typeValues.rsValues {
						if v.cardType == CT_DOUBLE {
							pairs = v.cardValue
							break
						}
					}

					var wings []byte
					if len(pairs) >= lenWing {
						for i := len(pairs) - 1; i >= 0; i-- {
							wings = append(wings, pairs[i]...)
							if len(wings) == lenWing {
								break
							}
						}

						var rsFor [][]byte
						for _, rgCards := range rs {
							rgCards = append(rgCards, wings...)
							rsFor = append(rsFor, rgCards)
						}
						rs = rsFor
					}
				}
			}
		}
	} else if req.Type == CT_THREE || req.Type == CT_THREE_LINE_TAKE_ONE || req.Type == CT_THREE_LINE_TAKE_TWO {

		//三张，三带一，三带二 √
		maxTypeValues := R_GetBestCalc(req.Cards)
		for _, v := range maxTypeValues.rsValues {
			if v.cardType == CT_THREE {
				req.Cards = v.cardValue[0]
				break
			}
		}

		//三张
		stra := CalcTipsThrees(req.Cards, myCards)
		if len(stra) != 0 {
			rs = append(rs, stra...)
			if req.Type == CT_THREE_LINE_TAKE_ONE {
				//三带一
				lastCards := myCards
				for _, rgCards := range rs {
					lastCards = ListDelListByByte(lastCards, rgCards)
				}

				var singleCards []byte
				//获取单牌
				typeValues := R_GetBestCalc(lastCards)
				for _, v := range typeValues.rsValues {
					if v.cardType == CT_SINGLE {
						singleCards = v.cardValue[0]
						break
					}
				}
				if len(singleCards) == 0 && len(lastCards) != 0 {
					singleCards = []byte{lastCards[len(lastCards)-1]}
				}

				if len(singleCards) > 0 {
					var rsFor [][]byte
					for _, rgCards := range rs {
						rgCards = append(rgCards, singleCards[len(singleCards)-1])
						rsFor = append(rsFor, rgCards)
					}
					rs = rsFor
				} else {
					rs = [][]byte{}
				}
			}
			if req.Type == CT_THREE_LINE_TAKE_TWO {
				//三带二
				lastCards := myCards
				for _, rgCards := range rs {
					lastCards = ListDelListByByte(lastCards, rgCards)
				}

				var pairs [][]byte
				//获取对子
				typeValues := R_GetBestCalc(lastCards)
				for _, v := range typeValues.rsValues {
					if v.cardType == CT_DOUBLE {
						pairs = v.cardValue
						break
					}
				}
				if len(pairs) > 0 {
					var rsFor [][]byte
					for _, rgCards := range rs {
						rgCards = append(rgCards, pairs[len(pairs)-1]...)
						rsFor = append(rsFor, rgCards)
					}
					rs = rsFor
				} else {
					for _, v := range maxTypeValues.rsValues {
						if v.cardType == CT_DOUBLE {
							req.Cards = v.cardValue[0]
							break
						}
					}
					stra := CalcTipsPairs(req.Cards, lastCards)
					if len(stra) != 0 {
						var rsFor [][]byte
						for _, rgCards := range rs {
							rgCards = append(rgCards, stra[0]...)
							rsFor = append(rsFor, rgCards)
						}
						rs = rsFor
					} else {
						rs = [][]byte{}
					}
				}
			}
		}

	} else if req.Type == CT_DOUBLE {
		//对子
		stra := CalcTipsPairs(req.Cards, myCards)
		if len(stra) != 0 {
			rs = append(rs, stra...)
		}
	} else if req.Type == CT_SINGLE {
		//单牌
		stra := CalcTipsOnes(req.Cards, myCards)
		if len(stra) != 0 {
			rs = append(rs, stra...)
		}
	}

	calcCards := R_GetCalcTypes1(myCards)
	//炸弹
	for _, typeValue := range calcCards[0].rsValues {
		if typeValue.cardType == CT_BOMB_FOUR {
			for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
				if req.Type != CT_BOMB_FOUR || GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(req.Cards[0]) {
					rs = append(rs, typeValue.cardValue[i])
				}
			}
			break
		}
	}
	//王炸 √
	for _, typeValue := range calcCards[0].rsValues {
		if typeValue.cardType == CT_TWOKING {
			for i := 0; i < len(typeValue.cardValue); i++ {
				rs = append(rs, typeValue.cardValue[i])
			}
			break
		}
	}
	logs.Debug("提示2", rs)
	return rs
}

//提示：计算单牌 √
func CalcTipsOnes(playCards []byte, myCards []byte) [][]byte {

	rs := [][]byte{}

	playSort := []byte{
		CT_SINGLE,              //单牌类型
		CT_SINGLE_CONNECT,      //单龙
		CT_DOUBLE,              //对子类型
		CT_DOUBLE_CONNECT,      //双龙
		CT_THREE_LINE_TAKE_ONE, //三带一单
		CT_THREE_LINE_TAKE_TWO, //三带一对
		CT_AIRCRAFT_ONE,        //飞机带单
		CT_AIRCRAFT_TWO,        //飞机带对
		CT_THREE,               //三张
		CT_AIRCRAFT,            //飞机
		CT_BOMB_FOUR,           //炸弹
		CT_TWOKING,             //对王类型
		CT_FOUR_LINE_TAKE_TWO,  //四带两对
		CT_FOUR_LINE_TAKE_ONE,  //四带两单
	}

	calcCards := R_GetCalcTypes1(myCards)
	//计算
	if len(calcCards) > 0 && len(calcCards[0].rsValues) > 0 {
		straCards := playCards
		for _, sortType := range playSort {
			for _, typeValue := range calcCards[0].rsValues {
				if typeValue.cardType == sortType {
					for k := len(typeValue.cardValue) - 1; k >= 0; k-- {
						ctStraCards := typeValue.cardValue[k]
						for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
							if GetLogicValue(ctStraCards[i]) > GetLogicValue(straCards[len(straCards)-1]) {
								rs = append(rs, []byte{ctStraCards[i]})
							}
						}
					}
				}
			}
		}
	}

	return rs
}

//提示：计算对子 √
func CalcTipsPairs(playCards []byte, myCards []byte) [][]byte {

	rs := [][]byte{}

	playSort := []byte{
		CT_DOUBLE,              //对子类型
		CT_DOUBLE_CONNECT,      //双龙
		CT_THREE_LINE_TAKE_ONE, //三带一单
		CT_THREE_LINE_TAKE_TWO, //三带一对
		CT_AIRCRAFT_ONE,        //飞机带单
		CT_AIRCRAFT_TWO,        //飞机带对
		CT_THREE,               //三张
		CT_AIRCRAFT,            //飞机
	}

	calcCards := R_GetCalcTypes1(myCards)
	//计算
	if len(calcCards) > 0 && len(calcCards[0].rsValues) > 0 {
		straCards := playCards
		for _, sortType := range playSort {
			for _, typeValue := range calcCards[0].rsValues {
				if typeValue.cardType == sortType {
					for k := len(typeValue.cardValue) - 1; k >= 0; k-- {
						ctStraCards := typeValue.cardValue[k]
						logs.Debug("牌型", ctStraCards)
						for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
							if GetLogicValue(ctStraCards[i]) > GetLogicValue(straCards[len(straCards)-1]) {
								start := (i + 1) - len(straCards)
								rsCards := ListGetByByte(ctStraCards, start, len(straCards))
								if len(rsCards) > 1 {
									if GetLogicValue(rsCards[0]) == GetLogicValue(rsCards[1]) {
										rs = append(rs, rsCards)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return rs
}

//提示：计算三条 √
func CalcTipsThrees(playCards []byte, myCards []byte) [][]byte {

	rs := [][]byte{}

	playSort := []byte{
		CT_THREE_LINE_TAKE_ONE, //三带一单
		CT_THREE_LINE_TAKE_TWO, //三带一对
		CT_AIRCRAFT_ONE,        //飞机带单
		CT_AIRCRAFT_TWO,        //飞机带对
		CT_THREE,               //三张
		CT_AIRCRAFT,            //飞机
	}

	calcCards := R_GetCalcTypes1(myCards)
	//计算
	if len(calcCards) > 0 && len(calcCards[0].rsValues) > 0 {
		straCards := playCards
		for _, sortType := range playSort {
			for _, typeValue := range calcCards[0].rsValues {
				if typeValue.cardType == sortType {
					for k := len(typeValue.cardValue) - 1; k >= 0; k-- {
						ctStraCards := typeValue.cardValue[k]
						logs.Debug("牌型", ctStraCards)
						for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
							if GetLogicValue(ctStraCards[i]) > GetLogicValue(straCards[len(straCards)-1]) {
								start := (i + 1) - len(straCards)
								rsCards := ListGetByByte(ctStraCards, start, len(straCards))
								logs.Debug("提示结果", start, rsCards)
								if len(rsCards) > 1 {
									if GetLogicValue(rsCards[0]) == GetLogicValue(rsCards[len(rsCards)-1]) {
										rs = append(rs, rsCards)
										logs.Debug("提示结果", rs)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return rs
}

//提示：计算顺子 √
func CalcTipsOneStraight(playCards []byte, myCards []byte) []byte {

	straCards := playCards
	ctStraCards := CalcStraight(myCards)
	if len(ctStraCards) >= len(straCards) {
		for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
			if GetLogicValue(ctStraCards[i]) > GetLogicValue(straCards[len(straCards)-1]) {
				start := (i + 1) - len(straCards)
				return ListGetByByte(ctStraCards, start, len(straCards))
			}

		}
	}
	return []byte{}
}

//提示：计算连对 √
func CalcTipsPairStraight(playCards []byte, myCards []byte) []byte {

	calcCards := R_GetCalcTypes1(myCards)

	if len(calcCards) > 0 && len(calcCards[0].rsValues) > 0 {
		straCards := playCards
		for _, typeValue := range calcCards[0].rsValues {
			if typeValue.cardType == CT_DOUBLE_CONNECT {
				ctStraCards := typeValue.cardValue[0]
				logs.Debug("出牌提示连对牌组", ctStraCards)
				if len(ctStraCards) >= len(straCards) {
					for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
						if GetLogicValue(ctStraCards[i]) > GetLogicValue(straCards[len(straCards)-1]) {
							start := (i + 1) - len(straCards)
							return ListGetByByte(ctStraCards, start, len(straCards))
						}
					}
				}
			}
		}
	}
	return []byte{}
}

//提示：计算飞机 √
func CalcTipsPlane(playCards []byte, myCards []byte) []byte {

	calcCards1 := R_GetCalcTypes1(playCards)

	var airCards []byte
	if len(calcCards1) > 0 && len(calcCards1[0].rsValues) > 0 {
		for _, typeValue := range calcCards1[0].rsValues {
			if typeValue.cardType == CT_AIRCRAFT {
				airCards = typeValue.cardValue[0]
				break
			}
		}
	}
	//计算
	calcCards2 := R_GetCalcTypes1(myCards)

	if len(calcCards2) > 0 && len(calcCards2[0].rsValues) > 0 {
		straCards := airCards
		for _, typeValue := range calcCards2[0].rsValues {
			if typeValue.cardType == CT_AIRCRAFT {
				ctStraCards := typeValue.cardValue[0]
				logs.Debug("出牌提示飞机牌组", ctStraCards)
				if len(ctStraCards) >= len(straCards) {
					for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
						if GetLogicValue(ctStraCards[i]) > GetLogicValue(straCards[len(straCards)-1]) {
							start := (i + 1) - len(straCards)
							return ListGetByByte(ctStraCards, start, len(straCards))
						}
					}
				}
			}
		}
	}
	return []byte{}
}

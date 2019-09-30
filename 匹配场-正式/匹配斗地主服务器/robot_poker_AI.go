package main

import (
	"logs"
)

//机器人叫分
func R_CallBank(cards []byte) int32 {

	//排序
	cards = Sort(cards)

	var callFen, callFen2 int32

	var maxCardNum int
	//先判断有没有大小王
	for _, card := range cards {
		if GetLogicValue(card) > 15 {
			maxCardNum += 1
		}
	}

	//没有大小王-1
	if maxCardNum > 0 {
		callFen2 += 1
	} else {
		callFen2 -= 1
	}

	maxCardNum = 0
	//后判断2的数量
	for _, card := range cards {
		if GetLogicValue(card) > 14 {
			maxCardNum += 1
		}
	}
	//大于两张2
	if maxCardNum > 1 {
		callFen2 += 1
	}

	if callFen2 > 0 {
		//步骤
		typeValues := R_GetBestCalc(cards)
		if typeValues.num < 8 {
			callFen += 1
		} else if typeValues.num < 6 {
			callFen += 2
		}
	}

	callFen += callFen2

	if callFen <= 0 {
		return -1
	}
	if callFen > 3 {
		return 3
	}
	return callFen
}

//机器人先手
func R_OTOffensive(cards []byte, foeCards []byte) []byte {

	//排序
	cards = Sort(cards)

	//步骤
	typeValues := R_GetBestCalc(cards)
	foeTypeValues := R_GetBestCalc(foeCards)

	//对手少于3步时
	var playSort []byte
	if len(foeCards) == 1 {
		//对手剩1张
		//出牌顺序
		playSort = []byte{
			CT_DOUBLE_CONNECT,      //双龙
			CT_SINGLE_CONNECT,      //单龙
			CT_AIRCRAFT_ONE,        //飞机带单
			CT_AIRCRAFT_TWO,        //飞机带对
			CT_THREE_LINE_TAKE_ONE, //三带一单
			CT_THREE_LINE_TAKE_TWO, //三带一对
			CT_THREE,               //三张
			CT_AIRCRAFT,            //飞机
			CT_BOMB_FOUR,           //炸弹
			CT_TWOKING,             //对王类型
			CT_FOUR_LINE_TAKE_ONE,  //四带两单
			CT_FOUR_LINE_TAKE_TWO,  //四带两对
			CT_DOUBLE,              //对子类型
			CT_SINGLE,              //单牌类型
		}
	} else if foeTypeValues.num < 3 && len(foeCards) < 4 {
		//对手少于3步结束，剩余牌小于4张
		//出牌顺序
		playSort = []byte{
			CT_DOUBLE_CONNECT,      //双龙
			CT_SINGLE_CONNECT,      //单龙
			CT_AIRCRAFT_TWO,        //飞机带对
			CT_AIRCRAFT_ONE,        //飞机带单
			CT_THREE_LINE_TAKE_TWO, //三带一对
			CT_THREE_LINE_TAKE_ONE, //三带一单
			CT_THREE,               //三张
			CT_AIRCRAFT,            //飞机
			CT_BOMB_FOUR,           //炸弹
			CT_TWOKING,             //对王类型
			CT_FOUR_LINE_TAKE_TWO,  //四带两对
			CT_FOUR_LINE_TAKE_ONE,  //四带两单
			CT_DOUBLE,              //对子类型
			CT_SINGLE,              //单牌类型
		}
		logs.Debug("对手", foeTypeValues)
	} else {

		//第一步
		var singleCard []byte
		var singleCardLenth int
		var pairCard [][]byte
		var pairCardLenth int

		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_SINGLE {
				singleCard = typeValue.cardValue[0]
				singleCardLenth = len(typeValue.cardValue[0])
			} else if typeValue.cardType == CT_DOUBLE {
				pairCard = typeValue.cardValue
				pairCardLenth = len(typeValue.cardValue)
			}
		}

		// logs.Debug("正常走牌")
		//对子，单牌都很多的情况
		if singleCardLenth > 1 && pairCardLenth > 1 {
			if len(pairCard) > 0 && GetLogicValue(singleCard[len(singleCard)-1]) > GetLogicValue(pairCard[len(pairCard)-1][0]) {
				if len(cards) > 3 {
					logs.Debug("单牌比较大，对牌先出")
					playSort = []byte{
						CT_DOUBLE_CONNECT,      //双龙
						CT_SINGLE_CONNECT,      //单龙
						CT_AIRCRAFT_TWO,        //飞机带对
						CT_AIRCRAFT_ONE,        //飞机带单
						CT_THREE_LINE_TAKE_TWO, //三带一对
						CT_THREE_LINE_TAKE_ONE, //三带一单
						CT_DOUBLE,              //对子类型
						CT_SINGLE,              //单牌类型
						CT_THREE,               //三张
						CT_AIRCRAFT,            //飞机
						CT_FOUR_LINE_TAKE_TWO,  //四带两对
						CT_FOUR_LINE_TAKE_ONE,  //四带两单
						CT_BOMB_FOUR,           //炸弹
						CT_TWOKING,             //对王类型
					}
				} else {
					logs.Debug("单牌比较大，单牌先出")
					playSort = []byte{
						CT_SINGLE_CONNECT,      //单龙
						CT_DOUBLE_CONNECT,      //双龙
						CT_AIRCRAFT_ONE,        //飞机带单
						CT_AIRCRAFT_TWO,        //飞机带对
						CT_THREE_LINE_TAKE_ONE, //三带一单
						CT_THREE_LINE_TAKE_TWO, //三带一对
						CT_SINGLE,              //单牌类型
						CT_DOUBLE,              //对子类型
						CT_THREE,               //三张
						CT_AIRCRAFT,            //飞机
						CT_FOUR_LINE_TAKE_ONE,  //四带两单
						CT_FOUR_LINE_TAKE_TWO,  //四带两对
						CT_BOMB_FOUR,           //炸弹
						CT_TWOKING,             //对王类型
					}
				}
			} else {
				logs.Debug("对牌比较大，单牌先出")
				playSort = []byte{
					CT_SINGLE_CONNECT,      //单龙
					CT_DOUBLE_CONNECT,      //双龙
					CT_AIRCRAFT_ONE,        //飞机带单
					CT_AIRCRAFT_TWO,        //飞机带对
					CT_THREE_LINE_TAKE_ONE, //三带一单
					CT_THREE_LINE_TAKE_TWO, //三带一对
					CT_SINGLE,              //单牌类型
					CT_DOUBLE,              //对子类型
					CT_THREE,               //三张
					CT_AIRCRAFT,            //飞机
					CT_FOUR_LINE_TAKE_ONE,  //四带两单
					CT_FOUR_LINE_TAKE_TWO,  //四带两对
					CT_BOMB_FOUR,           //炸弹
					CT_TWOKING,             //对王类型
				}
			}
		} else if pairCardLenth > 1 {
			logs.Debug("没有单牌，对牌先出")
			playSort = []byte{
				CT_DOUBLE_CONNECT,      //双龙
				CT_SINGLE_CONNECT,      //单龙
				CT_AIRCRAFT_TWO,        //飞机带对
				CT_AIRCRAFT_ONE,        //飞机带单
				CT_THREE_LINE_TAKE_TWO, //三带一对
				CT_THREE_LINE_TAKE_ONE, //三带一单
				CT_DOUBLE,              //对子类型
				CT_SINGLE,              //单牌类型
				CT_THREE,               //三张
				CT_AIRCRAFT,            //飞机
				CT_FOUR_LINE_TAKE_TWO,  //四带两对
				CT_FOUR_LINE_TAKE_ONE,  //四带两单
				CT_BOMB_FOUR,           //炸弹
				CT_TWOKING,             //对王类型
			}
		} else if singleCardLenth > 1 {
			logs.Debug("没对子，单牌先出")
			playSort = []byte{
				CT_SINGLE_CONNECT,      //单龙
				CT_DOUBLE_CONNECT,      //双龙
				CT_AIRCRAFT_ONE,        //飞机带单
				CT_AIRCRAFT_TWO,        //飞机带对
				CT_THREE_LINE_TAKE_ONE, //三带一单
				CT_THREE_LINE_TAKE_TWO, //三带一对
				CT_SINGLE,              //单牌类型
				CT_DOUBLE,              //对子类型
				CT_THREE,               //三张
				CT_AIRCRAFT,            //飞机
				CT_FOUR_LINE_TAKE_ONE,  //四带两单
				CT_FOUR_LINE_TAKE_TWO,  //四带两对
				CT_BOMB_FOUR,           //炸弹
				CT_TWOKING,             //对王类型
			}
		} else {
			if typeValues.num < 2 {
				logs.Debug("默认炸弹先出")
				playSort = []byte{
					CT_SINGLE_CONNECT,      //单龙
					CT_DOUBLE_CONNECT,      //双龙
					CT_AIRCRAFT_ONE,        //飞机带单
					CT_AIRCRAFT_TWO,        //飞机带对
					CT_THREE_LINE_TAKE_ONE, //三带一单
					CT_THREE_LINE_TAKE_TWO, //三带一对
					CT_THREE,               //三张
					CT_AIRCRAFT,            //飞机
					CT_FOUR_LINE_TAKE_ONE,  //四带两单
					CT_FOUR_LINE_TAKE_TWO,  //四带两对
					CT_BOMB_FOUR,           //炸弹
					CT_TWOKING,             //对王类型
					CT_DOUBLE,              //对子类型
					CT_SINGLE,              //单牌类型
				}
			} else {

				if len(singleCard) > 0 && len(pairCard) > 0 && GetLogicValue(singleCard[len(singleCard)-1]) > GetLogicValue(pairCard[len(pairCard)-1][0]) {
					if len(cards) > 3 {
						logs.Debug("单牌比较大，对牌先出")
						playSort = []byte{
							CT_DOUBLE_CONNECT,      //双龙
							CT_SINGLE_CONNECT,      //单龙
							CT_AIRCRAFT_TWO,        //飞机带对
							CT_AIRCRAFT_ONE,        //飞机带单
							CT_THREE_LINE_TAKE_TWO, //三带一对
							CT_THREE_LINE_TAKE_ONE, //三带一单
							CT_DOUBLE,              //对子类型
							CT_SINGLE,              //单牌类型
							CT_THREE,               //三张
							CT_AIRCRAFT,            //飞机
							CT_FOUR_LINE_TAKE_TWO,  //四带两对
							CT_FOUR_LINE_TAKE_ONE,  //四带两单
							CT_BOMB_FOUR,           //炸弹
							CT_TWOKING,             //对王类型
						}
					} else {
						logs.Debug("单牌比较大，单牌先出")
						playSort = []byte{
							CT_SINGLE_CONNECT,      //单龙
							CT_DOUBLE_CONNECT,      //双龙
							CT_AIRCRAFT_ONE,        //飞机带单
							CT_AIRCRAFT_TWO,        //飞机带对
							CT_THREE_LINE_TAKE_ONE, //三带一单
							CT_THREE_LINE_TAKE_TWO, //三带一对
							CT_SINGLE,              //单牌类型
							CT_DOUBLE,              //对子类型
							CT_THREE,               //三张
							CT_AIRCRAFT,            //飞机
							CT_FOUR_LINE_TAKE_ONE,  //四带两单
							CT_FOUR_LINE_TAKE_TWO,  //四带两对
							CT_BOMB_FOUR,           //炸弹
							CT_TWOKING,             //对王类型
						}
					}
				} else {
					logs.Debug("对牌比较大，单牌先出")
					playSort = []byte{
						CT_SINGLE_CONNECT,      //单龙
						CT_DOUBLE_CONNECT,      //双龙
						CT_AIRCRAFT_ONE,        //飞机带单
						CT_AIRCRAFT_TWO,        //飞机带对
						CT_THREE_LINE_TAKE_ONE, //三带一单
						CT_THREE_LINE_TAKE_TWO, //三带一对
						CT_SINGLE,              //单牌类型
						CT_DOUBLE,              //对子类型
						CT_THREE,               //三张
						CT_AIRCRAFT,            //飞机
						CT_FOUR_LINE_TAKE_ONE,  //四带两单
						CT_FOUR_LINE_TAKE_TWO,  //四带两对
						CT_BOMB_FOUR,           //炸弹
						CT_TWOKING,             //对王类型
					}
				}
			}
		}
	}

	//根据规则出手

	var palyCards []byte
	for _, playType := range playSort {
		// logs.Debug("牌型", playType)
		if playType == CT_AIRCRAFT_ONE {
			palyCards = []byte{}
			//飞机带单
			var airPlaneLength int
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_AIRCRAFT {
					airPlane := typeValue.cardValue[len(typeValue.cardValue)-1]
					palyCards = append(palyCards, airPlane...)
					airPlaneLength = len(airPlane)
					break
				}
			}
			if airPlaneLength > 0 {
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_SINGLE {
						if len(typeValue.cardValue[len(typeValue.cardValue)-1]) >= 2 {
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-2])
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-1])
						}
						break
					}
				}
				if len(palyCards) == (airPlaneLength + airPlaneLength/3) {
					return palyCards
				}
			}
		} else if playType == CT_AIRCRAFT_TWO {
			palyCards = []byte{}
			//飞机带队
			var airPlaneLength int
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_AIRCRAFT {
					airPlane := typeValue.cardValue[len(typeValue.cardValue)-1]
					palyCards = append(palyCards, airPlane...)
					airPlaneLength = len(airPlane)
					break
				}
			}
			if airPlaneLength > 0 {
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_DOUBLE {
						if len(typeValue.cardValue) >= 2 {
							if GetLogicValue(typeValue.cardValue[len(typeValue.cardValue)-2][0]) >= 14 {
								break
							}
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-2]...)
							if GetLogicValue(typeValue.cardValue[len(typeValue.cardValue)-1][0]) >= 14 {
								break
							}
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
						}
						break
					}
				}
				if len(palyCards) == (airPlaneLength + (airPlaneLength/3)*2) {
					return palyCards
				}
			}
		} else if playType == CT_THREE_LINE_TAKE_ONE {
			palyCards = []byte{}
			//三带一
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_THREE {
					palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
					break
				}
			}

			if len(palyCards) > 0 {
				if GetLogicValue(palyCards[0]) > 12 && typeValues.num > 2 {
					//判断 三张2 && 大于 6张，跳过
				} else {
					for _, typeValue := range typeValues.rsValues {
						if typeValue.cardType == CT_SINGLE {
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-1])
							break
						}
					}
				}
				if len(palyCards) == 4 {
					return palyCards
				}
			}
		} else if playType == CT_THREE_LINE_TAKE_TWO {
			palyCards = []byte{}
			//三带二
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_THREE {
					palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
					break
				}
			}
			if len(palyCards) > 0 {
				if GetLogicValue(palyCards[0]) > 12 && typeValues.num > 2 {
					//判断 三张2 && 大于 6张，跳过
				} else {
					for _, typeValue := range typeValues.rsValues {
						if typeValue.cardType == CT_DOUBLE {
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
							break
						}
					}
				}
				if len(palyCards) == 5 {
					return palyCards
				}
			}
		} else if playType == CT_FOUR_LINE_TAKE_ONE {
			palyCards = []byte{}
			//四带2
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_BOMB_FOUR {
					palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
					break
				}
			}
			if len(palyCards) > 0 {
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_SINGLE {
						if len(typeValue.cardValue[len(typeValue.cardValue)-1]) >= 2 {
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-2])
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-1])
						}
						break
					}
				}

				if len(palyCards) == 6 {
					return palyCards
				}
			}
		} else if playType == CT_FOUR_LINE_TAKE_TWO {
			palyCards = []byte{}
			//四带两对
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_BOMB_FOUR {
					palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
					break
				}
			}
			if len(palyCards) > 0 {
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_DOUBLE {
						if len(typeValue.cardValue) >= 2 {
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-2]...)
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
						}
						break
					}
				}

				if len(palyCards) == 8 {
					return palyCards
				}
			}
		} else if playType == CT_SINGLE {
			palyCards = []byte{}
			//单牌
			var sigleCards []byte
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_SINGLE {
					sigleCards = typeValue.cardValue[0]
					break
				}
			}

			if (len(cards) <= 2 && len(foeCards) <= 2) || len(foeCards) == 1 {
				palyCards = append(palyCards, cards[0])
			} else {
				//自己手牌
				if len(sigleCards) > 0 {
					palyCards = append(palyCards, sigleCards[len(sigleCards)-1])
				} else {
					palyCards = append(palyCards, cards[len(cards)-1])
				}
			}

			if len(palyCards) == 1 {
				return palyCards
			}
		} else {
			//三条、、、对子
			palyCards = []byte{}
			for _, typeValue := range typeValues.rsValues {
				if playType == typeValue.cardType && playType != CT_SINGLE {
					tempCards := typeValue.cardValue[len(typeValue.cardValue)-1]
					if playType == CT_DOUBLE && typeValues.num > 2 && GetLogicValue(tempCards[0]) == 15 {
						continue
					}
					if playType == CT_THREE && typeValues.num > 2 && GetLogicValue(tempCards[0]) > 12 {
						continue
					}
					return tempCards
				}
			}

		}
	}

	//过滤规则，最后还没出牌。出单牌
	if len(palyCards) == 0 {
		palyCards = append(palyCards, cards[len(cards)-1])
	}

	return palyCards
}

//机器人后手
func R_DefPosition1(maxPlay *GOutCard, cards []byte) (byte, []byte) {
	logs.Debug("机器人后手1")
	//排序
	cards = Sort(cards)

	//步骤
	typeValues := R_GetBestCalc(cards)

	//单张
	if maxPlay.Type == CT_SINGLE {

		//优化牌组里面判断
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == byte(maxPlay.Type) {
				for i := len(typeValue.cardValue[0]) - 1; i >= 0; i-- {
					if GetLogicValue(typeValue.cardValue[0][i]) > GetLogicValue(maxPlay.Max) {
						return CT_SINGLE, []byte{typeValue.cardValue[0][i]}
					}
				}
			}
		}
	}

	//对子
	if maxPlay.Type == CT_DOUBLE {
		//优化牌组里面判断
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_DOUBLE {
				for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
					if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxPlay.Max) {
						return CT_DOUBLE, typeValue.cardValue[i]
					}
				}
				break
			}
		}

		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_DOUBLE_CONNECT {
				for i := 0; i < len(typeValue.cardValue); i++ {
					for j := 0; j < len(typeValue.cardValue[i]); j++ {
						if GetLogicValue(typeValue.cardValue[i][j]) > GetLogicValue(maxPlay.Max) {
							return CT_DOUBLE, ListGetByByte(typeValue.cardValue[i], j, 2)
						}
					}
				}
				break
			}
		}

	}

	//三条
	if maxPlay.Type == CT_THREE || maxPlay.Type == CT_THREE_LINE_TAKE_ONE || maxPlay.Type == CT_THREE_LINE_TAKE_TWO {

		palyCards := []byte{}
		var maxTypeValue GTypeValue
		maxTypeValues := R_GetBestCalc(maxPlay.Cards)
		for _, v := range maxTypeValues.rsValues {
			if v.cardType == CT_THREE {
				maxTypeValue = v
			}
		}
		// logs.Debug("出的最大三条", R_GetValues(maxTypeValue.cardValue[0]))

		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_THREE {
				for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
					if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxTypeValue.cardValue[0][0]) {
						if typeValues.num <= 2 {
							palyCards = append(palyCards, typeValue.cardValue[i]...)
							break
						}
						if GetLogicValue(typeValue.cardValue[i][0]) < 14 {
							palyCards = append(palyCards, typeValue.cardValue[i]...)
							break
						}
					}
				}
				break
			}
		}
		// logs.Debug("三条", palyCards)
		if len(palyCards) > 0 {
			if maxPlay.Type == CT_THREE_LINE_TAKE_ONE {
				//三带一
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_SINGLE {
						palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-1])
						break
					}
				}

				if len(palyCards) == 4 {
					return CT_THREE_LINE_TAKE_ONE, palyCards
				}

			} else if maxPlay.Type == CT_THREE_LINE_TAKE_TWO {
				//三带二
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_DOUBLE {
						// logs.Debug("三条2", typeValue.cardValue)
						if GetLogicValue(typeValue.cardValue[len(typeValue.cardValue)-1][0]) >= 14 && typeValues.num > 2 {
							break
						}
						palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
						break
					}
				}
				if len(palyCards) == 5 {
					return CT_THREE_LINE_TAKE_TWO, palyCards
				}
			} else {
				return CT_THREE, palyCards
			}
		}
	}

	//飞机
	if maxPlay.Type == CT_AIRCRAFT || maxPlay.Type == CT_AIRCRAFT_ONE || maxPlay.Type == CT_AIRCRAFT_TWO {
		palyCards := []byte{}
		straCards := Sort(maxPlay.Cards)
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_AIRCRAFT {
				for _, ctStraCards := range typeValue.cardValue {
					if len(ctStraCards) >= len(straCards) {
						for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
							if GetCardValue(ctStraCards[i]) > GetCardValue(straCards[len(straCards)-1]) {
								start := (i + 1) - len(straCards)
								palyCards = append(palyCards, ListGetByByte(ctStraCards, start, len(straCards))...)
								break
							}
						}
					}
				}
				break
			}
		}

		if len(palyCards) > 0 {
			if maxPlay.Type == CT_AIRCRAFT_TWO {
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_DOUBLE {
						if len(typeValue.cardValue) >= 2 {
							if GetLogicValue(typeValue.cardValue[len(typeValue.cardValue)-2][0]) >= 14 {
								break
							}
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-2]...)
							if GetLogicValue(typeValue.cardValue[len(typeValue.cardValue)-1][0]) >= 14 {
								break
							}
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
						}
						break
					}
				}
				if len(palyCards) == len(straCards) {
					return CT_AIRCRAFT_TWO, palyCards
				}
			} else if maxPlay.Type == CT_AIRCRAFT_ONE {
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_SINGLE {
						if len(typeValue.cardValue[len(typeValue.cardValue)-1]) >= 2 {
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-2])
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-1])
						}
						break
					}
				}
				if len(palyCards) == len(straCards) {
					return CT_AIRCRAFT_ONE, palyCards
				}
			} else {
				return CT_AIRCRAFT, palyCards
			}
		}
	}

	//顺子
	if maxPlay.Type == CT_SINGLE_CONNECT {

		straCards := Sort(maxPlay.Cards)
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_SINGLE_CONNECT {
				ctStraCards := typeValue.cardValue[0]
				if len(ctStraCards) >= len(straCards) {
					for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
						if GetCardValue(ctStraCards[i]) > GetCardValue(straCards[len(straCards)-1]) {
							start := (i + 1) - len(straCards)
							return CT_SINGLE_CONNECT, ListGetByByte(ctStraCards, start, len(straCards))
						}
					}
				}
			}
		}
	}

	//连对
	if maxPlay.Type == CT_DOUBLE_CONNECT {

		straCards := Sort(maxPlay.Cards)
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_DOUBLE_CONNECT {
				for _, ctStraCards := range typeValue.cardValue {
					if len(ctStraCards) >= len(straCards) {
						for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
							if GetCardValue(ctStraCards[i]) > GetCardValue(straCards[len(straCards)-1]) {
								start := (i + 1) - len(straCards)
								return CT_DOUBLE_CONNECT, ListGetByByte(ctStraCards, start, len(straCards))
							}
						}
					}
				}
			}
		}
	}

	//小等于5步，可以炸
	if typeValues.num <= 3 {

		//炸弹
		if maxPlay.Type == CT_BOMB_FOUR {
			logs.Debug("有炸弹1")
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_BOMB_FOUR {
					for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
						if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxPlay.Max) {
							return CT_BOMB_FOUR, typeValue.cardValue[i]
						}
					}
				}
			}
			//王炸
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_TWOKING {
					return CT_TWOKING, typeValue.cardValue[0]
				}
			}
		} else {
			//炸弹
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_BOMB_FOUR {
					logs.Debug("有炸弹2", typeValue.cardValue)
					for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
						logs.Debug("有炸弹2-1", typeValue.cardValue[i])
						return CT_BOMB_FOUR, typeValue.cardValue[i]
					}
				}
			}
		}

		//王炸
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_TWOKING {
				return CT_TWOKING, typeValue.cardValue[0]
			}
		}
	}

	if maxPlay.Type == CT_SINGLE {
		if len(cards) < 4 {
			logs.Debug("机器人后手1,随便取张比他大的")
			//随便取张比他大的
			for i := len(cards) - 1; i >= 0; i-- {
				if GetLogicValue(cards[i]) > GetLogicValue(maxPlay.Max) {
					return CT_SINGLE, []byte{cards[i]}
				}
			}
		}
	}

	return CT_ERROR, []byte{}
}

//机器人后手2，火烧屁股了
func R_DefPosition2(maxPlay *GOutCard, cards []byte) (byte, []byte) {

	// fmt.Println("火烧屁股")
	//排序
	cards = Sort(cards)

	//步骤
	typeValues := R_GetBestCalc(cards)

	//单张
	if maxPlay.Type == CT_SINGLE {
		//火烧屁股了
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == byte(maxPlay.Type) {
				for i := 0; i < len(typeValue.cardValue[0]); i++ {
					if GetLogicValue(typeValue.cardValue[0][i]) > GetLogicValue(maxPlay.Max) {
						return CT_SINGLE, []byte{typeValue.cardValue[0][i]}
					}
				}
			}
		}

		for i := 0; i < len(cards); i++ {
			if GetLogicValue(cards[i]) > GetLogicValue(maxPlay.Max) {
				return CT_SINGLE, []byte{cards[i]}
			}
		}
	}

	//对子
	if maxPlay.Type == CT_DOUBLE {
		//火烧屁股了

		//优化牌组里面判断
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_DOUBLE {
				for i := 0; i < len(typeValue.cardValue); i++ {
					if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxPlay.Max) {
						return CT_DOUBLE, typeValue.cardValue[i]
					}
				}
			}
		}

		palyCards := []byte{}
		var maxTypeValue GTypeValue
		maxTypeValues := R_GetBestCalc(cards)
		for _, v := range maxTypeValues.rsValues {
			if v.cardType == CT_THREE {
				maxTypeValue = v
			}
		}

		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_THREE {
				for i := len(typeValue.cardValue) - 1; i >= 1; i-- {
					if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxTypeValue.cardValue[0][0]) {
						return CT_DOUBLE, ListGetByByte(typeValue.cardValue[len(typeValue.cardValue)-1], 0, 2)
					}
				}
				break
			}
		}

		if len(palyCards) == 0 {
			//从连对立面拆牌
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_DOUBLE_CONNECT {
					for i := 0; i < len(typeValue.cardValue); i++ {
						for j := 0; j < len(typeValue.cardValue[i]); j++ {
							if GetLogicValue(typeValue.cardValue[i][j]) > GetLogicValue(maxPlay.Max) {
								return CT_DOUBLE, ListGetByByte(typeValue.cardValue[i], j, 2)
							}
						}
					}
					break
				}
			}
		}

		if len(palyCards) == 0 {
			//获取算法1
			totalCalc2 := R_GetCalcTypes1(cards)
			if len(totalCalc2) > 0 {
				for _, typeValue := range totalCalc2[0].rsValues {
					if typeValue.cardType == CT_DOUBLE {
						for i := 0; i < len(typeValue.cardValue); i++ {
							if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxPlay.Max) {
								return CT_DOUBLE, typeValue.cardValue[i]
							}
						}
					}
				}
			}
		}
	}

	//三条
	if maxPlay.Type == CT_THREE || maxPlay.Type == CT_THREE_LINE_TAKE_ONE || maxPlay.Type == CT_THREE_LINE_TAKE_TWO {

		palyCards := []byte{}
		var maxTypeValue GTypeValue
		maxTypeValues := R_GetBestCalc(maxPlay.Cards)
		for _, v := range maxTypeValues.rsValues {
			if v.cardType == CT_THREE {
				maxTypeValue = v
			}
		}

		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_THREE {
				for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
					if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxTypeValue.cardValue[0][0]) {
						palyCards = append(palyCards, typeValue.cardValue[i]...)
						break
					}
				}
				break
			}
		}

		if len(palyCards) > 0 {
			if maxPlay.Type == CT_THREE_LINE_TAKE_ONE {
				//三带一
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_SINGLE {
						palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-1])
						break
					}
				}

				if len(palyCards) == 4 {
					return CT_THREE_LINE_TAKE_ONE, palyCards
				}

			} else if maxPlay.Type == CT_THREE_LINE_TAKE_TWO {
				//三带二
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_DOUBLE {
						if GetLogicValue(typeValue.cardValue[len(typeValue.cardValue)-1][0]) >= 14 {
							break
						}
						palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
						break
					}
				}

				if len(palyCards) == 5 {
					return CT_THREE_LINE_TAKE_TWO, palyCards
				}
			} else {
				return CT_THREE, palyCards
			}
		}
	}

	//飞机
	if maxPlay.Type == CT_AIRCRAFT || maxPlay.Type == CT_AIRCRAFT_ONE || maxPlay.Type == CT_AIRCRAFT_TWO {
		palyCards := []byte{}
		straCards := Sort(maxPlay.Cards)
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_AIRCRAFT {
				for _, ctStraCards := range typeValue.cardValue {
					if len(ctStraCards) >= len(straCards) {
						for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
							if GetCardValue(ctStraCards[i]) > GetCardValue(straCards[len(straCards)-1]) {
								start := (i + 1) - len(straCards)
								palyCards = append(palyCards, ListGetByByte(ctStraCards, start, len(straCards))...)
								break
							}
						}
					}
				}
				break
			}
		}

		if len(palyCards) > 0 {
			if maxPlay.Type == CT_AIRCRAFT_TWO {
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_DOUBLE {
						if len(typeValue.cardValue) >= 2 {
							if GetLogicValue(typeValue.cardValue[len(typeValue.cardValue)-2][0]) >= 14 {
								break
							}
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-2]...)
							if GetLogicValue(typeValue.cardValue[len(typeValue.cardValue)-1][0]) >= 14 {
								break
							}
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1]...)
						}
						break
					}
				}
				if len(palyCards) == len(straCards) {
					return CT_AIRCRAFT_TWO, palyCards
				}
			} else if maxPlay.Type == CT_AIRCRAFT_ONE {
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_SINGLE {
						if len(typeValue.cardValue[len(typeValue.cardValue)-1]) >= 2 {
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-2])
							palyCards = append(palyCards, typeValue.cardValue[len(typeValue.cardValue)-1][len(typeValue.cardValue[len(typeValue.cardValue)-1])-1])
						}
						break
					}
				}
				if len(palyCards) == len(straCards) {
					return CT_AIRCRAFT_ONE, palyCards
				}
			} else {
				return CT_AIRCRAFT, palyCards
			}
		}
	}

	//顺子
	if maxPlay.Type == CT_SINGLE_CONNECT {

		straCards := Sort(maxPlay.Cards)
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_SINGLE_CONNECT {
				ctStraCards := typeValue.cardValue[0]
				if len(ctStraCards) >= len(straCards) {
					for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
						if GetCardValue(ctStraCards[i]) > GetCardValue(straCards[len(straCards)-1]) {
							start := (i + 1) - len(straCards)
							return CT_SINGLE_CONNECT, ListGetByByte(ctStraCards, start, len(straCards))
						}
					}
				}
			}
		}
	}

	//连对
	if maxPlay.Type == CT_DOUBLE_CONNECT {

		straCards := Sort(maxPlay.Cards)
		for _, typeValue := range typeValues.rsValues {
			if typeValue.cardType == CT_DOUBLE_CONNECT {
				for _, ctStraCards := range typeValue.cardValue {
					if len(ctStraCards) >= len(straCards) {
						for i := len(ctStraCards) - 1; i >= len(straCards)-1; i-- {
							if GetCardValue(ctStraCards[i]) > GetCardValue(straCards[len(straCards)-1]) {
								start := (i + 1) - len(straCards)
								return CT_DOUBLE_CONNECT, ListGetByByte(ctStraCards, start, len(straCards))
							}
						}
					}
				}
			}
		}
	}

	//炸弹
	for _, typeValue := range typeValues.rsValues {
		if typeValue.cardType == CT_BOMB_FOUR {
			if maxPlay.Type == CT_BOMB_FOUR {
				for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
					if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxPlay.Max) {
						return CT_BOMB_FOUR, typeValue.cardValue[i]
					}
				}
			} else {
				return CT_BOMB_FOUR, typeValue.cardValue[len(typeValue.cardValue)-1]
			}
		}
	}
	//王炸
	for _, typeValue := range typeValues.rsValues {
		if typeValue.cardType == CT_TWOKING {
			return CT_TWOKING, typeValue.cardValue[0]
		}
	}

	return CT_ERROR, []byte{}
}

//机器人后手顶牌
func R_Trip(maxPlay *GOutCard, cards []byte, foeCards []byte, friCards []byte) []byte {

	//排序
	cards = Sort(cards)
	typeValues := R_GetBestCalc(cards)

	//地主牌
	foeCards = Sort(foeCards)
	foeTypeValues := R_GetBestCalc(foeCards)
	friTypeValues := R_GetBestCalc(friCards)

	if foeTypeValues.num < 3 && len(foeCards) < 5 {
		logs.Debug("顶牌 --> 对手剩两步")
		//对手剩两步
		if maxPlay.Type == CT_SINGLE {
			//单牌
			for _, typeValue := range typeValues.rsValues {
				//出最大的压死
				if typeValue.cardType == CT_SINGLE {
					for i := len(typeValue.cardValue[0]) - 1; i >= 0; i-- {
						if GetLogicValue(typeValue.cardValue[0][i]) >= GetLogicValue(foeCards[0]) && GetLogicValue(typeValue.cardValue[0][i]) > GetLogicValue(maxPlay.Max) {
							return []byte{typeValue.cardValue[0][i]}
						}
					}
				}
			}
		}

		if maxPlay.Type == CT_DOUBLE {
			//知道地主最大对子
			var guessCards [][]byte
			//判断地主会出什么牌
			for _, typeValue := range foeTypeValues.rsValues {
				if typeValue.cardType == CT_DOUBLE {
					for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
						if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxPlay.Max) {
							guessCards = append(guessCards, typeValue.cardValue[i])
						}
					}
					break
				}
			}
			if len(guessCards) > 0 {
				//对子
				for _, typeValue := range typeValues.rsValues {
					//出最大的压死
					if typeValue.cardType == CT_DOUBLE {
						for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
							if GetLogicValue(typeValue.cardValue[i][0]) >= GetLogicValue(guessCards[0][0]) && GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxPlay.Max) {
								return typeValue.cardValue[i]
							}
						}
					}
				}
			}

			if foeTypeValues.num < 2 && typeValues.num < 3 && friTypeValues.num < 4 {
				//炸弹
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_BOMB_FOUR {
						for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
							if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxPlay.Max) {
								return typeValue.cardValue[i]
							}
						}
						break
					}
				}

				//王炸
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_TWOKING {
						return typeValue.cardValue[0]
					}
				}
			}
		}

	} else {
		logs.Debug("顶牌 --> 计算出牌")
		//计算玩家的手牌，让玩家难受一批
		if maxPlay.Type == CT_SINGLE {

			var guessCards []byte
			//判断地主会出什么牌
			for _, typeValue := range foeTypeValues.rsValues {
				if typeValue.cardType == CT_SINGLE {
					for i := len(typeValue.cardValue[0]) - 1; i >= 0; i-- {
						if GetLogicValue(typeValue.cardValue[0][i]) > GetLogicValue(maxPlay.Max) {
							guessCards = append(guessCards, typeValue.cardValue[0][i])
						}
					}
					break
				}
			}
			//让玩家有路可走
			if len(guessCards) > 0 {
				//机器人出牌
				for _, typeValue := range typeValues.rsValues {
					if typeValue.cardType == CT_SINGLE {
						for i := len(typeValue.cardValue[0]) - 1; i >= 0; i-- {
							if GetLogicValue(typeValue.cardValue[0][i]) >= GetLogicValue(guessCards[0]) && GetLogicValue(typeValue.cardValue[0][i]) < 15 && GetLogicValue(typeValue.cardValue[0][i]) > GetLogicValue(maxPlay.Max) {
								return []byte{typeValue.cardValue[0][i]}
							}
						}
					}
				}
			}

			//溜牌
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_SINGLE {
					for i := len(typeValue.cardValue[0]) - 1; i >= 0; i-- {
						if GetLogicValue(typeValue.cardValue[0][i]) < 15 && GetLogicValue(typeValue.cardValue[0][i]) > GetLogicValue(maxPlay.Max) {
							return []byte{typeValue.cardValue[0][i]}
						}
					}
				}
			}
		}
		//对子
		if maxPlay.Type == CT_DOUBLE {
			logs.Debug("顶牌 --> 出对子")
			typeValues = R_GetCalcTypes1(cards)[0]
			for _, typeValue := range typeValues.rsValues {
				if typeValue.cardType == CT_DOUBLE {
					for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
						if GetLogicValue(typeValue.cardValue[i][0]) > GetLogicValue(maxPlay.Max) && GetLogicValue(typeValue.cardValue[i][0]) <= 10 {
							return typeValue.cardValue[i]
						}
					}
				}
			}

		}
	}

	return []byte{}
}

//计算地主最后一张牌，能否赢
func R_CalcLast(farmer1Cards []byte, landlordCards []byte) int {

	typeV1 := R_GetBestCalc(farmer1Cards)
	var count int

	if len(landlordCards) > 1 {
		//判断对子
		for _, typeValue := range typeV1.rsValues {
			if typeValue.cardType == CT_DOUBLE {
				for i := len(typeValue.cardValue) - 1; i >= 0; i-- {
					if GetLogicValue(typeValue.cardValue[i][0]) < GetLogicValue(landlordCards[0]) {
						break
					}
					count++
				}
				break
			}
		}
	}

	//单张
	for _, typeValue := range typeV1.rsValues {
		if typeValue.cardType == CT_SINGLE {
			for i := len(typeValue.cardValue[0]) - 1; i >= 0; i-- {
				if GetLogicValue(typeValue.cardValue[0][i]) < GetLogicValue(landlordCards[0]) {
					break
				}
				count++
			}
			break
		}
	}
	if count > 1 {
		return 0
	}
	return 1
}

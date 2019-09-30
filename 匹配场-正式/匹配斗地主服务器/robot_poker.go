package main

import (
	"logs"
)

//显示手牌值
func R_GetValues(cards []byte) []byte {
	var cardsValue []byte
	for _, card := range cards {
		cardsValue = append(cardsValue, GetLogicValue(card))
	}
	return cardsValue
}

type GTypeValue struct {
	cardType  byte     //类型
	cardValue [][]byte //手牌
}

type GTypeValues struct {
	num      int          //步数
	rsValues []GTypeValue //结果
}

//计算最优步数
func R_GetBestCalc(cards []byte) GTypeValues {

	var totalCalc []GTypeValues
	//获取算法1
	totalCalc = append(totalCalc, R_GetCalcTypes1(cards)...)
	//获取算法2
	totalCalc = append(totalCalc, R_GetCalcTypes2(cards)...)

	maxCalc := totalCalc[0]
	//排序所有结果
	for i := 0; i < len(totalCalc)-1; i++ {
		for j := i + 1; j < len(totalCalc); j++ {
			if totalCalc[i].num > totalCalc[j].num {
				// 1.步数排序
				maxCalc = totalCalc[j]
				totalCalc[j] = totalCalc[i]
				totalCalc[i] = maxCalc
			} else if totalCalc[i].num == totalCalc[j].num {
				// 2.步数一样，按单排排序
				var num1, num2 int
				for _, typeValue := range totalCalc[i].rsValues {
					if typeValue.cardType == CT_SINGLE {
						num1 = len(typeValue.cardValue)
					}
				}
				for _, typeValue := range totalCalc[j].rsValues {
					if typeValue.cardType == CT_SINGLE {
						num2 = len(typeValue.cardValue)
					}
				}
				if num1 > num2 {
					maxCalc = totalCalc[j]
					totalCalc[j] = totalCalc[i]
					totalCalc[i] = maxCalc
				}
			}
		}
	}
	// fmt.Println("最优步骤", totalCalc[0])
	return totalCalc[0]
}

//更新最后步数
func R_GetLastBestCalc(cards []byte) GTypeValues {
	playTypeValues := R_GetBestCalc(cards)
	//返回结果参数
	var rs GTypeValues
	//设置计算参数
	var onePlays, pairPlays, threePlays, airPlanes [][]byte
	for _, vartype := range playTypeValues.rsValues {
		if vartype.cardType == CT_THREE {
			//有三条
			threePlays = vartype.cardValue
		} else if vartype.cardType == CT_AIRCRAFT {
			//有飞机
			airPlanes = vartype.cardValue
		} else if vartype.cardType == CT_SINGLE {
			//单牌
			onePlays = vartype.cardValue
		} else if vartype.cardType == CT_DOUBLE {
			//对子
			pairPlays = vartype.cardValue
		} else {
			//其他
			rs.num += len(vartype.cardValue)
			rs.rsValues = append(rs.rsValues, vartype)
		}
	}

	//如果有三条
	if len(threePlays) > 0 {

		var ctThreeltt, ctThreelto, ctThree GTypeValue
		//三带二
		ctThreeltt.cardType = CT_THREE_LINE_TAKE_TWO
		//三带一
		ctThreelto.cardType = CT_THREE_LINE_TAKE_ONE
		//三张
		ctThree.cardType = CT_THREE

		for i := 0; i < len(threePlays); i++ {
			//判断单张对子比大小
			if (len(onePlays) > 0 && len(onePlays[0]) > 0) && len(pairPlays) > 0 {
				oneCardValue := GetLogicValue(onePlays[0][len(onePlays[0])-1])
				pairCardValue := GetLogicValue(pairPlays[len(pairPlays)-1][0])
				if oneCardValue > pairCardValue {
					//三带二
					rsCards := threePlays[i]
					rsCards = append(rsCards, pairPlays[len(pairPlays)-1]...)
					ctThreeltt.cardValue = append(ctThreeltt.cardValue, rsCards)
					logs.Debug("pairPlays", pairPlays)
					//移除最后一个元素
					pairPlays = pairPlays[:len(pairPlays)-1]
					logs.Debug("pairPlays", pairPlays)
				} else {
					//三带一
					rsCards := threePlays[i]
					rsCards = append(rsCards, onePlays[0][len(onePlays[0])-1])
					ctThreelto.cardValue = append(ctThreelto.cardValue, rsCards)
					logs.Debug("onePlays", onePlays)
					//移除最后一个元素
					onePlays[0] = onePlays[0][:len(onePlays[0])-1]
					logs.Debug("onePlays", onePlays)
				}
			} else if len(onePlays) > 0 && len(onePlays[0]) > 0 {
				//三带一
				rsCards := threePlays[i]
				rsCards = append(rsCards, onePlays[0][len(onePlays[0])-1])
				ctThreelto.cardValue = append(ctThreelto.cardValue, rsCards)
				logs.Debug("onePlays", onePlays)
				//移除最后一个元素
				onePlays[0] = onePlays[0][:len(onePlays[0])-1]
				logs.Debug("onePlays", onePlays)
			} else if len(pairPlays) > 0 {
				//三带二
				rsCards := threePlays[i]
				rsCards = append(rsCards, pairPlays[len(pairPlays)-1]...)
				ctThreeltt.cardValue = append(ctThreeltt.cardValue, rsCards)
				logs.Debug("pairPlays", pairPlays)
				//移除最后一个元素
				pairPlays = pairPlays[:len(pairPlays)-1]
				logs.Debug("pairPlays", pairPlays)
			} else {
				//三张
				ctThree.cardValue = append(ctThree.cardValue, threePlays[i])
			}
		}
		//添加结果
		if len(ctThreeltt.cardValue) > 0 {
			rs.num += len(ctThreeltt.cardValue)
			rs.rsValues = append(rs.rsValues, ctThreeltt)
		}
		if len(ctThreelto.cardValue) > 0 {
			rs.num += len(ctThreelto.cardValue)
			rs.rsValues = append(rs.rsValues, ctThreelto)
		}
		if len(ctThree.cardValue) > 0 {
			rs.num += len(ctThree.cardValue)
			rs.rsValues = append(rs.rsValues, ctThree)
		}
	}
	//如果有飞机
	if len(airPlanes) > 0 {
		var ctAirt, ctAiro, ctAir GTypeValue
		//双飞
		ctAirt.cardType = CT_AIRCRAFT_TWO
		//单飞
		ctAiro.cardType = CT_AIRCRAFT_ONE
		//飞机
		ctAir.cardType = CT_AIRCRAFT

		for i := 0; i < len(airPlanes); i++ {
			if (len(onePlays) > 0 && len(onePlays[0]) > (len(airPlanes)/3)) && len(pairPlays) > (len(airPlanes)/3) {
				oneCardValue := GetLogicValue(onePlays[0][len(onePlays[0])-1])
				pairCardValue := GetLogicValue(pairPlays[len(pairPlays)-1][0])
				//判断单张对子比大小
				if oneCardValue > pairCardValue {
					//双飞
					rsCards := airPlanes[i]
					for j := 0; j < len(airPlanes)/3; j++ {
						rsCards = append(rsCards, pairPlays[len(pairPlays)-1]...)
						//移除最后一个元素
						pairPlays = pairPlays[:len(pairPlays)-1]
					}
					ctAirt.cardValue = append(ctAirt.cardValue, rsCards)
				} else {
					//单飞
					rsCards := airPlanes[i]
					for j := 0; j < len(airPlanes)/3; j++ {
						rsCards = append(rsCards, onePlays[0][len(onePlays[0])-1])
						//移除最后一个元素
						onePlays[0] = onePlays[0][:len(onePlays[0])-1]
					}
					ctAirt.cardValue = append(ctAirt.cardValue, rsCards)
				}
			} else if len(onePlays) > 0 && len(onePlays[0]) > (len(airPlanes)/3) {
				//单飞
				rsCards := airPlanes[i]
				for j := 0; j < len(airPlanes)/3; j++ {
					rsCards = append(rsCards, onePlays[0][len(onePlays[0])-1])
					//移除最后一个元素
					onePlays[0] = onePlays[0][:len(onePlays[0])-1]
				}
				ctAirt.cardValue = append(ctAirt.cardValue, rsCards)
			} else if len(pairPlays) > (len(airPlanes) / 3) {
				//双飞
				rsCards := airPlanes[i]
				for j := 0; j < len(airPlanes)/3; j++ {
					rsCards = append(rsCards, pairPlays[len(pairPlays)-1]...)
					//移除最后一个元素
					pairPlays = pairPlays[:len(pairPlays)-1]
				}
				ctAirt.cardValue = append(ctAirt.cardValue, rsCards)
			} else {
				//飞机
				ctAir.cardValue = append(ctAir.cardValue, airPlanes[i])
			}
		}
		//添加结果
		if len(ctAirt.cardValue) > 0 {
			rs.num += len(ctAirt.cardValue)
			rs.rsValues = append(rs.rsValues, ctAirt)
		}
		if len(ctAiro.cardValue) > 0 {
			rs.num += len(ctAiro.cardValue)
			rs.rsValues = append(rs.rsValues, ctAiro)
		}
		if len(ctAir.cardValue) > 0 {
			rs.num += len(ctAir.cardValue)
			rs.rsValues = append(rs.rsValues, ctAir)
		}
	}
	if len(onePlays) > 0 && len(onePlays[0]) > 0 {
		rs.num += len(onePlays[0])
		var onePlayValues GTypeValue
		onePlayValues.cardType = CT_SINGLE
		onePlayValues.cardValue = onePlays
		rs.rsValues = append(rs.rsValues, onePlayValues)
	}
	if len(pairPlays) > 0 {
		rs.num += len(pairPlays)
		var twoPlayValues GTypeValue
		twoPlayValues.cardType = CT_DOUBLE
		twoPlayValues.cardValue = pairPlays
		rs.rsValues = append(rs.rsValues, twoPlayValues)
	}

	return rs
}

//方法1：计算手牌数量 先计算对子以上，再计算顺子
func R_GetCalcTypes1(cards []byte) []GTypeValues {

	opCards := append([]byte{}, cards...)

	//返回结果参数
	var rs GTypeValues
	var rsValues []GTypeValue
	var num int

	//自定义参数
	var kingTV GTypeValue   //王炸
	var singleTV GTypeValue //单张
	//------------------------------------

	//判断王炸
	var kingCards []byte //王炸牌
	for _, card := range opCards {
		if card == Card_King_1 {
			kingCards = append(kingCards, card)
		} else if card == Card_King_2 {
			kingCards = append(kingCards, card)
		}
	}

	if len(kingCards) > 1 {
		//有王炸，进入结果集
		kingTV.cardType = CT_TWOKING
		kingTV.cardValue = append(kingTV.cardValue, kingCards)
		rsValues = append(rsValues, kingTV)
		for _, useCards := range kingTV.cardValue {
			opCards = ListDelListByByte(opCards, useCards)
		}
	} else if len(kingCards) > 0 {
		//没有王炸，进入单张
		var singleCards []byte //单张牌
		singleCards = append(singleCards, kingCards...)
		singleTV.cardValue = append(singleTV.cardValue, singleCards)
	}

	//获取炸弹，飞机，连对，三条，对子
	opValues, opCards1 := GetKinds(opCards)
	rsValues = append(rsValues, opValues...)
	opCards = opCards1

	//判断顺子
	straTV := GetStraight(opCards)
	if straTV.cardType != CT_ERROR {
		rsValues = append(rsValues, straTV)
		for _, useCards := range straTV.cardValue {
			opCards = ListDelListByByte(opCards, useCards)
		}
	}

	if len(opCards) > 0 {
		var singleTV GTypeValue
		singleTV.cardType = CT_SINGLE
		singleTV.cardValue = append(singleTV.cardValue, opCards)
		rsValues = append(rsValues, singleTV)
	}

	//统计结果
	for _, v := range rsValues {
		if v.cardType == CT_THREE || v.cardType == CT_AIRCRAFT {
			continue
		} else if v.cardType == CT_SINGLE {
			num += len(v.cardValue[0])
		} else {
			num += len(v.cardValue)
		}
	}
	// fmt.Println("临时结果1     ", num, rsValues)
	// fmt.Println()
	// fmt.Println("======================================")
	rs.num = num
	rs.rsValues = rsValues

	var tempRs []GTypeValues
	tempRs = append(tempRs, rs)
	return tempRs
}

//方法2：计算手牌数量 先计算顺子，再计算对子以上
func R_GetCalcTypes2(cards []byte) []GTypeValues {

	opCards := append([]byte{}, cards...)
	//返回结果参数

	var rsValues []GTypeValue
	var tempRs []GTypeValues

	//自定义参数
	var kingTV GTypeValue   //王炸
	var singleTV GTypeValue //单张
	//------------------------------------

	//判断王炸
	var kingCards []byte //王炸牌
	for _, card := range opCards {
		if card == Card_King_1 {
			kingCards = append(kingCards, card)
		} else if card == Card_King_2 {
			kingCards = append(kingCards, card)
		}
	}

	if len(kingCards) > 1 {
		//有王炸，进入结果集
		kingTV.cardType = CT_TWOKING
		kingTV.cardValue = append(kingTV.cardValue, kingCards)
		rsValues = append(rsValues, kingTV)
		for _, useCards := range kingTV.cardValue {
			opCards = ListDelListByByte(opCards, useCards)
		}
	} else if len(kingCards) > 0 {
		//没有王炸，进入单张
		var singleCards []byte //单张牌
		singleCards = append(singleCards, kingCards...)
		singleTV.cardValue = append(singleTV.cardValue, singleCards)
	}

	//判断顺子
	straTV := GetStraight(opCards)
	if straTV.cardType != CT_ERROR {
		straList := straTV.cardValue[0]
		for i := 0; i <= len(straList)-5; i++ {
			for j := 0; j <= len(straList)-(5+i); j++ {
				var tempNum int
				var tempRsValues []GTypeValue
				if len(rsValues) > 0 {
					tempRsValues = append(tempRsValues, rsValues...)
				}
				var tempCards []byte
				tempCards = append([]byte{}, opCards...)
				//临时顺子数组
				tempList := ListGetByByte(straList, i, j+5)

				tempTV := GTypeValue{
					cardType:  CT_SINGLE_CONNECT,
					cardValue: [][]byte{tempList},
				}

				tempRsValues = append(tempRsValues, tempTV)
				for _, useCards := range tempTV.cardValue {
					tempCards = ListDelListByByte(tempCards, useCards)
				}

				// fmt.Println("顺子", tempTV)

				//获取炸弹，飞机，连对，三条，对子
				opValues, opCards1 := GetKinds(tempCards)
				tempRsValues = append(tempRsValues, opValues...)
				tempCards = opCards1

				//单牌
				if len(tempCards) > 0 {
					var singleTV GTypeValue
					singleTV.cardType = CT_SINGLE
					singleTV.cardValue = append(singleTV.cardValue, tempCards)
					tempRsValues = append(tempRsValues, singleTV)
				}

				//统计结果
				for _, v := range tempRsValues {
					if v.cardType == CT_THREE || v.cardType == CT_AIRCRAFT {
						continue
					} else if v.cardType == CT_SINGLE {
						tempNum += len(v.cardValue[0])
					} else {
						tempNum += len(v.cardValue)
					}
				}

				// fmt.Println("临时结果2     ", tempNum, tempRsValues)
				// fmt.Println()
				// fmt.Println("======================================")
				tempRs = append(tempRs, GTypeValues{
					num:      tempNum,
					rsValues: tempRsValues,
				})
			}
		}
	} else {
		var rs GTypeValues
		var num int
		//单牌
		if len(opCards) > 0 {
			var singleTV GTypeValue
			singleTV.cardType = CT_SINGLE
			singleTV.cardValue = append(singleTV.cardValue, opCards)
			rsValues = append(rsValues, singleTV)
		}
		//统计结果
		for _, v := range rsValues {
			if v.cardType == CT_THREE || v.cardType == CT_AIRCRAFT {
				continue
			} else if v.cardType == CT_SINGLE {
				num += len(v.cardValue[0])
			} else {
				num += len(v.cardValue)
			}
		}

		// fmt.Println("临时结果3     ", num, rsValues)
		rs.num = num
		rs.rsValues = rsValues

		tempRs = append(tempRs, rs)
	}

	return tempRs
}

//获取炸弹，飞机，连对，三条，对子
func GetKinds(opCards []byte) ([]GTypeValue, []byte) {

	var rsValues []GTypeValue

	//判断炸弹
	bombTV := CalcKinds(opCards, CT_BOMB_FOUR)
	if bombTV.cardType != CT_ERROR {
		rsValues = append(rsValues, bombTV)
		for _, useCards := range bombTV.cardValue {
			opCards = ListDelListByByte(opCards, useCards)
		}
	}

	//判断三张
	threeTV := CalcKinds(opCards, CT_THREE)
	if threeTV.cardType != CT_ERROR {
		//判断飞机
		var planeTV, delThreeTV GTypeValue
		for i := 0; i < len(threeTV.cardValue)-1; i++ {
			card1 := GetLogicValue(threeTV.cardValue[i][0])
			if card1 >= 15 {
				continue
			}
			plane := append([]byte{}, threeTV.cardValue[i]...)
			tempPlane := append([][]byte{}, threeTV.cardValue[i])
			for j := i + 1; j < len(threeTV.cardValue); j++ {
				card2 := GetLogicValue(threeTV.cardValue[j][0])
				if card1-card2 == 1 {
					plane = append(plane, threeTV.cardValue[j]...)
					tempPlane = append(tempPlane, threeTV.cardValue[j])
					card1 = card2
				}
				if card1 != card2 || len(threeTV.cardValue)-j == 1 {
					if len(plane) > 3 {
						planeTV.cardType = CT_AIRCRAFT
						delThreeTV.cardType = CT_AIRCRAFT
						planeTV.cardValue = append(planeTV.cardValue, plane)
						delThreeTV.cardValue = append(delThreeTV.cardValue, tempPlane...)
					}
					i = j - 1
					break
				}
			}
		}
		// fmt.Println("飞机1", planeTV)
		// fmt.Println("飞机2", delThreeTV)
		//先添加飞机
		if planeTV.cardType != CT_ERROR {
			rsValues = append(rsValues, planeTV)
			for _, useCards := range planeTV.cardValue {
				opCards = ListDelListByByte(opCards, useCards)
			}
		}
		//后添加三条（移除飞机）
		var tempThree [][]byte
		for _, three := range threeTV.cardValue {
			iflag := true
			for _, delThree := range delThreeTV.cardValue {
				if three[0] == delThree[0] {
					iflag = false
				}
			}
			if iflag {
				tempThree = append(tempThree, three)
			}
		}
		if len(tempThree) > 0 {
			threeTV.cardValue = tempThree
			rsValues = append(rsValues, threeTV)
			for _, useCards := range threeTV.cardValue {
				opCards = ListDelListByByte(opCards, useCards)
			}
		}
	}

	//判断对子
	twoTV := CalcKinds(opCards, CT_DOUBLE)
	if twoTV.cardType != CT_ERROR {
		//判断连对 3个对子以上
		var lianDuiTV, delTwoTV GTypeValue
		for i := 0; i < len(twoTV.cardValue)-1; i++ {
			card1 := GetLogicValue(twoTV.cardValue[i][0])
			//排除2
			if card1 == 15 {
				continue
			}
			lianDui := append([]byte{}, twoTV.cardValue[i]...)
			tempLianDui := append([][]byte{}, twoTV.cardValue[i])
			for j := i + 1; j < len(twoTV.cardValue); j++ {
				card2 := GetLogicValue(twoTV.cardValue[j][0])
				if card1-card2 == 1 {
					lianDui = append(lianDui, twoTV.cardValue[j]...)
					tempLianDui = append(tempLianDui, twoTV.cardValue[j])
					card1 = card2
				}

				if card1 != card2 || len(twoTV.cardValue)-j == 1 {
					if len(tempLianDui) > 2 {
						lianDuiTV.cardType = CT_DOUBLE_CONNECT
						delTwoTV.cardType = CT_DOUBLE_CONNECT
						lianDuiTV.cardValue = append(lianDuiTV.cardValue, lianDui)
						delTwoTV.cardValue = append(delTwoTV.cardValue, tempLianDui...)
					}
					i = j - 1
					break
				}
			}
		}

		//先添加连对
		if lianDuiTV.cardType != CT_ERROR {
			rsValues = append(rsValues, lianDuiTV)
			for _, useCards := range lianDuiTV.cardValue {
				opCards = ListDelListByByte(opCards, useCards)
			}
		}
		//后添加对子
		var tempTwo [][]byte
		for _, pair := range twoTV.cardValue {
			iflag := true
			for _, delPair := range delTwoTV.cardValue {
				if pair[0] == delPair[0] {
					iflag = false
				}
			}
			if iflag {
				tempTwo = append(tempTwo, pair)
			}
		}
		if len(tempTwo) > 0 {
			twoTV.cardValue = tempTwo
			rsValues = append(rsValues, twoTV)
			for _, useCards := range twoTV.cardValue {
				opCards = ListDelListByByte(opCards, useCards)
			}
		}
	}

	return rsValues, opCards

}

//获取顺子
func GetStraight(opCards []byte) GTypeValue {
	var tempTV GTypeValue

	rsCards := CalcStraight(opCards)
	if len(rsCards) > 0 {
		// fmt.Println("顺子", R_GetValues(rsCards))
		tempTV.cardType = CT_SINGLE_CONNECT
		tempTV.cardValue = append(tempTV.cardValue, rsCards)
	}
	return tempTV
}

//获取对子,三条，炸弹
func CalcKinds(opCards []byte, ctType byte) GTypeValue {

	var tempTV GTypeValue
	var kindCount int
	if ctType == CT_BOMB_FOUR {
		kindCount = 4
	} else if ctType == CT_THREE {
		kindCount = 3
	} else if ctType == CT_DOUBLE {
		kindCount = 2
	}
	if kindCount != 0 {
		for {
			rsCards := CalcKindsBase(opCards, kindCount)
			if len(rsCards) == 0 {
				break
			}
			opCards = ListDelListByByte(opCards, rsCards)
			tempTV.cardType = ctType
			tempTV.cardValue = append(tempTV.cardValue, rsCards)
		}
	}
	return tempTV
}

//计算顺子
func CalcStraight(cards []byte) []byte {
	if len(cards) < 5 {
		return []byte{}
	}
	var rsCards []byte
	//排序
	sortHdValues := Sort(cards)
	for i := 0; i < len(sortHdValues)-1; i++ {
		card1 := GetLogicValue(sortHdValues[i])
		if card1 >= 15 {
			continue
		}
		rsCards = append([]byte{}, sortHdValues[i])
		for j := i + 1; j < len(sortHdValues); j++ {
			card2 := GetLogicValue(sortHdValues[j])
			if card1 == card2 && len(sortHdValues)-j != 1 {
				continue
			} else if card1-card2 == 1 {
				// fmt.Println("card2", card2)
				rsCards = append(rsCards, sortHdValues[j])
				card1 = card2
			}
			// fmt.Println(rsCards)
			if len(sortHdValues)-j == 1 || card1 != card2 {
				if len(rsCards) >= 5 {
					return rsCards
				}
			}
		}
	}

	return []byte{}
}

//计算条数（对子，三条，炸弹）
func CalcKindsBase(cards []byte, kingType int) []byte {

	if len(cards) < 2 {
		return []byte{}
	}
	//排序
	sortHdValues := Sort(cards)
	//参数定义
	var rsCards []byte
	var kingCount int = 1

	for i := 0; i < len(sortHdValues)-1; i++ {
		kingCount = 1
		rsCards = append([]byte{}, sortHdValues[i])
		for j := i + 1; j < len(sortHdValues); j++ {
			if GetLogicValue(sortHdValues[i]) == GetLogicValue(sortHdValues[j]) {
				rsCards = append(rsCards, sortHdValues[j])
				kingCount++
				if kingCount == kingType {
					return rsCards
				}
			}
		}
	}

	return []byte{}
}

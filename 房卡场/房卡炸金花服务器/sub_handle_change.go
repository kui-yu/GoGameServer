package main

import (
	"encoding/json"
	"fmt"
	"logs"
	"math/rand"
	"time"
)

//换牌
func (this *ExtDesk) HandleChangeRobotCard(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("机器人请求换牌")
	if GetCostType() != 1 {
		return
	}
	data := GAMaxCard{}
	json.Unmarshal([]byte(d.Data), &data)
	if !p.Robot || p.CardType == 2 {
		info := GSChangeCard{
			Id:       MSG_GAME_INFO_CHANGE_CARD,
			CardLv:   data.CardLv,
			HandCard: data.HandCard,
			Result:   0,
		}
		p.SendNativeMsg(MSG_GAME_INFO_CHANGE_CARD, &info)
		return
	}

	//所有未弃牌真实玩家数组 pal
	var realPlayers []*ExtPlayer
	for _, v := range this.Players {
		if v.Robot {
		} else {
			if v.CardType != 2 {
				realPlayers = append(realPlayers, v)
			}
		}
	}

	//白名单
	var hierarchyWinPlayers []int32
	//黑名单
	var hierarchyLostPlayers []int32

	//1.层级概率
	//判断每个玩家是否有层级概率
	for _, v := range realPlayers {
		rand.Seed(time.Now().UnixNano())
		hierarchyRate := GetRateByHierarchyId(v.HierarchyId)

		logs.Debug("黑名单", hierarchyRate, v.HierarchyId)
		if hierarchyRate > 0 {
			//走玩家白名单层级判断
			if rand.Perm(100)[0] < int(hierarchyRate*10000/100) {
				hierarchyWinPlayers = append(hierarchyWinPlayers, v.ChairId)
			}
		} else if hierarchyRate < 0 {
			hierarchyRate = -hierarchyRate
			//走玩家黑名单层级判断
			if rand.Perm(100)[0] < int(hierarchyRate*10000/100) {
				hierarchyLostPlayers = append(hierarchyLostPlayers, v.ChairId)
			}
		}
	}

	//换牌开关
	var isChangeFalg bool = false
	// 1.有黑名单
	if len(hierarchyLostPlayers) > 0 {
		isChangeFalg = true
	}
	// 2.黑名单没管控到，走库存判断
	if !isChangeFalg {
		//2.库存概率
		intervalRate := GetRateByInterval()
		if intervalRate > 0 {
			//有库存概率，使用库存概率，否则默认
			this.RobotRate = int(intervalRate * 10000 / 100)
		}
		logs.Debug("库存管控", intervalRate, this.RobotRate, G_DbGetGameServerData.GameConfig.CurrentStock, GDeskMgr.ChangeStock)

		rand.Seed(time.Now().UnixNano())
		if rand.Perm(100)[0] < this.RobotRate || this.RobotRate == 0 {
			//robot win
			isChangeFalg = true
			logs.Debug("库存管控成功", intervalRate)
		}
	}

	// 3.桌面筹码管控
	if !isChangeFalg {
		//桌面筹码 tableTotalBet
		var tableTotalBet int64
		for i := 0; i < len(this.CoinList); i++ {
			tableTotalBet += this.CoinList[i]
		}

		tableBets := GExtRobot.RobotBets
		tableRates := GExtRobot.RobotRates

		logs.Debug("桌面管控1", tableBets, tableRates)
		if len(tableBets) > 0 {
			for i := (len(tableBets) - 1); i >= 0; i-- {
				//桌面筹码 大等于 配置筹码
				if tableTotalBet >= int64(tableBets[i]) {
					logs.Debug("桌面管控2", tableTotalBet, tableBets[i])
					//判断是否对应有筹码概率
					if len(tableBets) >= len(tableRates) {
						//判断筹码概率
						rand.Seed(time.Now().UnixNano())
						rate := rand.Perm(100)[0]
						logs.Debug("桌面管控3", rate, tableRates[i])
						if rate < tableRates[i] {
							//robot win
							isChangeFalg = true
							logs.Debug("桌面管控成功")
							break
						}
					}
				}
			}
		}
	}
	if isChangeFalg {
		cHandCard := this.ChangeCard(data.CardLv, data.HandCard, 0)

		info := GSChangeCard{
			Id:     MSG_GAME_INFO_CHANGE_CARD,
			Result: 0,
		}

		if len(cHandCard) != 3 || len(cHandCard) == 0 || cHandCard == nil {
			// info.Result = 1
			// p.SendNativeMsg(MSG_GAME_INFO_CHANGE_CARD, &info)
			this.MaxCard[0] = 1
			return
		}

		this.CardMgr.MVSourceCard = append(this.CardMgr.MVSourceCard, p.OldHandCard[:]...) //还牌
		this.CardMgr.MVSourceCard = DuplicateRemoval(this.CardMgr.MVSourceCard)

		p.HandCards, p.HandColor = SortHandCard(cHandCard)
		p.OldHandCard = cHandCard
		p.CardLv, _ = GetCardType(p.HandCards, p.HandColor)

		info.HandCard = p.OldHandCard
		info.CardLv = p.CardLv

		this.SendMaxCardRobot(p)
		p.SendNativeMsg(MSG_GAME_INFO_CHANGE_CARD, &info)

		this.MaxCard[0] = 1
	} else {
		info := GSChangeCard{
			Id:     MSG_GAME_INFO_CHANGE_CARD,
			Result: 1,
		}
		p.SendNativeMsg(MSG_GAME_INFO_CHANGE_CARD, &info)
	}
}

func (this *ExtDesk) SendMaxCardRobot(p *ExtPlayer) {
	// logs.Debug("发送换牌消息")
	info := GSMaxCard{}
	info.CardLv = p.CardLv
	info.HandCard = p.HandCards
	info.ChairId = p.ChairId
	info.IsRobot = 1

	for _, v := range this.Players {
		if v.Robot {
			v.SendNativeMsg(MSG_GAME_INFO_MAX, &info)
		}
	}
}

//
func (this *ExtDesk) ChangeCard(lv int, handcard []int, getLv int) []int {

	allCard := DuplicateRemoval(this.CardMgr.MVSourceCard)
	changeHandcard := make([]int, 0, 3)
	var ok bool
	if getLv != 0 {
		if getLv == 2 {
			ok, changeHandcard = this.GetPair(handcard, 0, allCard)
			if !ok {
				fmt.Println("获取对子失败")
			}
		}
		if getLv == 3 {
			ok, changeHandcard = this.GetShunZi(handcard, 0, allCard)
			if !ok {
				fmt.Println("获取顺子失败")
			}
		}
	} else {
		var getfunc map[int]func([]int, int, []int) (bool, []int)
		var ok bool
		getfunc = make(map[int]func([]int, int, []int) (bool, []int), 7)
		getfunc[2] = this.GetPair
		getfunc[3] = this.GetShunZi
		getfunc[4] = this.GetTongHua
		getfunc[5] = this.GetFlush
		getfunc[6] = this.GetBoom

		for i := 0; i < 7; i++ {
			if i >= lv {
				if i == 1 {
					continue
				}
				ok, changeHandcard = getfunc[i](handcard, lv, allCard)
				if ok {
					break
				}
			}
		}
	}
	//移出牌
	// this.CardMgr.MVSourceCard = ListDelList(this.CardMgr.MVSourceCard, changeHandcard)
	return changeHandcard
}

//取对子
func (this *ExtDesk) GetPair(handcard []int, cardLv int, allCard []int) (bool, []int) {
	changeHandcard := make([]int, 0, 3)
	var pairAry [][]int
	for i := 0; i < len(allCard); i++ {
		var pair []int
		pair = append(pair, allCard[i])
		for j := 0; j < len(allCard); j++ {
			if i == j {
				continue
			} else if allCard[i]%16 == allCard[j]%16 {
				pair = append(pair, allCard[j])
				pairAry = append(pairAry, pair)
				break
			}
		}
	}

	// fmt.Println("对子组：", pairAry)
	if cardLv == 0 || cardLv < 2 {
		if len(pairAry) == 0 {
			return false, nil
		}
		changeHandcard = append(changeHandcard, pairAry[rand.Perm(len(pairAry))[0]]...)
	} else if cardLv == 2 {
		for i := 0; i < len(pairAry); i++ {
			if pairAry[i][0]&0x0F > handcard[2] {
				changeHandcard = append(changeHandcard, pairAry[i]...)
				break
			}
		}
		if len(changeHandcard) == 0 {
			return false, nil
		}
	} else {
		return false, nil
	}
	for i := 0; i < len(allCard); i++ {
		if allCard[i]%16 != changeHandcard[0]%16 {
			changeHandcard = append(changeHandcard, allCard[i])
			break
		}
	}

	this.CardMgr.MVSourceCard = ListDelList(this.CardMgr.MVSourceCard, changeHandcard)
	return true, changeHandcard
}

//取顺子
func (this *ExtDesk) GetShunZi(handcard []int, cardLv int, allCard []int) (bool, []int) {
	changeHandcard := make([]int, 0, 3)
	shunziAry := make([][]int, 0)
	for i := 0; i < len(allCard); i++ {
		var shunzi []int
		shunzi = append(shunzi, allCard[i])
		for j := 0; j < len(allCard); j++ {
			if allCard[j]%16 != allCard[i]%16+1 {
				continue
			} else {
				if len(shunzi) == 3 {
					break
				}
				shunzi = append(shunzi, allCard[j])
			}
			for k := 0; k < len(allCard); k++ {
				if allCard[k]%16 != allCard[j]%16+1 {
					continue
				} else {
					shunzi = append(shunzi, allCard[k])
					if shunzi[0]+1 == shunzi[1] && shunzi[1]+1 == shunzi[2] {
						// logs.Debug("顺金跳出")
						break
					} else {
						shunziAry = append(shunziAry, shunzi)
						break
					}
				}
			}
		}
	}
	// fmt.Println("顺子组：", shunziAry)
	if cardLv == 0 || cardLv < 3 {
		if len(shunziAry) == 0 {
			return false, nil
		}
		rand.Seed(time.Now().UnixNano())
		changeHandcard = shunziAry[rand.Perm(len(shunziAry))[0]]
	} else if cardLv == 3 {
		for i := 0; i < len(shunziAry); i++ {
			if shunziAry[i][2]&0x0F > handcard[2] {
				changeHandcard = shunziAry[i]
			}
		}
		if len(changeHandcard) == 0 {
			return false, nil
		}
	} else {
		return false, nil
	}

	this.CardMgr.MVSourceCard = ListDelList(this.CardMgr.MVSourceCard, changeHandcard)
	return true, changeHandcard
}

//取金花
func (this *ExtDesk) GetTongHua(handcard []int, cardLv int, allCard []int) (bool, []int) {
	changeHandcard := make([]int, 0, 3)
	tonghuaAry := make([][]int, 0)

	for i := 0; i < len(allCard); i++ {
		tonghua := make([]int, 0)
		tonghua = append(tonghua, allCard[i])
		for j := 0; j < len(allCard); j++ {
			if len(tonghua) == 3 {
				break
			}
			if j == i {
				continue
			}
			if allCard[i]&0xF0 == allCard[j]&0xF0 &&
				allCard[i]&0x0F < allCard[j]&0x0F &&
				allCard[i]&0x0F != allCard[j]&0x0F-1 &&
				allCard[i]&0x0F != allCard[j]&0x0F+1 {
				tonghua = append(tonghua, allCard[j])
			} else {
				continue
			}
			for k := 0; k < len(allCard); k++ {
				if len(tonghua) == 3 {
					break
				}
				if j == k {
					continue
				}
				if allCard[j]&0xF0 == allCard[k]&0xF0 && allCard[j]&0x0F < allCard[k]&0x0F {
					if allCard[i]&0x0F != allCard[j]&0x0F-1 || allCard[i]&0x0F != allCard[j]&0x0F+1 {
						tonghua = append(tonghua, allCard[k])
						tonghuaAry = append(tonghuaAry, tonghua)
						break
					} else {
						continue
					}
				} else {
					continue
				}
			}
		}
	}

	// fmt.Println("同花数组：", tonghuaAry)
	if cardLv == 0 || cardLv < 4 {
		if len(tonghuaAry) == 0 {
			return false, nil
		}
		changeHandcard = append(changeHandcard, tonghuaAry[0]...)
	} else if cardLv == 4 {
		for i := 0; i < len(tonghuaAry); i++ {
			if tonghuaAry[i][2]&0x0F > handcard[2] {
				changeHandcard = append(changeHandcard, tonghuaAry[i]...)
				break
			}
		}
		if len(changeHandcard) == 0 {
			return false, nil
		}
	} else {
		return false, nil
	}
	this.CardMgr.MVSourceCard = ListDelList(this.CardMgr.MVSourceCard, changeHandcard)
	return true, changeHandcard
}

//取同花顺
func (this *ExtDesk) GetFlush(handcard []int, cardLv int, allCard []int) (bool, []int) {
	changeHandcard := make([]int, 0, 3)
	shunziAry := make([][]int, 0)
	for i := 0; i < len(allCard); i++ {
		var shunzi []int
		shunzi = append(shunzi, allCard[i])
		for j := 0; j < len(allCard); j++ {
			if allCard[j]%16 != allCard[i]%16+1 {
				continue
			} else {
				if len(shunzi) == 3 {
					break
				}
				shunzi = append(shunzi, allCard[j])
			}
			for k := 0; k < len(allCard); k++ {
				if allCard[k]%16 != allCard[j]%16+1 {
					continue
				} else {
					shunzi = append(shunzi, allCard[k])
					if shunzi[0]+1 == shunzi[1] && shunzi[1]+1 == shunzi[2] {
						shunziAry = append(shunziAry, shunzi)
						break
					} else {
						break
					}
				}
			}
		}
	}

	// fmt.Println("顺金组：", shunziAry)
	if cardLv == 0 || cardLv < 5 {
		if len(shunziAry) == 0 {
			return false, nil
		}
		changeHandcard = append(changeHandcard, shunziAry[0]...)
	} else if cardLv == 5 {
		for i := 0; i < len(shunziAry); i++ {
			if shunziAry[i][2]&0x0F > handcard[2] {
				changeHandcard = append(changeHandcard, shunziAry[i]...)
				break
			}
		}
		if len(changeHandcard) == 0 {
			return false, nil
		}
	} else {
		return false, nil
	}

	this.CardMgr.MVSourceCard = ListDelList(this.CardMgr.MVSourceCard, changeHandcard)
	return true, changeHandcard
}

//取炸弹
func (this *ExtDesk) GetBoom(handcard []int, cardLv int, allCard []int) (bool, []int) {
	changeHandcard := make([]int, 0, 3)
	boomAry := make([][]int, 0)
	for i := 0; i < len(allCard); i++ {
		var boom []int
		boom = append(boom, allCard[i])
		for j := 0; j < len(allCard); j++ {
			if len(boom) == 3 {
				break
			}
			if i == j || allCard[i]%16 != allCard[j]%16 {
				continue
			} else {
				boom = append(boom, allCard[j])
			}
			for k := 0; k < len(allCard); k++ {
				if len(boom) == 3 {
					break
				}
				if i == k || j == k || allCard[k]%16 != allCard[j]%16 {
					continue
				} else {
					boom = append(boom, allCard[k])
					boomAry = append(boomAry, boom)
					break
				}
			}
		}
	}

	// fmt.Println("炸弹组：", boomAry)
	if cardLv == 0 || cardLv < 6 {
		if len(boomAry) == 0 {
			return false, nil
		}
		changeHandcard = append(changeHandcard, boomAry[0]...)
	} else if cardLv == 6 {
		for i := 0; i < len(boomAry); i++ {
			if boomAry[i][0]&0x0F > handcard[0] {
				changeHandcard = append(changeHandcard, boomAry[i]...)
				break
			}
		}
		if len(changeHandcard) == 0 {
			return false, nil
		}
	} else {
		return false, nil
	}

	this.CardMgr.MVSourceCard = ListDelList(this.CardMgr.MVSourceCard, changeHandcard)
	return true, changeHandcard
}

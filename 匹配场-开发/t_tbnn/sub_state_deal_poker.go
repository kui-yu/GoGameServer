package main

import (
	"logs"
	"math/rand"
	"time"
)

func (this *ExtDesk) GameStateDeal() {

	this.CardMgr.Shuffle()

	var cardTypes []GCardType //手牌数组
	for i := 0; i < len(this.Players); i++ {
		handCard := this.CardMgr.SendHandCard(5)
		niuPoint, niuCards := GetResult(handCard)
		cardType := GCardType{
			HandCard: handCard,
			NiuPoint: niuPoint,
			NiuCards: niuCards,
		}
		cardTypes = append(cardTypes, cardType)
	}
	if GetCostType() == 1 { //如果不是体验场再进行排序
		cardTypes = this.ControlPoker(cardTypes)
	}
	//排序

	winChairs := this.ControlResult()

	for i := 0; i < len(winChairs); i++ {
		this.Players[winChairs[i]].HandCard = cardTypes[i].HandCard
		this.Players[winChairs[i]].NiuPoint = cardTypes[i].NiuPoint
		this.Players[winChairs[i]].NiuCards = cardTypes[i].NiuCards
		this.Players[winChairs[i]].NiuMultiple = GetNiuMultiple(cardTypes[i].NiuPoint)

		handcard := GHandNiuReply{
			Id:       MSG_GAME_INFO_DEAL_REPLY,
			ChairId:  this.Players[winChairs[i]].ChairId,
			NiuPoint: this.Players[winChairs[i]].NiuPoint,
			NiuCards: this.Players[winChairs[i]].HandCard,
		}
		this.Players[winChairs[i]].SendNativeMsg(MSG_GAME_INFO_DEAL_REPLY, &handcard)
	}

	//进入发牌动画
	this.BroadStageTime(TIME_STAGE_START_NUM)
	//进入倒计时
	this.runTimer(TIME_STAGE_START_NUM, this.GameStartEnd)
}

//发牌动画-结束
func (this *ExtDesk) GameStartEnd(d interface{}) {
	// logs.Debug("发牌动画")
	//进入玩牌
	this.nextStage(STAGE_PLAY)
}

func (this *ExtDesk) ControlPoker(cardTypes []GCardType) []GCardType {
	// logs.Debug("排序前", cardTypes)
	for i := 0; i < len(cardTypes); i++ {
		maxCardTypes := cardTypes[i]
		for j := i + 1; j < len(cardTypes); j++ {
			rs := SoloResult(cardTypes[i], cardTypes[j])
			if rs == 2 {
				maxCardTypes = cardTypes[j]
				cardTypes[j] = cardTypes[i]
				cardTypes[i] = maxCardTypes
			}
		}
	}
	// logs.Debug("排序后", cardTypes)
	return cardTypes
}

func (this *ExtDesk) ControlResult() []int32 {

	//玩家随机
	var winChairs []int32

	//robot
	var robotPlayers []int32
	//真实玩家
	var realPlayers []int32

	for _, v := range this.Players {
		if v.Robot {
			robotPlayers = append(robotPlayers, v.ChairId)
		} else {
			realPlayers = append(realPlayers, v.ChairId)
		}
	}

	if len(robotPlayers) > 0 {
		//白名单
		var hierarchyWinPlayers []int32
		//黑名单
		var hierarchyLostPlayers []int32

		//1.层级概率
		//判断每个玩家是否有层级概率
		for _, v := range this.Players {
			if !v.Robot {
				rand.Seed(time.Now().UnixNano())
				hierarchyRate := GetRateByHierarchyId(v.HierarchyId)
				logs.Debug("黑名单", hierarchyRate, v.HierarchyId)
				if hierarchyRate > 0 {
					logs.Debug("进入白名单", hierarchyRate, v.HierarchyId)
					//走玩家白名单层级判断
					if rand.Perm(100)[0] < int(hierarchyRate*10000/100) {
						hierarchyWinPlayers = append(hierarchyWinPlayers, v.ChairId)
					}
				} else if hierarchyRate < 0 {
					logs.Debug("进入黑名单", hierarchyRate, v.HierarchyId)
					hierarchyRate = -hierarchyRate
					//走玩家黑名单层级判断
					if rand.Perm(100)[0] < int(hierarchyRate*10000/100) {
						hierarchyLostPlayers = append(hierarchyLostPlayers, v.ChairId)
					}
				}
			}
		}

		//执行判断
		if len(hierarchyWinPlayers) > 0 {
			//1.判断白名单，有结束
			if len(hierarchyWinPlayers) > 1 {
				hierarchyWinPlayers = ListShuffle(hierarchyWinPlayers)
			}
			winChairs = append(winChairs, hierarchyWinPlayers...)
		} else if len(hierarchyLostPlayers) == len(realPlayers) {
			//黑名单玩家 都是真实玩家，结束
			//robot win
			var robotList []int32
			robotList = append(robotList, robotPlayers...)
			robotList = ListShuffle(robotList)

			winChairs = append(winChairs, robotList[0])
		} else {
			//库存判断
			intervalRate := GetRateByInterval()
			logs.Debug("库存管控", intervalRate, this.RobotRate, G_DbGetGameServerData.GameConfig.CurrentStock, GDeskMgr.ChangeStock)
			if intervalRate > 0 {
				//有库存概率，使用库存概率，否则默认
				this.RobotRate = int(intervalRate * 10000 / 100)
			}

			rand.Seed(time.Now().UnixNano())
			if rand.Intn(100) < this.RobotRate {
				logs.Debug("进入库存")
				//robot win
				var robotList []int32
				robotList = append(robotList, robotPlayers...)
				robotList = ListShuffle(robotList)
				winChairs = append(winChairs, robotList[0])
			} else {
				// 排除黑名单玩家
				var lastThisPlayers []*ExtPlayer
				for _, v := range this.Players {
					flag := true
					for j := 0; j < len(hierarchyLostPlayers); j++ {
						if v.ChairId == hierarchyLostPlayers[j] {
							flag = false
							break
						}
					}
					if flag {
						lastThisPlayers = append(lastThisPlayers, v)
					}
				}
				// player win
				for _, v := range lastThisPlayers {
					if !v.Robot {
						winChairs = append(winChairs, v.ChairId)
					}
				}
			}
			//随机Win玩家
			winChairs = ListShuffle(winChairs)
		}

		//剩余玩家
		if len(winChairs) != len(this.Players) {
			var lastPlayers []int32
			for i := 0; i < len(this.Players); i++ {
				flag := true
				for j := 0; j < len(winChairs); j++ {
					if this.Players[i].ChairId == winChairs[j] {
						flag = false
						break
					}
				}
				if flag {
					lastPlayers = append(lastPlayers, this.Players[i].ChairId)
				}
			}
			//随机LAST玩家
			lastPlayers = ListShuffle(lastPlayers)
			winChairs = append(winChairs, lastPlayers...)
		}
	} else {
		for _, v := range this.Players {
			winChairs = append(winChairs, v.ChairId)
		}
		winChairs = ListShuffle(winChairs)
	}
	// logs.Debug("玩家排序", winChairs)
	return winChairs
}

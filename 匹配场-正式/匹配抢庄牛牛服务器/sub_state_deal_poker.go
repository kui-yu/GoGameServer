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
			NiuPoint: niuPoint,
			NiuCards: niuCards,
			HandCard: handCard,
		}
		cardTypes = append(cardTypes, cardType)
	}
	//排序
	if GetCostType() == 1 { //如果不是体验场再进行大小排序
		cardTypes = this.ControlPoker(cardTypes)
	}
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

	this.BroadStageTime(TIME_STAGE_START_NUM)
	//进入倒计时
	this.runTimer(TIME_STAGE_START_NUM, this.GameStateDealEnd)
}

//发牌阶段-结束
func (this *ExtDesk) GameStateDealEnd(d interface{}) {
	//进入玩牌
	this.nextStage(STAGE_PLAY)
}

func (this *ExtDesk) ControlPoker(cardTypes []GCardType) []GCardType {
	logs.Debug("排序前", cardTypes)
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
	logs.Debug("排序后", cardTypes)
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

	//robot win
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
				if hierarchyRate > 0 {
					logs.Debug("进入白名单", hierarchyRate)
					//走玩家白名单层级判断
					if rand.Perm(100)[0] < int(hierarchyRate*10000/100) {
						hierarchyWinPlayers = append(hierarchyWinPlayers, v.ChairId)
					}
				} else if hierarchyRate < 0 {
					hierarchyRate = -hierarchyRate
					//走玩家黑名单层级判断
					if rand.Perm(100)[0] < int(hierarchyRate*10000/100) {
						hierarchyLostPlayers = append(hierarchyLostPlayers, v.ChairId)
						logs.Debug("进入黑名单", hierarchyRate)
					}
				}
			}
		}

		//执行判断
		if len(hierarchyWinPlayers) == len(realPlayers) {
			//1.判断白名单，都是真实玩家，结束
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

			winChairs = append(winChairs, robotList...)
		} else {
			//1.判断白名单，都是真实玩家，结束
			if len(hierarchyWinPlayers) > 1 {
				hierarchyWinPlayers = ListShuffle(hierarchyWinPlayers)
			}
			winChairs = append(winChairs, hierarchyWinPlayers...)

			//库存判断
			intervalRate := GetRateByInterval()
			if intervalRate > 0 {
				//有库存概率，使用库存概率，否则默认
				this.RobotRate = int(intervalRate * 10000 / 100)
			}

			rand.Seed(time.Now().UnixNano())
			if rand.Intn(100) < this.RobotRate {
				//robot win
				RobotBanker := false
				for _, v := range this.Players {
					if v.Robot {
						if v.ChairId == this.Banker {
							RobotBanker = true
							break
						}
					}
				}

				var robotList []int32
				robotList = append(robotList, robotPlayers...)
				robotList = ListShuffle(robotList)

				for _, v := range robotList {
					winChairs = append(winChairs, v)
					if RobotBanker {
						break
					}
				}
			} else {

				//player allwin
				if len(winChairs) == 0 {
					var robotList []int32
					robotList = append(robotList, robotPlayers...)
					robotList = ListShuffle(robotList)

					rand.Seed(time.Now().UnixNano())
					randNum := rand.Intn(100)
					//player allWin
					if randNum < GExtRobot.PlayerRate && GExtRobot.PlayerRate != 0 {
					} else {
						winChairs = append(winChairs, robotList[0])
					}
				}

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
				// 排除白名单玩家
				var lastAllPlayers []*ExtPlayer
				for _, v := range lastThisPlayers {
					flag := true
					for j := 0; j < len(hierarchyWinPlayers); j++ {
						if v.ChairId == hierarchyWinPlayers[j] {
							flag = false
							break
						}
					}
					if flag {
						lastAllPlayers = append(lastAllPlayers, v)
					}
				}

				var tempChairs []int32
				//player win
				for _, v := range lastAllPlayers {
					if !v.Robot {
						tempChairs = append(tempChairs, v.ChairId)
					}
				}
				tempChairs = ListShuffle(tempChairs)
				winChairs = append(winChairs, tempChairs...)
			}
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
			lastPlayers = ListShuffle(lastPlayers)
			winChairs = append(winChairs, lastPlayers...)
		}
	} else {
		for _, v := range this.Players {
			winChairs = append(winChairs, v.ChairId)
		}
		winChairs = ListShuffle(winChairs)
	}
	logs.Debug("玩家排序", winChairs)
	return winChairs
}

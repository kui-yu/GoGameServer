package main

import (
	"logs"
	"math/rand"
	"time"
)

func (this *ExtDesk) GameStateDeal() {
	//摇双骰
	diceList := []int{1, 2, 3, 4, 5, 6}
	temp1 := ListShuffle(diceList)
	temp2 := ListShuffle(diceList)

	var cardInfo GSCardInfo
	ListAdd(&cardInfo.Dices, temp1[0])
	ListAdd(&cardInfo.Dices, temp2[0])

	var playCards [][]int
	//发牌
	for i := 0; i < len(this.Players); i++ {
		if len(this.CardMgr.MVSourceCard) < 2 {
			//没牌，直接结束当前桌子信息
			this.SysTableEnd()
			logs.Error("当前局游戏信息错误，手牌不够。", "当前局数", this.Round)
			return
		} else {
			cards := this.CardMgr.SendHandCard(2)
			playCards = append(playCards, cards)
			// logs.Debug("发牌", cards)
			this.AddPutInfos(cards)
		}
	}
	winChairs := make([]int, 0)
	for _, v := range this.Players {
		winChairs = append(winChairs, int(v.ChairId))
	}
	if GetCostType() == 1 { //如果不是体验场再进入层级、库存概率判断
		playCards = this.ControlPoker(playCards)
		winChairs = this.ControlResult()
	}

	for i := 0; i < len(winChairs); i++ {
		chairId := int32(winChairs[i])
		//*控制赢家
		this.Players[chairId].HandCards = playCards[i]
		//通知用户
		cardInfo.HandCards = this.Players[chairId].HandCards
		cardInfo.Id = MSG_GAME_INFO_CARD_INFO_REPLY
		this.Players[chairId].SendNativeMsg(MSG_GAME_INFO_CARD_INFO_REPLY, &cardInfo)
	}

	//发牌阶段
	this.BroadStageTime(STAGE_DEAL_TIME)
	this.runTimer(STAGE_DEAL_TIME, this.GameStateDealEnd)
}

//阶段-发牌
func (this *ExtDesk) GameStateDealEnd(d interface{}) {
	//发牌时间到，进入结算阶段
	this.nextStage(STAGE_SETTLE)
}

func (this *ExtDesk) ControlPoker(playCards [][]int) [][]int {
	// logs.Debug("排序前", playCards)
	//牌组比大小
	for i := 0; i < len(playCards)-1; i++ {
		maxCards := playCards[i]
		for j := i + 1; j < len(playCards); j++ {
			rs := compareCards(playCards[i], playCards[j])
			if rs == 2 {
				maxCards = playCards[j]
				playCards[j] = playCards[i]
				playCards[i] = maxCards
			}
		}
	}
	// logs.Debug("排序后", playCards)
	return playCards
}

//game-control
func (this *ExtDesk) ControlResult() []int {

	var winChairs []int
	//robot
	var robotPlayers []int
	for _, v := range this.Players {
		if v.Robot {
			robotPlayers = append(robotPlayers, int(v.ChairId))
		}
	}

	logs.Debug("机器数组", robotPlayers)

	// 有机器人
	if len(robotPlayers) > 0 {
		var hierarchyWinPlayers []int
		var hierarchyLostPlayers []int

		//1.层级概率
		//判断每个玩家是否有层级概率
		for _, v := range this.Players {
			if !v.Robot {
				rand.Seed(time.Now().UnixNano())
				hierarchyRate := GetRateByHierarchyId(v.HierarchyId)
				logs.Debug("黑白名单", hierarchyRate, v.HierarchyId)
				if hierarchyRate > 0 {
					//走玩家白名单层级判断
					if rand.Perm(100)[0] < int(hierarchyRate*10000/100) {
						hierarchyWinPlayers = append(hierarchyWinPlayers, int(v.ChairId))
					}
				} else if hierarchyRate < 0 {
					hierarchyRate = -hierarchyRate
					//走玩家黑名单层级判断
					if rand.Perm(100)[0] < int(hierarchyRate*10000/100) {
						//*取得黑名单的玩家
						hierarchyLostPlayers = append(hierarchyLostPlayers, int(v.ChairId))
						logs.Debug("进入黑名单")
					}
				}
			}
		}
		if len(hierarchyWinPlayers) > 0 {
			if len(hierarchyWinPlayers) > 1 {
				hierarchyWinPlayers = ListShuffle(hierarchyWinPlayers)
			}
			winChairs = append(winChairs, hierarchyWinPlayers...)
		}

		//去黑名单
		var quBlackPlayers []*ExtPlayer
		for _, v := range this.Players {
			flag := true
			for j := 0; j < len(hierarchyLostPlayers); j++ {
				if int(v.ChairId) == hierarchyLostPlayers[j] {
					flag = false
					break
				}
			}
			if flag {
				//*取得黑名单外的玩家
				quBlackPlayers = append(quBlackPlayers, v)
			}
		}

		//剩余玩家归组
		var lastThisPlayers []*ExtPlayer //正常玩家 不在winChairs组里的
		for _, v := range quBlackPlayers {
			flag := true
			for j := 0; j < len(winChairs); j++ {
				if int(v.ChairId) == winChairs[j] {
					flag = false
					break
				}
			}
			if flag {
				//*取得黑名单白名单外的玩家
				lastThisPlayers = append(lastThisPlayers, v)
			}
		}

		//2.库存概率
		intervalRate := GetRateByInterval()
		if intervalRate > 0 {
			//有库存概率，使用库存概率，否则默认
			this.RobotRate = int(intervalRate * 10000 / 100)
		}
		logs.Debug("库存概率", intervalRate, this.RobotRate)
		rand.Seed(time.Now().UnixNano())
		if rand.Perm(100)[0] < this.RobotRate {

			logs.Debug("进入库存概率")
			//robot win   robotplayers->机器人数组
			if len(robotPlayers) == 1 {
				winChairs = append(winChairs, robotPlayers...)
			} else if this.Players[this.Banker].Robot {
				//倍数大小数组玩家
				var multPlayers []int
				//机器人坐庄
				for _, v := range lastThisPlayers {
					if v.ChairId != this.Banker {
						//玩家赔率
						multPlayers = append(multPlayers, int(v.ChairId))
					}
				}
				//排序倍数大小
				for i := 0; i < len(multPlayers)-1; i++ {
					maxBetsPlayer := multPlayers[i]
					for j := i + 1; j < len(multPlayers); j++ {
						if this.Players[multPlayers[i]].PlayMultiple < this.Players[multPlayers[j]].PlayMultiple {
							maxBetsPlayer = multPlayers[j]
							multPlayers[j] = multPlayers[i]
							multPlayers[i] = maxBetsPlayer
						}
					}
				}

				// 1 通杀 ,2 赢两家 3 赢一家
				rand.Seed(time.Now().UnixNano())
				if rand.Perm(100)[0] < 25 {
					//通杀
					winChairs = append(winChairs, int(this.Banker))
				} else {

					var tempPlayers []int
					for i := 0; i < len(multPlayers); i++ {
						if i == 0 {
							continue
						}
						tempPlayers = append(tempPlayers, multPlayers[i])
					}
					if len(tempPlayers) > 1 {
						tempPlayers = ListShuffle(tempPlayers)
					}

					rand.Seed(time.Now().UnixNano())
					if len(multPlayers) > 2 && this.Players[multPlayers[0]].PlayMultiple > this.Players[multPlayers[1]].PlayMultiple+this.Players[multPlayers[2]].PlayMultiple &&
						rand.Perm(100)[0] < 55 {
						//赢一家
						winChairs = append(winChairs, tempPlayers...)
						winChairs = append(winChairs, int(this.Banker))
					} else {
						//赢两家
						winChairs = append(winChairs, tempPlayers[0])
						winChairs = append(winChairs, int(this.Banker))
					}
				}
			} else {
				//玩家坐庄
				// logs.Debug("玩家坐庄")
				//倍数大小数组玩家
				var multPlayers []int
				multPlayers = append(multPlayers, robotPlayers...)

				rand.Seed(time.Now().UnixNano())
				if len(multPlayers) > 1 && rand.Perm(100)[0] < 50 {
					//排序倍数大小
					for i := 0; i < len(multPlayers)-1; i++ {
						maxBetsPlayer := multPlayers[i]
						for j := i + 1; j < len(multPlayers); j++ {
							if this.Players[multPlayers[i]].PlayMultiple < this.Players[multPlayers[j]].PlayMultiple {
								maxBetsPlayer = multPlayers[j]
								multPlayers[j] = multPlayers[i]
								multPlayers[i] = maxBetsPlayer
							}
						}
					}
					// 输所有
					winChairs = append(winChairs, multPlayers...)
				} else {
					winChairs = append(winChairs, multPlayers[0])
				}
			}
		} else {
			//player win
			var robotList []int
			robotList = append(robotList, robotPlayers...)
			if len(robotList) > 1 {
				robotList = ListShuffle(robotList)
			}

			rand.Seed(time.Now().UnixNano())
			randNum := rand.Intn(100)
			//player allWin
			if randNum < GExtRobot.PlayerRate && GExtRobot.PlayerRate != 0 {
			} else {
				winChairs = append(winChairs, robotList[0])
			}

			var tempChairs []int
			//player win
			for _, v := range lastThisPlayers {
				if !v.Robot {
					tempChairs = append(tempChairs, int(v.ChairId))
				}
			}
			if len(tempChairs) > 0 {
				if len(tempChairs) > 1 {
					tempChairs = ListShuffle(tempChairs)
				}
				winChairs = append(winChairs, tempChairs...)
			}

			// logs.Debug("WIN玩家", winChairs)
		}

		//剩余玩家
		if len(this.Players) != len(winChairs) {
			var lastPlayers []int
			for i := 0; i < len(this.Players); i++ {
				flag := true
				for j := 0; j < len(winChairs); j++ {
					if int(this.Players[i].ChairId) == winChairs[j] {
						flag = false
						break
					}
				}
				if flag {
					lastPlayers = append(lastPlayers, int(this.Players[i].ChairId))
				}
			}
			if len(lastPlayers) > 0 {
				if len(lastPlayers) > 1 {
					lastPlayers = ListShuffle(lastPlayers)
				}
				winChairs = append(winChairs, lastPlayers...)
			}

			// logs.Debug("剩余玩家", winChairs)
		}
	} else {
		for _, v := range this.Players {
			winChairs = append(winChairs, int(v.ChairId))
		}
		winChairs = ListShuffle(winChairs)
	}

	logs.Debug("剩余玩家1111", winChairs)
	return winChairs
}

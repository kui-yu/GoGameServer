package main

import (
	"logs"
	"math/rand"
	"time"
)

func (this *ExtDesk) ControlResult() []int {

	//玩家随机
	var winChairs []int

	//robot
	var robotPlayers, realPlayers []int

	for index, v := range this.Players {
		if v.Robot {
			robotPlayers = append(robotPlayers, index)
		} else {
			realPlayers = append(realPlayers, index)
		}
	}

	if len(robotPlayers) > 0 && GetCostType() == 1 { //如果不是体验场再进行排序

		//黑名单
		var hierarchyLostPlayers []int

		//1.层级概率
		//判断每个玩家是否有层级概率
		for _, pp := range realPlayers {
			v := this.Players[pp]
			rand.Seed(time.Now().UnixNano())
			hierarchyRate := GetRateByHierarchyId(v.HierarchyId)
			logs.Debug("黑名单概率 ", hierarchyRate, v.HierarchyId)
			if hierarchyRate < 0 {
				hierarchyRate = -hierarchyRate
			}
			if hierarchyRate > 0 {
				logs.Debug("进入黑名单", hierarchyRate, v.HierarchyId)
				//走玩家黑名单层级判断
				if rand.Perm(100)[0] < int(hierarchyRate*10000/100) {
					hierarchyLostPlayers = append(hierarchyLostPlayers, pp)
				}
			}
		}

		//执行判断
		if len(hierarchyLostPlayers) == len(realPlayers) {
			//黑名单玩家 都是真实玩家，结束
			//robot win
			var robotList []int
			robotList = append(robotList, robotPlayers...)
			robotList = ListShuffleByInt(robotList)

			winChairs = append(winChairs, robotList[0])
		} else {
			//库存判断
			intervalRate := GetRateByInterval()
			logs.Debug("库存管控", intervalRate, this.RobotRate, G_DbGetGameServerData.GameConfig.CurrentStock, GDeskMgr.ChangeStock)
			if intervalRate < 0 {
				intervalRate = -intervalRate
			}
			if intervalRate > 0 {
				//有库存概率，使用库存概率，否则默认
				this.RobotRate = int(intervalRate * 10000 / 100)
			}

			rand.Seed(time.Now().UnixNano())
			if rand.Intn(100) < this.RobotRate {
				logs.Debug("进入库存")
				//robot win
				var robotList []int
				robotList = append(robotList, robotPlayers...)
				robotList = ListShuffleByInt(robotList)
				winChairs = append(winChairs, robotList[0])
			} else {
				// 排除黑名单玩家
				for _, pp := range realPlayers {
					flag := true
					for j := 0; j < len(hierarchyLostPlayers); j++ {
						if pp == hierarchyLostPlayers[j] {
							flag = false
							break
						}
					}
					if flag {
						winChairs = append(winChairs, pp)
					}
				}
			}
			//随机Win玩家
			winChairs = ListShuffleByInt(winChairs)
		}

		//剩余玩家
		if len(winChairs) != len(this.Players) {
			var lastPlayers []int
			for i := 0; i < len(this.Players); i++ {
				flag := true
				for j := 0; j < len(winChairs); j++ {
					if i == winChairs[j] {
						flag = false
						break
					}
				}
				if flag {
					lastPlayers = append(lastPlayers, i)
				}
			}
			//随机LAST玩家
			lastPlayers = ListShuffleByInt(lastPlayers)
			winChairs = append(winChairs, lastPlayers...)
		}
	} else {

		for index, _ := range this.Players {
			winChairs = append(winChairs, index)
		}
		winChairs = ListShuffleByInt(winChairs)
	}
	logs.Debug("玩家排序", winChairs)
	return winChairs
}

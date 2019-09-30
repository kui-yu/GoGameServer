package main

///////////////////////////////////////////////////
//此文件定义控制管理机器人的添加，删除，闲置，改变接口
///////////////////////////////////////////////////
type ExtController struct {
	Controller
}

// 接收任务， 机器人闲置状态通知
// 机器人结算的时候，发送这条协议，人数太多就归还，金币太多，太少就归还
// 一般不需要自己实现，如有需要再改，所以这边不放入公共地方
func (this *ExtController) onRobotIdle(id int32, d interface{}) {
	robot := d.(*ExtRobotClient)

	num := gameConfig.getGameConfigInt("num")
	curNum := len(this.robotClients)

	DebugLog("机器人当前数量:", curNum, num)

	forceOffline := false
	for i, token := range this.offlineRobotTokens {
		if token == robot.HallToken {
			forceOffline = true
			this.offlineRobotTokens = append(this.offlineRobotTokens[:i], this.offlineRobotTokens[i+1:]...)
			break
		}
	}

	if curNum > num {
		robot.Stop()
		this.deleteRobot(robot)

		curNum = len(this.robotClients)
		// TestLog("机器人删除之后数量:", curNum)
	} else if robot.Coin < int64(gameConfig.getGameConfigInt("shiftRobotCoinMin")) ||
		robot.Coin > int64(gameConfig.getGameConfigInt("shiftRobotCoinMax")) ||
		forceOffline == true {
		robot.Stop()
		DebugLog("替换机器人 %d %d", curNum, num)
		this.deleteRobot(robot)
		if curNum <= num {
			this.sendEvent(EVENT_CONT_ADDROBOT, 1)
		}
	} else {
		robot.sendGameAuto()
	}
}

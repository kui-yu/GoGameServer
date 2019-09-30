package main

// import (
// 	"strconv"
// )

///////////////////////////////////////////////////
//此文件定义控制管理机器人的添加，删除，闲置，改变接口
///////////////////////////////////////////////////
type ExtController struct {
	Controller
}

// 接收任务， 机器人闲置状态通知
// 机器人结算的时候，发送这条协议，人数太多就归还，金币太多，太少就归还
// 这边需要自己实现
func (this *ExtController) onRobotIdle(id int32, d interface{}) {
	robot := d.(*ExtRobotClient)

	num := gameConfig.getGameConfigInt("num")
	curNum := len(this.robotClients)

	DebugLog("机器人当前数量:", curNum, num)

	if curNum > num {
		robot.Stop()
		this.deleteRobot(robot)

		curNum = len(this.robotClients)
		DebugLog("机器人删除之后数量:", curNum)
	} else if robot.CarryCoin < robot.SmallBlind*3 ||
		robot.CarryCoin > gameConfig.getGameConfigCoin("maxcoin") {
		robot.Stop()
		DebugLog("替换机器人 %d %d", curNum, num)
		this.deleteRobot(robot)
		if curNum <= num {
			this.sendEvent(EVENT_CONT_ADDROBOT, 1)
		}
	}
}

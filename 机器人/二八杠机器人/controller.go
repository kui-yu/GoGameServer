package main

import (
	"strconv"
)

type Controller struct {
	robotClients []*ExtRobotClient // 所有机器人
	events       chan *EventMsg
	handlers     map[int32]func(int32, interface{})

	offlineRobotTokens []string // 等待下线的机器人
}

// 获得机器人集合
func (this *Controller) getRobotClients() []*ExtRobotClient {
	return this.robotClients
}

// 查找机器人
func (this *Controller) findRobotClient(uid int64) *ExtRobotClient {
	clen := len(this.robotClients)
	for i := 0; i < clen; i++ {
		if this.robotClients[i].UserInfo.Uid == uid {
			return this.robotClients[i]
		}
	}
	return nil
}

// 删除机器人
func (this *Controller) deleteRobot(robot *ExtRobotClient) bool {
	clen := len(this.robotClients)
	for i := 0; i < clen; i++ {
		if this.robotClients[i] == robot {
			this.robotClients = append(this.robotClients[:i], this.robotClients[i+1:]...)
			break
		}
	}

	GRabbitClient.AddMsgNative(struct {
		Coin  int64  `json:"coin"`
		Token string `json:"token"`
	}{
		Coin:  robot.Coin,
		Token: robot.HallToken,
	})

	return true
}

func (this *Controller) sendEvent(id int32, d interface{}) {
	select {
	case this.events <- &EventMsg{
		Id:   id,
		Data: d,
	}:

	default:
		ErrorLog("丢失信息")
	}

}

// 接收消息
func (this *Controller) RecvEventHandler() {
	go func() {
		for s := range this.events {
			if hand, ok := this.handlers[s.Id]; ok {
				hand(s.Id, s.Data)
			}
		}
	}()
}

// 查询机器人
func (this *Controller) findRobotFromToken(token string) *ExtRobotClient {
	clients := this.robotClients
	for _, client := range clients {
		if client.HallToken == token {
			return client
		}
	}

	return nil
}

func (this *Controller) InitBase() {
	this.events = make(chan *EventMsg, 100000)
	this.handlers = make(map[int32]func(int32, interface{}), 1000)

	this.RecvEventHandler()
}

// 机器人管理中心发送过来的控制协议
// 注册监听
func (this *ExtController) addEventListener() {
	this.handlers[EVENT_CONT_ADDROBOT] = this.onRobotAdd          //添加机器人通知
	this.handlers[EVENT_CONT_CONFIG_CHANGE] = this.onConfigChange // 配置文件改变通知
	this.handlers[EVENT_CONT_ROBOTIDLE] = this.onRobotIdle        // 机器人当前闲置通知
	this.handlers[EVENT_CONT_ROBOTSHIFT] = this.onRobotShift      // 机器人替换
	this.handlers[EVENT_CONT_OFFLINEROBOT] = this.onOfflineRobot  // 设置机器人下线
}

func (this *ExtController) Init() {
	this.addEventListener()
}

//公共接口
// 接收任务，添加机器人
func (this *ExtController) onRobotAdd(id int32, d interface{}) {
	num, _ := d.(int)
	DebugLog("接收任务，添加机器人 num:", num)

	for i := 0; i < num; i++ {
		robot := &ExtRobotClient{}
		robot.Start()
		this.robotClients = append(this.robotClients, robot)
	}

	DebugLog("添加机器人 完成")
}

// 接收任务， 配置文件改变
func (this *ExtController) onConfigChange(id int32, d interface{}) {
	DebugLog("接收任务， 配置文件改变")

	h := d.(struct {
		Key      string
		OldValue string
		NewValue string
	})

	// 机器人数量改变
	if h.Key == "num" {
		DebugLog("机器人数量改变")
		robotLen := len(this.robotClients)
		newLen, _ := strconv.Atoi(h.NewValue)

		// 通知机器人管理，当前数量变化
		GRobotServer.AddMsgNative(MSG_GAME_ROBOT_NUMCHANGE, struct {
			Id         int
			RobotCount int
		}{
			Id:         MSG_GAME_ROBOT_NUMCHANGE,
			RobotCount: newLen,
		})

		if newLen > robotLen { // 添加机器人
			this.sendEvent(EVENT_CONT_ADDROBOT, newLen-robotLen)
		}
		if newLen < robotLen { // 删除机器人
			dellen := newLen - robotLen
			removeArr := []*ExtRobotClient{}

			for _, client := range this.robotClients {
				if client.GameIn == false {
					removeArr = append(removeArr, client)
					dellen--
					if dellen == 0 {
						break
					}
				}
			}

			for {
				if len(removeArr) == 0 {
					break
				}
				this.sendEvent(EVENT_CONT_ROBOTIDLE, removeArr[0])
				removeArr = removeArr[1:]
			}
		}
	}
}

// 替换机器人事件
func (this *ExtController) onRobotShift(id int32, d interface{}) {
	robot := d.(*ExtRobotClient)
	succ := robot.Stop()
	if !succ {
		return
	}

	this.deleteRobot(robot)

	num := gameConfig.getGameConfigInt("num")
	curNum := len(this.robotClients)

	if curNum < num {
		this.sendEvent(EVENT_CONT_ADDROBOT, 1)
	}
	DebugLog("替换机器人")
}

// 设置机器人下线
func (this *ExtController) onOfflineRobot(id int32, d interface{}) {
	token := d.(string)

	robotClient := this.findRobotFromToken(token)
	if robotClient == nil {
		return
	}

	if robotClient.GameIn == false {
		this.sendEvent(EVENT_CONT_ROBOTSHIFT, robotClient)
	} else {
		this.offlineRobotTokens = append(this.offlineRobotTokens, token)
	}
}

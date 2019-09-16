package main

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

// 接收到配置信息
func (this *RobotClient) onRecvRobotOnline(data string) {
	fmt.Println("接收到配置信息")

	this.IsLogin = true

	info := GameInfo{}
	json.Unmarshal([]byte(data), &info)

	fmt.Println(info)
	this.GroupId = info.GroupId
	this.GameId = info.GameId
	this.GradeId = info.GradeId
	this.Name = info.Name
	this.GetRobotConfigUrl = info.GetRobotConfigUrl
	this.PutRobotConfigUrl = info.PutRobotConfigUrl
	this.CheckRobotUrl = info.CheckRobotUrl
	this.OfflineRobotUrl = info.OfflineRobotUrl
	this.Forceoffroboturl = info.Forceoffroboturl
	this.Forceonroboturl = info.Forceonroboturl
	this.RobotCount = info.RobotCount

	this.AddMsgNative(MSG_GAME_ROBOT_ONLINE, &RepBgInfo{
		Id:                MSG_GAME_ROBOT_ONLINE,
		BgRobotGetUrl:     GCONFIG.BgRobotGetUrl,
		BgRobotRestUrl:    GCONFIG.BgRobotRestUrl,
		BgRobotRestAllUrl: GCONFIG.BgRobotRestAllUrl,
		BgRobotAddCoinUrl: GCONFIG.BgRobotAddCoinUrl,
		BgRobotToken:      GCONFIG.BgRobotToken,
		BgRobotHallId:     GCONFIG.BgRobotHallId,
	})
}

// 接收到机器人数量变化
func (this *RobotClient) onRecvRobotNumChange(data string) {
	this.RobotCount = uint32(gjson.Get(data, "RobotCount").Int())
}

func (this *RobotClient) addHandler() {
	this.Handle[MSG_GAME_ROBOT_ONLINE] = this.onRecvRobotOnline
	this.Handle[MSG_GAME_ROBOT_NUMCHANGE] = this.onRecvRobotNumChange
}

func (this *RobotClient) removeHandler() {
	delete(this.Handle, MSG_GAME_ROBOT_ONLINE)
	delete(this.Handle, MSG_GAME_ROBOT_NUMCHANGE)
}

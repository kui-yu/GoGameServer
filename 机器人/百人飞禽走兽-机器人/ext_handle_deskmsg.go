package main

import (
	"encoding/json"
	"logs"
	"math/rand"
	"time"
)

func (this *ExtRobotClient) HandleAutoReply(d string) {
	data := GAutoGameReply{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("在监听匹配响应时json转换错误:", err)
		return
	}
	if data.CostType == 1 {
		logs.Debug("现在操作的是金币模式")
	} else {
		logs.Debug("现在操作的是积分模式")
	}
	//请求桌子信息
	this.AddMsgNative(MSG_GAME_INFO_QDESKINFO, struct {
		Id int
	}{
		Id: MSG_GAME_INFO_QDESKINFO,
	})
}
func (this *ExtRobotClient) HandleStatusChange(d string) {
	logs.Debug("接收到状态改变通知")
	data := GameStatuInfo{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("在接收状态时json转换错误：", err)
		return
	}
	this.Status = data.Status
	this.StatusTime = data.StatusTime
	if this.Status == GAME_STATUS_DOWNBET {
		logs.Debug("现在是下注状态，机器人模拟下注")
		this.bet()
	}
}

func (this *ExtRobotClient) HandleDeskReply(d string) {
	logs.Debug("处理桌子响应")
	data := GClientDeskInfo{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("出错")
	}
	this.Index = data.BetAbleIndex
	this.Status = data.GameStatus
	this.StatusTime = int(data.GameStatusDuration)
	if this.Status == GAME_STATUS_DOWNBET {
		this.bet()
	}
}
func (this *ExtRobotClient) HandleDownBetReplay(d string) {
	logs.Debug("正在处理下注返回")
	data := DownBetReplay{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("json转换错误")
	}
	this.Index = data.BetAbleIndex
}
func (this *ExtRobotClient) bet() {
	ra := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < this.StatusTime; i++ {
		if this.Status == GAME_STATUS_DOWNBET {
			logs.Debug("index;", this.Index)
			if this.Index != -1 {
				chipindex := ra.Intn(this.Index + 1)
				var areaindex int
				if this.Tobet != -1 {
					bil := ra.Int31n(100)
					if bil < 75 {
						areaindex = this.Tobet
					} else {
						areaindex = int(ra.Int31n(12))
					}
				} else {
					areaindex = int(ra.Int31n(12))
				}
				this.Tobet = -1
				this.AddMsgNative(MSG_GAME_INFO_DOWNBET, struct {
					Id        int
					ChipIndex int
					AreaIndex int
				}{
					Id:        MSG_GAME_INFO_DOWNBET,
					ChipIndex: chipindex,
					AreaIndex: areaindex,
				})
				time.Sleep(time.Second)

			} else {
				break
			}
		} else {
			break
		}
	}
}
func (this *ExtRobotClient) HandleToRobot(d string) {
	logs.Debug("接收到下注提示:")
	tor := ToRobot{}
	err := json.Unmarshal([]byte(d), &tor)
	if err != nil {
		logs.Debug("接收下注提示的时候json解析错误", tor)
	}
	this.Tobet = tor.Index
}

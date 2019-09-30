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
	data := GameStatuInfo{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("在接收状态时json转换错误：", err)
		return
	}
	this.Status = data.Status
	this.Time = data.StatusTime
	if this.Status == 12 {
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
	this.Time = int(data.GameStatusDuration)
	if data.GameStatus == 12 {
		this.bet()
	}

}
func (this *ExtRobotClient) HandleDownBetReplay(d string) {
	data := DownBetReplay{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("json转换错误")
	}
	this.Index = data.BetAbleIndex
}
func (this *ExtRobotClient) bet() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < this.Time; i++ {
		if this.Index != -1 {
			chipindex := rand.Intn(this.Index + 1)
			areaindex := rand.Intn(4)
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
	}
}

package main

import (
	"encoding/json"
	// "fmt"
	"math/rand"
	"time"

	"github.com/tidwall/gjson"
)

//匹配数据
func (this *ExtRobotClient) HandleAuto(d string) {
	// DebugLog("匹配数据", d)
	this.GameIn = true
	gameReply := GInfoAutoGameReply{}
	json.Unmarshal([]byte(d), &gameReply)
	for _, seat := range gameReply.Seat {
		if seat.Uid == this.UserInfo.Uid {
			this.Coin = seat.Coin
		}
	}
}

//游戏状态
func (this *ExtRobotClient) HandlerGameStatus(d string) {
	// DebugLog("当前阶段消息", d)
	if d != "" {
		Stage := gjson.Get(d, "Stage").Int()
		// fmt.Println("消息", this.PlayTime)
		if Stage == GAME_STATUS_START {

		} else if Stage == STAGE_PLAY {

			time.Sleep(time.Millisecond * 100)
			//随机休眠
			rand.Seed(time.Now().UnixNano())
			randTime := rand.Perm(this.PlayTime)[0] + 2
			time.Sleep(time.Second * time.Duration(randTime))
			//发送开牌信息
			this.AddMsgNative(MSG_GAME_INFO_PLAY, struct {
				Id int
			}{
				Id: MSG_GAME_INFO_PLAY,
			})
		}
	}
}

func (this *ExtRobotClient) getBet() int {
	// DebugLog("选定庄家", this.PlayerBets)
	betSize := len(this.PlayerBets)
	if betSize > 0 {
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Perm(betSize)[0]
		return this.PlayerBets[randNum]
	}
	return 1
}

//下注
func (this *ExtRobotClient) HandlerBetList(d string) {
	betList := GCallListMsg{}
	json.Unmarshal([]byte(d), &betList)

	this.PlayerBets = betList.BetList

	time.Sleep(time.Millisecond * 100)
	//随机休眠
	rand.Seed(time.Now().UnixNano())
	randTime := rand.Perm(this.BetTime)[0] + 2
	time.Sleep(time.Second * time.Duration(randTime))
	//发送下注信息
	this.AddMsgNative(MSG_GAME_INFO_CALL, struct {
		Id       int
		Multiple int //叫庄倍数
	}{
		Id:       MSG_GAME_INFO_CALL,
		Multiple: this.getBet(),
	})
}

//结算
func (this *ExtRobotClient) HandlerSettle(d string) {
	// // 第二种方式，需要先定义结构体
	win := GWinInfosReply{}
	json.Unmarshal([]byte(d), &win)

	this.GameIn = false

	for _, info := range win.Infos {
		if info.Uid == this.UserInfo.Uid {
			this.Coin = info.Coins
			//通知可以离开
			controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
			return
		}
	}
}

package main

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/tidwall/gjson"
)

//匹配数据
func (this *ExtRobotClient) HandleAuto(d string) {
	// DebugLog("匹配数据", d)
	gameReply := GInfoAutoGameReply{}
	json.Unmarshal([]byte(d), &gameReply)
	for _, seat := range gameReply.Seat {
		if seat.Uid == this.UserInfo.Uid {
			this.Coin = seat.Coin
		}
	}
	this.GameIn = true

}

//游戏状态
func (this *ExtRobotClient) HandlerGameStatus(d string) {
	// DebugLog("当前阶段消息", d)
	if d != "" {
		Stage := gjson.Get(d, "Stage").Int()

		// DebugLog("当前阶段消息", Stage)
		if Stage == GAME_STATUS_START {
			//抢庄
		} else if Stage == STAGE_CALL {
			//叫分
			time.Sleep(time.Millisecond * 100)
			//随机休眠
			rand.Seed(time.Now().UnixNano())
			randTime := rand.Perm(this.BetTime)[0] + 5
			time.Sleep(time.Second * time.Duration(randTime))

			//发送叫庄信息
			this.AddMsgNative(MSG_GAME_INFO_CALL, struct {
				Id       int
				Multiple int //叫庄倍数
			}{
				Id:       MSG_GAME_INFO_CALL,
				Multiple: this.getBet(),
			})
		} else if Stage == STAGE_PLAY {
			time.Sleep(time.Millisecond * 100)
			//随机休眠
			rand.Seed(time.Now().UnixNano())
			randTime := rand.Perm(this.PlayTime)[0] + 3
			time.Sleep(time.Second * time.Duration(randTime))
			//发送下注信息
			this.AddMsgNative(MSG_GAME_INFO_PLAY, struct {
				Id int
			}{
				Id: MSG_GAME_INFO_PLAY,
			})
		}
	}
}

func (this *ExtRobotClient) getCall() int {
	// DebugLog("下注列表2", this.CallBets)
	callSize := len(this.CallBets)
	if callSize > 0 {
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Perm(callSize)[0]
		return this.CallBets[randNum]
	}
	return 0
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

func (this *ExtRobotClient) HanlderBanker(d string) {
	banker := GCallBankReply{}
	json.Unmarshal([]byte(d), &banker)

	this.PlayerBets = []int{}
	this.PlayerBets = banker.BetList
}

func (this *ExtRobotClient) HandlerCallList(d string) {
	call := GSCallList{}
	json.Unmarshal([]byte(d), &call)

	this.CallBets = call.CallList
	time.Sleep(time.Millisecond * 100)
	//随机休眠
	rand.Seed(time.Now().UnixNano())
	randTime := rand.Perm(this.CallTime)[0] + 3
	time.Sleep(time.Second * time.Duration(randTime))
	//发送叫庄信息
	this.AddMsgNative(MSG_GAME_INFO_CALL_BANKER, struct {
		Id       int
		Multiple int //叫庄倍数
	}{
		Id:       MSG_GAME_INFO_CALL_BANKER,
		Multiple: this.getCall(),
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

package main

import (
	"encoding/json"
	"math/rand"
	"time"
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
	stage := GSStageInfo{}
	json.Unmarshal([]byte(d), &stage)

	if d != "" {
		Stage := stage.Stage
		// DebugLog("当前阶段消息", Stage)
		if Stage == STAGE_PLAY {
			time.Sleep(time.Millisecond * 100)
			//随机休眠
			rand.Seed(time.Now().UnixNano())
			randTime := rand.Perm(this.BetTime)[0] + 4
			time.Sleep(time.Second * time.Duration(randTime))
			//发送下注信息
			this.AddMsgNative(410006, struct {
				Id           int
				PlayMultiple int //下注倍数
			}{
				Id:           410006,
				PlayMultiple: this.getBet(),
			})
		}
	}
}

//获取抢庄筹码
func (this *ExtRobotClient) getCall() int {
	// DebugLog("下注列表2", this.CallBets)
	callSize := len(this.CallBets)
	if callSize > 0 {
		return this.CallBets[rand.Intn(callSize)]
	}
	return 0
}

//获取下注筹码
func (this *ExtRobotClient) getBet() int {
	// DebugLog("选定庄家", this.PlayerBets)
	betSize := len(this.PlayerBets)
	if betSize > 0 {
		return this.PlayerBets[rand.Intn(betSize)]
	}
	return 1
}

//抢庄
func (this *ExtRobotClient) HandlerCallList(d string) {
	callList := GSCallList{}
	json.Unmarshal([]byte(d), &callList)
	this.CallBets = callList.CallList
	// DebugLog("下注列表", d, this.CallBets)
	time.Sleep(time.Millisecond * 100)
	//随机休眠
	rand.Seed(time.Now().UnixNano())
	randTime := rand.Perm(this.CallTime)[0] + 3
	time.Sleep(time.Second * time.Duration(randTime))
	//发送叫庄信息
	this.AddMsgNative(410004, struct {
		Id           int
		CallMultiple int //叫庄倍数
	}{
		Id:           410004,
		CallMultiple: this.getCall(),
	})
}

//选定庄家
func (this *ExtRobotClient) HandlerBanker(d string) {
	callInfo := GSPlayerCallBank{}
	json.Unmarshal([]byte(d), &callInfo)

	this.Banker = callInfo.Banker
	this.BankerBet = callInfo.BankerMultiples
	this.PlayerBets = callInfo.BetList

	// DebugLog("选定庄家", d, this.PlayerBets)
}

//结算
func (this *ExtRobotClient) HandlerSettle(d string) {
	win := GSSettleInfoEnd{}
	json.Unmarshal([]byte(d), &win)

	this.GameIn = false

	for _, info := range win.PlayInfos {
		if info.Uid == this.UserInfo.Uid {
			this.Coin = info.Coins
			//通知可以离开
			controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
			return
		}
	}
}

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
		if Stage == STAGE_PLAY {
			//随机休眠
			rand.Seed(time.Now().UnixNano())
			randTime := rand.Perm(this.PlayTime)[0] + 5
			time.Sleep(time.Second * time.Duration(randTime))
			//发送bai牌信息
			this.AddMsgNative(MSG_GAME_INFO_PLAY, struct {
				Id        int
				PlayType  int //0 自己摆牌 ；摆特殊牌型
				PlayCards []int
			}{
				Id:       MSG_GAME_INFO_PLAY,
				PlayType: this.SpecialType,
			})
		}
	}
}

func (this *ExtRobotClient) HandlerCardInfo(d string) {
	handInfo := GSHandInfo{}
	json.Unmarshal([]byte(d), &handInfo)

	if handInfo.SpecialType > 9 {
		this.SpecialType = handInfo.SpecialType
	} else {
		this.SpecialType = 1
	}
}

//结算
func (this *ExtRobotClient) HandlerSettle(d string) {
	// // 第二种方式，需要先定义结构体
	win := GSSettleInfos{}
	json.Unmarshal([]byte(d), &win)

	this.GameIn = false

	for _, info := range win.PlayInfo {
		if info.Uid == this.UserInfo.Uid {
			this.Coin = info.Coins
			//通知可以离开
			controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
			return
		}
	}
}

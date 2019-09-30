package main

import (
	"encoding/json"
)

//匹配成功发回的座位信息
func (this *ExtRobotClient) HandleGameAutoReply(d string) {
	rsp := GInfoAutoGameReply{}
	json.Unmarshal([]byte(d), &rsp)
	for _, v := range rsp.Seat {
		if v.Uid == this.UserInfo.Uid {
			this.SeatId = v.Cid
		}
	}
	//已经在游戏中
	this.GameIn = true
}

func (this *ExtRobotClient) GameEndNotify(data string) {
	//
	rsp := GInfoGameEnd{}
	json.Unmarshal([]byte(data), &rsp)
	//
	this.GameIn = false

	for i, v := range rsp.Coins {
		if int32(i) == this.SeatId {
			this.Coin = v
			break
		}
	}
	//
	controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
}

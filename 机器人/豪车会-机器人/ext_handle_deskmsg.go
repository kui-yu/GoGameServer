package main

import (
	"encoding/json"
	"logs"
	"math/rand"
	"time"
)

func (this *ExtRobotClient) HandleStatusChange(d string) {
	data := GSGameStatusInfo{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		return
	}
	this.Status = data.GameStatus
	this.Time = int(data.GameStatusDuration)
	if this.Status == GAME_STATUS_DOWNBET {
		this.bet()
	}
}

func (this *ExtRobotClient) HandleDeskReply(d string) {
	data := GClientDeskInfo{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("出错")
	}
	this.Index = data.CanUserChip
	this.Time = int(data.GameStatusDuration)
	if data.GameStatus == GAME_STATUS_DOWNBET {
		this.bet()
	}
}
func (this *ExtRobotClient) HandleDownBetReplay(d string) {
	data := GSDownBet{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("json转换错误")
	}
	this.Index = data.CanUserChip
}

func (this *ExtRobotClient) bet() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < this.Time; i++ {
		if this.Index != -1 {
			chipindex := rand.Intn(this.Index + 1)
			areaindex := rand.Intn(8)
			this.AddMsgNative(MSG_GAME_INFO_QDOWNBET, struct {
				Id      int
				CoinIdx int
				BetsIdx int
			}{
				Id:      MSG_GAME_INFO_QDOWNBET,
				CoinIdx: chipindex,
				BetsIdx: areaindex,
			})
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

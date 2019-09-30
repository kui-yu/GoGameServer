/**
* 结算
**/
package main

import (
	"encoding/json"
)

type FsmResults struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FsmResults) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FsmResults) GetMark() int {
	return this.Mark
}

func (this *FsmResults) Run(upMark int, overtime int) {
	DebugLog("进入游戏状态：结算")
	this.UpMark = upMark

	this.addListener() // 添加监听
}

func (this *FsmResults) Leave() {
	this.removeListener()
}

func (this *FsmResults) onEvent(interface{}) {

}

// 添加网络监听
func (this *FsmResults) addListener() {
	this.RC.Handle[MSG_GAME_NGameResult] = this.onGameResults
}

// 删除网络监听
func (this *FsmResults) removeListener() {
	delete(this.RC.Handle, MSG_GAME_NGameResult)
}

type Results struct {
	Uid           int64
	Sid           int
	Value         int64
	WaterProfit   float64
	Cards         []int
	CardGroupType int
}

func (this *FsmResults) onGameResults(str string) {
	DebugLog(str)
	DebugLog("手牌", this.RC.Cards)

	data := struct {
		JackpotVal int64
		Results    []Results
	}{}
	json.Unmarshal([]byte(str), &data)
	for _, result := range data.Results {
		if result.Uid == this.RC.UserInfo.Uid {
			this.RC.Coin += result.Value
			this.RC.CarryCoin += result.Value
			break
		}
	}

	this.RC.Reset()
}

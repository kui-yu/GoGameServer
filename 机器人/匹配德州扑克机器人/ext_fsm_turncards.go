/**
* 第四张公共牌
**/
package main

import (
	"encoding/json"
)

type FsmTurnCards struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FsmTurnCards) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FsmTurnCards) GetMark() int {
	return this.Mark
}

func (this *FsmTurnCards) Run(upMark int, overtime int) {
	DebugLog("进入游戏状态：第四张公共牌")
	this.UpMark = upMark

	this.addListener() // 添加监听
}

func (this *FsmTurnCards) Leave() {
	this.removeListener()
}

func (this *FsmTurnCards) onEvent(interface{}) {

}

// 添加网络监听
func (this *FsmTurnCards) addListener() {
	this.RC.Handle[MSG_GAME_NGamePublicCards] = this.onGamePublic
}

// 删除网络监听
func (this *FsmTurnCards) removeListener() {
	delete(this.RC.Handle, MSG_GAME_NGamePublicCards)
}

func (this *FsmTurnCards) onGamePublic(s string) {
	data := struct {
		Type  int
		Cards []int
	}{}

	json.Unmarshal([]byte(s), &data)

	this.RC.Cards = append(this.RC.Cards, data.Cards...)
}

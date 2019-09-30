/**
* 发给玩家的两张牌
**/
package main

import (
	"encoding/json"
)

type FsmHoleCards struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FsmHoleCards) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FsmHoleCards) GetMark() int {
	return this.Mark
}

func (this *FsmHoleCards) Run(upMark int, overtime int) {
	DebugLog("进入游戏状态：发给玩家的两张牌")
	this.UpMark = upMark

	this.addListener() // 添加监听
}

func (this *FsmHoleCards) Leave() {
	this.removeListener()
}

func (this *FsmHoleCards) onEvent(interface{}) {

}

// 添加网络监听
func (this *FsmHoleCards) addListener() {
	this.RC.Handle[MSG_GAME_NGameHoleCards] = this.onGameHoleCards
}

// 删除网络监听
func (this *FsmHoleCards) removeListener() {
	delete(this.RC.Handle, MSG_GAME_NGameHoleCards)
}

func (this *FsmHoleCards) onGameHoleCards(str string) {
	data := struct {
		Id    int
		Cards []int
	}{}
	json.Unmarshal([]byte(str), &data)

	this.RC.Cards = append(this.RC.Cards, data.Cards...)
}

/**
* 头三张公共牌
**/
package main

import (
	"encoding/json"
)

type FsmFlopCards struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FsmFlopCards) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FsmFlopCards) GetMark() int {
	return this.Mark
}

func (this *FsmFlopCards) Run(upMark int, overtime int) {
	DebugLog("进入游戏状态：发头三张公共牌")
	this.UpMark = upMark

	this.addListener() // 添加监听
}

func (this *FsmFlopCards) Leave() {
	this.removeListener()
}

func (this *FsmFlopCards) onEvent(interface{}) {

}

// 添加网络监听
func (this *FsmFlopCards) addListener() {
	this.RC.Handle[MSG_GAME_NGamePublicCards] = this.onPublicCards
}

// 删除网络监听
func (this *FsmFlopCards) removeListener() {
	delete(this.RC.Handle, MSG_GAME_NGamePublicCards)
}

func (this *FsmFlopCards) onPublicCards(s string) {
	data := struct {
		Type  int
		Cards []int
	}{}

	json.Unmarshal([]byte(s), &data)

	this.RC.Cards = append(this.RC.Cards, data.Cards...)
}

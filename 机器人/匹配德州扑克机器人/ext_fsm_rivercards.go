/**
* 第五张公共牌
**/
package main

import (
	"encoding/json"
)

type FsmRiverCards struct {
	UpMark int
	Mark   int
	RC     *ExtRobotClient
}

func (this *FsmRiverCards) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FsmRiverCards) GetMark() int {
	return this.Mark
}

func (this *FsmRiverCards) Run(upMark int, overtime int) {
	DebugLog("进入游戏状态：第五张公共牌")
	this.UpMark = upMark

	this.addListener() // 添加监听
}

func (this *FsmRiverCards) Leave() {
	this.removeListener()
}

func (this *FsmRiverCards) onEvent(interface{}) {

}

// 添加网络监听
func (this *FsmRiverCards) addListener() {
	this.RC.Handle[MSG_GAME_NGamePublicCards] = this.onGamePublic
}

// 删除网络监听
func (this *FsmRiverCards) removeListener() {
	delete(this.RC.Handle, MSG_GAME_NGamePublicCards)
}

func (this *FsmRiverCards) onGamePublic(s string) {
	data := struct {
		Type  int
		Cards []int
	}{}

	json.Unmarshal([]byte(s), &data)

	this.RC.Cards = append(this.RC.Cards, data.Cards...)
}

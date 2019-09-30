package main

import (
	"encoding/json"
	//"logs"
)

//玩家准备
func (this *ExtDesk) HandleReady(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != GAME_STATUS_FREE && this.GameState != STAGE_INIT {
		return
	}

	if this.TableConfig.GameModule == 2 {
		//判断金币是否充足
		if p.Coins < int64(this.Bscore*108) {
			p.SendNativeMsg(MSG_GAME_INFO_ERR, GSInfoErr{
				Id:  MSG_GAME_INFO_ERR,
				Err: "金币不足，请充值",
			})
			this.BroadcastOther(p, MSG_GAME_INFO_ERR, GSInfoErr{
				Id:  MSG_GAME_INFO_ERR,
				Err: "玩家" + p.Nick + " 金币不足，请等待",
			})
			return
		}
	}
	data := GAPlayerReady{}
	json.Unmarshal([]byte(d.Data), &data)

	p.IsReady = 1

	this.BroadcastAll(MSG_GAME_INFO_READY_REPLY, &GSPlayerReady{
		Id:      MSG_GAME_INFO_READY_REPLY,
		ChairId: p.ChairId,
		IsReady: p.IsReady,
	})

	flag := true
	for _, v := range this.Players {
		if v.IsReady == 0 {
			flag = false
		}
	}
	if flag && len(this.Players) >= this.TableConfig.PlayerNumber {
		//初始化卡牌数量列表
		this.InitGame()
		//进入游戏开始
		this.nextStage(GAME_STATUS_START)
	}
}

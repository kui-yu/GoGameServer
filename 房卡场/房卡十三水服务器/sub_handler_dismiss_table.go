package main

import (
	"encoding/json"
	"logs"
)

//玩家解散
func (this *ExtDesk) HandleDisMiss(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("玩家解散")
	data := GADismiss{}
	json.Unmarshal([]byte(d.Data), &data)

	if data.IsDismiss == 1 {
		//同意解散
		if this.GameState != STAGE_DISMISS && this.GameState == STAGE_INIT {
			//跳到游戏解散阶段
			this.nextStage(STAGE_DISMISS)
		}
		if this.GameState != STAGE_DISMISS {
			p.SendNativeMsg(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
				Id:        MSG_GAME_INFO_DISMISS_REPLY,
				DisPlayer: this.DisPlayer,
				IsDismiss: 2,
				Message:   "请在准备阶段发起解散",
			})
			return
		}
		this.DisPlayer = append(this.DisPlayer, p.ChairId)

	} else {
		this.ClearTimer()
		//解散清空
		this.DisPlayer = []int32{}
		//回到游戏准备阶段
		this.GameState = STAGE_INIT
		this.BroadStageTime(0)

		this.BroadcastAll(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
			Id:        MSG_GAME_INFO_DISMISS_REPLY,
			DisPlayer: this.DisPlayer,
			IsDismiss: 0,
			Message:   p.Nick + "不同意解散",
		})
		return
	}
	//

	if len(this.DisPlayer) == len(this.Players) {
		this.BroadcastAll(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
			Id:        MSG_GAME_INFO_DISMISS_REPLY,
			DisPlayer: this.DisPlayer,
			IsDismiss: 3,
		})
		//同意人数已满，解散
		this.TimerOver()
	} else {
		//部分同意
		this.BroadcastAll(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
			Id:        MSG_GAME_INFO_DISMISS_REPLY,
			DisPlayer: this.DisPlayer,
			IsDismiss: 1,
		})
	}
}

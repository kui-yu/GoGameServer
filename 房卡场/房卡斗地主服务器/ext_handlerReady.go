//玩家准备状态
package main

import (
	"encoding/json"
	"fmt"
	"logs"
)

func (this *ExtDesk) HandleReady(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("玩家进入准备状态")
	//如果该桌子不是 空闲状态 或者是游戏准备状态的话，那么是无法准备的
	if this.GameState != GAME_STATUS_READ && this.GameState != GAME_STATUS_FREE {
		logs.Debug("该桌子不是准备状态!")
		return
	}
	data := GAPlayerReady{}
	json.Unmarshal([]byte(d.Data), &data)
	p.isReady = data.IsReady
	fmt.Println(p.isReady)
	this.BroadcastAll(MSG_GAME_INFO_READY_REPLY, &GSPlayerReady{
		Id:      MSG_GAME_INFO_READY_REPLY,
		ChairId: p.ChairId,
		IsReady: p.isReady,
	})
	flag := true

	for _, v := range this.Players {
		if v.isReady == 0 {
			flag = false
		}
	}
	if flag && len(this.Players) >= this.TableConfig.PlayerNum {
		logs.Debug("所有玩家准备完毕，进入游戏")
		this.Round++
		if this.Round == 1 {
			for _, v := range this.Players {
				if this.FkOwner == v.Uid {
					this.GToHAddRoomCard(v.Uid, -this.getPayMoney())
				}
			}

		}
		logs.Debug("当前回合:", this.Round)
		this.GameState = GAME_STATUS_START
		this.BroadStageTime(TIMER_START_NUM)
		//发送游戏开始通知 并附带额外信息
		this.BroadcastAll(MSG_GAME_INFO_START, &GGameStartNotify{
			Id:    MSG_GAME_INFO_START,
			Round: this.Round,
		})
		//游戏开始，进入发牌阶段
		this.TimerSendCard(nil)
		return
	}
	logs.Debug("玩家未准备或者玩家不够")
}

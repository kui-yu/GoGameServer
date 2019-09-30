package main

import (
	"logs"
)

func (this *ExtDesk) HandleDownZhuang(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("接收到玩家下庄请求", p.Nick)
	if this.GameState != GAME_STATUS_READY {
		p.SendNativeMsg(MSG_GAME_INFO_DOWNZHUANG_REPLY, ChangZhuangReply{
			Id:     MSG_GAME_INFO_DOWNZHUANG_REPLY,
			Result: 2,
			Err:    "目前不是准备状态..",
		})
		return
	}
	this.DownZhuang(p)
}

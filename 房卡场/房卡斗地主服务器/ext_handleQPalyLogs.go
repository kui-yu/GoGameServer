package main

import (
	"fmt"
	"logs"
)

func (this *ExtDesk) handleQPalyLogs(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!接收到 游戏记录请求")
	p.GGameLogs.Id = MSG_GAME_INFO_QPLAYERLOGS_REPLY
	gl := p.GGameLogs
	p.SendNativeMsg(MSG_GAME_INFO_QPLAYERLOGS_REPLY, gl)
	fmt.Println("ID：", MSG_GAME_INFO_QPLAYERLOGS_REPLY, "LOGS:", gl)
}

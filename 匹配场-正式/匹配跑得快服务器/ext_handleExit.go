package main

import "logs"

func (this *ExtDesk) HandleExit(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("接收到玩家离开请求")
	if this.GameState != GAME_STATUS_FREE {
		logs.Error("无法退出，游戏已经开始")
		p.SendNativeMsg(MSG_GAME_INFO_EXIT_REPLY, ExitReply{
			Id:     MSG_GAME_INFO_EXIT_REPLY,
			Result: 1,
			Err:    "正在游戏中，无法退出",
		})
		return
	}
	p.SendNativeMsg(MSG_GAME_INFO_EXIT_REPLY, ExitReply{
		Id:     MSG_GAME_INFO_EXIT_REPLY,
		Result: 0,
	})
	this.LeaveByForce(p)
}

package main

func (this *ExtDesk) GameStateDismiss() {
	this.BroadStageTime(STAGE_DISMISS_TIME)
	//进入倒计时
	this.runTimer(STAGE_DISMISS_TIME, this.GameStateDismissEnd)
}

//阶段-解散
func (this *ExtDesk) GameStateDismissEnd(d interface{}) {
	//系统定义结束
	this.GameState = GAME_STATUS_END
	this.BroadStageTime(TIMER_OVER_NUM)
	if GetCostType() != 1 {
		for _, p := range this.Players {
			p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:      MSG_GAME_LEAVE_REPLY,
				Result:  0,
				Cid:     p.ChairId,
				Uid:     p.Uid,
				Robot:   p.Robot,
				NoToCli: true,
			})
		}
	}
	this.BroadcastAll(MSG_GAME_INFO_DISMISS_REPLY, &GSDismiss{
		Id:        MSG_GAME_INFO_DISMISS_REPLY,
		DisPlayer: this.DisPlayer,
		IsDismiss: 3,
	})
	//同意人数已满，解散
	//房间结束
	this.TimerOver("")
}

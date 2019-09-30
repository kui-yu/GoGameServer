package main

import (
	"encoding/json"
	"logs"
)

//玩家离开
func (this *ExtDesk) HandleLeave(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("玩家离开")
	//广播给其他人，掉线
	if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
		if p.ChairId == 0 {
			this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:     MSG_GAME_LEAVE_REPLY,
				Cid:    p.ChairId,
				Uid:    p.Uid,
				Result: 1,
				Token:  p.Token,
				Err:    "房主解散该房间",
			})
			logs.Debug("房主解散该房间")
			this.TimerOver()
		} else {
			this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:     MSG_GAME_LEAVE_REPLY,
				Cid:    p.ChairId,
				Uid:    p.Uid,
				Result: 0,
				Token:  p.Token,
			})
			this.DelPlayer(p.Uid)
			this.DeskMgr.LeaveDo(p.Uid)
		}
	} else {
		this.BroadcastAll(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Cid:    p.ChairId,
			Uid:    p.Uid,
			Result: 2,
			Token:  p.Token,
			Err:    "游戏已开始，请发起解散",
		})
	}
}

//玩家解散
func (this *ExtDesk) HandleDisMiss(p *ExtPlayer, d *DkInMsg) {
	// logs.Debug("玩家发起解散")
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
	// fmt.Println("同意解散玩家人数：", len(this.DisPlayer))
	// fmt.Println("总的玩家", len(this.Players))
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

package main

import (
	"logs"
)

func (this *ExtDesk) HandleGameAutoFinal(p *ExtPlayer, d *DkInMsg) {
	// 群发用户房间人数和座位信息
	gameReply := GSInfoAutoGame{
		Id: MSG_GAME_INFO_AUTO_REPLY,
	}
	for _, v := range this.Players {
		seat := GSeatInfos{
			Uid:    v.Uid,
			Nick:   v.Nick,
			Cid:    v.ChairId,
			Sex:    v.Sex,
			Head:   v.Head,
			Lv:     v.Lv,
			Coin:   v.Coins,
			Bscore: int64(G_DbGetGameServerData.Bscore),
		}
		gameReply.Seat = append(gameReply.Seat, seat)
	}
	this.BroadcastAll(MSG_GAME_INFO_AUTO_REPLY, &gameReply)

	//判断人员是否已满，开启游戏
	if len(this.Players) >= GCONFIG.PlayerNum {
		//发送房间信息,底分
		sd := GGameStartNotify{
			Id: MSG_GAME_INFO_START,
		}
		this.BroadcastAll(MSG_GAME_INFO_START, &sd)
	}

	this.GameState = GAME_STATUS_START
}

func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("进入匹配")
	if this.GameState != GAME_STATUS_FREE {
		return
	}

	// 初始化
	p.ResetExtPlayer()

	//发送匹配成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})

	this.HandleGameAutoFinal(p, d)
}

//重连消息
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:     MSG_GAME_RECONNECT_REPLY,
		Result: 1,
		Err:    "不支持重连",
	})

	this.GameOverLeave()
	//归还桌子
	this.GameState = GAME_STATUS_FREE
	this.ResetTable()
	this.DeskMgr.BackDesk(this)
	return
}

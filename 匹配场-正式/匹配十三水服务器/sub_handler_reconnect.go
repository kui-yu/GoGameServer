package main

//重连消息
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	if this.GameState == GAME_STATUS_FREE || this.GameState == GAME_STATUS_END {
		p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
			Id:     MSG_GAME_RECONNECT_REPLY,
			Result: 1,
			Err:    "本桌子没有正在的游戏",
		})
		return
	}
	//400011 重连消息
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:       MSG_GAME_RECONNECT_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})

	//群发用户信息
	for _, v := range this.Players {
		gameReply := GSInfoAutoGame{
			Id: MSG_GAME_INFO_AUTO_REPLY,
		}
		for _, p := range this.Players {
			seat := GSeatInfo{
				Uid:  p.Uid,
				Nick: p.Nick,
				Cid:  p.ChairId,
				Sex:  p.Sex,
				Head: p.Head,
				Lv:   p.Lv,
				Coin: p.Coins,
			}
			if p.Uid != v.Uid {
				seat.Nick = "***" + seat.Nick[len(seat.Nick)-4:]
			}
			gameReply.Seat = append(gameReply.Seat, seat)
		}
		v.SendNativeMsg(MSG_GAME_INFO_AUTO_REPLY, &gameReply)
	}

	//房间消息
	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &GSTableInfo{
		Id:      MSG_GAME_INFO_ROOM_NOTIFY,
		TableId: this.JuHao,
		BScore:  this.Bscore,
	})

	//桌子消息
	reconnectInfo := GSReconnectInfo{}
	for _, v := range this.Players {
		reconnectInfo.ChairIds = append(reconnectInfo.ChairIds, v.ChairId)
		//当前状态
		if v.LiXian {
			reconnectInfo.States = append(reconnectInfo.States, 2)
		} else {
			reconnectInfo.States = append(reconnectInfo.States, 1)
		}
		//已摆牌
		if v.IsPlay > 0 {
			reconnectInfo.PlayChairIds = append(reconnectInfo.PlayChairIds, v.ChairId)
			reconnectInfo.PlayNum += 1
		}
	}
	reconnectInfo.HandCards = p.HandCards
	reconnectInfo.SpecialType = p.SpecialType
	reconnectInfo.Stage = this.GameState
	reconnectInfo.StageTime = this.TList[0].T
	reconnectInfo.Id = MSG_GAME_INFO_RECONNECT
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT, &reconnectInfo)
	//在线状态更新
	p.LiXian = false
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 1,
	})
}

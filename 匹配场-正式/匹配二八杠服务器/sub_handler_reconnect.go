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
		Result:   0,
		CostType: GetCostType(),
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

	//返回玩家状态
	status := GSeatInfoReconnect{}
	for _, v := range this.Players {
		status.ChairIds = append(status.ChairIds, v.ChairId)
		//当前状态
		if v.LiXian {
			status.States = append(status.States, 2)
		} else {
			status.States = append(status.States, 1)
		}
		ListAdd(&status.CallMultiples, v.CallMultiple)
		ListAdd(&status.PlayMultiples, v.PlayMultiple)
	}
	status.MyCard = p.HandCards
	if this.GameState == GAME_STATUS_START || this.GameState == STAGE_CALL {
		status.BankerId = -1
		status.BankerMultiple = -1
	} else {
		status.BankerId = this.Banker
		status.BankerMultiple = this.Players[this.Banker].CallMultiple
	}
	status.Round = this.Round
	status.PutInfos = this.PutInfos
	status.Stage = this.GameState
	status.StageTime = this.TList[0].T
	status.Id = MSG_GAME_INFO_RECONNECT
	status.RsInfo = this.RsInfo
	status.CallListCnt = len(p.PlayerCalls)
	status.CallList = p.PlayerCalls
	status.BetListCnt = len(p.PlayerBets)
	status.BetList = p.PlayerBets
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT, &status)
	//在线状态更新
	p.LiXian = false
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 1,
	})
}

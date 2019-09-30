package main

//重连
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	// logs.Debug("重连")
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
	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &GTableInfoReply{
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
		//倍数
		status.Multiples = append(status.Multiples, v.BetMultiple)
		if v.IsLook {
			status.PlayNum += 1
			status.PlayChairIds = append(status.PlayChairIds, v.ChairId)
		}
		if v.CallBankFlag {
			status.CallMultiples = append(status.CallMultiples, v.CallMultiple)
		} else {
			status.CallMultiples = append(status.CallMultiples, -1)
		}
	}
	status.Id = MSG_GAME_INFO_RECONNECT
	var handCards []int32
	for _, sv := range p.HandCard {
		handCards = append(handCards, int32(sv))
	}
	status.MyCard = handCards
	//庄家消息
	status.Banker = this.Banker
	status.BankerMultiples = this.Players[this.Banker].CallMultiple
	//返回阶段消息
	status.Stage = this.GameState
	status.StageTime = this.TList[0].T

	status.CallListCnt = len(p.PlayerCalls)
	status.CallList = p.PlayerCalls
	status.BetListCnt = len(p.PlayerBets)
	status.BetList = p.PlayerBets

	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT, status)

	//在线状态更新
	p.LiXian = false
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 1,
	})

}

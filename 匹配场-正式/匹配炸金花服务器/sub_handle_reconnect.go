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
		Round:   GameRound - 1,
	})

	//返回玩家状态
	status := GSeatInfoReconnect{}
	for _, v := range this.Players { // 分配椅子
		status.ChairIds = append(status.ChairIds, v.ChairId)
		//当前状态
		if v.LiXian {
			status.States = append(status.States, 2)
		} else if v.IsLeave == 1 {
			status.States = append(status.States, 2)
		} else {
			status.States = append(status.States, 1)
		}
		sum := int64(0)
		for i := 0; i < len(v.PayCoin); i++ {
			sum += v.PayCoin[i]
		}
		status.PayCoin = append(status.PayCoin, sum)
		if v.CardType != 2 {
			status.CardType = append(status.CardType, v.CardType)
		} else {
			if v.IsGU {
				status.CardType = append(status.CardType, v.CardType)
			} else {
				status.CardType = append(status.CardType, 3)
			}
		}

	}
	for i := 0; i < len(status.PayCoin); i++ {
		status.CoinList += status.PayCoin[i]
	}
	status.CallPlayer = this.CallPlayer
	status.MinCoin = this.MinCoin
	status.Round = this.Round
	status.Stage = this.GameState
	if len(this.TList) != 0 {
		status.TimeRemaining = this.TList[0].T
	}
	if this.GameState == STAGE_PLAY_OPERATION {
		status.StageTime = 15
	} else if this.GameState == STAGE_START_TIME {
		status.StageTime = 2
	}

	status.Id = MSG_GAME_INFO_RECONNECT
	status.ReconnectPlayer = GSPlayerConnect{ChairId: p.ChairId, AutoFollowUp: p.AutoFollowUp, ProtectGU: p.ProtectGU}
	if p.CardType != 2 {
		status.ReconnectPlayer.CardType = p.CardType
		if p.CardType == 1 {
			status.ReconnectPlayer.HandCard = p.OldHandCard
			status.ReconnectPlayer.CardLv = p.CardLv
		}
	} else if p.IsGU && p.CardType == 2 {
		status.ReconnectPlayer.CardType = 2
	} else {
		status.ReconnectPlayer.CardType = 3
	}
	status.ReconnectPlayer.CoinEnough = IsCoinEnough(p.Coins, p.PayCoin, this.Bscore, this.MinCoin, p.CardType)
	// status.RsInfo = this.RsInfo
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT, &status)
	//在线状态更新
	p.LiXian = false
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 1,
	})
}

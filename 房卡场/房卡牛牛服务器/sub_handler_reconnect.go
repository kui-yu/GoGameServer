package main

//重连
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	// logs.Debug("重连")
	//400011 重连消息
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:     MSG_GAME_RECONNECT_REPLY,
		Result: 0,
	})
	//群发用户信息
	gameReply := GSInfoAutoGame{
		Id: MSG_GAME_INFO_AUTO_REPLY,
	}
	//返回所有用户信息
	for _, v := range this.Players {
		// var coin int64
		// if this.TableConfig.GameModule == 2 {
		// 	coin = v.Coins
		// }
		seat := GSSeatInfo{
			Uid:     v.Uid,
			Nick:    v.Nick,
			Cid:     v.ChairId,
			Sex:     v.Sex,
			Head:    v.Head,
			Lv:      v.Lv,
			Coin:    v.TotalCoins,
			IsReady: v.IsReady,
		}
		gameReply.Seat = append(gameReply.Seat, seat)
	}
	p.SendNativeMsg(MSG_GAME_INFO_AUTO_REPLY, &gameReply)

	//房间消息
	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &GTableInfoReply{
		Id:      MSG_GAME_INFO_ROOM_NOTIFY,
		TableId: this.FkNo,
		Config:  this.TableConfig,
	})

	//返回玩家状态
	status := GSSeatInfoReconnect{}
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

		status.CallMultiples = append(status.CallMultiples, v.CallMultiple)
	}
	status.Id = MSG_GAME_INFO_RECONNECT
	var handCards []int
	for _, sv := range p.HandCard {
		handCards = append(handCards, sv)
	}
	status.MyCard = handCards
	//庄家消息
	status.Banker = this.Banker
	status.Round = this.Round
	if this.Banker != -1 {
		status.BankerMultiples = this.Players[this.Banker].CallMultiple
	}
	//返回阶段消息
	status.Stage = this.GameState
	if len(this.TList) == 0 {
		status.StageTime = 0
	} else {
		status.StageTime = this.TList[0].T
	}
	status.DisPlayer = this.DisPlayer
	status.BetListCnt = len(p.PlayerBets)
	status.BetList = p.PlayerBets
	status.CallListCnt = len(p.PlayerCalls)
	status.CallList = p.PlayerCalls

	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT, status)

	//在线状态更新
	p.LiXian = false
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 1,
	})

}

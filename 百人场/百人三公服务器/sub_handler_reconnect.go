package main

//掉线重连
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	//断线重连应答
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:       MSG_GAME_RECONNECT_REPLY,
		Result:   0,
		CostType: GetCostType(),
	})
	p.LiXian = false
	//桌面玩家更新
	this.UadatePlayer(6)
	//发送房间信息给刚进场的玩家
	roomInfo := &GSPlayerIn{
		Id:         MSG_GAME_RECONNECT_TABLE_REPLY,
		GameId:     GetJuHao(),
		AllPlayer:  len(this.Players),
		ManyPlayer: p.GetDeskPlayerInfo(),                        //其他玩家
		Stage:      this.Stage,                                   //当前游戏阶段
		MaxBet:     G_DbGetGameServerData.GameConfig.LimitRedMax, //限红
		GameTrend:  this.GameTrend,                               //游戏走势，根据索引分别为庄、黑、红、梅、方， Trend.Player 0为闲家，1为庄家
	}
	roomInfo.DeskInfos = DeskInfo{
		Time:        this.TList[0].T,
		BetArr:      G_DbGetGameServerData.GameConfig.TenChips, //下注筹码
		PlaceBetAll: this.PlaceBet,                             //区域下注，0到3分别是黑红梅方
		HandCards:   this.HandCards,
		Players:     p.GetDeskPlayerInfo(),
	}
	roomInfo.PlayerInfos = PlayerInfo{
		Account:    p.Account,
		Uid:        p.Uid,
		Head:       p.Head,
		Coin:       p.Coins,
		IsDouble:   p.IsDouble,   //是否翻倍
		PlaceBet:   p.PlaceBet,   //自己区域下注，0到3分别是黑红梅方
		BetArrAble: p.BetArrAble, //可下注筹码
	}
	p.SendNativeMsg(MSG_GAME_RECONNECT_TABLE_REPLY, roomInfo)
}

//掉线信息
func (this *ExtDesk) HandleDisConnect(p *ExtPlayer, d *DkInMsg) {
	if !p.IsBet {
		this.LeaveByForce(p)
	} else {
		p.LiXian = true
	}
}

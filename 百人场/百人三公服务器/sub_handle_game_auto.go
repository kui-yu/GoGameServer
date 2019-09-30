package main

func (this *ExtDesk) PlayerIn(p *ExtPlayer, m *DkInMsg) {
	// 初始化
	p.init()
	// 发送匹配成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})
	//返送桌子信息
	this.SendDeskMSG(p)
}

//请求桌子信息
func (this *ExtDesk) SendDeskMSG(p *ExtPlayer) {
	p.LiXian = false
	//桌面玩家更新
	this.UadatePlayer(6)
	//发送房间信息给刚进场的玩家
	roomInfo := &GSPlayerIn{
		Id:         MSG_GAME_INFO_DESKINFO_REPLAY,
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
	p.SendNativeMsg(MSG_GAME_INFO_DESKINFO_REPLAY, roomInfo)
	//桌面玩家更新
	if len(this.ManyPlayer) >= 6 {
		return
	}
	this.UadatePlayer(6)
	//发送桌面玩家更新
	for _, v := range this.Players {
		if v.Uid == p.Uid {
			continue
		}
		v.SendNativeMsg(MSG_GAME_INFO_DESKPLAYER_REPLAY, &GSManyPlayer{
			Id:        MSG_GAME_INFO_DESKPLAYER_REPLAY,
			Players:   v.GetDeskPlayerInfo(), //获取赢金币最多的6个玩家
			AllPlayer: len(this.Players),
			JuHao:     this.JuHao,
		})
	}
}

package main

//重连消息
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {

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
		var coin int64
		if this.TableConfig.GameModule == 2 {
			coin = v.Coins
		} else {
			coin = v.TotalCoins
		}
		seat := GSSeatInfo{
			Uid:     v.Uid,
			Nick:    v.Nick,
			Cid:     v.ChairId,
			Sex:     v.Sex,
			Head:    v.Head,
			Lv:      v.Lv,
			Coin:    coin,
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
	if len(this.TList) != 0 {
		reconnectInfo.StageTime = this.TList[0].T
	}
	reconnectInfo.Id = MSG_GAME_INFO_RECONNECT
	reconnectInfo.Round = this.Round
	reconnectInfo.DisPlayer = this.DisPlayer
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT, &reconnectInfo)
	//在线状态更新
	p.LiXian = false
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 1,
	})
}

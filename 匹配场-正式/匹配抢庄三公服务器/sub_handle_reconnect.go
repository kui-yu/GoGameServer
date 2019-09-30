package main

import (
	"logs"
)

func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	if this.GameState == GAME_STATUS_FREE {
		p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
			Id:     MSG_GAME_RECONNECT_REPLY,
			Result: 1,
			Err:    "本桌子没有正在的游戏",
		})
		return
	}
	//返回重连成功
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:       MSG_GAME_RECONNECT_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})
	//发送座位玩家信息
	seatInfos := GPlayerSeatInfos{
		Id: MSG_GAME_INFO_SEAT,
	}
	for _, l := range this.Players {
		seatInfo := SeatInfo{
			Head:    l.Head,
			Coins:   l.Coins,
			Name:    l.Nick,
			Uid:     l.Uid,
			ChairId: l.ChairId,
		}
		if p.Uid != l.Uid {
			seatInfo.Name = "***" + l.Nick[len(p.Nick)-4:]
		}
		seatInfos.Data = append(seatInfos.Data, seatInfo)
	}
	p.SendNativeMsg(MSG_GAME_INFO_SEAT, seatInfos)
	//房间信息
	this.BScore = int64(G_DbGetGameServerData.Bscore) //设置底注
	//发送房间信息
	if this.JuHao == "" { //设置局号
		this.JuHao = GetJuHao()
	}
	p.SendNativeMsg(MSG_GAME_INFO_ROOM, GRoomInfo{
		Id:             MSG_GAME_INFO_ROOM,
		RoomNumber:     this.JuHao,
		MaxMultiple:    gameConfig.Max_Multiple,
		BScore:         this.BScore,
		PlayerMultiple: gameConfig.Idle_Multiple,
	})
	//玩家在线通知
	this.BroadcastOther(p, MSG_GAME_ONLINE_NOTIFY, GOnLineNotify{
		Id:    MSG_GAME_ONLINE_NOTIFY,
		Cid:   p.ChairId,
		State: 1,
	})
	//桌子信息
	info := GRDeskInfo{
		Id:              MSG_GAME_INFO_RECONNECT_DESK_INFO,
		Stage:           this.GameState,
		ChoiceBankerArr: make([]int, len(this.Players)),
		PlayerCards:     make([]PlayerCard, len(this.Players)),
		BankerChairId:   -1,
	}
	if len(this.TList) > 0 {
		info.StageTime = this.TList[0].T
	}
	switch this.GameState {
	case 13:
		info.BankerChairId = this.DeskBankerInfos.BankerId
		for k, v := range this.Players {
			info.ChoiceBankerArr[k] = v.BankerInfos.IsChoice
		}
	case 14:
		info.BankerChairId = this.DeskBankerInfos.BankerId
		//下注信息
		idleBet := IdleBet{}
		for _, v := range this.Players {
			if v.BankerInfos.IsBanker {
				continue
			}
			idleBet.Uid = v.Uid
			idleBet.Coins = v.BankerInfos.Multiple
			info.IdleBets = append(info.IdleBets, idleBet)
		}
	case 15:
		info.BankerChairId = this.DeskBankerInfos.BankerId
		//下注信息
		idleBet := IdleBet{}
		for _, v := range this.Players {
			if v.BankerInfos.IsBanker {
				continue
			}
			idleBet.Uid = v.Uid
			idleBet.Coins = v.BankerInfos.Multiple
			info.IdleBets = append(info.IdleBets, idleBet)
		}
		//牌信息
		playerCard := PlayerCard{}
		for k, v := range this.Players {
			if !v.IsOpenCards {
				continue
			}
			playerCard.Uid = v.Uid
			playerCard.Cards = v.HandCards
			info.PlayerCards[k] = playerCard
		}
	case 16:
		info.BankerChairId = this.DeskBankerInfos.BankerId
		//下注信息
		idleBet := IdleBet{}
		for _, v := range this.Players {
			if v.BankerInfos.IsBanker {
				continue
			}
			idleBet.Uid = v.Uid
			idleBet.Coins = v.BankerInfos.Multiple
			info.IdleBets = append(info.IdleBets, idleBet)
		}
		//牌信息
		playerCard := PlayerCard{}
		for k, v := range this.Players {
			if !v.IsOpenCards {
				continue
			}
			playerCard.Uid = v.Uid
			playerCard.Cards = v.HandCards
			info.PlayerCards[k] = playerCard
		}
		info.Settle = this.SettleInfo
	}
	logs.Debug("******用户正在进行断线重连：%v;%v", p.Nick, MSG_GAME_INFO_RECONNECT_DESK_INFO)
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT_DESK_INFO, info)
}

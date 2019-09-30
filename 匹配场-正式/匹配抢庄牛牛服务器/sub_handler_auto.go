package main

//玩家匹配
func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {
	// logs.Debug("玩家匹配")
	if this.GameState != GAME_STATUS_FREE {
		return
	}

	var playerCnt int = 0
	//机器人不能超过4个
	for _, v := range this.Players {
		if !v.Robot {
			playerCnt++
		}
	}
	if playerCnt < (GCONFIG.PlayerNum-this.MaxRobot) && p.Robot {
		//发送匹配成功
		p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
			Id:     MSG_GAME_AUTO_REPLY,
			Result: 13,
		})
		//踢出
		this.LeaveByForce(p)
		return
	}

	//发送匹配成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
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

	//发送房间信息
	if this.JuHao == "" {
		this.JuHao = GetJuHao()
		this.Bscore = G_DbGetGameServerData.Bscore
		this.Rate = G_DbGetGameServerData.Rate
	}

	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &GTableInfoReply{
		Id:      MSG_GAME_INFO_ROOM_NOTIFY,
		TableId: this.JuHao,
		BScore:  this.Bscore,
	})

	//人满开局
	if len(this.Players) >= GCONFIG.PlayerNum {
		this.nextStage(GAME_STATUS_START)
	}
}

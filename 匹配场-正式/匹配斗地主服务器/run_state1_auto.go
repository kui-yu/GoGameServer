package main

// 玩家匹配进入游戏
// 成功 进入发牌阶段 run_state_dealPoker
func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {

	//加入成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})
	//群发用户信息
	for _, v := range this.Players {
		gameReply := GInfoAutoGameReply{
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
				seat.Nick = "***" + seat.Nick[len(seat.Nick)-4:] //不是自己的玩家的名称将会被隐藏
			}
			gameReply.Seat = append(gameReply.Seat, seat)
		}
		v.SendNativeMsg(MSG_GAME_INFO_AUTO_REPLY, &gameReply)
	}
	//发送房间信息
	this.MaxDouble = int32(G_DbGetGameServerData.MaxTimes)
	this.Bscore = int32(G_DbGetGameServerData.Bscore)
	if this.JuHao == "" { //如果局号为空
		this.JuHao = GetJuHao()
	}
	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &GGameInfoNotify{
		Id:     MSG_GAME_INFO_ROOM_NOTIFY,
		BScore: this.Bscore,
		MaxBei: this.MaxDouble,
		JuHao:  this.JuHao,
	})

	//判断人员是否已满，开启游戏
	if len(this.Players) >= GCONFIG.PlayerNum {

		this.GameState = GAME_STATUS_START
		this.BroadStageTime(TIMER_START_NUM)
		//发送游戏开始通知
		sd := GGameStartNotify{
			Id: MSG_GAME_START,
		}
		this.BroadcastAll(MSG_GAME_START, &sd)

		//直接进入发牌
		this.TimerSendCard(nil)
	}
}

// 重写，添加玩家
func (this *ExtDesk) AddPlayer(p *ExtPlayer) int {
	var playerCnt int = 0
	//机器人不能超过4个
	for _, v := range this.Players {
		if !v.Robot {
			playerCnt++
		}
	}

	if p.Robot {
		var leave bool
		if playerCnt < (GCONFIG.PlayerNum - 2) {
			//没有真实玩家，机器人不加人
			leave = true
		} else if this.NotAllowRobotInRoom {
			//按时间不让机器人进入
			leave = true
		}

		if leave {
			//踢出
			p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
				Id:     MSG_GAME_AUTO_REPLY,
				Result: 13,
			})
			// this.LeaveByForce(p)
			return -2
		}

		this.NotAllowRobotInRoom = true
		//进入倒计时4s
		this.AddTimer(10, 2, this.TimerRobotInRoom, nil)
	}

	//正常进入
	if len(this.Players) >= GCONFIG.PlayerNum {
		return -1
	}
	//设置chairid
	doinsert := false
	for i, v := range this.Players {
		if i != int(v.ChairId) {
			doinsert = true
			p.ChairId = int32(i)
			nl := append([]*ExtPlayer{}, this.Players[:i]...)
			nl = append(nl, p)
			nl = append(nl, this.Players[i:]...)
			this.Players = nl
			break
		}
	}
	if !doinsert {
		p.ChairId = int32(len(this.Players))
		this.Players = append(this.Players, p)
	}
	//
	return len(this.Players)
}

//允许 机器人进入房间
func (this *ExtDesk) TimerRobotInRoom(d interface{}) {
	this.NotAllowRobotInRoom = false
}

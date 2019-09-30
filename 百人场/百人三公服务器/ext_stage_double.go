package main

func (this *ExtDesk) MultipleChoice(d interface{}) {
	//判断金币不足变会平倍
	for _, v := range this.Players {
		if v.IsDouble {
			v.GetBetArr(true, 0)
			if v.IsChangeDouble() {
				v.IsDouble = false
				v.SendNativeMsg(MSG_GAME_INFO_CHANGE_DOUBLE_REPLAY, &GSChangedouble{
					Id:         MSG_GAME_INFO_CHANGE_DOUBLE_REPLAY,
					BetArrAble: v.BetArrAble,
				})
			}
		}
	}
	this.Stage = STAGE_GAME_DOUBLE
	this.BroadStageTime(gameConfigInfo.Double_Timer)
	this.runTimer(gameConfigInfo.Double_Timer, this.GameStart)
}
func (this *ExtDesk) ThreeTimesLeave() {
	for _, v := range this.Players {
		if v.NotBet == 3 {
			v.SendNativeMsg(MSG_GAME_INFO_THREE_THMES, &struct {
				Id     int
				ErrStr string
				Code   int
			}{
				Id:     MSG_GAME_INFO_THREE_THMES,
				ErrStr: "三局没下注提示",
				Code:   2,
			})
		} else if v.NotBet >= 5 {
			v.SendNativeMsg(MSG_GAME_INFO_FIVE_THMES_LEAVE, &struct {
				Id int
			}{
				Id: MSG_GAME_INFO_FIVE_THMES_LEAVE,
			})
			v.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
				Id:      MSG_GAME_LEAVE_REPLY,
				Result:  0,
				Cid:     v.ChairId,
				Uid:     v.Uid,
				Token:   v.Token,
				Robot:   v.Robot,
				NoToCli: true,
			})
			this.DelPlayer(v.Uid)
			this.DeskMgr.LeaveDo(v.Uid)
		}
	}
}

//桌面玩家更新
func (this *ExtDesk) UadatePlayer(count int) {
	this.ManyPlayer = make([]ManyPlayer, 0)
	if len(this.Players) == 0 {
		return
	}
	for _, v := range this.Players {
		if len(this.Players) == count && v.AccumulateCoins == 0 { //人数多于要求的人数后不添加0输赢的玩家
			continue
		}
		this.ManyPlayer = append(this.ManyPlayer, ManyPlayer{
			Head:    v.Head,
			Account: v.Account,
			Uid:     v.Uid,
			Coins:   v.Coins,
			//Round:           v.Round,
			Round:           20, //20局记录
			AccumulateBet:   v.AccumulateBet,
			AccumulateCoins: v.AccumulateCoins,
		})
	}
}

//发送桌面玩家信息
func (this *ExtPlayer) GetDeskPlayerInfo() []ManyPlayer {
	manyPlayer := append([]ManyPlayer{}, this.Dk.ManyPlayer...)
	manyPlayer = SortForPlayers(manyPlayer, 6) //获取赢钱最多的前6个
	for k, v := range manyPlayer {             //其他人的名字要隐藏
		if this.Account != v.Account {
			if len(manyPlayer[k].Account) > 4 {
				manyPlayer[k].Account = "***" + manyPlayer[k].Account[len(manyPlayer[k].Account)-4:]
			}
		}
	}
	return manyPlayer
}
func (this *ExtDesk) IsUpdate(p *ExtPlayer) {
	//更新桌面玩家
	isUpdate := false
	for _, v := range this.ManyPlayer {
		if v.Uid == p.Uid {
			isUpdate = true
		}
	}
	if isUpdate {
		//桌面玩家更新
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
}

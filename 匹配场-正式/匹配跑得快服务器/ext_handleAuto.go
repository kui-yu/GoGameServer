package main

import "logs"

func (this *ExtDesk) HandleAuto(p *ExtPlayer, d *DkInMsg) {

	//加入成功	logs.Debug("接收到玩家匹配请求", p)
	logs.Debug("...............11111", this.DeskMgr.MapPlayers[p.Uid])
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})
	// if p.Nick == "ahao101" {
	// 	p.Coins = 20000
	// }

	//群发用户信息
	for _, v := range this.Players {
		gameReply := GInfoAutoGameReply{
			Id:        MSG_GAME_INFO_AUTO_REPALY,
			GameState: this.GameState,
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
				seat.Nick = "****" + seat.Nick[len(seat.Nick)-4:]
			}
			gameReply.Seat = append(gameReply.Seat, seat)
		}
		v.SendNativeMsg(MSG_GAME_INFO_AUTO_REPALY, &gameReply) //因为每个玩家得到的座位信息不同，所以需要要意义赋值
	}
	//发送房间信息
	this.Bscore = G_DbGetGameServerData.Bscore
	// this.Bscore = 1000
	logs.Debug("底分是:::::", this.Bscore)
	if this.JuHao == "" { //如果局号为空，为局号赋值
		this.JuHao = GetJuHao()
	}

	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, GGameInfoNotify{
		Id:     MSG_GAME_INFO_ROOM_NOTIFY,
		Bscore: this.Bscore,
		JuHao:  this.JuHao,
	})
	p.SendNativeMsg(MSG_GAME_INFO_STAGE, GStageInfo{
		Id:        MSG_GAME_INFO_STAGE,
		Stage:     0,
		StageTime: 0,
	})
	//判断人员是否满员，如果满员，进入游戏
	if len(this.Players) >= GCONFIG.PlayerNum {
		logs.Debug("玩家人数已经凑齐，可以开始游戏!!")
		//初始化玩家部分信息
		for _, v := range this.Players {
			v.WinForMap = make(map[int32]int)
			v.LoseForMap = make(map[int32]int)
		}
		//发送游戏开始通知
		this.GameState = GAME_STATUS_SENDCAR
		this.BroadStageTime(GAME_STATUS_SENDCAR_TIME)
		//直接进入发牌阶段
		this.SendCards()
	}
}

func (this *ExtDesk) SendCards() {
	logs.Debug("已经进入发牌阶段")
	//洗牌
	this.CardMgr.Shuffle()
	//发牌16张
	for _, v := range this.Players {
		v.HandCards = []byte{}
		v.HandCards = this.CardMgr.SendHandCard(16)
	}
	//发完牌后从玩家手牌中查找黑桃3玩家，将其设置为下一个操作玩家，并将其椅子ID发送到客户端
	this.SetH3()
	//发送牌消息
	sd := GGameSendCardNotify{
		Id:  MSG_GAME_INFO_SENDCARD,
		Cid: this.GetH3,
	}
	//因为每一个玩家的手牌都是不一样的，所以我们需要分开来发送手牌信息
	for _, v := range this.Players {
		for _, h := range v.HandCards {
			sd.HandCards = append(sd.HandCards, int(h))
		}
		logs.Debug("玩家", v.Nick, "的手牌是:", Sort(v.HandCards))
		if v.ChairId == this.CurCid {
			tishi, _, _ := this.CanOutCards(v)
			for _, v := range tishi {
				list := []int{}
				for _, v1 := range v {
					list = append(list, int(v1))
				}
				sd.Hint = append(sd.Hint, list)
			}
		}
		v.SendNativeMsg(MSG_GAME_INFO_SENDCARD, &sd)
		sd.HandCards = []int{} //初始化，否则下一个玩家的手牌将会与上一个玩家的手牌进行叠加
	}
	//定时器进入出牌阶段
	this.AddTimer(GAME_STATUS_SENDCAR, GAME_STATUS_SENDCAR_TIME, this.OutCard, nil)
}
func (this *ExtDesk) OutCard(d interface{}) {
	logs.Debug("进入出牌阶段")
	this.GameState = GAME_STATUS_OUTCARD
	this.BroadStageTime(GAME_STATUS_OUTCARD_TIME)
	//定时器进入出牌阶段
	this.AddTimer(GAME_STATUS_OUTCARD, GAME_STATUS_OUTCARD_TIME, this.TimerOutCard, nil)
}

//设置 黑桃三持有者
func (this *ExtDesk) SetH3() {
	//循环玩家手牌 查找黑桃3
	for _, v := range this.Players {
		for _, hc := range v.HandCards {
			if hc == Card_Hei_3 {
				this.GetH3 = v.ChairId
				this.CurCid = this.GetH3
				return
			}
		}
	}
}

// 重写，添加玩家 （由底层调用)
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
			//没有真实玩家，机器人不加入
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

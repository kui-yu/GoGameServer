package main

// import "logs"

func (this *ExtDesk) GameStateOperation(filter bool) {
	//接受玩家消息，广播消息。桌子属性，指针，下注信息
	if filter {
	} else {
		// logs.Debug("当前正在游戏桌子数：", len(this.DeskMgr.All)-len(this.DeskMgr.Free))
		this.BroadStageTime(STAGE_PLAY_OPERATION_TIME)
	}
	if this.Round == 0 {
		this.ChooseCallPlayer()
	}
	this.runTimer(STAGE_PLAY_OPERATION_TIME, this.HandleGameOperation)
}

//指定优先叫牌(\/) //阶段操作
func (this *ExtDesk) ChooseCallPlayer() {
	// fmt.Println("指定玩家优先叫牌")
	list := []int{}
	for i := 0; i < len(this.Players); i++ {
		list = append(list, i)
		this.ChairList = append(this.ChairList, this.Players[i].ChairId)
	}

	//随机一个庄家
	tempList := ListShuffle(list)
	//取第一个
	this.CallPlayer = this.Players[tempList[0]].ChairId
	this.Pround = tempList[0]
	this.Round = 1

	//谁叫牌
	bankerInfo := GSPlayerCallPlayer{
		Id:         MSG_GAME_INFO_CALLPLAYER_REPLY,
		Player:     this.CallPlayer,
		Round:      this.Round,
		MinCoin:    this.MinCoin,
		CoinEnough: IsCoinEnough(this.Players[tempList[0]].Coins, this.Players[tempList[0]].PayCoin, this.Bscore, this.MinCoin, this.Players[tempList[0]].CardType),
	}
	// fmt.Println(this.CallPlayer, "叫牌玩家")

	this.BroadcastAll(MSG_GAME_INFO_CALLPLAYER_REPLY, &bankerInfo)
	//首次机器人操作
	if this.Players[tempList[0]].Robot {
		//随机秒数
		//跳转机器人操作
		this.AddTimer(11, 3, this.RobotPlay, this.Players[tempList[0]])
	}
}

//阶段-操作
func (this *ExtDesk) HandleGameOperation(d interface{}) {
	// logs.Debug("阶段-操作")
	//玩家如果没有操作，大于顺子跟注，不然弃牌

	next := true
	for _, v := range this.Players {
		if this.CallPlayer == v.ChairId {
			if v.ProtectGU == 1 {
				// fmt.Println("防弃牌跟注")
				next = false
				this.GetGamePlay(4, v)
			} else if v.CardType == 2 {
				return
			} else if v.CardType == 1 {
				if v.CardLv >= CARD_SHUNZI { //跟注
					next = false
					this.GetGamePlay(4, v)
				} else {
					v.CardType = 2
					v.IsGU = true
					msg := GSCardType{
						Id:      MSG_GAME_INFO_GIVE_UP,
						ChairId: v.ChairId,
					}
					this.BroadcastAll(MSG_GAME_INFO_GIVE_UP, &msg)
				}
			} else {
				v.CardType = 2
				v.IsGU = true
				msg := GSCardType{
					Id:      MSG_GAME_INFO_GIVE_UP,
					ChairId: v.ChairId,
				}
				this.BroadcastAll(MSG_GAME_INFO_GIVE_UP, &msg)
			}
			break
		}
	}

	if next {
		this.MsgCallPlayer()
	}
}

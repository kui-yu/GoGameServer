package main

import (
	"encoding/json"
	// "bl.com/util"
)

func (this *ExtDesk) HandleTimeOutBet(d interface{}) {
	this.Lock()
	defer this.Unlock()

	if this.GameState != MSG_GAME_INFO_BET_NOTIFY {
		return
	}

	if this.NewBet {
		msg := GGameNewBetNotify{
			Id: MSG_GAME_INFO_BET_NEW_NOTIFY,
		}

		// 座位玩家下注信息
		// msg.SeatBetList = this.SeatMgr.GetSeatNewBetList()

		// 新下注，排除座位玩家
		// otherNewBetList := this.SeatMgr.GetOtherNewBetList()

		// 总下注
		msg.TAreaCoins = this.GetAreaCoinsList()

		for _, p := range this.Players {
			// 玩家新下注
			// newBetList := p.GetNewBetList()
			// 下注归集
			// p.ColAreaCoins()
			// 玩家总下注
			msg.PAreaCoins = p.GetNTAreaCoinsList()
			msg.OtherBetList = this.SeatMgr.GetOtherNewBetList2(p.Uid)
			// if this.SeatMgr.IsOnSeat(p) {
			// 	msg.OtherBetList = otherNewBetList
			// } else {
			// 	msg.OtherBetList = util.LessInt64List(otherNewBetList, newBetList)
			// }

			p.SendNativeMsg(MSG_GAME_INFO_BET_NEW_NOTIFY, msg)
		}

		for _, user := range this.Players {
			user.ColAreaCoins()
		}

		this.NewBet = false
	}

	if this.GameState == MSG_GAME_INFO_BET_NOTIFY {
		this.AddTimer(gameConfig.Timer.NewBet, gameConfig.Timer.NewBetNum, this.HandleTimeOutBet, nil)
	}
}

func (this *ExtDesk) HandleGameBet(p *ExtPlayer, d *DkInMsg) {
	// 同一个玩家不允许多线程操作
	this.RLock()
	p.Lock()

	defer this.RUnlock()
	defer p.Unlock()

	// 非下注状态 不处理
	if this.GameState != MSG_GAME_INFO_BET_NOTIFY {
		return
	}
	// if p.Robot {
	// 	//判断机器人的金币是不是最低或最高
	// 	if p.GetCoins() <= int64(G_DbGetGameServerData.Restrict) || p.GetCoins() > int64(G_DbGetGameServerData.LimitHigh) {
	// 		this.SeatMgr.DelPlayer(p)
	// 		this.LeaveByForce(p)
	// 	}
	// }
	// 下注信息检查
	data := GGameBet{}
	json.Unmarshal([]byte(d.Data), &data)
	msg := GGameBetReply{
		Id:     MSG_GAME_INFO_BET_REPLY,
		AreaId: data.AreaId,
		CoinId: data.CoinId,
	}

	// 下注区域检测
	switch int(data.AreaId) {
	// 闲、庄、和
	case INDEX_IDLE:
	case INDEX_BANKER:
	case INDEX_DRAW:

	// 小、大
	case INDEX_SMALL:
	case INDEX_BIG:

	// 闲对、庄对
	case INDEX_IDLEPAIR:
	case INDEX_BANKERPAIR:

	// 庄赢、庄输
	case INDEX_BANKERWIN:
	case INDEX_BANKERLOSE:
	default:
		msg.MsgId = ERR_AREAID
		goto end
	}

	if !this.BetArea[data.AreaId-1] {
		msg.MsgId = ERR_AREAID
		goto end
	}
	//if GetCostType() == 1 { //如果不是体验场再进行金币限制
	// 下注金币检测
	data.CoinId -= 1
	if data.CoinId < 0 || data.CoinId >= int32(len(this.BetList)) || p.Coins < this.BetList[data.CoinId] {
		msg.MsgId = ERR_COINID
		goto end
	}
	//} else {
	//data.CoinId -= 1
	//}

	if data.MsgId == p.GetMsgId()+1 {
		area := int(data.AreaId) - 1
		// 添加下注金额
		Coins := p.GetNewAreaCoin(area) + p.GetTotAreaCoin(area) + this.BetList[data.CoinId]
		// 区域上限检测
		if Coins > this.GameLimit.High || Coins < this.GameLimit.Low {
			msg.MsgId = ERR_LIMITCOIN
			goto end
		}

		// 下注
		p.AddNewAreaCoins(area, this.BetList[data.CoinId])

		// 桌子区域添加金币
		this.AddAreaCoins(area, this.BetList[data.CoinId])
		if !p.Robot {
			this.AddUserAreaCoins(area, this.BetList[data.CoinId])
		}

		p.SetMsgId(data.MsgId)
	}

	msg.MsgId = p.GetMsgId()

end:
	msg.PAreaCoins = p.GetNTAreaCoinsList()
	msg.Coins = p.Coins
	p.IsBet = true //赋值有下注
	p.SendNativeMsg(MSG_GAME_INFO_BET_REPLY, msg)
}

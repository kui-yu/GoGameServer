package main

import (
	"encoding/json"
	"logs"
)

//处理下注超时
func (this *ExtDesk) HandleTimeOutBet(d interface{}) {
	this.Lock()
	defer this.Unlock()
	//当游戏不是处于下注状态时
	if this.GameState != MSG_GAME_INFO_BET_NOTIFY {
		return
	}

	//	如果是新下注
	if this.NewBet {
		msg := GGameNewBetNotify{
			Id: MSG_GAME_INFO_BET_NEW_NOTIFY,
		}
		// //座位玩家下注信息
		// msg.SeatBetList = this.SeatMgr.GetSeatNewBetList()
		// //	新下注，排除座位玩家
		// otherNewBetList := this.SeatMgr.GetOtherNewBetList()

		//	总下注
		msg.TAreaCoins = this.GetAreaCoinsList()

		for _, p := range this.Players {
			//	玩家新下注
			// newBetList := p.GetNewBetList()
			//下注归集
			// p.ColAreaCoins()
			msg.PAreaCoins = p.GetNTAreaCoinsList()
			msg.OtherBetList = this.SeatMgr.GetOtherNewBetList2(p.Uid)
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

// //处理游戏下注操作
func (this *ExtDesk) HandleGameBet(p *ExtPlayer, d *DkInMsg) {
	//非下注状态 不处理
	if this.GameState != MSG_GAME_INFO_BET_NOTIFY {
		return
	}
	if p.Robot {
		//判断机器人的金币是不是最低
		if p.GetCoins() <= int64(G_DbGetGameServerData.Restrict) {
			this.SeatMgr.DelPlayer(p)
			this.LeaveByForce(p)
		}
	}
	//下注信息检查
	a := -1
	data := GGameBet{}
	json.Unmarshal([]byte(d.Data), &data)

	msg := GGameBetReply{
		Id:     MSG_GAME_INFO_BET_REPLY,
		AreaId: data.AreaId,
		CoinId: data.CoinId,
	}
	betlist := G_DbGetGameServerData.GameConfig.TenChips
	if data.CoinId < 0 || data.CoinId >= int32(len(betlist)) || p.Coins < betlist[data.CoinId] {
		logs.Debug("账户余额不足")
		msg.MsgId = ERR_COINID
		for i := len(betlist) - 1; i >= 0; i-- {
			if p.Coins >= betlist[i] {
				a = i
				break
			}
		}
		p.Limitcoinid = int32(a)
		msg.LimitCoinId = int32(a)
		goto end
	}

	//	下注区域检测
	switch int(data.AreaId) {
	case INDEX_RED:
	case INDEX_BLACK:
	case INDEX_LUCKYBLOW:
	default:
		msg.MsgId = ERR_AREAID
		goto end
	}
	//限制不能红黑都下注
	if data.AreaId == INDEX_RED {
		if p.PBetArea[INDEX_RED-1] == false {
			msg.MsgId = ERR_AREAID
			msg.LimitCoinId = p.Limitcoinid
			p.SendNativeMsg(MSG_GAME_INFO_BET_REPLY, msg)
			return
		}
	} else if data.AreaId == INDEX_BLACK {
		if p.PBetArea[INDEX_BLACK-1] == false {
			msg.MsgId = ERR_AREAID
			msg.LimitCoinId = p.Limitcoinid
			p.SendNativeMsg(MSG_GAME_INFO_BET_REPLY, msg)
			return
		}
	}
	////
	p.betArealist = append(p.betArealist, int(data.AreaId))
	if p.betArealist[len(p.betArealist)-1] == INDEX_RED {
		p.PBetArea[INDEX_BLACK-1] = false
	} else if p.betArealist[len(p.betArealist)-1] == INDEX_BLACK {
		p.PBetArea[INDEX_RED-1] = false
	}
	for i, _ := range p.PBetArea[:2] {
		if p.PBetArea[0] != p.PBetArea[1] {
			if p.PBetArea[i] {
				this.betId = int32(i)
			}
		} else {
			this.betId = 2
		}
	}
	// fmt.Println("下注ID：", this.betId)
	//下注金币检测
	// data.CoinId -= 1
	if data.MsgId == p.GetMsgId() {
		area := int(data.AreaId) - 1
		//	添加下注金额
		Coins := p.GetNewAreaCoin(area) + p.GetTotAreaCoin(area) + betlist[data.CoinId]
		//	区域上限检测
		if Coins > int64(G_DbGetGameServerData.GameConfig.LimitRedMax) {
			msg.MsgId = ERR_LIMITCOIN
			msg.LimitCoinId = p.Limitcoinid
			// logs.Debug("**************限制ID:%v", p.Limitcoinid)
			goto end
		}
		//下注
		p.AddNewAreaCoins(area, betlist[data.CoinId])
		//
		for i := len(betlist) - 1; i >= 0; i-- {
			if p.Coins >= betlist[i] {
				a = i
				break
			}
		}
		p.Limitcoinid = int32(a)
		msg.LimitCoinId = int32(a)
		//	添加桌子区域金币
		this.AddAreaCoins(area, betlist[data.CoinId])
		if !p.Robot {
			this.AddUserAreaCoins(area, betlist[data.CoinId])
		}
		p.SetMsgId(data.MsgId)
	}
	logs.Debug("玩家昵称:%v,GetNTAreaCoinsList:%v,", p.Nick, p.GetNTAreaCoinsList())
	msg.MsgId = p.GetMsgId()
end:
	msg.PAreaCoins = p.GetNTAreaCoinsList()
	msg.Coins = p.Coins
	p.SendNativeMsg(MSG_GAME_INFO_BET_REPLY, msg)
}

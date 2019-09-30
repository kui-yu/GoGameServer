package main

import (
	"encoding/json"
	"logs"
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
	logs.Debug("接收到玩家下注信息", p.Nick)
	// 同一个玩家不允许多线程操作
	this.RLock()
	p.Lock()

	defer this.RUnlock()
	defer p.Unlock()

	// 非下注状态 不处理
	if this.GameState != MSG_GAME_INFO_BET_NOTIFY {
		return
	}

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
	// 龙、虎、和
	case INDEX_DRAGON:
	case INDEX_TIGER:
	case INDEX_DRAW:

	// 龙   方、梅、红、黑
	case INDEX_DRAGONSPADE:
	case INDEX_DRAGONPLUM:
	case INDEX_DRAGONRED:
	case INDEX_DRAGONBLOCK:

	// 虎   方、梅、红、黑
	case INDEX_TIGERSPADE:
	case INDEX_TIGERPLUM:
	case INDEX_TIGERRED:
	case INDEX_TIGERBLOCK:

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

	// 下注金币检测
	data.CoinId -= 1
	if data.CoinId < 0 || data.CoinId >= int32(len(this.BetList)) || p.Coins < this.BetList[data.CoinId] {
		logs.Debug("发现账户没有钱")
		msg.MsgId = ERR_COINID
		goto end
	}
	logs.Debug("data.Msgid:::", data.MsgId)
	logs.Debug("p.MSGid", p.msgId)
	if data.MsgId == p.GetMsgId()+1 || true {
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
	p.SendNativeMsg(MSG_GAME_INFO_BET_REPLY, msg)
}

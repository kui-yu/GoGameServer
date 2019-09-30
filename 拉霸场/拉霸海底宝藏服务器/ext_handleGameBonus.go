package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) BonusStart(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("BonusStart")
	if !this.IsBonusStart && this.BonusCount >= int64(gameConfig.DeskInfo.OpenDimension) {
		this.AddTimer(TIMER_BONUS, TIMER_BONUS_NUM, this.BonusTimeout, p)
		this.IsBonusStart = true

		sd := GMsgBonusStartNotify{
			Id:        MSG_GAME_INFO_BONUS_START_NOTIFY,
			BoxCount:  this.BoxCount,
			MaxChoose: int64(gameConfig.DeskInfo.DimensionNum),
			Times:     TIMER_BONUS_NUM,
		}
		this.InitBonus()

		p.SendNativeMsg(MSG_GAME_INFO_BONUS_START_NOTIFY, sd)
	}
}

func (this *ExtDesk) BonusTimeout(d interface{}) {
	p := d.(*ExtPlayer)
	this.BonusEnd(p)
}

func (this *ExtDesk) BonusEnd(p *ExtPlayer) {
	if !this.IsBonusGame {
		return
	}

	if len(this.NormalCoins) != int(this.BoxCount)-gameConfig.DeskInfo.DimensionNum {
		for i := 0; i < gameConfig.DeskInfo.DimensionNum; i++ {
			if len(this.NormalCoins) == int(this.BoxCount)-gameConfig.DeskInfo.DimensionNum {
				break
			}
			msgBound := GMsgBonus{
				BoxIndex: int32(i),
			}
			m, _ := json.Marshal(msgBound)

			msg := DkInMsg{
				Uid:  p.Uid,
				Id:   MSG_GAME_INFO_BONUS,
				Data: string(m),
			}

			this.HandleGameBonus(p, &msg)
		}

		return
	}

	winCoins := int64(0)
	for i, v := range this.BoxCoins {
		if v == 0 {
			if len(this.NormalCoins) == 0 {
				continue
			}

			v = this.NormalCoins[0]

			this.BoxCoins[i] = int64(v)

			this.NormalWeight = append(this.NormalWeight[:0], this.NormalWeight[1:]...)
			this.NormalCoins = append(this.NormalCoins[:0], this.NormalCoins[1:]...)
		} else {
			winCoins += v
		}
	}

	sd := GMsgBonusEndNotify{
		Id:    MSG_GAME_INFO_BONUS_END,
		Bonus: append([]int64{}, this.BoxCoins...),
	}

	//发送结算消息给数据库, 简单记录
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.GameId,
		Mini:        true,
		SetLeave:    1, //是否设置离开，0离开，1不离开
		NoSaveCoin:  0,
	}

	betCoins := int64(0)
	dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
		UserId:      p.Uid,
		UserAccount: p.Account,
		BetCoins:    betCoins,
		ValidBet:    betCoins,
		PrizeCoins:  winCoins - betCoins,
		Robot:       p.Robot,
	})
	if GetCostType() == 1 {
		logs.Debug("发了数据库消息")
		p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
	}

	//发送消息给大厅去记录游戏记录
	rdreq := GGameRecord{
		Id:          MSG_GAME_END_RECORD,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.GameId,
	}
	rddata := GGameRecordInfo{
		UserId:      p.Uid,
		UserAccount: p.Account,
		BetCoins:    betCoins,
		PrizeCoins:  winCoins,
		CoinsAfter:  p.Coins,
		Robot:       p.Robot,
	}
	rddata.CoinsBefore = rddata.CoinsAfter - rddata.PrizeCoins + rddata.BetCoins
	rdreq.UserRecord = append(rdreq.UserRecord, rddata)
	if GetCostType() == 1 {
		logs.Debug("发送游戏大厅记录")
		p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
	}

	this.IsBonusGame = false
	this.IsBonusStart = false
	p.SendNativeMsg(MSG_GAME_INFO_BONUS_END, sd)

	this.ResetTable()
}

func (this *ExtDesk) HandleGameBonus(p *ExtPlayer, d *DkInMsg) {
	if !this.IsBonusGame || !this.IsBonusStart {
		return
	}

	if len(this.NormalCoins) == int(this.BoxCount)-gameConfig.DeskInfo.DimensionNum {
		return
	}

	//
	jmsg := GMsgBonus{}
	err := json.Unmarshal([]byte(d.Data), &jmsg)
	if err != nil {
		logs.Error("HandleGameBonus 消息错误:", err, d, p)
		return
	}

	if jmsg.BoxIndex >= int32(this.BoxCount) || jmsg.BoxIndex < 0 || this.BoxCoins[jmsg.BoxIndex] != 0 {
		logs.Error("HandleGameBonus 开启箱子错误：", jmsg.BoxIndex)
		return
	}

	var max int32 = 0
	for _, v := range this.NormalWeight {
		max = max + v
	}

	rand := RandInt64(int64(max)) + 1
	max = int32(rand)

	for i, v := range this.NormalWeight {
		max = max - v
		if len(this.NormalCoins) <= i {
			break
		}
		if max <= 0 {
			this.BoxCoins[jmsg.BoxIndex] = this.NormalCoins[i]

			this.NormalWeight = append(this.NormalWeight[:i], this.NormalWeight[i+1:]...)
			this.NormalCoins = append(this.NormalCoins[:i], this.NormalCoins[i+1:]...)
			break
		}
	}

	// 添加中奖金币
	p.Coins = p.Coins + this.BoxCoins[jmsg.BoxIndex]

	sd := GMsgBonusNotify{
		Id:       MSG_GAME_INFO_BONUS_NOTIFY,
		BoxIndex: jmsg.BoxIndex,
		Coins:    p.Coins,
		Bonus:    this.BoxCoins[jmsg.BoxIndex],
	}

	p.SendNativeMsg(MSG_GAME_INFO_BONUS_NOTIFY, sd)

	if len(this.NormalCoins) == int(this.BoxCount)-gameConfig.DeskInfo.DimensionNum {
		this.DelTimer(TIMER_BONUS)
		this.BonusEnd(p)
	}
}

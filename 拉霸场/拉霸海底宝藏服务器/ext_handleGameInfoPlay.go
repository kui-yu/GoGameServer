package main

import (
	"encoding/json"
	"fmt"
	"logs"
)

func (this *ExtDesk) HandleGameInfoPlay(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("开始扭蛋？")
	//
	jmsg := GMsgPlay{}
	err := json.Unmarshal([]byte(d.Data), &jmsg)
	if err != nil {
		logs.Error("HandleStart 消息错误:", err, d, p)
		return
	}

	if jmsg.Lines <= 0 || jmsg.Lines > 9 {
		logs.Error("HandleStart 线数错误", jmsg.Lines)
		return
	}

	if jmsg.Coins <= 0 || jmsg.Coins > 10 {
		logs.Error("HandleStart 分数错误", jmsg.Coins)
		return
	}
	if jmsg.Coins*jmsg.Lines*int64(G_DbGetGameServerData.Bscore) > p.Coins {
		logs.Debug("玩家金币不足")
		logs.Error("HandleStart 玩家金币不足", jmsg, p.Account, p.Coins)
		return
	}

	if this.BonusCount >= int64(gameConfig.DeskInfo.OpenDimension) {
		logs.Error("HandleStart 请触发奖金游戏", p.Account)
		return
	}

	if p.Coins >= int64(G_DbGetGameServerData.LimitHigh) {
		p.SendNativeMsg(MSG_GAME_INFO_EXIT_LIMIT_HIGHT, &GLeaveReply{
			Id: MSG_GAME_INFO_EXIT_LIMIT_HIGHT,
		})

		this.LeaveByForce(p)

		return
	}

	this.Times++
	this.TotLines += jmsg.Lines
	this.Lines = this.TotLines / this.Times
	if this.Lines < 1 {
		this.Lines = 1
	}

	var scenes []byte
	var lines [18][]byte
	var Pow int64
	var needShow []byte

	win := int64(gameConfig.DeskInfo.Win) / 9 * jmsg.Lines
	isWin := RandInt64(100) > win

	for i := 0; i < 20; i++ {
		// 随机15个图标, 中奖线, 中间倍率, 中奖图标位置
		scenes, lines, Pow, needShow = this.imgMgr.GetPlayResult(jmsg.Lines)

		realWin := jmsg.Coins*jmsg.Lines*int64(G_DbGetGameServerData.Bscore) > jmsg.Coins*Pow*int64(G_DbGetGameServerData.Bscore)
		if realWin == isWin || GetCostType() != 1 {
			break
		}
	}

	// 玩家金币 = 当前金币 - 下注金币 + 中奖金币
	p.Coins = p.Coins - jmsg.Coins*jmsg.Lines*int64(G_DbGetGameServerData.Bscore) + jmsg.Coins*Pow*int64(G_DbGetGameServerData.Bscore)

	// 计算Bonus数量
	BonusCount := this.imgMgr.GetBonusCount(scenes)
	if BonusCount != 0 {
		for i := 0; i < int(BonusCount); i++ {
			this.LineCount += int(jmsg.Lines)
		}
	}
	fmt.Println("bns", BonusCount)
	if BonusCount >= this.MinBonus {
		this.BonusCount = int64(gameConfig.DeskInfo.OpenDimension)
	} else {
		this.BonusCount = this.BonusCount + BonusCount
	}

	sd := GMsgPlayNotify{
		Id:         MSG_GAME_INFO_PLAY_NOTIFY,
		Scenes:     append([]byte{}, scenes...),
		Win:        jmsg.Coins * Pow * int64(G_DbGetGameServerData.Bscore),
		Coins:      p.Coins,
		IsShow:     append([]byte{}, needShow...),
		BonusCount: this.BonusCount,
		MinBonus:   int32(gameConfig.DeskInfo.OpenDimension),
	}

	for i, _ := range sd.Lines {
		sd.Lines[i] = append(sd.Lines[i], lines[i]...)
	}

	p.SendNativeMsg(MSG_GAME_INFO_PLAY_NOTIFY, sd)

	if this.BonusCount >= int64(gameConfig.DeskInfo.OpenDimension) {
		// // 累计图标触发
		// sd := GMsgBonusStartNotify{
		// 	Id:        MSG_GAME_INFO_BONUS_START_NOTIFY,
		// 	BoxCount:  this.BoxCount,
		// 	MaxChoose: int64(gameConfig.DeskInfo.DimensionNum),
		// }
		this.IsBonusStart = false

		// p.SendNativeMsg(MSG_GAME_INFO_BONUS_START_NOTIFY, sd)
	}

	//发送结算消息给数据库, 简单记录
	this.GameId = GetJuHao()
	dbreq := GGameEnd{
		Id:          MSG_GAME_END_NOTIFY,
		GameId:      GCONFIG.GameType,
		GradeId:     GCONFIG.GradeType,
		RoomId:      GCONFIG.RoomType,
		GameRoundNo: this.GameId,
		Mini:        false,
		SetLeave:    1, //是否设置离开，0离开，1不离开
	}

	betCoins := jmsg.Coins * jmsg.Lines * int64(G_DbGetGameServerData.Bscore)
	dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
		UserId:      p.Uid,
		UserAccount: p.Account,
		BetCoins:    betCoins,
		ValidBet:    betCoins,
		PrizeCoins:  sd.Win - betCoins,
		Robot:       p.Robot,
	})
	if GetCostType() == 1 {
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
		Scenes:      append([]byte{}, scenes...),
		Lines:       sd.Lines,
		Pow:         Pow,
		BetLines:    jmsg.Lines,
		Bcoins:      int64(G_DbGetGameServerData.Bscore),
		BetCoins:    betCoins,
		PrizeCoins:  sd.Win - betCoins,
		CoinsAfter:  p.Coins,
		Robot:       p.Robot,
	}
	rddata.CoinsBefore = rddata.CoinsAfter - rddata.PrizeCoins + rddata.BetCoins
	rdreq.UserRecord = append(rdreq.UserRecord, rddata)
	if GetCostType() == 1 {
		p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
	}
}

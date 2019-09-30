/**
* 玩家操作
**/

package main

import (
	"encoding/json"
)

type FsmUserOperate struct {
	Mark           int
	EDesk          *ExtDesk
	EndDateTime    int64
	CurrMinDownBet int64
	EnableOperate  bool
	OperateCode    int
	OperateData    interface{}
}

func (this *FsmUserOperate) InitFsm(mark int, extDesk *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDesk

	this.EnableOperate = false
	this.CurrMinDownBet = gameConfig.BigBlindCoin

	this.addListerner()
}

func (this *FsmUserOperate) GetMark() int {
	return this.Mark
}

func (this *FsmUserOperate) Run(upMark int, args ...interface{}) {
	DebugLog("状态 玩家操作")
	timerMs := gameConfig.GameTimer.UserOperate
	this.EndDateTime = GetTimeMS() + int64(timerMs)

	this.EDesk.SendDeskStatus(this.Mark, timerMs)
	this.EDesk.AddUniueTimer(TimerId, int(timerMs/1000), this.TimerCall, nil)

	// 重置初始金额
	if upMark == GameStatusHoleCards {
		this.CurrMinDownBet = gameConfig.BigBlindCoin
	} else if upMark != this.GetMark() {
		this.CurrMinDownBet = 0
	}

	// 状态逻辑
	operateAuth := OperateAuthQP

	sid := this.EDesk.SeatOperateId
	p := this.EDesk.GetPlayerFromSid(sid)
	limitMaxBet := this.EDesk.GetStageMaxBet() - this.EDesk.GetPlayerStageDownBet(p)
	limitMinBet := this.CurrMinDownBet - this.EDesk.GetPlayerStageDownBet(p)

	if p.CarryCoin <= limitMaxBet {
		operateAuth |= OperateAuthSH
	}

	if this.EDesk.IsOperateOpen {
		operateAuth |= OperateAuthKP
	}

	if p.CarryCoin >= limitMinBet {
		if !this.EDesk.IsOperateOpen {
			operateAuth |= OperateAuthGZ
		}

		if limitMinBet < limitMaxBet {
			operateAuth |= OperateAuthJZ
		}
	}
	if this.EDesk.IsOperateOpen {
		operateAuth |= OperateAuthKP
	}

	data := struct {
		Sid         int
		OperateAuth int
		MinCoin     int64
		MaxCoin     int64
	}{
		Sid:         sid,
		OperateAuth: operateAuth,
		MinCoin:     limitMinBet,
		MaxCoin:     limitMaxBet,
	}

	this.OperateData = data

	this.EnableOperate = true
	this.EDesk.SendNetMessage(MSG_GAME_NGameOperate, data)
}

func (this *FsmUserOperate) TimerCall(d interface{}) {
	sid := this.EDesk.SeatOperateId
	p := this.EDesk.GetPlayerFromSid(sid)

	p.IsFold = true
	p.StageOperate = OperateAuthQP

	this.EDesk.SendNetMessage(MSG_GAME_NGameOperateResult, struct {
		Sid         int
		OperateAuth int
		RaiseValue  int64
	}{
		Sid:         sid,
		OperateAuth: OperateAuthQP,
		RaiseValue:  0,
	})
	this.GotoFsm()
}

func (this *FsmUserOperate) Leave() {

}

func (this *FsmUserOperate) Reset() {
	this.EnableOperate = false
	this.CurrMinDownBet = gameConfig.BigBlindCoin
}

func (this *FsmUserOperate) GetRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FsmUserOperate) OnUserOnline(p *ExtPlayer) {
	if this.EnableOperate {
		p.SendNetMessage(MSG_GAME_NGameOperate, this.OperateData)
	}
}

func (this *FsmUserOperate) OnUserOffline(p *ExtPlayer) {
}

func (this *FsmUserOperate) addListerner() {
	// 玩家操作返回
	this.EDesk.Handle[MSG_GAME_QGameOperate] = this.RecvOperate
}

func (this *FsmUserOperate) RecvOperate(p *ExtPlayer, d *DkInMsg) {
	type Reply struct {
		Result int
		Err    string
	}
	if this.EnableOperate == false || this.EDesk.GetFSM(0).GetMark() != this.Mark {
		p.SendNetMessage(MSG_GAME_RGameOperate, Reply{
			Result: 1,
			Err:    "当前还不是下注状态，请稍后",
		})
		return
	}
	if this.EDesk.SeatOperateId != p.Sid {
		p.SendNetMessage(MSG_GAME_RGameOperate, Reply{
			Result: 2,
			Err:    "请等待其他人下注",
		})
		return
	}

	req := struct {
		OperateAuth int
		RaiseValue  int64
	}{}

	json.Unmarshal([]byte(d.Data), &req)
	//弃牌
	if (req.OperateAuth & OperateAuthQP) != 0 {
		p.IsFold = true
		req.RaiseValue = 0
	}
	//过牌
	if (req.OperateAuth & OperateAuthKP) != 0 {
		if this.EDesk.IsOperateOpen != true {
			p.SendNetMessage(MSG_GAME_RGameOperate, Reply{
				Result: 3,
				Err:    "过牌操作错误",
			})
			return
		}
		this.EDesk.AddPlayerDownBet(p, 0)
	}
	//allin
	if (req.OperateAuth & OperateAuthSH) != 0 {
		if p.CarryCoin > this.EDesk.GetStageMaxBet() {
			p.SendNetMessage(MSG_GAME_RGameOperate, Reply{
				Result: 4,
				Err:    "ALLIN操作错误",
			})
			return
		}

		if this.CurrMinDownBet < p.CarryCoin {
			this.CurrMinDownBet = p.CarryCoin
		}
		this.EDesk.AddPlayerDownBet(p, p.CarryCoin)
		p.AllInStage = this.EDesk.CurrStage
	}

	minBet := this.CurrMinDownBet - this.EDesk.GetPlayerStageDownBet(p)

	//跟注
	if (req.OperateAuth & OperateAuthGZ) != 0 {
		if p.CarryCoin < minBet {
			p.SendNetMessage(MSG_GAME_RGameOperate, Reply{
				Result: 5,
				Err:    "金币不足",
			})
			return
		}

		this.EDesk.AddPlayerDownBet(p, minBet)
	}
	//加注
	if (req.OperateAuth & OperateAuthJZ) != 0 {
		if p.CarryCoin < req.RaiseValue {
			p.SendNetMessage(MSG_GAME_RGameOperate, Reply{
				Result: 6,
				Err:    "金币不足",
			})
			return
		}

		if minBet >= req.RaiseValue {
			p.SendNetMessage(MSG_GAME_RGameOperate, Reply{
				Result: 7,
				Err:    "加注额不足",
			})
			return
		}

		if req.RaiseValue+this.EDesk.GetPlayerStageDownBet(p) > this.EDesk.GetStageMaxBet() {
			p.SendNetMessage(MSG_GAME_RGameOperate, Reply{
				Result: 8,
				Err:    "加注额太大",
			})
			return
		}

		this.EDesk.AddPlayerDownBet(p, req.RaiseValue)
		this.CurrMinDownBet = this.EDesk.GetPlayerStageDownBet(p)
	}

	if this.EDesk.IsOperateOpen == true && (req.OperateAuth == OperateAuthJZ ||
		req.OperateAuth == OperateAuthGZ || req.OperateAuth == OperateAuthSH) {
		this.EDesk.IsOperateOpen = false
	}

	p.SendNetMessage(MSG_GAME_RGameOperate, Reply{
		Result: 0,
		Err:    "",
	})

	p.StageOperate = req.OperateAuth

	this.EDesk.SendNetMessage(MSG_GAME_NGameOperateResult, struct {
		Sid         int
		OperateAuth int
		RaiseValue  int64
	}{
		Sid:         p.Sid,
		OperateAuth: req.OperateAuth,
		RaiseValue:  req.RaiseValue,
	})

	this.GotoFsm()
}

func (this *FsmUserOperate) GotoFsm() {
	this.EnableOperate = false
	// 等待跳转
	this.EDesk.AddUniueTimer(TimerId, 1, func(d interface{}) {
		nextP := this.FindNextOperatePlayer()

		if nextP != nil {
			this.EDesk.SeatOperateId = nextP.Sid
			this.EDesk.RunFSM(this.GetMark())
		} else {
			// 重置玩家的操作
			for _, p := range this.EDesk.Players {
				p.StageOperate = 0
			}
			// 跳转
			this.EDesk.RunFSM(this.EDesk.GetNextStageMark())
		}
	}, nil)
}

// 获得当前用户的当前状态下注额
func (this *FsmUserOperate) GetPlayerCurrStageDownBet(p *ExtPlayer) int64 {
	plen := len(p.DownCoins)
	stageIdx := this.EDesk.GetStageIdx(0xFF)
	if plen < stageIdx+1 {
		return -1
	}
	return p.DownCoins[stageIdx]
}

// 查找下一个操作的用户 返回下一个玩家
func (this *FsmUserOperate) FindNextOperatePlayer() *ExtPlayer {
	// 有效玩家只有一个，直接返回nil
	validNum := 0
	for _, p := range this.EDesk.Players {
		if p.State == UserStateGameIn && p.IsFold == false && p.AllInStage == 0xFF {
			validNum += 1
		}
	}
	if validNum == 1 {
		this.EDesk.IsExistOperateFsm = false
		return nil
	} else {
		this.EDesk.IsExistOperateFsm = true
	}

	firstSid, nextSid := this.EDesk.SeatOperateId, this.EDesk.SeatOperateId
	//查找下一个没有操作过的玩家
	for {
		nextP := this.EDesk.GetNextPlayer(nextSid)
		nextSid = nextP.Sid
		if nextSid == firstSid {
			nextSid = -1
			break
		}
		if nextP.State != UserStateGameIn {
			continue
		}

		if nextP.IsFold == true || nextP.AllInStage != 0xFF {
			continue
		}

		// 他没有下注过
		if this.GetPlayerCurrStageDownBet(nextP) == -1 {
			break
		}

		// 他下注的钱比其他人少
		if nextP.CarryCoin > 0 && this.GetPlayerCurrStageDownBet(nextP) < this.CurrMinDownBet {
			break
		}
	}
	if nextSid != -1 {
		return this.EDesk.GetPlayerFromSid(nextSid)
	}

	return nil
}

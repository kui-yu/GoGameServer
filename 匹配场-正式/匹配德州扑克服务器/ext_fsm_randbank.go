/**
* 随机庄家，下盲注
**/

package main

type FsmRandBank struct {
	Mark        int
	EDesk       *ExtDesk
	EndDateTime int64
}

func (this *FsmRandBank) InitFsm(mark int, extDesk *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDesk
}

func (this *FsmRandBank) GetMark() int {
	return this.Mark
}

func (this *FsmRandBank) Run(upMark int, args ...interface{}) {
	DebugLog("状态 随机庄家，下盲注")
	timerMs := gameConfig.GameTimer.RandBank
	this.EndDateTime = GetTimeMS() + int64(timerMs)

	this.EDesk.SendDeskStatus(this.Mark, timerMs)
	this.EDesk.AddUniueTimer(TimerId, int(timerMs/1000), this.TimerCall, nil)

	// 设置等待下局的玩家状态到游戏进行中
	for _, p := range this.EDesk.Players {
		if p.State == UserStateWaitStart {
			p.State = UserStateGameIn
		}

		// 通知玩家改变状态
		changeInfo := struct {
			Sid   int
			State int
		}{
			Sid:   p.Sid,
			State: p.State,
		}
		ndata := struct {
			ChangeType int
			UserInfo   interface{}
		}{
			ChangeType: 3,
			UserInfo:   changeInfo,
		}
		this.EDesk.SendNetMessage(MSG_GAME_NGameUserChange, ndata)
	}

	// 设置庄家位置
	this.SetBankSid()
	//盲注
	this.SetBlindBet()
}

func (this *FsmRandBank) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GameStatusHoleCards)
}

func (this *FsmRandBank) Leave() {

}

func (this *FsmRandBank) Reset() {
}

func (this *FsmRandBank) GetRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FsmRandBank) OnUserOnline(p *ExtPlayer) {
}

func (this *FsmRandBank) OnUserOffline(p *ExtPlayer) {
}

//设置庄家位置
func (this *FsmRandBank) SetBankSid() {
	var bankPlayer *ExtPlayer

	if this.EDesk.SeatBank == 0xFF {
		ulen := this.EDesk.GetUserStateTotal(UserStateGameIn)

		rank, _ := GetRandomNum(0, ulen)

		curr := 0
		for _, p := range this.EDesk.Players {
			if p.State == UserStateGameIn {
				if curr == rank {
					bankPlayer = p
					break
				}
				curr += 1
			}
		}
	} else {
		oldSBank := this.EDesk.SeatBank

		oldP := this.EDesk.GetPlayerFromSid(oldSBank)
		if oldP != nil {
			oldP.IsBank = false
		}

		nextSid := this.EDesk.GetNextPlayer(oldSBank, UserStateGameIn).Sid
		bankPlayer = this.EDesk.GetPlayerFromSid(nextSid)
	}

	bankPlayer.IsBank = true
	this.EDesk.SeatBank = bankPlayer.Sid
	this.EDesk.SeatOperateId = bankPlayer.Sid

	//发送庄家
	this.EDesk.SendNetMessage(MSG_GAME_NGameRandRank, struct {
		BankSid int
	}{
		BankSid: this.EDesk.SeatBank,
	})
}

//设置大小盲注
func (this *FsmRandBank) SetBlindBet() {
	p1 := this.EDesk.GetNextPlayer(0xFF)
	p2 := this.EDesk.GetNextPlayer(p1.Sid)
	this.EDesk.AddPlayerDownBet(p1, gameConfig.SmallBlindCoin)
	this.EDesk.AddPlayerDownBet(p2, gameConfig.BigBlindCoin)
	this.EDesk.SeatOperateId = this.EDesk.GetNextPlayer(p2.Sid).Sid

	type BlindCoinInfo struct {
		Sid  int
		Coin int64
	}

	blindInfo := struct {
		SmallBlindCoin BlindCoinInfo
		BigBlindCoin   BlindCoinInfo
	}{
		SmallBlindCoin: BlindCoinInfo{
			Sid:  p1.Sid,
			Coin: gameConfig.SmallBlindCoin,
		},
		BigBlindCoin: BlindCoinInfo{
			Sid:  p2.Sid,
			Coin: gameConfig.BigBlindCoin,
		},
	}
	//发送盲注
	this.EDesk.SendNetMessage(MSG_GAME_NGameBlind, blindInfo)
}

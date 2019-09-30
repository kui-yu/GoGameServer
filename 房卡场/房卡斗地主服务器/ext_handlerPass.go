package main

// "encoding/json"
// "logs"
//	"sort"

func (this *ExtDesk) HandlePass(p *ExtPlayer, d *DkInMsg) {
	// logs.Debug("玩家过", p.ChairId)
	if this.GameState != GAME_STATUS_PLAY {
		return
	}
	if this.CurCid != p.ChairId {
		return
	}
	if this.MaxChuPai == nil {
		return
	}
	if (p.ChairId+1)%int32(len(this.Players)) == this.MaxChuPai.Cid {
		this.MaxChuPai = nil
	}
	this.RdChuPai = append(this.RdChuPai, &GOutCard{
		Cid:  p.ChairId,
		Type: -1,
	})
	this.CurCid = (this.CurCid + 1) % int32(len(this.Players))
	//广播
	this.BroadcastAll(MSG_GAME_INFO_PASS_REPLY, &GGamePassReply{
		Id:      MSG_GAME_INFO_PASS_REPLY,
		Cid:     p.ChairId,
		NextCid: (this.CurCid) % int32(len(this.Players)),
	})
	this.DelTimer(TIMER_OUTCARD)
	nextplayer := this.Players[this.CurCid]
	if nextplayer.TuoGuan {
		this.AddTimer(TIMER_OUTCARD, 1, this.TuoGuanOut, nil)
	} else {
		this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCard, nil)
	}

}

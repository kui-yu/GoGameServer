package main

import (
	"encoding/json"
	// "logs"
)

//托管
func (this *ExtDesk) HandleTuoGuan(p *ExtPlayer, d *DkInMsg) {

	if this.GameState == GAME_STATUS_END || this.GameState == GAME_STATUS_FREE {
		return
	}
	re := GTuoGuan{}
	json.Unmarshal([]byte(d.Data), &re)
	//
	if (re.Ctl == 1 && !p.TuoGuan) || (re.Ctl == 2 && p.TuoGuan) {
		p.TuoGuan = !p.TuoGuan
		this.BroadcastAll(MSG_GAME_INFO_TUOGUAN_REPLY, &GTuoGuanReply{
			Id:     MSG_GAME_INFO_TUOGUAN_REPLY,
			Ctl:    re.Ctl,
			Result: 0,
			Cid:    p.ChairId,
		})
		if this.GameState == GAME_STATUS_PLAY && p.TuoGuan {
			if this.CurCid == p.ChairId {
				this.DelTimer(TIMER_OUTCARD) //删除出牌定时器
				this.AddTimer(TIMER_OUTCARD, 1, this.TimerOutCard, nil)
			}
		}
		return
	}

}

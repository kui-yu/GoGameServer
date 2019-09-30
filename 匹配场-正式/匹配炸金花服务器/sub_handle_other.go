package main

import "encoding/json"

//玩家属性操作
func (this *ExtDesk) HandlePlayWithSys(p *ExtPlayer, d *DkInMsg) {
	if this.GameState == GAME_STATUS_END {
		return
	}

	data := GAProtectGiveUp{}
	json.Unmarshal([]byte(d.Data), &data)

	if data.PAttribute == 1 { //自动跟注操作
		if p.AutoFollowUp == 0 {
			p.AutoFollowUp = 1

			if p.ChairId == this.CallPlayer {
				this.GetGamePlay(4, p)
			}
		} else {
			p.AutoFollowUp = 0
		}
		info := GSSystemOpertion{
			Id:         MSG_GAME_INFO_PLAY_WITH_SYS,
			PAttribute: data.PAttribute,
			OpSuccess:  p.AutoFollowUp,
		}
		p.SendNativeMsg(MSG_GAME_INFO_PLAY_WITH_SYS, &info)
	} else if data.PAttribute == 2 { //防弃牌
		if p.ProtectGU == 0 {
			p.ProtectGU = 1
		} else {
			p.ProtectGU = 0
		}
		info := GSSystemOpertion{
			Id:         MSG_GAME_INFO_PLAY_WITH_SYS,
			PAttribute: data.PAttribute,
			OpSuccess:  p.ProtectGU,
		}
		p.SendNativeMsg(MSG_GAME_INFO_PLAY_WITH_SYS, &info)
	}

}

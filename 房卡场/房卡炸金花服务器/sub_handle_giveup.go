package main

import "encoding/json"

//玩家弃牌(\/)
func (this *ExtDesk) HandleGiveUp(p *ExtPlayer, d *DkInMsg) {
	// fmt.Println("弃牌")
	if this.GameState == STAGE_SETTLE || this.GameState == GAME_STATUS_END || this.GameState == STAGE_INIT {
		return
	}
	data := GAPlayerOperation{}
	json.Unmarshal([]byte(d.Data), &data)

	if data.Operation == 0 {
		p.CardType = 2
		p.IsGU = true
		info := GSCardType{
			Id:      MSG_GAME_INFO_GIVE_UP,
			ChairId: p.ChairId,
		}
		// fmt.Println(info)
		this.BroadcastAll(MSG_GAME_INFO_GIVE_UP, &info)
	}

	count := 0
	for _, v := range this.Players {
		if v.CardType != 2 {
			count++
		}
	}
	if count == 1 { //进入结算阶段
		this.nextStage(GAME_STATUS_END)
		return
	}

	if p.ChairId == this.CallPlayer {
		this.MsgCallPlayer()
	}
}

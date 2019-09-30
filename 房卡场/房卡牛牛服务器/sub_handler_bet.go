package main

import (
	"encoding/json"
)

//玩家叫分
func (this *ExtDesk) HandleGameCall(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != STAGE_CALL_SCORE {
		return
	}
	//抢庄模式，庄家不叫分
	if this.TableConfig.GameType == 1 && p.ChairId == this.Banker {
		return
	}

	data := GACallMsg{}
	json.Unmarshal([]byte(d.Data), &data)
	// logs.Debug("玩家叫分", d)

	betFlag := true
	for _, bet := range p.PlayerBets {
		if bet == data.Multiple {
			betFlag = false
			break
		}
	}
	if betFlag {
		return
	}

	//倍数
	p.BetMultiple = data.Multiple
	//广播回复
	this.BroadcastAll(MSG_GAME_INFO_CALL_REPLY, &GSCallMsg{
		Id:       MSG_GAME_INFO_CALL_REPLY,
		ChairId:  p.ChairId,
		Multiple: data.Multiple,
	})
	//标志
	flag := true

	for _, v := range this.Players {
		if this.TableConfig.GameType == 1 {
			if v.BetMultiple == 0 && v.ChairId != this.Banker {
				flag = false
			}
		} else {
			if v.BetMultiple == 0 {
				flag = false
			}
		}
	}

	if flag {
		//进入发牌
		this.nextStage(STAGE_DEAL)
	}

	return
}

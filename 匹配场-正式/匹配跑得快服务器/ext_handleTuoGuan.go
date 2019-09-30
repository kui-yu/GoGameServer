package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) HandleTuoGuan(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("接收到玩家托管请求")
	re := TuoGuan{}
	err := json.Unmarshal([]byte(d.Data), &re)
	if err != nil {
		logs.Error("处理玩家托管请求-----数据转换格式失败!!", err)
		return
	}
	if re.Ctl == 1 && !p.TuoGuan || re.Ctl == 2 && p.TuoGuan {
		p.TuoGuan = !p.TuoGuan
		this.BroadcastAll(MSG_GAME_INFO_TUOGUAN_BRO, TuoGuanReply{
			Id:  MSG_GAME_INFO_TUOGUAN_BRO,
			Cid: p.ChairId,
			Ctl: re.Ctl,
		})
		//如果现在为出牌状态，并且玩家托管
		if this.GameState == GAME_STATUS_OUTCARD && p.TuoGuan {
			logs.Debug("玩家在出牌阶段进行托管:", p.TuoGuan)
			//如果轮到托管玩家出牌
			if this.CurCid == p.ChairId {
				//删除之前的定时器，设置托管定时器
				this.DelTimer(GAME_STATUS_OUTCARD)
				this.AddTimer(GAME_STATUS_OUTCARD, 1, this.TuoGuanOut, nil)
			}
		}
		return
	}
}

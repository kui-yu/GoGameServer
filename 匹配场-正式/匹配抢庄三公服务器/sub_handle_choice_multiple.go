package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) ChoiceMultiple(p *ExtPlayer, d *DkInMsg) {
	req := new(GChoiceMultipleREQ)
	err := json.Unmarshal([]byte(d.Data), req)
	if err != nil {
		logs.Debug("下注json解析失败:%v", err)
	}
	res := GChoiceMultipleRES{
		Id:       MSG_GAME_INFO_IDLE_MULTIPLE_REPLY,
		Uid:      p.Uid,
		Multiple: req.Multiple,
	}
	if req.Multiple < 0 || req.Multiple > int64(len(gameConfig.Idle_Multiple)-1) {
		return
	}
	if this.GameState != STAGE_IDLE_MULTIPLE {
		res.Err = "不是下注阶段"
		res.Result = 3
		p.SendNativeMsg(MSG_GAME_INFO_IDLE_MULTIPLE_REPLY, res)
		return
	}
	if p.BankerInfos.Multiple != 0 {
		res.Err = "已经下注过了"
		res.Result = 4
		p.SendNativeMsg(MSG_GAME_INFO_IDLE_MULTIPLE_REPLY, res)
		return
	}
	//成功下注
	p.BankerInfos.Multiple = gameConfig.Idle_Multiple[req.Multiple] //赋值倍数
	this.BroadcastAll(MSG_GAME_INFO_IDLE_MULTIPLE_REPLY, res)
	count := 0
	for _, v := range this.Players {
		if v.BankerInfos.Multiple > 0 {
			count++
		}
	}
	//所有人都下注了进入开牌阶段
	if count == len(this.Players)-1 {
		this.DelTimer(2) //删除计时器
		this.StageOpenCards("")
	}
}

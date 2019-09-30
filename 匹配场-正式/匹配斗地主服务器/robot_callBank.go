package main

import (
	"encoding/json"
	// "logs"
)

//机器人叫分
func (this *ExtDesk) robotCallBank(d interface{}) {

	p := d.(*ExtPlayer)
	//判断机器人 是否要叫分
	pCallFen := R_CallBank(p.HandCard)
	// logs.Debug("机器人叫的分数", pCallFen)

	//当前分数
	if this.CallFen >= pCallFen {
		pCallFen = -1
	}

	//叫分结束
	dv, _ := json.Marshal(GCallMsg{
		Coins: pCallFen,
	})
	this.HandleGameCall(p, &DkInMsg{
		Uid:  p.Uid,
		Data: string(dv),
	})
}

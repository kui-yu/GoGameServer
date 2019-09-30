package main

import (
	"logs"
)

func (this *ExtDesk) Leave(p *ExtPlayer) {
	logs.Debug("接收到玩家请求返回", p.Nick)
	//判断玩家该局是否有下注
	var allbet int64
	for _, v := range p.DownBet {
		allbet += v
	}
	if allbet > 0 {
		logs.Error("玩家正在游戏中，无法退出")
		p.SendNativeMsg(MSG_GAME_INFO_WARNING, MsgWarning{
			Id:     MSG_GAME_INFO_WARNING,
			Result: 3,
		})
		return
	}
	if p.Uid == this.Zhuang.Uid {
		logs.Error("玩家是庄无法退出")
		p.SendNativeMsg(MSG_GAME_INFO_WARNING, MsgWarning{
			Id:     MSG_GAME_INFO_WARNING,
			Result: 7,
		})
		return
	}

	logs.Debug("玩家成功退出")
	this.UpChair(p)       //玩家离座
	this.BroChairChange() //座位变更通知
	this.LeaveByForce(p)  //玩家离开
}

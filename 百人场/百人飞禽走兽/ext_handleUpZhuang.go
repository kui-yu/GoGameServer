package main

import (
	"logs"
)

func (this *ExtDesk) HandleUpZhuang(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("接收到玩家请求上庄请求")
	//判断玩家金币是否足够
	if p.Coins < gameConfig.LimitInfo.UpZhuangNeed {
		p.SendNativeMsg(MSG_GAME_INFO_UPZHUANG_REPLY, ChangZhuangReply{
			Id:     MSG_GAME_INFO_UPZHUANG_REPLY,
			Result: 3,
			Err:    "对不起，您的金币不足，无法上庄",
		})
		return
	}
	this.UpZhuang(p)
}

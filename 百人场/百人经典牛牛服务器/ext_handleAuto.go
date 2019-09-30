package main

import (
	"logs"
)

//处理匹配请求
func (this *ExtDesk) HandleAuto(p *ExtPlayer, d *DkInMsg) {
	p.LiXian = false
	//初始化玩家
	p.InitExtData()
	logs.Debug("接收到玩家", p.Nick, "的匹配请求")
	logs.Debug("该玩家的头像是：：：：：：：", p.Head)
	logs.Debug("下面将展示所有玩家!!!!!!")

	if p.Nick == "app304" {
		p.Robot = true
	}

	for _, v := range this.Players {
		logs.Debug("玩家", v.Coins)
	}
	//发送匹配成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
		Result:   0,
		CostType: GetCostType(),
	})
	//判断该玩家进入房间时，房间椅子是否有空位，如果有空位，就入座
	this.OnChair(p)
	this.BroChairChangeNoto(p)
}

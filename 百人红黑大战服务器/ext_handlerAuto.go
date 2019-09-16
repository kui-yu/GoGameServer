package main

import (
	"logs"
)

//玩家自由匹配
func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("匹配中")
	//	初始化
	p.Init()
	//	发送匹配成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})
	// logs.Debug("下注列表：", G_DbGetGameServerData.GameConfig.TenChips)
}

//匹配完成后续处理
func (this *ExtDesk) HandleGameAutoFinal(p *ExtPlayer, d *DkInMsg) {
	//	将玩家添加到座位管理器中
	this.SeatMgr.AddPlayer(p)

	//	群发用户房间人数和座位信息
	gameReply := GInfoAutoGameReply{
		Id:        MSG_GAME_INFO_AUTO_REPLY,
		PlayerNum: int32(len(this.Players)),
	}

	allUser := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range allUser {
		gameReply.SeatList = this.GetSeatInfo(v.(*ExtPlayer))
		v.(*ExtPlayer).SendNativeMsg(MSG_GAME_INFO_AUTO_REPLY, gameReply)
	}
	this.HandleRoomInfo(p)
}

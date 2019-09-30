package main

import (
	"logs"
)

//用户离线处理
func (this *ExtDesk) HandleDisconnect(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("玩家", p.Nick, "掉线")
	//判断玩家是否下注，如果没有下注，则踢掉
	var allbet int64
	for _, v := range p.DownBet {
		allbet += v
	}
	if allbet > 0 {
		logs.Debug("设置未离线状态")
		p.LiXian = true
	} else {
		this.UpChair(p)
		this.BroChairChange()
		this.LeaveByForce(p)
	}
}

//用户断线重连
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("玩家", p.Nick, "断线重连")
	p.LiXian = false
	//检测是否存在空位置，如果存在空位置，则自动将该玩家存入座位中
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:       MSG_GAME_RECONNECT_REPLY,
		Result:   0,
		CostType: GetCostType(),
	})
	//发送桌子信息
	deskInfo := &GClientDeskInfo{
		Id:        MSG_GAME_INFO_RECONNECT_REPLY,
		Result:    0,
		FangHao:   this.GetRoomName(),
		JuHao:     this.JuHao,
		BetLevels: G_DbGetGameServerData.GameConfig.TenChips,
		PlayerMassage: PlayerMsg{
			Uid:          p.Uid,
			MyUserAvatar: p.Head,
			MyUserName:   p.Nick,
			MyUserCoin:   p.Coins,
		},
		AreaCoin:           this.DownBet,
		AreaMaxCoin:        G_DbGetGameServerData.GameConfig.LimitRedMax,
		GameStatus:         this.GameState,
		GameStatusDuration: int64(this.GetTimerNum(this.GameState)),
		CardGroupArray:     this.CardGroupArray,
		ChairList:          p.getChairList(),
		Zoushi:             this.GameZs.Zoushi,
		BetAbleIndex:       this.CanUseChip(p),
		MyDownBets:         p.DownBet,
	}
	logs.Debug("用户断线重连给其发的东西", this.CardGroupArray)
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT_REPLY, deskInfo)
	//判断该玩家进入房间时，房间椅子是否有空位，如果有空位，就入座
	this.OnChair(p)
	this.BroChairChangeNoto(p)
}

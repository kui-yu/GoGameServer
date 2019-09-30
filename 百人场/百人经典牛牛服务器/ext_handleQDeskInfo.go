package main

import (
	"logs"
)

func (this *ExtDesk) HandleQDeskInfo(p *ExtPlayer, d *DkInMsg) {
	//处理玩家桌子信息请求
	logs.Debug("接收到玩家请求桌子信息===========", p.Nick)
	//检测是否存在空位置，如果存在空位置，则自动将该玩家存入座位中

	//发送桌子信息
	deskInfo := &GClientDeskInfo{
		Id:        MSG_GAME_INFO_QDESKINFO_REPLY,
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
	//初始化筹码列表
	p.SendNativeMsg(MSG_GAME_INFO_QDESKINFO_REPLY, deskInfo)
}

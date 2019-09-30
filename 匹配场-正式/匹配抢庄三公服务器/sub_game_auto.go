package main

import (
	"logs"
)

func (this *ExtDesk) GameAuto(p *ExtPlayer, d *DkInMsg) {
	//匹配成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
		CostType: GetCostType(),
	})
	logs.Debug("匹配玩家成功信息:%v,%v", p.Nick, MSG_GAME_AUTO_REPLY)
	this.SendRoomInfo(p) //发送房间消息
	//时间到开始游戏，每有玩家进入刷新计时器
	var isHave bool
	for _, v := range this.TList {
		if v.Id == 7 {
			isHave = true
			v.T = gameConfig.Game_Auto_Timer
		}
	}
	if !isHave && len(this.Players) == 2 {
		this.AddTimer(7, gameConfig.Game_Auto_Timer, this.GameStart, "")
	}
	//如果人满就开始游戏
	if len(this.Players) == GCONFIG.PlayerNum {
		this.GameStart("")
		this.DelTimer(7)
	}

}

//发送桌子消息
func (this *ExtDesk) SendRoomInfo(p *ExtPlayer) {
	//发送所有玩家信息
	for _, v := range this.Players {
		seatInfos := GPlayerSeatInfos{
			Id: MSG_GAME_INFO_SEAT,
		}
		for _, l := range this.Players {
			seatInfo := SeatInfo{
				Head:    l.Head,
				Coins:   l.Coins,
				Name:    l.Nick,
				Uid:     l.Uid,
				ChairId: l.ChairId,
			}
			if v.Uid != l.Uid {
				seatInfo.Name = "***" + l.Nick[len(p.Nick)-4:]
			}
			seatInfos.Data = append(seatInfos.Data, seatInfo)
		}
		v.SendNativeMsg(MSG_GAME_INFO_SEAT, seatInfos)
	}
	this.BScore = int64(G_DbGetGameServerData.Bscore) //设置底注
	//发送房间信息
	if this.JuHao == "" { //设置局号
		this.JuHao = GetJuHao()
	}
	p.SendNativeMsg(MSG_GAME_INFO_ROOM, GRoomInfo{
		Id:             MSG_GAME_INFO_ROOM,
		RoomNumber:     this.JuHao,
		MaxMultiple:    gameConfig.Max_Multiple,
		BScore:         this.BScore,
		PlayerMultiple: gameConfig.Idle_Multiple,
	})
}

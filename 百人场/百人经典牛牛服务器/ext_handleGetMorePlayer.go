package main

import (
	"logs"
)

func (this *ExtDesk) HandleGetMorePlayer(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("接收到请求跟多玩家信息", p.Nick)
	p.SendNativeMsgForce(MSG_GAME_INFO_GETMOREPLAYER_REPLAY, MorePlayer{
		Id:        MSG_GAME_INFO_GETMOREPLAYER_REPLAY,
		PlayerMsg: p.getMorePlayer(),
	})
	logs.Debug("发送过去的更多玩家信息:", p.getMorePlayer())
}

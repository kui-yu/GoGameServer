package main

// "github.com/astaxie/beego/logs"
// . "MaJiangTool"
// "encoding/json"
// "logs"

func (this *ExtDesk) HandleGiveUp(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != GAME_STATE_PLAY {
		return
	}
	//
	if p.GiveUp {
		return
	}
	//
	p.GiveUp = true
	//广播认输
	this.BroadcastAll(MSG_GAME_INFO_GIVEUP, &GGiveUpNotify{
		Id:  MSG_GAME_INFO_GIVEUP,
		Cid: int(p.ChairId),
	})
	//判断游戏是否继续
	cnt := 0
	for _, v := range this.Players {
		if v.GiveUp {
			cnt++
		}
	}
	if cnt >= len(this.Players)-1 {
		this.ClearTimer()
		//游戏结束
	}
}

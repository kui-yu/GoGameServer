package main

//开牌阶段
func (this *ExtDesk) StageOpenCards(d interface{}) {
	this.GameState = STAGE_OPEN_CARDS
	this.BroadStageTime(gameConfig.Stage_Open_Cards_Timer)
	this.AddTimer(3, gameConfig.Stage_Open_Cards_Timer, this.StageOpenCardsEnd, "")
}
func (this *ExtDesk) StageOpenCardsEnd(d interface{}) {
	req := GShowCardsREQ{Id: MSG_GAME_INFO_SHOW_CARDS}
	for _, v := range this.Players {
		if !v.IsOpenCards {
			req.Uid = v.Uid
			GDeskMgr.AddNativeMsg(MSG_GAME_INFO_SHOW_CARDS, v.Uid, req)
		}
	}
}

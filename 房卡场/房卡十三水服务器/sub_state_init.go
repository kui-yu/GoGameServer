package main

func (this *ExtDesk) GameStateInit() {
	//准备阶段
	this.BroadStageTime(0)
	//重置玩家信息
	for _, v := range this.Players {
		this.ResetPlayer(v)
	}
}

package main

// //初始化展示：
// type ExtDesk struct {
// 	Desk
// }

type ExtDesk struct {
	Desk

	Bscore      int32
	CardMgr     MgrCard   // 扑克牌牌管理
	CurCid      int32     // 当前用户的椅子id
	DiPai       []byte    // 底牌
	CallFen     int32     //叫分
	Banker      int32     //庄家
	DiPaiDoulbe int32     //底牌倍数
	MaxDouble   int32     //最大倍率
	MaxChuPai   *GOutCard //当前最大出牌
	RdChuPai    []*GOutCard
	CallTimes   int32 //已经都不叫分几次，超过3次，游戏结束
}

func (this *ExtDesk) InitExtData() {
	//牌内容初始化
	this.CardMgr.InitCards()
	this.CardMgr.InitNormalCards()
	//
	this.Handle[MSG_GAME_AUTO] = this.HandleGameAuto
	this.Handle[MSG_GAME_INFO_CALL] = this.HandleGameCall
	this.Handle[MSG_GAME_INFO_OUTCARD] = this.HandleGameOutCard
	this.Handle[MSG_GAME_INFO_PASS] = this.HandlePass
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
	this.Handle[MSG_GAME_INFO_TUOGUAN] = this.HandleTuoGuan
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
}

//广播阶段
func (this *ExtDesk) BroadStageTime(time int32) {
	stage := GStageInfo{
		Id:        MSG_GAME_INFO_STAGE,
		Stage:     int32(this.GameState),
		StageTime: time,
	}
	this.BroadcastAll(MSG_GAME_INFO_STAGE, &stage)
}

package main

type ExtDesk struct {
	Desk
	BScore          int64            //底注
	DeskCards       []int            //卡组
	ChoiceBankerArr []PlayerMultiple //抢庄信息
	DeskBankerInfos DeskBankerInfo
	HandlePlayer    int         //抢庄的玩家数
	SettleInfo      GSettleInfo //结算信息
}
type PlayerMultiple struct {
	Player *ExtPlayer
	IsJoin int
}
type DeskBankerInfo struct {
	BankerId       int32 //庄家的位置
	BankerMultiple int64 //庄家的位置
	//CardMultiple   int64 //庄家卡倍数
}

func (this *ExtDesk) InitExtData() {
	this.Handle[MSG_GAME_AUTO] = this.GameAuto                     //匹配
	this.Handle[MSG_GAME_INFO_BANKER_MULTIPLE] = this.ChoiceBanker //抢庄
	this.Handle[MSG_GAME_INFO_IDLE_MULTIPLE] = this.ChoiceMultiple //下注倍数
	this.Handle[MSG_GAME_INFO_SHOW_CARDS] = this.ShowCards         //亮牌
	this.Handle[MSG_GAME_INFO_GET_RECORD] = this.HandleGetRecord   //获取游戏记录
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
}

//初始化桌子
func (this *ExtDesk) initDesk() {
	this.GameState = GAME_STATUS_FREE
	this.JuHao = ""
	this.ChoiceBankerArr = make([]PlayerMultiple, 0)
	this.DeskBankerInfos = DeskBankerInfo{}
	this.SettleInfo = GSettleInfo{}
	this.HandlePlayer = 0
	// this.BroadStageTime(0)
}

//群发当前阶段
func (this *ExtDesk) BroadStageTime(time int) {
	for _, v := range this.Players {
		v.SendNativeMsg(MSG_GAME_INFO_STAGE_INFO, &GSGameStageInfo{
			Id:        MSG_GAME_INFO_STAGE_INFO,
			Stage:     this.GameState,
			StageTime: time,
		})
	}
}

//自封装定时器
func (this *ExtDesk) runTimer(t int, h func(interface{})) {
	//定时器ID，定时器时间，可执行函数，可执行参数
	this.AddTimer(10, t, h, nil)
}

package main

type ExtDesk struct {
	Desk
	Stage       int          //阶段
	CardsRound  int          //控制每3局洗一次牌
	PlaceBet    []int64      //区域下注，0到3分别是黑红梅方
	DeskCards   []int        //桌面牌
	HandCards   []Card       //庄家闲家手牌，索引0为庄家手牌
	MaxBet      int64        //限红
	BetArr      []int64      //下注筹码
	GameTrend   [][]Trend    //游戏走势，根据索引分别为庄、黑、红、梅、方, Trend.Player字段 =>  0为闲家，1为庄家
	RoundResult []AreaRes    //赢WinCoins为正数，Multiple为闲家牌倍数，输WinCoins为负数，Multiple为庄家牌倍数
	ManyPlayer  []ManyPlayer //桌面围观玩家
	AreaRes     []int        //黑红梅方4区域的输赢，0为输，1为赢
}

//游戏入口
func (this *ExtDesk) InitExtData() {
	//初始化桌子
	this.initDesk()
	//操作
	this.Handle[MSG_GAME_AUTO] = this.PlayerIn                   //玩家进场
	this.Handle[MSG_GAME_INFO_CHANGE_DOUBLE] = this.ChangeDouble //改变倍数
	this.Handle[MSG_GAME_INFO_PLAYER_BET] = this.PlayerBet       //下注
	//this.Handle[MSG_GAME_INFO_ROUND_SETTLE] = this.GetRoundSettle //获取玩家历史输赢
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect       //掉线重连
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect     //掉线信息
	this.Handle[MSG_GAME_LEAVE] = this.HandleGameBack            //返回大厅
	this.Handle[MSG_GAME_INFO_OTHER_PLAYER] = this.OtherPlayer   //请求玩家列表
	this.Handle[MSG_GAME_INFO_GET_RECORD] = this.HandleGetRecord //瀑布流获取玩家历史输赢
	//游戏开始
	this.ShuffleStage("")
}

//初始化桌子
func (this *ExtDesk) initDesk() {
	this.MaxBet = maxBet
	this.PlaceBet = make([]int64, 4)
	this.HandCards = make([]Card, 5)
	this.GameTrend = make([][]Trend, 5)
	this.RoundResult = make([]AreaRes, 4)
	this.AreaRes = make([]int, 4)
}

//群发当前阶段
func (this *ExtDesk) BroadStageTime(time int) {
	for _, v := range this.Players {
		v.SendNativeMsg(MSG_GAME_INFO_STAGE_INFO, &GSGameStageInfo{
			Id:    MSG_GAME_INFO_STAGE_INFO,
			Stage: this.Stage,
			Time:  time,
		})
	}
}

//自封装定时器
func (this *ExtDesk) runTimer(t int, h func(interface{})) {
	//定时器ID，定时器时间，可执行函数，可执行参数
	this.AddTimer(10, t, h, nil)
}

//消息协议函数
func (this *ExtDesk) Leave(p *ExtPlayer) bool { //重写的方法
	this.HandleGameBack(p, nil)
	return true
}

package main

//import (
//	"sort"
//)

type ExtPlayer struct {
	Player
	//以下是用户自己定义的变量
	HandCards    []int   //初始手牌
	HandColor    []int   //卡牌花色
	OldHandCard  []int   //发给前端的牌
	PayCoin      []int64 //下注数量
	CardType     int     //牌面类型 0未看,1已看,2弃牌(或同等状态的已无效牌)
	IsGU         bool    //是否是弃牌操作
	CardLv       int     //牌等级
	WinCoins     int64   //赢得金币数
	AutoFollowUp int     //是否自动跟注 0否，1是
	ProtectGU    int     //防超时弃牌 0否，1是
	RateCoins    float64 //抽水
	IsLeave      int     //0没有离开 1已离开
	IsReady      int
	TotalCoins   int64          //总结算金币
	RecordInfos  []GSRecordInfo //游戏记录
}

//重置玩家信息
func (this *ExtDesk) ResetPlayer(p *ExtPlayer) {
	p.HandCards = []int{}
	p.HandColor = []int{}
	p.OldHandCard = []int{}
	p.CardType = 0
	p.CardLv = 0
	p.PayCoin = []int64{}
	p.WinCoins = 0
	p.AutoFollowUp = 0
	p.IsGU = false
	p.RateCoins = 0
	p.ProtectGU = 0
	this.SettleContest = []PlayerContest{}
}

//个人战绩回复
func (this *ExtDesk) HandleRecord(p *ExtPlayer, d *DkInMsg) {

	info := GSRecordInfos{
		Id:    MSG_GAME_INFO_RECORD_INFO_REPLY,
		Infos: p.RecordInfos,
	}
	p.SendNativeMsg(MSG_GAME_INFO_RECORD_INFO_REPLY, &info)
}

//注册玩家消息
func (this *ExtDesk) InitExtData() {
	//初始化桌子信息
	this.InitAttribute()
	//牌内容初始化
	this.CardMgr.InitCards()
	//玩家匹配 400001
	this.Handle[MSG_GAME_FK_JOIN] = this.HandleGameAuto
	//玩家看牌 410005
	this.Handle[MSG_GAME_INFO_LOOK_CARD] = this.HandleLookCard
	//玩家弃牌 410006
	this.Handle[MSG_GAME_INFO_GIVE_UP] = this.HandleGiveUp
	// //玩家金币不足比牌
	// this.Handle[MSG_GAME_CONTEST] = this.HandleConinEnd
	//玩家下注 410008
	this.Handle[MSG_GAME_INFO_PLAY_INFO] = this.HandleGamePlay
	//玩家属性操作 410010
	this.Handle[MSG_GAME_INFO_PLAY_WITH_SYS] = this.HandlePlayWithSys
	//断线重连 400010
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
	//断线消息 400013断线
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
	//离开消息
	this.Handle[MSG_GAME_INFO_LEAVE] = this.HandleIsLeave
	//机器人换牌
	this.Handle[MSG_GAME_INFO_CHANGE_CARD] = this.HandleChangeRobotCard
	//玩家准备
	this.Handle[MSG_GAME_INFO_READY] = this.HandleReady
	//玩家解散
	this.Handle[MSG_GAME_INFO_DISMISS] = this.HandleDisMiss
	//个人战绩查询
	this.Handle[MSG_GAME_INFO_RECORD_INFO] = this.HandleRecord
}

package main

// //初始化展示：
// type ExtDesk struct {
// 	Desk
// }
import (
	"logs"
)

type ExtDesk struct {
	Desk

	//房间配置信息
	TableConfig            GATableConfig
	Bscore                 int32      //底分
	Round                  int        //回合
	CardMgr                MgrCard    // 扑克牌牌管理
	Lz_CardMgr             Lz_MgrCard //赖子模式扑克牌管理
	CurCid                 int32      // 当前用户的椅子id
	DiPai                  []byte     // 底牌
	CallFen                int32      //叫分
	Double                 int32      //当前桌子倍数
	Banker                 int32      //庄家
	DiPaiDoulbe            int32      //底牌倍数
	MaxDouble              int32      //最大倍率
	MaxChuPai              *GOutCard  //当前最大出牌
	RdChuPai               []*GOutCard
	CallTimes              int32 //已经都不叫分几次，超过3次，游戏结束
	AgreetBreakRoomCount   int   //同意房间解散人数
	NoAgreetBreakRoomCount int   //不同意人数
	CallOrGet              int   //如果为抢地主 ，代表下次是 叫地主 还是抢地主
}

//重置桌子
func (this *ExtDesk) ResetTable() {
	this.JuHao = ""
	this.TableConfig = GATableConfig{}
}

func (this *ExtDesk) ResetExtDest() {
	this.CardMgr.InitCards()
	this.CardMgr.InitNormalCards()
	this.Lz_CardMgr.InitCard()
	this.Lz_CardMgr.InitNormalCard()
	this.DiPai = []byte{}
	this.CallFen = 0
	this.DiPaiDoulbe = 0
	this.CallTimes = 0
	this.Banker = -1
	this.CurCid = -1
	this.AgreetBreakRoomCount = 0
}
func (this *ExtDesk) InitExtData() {
	this.Round = 0
	//牌内容初始化
	this.CardMgr.InitCards()
	this.CardMgr.InitNormalCards()
	this.Lz_CardMgr.InitCard()
	this.Lz_CardMgr.InitNormalCard()
	this.CallOrGet = 1
	//
	this.Handle[MSG_GAME_FK_JOIN] = this.HandleGameAuto
	this.Handle[MSG_GAME_INFO_CALL] = this.HandleGameCall
	this.Handle[MSG_GAME_INFO_OUTCARD] = this.HandleGameOutCard
	this.Handle[MSG_GAME_INFO_PASS] = this.HandlePass
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
	this.Handle[MSG_GAME_INFO_TUOGUAN] = this.HandleTuoGuan
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
	this.Handle[MSG_GAME_INFO_READY] = this.HandleReady
	this.Handle[MSG_GAME_INFO_GETMSG] = this.HandleGameGetMsg
	this.Handle[MSG_GAME_INFO_BREAKROOM] = this.handleBreakRomm
	this.Handle[MSG_GAME_INFO_OUTCARD_Lz] = this.handleGameOutCard_lz
	this.Handle[MSG_GAME_INFO_SELECT] = this.handleOutCardSelect_lz
	this.Handle[MSG_GAME_INFO_QPLAYERLOGS] = this.handleQPalyLogs
}

//广播阶段
func (this *ExtDesk) BroadStageTime(time int32) {
	stage := GStageInfo{
		Id:        MSG_GAME_INFO_STAGE,   //状态Id
		Stage:     int32(this.GameState), //状态
		StageTime: time,                  //状态时间
	}
	this.BroadcastAll(MSG_GAME_INFO_STAGE, &stage)
}

//计算房费
func (this *ExtDesk) getPayMoney() int64 {
	if this.TableConfig.MatchCount <= 0 {
		logs.Error("配置局数出错:", this.TableConfig.MatchCount)
	}
	money := int64(this.TableConfig.MatchCount) / 5
	if money == 0 {
		money = 1
	}
	return money
}

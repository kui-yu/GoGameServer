package main

import (
	. "MaJiangTool"
)

// //初始化展示：
// type ExtDesk struct {
// 	Desk
// }

type ExtDesk struct {
	Desk

	CardMgr CardManager108 // 牌管理
	Banker  int            // 庄家
	Dice1   int
	Dice2   int
	//动作
	Peng     ActionPeng
	MingGang ActionMingGang
	AnGang   ActionAnGang
	PuBuGang ActionPuBuGang
	HuPai    ExtActionHuPai
	//出牌，动作之后触发的动作
	ActionsAfterSendCard    []ActionIe
	ActionsAfterOutCard     []ActionIe
	ActionsAfterActionOther []ActionIe
	ActionsAfterActionSelf  []ActionIe
	//事件管理器
	EventManager EventMgr //事件管理器
	//
	CurCid int // 当前用户的椅子id
	//
	MaxDouble int
	Bscore    int
}

func (this *ExtDesk) InitExtData() {
	//牌内容初始化
	this.CardMgr.Initialize()
	//
	this.PuBuGang.InitData(nil, ActionType_Gang_PuBuGang)
	this.Peng.InitData(nil, ActionType_Peng)
	this.AnGang.InitData(nil, ActionType_Gang_An)
	this.MingGang.InitData(nil, ActionType_Gang_Ming)
	//
	this.ActionsAfterSendCard = []ActionIe{&this.AnGang, &this.PuBuGang}
	this.ActionsAfterOutCard = []ActionIe{&this.Peng, &this.MingGang}
	this.ActionsAfterActionSelf = this.ActionsAfterSendCard
	//
	this.Handle[MSG_GAME_AUTO] = this.HandleGameAuto
	this.Handle[MSG_GAME_INFO_OUTCARD] = this.HandleGameOutCard
	this.Handle[MSG_GAME_INFO_DINGQUE] = this.HandleDingQue
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect
	this.Handle[MSG_GAME_INFO_TUOGUAN] = this.HandleTuoGuan
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect
	this.Handle[MSG_GAME_INFO_GIVEUP] = this.HandleGiveUp
	this.Handle[MSG_GAME_INFO_HUANPAI] = this.HandleHuanPai
	this.Handle[MSG_GAME_INFO_ACTION] = this.HandleAction
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

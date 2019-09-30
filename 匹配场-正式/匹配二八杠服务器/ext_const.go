package main

//阶段
const (
	//匹配
	//开始
	STAGE_CALL    = 11 //抢庄阶段
	STAGE_BET     = 12 //下注阶段
	STAGE_DEAL    = 13 //发牌阶段
	STAGE_SETTLE  = 14 //结算阶段
	STAGE_RESTART = 15 //重新开始阶段
)

const (
	STAGE_START_TIME   = 3  //开始阶段时间
	STAGE_CALL_TIME    = 10 //抢庄阶段时间
	STAGE_PLAY_TIME    = 12 //下注阶段时间
	STAGE_DEAL_TIME    = 5  //发牌阶段时间
	STAGE_SETTLE_TIME  = 10 //结算阶段时间
	STAGE_RESTART_TIME = 5  //重新开始阶段时间
)

//系统定义
const (
	TIMER_OVER     = 8
	TIMER_OVER_NUM = 1 //归还桌子时间
)

//消息定义
const (
	MSG_GAME_INFO_STAGE                 = 410001 //阶段消息
	MSG_GAME_INFO_AUTO_REPLY            = 410002 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_ROOM_NOTIFY           = 410003 //房间信息通知410003
	MSG_GAME_INFO_CALL_INFO             = 410004 //玩家叫庄
	MSG_GAME_INFO_CALL_INFO_REPLY       = 410005 //返回玩家叫庄
	MSG_GAME_INFO_PLAY_INFO             = 410006 //玩家下注
	MSG_GAME_INFO_PLAY_INFO_REPLY       = 410007 //返回玩家下注
	MSG_GAME_INFO_BANKER_REPLY          = 410008 //庄家通知
	MSG_GAME_INFO_CARD_INFO_REPLY       = 410009 //发牌
	MSG_GAME_INFO_SETTLE_INFO_REPLY     = 410010 //结算
	MSG_GAME_INFO_RECONNECT             = 410011 //重连
	MSG_GAME_INFO_START_INFO_REPLY      = 410012 //开始信息
	MSG_GAME_INFO_RECORD_INFO           = 410013 //个人游戏战绩
	MSG_GAME_INFO_RECORD_INFO_REPLY     = 410014 //返回个人游戏战绩
	MSG_GAME_INFO_SETTLE_INFO_END_REPLY = 410015 //总结算
	MSG_GAME_INFO_CALL_LIST             = 410016 //叫庄列表
)

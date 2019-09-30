package main

// 阶段
const (
	STAGE_INIT           = 10 //准备阶段
	STAGE_CONTEST        = 13 //比牌阶段
	STAGE_PLAY_OPERATION = 14 //玩家操作阶段
	STAGE_SETTLE         = 15 //结算阶段
	STAGE_DISMISS        = 16 //解散阶段

	STAGE_CONTEST_TIME        = 1  //比牌阶段时间
	STAGE_START_TIME          = 2  //开始阶段时间
	STAGE_PLAY_OPERATION_TIME = 15 //玩家操作阶段时间
	STAGE_SETTLE_TIME         = 0  //结算阶段时间
	STAGE_END_TIME            = 0  //结束时间

	STAGE_DISMISS_TIME = 60 //解散时间
)

//系统定义
const (
	TIMER_OVER     = 16
	TIMER_OVER_NUM = 1 //归还桌子时间
)

const (
	MSG_GAME_INFO_STAGE             = 410001 //阶段消息
	MSG_GAME_INFO_AUTO_REPLY        = 410002 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_ROOM_NOTIFY       = 410003 //房间信息通知410003
	MSG_GAME_INFO_CALLPLAYER_REPLY  = 410004 //叫牌玩家通知
	MSG_GAME_INFO_LOOK_CARD         = 410005 //玩家看牌
	MSG_GAME_INFO_GIVE_UP           = 410006 //玩家弃牌
	MSG_GAME_INFO_CONTEST           = 410007 //金币不足玩家比牌
	MSG_GAME_INFO_PLAY_INFO         = 410008 //玩家下注
	MSG_GAME_INFO_PLAY_INFO_REPLY   = 410009 //返回玩家下注
	MSG_GAME_INFO_PLAY_WITH_SYS     = 410010 //玩家属性判断操作
	MSG_GAME_INFO_RECONNECT         = 410011 //重连
	MSG_GAME_INFO_SETTLE            = 410012 //结算
	MSG_GAME_INFO_COIN              = 410025 //金币消息
	MSG_GAME_INFO_LEAVE             = 410018 //中途离开
	MSG_GAME_INFO_LEAVE_REPLY       = 410015 //离开应答
	MSG_GAME_INFO_MAX               = 410016 //控牌消息
	MSG_GAME_INFO_CHANGE_CARD       = 410017 //换牌
	MSG_GAME_INFO_ERR               = 410500 // 错误信息
	MSG_GAME_INFO_READY_REPLY       = 410014 //玩家准备消息返回
	MSG_GAME_INFO_READY             = 410013 //玩家准备
	MSG_GAME_INFO_LUMPSUM_READY     = 410020 //总结算返回
	MSG_GAME_INFO_DISMISS           = 410021 //玩家离开
	MSG_GAME_INFO_DISMISS_REPLY     = 410022 //玩家离开消息返回
	MSG_GAME_INFO_RECORD_INFO       = 410023 //个人游戏战绩
	MSG_GAME_INFO_RECORD_INFO_REPLY = 410024 //返回个人游戏战绩
	MSG_GAME_INFO_START_INFO_REPLY  = 410026 //游戏开始消息
)

var GameRound = 7 //最大轮数+1

package main

//阶段
const (
	//匹配
	STAGE_INIT    = 10 //准备阶段
	STAGE_PLAY    = 12 //玩牌阶段
	STAGE_SETTLE  = 13 //结算阶段
	STAGE_DISMISS = 16 //解散阶段
)

const (
	STAGE_START_TIME   = 5  //开始
	STAGE_SETTLE_TIME  = 0  //结算阶段时间
	STAGE_DISMISS_TIME = 60 //解散时间
)

//消息定义
const (
	MSG_GAME_INFO_STAGE                 = 410001 //阶段消息
	MSG_GAME_INFO_AUTO_REPLY            = 410002 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_ROOM_NOTIFY           = 410003 //房间信息通知410003
	MSG_GAME_INFO_HANDINFO_REPLY        = 410004 //发送手牌信息
	MSG_GAME_INFO_PLAY                  = 410005 //玩家摆牌
	MSG_GAME_INFO_PLAY_REPLY            = 410006 //玩家摆牌通知
	MSG_GAME_INFO_SETTLE_INFO_REPLY     = 410007 //结算
	MSG_GAME_INFO_RECONNECT             = 410010 //重连
	MSG_GAME_INFO_READY                 = 410013 //玩家准备
	MSG_GAME_INFO_READY_REPLY           = 410014 //玩家准备消息返回
	MSG_GAME_INFO_START_INFO            = 410015 //游戏开始消息
	MSG_GAME_INFO_DISMISS               = 410016 //玩家解散
	MSG_GAME_INFO_DISMISS_REPLY         = 410017 //玩家解散消息返回
	MSG_GAME_INFO_RECORD_INFO           = 410018 //个人游戏战绩
	MSG_GAME_INFO_RECORD_INFO_REPLY     = 410019 //返回个人游戏战绩
	MSG_GAME_INFO_LEAVE                 = 410020 //玩家离开
	MSG_GAME_INFO_SETTLE_INFO_END_REPLY = 410021 //总结算
	MSG_GAME_INFO_ERR                   = 410500 // 错误信息
)

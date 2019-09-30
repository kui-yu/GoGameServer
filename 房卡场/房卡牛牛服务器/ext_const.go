package main

//阶段
const (
	STAGE_INIT        = 10 //准备阶段
	STAGE_CALL_BANKER = 11 //抢庄阶段
	STAGE_CALL_SCORE  = 12 //叫分阶段
	STAGE_DEAL        = 13 //发牌阶段
	STAGE_PLAY        = 14 //玩牌阶段
	STAGE_SETTLE      = 15 //结算阶段
	STAGE_DISMISS     = 16 //解散阶段
)

//阶段时间
const (
	STAGE_START_TIME       = 2  //开始时间
	STAGE_CALL_BANKER_TIME = 15 //叫庄时间
	STAGE_CALL_SCORE_TIME  = 15 //叫分时间
	STAGE_DEAL_POKER_TIME  = 3  //发牌时间
	STAGE_PLAYCARD_TIME    = 10 //玩牌
	STAGE_SETTLE_TIME      = 0  //结算
	STAGE_DISMISS_TIME     = 60 //解散时间
)

//消息定义
const (
	MSG_GAME_INFO_STAGE                 = 410030 //阶段消息
	MSG_GAME_INFO_CALL                  = 410001 //叫分
	MSG_GAME_INFO_CALL_REPLY            = 410002 //叫分返回
	MSG_GAME_INFO_DEAL_REPLY            = 410003 //发牌
	MSG_GAME_INFO_PLAY_REPLY            = 410004 //摆牌结果
	MSG_GAME_INFO_PLAY                  = 410005 //开牌
	MSG_GAME_INFO_SETTLE                = 410006 //结算
	MSG_GAME_INFO_CALL_BANKER           = 410007 //抢庄
	MSG_GAME_INFO_CALL_BANKER_NOTIFY    = 410008 //抢庄通知
	MSG_GAME_INFO_ROOM_NOTIFY           = 410009 //房间信息通知410009
	MSG_GAME_INFO_RECONNECT             = 410010 //重连
	MSG_GAME_INFO_CHOOSE_BANKER         = 410011 //选庄通知
	MSG_GAME_INFO_AUTO_REPLY            = 410012 //410012,游戏随机匹配成功的数据
	MSG_GAME_INFO_READY                 = 410013 //玩家准备
	MSG_GAME_INFO_READY_REPLY           = 410014 //玩家准备消息返回
	MSG_GAME_INFO_START_INFO            = 410015 //游戏开始消息
	MSG_GAME_INFO_DISMISS               = 410016 //玩家离开
	MSG_GAME_INFO_DISMISS_REPLY         = 410017 //玩家离开消息返回
	MSG_GAME_INFO_RECORD_INFO           = 410018 //个人游戏战绩
	MSG_GAME_INFO_RECORD_INFO_REPLY     = 410019 //返回个人游戏战绩
	MSG_GAME_INFO_SETTLE_INFO_END_REPLY = 410020 //总结算
	MSG_GAME_INFO_LEAVE                 = 410021 //玩家离开
	MSG_GAME_INFO_CALL_LIST             = 410022 //庄家筹码列表
	MSG_GAME_INFO_BET_LIST              = 410023 //下注筹码列表
)

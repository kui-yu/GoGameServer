package main

//阶段
const (
	//匹配
	STAGE_PLAY   = 12 //玩牌阶段
	STAGE_SETTLE = 13 //结算阶段
)

const (
	STAGE_START_TIME  = 5  //开始
	STAGE_PLAY_TIME   = 60 //玩牌阶段时间
	STAGE_SETTLE_TIME = 0  //结算阶段时间
)

const (
	TIMER_ROBOT     = 17
	TIMER_ROBOT_NUM = 5
	TIMER_OVER      = 8
	TIMER_OVER_NUM  = 1 //归还桌子时间
)

//消息定义
const (
	MSG_GAME_INFO_STAGE             = 410001 //阶段消息
	MSG_GAME_INFO_AUTO_REPLY        = 410002 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_ROOM_NOTIFY       = 410003 //房间信息通知410003
	MSG_GAME_INFO_HANDINFO_REPLY    = 410004 //发送手牌信息
	MSG_GAME_INFO_PLAY              = 410005 //玩家摆牌
	MSG_GAME_INFO_PLAY_REPLY        = 410006 //玩家摆牌通知
	MSG_GAME_INFO_SETTLE_INFO_REPLY = 410007 //结算
	MSG_GAME_INFO_RECONNECT         = 410010 //重连
	MSG_GAME_INFO_ERR               = 410500 // 错误信息
)

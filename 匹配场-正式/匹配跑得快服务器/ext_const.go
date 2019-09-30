package main

const (
	MSG_GAME_INFO_AUTO_REPALY   = 410000 + iota //随机匹配成功游戏数据 410000
	MSG_GAME_INFO_ROOM_NOTIFY                   //房间信息 410001
	MSG_GAME_INFO_STAGE                         //状态通知 410002
	MSG_GAME_INFO_SENDCARD                      //发牌通知 410003
	MSG_GAME_INFO_OUTCARD                       //玩家出牌 410004
	MSG_GAME_INFO_OUTCARD_REPLY                 //玩家出牌应答 410005
	MSG_GAME_INFO_OUTCARD_BRO                   //玩家出牌广播 410006
	MSG_GAME_INFO_TUOGUAN                       //玩家托管请求 410007
	MSG_GAME_INFO_TUOGUAN_BRO                   //玩家托管广播 410008
	MSG_GAME_INFO_PASS                          //玩家pass请求 410009
	MSG_GAME_INFO_PASS_BRO                      //玩家pass广播 410010
	MSG_GAME_INFO_EXIT                          //玩家请求退出 410011
	MSG_GAME_INFO_EXIT_REPLY                    //玩家退出请求应答  410012
	MSG_GAME_INFO_BALANCE_BRO                   //结算广播 410013
	MSG_GAME_INFO_DISORREC_BRO                  //玩家离线/上线广播 410014
	MSG_GAME_INFO_RECONNECT                     //玩家重新连接 410015
)

const (
	GAME_STATUS_SENDCAR      = 10 //发牌状态
	GAME_STATUS_SENDCAR_TIME = 5  //发牌状态时间
	GAME_STATUS_OUTCARD      = 11 //出牌阶段
	GAME_STATUS_OUTCARD_TIME = 20 //出牌时间
	GAME_STATUS_BALANCE      = 12 //结算状态
)

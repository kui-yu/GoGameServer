package main

const (
	MSG_GAME_INFO_STAGE_INFO           = 410000 + iota //410000//阶段消息
	MSG_GAME_RECONNECT_TABLE_REPLY                     //410001//加入百人场
	MSG_GAME_INFO_DESKINFO_REPLAY                      //410002//回复游戏桌子信息
	MSG_GAME_INFO_PLAYER_BET                           //410003//玩家下注
	MSG_GAME_INFO_PLAYER_BET_REPLAY                    //410004//玩家下注返回
	MSG_GAME_INFO_PLAYER_BET_MASS                      //410005//玩家下注群发
	MSG_GAME_INFO_CHANGE_DOUBLE                        //410006//倍数改变
	MSG_GAME_INFO_CHANGE_DOUBLE_REPLAY                 //410007//返回倍数改变
	MSG_GAME_INFO_OPEN_CARD_REPLY                      //410008//开牌结果
	MSG_GAME_INFO_SETTLE_REPLY                         //410009//结算信息
	MSG_GAME_INFO_DESKPLAYER_REPLAY                    //410010//桌面玩家
	MSG_GAME_INFO_ROUND_SETTLE                         //410011//获取玩家历史输赢
	MSG_GAME_INFO_ROUND_SETTLE_REPLY                   //410012//返回玩家历史输赢
	MSG_GAME_INFO_BACK                                 //410013//请求返回大厅
	MSG_GAME_INFO_BACK_REPLAY                          //410014//返回大厅
	MSG_GAME_INFO_RECONNECT                            //410015//重连
	MSG_GAME_INFO_OTHER_PLAYER                         //410016//请求玩家列表
	MSG_GAME_INFO_OTHER_PLAYER_REPLY                   //410017//返回玩家列表
	MSG_GAME_INFO_GET_RECORD                           //410018//获取游戏记录
	MSG_GAME_INFO_GET_RECORD_REPLY                     //410019//获取游戏记录
	MSG_GAME_INFO_THREE_THMES                          //410020//提醒三次没下注或返回大厅错误(双用)
	MSG_GAME_INFO_FIVE_THMES_LEAVE                     //410021//五次到踢出
)

const (
	STAGE_GAME_DOUBLE    = 10 + iota //10选择倍数
	STAGE_GAME_START                 //11//开始
	STAGE_GAME_BET                   //12//下注
	STAGE_GAME_STOP_BET              //13停止下注
	STAGE_GAME_SEND_CARD             //14发牌
	STAGE_GAME_OPEN_CARD             //15//开牌
	STAGE_SETTLE                     //16//派奖
	STAGE_SHUFFLE                    //17//洗牌
)

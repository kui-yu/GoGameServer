package main

const (
	MSG_GAME_INFO_STAGE_BET         = 410000 + iota //阶段消息=>下注阶段
	MSG_GAME_INFO_STAGE_STOP_BET            //410001//阶段消息=>停止下注
	MSG_GAME_INFO_STAGE_GAME_RESULT         //410002//阶段消息=>开奖阶段
	MSG_GAME_INFO_STAGE_SETTLE              //410003//阶段消息=>结算阶段
	MSG_GAME_INFO_PLAYER_IN                 //410004//玩家进场
	MSG_GAME_INFO_PLAYER_IN_REPLY           //410005//玩家进场回复
	MSG_GAME_INFO_PLAYER_BET_MASS           //410006//玩家下注群发
	MSG_GAME_INFO_PLAYER_BET                //410007//玩家下注
	MSG_GAME_INFO_PLAYER_BET_REPLAY         //410008//玩家下注返回
	MSG_GAME_INFO_BACK                      //410009//返回大厅
	MSG_GAME_INFO_BACK_REPLAY              //4100010//请求返回大厅
	MSG_GAME_INFO_ONLINE_PLAYER            //4100011 //在线玩家
)

const (
	STAGE_GAME_START_BET = 10 + iota //游戏下注阶段
	STAGE_BET_STOP               //11//停止下注
	STAGE_GAME_RESULT            //12//开奖阶段
	STAGE_SETTLE                 //13//结算阶段
)

const (
	GAME_RATE     = 0.05 //费率
	BET_TIMER     = 15   //投注时间
	LOTTERY_TIMER = 15   //开奖结算时间
)

const (
	PLACE_BIG   = iota    //大
	PLACE_SMALL        //1//小
	PLACE_ODD          //2//单
	PLACE_EVEN         //3//双

)

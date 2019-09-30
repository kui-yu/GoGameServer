package main

//自己定义的游服Id从410000开始
const (
	MSG_GAME_INFO_STATUSCHANGE         = 410000 + iota //410000  状态改变通知。
	MSG_GAME_INFO_QDESKINFO                            //410001  请求桌子信息。
	MSG_GAME_INFO_QDESKINFO_REPLY                      //410002  请求桌子信息返回。
	MSG_GAME_INFO_JUHAOCHANGE                          //410003  局号变更通知。
	MSG_GAME_INFO_CHAIRCHANGE                          //410004  座位玩家改变通知。
	MSG_GAME_INFO_DOWNBET                              //410005  玩家下注请求。
	MSG_GAME_INFO_DOWNBET_REPLAY                       //410006  玩家下注请求应答。
	MSG_GAME_INFO_DOWNBET_BRO                          //410007  玩家下注通知。
	MSG_GAME_INFO_FACARD_BRO                           //410008  发牌通知。
	MSG_GAME_INFO_BALANCE                              //410009  结算。
	MSG_GAME_INFO_GETMOREPLAYER                        //410010  请求更多玩家。
	MSG_GAME_INFO_GETMOREPLAYER_REPLAY                 //410011  请求更多玩家应答。
	MSG_GAME_INFO_BACK                                 //410012  玩家请求返回。
	MSG_GAME_INFO_BACK_REPLAY                          //410013  玩家请求返回应答。
	MSG_GAME_INFO_WARNING                              //410014  警告。
	MSG_GAME_INFO_RECONNECT_REPLY                      //410015  玩家断线重连应答。
	MSG_GAME_INFO_GET_RECORD                           //410016  获取游戏记录
	MSG_GAME_INFO_GET_RECORD_REPLY                     //410017  获取游戏记录应答
)

//游戏状态 ： 游戏状态在yaml配置文件中修改,查看
// const (
// 	GAME_STATUS_READY       = 10 + iota //准备状态
// 	GAME_STATUS_SHUFFLECARD             //洗牌
// 	GAME_STATUS_DOWNBET                 //下注状态
// 	GAME_STATUS_SENDCARD                //发牌
// 	GAME_STATUS_OPENCARD                //开牌状态
// 	GAME_STATUS_BALANCE                 //结算
// )

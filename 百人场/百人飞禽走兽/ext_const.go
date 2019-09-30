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
	MSG_GAME_INFO_LOTTERY_BRO                          //410008  开奖通知。
	MSG_GAME_INFO_BALANCE                              //410009  结算。
	MSG_GAME_INFO_GETMOREPLAYER                        //410010  请求更多玩家。
	MSG_GAME_INFO_GETMOREPLAYER_REPLAY                 //410011  请求更多玩家应答。
	MSG_GAME_INFO_BACK                                 //410012  玩家请求返回。
	MSG_GAME_INFO_BACK_REPLAY                          //410013  玩家请求返回应答。
	MSG_GAME_INFO_WARNING                              //410014  警告。
	MSG_GAME_INFO_RECONNECT_REPLY                      //410015  玩家断线重连应答。
	MSG_GAME_INFO_GET_RECORD                           //410016  获取游戏记录
	MSG_GAME_INFO_GET_RECORD_REPLY                     //410017  获取游戏记录应答
	MSG_GAME_INFO_UPZHUANG                             //410018  玩家请求上庄
	MSG_GAME_INFO_UPZHUANG_REPLY                       //410019  玩家请求上庄应答
	MSG_GAME_INFO_DOWNZHUANG                           //410020  玩家请求下庄
	MSG_GAME_INFO_DOWNZHUANG_REPLY                     //410021  玩家请求下庄应答
	MSG_GAME_INFO_CHANGEZHUANG                         //410022  庄家改变通知
	MSG_GAME_INFO_TOROBOTRESULT                        //410023  通知机器人下注正确区域
)

// 游戏状态 ： 游戏状态在yaml配置文件中修改,查看
const (
	GAME_STATUS_READY    = 10 + iota //准备状态
	GAME_STATUS_STARTBET             //开始下注状态
	GAME_STATUS_DOWNBET              //下注状态
	GAME_STATUS_ENDBET               //结束下注状态
	GAME_STATUS_LOTTERY              //开奖状态
	GAME_STATUS_BALANCE              //结算
)
const (
	GAME_STATUS_READY_TIME    = 3000  //准备状态时间
	GAME_STATUS_STARTBET_TIME = 2000  //开始下注时间
	GAME_STATUS_DOWNBET_TIME  = 15000 //下注时间
	GAME_STATUS_ENDBET_TIME   = 2000  //结束下注时间
	GAME_STATUS_LOTTERY_TIME  = 10000 //开奖状态时间
	GAME_STATUS_BALANCE_TIME  = 5000  //结算状态时间
)

//开奖将项
const (
	Feiqing = iota
	Zoushou
	Tuzi
	Yanzi
	Gezi
	Houzi
	Xiongmao
	Kongque
	Shizi
	Laoying
	Shayu
	JinShayu
)

var LotteryDouble map[int]float32 //倍数集合
func InitMultiple() {
	LotteryDouble = make(map[int]float32)
	LotteryDouble[Feiqing] = gameConfig.AreaDouble.Feiqing
	LotteryDouble[Zoushou] = gameConfig.AreaDouble.Zoushou
	LotteryDouble[Tuzi] = gameConfig.AreaDouble.Tuzi
	LotteryDouble[Yanzi] = gameConfig.AreaDouble.Yanzi
	LotteryDouble[Gezi] = gameConfig.AreaDouble.Gezi
	LotteryDouble[Houzi] = gameConfig.AreaDouble.Houzi
	LotteryDouble[Xiongmao] = gameConfig.AreaDouble.Xiongmao
	LotteryDouble[Kongque] = gameConfig.AreaDouble.Kongque
	LotteryDouble[Shizi] = gameConfig.AreaDouble.Shizi
	LotteryDouble[Laoying] = gameConfig.AreaDouble.Laoying
	LotteryDouble[Shayu] = gameConfig.AreaDouble.Shayu
	LotteryDouble[JinShayu] = gameConfig.AreaDouble.Jinshayu
}

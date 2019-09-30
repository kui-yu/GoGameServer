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
	MSG_GAME_INFO_TOROBOTRESULT                        //410021  通知机器人下注正确区域
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

//阶段信息
type GameStatuInfo struct {
	Id         int
	Status     int //当前状态
	StatusTime int //状态持续时间
}

//下注请求附带结构体
type DownBet struct {
	Id        int
	ChipIndex int //筹码Id
	AreaIndex int //区域Id
}

// 发送给客户端的桌子信息
type GClientDeskInfo struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string

	JuHao     string  // 局号
	FangHao   string  // 房号
	BetLevels []int64 // 下注级别

	MyDownBets    []int64   // 自己下注的集合
	PlayerMassage PlayerMsg //用户信息

	AreaCoin           []int64     //区域金币
	GameStatus         int         // 游戏状态
	GameStatusDuration int64       // 当前状态持续时间毫秒
	CardGroupArray     map[int]int // 玩家和庄家的牌 庄家牌索引最后一个4

	AreaMaxCoin int64 // 限制区域最大下注

	ChairList    []PlayerInfoByChair //座位玩家信息
	Zoushi       [][]Pzshi
	BetAbleIndex int
}

//用户信息
type PlayerMsg struct {
	Uid          int64  //玩家uid
	MyUserAvatar string // 用户头像
	MyUserName   string // 用户昵称
	MyUserCoin   int64  // 用户金币
}

//座位玩家信息
type PlayerInfoByChair struct {
	Uid    int64  //玩家Id    (uid 为 0时代表 该座位上没有人)
	Nick   string //玩家昵称
	Avatar string //头像
	Coins  int64  //玩家金币
}

//手牌信息
type CardGroupInfo struct {
	CardGroupType CardGroupType //牌型
	Cards         []int         //牌组
	MaxCard       uint8         //最大牌
}

//走势
type ZouShiToClient struct {
	Id     int32     //协议号
	Zoushi [][]Pzshi //区域走势
}

//走势详情
type Pzshi struct {
	ZouShi string //庄或区域开奖结果
	Mark   int    //0，代表庄家  1，代表闲家（用来变颜色） 2,-1代表庄家
}

// //手牌信息
// type CardGroupInfo struct {
// 	CardGroupType CardGroupType //牌型
// 	Cards         []int         //牌组
// 	MaxCard       uint8         //最大牌
// }
type CardGroupType int

const (
	_ CardGroupType = iota
	CardGroupType_Cattle_1
	CardGroupType_Cattle_2
	CardGroupType_Cattle_3
	CardGroupType_Cattle_4
	CardGroupType_Cattle_5
	CardGroupType_Cattle_6
	CardGroupType_Cattle_7     // 2倍
	CardGroupType_Cattle_8     // 3倍
	CardGroupType_Cattle_9     // 3倍
	CardGroupType_Cattle_C     // 4倍
	CardGroupType_Cattle_BOMB  // 炸弹 5倍
	CardGroupType_Cattle_WUHUA // 五花牛 不包括10 6倍
	CardGroupType_None
	CardGroupType_NotCattle
)

//下注应答
type DownBetReplay struct {
	Id           int    //协议号
	Result       int    //0代表成功，其他代表失败
	ErrStr       string //错误信息，成功时为空
	BetAbleIndex int    //可用筹码下标
	CoinsIndex   int    //金币下标
	AreaIndex    int    //区域下标
	SelfAllCoins int64  //自己区域总下注
	AllCoins     int64  //区域总下注
	Coins        int64  //身上金币
}

//自由匹配应答，此外还有一个匹配消息和游戏相关的（斗地主为GInfoAutoGameReply）
type GAutoGameReply struct {
	Id       int32
	Result   int32 //0成功，其他失败
	CostType int   //1金币，2代币
	Err      string
}

//发送给机器人开奖结果
type ToRobot struct {
	Id    int
	Index int //开奖区域
}

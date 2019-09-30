package main

// 定义游戏消息
const (
	MSG_GAME_INFO_STAGE            = 410001 //阶段消息
	MSG_GAME_INFO_AUTO_REPLY       = 410002 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_ROOM_NOTIFY      = 410003 //房间信息通知410003
	MSG_GAME_INFO_CALLPLAYER_REPLY = 410004 //叫牌玩家通知
	MSG_GAME_LOOK_CARD             = 410005 //玩家看牌
	MSG_GAME_GIVE_UP               = 410006 //玩家弃牌
	MSG_GAME_CONTEST               = 410007 //金币不足玩家比牌
	MSG_GAME_INFO_PLAY_INFO        = 410008 //玩家下注
	MSG_GAME_INFO_PLAY_INFO_REPLY  = 410009 //返回玩家下注
	MSG_GAME_PLAY_WITH_SYS         = 410010 //玩家属性判断操作
	MSG_GAME_INFO_RECONNECT        = 410011 //重连
	MSG_GAME_SETTLE                = 410012 //结算
	MSG_GAME_COIN                  = 410013 //金币消息
	MSG_GAME_INFO_LEAVE_REPLY      = 410015 //离开应答
	MSG_GAME_INFO_MAX              = 410016 //控牌消息
	MSG_GAME_INFO_CHANGE_CARD      = 410017 //换牌
)

// 阶段
const (
	STAGE_CONTEST        = 13 //比牌阶段
	STAGE_PLAY_OPERATION = 14 //玩家操作阶段
	STAGE_SETTLE         = 15 //结算阶段
	GAME_STATUS_START    = 1
	GAME_STATUS_END      = 2
)

//
const (
	LookCardCoin = 3
	CriticalCoin = 5
	MaxCoin      = 7

	LookCardRound = 3
	MaxRound      = 5 //
)

//阶段时间
type GSStageInfo struct {
	Id        int
	Stage     int
	StageTime int
}

//座位信息
type GSeatInfo struct {
	Uid   int64
	Nick  string
	Ready bool
	Cid   int32 //椅子号
	Sex   int32
	Head  string
	Lv    int32
	Coin  int64
}

//所有匹配玩家信息
type GInfoAutoGameReply struct {
	Id   int32
	Seat []GSeatInfo
}

//叫牌玩家
type GSPlayerCallPlayer struct {
	Id         int
	Player     int32 //叫牌玩家ChairId
	Round      int   //当前轮数
	CoinEnough int   //0不足 1充足
	MinCoin    int64 //当前最低注
}

//玩家操作
type GAPlayerOperation struct {
	Id        uint32
	ChairId   int32 //座位号
	PlayCoin  int64 //下注金币
	Operation int   //操作(0弃牌，1看牌，2比牌，3加注，4跟注)
}

//返回下注操作
type GSPlayerPayCoin struct {
	Id        int
	PChairId  int32   //当前操作玩家座位号
	PlayCoin  int64   //下注金币
	ChairId   []int32 //比牌玩家座位号
	Winner    int32   //获胜者chairid
	Operation int     //玩家操作返回(2比牌，3加注，4跟注)
}

//结算玩家信息
type GSSettlePlayInfo struct {
	Id       int
	Count    int          //比牌数
	PContest []Contest    //比牌玩家
	SCard    []SettleCard //结算后玩家看牌
	CList    []CoinList   //所有玩家结算金币
}

type Contest struct {
	Person_1 int32
	Person_2 int32
	Winner   int32
}

type SettleCard struct {
	ChairId  int32 //座位
	Identity int   //0 winner  1 loser
	HandCard []int //手牌
	Lv       int   //牌等级
}

type CoinList struct {
	ChairId  int32 //座位
	WinCoins int64 //输赢金币
	Coins    int64 //身上金币
}

//看牌
type GSCardInfo struct {
	Id         int
	HandCards  []int //手牌
	Lv         int   //牌等级
	ChairId    int32 //座位号
	Model      int   // 0 主动看牌，1 失败看牌
	CoinEnough int   //0不足 1充足
}

//玩家牌型(是否已看牌，弃牌)
type GSCardType struct {
	Id      int
	ChairId int32 //座位号
}

//通知牌
type GSMaxCard struct {
	Id             int
	IsRobot        int         //0玩家 1机器人
	CardLv         int         //最大牌型
	HandCard       []int       //手牌
	ChairId        int32       //座位号
	WinnerRole     int         //0玩家，1机器人
	PlayerHandCard []PHandCard //真实玩家手牌
}

type PHandCard struct {
	HandCards []int
	CardLv    int
	ChairId   int32 //座位号
}

//接收最大牌
type GAMaxCard struct {
	Id       int
	CardLv   int //
	HandCard []int
}

//换牌消息返回
type GSChangeCard struct {
	Id       int
	HandCard []int //手牌
	CardLv   int   //牌等级
	Result   int   //0成功，其他失败
}

//金币消息
type GSCoinMsg struct {
	Id      int
	AllCoin int64   //场内总金币
	PCoin   []int64 //玩家金币
	Round   int     //轮数
}

//离开应答
type GSPlayerLeave struct {
	Id        int32
	ChairId   int32 // 座位号
	LeaveType int   // 是否已从大厅离开 0否 1是 2继续游戏
}

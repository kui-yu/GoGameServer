package main

// 自己定义的游服id从410000开始
// 游戏状态使用id表示，方便客户端使用 前缀N通知 Q请求 R回复
const (
	MSG_GAME_INFO_START            = 410000 + iota // 游戏开始额外信息
	MSG_GAME_NSTATUS_CHANGE                        // 410001 游戏状态改变
	MSG_GAME_QSEATDOWN                             // 410002 玩家请求坐下
	MSG_GAME_RSEATDOWN                             // 410003 玩家请求坐下的回复
	MSG_GAME_NSEATDOWN                             // 410004 座位信息改变通知
	MSG_GAME_QDOWNBET                              // 410005 玩家请求下注
	MSG_GAME_RDOWNBET                              // 410006 玩家请求下注回复
	MSG_GAME_NDOWNBET                              // 410007 玩家下注通知
	MSG_GAME_FACARD                                // 410008 系统发牌
	MSG_GAME_OPENCARD                              // 410009 系统开牌
	MSG_GAME_BALANCE                               // 410010 结算
	MSG_GAME_NTIPS                                 // 410011 提示信息
	MSG_GAME_QHISTORY                              // 410012 请求走势
	MSG_GAME_RHISTORY                              // 410013 请求走势回复
	MSG_GAME_QMANYUSER                             // 410014 请求更多玩家信息
	MSG_GAME_RMANYUSER                             // 410015 请求更多玩家信息回复
	MSG_GAME_QBACK                                 // 410016 请求返回
	MSG_GAME_RBACK                                 // 410017 请求返回回复
	MSG_GAME_QDESKINFO                             // 410018 请求游戏桌子信息
	MSG_GAME_RDESKINFO                             // 410019 回复游戏桌子信息
	MSG_GAME_QSEATUP                               // 410020 玩家请求站立
	MSG_GAME_RSEATUP                               // 410021 玩家请求站立回复
	MSG_GAME_NDESKCHANGE                           // 410022 房间信息改变
	MSG_GAME_NRANKLIST                             // 410023 玩家排行列表
	MSG_GAME_QSEATINFO                             // 410024 请求座位信息
	MSG_GAME_RSEATINFO                             // 410025 请求座位信息回复
	MSG_GAME_INFO_RECONNECT                        // 410026 重连
	MSG_GAME_INFO_GET_RECORD                       // 410027 获取游戏记录
	MSG_GAME_INFO_GET_RECORD_REPLY                 // 410028 获取游戏记录应答
)

// 无效的牌
const Card_Invalid = 0xFF

// 牌的类型
const (
	Card_Fang = 0x10
	Card_Mei  = 0x20
	Card_Hong = 0x30
	Card_Hei  = 0x40
	Card_King = 0x50 // |14,15 小王，大王
)

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

// 牌类型倍数
var BetDoubleMap map[CardGroupType]int32

// 随机牌类型集合
var RandCardTypeMap []CardGroupType

func init() {
	BetDoubleMap = make(map[CardGroupType]int32)
	BetDoubleMap[CardGroupType_NotCattle] = 1
	BetDoubleMap[CardGroupType_Cattle_1] = 1
	BetDoubleMap[CardGroupType_Cattle_2] = 1
	BetDoubleMap[CardGroupType_Cattle_3] = 1
	BetDoubleMap[CardGroupType_Cattle_4] = 1
	BetDoubleMap[CardGroupType_Cattle_5] = 1
	BetDoubleMap[CardGroupType_Cattle_6] = 1
	BetDoubleMap[CardGroupType_Cattle_7] = 2
	BetDoubleMap[CardGroupType_Cattle_8] = 3
	BetDoubleMap[CardGroupType_Cattle_9] = 3
	BetDoubleMap[CardGroupType_Cattle_C] = 4
	BetDoubleMap[CardGroupType_Cattle_BOMB] = 5
	BetDoubleMap[CardGroupType_Cattle_WUHUA] = 6

	RandCardTypeMap = []CardGroupType{
		CardGroupType_NotCattle,
		CardGroupType_Cattle_1,
		CardGroupType_Cattle_2,
		CardGroupType_Cattle_3,
		CardGroupType_Cattle_4,
		CardGroupType_Cattle_5,
		CardGroupType_Cattle_6,
		CardGroupType_Cattle_7,
		CardGroupType_Cattle_8,
		CardGroupType_Cattle_9,
		CardGroupType_Cattle_C,
	}
}

//游戏状态
const (
	GAME_STATUS_WAITSTART   = 10 + iota // 等待开始
	GAME_STATUS_SEATBET                 // 抢坐和坐下的人下注状态
	GAME_STATUS_FACARD                  // 发牌状态
	GAME_STATUS_DOWNBTES                // 下注状态
	GAME_STATUS_OPENCARD                // 开牌状态
	GAME_STATUS_BALANCE                 // 结算
	GAME_STATUS_SHUFFLECARD             // 洗牌
)

// 发送给客户端的座位信息
type GClientSeatInfo struct {
	Id            int    // 座位id
	UserId        int64  // 用户id
	Name          string // 名称
	Avatar        string // 头像
	SeatDownCount int    // 坐下的次数
	DownBetTotal  int64  // 总下注金额
}

// 展示的用户
type GClientRankInfo struct {
	UserId int64  // 用户Id
	Avatar string // 头像
}

// 发送给客户端的桌子信息
type GClientDeskInfo struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string

	JuHao     string            // 局号
	FangHao   string            // 房号
	Seats     []GClientSeatInfo // 座位信息
	RankUsers []GClientRankInfo // 展示用户
	BetLevels []int64           // 下注级别

	MyUserAvatar string          // 用户头像
	MyUserName   string          // 用户昵称
	MyUserCoin   int64           // 用户金币
	MyDownBets   map[uint8]int64 // 自己下注的集合

	GameStatus          int                   // 游戏状态
	GameStatusDuration  int64                 // 当前状态持续时间毫秒
	CardGroupArray      map[int]CardGroupInfo // 玩家和庄家的牌 庄家牌索引最后一个4
	SeatDownMinCoinCond int                   // 坐下条件
	SeatDownMinBetCond  int                   // 坐下的人最低下注

	AreaMaxCoin         int // 限制区域最大下注
	AreaMaxCoinDownSeat int // 限制区域最大下注
	SeatUpTotalCount    int // 站立条件
}

// 游戏状态
type GClientGameStatusInfo struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string

	GameStatus         int   // 游戏状态
	GameStatusDuration int64 // 当前状态持续时间
}

// 玩家请求坐下
type GClientQSeatDown struct {
	Id      int32 //协议号
	SeatIdx int   // 座位索引
}

// 玩家请求坐下返回信息
type GClientRSeatDown struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string
}

// 玩家请求站立
type GClientQSeatUp struct {
	Id int32 //协议号
}

// 玩家请求站立回复
type GClientRSeatUp struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string
}

// 座位信息改变通知
type GClientSeatDownChange struct {
	Id            int32
	Type          int    // 0添加 1修改 2删除
	SeatId        int    // 座位号
	GameStatus    int    // 当前状态
	OldUserId     int64  // 玩家Id
	NewUserId     int64  // 玩家Id
	NewUserAvatar string // 玩家头像
	NewUserName   string // 玩家昵称
}

// 玩家请求下注信息
type GClientQDownBet struct {
	Id      int32
	SeatIdx int
	CoinIdx int
}

// 玩家请求下注返回
type GClientRDownBet struct {
	Id     int32
	Result int32 //0成功，其他失败
	Coins  int64 //玩家剩余金币
	Err    string
}

// 玩家下注通知
type GClientNDownBet struct {
	Id      int32
	Uid     int64
	SeatIdx int
	CoinIdx int
	Coin    uint32
}

type CardGroupInfo struct {
	CardGroupType CardGroupType
	Cards         []int
	MaxCard       uint8
}

// 系统发牌通知
type GClientNFaCard struct {
	Id    int32
	Cards map[int]CardGroupInfo
}

// 单个位置结算
type GClientSeatBalance struct {
	Bottom   int64
	Result   int64
	MyBottom int64
	MyResult int64
}

// 系统结算
type GClientBalance struct {
	Id        int32
	Results   map[int]GClientSeatBalance
	RankUsers []GClientRankInfo // 展示用户
	MyCoin    int64             // 用户金币
}

//
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	BankerCard  []int32           `json:"bankerCard"`
	IdleCard    [][]int32         `json:"idleCard"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId      int64    `json:"userId"`
	UserAccount string   `json:"userAccount"`
	BetCoins    int64    `json:"betCoins"`    // 下注金币
	BetArea     [4]int32 `json:"betArea"`     // 区域下注情况
	PrizeCoins  int64    `json:"prizeCoins"`  // 赢取金币
	CoinsBefore int64    `json:"coinsBefore"` // 下注前金币
	CoinsAfter  int64    `json:"coinsAfter"`  // 结束后金币
	Robot       bool     `json:"robot"`
}

const (
	GCLIENT_TIPS_NOTBET = iota + 1 // 几局未下注提示
)

// 玩家提示信息
type GClientTips struct {
	Id   int32
	Code int32
	Msg  string
}

// 请求走势回复 MSG_GAME_RHISTORY
type GClientRGameHistory struct {
	Id       int32
	Result   int32 //0成功，其他失败
	Err      string
	Historys map[int][]CardGroupType
}

type GClientManyUserInfo struct {
	Uid       int64
	NickName  string
	Avatar    string
	Coin      int64
	GameCount int32
	Victory   int32
	DownBet   int64
}

// 请求更多玩家信息回复 MSG_GAME_RMANYUSER
type GClientRManyUser struct {
	Id        int32
	Result    int32 //0成功，其他失败
	Err       string
	ManyUsers []GClientManyUserInfo
}

// 请求返回回复
type GClientRGameBack struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string
}

//大厅走势图
type GClientHallHistory struct {
	HomeName       int     //房间名称
	HomeOdds       int     //房间倍率
	LimitRed       int     //限红
	GameState      string  //状态
	Trend          [][]int //四个位置输赢走势图	 0 输，1 赢
	WinCount       []int   //四个位置赢的次数
	StageTimeCount int     //阶段总时间
	RemainTime     int     //剩余时间
}

package main

//座位玩家信息
type PlayerInfoByChair struct {
	Uid    int64  //玩家Id    (uid 为 0时代表 该座位上没有人)
	Nick   string //玩家昵称
	Avatar string //头像
	Coins  int64  //玩家金币
}

//游戏状态与持续时间
type StatuAndTimes struct {
	ShuffleId       int
	ShuffleMs       int
	StartdownbetsId int
	StartdownbetsMs int
	DownBetsId      int
	DownBetsMS      int
	StopdownbetsId  int
	StopdownbetsMs  int
	FaCardId        int
	FaCardMS        int
	OpenCardId      int
	OpenCardMS      int
	BalanceId       int
	BalanceMS       int
}

//阶段信息
type GameStatuInfo struct {
	Id         int
	Status     int //当前状态
	StatusTime int //状态持续时间
}

// 玩家历史
type BetHistory struct {
	IsVictory bool  // 是否胜利
	DownBet   int64 // 下注的金额
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

	AreaCoin           []int64               //区域金币
	GameStatus         int                   // 游戏状态
	GameStatusDuration int64                 // 当前状态持续时间毫秒
	CardGroupArray     map[int]CardGroupInfo // 玩家和庄家的牌 庄家牌索引最后一个4

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

//手牌信息
type CardGroupInfo struct {
	CardGroupType CardGroupType //牌型
	Cards         []int         //牌组
	MaxCard       uint8         //最大牌
}

// //局号变更通知结构体
// type JuHaoChange struct {
// 	Id    int
// 	JuHao string
// }

//下注请求附带结构体
type DownBet struct {
	ChipIndex int //筹码Id
	AreaIndex int //区域Id
}

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

//下注广播结构体
type DownBetBro struct {
	Id           int     //协议号
	DownBet      []int64 //总下注
	MyDownBet    []int64 //自己下注情况
	OtherDownBet []int64 //其他玩家下注情况
}

//
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	BankerCard  []int             `json:"bankerCard"`
	IdleCard    [][]int           `json:"idleCard"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId      int64   `json:"userId"`
	UserAccount string  `json:"userAccount"`
	BetCoins    int64   `json:"betCoins"`    // 下注金币
	BetArea     []int64 `json:"betArea"`     // 区域下注情况
	PrizeCoins  int64   `json:"prizeCoins"`  // 赢取金币
	CoinsBefore int64   `json:"coinsBefore"` // 下注前金币
	CoinsAfter  int64   `json:"coinsAfter"`  // 结束后金币
	Robot       bool    `json:"robot"`
}

type ToClientBalance struct {
	Id            int     //协议号
	MyCoins       int64   //玩家金币
	MyResult      []int64 //玩家自身结算集合
	AllResult     []int64 //所有区域结算
	MyGetCoins    int64   //玩家自身盈利
	Zoushi        [][]Pzshi
	BetAbleIndex  int
	AreaWinDouble []int //倍数
	WinArea       []int //胜利区域
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

//请求更多玩家应答
type MorePlayer struct {
	Id        int               //协议
	PlayerMsg []PlayerMsgByMore //更多玩家信息
}
type PlayerMsgByMore struct {
	Head       string //头像
	Nick       string //昵称
	Coins      int64  //金币
	MatchCount int    //记录局数
	BetAll     int64  //总下注
	WinCount   int    //获胜局数
}

//玩家请求离开应答
type GoBackReply struct {
	Id     int    //协议
	Result int    //0为成功，其他为失败
	Err    string //成功时为空
}

type MsgWarning struct {
	Id     int //协议
	Result int //1，未下注警告 2,未下注踢出  3，本局已下注  无法退出
}

package main

//座位玩家信息
type PlayerInfoByChair struct {
	Uid    int64  //玩家Id    (uid 为 0时代表 该座位上没有人)
	Nick   string //玩家昵称
	Avatar string //头像
	Coins  int64  //玩家金币
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

	AreaCoin           []int64 //区域金币
	GameStatus         int     // 游戏状态
	GameStatusDuration int64   // 当前状态持续时间毫秒

	AreaMaxCoin int64 // 限制区域最大下注

	ChairList    []PlayerInfoByChair //座位玩家信息
	History      []int               //历史记录
	BetAbleIndex int                 //玩家可下注筹码id
	Zhuang       PlayerMsg           //当前庄家
	WaitList     []PlayerMsg         //等待玩家信息
}

//用户信息
type PlayerMsg struct {
	Uid          int64  //玩家uid
	MyUserAvatar string // 用户头像
	MyUserName   string // 用户昵称
	MyUserCoin   int64  // 用户金币
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
	IsZhuang     bool    //该玩家是否为庄
}

//
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId      int64   `json:"userId"`
	UserAccount string  `json:"userAccount"`
	Robot       bool    `json:"gradeNumber"`
	BetCoins    int64   `json:"betCoins"`    // 下注金币
	BetArea     []int64 `json:"betArea"`     // 区域下注情况
	PrizeCoins  int64   `json:"prizeCoins"`  // 赢取金币
	CoinsBefore int64   `json:"coinsBefore"` // 下注前金币
	CoinsAfter  int64   `json:"coinsAfter"`  // 结束后金币
	Banker      bool    `json:"banker"`      // 玩家是否为庄
}

type ToClientBalance struct {
	Id              int     //协议号
	MyCoins         int64   //玩家金币
	MyResult        []int64 //玩家自身结算集合
	AllResult       []int64 //所有区域结算
	OtherResult     []int64 //所有其他玩家结算
	MyGetCoins      int64   //玩家自身盈利
	History         []int   //开奖历史
	BetAbleIndex    int
	AreaWinDouble   int   //倍数
	WinArea         int   //胜利区域
	IsZhuang        bool  //是否为庄
	ZhuangWinOrLose int64 //庄家输赢
	ZhuangCoins     int64 //庄家金币
	IsHasZhuang     bool  //本局是否有庄
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
	Result int //1，未下注警告 2,未下注踢出  3，本局已下注 无法退出   4,金币不足 下庄  5,局数已经到了 下庄  6，金币不足上庄，将其移除队列  7，您现在是庄家 无法退出
}

type ChangZhuangReply struct {
	Id     int    //协议号
	Result int    //0成功，其他失败
	Err    string //错误信息，成功时为0
}
type ZhuangInfo struct {
	Id       int         //协议号
	Info     PlayerMsg   //庄家信息
	WaitList []PlayerMsg //等待玩家上庄列表
}

//发送给机器人开奖结果
type ToRobot struct {
	Id    int
	Index int //开奖区域
}

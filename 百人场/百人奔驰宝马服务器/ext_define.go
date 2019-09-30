package main

type GBetCountInfo struct {
	Id           int   // 座位id
	DownBetValue int64 // 总下注金额
	UserBetValue int64 // 用户下注金额
}

type GClientBetCountInfo struct {
	Id           int   // 座位id
	DownBetTotal int64 // 总下注金额
}

// 发送给客户端的桌子信息
type GClientDeskInfo struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string

	JuHao     string                // 局号
	FangHao   string                // 房号
	Bets      []GClientBetCountInfo // 座位信息
	BetLevels []int64               // 下注级别

	PlayerMassage PlayerMsg //用户信息

	AreaCoin           []int64 //区域金币
	GameStatus         int     // 游戏状态
	GameStatusDuration int64   // 当前状态持续时间毫秒

	AreaMaxCoin int // 限制区域最大下注
}

// 游戏状态
type GSGameStatusInfo struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string

	GameStatus         int   // 游戏状态
	GameStatusDuration int64 // 当前状态持续时间
}

type PlayerMsg struct {
	Uid          int64   //玩家uid
	MyUserAvatar string  // 用户头像
	MyUserName   string  // 用户昵称
	MyUserCoin   int64   // 用户金币
	MyDownBets   []int64 // 自己下注的集合
}

//玩家请求下注
type GADownBet struct {
	Id      int
	BetsIdx int // 下注区域索引(0-7)
	CoinIdx int // 下注金额索引(0-4)
}

//玩家请求下注返回
type GSDownBet struct {
	Id     int
	Result int32 // 0 成功，其他失败
	Err    string
}

//玩家请求下注通知
type GNDownBet struct {
	Id             int
	Uid            int64 // 玩家uid
	BetIdx         int   // 区域索引
	Coin           int64 // 下注金币数
	CoinIdx        int   // 下注金币索引
	AreaCoin       int64 // 下注区域金币数
	PlayerAreaCoin int64 // 玩家在该下注区域金币数
}

//开奖通知
type GNLottery struct {
	Id  int
	Car int // 开奖结果
}

// 系统结算
type GNBalance struct {
	Id      int32
	Results map[int]GBetBalance
	MyCoin  int64 // 用户金币
}

// 单个位置结算
type GBetBalance struct {
	Bottom   int64
	Result   int64 //0 lose 1 win
	MyBottom int64 //玩家下注金币数
	MyResult int64
}

// 玩家提示信息
type GSTips struct {
	Id  int32
	Msg string
}

//游戏开奖记录通知
type GNRecord struct {
	Id              int
	Record          []int // 游戏开奖记录
	OnlinePlayerNum int   // 在线玩家数
}

//
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	BankerCard  int               `json:"bankerCard"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId      int64    `json:"userId"`
	UserAccount string   `json:"userAccount"`
	BetCoins    int64    `json:"betCoins"`    // 下注金币
	BetArea     [8]int32 `json:"betArea"`     // 区域下注情况
	PrizeCoins  int64    `json:"prizeCoins"`  // 赢取金币
	CoinsBefore int64    `json:"coinsBefore"` // 下注前金币
	CoinsAfter  int64    `json:"coinsAfter"`  // 结束后金币
	Robot       bool     `json:"robot"`
}

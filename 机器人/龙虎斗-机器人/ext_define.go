package main

// 定义游戏消息
const (
	MSG_GAME_INFO_READY_NOTIFY    = 410002 // 准备中410002
	MSG_GAME_INFO_BET_NOTIFY      = 410004 // 下注中410004
	MSG_GAME_INFO_BET             = 410005 // 下注410005
	MSG_GAME_INFO_STOP_BET_NOTIFY = 410008 // 停止下注410008
	MSG_GAME_INFO_AWARD_NOTIFY    = 410010 // 派奖中410010
	MSG_GAME_INFO_RECONNECT_REPLY = 410014 // 410014房间信息
	MSG_GAME_INFO_INTO            = 410020 // 410020请求AUTO_REPLY剩余消息
)

// 下注区域
const (
	// 龙、虎、和
	INDEX_DRAGON int = 1 + iota
	INDEX_TIGER
	INDEX_DRAW

	// 龙   方、梅、红、黑
	INDEX_DRAGONSPADE
	INDEX_DRAGONPLUM
	INDEX_DRAGONRED
	INDEX_DRAGONBLOCK

	// 虎   方、梅、红、黑
	INDEX_TIGERSPADE
	INDEX_TIGERPLUM
	INDEX_TIGERRED
	INDEX_TIGERBLOCK

	// 庄赢、庄输
	INDEX_BANKERWIN
	INDEX_BANKERLOSE

	// 错误下标
	INDEX_ERROR
)

//-------------龙虎斗--------------

// 座位信息
type GSInfo struct {
	Nick string
	Head string
}

type GInfoReConnectReply struct {
	Id         int32    // 协议号
	GameState  int32    // 游戏状态
	RoomId     string   // 房号
	GameId     string   // 局号
	GameLimit  []int64  // 限红
	BetList    []int64  // 下注金币限制
	LeftCount  int32    // 剩余牌
	RightCount int32    // 废牌
	GameCount  int32    // 当前牌局
	TAreaCoins []int64  // 总下注金币
	SeatList   []GSInfo // 座位玩家
	PAreaCoins []int64  // 当前玩家下注金币
	PCoins     int64    // 当前玩家当前金币
	BetArea    []bool   // 可下注区域
	OpenCard   int32    // 翻开的牌
	DragonCard int32    // 龙牌
	TigerCard  int32    // 虎牌
	WinArea    []bool   // 赢取区域
	RunChart   []int32  // 走势
	Timer      int32    // 阶段时间（毫秒）
}

// 游戏下注
type GGameBetNotify struct {
	Id      int32  // 协议号
	Timer   int32  // 阶段时间（毫秒）
	BetArea []bool // 可下注区域
}

// 游戏派奖
type GGameAwardNotify struct {
	Id           int32   // 协议号
	Timer        int32   // 阶段时间（毫秒）
	BankRunChart []int32 // 庄走势
	WinArea      []bool  // 区域输赢情况
	TWinArea     []int64 // 总输赢
	SeatWinCoins []int64 // 座位玩家输赢
	OtherWinArea []int64 // 其他玩家输赢
	PWinArea     []int64 // 自己输赢
	PWin         int64   // 自己赢取值
	PCoins       int64   // 自己最终金币
}

// 游戏准备
type GGameReadyNotify struct {
	Id         int32    // 协议号
	LeftCount  int32    // 剩余牌
	RightCount int32    // 废牌
	GameCount  int32    // 当前局数
	Timer      int32    // 阶段时间（毫秒）
	SeatList   []GSInfo // 座位玩家
	GameId     string   // 局号
}

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
	// 东西南北
	INDEX_DONG int = 1 + iota
	INDEX_XI
	INDEX_NAN
	INDEX_BEI

	// 对、顺、金
	INDEX_BANKER_DOUBLE
	INDEX_BANKER_SHUNZI
	INDEX_BANKER_JINHUA

	// 顺金、炸弹
	INDEX_BANKER_SHUNJIN
	INDEX_BANKER_BAOZI

	// 错误下标
	INDEX_ERROR
)

//-------------百人炸金花--------------

// 座位信息
type GSInfo struct {
	Nick string
	Head string
}

type GInfoReConnectReply struct {
	Id           int32     // 协议号
	GameState    int32     // 游戏状态
	RoomId       string    // 房号
	GameId       string    // 局号
	GameLimit    []int64   // 限红
	BetList      []int64   // 下注金币限制
	TAreaCoins   []int64   // 总下注金币
	SeatList     []GSInfo  // 座位玩家
	PAreaCoins   []int64   // 当前玩家下注金币
	PCoins       int64     // 当前玩家当前金币
	BankerCard   []int32   // 庄牌
	MBankerCard  []int32   // 庄最大牌型
	TBankerCard  int32     // 庄牌型
	IdleCard     [][]int32 // 东西南北牌
	MIdleCard    [][]int32 // 闲最大牌型（东西南北）
	TypeList     []int32   // 牌型（东西南北）
	WinArea      []bool    // 赢取区域
	BankRunChart []int32   // 庄走势
	Timer        int32     // 阶段时间（毫秒）
	Dices        []int     // 骰子
}

// 游戏下注
type GGameBetNotify struct {
	Id    int32 // 协议号
	Timer int32 // 阶段时间（毫秒）
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
	Id       int32    // 协议号
	Timer    int32    // 阶段时间（毫秒）
	SeatList []GSInfo // 座位玩家
	GameId   string   // 局号
	Dices    []int    // 骰子
}

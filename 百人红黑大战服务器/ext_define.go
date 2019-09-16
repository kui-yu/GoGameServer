package main

const (
	MSG_GAME_INFO_START            = 410000 + iota // 游戏开始额外信息
	MSG_GAME_INFO_SHUFFLE_NOTIFY                   // 洗牌中410001
	MSG_GAME_INFO_READY_NOTIFY                     // 准备中410002
	MSG_GAME_INFO_SEND_NOTIFY                      // 发牌中410003
	MSG_GAME_INFO_BET_NOTIFY                       // 下注中410004
	MSG_GAME_INFO_BET                              // 下注410005
	MSG_GAME_INFO_BET_REPLY                        // 下注反馈410006
	MSG_GAME_INFO_BET_NEW_NOTIFY                   // 新下注410007
	MSG_GAME_INFO_STOP_BET_NOTIFY                  // 停止下注410008
	MSG_GAME_INFO_OPEN_NOTIFY                      // 开牌中410009
	MSG_GAME_INFO_AWARD_NOTIFY                     // 派奖中410010
	MSG_GAME_INFO_UNDO_NOTIFY                      // 410011用户超过多少局没下注通知
	MSG_GAME_INFO_EXIT_NOTIFY                      // 410012强制用户离开通知
	MSG_GAME_INFO_AUTO_REPLY                       // 410013座位信息
	MSG_GAME_INFO_RECONNECT_REPLY                  // 410014房间信息
	MSG_GAME_INFO_END_NOTIFY                       // 410015游戏结束应答
	MSG_GAME_INFO_RUN_CHART                        // 410016游戏走势
	MSG_GAME_INFO_RUN_CHART_REPLY                  // 410017游戏走势反馈
	MSG_GAME_INFO_USER_LIST                        // 410018游戏其他玩家
	MSG_GAME_INFO_USER_LIST_REPLY                  // 410019游戏其他玩家反馈
	MSG_GAME_INFO_INTO                             // 410020请求AUTO_REPLY剩余消息
	MSG_GAME_INFO_EXIT                             // 410021请求退出房间
	MSG_GAME_INFO_EXIT_REPLY                       // 410022退出应答
	MSG_GAME_INFO_EXIT_LIMIT_LOW                   // 410023 金币过少
	MSG_GAME_INFO_EXIT_LIMIT_HIGHT                 // 410024 金币过多
	MSG_GAME_INFO_GET_RECORD                       // 410025 获取游戏记录
	MSG_GAME_INFO_GET_RECORD_REPLY                 // 410026 游戏结束后房间的游戏记录
)
const (
	ERR_AREAID    = -1 //下注区域错误
	ERR_COINID    = -2 //下注金额错误
	ERR_LIMITCOIN = -3 //下注金额上限
	// ERR_SPECBET   = -4 //下注幸运一击区域错误
)

//座位信息
type GSInfo struct {
	Nick  string //昵称
	Head  string //头像
	Coins int64  //金币
}

//用户信息
type GUserInfo struct {
	Uid       int64
	Nick      string //昵称
	Head      string //头像
	TotBet    int64  //总下注
	WinCount  int64  //赢取次数
	Coins     int64  //当前金币
	UserCount int    //用户人数
}

//走势图
type GRunChartReply struct {
	Id              int32
	PRunchart       []int32 //红黑走势
	CardTypeChart   []int   //牌型记录走势
	RunChartTwenty  []int32 //近20局红黑走势
	ChartCount      int     //游戏局数
	ChartRedCount   int     //走势中红方赢的局数
	ChartBlackCount int     //走势中黑方赢的局数
	Rcount          int     //近20局红方赢的局数
	Bcount          int     //近20局黑方赢的局数
}

//退出反馈
type GGameExitReply struct {
	Id     int32
	Result int32
}

//玩家列表
type GUserInfoReply struct {
	Id        int32
	UserInfo  []GUserInfo
	UserCount int
}

// 游戏洗牌
type GGameShuffleNotify struct {
	Id        int32 // 协议号
	Timer     int32 // 阶段时间（毫秒）
	GameCount int32 // 当前局数
}

// 游戏准备
type GGameReadyNotify struct {
	Id          int32    // 协议号
	GameCount   int32    // 当前局数
	Timer       int32    // 阶段时间（毫秒）
	SeatList    []GSInfo // 座位玩家
	GameId      string   // 局号
	Limitcoinid int32    //下注限制ID
}

//游戏发牌
type GGameSendCardNotify struct {
	Id             int32   // 协议号
	Timer          int32   // 阶段时间（毫秒）
	RedCard        []int32 //红方牌
	RedCardColor   []int32 //红方牌花色
	RGrade         int     //红方牌等级
	BlackCard      []int32 //黑方牌
	BlackCardColor []int32 //黑方牌花色
	BGrade         int     //黑方牌等级

}

// 游戏下注定时广播
type GGameTotBetNotify struct {
	Id         int32   // 协议号
	TAreaCoins []int64 // 区域总下注
	PAreaCoins []int64 // 自己区域总下注
	OAreaCoins []int64 // 其他玩家区域新下注
}

// 游戏下注
type GGameBetNotify struct {
	Id      int32  // 协议号
	Timer   int32  // 阶段时间（毫秒）
	BetArea []bool // 可下注区域
}

// 下注
type GGameBet struct {
	Id     int32 // 协议号
	MsgId  int32 // 消息系号，防止重复(新开局系号1开始（断线重连也一样）)
	AreaId int32 // 下注区域Id
	CoinId int32 // 下注金额Id
}

// 下注应答
type GGameBetReply struct {
	Id          int32   // 协议号
	MsgId       int32   // 反馈消息
	PAreaCoins  []int64 // 自己区域总下注
	Coins       int64   // 自己剩余金币
	AreaId      int32   // 下注区域
	CoinId      int32   // 下注金额Id
	LimitCoinId int32   //当用户金币小于下注池的某一值时，获得此值得ID
}

// 新下注广播
type GGameNewBetNotify struct {
	Id           int32     // 协议号
	SeatBetList  [][]int64 // 座位玩家下注情况
	OtherBetList []int64   // 除自己以外，其他玩家下注情况
	PAreaCoins   []int64   // 自己总下注情况
	TAreaCoins   []int64   // 区域总下注
}

// 游戏停止下注
type GGameStopBetNotify struct {
	Id           int32     // 协议号
	Timer        int32     // 阶段时间（毫秒）
	TAreaCoins   []int64   // 区域总下注
	PAreaCoins   []int64   // 自己区域总下注
	SeatBetList  [][]int64 // 座位玩家下注情况
	OtherBetList []int64   // 其他玩家区域新下注
}

//游戏开牌
type GGameOpenNotify struct {
	Id        int32   //协议号
	Timer     int32   //阶段时间
	RedCard   []int32 //红方牌
	RedType   int     //红方牌类型
	BlackCard []int32 //黑方牌
	BlackType int     //黑方牌类型
}

//游戏派奖
type GGameAwardNotify struct {
	Id            int32   // 协议号
	Timer         int32   // 阶段时间（毫秒）
	RunChart      []int32 // 输赢走势
	CardTypeChart []int   // 牌型记录走势
	WinArea       []bool  // 区域输赢情况
	WinAreaId     int32   // 红黑方赢得ID值
	TWinArea      []int64 // 总输赢
	SeatWinCoins  []int64 // 座位玩家输赢
	OtherWinArea  []int64 // 其他玩家输赢
	PWinArea      []int64 // 自己输赢
	PWin          int64   // 自己赢取值
	PCoins        int64   // 自己最终金币
	PrizeCoins    int64   // 真实输赢
	Rtype         int     //红方牌型
	Btype         int     //黑方牌型
	Limitcoinid   int32   //下注限制ID
	PairValue     int32   //对子的值
}

//断线重连失败，需要清除游戏ID
type GReConnectFailedNotify struct {
	Id int32
}

type GWinCard struct {
	Index    []int
	WinScale float64
}

//游戏重连信息返回
type GInfoReConnectReply struct {
	Id            int32    // 协议号
	GameState     int32    // 游戏状态
	RoomId        string   // 房号
	GameId        string   // 局号
	GameLimit     int64    // 限红
	BetList       []int64  // 下注金币限制
	GameCount     int32    // 当前牌局
	TAreaCoins    []int64  // 总下注金币
	SeatList      []GSInfo // 座位玩家
	PAreaCoins    []int64  // 当前玩家下注金币
	PCoins        int64    // 当前玩家当前金币
	BetArea       []bool   // 可下注区域
	CardList      []int32  // 展示牌列表
	RedCard       []int32  // 红方牌
	BlackCard     []int32  // 黑方牌
	WinArea       []bool   // 赢取区域
	RunChart      []int32  // 输赢走势
	CardTypeChart []int    // 牌型记录走势
	Rtype         int      //红方牌型
	Btype         int      //黑方牌型
	Timer         int32    // 阶段时间（毫秒）
	LimitCoinId   int32    //当用户金币小于下注池的某一值时，获得此值得ID
}

//游戏开始自由匹配阶段返回
type GInfoAutoGameReply struct {
	Id        int32    // 协议号
	PlayerNum int32    // 房间人数
	SeatList  []GSInfo // 座位玩家
}

//游戏记录
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`  //游戏编号
	GradeId     int               `json:"gradeId"` //游戏级别ID
	RoomId      int               `json:"roomId"`  //房间编号
	GradeNumber int               `json:"gradeNumber"`
	GameRoundNo string            `json:"gameRoundNo"` //游戏回合数
	RedCard     []int32           `json:"redCard"`     //红方牌
	BlackCard   []int32           `json:"blackCard"`   //黑方牌
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}
type GGameRecordInfo struct {
	UserId      int64   `json:"userId"`
	UserAccount string  `json:"userAccount"`
	Robot       bool    `json:"robot"`
	BetCoins    int64   `json:"betCoins"`    // 下注金币总额
	BetArea     []int64 `json:"betArea"`     // 区域下注情况
	PrizeCoins  int64   `json:"prizeCoins"`  // 赢取金币
	CoinsBefore int64   `json:"coinsBefore"` // 开始时候多少金币
	CoinsAfter  int64   `json:"coinsAfter"`  // 结束时候多少金币

}

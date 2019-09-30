package main

//阶段
const (
	STAGE_DEAL   = 11 //发牌阶段
	STAGE_PLAY   = 12 //玩牌阶段
	STAGE_SETTLE = 14 //结算阶段
)

//消息定义
const (
	MSG_GAME_INFO_STAGE       = 410000 //阶段消息
	MSG_GAME_INFO_CALL        = 410001 //通比牛牛 叫分
	MSG_GAME_INFO_CALL_REPLY  = 410002 //通比牛牛 叫分返回
	MSG_GAME_INFO_DEAL_REPLY  = 410003 //通比牛牛 发牌
	MSG_GAME_INFO_PLAY_REPLY  = 410004 //通比牛牛 结果
	MSG_GAME_INFO_PLAY        = 410005 //通比牛牛 开牌
	MSG_GAME_INFO_SETTLE      = 410006 //通比牛牛 结算
	MSG_GAME_INFO_BET_LIST    = 410007 //筹码列表
	MSG_GAME_INFO_ROOM_NOTIFY = 410009 //房间信息通知410009
	MSG_GAME_INFO_RECONNECT   = 410010 //重连
	MSG_GAME_INFO_AUTO_REPLY  = 410012 //410012,游戏随机匹配成功的数据
)

//阶段时间
const (
	TIME_STAGE_ZERO     = 1
	TIME_STAGE_ZERO_NUM = 0
	//匹配时间
	TIME_STAGE_INIT     = 1
	TIME_STAGE_INIT_NUM = 10 //10秒中匹配
	//叫分时间
	TIME_STAGE_CALL     = 2
	TIME_STAGE_CALL_NUM = 12
	//开牌动画
	TIME_STAGE_START     = 3
	TIME_STAGE_START_NUM = 3
	//玩牌
	TIME_STAGE_PLAYCARD     = 4
	TIME_STAGE_PLAYCARD_NUM = 10

	TIMER_OVER     = 8
	TIMER_OVER_NUM = 1 //归还桌子时间
)

//所有匹配玩家信息
type GSInfoAutoGame struct {
	Id   int32
	Seat []GSeatInfo
}

//筹码列表
type GCallListMsg struct {
	Id         int
	BetListCnt int   //下注数量
	BetList    []int //下注列表
}

//玩家叫分
type GCallMsg struct {
	Id       int32
	Multiple int
}

//玩家叫分返回
type GCallMsgReply struct {
	Id       int32
	ChairId  int32
	Multiple int
}

//阶段时间
type GStageInfo struct {
	Id        int32
	Stage     int32
	StageTime int32
}

//玩家点击结果
type GPlayCard struct {
	Id    int32
	Value int32
}

//返回结果
type GHandNiuReply struct {
	Id       int32
	ChairId  int32
	NiuPoint int32
	NiuCards []int32
}

type GWinInfo struct {
	Uid     int64
	ChairId int32
	WinCoin int64
	Coins   int64
}

type GWinInfosReply struct {
	Id         int32
	WinChairId int32
	InfoCount  int32
	Infos      []GWinInfo
}

type GTableInfoReply struct {
	Id      int32
	TableId string //房间号
	BScore  int    //底分
}

//座位信息-重连
type GSeatInfoReconnect struct {
	Id           int32
	ChairIds     []int32   //所有玩家位置[0,1,2,3,4]
	States       []int32   //所有玩家状态[1,1,1,1,2]
	Multiples    []int     //所有玩家倍数[0,0,1,2,3]
	PlayNum      int32     //已出牌玩家数量
	PlayChairIds []int32   //已出牌玩家位置[0,2,3]
	PlayPoints   []int32   //已出牌结果[8,6,7]
	PlayCards    [][]int32 //已出牌玩家手牌[[1,2,3,4,5],[2,3,4,5,6],[5,5,5,5,5]]
	MyCard       []int32   //手牌
	Stage        int
	StageTime    int
	BetListCnt   int   //下注数量
	BetList      []int //下注列表
}

//游戏记录
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

//游戏详细记录
type GGameRecordInfo struct {
	UserId        int64  `json:"userId"`
	UserAccount   string `json:"userAccount"`
	Robot         bool   `json:"robot"`
	CoinsBefore   int64  `json:"coinsBefore"`
	BetCoins      int64  `json:"betCoins"`
	Coins         int64  `json:"coins"`
	CoinsAfter    int64  `json:"coinsAfter"`
	Cards         []int  `json:"cards"`
	BrandMultiply int    `json:"brandMultiple"`
	BetMultiple   int    `json:"betMultiple"`
	Multiple      int    `json:"multiple"`
	Score         int    `json:"score"`
}

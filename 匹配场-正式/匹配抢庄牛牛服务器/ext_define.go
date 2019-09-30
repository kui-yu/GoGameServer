package main

//阶段
const (
	STAGE_BET    = 10 //抢庄阶段
	STAGE_DEAL   = 11 //发牌阶段
	STAGE_PLAY   = 12 //玩牌阶段
	STAGE_SETTLE = 14 //结算阶段
)

//消息定义
const (
	MSG_GAME_INFO_STAGE              = 410000 //阶段消息
	MSG_GAME_INFO_CALL               = 410001 //通比牛牛 叫分
	MSG_GAME_INFO_CALL_REPLY         = 410002 //通比牛牛 叫分返回
	MSG_GAME_INFO_DEAL_REPLY         = 410003 //通比牛牛 发牌
	MSG_GAME_INFO_PLAY_REPLY         = 410004 //通比牛牛 结果
	MSG_GAME_INFO_PLAY               = 410005 //通比牛牛 开牌
	MSG_GAME_INFO_SETTLE             = 410006 //通比牛牛 结算
	MSG_GAME_INFO_CALL_BANKER        = 410007 //抢庄
	MSG_GAME_INFO_CALL_BANKER_NOTIFY = 410008 //抢庄通知
	MSG_GAME_INFO_ROOM_NOTIFY        = 410009 //房间信息通知410009
	MSG_GAME_INFO_RECONNECT          = 410010 //重连
	MSG_GAME_INFO_CHOOSE_BANKER      = 410011 //选庄通知
	MSG_GAME_INFO_AUTO_REPLY         = 410012 //410012,游戏随机匹配成功的数据
	MSG_GAME_INFO_CALL_LIST          = 410013 //庄家筹码
)

//阶段时间
const (
	TIME_STAGE_ZERO_NUM     = 0
	TIME_STAGE_CALL_NUM     = 15 //叫分时间
	TIME_STAGE_BET_NUM      = 15 //下注时间
	TIME_STAGE_START_NUM    = 3  //发牌动画
	TIME_STAGE_PLAYCARD_NUM = 10 //玩牌
)

const (
	TIMER_OVER     = 8
	TIMER_OVER_NUM = 1 //归还桌子时间
)

//所有匹配玩家信息
type GSInfoAutoGame struct {
	Id   int32
	Seat []GSeatInfo
}

type GSCallList struct {
	Id          int
	CallListCnt int
	CallList    []int
}

//玩家下注
type GCallMsg struct {
	Id       int32
	Multiple int
}

//玩家下注返回
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
	Id      int32
	ChairId int32
}

//玩家抢庄结果
type GCallBankReply struct {
	Id              int32
	Banker          int32
	BankerList      []int
	BankerMultiples int   //庄家倍数
	BetListCnt      int   //下注数量
	BetList         []int //下注列表
}

//返回发牌结果
type GHandNiuReply struct {
	Id       int32
	ChairId  int32
	NiuPoint int32
	NiuCards []int32
}

//结算子集
type GWinInfo struct {
	Uid      int64
	ChairId  int32
	WinCoin  int64
	Coins    int64
	NiuPoint int32
	NiuCards []int32
}

//结算结果
type GWinInfosReply struct {
	Id        int32
	InfoCount int32
	Infos     []GWinInfo
}

type GTableInfoReply struct {
	Id      int32
	TableId string //房间号
	BScore  int    //底分
}

//座位信息-重连
type GSeatInfoReconnect struct {
	Id              int32
	ChairIds        []int32 //所有玩家位置[0,1,2,3,4]
	States          []int32 //所有玩家状态[1,1,1,1,2]
	Multiples       []int   //所有玩家倍数[0,0,1,2,3]
	Banker          int32   //庄家
	BankerMultiples int     //庄家倍数
	CallMultiples   []int   //玩家叫庄倍数 -1,没叫；0,不抢
	PlayNum         int32   //已出牌玩家数量
	PlayChairIds    []int32 //已出牌玩家位置[0,2,3]
	MyCard          []int32 //手牌
	Stage           int
	StageTime       int
	CallListCnt     int
	CallList        []int
	BetListCnt      int   //下注数量
	BetList         []int //下注列表
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
	Banker        int    `json:"banker"`
	Score         int    `json:"score"`
}

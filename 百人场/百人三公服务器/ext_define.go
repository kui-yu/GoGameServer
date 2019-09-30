package main

//请求是否翻倍
type GAChangedouble struct {
	Id         int
	DoubleMode int //0平倍，1翻倍
}

//返回是否翻倍
type GSChangedouble struct {
	Id         int
	DoubleMode int
	ErrStr     string //错误
	Code       int    //0没错，1有错
	BetArrAble int    //可下注筹码
}

//阶段消息
type GSGameStageInfo struct {
	Id    int
	Stage int
	Time  int
}

//玩家下注信息
type GABetInfo struct {
	Id    int
	Place int
	Coin  int
}

//返回玩家下注
type GSPlayerBet struct {
	Id         int
	ErrStr     string
	Code       int
	CoinIndex  int   //下注金币
	PlaceIndex int   //下注区域
	MeAllCoin  int64 //自己累积下注的金币
	AllCoin    int64 //区域累积下注的金币
	Coins      int64 //自己本身金币
	BetArrAble int
}

//玩家下注群发
type GSPlayerBetMass struct {
	Id      int
	Place   []int64 //下注区域,黑、红、梅、方
	AllCoin []int64 //区域累积下注
}

//每秒群发
type SecondToCli struct {
	Desk    *ExtDesk
	Player  *ExtPlayer
	Uid     int64
	AllCoin []int64
	Place   []int64
}

//开牌结果
type GSGameResult struct {
	Id          int
	Time        int
	CardsResult []Card //索引0是庄家手牌
}

//区域输赢信息
type AreaRes struct {
	Multiple int64
	WinCoins int64
}

//结算信息
type GSSettleInfo struct {
	Id          int
	Time        int
	WinCoins    int64     //输赢的金币
	Coin        int64     //结算后的玩家金币
	RoundResult []AreaRes //区域输赢，黑、红、梅、方
	CoinResult  []int64   //自己区域输赢
	GameTrend   [][]Trend //游戏走势，根据索引分别为庄、黑、红、梅、方
	BetArrAble  int
	AreaRes     []int //黑红梅方4区域的输赢，-1为输，1为赢
}

//玩家进场=>客户端
type GSPlayerIn struct {
	Id          int
	GameId      string       //局号
	MaxBet      int64        //限红
	ManyPlayer  []ManyPlayer //桌面围观玩家
	AllPlayer   int          //在线玩家
	Stage       int          //当前游戏阶段
	GameTrend   [][]Trend    //游戏走势，根据索引分别为庄、黑、红、梅、方， Trend.Player 0为闲家，1为庄家
	DeskInfos   DeskInfo
	PlayerInfos PlayerInfo
}
type DeskInfo struct {
	Time        int          //房间的倒计时时间
	BetArr      []int64      //下注筹码
	PlaceBetAll []int64      //区域下注，0到3分别是黑红梅方
	HandCards   []Card       //庄家闲家手牌，索引0为庄家手牌
	Players     []ManyPlayer //桌面的6个玩家
}
type PlayerInfo struct {
	Account    string  // 账号
	Uid        int64   // 用户ID
	Head       string  // 头像
	Coin       int64   //玩家金币
	IsDouble   bool    //是否翻倍
	BetArrAble int     //可下注筹码
	PlaceBet   []int64 //自己区域下注，0到3分别是黑红梅方
}

//桌面玩家
type ManyPlayer struct {
	Head            string //头像
	Account         string //账号
	Uid             int64
	Coins           int64 //金币
	Round           int   //累积下注局数
	AccumulateBet   int64 //累积下注
	AccumulateCoins int64 //累积输赢
}

//桌面玩家
type GSManyPlayer struct {
	Id        int
	Players   []ManyPlayer
	AllPlayer int //总玩家数
	JuHao     string
}

//请求返回
type GSGameBack struct {
	Id     int
	Result int //0成功
	Err    string
}

//返回玩家每局输赢信息
// type GSRoundSettleInfo struct {
// 	Id         int
// 	SettleInfo []RoundSettleInfo
// }

//每局输赢信息
type RoundSettleInfo struct {
	GradeType int    //游戏场次
	AllBet    int64  //总下注
	WinCoins  int64  //输赢金币
	BetArea   []Area //押注区域
	CardType  int    //牌型
	Time      string //阶段时间
}
type Area struct {
	BetArea  int //0,1,2,3分别对应，黑红梅方
	CardType int //牌型
}

//请求玩家列表//410016
type GSPlayerList struct {
	Id         int
	PlayerInfo []ManyPlayer
}

//发送消息给大厅去记录游戏记录
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GradeNumber int               `json:"gradeNumber"`
	GameRoundNo string            `json:"gameRoundNo"`
	BankerCard  []int             `json:"bankerCard"`
	IdleCard    [][]int           `json:"idleCard"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId       int64   `json:"userId"`
	UserAccount  string  `json:"userAccount"`
	Robot        bool    `json:"robot"`
	CoinsBefore  int64   `json:"coinsBefore"`
	BetCoins     int64   `json:"betCoins"` // 下注金币总额
	PrizeCoins   int64   `json:"prizeCoins"`
	CoinsAfter   int64   `json:"coinsAfter"`
	MultipleType int     `json:"multipleType"`
	BetArea      []int64 `json:"betArea"` // 区域下注情况
}

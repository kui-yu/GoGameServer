package main

//命名规则 G + S(Send)/A(Accept) + 自定义
//所有匹配玩家信息
type GSInfoAutoGame struct {
	Id   int32
	Seat []GSSeatInfo
}

//座位信息
type GSSeatInfo struct {
	Uid     int64
	Nick    string
	Ready   bool
	Cid     int32 //椅子号
	Sex     int32
	Head    string
	Lv      int32
	Coin    int64
	IsReady int //是否准备
}

//房卡配置信息
type GATableConfig struct {
	GameModule   int   //配置-专场 1.积分模式 2.金币模式
	BaseScore    int   //配置-底分
	TotalRound   int   //配置-局数
	LimitMoney   int64 //配置-携带金币
	PlayerNumber int   //配置-牌局人数
	GameType     int   //配置-玩法模式 1.普通模式 2.加一色模式
	PayType      int   //配置-支付类型 1.房主支付 2.AA支付
	TimeType     int   //配置-玩牌时间 1.15s 2.30s 3.60s 4.90s 5.120s
	FailType     int   //配置-相公模式 1.是 2.否
}

//房卡信息返回
type GTableInfoReply struct {
	Id      int32
	TableId string        //房间号
	Config  GATableConfig //房间配置
}

//玩家准备
type GAPlayerReady struct {
	Id      int
	IsReady int //0.为准备 1.准备
}

//玩家准备
type GSPlayerReady struct {
	Id      int
	ChairId int32 //座位ID
	IsReady int   //0.为准备 1.准备
}

//开始消息
type GSStartInfo struct {
	Id    int
	Round int
}

//阶段时间
type GSStageInfo struct {
	Id        int
	Stage     int
	StageTime int
}

//手牌消息
type GSHandInfo struct {
	Id           int
	ChairId      int32
	HandCards    []int
	SpecialType  int
	SpecialCards []int
}

//玩家玩牌
type GAPlayInfo struct {
	Id        int
	PlayType  int //0 自己摆牌 ；摆特殊牌型
	PlayCards []int
}

//错误信息
type GSInfoErr struct {
	Id  int
	Err string
}

//摆牌结束
type GSPlayInfo struct {
	Id      int
	ChairId int32
}

//结算玩家消息
type GSettlePlayerInfo struct {
	Uid          int64
	ChairId      int32
	PlayCards    []int   //结算摆牌
	SpecialType  int     //特殊牌型
	SpecialScore int     //特殊得分
	NormalTypes  []int   //牌型数组[头墩牌型，中墩牌型，底分牌型]
	WinCoinList  []int   //比分总得分数组 [头墩得分，中墩得分，底分得分，总得分/特殊得分]
	WinCoins     int64   //总输赢
	Coins        int64   //身上金币
	NormalScores []int   //普通得分
	ShootList    []int32 //打枪{位置1，位置2}
	ShootScoress [][]int //打枪分数{[头墩得分，中墩得分，底分得分]，[头墩得分，中墩得分，底分得分]}
}

//结算消息
type GSSettleInfos struct {
	Id            int
	AllWinChairId int32 //全垒打
	PlayInfo      []GSettlePlayerInfo
}

//总结算玩家信息
type GSSettlePlayInfoEnd struct {
	ChairId  int32 //位置信息
	WinCoins int64 //输赢得分
}

//总结算
type GSSettleInfoEnd struct {
	Id        int
	PlayInfos []GSSettlePlayInfoEnd //玩家输赢集合
}

//重连消息
type GSReconnectInfo struct {
	Id           int
	ChairIds     []int32 //所有玩家位置[0,1,2,3,4]
	States       []int   //所有玩家状态[1,1,1,1,2]
	PlayNum      int     //已出牌玩家数量
	PlayChairIds []int32 //已出牌玩家位置[0,2,3]
	HandCards    []int   //手牌
	SpecialType  int     //特殊牌型
	Stage        int     //状态
	StageTime    int     //状态时间
	Round        int
	DisPlayer    []int32 //解散数组
}

// 游戏解散请求
type GADismiss struct {
	Id        int
	IsDismiss int //0.不同意 1.同意
}

// 游戏解散信息
type GSDismiss struct {
	Id        int
	DisPlayer []int32 //解散数组
	IsDismiss int     //0.不同意 1.部分同意 2.不能点击解散按钮 3.所有人同意
	Message   string  //返回消息
}

//游戏记录
type GSRecordInfo struct {
	WinCoins int64
	WinDate  string
}

//玩家游戏记录返回
type GSRecordInfos struct {
	Id    int
	Infos []GSRecordInfo
}

//----------------游戏内参数使用--------------------
type GCardsType struct {
	Type  int
	Cards []int
}

type GRecommendPoker struct {
	Types []int
	Cards []int
}

//----------------游戏通讯记录---------------------
//游戏记录
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	GameType    int               `json:"gameType"`
	PayType     int               `json:"payType"`
	Round       int               `json:"round"`
	GameModule  int               `json:"gameModule"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
	RoomNo      string            `json:"roomNo"`
}

//游戏详细记录
type GGameRecordInfo struct {
	UserId      int64  `json:"userId"`
	UserAccount string `json:"userAccount"`
	Robot       bool   `json:"robot"`
	CoinsBefore int64  `json:"coinsBefore"`
	BetCoins    int64  `json:"betCoins"`
	Coins       int64  `json:"coins"`
	CoinsAfter  int64  `json:"coinsAfter"`
	HeadCards   []int  `json:"headCards"`
	MiddleCards []int  `json:"middleCards"`
	BottomCards []int  `json:"bottomCards"`
	Multiple    int    `json:"multiple"`
	Score       int    `json:"score"`
}

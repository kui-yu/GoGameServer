package main

//命名规则 G + S(Send)/A(Accept) + 自定义
//所有匹配玩家信息
type GSInfoAutoGame struct {
	Id   int32
	Seat []GSeatInfo
}

//桌子信息
type GSTableInfo struct {
	Id      int32
	TableId string //房间号
	BScore  int    //底分
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

type GCardsType struct {
	Type  int
	Cards []int
}

type GRecommendPoker struct {
	Types []int
	Cards []int
}

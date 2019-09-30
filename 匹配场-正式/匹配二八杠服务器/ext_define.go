package main

//命名规则 G + S(Send)/A(Accept) + 自定义
//所有匹配玩家信息
type GSInfoAutoGame struct {
	Id   int32
	Seat []GSeatInfo
}

//桌子信息
type GSTableInfo struct {
	Id      int
	TableId string //房间号
	BScore  int    //底分
}

//阶段时间
type GSStageInfo struct {
	Id        int
	Stage     int
	StageTime int
}

//开始消息
type GSStartInfo struct {
	Id    int
	Round int //开始回合
}

type GSCallList struct {
	Id          int
	CallListCnt int
	CallList    []int
}

//玩家叫庄
type GAPlayerCallInfo struct {
	Id           int
	CallMultiple int //叫庄倍数
}

//玩家叫庄返回
type GSPlayerCallInfo struct {
	Id           int
	ChairId      int32 //位置ID
	CallMultiple int   //叫庄倍数
}

//玩家抢庄结果
type GSPlayerCallBank struct {
	Id              int
	Banker          int32
	BankerList      []int
	BankerMultiples int   //庄家倍数
	BetListCnt      int   //下注数量
	BetList         []int //下注列表
}

//玩家下注
type GAPlayerPlayInfo struct {
	Id           int
	PlayMultiple int //下注倍数
}

//玩家下注返回
type GSPlayerPlayInfo struct {
	Id           int
	ChairId      int32 //位置ID
	PlayMultiple int   //下注倍数
}

//发牌
type GSCardInfo struct {
	Id        int
	Dices     []int //双骰
	HandCards []int //手牌
}

//结算玩家信息
type GSSettlePlayInfo struct {
	ChairId         int32   //位置信息
	HandCard        []int   //玩家手牌
	BankerMultiples int     //庄家倍数
	PlayerMultiples int     //自己倍数
	WinCoins        int64   //输赢得分
	Coins           int64   //身上金币
	WinList         []int32 //输赢玩家列表
}

//结算
type GSSettleInfo struct {
	Id        int
	Round     int
	PutInfos  []G_PutInfo        //历史纪录
	PlayInfos []GSSettlePlayInfo //玩家输赢集合
}

//座位信息-重连
type GSeatInfoReconnect struct {
	Id             int32
	ChairIds       []int32     //所有玩家位置[0,1,2,3,4]
	States         []int32     //所有玩家状态[1,1,1,1,2]
	CallMultiples  []int       //所有玩家倍数[0,0,1,2,3]
	PlayMultiples  []int       //所有玩家倍数[0,0,1,2,3]
	MyCard         []int       //手牌
	BankerId       int32       //庄家
	BankerMultiple int         //庄家倍数
	Round          int         //当前回合
	PutInfos       []G_PutInfo //历史纪录
	RsInfo         GSSettleInfo
	Stage          int
	StageTime      int
	CallListCnt    int
	CallList       []int
	BetListCnt     int   //下注数量
	BetList        []int //下注列表
}

type GSRecordInfo struct {
	WinCoins int64
	WinDate  string
}

type GSRecordInfos struct {
	Id    int
	Infos []GSRecordInfo
}

//结算玩家信息
type GSSettlePlayInfoEnd struct {
	Uid      int64
	ChairId  int32 //位置信息
	WinCoins int64 //输赢得分
	Coins    int64
}

//结算
type GSSettleInfoEnd struct {
	Id        int
	PlayInfos []GSSettlePlayInfoEnd //玩家输赢集合
}

//游戏记录
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	Round       int               `json:"round"`
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

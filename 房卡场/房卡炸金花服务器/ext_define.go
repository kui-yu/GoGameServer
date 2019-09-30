package main

//命名规则 G + S(Send)/A(Accept) + 自定义
//所有匹配玩家信息
type GSInfoAutoGame struct {
	Id   int32
	Seat []GSSeatInfo
}

//桌子信息
type GSTableInfo struct {
	Id      int
	TableId string //房间号
	BScore  int64  //底分
	Round   int    //总轮数
}

//阶段时间
type GSStageInfo struct {
	Id        int
	Stage     int
	StageTime int
}

//游戏开始
type GSGameStart struct {
	Id    int
	Round int
}
type GSCoinMsg struct {
	Id      int
	AllCoin int64   //场内总金币
	PCoin   []int64 //玩家金币
	Round   int     //轮数
}

//底注通知
type GSPlayerBscore struct {
	Id     int
	Bscore int64 //底注
}

//叫牌玩家
type GSPlayerCallPlayer struct {
	Id         int
	Player     int32 //叫牌玩家ChairId
	Round      int   //当前轮数
	CoinEnough int   //0不足 1充足
	MinCoin    int64 //当前最低注
}

//玩家操作
type GAPlayerOperation struct {
	Id        int
	ChairId   int32 //座位号
	PlayCoin  int64 //下注金币
	Operation int   //操作(0弃牌，1看牌，2比牌，3加注，4跟注)
}

//看牌
type GSCardInfo struct {
	Id         int
	HandCards  []int //手牌
	Lv         int   //牌等级
	ChairId    int32 //座位号
	Model      int   // 0 主动看牌，1 失败看牌
	CoinEnough int   //0不足 1充足
}

//玩家牌型(是否已看牌，弃牌)
type GSCardType struct {
	Id      int
	ChairId int32 //座位号
}

//玩家后台操作
type GAProtectGiveUp struct {
	Id         int
	PAttribute int //玩家托付后台属性 1自动跟注 2防弃牌
}

//玩家 后台操作信息
type GSSystemOpertion struct {
	Id         int
	PAttribute int //玩家托付后台属性 1自动跟注 2防弃牌
	OpSuccess  int //0未启动，1启动
}

//返回下注操作
type GSPlayerPayCoin struct {
	Id        int
	PChairId  int32   //当前操作玩家座位号
	PlayCoin  int64   //下注金币
	ChairId   []int32 //比牌玩家座位号
	Winner    int32   //获胜者chairid
	MinCoin   int64   //当前最低注
	Operation int     //玩家操作返回(2比牌，3加注，4跟注)
}

//金币不足比牌
type GSPlayerContest struct {
	Id       int             //
	Count    int             //比牌次数
	PContest []PlayerContest //比牌玩家
	PlayCoin int64           //玩家下注
}

type PlayerContest struct {
	Person_1  int32
	Person_2  int32
	Winner    int32 //获胜玩家chairid
	LoserCard []int //输家手牌
	CardLv    int   //牌等级
}

//结算玩家信息
type GSSettlePlayInfo struct {
	Id       int
	Count    int          //比牌数
	PContest []Contest    //比牌玩家
	SCard    []SettleCard //结算后玩家看牌
	CList    []CoinList   //所有玩家结算金币
}

type Contest struct {
	Person_1 int32
	Person_2 int32
	Winner   int32
}

type SettleCard struct {
	ChairId  int32 //座位
	Identity int   //0 winner  1 loser
	HandCard []int //手牌
	Lv       int   //牌等级
}

type CoinList struct {
	ChairId  int32 //座位
	WinCoins int64 //输赢金币
	Coins    int64 //身上金币
}

//玩家离线通知
type GSLeave struct {
	Id      int
	ChairId int32
}

//座位信息-重连
type GSeatInfoReconnect struct {
	Id              int32
	ChairIds        []int32 //所有玩家位置[0,1,2,3,4,5]
	States          []int32 //所有玩家状态 1在线 2离线
	CardType        []int   //所有玩家牌状态 0未看牌，1已看牌，2已弃牌，3比牌输家
	CallPlayer      int32   //当前叫牌玩家
	Round           int     //当前回合
	GameRound       int     //局数
	CoinList        int64   //场上总金币数
	PayCoin         []int64 //玩家对应下注金币数
	MinCoin         int64   //最小下注金币
	Stage           int
	StageTime       int
	TimeRemaining   int             //剩余时间
	ReconnectPlayer GSPlayerConnect //重连玩家信息
	DisPlayer       []int32         //解散玩家
}

type GSPlayerConnect struct {
	ChairId      int32 //重连玩家座位
	AutoFollowUp int   //是否自动跟注 0否，1是
	ProtectGU    int   //防超时弃牌 0否，1是
	CardType     int   //牌状态 0未看牌，1已看牌，2已弃牌，3比牌输家
	CardLv       int   //牌等级
	HandCard     []int //cardtype为1时手牌
	CoinEnough   int   //0不足 1充足
}

//离开
type GAPlayerLeave struct {
	Id        int32
	ChairId   int32 // 座位号
	LeaveType int   // 离开类型 0否 1是 2继续游戏
}

//离开应答
type GSPlayerLeave struct {
	Id        int32
	ChairId   int32 // 座位号
	LeaveType int   // 是否已从大厅离开 0否 1是 2继续游戏
}

//通知牌
type GSMaxCard struct {
	Id             int
	IsRobot        int         //0玩家 1机器人
	CardLv         int         //最大牌型
	HandCard       []int       //手牌
	ChairId        int32       //座位号
	WinnerRole     int         //0玩家，1机器人
	PlayerHandCard []PHandCard //真实玩家手牌
}

type PHandCard struct {
	HandCards []int
	CardLv    int
	ChairId   int32 //座位号
}

//接收最大牌
type GAMaxCard struct {
	Id       int
	CardLv   int //
	HandCard []int
}

//换牌消息返回
type GSChangeCard struct {
	Id       int
	HandCard []int //手牌
	CardLv   int   //牌等级
	Result   int   //0成功，其他失败
}

////////////////////////////////////////////

//游戏记录
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameModule  int               `json:"gameModule"`
	PayType     int               `json:"payType"`
	GameRoundNo string            `json:"gameRoundNo"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
	GameType    int               `json:"gameType"`
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
	BrandMultiple int    `json:"brandMultiple"`
	BetMultiple   int    `json:"betMultiple"`
	Multiple      int    `json:"multiple"`
	Score         int64  `json:"score"`
}

//房卡配置信息
// type GATableConfig struct {
// 	GameModule   int   //配置-专场 1.积分模式 2.金币模式
// 	BaseScore    int   //配置-底分
// 	MaxCall      int   //最大跟注轮数
// 	TotalRound   int   //配置-局数
// 	LimitMoney   int64 //配置-携带金币
// 	PlayerNumber int   //配置-牌局人数
// 	GameType     int   //配置-玩法模式 1.普通模式 2.加一色模式
// 	PayType      int   //配置-支付类型 1.房主支付 2.AA支付
// 	TimeType     int   //配置-玩牌时间 1.15s 2.30s 3.60s 4.90s 5.120s
// 	FailType     int   //配置-相公模式 1.是 2.否
// }
type GATableConfig struct {
	GameModule   int   //配置-专场 1.积分模式 2.金币模式
	BaseScore    int   //配置-底分
	TotalRound   int   //配置-局数
	PlayerNumber int   //配置-牌局人数
	BetRound     int   //最大跟注轮数
	LimitMoney   int64 //配置-携带金币
	PayType      int   //配置-支付类型 1.房主支付 2.AA支付
}

//房卡信息返回
type GTableInfoReply struct {
	Id      int32
	TableId string        //房间号
	Config  GATableConfig //房间配置
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

//错误信息
type GSInfoErr struct {
	Id  int
	Err string
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

//总结算信息
type GSLumpSum struct {
	Id   int
	Info []GSPlayerSum
}
type GSPlayerSum struct {
	Coin    int64
	ChairId int32
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

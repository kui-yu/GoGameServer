package main

//定义

//命名规则 G + S(Send)/A(Accept) + 自定义
//所有匹配玩家信息
type GSInfoAutoGame struct {
	Id   int32
	Seat []GSSeatInfo
}

//座位信息
type GSSeatInfo struct {
	Uid        int64  //用户Id
	Nick       string //用户个性签名
	Ready      bool   //是否准备
	Cid        int32  //椅子号
	Sex        int32  //性别
	Head       string
	Lv         int32 //等级
	Coin       int64 //金币
	TotalCoins int64 //总分数
	IsReady    int   //是否准备
}

//房卡配置信息
type GATableConfig struct {
	GameModule int //配置-专场 1.积分模式 2.金币模式
	BaseScore  int //配置-底分
	TotalRound int //配置-局数
	// LimitMoney   int64 //配置-携带金币
	PlayerNumber int //配置-牌局人数
	GameType     int //配置-玩法模式 1.抢庄 2.通比
	// PayType      int   //配置-支付类型 1.房主支付 2.AA支付
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

//玩家下注
type GACallMsg struct {
	Id       int32
	Multiple int
}

//玩家下注返回
type GSCallMsg struct {
	Id       int32
	ChairId  int32
	Multiple int
}

//阶段时间
type GSStageInfo struct {
	Id        int32
	Stage     int32
	StageTime int32
}

//玩家摆牌结果
type GSPlayCard struct {
	Id      int32
	ChairId int32
}

//玩家抢庄结果
type GSCallBank struct {
	Id                int32
	Banker            int32
	BankerList        []int
	BankerMultiples   int   //庄家倍数
	CallMultiplesList []int //下注倍数集合
}

//返回发牌结果
type GSSendHandCards struct {
	Id       int32
	ChairId  int32
	NiuPoint int
	NiuCards []int
}

//结算子集
type GSWinInfo struct {
	Uid      int64
	ChairId  int32
	WinCoin  int64
	Coins    int64
	NiuPoint int
	NiuCards []int
}

//结算结果
type GSWinInfosReply struct {
	Id        int32
	InfoCount int
	Infos     []GSWinInfo
}

//游戏记录
type GSRecordInfo struct {
	WinCoins int64
	WinDate  string
}

//座位信息-重连
type GSSeatInfoReconnect struct {
	Id              int32
	ChairIds        []int32 //所有玩家位置[0,1,2,3,4]
	States          []int32 //所有玩家状态[1,1,1,1,2]
	Multiples       []int   //所有玩家倍数[0,0,1,2,3]
	Banker          int32   //庄家
	BankerMultiples int     //庄家倍数
	CallMultiples   []int   //玩家叫庄倍数 -1,没叫；0,不抢
	PlayNum         int     //已出牌玩家数量
	PlayChairIds    []int32 //已出牌玩家位置[0,2,3]
	DisPlayer       []int32 //解散数组
	MyCard          []int   //手牌
	Round           int
	Stage           int
	StageTime       int
	CallListCnt     int
	CallList        []int
	BetListCnt      int   //下注数量
	BetList         []int //下注列表
}

//游戏记录  发送给数据库
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	GameType    int               `json:"gameType"`
	Round       int               `json:"round"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
	RoomNo      string            `json:"roomNo"`
}

//游戏详细记录 发送给大厅
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
	Banker        int32  `json:"banker"`
	Score         int    `json:"score"`
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

//玩家游戏记录返回
type GSRecordInfos struct {
	Id    int
	Infos []GSRecordInfo
}

//房卡信息返回
type GTableInfoReply struct {
	Id      int32
	TableId string        //房间号
	Config  GATableConfig //房间配置
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

//抢庄筹码
type GSCallList struct {
	Id          int
	CallListCnt int
	CallList    []int
}

//下注筹码列表
type GCallListMsg struct {
	Id         int
	BetListCnt int   //下注数量
	BetList    []int //下注列表
}

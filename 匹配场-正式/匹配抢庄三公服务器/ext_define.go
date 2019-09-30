package main

const (
	MSG_GAME_INFO_STAGE_INFO            = 410000 + iota //410000//阶段消息
	MSG_GAME_INFO_SEAT                                  //410001//群发座位信息
	MSG_GAME_INFO_ROOM                                  //410002//群发房间信息
	MSG_GAME_INFO_BANKER_MULTIPLE                       //410003//抢庄
	MSG_GAME_INFO_BANKER_MULTIPLE_REPLY                 //410004//抢庄返回
	MSG_GAME_INFO_IDLE_MULTIPLE                         //410005//下注倍数
	MSG_GAME_INFO_IDLE_MULTIPLE_REPLY                   //410006//下注倍数返回
	MSG_GAME_INFO_CHOICE_BANKER_REPLY                   //410007//选庄结果
	MSG_GAME_INFO_SHOW_CARDS                            //410008//亮牌
	MSG_GAME_INFO_SHOW_CARDS_REPLY                      //410009//亮牌返回
	MSG_GAME_INFO_GET_RECORD                            //410010//获取游戏记录
	MSG_GAME_INFO_GET_RECORD_REPLY                      //410011//获取游戏记录返回
	MSG_GAME_INFO_MAX_MULTIPLE_REPLY                    //410012//最大下注倍数返回
	MSG_GAME_INFO_SETTLE_REPLY                          //410013//结算信息
	MSG_GAME_INFO_RECONNECT_DESK_INFO                   //410014//断线重连，桌子消息
)
const (
	STAGE_START           = iota + 10 //10//开始
	STAGE_SHUFFLE_CARDS               //11//洗牌
	STAGE_SEND_CARDS                  //12//发牌
	STAGE_BANKER_MULTIPLE             //13//选庄
	STAGE_IDLE_MULTIPLE               //14//下注
	STAGE_OPEN_CARDS                  //15//开牌
	STAGE_GAME_SETTLE                 //16//结算
)

//座位信息
type SeatInfo struct {
	Head    string
	Name    string
	Coins   int64
	Uid     int64
	ChairId int32
}

//群发的座位玩家信息
type GPlayerSeatInfos struct {
	Id   int
	Data []SeatInfo
}

//群发房间消息
type GRoomInfo struct {
	Id             int
	RoomNumber     string  //局号
	MaxMultiple    int64   //最大倍数
	BScore         int64   //底分
	PlayerMultiple []int64 //下注倍数
}

//阶段消息
type GSGameStageInfo struct {
	Id        int
	Stage     int
	StageTime int
}

//抢庄及下注倍数信息
type BankerInfo struct {
	IsBanker bool  //是不是庄家
	IsChoice int   //是否抢庄,0=未操作，1=不抢，2=抢
	Multiple int64 //抢庄倍数或下注倍数
}

//抢庄或下注倍数请求
type GChoiceMultipleREQ struct {
	Id       int
	Uid      int64
	Multiple int64 //抢庄时1为抢，0为不抢，下注时对应倍数
}

//抢庄或下注倍数响应
type GChoiceMultipleRES struct {
	Id       int
	Uid      int64
	Multiple int64  //倍数，0为不抢
	Err      string //错误信息
	Result   int    //错误码
}

//返回最大可下注倍数
type GStageMultipleRES struct {
	Id       int
	Multiple int //可下注倍数索引
}

//选出的庄家
type GBankerInfo struct {
	Id      int
	Players []int32 //参与抢庄的玩家
	ChairId int32
}

//请求亮牌
type GShowCardsREQ struct {
	Id  int
	Uid int64
}

//响应亮牌
type GShowCardsRES struct {
	Id     int
	Uid    int64
	Cards  Card
	Err    string
	Result int
}

//发送结算信息
type GSettleInfo struct {
	Id           int
	BankerResult PlayerResult   //庄家输赢结果
	IdleResult   []PlayerResult //闲家输赢结果
}

//结算信息
type PlayerResult struct {
	Uid      int64
	WinCoins int64
	Coins    int64 //结算后金币
}

//断线重连，桌子消息
type GRDeskInfo struct {
	Id              int
	Stage           int
	StageTime       int
	BankerChairId   int32
	ChoiceBankerArr []int        //2=抢，1不抢，索引的ChairId
	PlayerCards     []PlayerCard //玩家卡
	IdleBets        []IdleBet    //闲家下注
	Settle          GSettleInfo
}
type PlayerCard struct {
	Uid   int64
	Cards Card
}
type IdleBet struct {
	Uid   int64
	Coins int64
}

//发送消息给大厅去记录游戏记录
type GGameRecord struct {
	Id             int32             //协议号
	GameId         int               `json:"gameId"`
	GradeId        int               `json:"gradeId"`
	RoomId         int               `json:"roomId"`
	GradeNumber    int               `json:"gradeNumber"`
	GameRoundNo    string            `json:"gameRoundNo"`
	PlayerCard     []int             `json:"playerCard"`
	SettlementCard [][]int           `json:"settlementCard"`
	UserRecord     []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId         int64  `json:"userId"`         //玩家id
	UserAccount    string `json:"userAccount"`    //玩家名字
	Robot          bool   `json:"robot"`          //是否机器人
	CoinsBefore    int64  `json:"coinsBefore"`    //下注前金币
	BetCoins       int64  `json:"betCoins"`       //下注金币
	Coins          int64  `json:"coins"`          //金币
	CoinsAfter     int64  `json:"coinsAfter"`     //下注后金币
	BankerMultiple int64  `json:"bankerMultiple"` //庄家倍数
	BetMultiple    int64  `json:"betMultiple"`    //下注倍数
	Banker         bool   `json:"banker"`         //是否是庄家
	BaseScore      int64  `json:"baseScore"`      //底注
	PrizeCoins     int64  `json:"PrizeCoins"`     //输赢金币
}

//上线和掉线通知
type Ext_GOnLineNotify struct {
	Id    int32
	Uid   int32
	State int32 //1上线，2掉线
}

package main

//大厅消息
const (
	MSG_HALL_START = 300000 + iota
	MSG_HALL_LOGIN
	MSG_HALL_LOGIN_REPLY
	MSG_HALL_LEAVE
	MSG_HALL_LEAVE_REPLY
)

//每个游戏自己的状态请添加到ext_define.go中，以10开始
const (
	GAME_STATUS_FREE  = 0 + iota // 桌子空闲状态
	GAME_STATUS_START            // 游戏开始
	GAME_STATUS_END              //游戏结束
)

//自己定义的状态，从10开始
const (
	GAME_STATE_SENDCARD        = 10 + iota //发牌阶段
	GAME_STATE_CHANGECARD                  //换牌阶段
	GAME_STATE_CHANGECARD_OVER             //换牌结束阶段
	GAME_STATE_DINGQUE                     //定缺阶段
	GAME_STATE_PLAY                        //玩游戏阶段
)

// 系统消息，desk 相关
const (
	MSG_GAME_START = 400000 + iota
	MSG_GAME_AUTO
	MSG_GAME_AUTO_REPLY
	MSG_GAME_JOIN
	MSG_GAME_JOIN_REPLY
	MSG_GAME_CREATE
	MSG_GAME_CREATE_REPLY
	MSG_GAME_LEAVE
	MSG_GAME_LEAVE_REPLY
	MSG_GAME_END_NOTIFY
	MSG_GAME_RECONNECT //10
	MSG_GAME_RECONNECT_REPLY
)

//消息自定义
const (
	MSG_GAME_INFO_STAGE              = 410030 //阶段消息
	MSG_GAME_INFO_AUTO_REPLY         = 410018 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_ROOM_NOTIFY        = 410015 //房间信息通知410003
	MSG_GAME_INFO_SEND_NOTIFY        = 410001 //发牌410001
	MSG_GAME_INFO_HUANPAI            = 410002 //换牌410002
	MSG_GAME_INFO_HUANPAI_NOTIFY     = 410003 //换牌通知410003
	MSG_GAME_INFO_HUANPAIOVER_NOTIFY = 410004
	MSG_GAME_INFO_DINGQUE            = 410008 //410008
	MSG_GAME_INFO_DINGQUE_NOTIFY     = 410009
	MSG_GAME_INFO_HAVEACTION_NOTIFY  = 410012 //410012
	MSG_GAME_INFO_OUTCARD            = 410005 //出牌410005
	MSG_GAME_INFO_OUTCARD_NOTIFY     = 410006 //410006
	MSG_GAME_INFO_ACTION             = 410010 //410010
	MSG_GAME_INFO_ACTION_NOTIFY      = 410011 //410011
	MSG_GAME_INFO_SENDCARD_NOTIFY    = 410007 //410007
)

// 登录大厅需要的结构数据
type HMsgHallLogin struct {
	Id      int32
	Account string
	Gid     string
}

// 登录大厅返回的结构数据
type HMsgHallLoginReply struct {
	Id      int32
	Account string
	Uid     int64
	Nick    string
	Sex     int32
	Head    string
	Lv      int32
	Coin    int32
	GameId  int32
}

type GSeatInfo struct {
	Uid   int64
	Nick  string
	Ready bool
	Cid   int32 //椅子号
	Sex   int32
	Head  string
	Lv    int32
	Coin  int32
}

//自由匹配
type GAutoGame struct {
	Id      int32
	Account string
	Uid     int64
	Nick    string
	Sex     int32
	Head    string
	Lv      int32
	Coin    int32
}

type GAutoGame2 struct {
	Id int32
}

type GReconnect struct {
	Id int32
}

type GAutoGameReply struct {
	Id     int32
	Result int32
	Seat   []GSeatInfo
}

type GInfoReConnectReply struct {
	Id        int32
	GameState int32       //游戏状态
	Cid       int32       //座位号id
	Seats     []GSeatInfo //所有玩家信息
}

type GInfoAutoGameReply struct {
	Id   int32
	Seat []GSeatInfo
}

//阶段时间
type GStageInfo struct {
	Id        int32
	Stage     int32
	StageTime int32
}

//结算消息
type GSSettleInfos struct {
	Id            int
	AllWinChairId int32 //全垒打
	PlayInfo      []GSettlePlayerInfo
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

//玩家玩牌
type GAPlayInfo struct {
	Id        int
	PlayType  int //0 自己摆牌 ；摆特殊牌型
	PlayCards []int
}

type GGameSendCardNotify struct {
	Id         int32
	Banker     int //庄
	HandsCards []int
}

type GHuanPai struct {
	Id    int32
	Cards []int
}

type GHuanPaiOver struct {
	Id    int32
	Style int //0顺时针，1逆时针
	Cards []int
}

type GDingQue struct {
	Id    int32
	Color int //0玩，1条，2饼
}

type GDingQueNotify struct {
	Id  int32
	Cid int //玩家的座位号id
}

type HaveAction struct {
	Style   int
	Card    int
	HuTypes []int
}
type GHaveActionNotify struct {
	Id   int32
	Data []HaveAction
}

type GOutCard struct {
	Id   int32 //协议号
	Card int   //出的牌
}

type GOutCardNotify struct {
	Id   int32
	Cid  int
	Card int
}

type GAction struct {
	Id    int32 //协议号
	Style int   //动作类型
	Card  int   //出的牌
}

type GActionDoNotify struct {
	Id         int32
	Cid        int
	ActionType int
	Cards      []int
}

//
type GGiveUpNotify struct {
	Id  int32
	Cid int
}

type GSendCardNofify struct {
	Id   int32 //协议号
	Cid  int32 //哪个玩家叫的分
	Card int   //发的牌
	Gang bool  //是否补杠发的牌
}

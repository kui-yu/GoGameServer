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

//阶段
const (
	//匹配
	STAGE_PLAY   = 12 //玩牌阶段
	STAGE_SETTLE = 13 //结算阶段
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
	MSG_GAME_INFO_STAGE             = 410001 //阶段消息
	MSG_GAME_INFO_AUTO_REPLY        = 410002 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_ROOM_NOTIFY       = 410003 //房间信息通知410003
	MSG_GAME_INFO_PLAY              = 410005 //玩家摆牌
	MSG_GAME_INFO_SETTLE_INFO_REPLY = 410007 //结算
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

package main

// import (
// 	. "MaJiangTool"
// )

//自己定义的游服id从410000开始
const (
	MSG_GAME_INFO_START              = 410000 + iota //游戏开始额外信息
	MSG_GAME_INFO_SEND_NOTIFY                        //发牌410001
	MSG_GAME_INFO_HUANPAI                            //换牌410002
	MSG_GAME_INFO_HUANPAI_NOTIFY                     //换牌通知410003
	MSG_GAME_INFO_HUANPAIOVER_NOTIFY                 //换牌结束
	MSG_GAME_INFO_OUTCARD                            //出牌410005
	MSG_GAME_INFO_OUTCARD_NOTIFY                     //410006
	MSG_GAME_INFO_SENDCARD_NOTIFY                    //410007
	MSG_GAME_INFO_DINGQUE                            //410008
	MSG_GAME_INFO_DINGQUE_NOTIFY                     //410009
	MSG_GAME_INFO_ACTION                             //410010
	MSG_GAME_INFO_ACTION_NOTIFY                      //410011
	MSG_GAME_INFO_HAVEACTION_NOTIFY                  //410012
	MSG_GAME_INFO_GIVEUP                             //410013
	MSG_GAME_INFO_GIVEUP_NOTIFY                      //410014
	MSG_GAME_INFO_ROOM_NOTIFY                        ////410015
	MSG_GAME_INFO_TUOGUAN                            //410016
	MSG_GAME_INFO_TUOGUAN_REPLY                      //410017
	MSG_GAME_INFO_AUTO_REPLY                         //410018,游戏随机匹配成功的数据
	MSG_GAME_INFO_RECONNECT_REPLY                    //410019
	MSG_GAME_INFO_END_NOTIFY                         //410020
	MSG_GAME_INFO_LEAVE_NOTIFY                       //410021
	MSG_GAME_INFO_GAMEOVER                           //410022
)

const (
	MSG_GAME_INFO_STAGE = 410030 //阶段消息
)

//自己定义的状态，从10开始
const (
	GAME_STATE_SENDCARD        = 10 + iota //发牌阶段
	GAME_STATE_CHANGECARD                  //换牌阶段
	GAME_STATE_CHANGECARD_OVER             //换牌结束阶段
	GAME_STATE_DINGQUE                     //定缺阶段
	GAME_STATE_PLAY                        //玩游戏阶段
)

//////////////////////////////////////////////////////////
const (
	TIMER_START     = 1
	TIMER_START_NUM = 3 //开始动画
	//
	TIMER_SENDCARD     = 2
	TIMER_SENDCARD_NUM = 3 //3秒后自动发牌
	//
	TIMER_HUANPAI     = 3
	TIMER_HUANPAI_NUM = 3
	//
	TIMER_HUANPAIOVER     = 4
	TIMER_HUANPAIOVER_NUM = 3
	//
	TIMER_DINGQUE     = 5
	TIMER_DINGQUE_NUM = 3 //叫分超时时间
	//
	TIMER_OUTCARD     = 6
	TIMER_OUTCARD_NUM = 5 //出牌超时时间

	TIMER_ACTION     = 7
	TIMER_ACTION_NUM = 5

	TIMER_OVER     = 8
	TIMER_OVER_NUM = 1 //归还桌子时间
)

//游戏开始
type GGameStartNotify struct {
	Id int32
}

//出牌
type GGameOutCard struct {
	Id    int32   //协议号
	Cards []int32 //牌
}

//出牌应答
type GGameOutCardReply struct {
	Id     int32  //协议号
	Cid    int32  //谁出的牌
	Type   int32  //出牌类型
	Cards  []byte //牌
	Max    byte   //最大牌
	Double int32  //倍数,炸弹翻倍
}

//过
type GGamePass struct {
	Id int32 //协议号
}

//过应答
type GGamePassReply struct {
	Id  int32 //协议号
	Cid int32
}

//发牌通知
type GGameSendCardNotify struct {
	Id         int32
	Banker     int //庄
	HandsCards []int
}

//
type GHuanPai struct {
	Id    int32
	Cards []int
}

type GHuanPaiNotify struct {
	Id  int32
	Cid int //换牌的座位号id
}

type GHuanPaiOver struct {
	Id    int32
	Style int //0顺时针，1逆时针
	Cards []int
}

//
type GDingQue struct {
	Id    int32
	Color int //0玩，1条，2饼
}

type GDingQueNotify struct {
	Id  int32
	Cid int //玩家的座位号id
}

type GSendCardNofify struct {
	Id   int32 //协议号
	Cid  int32 //哪个玩家叫的分
	Card int   //发的牌
	Gang bool  //是否补杠发的牌
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

type HaveAction struct {
	Style   int
	Card    int
	HuTypes []int
}
type GHaveActionNotify struct {
	Id   int32
	Data []HaveAction
}

//
type GGiveUpNotify struct {
	Id  int32
	Cid int
}

//
type GActionDoNotify struct {
	Id         int32
	Cid        int
	ActionType int
	Cards      []int
}

//定庄通知
type GBankerNotify struct {
	Id     int32
	Banker int32
	DiPai  []byte
	Double int32
}

//底牌
type GInfoGameEndPlayInfo struct {
	Cid       int32
	HandCards []byte
}

//游戏结束
type GInfoGameEnd struct {
	Id        int32    //协议号
	EndType   int32    //结束类型,0地主赢了，1,农民赢了
	ChunTian  int32    //是否春天
	Double    []int32  //加倍
	Scores    []int64  //得分情况
	Coins     []int64  //结算后的分数
	Accouts   []string //账号
	PlayInfos []GInfoGameEndPlayInfo
}

//游戏房间信息
type GGameInfoNotify struct {
	Id     int32
	BScore int
	MaxBei int
	JuHao  string
}

type GTuoGuan struct {
	Id  int32 //协议号
	Ctl int32 //1托管，2取消托管
}

type GTuoGuanReply struct {
	Id     int32 //协议号
	Ctl    int32
	Result int32 //结果
	Err    string
	Cid    int32 //谁托管
}

//断线重连失败，需要清除gameid
type GReConnectFailedNotify struct {
	Id int32
}

type GInfoReConnectReply struct {
	Id            int32
	GameState     int32       //游戏状态
	GameStateTime int32       //状态时间
	Cid           int32       //座位号id
	Seats         []GSeatInfo //所有玩家信息
	BScore        int32
	MaxBei        int32
	JuHao         string
	TimerNum      int32 //定时器时间
	//叫分阶段
	Cards    []byte
	CardNum  []int32
	CurCid   int32
	CallFens []int32
	TuoGuans []int32 //是否托管
	LiXians  []int32 //是否在线
	//游戏开始阶段
	Banker   int32
	LastCall int32
	DiPai    []byte
	Double   int32
	OutEd    []GOutCard
}

type GInfoAutoGameReply struct {
	Id   int32
	Seat []GSeatInfo
}

//
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId      int64  `json:"userId"`
	UserAccount string `json:"userAccount"`
	Robot       bool   `json:"robot"`
	CoinsBefore int64  `json:"coinsBefore"`
	BetCoins    int64  `json:"betCoins"`
	Coins       int64  `json:"coins"`
	CoinsAfter  int64  `json:"coinsAfter"`
	Score       int    `json:"score"`
	Multiple    int    `json:"multiple"`
	Landlord    bool   `json:"landlord"`
}

//阶段时间
type GStageInfo struct {
	Id        int32
	Stage     int32
	StageTime int32
}

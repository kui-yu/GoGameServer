package main

//自己定义的游服id从410000开始
const (
	MSG_GAME_INFO_START           = 410000 + iota //游戏开始额外信息
	MSG_GAME_INFO_SEND_NOTIFY                     //发牌410001
	MSG_GAME_INFO_CALL                            //叫分410002
	MSG_GAME_INFO_CALL_REPLY                      //410003
	MSG_GAME_INFO_BANKER_NOTIFY                   //定庄410004
	MSG_GAME_INFO_OUTCARD                         //出牌410005
	MSG_GAME_INFO_OUTCARD_REPLY                   //410006
	MSG_GAME_INFO_PASS                            //过410007
	MSG_GAME_INFO_PASS_REPLY                      //410008
	MSG_GAME_INFO_ROOM_NOTIFY                     //房间信息通知410009
	MSG_GAME_INFO_TUOGUAN                         //410010托管
	MSG_GAME_INFO_TUOGUAN_REPLY                   //410011
	MSG_GAME_INFO_AUTO_REPLY                      //410012,游戏随机匹配成功的数据
	MSG_GAME_INFO_RECONNECT_REPLY                 //断线重连应答
	MSG_GAME_INFO_END_NOTIFY                      //410014游戏结束应答
	MSG_GAME_INFO_LEAVE_NOTIFY                    //玩家离开通知410015
)

//lgh --2019/03/04
const (
	MSG_GAME_INFO_STAGE      = 410030 //阶段消息
	MSG_GAME_INFO_TIPS       = 410031 //出牌提示
	MSG_GAME_INFO_TIPS_REPLY = 410032 //出牌消息回复
)

//自己定义的阶段，从10开始
const (
	GAME_STATUS_CALL = 10 + iota // 叫分阶段
	GAME_STATUS_PLAY             // 游戏操作阶段11

)

//////////////////////////////////////////////////////////
const (
	TIMER_SENDCARD     = 1
	TIMER_SENDCARD_NUM = 1 //3秒后自动发牌
	//
	TIMER_CALL     = 2
	TIMER_CALL_NUM = 10 //叫分超时时间
	//
	TIMER_OUTCARD     = 3
	TIMER_OUTCARD_NUM = 20 //出牌超时时间

	TIMER_OVER     = 8
	TIMER_OVER_NUM = 1 //归还桌子时间
)

const (
	TIMER_START     = 11
	TIMER_START_NUM = 3 //开始动画
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

//出牌提示回复
type GGameTips struct {
	Id   int32    //协议号
	Tips [][]byte //出牌数组
}

//牌型记录，桌子记录出牌
type GOutCard struct {
	Cid   int32
	Type  int32
	Max   byte
	Cards []byte
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
	Cid        int32 //此人开始叫地主
	HandsCards []byte
}

//叫地主
type GCallMsg struct {
	Id    int32
	Coins int32
}

type GCallMsgReply struct {
	Id    int32 //协议号
	Cid   int32 //哪位玩家叫得分
	Coins int32 //叫几分
	End   int32 //叫分是否结束
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
	BScore int32
	MaxBei int32
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

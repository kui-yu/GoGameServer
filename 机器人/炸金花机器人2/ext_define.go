package main

// 定义游戏消息，desk 相关
const (
	MSG_GAME_INFO_AUTO_REPLY = 410002 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_END_NOTIFY = 410012 //结算
)

const (
	GAME_STATUS_FREE  = 0 + iota // 结束
	GAME_STATUS_START            // 游戏开始
	GAME_STATUS_END
)

const (
	// GAME_STATUS_END      = 0 + iota // 结束
	// GAME_STATUS_START               // 游戏开始
	// GAME_STATUS_SENDCARD = 10 + iota // 发牌阶段
	GAME_STATUS_CALL   = 10 + iota // 叫分阶段
	GAME_STATUS_BANKER             // 定地主
	GAME_STATUS_PLAY               // 游戏操作阶段

)

//////////////////////////////////////////////////////////
const (
	CT_ERROR               = 0  //错误类型
	CT_SINGLE              = 1  //单牌类型
	CT_DOUBLE              = 2  //对子类型
	CT_SINGLE_CONNECT      = 3  //单龙
	CT_DOUBLE_CONNECT      = 4  //双龙
	CT_THREE               = 5  //三张
	CT_THREE_LINE_TAKE_ONE = 6  //三带一单
	CT_THREE_LINE_TAKE_TWO = 7  //三带一对
	CT_FOUR_LINE_TAKE_ONE  = 8  //四带两单
	CT_FOUR_LINE_TAKE_TWO  = 9  //四带两对
	CT_AIRCRAFT            = 10 //飞机
	CT_AIRCRAFT_ONE        = 11 //飞机带单
	CT_AIRCRAFT_TWO        = 12 //飞机带对
	CT_BOMB_FOUR           = 13 //炸弹
	CT_TWOKING             = 14 //对王类型
)

const (
	CARD_COLOR   = 0xF0 //花色掩码
	CARD_VALUE   = 0x0F //数值掩码
	Card_Invalid = 0x00
	Card_Rear    = 0xFF
)

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
	Id        int32
	GameType  int32
	RoomType  int32
	GradeType int32
}

type GReconnect struct {
	Id int32
}

// type GAutoGameReply struct {
// 	Id     int32
// 	Result int32
// 	Seat   []GSeatInfo
// }

//游戏开始
type GGameStartNotify struct {
	Id     int32
	Bscore int32
}

//出牌
type GGameOutCard struct {
	Id    int32
	Type  int32
	Cards []byte
	Max   byte
}

type GRealOutCard struct {
	Id    int32
	Type  int32
	Cards []int32
}

type GGameOutCardReply struct {
	Id     int32  //协议号
	Cid    int32  //谁出的牌
	Type   int32  //出牌类型
	Cards  []byte //牌
	Max    byte
	Double int32
}

type GOutCard struct {
	Cid   int32
	Type  int32
	Max   byte
	Cards []byte
}

//过
type GGamePass struct {
	Id int32
}

type GGamePassReply struct {
	Id  int32
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
	Cid   int32 //哪个玩家叫的分
	Coins int32 //叫几分
}

//定庄通知
type GBankerNotify struct {
	Id     int32
	Banker int32
	DiPai  []byte
	Double int32
}

//游戏结束
type GGameEnd struct {
	Id       int32   //协议号
	EndType  int32   //结束类型,0地主赢了，1,农民赢了
	ChunTian int32   //是否春天
	Double   []int32 //加倍
	Scores   []int32 //得分情况
	Accouts  []string
}

type GInfoReConnectReply struct {
	Id        int32
	GameState int32       //游戏状态
	Cid       int32       //座位号id
	Seats     []GSeatInfo //所有玩家信息
	//叫分阶段
	Cards    []byte
	CardNum  []int32
	CurCid   int32
	CallFens []int32
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

type GAutoGameReply struct {
	Id     int32
	ToUid  int64
	Result int32 //0成功，其他失败
	Err    string
}

type GInfoGameEnd struct {
	Id       int32    //协议号
	EndType  int32    //结束类型,0地主赢了，1,农民赢了
	ChunTian int32    //是否春天
	Double   []int32  //加倍
	Scores   []int64  //得分情况
	Coins    []int64  //结算后的分数
	Accouts  []string //账号
}

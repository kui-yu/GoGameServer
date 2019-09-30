package main

//自己定义的游服id从410000开始
const (
	MSG_GAME_INFO_START             = 410000 + iota //游戏开始额外信息
	MSG_GAME_INFO_READY                             //玩家准备410001
	MSG_GAME_INFO_READY_REPLY                       //玩家准备消息返回410002
	MSG_GAME_INFO_SEND_NOTIFY                       //发牌410003
	MSG_GAME_INFO_CALL                              //叫分410004
	MSG_GAME_INFO_CALL_REPLY                        //410005
	MSG_GAME_INFO_BANKER_NOTIFY                     //定庄410006
	MSG_GAME_INFO_OUTCARD                           //出牌410007
	MSG_GAME_INFO_OUTCARD_REPLY                     //410008
	MSG_GAME_INFO_PASS                              //过410009
	MSG_GAME_INFO_PASS_REPLY                        //4100010
	MSG_GAME_INFO_ROOM_NOTIFY                       //房间信息通知410011
	MSG_GAME_INFO_TUOGUAN                           //410012托管
	MSG_GAME_INFO_TUOGUAN_REPLY                     //410013
	MSG_GAME_INFO_AUTO_REPLY                        //410014,游戏随机匹配成功的数据
	MSG_GAME_INFO_RECONNECT_REPLY                   //断线重连应答 410015
	MSG_GAME_INFO_END_NOTIFY                        //410016游戏结束应答
	MSG_GAME_INFO_LEAVE_NOTIFY                      //玩家离开通知410017
	MSG_GAME_INFO_ERR                               //错误信息410018
	MSG_GAME_INFO_GETMSG                            //抢地主410019
	MSG_GAME_INFO_GETMSG_REPLY                      //抢地主应答 410020
	MSG_GAME_INFO_BREAKROOM                         //解散房间请求410021
	MSG_GAME_INFO_BREAKROOM_REPLY                   //解散房间应答410022
	MSG_GAME_INFO_OUTCARD_Lz                        //癞子出牌(如果只有一种出牌情况，使用 410008 应答）410023
	MSG_GAME_INFO_OUTCARD_LZ_SELECT                 //返回多种牌型供选择410024
	MSG_GAME_INFO_SELECT                            //选择牌型（牌型）出牌410025
	MSG_GAME_INFO_QPLAYERLOGS                       //请求玩家游戏记录410026
	MSG_GAME_INFO_QPLAYERLOGS_REPLY                 // 游戏记录响应410027
	MSG_GAME_INFO_ALLBALANCE                        //所有玩家总结算410028
)

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

//lgh --2019/03/04
const (
	MSG_GAME_INFO_STAGE = 410030 //阶段消息
)

//自己定义的状态，从10开始
const (
	GAME_STATUS_CALL          = 10 + iota // 叫分阶段
	GAME_STATUS_READ                      //准备阶段 11
	GAME_STATUS_PLAY                      //游戏操作阶段 12
	GAME_STATUS_BreakRoomVote             //房间解散投票阶段 13
	GAME_STATUS_BALANCE                   //游戏结算阶段 14(当局)
)

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

//房卡配置信息
type GATableConfig struct {
	GameModle  int   //配置-专场  1,娱乐专场（积分模式）
	PlayerNum  int   //配置-牌局人数			（二人斗地主只有两个人）
	MatchCount int   //配置-局数
	GameType   int   //配置-玩法模式  1，正常模式  2,，二人斗地主 3，闪电斗地主  4，癞子模式
	PayType    int   //配置-支付条件 1，房主支付
	BaseScore  int   //配置-底分
	CallType   int   //配置-叫分方式  1，叫分  2，抢地主  (两人斗地主只有叫分）
	Boom       int   //配置-炸弹  3，3炸   4，4炸   5，5炸   999，无上限
	CanSelect  []int //配置-可选， 0， 不可三带一对  1，不可4带两队  2，显示剩余手牌数量
}

//房卡信息返回
type GTableInfoReply struct {
	Id      int32
	TableId string        //房间号
	Config  GATableConfig //房间配置
}

//////////////////////////////////////////////////////////
const (
	TIMER_SENDCARD     = 1
	TIMER_SENDCARD_NUM = 1 //3秒后自动发牌
	//
	TIMER_CALL     = 2
	TIMER_CALL_NUM = 10 //叫分超时时间

	TIMER_GETGMS     = 5
	TIMER_GETGMS_NUM = 20 //抢地主超时时间
	//
	TIMER_OUTCARD     = 3
	TIMER_OUTCARD_NUM = 20 //出牌超时时间

	TIMER_READ     = 4
	TIMER_READ_NUM = 999 //准备超时时间

	TIMER_OVER     = 8
	TIMER_OVER_NUM = 1 //归还桌子时间

	TIMER_LZ_START     = 14
	ITMER_LZ_START_NUM = 6 //开始动画 （包括选癞子）

	TIMER_START     = 11
	TIMER_START_NUM = 3 //开始动画

	TIMER_BREAKROOM     = 12
	TIMER_BREAKROOM_NUM = 60 //房间解散投票时间
)

//游戏开始
type GGameStartNotify struct {
	Id    int32
	Round int
}

//出牌
type GGameOutCard struct {
	Id    int32   //协议号
	Cards []int32 //牌
}

//出牌应答
type GGameOutCardReply struct {
	Id      int32 //协议号
	Ishas   bool  //是否含有癞子
	Cid     int32 //谁出的牌
	Type    int32 //出牌类型
	Cards   []int //牌
	CardsLz []int //癞子原来的牌
	Max     byte  //最大牌
	Double  int32 //倍数,炸弹翻倍
	// LzAndBecome []GGameOutLzBecome //癞子和癞子所变换的牌集合，没有出的时候集合为空
	// Ptcon       []byte             //普通牌集合 （有癞子时使用）
	NextCid int32 // 下一位
}

//癞子牌型选择
type GGameSelectOut struct {
	Id     int32
	Canout []CanOutType
}

//癞子和癞子所变换的牌
type GGameOutLzBecome struct {
	Lz     byte
	Become byte //当Become 为0时，则表示代表原牌
}

//牌型记录，桌子记录出牌
type GOutCard struct {
	Cid   int32
	Type  int32
	Max   byte
	Cards []byte
}

type GOutCard1 struct {
	Cid   int32
	Type  int32
	Max   byte
	Cards []int
}

//过应答
type GGamePassReply struct {
	Id      int32 //协议号
	Cid     int32
	NextCid int32
}

//发牌通知
type GGameSendCardNotify struct {
	Id         int32
	Rount      int
	Cid        int32 //此人开始叫地主（抢地主）
	HandsCards []int
	Lz         Lz //本局癞子，如果不是癞子模式则为空
}

//叫地主
type GCallMsg struct {
	Id    int32
	Coins int32
}

//叫地主应答
type GCallMsgReply struct {
	Id      int32 //协议号
	Cid     int32 //哪个玩家叫的分
	Coins   int32 //叫几分
	End     int32 //叫分是否结束
	NextCid int32
}

//抢地主
type GGetMsg struct {
	Id     int32
	GetMsg int32 //是否抢（叫）地主  1，叫地主 2，抢地主 3，不叫 4，不抢
}

//抢地主应答
type GGetMsgReply struct {
	Id        int32
	Cid       int32 //哪一个玩家叫分
	IsGet     int32 //是否抢（叫）地主  1，叫地主 2，抢地主 3，不叫 4，不抢
	CallOrGet int32 //接下来是抢地主 还是 叫地主  1,叫地主 2，抢地主
	End       int32 //抢地主是否结束  1,未结束  2，结束
	NextCid   int32 //下一位玩家的 椅子ID
	Double    int32 //倍数
}

//定庄通知
type GBankerNotify struct {
	Id        int32
	Banker    int32
	DiPai     []int
	Double    int32 //底牌倍数
	AllDouble int32
}

//底牌
type GInfoGameEndPlayInfo struct {
	Cid       int32
	HandCards []int
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
	Ctl    int32 //1，代表托管 2，代表取消托管
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
	Cards     []int
	CardNum   []int32
	CurCid    int32
	CallFens  []int32
	TuoGuans  []int32 //是否托管
	LiXians   []int32 //是否在线
	CallOrGet int     //代表下次是抢地主 还是叫地主
	//游戏开始阶段
	Banker   int32
	LastCall int32
	DiPai    []int
	Double   int32
	OutEd    []GOutCard1
	Round    int //局数
	//玩家解散
	DisPlayer []int32 //同意玩家  椅子号 数组
}

//匹配信息返回
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
	GameModule  int               `json:"gameModule"`
	PayType     int               `json:"payType"`
	GameType    int               `json:"gameType"`
	GameRoundNo string            `json:"gameRoundNo"`
	Round       int               `json:"round"`
	RoomNo      string            `json:"roomNo"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId      int64  `json:"userId"`
	UserAccount string `json:"userAccount"`
	Robot       bool   `json:"robot"`
	CoinsBefore int64  `json:"coinsBefore"`
	// BetCoins    int64  `json:"betCoins"`
	Coins      int64 `json:"coins"`
	CoinsAfter int64 `json:"coinsAfter"`
	Score      int   `json:"score"`
	Multiple   int   `json:"multiple"`
	Landlord   bool  `json:"landlord"`
}

//阶段时间
type GStageInfo struct {
	Id        int32
	Stage     int32
	StageTime int32
}

//解散房间请求
type GBreakRoom struct {
	Id int32
}

//解散房间投票
type GBreakRoomVote struct {
	Id         int32
	AgreetOrNo int32 //1，表示同意 2，表示不同意
	Cid        int32
}

//解散房间投票应答
type GBreakRommVotePlay struct {
	Id         int32
	AgreetOrNo int32
	Name       string
}

//解散房间应答
type GBreakRommReplay struct {
	Id        int
	IsDismiss int     // 0，不同意   1，部分同意 ，2，不能解散  ，  3，  全部同意
	Message   string  //错误信息
	DisPlayer []int32 //同意玩家  椅子号 数组
}

//游戏记录
type GGameLog struct {
	EndTime string //游戏结算时间
	Coins   int64  //输赢金币
}

//游戏记录请求响应
type GGameLogReplay struct {
	Id        int
	GGameLogs []GGameLog //游戏记录
}

type DisMiss struct {
	IsDismiss int // 玩家投票情况
}

type AllBalance struct {
	Id       int     //协议号
	AllConis []int64 //所有玩家的输赢得分
}

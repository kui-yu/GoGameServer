package main

type GInfoAutoGameReply struct {
	Id        int32
	GameState int //游戏状态
	Seat      []GSeatInfo
}

//游戏房间信息
type GGameInfoNotify struct {
	Id     int32
	Bscore int    //底分
	JuHao  string //局号
}

//阶段时间
type GStageInfo struct {
	Id        int32
	Stage     int32 //状态
	StageTime int32 //状态持续时间
}

//发牌通知
type GGameSendCardNotify struct {
	Id        int32
	Cid       int32   //此人开始叫地主（抢地主）
	HandCards []int   //手牌
	Hint      [][]int //能出的牌
}

//玩家出牌结构体
type OutCard struct {
	Id    int     //协议
	Cards []int32 //出牌牌组
}

//玩家出牌应答
type OutCardsReply struct {
	Id     int    //协议
	Result int    //0代表成功，其他代表失败
	Err    string //错误信息，输的代表为空
}

//玩家出牌广播结构体
type OutCardsBro struct {
	Id      int32   //协议号
	Cid     int32   //谁出的牌
	Type    int     //出牌类型
	Cards   []int   //牌
	Max     byte    //最大牌
	NextCid int32   //下一位出牌玩家
	IsDan   bool    //出牌玩家是否报单
	Hint    [][]int //提示 （牌型从大到小排列）
}

//玩家托管请求携带结构体
type TuoGuan struct {
	Id  int //协议号
	Ctl int //操作 1，托管  2,取消托管
}

//玩家托管广播结构体
type TuoGuanReply struct {
	Id  int32 //协议号
	Ctl int
	Cid int32 //谁托管
}

//玩家过牌广播结构体
type PassBro struct {
	Id   int32 //协议号
	Cid  int32
	Next int32   //下一位操作玩家 椅子id
	Hint [][]int //下一位玩家提示
}

//结算发送信息
type BalanceToClient struct {
	Id int //协议号
	// Booms         map[int32]int   //玩家炸弹数量集合（key为椅子id）
	// Balance       map[int32]int   //结算集合（key为椅子id）
	// QuanGuan      []int32         //被全关人椅子id
	// BaoPei        int32           //包赔玩家
	// Coins         map[int32]int64 //玩家金币（key为椅子id）
	// Handcards     map[int32][]int //剩余手牌展示（key为椅子id）
	PlayerMsgToBa []PlayerMsgToBa
}
type PlayerMsgToBa struct {
	Booms      int   //玩家炸弹数量集合（key为椅子id）
	Balance    int   //结算集合（key为椅子id）
	Coins      int64 //玩家金币（key为椅子id）
	Handcards  []int //剩余手牌展示（key为椅子id）
	IsQuanGuan bool
	BaoPei     bool
}

//游戏详情记录
type GameRecord struct {
	Id          int              //协议号
	GameId      int              `json:"gameId"`      //游戏id
	GradeId     int              `json:"gradeId"`     //场次id
	RoomId      int              `json:"roomId"`      //房间id
	GradeNumber int              `json:"gradeNumber"` //场次编号
	GameRoundNo string           `json:"gameRoundNo"` //游戏局号
	UserRecord  []GameRecordInfo `json:"userRecord"`
}
type GameRecordInfo struct {
	UserId          int64  `json:"userId"`          //用户id
	UserAccount     string `json:"userAccount"`     //用户名称
	Robot           bool   `json:"robot"`           //是否机器人
	CoinsBefore     int64  `json:"coinsBefore"`     //下注前金币
	PrizeCoins      int64  `json:"prizeCoins"`      //总盈亏
	CoinsAfter      int64  `json:"coinsAfter"`      //下注后金币
	BaseScore       int    `json:"baseScore"`       //底分
	SurPlusCardsNum int    `json:"surPlusCardsNum"` //剩余张数
	CoverBombNum    int    `json:"coverBombNum`     //被炸弹数
	BombNum         int    `json:"bombNum"`         //所处炸弹数
	CompensateNum   int64  `json:"compenstateNum"`  //包赔数
}

//重新连接结构体
type GInfoReConnectReply struct {
	Id            int32
	GameState     int         //游戏状态
	TimerNum      int         //状态剩余时间
	Cid           int32       //座位号id
	Seats         []GSeatInfo //所有玩家信息
	BScore        int
	JuHao         string
	Cards         []int             //手牌
	LiXian        []bool            //桌子中玩家托管情况
	Curcid        int32             //下一次出牌玩家
	TuoGuans      []bool            //桌子玩家托管情况
	CardsNum      []int             //手牌数
	StateTime     int               //状态时间
	DeskCardInfos []DeskCardInfo    //玩家桌面
	LastOutCards  LastOutCardsToCli // 上一轮出牌
	Hint          [][]int
}
type DeskCardInfo struct {
	DeskOutCard LastOutCardsToCli //桌面出的牌
	PlayerDo    int               //0,代表未操作  1，代表正在出牌 2,代表pass
}

type LastOutCardsToCli struct {
	Max   byte
	Type  int
	Cards []int
	Cid   int32
}

//玩家离开应答
type ExitReply struct {
	Id     int    //协议
	Result int    //0为成功，其他失败
	Err    string //成功时 为空
}

type DisOrRecBro struct {
	Id  int
	Cid int32
	Ctl int //1，代表离线  2，代表上线
}

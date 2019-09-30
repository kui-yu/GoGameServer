package main

type GBetCountInfo struct {
	Id           int   // 座位id
	DownBetValue int64 // 总下注金额
	UserBetValue int64 // 用户下注金额
}

//座位信息
type GClientBetCountInfo struct {
	Id           int   // 座位id
	DownBetTotal int64 // 总下注金额
}

// 发送给客户端的桌子信息
type GClientDeskInfo struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string

	JuHao     string                // 局号
	FangHao   string                // 房号
	Bets      []GClientBetCountInfo // 座位信息
	BetLevels []int64               // 下注级别

	PlayerMassage      PlayerMsg //用户信息
	AllBetCount        int64     //总下注金额
	AreaCoin           []int64   //区域金币
	GameStatus         int       // 游戏状态
	GameStatusDuration int64     // 当前状态持续时间毫秒
	Multiple           []float32 //倍数
	AreaMaxCoin        int64     // 限制区域最大下注
	Log                []int     //开奖记录
	Car                int       //上一把开奖结果
	ChangeCard         int
	Chair              map[int]*GUserInfoByChair //座位
	CanUserChip        int                       //可以使用的筹码
	// OldOther           []int64                   //其他玩家下注
	Head      string //玩家头像
	BetAll    int64  //玩家总下注
	Nick      string //名字
	CarName   string //开奖结果
	Index     int
	Double    float32
	ShengTime int //剩余时间
}

// 游戏状态
type GSGameStatusInfo struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string

	GameStatus         int   // 游戏状态
	GameStatusDuration int64 // 当前状态持续时间
}

//用户信息
type PlayerMsg struct {
	Uid          int64   //玩家uid
	MyUserAvatar string  // 用户头像
	MyUserName   string  // 用户昵称
	MyUserCoin   int64   // 用户金币
	MyDownBets   []int64 // 自己下注的集合
}

// 玩家列表
type GUserInfoReply struct {
	Id       int32
	UserInfo []GUserInfo
}

// 用户信息 (更多玩家请求)
type GUserInfo struct {
	Uid      int64
	Nick     string // 昵称
	Head     string // 头像
	TotBet   int64  // 总下注
	WinCount int32  // 赢取次数
	Coins    int64  // 当前金币
	Index    int    //当前用户排名
	Match    int    //局数
}

// 用户信息 (座位)
type GUserInfoByChair struct {
	Uid  int64
	Nick string // 昵称
	Head string // 头像
	// WinCount int64  // 赢取次数
	Coins int64 // 当前金币
}
type BetAgainReply struct {
	Id          int
	BetArea     []int64
	Result      int    //0成功，其他失败
	Err         string //错误信息
	Coins       int64  //剩余金币
	CanUserChip int    //可使用筹码
}

//玩家请求下注
type GADownBet struct {
	BetsIdx int // 下注区域索引(0-7)
	CoinIdx int // 下注金额索引(0-4)
}

//玩家请求下注返回
type GSDownBet struct {
	Id          int
	Result      int32 // 0 成功，其他失败
	Err         string
	PAreaCoins  []int64 //自己区域下注金币
	Coins       int64   //玩家剩余金币
	AreaId      int     //下注区域
	CoinId      int     //下注筹码Id
	CanUserChip int     //可以使用的筹码
	DownCoins   int64   //下注金额
}

//玩家请求下注通知
type GNDownBet struct {
	Id   int
	Bets [8]int64 //各区域下注
	// SeatBetList  [][]int64 // 座位玩家下注情况
	OtherBetList    []int64 // 除自己以外，其他玩家下注情况
	OldOtherBetList []int64 //老玩家下注
	PAreaCoins      []int64 // 自己总下注情况
	AllBets         int64   //区域总金币
}

//开奖通知
type GNLottery struct {
	Id       int
	Car      int     // 开奖结果
	Index    int     // 开奖结果下标（为转换)
	Double   float32 //开奖结果倍数
	DataTime string  //时间
}

// 系统结算
type GNBalance struct {
	Id             int32
	Results        map[int]GBetBalance //用户结算集合
	MyCoin         int64               // 用户金币
	ElseWinAndLose []int64             //其他用户赢取金币
	WinOrLoseCoins int64               //玩家输赢金币
	CanUserChip    int                 //可以使用的筹码
	// SeatWinCoins   []int64             // 座位玩家输赢
	Head    string //玩家头像
	BetAll  int64  //玩家总下注
	Nick    string //名字
	CarName string //开奖结果
}

// 单个位置结算
type GBetBalance struct {
	Bottom   int64 //区域下注金额
	Result   int64 //区域输赢
	MyBottom int64 //玩家下注金币数
	MyResult int64 //玩家输赢
	Win      bool  //玩家输赢  ,true  代表胜利，false 代表失败
}

// 玩家提示信息
type GSTips struct {
	Id  int32
	Msg string
}

//游戏开奖记录通知
type GNRecord struct {
	Id              int
	Record          []int // 游戏开奖记录
	OnlinePlayerNum int   // 在线玩家数
}

//
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GradeNumber string            `json:"gradeNumber"`
	GameRoundNo string            `json:"gameRoundNo"`
	LotteryCard int               `json:"lotteryCard"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId      int64    `json:"userId"`
	UserAccount string   `json:"userAccount"`
	Robot       bool     `json:"gradeNumber"`
	BetCoins    int64    `json:"betCoins"`    // 下注金币
	BetArea     [8]int32 `json:"betArea"`     // 区域下注情况
	PrizeCoins  int64    `json:"prizeCoins"`  // 赢取金币
	CoinsBefore int64    `json:"coinsBefore"` // 下注前金币
	CoinsAfter  int64    `json:"coinsAfter"`  // 结束后金币
}
type ChairUpdate struct { //座位变更通知结构体
	Id    int                      //协议号
	Chair map[int]GUserInfoByChair //座位
}

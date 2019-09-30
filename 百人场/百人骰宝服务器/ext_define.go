package main

//游戏开始阶段
type GSGameStart struct {
	Id int
	Time  int
	Round int
}

//停止下注
type GSGameStop struct {
	Id int
	Time int
}

//游戏结果以及历史记录的更新
type GSGameHistory struct {
	Id int
	Info    GameInfo
	InfoArr []GameInfo
}

//结算信息
type GSSettleInfo struct {
	Id int
	Uid        int64    //玩家uid
	Count      int64    //输赢的金币
	RealCount  int64    //连本金算的钱
	CountRate  int64    //扣除的费率
	Coin       int64    //结算后的玩家金币
	Place      []int    //赢的位置
	MaxBet     []int    //可下注范围
	GameResult GameInfo //游戏结果
}

//玩家进场=>客户端
type GSPlayerIn struct {
	Id      int
	Account string // 账号
	Uid     int64  // 用户ID
	Head    string // 头像

	Coin      int64      //玩家金币
	Round     int        //当前局数
	OnLine    int        //在线玩家
	History   []GameInfo //游戏历史
	DeskMoney int64      //桌面总投注
	MaxBet    []int      //最大下注区间
	Time      int        //房间的倒计时时间
	Stage     int        //当前游戏阶段
	Big       int64      //大
	Small     int64      //小
	Odd       int64      //单
	Even      int64      //双
}

//在线玩家
type GSOnLine struct {
	Id    int
	Count int
}

//玩家下注信息
type GABetInfo struct {
	Id    int64
	Big   int64
	Small int64
	Odd   int64
	Even  int64
}

//玩家下注群发=>客户端
type GSPlayerBetMass struct {
	Id        int
	DeskMoney int64
	Big       int64
	Small     int64
	Odd       int64
	Even      int64
}

//玩家下注 => 客户端
type GSPlayerBet struct {
	Id        int
	Big       int64
	Small     int64
	Odd       int64
	Even      int64
	Coin      int64
	DeskMoney int64
	MaxBet    []int //可下注数组
	PlayerBet map[int]int64
}

//开奖阶段
type GSGameResult struct {
	Id         int
	History    []GameInfo
	GameResult GameInfo
	IsWin      bool
	Count      int64
}

//可投注范围=>客户端
type GSMaxBet struct {
	Id   int
	Coin int64
	Bet  []int //可下注数组
}

//掉线重连
type LostConnection struct {
	Id        int
	History   []GameInfo
	DeskMoney int64
	MaxBet    []int
	Time      int
	Stage     int
	Big       int64
	Small     int64
	Odd       int64
	Even      int64
}

//请求返回
type GSGameBack struct {
	Id     int
	Result bool
	Err    string
}

// //数据库
// type SqlUpdate struct {
// 	Id    int64
// 	Coins int64
// }

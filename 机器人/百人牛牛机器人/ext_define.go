package main

// 定义游戏消息
const (
	MSG_GAME_NSTATUS_CHANGE = 410001 // 游戏状态改变
	MSG_GAME_QSEATDOWN      = 410002 // 410002 玩家请求坐下
	MSG_GAME_RSEATDOWN      = 410003 // 410003 玩家请求坐下的回复
	MSG_GAME_NSEATDOWN      = 410004 // 410004 座位信息改变通知
	MSG_GAME_BALANCE        = 410010 // 410010 结算
	MSG_GAME_QMANYUSER      = 410014 // 410014 请求更多玩家信息
	MSG_GAME_RMANYUSER      = 410015 // 410015 请求更多玩家信息回复
	MSG_GAME_QSEATUP        = 410020 // 玩家请求站立
	MSG_GAME_RSEATUP        = 410021 // 玩家请求站立回复
	MSG_GAME_QDOWNBET       = 410005 // 410005 玩家请求下注
	MSG_GAME_RDOWNBET       = 410006 // 410006 玩家请求下注回复
	MSG_GAME_NDOWNBET       = 410007 // 410007 玩家下注通知
	MSG_GAME_QDESKINFO      = 410018 // 请求游戏桌子信息
	MSG_GAME_RDESKINFO      = 410019 // 回复游戏桌子信息
)

//机器人游戏状态
const (
	GAME_STATUS_WAITSTART = 10 + iota // 等待开始
	GAME_STATUS_SEATBET               // 抢坐和坐下的人下注状态
	GAME_STATUS_FACARD                // 发牌状态
	GAME_STATUS_DOWNBTES              // 下注状态
	GAME_STATUS_OPENCARD              // 开牌状态
	GAME_STATUS_BALANCE               // 结算

	GAME_STATUS_WEB_ACCOUNT  // 请求登录账号
	GAME_STATUS_HALL_CONNECT // 连接大厅
	GAME_STATUS_HALL_LOGIN   // 请求登录大厅
	GAME_STATUS_GAME_ENTER   // 请求进入游戏
	GAME_STATUS_GAME_REENTER // 请求重新进入游戏
	GAME_STATUS_GAME_AUTO    // 请求匹配房间
	GAME_STATUS_GAME_DESK    // 房间信息
)

// 机器人事件通知
const (
	EVENT_ROBOT_SEATDOWN = 20000 + iota // 机器人坐下
)

// 桌子座位信息
// 发送给客户端的座位信息
type GClientSeatInfo struct {
	Id            int    // 座位id
	UserId        int64  // 用户id
	Name          string // 名称
	Avatar        string // 头像
	SeatDownCount int    // 坐下的次数
	DownBetTotal  int64  // 总下注金额
}

// 发送给客户端的桌子信息
type GClientDeskInfo struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string

	JuHao     string            // 局号
	FangHao   string            // 房号
	Seats     []GClientSeatInfo // 座位信息
	BetLevels []int64           // 下注级别

	MyUserAvatar string           // 用户头像
	MyUserName   string           // 用户昵称
	MyUserCoin   int64            // 用户金币
	MyDownBets   map[uint8]uint32 // 自己下注的集合

	GameStatus          int   // 游戏状态
	GameStatusDuration  int64 // 当前状态持续时间毫秒
	SeatDownMinCoinCond int   // 坐下条件
	SeatDownMinBetCond  int   // 坐下的人最低下注

	AreaMaxCoin         int // 限制区域最大下注
	AreaMaxCoinDownSeat int // 限制区域最大下注
	SeatUpTotalCount    int // 站立条件
}

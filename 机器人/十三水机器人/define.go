package main

// 定义机器人管理 消息
const (
	MSG_GAME_ROBOT_ONLINE    = 500000 + iota // 机器人上线
	MSG_GAME_ROBOT_NUMCHANGE                 // 机器人数量变化
)

const (
	MSG_HALL_ROBOT_LOGIN       = 300010 // 机器人登录
	MSG_HALL_ROBOT_LOGIN_REPLY = 300011 // 机器人登录应答
	MSG_HALL_LOGIN_REPLY       = 300002 // 大厅登录应答
	MSG_HALL_JOIN_GMAE         = 300005 // 加入游戏
	MSG_HALL_JOIN_GAME_REPLY   = 300006 // 加入游戏应答
	MSG_GAME_AUTO              = 400001 // 请求匹配房间
	MSG_GAME_AUTO_REPLY        = 400002 // 匹配房间回复
	MSG_GAME_RECONNECT         = 400010 // 400010断线重连
	MSG_GAME_RECONNECT_REPLY   = 400011 // 断线重连应答
	MSG_HALL_HEART             = 300007 // 心跳
	MSG_HALL_HEART_REPLY       = 300017 // 心跳应答
)

// 机器人事件通信
const (
	EVENT_CONT_ADDROBOT      = 10000 + iota // 机器人控制器添加机器人
	EVENT_CONT_DELROBOT                     // 机器人控制器删除机器人
	EVENT_CONT_CONFIG_CHANGE                // 机器人配置改变
	EVENT_CONT_ROBOTIDLE                    // 机器人闲置通知，代表可以退出房间
	EVENT_CONT_ROBOTSHIFT                   // 机器人替换
	EVENT_CONT_OFFLINEROBOT                 // 设置机器人下线
	EVENT_ROBOT_STOP                        // 机器人关闭
	EVENT_CONNECT_SUCCESS                   // 连接成功通知
)

// 发送的消息格式
type SendMsg struct {
	Id   uint32
	Data []byte
}

// 接收的消息格式
type RecvMsg struct {
	Id   uint32
	Data string
}

// 机器人发管理
type ResGameInfo struct {
	Id                uint32
	GroupId           uint32
	GameId            uint32
	GradeId           uint32
	Name              string
	GetRobotConfigUrl string
	PutRobotConfigUrl string
	CheckRobotUrl     string
	OfflineRobotUrl   string
	Forceoffroboturl  string
	Forceonroboturl   string
	RobotCount        uint32
}

// 管理发机器人
type RepBgInfo struct {
	Id                uint32
	BgRobotGetUrl     string // 请求获得机器人url
	BgRobotRestUrl    string // 请求归还机器人url
	BgRobotRestAllUrl string // 请求归还所有机器人url
	BgRobotAddCoinUrl string // 请求添加金币接口
	BgRobotToken      string // 秘钥
	BgRobotHallId     int    // 大厅
}

// 登录账号信息
type LoginInfo struct {
	Account  string
	Password string
}

// 控制器事件
type EventMsg struct {
	Id   int32       // 事件Id
	Data interface{} // 事件参数Event
}

// 用户信息
type UserInfo struct {
	Id           uint64 `json:"id"`
	Account      string `json:"account"`
	GameServerId int32  `json:"gameServerId"`
	Uid          int64
}

// 检查机器人单个信息
type CheckRobotItemInfo struct {
	Coin       int64  `json:"coin"`
	CreateTime int64  `json:"createTime"`
	Enable     bool   `json:"enable"`
	GameId     int    `json:"gameId"`
	GradeId    int    `json:"gradeId"`
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Portrait   string `json:"portrait"`
	ProfitCoin int64  `json:"profitCoin"`
	Sex        bool   `json:"sex"`
	Status     bool   `json:"status"`
	Token      string `json:"token"`
	UpdateTime int64  `json:"updateTime"`
	version    int    `json:"version"`
}

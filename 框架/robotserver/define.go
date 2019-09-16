package main

// 定义游戏消息
const (
	MSG_GAME_ROBOT_ONLINE    = 500000 + iota // 机器人上线
	MSG_GAME_ROBOT_NUMCHANGE                 // 机器人数量变化
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

// 发送给web后台的消息的游戏信息
type GameInfo struct {
	GroupId           uint32 `json:"groupId"`
	GameId            uint32 `json:"gameId"`
	GradeId           uint32 `json:"gradeId"`
	Name              string `json:"name"`
	GetRobotConfigUrl string `json:"getRobotConfigUrl"`
	PutRobotConfigUrl string `json:"putRobotConfigUrl"`
	CheckRobotUrl     string `json:"checkRobotUrl"`
	OfflineRobotUrl   string `json:"offlineRobotUrl"`
	Forceoffroboturl  string `json:"forceoffroboturl"`
	Forceonroboturl   string `json:"forceonroboturl"`
	RobotCount        uint32 `json:"robotCount"`
}

// 发送到机器人的配置信息
type RepBgInfo struct {
	Id                uint32
	BgRobotGetUrl     string // 请求获得机器人url
	BgRobotRestUrl    string // 请求归还机器人url
	BgRobotRestAllUrl string // 请求归还所有机器人url
	BgRobotAddCoinUrl string // 请求添加金币接口
	BgRobotToken      string // 秘钥
	BgRobotHallId     int    // 大厅
}

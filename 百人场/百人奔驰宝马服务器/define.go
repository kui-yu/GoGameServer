package main

import (
	"bytes"
	"encoding/json"

	// "fmt"
	"io/ioutil"
	"net/http"
)

const (
	MSG_HALL_START               = 300000
	MSG_HALL_LOGIN               = 300001 //大厅登录
	MSG_HALL_LOGIN_REPLY         = 300002 //大厅登录应答
	MSG_HALL_LEAVE               = 300003 //退出大厅
	MSG_HALL_LEAVE_REPLY         = 300004 //退出大厅应答
	MSG_HALL_JOIN_GMAE           = 300005 //加入游戏
	MSG_HALL_JOIN_GAME_REPLY     = 300006 //加入游戏应答
	MSG_HALL_HEART               = 300007 //心跳
	MSG_HALL_GETNEWESTCOIN       = 300008 //玩家获取最新的金币
	MSG_HALL_GETNEWESTCOIN_REPLY = 300009 //玩家更新金币应答
	MSG_HALL_ROBOT_LOGIN         = 300010 //机器人登录
	MSG_HALL_ROBOT_LOGIN_REPLY   = 300011 //机器人登录应答
	MSG_HALL_CREATE_FKROOM       = 300012 //创建房卡
	MSG_HALL_CREATE_FKROOM_REPLY = 300013 //创建房卡应答
	MSG_HALL_JOIN_FKROOM         = 300014 //输入房间号加入房间
	MSG_HALL_JOIN_FKROOM_REPLY   = 300015 //加入房间应答
	MSG_HALL_OTHERLOGIN_NOTIFY   = 300016 //账号被抢登陆
	MSG_HALL_HEART_REPLY         = 300017 //心跳应答
)

// 定义游戏消息，desk 相关
//以下为公共的id，每个游戏相关的id不可以在这里添加，请添加到ext_define.go中
const (
	MSG_GAME_START               = 400000 //游戏开始
	MSG_GAME_AUTO                = 400001 //自由匹配
	MSG_GAME_AUTO_REPLY          = 400002 //自由匹配应答
	MSG_GAME_JOIN                = 400003 //游戏加入（房卡）
	MSG_GAME_JOIN_REPLY          = 400004 //游戏加入应答（房卡）
	MSG_GAME_CREATE              = 400005 //游戏创建（房卡）
	MSG_GAME_CREATE_REPLY        = 400006 //游戏创建应答（房卡）
	MSG_GAME_LEAVE               = 400007 //离开游戏
	MSG_GAME_LEAVE_REPLY         = 400008
	MSG_GAME_END_NOTIFY          = 400009 //400009游戏结束通知
	MSG_GAME_RECONNECT           = 400010 //400010断线重连
	MSG_GAME_RECONNECT_REPLY     = 400011 //断线重连应答
	MSG_GAME_ONLINE_NOTIFY       = 400012 //400012离线和在线通知
	MSG_GAME_DISCONNECT          = 400013 //400013断线
	MSG_GAME_END_RECORD          = 400014 //400014游戏详情记录
	MSG_GAME_UPDATE_PLAYER_INFO  = 400016 //400016通知游服更改金币
	MSG_GAME_FK_CREATEDESK       = 400017 //创建房卡桌子
	MSG_GAME_FK_CREATEDESK_REPLY = 400018 //创建房卡桌子应答
	MSG_GAME_FK_JOIN             = 400019 //加入房卡场
	MSG_GAME_FK_JOIN_REPLY       = 400020 //加入房卡场应答
	MSG_GAME_UPDATEPLAYER_NOTIFY = 400021 //游戏内充值通知
	MSG_GAME_ADDCOIN_NOTIFY      = 400022 //修改金币
)

//每个游戏自己的状态请添加到ext_define.go中，以10开始
const (
	GAME_STATUS_FREE  = 0 + iota // 桌子空闲状态
	GAME_STATUS_START            // 游戏开始
	GAME_STATUS_END              //游戏结束
)

///////////////////////////////////////////////////////
// 游戏内部消息
type DkInMsg struct {
	Id   int32
	Uid  int64
	Data string
}

type DkOutMsg struct {
	Id   int32
	Uid  int64
	Data []byte
}

//数据库内部消息
type DbInMsg struct {
	Id   int32
	Data string
}

type DbOutMsg struct {
	Id   int32
	Data []byte
}

//大厅内部消息
type OutMsg struct {
	Id   int32
	Data []byte
}

type InMsg struct {
	Id   int32
	Uid  int64
	Data string
	Col  *Contoller
}

// 登录大厅需要的结构数据
type HMsgHallLogin struct {
	Id      int32
	Account string
	Gid     string
}

// 登录大厅返回的结构数据
type HMsgHallLoginReply struct {
	Id              int32   //协议id
	Result          int     //错误号
	Err             string  //错误信息
	Account         string  //账号
	Uid             int64   //用户uid
	Nick            string  //昵称
	Sex             int32   //性别
	Head            string  //头像
	Coin            int64   //金币
	GameSerId       int32   //加入的游戏信息
	AliPayAccount   string  //支付宝账号
	BankCard        string  //银行卡号
	Commission      float32 //佣金
	Money           float64 //金钱
	UnReadNum       int     //未读信息数量
	UserName        string  //真实姓名
	ForbiddenEnable bool    //用户冻结
	FrozenEnable    bool    //用户禁用
	BindPassword    bool    // 是否绑定密码
}

type HMsgHallRobotLogin struct {
	Id  int32
	Gid string
}

type HMsgHallRobotLoginReply struct {
	Id      int32
	Result  int32
	Err     string
	Uid     int64
	Coin    int64
	Account string
	Head    string
	Sex     bool
}

type HMsgHallJoinGame struct {
	Id        int32 //协议号
	GameType  int32 //游戏类型
	RoomType  int32 //房间类型
	GradeType int32 //场次类型
}

type HMsgHallJoinGameReply struct {
	Id     int32  //协议号
	Result int32  //0成功，其余失败
	Err    string //失败原因
}

type HMsgHallGetNewestCoin struct {
	Id int32
}

type HMsgHallGetNewestCoinReply struct {
	Id     int32 //协议号
	Result int32 //0成功，其他失败
	Err    string
	Coin   int64 //更新的金币数量
}

type HMsgHallCreateFkRoom struct {
	Id        int32  //协议号
	GameType  int32  //游戏类型
	RoomType  int32  //房间类型
	GradeType int32  //场次类型
	FkInfo    string //房间信息
}

type HMsgHallCreateFkRoomReply struct {
	Id     int32  //协议号
	Result int32  //0成功，其他失败
	Err    string //失败原因
}

type HMsgHallJoinFkRoom struct {
	Id        int32  //协议号
	GameType  int32  //游戏类型
	RoomType  int32  //房间类型
	GradeType int32  //场次类型
	FkNo      string //房号
}

type HMsgHallJoinFkRoomReply struct {
	Id        int32  //协议号
	Result    int32  //0成功，其他失败
	Err       string //失败原因
	FkNo      string //房间号
	GameType  int32  //游戏类型
	RoomType  int32  //房间类型
	GradeType int32  //场次类型
}

type HMsgHallQiangDengNotify struct {
	Id int32 //协议号
}

type HMsgHallHeart struct {
	Id int32
}

type HMsgHallHeartReply struct {
	Id int32
}

////////////////////////////////////////////
//座位信息
type GSeatInfo struct {
	Uid   int64
	Nick  string
	Ready bool
	Cid   int32 //椅子号
	Sex   int32
	Head  string
	Lv    int32
	Coin  int64
}

//自由匹配
//大厅-》游服
type GAutoGame struct {
	Id          int32
	Account     string
	Uid         int64
	Nick        string
	Sex         int32
	Head        string
	Lv          int32
	Coin        int64
	Token       string
	Robot       bool
	HierarchyId int32
}

//自由匹配
//客户端-》大厅
type GAutoGameToHall struct {
	Id        int32
	IsRobot   bool
	GameType  int32
	RoomType  int32
	GradeType int32
}

//自由匹配应答，此外还有一个匹配消息和游戏相关的（斗地主为GInfoAutoGameReply）
type GAutoGameReply struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string
}

type GFkJoinToGame struct {
	Id      int32
	Account string
	Uid     int64
	Nick    string
	Sex     int32
	Head    string
	Lv      int32
	Coin    int64
	Token   string
	Robot   bool
	FkNo    string //房间号
}

type GFkJoinToHall struct {
	Id        int32
	GameType  int32
	RoomType  int32
	GradeType int32
	FkNo      string //房间号
}

type GFkJoinReply struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string
}

//游戏离开请求
type GLeave struct {
	Id int32
}

//游戏离开应答
type GLeaveReply struct {
	Id     int32
	Result int32 //0成功，其他失败
	Cid    int32
	Uid    int64
	Err    string
	Token  string
	Robot  bool // 是否机器人
}

type GReConnect struct {
	Id int32
}

//短线重连应答
type GReConnectReply struct {
	Id     int32
	Result int32
	Err    string
}

//游戏结束发送的消息，用于存储玩家分数局号，此外还有其他结束信息通过ext_define.go中定义应答
type GGameEnd struct {
	Id          int32          //协议号
	GameId      int            `json:"gameId"`
	GradeId     int            `json:"gradeId"`
	RoomId      int            `json:"roomId"`
	GameRoundNo string         `json:"gameRoundNo"`
	Mini        bool           `json:"mini"`
	Round       int            `json:"round"`
	UserCoin    []GGameEndInfo `json:"userCoin"`
	SetLeave    int32          //是否设置离开，0离开，1不离开
	ActiveUid   int64          //主动保存这些数据的用户
}

type GGameEndInfo struct {
	UserId      int64   `json:"userId"`
	UserAccount string  `json:"userAccount"`
	BetCoins    int64   `json:"betCoins"`
	ValidBet    int64   `json:"validBet"`
	PrizeCoins  int64   `json:"prizeCoins"`
	Robot       bool    `json:"robot"`
	WaterProfit float64 `json:"waterProfit"`
	WaterRate   float64 `json:"waterRate"`
}

//上线和掉线通知
type GOnLineNotify struct {
	Id    int32
	Cid   int32
	State int32 //1上线，2掉线
}

type GUpdatePlayerInfo struct {
	Id         int32
	PlayerInfo Player
}

type GFkCreateDesk struct {
	Id     int32
	FkNo   string //房间号
	FkInfo string //房间规则信息
}

type GFkCreateDeskReply struct {
	Id        int32
	Result    int32
	Err       string
	FkNo      string //房间号
	GameType  int32
	RoomType  int32
	GradeType int32
}

type GUpdatePlayerNotify struct {
	Id   int32
	Coin int64
	Head string
}

//////////////////////////////////////////////////////////////////////////
//数据库代理应答
type DbGetPlayerRsp struct {
	Msg  string       `json:"msg"`
	Code int          `json:"code"`
	Data DbPlayerData `json:"data"`
}

type DbPlayerData struct {
	Account         string  `json:"account"`
	AliPayAccount   string  `json:"alipayAccount"`
	BankCard        string  `json:"bankCard"`
	Coin            int64   `json:"coin"`
	Commission      float32 `json:"commission"`
	GameInfoId      int     `json:"gameInfoId"`
	GameServerId    int     `json:"gameServerId"`
	Uid             int64   `json:"id"`
	Money           float64 `json:"money"`
	Portrait        string  `json:"portrait"`
	Sex             bool    `json:"sex"`
	UnreadNum       int     `json:"unreadNum"`
	UserName        string  `json:"userName"`
	ForbiddenEnable bool    `json:"forbiddenEnable"`
	FrozenEnable    bool    `json:"frozenEnable"`
	BindPassword    bool    `json:"bindPassword"`
	HierarchyId     int     `json:"hierarchyId"`
}

type DbRobotData struct {
	Coin int64  `json:"coin"`
	Uid  int64  `json:"id"`
	Name string `json:"name"`
	Head string `json:"portrait"`
	Sex  bool   `json:"sex"`
}

type DbGetRobotRsp struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data DbRobotData `json:"data"`
}

type DbGetGameServerRsp struct {
	Msg  string         `json:"msg"`
	Code int            `json:"code"`
	Data DbGetGamesData `json:"data"`
}

type DbGetGamesData struct {
	Game []DbGetGameServerData `json:"game"`
}

type DbGetGameServerData struct {
	Bscore      int     `json:"bscore"`
	Display     bool    `json:"display"`
	Enable      bool    `json:"enable"`
	GameId      int     `json:"gameId"`
	GameName    string  `json:"gameName"`
	GradeId     int     `json:"gradeId"`
	GradeName   string  `json:"gradeName"`
	Sid         int     `json:"id"`
	LimitHigh   int     `json:"limitHigh"`
	LimitLower  int     `json:"limitLower"`
	Maintenance bool    `json:"maintenance"`
	MaxLines    int     `json:"maxLines"`
	MaxTimes    int     `json:"maxTimes"`
	OnlineMax   int     `json:"onlineMax"`
	OnlineMin   int     `json:"onlineMin"`
	Rate        float64 `json:"rate"`
	Restrict    int     `json:"restrict"`
	RoomId      int     `json:"roomId"`
	RoomName    string  `json:"roomName"`
	Sence       string  `json:"sence"`

	GameConfig DbGameServerConfig `json:"gameConfig"`
}

type DbGameServerConfig struct {
	// HierarchyGameRate []DbHierarchyGameRate `json:"hierarchyGameRate"` //玩家层级概率
	IntervalGameRate []DbIntervalGameRate `json:"intervalGameRate"` //库存层级概率
	CurrentStock     int64                `json:"currentStock"`     //当前库存值
	InitialStock     int64                `json:"initialStock"`     //初始库存
	GoalStock        int64                `json:"goalStock"`        //目标库存
	MiniRate         float64              `json:"miniRate"`         //小游戏胜率
	MaxRobot         int                  `json:"maxRobot"`         //单桌最大的机器人数量
	RobotWait        int                  `json:"robotWait"`        //等待几秒机器人进入
}

type DbHierarchyGameRate struct {
	Rate        float32 `json:"rate"`
	HierarchyId int     `json:"hierarchyId"`
}

type DbIntervalGameRate struct {
	Rate          float32 `json:"rate"`
	IntervalEnd   int64   `json:"intervalEnd"`
	IntervalStart int64   `json:"intervalStart"`
}

type DbGetRoleRateListRsp struct {
	Msg  string                `json:"msg"`
	Code int                   `json:"code"`
	Data []DbHierarchyGameRate `json:"data"`
}

type GameTypeDetail struct {
	GameType  int32 //游戏类型
	RoomType  int32 //房间类型
	GradeType int32 //场次类型
}

type GToHAddCoin struct {
	Id   int32
	Uid  int64
	Coin int64
}

///////////////////////////////////////////
func SendRequest(url string, data interface{}, method string, token string) (string, error) {
	//
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer([]byte(b))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	if method == "GET" {
		req.Header.Set("Content-Type", "text/plain")
	} else {
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
	}

	if token != "" {
		req.Header.Set("token", token)
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func SendRequestByString(url string, data string, method string, token string) (string, error) {
	body := bytes.NewBuffer([]byte(data))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	if method == "GET" {
		req.Header.Set("Content-Type", "text/plain")
	} else {
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
	}

	if token != "" {
		req.Header.Set("token", token)
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	return string(result), nil
}

//
func MarshalGameAllType(roomtype, gradetype, gametype int32) int32 {
	return roomtype<<23 + gradetype<<14 + gametype
}

func UnmarshalGameAllType(id int32) (int32, int32, int32) {
	return id >> 23, (id >> 14) & 0x1FF, id & 0x3FFF
}

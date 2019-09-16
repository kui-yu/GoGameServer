package main

//大厅消息
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

//每个游戏自己的状态请添加到ext_define.go中，以10开始
const (
	GAME_STATUS_FREE  = 0 + iota // 桌子空闲状态
	GAME_STATUS_START            // 游戏开始
	GAME_STATUS_END              //游戏结束
)

//阶段消息
//自己定义的状态，从10开始
const (
	GAME_STATUS_CALL   = 10 + iota // 叫分阶段
	GAME_STATUS_GETMSG             //抢地主阶段
	GAME_STATUS_READ               //准备阶段
	GAME_STATUS_PLAY               // 游戏操作阶段
)

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

const (
	MSG_GAME_INFO_STAGE = 410030 //阶段消息
)

// 系统消息，desk 相关
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

//消息自定义
const (
	MSG_GAME_INFO_START           = 410000 + iota //游戏开始额外信息
	MSG_GAME_INFO_READY                           //玩家准备410001
	MSG_GAME_INFO_READY_REPLY                     //玩家准备消息返回410002
	MSG_GAME_INFO_SEND_NOTIFY                     //发牌410003
	MSG_GAME_INFO_CALL                            //叫分410004
	MSG_GAME_INFO_CALL_REPLY                      //410005
	MSG_GAME_INFO_BANKER_NOTIFY                   //定庄410006
	MSG_GAME_INFO_OUTCARD                         //出牌410007
	MSG_GAME_INFO_OUTCARD_REPLY                   //410008
	MSG_GAME_INFO_PASS                            //过410009
	MSG_GAME_INFO_PASS_REPLY                      //4100010
	MSG_GAME_INFO_ROOM_NOTIFY                     //房间信息通知4100011
	MSG_GAME_INFO_TUOGUAN                         //410012托管
	MSG_GAME_INFO_TUOGUAN_REPLY                   //410013
	MSG_GAME_INFO_AUTO_REPLY                      //410014,游戏随机匹配成功的数据
	MSG_GAME_INFO_RECONNECT_REPLY                 //断线重连应答 410015
	MSG_GAME_INFO_END_NOTIFY                      //410016游戏结束应答
	MSG_GAME_INFO_LEAVE_NOTIFY                    //玩家离开通知410017
	MSG_GAME_INFO_ERR                             // 错误信息410018
	MSG_GAME_INFO_GETMSG                          //抢地主410019
	MSG_GAME_INFO_GETMSG_REPLY                    //抢地主应答 410020
)

//lgh --2019/03/04
const (
	MSG_GAME_INFO_STAGE = 410030 //阶段消息
)

//自己定义的状态，从10开始
const (
	GAME_STATUS_CALL = 10 + iota // 叫分阶段
	GAME_STATUS_PLAY             // 游戏操作阶段11

)

// 登录大厅需要的结构数据
type HMsgHallLogin struct {
	Id      int32
	Account string
	Gid     string
}

// 登录大厅返回的结构数据
type HMsgHallLoginReply struct {
	Id      int32
	Account string
	Uid     int64
	Nick    string
	Sex     int32
	Head    string
	Lv      int32
	Coin    int32
	GameId  int32
}

type HMsgHallCreateFkRoom struct {
	Id        int32  //协议号
	GameType  int32  //游戏类型
	RoomType  int32  //房间类型
	GradeType int32  //场次类型
	FkInfo    string //房间信息
}

//房卡配置信息
type GATableConfig struct {
	GameModle  int   //配置-专场  1,娱乐专场（积分模式）  2，土豪专场（金币模式）
	PlayerNum  int   //配置-牌局人数			（二人斗地主只有两个人）
	MatchCount int   //配置-局数
	GameType   int   //配置-玩法模式  1，正常模式  2,，二人斗地主 3，闪电斗地主  2，癞子模式
	PayType    int   //配置-支付条件 1，房主支付
	BaseScore  int   //配置-底分
	CallType   int   //配置-叫分方式  1，叫分  2，抢地主
	Boom       int   //配置-炸弹  3，3炸   4，4炸   5，5炸   999`，无上限
	CanSelect  []int //可选， 里面有1， 不可三带一对  2，不可4带两队  3，显示剩余手牌数量   没有选择默认为0
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

type GFkJoinReply struct {
	Id     int32
	Result int32 //0成功，其他失败
	Err    string
}

//所有匹配玩家信息
type GSInfoAutoGame struct {
	Id   int32
	Seat []GSSeatInfo
}

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

//阶段时间
type GSStageInfo struct {
	Id        int32
	Stage     int32
	StageTime int32
}

//叫牌玩家通知
type GNCallPlayer struct {
	Id      int
	ChairId int32
}

//游戏内部消息
type DkInMsg struct {
	Id   int32
	UID  int64
	Data string
}

//叫地主
type GCallMsg struct {
	Id    int32
	Coins int32
}

//游戏开始
type GGameStartNotify struct {
	Id    int32
	Round int
}

//定庄通知
type GBankerNotify struct {
	Id     int32
	Banker int32
	DiPai  []byte
	Double int32
}

//发牌通知
type GGameSendCardNotify struct {
	Id         int32
	Rount      int
	Cid        int32 //此人开始叫地主
	HandsCards []byte
}

//抢地主
type GGetMsg struct {
	Id     int32
	GetMsg int32 //是否抢（叫）地主  1，叫地主 2，抢地主 3，不叫 4，不抢
}

//抢地主应答
type GGetMsgReply struct {
	Id        int32
	Cid       int32 //轮到哪一个玩家叫分
	CallOrGet int32 //接下来是抢地主 还是 叫地主  1,抢地主 2，叫地主
	End       int32 //抢地主是否结束  1,未结束  2，结束
}

//匹配信息返回
type GInfoAutoGameReply struct {
	Id   int32
	Seat []GSeatInfo
}

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

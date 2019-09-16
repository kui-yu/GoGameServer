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
	MSG_GAME_SETCURSTOCK_NOTIFY  = 400023 //设置当前库存值
	MSG_GAME_ADDROOMCARD_NOTIFY  = 400024 //消耗房卡
)

//消息自定义
const (
	MSG_GAME_INFO_START           = 410000 + iota //游戏开始额外信息
	MSG_GAME_INFO_SEND_NOTIFY                     //发牌410001
	MSG_GAME_INFO_CALL                            //叫分410002
	MSG_GAME_INFO_CALL_REPLY                      //410003
	MSG_GAME_INFO_BANKER_NOTIFY                   //定庄410004
	MSG_GAME_INFO_OUTCARD                         //出牌410005
	MSG_GAME_INFO_OUTCARD_REPLY                   //410006
	MSG_GAME_INFO_PASS                            //过410007
	MSG_GAME_INFO_PASS_REPLY                      //410008
	MSG_GAME_INFO_ROOM_NOTIFY                     //房间信息通知410009
	MSG_GAME_INFO_TUOGUAN                         //410010托管
	MSG_GAME_INFO_TUOGUAN_REPLY                   //410011
	MSG_GAME_INFO_AUTO_REPLY                      //410012,游戏随机匹配成功的数据
	MSG_GAME_INFO_RECONNECT_REPLY                 //断线重连应答
	MSG_GAME_INFO_END_NOTIFY                      //410014游戏结束应答
	MSG_GAME_INFO_LEAVE_NOTIFY                    //玩家离开通知410015
)
const (
	CARD_COLOR = 0xF0 //花色掩码
	CARD_VALUE = 0x0F //数值掩码
	// 王
	Card_King_1 = iota + 0x41
	Card_King_2
)

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
	RoomCard        int64   //房卡数量
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

//lgh --2019/03/04
const (
	MSG_GAME_INFO_STAGE = 410030 //阶段消息
)

//自己定义的状态，从10开始
const (
	GAME_STATUS_CALL = 10 + iota // 叫分阶段
	GAME_STATUS_PLAY             // 游戏操作阶段11

)

//自由匹配应答，此外还有一个匹配消息和游戏相关的（斗地主为GInfoAutoGameReply）
type GAutoGameReply struct {
	Id       int32
	Result   int32 //0成功，其他失败
	CostType int   //1金币，2代币
	Err      string
}

//玩家匹配响应
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

//游戏房间信息
type GGameInfoNotify struct {
	Id     int32
	BScore int32
	MaxBei int32
	JuHao  string
}

//发牌通知
type GGameSendCardNotify struct {
	Id         int32
	Cid        int32  //此人开始叫地主
	HandsCards []byte //手牌
}

//阶段和阶段时间
type GStageInfo struct {
	Id        int32
	Stage     int32
	StageTime int32
}

//叫地主
type GCallMsg struct {
	Id    int32
	Coins int32
}

type GCallMsgReply struct {
	Id    int32 //协议号
	Cid   int32 //哪位玩家叫得分
	Coins int32 //叫几分
	End   int32 //叫分是否结束
}

//托管
type GTuoGuan struct {
	Id  int32 //协议号
	Ctl int32 //1托管，2取消托管
}

//托管应答
type GTuoGuanReply struct {
	Id     int32 //协议号
	Ctl    int32
	Result int32 //结果
	Err    string
	Cid    int32 //谁托管
}

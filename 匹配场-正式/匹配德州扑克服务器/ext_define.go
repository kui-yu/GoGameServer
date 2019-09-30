package main

// 自己定义的游服id从410000开始
// 游戏状态使用id表示，方便客户端使用 前缀N通知 Q请求 R回复
const (
	MSG_GAME_INFO_START         = 410000 + iota // 游戏开始额外信息
	MSG_GAME_NGameStatus                        // 410001 通知游戏状态
	MSG_GAME_NGameInfo                          // 410002 通知桌子信息 #当前消息废弃，走断线重连
	MSG_GAME_NGameUserChange                    // 410003 玩家信息改变
	MSG_GAME_NGameRandRank                      // 410004 随机庄家
	MSG_GAME_NGameBlind                         // 410005 盲注
	MSG_GAME_NGameHoleCards                     // 410006 发玩家的两张底牌
	MSG_GAME_NGamePublicCards                   // 410007 发公共牌
	MSG_GAME_NGameOperate                       // 410008 通知玩家操作
	MSG_GAME_QGameOperate                       // 410009 玩家请求操作
	MSG_GAME_RGameOperate                       // 410010 回复玩家请求操作
	MSG_GAME_NGameOperateResult                 // 410011 通知玩家操作结果
	MSG_GAME_NGameJackpotChange                 // 410012 奖池变化
	MSG_GAME_NGameResult                        // 410013 结算
	MSG_GAME_NGameReconnectInfo                 // 410014 断线重连信息
	MSG_GAME_NGameSetCoin                       // 410015 通知玩家设置携带金币
	MSG_GAME_QGameSetCoin                       // 410016 玩家设置携带金币
	MSG_GAME_RGameSetCoin                       // 410017 回复玩家携带金币
	MSG_GAME_NDeskUpdate                        // 410018 桌子更新
)

// 牌的类型 0x11=方块A
const (
	Card_Fang    = 0x10
	Card_Mei     = 0x20
	Card_Hong    = 0x30
	Card_Hei     = 0x40
	Card_King    = 0x50 // |14,15 小王，大王
	Card_Invalid = 0xFF // 无效的牌
)

// 游戏状态
const (
	GameStatusWaitStart   = 10 + iota // 等待开始游戏
	GameStatusRandBank                // 随机庄家，下盲注
	GameStatusHoleCards               // 发给玩家的两张牌
	GameStatusFlopCards               // 头三张公共牌
	GameStatusTurnCards               // 第四张公共牌
	GameStatusRiverCards              // 第五张公共牌
	GameStatusUserOperate             // 玩家操作
	GameStatusResults                 // 结算
)

// 牌组合类型
const (
	CardGroupHighCard   = 1 + iota //高牌
	CardGroupOnePair               //一对
	CardGroupTwoPair               //两对
	CardGroupThreeT                //三条
	CardGroupStraight              //顺子
	CardGroupFlush                 //同花
	CardGroupFullhouse             //三张+一对
	CardGroupFourT                 //四条
	CardGroupSFlush                //同花顺
	CardGroupRoyalFlush            //皇家同花顺
)

// 用户状态
const (
	UserStateGameIn    = 1 + iota //游戏中
	UserStateWaitStart            //等待下一局开始
	UserStateSettCoin             //设置携带金币
)

// 操作权限
const (
	OperateAuthQP = 0x1     //弃牌
	OperateAuthJZ = 0x10    //加注
	OperateAuthGZ = 0x100   //跟注
	OperateAuthKP = 0x1000  //开牌
	OperateAuthSH = 0x10000 //梭哈
)
const TimerId = 0x0F

var StageDefines []int = []int{
	GameStatusHoleCards,
	GameStatusFlopCards,
	GameStatusTurnCards,
	GameStatusRiverCards,
	GameStatusResults,
}

type GCUserInfo struct {
	Uid          int64
	NickName     string
	Avatar       string
	State        int
	Sid          int
	Coin         int64
	DownCoins    int64
	IsBank       bool
	IsFold       bool
	IsAllIn      bool
	Online       bool
	Cards        []int
	CurrStageOpt int
}

type GCGameStatusInfo struct {
	GameStatus int //游戏状态
	OverTime   int //结束时间
}

type GCardGroupInfo struct {
	GroupType int
	Cards     []int
}

type GGameRecordInfo struct {
	UserId      int64  `json:"userId"`
	UserAccount string `json:"userAccount"`
	Robot       bool   `json:"robot"`
	CoinsBefore int64  `json:"coinsBefore"`
	BetCoins    int64  `json:"betCoins"`
	Coins       int64  `json:"coins"`
	CoinsAfter  int64  `json:"coinsAfter"`
	Score       int    `json:"score"`
	Multiple    int    `json:"multiple"`
	Landlord    bool   `json:"landlord"`
}

type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

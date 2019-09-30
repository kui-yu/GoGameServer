package main

const (
	MSG_GAME_INFO_START              = 410000 + iota
	MSG_GAME_INFO_PLAY               // 410001 启动
	MSG_GAME_INFO_PLAY_NOTIFY        // 410002 启动反馈
	MSG_GAME_INFO_BONUS_START        // 410003 进入小游戏
	MSG_GAME_INFO_BONUS_START_NOTIFY // 410004 进入结果
	MSG_GAME_INFO_BONUS              // 410005 小游戏开箱子
	MSG_GAME_INFO_BONUS_NOTIFY       // 410006 开箱子结果
	MSG_GAME_INFO_BONUS_END          // 410007 小游戏结束
	MSG_GAME_INFO_END_NOTIFY         // 410008 游戏结束
	MSG_GAME_INFO_AUTO_REPLY         // 410009 匹配广播
	MSG_GAME_INFO_LEAVE              // 410010 玩家离开
	MSG_GAME_INFO_LEAVE_REPLY        // 410011 玩家离开
	MSG_GAME_INFO_EXIT_LIMIT_HIGHT   // 410012 金币过多
)

const (
	TIMER_BONUS = 10000 + iota // 奖金游戏倒计时
)

const (
	TIMER_BONUS_NUM = 10 // 奖金游戏时间
)

//所有匹配玩家信息
type GSInfoAutoGame struct {
	Id   int32
	Seat []GSeatInfos
}

type GSeatInfos struct {
	Uid    int64
	Nick   string
	Ready  bool
	Cid    int32 //椅子号
	Sex    int32
	Head   string
	Lv     int32
	Coin   int64
	Bscore int64
}

//////////////////////////////////////////////////////////
type GGameStartNotify struct {
	Id int32
}

//////////////////////////////////////////////////////////
type GGameExitNotify struct {
	Id int32
}

// 开始：线数和每条线的金币底分倍数
type GMsgPlay struct {
	Lines int64
	Coins int64
}

type GMsgPlayNotify struct {
	Id         int32
	Scenes     []byte
	Lines      [18][]byte
	IsShow     []byte
	Coins      int64
	Win        int64
	BonusCount int64
	MinBonus   int32
}

type GMsgBonusStartNotify struct {
	Id        int32 // 协议号
	BoxCount  int64 // 宝箱个数
	MaxChoose int64 // 最大开启宝箱数
	Times     int32 // 时间
}

type GMsgBonus struct {
	BoxIndex int32 // 宝箱下标
}

type GMsgBonusNotify struct {
	Id       int32 // 协议号
	BoxIndex int32 // 宝箱下标
	Coins    int64 // 用户金币
	Bonus    int64 // 宝箱金币
}

type GMsgBonusEndNotify struct {
	Id    int32   // 协议号
	Bonus []int64 // 宝箱金币列表
}

//
type GGameRecord struct {
	Id          int32             //协议号
	GameId      int               `json:"gameId"`
	GradeId     int               `json:"gradeId"`
	RoomId      int               `json:"roomId"`
	GameRoundNo string            `json:"gameRoundNo"`
	UserRecord  []GGameRecordInfo `json:"userRecord"`
}

type GGameRecordInfo struct {
	UserId      int64      `json:"userId"`
	UserAccount string     `json:"userAccount"`
	BetCoins    int64      `json:"betCoins"`    // 下注金币总额
	Scenes      []byte     `json:"scenes"`      // 图标情况
	Lines       [18][]byte `json:"lines"`       // 中奖线情况
	Pow         int64      `json:"pow"`         // 中奖倍率
	BetLines    int64      `json:"betLines"`    // 押线数
	Bcoins      int64      `json:"bcoins"`      // 底分
	PrizeCoins  int64      `json:"prizeCoins"`  // 赢取金币
	CoinsBefore int64      `json:"coinsBefore"` // 开始时候多少金币
	CoinsAfter  int64      `json:"coinsAfter"`  // 结束时候多少金币
	Robot       bool       `json:"robot"`
}

package main

// 定义游戏消息
const (
	MSG_GAME_INFO_STAGE             = 410001 //阶段消息
	MSG_GAME_INFO_AUTO_REPLY        = 410002 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_HANDINFO_REPLY    = 410004 //发送手牌信息
	MSG_GAME_INFO_PLAY              = 410005 //玩家摆牌
	MSG_GAME_INFO_SETTLE_INFO_REPLY = 410007 //结算
)

const (
	GAME_STATUS_START = 1
	STAGE_PLAY        = 12 //玩牌阶段
	STAGE_SETTLE      = 13 //结算阶段
)

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

//所有匹配玩家信息
type GInfoAutoGameReply struct {
	Id   int32
	Seat []GSeatInfo
}

//手牌消息
type GSHandInfo struct {
	Id           int
	ChairId      int32
	HandCards    []int
	SpecialType  int
	SpecialCards []int
}

//结算玩家消息
type GSettlePlayerInfo struct {
	Uid          int64
	ChairId      int32
	PlayCards    []int   //结算摆牌
	SpecialType  int     //特殊牌型
	SpecialScore int     //特殊得分
	NormalTypes  []int   //牌型数组[头墩牌型，中墩牌型，底分牌型]
	WinCoinList  []int   //比分总得分数组 [头墩得分，中墩得分，底分得分，总得分/特殊得分]
	WinCoins     int64   //总输赢
	Coins        int64   //身上金币
	NormalScores []int   //普通得分
	ShootList    []int32 //打枪{位置1，位置2}
	ShootScoress [][]int //打枪分数{[头墩得分，中墩得分，底分得分]，[头墩得分，中墩得分，底分得分]}
}

//结算消息
type GSSettleInfos struct {
	Id            int
	AllWinChairId int32 //全垒打
	PlayInfo      []GSettlePlayerInfo
}

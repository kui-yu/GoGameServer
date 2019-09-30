package main

// 定义游戏消息
const (
	MSG_GAME_INFO_AUTO_REPLY = 410012 //410012,游戏随机匹配成功的数据
	MSG_GAME_INFO_STAGE      = 410000 //阶段消息
	MSG_GAME_INFO_SETTLE     = 410006 //通比牛牛 结算
	MSG_GAME_INFO_CALL       = 410001 //通比牛牛 叫分
	MSG_GAME_INFO_PLAY       = 410005 //通比牛牛 开牌
	MSG_GAME_INFO_BET_LIST   = 410007 //筹码列表
)

const (
	GAME_STATUS_START = 1
	STAGE_DEAL        = 11 //发牌阶段
	STAGE_PLAY        = 12 //玩牌阶段
	STAGE_SETTLE      = 14 //结算阶段
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

//筹码列表
type GCallListMsg struct {
	Id         int
	BetListCnt int   //下注数量
	BetList    []int //下注列表
}

//结算子结构体2
type GWinInfo struct {
	Uid     int64
	ChairId int32
	WinCoin int64
	Coins   int64
}

//结算结构体
type GWinInfosReply struct {
	Id         int32
	WinChairId int32
	InfoCount  int32
	Infos      []GWinInfo
}

package main

// 定义游戏消息
const (
	MSG_GAME_INFO_STAGE                 = 410001 //阶段消息
	MSG_GAME_INFO_AUTO_REPLY            = 410002 //410002,游戏随机匹配成功的数据
	MSG_GAME_INFO_BANKER_REPLY          = 410008 //庄家通知
	MSG_GAME_INFO_SETTLE_INFO_END_REPLY = 410015 //总结算
	MSG_GAME_INFO_CALL_LIST             = 410016 //叫庄列表
)

const (
	STAGE_CALL = 11 //抢庄阶段
	STAGE_PLAY = 12 //下注阶段
)

//阶段时间
type GSStageInfo struct {
	Id        int
	Stage     int
	StageTime int
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

//所有匹配玩家信息
type GInfoAutoGameReply struct {
	Id   int32
	Seat []GSeatInfo
}

//结算玩家信息
type GSSettlePlayInfoEnd struct {
	Uid      int64
	ChairId  int32 //位置信息
	WinCoins int64 //输赢得分
	Coins    int64
}

//结算
type GSSettleInfoEnd struct {
	Id        int
	PlayInfos []GSSettlePlayInfoEnd //玩家输赢集合
}

//玩家抢庄结果
type GSPlayerCallBank struct {
	Id              int
	Banker          int32
	BankerList      []int
	BankerMultiples int   //庄家倍数
	BetListCnt      int   //下注数量
	BetList         []int //下注列表
}

type GSCallList struct {
	Id          int
	CallListCnt int
	CallList    []int
}

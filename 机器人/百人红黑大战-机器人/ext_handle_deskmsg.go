package main

import (
	"encoding/json"
	"fmt"
	"logs"
	"math/rand"
	"time"
)

const (
	MSG_GAME_INFO_BET_NOTIFY      = 410004 //下注中410004
	MSG_GAME_INFO_BET             = 410005 //下注410005
	MSG_GAME_INFO_AUTO_REPLY      = 410013 // 410013座位信息
	MSG_GAME_INFO_RECONNECT_REPLY = 410014 // 410014房间信息
	MSG_GAME_INFO_INTO            = 410020 // 410020请求AUTO_REPLY剩余消息
)

//下注区间
var BetList = []int64{100, 500, 1000, 2000, 3000, 5000, 6000, 10000, 50000, 100000}

// type AutoGameReply struct {
// 	Id       int32
// 	Result   int32 //0成功，其他失败
// 	CostType int   //1金币，2代币
// 	Err      string
// }

//匹配完成后续处理
type GInfoAutoGameReply1 struct {
	Id        int32    // 协议号
	PlayerNum int32    // 房间人数
	SeatList  []GSInfo // 座位玩家
}

//房间信息
type GInfoReConnectReply1 struct {
	Id            int32    // 协议号
	GameState     int32    // 游戏状态
	RoomId        string   // 房号
	GameId        string   // 局号
	GameLimit     int64    // 限红
	BetList       []int64  // 下注金币限制
	GameCount     int32    // 当前牌局
	TAreaCoins    []int64  // 总下注金币
	SeatList      []GSInfo // 座位玩家
	PAreaCoins    []int64  // 当前玩家下注金币
	PCoins        int64    // 当前玩家当前金币
	BetArea       []bool   // 可下注区域
	CardList      []int32  // 展示牌列表
	RedCard       []int32  // 红方牌
	BlackCard     []int32  // 黑方牌
	WinArea       []bool   // 赢取区域
	RunChart      []int32  // 输赢走势
	CardTypeChart []int    // 牌型记录走势
	Rtype         int      //红方牌型
	Btype         int      //黑方牌型
	Timer         int32    // 阶段时间（毫秒）
	LimitCoinId   int32    //当用户金币小于下注池的某一值时，获得此值得ID
}

type GSInfo struct {
	Nick  string //昵称
	Head  string //头像
	Coins int64  //金币
}

// 游戏下注
type GGameBetNotify struct {
	Id      int32  // 协议号
	Timer   int32  // 阶段时间（毫秒）
	BetArea []bool // 可下注区域
}

// 下注
type GGameBet struct {
	Id     int32 // 协议号
	MsgId  int32 // 消息系号，防止重复(新开局系号1开始（断线重连也一样）)
	AreaId int32 // 下注区域Id
	CoinId int32 // 下注金额Id
}

// //机器人自由匹配进场
// func (this *ExtRobotClient) AutoReply(msg string) {
// 	res := &AutoGameReply{}
// 	err := json.Unmarshal([]byte(msg), res)
// 	if err != nil {
// 		logs.Debug("json:AutoReply()错误:%v", err)
// 	}
// 	logs.Debug("AutoReply!:%v", res)
// }

//匹配完成后续处理
func (this *ExtRobotClient) AutoReplyFinal(msg string) {
	res := &GInfoAutoGameReply1{}
	err := json.Unmarshal([]byte(msg), res)
	if err != nil {
		logs.Debug("json:AutoReplyFinal()错误:%v", err)
	}
	// logs.Debug("AutoReplyFinal:%v", res)

}

//房间信息
func (this *ExtRobotClient) RoomInfo(msg string) {
	res := &GInfoReConnectReply1{}
	err := json.Unmarshal([]byte(msg), res)
	if err != nil {
		logs.Debug("json:RoomInfo()错误:%v", err)
	}
	// logs.Debug("RoomInfo:%v", res)
}
func (this *ExtRobotClient) BetTime(msg string) {
	res := &GGameBetNotify{}
	err := json.Unmarshal([]byte(msg), res)
	if err != nil {
		logs.Debug("json:BetTime()错误:%v", err)
	}
	timer := res.Timer //在规定时间内进行下注
	for timer > 0 {
		time.Sleep(time.Second * 2)
		if this.Coin < 0 {
			break
		}
		GGB := GGameBet{
			Id: MSG_GAME_INFO_BET,
		}
		rand.Seed(time.Now().UnixNano())
		areaId := int32(rand.Intn(3)) + 1
		coinid := int32(rand.Intn(4))
		GGB.AreaId = areaId
		GGB.CoinId = coinid
		this.AddMsgNative(MSG_GAME_INFO_BET, GGB)
		timer -= 1000
		fmt.Println("下注成功!", areaId, coinid)
	}
}

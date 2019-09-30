package main

import (
	"encoding/json"
	"logs"
	"math/rand"
	"time"
)

const (
	MSG_GAME_INFO_STAGE_INFO      = 410000 //410000//阶段消息
	MSG_GAME_INFO_DESKINFO_REPLAY = 410002
	MSG_GAME_INFO_PLAYER_BET      = 410003 //410003//玩家下注
	STAGE_GAME_BET                = 12     //下注阶段
	STAGE_GAME_STOP_BET           = 13     //13停止下注
)

//阶段消息
type GSGameStageInfo struct {
	Id    int
	Stage int
	Time  int
}

//下注区间
var betArr = []int{100, 500, 1000, 2000, 3000, 5000, 6000, 10000, 50000, 100000, 500000}

//玩家进场=>客户端
type GSPlayerIn struct {
	Id          int
	GameId      string       //局号
	MaxBet      int64        //限红
	ManyPlayer  []ManyPlayer //桌面围观玩家
	AllPlayer   int          //在线玩家
	Stage       int          //当前游戏阶段
	GameTrend   [][]Trend    //游戏走势，根据索引分别为庄、黑、红、梅、方， Trend.Player 0为闲家，1为庄家
	DeskInfos   DeskInfo
	PlayerInfos PlayerInfo
}
type DeskInfo struct {
	Time        int          //房间的倒计时时间
	BetArr      []int64      //下注筹码
	PlaceBetAll []int64      //区域下注，0到3分别是黑红梅方
	HandCards   []Card       //庄家闲家手牌，索引0为庄家手牌
	Players     []ManyPlayer //桌面的6个玩家
}
type PlayerInfo struct {
	Account    string  // 账号
	Uid        int64   // 用户ID
	Head       string  // 头像
	Coin       int64   //玩家金币
	IsDouble   bool    //是否翻倍
	BetArrAble int     //可下注筹码
	PlaceBet   []int64 //自己区域下注，0到3分别是黑红梅方
}

//桌面玩家
type ManyPlayer struct {
	Head            string //头像
	Account         string //账号
	Uid             int64
	Coins           int64 //金币
	Round           int   //累积下注局数
	AccumulateBet   int64 //累积下注
	AccumulateCoins int64 //累积输赢
}

//走势记录
type Trend struct {
	CardType int //牌型
	Player   int //1为庄家，0为闲家
}

//牌结构
type Card struct {
	CardValue []int
	CardType  int
	Multiple  int64
}

//玩家下注信息
type GABetInfo struct {
	Id    int
	Place int
	Coin  int
}

func (this *ExtRobotClient) DeskInfo(msg string) {
	res := &GSPlayerIn{}
	err := json.Unmarshal([]byte(msg), res)
	if err != nil {
		logs.Debug("json:DeskInfo()错误:%v", err)
	}

}
func (this *ExtRobotClient) StageInfo(msg string) {
	res := &GSGameStageInfo{}
	err := json.Unmarshal([]byte(msg), res)
	if err != nil {
		logs.Debug("json:StageInfo()错误:%v", err)
	}
	if res.Stage == STAGE_GAME_BET {
		this.DownBet(res.Time) //下注
	}
}
func (this *ExtRobotClient) DownBet(t int) {
	rand.Seed(time.Now().UnixNano())
	d := GABetInfo{
		Id: MSG_GAME_INFO_PLAYER_BET,
	}
	for t > 0 {
		time.Sleep(time.Second * 2)
		d.Coin = rand.Intn(int(this.RobotClient.UserInfo.Uid)) % len(betArr)
		d.Place = rand.Intn(int(this.RobotClient.UserInfo.Uid)) % 4
		this.AddMsgNative(MSG_GAME_INFO_PLAYER_BET, d)
		t2 := (rand.Intn(int(this.RobotClient.UserInfo.Uid)) % 3) + 1
		t -= t2
	}
}

//获取随机下注金额和区域
func getRandBet() (d *GABetInfo, t int) {
	rand.Seed(time.Now().UnixNano())
	place := rand.Intn(4)          //区域
	coin := rand.Intn(len(betArr)) //金额
	//返回下注结构
	return &GABetInfo{
		Id:    MSG_GAME_INFO_PLAYER_BET,
		Place: place,
		Coin:  coin,
	}, rand.Intn(3) + 1
}

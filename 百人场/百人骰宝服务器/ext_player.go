package main

import (
	"encoding/json"
	"fmt"
	"log"
	"logs"
	"sync"
)

type ExtPlayer struct {
	Player
	Coin      int64 //玩家金币
	MaxBetArr []int //玩家最大下注区间
	IsBet     bool  //是否有下注
	LostBet    int64    //输的位置的下注总和
	WinBets    int64    //赢的位置的下注总和
	BetInfo    []int64  //下注的位置
	GameResult GameInfo //游戏结果
	Lock       sync.RWMutex
}

//游戏入口
func (this *ExtDesk) InitExtData() {
	//操作逻辑
	this.Handle[MSG_GAME_INFO_PLAYER_IN] = this.PlayerIn     //玩家进场
	this.Handle[MSG_GAME_INFO_PLAYER_BET] = this.PlayerBet   //玩家下注
	this.Handle[MSG_GAME_RECONNECT] = this.HandleReconnect   //掉线重连
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisConnect //掉线信息
	this.Handle[MSG_GAME_INFO_BACK] = this.HandleGameBack    //返回大厅
	//游戏开始
	this.GameServeLoop(BET_TIMER)
}

//初始化玩家
func (this *ExtPlayer) InitPlayer() {
	this.IsBet = false
	this.WinBets = 0
	this.LostBet = 0
	this.BetInfo = make([]int64, 4)
}

//玩家进场
func (this *ExtDesk) PlayerIn(p *ExtPlayer, m *DkInMsg) {
	fmt.Printf("玩家%d进场\n", p.Uid)
	p.BetInfo = make([]int64, 4)
	p.Coin = p.Coins
	//发送房间信息给刚进场的玩家
	obj := new(GSPlayerIn)
	obj.Account = p.Account
	obj.Head = p.Head
	obj.Uid = p.Uid
	obj.Id = MSG_GAME_INFO_PLAYER_IN_REPLY
	obj.Coin = p.Coin
	//读写锁
	p.Lock.RLock()
	obj.OnLine = len(this.Players)
	obj.Round = this.Round
	obj.Big = this.Big
	obj.Small = this.Small
	obj.Odd = this.Odd
	obj.Even = this.Even
	obj.DeskMoney = this.DeskMoney
	obj.History = this.History
	obj.Stage = this.Stage
	obj.Time = this.TList[0].T
	//解锁
	p.Lock.RUnlock()
	//获取可下注区间
	arr := p.MaxBet(0)
	p.MaxBetArr = arr
	obj.MaxBet = arr
	p.SendNativeMsg(MSG_GAME_INFO_PLAYER_IN_REPLY, obj)
	//更新在线玩家,跳过自己
	online := new(GSOnLine)
	online.Id = MSG_GAME_INFO_ONLINE_PLAYER
	online.Count = obj.Round
	for _, v := range this.Players {
		if v.Uid == p.Uid {
			continue
		}
		v.SendNativeMsg(MSG_GAME_INFO_ONLINE_PLAYER, online)
	}
}

//判断玩家可下注范围
func (this *ExtPlayer) MaxBet(count int64) []int {
	arr := make([]int, 0)
	if this.Coin-count >= 100 {
		arr = []int{1, 5, 10, 50, 100}
	} else if this.Coin-count >= 50 {
		arr = []int{1, 5, 10, 50}
	} else if this.Coin-count >= 10 {
		arr = []int{1, 5, 10}
	} else if this.Coin-count >= 5 {
		arr = []int{1, 5}
	} else if this.Coin-count >= 1 {
		arr = []int{1}
	}
	return arr
}

//玩家下注
func (this *ExtDesk) PlayerBet(p *ExtPlayer, m *DkInMsg) {
	//不是下注阶段返回
	if this.Stage != STAGE_GAME_START_BET {
		return
	}
	//携带金额达不到最小下注金额返回
	if len(p.MaxBetArr) <= 0 && !p.Robot {
		return
	}
	p.IsBet = true
	//解析玩家下注信息
	obj := new(GABetInfo)
	err := json.Unmarshal([]byte(m.Data), obj)
	if err != nil {
		log.Panic("json解析失败")
	}
	fmt.Printf("玩家%d下注:================\n", p.Uid)
	fmt.Println("big: ", obj.Big)
	fmt.Println("small: ", obj.Small)
	fmt.Println("odd: ", obj.Odd)
	fmt.Println("even: ", obj.Even)
	//玩家下注结构分析
	p.Lock.RLock()
	defer p.Lock.RUnlock()
	if obj.Big > 0 {
		p.BetInfo[PLACE_BIG] += obj.Big
		p.Coin -= obj.Big
		this.Big += obj.Big
	} else if obj.Small > 0 {
		p.BetInfo[PLACE_SMALL] += obj.Small
		p.Coin -= obj.Small
		this.Small += obj.Small
	} else if obj.Odd > 0 {
		p.BetInfo[PLACE_ODD] += obj.Odd
		p.Coin -= obj.Odd
		this.Odd += obj.Odd
	} else if obj.Even > 0 {
		p.BetInfo[PLACE_EVEN] += obj.Even
		p.Coin -= obj.Even
		this.Even += obj.Even
	} else {
		logs.Debug("玩家没有下注")
		return
	}
	bet := obj.Big + obj.Small + obj.Even + obj.Odd
	this.DeskMoney += bet

	//群发下注信息给其他客户端
	for _, v := range this.Players {
		if v.Uid == p.Uid {
			continue
		}
		v.SendNativeMsg(MSG_GAME_INFO_PLAYER_BET_MASS, &GSPlayerBetMass{
			Id:        MSG_GAME_INFO_PLAYER_BET_MASS,
			Big:       this.Big,
			Small:     this.Small,
			Odd:       this.Odd,
			Even:      this.Even,
			DeskMoney: this.DeskMoney})
	}
	//发送下注信息给客户端
	maxarr := p.MaxBet(bet)
	info := new(GSPlayerBet)
	info.Id = MSG_GAME_INFO_PLAYER_BET_REPLAY
	info.Big = this.Big
	info.Small = this.Small
	info.Odd = this.Odd
	info.Even = this.Even
	info.Coin = p.Coin
	info.MaxBet = maxarr
	info.DeskMoney = this.DeskMoney
	info.PlayerBet = map[int]int64{
		PLACE_BIG:   p.BetInfo[PLACE_BIG],
		PLACE_SMALL: p.BetInfo[PLACE_SMALL],
		PLACE_ODD:   p.BetInfo[PLACE_ODD],
		PLACE_EVEN:  p.BetInfo[PLACE_EVEN],
	}
	p.MaxBetArr = maxarr
	p.SendNativeMsg(MSG_GAME_INFO_PLAYER_BET_REPLAY, info)
}

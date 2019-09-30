package main

import (
	"fmt"
	"logs"
	"math/rand"
	"time"
)

type FSMLottery struct {
	Mark int

	EDesk       *ExtDesk
	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMLottery) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMLottery) Run() {
	DebugLog("游戏状态-开奖")

	this.EndDateTime = GetTimeMS() + int64(gameConfig.StateInfo.RunCarLogoTime)

	this.addListen() // 添加监听
	this.EDesk.GameState = GAME_STATUS_LOTTERY
	this.EDesk.SendGameState(GAME_STATUS_LOTTERY, int64(gameConfig.StateInfo.RunCarLogoTime))

	this.EDesk.AddTimer(GAME_STATUS_LOTTERY, gameConfig.StateInfo.RunCarLogoTime/1000, this.TimerCall, nil)

	this.RunLogo()
}

func (this *FSMLottery) TimerCall(d interface{}) {
	fmt.Println("开奖玩进入结算")
	this.EDesk.RunFSM(GAME_STATUS_BALANCE)
}

func (this *FSMLottery) GetMark() int {
	return this.Mark
}

func (this *FSMLottery) Leave() {
	this.removeListen()
}

func (this *FSMLottery) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()
	return remainTimeMS
}

func (this *FSMLottery) addListen() {}

func (this *FSMLottery) removeListen() {}

func (this *FSMLottery) RunLogo() {
	var result int
	//计算玩家总下注
	var totalbetcoin int64
	for _, p := range this.EDesk.Players {
		if p.Robot {
			continue
		}
		totalbetcoin += p.TotalBet
	}
	//风控判断
	logs.Debug("当前库存：", CD, "目标库存：", CalPkAll(StartControlTime, time.Now().Unix()))
	if GetCostType() == 2 || totalbetcoin == 0 { //控制机器人75%的胜率(3/4)
		w := RandInt64(4)
		if w >= 1 { //123庄输
			result = this.EDesk.allotCard(false)
		} else { //0赢
			result = this.EDesk.allotCard(true)
		}
	} else if CD-CalPkAll(StartControlTime, time.Now().Unix()) >= 0 { //纯随机
		result = int(RandInt64(8))
	} else { //控制玩家输
		result = this.EDesk.allotCard(true)
	}
	this.EDesk.GameResult = result
	fmt.Println("开奖结果::", result)
	suiji := rand.Intn(3)
	this.EDesk.Car = result
	res := []int{0, 2, 4, 6, 1, 3, 5, 7}
	for k, v := range res {
		if result == k {
			result = v
			break
		}
	}

	info := GNLottery{
		Id:       MSG_GAME_INFO_NLOTTERY,
		Car:      result + suiji*8,
		Index:    this.EDesk.GameResult,
		Double:   CarTypeMultiple[this.EDesk.GameResult],
		DataTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	this.EDesk.ChangeCar = info.Car
	DebugLog("开奖结果：", result)
	this.EDesk.BroadcastAll(MSG_GAME_INFO_NLOTTERY, &info)
}

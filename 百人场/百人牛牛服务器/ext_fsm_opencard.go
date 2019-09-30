/**
* 开牌状态
**/
package main

import (
	"encoding/json"
	"logs"
	"time"
)

type FSMOpenCard struct {
	UpMark int
	Mark   int
	EDesk  *ExtDesk

	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMOpenCard) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMOpenCard) GetMark() int {
	return this.Mark
}

func (this *FSMOpenCard) Run(upMark int) {
	DebugLog("游戏状态：开牌")
	logs.Debug("游戏状态：开牌")

	this.UpMark = upMark

	timeId := gameConfig.GameStatusTimer.OpenCardId
	timeMs := int64(gameConfig.GameStatusTimer.OpenCardMS)

	this.EndDateTime = GetTimeMS() + timeMs

	this.addListener()                          // 添加监听
	this.EDesk.SendGameState(this.Mark, timeMs) // 发送桌子状态

	this.EDesk.AddTimer(timeId, int(timeMs)/1000, this.TimerCall, nil)

	this.openCard() // 开牌
}

func (this *FSMOpenCard) Leave() {
	this.removeListener()
}

func (this *FSMOpenCard) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_BALANCE)
}

func (this *FSMOpenCard) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()

	return remainTimeMS
}

// 添加网络监听
func (this *FSMOpenCard) addListener() {
	this.EDesk.Handle[MSG_GAME_QDOWNBET] = this.recvRSeatBet
}

// 接收到玩家下注
func (this *FSMOpenCard) recvRSeatBet(p *ExtPlayer, d *DkInMsg) {
	req := GClientQDownBet{}
	json.Unmarshal([]byte(d.Data), &req)
	this.EDesk.UserDownBet(p, req.SeatIdx, req.CoinIdx, true)
}

// 删除网络监听
func (this *FSMOpenCard) removeListener() {
}

func (this *FSMOpenCard) onUserOnline(p *ExtPlayer) {

}

func (this *FSMOpenCard) onUserOffline(p *ExtPlayer) {

}

// 开牌
func (this *FSMOpenCard) openCard() {
	// 分配牌（原始）
	//this.EDesk.allotCard()
	// 分配牌（新的）
	this.controlCard()

	// 发送牌信息
	data := &GClientNFaCard{
		Id:    MSG_GAME_OPENCARD,
		Cards: this.EDesk.CardGroupArray,
	}
	this.EDesk.SendNotice(MSG_GAME_OPENCARD, data, true, nil)
}

// 控牌
func (this *FSMOpenCard) controlCard() {
	//统计本局真实玩家下注量
	var hasRealBet bool
	for i := 0; i < gameConfig.GameLimtInfo.SeatCount; i++ {
		if this.EDesk.Seats[i].UserBetValue > 0 {
			hasRealBet = true
			break
		}
	}

	//先判断本房间是否为体验房或是本桌是否有真实玩家下注，如果是体验房或是没有真实玩家下注，则可随机分牌（以庄家多输为主），
	//反之，则进入库存控制
	if GetCostType() == 2 || !hasRealBet {
		//控制庄家输，但有一定概率不控制。
		var killRate int = gameConfig.GameCtrlInfo.BankerLoseProb
		rand, _ := GetRandomNum(0, 100)
		if rand < killRate {
			this.EDesk.AllotCardBankerLose()
		} else {
			this.EDesk.AllotCardRand()
		}
		return
	}

	//计算目标库存值
	var currTime int64 = time.Now().Unix()
	var destStock int64 = CalPkAll(StartControlTime, currTime)

	logs.Debug("目标库存:%v", destStock)
	logs.Debug("当前库存:%v", CD)

	if CD >= destStock {
		//如果当前库存大于或等于目标库存，则无控制必要
		this.EDesk.AllotCardRand()
	} else {
		//控制庄家赢，但有一定概率不控制。
		var killRate int = gameConfig.GameCtrlInfo.BankerWinProb
		rand, _ := GetRandomNum(0, 100)
		if rand < killRate {
			this.EDesk.InAllWinAllot()
		} else {
			this.EDesk.AllotCardRand()
		}
	}
}

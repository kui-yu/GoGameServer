/**
* 发牌状态
**/
package main

import (
	"encoding/json"
	"logs"
)

type FSMFaCard struct {
	UpMark int
	Mark   int
	EDesk  *ExtDesk

	EndDateTime int64 // 当前状态的结束时间
}

func (this *FSMFaCard) InitFSM(mark int, extDest *ExtDesk) {
	this.Mark = mark
	this.EDesk = extDest
}

func (this *FSMFaCard) GetMark() int {
	return this.Mark
}

func (this *FSMFaCard) Run(upMark int) {
	DebugLog("游戏状态：发牌")
	logs.Debug("游戏状态：发牌")

	this.UpMark = upMark

	timeId := gameConfig.GameStatusTimer.FaCardId
	timeMs := int64(gameConfig.GameStatusTimer.FaCardMS)

	this.EndDateTime = GetTimeMS() + timeMs

	this.addListener()                          // 添加监听
	this.EDesk.SendGameState(this.Mark, timeMs) // 发送桌子状态

	this.EDesk.AddTimer(timeId, int(timeMs)/1000, this.TimerCall, nil)

	this.faCard() // 发牌
}

func (this *FSMFaCard) Leave() {
	this.removeListener()
}

func (this *FSMFaCard) TimerCall(d interface{}) {
	this.EDesk.RunFSM(GAME_STATUS_DOWNBTES)
}

func (this *FSMFaCard) getRestTime() int64 {
	remainTimeMS := this.EndDateTime - GetTimeMS()

	return remainTimeMS
}

// 添加网络监听
func (this *FSMFaCard) addListener() {
	this.EDesk.Handle[MSG_GAME_QDOWNBET] = this.recvRSeatBet
}

// 接收到玩家下注
func (this *FSMFaCard) recvRSeatBet(p *ExtPlayer, d *DkInMsg) {
	req := GClientQDownBet{}
	json.Unmarshal([]byte(d.Data), &req)
	this.EDesk.UserDownBet(p, req.SeatIdx, req.CoinIdx, true)
}

// 删除网络监听
func (this *FSMFaCard) removeListener() {
}

func (this *FSMFaCard) onUserOnline(p *ExtPlayer) {

}

func (this *FSMFaCard) onUserOffline(p *ExtPlayer) {

}

// 发牌
func (this *FSMFaCard) faCard() {
	seatCount := gameConfig.GameLimtInfo.SeatCount

	this.EDesk.ShuffleCard()
	for i := 0; i < seatCount; i++ {
		//发牌
		cards := this.EDesk.SendCard(2)
		this.EDesk.CardGroupArray[i] = CardGroupInfo{
			MaxCard:       0,
			CardGroupType: CardGroupType_None,
			Cards:         []int{int(cards[0]), int(cards[1])},
		}
	}

	//-----old-----
	// cards := ShuffleCard()
	// this.EDesk.DownCards = cards
	// this.EDesk.DownCardIdx = uint8(seatCount * 2)

	// cardIdx := 0
	// for i := 0; i < seatCount; i++ {
	// 	this.EDesk.CardGroupArray[i] = CardGroupInfo{
	// 		MaxCard:       0,
	// 		CardGroupType: CardGroupType_None,
	// 		Cards:         []int{int(cards[cardIdx]), int(cards[cardIdx+1])},
	// 	}
	// 	cardIdx += 2
	// }

	data := &GClientNFaCard{
		Id:    MSG_GAME_FACARD,
		Cards: this.EDesk.CardGroupArray,
	}

	this.EDesk.SendNotice(MSG_GAME_FACARD, data, true, nil)
}

package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

//匹配数据
func (this *ExtRobotClient) HandlerGameInfos(d string) {
	// DebugLog("匹配数据", d)
	gameInfo := GInfoReConnectReply{}
	json.Unmarshal([]byte(d), &gameInfo)

	//金币赋值
	this.Coin = gameInfo.PCoins
	this.BetList = gameInfo.BetList

	this.ResetBetArea()

	this.GameIn = true
}

//准备阶段
func (this *ExtRobotClient) HandlerGameInit(d string) {
	// DebugLog("匹配数据", d)
	gameInfo := GGameReadyNotify{}
	json.Unmarshal([]byte(d), &gameInfo)

	this.ResetBetArea()
}

func (this *ExtRobotClient) ResetBetArea() {

	this.AreaList = []int{}
	//下注区域比率
	for i := 0; i < 60; i++ {
		this.AreaList = append(this.AreaList, INDEX_TIAN_WIN)
		this.AreaList = append(this.AreaList, INDEX_TIAN_LOSS)
		this.AreaList = append(this.AreaList, INDEX_DI_WIN)
		this.AreaList = append(this.AreaList, INDEX_DI_LOSS)
		this.AreaList = append(this.AreaList, INDEX_REN_WIN)
		this.AreaList = append(this.AreaList, INDEX_REN_LOSS)
	}

	for i := 0; i < 10; i++ {
		this.AreaList = append(this.AreaList, INDEX_BANKER_TIAN)
		this.AreaList = append(this.AreaList, INDEX_BANKER_ZHIZUN)
	}

	this.PlayList = []int{}
	//下注筹码比率
	for i := 0; i < 35; i++ {
		this.PlayList = append(this.PlayList, 1)
	}
	for i := 0; i < 35; i++ {
		this.PlayList = append(this.PlayList, 2)
	}
	for i := 0; i < 20; i++ {
		this.PlayList = append(this.PlayList, 3)
	}
	for i := 0; i < 10; i++ {
		this.PlayList = append(this.PlayList, 4)
	}
	for i := 0; i < 10; i++ {
		this.PlayList = append(this.PlayList, 5)
	}
}

//下注阶段
func (this *ExtRobotClient) HandlerGameBet(d string) {
	// DebugLog("下注阶段", d)
	betInfo := GGameBetNotify{}
	json.Unmarshal([]byte(d), &betInfo)
	// areaList := ListShuffle(this.AreaList)
	// coinList := ListShuffle(this.PlayList)
	rand.Seed(time.Now().UnixNano())

	second := int(betInfo.Timer / 1000)
	if len(this.AreaList) != 0 && len(this.PlayList) != 0 {

		var msgId int = 1
		for i := 1; i < second; i++ {
			//休眠1秒
			time.Sleep(time.Second * time.Duration(1))
			randNum := rand.Perm(100)
			// DebugLog("随机数", randNum[0])
			if randNum[0] > 50 {
				continue
			}
			indexArea := rand.Int63n(int64(len(this.AreaList)))
			indexCoin := rand.Int63n(int64(len(this.PlayList)))
			areaId := this.AreaList[indexArea]
			coinId := this.PlayList[indexCoin]
			// DebugLog("下注区域ID", areaId, "；下注金币ID", coinId)
			if coinId-1 < 0 || len(this.BetList) <= coinId-1 {
				continue
			}
			//机器人金币充足
			if this.Coin >= this.BetList[coinId-1] {
				//发送下注信息
				this.AddMsgNative(MSG_GAME_INFO_BET, struct {
					Id     int // 协议号
					MsgId  int // 消息系号，防止重复(新开局系号1开始（断线重连也一样）)
					AreaId int // 下注区域Id
					CoinId int // 下注金额Id
				}{
					Id:     MSG_GAME_INFO_BET,
					MsgId:  msgId,
					AreaId: areaId,
					CoinId: coinId,
				})
				this.Coin -= this.BetList[coinId-1]
			}
			msgId++

			if msgId > 5 {
				break
			}
		}
	}
}

//游戏结算
func (this *ExtRobotClient) HandlerGameEnd(d string) {
	// DebugLog("匹配数据", d)
	win := GGameAwardNotify{}
	json.Unmarshal([]byte(d), &win)
	//金币赋值
	this.Coin = win.PCoins
	this.GameIn = false
	//通知可以离开
	controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
}

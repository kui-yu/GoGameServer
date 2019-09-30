package main

import (
	"math/rand"
	"time"
)

//抢庄阶段
func (this *ExtDesk) StageChoiceBanker(d interface{}) {
	this.GameState = STAGE_BANKER_MULTIPLE
	this.BroadStageTime(gameConfig.Stage_Banker_Multiple_Timer)
	this.AddTimer(1, gameConfig.Stage_Banker_Multiple_Timer, this.StageChoiceBankerEnd, "")
}

//抢庄阶段结束
func (this *ExtDesk) StageChoiceBankerEnd(d interface{}) {
	//判断是否抢庄超时，没操作默认不抢
	req := GChoiceMultipleREQ{Id: MSG_GAME_INFO_BANKER_MULTIPLE}
	for _, v := range this.Players {
		if v.BankerInfos.IsChoice == 0 {
			req.Uid = v.Uid
			GDeskMgr.AddNativeMsg(MSG_GAME_INFO_BANKER_MULTIPLE, v.Uid, req)
		}
	}
}

//选出庄家，倍数一样金币多的概率大
func (this *ExtDesk) ChoiceBankerByCoins() *ExtPlayer {
	rand.Seed(time.Now().UnixNano())    //随机数种子
	if len(this.ChoiceBankerArr) == 1 { //只有一个人抢
		return this.ChoiceBankerArr[0].Player
	} else if len(this.ChoiceBankerArr) == 0 { //都不抢的话就随机
		return this.Players[rand.Intn(len(this.Players))]
	}
	//都抢庄的话，金币最多的玩家得庄的概率大
	coinsArr := make([]int, 0)
	//生成金币区间
	for _, v := range this.ChoiceBankerArr {
		if len(coinsArr) == 0 {
			coinsArr = append(coinsArr, int(v.Player.Coins))
		} else {
			coinsArr = append(coinsArr, int(v.Player.Coins)+coinsArr[len(coinsArr)-1])
		}
	}
	randNum := rand.Intn(coinsArr[len(coinsArr)-1] + 1)
	for k, v := range coinsArr {
		if randNum <= v {
			return this.ChoiceBankerArr[k].Player
		}
	}
	return nil
}

package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"time"

	"bl.com/paigow"
	//"bl.com/util"
)

// 停止下注
func (this *ExtDesk) TimerStopBet(d interface{}) {
	this.Lock()
	defer this.Unlock()

	tAreaCoins := this.GetAreaCoinsList()

	// 消息公共部分
	sd := GGameStopBetNotify{
		Id:         MSG_GAME_INFO_STOP_BET_NOTIFY,
		Timer:      int32(gameConfig.Timer.StopBetNum) * 1000,
		TAreaCoins: tAreaCoins,
	}

	// 座位玩家下注信息
	// sd.SeatBetList = this.SeatMgr.GetSeatNewBetList()

	// 新下注，排除座位玩家
	// otherNewBetList := this.SeatMgr.GetOtherNewBetList()

	// 每个人的下注都不一样  需要单独处理
	players := this.SeatMgr.GetUserList(len(this.Players))
	for _, v := range players {
		p := v.(*ExtPlayer)
		// newBetList := p.GetNewBetList()
		// p.ColAreaCoins()
		// TPAreaCoins := p.GetTotBetList()

		sd.PAreaCoins = p.GetNTAreaCoinsList()

		sd.OtherBetList = this.SeatMgr.GetOtherNewBetList2(p.Uid)

		// if this.SeatMgr.IsOnSeat(p) {
		// 	sd.OtherBetList = otherNewBetList
		// } else {
		// 	sd.OtherBetList = util.LessInt64List(otherNewBetList, newBetList)
		// }
		// logs.Debug("停止下注: ", gameConfig.Timer.StopBetNum)
		p.SendNativeMsg(MSG_GAME_INFO_STOP_BET_NOTIFY, sd)
	}

	for _, user := range this.Players {
		user.ColAreaCoins()
	}

	var areaCards [][]int32
	//发四组牌
	for areaId := 0; areaId < 4; areaId++ {
		areaCards = append(areaCards, this.CardMgr.SendCard(2))
	}
	//fmt.Printf("4组牌:%v\n", areaCards)
	this.IdleCard = areaCards[:3]
	this.BankerCard = areaCards[len(areaCards)-1]
	//fmt.Printf("庄牌:%v,闲牌:%v\n", this.BankerCard, this.IdleCard)
	this.GameState = MSG_GAME_INFO_STOP_BET_NOTIFY
	//fmt.Println("wincoins:", this.test())
	//如果没真实玩家下注，有一点几率控制庄家输
	if !this.IsValidBet() && this.IsValidBetByRobot() && BankerLose() && GetCostType() == 1 {
		fmt.Println("进入庄家百分75概率")
		this.GetWinOrLoseResult(areaCards, 0, 3)
		this.AddTimer(gameConfig.Timer.StopBet, gameConfig.Timer.StopBetNum, this.TimerOpen, nil)
		return
	}
	//fmt.Printf("目标库存:%v累积库存:%v\n", G_DbGetGameServerData.GameConfig.GoalStock, CD)
	curCd := CalPkAll(StartControlTime, time.Now().Unix())
	if CD-curCd < 0 && GetCostType() == 1 && this.IsValidBet() { //进入风控换牌
		//fmt.Println("进入风控")
		this.GetWinOrLoseResult(areaCards, 1, 3)
	}
	//logs.Debug("停止下注结束")
	this.AddTimer(gameConfig.Timer.StopBet, gameConfig.Timer.StopBetNum, this.TimerOpen, nil)
}

//递归换牌
func (this *ExtDesk) GetWinOrLoseResult(cards [][]int32, c int, flag int) int64 {
	this.IdleCard = cards[:3]
	this.BankerCard = cards[len(cards)-1]
	if flag <= 0 {
		return 0
	}
	if c == 0 {
		if res := this.testResult(); res > 0 {
			return res
		} else {
			newCards := append([][]int32{}, cards[1:]...)
			newCards = append(newCards, cards[0])
			flag--
			this.GetWinOrLoseResult(newCards, 0, flag)
		}
	} else if c == 1 {
		if res := this.testResult(); res < 0 {
			return res
		} else {
			newCards := append([][]int32{}, cards[1:]...)
			newCards = append(newCards, cards[0])
			flag--
			this.GetWinOrLoseResult(newCards, 1, flag)
		}
	}
	return 0
}

//判断是否有有效下注
func (this *ExtDesk) IsValidBet() bool {
	for _, v := range this.Players {
		if !v.Robot && v.IsBet {
			return true
		}
	}
	return false
}

//判断是否只有机器人下注
func (this *ExtDesk) IsValidBetByRobot() bool {
	for _, v := range this.Players {
		if v.Robot && v.IsBet {
			return true
		}
	}
	return false
}
func (this *ExtDesk) ControlPoker(allCards [][]int32) [][]int32 {

	var maxCards []int32
	for i := 0; i < len(allCards)-1; i++ {
		maxCards = allCards[i]
		for j := i + 1; j < len(allCards); j++ {
			if paigow.CompareCard(allCards[j], allCards[i]) {
				maxCards = allCards[j]
				allCards[j] = allCards[i]
				allCards[i] = maxCards
			}
		}
	}

	return allCards
}

//数组随机
func ListShuffle(list []int) []int {

	var sourceList []int

	MVCard := append([]int{}, list...)
	// 随机打乱牌型
	for i := 0; i < len(list); i++ {
		//打乱
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(MVCard))))
		//添加
		sourceList = append(sourceList, MVCard[randIndex.Int64()])
		// 移除已经添加的牌
		MVCard = ListDelOne(MVCard, int(randIndex.Int64()))
	}
	return sourceList
}

//删除某个元素
func ListDelOne(list []int, num int) []int {

	sourceList := append([]int{}, list...)
	sourceList = append(sourceList[:num], sourceList[num+1:]...)

	return sourceList
}

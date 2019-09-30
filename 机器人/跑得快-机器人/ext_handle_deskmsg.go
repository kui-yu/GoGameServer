package main

import (
	"encoding/json"
	"fmt"
	"logs"
	"math/rand"
	"time"
)

//处理匹配消息返回
func (this *ExtRobotClient) HandleAutoReply(d string) {
	data := GAutoGameReply{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("在监听匹配响应时json转换错误:", err)
		return
	}
	if data.CostType == 1 {
		logs.Debug("现在操作的是金币模式")
	} else {
		logs.Debug("现在操作的是积分模式")
	}
	this.IsDan = -1
	this.Rest()
}

//处理更多匹配信息
func (this *ExtRobotClient) HandleAutoInfo(d string) {
	data := GInfoAutoGameReply{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("处理更多匹配信息时错误")
	}
	for _, v := range data.Seat {
		if this.UserInfo.Uid == v.Uid {
			this.UserInfo.Account = v.Nick
			this.SeatId = v.Cid
		}
	}
}

//处理状态改变通知
func (this *ExtRobotClient) HandleStage(d string) {
	fmt.Println("接收到状态改变信息：")
	data := GStageInfo{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
	}
	if data.Stage == 12 {
		this.Rest()
		this.sendGameAuto()
	}
}

//处理发牌通知
func (this *ExtRobotClient) HandleSendCards(d string) {
	fmt.Println("处理发牌通知:::")
	data := GGameSendCardNotify{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
		logs.Debug("处理发牌通知时发生错误", err)
	}
	for _, v := range data.HandCards {
		this.HandCard = append(this.HandCard, byte(v))
	}
	this.HandCard = Sort(this.HandCard)
	//判断是否是自己出牌 ，如果是自己出牌，那么就出第一张
	this.NextPlayer = data.Cid
	if this.SeatId == this.NextPlayer {
		time.Sleep(time.Second * 5)
		this.OutCard()
	}
	fmt.Println("名称:", this.UserInfo.Account, "的手牌是:", this.HandCard)
}

//处理玩家出牌广播
func (this *ExtRobotClient) HandleOutCards(d string) {
	fmt.Println("接收到玩家出牌")
	data := OutCardsBro{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
	}
	this.LastOutCards = LastOutCards{}
	this.NextPlayer = data.NextCid
	this.LastOutCards.Max = data.Max
	this.LastOutCards.Cid = data.Cid
	for _, v := range data.Cards {
		this.LastOutCards.Cards = append(this.LastOutCards.Cards, byte(v))
	}
	this.LastOutCards.Type = data.Type
	logs.Debug("出牌玩家是否报单：>>>>>>>>>>>>>>>>>>>>>>", data.IsDan)
	if data.IsDan && (this.NextPlayer != this.SeatId) {
		logs.Debug("玩家>>>>>>>>>>>>>>>>:", this.UserInfo.Account, "发现下家报单")
		this.IsDan = data.Cid
	}
	if this.NextPlayer == this.SeatId {
		this.OutCard()
	}
}

//处理过请求
func (this *ExtRobotClient) HandlePass(d string) {
	//处理过通知
	data := PassBro{}
	err := json.Unmarshal([]byte(d), &data)
	if err != nil {
	}
	this.NextPlayer = data.Next
	if this.LastOutCards.Cid == this.NextPlayer {
		this.LastOutCards.Max = 0
	}
	if this.NextPlayer == this.SeatId {
		this.OutCard()
	}
}

//出牌
func (this *ExtRobotClient) OutCard() {
	out := this.TuoGuanOutCards()
	rand.Seed(time.Now().Unix())
	sec := rand.Intn(5) + 1
	fmt.Println("名字为:", this.UserInfo.Account, "的可出牌组为:", out)
	if len(out) > 0 {
		outcard := OutCard{
			Id: MSG_GAME_INFO_OUTCARD,
		}
		for _, v := range out {
			outcard.Cards = append(outcard.Cards, int32(v))
		}
		fmt.Println("机器人装出牌秒数", sec)
		time.Sleep(time.Second * time.Duration(sec))
		this.AddMsgNative(MSG_GAME_INFO_OUTCARD, outcard)
		this.HandCard, _ = VecDelMulti(this.HandCard, out)
	} else {
		fmt.Println("玩家没有合适的牌，等待系统自动帮忙过牌！")
	}
}

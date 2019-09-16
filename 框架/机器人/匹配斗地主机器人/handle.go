package main

import (
	"encoding/json"

	"logs"
	// "logs"
	// "math/rand"
	// "strconv"
	"fmt"
	// "time"
)

//协议监听
func (this *Roboter) InitHandle() {
	//接收大厅登陆应答
	this.Handle[MSG_HALL_LOGIN_REPLY] = this.HandleHallLoginReply
	//接收匹配应答
	this.Handle[MSG_GAME_INFO_AUTO_REPLY] = this.HandleAutoInfoReply
	this.Handle[MSG_GAME_AUTO_REPLY] = this.HandleAutoReply
	//接收游戏房间信息
	this.Handle[MSG_GAME_INFO_ROOM_NOTIFY] = this.HandleRoomNotify
	//接收游戏开始通知
	this.Handle[MSG_GAME_START] = this.HandleGameStart
	//接收发牌通知
	this.Handle[MSG_GAME_INFO_SEND_NOTIFY] = this.HandleSendNotify
	//接收阶段消息
	this.Handle[MSG_GAME_INFO_STAGE] = this.HandleStage
	//接收叫地主广播
	this.Handle[MSG_GAME_INFO_CALL_REPLY] = this.HandlePlayerCallReply
	//接收托管广播
	this.Handle[MSG_GAME_INFO_TUOGUAN_REPLY] = this.HandleTuoGuan

}

// 牌值获取函数
func GetCardColor(card byte) byte {
	return (card & CARD_COLOR) >> 4
}
func GetCardValue(card byte) byte {
	return (card & CARD_VALUE)
}
func GetLogicValue(card byte) byte {
	d := GetCardValue(card)
	if card == 0x41 {
		return 16
	}
	if card == 0x42 {
		return 17
	}
	if d <= 2 {
		return d + 13
	}
	return d
}

//大厅登陆
func (this *Roboter) HandleHallLoginReply(data []byte) {
	logs.Debug("接收到大厅登陆应答")
	var hmloginreply = HMsgHallLoginReply{}
	err := json.Unmarshal(data, &hmloginreply)
	if err != nil {
		logs.Error("大厅应答解析错误")
		fmt.Print(err)
	}
	this.Uid = hmloginreply.Uid
	this.Coins = int32(hmloginreply.Coin)
	this.Account = hmloginreply.Account
	this.Sex = hmloginreply.Sex
	//发送匹配消息
	this.SendToServer(GAutoGame{
		Id:      MSG_GAME_AUTO,
		Account: this.Account,
		Uid:     this.Uid,
		Nick:    this.Name,
		Sex:     this.Sex,
		Coin:    int64(this.Coins),
		Robot:   true,
	})
}

//匹配响应 ，返回额外信息
func (this *Roboter) HandleAutoInfoReply(data []byte) {
	// logs.Debug("接收到匹配响应!")
	autoInfoReply := GInfoAutoGameReply{}
	json.Unmarshal(data, &autoInfoReply)
	for _, v := range autoInfoReply.Seat {
		if this.Uid == v.Uid {
			this.ChairId = v.Cid
			this.Name = v.Nick
		}
	}
}

//匹配响应 ，返回是否匹配成功，是否使用代币
func (this *Roboter) HandleAutoReply(data []byte) {
	autoReply := GAutoGameReply{}
	json.Unmarshal(data, &autoReply)
	if autoReply.Result == 0 {
		// logs.Debug("匹配成功!")
	}
	if autoReply.CostType == 2 {
		logs.Debug("体验场")
	} else {
		logs.Debug("普通场")
	}
}

//房间信息通知
func (this *Roboter) HandleRoomNotify(data []byte) {
	logs.Debug("接收到房间信息")
	infoNotify := GGameInfoNotify{}
	json.Unmarshal(data, &infoNotify)
	fmt.Println(infoNotify)
}

//游戏开始通知
func (this *Roboter) HandleGameStart(data []byte) {
	logs.Debug("游戏开始！")
}

//发牌通知
func (this *Roboter) HandleSendNotify(data []byte) {
	logs.Debug("接收到发牌通知")
	sd := GGameSendCardNotify{}
	json.Unmarshal(data, &sd)
	this.HandCard = sd.HandsCards
	var CardsValue []byte
	for _, v := range this.HandCard {
		CardsValue = append(CardsValue, GetLogicValue(v))
	}
	fmt.Println(this.ChairId, "座位玩家的手牌:", CardsValue)
	//开始根据服务端返回过来的Cid 判断是由谁第一个叫地主
	if this.ChairId == sd.Cid {
		fmt.Println("由", this.ChairId, "座位玩家第一个叫地主！")
		this.firstCall = this.ChairId
	}
}

//阶段通知处理
func (this *Roboter) HandleStage(data []byte) {
	logs.Debug("接收到阶段通知")
	stageInfo := GStageInfo{}
	json.Unmarshal(data, &stageInfo)
	//判断现在到底是什么阶段
	switch stageInfo.Stage {
	case GAME_STATUS_CALL:
		logs.Debug("现在到了叫分阶段!")
		this.PlayerCall() //玩家叫分
		break
	case GAME_STATUS_PLAY:
		logs.Debug("现在到了游戏操作阶段!")
		this.PlayerPlay()
		break
	}
}

//阶段执行方法--叫地主
func (this *Roboter) PlayerCall() {
	if this.ChairId == this.firstCall {
		this.SendToServer(GCallMsg{
			Id:    MSG_GAME_INFO_CALL,
			Coins: 3,
		})
	} else {
		return
	}
}

//阶段执行方法--游戏操作
func (this *Roboter) PlayerPlay() {
	//为了方便 ，我们执行托管
	this.SendToServer(GTuoGuan{
		Id:  MSG_GAME_INFO_TUOGUAN,
		Ctl: 1,
	})
	fmt.Println(this.ChairId, "托管")
}

//叫地主应答处理
func (this *Roboter) HandlePlayerCallReply(data []byte) {
	logs.Debug("接收到叫地主广播")
	callMsgReply := GCallMsgReply{}
	json.Unmarshal(data, &callMsgReply)
	fmt.Println(callMsgReply.Cid, "座位玩家叫了", callMsgReply.Coins, "分", "是否结束:", callMsgReply.End)
}

//托管应答处理
func (this *Roboter) HandleTuoGuan(data []byte) {
	tuoGuanReply := GTuoGuanReply{}
	json.Unmarshal(data, &tuoGuanReply)
}

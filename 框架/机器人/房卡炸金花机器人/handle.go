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

func (this *Roboter) InitHandle() {
	//大厅登陆
	this.Handle[MSG_HALL_LOGIN_REPLY] = this.HandleHallLoginReply
	//房间创建回复消息
	this.Handle[MSG_GAME_FK_CREATEDESK_REPLY] = this.HandleFk
	//房间加入消息
	this.Handle[MSG_GAME_FK_JOIN_REPLY] = this.HandleAutoReply
	//接收赋值chairid
	this.Handle[MSG_GAME_INFO_AUTO_REPLY] = this.HandleGetChairId
	//准备应答
	this.Handle[MSG_GAME_INFO_READY_REPLY] = this.handleReadReplay
	//游戏开始应答
	this.Handle[MSG_GAME_INFO_START] = this.HandleGameStart
	//发牌
	this.Handle[MSG_GAME_INFO_SEND_NOTIFY] = this.handleSendCard
	//游戏叫分阶段
	this.Handle[MSG_GAME_INFO_STAGE] = this.HandleGameStage
	//游戏 抢地主应答
	this.Handle[MSG_GAME_INFO_GETMSG_REPLY] = this.HanleGameGetMsgReply
	//游戏顶庄阶段
	this.Handle[MSG_GAME_INFO_BANKER_NOTIFY] = this.HandleGameBanker

}

var hasRoom chan bool // 无缓冲通道
var FkNo string
var CallNum int
var GetMsgCid int32 //记录第一个抢地主人的ID

func init() {
	hasRoom = make(chan bool, 1)
	hasRoom <- false
}

//大厅登陆
func (this *Roboter) HandleHallLoginReply(data []byte) {
	d := HMsgHallLoginReply{}
	json.Unmarshal(data, &d)

	this.Coins = d.Coin
	this.Uid = d.Uid

	room := <-hasRoom
	fmt.Println("hasroom：", room)

	if !room { //如果该房间不存在的话，那么我们需要创建该房间
		logs.Debug("正在创建房间", this.Uid)
		fkinfo := GATableConfig{
			GameModle:  1,
			PlayerNum:  3,
			MatchCount: 5,
			GameType:   1,
			PayType:    1,
			BaseScore:  100,
			CallType:   1,
			Boom:       3,
			CanSelect:  []int{1, 2, 3},
		}
		fkinfos, _ := json.Marshal(fkinfo)
		info := HMsgHallCreateFkRoom{ // 300012 -> 400018
			Id:        MSG_HALL_CREATE_FKROOM,
			GameType:  14,
			RoomType:  2,
			GradeType: 5,
			FkInfo:    string(fkinfos),
		}
		this.SendToServer(info)

	} else { //如果存在该房间，那么就可以直接进入了。

		info1 := GFkJoinToGame{
			Id:      MSG_GAME_FK_JOIN,
			Account: this.Account,
			Uid:     this.Uid,
			Nick:    "noname",
			Sex:     this.Sex,
			Head:    "nohead",
			Lv:      0,
			Coin:    int64(this.Coins),
			Token:   "notoken",
			Robot:   true,
			FkNo:    FkNo, //房间号
		}
		logs.Debug("房间号", FkNo)
		this.SendToServer(info1)
	}
}

func (this *Roboter) HandleFk(data []byte) {
	d := GFkCreateDeskReply{}
	json.Unmarshal(data, &d)
	fmt.Println("收到消息：", d, this.Coins)

	info := GFkJoinToGame{
		Id:      MSG_GAME_FK_JOIN,
		Account: this.Account,
		Uid:     this.Uid,
		Nick:    "noname",
		Sex:     this.Sex,
		Head:    "nohead",
		Lv:      0,
		Coin:    int64(this.Coins),
		Token:   "notoken",
		Robot:   true,
		FkNo:    d.FkNo, //房间号
	}
	FkNo = d.FkNo
	this.SendToServer(info)
	logs.Debug("房间主人进入房间")
	fmt.Println("加入房间：", d.FkNo)
}

func (this *Roboter) HandleGetChairId(data []byte) {
	d := GInfoAutoGameReply{}
	json.Unmarshal(data, &d)

	fmt.Println(d)
	for _, v := range d.Seat {
		if this.Uid == v.Uid {
			this.ChairId = v.Cid
		}
	}
	fmt.Println("赋值chairID：", this.ChairId)
}

func (this *Roboter) HandleAutoReply(data []byte) {
	d := GFkJoinReply{}
	json.Unmarshal(data, &d)
	if d.Result == 0 {
		if hasRoom != nil {
			hasRoom <- true
		} else {
			fmt.Println("发送消息到通道失败")
		}
		fmt.Println("加入房间：", d)
	} else {
		fmt.Println("房间加入失败应答：", d)
	}
}

//匹配应答
func (this *Roboter) HandleGameAutoReply(data []byte) {
	d := GSInfoAutoGame{}
	json.Unmarshal(data, &d)

	for _, v := range d.Seat {
		if v.Uid == this.Uid {
			this.ChairId = v.Cid
			break
		}
	}
	fmt.Println(this.ChairId, "座位号")
}

var getmsg int

//阶段消息
func (this *Roboter) HandleGameStage(data []byte) {
	d := GSStageInfo{}
	json.Unmarshal(data, &d)
	fmt.Println("阶段消息：", d.Stage)
	if d.Stage == GAME_STATUS_CALL { //判断该状态是否为叫分状态
		fmt.Println("叫分阶段可以叫分")
		CallNum++
		//叫分
		call := GCallMsg{
			Id:    MSG_GAME_INFO_CALL,
			Coins: int32(CallNum),
		}
		fmt.Println("玩家", this.ChairId, "叫", CallNum, "分")
		//发送信息
		this.SendToServer(call)
	} else if d.Stage == GAME_STATUS_PLAY { //判断该状态是否为出牌状态
		fmt.Println("出牌阶段，可以出牌")
	} else if d.Stage == GAME_STATUS_READ {
		fmt.Println("准备阶段到了")
		msg := GAPlayerReady{
			Id:      MSG_GAME_INFO_READY,
			IsReady: 1,
		}
		this.SendToServer(msg)
		fmt.Println("玩家准备成功!", msg)
	} else if d.Stage == GAME_STATUS_GETMSG {
		fmt.Println("进入抢地主阶段")
		if this.ChairId != GetMsgCid {

		} else {

			this.SendToServer(GGetMsg{
				Id:     MSG_GAME_INFO_GETMSG,
				GetMsg: 1,
			})
		}
	}
}

//游戏开始
func (this *Roboter) HandleGameStart(data []byte) {
	d := GGameStartNotify{}
	json.Unmarshal(data, &d)
	fmt.Println("第", d.Round, "局游戏")
}

//发牌
func (this *Roboter) handleSendCard(data []byte) {
	fmt.Println("进入发牌阶段")
	//接收发牌发过来的值
	d := GGameSendCardNotify{}
	err := json.Unmarshal(data, &d)
	fmt.Println("第一个抢地主人的ChairID:", d.Cid)
	GetMsgCid = d.Cid
	if err != nil {
		logs.Debug("发牌信息解析错误")
	}
}

//定庄通知
func (this *Roboter) HandleGameBanker(data []byte) {
	fmt.Println("定庄完毕!")
	d := GBankerNotify{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		fmt.Println("定庄解析错误!")
		return
	}
	fmt.Println("底牌:", d.DiPai)
}

//准备应答
func (this *Roboter) handleReadReplay(data []byte) {
	d := GSPlayerReady{}
	json.Unmarshal(data, &d)
	if d.IsReady == 1 {
		fmt.Println("该玩家已经准备，座位号为", d.ChairId)
	} else {
		fmt.Println("该玩家还没准备")
	}
}

var count int

//抢地主应答
func (this *Roboter) HanleGameGetMsgReply(data []byte) {
	fmt.Println("接收到抢地主应答")
	d := GGetMsgReply{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		logs.Debug("抢地主应答解析错误！", err)
	}
	fmt.Println("服务器让下一位叫庄人的椅子Id", d.Cid)
	fmt.Println("我的椅子Id", this.ChairId)

	if this.ChairId == d.Cid {
		this.SendToServer(GGetMsg{
			Id:     MSG_GAME_INFO_GETMSG,
			GetMsg: 2,
		})
		fmt.Println("再次发送抢地主:", this.ChairId)
	}
	count++
}

package main

import (
	// "crypto/rand"
	"encoding/json"
	"logs"

	"math/rand"
	// "strconv"
	"time"
)

func (this *Roboter) InitHandle() {
	//大厅登陆
	this.Handle[MSG_HALL_LOGIN_REPLY] = this.HandleHallLoginReply
	//匹配应答
	this.Handle[MSG_GAME_INFO_AUTO_REPLY] = this.HandleGameAutoReply
	//
	this.Handle[MSG_GAME_INFO_STAGE] = this.HandleGameStage
	//
	this.Handle[MSG_GAME_INFO_SEND_NOTIFY] = this.HandleSenddCard
	//
	this.Handle[MSG_GAME_INFO_HUANPAIOVER_NOTIFY] = this.HuanPaiOver
	//
	this.Handle[MSG_GAME_INFO_HAVEACTION_NOTIFY] = this.HaveAction
	//
	this.Handle[MSG_GAME_INFO_SENDCARD_NOTIFY] = this.SendCard
	//
	this.Handle[MSG_GAME_INFO_ACTION_NOTIFY] = this.ActionNotify
}

//状态消息
func (this *Roboter) HandleGameStage(data []byte) {
	d := GStageInfo{}
	json.Unmarshal(data, &d)
	logs.Debug("阶段", d.Stage)
	if d.Stage == GAME_STATE_CHANGECARD {
		time.Sleep(time.Second)
		huan := this.HandCard[len(this.HandCard)-3:]
		this.HandCard = append([]int{}, this.HandCard[:len(this.HandCard)-3]...)
		logs.Debug("换牌", huan, this.HandCard)
		sd := GHuanPai{
			Id:    MSG_GAME_INFO_HUANPAI,
			Cards: huan,
		}
		this.SendToServer(sd)
	} else if d.Stage == GAME_STATE_DINGQUE {
		time.Sleep(time.Second)
		co := rand.Intn(10) % 3
		this.SendToServer(GDingQue{
			Id:    MSG_GAME_INFO_DINGQUE,
			Color: co,
		})
		this.QueColor = co
		logs.Debug("定缺")
	} else if d.Stage == GAME_STATE_PLAY { //游戏开始，庄家出牌或者动作
		//定时器延迟3秒出牌，因为有动作要先动作,动作后，清除定时器
		if this.SeatId == this.CurId {
			logs.Debug("添加出牌定时器")
			this.AddTimer(1, 3, this.OutCard, nil)
		}
	}
}

func (this *Roboter) OutCard(d interface{}) {
	if len(this.HandCard) == 0 {
		return
	}
	c := this.HandCard[len(this.HandCard)-1]
	this.SendToServer(GOutCard{
		Id:   MSG_GAME_INFO_OUTCARD,
		Card: c,
	})
	logs.Debug("玩家出牌", this.SeatId)
	this.HandCard = this.HandCard[:len(this.HandCard)-1]
}

//大厅登陆
func (this *Roboter) HandleHallLoginReply(data []byte) {
	d := HMsgHallLoginReply{}
	json.Unmarshal(data, &d)
	// logs.Debug("登录大厅收到的消息", d)
	this.Coins = d.Coin
	this.Uid = d.Uid

	if d.GameId == 0 {
		sd := GAutoGame2{
			Id: MSG_GAME_AUTO,
		}
		logs.Debug("发送匹配信息")
		this.SendToServer(sd)
	} else {
		sd := GReconnect{
			Id: MSG_GAME_RECONNECT,
		}
		this.SendToServer(sd)
	}
}

//匹配应答
func (this *Roboter) HandleGameAutoReply(data []byte) {
	d := GInfoAutoGameReply{}
	json.Unmarshal(data, &d)
	for _, v := range d.Seat {
		if v.Uid == this.Uid {
			this.SeatId = v.Cid
		}
	}
}

func (this *Roboter) HandleSenddCard(data []byte) {
	d := GGameSendCardNotify{}
	json.Unmarshal(data, &d)
	//
	this.Banker = d.Banker
	this.CurId = int32(d.Banker)
	this.HandCard = d.HandsCards
	logs.Debug("手牌：", d.HandsCards, d.Banker)
}

func (this *Roboter) HuanPaiOver(data []byte) {
	d := GHuanPaiOver{}
	json.Unmarshal(data, &d)
	//
	this.HandCard = append(this.HandCard, d.Cards...)
	logs.Debug("当前手牌:", this.HandCard)
}

func (this *Roboter) HaveAction(data []byte) {
	d := GHaveActionNotify{}
	json.Unmarshal(data, &d)
	//
	logs.Debug("当前有动作", d.Data, this.SeatId)
	ac := d.Data[0]
	time.Sleep(time.Second)
	this.SendToServer(GAction{
		Id:    MSG_GAME_INFO_ACTION,
		Style: ac.Style,
		Card:  ac.Card,
	})
	logs.Debug("操作动作", d.Data[0], this.SeatId)
	this.ClearTimer()
}

func (this *Roboter) SendCard(data []byte) {
	d := GSendCardNofify{}
	json.Unmarshal(data, &d)
	//
	if d.Cid == this.SeatId {
		this.HandCard = append(this.HandCard, d.Card)
		SortByQue(&this.HandCard, this.QueColor)
		logs.Debug("发牌，3秒后没动作才出牌", d.Cid, d.Card)
		this.CurId = d.Cid
		this.AddTimer(1, 3, this.OutCard, nil)
	}
}

//游戏结束
func (this *Roboter) GameEndNotify(data []byte) {
	d := GSSettleInfos{}
	json.Unmarshal(data, &d)
	// logs.Debug("收到结束的消息", d)

	//初始化数据
	this.HandCard = []int{}
	this.SeatId = -2
	this.CurId = -2
	//
	// time.Sleep(time.Second * 3)

	sd := GAutoGame{
		Id:      MSG_GAME_AUTO,
		Account: this.Account,
		Uid:     this.Uid,
		Nick:    this.Name,
		Sex:     this.Sex,
		Head:    "",
		Lv:      1,
		Coin:    this.Coins,
	}
	this.SendToServer(sd)
}

func (this *Roboter) ActionNotify(data []byte) {
	d := GActionDoNotify{}
	json.Unmarshal(data, &d)
	this.ClearTimer()
	//
	if d.Cid == int(this.SeatId) {
		logs.Debug("玩家：", d.Cid, "操作了动作:", d.ActionType, d.Cards)
		//删除手牌，然后判断是否需要出牌
		DeleteCard(&this.HandCard, d.Cards)
		if d.ActionType == 5 {
			this.AddTimer(1, 3, this.OutCard, nil)
		}
	}
}

////////////////////////////////////////////////
func SortByQue(cs *[]int, color int) {
	cs1 := []int{}
	cs2 := []int{}
	for _, v := range *cs {
		if GetCardColor(byte(v)) == byte(color) {
			cs1 = append(cs1, v)
		} else {
			cs2 = append(cs2, v)
		}
	}
	Sort(cs1)
	Sort(cs2)
	*cs = append([]int{}, cs2...)
	*cs = append(*cs, cs1...)

}

func GetCardColor(card byte) byte {
	return (card & 0xf0) >> 4
}

func GetCardValue(card byte) byte {
	return card & 0x0f
}

func Sort(cs []int) {
	for i := 0; i < len(cs)-1; i++ {
		for j := i + 1; j < len(cs); j++ {
			if cs[i] > cs[j] {
				ad := cs[j]
				cs[j] = cs[i]
				cs[i] = ad
			}
		}
	}
}

func DeleteCard(ShouPai *[]int, Cards []int) bool {
	nh := append([]int{}, (*ShouPai)...)
	delcount := 0
	for _, c := range Cards {
		for i, s := range nh {
			if c == s {
				delcount++
				nh = append(nh[:i], nh[i+1:]...)
				break
			}
		}
	}
	if delcount == len(Cards) {
		*ShouPai = nh
		return true
	}
	return false
}

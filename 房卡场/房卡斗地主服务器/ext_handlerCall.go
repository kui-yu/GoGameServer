package main

import (
	"encoding/json"
	"fmt"
	"logs"
)

var reSendCount int   //计算重新发牌次数
var netxCallCount int //计算下一轮叫地主的次数
var callPlayer int32  //下一个叫分的玩家椅子号

//叫地主方法
func (this *ExtDesk) HandleGameCall(p *ExtPlayer, d *DkInMsg) {
	//判断是否叫分阶段
	if this.GameState != GAME_STATUS_CALL {
		logs.Error("叫分游戏状态错误:", this.GameState, GAME_STATUS_CALL)
		return
	}

	logs.Debug("玩家叫分")
	data := GCallMsg{}
	json.Unmarshal([]byte(d.Data), &data)

	//判断是否玩家叫分
	if p.ChairId != this.CurCid || data.Coins == 0 {
		logs.Error("不是当前叫分用户或叫分为0:", p.ChairId, this.CurCid, data.Coins)
		return
	}
	//叫分不能小于最大叫分
	if data.Coins > 0 && data.Coins <= this.CallFen {
		logs.Error("叫分小于当前叫分:", data.Coins, this.CallFen)
		return
	}
	p.CFen = data.Coins
	logs.Debug("玩家叫分：", p.CFen)
	if p.CFen > 0 {
		this.CallFen = p.CFen
		this.Banker = p.ChairId
	}
	//删除定时器
	this.DelTimer(TIMER_CALL)
	//是否还有下一个
	isover := false
	if p.CFen < 3 {
		nextid := (p.ChairId + 1) % int32(len(this.Players))
		nextPlayer := this.Players[nextid]
		if nextPlayer.CFen != 0 { //如果下一位玩家已经叫过分了，证明叫地主已经轮流一圈了，所以地主直接归该玩家
			isover = true
		} else {
			this.CurCid = nextid
			if p.TuoGuan {
				this.AddTimer(TIMER_CALL, 1, this.TimerCall, nil)
			} else {
				this.AddTimer(TIMER_CALL, TIMER_CALL_NUM, this.TimerCall, nil)
			}

		}
	} else {
		isover = true
	}
	//
	logs.Debug("下一个椅子id：", this.CurCid)
	var dd int32 = -1
	if data.Coins > 0 {
		dd = data.Coins
	}
	re := GCallMsgReply{
		Id:      MSG_GAME_INFO_CALL_REPLY,
		Cid:     p.ChairId,
		Coins:   dd,
		NextCid: (this.CurCid) % int32(len(this.Players)),
	}

	if isover {
		re.End = 2
	} else {
		re.End = 1
	}
	this.BroadcastAll(MSG_GAME_INFO_CALL_REPLY, &re)
	fmt.Println("广播玩家叫分：", re.Coins)
	logs.Debug("是否结束:", re.End)
	if !isover {
		logs.Debug("走不到下面")
		return
	}
	//判断所有人是不是都不叫。如果都不叫就重新进去发牌阶段
	allnocall := true
	for _, v := range this.Players {
		if v.CFen != -1 {
			allnocall = false
			break
		}
	}
	if allnocall {

		re.End = 2
		this.BroadcastAll(MSG_GAME_INFO_CALL_REPLY, re)
		logs.Debug("是否结束:", re.End)
		for _, v := range this.Players {
			v.CFen = 0
		}
		if this.CallTimes >= 3 {
			this.GameOverByNoCall(p)
		} else {
			//重新发牌
			logs.Debug("检测到3个玩家未叫地主")
			for _, v := range this.Players {
				v.GetMSG = 0
			}
			this.CallTimes++
			fmt.Println("叫分CallTime:", this.CallTimes)
			if this.CallTimes >= 3 {
				this.GameState = GAME_STATUS_BALANCE
				this.BroadStageTime(0)
				this.GameOverByNoCall(p)
				return
			} else {
				logs.Debug("发送状态：", GAME_STATUS_START)
				this.GameState = GAME_STATUS_START
				this.BroadStageTime(TIMER_START_NUM)
				//发送游戏开始通知 并附带额外信息
				this.BroadcastAll(MSG_GAME_INFO_START, &GGameStartNotify{
					Id:    MSG_GAME_INFO_START,
					Round: this.Round,
				})
				//游戏开始，进入发牌阶段
				logs.Debug("并且进入发牌阶段")
				this.TimerSendCard(nil)
				return
			}
		}
		return
	}
	//设置倍数
	for _, v := range this.Players {
		if v.ChairId == this.Banker {
			v.Double = this.CallFen * 2
		} else {
			v.Double = this.CallFen
		}
	}
	//
	banker := this.Players[this.Banker]
	banker.HandCard = append(banker.HandCard, this.DiPai...)
	this.DiPaiDoulbe = int32(this.CalDiPaiDouble(this.DiPai))
	//都叫分过了。定庄
	this.GameState = GAME_STATUS_PLAY
	//出牌阶段消息
	this.BroadStageTime(TIMER_OUTCARD_NUM)

	this.CurCid = banker.ChairId
	notify := GBankerNotify{
		Id:        MSG_GAME_INFO_BANKER_NOTIFY,
		Banker:    banker.ChairId,
		Double:    this.DiPaiDoulbe,
		AllDouble: banker.CFen,
	}
	for _, v := range this.DiPai {
		notify.DiPai = append(notify.DiPai, int(v))
	}
	fmt.Println("ALLDOUBLE:", notify.AllDouble)
	notify.AllDouble = notify.AllDouble * this.DiPaiDoulbe
	this.Double = notify.AllDouble
	this.BroadcastAll(MSG_GAME_INFO_BANKER_NOTIFY, &notify)

	//添加定时器，进入出牌阶段
	nextplayer := this.Players[this.CurCid]
	if nextplayer.TuoGuan {
		this.AddTimer(TIMER_OUTCARD, 1, this.TuoGuanOut, nil)
	} else {
		this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCard, nil)
	}
}

//叫地主阶段超时处理
func (this *ExtDesk) TimerCall(d interface{}) {
	logs.Debug("进入叫分")
	data := GCallMsg{
		Coins: -1,
	}
	p := this.Players[this.CurCid]
	dv, _ := json.Marshal(data)
	this.HandleGameCall(p, &DkInMsg{
		Uid:  p.Uid,
		Data: string(dv),
	})
}

//抢地主方法
func (this *ExtDesk) HandleGameGetMsg(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("进入抢地主判断")
	data := GGetMsg{}
	json.Unmarshal([]byte(d.Data), &data)
	this.CurCid = (this.CurCid + 1) % int32(len(this.Players))
	gmr := GGetMsgReply{
		Id:      MSG_GAME_INFO_GETMSG_REPLY,
		Cid:     p.ChairId,
		IsGet:   data.GetMsg,
		NextCid: this.CurCid,
	}
	this.DelTimer(TIMER_GETGMS)
	p.GetMSG = data.GetMsg
	fmt.Println("玩家抢叫：", p.GetMSG)
	if data.GetMsg == 3 {
		gmr.CallOrGet = 1
		this.CallOrGet = 1
	} else {
		gmr.CallOrGet = 2
		this.CallOrGet = 2
	}
	var cheng int32
	if p.GetMSG == 1 {
		p.Double = 1
		cheng = 1
	}
	if p.GetMSG == 2 {
		p.Double = 2
		cheng = 1
	}
	if cheng*p.Double != 0 {
		this.CallFen = cheng * p.Double
	}
	//如果3个玩家都不叫，那么会重新进入游戏开始阶段
	if p.GetMSG == 3 {
		var pd bool = true
		for _, v := range this.Players {
			if v.GetMSG == 0 {
				pd = false
			}
		}
		if pd {
			//重新发牌
			logs.Debug("检测到3个玩家未叫地主")
			for _, v := range this.Players {
				v.GetMSG = 0
			}
			this.CallTimes++
			gmr.End = 2
			this.BroadcastAll(MSG_GAME_INFO_GETMSG_REPLY, gmr)
			if this.CallTimes >= 3 {
				this.GameState = GAME_STATUS_BALANCE
				this.BroadStageTime(0)
				this.GameOverByNoCall(p)
				return
			} else {
				logs.Debug("发送状态：", GAME_STATUS_START)
				this.GameState = GAME_STATUS_START
				this.BroadStageTime(TIMER_START_NUM)
				//发送游戏开始通知 并附带额外信息
				this.BroadcastAll(MSG_GAME_INFO_START, &GGameStartNotify{
					Id:    MSG_GAME_INFO_START,
					Round: this.Round,
				})
				//游戏开始，进入发牌阶段
				logs.Debug("并且进入发牌阶段")
				this.TimerSendCard(nil)
				return
			}

		}
	}
	//选则庄位
	for _, v := range this.Players {
		if v.GetMSG == 0 {
			fmt.Println(v.Nick, "玩家抢地主：", v.GetMSG)
			gmr.End = 1
		}
	}
	var ban int32 //如果产生庄，记录庄的座位号
	var d1 int32 = 1
	for _, v := range this.Players {
		d1 = d1 * v.Double
	}
	fmt.Println("是否结束：", gmr.End)
	gmr.Double = d1

	if data.GetMsg == 3 {
		gmr.Double = 0
	}
	if gmr.End == 1 {
		this.BroadcastAll(MSG_GAME_INFO_GETMSG_REPLY, gmr)
		this.AddTimer(TIMER_GETGMS, TIMER_GETGMS_NUM, this.TimerGet, nil)
		fmt.Print("倍数：", gmr.Double)
		return
	} else {
		gmr.End = 2
		fmt.Println("zhe!!!!!!!!!!!!!!!!!!是否结束：", gmr.End)
		if p.GetMSG == 1 || p.GetMSG == 2 {
			ban = p.ChairId
		} else if p.GetMSG == 4 {
			chaird := ((p.ChairId+1)%int32(len(this.Players)) + 1) % int32(len(this.Players))
			if this.Players[chaird].GetMSG == 1 || this.Players[chaird].GetMSG == 2 {
				ban = chaird
			} else {
				ban = this.CurCid
			}
		}
		this.BroadcastAll(MSG_GAME_INFO_GETMSG_REPLY, gmr)
		fmt.Print("倍数：", gmr.Double)
	}
	this.Banker = ban
	//定庄结束
	logs.Debug("定庄结束 ，庄位", this.Banker)
	netxCallCount = 0
	banker := this.Players[this.Banker]
	banker.HandCard = append(banker.HandCard, this.DiPai...)
	this.DiPaiDoulbe = int32(this.CalDiPaiDouble(this.DiPai))
	this.GameState = GAME_STATUS_PLAY
	logs.Debug("目前状态", this.GameState)
	//出牌阶段消息
	this.BroadStageTime(TIMER_OUTCARD_NUM)
	this.CurCid = banker.ChairId
	notify := GBankerNotify{
		Id:        MSG_GAME_INFO_BANKER_NOTIFY,
		Banker:    banker.ChairId,
		Double:    this.DiPaiDoulbe,
		AllDouble: 1,
	}
	for _, v := range this.DiPai {
		notify.DiPai = append(notify.DiPai, int(v))
	}
	var dou int32 = 1
	for _, v := range this.Players {
		dou = dou * v.Double
	}
	fmt.Println("doddddddddddddddddddddddddddddddddddd:", dou)
	notify.AllDouble = dou
	notify.AllDouble = notify.AllDouble * this.DiPaiDoulbe
	this.Double = notify.AllDouble
	fmt.Println("底牌倍数:", this.DiPaiDoulbe)
	this.BroadcastAll(MSG_GAME_INFO_BANKER_NOTIFY, &notify)
	fmt.Println("定庄倍数:", notify.AllDouble)
	this.Players[this.Banker].Double = this.Players[this.Banker].Double * 4
	for _, v := range this.Players {
		if v.Uid != banker.Uid {
			v.Double = banker.Double / 2
		}
	}
	//添加定时器，进入出牌阶段
	nextplayer := this.Players[this.CurCid]
	if nextplayer.TuoGuan {
		this.AddTimer(TIMER_OUTCARD, 1, this.TuoGuanOut, nil)
	} else {
		this.AddTimer(TIMER_OUTCARD, TIMER_OUTCARD_NUM, this.TimerOutCard, nil)
	}
	fmt.Println("定庄彻底结束")
}

//抢地主阶段超时处理
func (this *ExtDesk) TimerGet(d interface{}) {
	logs.Debug("定时触发？？？？？？？？11111111111111111111111111111111111111111111111111111111111111111111")
	data := GGetMsg{}
	if this.CallOrGet == 1 {
		data.GetMsg = 3
	} else {
		data.GetMsg = 4
	}
	p := this.Players[this.CurCid]
	dv, _ := json.Marshal(data)
	this.HandleGameGetMsg(p, &DkInMsg{
		Uid:  p.Uid,
		Data: string(dv),
	})
}

//计算底牌倍数，暂时不用
func (this *ExtDesk) CalDiPaiDouble(cards []byte) int {
	vdipai := Sort(cards)
	if vdipai[0] == 0x42 && vdipai[1] == 0x41 {
		return 4
	} else if vdipai[0] == 0x42 {
		return 2
	} else if vdipai[0] == 0x41 {
		return 2
	} else if GetLogicValue(vdipai[0]) == GetLogicValue(vdipai[1]) &&
		GetLogicValue(vdipai[0]) == GetLogicValue(vdipai[2]) {
		return 4
	} else if GetCardColor(vdipai[0]) == GetCardColor(vdipai[1]) &&
		GetCardColor(vdipai[0]) == GetCardColor(vdipai[2]) {
		if this.IsShunZiType(vdipai) {
			return 4
		} else {
			return 3
		}
	} else if this.IsShunZiType(vdipai) {
		return 3
	}
	return 1
}

func (this *ExtDesk) IsShunZiType(cards []byte) bool {
	for i := 0; i < len(cards)-1; i++ {
		if GetLogicValue(cards[i]) != GetLogicValue(cards[i+1])+1 {
			return false
		}
	}
	return true
}

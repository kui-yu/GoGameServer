package main

import (
	. "MaJiangTool"
	"encoding/json"
	"logs"
	"math/rand"
)

func (this *ExtDesk) HandleHuanPai(p *ExtPlayer, d *DkInMsg) {
	//判断是否叫分阶段
	if this.GameState != GAME_STATE_CHANGECARD {
		logs.Error("还没到换牌阶段:", this.GameState)
		return
	}

	data := GHuanPai{}
	json.Unmarshal([]byte(d.Data), &data)
	//
	if len(p.HuanPais) != 0 || len(data.Cards) != 3 {
		return
	}
	for _, v := range data.Cards {
		p.HuanPais = append(p.HuanPais, byte(v))
	}
	//删除手牌
	hand := append([]byte{}, p.HandCard...)
	hand, ok := VecDelMulti(hand, p.HuanPais)
	if !ok {
		return
	}
	p.HandCard = hand
	//广播谁提交了交换牌
	this.BroadcastAll(MSG_GAME_INFO_HUANPAI_NOTIFY, &GHuanPaiNotify{
		Id:  MSG_GAME_INFO_HUANPAI_NOTIFY,
		Cid: int(p.ChairId),
	})
	//判断都否都提交了换牌
	for _, p := range this.Players {
		if len(p.HuanPais) == 0 {
			return
		}
	}
	//////////////////////////////////////////////
	//广播进入换牌结束阶段
	this.GameState = GAME_STATE_CHANGECARD_OVER
	this.BroadStageTime(TIMER_HUANPAIOVER_NUM)

	//交换牌
	changeType := rand.Intn(10) % 2
	if changeType == 0 { //顺时针
		vcards := append([]byte{}, this.Players[0].HuanPais...)
		for i := 0; i < len(this.Players); i++ {
			if i == len(this.Players)-1 {
				this.Players[i].HuanPais = vcards
			} else {
				this.Players[i].HuanPais = append([]byte{}, this.Players[i+1].HuanPais...)
			}
		}
	} else if changeType == 1 { //逆时针
		vcards := append([]byte{}, this.Players[len(this.Players)-1].HuanPais...)
		for i := len(this.Players) - 1; i >= 0; i-- {
			if i == 0 {
				this.Players[i].HuanPais = vcards
			} else {
				this.Players[i].HuanPais = append([]byte{}, this.Players[i-1].HuanPais...)
			}
		}
	}
	//对手牌重新排序
	for _, v := range this.Players {
		v.HandCard = append(v.HandCard, v.HuanPais...)
		logs.Debug("当前手牌1:", v.HandCard, v.HuanPais)
		Sort(v.HandCard)
	}
	//删除定时器
	this.DelTimer(TIMER_HUANPAI)
	//广播交换方式
	for _, v := range this.Players {
		huanover := GHuanPaiOver{
			Id:    MSG_GAME_INFO_HUANPAIOVER_NOTIFY,
			Style: changeType,
		}
		for _, h := range v.HuanPais {
			huanover.Cards = append(huanover.Cards, int(h))
		}
		v.SendNativeMsg(MSG_GAME_INFO_HUANPAIOVER_NOTIFY, &huanover)
	}
	// this.BroadcastAll(MSG_GAME_INFO_HUANPAIOVER_NOTIFY, &GHuanPaiOver{
	// 	Id:    MSG_GAME_INFO_HUANPAIOVER_NOTIFY,
	// 	Style: changeType,
	// })
	//广播几秒后进入定缺状态
	this.AddTimer(TIMER_HUANPAIOVER, TIMER_HUANPAIOVER_NUM, this.TimerDingQue, nil)
}

//进入和广播换牌阶段
func (this *ExtDesk) TimerHuanPai(d interface{}) {
	//进入抢庄
	this.GameState = GAME_STATE_CHANGECARD
	this.BroadStageTime(TIMER_HUANPAI_NUM)
	//发牌后进入叫分阶段，开启叫分阶段的定时器
	this.AddTimer(TIMER_HUANPAI, TIMER_HUANPAI_NUM, this.TimerHuanPaiDo, nil)
}

//换牌超时定时器
func (this *ExtDesk) TimerHuanPaiDo(d interface{}) {
	for cid := 0; cid < len(this.Players); cid++ {
		if len(this.Players[cid].HuanPais) == 0 {
			//获取花色对应的牌
			wan := []byte{}
			tiao := []byte{}
			tong := []byte{}
			for _, c := range this.Players[cid].HandCard {
				color := GetCardColor(c)
				if color == CARD_COLOR_Wan {
					wan = append(wan, c)
				} else if color == CARD_COLOR_Tiao {
					tiao = append(tiao, c)
				} else if color == CARD_COLOR_Tiao {
					tong = append(tong, c)
				}
			}
			//
			msg := GHuanPai{
				Id: MSG_GAME_INFO_HUANPAI,
			}
			if len(wan) >= 3 {
				for n := 0; n < 3; n++ {
					msg.Cards = append(msg.Cards, int(wan[n]))
				}
			} else if len(tiao) >= 3 {
				for n := 0; n < 3; n++ {
					msg.Cards = append(msg.Cards, int(tiao[n]))
				}
			} else {
				for n := 0; n < 3; n++ {
					msg.Cards = append(msg.Cards, int(tong[n]))
				}
			}
			//调用信息处理函数
			js, _ := json.Marshal(&msg)
			InMsg := DkInMsg{
				Id:   MSG_GAME_INFO_HUANPAI,
				Data: string(js),
			}
			this.HandleHuanPai(this.Players[cid], &InMsg)
		}
	}
}

/////////////////////////////////////////////////////////////////

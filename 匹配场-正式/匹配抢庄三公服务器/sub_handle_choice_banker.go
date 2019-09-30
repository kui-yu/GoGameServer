package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) ChoiceBanker(p *ExtPlayer, d *DkInMsg) {
	req := new(GChoiceMultipleREQ)
	err := json.Unmarshal([]byte(d.Data), req)
	if err != nil {
		logs.Debug("抢庄json解析失败:%v", err)
	}
	res := GChoiceMultipleRES{
		Id:       MSG_GAME_INFO_BANKER_MULTIPLE_REPLY,
		Uid:      p.Uid,
		Multiple: req.Multiple,
	}
	if this.GameState != STAGE_BANKER_MULTIPLE {
		res.Err = "不是抢庄阶段"
		res.Result = 1
		p.SendNativeMsg(MSG_GAME_INFO_BANKER_MULTIPLE_REPLY, res)
		return
	}
	if p.BankerInfos.IsChoice != 0 {
		res.Err = "请勿重复操作"
		res.Result = 2
		p.SendNativeMsg(MSG_GAME_INFO_BANKER_MULTIPLE_REPLY, res)
		return
	}
	if req.Multiple == 1 { //如果抢庄
		this.ChoiceBankerArr = append(this.ChoiceBankerArr, PlayerMultiple{ //添加抢庄信息
			Player: p,
			IsJoin: int(req.Multiple),
		})
		p.BankerInfos.IsChoice = 2
	} else {
		p.BankerInfos.IsChoice = 1
	}
	this.HandlePlayer++
	this.BroadcastAll(MSG_GAME_INFO_BANKER_MULTIPLE_REPLY, res)
	if this.HandlePlayer == len(this.Players) { //所有人都抢庄了，执行选择庄家
		this.DelTimer(1) //删除计时器
		//选出庄家
		banker := this.ChoiceBankerByCoins()
		banker.BankerInfos.IsBanker = true
		this.DeskBankerInfos.BankerId = banker.ChairId //赋值庄家id
		//群发选庄结果
		joins := []int32{}
		for _, v := range this.ChoiceBankerArr {
			joins = append(joins, v.Player.ChairId)
		}
		this.BroadcastAll(MSG_GAME_INFO_CHOICE_BANKER_REPLY, GBankerInfo{
			Id:      MSG_GAME_INFO_CHOICE_BANKER_REPLY,
			ChairId: banker.ChairId,
			Players: joins,
		})
		//进入选择倍数阶段
		this.runTimer(len(this.Players)+1, this.StageChoiceMultiple)
	}
}

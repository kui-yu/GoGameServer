package main

import (
	"fmt"
)

//选倍数模式
func (this *ExtDesk) StageChoiceMultiple(d interface{}) {
	this.GameState = STAGE_IDLE_MULTIPLE
	this.BroadStageTime(gameConfig.Stage_Idle_Multiple_Timer)
	this.GetMultipleArr() //找出可下注倍数范围
	this.AddTimer(2, gameConfig.Stage_Idle_Multiple_Timer, this.StageChoiceMultipleEnd, "")
}

//找出可下注倍数范围
func (this *ExtDesk) GetMultipleArr() {
	res := GStageMultipleRES{
		Id:       MSG_GAME_INFO_MAX_MULTIPLE_REPLY,
		Multiple: -1,
	}
	for k, v := range this.Players {
		//跳过庄家
		if int32(k) == this.DeskBankerInfos.BankerId {
			continue
		}
		//如果够赔
		for i := len(gameConfig.Idle_Multiple) - 1; i >= 0; i-- {
			num := int64(gameConfig.Idle_Multiple[i]) * this.BScore * 5
			//闲家金币够输并且庄家金币够输的情况
			fmt.Printf("庄家ChairId:%v,当前玩家人数:%v\n", this.DeskBankerInfos.BankerId, len(this.Players))
			if v.Coins >= num && num*int64(len(this.Players)) <= this.Players[this.DeskBankerInfos.BankerId].Coins {
				res.Multiple = i
			}
		}
		v.SendNativeMsg(MSG_GAME_INFO_MAX_MULTIPLE_REPLY, res)
	}
}
func (this *ExtDesk) StageChoiceMultipleEnd(d interface{}) {
	//判断下注倍数超时，没操作默认最小倍数
	req := GChoiceMultipleREQ{Id: MSG_GAME_INFO_IDLE_MULTIPLE}
	for _, v := range this.Players {
		if v.BankerInfos.IsBanker {
			continue
		}
		if v.BankerInfos.Multiple == 0 {
			req.Uid = v.Uid
			GDeskMgr.AddNativeMsg(MSG_GAME_INFO_IDLE_MULTIPLE, v.Uid, req)
		}
	}
}

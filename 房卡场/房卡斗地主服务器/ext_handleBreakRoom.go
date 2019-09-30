package main

import (
	"encoding/json"
	"fmt"
)

//处理房间解散请求
func (this *ExtDesk) handleBreakRomm(p *ExtPlayer, d *DkInMsg) {
	fmt.Println("接收到请求")
	if this.GameState != GAME_STATUS_READ && this.GameState != GAME_STATUS_BreakRoomVote {
		fmt.Println("不是准备阶段")
		//必须在准备状态使才能够发起解散
		rr := GBreakRommReplay{
			Id:        MSG_GAME_INFO_BREAKROOM_REPLY,
			IsDismiss: 2,
			Message:   "不是准备阶段！",
		}
		p.SendNativeMsg(MSG_GAME_INFO_BREAKROOM_REPLY, rr)
		fmt.Println("rr:", rr)
		return
	}
	// if this.FkOwner == p.ChairId {
	// 	//必须是房主才能够发起房间解散投票
	// 	p.SendNativeMsg(MSG_GAME_INFO_BREAKROOM_REPLY, GBreakRommReplay{
	// 		Id:     MSG_GAME_INFO_BREAKROOM_REPLY,
	// 		Result: 2,
	// 		Err:    "对不起您不是房主！无法发起房间解散投票!",
	// 	})
	// 	return
	// }
	//可以发起房间解散投票
	// p.SendNativeMsg(MSG_GAME_INFO_BREAKROOM_REPLY, GBreakRommReplay{
	// 	Id:     MSG_GAME_INFO_BREAKROOM_REPLY,
	// 	Result: 0,
	// })
	dv1 := DisMiss{}
	json.Unmarshal([]byte(d.Data), &dv1)
	res := GBreakRommReplay{
		Id: MSG_GAME_INFO_BREAKROOM_REPLY,
	}
	fmt.Println("dv1", dv1.IsDismiss)
	if dv1.IsDismiss == 1 {
		fmt.Println("发现是一")
		p.IsDimiss = 1
		var pd bool
		for _, v := range this.Players {
			if v.Uid != p.Uid {
				if v.IsDimiss != -1 {
					pd = true
				}
			}
		}
		if pd {
			fmt.Println("不是第一个")
			var pd bool = true
			for _, v := range this.Players {
				if v.IsDimiss != 1 {
					pd = false
				}
			}

			if pd {
				fmt.Println("全部同意")
				res.IsDismiss = 3
				this.BroadcastAll(MSG_GAME_INFO_BREAKROOM_REPLY, res)
				this.TList = []*Timer{}
				this.allB()
			} else {
				fmt.Println("部分同意")
				res.IsDismiss = 1
				var agreetCha []int32
				for _, v := range this.Players {
					if v.IsDimiss == 1 {
						agreetCha = append(agreetCha, v.ChairId)
					}
				}
				res.DisPlayer = agreetCha
				this.BroadcastAll(MSG_GAME_INFO_BREAKROOM_REPLY, res)
			}
		} else {
			res.IsDismiss = 1
			var agreetCha []int32
			for _, v := range this.Players {
				if v.IsDimiss == 1 {
					agreetCha = append(agreetCha, v.ChairId)
				}
			}
			res.DisPlayer = agreetCha
			//改变状态
			this.GameState = GAME_STATUS_BreakRoomVote
			this.BroadStageTime(TIMER_BREAKROOM_NUM) //添加计时器
			//添加计时器
			this.AddTimer(TIMER_BREAKROOM, TIMER_BREAKROOM_NUM, this.TimerBreamRoom, nil)
			this.BroadcastAll(MSG_GAME_INFO_BREAKROOM_REPLY, res)
		}
	} else {
		res.IsDismiss = 0
		res.Message = "玩家" + p.Nick + "不同意，解散房间失败！"
		this.BroadcastAll(MSG_GAME_INFO_BREAKROOM_REPLY, res)
		for _, v := range this.Players {
			v.IsDimiss = -1
		}
		this.TList = []*Timer{}
		this.GameState = GAME_STATUS_READ
		this.BroadStageTime(TIMER_READ_NUM) //添加计时器
	}

}

//超时后 将剩余还没有进行投票的玩家 默认同意
func (this *ExtDesk) TimerBreamRoom(d interface{}) {
	for _, v := range this.Players {
		if v.IsDimiss == -1 {
			d := DisMiss{
				IsDismiss: 1,
			}
			dv, _ := json.Marshal(d)
			this.handleBreakRomm(v, &DkInMsg{
				Uid:  v.Uid,
				Data: string(dv),
			})
		}
	}
}

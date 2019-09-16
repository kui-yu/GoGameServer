package main

import (
	"encoding/json"
	"fmt"
)

var G_MgrPushMsg MgrPushMsg

type MgrPushMsg struct {
}

//////////////////////////
//以下为推送处理函数
//推送金币改变给用户
func (this *MgrPushMsg) PushPlayerInfoChange(id int64, msg map[string]interface{}) {
	GHall.UpdatePlayerInfo(id, msg)
}

func (this *MgrPushMsg) PushMsg(id int64, msg interface{}) error {
	return GHall.SendMsg(id, msg)
}

//////////////////////////
//以下为配合推动封装的大厅接口
//修改大厅中用户缓存的金币
func (this *Hall) UpdatePlayerInfo(uid int64, playerInfo map[string]interface{}) {
	var player Player

	this.Lk.Lock()
	p, ok := this.Players[uid]
	if ok {
		pj, _ := json.Marshal(p)
		_ = json.Unmarshal(pj, player)
		var coin int64
		var head string

		for k, v := range playerInfo {
			switch k {
			case "coin":
				coin = int64(v.(float64))
				player.Coins = coin
			case "head":
				head = v.(string)
				player.Head = head
			}
		}

		if p.GameId == 0 {
			p.Coin += coin
			p.Head = head
		}
	} else {
		this.Lk.Unlock()
		return
	}
	this.Lk.Unlock()

	if p.GameId != 0 {
		GDeskMgr.AddNativeMsg(MSG_GAME_UPDATE_PLAYER_INFO, uid, &GUpdatePlayerInfo{
			Id:         MSG_GAME_UPDATE_PLAYER_INFO,
			PlayerInfo: player,
		})
	}
}

func (this *Hall) SendMsg(uid int64, msg interface{}) error {
	this.Lk.Lock()
	p, ok := this.Players[uid]
	if !ok {
		this.Lk.Unlock()
		return fmt.Errorf("player no exist")
	}
	this.Lk.Unlock()
	//
	p.AddMsgNative(MSG_HALL_PUSH_CLIENT, &PMsgToClientWebMsg{
		Id:  MSG_HALL_PUSH_CLIENT,
		Msg: msg,
	}, false)
	return nil
}

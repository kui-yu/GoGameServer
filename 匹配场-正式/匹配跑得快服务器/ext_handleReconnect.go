package main

import (
	"encoding/json"
	"logs"
)

func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d *DkInMsg) {
	p.LiXian = false
	logs.Debug("接收到玩家的重新连接请求", p.Nick)
	if this.GameState == GAME_STATUS_FREE {
		p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
			Id:     MSG_GAME_RECONNECT_REPLY,
			Result: 1,
			Err:    "本桌子没有正在的游戏",
		})
		return
	}
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{ //先返回是否重连成功
		Id:       MSG_GAME_RECONNECT_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})
	//广播玩家上线
	this.BroadcastAll(MSG_GAME_INFO_DISORREC_BRO, DisOrRecBro{
		Id:  MSG_GAME_INFO_DISORREC_BRO,
		Cid: p.ChairId,
		Ctl: 2,
	})
	//发用户信息
	for _, v := range this.Players {
		if v.Uid == p.Uid {
			gameReply := GInfoAutoGameReply{
				Id:        MSG_GAME_INFO_AUTO_REPALY,
				GameState: this.GameState,
			}
			for _, p := range this.Players {
				seat := GSeatInfo{
					Uid:  p.Uid,
					Nick: p.Nick,
					Cid:  p.ChairId,
					Sex:  p.Sex,
					Head: p.Head,
					Lv:   p.Lv,
					Coin: p.Coins,
				}
				if p.Uid != v.Uid {
					seat.Nick = "****" + seat.Nick[len(seat.Nick)-4:]
				}
				gameReply.Seat = append(gameReply.Seat, seat)
			}
			v.SendNativeMsg(MSG_GAME_INFO_AUTO_REPALY, &gameReply) //因为每个玩家得到的座位信息不同，所以需要要意义赋值
		}
	}
	//创建重连结构体
	reC := GInfoReConnectReply{
		Id:        MSG_GAME_INFO_RECONNECT,
		GameState: this.GameState,
		TimerNum:  this.GetTimerNum(this.GameState),
		Cid:       p.ChairId,
		BScore:    this.Bscore,
		JuHao:     this.JuHao,
		Curcid:    this.CurCid,
	}
	if this.CurCid == p.ChairId {
		tishi, _, _ := this.CanOutCards(p)
		for _, v := range tishi {
			list := []int{}
			for _, v1 := range v {
				list = append(list, int(v1))
			}
			reC.Hint = append(reC.Hint, list)
		}
	}
	reC.LastOutCards = LastOutCardsToCli{
		Cid:  this.LastOutCards.Cid,
		Type: this.LastOutCards.Type,
		Max:  this.LastOutCards.Max,
	}
	if this.GameState == GAME_STATUS_SENDCAR {
		reC.StateTime = GAME_STATUS_SENDCAR_TIME
	} else if this.GameState == GAME_STATUS_OUTCARD {
		reC.StateTime = GAME_STATUS_OUTCARD_TIME
	}
	p.HandCards = Sort(p.HandCards)
	for _, v := range p.HandCards {
		reC.Cards = append(reC.Cards, int(v))
	}
	for _, v := range this.Players {
		deskCardInfo := DeskCardInfo{
			DeskOutCard: LastOutCardsToCli{
				Max:  this.OutCards[v.ChairId].Max,
				Cid:  this.OutCards[v.ChairId].Cid,
				Type: this.OutCards[v.ChairId].Type,
			},
		}
		for _, v1 := range this.OutCards[v.ChairId].Cards {
			deskCardInfo.DeskOutCard.Cards = append(deskCardInfo.DeskOutCard.Cards, int(v1))
		}
		if v.Pass {
			deskCardInfo.PlayerDo = 2
		} else {
			if len(v.HandCards) < 16 {
				deskCardInfo.PlayerDo = 1
			} else {
				deskCardInfo.PlayerDo = 0
			}
		}
		reC.Seats = append(reC.Seats, GSeatInfo{
			Uid:  v.Uid,
			Nick: v.Nick,
			Cid:  v.ChairId,
			Sex:  v.Sex,
			Head: v.Head,
			Lv:   v.Lv,
			Coin: v.Coins,
		})
		reC.TuoGuans = append(reC.TuoGuans, v.TuoGuan)
		reC.LiXian = append(reC.LiXian, v.LiXian)
		reC.CardsNum = append(reC.CardsNum, len(v.HandCards))
		reC.DeskCardInfos = append(reC.DeskCardInfos, deskCardInfo)
	}
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT, reC)
}
func (this *ExtDesk) HandleDiconnect(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("玩家断线")
	if this.GameState == GAME_STATUS_FREE {
		this.LeaveByForce(p)
		logs.Debug("已经成功删除玩家")
		return
	}
	p.LiXian = true
	//广播玩家 离线
	this.BroadcastAll(MSG_GAME_INFO_DISORREC_BRO, DisOrRecBro{
		Id:  MSG_GAME_INFO_DISORREC_BRO,
		Cid: p.ChairId,
		Ctl: 1,
	})
	if this.GameState == GAME_STATUS_OUTCARD {
		logs.Debug("该玩家离线的时候，游戏阶段是出牌阶段")
		tuog := TuoGuan{
			Id:  MSG_GAME_INFO_TUOGUAN,
			Ctl: 1,
		}
		data, _ := json.Marshal(tuog)
		this.HandleTuoGuan(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(data),
		})
	}
}

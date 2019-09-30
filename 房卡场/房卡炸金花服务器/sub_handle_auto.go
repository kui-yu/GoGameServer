package main

import (
	"encoding/json"
	"logs"
	"strconv"
)

// import "logs"

//玩家匹配
func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {
	p.Coins = 0 //房卡场玩家金币为0
	//判断房间是否是空闲阶段
	if this.GameState != GAME_STATUS_FREE {
		p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
			Id:     MSG_GAME_FK_JOIN_REPLY,
			Result: 1,
			Err:    "游戏房间已关闭",
		})
		this.DeskMgr.LeaveDo(p.Uid)
		return
	}
	//判断人数是否满
	if len(this.Players) >= this.TableConfig.PlayerNumber && this.TableConfig.PlayerNumber != 0 {
		p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
			Id:     MSG_GAME_FK_JOIN_REPLY,
			Result: 1,
			Err:    "房间人数已满",
		})
		this.DeskMgr.LeaveDo(p.Uid)
		return
	}
	//设置chairid
	for i := 0; i < this.TableConfig.PlayerNumber; i++ {
		exist := false
		for _, v := range this.Players {
			if v.ChairId == int32(i) {
				exist = true
				break
			}
		}
		if !exist {
			//加入玩家队列
			p.ChairId = int32(i)
			break
		}
	}
	this.Players = append(this.Players, p)
	if p.ChairId == 0 { //判断是房主
		if p.Uid == this.FkOwner { //判断是房主
			this.TableConfig = GATableConfig{}
			err := json.Unmarshal([]byte(this.FkInfo), &this.TableConfig)
			if err != nil {
				logs.Debug("解析房卡信息失败: ", err)
				this.ReSet()
				return
			}
			GameRound = this.TableConfig.BetRound + 1 //最大跟注
			logs.Debug("配置: ", this.TableConfig)
			//判断房卡是否够
			minGoal := this.TableConfig.TotalRound / 5
			if int64(minGoal) > p.RoomCard {
				p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
					Id:     MSG_GAME_FK_JOIN_REPLY,
					Result: 1,
					Err:    "携带房卡少于" + strconv.Itoa(minGoal) + ",请充值",
				})
				this.ReSet()
				return
			}
			//配置玩牌时间
			// if this.TableConfig.TimeType == 1 {
			// 	this.PlayTime = 15
			// } else if this.TableConfig.TimeType == 2 {
			// 	this.PlayTime = 30
			// } else if this.TableConfig.TimeType == 3 {
			// 	this.PlayTime = 60
			// } else {
			// 	this.PlayTime = 90
			// }
			//积分专场
			if this.TableConfig.GameModule == 1 {
				this.TableConfig.BaseScore = 100
			}
			//底分
			this.Bscore = int64(this.TableConfig.BaseScore)
			//清除定时器
			this.ClearTimer()
		} else {
			//不是房主
			p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
				Id:     MSG_GAME_FK_JOIN_REPLY,
				Result: 1,
				Err:    "你不是该房间主人",
			})
			this.ReSet()
			return
		}
	}
	//发送房卡匹配成功
	p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
		Id:     MSG_GAME_FK_JOIN_REPLY,
		Result: 0,
	})

	//发送房间信息
	if this.JuHao == "" {
		this.JuHao = GetJuHao()
		this.Rate = G_DbGetGameServerData.Rate
	}

	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &GTableInfoReply{
		Id:      MSG_GAME_INFO_ROOM_NOTIFY,
		TableId: this.FkNo,
		Config:  this.TableConfig,
	})

	//群发用户信息
	gameReply := GSInfoAutoGame{
		Id: MSG_GAME_INFO_AUTO_REPLY,
	}
	for _, k := range this.Players {
		for _, v := range this.Players {
			var coin int64
			if this.TableConfig.GameModule == 2 {
				coin = v.Coins
			}
			seat := GSSeatInfo{
				Uid:     v.Uid,
				Nick:    v.Nick,
				Cid:     v.ChairId,
				Sex:     v.Sex,
				Head:    v.Head,
				Lv:      v.Lv,
				Coin:    coin,
				IsReady: v.IsReady,
			}
			if k.Account != v.Account {
				seat.Nick = MarkNickName(v.Nick)
			}
			gameReply.Seat = append(gameReply.Seat, seat)
		}
		k.SendNativeMsg(MSG_GAME_INFO_AUTO_REPLY, &gameReply)
		gameReply.Seat = []GSSeatInfo{}
	}
	this.BroadStageTime(0) //阶段时间
}

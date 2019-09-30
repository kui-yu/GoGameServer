package main

import (
	"encoding/json"
	"logs"
	"strconv"
)

//玩家匹配
func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {
	//重置玩家数组
	this.ResetPlayer(p)

	if this.GameState != GAME_STATUS_FREE {
		p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
			Id:     MSG_GAME_FK_JOIN_REPLY,
			Result: 1,
			Err:    "游戏房间已关闭",
		})
		this.DeskMgr.LeaveDo(p.Uid)
		return
	}
	//超过人不让进
	if len(this.Players) >= this.TableConfig.PlayerNumber && this.TableConfig.PlayerNumber != 0 {
		//发送房卡匹配成功
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

	if p.ChairId == 0 {
		logs.Debug("房间号码", this.FkNo)
		//房主进入，房间配置
		if p.Uid == this.FkOwner {
			//底分，局数配置，携带金币，牌局人数，玩法模式，金币消耗
			this.TableConfig = GATableConfig{}
			err := json.Unmarshal([]byte(this.FkInfo), &this.TableConfig)
			if err != nil {
				//解析失败
				this.TimerOver()
				return
			}
			if this.TableConfig.GameType == 2 {
				this.Banker = -1
			}
			minGoal := this.TableConfig.TotalRound / 5
			if int64(minGoal) > p.RoomCard {
				logs.Debug("房卡不足")
				p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
					Id:     MSG_GAME_FK_JOIN_REPLY,
					Result: 1,
					Err:    "携带房卡少于" + strconv.Itoa(minGoal) + ",请充值",
				})
				this.TimerOver()
				return
			}
			//积分专场
			if this.TableConfig.GameModule == 1 {
				this.TableConfig.BaseScore = 100
			}
		} else {
			//不是房主
			p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
				Id:     MSG_GAME_FK_JOIN_REPLY,
				Result: 1,
				Err:    "你不是该房间主人",
			})
			this.TimerOver()
			return
		}
	}

	//发送房卡匹配成功
	p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
		Id:     MSG_GAME_FK_JOIN_REPLY,
		Result: 0,
	})

	//群发用户信息
	gameReply := GSInfoAutoGame{
		Id: MSG_GAME_INFO_AUTO_REPLY,
	}
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
		gameReply.Seat = append(gameReply.Seat, seat)
	}
	this.BroadcastAll(MSG_GAME_INFO_AUTO_REPLY, &gameReply)

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

	this.BroadStageTime(0)
}

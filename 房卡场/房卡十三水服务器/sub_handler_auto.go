package main

import (
	"encoding/json"
	"fmt"
	"logs"
	"strconv"
)

//玩家匹配
func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {

	logs.Debug("所有玩家", this.Players)

	// if this.TableConfig.BaseScore != 0 {

	// 	if this.TableConfig.GameModule == 2 {
	// 		minGoal := int64(this.TableConfig.BaseScore * this.TableConfig.TotalRound * 108)
	// 		if minGoal > p.Coins {
	// 			p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
	// 				Id:     MSG_GAME_FK_JOIN_REPLY,
	// 				Result: 1,
	// 				Err:    "携带金币少于" + strconv.Itoa(int(minGoal)/100+1) + ",请充值",
	// 			})
	// 			this.DeskMgr.LeaveDo(p.Uid)
	// 			return
	// 		}
	// 	} else {
	// 		//积分模式 aa制
	// 		if this.TableConfig.PayType == 2 {
	// 			minGoal := this.getPayMoney()
	// 			if minGoal > p.Coins {
	// 				p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
	// 					Id:     MSG_GAME_FK_JOIN_REPLY,
	// 					Result: 1,
	// 					Err:    "携带金币少于" + strconv.Itoa(int(minGoal)/100) + ",请充值",
	// 				})
	// 				this.DeskMgr.LeaveDo(p.Uid)
	// 				return
	// 			}
	// 		}
	// 	}
	// }

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
	fmt.Println("配置文件游戏人数:", this.TableConfig.PlayerNumber)

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
			// //底分小于100
			// if this.TableConfig.BaseScore < 10000 && this.TableConfig.GameModule == 2 {
			// 	p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
			// 		Id:     MSG_GAME_FK_JOIN_REPLY,
			// 		Result: 1,
			// 		Err:    "底分不小于100金币",
			// 	})
			// 	this.TimerOver()
			// 	return
			// }

			logs.Debug("配置", this.TableConfig)

			minGoal := this.TableConfig.TotalRound / 5
			// var minGoal int64
			// if this.TableConfig.GameModule == 2 {
			// 	minGoal = int64(this.TableConfig.BaseScore * this.TableConfig.TotalRound * 108)
			// } else {
			// 	minGoal = this.getPayMoney()
			// }
			if int64(minGoal) > p.RoomCard {
				p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
					Id:     MSG_GAME_FK_JOIN_REPLY,
					Result: 1,
					Err:    "携带房卡少于" + strconv.Itoa(minGoal) + ",请充值",
				})
				this.TimerOver()
				return
			}

			//配置-玩牌时间 1.15s 2.30s 3.60s 4.90s
			if this.TableConfig.TimeType == 1 {
				this.PlayTime = 15
			} else if this.TableConfig.TimeType == 2 {
				this.PlayTime = 30
			} else if this.TableConfig.TimeType == 3 {
				this.PlayTime = 60
			} else {
				this.PlayTime = 90
			}
			//积分专场
			if this.TableConfig.GameModule == 1 {
				this.TableConfig.BaseScore = 100
			}
			//底分
			this.Bscore = this.TableConfig.BaseScore
			this.ClearTimer()

			//牌内容初始化
			this.CardMgr.InitCards()
			this.CardMgr.InitNormalCards(this.TableConfig.GameType)
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

	for _, v1 := range this.Players {
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
			if v1.Uid != v.Uid { //不是自己的玩家 名字隐藏
				seat.Nick = "***" + seat.Nick[len(seat.Nick)-4:]
			}
			gameReply.Seat = append(gameReply.Seat, seat)
		}
		v1.SendNativeMsg(MSG_GAME_INFO_AUTO_REPLY, &gameReply)
	}
	this.BroadStageTime(0)
}

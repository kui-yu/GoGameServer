package main

import (
	"logs"
	"math/rand"
	"time"
)

//玩家匹配
func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {

	if this.GameState != GAME_STATUS_FREE {
		return
	}

	//第一个玩家进入，初始化机器人
	if this.MaxRobot == 0 {
		//3秒后加入机器人，机器人人数不超过3人
		var maxRobotNum int
		maxRobotNum = G_DbGetGameServerData.GameConfig.MaxRobot
		logs.Debug("机器人数", maxRobotNum)
		if GetCostType() == 1 {
			//比赛场
			if maxRobotNum == 0 {
				maxRobotNum = 3
			}
			rand.Seed(time.Now().UnixNano())
			this.MaxRobot = rand.Perm(maxRobotNum)[0]
			this.MaxRobot += 1
		} else {
			//体验场
			if maxRobotNum == 0 {
				maxRobotNum = 5
			}
			this.MaxRobot = maxRobotNum
		}
		logs.Debug("机器人数 ", this.MaxRobot)

	}

	//发送匹配成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:     MSG_GAME_AUTO_REPLY,
		Result: 0,
	})
	//群发用户信息
	for _, v := range this.Players {
		gameReply := GSInfoAutoGame{
			Id: MSG_GAME_INFO_AUTO_REPLY,
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
				seat.Nick = "***" + seat.Nick[len(seat.Nick)-4:]
			}
			gameReply.Seat = append(gameReply.Seat, seat)
		}
		v.SendNativeMsg(MSG_GAME_INFO_AUTO_REPLY, &gameReply)
	}
	//发送房间信息
	if this.JuHao == "" {
		this.JuHao = GetJuHao()
		this.Bscore = G_DbGetGameServerData.Bscore
		this.Rate = G_DbGetGameServerData.Rate
	}

	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &GTableInfoReply{
		Id:      MSG_GAME_INFO_ROOM_NOTIFY,
		TableId: this.JuHao,
		BScore:  this.Bscore,
	})

	if len(this.Players) >= 2 {
		//游戏开始倒计时10秒
		this.runTimer(10, this.TimerGameStart)
	}
}

//定时器 ， 游戏开始
func (this *ExtDesk) TimerGameStart(d interface{}) {
	var playerCnt int = 0
	//机器人不能超过4个
	for _, v := range this.Players {
		if !v.Robot {
			playerCnt++
		}
	}

	if playerCnt == 0 {
		logs.Debug("没有真实玩家，游戏结束")
		this.GameOver()
	}

	if len(this.Players) >= 2 {
		logs.Debug("游戏开始")
		this.nextStage(GAME_STATUS_START)
	}
}

// 重写，添加玩家
func (this *ExtDesk) AddPlayer(p *ExtPlayer) int {

	if len(this.Players) >= GCONFIG.PlayerNum {
		return -1
	}

	var playerCnt, robotCnt int
	//机器人不能超过4个
	for _, v := range this.Players {
		if !v.Robot {
			playerCnt++
		} else {
			robotCnt++
		}
	}

	if p.Robot {
		var leave bool
		if this.NotAllowRobotInRoom {
			//不允许机器人进入
			leave = true
		}

		if robotCnt >= this.MaxRobot || this.MaxRobot == 0 || playerCnt == 0 || robotCnt >= playerCnt*2 {
			//超过机器人限制人数 || 最大机器人数为0 || 真实玩家人数为0 || 超过玩家人数的2倍
			leave = true
		}

		if leave {
			p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
				Id:     MSG_GAME_AUTO_REPLY,
				Result: 13,
			})
			return -2
		}
		//不允许机器人进入房间
		this.NotAllowRobotInRoom = true
		//间隔1~3秒匹配一次机器人
		rand.Seed(time.Now().UnixNano())
		this.runTimer((rand.Perm(3)[0] + 1), this.TimerRobotNotAllowInRoom)
	} else {
		//真实玩家进入
		if len(this.Players) == 0 && !this.NotAllowRobotInRoom {
			//首次真实玩家加入房间
			this.NotAllowRobotInRoom = true
			//3秒后加入机器人
			this.runTimer(3, this.TimerRobotNotAllowInRoom)
		}
	}

	//设置chairid
	doinsert := false
	for i, v := range this.Players {
		if i != int(v.ChairId) {
			doinsert = true
			p.ChairId = int32(i)
			nl := append([]*ExtPlayer{}, this.Players[:i]...)
			nl = append(nl, p)
			nl = append(nl, this.Players[i:]...)
			this.Players = nl
			break
		}
	}
	if !doinsert {
		p.ChairId = int32(len(this.Players))
		this.Players = append(this.Players, p)
	}
	//
	return len(this.Players)
}

//不允许 机器人进入房间
func (this *ExtDesk) TimerRobotNotAllowInRoom(d interface{}) {
	this.NotAllowRobotInRoom = false
}

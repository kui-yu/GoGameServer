package main

import (
	"encoding/json"
	"logs"
)

//机器人出牌
func (this *ExtDesk) RobotPlay(d interface{}) {

	p := d.(*ExtPlayer)

	if p.CardType == 0 {
		//看牌
		req := GAPlayerOperation{
			Id:        MSG_GAME_INFO_LOOK_CARD,
			ChairId:   p.ChairId, //座位号
			PlayCoin:  0,         //下注金币
			Operation: 1,         //操作(0弃牌，1看牌，2比牌，3加注，4跟注)
		}

		dv, _ := json.Marshal(req)
		this.HandleLookCard(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(dv),
		})
		logs.Debug("机器人看牌")
	}

	//判断机器人动作
	this.AddTimer(12, 3, this.RobotPlayAction, p)
}

//允许 机器人进入房间
func (this *ExtDesk) RobotPlayAction(d interface{}) {

	p := d.(*ExtPlayer)
	logs.Debug("机器人", p.ChairId)

	//判断手牌类型,返回操作指令
	playCmd := R_PlayAction(p)

	if playCmd == 0 {
		logs.Debug("放弃")
		req := GAPlayerOperation{
			Id:        MSG_GAME_INFO_GIVE_UP,
			ChairId:   p.ChairId, //座位号
			PlayCoin:  0,         //下注金币
			Operation: 0,         //操作(0弃牌，1看牌，2比牌，3加注，4跟注)
		}
		dv, _ := json.Marshal(req)
		this.HandleGiveUp(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(dv),
		})
	} else if playCmd == 1 {
		logs.Debug("看牌")
	} else {
		var playCoin int64    //下注筹码
		var soloChairId int32 //
		if playCmd == 3 {
			playCoin = this.MinCoin + 1
			logs.Debug("加注")
		} else if playCmd == 4 {
			playCoin = this.MinCoin
			logs.Debug("跟注")
		} else {
			logs.Debug("比牌")
			//solo玩家池子
			var soloPlayers []int
			for i := 0; i < len(this.Players); i++ { //取出比牌对象顺序
				if p.ChairId == this.Players[i].ChairId || this.Players[i].CardType == 2 {
					//跳过自己，跳过弃牌玩家
					continue
				}
				soloPlayers = append(soloPlayers, i)
			}
			soloPlayers = ListShuffle(soloPlayers)
			//solo 玩家位置
			soloChairId = this.Players[soloPlayers[0]].ChairId
		}

		req := GAPlayerOperation{
			Id:        MSG_GAME_INFO_PLAY_INFO,
			ChairId:   soloChairId, //座位号
			PlayCoin:  playCoin,    //下注金币
			Operation: playCmd,     //操作(0弃牌，1看牌，2比牌，3加注，4跟注)
		}

		dv, _ := json.Marshal(req)
		this.HandleGamePlay(p, &DkInMsg{
			Uid:  p.Uid,
			Data: string(dv),
		})
		p.PlayActions = append(p.PlayActions, playCmd)
	}

}

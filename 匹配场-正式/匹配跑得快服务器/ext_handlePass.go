package main

import (
	"logs"
)

func (this *ExtDesk) HandlePass(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("接收到玩家的过请求")
	if this.GameState != GAME_STATUS_OUTCARD {
		logs.Error("不是出牌状态")
		return
	}
	if this.CurCid != p.ChairId {
		logs.Error("不是该玩家出牌")
		return
	}
	if this.LastOutCards.Max == 0 {
		logs.Error("没有上家出牌，不能过")
		return
	}
	//如果这个玩家过之后又轮到上一个出牌的人，则清除 桌子lastout属性，由下家第一个出牌
	this.CurCid = (p.ChairId + 1) % int32(len(this.Players))
	var pd bool
	if this.CurCid == this.LastOutCards.Cid {
		pd = true
		//判断该玩家出的是不是炸弹
		if this.LastOutCards.Type == CT_BOMB_FOUR {
			//更新炸弹结算集合
			for _, v := range this.Players {
				if v.ChairId == this.CurCid {
					v.Booms += 1
					v.BoomBalance += this.Bscore * 10 * 2
					v.WinForMap[(v.ChairId+1)%int32(len(this.Players))] += this.Bscore * 10
					v.WinForMap[(v.ChairId+2)%int32(len(this.Players))] += this.Bscore * 10
				} else {
					v.BeBooms += 1
					v.BoomBalance -= this.Bscore * 10
					v.BeBoomPlayer = append(v.BeBoomPlayer, this.CurCid)
					v.LoseForMap[this.CurCid] += this.Bscore * 10
				}
			}
		}
		for i, _ := range this.OutCards {
			this.OutCards[i] = LastOutCards{}
		}
		this.LastOutCards.Max = 0
	}
	//将上家的isbaopei设置为false
	this.Players[(p.ChairId+2)%int32(len(this.Players))].IsBaoPei = false
	//广播玩家过
	passbro := PassBro{
		Id:   MSG_GAME_INFO_PASS_BRO,
		Cid:  p.ChairId,
		Next: this.CurCid,
	}
	ts := [][]int{}
	for _, v := range this.Players {
		if v.ChairId == this.CurCid {
			tishi, _, _ := this.CanOutCards(v)
			for _, v := range tishi {
				list := []int{}
				for _, v1 := range v {
					list = append(list, int(v1))
				}
				passbro.Hint = append(passbro.Hint, list)
			}
			ts = append(ts, passbro.Hint...)
		}
		v.SendNativeMsg(MSG_GAME_INFO_PASS_BRO, passbro)
		passbro.Hint = [][]int{}
	}

	if len(ts) > 0 || pd {
		//删除玩家出牌定时器
		this.DelTimer(GAME_STATUS_OUTCARD)
		nextplayer := this.Players[this.CurCid]
		if nextplayer.TuoGuan {
			this.AddTimer(GAME_STATUS_OUTCARD, 1, this.TuoGuanOut, nil)
		} else {
			this.AddTimer(GAME_STATUS_OUTCARD, GAME_STATUS_OUTCARD_TIME, this.TimerOutCard, nil)
		}
	} else {
		this.AddTimer(GAME_STATUS_OUTCARD, 2, this.TimerPass, nil)
	}
	p.Pass = true
}

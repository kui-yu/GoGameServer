package main

import (
	"fmt"
	"logs"
	"time"
)

func (this *ExtDesk) GameStateSettle() {
	this.BroadStageTime(STAGE_SETTLE_TIME)
	//进入倒计时
	this.runTimer(STAGE_SETTLE_TIME, this.GameStateSettleEnd)
}

//阶段-结算
func (this *ExtDesk) GameStateSettleEnd(d interface{}) {

	if this.TableConfig.GameType == 1 {
		//先比牌
		for i := 0; i < len(this.Players); i++ {
			if this.Banker != this.Players[i].ChairId {
				rs := SoloResult(this.Players[this.Banker], this.Players[i])
				if rs > 1 {
					//玩家赢
					winMultiple := GetNiuMultiple(this.Players[i].NiuPoint) * this.Players[i].BetMultiple
					this.Players[i].WinMultiple += winMultiple
					this.Players[this.Banker].WinMultiple -= winMultiple
				} else {
					//庄家赢
					winMultiple := GetNiuMultiple(this.Players[this.Banker].NiuPoint) * this.Players[i].BetMultiple
					this.Players[this.Banker].WinMultiple += winMultiple
					this.Players[i].WinMultiple -= winMultiple
				}
			}
		}
	} else {
		var tempWin int = 0
		//先比牌
		for i := 1; i < len(this.Players); i++ {
			rs := SoloResult(this.Players[tempWin], this.Players[i])
			if rs > 1 {
				tempWin = i
			}
		}
		// logs.Debug("this.Players", len(this.Players))
		//算分
		for i := 0; i < len(this.Players); i++ {
			if i != tempWin {
				winMultiple := this.Players[tempWin].NiuMultiple * this.Players[tempWin].BetMultiple * this.Players[i].BetMultiple
				this.Players[tempWin].WinMultiple += winMultiple
				this.Players[i].WinMultiple -= winMultiple
			}
			// logs.Debug("准备离开玩家", this.Players[i].Uid)
		}
	}

	if this.TableConfig.GameModule == 2 {
		//金币模式
		var tableMultiple int64 //桌面池子总倍数
		var tableMoney int64    //桌面池子总金额
		for _, v := range this.Players {
			winCoins := int64(v.WinMultiple * this.TableConfig.BaseScore)
			if winCoins < 0 {
				winCoins = -winCoins
				if winCoins > v.Coins {
					tableMoney += v.Coins
					v.WinCoins = -v.Coins
				} else {
					tableMoney += winCoins
					v.WinCoins = -winCoins
				}
			} else {
				tableMultiple += int64(v.WinMultiple)
			}
		}

		//平摊池子
		for _, v := range this.Players {

			if v.WinMultiple > 0 {
				v.WinCoins = int64(int64(v.WinMultiple) * tableMoney / tableMultiple)
				//手续费
				v.RateCoins = float64(v.WinCoins) * this.Rate
				v.WinCoins = v.WinCoins - int64(v.RateCoins)
			}
		}
	} else {
		for _, v := range this.Players {
			v.WinCoins = int64(v.WinMultiple * this.TableConfig.BaseScore)
		}
	}

	//添加游戏记录
	for _, v := range this.Players {
		recordInfo := GSRecordInfo{
			WinCoins: v.WinCoins,
			WinDate:  time.Now().Format("2006-01-02 15:04:05"),
		}
		v.RecordInfos = append(v.RecordInfos, recordInfo)
	}

	//广播结算
	var winInfos GSWinInfosReply
	for _, v := range this.Players {
		var coins int64
		//土豪专场
		if this.TableConfig.GameModule == 2 {
			v.Coins += v.WinCoins
			coins = v.Coins
		} else {
			v.TotalCoins += v.WinCoins
			coins = v.TotalCoins
		}
		// v.Coins += int64(v.WinCoins)
		info := GSWinInfo{
			Uid:      v.Uid,
			ChairId:  v.ChairId,
			WinCoin:  v.WinCoins,
			Coins:    coins,
			NiuCards: v.NiuCards,
			NiuPoint: v.NiuPoint,
		}
		winInfos.Infos = append(winInfos.Infos, info)
		fmt.Println(v.Nick, ":", v.WinCoins)
	}
	winInfos.Id = MSG_GAME_INFO_SETTLE
	winInfos.InfoCount = len(this.Players)
	var isLeave int32 = 0
	if this.TableConfig.TotalRound > this.Round {
		isLeave = 1
	}
	fmt.Println(winInfos)
	//数据交互
	this.PutSqlData(isLeave)
	//发送结算信息
	this.BroadcastAll(MSG_GAME_INFO_SETTLE, winInfos)

	for _, v := range this.Players {
		recordInfo := GSRecordInfos{
			Id:    MSG_GAME_INFO_RECORD_INFO_REPLY,
			Infos: v.RecordInfos,
		}
		v.SendNativeMsg(MSG_GAME_INFO_RECORD_INFO_REPLY, &recordInfo)
	}

	if this.TableConfig.TotalRound <= this.Round {
		//结束
		this.TimerOver()
	} else {
		this.nextStage(STAGE_INIT)
	}
}

//获取牛点倍数
func GetNiuMultiple(niuPoint int) int {
	if niuPoint > 10 {
		return 5
	} else if niuPoint == 10 {
		return 4
	} else if niuPoint > 7 {
		return 3
	} else if niuPoint == 7 {
		return 2
	} else {
		return 1
	}
}

//房间结束
func (this *ExtDesk) TimerOver() {
	logs.Debug("房间结束")
	//总结算消息
	if this.Round > 0 {
		var rsInfoEnd GSSettleInfoEnd
		for _, v := range this.Players {
			var coins int64 = 0
			for _, c := range v.RecordInfos {
				coins += c.WinCoins
			}
			infoEnd := GSSettlePlayInfoEnd{
				ChairId:  v.ChairId,
				WinCoins: coins,
			}
			rsInfoEnd.PlayInfos = append(rsInfoEnd.PlayInfos, infoEnd)
		}
		rsInfoEnd.Id = MSG_GAME_INFO_SETTLE_INFO_END_REPLY
		this.BroadcastAll(MSG_GAME_INFO_SETTLE_INFO_END_REPLY, &rsInfoEnd)
	}
	//游戏结束
	// logs.Debug("游戏结束", this.Players)
	this.ClearTimer()

	this.GameState = GAME_STATUS_END
	this.BroadStageTime(0)

	//玩家离开
	for _, p := range this.Players {
		p.SendNativeMsgForce(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:      MSG_GAME_LEAVE_REPLY,
			Result:  0,
			Cid:     p.ChairId,
			Uid:     p.Uid,
			Token:   p.Token,
			NoToCli: true,
		})
	}
	this.GameOverLeave()

	//归还桌子
	this.GameState = GAME_STATUS_FREE
	this.ResetTable()
	this.DeskMgr.BackDesk(this)
}

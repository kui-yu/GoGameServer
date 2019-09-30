package main

import (
	"fmt"
	"logs"
)

type ExtDesk struct {
	Desk
	Bscore              int                    //底分
	CardMgr             MgrCard                //扑克牌卡牌管理器
	CurCid              int32                  //当前操作玩家
	GetH3               int32                  //获取黑桃3玩家
	LastOutCards        LastOutCards           //上一个玩家的出牌
	NotAllowRobotInRoom bool                   //机器人进入
	OutCards            map[int32]LastOutCards //玩家桌面出牌
	BalanceResult       map[int32]int          //玩家结算
}

func (this *ExtDesk) InitExtData() {
	//初始化卡牌管理器
	this.CardMgr.InitCards()
	this.CardMgr.InitNormalCards()
	//初始化桌面出牌
	this.OutCards = make(map[int32]LastOutCards)
	//初始化玩家结算
	this.BalanceResult = make(map[int32]int)
	//添加监听
	this.addListen()

}
func (this *ExtDesk) Rest() {
	this.Bscore = 0
	this.LastOutCards = LastOutCards{}
	this.OutCards = make(map[int32]LastOutCards)
}
func (this *ExtDesk) addListen() {
	this.Handle[MSG_GAME_AUTO] = this.HandleAuto
	this.Handle[MSG_GAME_INFO_OUTCARD] = this.HandleOutCard
	this.Handle[MSG_GAME_INFO_PASS] = this.HandlePass
	this.Handle[MSG_GAME_INFO_TUOGUAN] = this.HandleTuoGuan
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDiconnect
	this.Handle[MSG_GAME_INFO_EXIT] = this.HandleExit
}

//广播阶段
func (this *ExtDesk) BroadStageTime(time int32) {
	stage := GStageInfo{
		Id:        MSG_GAME_INFO_STAGE,
		Stage:     int32(this.GameState),
		StageTime: time,
	}
	this.BroadcastAll(MSG_GAME_INFO_STAGE, &stage)
}

//根据上一把玩家出牌判断可出牌组
func (this *ExtDesk) CanOutCards(p *ExtPlayer) ([][]byte, [][]byte, [][]byte) {
	boom, result := FindType(p.HandCards, this.LastOutCards)
	rre := [][]byte{}
	newboom := [][]byte{}
	newresult := [][]byte{}
	if this.LastOutCards.Max != 0 {
		if len(boom) != 0 || len(result) != 0 { //出牌数组不为空
			if this.LastOutCards.Type == CT_BOMB_FOUR {
				for _, v := range boom {
					if GetLogicValue(v[0]) > GetLogicValue(this.LastOutCards.Max) {
						rre = append(rre, v)
						newboom = append(newboom, v)
					}
				}
			} else {
				for _, v := range boom {
					rre = append(rre, v)
					newboom = append(newboom, v)
				}
				for _, v := range result {
					if GetLogicValue(v[0]) > GetLogicValue(this.LastOutCards.Max) {
						rre = append(rre, v)
						newresult = append(newresult, v)
					}
				}
			}
		}
	} else {
		if len(boom) != 0 {
			rre = append(rre, boom...)
		}
		if len(result) != 0 {
			rre = append(rre, result...)
		}
		newboom = boom
		newresult = result
	}
	fmt.Println("原本Result:", newresult)

	if len(rre) != 0 {
		nree := [][]byte{}
		for i := len(rre) - 1; i >= 0; i-- {
			nree = append(nree, rre[i])
		}
		rre = nree
	}

	if len(newresult) != 0 {
		nresults := [][]byte{}
		for i := len(newresult) - 1; i >= 0; i-- {
			nresults = append(nresults, newresult[i])
		}
		newresult = nresults
	}

	if len(newboom) != 0 {
		nboom := [][]byte{}
		for i := len(newboom) - 1; i >= 0; i-- {
			nboom = append(nboom, newboom[i])
		}
		newboom = nboom
	}
	return rre, newboom, newresult
}

//托管出牌
func (this *ExtDesk) TuoGuanOutCards(p *ExtPlayer) []byte {
	toreturn := []byte{}
	// pd := true
	// //判断是否存在玩家保单
	// for _, v := range this.Players {
	// 	if v.IsDan {
	// 		pd = false
	// 	}
	// }
	//判断自己下家是否报单
	nextPlayer := this.Players[(p.ChairId+1)%int32(len(this.Players))]
	fmt.Println("我是：：：：：：：：：", p.Nick, "下家玩家是否保单？？", nextPlayer.IsDan, nextPlayer.Nick)
	fmt.Println("现在桌面上的出牌格式:::::", this.LastOutCards.Max, this.LastOutCards.Type)
	if this.LastOutCards.Max == 0 {
		fmt.Println("玩家", p.Nick, "第一手出牌")
		//如果为第一手出牌
		if nextPlayer.IsDan {
			logs.Debug("托管发现玩家保单", this.Players[this.CurCid].Nick)
			_, _, newresult := this.CanOutCards(p)
			fmt.Println("newResult:", newresult)
			toreturn = append(toreturn, newresult[len(newresult)-1]...)
		} else {
			fmt.Println("安全到不行--1")
			p.HandCards = Sort(p.HandCards)
			toreturn = append(toreturn, p.HandCards[len(p.HandCards)-1])
			fmt.Println("出的牌???!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", toreturn)
		}

	} else {
		rre, newboom, newresult := this.CanOutCards(p)
		if len(rre) > 0 {
			if nextPlayer.IsDan && this.LastOutCards.Type == CT_SINGLE {
				logs.Debug("发现下架家报单，需要出最大牌")
				//判断result集合中是否有值，并且取出最大值，如果没有，就直接给出炸弹
				fmt.Println("newResult:", newresult)
				if len(rre) > 0 {
					if len(newresult) > 0 {
						toreturn = newresult[len(newresult)-1]
					} else if len(newboom) > 0 {
						toreturn = newboom[len(newboom)-1]
					}
				}
			} else {
				logs.Debug("安全到不行---22")
				if len(newresult) > 0 {
					toreturn = newresult[0]
				} else if len(newboom) > 0 {
					toreturn = newboom[0]
				}
			}
		}
	}
	return toreturn
}

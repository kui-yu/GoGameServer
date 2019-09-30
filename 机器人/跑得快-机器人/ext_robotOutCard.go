package main

import (
	"fmt"
	"logs"
)

//根据上一把玩家出牌判断可出牌组
func (this *ExtRobotClient) CanOutCards() ([][]byte, [][]byte, [][]byte) {
	boom, result := FindType(this.HandCard, this.LastOutCards)
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
	//倒序
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
func (this *ExtRobotClient) TuoGuanOutCards() []byte {
	toreturn := []byte{}
	nextPlayer := (this.NextPlayer + 1) % 3
	fmt.Println("maxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:", this.LastOutCards.Max)
	fmt.Println("this.isDan", this.IsDan)
	if (this.NextPlayer == this.SeatId && this.LastOutCards.Max == 0) && len(this.HandCard) > 0 {
		fmt.Println("发现第一手出牌")
		//如果为第一手出牌
		if nextPlayer == this.IsDan {
			boom, result := FindType(this.HandCard, this.LastOutCards)
			fmt.Println("boom::::::::::::::::::::::::", boom, "result::::::::::::::::::::", result)
			logs.Debug("发现下家保单")
			rre := [][]byte{}
			if len(boom) != 0 {
				rre = append(rre, boom...)
			}
			if len(result) != 0 {
				rre = append(rre, result...)
			}
			danpai := [][]byte{}
			for _, v := range rre {
				if len(v) == 1 {
					danpai = append(danpai, v)
				}
			}
			if len(danpai) == len(rre) {
				toreturn = danpai[0]
			} else {
				toreturn = rre[len(rre)-1]
			}
		} else {
			boom, result := FindType(this.HandCard, this.LastOutCards)
			rre := [][]byte{}
			if len(boom) != 0 {
				rre = append(rre, boom...)
			}
			if len(result) != 0 {
				rre = append(rre, result...)
			}
			toreturn = rre[len(rre)-1]
		}
	} else {
		rre, newboom, newresult := this.CanOutCards()
		if len(rre) > 0 {
			//判断自己下家是否报单
			if nextPlayer == this.IsDan && this.LastOutCards.Type == CT_SINGLE {
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
				if len(newresult) > 0 {
					toreturn = newresult[len(newresult)-1]
				} else {
					toreturn = newboom[len(newboom)-1]
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
		fmt.Println("发现玩家后手出牌")
	}
	return toreturn
}

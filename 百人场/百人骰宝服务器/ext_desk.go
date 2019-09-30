package main

import (
	"fmt"
	"logs"
	"math/rand"
	"time"
)

type ExtDesk struct {
	Desk
	Round     int        //当前局数
	OnLine    int        //在线玩家
	Stage     int        //游戏现阶段
	DeskMoney int64      //现桌面总的钱
	Big       int64      //猜大的累积下注
	Small     int64      //猜小的累积下注
	Odd       int64      //猜单的累积下注
	Even      int64      //猜双的累积下注
	History   []GameInfo //游戏历史记录,最多十局
}

//游戏结果
type GameInfo struct {
	NumberOne   int64
	NumberTwo   int64
	NumberThree int64
	Big         bool
	Small       bool
	Odd         bool
	Even        bool
}

//初始化桌面
func (this *ExtDesk) initDesk() {
	this.Stage = STAGE_GAME_START_BET
	this.DeskMoney = 0
	this.Big = 0
	this.Small = 0
	this.Odd = 0
	this.Even = 0
}

//游戏循环
func (this *ExtDesk) GameServeLoop(d interface{}) {
	this.initDesk()
	this.Round++
	logs.Debug("新的一局骰宝")
	//群发游戏开始
	for _, v := range this.Players {
		v.SendNativeMsg(MSG_GAME_INFO_STAGE_BET, &GSGameStart{
			Id:    MSG_GAME_INFO_STAGE_BET,
			Time:  BET_TIMER,
			Round: this.Round,
		})
	}
	//下注时间到进入开奖结算
	this.runTimer(LOTTERY_TIMER, this.GameRun)
}

//游戏结果控制
func (this *ExtDesk) GameResultControl() []int64 {
	r := []int64{1, -1}
	arr := []int64{}
	if this.DeskMoney == 0 {
		arr = GetCount()
		return arr
	}
	fmt.Println("hear")
	rand.Seed(time.Now().UnixNano())
	if this.Big >= this.Small && this.Odd >= this.Even {
		for {
			arr = GetCount()
			if arr[0]+arr[1]+arr[2] >= 11 {
				break
			}
		}
		if (arr[0]+arr[1]+arr[2])%2 == 0 {
			b := rand.Intn(3)
			if arr[0]+arr[1]+arr[2] == 11 {
				for {
					if arr[b] != 6 {
						break
					}
					b = rand.Intn(3)
				}
				arr[b] += 1
			} else {
				if arr[b] < 6 && arr[b] > 1 {
					i := rand.Intn(2)
					arr[b] += r[i]
				} else if arr[b] == 6 {
					arr[b] -= 1
				} else if arr[b] == 1 {
					arr[b] += 1
				}
			}
		}
	} else if this.Big >= this.Small && this.Odd <= this.Even {
		for {
			arr = GetCount()
			if arr[0]+arr[1]+arr[2] >= 11 {
				break
			}
		}
		if (arr[0]+arr[1]+arr[2])%2 != 0 {
			b := rand.Intn(3)
			if arr[0]+arr[1]+arr[2] == 11 {
				for {
					if arr[b] != 6 {
						break
					}
					b = rand.Intn(3)
				}
				arr[b] += 1
			} else {
				if arr[b] < 6 && arr[b] > 1 {
					i := rand.Intn(2)
					arr[b] += r[i]
				} else if arr[b] == 6 {
					arr[b] -= 1
				} else if arr[b] == 1 {
					arr[b] += 1
				}
			}
		}
	} else if this.Big < this.Small && this.Odd > this.Even {
		for {
			arr = GetCount()
			if arr[0]+arr[1]+arr[2] <= 10 {
				break
			}
		}
		if (arr[0]+arr[1]+arr[2])%2 == 0 {
			b := rand.Intn(3)
			if arr[0]+arr[1]+arr[2] == 10 {
				for {
					if arr[b] != 6 {
						break
					}
					b = rand.Intn(3)
				}
				arr[b] += 1
			} else {
				if arr[b] < 6 && arr[b] > 1 {
					i := rand.Intn(2)
					arr[b] += r[i]
				} else if arr[b] == 6 {
					arr[b] -= 1
				} else if arr[b] == 1 {
					arr[b] += 1
				}
			}
		}

	} else if this.Big < this.Small && this.Odd < this.Even {
		for {
			arr = GetCount()
			if arr[0]+arr[1]+arr[2] <= 10 {
				break
			}
		}
		if (arr[0]+arr[1]+arr[2])%2 != 0 {
			b := rand.Intn(3)
			if arr[0]+arr[1]+arr[2] == 10 {
				for {
					if arr[b] != 6 {
						break
					}
					b = rand.Intn(3)
				}
				arr[b] += 1
			} else {
				if arr[b] < 6 && arr[b] > 1 {
					i := rand.Intn(2)
					arr[b] += r[i]
				} else if arr[b] == 6 {
					arr[b] -= 1
				} else if arr[b] == 1 {
					arr[b] += 1
				}
			}
		}
	}
	return arr
}

//开出游戏结果并结算输赢
func (this *ExtDesk) GameRun(d interface{}) {
	//清空计算器
	this.ClearTimer()
	//停止下注，并返回阶段时间
	logs.Debug("停止下注")
	rand.Seed(time.Now().UnixNano())
	this.Stage = STAGE_BET_STOP
	for _, v := range this.Players {
		v.SendNativeMsg(MSG_GAME_INFO_STAGE_STOP_BET, &GSGameStop{
			Id: MSG_GAME_INFO_STAGE_STOP_BET,
			//Stage: this.Stage,
			Time: BET_TIMER,
		})
	}
	//时间到开始新的一局
	this.runTimer(BET_TIMER, this.GameServeLoop)
	//开奖时间,随机出游戏结果并返回客户端
	this.Stage = STAGE_GAME_RESULT
	obj := GameInfo{}
	var arr = make([]int64, 0)
	//获取库存概率
	intervalrate := GetRateByInterval()
	//有
	if intervalrate > 0 {
		robotrate := int(intervalrate * 10000 / 100)
		if rand.Perm(100)[0] < robotrate {
			arr = this.GameResultControl()
		}
	} else {
		//没有则随机
		arr = GetCount()
	}
	obj.NumberOne = arr[0]
	obj.NumberTwo = arr[1]
	obj.NumberThree = arr[2]
	if arr[0]+arr[1]+arr[2] >= 3 && arr[0]+arr[1]+arr[2] <= 10 {
		obj.Small = true
	} else {
		obj.Big = true
	}
	if (arr[0]+arr[1]+arr[2])%2 == 0 {
		obj.Even = true
	} else {
		obj.Odd = true
	}
	//添加历史记录
	if len(this.History) == 0 {
		this.History = append(this.History, obj)
	} else if len(this.History) <= 8 {
		this.History = PushArr(this.History, obj)
	}
	//返回开奖阶段和结果
	response := new(GSGameHistory)
	response.InfoArr = this.History
	response.Id = MSG_GAME_INFO_STAGE_GAME_RESULT
	response.Info = obj
	logs.Debug("开奖阶段")
	for _, v := range this.Players {
		v.GameResult = obj
		v.SendNativeMsg(MSG_GAME_INFO_STAGE_GAME_RESULT, response)
	}
	//结算输赢并返回客户端
	if this.DeskMoney <= 0 {
		logs.Debug("没人下注")
	} else {
		//计算输赢
		logs.Debug("结算阶段")
		this.GetGameResult()
		logs.Debug("结算完成")
	}
}

//计算输赢
func (this *ExtDesk) GetGameResult() {
	this.Stage = STAGE_SETTLE
	for _, val := range this.Players {
		//多线程处理
		func(play *ExtPlayer) {
			obj2 := new(GSSettleInfo)
			obj2.Id = MSG_GAME_INFO_STAGE_SETTLE
			obj2.Uid = play.Uid
			obj2.GameResult = val.GameResult
			//如果有下注在进入计算
			if val.IsBet {
				//如果开大
				if this.History[0].Big {
					//并且自己有买
					if play.BetInfo[PLACE_BIG] > 0 {
						obj2.Count += play.BetInfo[PLACE_BIG]
						play.WinBets += play.BetInfo[PLACE_BIG]
						obj2.Place = append(obj2.Place, PLACE_BIG)
					}
				} else {
					if play.BetInfo[PLACE_BIG] > 0 {
						obj2.Count -= play.BetInfo[PLACE_BIG]
						play.LostBet -= play.BetInfo[PLACE_BIG]
					}
				}
				//如果开小
				if this.History[0].Small {
					//并且自己有买
					if play.BetInfo[PLACE_SMALL] > 0 {
						obj2.Count += play.BetInfo[PLACE_SMALL]
						play.WinBets += play.BetInfo[PLACE_SMALL]
						obj2.Place = append(obj2.Place, PLACE_SMALL)
					}
				} else {
					if play.BetInfo[PLACE_SMALL] > 0 {
						obj2.Count -= play.BetInfo[PLACE_SMALL]
						play.LostBet -= play.BetInfo[PLACE_SMALL]
					}
				}
				//如果开单数
				if this.History[0].Odd {
					//并且自己有买
					if play.BetInfo[PLACE_ODD] > 0 {
						obj2.Count += play.BetInfo[PLACE_ODD]
						play.WinBets += play.BetInfo[PLACE_ODD]
						obj2.Place = append(obj2.Place, PLACE_ODD)
					}
				} else {
					if play.BetInfo[PLACE_ODD] > 0 {
						obj2.Count -= play.BetInfo[PLACE_ODD]
						play.LostBet -= play.BetInfo[PLACE_ODD]
					}
				}
				//如果开双数
				if this.History[0].Even {
					//并且自己有买
					if play.BetInfo[PLACE_EVEN] > 0 {
						obj2.Count += play.BetInfo[PLACE_EVEN]
						play.WinBets += play.BetInfo[PLACE_EVEN]
						obj2.Place = append(obj2.Place, PLACE_EVEN)
					}
				} else {
					if play.BetInfo[PLACE_EVEN] > 0 {
						obj2.Count -= play.BetInfo[PLACE_EVEN]
						play.LostBet -= play.BetInfo[PLACE_EVEN]
					}
				}
				fmt.Println(play.BetInfo[0])
				fmt.Println(play.BetInfo[1])
				fmt.Println(play.BetInfo[2])
				fmt.Println(play.BetInfo[3])
				//计算费率
				if obj2.Count > 0 {
					//有赢才扣除费率
					if obj2.Count > play.LostBet {
						obj2.CountRate = int64(float64(obj2.Count) * GAME_RATE) //扣除的费率
					}
					obj2.Count = obj2.Count - obj2.CountRate //扣除费率后的收益
					play.Coin -= obj2.CountRate              //扣掉费率
					play.Coin += play.WinBets * 2            //赢的钱加上低注
					obj2.Coin = play.Coin                    //结算后的玩家金币
					obj2.MaxBet = play.MaxBet(-obj2.Count)   //计算可下注范围
				} else {
					play.Coin += play.WinBets * 2
					obj2.Coin = play.Coin
					obj2.MaxBet = play.MaxBetArr
				}
				//数据库更新
				if !play.Robot {
					var validBet int64 = 0
					if obj2.Count < 0 {
						validBet = -validBet
					}
					dbreq := GGameEnd{
						Id:          MSG_GAME_END_NOTIFY,
						GameId:      GCONFIG.GameType,
						GradeId:     GCONFIG.GradeType,
						RoomId:      GCONFIG.RoomType,
						GameRoundNo: this.JuHao,
						Mini:        false,
						SetLeave:    1,
						//ActiveUid:   play.Uid,
					}
					dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
						UserId:      play.Uid,
						UserAccount: play.Account,
						BetCoins:    1,
						ValidBet:    validBet,
						PrizeCoins:  obj2.Count,
						Robot:       play.Robot,
					})
					play.SendNativeMsgForce(MSG_GAME_END_NOTIFY, &dbreq)
				}
			}
			//发送结算给客户端
			for _, v := range this.Players {
				v.SendNativeMsg(MSG_GAME_INFO_STAGE_SETTLE, obj2)
			}
			//初始化玩家
			play.InitPlayer()
		}(val)
	}
}

//自封装定时器
func (this *ExtDesk) runTimer(t int, h func(interface{})) {
	//定时器ID，定时器时间，可执行函数，可执行参数
	this.AddTimer(10, t, h, nil)
}

package main

import (
	"encoding/json"
	// "logs"
	"fmt"
	"math/rand"
	"time"
)

//匹配数据
func (this *ExtRobotClient) HandleAuto(d string) {
	// DebugLog("匹配数据", d)
	gameReply := GInfoAutoGameReply{}
	json.Unmarshal([]byte(d), &gameReply)
	for _, seat := range gameReply.Seat {
		if seat.Uid == this.UserInfo.Uid {
			this.Coin = seat.Coin
			this.ChairId = seat.Cid
		}
	}
	this.Alive = []int32{}
	this.Leaves = []int32{}
	for i := 0; i < 5; i++ {
		if int32(i) != this.ChairId {
			this.Alive = append(this.Alive, int32(i))
			this.Leaves = append(this.Leaves, int32(i))
		}
	}
	//去重
	// this.Alive = DelRepeat(this.Alive)
	// logs.Debug("座位：", this.Alive, this.ChairId)
	this.IsCardChange = 0
	this.GameIn = true
	// fmt.Println("在线玩家", this.Leaves)
}

//操作
func (this *ExtRobotClient) HandleCallPlayer(d string) {
	data := GSPlayerCallPlayer{}
	json.Unmarshal([]byte(d), &data)

	this.CallPlayer = data.Player
	// logs.Debug("下家叫牌：", this.CallPlayer, this.UserInfo.Uid)
	this.MinCoin = data.MinCoin
	this.Round = data.Round

	//发送下注信息
	if this.CallPlayer == this.ChairId {
		//随机休眠
		// logs.Debug("操作阶段")
		rand.Seed(time.Now().UnixNano())
		randTime := rand.Perm(1)[0] + 3
		time.Sleep(time.Second * time.Duration(randTime))
		this.OperationCard()
	}
}

func (this *ExtRobotClient) OperationCard() {

	// var aliveReal int
	// for _, v := range this.Leaves {
	// 	for _, p := range this.PlayerHandCard {
	// 		if v == p.ChairId {
	// 			aliveReal++
	// 		}
	// 	}
	// }
	// //没有真实玩家
	// if aliveReal == 0 {
	// 	fmt.Println("没有真实玩家")
	// 	rand.Seed(time.Now().UnixNano())
	// 	//有概率看牌
	// 	if rand.Intn(100) < 50 && this.CardLv > 1 {
	// 		//比牌
	// 		rand.Seed(time.Now().UnixNano())
	// 		i := rand.Perm(len(this.Alive))[0]
	// 		this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.Alive[i], 0, 2)
	// 		return
	// 	}

	// 	if this.CardLv < 2 {
	// 		//弃牌
	// 		this.SendOperation(MSG_GAME_GIVE_UP, this.ChairId, 0, 0)
	// 		return
	// 	}
	// }

	//概率换牌
	if len(this.Alive) == 1 {
		if this.CardType == 0 {
			this.SendOperation(MSG_GAME_LOOK_CARD, this.ChairId, 0, 1)
			this.CardType = 1
			time.Sleep(time.Second * time.Duration(1))
		}
		fmt.Println("请求换牌", this.IsCardChange)
		if this.IsCardChange == 0 {
			info := GAMaxCard{
				Id:       MSG_GAME_INFO_CHANGE_CARD,
				CardLv:   this.MaxCardLv,
				HandCard: this.MaxHandCard,
			}
			this.AddMsgNative(MSG_GAME_INFO_CHANGE_CARD, &info)
			fmt.Println("发送成功")
			return
		}
		// for _, v := range this.PlayerHandCard {
		// 	if v.CardLv > 1 {
		// 		rand.Seed(time.Now().UnixNano())
		// 		if rand.Intn(100) < this.RobotChange && this.RobotChange != 0 {
		// 			fmt.Println("换牌")
		// 			info := GAMaxCard{
		// 				Id:       MSG_GAME_INFO_CHANGE_CARD,
		// 				CardLv:   this.MaxCardLv,
		// 				HandCard: this.MaxHandCard,
		// 			}
		// 			this.AddMsgNative(MSG_GAME_INFO_CHANGE_CARD, &info)
		// 		}
		// 		break
		// 	}
		// }
	}

	//如果桌面筹码大于30
	if this.AllCoin > 20 && len(this.Alive) > 1 {
		if this.CardType == 0 {
			this.SendOperation(MSG_GAME_LOOK_CARD, this.ChairId, 0, 1)
			this.CardType = 1
			time.Sleep(time.Second * time.Duration(2))
			// this.TimeTicker.AddTimer(2, this.OperationOpenCard, nil)
			// return
		}
	}

	//看牌操作
	if this.CardType == 0 {
		var lookFlag = false
		if this.Round == 1 { //第一回合
			rand.Seed(time.Now().UnixNano())
			//有概率看牌
			if rand.Intn(100) < 10 {
				lookFlag = true
			} else if this.ChairId == 1 || this.ChairId == 4 {
				rand.Seed(time.Now().UnixNano())
				if rand.Intn(100) < 35 {
					lookFlag = true
				}
			}
		} else if this.Round == 2 { //第二回合
			rand.Seed(time.Now().UnixNano())
			//有概率看牌
			if rand.Intn(100) < 20 {
				lookFlag = true
			}
		} else {
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(100) < 10 {
				lookFlag = true
			}
		}

		if lookFlag {
			this.SendOperation(MSG_GAME_LOOK_CARD, this.ChairId, 0, 1)
			this.CardType = 1
			time.Sleep(time.Second * time.Duration(2))
		}
	}

	//1号位打手
	if this.ChairId == 1 {
		rand.Seed(time.Now().UnixNano())
		if this.CardLv == 1 && rand.Intn(100) < 50 && this.IsCardChange == 0 {
			//弃牌
			this.SendOperation(MSG_GAME_GIVE_UP, this.ChairId, 0, 0)
			return
		}
	} else {
		//其他打手
		if this.Round < 5 && len(this.Alive) > 1 {
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(100) < 50 && this.CardLv == 1 && this.IsCardChange == 0 {
				//弃牌
				this.SendOperation(MSG_GAME_GIVE_UP, this.ChairId, 0, 0)
				return
			}

		}
	}

	//feiyu 个人不能10毛
	if (this.PayCoin >= 10 || (this.CardType == 1 && this.MinCoin*2 > 10)) && this.CardLv == 1 && this.IsCardChange == 0 {
		//弃牌
		this.SendOperation(MSG_GAME_GIVE_UP, this.ChairId, 0, 0)
		return
	}

	//存活玩家超过2个
	if len(this.Alive) > 2 {
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(100) < 20 && this.CardLv == 1 && this.IsCardChange == 0 {
			//弃牌
			this.SendOperation(MSG_GAME_GIVE_UP, this.ChairId, 0, 0)
			return
		}
	}

	//随机加注
	if this.MinCoin <= 5 {
		//唬人
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(100) < 10 && this.Round > 1 {

			rand.Seed(time.Now().UnixNano())
			randCoin := int64(rand.Perm(6)[0]) + 1
			var betCoin int64
			if this.MinCoin > randCoin {
				betCoin = this.MinCoin + int64(rand.Perm(3)[0])
			} else {
				betCoin = randCoin
			}

			if this.CardType == 0 {
				//加注
				this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.ChairId, betCoin, 3)
			} else {
				//加注
				this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.ChairId, betCoin*2, 3)
			}
			return
		}

		//好牌
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(100) < 40 && this.CardLv > 1 {

			rand.Seed(time.Now().UnixNano())
			randCoin := int64(rand.Perm(6)[0]) + 1
			var betCoin int64
			if this.MinCoin > randCoin {
				betCoin = this.MinCoin + int64(rand.Perm(3)[0])
			} else {
				betCoin = randCoin
			}

			if this.CardType == 0 {
				//加注
				this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.ChairId, betCoin, 3)
			} else {
				//加注
				this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.ChairId, betCoin*2, 3)
			}
			return
		}
	}

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) < 50 && len(this.Alive) > 0 {
		if this.CardType == 0 {
			this.SendOperation(MSG_GAME_LOOK_CARD, this.ChairId, 0, 1)
			this.CardType = 1
			time.Sleep(time.Second * time.Duration(2))
		}
		if this.CardLv > 1 {
			//比牌判断
			for _, v := range this.PlayerHandCard {
				if v.CardLv > 1 {
					if this.CardLv < v.CardLv && this.IsCardChange == 0 {
						//弃牌
						this.SendOperation(MSG_GAME_GIVE_UP, this.ChairId, 0, 0)
						return
					}
				}
			}
			//比牌
			i := rand.Perm(len(this.Alive))[0]
			this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.Alive[i], 0, 2)
			return
		}
	}

	//剩余一个玩家
	if len(this.Alive) == 1 {
		if this.CardType == 0 {
			this.SendOperation(MSG_GAME_LOOK_CARD, this.ChairId, 0, 1)
			this.CardType = 1
			time.Sleep(time.Second * time.Duration(2))
		}

		rand.Seed(time.Now().UnixNano())
		if rand.Intn(100) < 40 && this.CardLv > 1 {
			//比牌判断
			for _, v := range this.PlayerHandCard {
				if v.CardLv > 1 {
					//
					if this.CardLv < v.CardLv && this.IsCardChange == 0 {
						//弃牌
						this.SendOperation(MSG_GAME_GIVE_UP, this.ChairId, 0, 0)
						return
					}
				}
			}
			//比牌
			i := rand.Perm(len(this.Alive))[0]
			this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.Alive[i], 0, 2)
			return
		}

		rand.Seed(time.Now().UnixNano())
		if rand.Intn(100) < 50 && this.MinCoin <= 5 && this.CardLv > 1 {
			rand.Seed(time.Now().UnixNano())
			if this.CardType == 0 {
				//加注
				this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.ChairId, this.MinCoin+int64(rand.Perm(3)[0]), 3)
			} else {
				//加注
				this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.ChairId, this.MinCoin*2+int64(rand.Perm(3)[0]*2), 3)
			}
		} else {
			//跟注
			this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.ChairId, 0, 4)
		}
	} else {
		//跟注
		this.SendOperation(MSG_GAME_INFO_PLAY_INFO, this.ChairId, 0, 4)
	}

}

func (this *ExtRobotClient) SendOperation(id uint32, chairId int32, coin int64, ope int) {

	if id == MSG_GAME_GIVE_UP {
		if this.CardType == 0 {
			this.SendOperation(MSG_GAME_LOOK_CARD, this.ChairId, 0, 1)
			this.CardType = 1
			//随机休眠
			rand.Seed(time.Now().UnixNano())
			randTime := rand.Perm(2)[0] + 2
			time.Sleep(time.Second * time.Duration(randTime))
		}
	}

	info := GAPlayerOperation{
		Id:        id,
		ChairId:   chairId,
		PlayCoin:  coin,
		Operation: ope,
	}

	this.AddMsgNative(id, &info)
}

//看牌
func (this *ExtRobotClient) HandleLookCard(d string) {
	data := GSCardInfo{}
	json.Unmarshal([]byte(d), &data)

	if data.HandCards[0] != 0 {
		this.HandCard = data.HandCards
		this.CardLv = data.Lv
		// logs.Debug("收到牌消息")
	}
}

//弃牌
func (this *ExtRobotClient) HandleGiveUp(d string) {
	data := GSCardType{}
	json.Unmarshal([]byte(d), &data)

	if this.ChairId != data.ChairId && len(this.Alive) > 1 {
		for i := 0; i < len(this.Alive); i++ {
			if this.Alive[i] == data.ChairId {
				if i == len(this.Alive)-1 {
					this.Alive = append(this.Alive[:i])
				} else {
					this.Alive = append(this.Alive[:i], this.Alive[i+1:]...)
				}
				break
			}
		}
	}
}

//剔除比牌输家
func (this *ExtRobotClient) HandleContestLoser(d string) {
	data := GSPlayerPayCoin{}
	json.Unmarshal([]byte(d), &data)

	if len(data.ChairId) == 2 {
		loser := int32(0)
		if data.ChairId[0] == data.Winner {
			loser = data.ChairId[1]
		} else {
			loser = data.ChairId[0]
		}

		if len(this.Alive) > 1 {
			for i := 0; i < len(this.Alive); i++ {
				if this.Alive[i] == loser {
					if i == len(this.Alive)-1 {
						this.Alive = append(this.Alive[:i])
					} else {
						this.Alive = append(this.Alive[:i], this.Alive[i+1:]...)
					}
					break
				}
			}
		}
		// logs.Debug("比牌数组处理：", this.Alive)
	}
}

//结算
func (this *ExtRobotClient) HandlerSettle(d string) {
	// // 第二种方式，需要先定义结构体
	data := GSSettlePlayInfo{}
	json.Unmarshal([]byte(d), &data)

	time.Sleep(time.Second * time.Duration(5))
	// logs.Debug("查看结算:", data)
	for _, info := range data.CList {
		if info.ChairId == this.ChairId {
			this.Coin = info.Coins
			this.GameIn = false
			//通知可以离开
			DebugLog("EVENT_CONT_ROBOTIDLE ext_robot_handler l:481")
			controller.sendEvent(EVENT_CONT_ROBOTIDLE, this)
		}
	}
}

//接收最大牌型及其位置
func (this *ExtRobotClient) HandleGetMax(d string) {
	data := GSMaxCard{}
	json.Unmarshal([]byte(d), &data)

	this.MaxIsRobot = data.IsRobot
	this.MaxCardLv = data.CardLv
	this.MaxHandCard = data.HandCard
	this.MaxChairId = data.ChairId
	this.WinnerRole = data.WinnerRole
	this.PlayerHandCard = append(data.PlayerHandCard[:])

	// fmt.Println(data.WinnerRole, "查看谁赢")
}

//换手牌
func (this *ExtRobotClient) HandleChangeCard(d string) {
	data := GSChangeCard{}
	json.Unmarshal([]byte(d), &data)
	this.HandCard = data.HandCard
	this.CardLv = data.CardLv
	fmt.Println("机器人换牌", data)

	if data.Result == 0 {
		this.IsCardChange = 1
	}
	//如果玩家是当前机器人，则重新发送操作
	if this.CallPlayer == this.ChairId {
		this.OperationCard()
	}
}

func (this *ExtRobotClient) HandleLeave(d string) {
	// fmt.Println("触发离开")
	data := GSPlayerLeave{}
	json.Unmarshal([]byte(d), &data)

	if data.LeaveType > 0 {
		for index, chairId := range this.Leaves {
			if chairId == data.ChairId {
				this.Leaves = append(this.Leaves[:index], this.Leaves[index+1:]...)
				break
			}
		}
	}
	// fmt.Println("剩余玩家", this.Leaves)
}

//获取下注金币
func (this *ExtRobotClient) HandleGetCoinMsg(d string) {
	data := GSCoinMsg{}
	json.Unmarshal([]byte(d), &data)

	// fmt.Println("查看数组：", data.PCoin, this.ChairId, this.Alive)
	this.PayCoin = data.PCoin[int(this.ChairId)]
	this.AllCoin = data.AllCoin
}

func (this *ExtRobotClient) HandlerGameStatus(d string) {
	data := GSStageInfo{}
	json.Unmarshal([]byte(d), &data)

	if data.Stage == 14 {

		rand.Seed(time.Now().UnixNano())
		//有概率看牌
		if rand.Intn(100) < 20 {
			if this.CardType == 0 {
				time.Sleep(time.Second * time.Duration(2))
				this.SendOperation(MSG_GAME_LOOK_CARD, this.ChairId, 0, 1)
				this.CardType = 1
			}
		}
	}
}

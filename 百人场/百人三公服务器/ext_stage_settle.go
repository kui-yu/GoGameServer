package main

import (
	"time"
)

func (this *ExtDesk) GameStageSettle(d interface{}) {
	this.Stage = STAGE_SETTLE
	this.BroadStageTime(gameConfigInfo.Settle_Timer)
	for _, p := range this.Players {
		gameSettleInfo := &GSSettleInfo{
			Id:          MSG_GAME_INFO_SETTLE_REPLY,
			Time:        gameConfigInfo.Settle_Timer,
			Coin:        p.Coins,
			GameTrend:   this.GameTrend,
			RoundResult: this.RoundResult,
			CoinResult:  make([]int64, 4),
			AreaRes:     this.AreaRes,
		}
		//没下注不用进行结算计算
		if !p.IsBet {
			p.NotBet++
			p.GetBetArr(p.IsDouble, 0)
			gameSettleInfo.BetArrAble = p.BetArrAble
			p.SendNativeMsg(MSG_GAME_INFO_SETTLE_REPLY, gameSettleInfo)
			continue
		} else {
			p.Round++ //下注局数增加
		}
		//结算
		for k, v := range this.RoundResult {
			if v.WinCoins > 0 { //区域赢
				if p.PlaceBet[k] > 0 { //自己有下注
					if p.IsDouble { //翻倍模式
						gameSettleInfo.CoinResult[k] += p.PlaceBet[k] * this.HandCards[k+1].Multiple //赢的钱
						p.Coins += p.PlaceBet[k] * (this.HandCards[k+1].Multiple + 1)                //自己要增加的钱
						//p.AccumulateCoins += p.PlaceBet[k] * this.HandCards[k+1].Multiple            //总输赢的钱(玩家列表字段)
						gameSettleInfo.WinCoins += p.PlaceBet[k] * this.HandCards[k+1].Multiple //本局总输赢的钱
					} else { //平倍模式
						gameSettleInfo.CoinResult[k] += p.PlaceBet[k]
						p.Coins += p.PlaceBet[k] * 2
						//p.AccumulateCoins += p.PlaceBet[k]
						gameSettleInfo.WinCoins += p.PlaceBet[k]
					}
				}
			} else { //区域输
				if p.PlaceBet[k] > 0 { //自己有下注
					if p.IsDouble { //翻倍模式
						gameSettleInfo.CoinResult[k] -= p.PlaceBet[k] * this.HandCards[0].Multiple //输的钱
						p.Coins -= p.PlaceBet[k] * (this.HandCards[0].Multiple - 1)                //自己要减少的钱
						//p.AccumulateCoins -= p.PlaceBet[k] * this.HandCards[0].Multiple            //总输赢的钱(玩家列表字段)
						gameSettleInfo.WinCoins -= p.PlaceBet[k] * this.HandCards[0].Multiple //本局总输赢的钱
					} else { //平倍模式
						//平倍模式下输钱p.Coins不用做改变，因为下注时已经扣钱了
						gameSettleInfo.CoinResult[k] -= p.PlaceBet[k]
						//p.AccumulateCoins -= p.PlaceBet[k]
						gameSettleInfo.WinCoins -= p.PlaceBet[k]
					}

				}
			}
		}
		if GetCostType() == 1 {
			//数据库记录
			var dbreq = GGameEnd{
				Id:          MSG_GAME_END_NOTIFY,
				GameId:      GCONFIG.GameType,
				GradeId:     GCONFIG.GradeType,
				RoomId:      GCONFIG.RoomType,
				GameRoundNo: this.JuHao,
				Mini:        false,
				SetLeave:    1,
			}
			dbreq.UserCoin = append(dbreq.UserCoin, GGameEndInfo{
				UserId:      p.Uid,
				UserAccount: p.Account,
				BetCoins:    p.AllBet,
				ValidBet:    p.AllBet,
				PrizeCoins:  gameSettleInfo.WinCoins,
				Robot:       p.Robot,
			})
			if !p.Robot {
				p.SendNativeMsgForce(MSG_GAME_END_NOTIFY, dbreq)
			}
			//大厅记录
			xCards := make([][]int, 0) //闲卡
			for k, v := range this.HandCards {
				if k != 0 {
					xCards = append(xCards, v.CardValue)
				}
			}
			if !p.Robot {
				var rdreq = GGameRecord{
					Id:          MSG_GAME_END_RECORD,
					GameId:      GCONFIG.GameType,
					GradeId:     GCONFIG.GradeType,
					RoomId:      GCONFIG.RoomType,
					GradeNumber: 1,
					GameRoundNo: this.JuHao,
					BankerCard:  this.HandCards[0].CardValue,
					IdleCard:    xCards,
				}
				//翻倍平倍字段
				double := 1
				if p.IsDouble {
					double = 2
				}
				rdreq.UserRecord = append(rdreq.UserRecord, GGameRecordInfo{
					UserId:       p.Uid,
					UserAccount:  p.Account,
					Robot:        p.Robot,
					CoinsBefore:  p.Coins - gameSettleInfo.WinCoins,
					BetCoins:     p.AllBet,
					PrizeCoins:   gameSettleInfo.WinCoins,
					CoinsAfter:   p.Coins,
					MultipleType: double,
					BetArea:      p.PlaceBet,
				})
				p.SendNativeMsgForce(MSG_GAME_END_RECORD, &rdreq)
			}
		}

		if gameSettleInfo.WinCoins > 0 {
			p.AccumulateCoins++
			//前20局
			if len(p.AccumulateCoinsArr) < 20 {
				p.AccumulateCoinsArr = append(p.AccumulateCoinsArr, 1)
			} else {
				if p.AccumulateCoinsArr[0] == 1 {
					p.AccumulateCoins--
				}
				p.AccumulateCoinsArr = append(p.AccumulateCoinsArr[1:], 1)
			}
			///
		} else {
			//前20局
			if len(p.AccumulateCoinsArr) < 20 {
				p.AccumulateCoinsArr = append(p.AccumulateCoinsArr, 0)
			} else {
				if p.AccumulateCoinsArr[0] == 1 {
					p.AccumulateCoins--
				}
				p.AccumulateCoinsArr = append(p.AccumulateCoinsArr[1:], 0)
			}
			////
		}
		//输赢记录结构
		var roundSettleInfo = RoundSettleInfo{
			GradeType: GCONFIG.GradeType,
			AllBet:    p.AllBet,
			WinCoins:  gameSettleInfo.WinCoins,
			CardType:  this.HandCards[0].CardType,
			Time:      time.Now().Format("2006-01-02 15:04:05"),
		}
		//添加下注区域和牌型的记录
		for k, v := range p.PlaceBet {
			if v > 0 {
				roundSettleInfo.BetArea = append(roundSettleInfo.BetArea, Area{BetArea: k, CardType: this.HandCards[k+1].CardType})
			}
		}
		//添加输赢记录
		p.AddRoundSettle(roundSettleInfo)
		//发送结算
		gameSettleInfo.Coin = p.Coins
		p.GetBetArr(p.IsDouble, 0)
		gameSettleInfo.BetArrAble = p.BetArrAble
		if GetCostType() != 2 && !p.Robot {
			AddCD(-gameSettleInfo.WinCoins) //更新库存
			AddLocalStock(-gameSettleInfo.WinCoins)
		}
		p.SendNativeMsg(MSG_GAME_INFO_SETTLE_REPLY, gameSettleInfo)
		//踢出离线玩家
		if p.LiXian {
			this.PlayerLeave(p)
		}
		p.init() //初始化玩家
	}
	this.runTimer(gameConfigInfo.Settle_Timer, this.ShuffleStage)
}

//洗牌阶段
func (this *ExtDesk) ShuffleStage(d interface{}) {
	//初始化
	this.HandCards = make([]Card, 5)
	this.PlaceBet = make([]int64, 4)
	this.RoundResult = make([]AreaRes, 4)
	this.JuHao = GetJuHao()
	//三局没下注提醒或者5局没下注踢出房间
	this.ThreeTimesLeave()
	//桌面玩家更新
	this.UadatePlayer(6)
	//发送桌面玩家更新
	for _, v := range this.Players {
		v.SendNativeMsg(MSG_GAME_INFO_DESKPLAYER_REPLAY, &GSManyPlayer{
			Id:        MSG_GAME_INFO_DESKPLAYER_REPLAY,
			Players:   v.GetDeskPlayerInfo(), //获取赢金币最多的6个玩家
			AllPlayer: len(this.Players),
			JuHao:     this.JuHao,
		})
	}
	this.Stage = STAGE_SHUFFLE
	this.BroadStageTime(gameConfigInfo.Shuffle_Timer)
	this.runTimer(gameConfigInfo.Shuffle_Timer, this.MultipleChoice)
}
func (this *Desk) PlayerLeave(p *ExtPlayer) {
	p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
		Id:      MSG_GAME_LEAVE_REPLY,
		Result:  0,
		Cid:     p.ChairId,
		Uid:     p.Uid,
		Token:   p.Token,
		Robot:   p.Robot,
		NoToCli: true,
	})
	this.DelPlayer(p.Uid)
	this.DeskMgr.LeaveDo(p.Uid)
}

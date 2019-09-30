package main

func (this *ExtDesk) GameStageOpenCard(d interface{}) {
	this.Stage = STAGE_GAME_OPEN_CARD
	this.BroadStageTime(gameConfigInfo.Open_Timer)
	//每三局洗一次牌
	if this.CardsRound == 0 || this.CardsRound == 3 {
		//洗牌
		this.DeskCards = InitCards()
		//打乱牌
		this.DeskCards = DisturbCards(this.DeskCards)
		//初始化牌轮数，每3局洗一次牌
		this.CardsRound = 0
	}
	//发送牌
	this.SendCards()
	this.CardsRound++ //牌轮数加1
	//按牌的花色和大小,从大到小排序
	for _, v := range this.HandCards {
		CardSort(v.CardValue)
	}
	//比牌，得出区域输赢
	this.GetGameResult()
	//返回开牌结果
	this.BroadcastAll(MSG_GAME_INFO_OPEN_CARD_REPLY, &GSGameResult{
		Id:          MSG_GAME_INFO_OPEN_CARD_REPLY,
		Time:        gameConfigInfo.Open_Timer,
		CardsResult: this.HandCards,
	})
	this.runTimer(gameConfigInfo.Open_Timer, this.GameStageSettle)
}

//比牌，获取区域输赢,并记录走势
func (this *ExtDesk) GetGameResult() {
	zCard := this.HandCards[0]
	for k, v := range this.HandCards[1:] {
		if v.CardType > zCard.CardType { //如果闲家赢
			this.AreaRes[k] = 1                                          //区域输赢
			this.RoundResult[k].WinCoins = this.PlaceBet[k] * v.Multiple //记录区域输赢结果和倍数
			this.RoundResult[k].Multiple = v.Multiple                    //记录区域输赢结果和倍数
			this.AddGameTrend(k+1, Trend{CardType: v.CardType}, 36)      //添加黑红梅方走势
			continue
		} else if v.CardType == zCard.CardType { //如果牌型相等
			if GetGongCardCount(v.CardValue) > GetGongCardCount(zCard.CardValue) {
				this.AreaRes[k] = 1
				this.RoundResult[k].WinCoins = this.PlaceBet[k] * v.Multiple
				this.RoundResult[k].Multiple = v.Multiple
				this.AddGameTrend(k+1, Trend{CardType: v.CardType}, 36)
				continue
			} else if GetGongCardCount(v.CardValue) == GetGongCardCount(zCard.CardValue) {
				if GetCradValue(v.CardValue[0]) > GetCradValue(zCard.CardValue[0]) {
					this.AreaRes[k] = 1
					this.RoundResult[k].WinCoins = this.PlaceBet[k] * v.Multiple
					this.RoundResult[k].Multiple = v.Multiple
					this.AddGameTrend(k+1, Trend{CardType: v.CardType}, 36)
					continue
				} else if GetCradValue(v.CardValue[0]) == GetCradValue(zCard.CardValue[0]) {
					if GetCardColr(v.CardValue[0]) > GetCardColr(zCard.CardValue[0]) {
						this.AreaRes[k] = 1
						this.RoundResult[k].WinCoins = this.PlaceBet[k] * v.Multiple
						this.RoundResult[k].Multiple = v.Multiple
						this.AddGameTrend(k+1, Trend{CardType: v.CardType}, 36)
						continue
					}
				}
			}
		}
		//有走到这说明闲家输
		this.AreaRes[k] = -1
		this.RoundResult[k].WinCoins = -this.PlaceBet[k] * this.HandCards[0].Multiple
		this.RoundResult[k].Multiple = this.HandCards[0].Multiple
		this.AddGameTrend(k+1, Trend{Player: 1, CardType: zCard.CardType}, 36)
	}
	this.AddGameTrend(0, Trend{Player: 1, CardType: zCard.CardType}, 10) //添加庄家走势
}

//添加走势记录
func (this *ExtDesk) AddGameTrend(index int, data Trend, count int) {
	this.GameTrend[index] = append(this.GameTrend[index], data)
	if len(this.GameTrend[index]) > count {
		if count == 10 { //庄走势只取10条
			this.GameTrend[index] = this.GameTrend[index][1:]
		}
		if count == 36 { //闲走势满36条清空
			this.GameTrend[index] = this.GameTrend[index][6:]
		}
	}
}

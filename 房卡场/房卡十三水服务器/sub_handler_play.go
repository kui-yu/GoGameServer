package main

import (
	"encoding/json"
	// "logs"
	"math/rand"
	"time"
)

//玩家摆牌
func (this *ExtDesk) HandlePlay(p *ExtPlayer, d *DkInMsg) {

	if this.GameState != STAGE_PLAY {
		return
	}
	//已摆牌，跳过
	if p.IsPlay > 0 {
		return
	}
	data := GAPlayInfo{}
	err := json.Unmarshal([]byte(d.Data), &data)
	if err != nil {
		return
	}

	p.PlayCards = data.PlayCards
	//自己摆牌
	if data.PlayType == 0 {

		//计算手牌类型
		firstCards := ListGet(p.PlayCards, 0, 3)
		firstType := ValidateThree(firstCards)
		// logs.Debug("手牌", p.PlayCards)
		secondCards := ListGet(p.PlayCards, 3, 5)
		secondType := ValidateFive(secondCards)
		threeCards := ListGet(p.PlayCards, 8, 5)
		threeType := ValidateFive(threeCards)

		p.PlayTypes[0] = firstType
		p.PlayTypes[1] = secondType
		p.PlayTypes[2] = threeType

		//判断倒水 相公
		if FailPoker(p) == SPECIAL_FAIL {
			if this.TableConfig.FailType == 1 {
				p.SpecialType = SPECIAL_FAIL
			} else {
				p.SendNativeMsg(MSG_GAME_INFO_ERR, GSInfoErr{
					Id:  MSG_GAME_INFO_ERR,
					Err: "牌型错误",
				})
			}
			return
		}

		p.IsPlay = 1

	} else if data.PlayType == 1 {
		// logs.Debug("机器人")
		p.IsPlay = 1
		//机器人数组
		var plays []GRecommendPoker
		maxTypes := NORMAL_FIVE_KIND
		for {
			types, cards := RecommendPoker(p.HandCards, maxTypes)
			// logs.Debug("机器人牌", types, cards, maxTypes)
			if types == nil {
				break
			}
			info := GRecommendPoker{
				Types: types,
				Cards: cards,
			}
			plays = append(plays, info)
			maxTypes = types[2] - 1
			if types == nil || maxTypes == 2 {
				break
			}
		}
		// logs.Debug("机器人牌组", plays)
		//机器人出牌
		if plays[0].Types[2] >= NORMAL_SAME_COLOR {
			p.PlayTypes = plays[0].Types
			p.PlayCards = plays[0].Cards
		} else {
			rand.Seed(time.Now().UnixNano())
			var playNum int = rand.Intn(len(plays))
			p.PlayTypes = plays[playNum].Types
			p.PlayCards = plays[playNum].Cards
		}
	} else if data.PlayType != p.SpecialType {
		//不是特殊牌
		return
	} else {
		p.IsPlay = 2
		p.PlayCards = p.SpecialCards
	}

	result := GSPlayInfo{
		Id:      MSG_GAME_INFO_PLAY_REPLY,
		ChairId: p.ChairId,
	}
	this.BroadcastAll(MSG_GAME_INFO_PLAY_REPLY, &result)

	//所有摆牌完成
	allPlayFlag := true
	for _, v := range this.Players {
		if v.IsPlay == 0 {
			allPlayFlag = false
		}
	}
	if allPlayFlag {
		this.nextStage(STAGE_SETTLE)
	}
}

//相公
func FailPoker(p *ExtPlayer) int {
	//头墩
	ct1 := GCardsType{
		Type:  p.PlayTypes[0],
		Cards: ListGet(p.PlayCards, 0, 3),
	}
	//中墩
	ct2 := GCardsType{
		Type:  p.PlayTypes[1],
		Cards: ListGet(p.PlayCards, 3, 5),
	}
	//底墩
	ct3 := GCardsType{
		Type:  p.PlayTypes[2],
		Cards: ListGet(p.PlayCards, 8, 5),
	}
	rs1 := MCardsCompare(ct1, ct2)
	rs2 := MCardsCompare(ct2, ct3)
	if rs1 == 1 {
		// p.SpecialType = SPECIAL_FAIL
		// p.SpecialCards = p.PlayCards
		// p.IsPlay = 2
		return SPECIAL_FAIL
	} else if rs2 == 1 {
		return SPECIAL_FAIL
	}
	return 0
}

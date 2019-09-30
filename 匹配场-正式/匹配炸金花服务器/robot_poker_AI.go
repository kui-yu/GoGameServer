package main

//0弃牌，1看牌，2比牌，3加注，4跟注
func R_PlayAction(p *ExtPlayer) int {

	types, _ := GetCardType(p.HandCards, p.HandColor)

	switch types {
	case CARD_SINGLE:
		return R_SingleAction(p)
	case CARD_PAIR:
		return R_PairAction(p)
	default:
		//豹子
		return R_BigCardAction(p)
	}
}

//机器人高牌操作
func R_SingleAction(p *ExtPlayer) int {
	//弃牌0
	return 0
}

//机器人对子操作
func R_PairAction(p *ExtPlayer) int {

	if len(p.PlayActions) > 0 && p.PlayActions[len(p.PlayActions)-1] != 2 {
		//比牌2
		return 2
	} else {
		//跟注4
		return 4
	}
}

//机器人顺子操作
func R_BigCardAction(p *ExtPlayer) int {
	var addCoinNum, folwCoinNum int
	for _, action := range p.PlayActions {
		if action == 3 {
			addCoinNum++
		} else if action == 4 {
			folwCoinNum++
		}
	}
	if len(p.PlayActions) > 0 && p.PlayActions[len(p.PlayActions)-1] != 2 && p.PlayActions[len(p.PlayActions)-1] != 3 {
		//比牌2
		return 2
	} else if folwCoinNum-addCoinNum > 1 {
		//加注3
		return 3
	} else {
		//跟注4
		return 4
	}
}

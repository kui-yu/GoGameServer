package main

import (
	"encoding/json"
)

//玩家看牌(\/)
func (this *ExtDesk) HandleLookCard(p *ExtPlayer, d *DkInMsg) {
	if this.GameState == STAGE_SETTLE || this.GameState == GAME_STATUS_END {
		return
	}
	data := GAPlayerOperation{}
	json.Unmarshal([]byte(d.Data), &data)

	if data.Operation == 1 {
		info := GSCardInfo{
			Id:         MSG_GAME_INFO_LOOK_CARD,
			HandCards:  p.OldHandCard,
			Lv:         p.CardLv,
			ChairId:    p.ChairId,
			Model:      0,
			CoinEnough: IsCoinEnough(p.Coins, p.PayCoin, this.Bscore, this.MinCoin, p.CardType),
		}
		p.Player.SendNativeMsg(MSG_GAME_INFO_LOOK_CARD, &info)

		p.CardType = 1
		info.HandCards = []int{0, 0, 0}
		this.BroadExceptOpPlay(p, MSG_GAME_INFO_LOOK_CARD, &info)
	}
}

package main

func (this *ExtDesk) OtherPlayer(p *ExtPlayer, m *DkInMsg) {
	this.UadatePlayer(20)
	manyPlayer := append([]ManyPlayer{}, this.ManyPlayer...)
	manyPlayer = SortForPlayers(manyPlayer, 20)
	for k, v := range manyPlayer {
		if p.Account != v.Account {
			if len(manyPlayer[k].Account) > 4 {
				manyPlayer[k].Account = "***" + manyPlayer[k].Account[len(manyPlayer[k].Account)-4:]
			}
		}
	}
	p.SendNativeMsg(MSG_GAME_INFO_OTHER_PLAYER_REPLY, &GSPlayerList{
		Id:         MSG_GAME_INFO_OTHER_PLAYER_REPLY,
		PlayerInfo: manyPlayer,
	})
}
func SortForPlayers(player []ManyPlayer, count int) []ManyPlayer {
	for i := 0; i < len(player)-1; i++ {
		if i == count {
			break
		}
		for j := 0; j < len(player)-i-1; j++ {
			if player[j].AccumulateBet < player[j+1].AccumulateBet {
				player[j], player[j+1] = player[j+1], player[j]
			}
		}
	}
	if len(player) <= count {
		return player
	}
	return player[:count]
}

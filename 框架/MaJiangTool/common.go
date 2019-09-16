package MaJiangTool

func DeleteCard(ShouPai *[]byte, Cards []byte) bool {
	nh := append([]byte{}, (*ShouPai)...)
	delcount := 0
	for _, c := range Cards {
		for i, s := range nh {
			if c == s {
				delcount++
				nh = append(nh[:i], nh[i+1:]...)
				break
			}
		}
	}
	if delcount == len(Cards) {
		*ShouPai = nh
		return true
	}
	return false
}

func SignHuiToCardVector(hui HuiIe, cards []byte) {
	if hui != nil {
		for i, v := range cards {
			if hui.IsHui(v) {
				cards[i] |= Hui_Mask
			}
		}
	}
}

func GetCardColor(card byte) byte {
	return (card & CARD_COLOR) >> 4
}

func GetCardValue(card byte) byte {
	return card & CARD_VALUE
}

func SwitchToCardIndex(card byte) int {
	if card == Card_Invalid {
		return Card_Invalid
	}
	return int(GetCardColor(card))*9 + int(GetCardValue(card)) - 1
}

func SwitchToCardData(d int) byte {
	if d < 27 {
		return byte(((d / 9) << 4) | (d%9 + 1))
	} else if d < 36 {
		return byte(0x30 | (d - 27 + 1))
	} else {
		return byte(0x40 | (d - 36 + 1))
	}
}

func HandCardSize(handCard []byte, fuZi []FuZi) int {
	cnt := len(handCard) + len(fuZi)*3
	return cnt
}

//从小到大
func Sort(cs []byte) {
	for i := 0; i < len(cs)-1; i++ {
		for j := i + 1; j < len(cs); j++ {
			if cs[i] > cs[j] {
				ad := cs[j]
				cs[j] = cs[i]
				cs[i] = ad
			}
		}
	}
}

func SortByQue(cs *[]byte, color byte) {
	cs1 := []byte{}
	cs2 := []byte{}
	for _, v := range *cs {
		if GetCardColor(v) == color {
			cs1 = append(cs1, v)
		} else {
			cs2 = append(cs2, v)
		}
	}
	Sort(cs1)
	Sort(cs2)
	*cs = append([]byte{}, cs2...)
	*cs = append(*cs, cs1...)

}

func CheckCardColorCount(fuzis []FuZi, shouPai []byte, hui HuiIe) int {
	colorCnt := []int{0, 0, 0, 0, 0, 0}
	for _, c := range shouPai {
		if hui != nil && hui.IsHui(c) {
			continue
		}
		colorCnt[GetCardColor(c)]++
	}
	for _, f := range fuzis {
		colorCnt[GetCardColor(f.OperateCard)]++
	}
	colorNum := 0
	for i := 0; i < 3; i++ {
		colorNum += colorCnt[i]
	}
	return colorNum
}

func CountSameCard(cards []byte, card byte) int {
	cnt := 0
	for _, v := range cards {
		if v == card {
			cnt++
		}
	}
	return cnt
}

func CheckAllJiang(fs []FuZi, shouPai []byte, hui HuiIe) bool {
	for _, f := range fs {
		for _, v := range f.CardData {
			cv := GetCardValue(v)
			if cv != 2 && cv != 5 && cv != 8 {
				return false
			}
		}
	}
	//
	for _, c := range shouPai {
		if hui != nil && hui.IsHui(c) {
			continue
		}
		cv := GetCardValue(c)
		if cv != 2 && cv != 5 && cv != 8 {
			return false
		}
	}
	return true
}

//手上只剩一张牌那种
func CheckShiBaLuoHan(fs []FuZi, shouPai []byte) bool {
	if len(fs) < 4 {
		return false
	}
	for _, f := range fs {
		if f.WeaveKind != ActionType_Gang_An &&
			f.WeaveKind != ActionType_Gang_Ming &&
			f.WeaveKind != ActionType_Gang_PuBuGang {
			return false
		}
	}
	return true
}

func CheckYaoJiuHu(fs []FuZi, analy []byte) bool {
	for _, f := range fs {
		c := GetCardValue(f.CardData[0])
		if c != 1 && c != 9 {
			return false
		}
	}
	//
	if GetCardValue(analy[0]) != 1 && GetCardValue(analy[0]) != 9 {
		return false
	}
	for i := 2; i < len(analy)-2; i += 3 {
		if GetCardValue(analy[i]) != 1 && GetCardValue(analy[i]) != 9 {
			return false
		}
	}
	//
	return true
}

func CheckPiaoHu(fs []FuZi, analy []byte) bool {
	if len(fs) != 0 {
		for _, f := range fs {
			if f.WeaveKind == ActionType_Chi_Left ||
				f.WeaveKind == ActionType_Chi_Center ||
				f.WeaveKind == ActionType_Chi_Right {
				return false
			}
		}
	}
	//剩余牌数小于3，一定是飘胡
	if len(analy) < 3 {
		return true
	}
	//头两张是将牌，接下来三张都要一样
	for i := 2; i < len(analy)-2; i += 3 {
		if analy[i]&(^byte(Hui_Mask)) != analy[i+2]&(^byte(Hui_Mask)) {
			return false
		}
	}
	return true
}

func Hu7Dui(ChairId int, SelfFuZi []FuZi, ShouPai *[]byte, LastEvent EventIe, hui HuiIe) bool {
	if len(SelfFuZi) > 0 || len(*ShouPai) == 0 {
		return false
	}
	//标记会子标志
	SignHuiToCardVector(hui, *ShouPai)
	//
	vHand := []byte{}
	vLaiZi := []byte{}
	laiZiCount := 0
	for _, v := range *ShouPai {
		if v&Hui_Mask > 0 {
			vLaiZi = append(vLaiZi, v)
			laiZiCount++
		} else {
			vHand = append(vHand, v)
		}
	}
	//判断手牌中相同牌并在临时牌组中进行删除,如果有会牌的话提前替换并删除
	*ShouPai = []byte{}
	for i := 0; i < len(vHand)-1; i++ {
		if vHand[i] == vHand[i+1] {
			*ShouPai = append(*ShouPai, vHand[i], vHand[i])
			i++
		} else if laiZiCount > 0 {
			laiZiCount--
			*ShouPai = append(*ShouPai, vHand[i]|Hui_Mask, vHand[i])
		} else {
			return false
		}
	}
	return true
}

func IntContain(a []int, c int) bool {
	for _, v := range a {
		if v == c {
			return true
		}
	}
	return false
}

//检查是否是中张,.中张：就是1-9中间的牌，即：4、5、6等
func CheckZhongZhang(fs []FuZi, hand []byte, hui HuiIe) bool {
	for _, f := range fs {
		if GetCardColor(f.CardData[0]) >= CARD_COLOR_FengZi {
			return false
		} else {
			for _, c := range f.CardData {
				value := GetCardValue(c)
				if value == 1 || value == 9 {
					return false
				}
			}
		}
	}

	for _, hc := range hand {
		if hui != nil && hui.IsHui(hc) {
			continue
		}
		if GetCardColor(hc) >= CARD_COLOR_FengZi {
			return false
		} else {
			value := GetCardValue(hc)
			if value == 1 || value == 9 {
				return false
			}
		}
	}
	return true
}

//门清，不能有任何动作附子
func CheckMenQing(fs []FuZi) bool {
	if len(fs) == 0 {
		return true
	}
	return false
}

//天胡:指的是庄家利用最初摸到的14张牌和牌的情况。只有庄家能够达成
//此时胡牌事件还没有加入事件管理器
func CheckTianHu(mgr *EventMgr) bool {
	if len(mgr.Events) == 1 && mgr.GetLastEvent().GetStyle() == EventType_SendCard {
		return true
	}
	return false
}

//地胡：闲家摸到的第一张牌便“和牌”才算地和，而在此和牌之前，不可以有任何家“吃，碰，杠（包括暗杠）”，否则不算
//此时胡牌事件还没有加入事件管理器
func CheckDiHu(mgr *EventMgr) bool {
	//庄发牌，出牌，玩家发牌，胡。只能有这个四个动作
	if len(mgr.Events) == 3 {
		if mgr.Events[0].GetStyle() == EventType_SendCard &&
			mgr.Events[1].GetStyle() == EventType_OutCard &&
			mgr.Events[2].GetStyle() == EventType_SendCard {
			return true
		}
	}
	return false
}

//杠上开花，杠后补一张牌，胡了
//此时胡牌事件还没有加入事件管理器
func CheckGangKai(mgr *EventMgr, huChairId int) bool {
	if len(mgr.Events) >= 2 {
		if mgr.GetBackEvent(2).GetStyle() == EventType_Action {
			ev := mgr.GetBackEvent(2).(*ActionEvent)
			if ev.Fu.WeaveKind == ActionType_Gang_An ||
				ev.Fu.WeaveKind == ActionType_Gang_Ming ||
				ev.Fu.WeaveKind == ActionType_Gang_PuBuGang {
				if mgr.GetBackEvent(1).GetStyle() == EventType_SendCard &&
					mgr.GetBackEvent(1).GetChairId() == huChairId {
					return true
				}
			}
		}
	}
	return false
}

//杠后炮，杠后出牌，别人胡了
func CheckGangHouPao(mgr *EventMgr, huChairId int) bool {
	if len(mgr.Events) >= 3 {
		if mgr.GetBackEvent(3).GetStyle() == EventType_Action {
			ev := mgr.GetBackEvent(3).(*ActionEvent)
			if ev.Fu.WeaveKind == ActionType_Gang_An ||
				ev.Fu.WeaveKind == ActionType_Gang_Ming ||
				ev.Fu.WeaveKind == ActionType_Gang_PuBuGang {
				if mgr.GetBackEvent(1).GetStyle() == EventType_OutCard &&
					mgr.GetBackEvent(1).GetChairId() != huChairId {
					return true
				}
			}
		}
	}
	return false
}

//抢杠胡
func CheckQiangGangHu(mgr *EventMgr, huChairId int) bool {
	if len(mgr.Events) >= 1 {
		if mgr.GetBackEvent(1).GetStyle() == EventType_Action &&
			mgr.GetBackEvent(1).GetChairId() != huChairId {
			ev := mgr.GetBackEvent(1).(*ActionEvent)
			if ev.Fu.WeaveKind == ActionType_Gang_PuBuGang {
				return true
			}
		}
	}
	return false
}

//自摸
func CheckZiMo(mgr *EventMgr, huChairId int) bool {
	if len(mgr.Events) >= 1 {
		if mgr.GetBackEvent(1).GetStyle() == EventType_SendCard &&
			mgr.GetBackEvent(1).GetChairId() == huChairId {
			return true
		}
	}
	return false
}

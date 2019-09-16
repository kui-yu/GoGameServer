package MaJiangTool

// import (
// 	"logs"
// )

type ActionHu struct {
	BaseAction
}

func (this *ActionHu) Init(Hui HuiIe) {
	this.InitData(Hui, ActionType_Hu)
}

func (this *ActionHu) GetResult(ChairId int, SelfFuZi []FuZi, ShouPai []byte, LastEvent EventIe, Out interface{}) bool {
	//检查动作前提条件
	if !this.CheckCondition(ChairId, SelfFuZi, ShouPai, LastEvent) {
		return false
	}

	pick := byte(Card_Invalid)
	//点炮
	if LastEvent.GetStyle() == EventType_OutCard {
		ev := LastEvent.(*OutCardEvent)
		pick = ev.GetCard() //点炮的牌不能为无效牌
		if pick == Card_Invalid {
			return false
		}
	} else if LastEvent.GetStyle() == EventType_Action {
		ev := LastEvent.(*ActionEvent)
		pick = ev.GetActionFuZi().OperateCard
	}

	//不是自摸
	if pick == Card_Invalid && HandCardSize(ShouPai, SelfFuZi) != 14 {
		return false
	}
	//
	vHandCards := append([]byte{}, ShouPai...)
	SignHuiToCardVector(this.Hui, vHandCards)
	//
	plstHuTypes := Out.(*[][]byte)
	//
	this.CheckHu(vHandCards, pick, plstHuTypes)
	if len(*plstHuTypes) > 0 {
		return true
	}
	return false
}

func (this *ActionHu) ReNew(SelfFuZi *[]FuZi, ShouPai *[]byte, Event *ActionEvent, DoCheck bool) bool {
	return true
}

func (this *ActionHu) RollBack(ChairId int, SelfFuZi *[]FuZi, ShouPai []byte, Event *ActionEvent, DoCheck bool) bool {
	return true
}

//
func (this *ActionHu) CheckHu(vHandCard []byte, byPickCard byte, pHuTypes *[][]byte) {
	vCards := []byte{}
	vLaiZi := []byte{}
	nLaiziCount := 0
	needInsert := false
	if byPickCard != Card_Invalid {
		needInsert = true
	}
	//分开赖子和牌
	for _, v := range vHandCard {
		if v != 0 {
			if v&Hui_Mask > 0 {
				vLaiZi = append(vLaiZi, v)
				nLaiziCount++
			} else {
				if needInsert && v > byPickCard {
					needInsert = false
					vCards = append(vCards, byPickCard)
				}
				vCards = append(vCards, v)
			}
		}
	}
	if needInsert {
		if byPickCard&Hui_Mask > 0 {
			vLaiZi = append(vLaiZi, byPickCard)
			nLaiziCount++
		} else {
			vCards = append(vCards, byPickCard)
		}
	}
	//排序
	Sort(vCards)
	//首先移除将牌
	if len(vCards) > 0 {
		byCardLast := byte(Card_Invalid)
		for i := 0; i < len(vCards); i++ {
			if byCardLast != vCards[i] {
				if i+1 < len(vCards) && vCards[i] == vCards[i+1] {
					vData := []byte{vCards[i], vCards[i]}
					vTemp := append([]byte{}, vCards[:i]...)
					vTemp = append(vTemp, vCards[i+2:]...)
					if len(vTemp) == 0 {
						*pHuTypes = append(*pHuTypes, vData)
						return
					} else {
						this.CheckKeziOrShunzi(vTemp, nLaiziCount, pHuTypes, &vData)
					}
					byCardLast = vCards[i]
					i++
				} else if nLaiziCount > 0 {
					vData := []byte{vCards[i], vCards[i] | Hui_Mask}
					vTemp := append([]byte{}, vCards[i+1:]...)
					if len(vTemp) == 0 {
						*pHuTypes = append(*pHuTypes, vData)
						return
					} else {
						this.CheckKeziOrShunzi(vTemp, nLaiziCount-1, pHuTypes, &vData)
					}
					byCardLast = vCards[i]
				}
			}
		}
	} else {
		*pHuTypes = append(*pHuTypes, []byte{vLaiZi[0], vLaiZi[1]})
	}
}

func (this *ActionHu) CheckKeziOrShunzi(vCards []byte, nLaiziCount int, pHuTypes *[][]byte, pData *[]byte) {
	vData := append([]byte{}, (*pData)...)
	this.RemoveKezi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, &vData)

	vData = append([]byte{}, (*pData)...)
	this.RemoveShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, &vData)
}

//移除刻字（3张一样的牌）
func (this *ActionHu) RemoveKezi(vCards []byte, nLaiziCount int, pHuTypes *[][]byte, pData *[]byte) {
	byCard := byte(Card_Invalid)
	byUsed := []byte{}
	//
	byCard = vCards[0]
	if len(vCards) > 2 && vCards[0] == vCards[2] {
		byUsed = []byte{byCard, byCard, byCard}
		vCards = vCards[3:]
	} else if len(vCards) > 1 && vCards[0] == vCards[1] && nLaiziCount > 1 { //两张牌一个赖子
		nLaiziCount--
		byUsed = []byte{byCard | Hui_Mask, byCard, byCard}
		vCards = vCards[2:]
	} else if nLaiziCount > 1 { //有两个赖子（和一张单牌）
		byUsed = []byte{byCard | Hui_Mask, byCard | Hui_Mask, byCard}
		vCards = vCards[1:]
	} else { //无效
		return
	}
	*pData = append(*pData, byUsed...)
	if len(vCards) == 0 {
		*pHuTypes = append(*pHuTypes, *pData)
	} else {
		this.CheckKeziOrShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, pData)
	}
}

func (this *ActionHu) RemoveShunzi(vCards []byte, nLaiziCount int, pHuTypes *[][]byte, pData *[]byte) {
	byCard := vCards[0]
	byValue := GetCardValue(byCard)
	byType := GetCardColor(byCard)
	//
	if byType >= CARD_COLOR_FengZi {
		return
	}
	if byValue == 1 {
		this.RemoveLeftShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, pData)
	} else if byValue == 2 {
		lstTempData := append([]byte{}, (*pData)...)
		this.RemoveCenterShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, &lstTempData)
		lstTempData = append([]byte{}, (*pData)...)
		this.RemoveLeftShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, &lstTempData)
	} else if byValue == 8 {
		lstTempData := append([]byte{}, (*pData)...)
		this.RemoveCenterShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, &lstTempData)
		lstTempData = append([]byte{}, (*pData)...)
		this.RemoveRightShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, &lstTempData)
	} else if byValue == 9 {
		this.RemoveRightShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, pData)
	} else {
		lstTempData := append([]byte{}, (*pData)...)
		this.RemoveCenterShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, &lstTempData)
		lstTempData = append([]byte{}, (*pData)...)
		this.RemoveRightShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, &lstTempData)
		lstTempData = append([]byte{}, (*pData)...)
		this.RemoveLeftShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, &lstTempData)
	}
}

func (this *ActionHu) RemoveLeftShunzi(vCards []byte, nLaiziCount int, pHuTypes *[][]byte, pData *[]byte) {
	byUsed := []byte{}
	i := 1
	//顺子第一张
	byUsed = append(byUsed, vCards[0])
	//找到第二张，并且删除
	for ; i < len(vCards); i++ {
		if vCards[i] == byUsed[0]+1 {
			byUsed = append(byUsed, vCards[i])
			vCards = append(vCards[:i], vCards[i+1:]...)
			break
		}
	}
	if len(byUsed) == 1 {
		if nLaiziCount > 0 {
			nLaiziCount--
			byUsed = append(byUsed, (byUsed[0]+1)|Hui_Mask)
		} else {
			return
		}
		i = 1
	}
	//找到第三张，并且删除
	for ; i < len(vCards); i++ {
		if vCards[i] == byUsed[0]+2 {
			byUsed = append(byUsed, vCards[i])
			vCards = append(vCards[:i], vCards[i+1:]...)
			break
		}
	}
	if len(byUsed) == 2 {
		if nLaiziCount > 0 {
			nLaiziCount--
			byUsed = append(byUsed, (byUsed[0]+2)|Hui_Mask)
		} else {
			return
		}
	}
	//
	vCards = vCards[1:]
	*pData = append(*pData, byUsed...)
	if len(vCards) == 0 {
		*pHuTypes = append(*pHuTypes, *pData)
	} else {
		this.CheckKeziOrShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, pData)
	}
}

//
func (this *ActionHu) RemoveCenterShunzi(vCards []byte, nLaiziCount int, pHuTypes *[][]byte, pData *[]byte) {
	if nLaiziCount < 1 {
		return
	}
	//
	byUsed := []byte{(vCards[0] - 1) | Hui_Mask, vCards[0]}
	i := 1
	nLaiziCount--
	//
	//找到第三张，并且删除
	for ; i < len(vCards); i++ {
		if vCards[i] == byUsed[1]+1 {
			byUsed = append(byUsed, vCards[i])
			vCards = append(vCards[:i], vCards[i+1:]...)
			break
		}
	}
	if len(byUsed) == 2 {
		if nLaiziCount > 0 {
			nLaiziCount--
			byUsed = append(byUsed, (byUsed[1]+1)|Hui_Mask)
		} else {
			return
		}
	}
	//
	vCards = vCards[1:]
	*pData = append(*pData, byUsed...)
	if len(vCards) == 0 {
		*pHuTypes = append(*pHuTypes, *pData)
	} else {
		this.CheckKeziOrShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, pData)
	}
}

func (this *ActionHu) RemoveRightShunzi(vCards []byte, nLaiziCount int, pHuTypes *[][]byte, pData *[]byte) {
	if nLaiziCount < 2 {
		return
	}
	nLaiziCount -= 2
	*pData = append(*pData, []byte{(vCards[0] - 2) | Hui_Mask, (vCards[0] - 1) | Hui_Mask, vCards[0]}...)
	vCards = vCards[1:]
	if len(vCards) == 0 {
		*pHuTypes = append(*pHuTypes, *pData)
	} else {
		this.CheckKeziOrShunzi(append([]byte{}, vCards...), nLaiziCount, pHuTypes, pData)
	}
}

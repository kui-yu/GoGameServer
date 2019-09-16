package MaJiangTool

type BaseHui struct {
	SetCard byte
	HuiList []byte
}

func (this *BaseHui) GetSetCard() byte {
	return this.SetCard
}

func (this *BaseHui) GetHui() []byte {
	return this.HuiList
}

func (this *BaseHui) IsHui(card byte) bool {
	for _, v := range this.HuiList {
		if v == card {
			return true
		}
	}
	return false
}

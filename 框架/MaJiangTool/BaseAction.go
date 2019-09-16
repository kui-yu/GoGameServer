package MaJiangTool

const (
	ActionType_None               = 0  //无动作
	ActionType_FangQi             = 1  //放弃
	ActionType_Chi_Left           = 2  //左吃
	ActionType_Chi_Center         = 3  //中吃
	ActionType_Chi_Right          = 4  //右吃
	ActionType_Peng               = 5  //碰
	ActionType_DingZhang          = 6  //定掌
	ActionType_Gang_Ming          = 7  //明杠
	ActionType_Gang_An            = 8  //暗杠
	ActionType_Gang_PuBuGang      = 9  //普补杠(由碰补成明杠)
	ActionType_Gang_Xi            = 10 //喜杠
	ActionType_Gang_Yao           = 11 //幺杠
	ActionType_Gang_Jiu           = 12 //九杠
	ActionType_Gang_3Feng_QueDong = 13 //三风杠(西南北)
	ActionType_Gang_3Feng_QueNan  = 14 //三风杠(东西北)
	ActionType_Gang_3Feng_QueXi   = 15 //三风杠(东南北)
	ActionType_Gang_3Feng_QueBei  = 16 //三风杠(东南西)
	ActionType_Gang_4Feng         = 17 //四风杠
	ActionType_Gang_HuiPi         = 18 //会皮杠
	ActionType_TiHui              = 19 //提会
	ActionType_Ting               = 20 //听牌
	ActionType_Chi_Left_Ting      = 21 //左吃听
	ActionType_Chi_Center_Ting    = 22 //中吃听
	ActionType_Chi_Right_Ting     = 23 //右吃听
	ActionType_Peng_Ting          = 24 //碰听
	ActionType_Gang_Ting          = 25
	ActionType_Hu                 = 26 //胡牌
	ActionType_Hu_BaoPai          = 27 //宝牌
	ActionType_Chu_Pai            = 28 //选择出牌牌值
	ActionType_Gang_Hui           = 29 //会杠
	ActionType_Gang_Hui_1Gang     = 30 //会幺杠
	ActionType_Gang_Hui_9Gang     = 31 //会九杠
	ActionType_Gang_Hui_XiGang    = 32 //会喜杠
	ActionType_Gang_Hui_FengGang  = 33 //会风杠
	ActionType_TiHui_1Gang        = 34 //提1
	ActionType_TiHui_9Gang        = 35 //提9
	ActionType_TiHui_FengGang     = 36 //提风
	ActionType_TiHui_XiGang       = 37 //提喜
	ActionType_XiGang_Zhong       = 38 //喜杠带中
	ActionType_XiGang_Fa          = 39 //喜杠带发
	ActionType_XiGang_Bai         = 40 //喜杠带白
	ActionType_PengBao            = 41 //碰宝牌（算明杠）
	ActionType_GangBao            = 42 //杠宝牌（算暗杠）
	ActionType_QueTuiYJGang       = 43 //瘸腿幺九杠
	ActionType_HuaGang            = 44 //花杠
	ActionType_Gang_Ming_Bu       = 45 //明杠(补)
	ActionType_Gang_An_Bu         = 46 //暗杠（补）
	ActionType_Gang_PuBuGang_Bu   = 47 //普补杠(由碰补成明杠)（补）
	ActionType_ZhiDui             = 48 //支对
)

// const (
// 	ActionType_None               = 0x0000000000000000 //无动作
// 	ActionType_FangQi             = 0x1000000000000000 //放弃
// 	ActionType_Chi_Left           = 0x0000000000000001 //左吃
// 	ActionType_Chi_Center         = 0x0000000000000002 //中吃
// 	ActionType_Chi_Right          = 0x0000000000000004 //右吃
// 	ActionType_Peng               = 0x0000000000000008 //碰
// 	ActionType_DingZhang          = 0x0000000000000010 //定掌
// 	ActionType_Gang_Ming          = 0x0000000000000020 //明杠
// 	ActionType_Gang_An            = 0x0000000000000040 //暗杠
// 	ActionType_Gang_PuBuGang      = 0x0000000000000080 //普补杠(由碰补成明杠)
// 	ActionType_Gang_Xi            = 0x0000000000000100 //喜杠
// 	ActionType_Gang_Yao           = 0x0000000000000200 //幺杠
// 	ActionType_Gang_Jiu           = 0x0000000000000400 //九杠
// 	ActionType_Gang_3Feng_QueDong = 0x0000000000000800 //三风杠(西南北)
// 	ActionType_Gang_3Feng_QueNan  = 0x0000000000001000 //三风杠(东西北)
// 	ActionType_Gang_3Feng_QueXi   = 0x0000000000002000 //三风杠(东南北)
// 	ActionType_Gang_3Feng_QueBei  = 0x0000000000004000 //三风杠(东南西)
// 	ActionType_Gang_4Feng         = 0x0000000000008000 //四风杠
// 	ActionType_Gang_HuiPi         = 0x0000000000010000 //会皮杠
// 	ActionType_TiHui              = 0x0000000000020000 //提会
// 	ActionType_Ting               = 0x0000000000040000 //听牌
// 	ActionType_Chi_Left_Ting      = 0x0000000000080000 //左吃听
// 	ActionType_Chi_Center_Ting    = 0x0000000000100000 //中吃听
// 	ActionType_Chi_Right_Ting     = 0x0000000000200000 //右吃听
// 	ActionType_Peng_Ting          = 0x0000000000400000 //碰听
// 	ActionType_Gang_Ting          = 0x0000000000800000
// 	ActionType_Hu                 = 0x0000000001000000 //胡牌
// 	ActionType_Hu_BaoPai          = 0x0000000002000000 //宝牌
// 	ActionType_Chu_Pai            = 0x0000000004000000 //选择出牌牌值
// 	ActionType_Gang_Hui           = 0x0000000008000000 //会杠
// 	ActionType_Gang_Hui_1Gang     = 0x0000000010000000 //会幺杠
// 	ActionType_Gang_Hui_9Gang     = 0x0000000020000000 //会九杠
// 	ActionType_Gang_Hui_XiGang    = 0x0000000040000000 //会喜杠
// 	ActionType_Gang_Hui_FengGang  = 0x0000000080000000 //会风杠
// 	ActionType_TiHui_1Gang        = 0x0000000100000000 //提1
// 	ActionType_TiHui_9Gang        = 0x0000000200000000 //提9
// 	ActionType_TiHui_FengGang     = 0x0000000400000000 //提风
// 	ActionType_TiHui_XiGang       = 0x0000000800000000 //提喜
// 	ActionType_XiGang_Zhong       = 0x0000001000000000 //喜杠带中
// 	ActionType_XiGang_Fa          = 0x0000002000000000 //喜杠带发
// 	ActionType_XiGang_Bai         = 0x0000004000000000 //喜杠带白
// 	ActionType_PengBao            = 0x0000008000000000 //碰宝牌（算明杠）
// 	ActionType_GangBao            = 0x0000010000000000 //杠宝牌（算暗杠）
// 	ActionType_QueTuiYJGang       = 0x0000020000000000 //瘸腿幺九杠
// 	ActionType_HuaGang            = 0x0000040000000000 //花杠
// 	ActionType_Gang_Ming_Bu       = 0x0000080000000000 //明杠(补)
// 	ActionType_Gang_An_Bu         = 0x0000100000000000 //暗杠（补）
// 	ActionType_Gang_PuBuGang_Bu   = 0x0000200000000000 //普补杠(由碰补成明杠)（补）
// 	ActionType_ZhiDui             = 0x0000400000000000 //支对
// )

type HuiIe interface {
	GetSetCard() byte
	GetHui() []byte
	IsHui(byte) bool
}

////////////////////////////////////////////
type ActionIe interface {
	GetResult(int, []FuZi, []byte, EventIe, interface{}) bool
	ReNew(*[]FuZi, *[]byte, *ActionEvent, bool) bool
	RollBack(int, *[]FuZi, []byte, *ActionEvent, bool) bool
	GetStyle() int
}

type ConditionIe interface {
	Satisfy(Action ActionIe, ChairID int, FuZis []FuZi, ShouPai []byte, LastEvent EventIe) bool
}

type EventIe interface {
	GetStyle() int
	GetChairId() int
	GetCard() byte
}

///////////////////////////////////////////

type BaseAction struct {
	Style         int
	Hui           HuiIe
	ReferAction   ActionIe
	ConditionList []ConditionIe
	Supper        ActionIe
}

//初始化类型
func (this *BaseAction) InitData(hui HuiIe, st int) {
	this.Style = st
	this.Hui = hui
}

//获取动作类型
func (this *BaseAction) GetStyle() int {
	return this.Style
}

//添加条件
func (this *BaseAction) AddCondition(cond ConditionIe) {
	this.ConditionList = append(this.ConditionList, cond)
}

//触发动作
func (this *BaseAction) SetReferAction(action ActionIe) {
	this.ReferAction = action
}

//回滚
func (this *BaseAction) RollBack(ChairId int, SelfFuZi *[]FuZi, ShouPai *[]byte) bool {
	return false
}

func (this *BaseAction) CheckCondition(ChairID int, SelfFuZi []FuZi, ShouPai []byte, LastEvent EventIe) bool {
	for _, v := range this.ConditionList {
		if !v.Satisfy(this.Supper, ChairID, SelfFuZi, ShouPai, LastEvent) {
			return false
		}
	}
	return true
}

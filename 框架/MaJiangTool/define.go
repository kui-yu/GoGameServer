package MaJiangTool

const (
	Card_Invalid = 0x00
	Card_Rear    = 0xFF
	//万
	Card_Wan_1 = 0x01
	Card_Wan_2 = 0x02
	Card_Wan_3 = 0x03
	Card_Wan_4 = 0x04
	Card_Wan_5 = 0x05
	Card_Wan_6 = 0x06
	Card_Wan_7 = 0x07
	Card_Wan_8 = 0x08
	Card_Wan_9 = 0x09
	//条
	Card_Tiao_1 = 0x11
	Card_Tiao_2 = 0x12
	Card_Tiao_3 = 0x13
	Card_Tiao_4 = 0x14
	Card_Tiao_5 = 0x15
	Card_Tiao_6 = 0x16
	Card_Tiao_7 = 0x17
	Card_Tiao_8 = 0x18
	Card_Tiao_9 = 0x19
	//筒
	Card_Bing_1 = 0x21
	Card_Bing_2 = 0x22
	Card_Bing_3 = 0x23
	Card_Bing_4 = 0x24
	Card_Bing_5 = 0x25
	Card_Bing_6 = 0x26
	Card_Bing_7 = 0x27
	Card_Bing_8 = 0x28
	Card_Bing_9 = 0x29
	//风牌（东南西北）
	Card_Dong = 0x31
	Card_Nan  = 0x32
	Card_Xi   = 0x33
	Card_Bei  = 0x34

	//箭牌（中发白）
	Card_zhong = 0x35
	Card_Fa    = 0x36
	Card_Bai   = 0x37
	//花牌（春夏秋冬梅兰竹菊）
	Card_Flower_Chun = 0x41
	Card_Flower_Xia  = 0x42
	Card_Flower_Qiu  = 0x43
	Card_Flower_Dong = 0x44
	Card_Flower_Mei  = 0x45
	Card_Flower_Lan  = 0x46
	Card_Flower_Zhu  = 0x47
	Card_Flower_Ju   = 0x48
)

//卡牌颜色和值掩码，会子掩码
const (
	CARD_COLOR = 0xF0
	CARD_VALUE = 0x0F
	Hui_Mask   = 0x80
)

//卡牌颜色值
const (
	CARD_COLOR_Wan    = 0 //万
	CARD_COLOR_Tiao   = 1 //条
	CARD_COLOR_Bing   = 2 //并
	CARD_COLOR_FengZi = 3 //风、字
	CARD_COLOR_Flower = 4 //花牌
)

//结束类型,结束的方式
const (
	GAME_END_TYPE_ZIMO               = 0  //自摸
	GAME_END_TYPE_DIANPAO            = 1  //点炮
	GAME_END_TYPE_QIANGGANGHU        = 2  //抢杠胡
	GAME_END_TYPE_HUANGZHUANG        = 3  //荒庄
	GAME_END_TYPE_TAOPAO             = 4  //逃跑
	GAME_END_TYPE_BAO                = 5  //宝牌
	GAME_END_TYPE_DUIBAO             = 6  //对宝
	GAME_END_TYPE_BAOZHONGBAO        = 7  //宝中宝
	GAME_END_TYPE_KANDUIBAO          = 8  //坎对宝
	GAME_END_TYPE_DISSOLUTION        = 9  //解散
	GAME_END_TYPE_TIANHU             = 10 //天胡
	GAME_END_TYPE_DIHU               = 11 //地胡
	GAME_END_TYPE_HAIDILAOYUE        = 12 //海底捞月
	GAME_END_TYPE_HAIDILAOYUEDIANPAO = 13 //海底捞月点炮
	GAME_END_TYPE_GANGSHANGHUA       = 14 //杠上花
	GAME_END_TYPE_SHUANGGANGSHAGNHUA = 15 //双杠上花
	GAME_END_TYPE_GANGSHANGPAO       = 16 //杠上炮
	GAME_END_TYPE_SHUANGGANGSHANGPAO = 17 //双杠上炮
)

//胡牌类型，牌型的方式
const (
	HuType_None            = 0
	HuType_PingHu          = 1  //平胡
	HuType_JiaHu           = 2  //夹胡
	HuType_DanDiao         = 3  //单调
	HuType_BianHu          = 4  //边胡
	HuType_BaYi            = 5  //把一
	HuType_PiaoBaYi        = 6  //飘把一
	HuType_PiaoHu          = 7  //飘胡
	HuType_7Dui            = 8  //七对
	HuType_7Dui_HaoHua     = 9  //豪华七对
	HuType_7Dui_ChaoHaoHua = 10 //超豪华七对
	HuType_GanBa           = 11 //干巴
	HuType_Jia5            = 12 //夹五
	HuType_XiaoDuHu        = 13 //小赌胡
	HuType_DaDuHu          = 14 //大赌胡
	HuType_QingYiSe        = 15 //清一色
	HuType_BiMenHu         = 16 //闭门胡
	HuType_GangKai         = 17 //杠开
	HuType_SanJiaBi        = 18 //三家闭
	HuType_SiGuiYi         = 19 //四归一
	HuType_GangKouHu       = 20 //胡杠口
	HuType_HunYiSe         = 21 //混一色
	HuType_QiangGangHu     = 22 //抢杠胡
	HuType_GangHouPao      = 23 //杠后炮
	HuType_DuiDuiHu        = 24 //对对胡
	HuType_ZiMo            = 25 //自摸
	HuType_ZhangHu         = 26 //庄胡
	HuType_JiuGangHu       = 27 //九杠胡
	HuType_YaoGangHu       = 28 //幺杠胡
	HuType_BaoPaiHu        = 29 //宝牌胡
	HuType_WuSeHu          = 30 //无色胡
	HuType_ShiBaLuoHan     = 31 //十八罗汉（手上只剩一张牌，其余全部附子，且都是杠）
	HuType_JinGouDiao      = 32 //金钩钓，手上只剩下一张牌
	HuType_JiangHu         = 33 //将胡（所有的牌都是2,5,8）
	HuType_YaoJiuHu        = 34 //幺久胡(附子的第一张牌必须要1或9)
	HuType_TianHu          = 35 //天胡(指的是庄家利用最初摸到的14张牌和牌的情况)
	HuType_DiHu            = 36 //地胡(闲家摸到的第一张牌便“和牌”才算地和，而在此和牌之前，不可以有任何家“吃，碰，杠（包括暗杠）”，否则不算。)
	HuType_MenQing         = 37 //玩家不吃、不碰、不明杠，全凭自己摸牌上听。听牌后胡别人点的炮，就叫门清；听牌自摸就叫不求人
	HuType_BuQiuRen        = 38 //玩家不吃、不碰、不明杠，全凭自己摸牌上听。听牌后胡别人点的炮，就叫门清；听牌自摸就叫不求人
	HuType_ZhongZhang      = 39 //中张：就是1-9中间的牌，即：4、5、6等
)

const (
	REQUEST_ActionType_None   = 0 //无动作
	REQUEST_ActionType_FangQi = 1 //放弃
	//吃
	REQUEST_ActionType_Chi_Left   = 10 //左吃
	REQUEST_ActionType_Chi_Center = 11 //中吃
	REQUEST_ActionType_Chi_Right  = 12 //右吃
	//碰
	REQUEST_ActionType_Peng = 30 //碰
	//杠
	REQUEST_ActionType_Gang_Ming          = 50 //明杠
	REQUEST_ActionType_Gang_An            = 51 //暗杠
	REQUEST_ActionType_Gang_PuBuGang      = 52 //普补杠(由碰补成明杠)
	REQUEST_ActionType_Gang_Xi            = 53 //喜杠
	REQUEST_ActionType_Gang_Yao           = 54 //幺杠
	REQUEST_ActionType_Gang_Jiu           = 55 //九杠
	REQUEST_ActionType_Gang_3Feng_QueDong = 56 //三风杠(西南北)
	REQUEST_ActionType_Gang_3Feng_QueNan  = 57 //三风杠(东西北)
	REQUEST_ActionType_Gang_3Feng_QueXi   = 58 //三风杠(东南北)
	REQUEST_ActionType_Gang_3Feng_QueBei  = 59 //三风杠(东南西)
	REQUEST_ActionType_Gang_4Feng         = 60 //四风杠
	REQUEST_ActionType_Gang_HuiPi         = 61 //会皮杠
	REQUEST_ActionType_Gang_Hui           = 62 //会杠
	REQUEST_ActionType_Gang_Hui_1Gang     = 63 //会幺杠
	REQUEST_ActionType_Deprecated         = 64 //旧服务器对这个值做了特殊处理，因此废弃
	REQUEST_ActionType_Gang_Hui_XiGang    = 65 //会喜杠
	REQUEST_ActionType_Gang_Hui_FengGang  = 66 //会风杠
	REQUEST_ActionType_XiGang_Zhong       = 67 //喜杠带中
	REQUEST_ActionType_XiGang_Fa          = 68 //喜杠带发
	REQUEST_ActionType_XiGang_Bai         = 69 //喜杠带白
	REQUEST_ActionType_PengBao            = 70 //碰宝牌（算明杠）
	REQUEST_ActionType_GangBao            = 71 //杠宝牌（算暗杠）
	REQUEST_ActionType_QueTuiYJGang       = 72 //瘸腿幺九杠
	REQUEST_ActionType_HuaGang            = 73 //花杠
	REQUEST_ActionType_Gang_Hui_9Gang     = 74 //会九杠
	REQUEST_ActionType_Gang_Ming_Bu       = 75 //明杠(补)
	REQUEST_ActionType_Gang_An_Bu         = 76 //暗杠(补)
	REQUEST_ActionType_Gang_PuBuGang_Bu   = 77 //普补杠(由碰补成明杠)(补)

	REQUEST_ActionType_Ting            = 100 //听牌
	REQUEST_ActionType_Chi_Left_Ting   = 101 //左吃听
	REQUEST_ActionType_Chi_Center_Ting = 102 //中吃听
	REQUEST_ActionType_Chi_Right_Ting  = 103 //右吃听
	REQUEST_ActionType_Peng_Ting       = 104 //碰听
	REQUEST_ActionType_Gang_Ting       = 105 //杠听
	//胡牌
	REQUEST_ActionType_Hu        = 120 //胡牌
	REQUEST_ActionType_Hu_BaoPai = 121 //宝牌

	REQUEST_ActionType_TiHui          = 140 //提会
	REQUEST_ActionType_TiHui_1Gang    = 141 //提1
	REQUEST_ActionType_TiHui_9Gang    = 142 //提9
	REQUEST_ActionType_TiHui_FengGang = 143 //提风
	REQUEST_ActionType_TiHui_XiGang   = 144 //提喜

	//
	REQUEST_ActionType_DingZhang = 160 //定掌
	REQUEST_ActionType_Chu_Pai   = 161 //选择出牌牌值
	//
	REQUEST_ActionType_Invalid = 255 //错误类型
)

//附子
type FuZi struct {
	WeaveKind   int    //动作类型
	ProvideUser int    //供应用户
	OperateCard byte   //操作牌，用于界面按钮显示
	TiHuiCard   byte   //有会杠的时候，提的会牌
	CardData    []byte //牌
}

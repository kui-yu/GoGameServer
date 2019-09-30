package main

import (
	"fmt"
	"logs"
	"math/rand"
	"time"
)

//癞子
type Lz struct {
	TianLz byte //天癞子
	DiLz   byte //地癞子
}

type CanoutCards struct {
	Ishas    bool         //是否包含癞子
	Canout   []CanOutType //可出牌型（包含癞子时）
	GoutCard GOutCard     //不包含癞子时
}

type CanOutType struct { //可出牌型
	CT       int
	Cards    []byte
	Ptcon    []byte
	Max      byte
	LzBecome []GGameOutLzBecome
}

//癞子卡牌管理器
type Lz_MgrCard struct {
	Lz_Cards      []byte //牌组（未洗牌的牌组，每一把当成一副新的牌，用来洗牌）
	Lz_SourceCard []byte //用来发牌的牌组
	Lz_Lz         Lz     //本局癞子
}

//初始化管理器
func (this *Lz_MgrCard) InitCard() {
	this.Lz_Cards = []byte{}
	this.Lz_SourceCard = []byte{}
}

//为LzCards切片中添加牌元素
func (this *Lz_MgrCard) InitNormalCard() {
	Newcards := []byte{Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1}
	for _, v := range Newcards {
		for j := byte(0); j < 13; j++ {
			this.Lz_Cards = append(this.Lz_Cards, j+v)
		}
	}
}

//发手牌
func (this *Lz_MgrCard) sendCard(num int) []byte {
	result := []byte{}
	//发牌
	result = append([]byte{}, this.Lz_SourceCard[0:num]...)
	//从牌组中减去发出回去的手牌
	this.Lz_SourceCard = append([]byte{}, this.Lz_SourceCard[num:]...)
	//返回发的牌数组
	return result
}

//洗牌 和 制作癞子
func (this *Lz_MgrCard) shuffle() {
	//赋值 Lz_SourceCard变量
	this.Lz_SourceCard = append([]byte{}, this.Lz_Cards...)
	//获取随机对象
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//利用rand.Perm方法获取一个打乱的 下标切片
	perm := r.Perm(len(this.Lz_Cards))
	//通过打乱的下标切片打乱Lz_SourceCards变量，从而达到洗牌的目的
	for i, randIndex := range perm {
		this.Lz_SourceCard[i] = this.Lz_Cards[randIndex]
	}
	//随机出癞子
	LzIndex1 := r.Intn(13) + 1
	for LzIndex1 == Card_King_1 || LzIndex1 == Card_King_2 {
		LzIndex1 = r.Intn(13) + 1
	}
	r2 := rand.New(rand.NewSource(time.Now().UnixNano()))
	LzIndex2 := -1
	for LzIndex1 == LzIndex2 || LzIndex2 == -1 || LzIndex1 == Card_King_1 || LzIndex1 == Card_King_2 {
		LzIndex2 = r2.Intn(13) + 1
	}
	this.Lz_Lz.TianLz = byte(LzIndex1 + 16*r2.Intn(4)) //赋值天癞子
	this.Lz_Lz.DiLz = byte(LzIndex2 + 16*r2.Intn(4))   //赋值地癞子
}

//////////////////////////////////////////

//判断出的牌是否包含癞子,如果有就返回排序过后的癞子集合和普通牌集合
func IsHasLz(cards []byte, lz Lz) (bool, []byte, []byte) {
	isHaslz := false
	ptcon := []byte{}
	lzcon := []byte{}
	for _, v := range cards {
		if GetLogicValue(v) == lz.DiLz || GetLogicValue(v) == lz.TianLz {
			isHaslz = true
			lzcon = append(lzcon, v)
			continue
		}
		ptcon = append(ptcon, v)
	}
	return isHaslz, lzcon, ptcon
}

//判断出牌是否为纯癞子
func IsAllLz(cards []byte, lz Lz) bool {
	var Ptcon []byte = []byte{}
	for _, v := range cards {
		if GetLogicValue(v) != lz.DiLz && GetLogicValue(v) != lz.TianLz {
			Ptcon = append(Ptcon, v)
			continue
		}
	}
	if len(Ptcon) == 0 { //如果普通牌集合长度等于0，则证明为出的牌为纯癞子
		return true
	}
	return false
}

//判断各种牌型是否符合 （癞子）

//对子
func Lz_Check_Duizi(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	canout := CanOutType{
		CT:    CT_DOUBLE,
		Ptcon: ptcon,
		Cards: cards,
	}
	if ptcon[0] < 16 {
		for _, v := range lzcon {
			canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
				Lz:     v,
				Become: GetLogicValue(ptcon[0]),
			})
		}
	}
}

//三个
func Lz_Check_Sange(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	canout := CanOutType{
		CT:    CT_THREE,
		Ptcon: ptcon,
		Cards: cards,
	}
	if GetLogicValue(ptcon[0])-GetLogicValue(ptcon[len(ptcon)-1]) == 0 {
		canout.Max = GetLogicValue(ptcon[0])
		for _, v := range lzcon {
			canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
				Lz:     v,
				Become: GetLogicValue(ptcon[0]),
			})
		}
		*Canout = append(*Canout, canout)
	}
}

//顺子
func Lz_Check_Shunzi(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	var chanum byte = 0
	var becom []byte
	canout2 := CanOutType{
		CT:    CT_SINGLE_CONNECT,
		Ptcon: ptcon,
		Cards: cards,
	}
	//先算出普通牌组中各牌之间总共相差多少
	//最大顺子
	if len(ptcon) != 1 {
		for i := 0; i < len(ptcon)-1; i++ {
			chanum += GetLogicValue(ptcon[i]) - GetLogicValue(ptcon[i+1]) - 1
			pd := GetLogicValue(ptcon[i]) - GetLogicValue(ptcon[i+1]) - 1
			for j := 1; byte(j) <= pd; i++ {
				becom = append(becom, GetLogicValue(ptcon[i+1])+byte(j))
			}
		}
		logs.Debug("判断顺子（癞子），需要几张癞子补缺口：", chanum)
		if int(chanum) > len(lzcon) {
			logs.Debug("普通牌差距太大，顺子判断失败！")
			return
		}

		logs.Debug("判断顺子时癞子切片：", lzcon)
		for i := 0; i < len(becom); i++ {
			canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
				Lz:     lzcon[i],
				Become: becom[i],
			})
		}
		lzcon = lzcon[len(becom):]
		logs.Debug("补完顺子的癞子切片:", lzcon)
		if len(lzcon) > 0 {
			if GetLogicValue(ptcon[0]) == 14 {
				canout2.Max = 14
			}
			var max byte
			//赋值
			for i, v := range lzcon {
				if GetLogicValue(ptcon[0])+byte(i+1) <= 14 {
					max = GetLogicValue(ptcon[0]) + byte(i)
					canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]) + byte(i+1),
					})
				} else {
					var k byte = 1
					for j := i; j < len(lzcon); j++ {
						canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
							Lz:     lzcon[j],
							Become: GetLogicValue(ptcon[len(ptcon)-1]) - k,
						})
						k++
					}
					break
				}
			}
			if canout2.Max != 14 {
				canout2.Max = max
			}

			//如果顺子中的数字为癞子本身则设置为0
			for _, v := range canout2.LzBecome {
				if v.Lz == v.Become {
					v.Become = 0
				}
			}
			*Canout = append(*Canout, canout2)
		} else {
			canout2.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout2)
		}
	} else {
		logs.Debug("只有一张普通牌!!（顺子）")
		if GetLogicValue(ptcon[0]) == 14 {
			canout2.Max = 14
		}
		var max byte
		//赋值
		for i, v := range lzcon {
			if GetLogicValue(ptcon[0])+byte(i+1) <= 14 {
				max = GetLogicValue(ptcon[0]) + byte(i)
				canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[0]) + byte(i+1),
				})
			} else {
				var k byte = 1
				for j := i; j < len(lzcon); j++ {
					canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - k,
					})
					k++
				}
				break
			}
		}
		if canout2.Max != 14 {
			canout2.Max = max
		}
		//如果顺子中的数字为癞子本身则设置为0
		for _, v := range canout2.LzBecome {
			if v.Lz == v.Become {
				v.Become = 0
			}
		}
		*Canout = append(*Canout, canout2)
	}
	//最小顺子
	if len(ptcon) != 1 {
		for i := 0; i < len(ptcon)-1; i++ {
			chanum += GetLogicValue(ptcon[i]) - GetLogicValue(ptcon[i+1]) - 1
			pd := GetLogicValue(ptcon[i]) - GetLogicValue(ptcon[i+1]) - 1
			for j := 1; byte(j) <= pd; i++ {
				becom = append(becom, GetLogicValue(ptcon[i+1])+byte(j))
			}
		}
		logs.Debug("判断顺子（癞子），需要几张癞子补缺口：", chanum)
		if int(chanum) > len(lzcon) {
			logs.Debug("普通牌差距太大，顺子判断失败！")
			return
		}

		logs.Debug("判断顺子时癞子切片：", lzcon)
		for i := 0; i < len(becom); i++ {
			canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
				Lz:     lzcon[i],
				Become: becom[i],
			})
		}
		lzcon = lzcon[len(becom):]
		logs.Debug("补完顺子的癞子切片:", lzcon)
		if len(lzcon) > 0 {
			var max byte
			max = GetLogicValue(ptcon[0])
			//赋值
			for i, v := range lzcon {
				if GetLogicValue(ptcon[len(ptcon)-1])-byte(i+1) >= 3 {
					canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(i+1),
					})
				} else {
					var k byte = 1
					for j := i; j < len(lzcon); j++ {
						canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
							Lz:     lzcon[j],
							Become: GetLogicValue(ptcon[0]) + k,
						})
						max = GetLogicValue(ptcon[0]) + k
						k++
					}
					break
				}
			}
			canout2.Max = max
			//如果顺子中的数字为癞子本身则设置为0
			for _, v := range canout2.LzBecome {
				if v.Lz == v.Become {
					v.Become = 0
				}
			}
			*Canout = append(*Canout, canout2)
		} else {
			canout2.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout2)
		}
	} else {
		logs.Debug("只有一张普通牌!!（顺子）")
		if GetLogicValue(ptcon[0]) == 14 {
			canout2.Max = 14
		}
		var max byte
		//赋值
		for i, v := range lzcon {
			if GetLogicValue(ptcon[len(ptcon)-1])-byte(i+1) >= 3 {
				canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(i+1),
				})
			} else {
				var k byte = 1
				for j := i; j < len(lzcon); j++ {
					canout2.LzBecome = append(canout2.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[0]) + k,
					})
					max = GetLogicValue(ptcon[0]) + k
					k++
				}
				break
			}
		}
		canout2.Max = max
		//如果顺子中的数字为癞子本身则设置为0
		for _, v := range canout2.LzBecome {
			if v.Lz == v.Become {
				v.Become = 0
			}
		}
		*Canout = append(*Canout, canout2)
	}

}

//炸弹
func Lz_Check_ZhaDan(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	canout := CanOutType{
		Cards: cards,
		Ptcon: ptcon,
	}
	if len(cards) == 4 {
		canout.CT = CT_BOMB_FOUR_SOFT
	} else {
		canout.CT = CT_BOME_FOUR_UP
	}
	canout.Max = ptcon[0]
	for _, v := range lzcon {
		canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
			Lz:     v,
			Become: GetLogicValue(ptcon[0]),
		})
	}
	*Canout = append(*Canout, canout)
}

//三带二
func Lz_Check_SanDaiEr(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	canout := CanOutType{
		CT:    CT_THREE_LINE_TAKE_TWO,
		Ptcon: ptcon,
		Cards: cards,
	}
	//一张普通牌
	if len(ptcon) == 1 {
		//最大三带二
		if GetLogicValue(ptcon[0]) != 15 {
			for i, v := range lzcon {
				if i <= 2 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 15,
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[0]),
				})
			}
			canout.Max = 15
			*Canout = append(*Canout, canout)
		} else {
			for i, v := range lzcon {
				if i <= 2 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 3,
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
		}

		//最小三带二
		if GetLogicValue(ptcon[0]) != 3 {
			for i, v := range lzcon {
				if i <= 2 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 3,
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[0]),
				})
			}
			canout.Max = 3
			*Canout = append(*Canout, canout)
		} else {
			for i, v := range lzcon {
				if i <= 2 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 4,
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
		}

	}
	//二张普通牌
	if len(ptcon) == 2 {
		if GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1]) {
			//获取最大的三带二
			if GetLogicValue(ptcon[0]) != 15 {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 15,
					})
				}
				canout.Max = 15
				*Canout = append(*Canout, canout)
			} else {
				for i, v := range lzcon {
					if i == 0 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 3,
					})
				}
				canout.Max = GetLogicValue(ptcon[0])
				*Canout = append(*Canout, canout)
			}
			///////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			///////////////////
			//获取最小的三带二
			if GetLogicValue(ptcon[0]) != 3 {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 3,
					})
				}
				canout.Max = 3
				*Canout = append(*Canout, canout)
			} else {
				for i, v := range lzcon {
					if i == 0 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 4,
					})
				}
				canout.Max = GetLogicValue(ptcon[0])
				*Canout = append(*Canout, canout)
			}
		} else {
			//最大三带二
			for i, v := range lzcon {
				if i <= 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[1]),
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
			//////////////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			//////////////////////////////////////
			//最小三带二
			for i, v := range lzcon {
				if i <= 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[1]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[0]),
				})
			}
			canout.Max = GetLogicValue(ptcon[1])
			*Canout = append(*Canout, canout)
		}
	}
	//三张普通牌
	if len(ptcon) == 3 {
		//判断是否3张一样
		var pd bool = true
		//如果三张一样
		if (GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1])) && (GetLogicValue(ptcon[1]) == GetLogicValue(ptcon[2])) {
			pd = false
			if lzcon[0] == lzcon[1] {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			} else {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     lzcon[0],
					Become: 1,
				})
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     lzcon[1],
					Become: lzcon[0],
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
		}
		//如果两张一样
		if (GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1])) || (GetLogicValue(ptcon[1]) == GetLogicValue(ptcon[2])) && pd {
			//最大值三带二
			if GetLogicValue(ptcon[0]) != GetLogicValue(ptcon[1]) {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 1 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[2]),
					})
				}
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)

			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}

			//最小值三带二
			if GetLogicValue(ptcon[2]) != GetLogicValue(ptcon[1]) {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[2]),
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 1 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[2]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
				}
			}
			canout.Max = GetLogicValue(ptcon[2])
			*Canout = append(*Canout, canout)
		}
	}
	//四张普通牌
	if len(ptcon) == 4 {
		logs.Debug("四张普通牌")
		if GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[2]) || (GetLogicValue(ptcon[1]) == GetLogicValue(ptcon[3])) {
			canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
				Lz:     lzcon[0],
				Become: GetLogicValue(ptcon[2]),
			})
			canout.Max = GetLogicValue(ptcon[2])
			*Canout = append(*Canout, canout)
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
		}

		if GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1]) && (GetLogicValue(ptcon[2]) == GetLogicValue(ptcon[3])) {
			//最大三带二
			canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
				Lz:     lzcon[0],
				Become: GetLogicValue(ptcon[0]),
			})
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			//最小三带二
			canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
				Lz:     lzcon[0],
				Become: GetLogicValue(ptcon[3]),
			})
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
		}
	}
}

//三带一
func Lz_Check_SanDaiYi(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	canout := CanOutType{
		CT:    CT_THREE_LINE_TAKE_ONE,
		Ptcon: ptcon,
		Cards: cards,
	}
	if len(ptcon) == 1 {
		//最大三带一
		if GetLogicValue(ptcon[0]) != 15 {
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 15,
				})
			}
			canout.Max = 15
			*Canout = append(*Canout, canout)
		} else {
			for i, v := range lzcon {
				if i <= 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
		}
		////////////////////////////////
		canout.Max = 0
		canout.LzBecome = []GGameOutLzBecome{}
		////////////////////////////////
		//最小三带一
		if GetLogicValue(ptcon[0]) != 3 {
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 3,
				})
			}
			canout.Max = 3
			*Canout = append(*Canout, canout)
		} else {
			for i, v := range lzcon {
				if i <= 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
		}

	}
	if len(ptcon) == 2 {
		if GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1]) {
			for i, v := range lzcon {
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
		} else {
			//最大三带一
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[0]),
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
			/////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			/////////////////////////////

			//最小三带一
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[1]),
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
		}
	}
	if len(ptcon) == 3 {
		var pd bool = true //用来判断到底3张牌是否全部相等
		if (GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1])) && (GetLogicValue(ptcon[1]) == GetLogicValue(ptcon[2])) {
			pd = false
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
		}

		if (GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1]) || GetLogicValue(ptcon[1]) == GetLogicValue(ptcon[2])) && pd {
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[1]),
				})
			}
			canout.Max = GetLogicValue(ptcon[1])
			*Canout = append(*Canout, canout)
		}
	}
}

//双顺
func Lz_Check_ShuangShun(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	var chanum int = 0 //需要几张癞子
	var becom []byte   //需要变成牌的集合
	canout := CanOutType{
		Cards: cards,
		CT:    CT_DOUBLE_CONNECT,
		Ptcon: ptcon,
	}
	//判断用户出的普通牌 如果需要将其补全需要多少的癞子
	//现将未组合成 一对的进行补全，查看需要多少癞子
	var Doublecards []byte //补全后的数组
	for i := 0; i < len(ptcon); i++ {
		var count int = 0
		for _, v := range ptcon {
			if GetLogicValue(v) == GetLogicValue(ptcon[i]) {
				count++
			}
		}
		if count > 2 {
			logs.Debug("超过两个，无法形成双顺")
			return
		}
		if count == 1 {
			chanum++
			becom = append(becom, GetLogicValue(ptcon[i]))
			Doublecards = append(Doublecards, ptcon[i])
		}
	}

	if chanum > len(lzcon) {
		logs.Debug("双顺判断中，癞子不足以凑对子")
		return
	}
	Doublecards = append(Doublecards, ptcon...)
	Doublecards = Sort(Doublecards)
	if len(Doublecards) > 2 {
		for i := 0; i < len(Doublecards)-2; i += 2 {
			num := 2 * (GetLogicValue(Doublecards[i]) - GetLogicValue(Doublecards[i+2]) - 1)
			if num > 0 {
				chanum += int(num)
				for j := 0; j < 2; j++ {
					becom = append(becom, GetLogicValue(ptcon[i]))
				}
			}
		}
	}
	if chanum > len(lzcon) {
		logs.Debug("双顺判断中，癞子不足以凑足双顺")
	}
	fmt.Println("双顺：替换之前的癞子：", lzcon)
	for i := 0; i < len(becom); i++ {
		canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
			Lz:     lzcon[i],
			Become: becom[i],
		})
	}
	lzcon = lzcon[len(becom):]
	fmt.Println("双顺：替换之后的癞子：", lzcon)
	//保留现在的Lzbecom
	tempLzbecom := canout.LzBecome
	if len(lzcon) > 0 {
		//最大值
		var max byte = 0
		var k int = 1
		for i := 0; i < len(lzcon)-1; i += 2 {
			if GetLogicValue(ptcon[0])+byte(i+1) < 14 {
				for j := 0; j < 2; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[i+j],
						Become: GetLogicValue(ptcon[0]) + byte(i+1),
					})
					max = GetLogicValue(ptcon[0]) + byte(i+1)
				}
			} else {
				for j := 0; j < 2; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[i+j],
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(k),
					})
				}
				k++
			}
		}
		if max != 0 {
			canout.Max = max
		} else {
			canout.Max = GetLogicValue(ptcon[0])
		}
		*Canout = append(*Canout, canout)
		////////////////////////////////////////////////
		canout.Max = 0
		max = 0
		k = 1
		canout.LzBecome = tempLzbecom
		////////////////////////////////////////////////

		//最小值
		for i := 0; i < len(lzcon)-1; i += 2 {
			if GetLogicValue(ptcon[len(ptcon)-1])-byte(i+1) > 3 {
				for j := 0; j < 2; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[i+j],
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(i+1),
					})
				}
			} else {
				for j := 0; j < 2; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[i+j],
						Become: GetLogicValue(ptcon[0]) + byte(k),
					})
					max = GetLogicValue(ptcon[0]) + byte(k)
				}
				k++
			}
		}
		if max != 0 {
			canout.Max = max
		} else {
			canout.Max = GetLogicValue(ptcon[0])
		}
		for _, v := range canout.LzBecome {
			if v.Lz == v.Become {
				v.Become = 0
			}
		}
		*Canout = append(*Canout, canout)

	} else {
		canout.Max = GetLogicValue(ptcon[0])
		for _, v := range canout.LzBecome {
			if v.Lz == v.Become {
				v.Become = 0
			}
		}
		*Canout = append(*Canout, canout)
	}

}

//飞机（可看成三顺）
func Lz_Check_FeiJi(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	canout := CanOutType{
		Cards: cards,
		Ptcon: ptcon,
		CT:    CT_AIRCRAFT,
	}
	for _, v := range ptcon {
		if GetLogicValue(v) > 14 {
			logs.Error("三顺判断： 牌中拥有超过 A的牌，无法组成三顺！")
			return
		}
	}
	single := make(map[byte]int)
	duizi := make(map[byte]int)
	sange := make(map[byte]int)
	sige := make(map[byte]int)
	numCout := make(map[byte]int)
	for v := 17; v <= 3; v++ {
		for _, v1 := range ptcon {
			if byte(v) == GetLogicValue(v1) {
				numCout[byte(v)] += 1
			}
		}
	}
	for key, v := range numCout {
		if v == 1 {
			single[key] += 1
		}
		if v == 2 {
			duizi[key] += 1
		}
		if v == 3 {
			sange[key] += 1
		}
		if v == 4 {
			sige[key] += 1
		}
	}
	if len(sige) == 1 {
		logs.Error("三顺中出现了四个相同，无法形成三顺！")
		return
	}
	//先将用户出的普通牌 未成3张 补齐
	//获取需要补齐的牌
	var needpolishing []byte
	if len(single) != 0 {
		for k := range single {
			for i := 0; i < 2; i++ {
				needpolishing = append(needpolishing, k)
			}
		}
	}
	if len(duizi) != 0 {
		for k := range duizi {
			needpolishing = append(needpolishing, k)
		}
	}
	if len(needpolishing) > len(lzcon) {
		logs.Error("三顺： 癞子无法补齐三张！")
		return
	}
	for i, v := range needpolishing {
		canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
			Lz:     lzcon[i],
			Become: v,
		})
	}
	lzcon = lzcon[len(needpolishing)+1:]
	var needpolishingnum []byte
	if len(ptcon) > 1 {
		for i := 0; i < len(ptcon)-1; i++ {
			cha := GetLogicValue(ptcon[i]) - GetLogicValue(ptcon[i+1]) - 1
			if cha > 0 {
				for j := 1; j < int(cha+1); i++ {
					needpolishing = append(needpolishing, GetLogicValue(ptcon[i+1])+byte(j))
				}
			}
		}
	}
	if len(lzcon) < len(needpolishingnum) {
		logs.Debug("三顺中第二补齐失败，牌之间差距太大")
		return
	}
	for i, v := range needpolishingnum {
		canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
			Lz:     lzcon[i],
			Become: v,
		})
	}
	lzcon = lzcon[len(needpolishingnum)+1:]
	lb := canout.LzBecome //记录当前 Lzbecome
	if len(lzcon) > 0 {
		//最大值
		if byte(len(lzcon)/3)+GetLogicValue(ptcon[0]) <= 14 {
			canout.Max = GetLogicValue(ptcon[0]) + byte(len(lzcon)/3)
			for i := 1; i < len(lzcon)/3+1; i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[0]) + byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
		} else {
			canout.Max = 14
			cha := 14 - GetLogicValue(ptcon[len(ptcon)-1]) - 1
			for i := 1; i < int(cha+1); i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[0]) + byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
			for i := 1; i < len(lzcon)/3+1; i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(i),
					})
				}
				lzcon = lzcon[4:]
			}

		}
		*Canout = append(*Canout, canout)
		//////////////////
		canout.Max = 0
		canout.LzBecome = lb
		/////////////////
		//最小值
		if GetLogicValue(ptcon[0])-byte(len(lzcon)/3) >= 3 {
			canout.Max = GetLogicValue(ptcon[0]) - 1
			for i := 1; i < len(lzcon)/3+1; i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
		} else {
			cha := GetLogicValue(ptcon[len(ptcon)-1]) - 3 - 1
			for i := 1; i < int(cha+1); i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
			canout.Max = GetLogicValue(ptcon[0]) + byte(len(lzcon)/3)

			for i := 1; i < len(lzcon)/3+1; i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[0]) + byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
		}
		*Canout = append(*Canout, canout)
	} else {
		canout.Max = GetLogicValue(ptcon[0])
		*Canout = append(*Canout, canout)
	}
}

//四带二单
func Lz_Check_SiDaiErDan(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	canout := CanOutType{
		Cards: cards,
		Ptcon: ptcon,
		CT:    CT_FOUR_LINE_TAKE_ONE,
	}
	if len(ptcon) == 1 {
		var max byte = 0
		//最大值
		if GetLogicValue(ptcon[0]) != 15 {
			max = 15
			for i, v := range lzcon {
				if i < 4 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 15,
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
		} else {
			max = GetLogicValue(ptcon[0])
			for i, v := range lzcon {
				if i < 3 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
		}
		canout.Max = max
		*Canout = append(*Canout, canout)
		/////////////////////////////////////////
		canout.LzBecome = []GGameOutLzBecome{}
		canout.Max = 0
		////////////////////////////////////////////
		//最小值
		if GetLogicValue(ptcon[0]) != 3 {
			max = 3
			for i, v := range lzcon {
				if i < 4 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 3,
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
		} else {
			max = GetLogicValue(ptcon[0])
			for i, v := range lzcon {
				if i < 3 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
		}
		canout.Max = max
		*Canout = append(*Canout, canout)

	}
	if len(ptcon) == 2 {
		var max byte
		var ishas3 bool = false
		var ishas15 bool = false
		for _, v := range ptcon {
			if GetLogicValue(v) == 3 {
				ishas3 = true
			}
			if GetLogicValue(v) == 4 {
				ishas3 = true
			}
		}
		max = 15
		//最大值
		if !ishas15 {
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 15,
				})
			}
		} else {
			if GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1]) {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			}
		}
		canout.Max = max
		*Canout = append(*Canout, canout)
		/////////////////////////////////
		canout.Max = 0
		canout.LzBecome = []GGameOutLzBecome{}
		////////////////////////////////////
		//最小值
		max = 3
		if !ishas3 {
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 3,
				})
			}
		} else {
			if GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1]) {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[len(ptcon)-1]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			}
		}
		canout.Max = max
		*Canout = append(*Canout, canout)
	}
	if len(ptcon) == 3 {
		var pd bool = true
		//三张一样
		if GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1]) && GetLogicValue(ptcon[1]) == GetLogicValue(ptcon[2]) {
			pd = false
			canout.Max = GetLogicValue(ptcon[0])
			for i, v := range lzcon {
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
			*Canout = append(*Canout, canout)
		}
		//两张一样
		if (GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1]) || GetLogicValue(ptcon[1]) == GetLogicValue(ptcon[2])) && pd {
			var dif byte
			if GetLogicValue(ptcon[0]) != GetLogicValue(ptcon[1]) {
				dif = ptcon[0]
			} else {
				dif = ptcon[2]
			}
			//最大值
			if dif != ptcon[0] {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			} else {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
				}
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
			////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			/////////////////////////////////
			canout.Max = GetLogicValue(ptcon[len(ptcon)-1])

			//最小值
			if dif != ptcon[len(ptcon)-1] {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[len(ptcon)-1]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			} else {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[len(ptcon)-1]),
					})
				}
			}
			*Canout = append(*Canout, canout)
		}
		//都不一样
		if GetLogicValue(ptcon[0]) != GetLogicValue(ptcon[1]) && GetLogicValue(ptcon[1]) != GetLogicValue(ptcon[2]) {
			//最大值
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[0]),
				})
			}
			canout.Max = GetLogicValue(ptcon[0])
			*Canout = append(*Canout, canout)
			/////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			/////////////////////////////

			//最小值
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[len(ptcon)-1]),
				})
			}
			canout.Max = GetLogicValue(ptcon[len(ptcon)-1])
			*Canout = append(*Canout, canout)
		}
	}
	if len(ptcon) == 4 {
		var duizi map[int]int
		var sange map[int]int
		var sige map[int]int
		var numCout map[int]int
		for v := 3; v <= 17; v++ {
			for _, v1 := range ptcon {
				if byte(v) == GetLogicValue(v1) {
					numCout[v]++
				}
			}
		}
		for key, v := range numCout {
			if v == 2 {
				duizi[key] += 1
			}
			if v == 3 {
				sange[key] += 1
			}
			if v == 4 {
				sige[key] += 1
			}
		}
		//一对
		if len(duizi) == 1 {
			var key int
			for k := range duizi {
				key = k
			}
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: byte(key),
				})
			}
			canout.Max = byte(key)
			*Canout = append(*Canout, canout)
		}
		//两对
		if len(duizi) == 2 {
			//最大值
			canout.Max = GetLogicValue(ptcon[0])
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[0]),
				})
			}
			*Canout = append(*Canout, canout)
			//////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			////////////////////////
			//最小值
			canout.Max = GetLogicValue(ptcon[len(ptcon)-1])
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[len(ptcon)-1]),
				})
			}
			*Canout = append(*Canout, canout)
		}
		//三一
		if len(sange) == 1 {
			var key int
			for k := range sange {
				key = k
			}
			canout.Max = byte(key)
			for i, v := range lzcon {
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: byte(key),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
			*Canout = append(*Canout, canout)
		}
		//四个
		if len(sige) == 1 {
			for k := range sange {
				canout.Max = byte(k)
			}
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
			*Canout = append(*Canout, canout)
		}
	}
	if len(ptcon) == 5 {
		var sange map[int]int
		var sige map[int]int
		var numCout map[int]int
		for v := 3; v <= 17; v++ {
			for _, v1 := range ptcon {
				if byte(v) == GetLogicValue(v1) {
					numCout[v]++
				}
			}
		}
		for key, v := range numCout {
			if v == 3 {
				sange[key]++
			}
			if v == 4 {
				sige[key]++
			}
		}
		if len(sange) == 1 {
			var key int
			for k := range sange {
				key = k
			}
			canout.Max = byte(key)
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: byte(key),
				})
			}
			*Canout = append(*Canout, canout)
		}
		if len(sige) == 1 {
			for k := range sige {
				canout.Max = byte(k)
			}
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
		}
	}
}

//四带二双
func Lz_Check_SiDaiErShuang(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	canout := CanOutType{
		Cards: cards,
		Ptcon: ptcon,
		CT:    CT_FOUR_LINE_TAKE_TWO,
	}
	if len(ptcon) == 1 {
		//最大值
		canout.Max = 15
		if GetLogicValue(ptcon[0]) != 15 {
			for i, v := range lzcon {
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				if i < 3 {
					var becom byte
					if GetLogicValue(ptcon[0]) > 4 {
						becom = GetLogicValue(ptcon[0]) - 1
					} else {
						becom = GetLogicValue(ptcon[0]) + 1
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: becom,
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 15,
				})
			}
		} else {
			for i, v := range lzcon {
				if i < 3 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				if i < 6 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 4,
					})
					continue
				}
				if i < 8 {

					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 3,
					})
					continue
				}
			}
		}
		*Canout = append(*Canout, canout)
		///////////////////////
		canout.Max = 0
		canout.LzBecome = []GGameOutLzBecome{}
		///////////////////////
		canout.Max = 3
		if GetLogicValue(ptcon[0]) != 3 {
			for i, v := range lzcon {
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				if i < 3 {
					var becom byte
					if GetLogicValue(ptcon[0]) > 4 {
						becom = GetLogicValue(ptcon[0]) - 1
					} else {
						becom = GetLogicValue(ptcon[0]) + 1
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: becom,
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 3,
				})
			}
		} else {
			for i, v := range lzcon {
				if i < 3 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				if i < 6 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 4,
					})
					continue
				}
				if i < 8 {

					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 5,
					})
					continue
				}
			}
		}
		*Canout = append(*Canout, canout)
	}
	if len(ptcon) == 2 {
		if GetLogicValue(ptcon[0]) == GetLogicValue(ptcon[1]) {
			//最大值
			canout.Max = 15
			if GetLogicValue(ptcon[0]) != 15 {
				var becom byte
				if GetLogicValue(ptcon[0]) > 4 {
					becom = GetLogicValue(ptcon[0]) - 1
				} else {
					becom = GetLogicValue(ptcon[0]) + 1
				}
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: becom,
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 15,
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}

					if i < 4 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: 3,
						})
					}
					if i < 6 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: 4,
						})
					}
				}
			}
			*Canout = append(*Canout, canout)
			//////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			//////////////////////////
			//最小值
			if GetLogicValue(ptcon[0]) != 3 {
				var becom byte
				if GetLogicValue(ptcon[0]) > 4 {
					becom = GetLogicValue(ptcon[0]) - 1
				} else {
					becom = GetLogicValue(ptcon[0]) + 1
				}
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: becom,
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 3,
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}

					if i < 4 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: 4,
						})
					}
					if i < 6 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: 5,
						})
					}
				}
			}
			*Canout = append(*Canout, canout)
		} else {

		}
	}
	if len(ptcon) == 3 {
		var single map[int]int
		var duizi map[int]int
		var sange map[int]int
		var sige map[int]int
		var numCout map[int]int
		for v := 3; v <= 17; v++ {
			for _, v1 := range ptcon {
				if byte(v) == GetLogicValue(v1) {
					numCout[v] += 1
				}
			}
		}
		for key, v := range numCout {
			if v == 1 {
				single[key] += 1
			}
			if v == 2 {
				duizi[key] += 1
			}
			if v == 3 {
				sange[key] += 1
			}
			if v == 4 {
				sige[key] += 1
			}
		}
		if len(single) == 3 {
			//最大值
			canout.Max = GetLogicValue(ptcon[0])
			for i, v := range lzcon {
				if i < 3 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
				if i < 4 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[1]),
					})
					continue
				}
				if i < 5 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[2]),
					})
					continue
				}
			}
			*Canout = append(*Canout, canout)
			///////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			//////////////////////////////////
			//最小值
			canout.Max = GetLogicValue(ptcon[len(ptcon)-1])
			for i, v := range lzcon {
				if i < 3 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[len(ptcon)-1]),
					})
					continue
				}
				if i < 4 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[1]),
					})
					continue
				}
				if i < 5 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
					continue
				}
			}
			*Canout = append(*Canout, canout)
		}
		if len(duizi) == 1 {
			//最大值
			canout.Max = 15
			var key int
			for k := range duizi {
				key = k
			}
			var become1 byte
			var become2 byte
			for _, v := range ptcon {
				if v != byte(key) {
					become1 = GetLogicValue(v)
					if become1 > 3 {
						become2 = become1 - 1
					} else {
						become2 = become1 + 1
					}
				}
			}
			if GetLogicValue(ptcon[0]) != 15 {
				for i, v := range lzcon {
					if i < 1 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: become1,
						})
						continue
					}

					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 15,
					})
				}
			} else {

				for i, v := range lzcon {
					if i < 1 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: become1,
						})
						continue
					}
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: become2,
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
				}
			}
			*Canout = append(*Canout, canout)
			/////////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			/////////////////////////////////
			//最小值
			canout.Max = 3
			if GetLogicValue(ptcon[0]) != 3 {
				for i, v := range lzcon {
					if i < 1 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: become1,
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 3,
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 1 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: become1,
						})
						continue
					}
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: become2,
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[len(ptcon)-1]),
					})
				}
			}
			*Canout = append(*Canout, canout)
		}
		if len(sange) == 1 {
			canout.Max = GetLogicValue(ptcon[0])
			var become1 byte
			var become2 byte
			if GetLogicValue(ptcon[0]) > 4 {
				become1 = GetLogicValue(ptcon[0]) - 1
				become2 = GetLogicValue(ptcon[0]) - 2
			} else {
				become1 = GetLogicValue(ptcon[0]) + 1
				become2 = GetLogicValue(ptcon[0]) + 2
			}
			for i, v := range lzcon {
				if i < 2 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: become1,
					})
					continue
				}
				if i < 4 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: become2,
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[0]),
				})
			}
			*Canout = append(*Canout, canout)
		}
	}
	if len(ptcon) == 4 {
		var single map[int]int
		var duizi map[int]int
		var sange map[int]int
		var sige map[int]int
		var numCout map[int]int
		for v := 3; v <= 17; v++ {
			for _, v1 := range ptcon {
				if byte(v) == GetLogicValue(v1) {
					numCout[v] += 1
				}
			}
		}
		for key, v := range numCout {
			if v == 1 {
				single[key] += 1
			}
			if v == 2 {
				duizi[key] += 1
			}
			if v == 3 {
				sange[key] += 1
			}
			if v == 4 {
				sige[key] += 1
			}
		}
		if len(duizi) == 2 {
			var become byte
			//最大值
			canout.Max = 15
			if GetLogicValue(ptcon[3]) > 3 {
				become = GetLogicValue(ptcon[3]) - 1
			} else {
				become = GetLogicValue(ptcon[3]) + 1
			}
			if GetLogicValue(ptcon[0]) != 15 {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 15,
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: become,
					})
				}
			}
			*Canout = append(*Canout, canout)
			//////////////////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			/////////////////////////////////////////
			//最小值
			canout.Max = 15
			if GetLogicValue(ptcon[0]) < 15 {
				become = GetLogicValue(ptcon[0]) + 1
			} else {
				become = GetLogicValue(ptcon[0]) - 1
			}
			if GetLogicValue(ptcon[3]) != 15 {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 3,
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[3]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: become,
					})
				}
			}
			*Canout = append(*Canout, canout)
		}
		if len(sige) == 1 {
			var become byte
			var become1 byte
			if GetLogicValue(ptcon[0]) > 4 {
				become = GetLogicValue(ptcon[0]) - 1
				become1 = GetLogicValue(ptcon[0]) - 2
			} else {
				become = GetLogicValue(ptcon[0]) - 1
				become1 = GetLogicValue(ptcon[0]) - 2
			}
			for i, v := range lzcon {
				if i < 2 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: become,
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: become1,
				})
			}
			*Canout = append(*Canout, canout)

		}
		if len(duizi) == 1 {
			var key int
			for k := range duizi {
				key = k
			}
			var indexCount []int
			for i, v := range ptcon {
				if byte(key) == GetLogicValue(v) {
					indexCount = append(indexCount, i)
				}
			}

			for _, i := range indexCount {
				ptcon = append(ptcon[:i], ptcon[i+1:]...)
			}
			//最大值
			canout.Max = GetLogicValue(canout.Ptcon[0])
			if GetLogicValue(canout.Ptcon[0]) == byte(key) {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: byte(key),
						})
						continue
					}
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[1]),
					})
				}

			} else {
				for i, v := range lzcon {
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[1]),
					})
				}
			}
			///////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			/////////////////////////////////
			//最小值
			canout.Max = GetLogicValue(canout.Ptcon[len(ptcon)-1])
			if GetLogicValue(canout.Ptcon[len(ptcon)-1]) == byte(key) {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: byte(key),
						})
						continue
					}
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[1]),
					})
				}

			} else {
				for i, v := range lzcon {
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[len(ptcon)-1]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
				}
			}
		}
	}
	if len(ptcon) == 5 {
		var singe map[int]int
		var duizi map[int]int
		var sange map[int]int
		var sige map[int]int
		var numCout map[int]int
		for v := 3; v <= 17; v++ {
			for _, v1 := range ptcon {
				if byte(v) == GetLogicValue(v1) {
					numCout[v] += 1
				}
			}
		}
		for key, v := range numCout {
			if v == 1 {
				singe[key] += 1
			}
			if v == 2 {
				duizi[key] += 1
			}
			if v == 3 {
				sange[key] += 1
			}
			if v == 4 {
				sige[key] += 1
			}
		}

		if len(sige) == 1 {
			var key int
			for k := range sige {
				key = k
			}
			canout.Max = byte(key)
			if GetLogicValue(ptcon[0]) == byte(key) {
				for i, v := range lzcon {
					if i < 1 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[len(ptcon)-1]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(lzcon[0]),
					})
				}
			} else {
				for i, v := range lzcon {
					if i < 1 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(lzcon[0]),
					})
				}
			}

		}
		if len(sange) == 1 && len(duizi) == 1 {
			var key int
			for k := range sange {
				key = k
			}
			for i, v := range lzcon {
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: byte(key),
					})
					continue
				}

				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(lzcon[0]),
				})
			}
		}
		if len(sange) == 1 && len(singe) == 2 {
			var keySange int
			var keySinge []int
			for k := range singe {
				keySinge = append(keySinge, k)
			}
			for k := range sange {
				keySange = k
			}
			canout.Max = GetLogicValue(ptcon[0])
			for i, v := range lzcon {
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: byte(keySange),
					})
					continue
				}
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: byte(keySinge[0]),
					})
					continue
				}
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: byte(keySinge[1]),
					})
					continue
				}
			}
			*Canout = append(*Canout, canout)
		}
		if len(duizi) == 2 {
			var key []int
			for k := range duizi {
				key = append(key, k)
			}
			var pd bool = true
			var dif byte
			for _, v := range ptcon {
				var pd1 bool = true
				for _, v1 := range key {
					if byte(v1) == GetLogicValue(v) {
						pd1 = false
						continue
					}
					dif = v
				}
				if pd1 {
					break
				}
			}

			//最大值
			canout.Max = GetLogicValue(ptcon[0])
			for _, v := range key {
				if byte(v) == GetLogicValue(ptcon[0]) {
					pd = true
				}
			}
			if pd {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[0]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(dif),
					})
				}
			} else {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[0]),
					})
				}
			}
			*Canout = append(*Canout, canout)
			////////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			//////////////////////////////

			canout.Max = GetLogicValue(ptcon[len(ptcon)-1])
			for _, v := range key {
				if byte(v) == GetLogicValue(ptcon[len(ptcon)-1]) {
					pd = true
				}
			}
			if pd {
				for i, v := range lzcon {
					if i < 2 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: GetLogicValue(ptcon[len(ptcon)-1]),
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(dif),
					})
				}
			} else {
				for _, v := range lzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: GetLogicValue(ptcon[len(ptcon)-1]),
					})
				}
			}
			*Canout = append(*Canout, canout)
		}
	}
	if len(ptcon) == 6 {
		var singe map[int]int
		var duizi map[int]int
		var sange map[int]int
		var sige map[int]int
		var numCout map[int]int
		for v := 3; v <= 17; v++ {
			for _, v1 := range ptcon {
				if byte(v) == GetLogicValue(v1) {
					numCout[v] += 1
				}
			}
		}
		for key, v := range numCout {
			if v == 1 {
				singe[key] += 1
			}
			if v == 2 {
				duizi[key] += 1
			}
			if v == 3 {
				sange[key] += 1
			}
			if v == 4 {
				sige[key] += 1
			}
		}
		if len(sige) == 1 && len(duizi) == 1 {
			canout.Max = GetLogicValue(ptcon[0])
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(lzcon[0]),
				})
			}
			*Canout = append(*Canout, canout)
		}
		if len(sige) == 1 && len(singe) == 2 {
			var key []int
			for k := range singe {
				key = append(key, k)
			}
			for i, v := range lzcon {
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: byte(key[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: byte(key[1]),
				})
			}
			*Canout = append(*Canout, canout)
		}
		if len(duizi) == 3 {
			//最大值
			canout.Max = GetLogicValue(ptcon[0])
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[0]),
				})
			}
			*Canout = append(*Canout, canout)
			////////////////////////////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			///////////////////////////////////
			//最小值
			canout.Max = GetLogicValue(ptcon[len(ptcon)-1])
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: GetLogicValue(ptcon[len(ptcon)-1]),
				})
			}
			*Canout = append(*Canout, canout)
		}
		if len(sange) == 1 && len(duizi) == 1 && len(singe) == 1 {
			var key []int
			for k := range sange {
				key = append(key, k)
			}
			for k := range singe {
				key = append(key, k)
			}
			for i, v := range lzcon {
				if i < 1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: byte(key[0]),
					})
					continue
				}
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: byte(key[1]),
				})
			}
		}
	}
	if len(ptcon) == 7 {
		var singe map[int]int
		var duizi map[int]int
		var sange map[int]int
		var sige map[int]int
		var numCout map[int]int
		for v := 3; v <= 17; v++ {
			for _, v1 := range ptcon {
				if byte(v) == GetLogicValue(v1) {
					numCout[v] += 1
				}
			}
		}
		for key, v := range numCout {
			if v == 1 {
				singe[key] += 1
			}
			if v == 2 {
				duizi[key] += 1
			}
			if v == 3 {
				sange[key] += 1
			}
			if v == 4 {
				sige[key] += 1
			}
		}
		if len(sige) == 1 && len(duizi) == 1 && len(singe) == 1 {
			var key int
			for k := range singe {
				key = k
			}
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: byte(key),
				})
			}
			for k := range sige {
				canout.Max = byte(k)
			}
			*Canout = append(*Canout, canout)
		}
		if len(sange) == 1 && len(duizi) == 2 {
			var key int
			for k := range sange {
				key = k
			}
			canout.Max = byte(key)
			for _, v := range lzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: byte(key),
				})
			}
			*Canout = append(*Canout, canout)
		}
	}
}

//飞机 鉴定三顺 (只取最大值)
func checkFeiji_sanShun_Max(cards []byte, ptcon []byte, lzcon []byte, canout *CanOutType) {
	for _, v := range ptcon {
		if GetLogicValue(v) > 14 {
			logs.Error("三顺判断： 牌中拥有超过 A的牌，无法组成三顺！")
			return
		}
	}
	single := make(map[byte]int)
	duizi := make(map[byte]int)
	sange := make(map[byte]int)
	sige := make(map[byte]int)
	numCout := make(map[byte]int)
	for v := 17; v <= 3; v++ {
		for _, v1 := range ptcon {
			if byte(v) == GetLogicValue(v1) {
				numCout[byte(v)] += 1
			}
		}
	}
	for key, v := range numCout {
		if v == 1 {
			single[key] += 1
		}
		if v == 2 {
			duizi[key] += 1
		}
		if v == 3 {
			sange[key] += 1
		}
		if v == 4 {
			sige[key] += 1
		}
	}
	if len(sige) == 1 {
		logs.Error("三顺中出现了四个相同，无法形成三顺！")
		return
	}
	//先将用户出的普通牌 未成3张 补齐
	//获取需要补齐的牌
	var needpolishing []byte
	if len(single) != 0 {
		for k := range single {
			for i := 0; i < 2; i++ {
				needpolishing = append(needpolishing, k)
			}
		}
	}
	if len(duizi) != 0 {
		for k := range duizi {
			needpolishing = append(needpolishing, k)
		}
	}
	if len(needpolishing) > len(lzcon) {
		logs.Error("三顺： 癞子无法补齐三张！")
		return
	}
	for i, v := range needpolishing {
		canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
			Lz:     lzcon[i],
			Become: v,
		})
	}
	lzcon = lzcon[len(needpolishing)+1:]
	var needpolishingnum []byte
	if len(ptcon) > 1 {
		for i := 0; i < len(ptcon)-1; i++ {
			cha := GetLogicValue(ptcon[i]) - GetLogicValue(ptcon[i+1]) - 1
			if cha > 0 {
				for j := 1; j < int(cha+1); i++ {
					needpolishing = append(needpolishing, GetLogicValue(ptcon[i+1])+byte(j))
				}
			}
		}
	}
	if len(lzcon) < len(needpolishingnum) {
		logs.Debug("三顺中第二补齐失败，牌之间差距太大")
		return
	}
	for i, v := range needpolishingnum {
		canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
			Lz:     lzcon[i],
			Become: v,
		})
	}
	lzcon = lzcon[len(needpolishingnum)+1:]
	if len(lzcon) > 0 {
		//最大值
		if byte(len(lzcon)/3)+GetLogicValue(ptcon[0]) <= 14 {
			canout.Max = GetLogicValue(ptcon[0]) + byte(len(lzcon)/3)
			for i := 1; i < len(lzcon)/3+1; i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[0]) + byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
		} else {
			canout.Max = 14
			cha := 14 - GetLogicValue(ptcon[len(ptcon)-1]) - 1
			for i := 1; i < int(cha+1); i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[0]) + byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
			for i := 1; i < len(lzcon)/3+1; i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
		}
	} else {
		canout.Max = GetLogicValue(ptcon[0])
	}
}

//只取最小值
func checkFeiji_sanShun_Min(cards []byte, ptcon []byte, lzcon []byte, canout *CanOutType) {
	for _, v := range ptcon {
		if GetLogicValue(v) > 14 {
			logs.Error("三顺判断： 牌中拥有超过 A的牌，无法组成三顺！--飞机")
			return
		}
	}
	single := make(map[byte]int)
	duizi := make(map[byte]int)
	sange := make(map[byte]int)
	sige := make(map[byte]int)
	numCout := make(map[byte]int)
	for v := 17; v <= 3; v++ {
		for _, v1 := range ptcon {
			if byte(v) == GetLogicValue(v1) {
				numCout[byte(v)] += 1
			}
		}
	}
	for key, v := range numCout {
		if v == 1 {
			single[key] += 1
		}
		if v == 2 {
			duizi[key] += 1
		}
		if v == 3 {
			sange[key] += 1
		}
		if v == 4 {
			sige[key] += 1
		}
	}
	if len(sige) == 1 {
		logs.Error("三顺中出现了四个相同，无法形成三顺！")
		return
	}
	//先将用户出的普通牌 未成3张 补齐
	//获取需要补齐的牌
	var needpolishing []byte
	if len(single) != 0 {
		for k := range single {
			for i := 0; i < 2; i++ {
				needpolishing = append(needpolishing, k)
			}
		}
	}
	if len(duizi) != 0 {
		for k := range duizi {
			needpolishing = append(needpolishing, k)
		}
	}
	if len(needpolishing) > len(lzcon) {
		logs.Error("三顺： 癞子无法补齐三张！")
		return
	}
	for i, v := range needpolishing {
		canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
			Lz:     lzcon[i],
			Become: v,
		})
	}
	lzcon = lzcon[len(needpolishing)+1:]
	var needpolishingnum []byte
	if len(ptcon) > 1 {
		for i := 0; i < len(ptcon)-1; i++ {
			cha := GetLogicValue(ptcon[i]) - GetLogicValue(ptcon[i+1]) - 1
			if cha > 0 {
				for j := 1; j < int(cha+1); i++ {
					needpolishing = append(needpolishing, GetLogicValue(ptcon[i+1])+byte(j))
				}
			}
		}
	}
	if len(lzcon) < len(needpolishingnum) {
		logs.Debug("三顺中第二补齐失败，牌之间差距太大")
		return
	}
	for i, v := range needpolishingnum {
		canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
			Lz:     lzcon[i],
			Become: v,
		})
	}
	lzcon = lzcon[len(needpolishingnum)+1:]
	if len(lzcon) > 0 {
		//最小值
		if GetLogicValue(ptcon[0])-byte(len(lzcon)/3) >= 3 {
			canout.Max = GetLogicValue(ptcon[0]) - 1
			for i := 1; i < len(lzcon)/3+1; i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
		} else {
			cha := GetLogicValue(ptcon[len(ptcon)-1]) - 3 - 1
			for i := 1; i < int(cha+1); i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[len(ptcon)-1]) - byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
			canout.Max = GetLogicValue(ptcon[0]) + byte(len(lzcon)/3)

			for i := 1; i < len(lzcon)/3+1; i++ {
				for j := 0; j < 3; j++ {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     lzcon[j],
						Become: GetLogicValue(ptcon[0]) + byte(i),
					})
				}
				lzcon = lzcon[4:]
			}
		}
	} else {
		canout.Max = GetLogicValue(ptcon[0])
	}
}

//飞机带单
func Lz_Check_FeiJiDaiErDan(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	canout := CanOutType{
		Cards: cards,
		Ptcon: ptcon,
		CT:    CT_AIRCRAFT_ONE,
	}

	//根据总牌数的长度 算出需要飞机的长度 ，在根据 lzcon 能否凑成飞机 判断该出牌类型是否合法
	ptlength := len(cards) / 4
	fjlenght := len(cards) - ptlength
	needpolish := fjlenght - len(lzcon)
	if len(lzcon) >= fjlenght {
		if len(lzcon) > fjlenght {
			//最大值
			var card []byte
			for _, v := range ptcon {
				if v == 14 && v == 13 {
					card = append(card, v)
				}
			}
			morelzcon := lzcon[:len(lzcon)-(fjlenght-len(lzcon))] //取出多余的癞子
			for _, v := range morelzcon {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
			lesslzcon := lzcon[len(morelzcon):] //取出刚刚好的飞机癞子
			if len(card) != 0 {
				Sort(card)
				newlzcon := lesslzcon[:len(lzcon)-len(card)] //使用普通牌 将癞子牌替换出来
				lastlzcon := lesslzcon[len(newlzcon):]       //多余出来得癞子 再次变成自己本身
				checkFeiji_sanShun_Max(card, card, newlzcon, &canout)
				for _, v := range lastlzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			} else {
				canout.Max = 14
				for i, v := range lesslzcon {
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: 14,
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 13,
					})
				}
			}
			*Canout = append(*Canout, canout)
			////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			////////
			//最小值
			var card1 []byte
			for _, v := range ptcon {
				if v == 3 && v == 4 {
					card1 = append(card1, v)
				}
			}
			morelzcon1 := lzcon[:len(lzcon)-(fjlenght-len(lzcon))] //取出多余的癞子
			for _, v := range morelzcon1 {
				canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
					Lz:     v,
					Become: 1,
				})
			}
			lesslzcon1 := lzcon[len(morelzcon1):] //取出刚刚好的飞机癞子
			if len(card1) != 0 {
				Sort(card)
				newlzcon1 := lesslzcon1[:len(lzcon)-len(card)] //使用普通牌 将癞子牌替换出来
				lastlzcon1 := lesslzcon1[len(newlzcon1):]      //多余出来得癞子 再次变成自己本身
				checkFeiji_sanShun_Min(card, card, newlzcon1, &canout)
				for _, v := range lastlzcon1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			} else {
				canout.Max = 4
				for i, v := range lesslzcon {
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: 4,
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 3,
					})
				}
			}
			*Canout = append(*Canout, canout)
		} else {
			//最大值
			var card []byte
			for _, v := range ptcon {
				if v == 14 && v == 13 {
					card = append(card, v)
				}
			}
			lesslzcon := append([]byte{}, lzcon...) //取出刚刚好的飞机癞子
			if len(card) != 0 {
				Sort(card)
				newlzcon := lesslzcon[:len(lzcon)-len(card)] //使用普通牌 将癞子牌替换出来
				lastlzcon := lesslzcon[len(newlzcon):]       //多余出来得癞子 再次变成自己本身
				checkFeiji_sanShun_Max(card, card, newlzcon, &canout)
				for _, v := range lastlzcon {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			} else {
				canout.Max = 14
				for i, v := range lesslzcon {
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: 14,
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 13,
					})
				}
			}
			*Canout = append(*Canout, canout)
			////////////
			canout.Max = 0
			canout.LzBecome = []GGameOutLzBecome{}
			////////
			//最小值
			var card1 []byte
			for _, v := range ptcon {
				if v == 3 && v == 4 {
					card1 = append(card1, v)
				}
			}
			lesslzcon1 := append([]byte{}, lzcon...) //取出刚刚好的飞机癞子
			if len(card1) != 0 {
				Sort(card1)
				newlzcon1 := lesslzcon[:len(lzcon)-len(card1)] //使用普通牌 将癞子牌替换出来
				lastlzcon1 := lesslzcon1[int(len(newlzcon1)):] //多余出来得癞子 再次变成自己本身
				checkFeiji_sanShun_Max(card, card, newlzcon1, &canout)
				for _, v := range lastlzcon1 {
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 1,
					})
				}
			} else {
				canout.Max = 4
				for i, v := range lesslzcon {
					if i < 3 {
						canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
							Lz:     v,
							Become: 3,
						})
						continue
					}
					canout.LzBecome = append(canout.LzBecome, GGameOutLzBecome{
						Lz:     v,
						Become: 4,
					})
				}
			}
			*Canout = append(*Canout, canout)
		}
	} else {
		//求最大值
		rememCanout := canout
		for i := 0; i < len(ptcon)-needpolish; i++ {
			var card []byte
			for j := 0; j < needpolish; j++ {
				card = append(card, ptcon[i+j])
			}
			checkFeiji_sanShun_Max(card, card, lzcon, &canout)
			if rememCanout.Max != canout.Max {
				break
			}
		}
		//////
		canout.Max = 0
		canout.LzBecome = []GGameOutLzBecome{}
		/////
		//求最小值
		rememCanout1 := canout
		needpolish := fjlenght - len(lzcon)
		for i := len(ptcon); i > 0+needpolish-2; i-- {
			var card []byte
			for j := 0; j < needpolish; j++ {
				card = append(card, ptcon[i-j])
			}
			Sort(card)
			checkFeiji_sanShun_Min(card, card, lzcon, &canout)
			if rememCanout1.Max != canout.Max {
				break
			}
		}
	}
}

//飞机带双
func Lz_Check_FeiJiDaiErShuang(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {
	// canout := CanOutType{
	// 	Cards: cards,
	// 	Ptcon: ptcon,
	// 	CT:    CT_AIRCRAFT_TWO,
	// }
	// single := make(map[byte]int)
	// duizi := make(map[byte]int)
	// sange := make(map[byte]int)
	// sige := make(map[byte]int)
	// numCout := make(map[byte]int)
	// for v := 17; v <= 3; v++ {
	// 	for _, v1 := range ptcon {
	// 		if byte(v) == GetLogicValue(v1) {
	// 			numCout[byte(v)] += 1
	// 		}
	// 	}
	// }
	// for key, v := range numCout {
	// 	if v == 1 {
	// 		single[key] += 1
	// 	}
	// 	if v == 2 {
	// 		duizi[key] += 1
	// 	}
	// 	if v == 3 {
	// 		sange[key] += 1
	// 	}
	// 	if v == 4 {
	// 		sige[key] += 1
	// 	}
	// }
	// //根据总牌数的长度 算出需要飞机的长度 ，在根据 lzcon 能否凑成飞机 判断该出牌类型是否合法
	// ptlength := len(cards) / 5
	// fjlenght := len(cards) - ptlength
	// needpolish := fjlenght - len(lzcon)
}

//检测牌型(有癞子的情况）
//////////////////////////////////

//检测牌型
func Lz_OutCard_Has(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType, SanDaiYiDui bool, SiDaiLiangDui bool) {
	cardlen := len(cards)
	switch cardlen {
	case 2:
		Lz_Check_Duizi(cards, ptcon, lzcon, Canout)
		break
	case 3:
		Lz_Check_Sange(cards, ptcon, lzcon, Canout)
		break
	case 4:
		Lz_Check_SanDaiYi(cards, ptcon, lzcon, Canout)
		break
	case 5:
		Lz_Check_Shunzi(cards, ptcon, lzcon, Canout)
		if SanDaiYiDui {
			Lz_Check_SanDaiEr(cards, ptcon, lzcon, Canout)
		}
		break
	case 6:
		Lz_Check_Shunzi(cards, ptcon, lzcon, Canout)
		Lz_Check_ShuangShun(cards, ptcon, lzcon, Canout)
		Lz_Check_FeiJi(cards, ptcon, lzcon, Canout)
		Lz_Check_SiDaiErDan(cards, ptcon, lzcon, Canout)
		break
	case 7:
		Lz_Check_Shunzi(cards, ptcon, lzcon, Canout)
		break
	case 8:
		Lz_Check_ShuangShun(cards, ptcon, lzcon, Canout)
		Lz_Check_FeiJiDaiErDan(cards, ptcon, lzcon, Canout)
		Lz_Check_Shunzi(cards, ptcon, lzcon, Canout)
		if SiDaiLiangDui {
			Lz_Check_SiDaiErShuang(cards, ptcon, lzcon, Canout)
		}
		break
	case 9:
		Lz_Check_Shunzi(cards, ptcon, lzcon, Canout)
		break
	case 10:
		Lz_Check_ShuangShun(cards, ptcon, lzcon, Canout)
		Lz_Check_FeiJiDaiErShuang(cards, ptcon, lzcon, Canout)
		Lz_Check_Shunzi(cards, ptcon, lzcon, Canout)
		break
	default:
		ChecTypeDefault(cards, ptcon, lzcon, Canout)
		break
	}
	for i := 0; i < len(*Canout)-1; i++ {
		if (*Canout)[i].CT == (*Canout)[i+1].CT {
			if (*Canout)[i].Max == (*Canout)[i+1].Max {
				*Canout = append((*Canout)[:i], (*Canout)[i+1:]...)
			}
		}
	}
	for i := 0; i < len(*Canout); i++ {
		if len((*Canout)[i].LzBecome) == 0 {
			*Canout = append((*Canout)[:i], (*Canout)[i+1:]...)
		}
	}
}
func ChecTypeDefault(cards []byte, ptcon []byte, lzcon []byte, Canout *[]CanOutType) {

}

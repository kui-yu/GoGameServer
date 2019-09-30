package main

// 图片
const (
	IMG_1 byte = iota + 1
	IMG_2
	IMG_3
	IMG_4
	IMG_5
	IMG_6
	IMG_7
	IMG_8
	IMG_9
	IMG_10
)

// 权重
const (
	WEIGHT_1  byte = 60
	WEIGHT_2  byte = 45
	WEIGHT_3  byte = 40
	WEIGHT_4  byte = 30
	WEIGHT_5  byte = 20
	WEIGHT_6  byte = 12
	WEIGHT_7  byte = 8
	WEIGHT_8  byte = 3
	WEIGHT_9  byte = 1
	WEIGHT_10 byte = 2
)

const BONUS = 10
const WILD = 9

// 中奖线
var Line_1 = []byte{5, 6, 7, 8, 9}
var Line_2 = []byte{0, 1, 2, 3, 4}
var Line_3 = []byte{10, 11, 12, 13, 14}
var Line_4 = []byte{0, 6, 12, 8, 4}
var Line_5 = []byte{10, 6, 2, 8, 14}
var Line_6 = []byte{0, 1, 7, 3, 4}
var Line_7 = []byte{10, 11, 7, 13, 14}
var Line_8 = []byte{5, 1, 2, 3, 9}
var Line_9 = []byte{5, 11, 12, 13, 9}

// 图标相同个数中奖
var IMG_1_DOUBLE = []int64{0, 0, 2, 5, 10}
var IMG_2_DOUBLE = []int64{0, 0, 3, 7, 20}
var IMG_3_DOUBLE = []int64{0, 0, 5, 12, 30}
var IMG_4_DOUBLE = []int64{0, 0, 7, 15, 50}
var IMG_5_DOUBLE = []int64{0, 0, 6, 15, 30}
var IMG_6_DOUBLE = []int64{0, 0, 9, 20, 60}
var IMG_7_DOUBLE = []int64{0, 0, 15, 35, 90}
var IMG_8_DOUBLE = []int64{0, 10, 20, 45, 300}

/////////////////////////////////////////////////////////
// 图标管理器，负责控制中奖
type MgrImg struct {
	Img     []byte
	Weight  []byte
	Bonus   []byte
	Wild    []byte
	MaxWild int32
	TWeight int32
	NWeight int32
	BWeight int32
}

func (this *MgrImg) Init() {
	this.Img = append([]byte{}, IMG_1, IMG_2, IMG_3, IMG_4, IMG_5, IMG_6, IMG_7, IMG_8, IMG_9, IMG_10)
	this.Weight = append([]byte{}, WEIGHT_1, WEIGHT_2, WEIGHT_3, WEIGHT_4, WEIGHT_5, WEIGHT_6, WEIGHT_7, WEIGHT_8, WEIGHT_9, WEIGHT_10)
	this.Bonus = []byte{1, 1, 1, 1, 1}
	this.Wild = []byte{1, 1, 1, 1, 1}
	this.MaxWild = 4

	for i, v := range this.Weight {
		this.TWeight = this.TWeight + int32(v)
		if i < len(this.Weight)-2 {
			this.NWeight = this.NWeight + int32(v)
		}
		if i < len(this.Weight)-1 {
			this.BWeight = this.BWeight + int32(v)
		}
	}
}

func (this *MgrImg) CheckLine(Scenes []byte, lineCount int64) ([18][]byte, int64, []byte) {
	AllLine := append([][]byte{}, Line_1, Line_2, Line_3, Line_4, Line_5, Line_6, Line_7, Line_8, Line_9)
	var Lines [18][]byte
	var Pow int64
	var isShow []byte
	// 从左到右
	for i, v := range AllLine {
		if i+1 > int(lineCount) {
			break
		}
		value := append([]byte{}, Scenes[v[0]], Scenes[v[1]], Scenes[v[2]], Scenes[v[3]], Scenes[v[4]])

		var wildCount = 0  // wild图标个数
		var index byte = 0 // 中奖图标
		var count = 0      // 中奖图标数

		for j := 0; j < len(value); j++ {
			if value[j] == WILD {
				// 百搭行首不计算百搭
				if j == 0 {
					break
				}
				count++
				wildCount++
				continue
			} else if value[j] == BONUS {
				break
			}

			// 未找到中奖图标Id时
			if index == 0 {
				index = value[j]
				count++
				continue
			} else {
				// 找到中奖图标Id后
				if index != value[j] {
					break
				}

				count++
			}
		}

		// 取到百搭符，但是没有取到其他图标，中奖图标为百搭
		if wildCount > 0 && index == 0 {
			index = WILD
		}

		var pow int64
		doubleIndex := count - 1
		switch index {
		case 1:
			pow = IMG_1_DOUBLE[doubleIndex]
		case 2:
			pow = IMG_2_DOUBLE[doubleIndex]
		case 3:
			pow = IMG_3_DOUBLE[doubleIndex]
		case 4:
			pow = IMG_4_DOUBLE[doubleIndex]
		case 5:
			pow = IMG_5_DOUBLE[doubleIndex]
		case 6:
			pow = IMG_6_DOUBLE[doubleIndex]
		case 7:
			pow = IMG_7_DOUBLE[doubleIndex]
		case 8:
			pow = IMG_8_DOUBLE[doubleIndex]
		}

		if pow == 0 {
			continue
		}

		// 添加需要显示的图标位置
		for i := 0; i < count; i++ {
			isNeedAdd := true
			for _, value := range isShow {
				if value == v[i] {
					isNeedAdd = false
					break
				}
			}
			if isNeedAdd {
				isShow = append(isShow, v[i])
			}
		}

		Pow = Pow + pow
		Lines[i] = append(Lines[i], v[0]+1, v[1]+1, v[2]+1, v[3]+1, v[4]+1)
		Lines[i] = append(Lines[i], byte(count))
	}

	// 从右到左
	for i, v := range AllLine {
		if i+1 > int(lineCount) {
			break
		}

		value := append([]byte{}, Scenes[v[4]], Scenes[v[3]], Scenes[v[2]], Scenes[v[1]], Scenes[v[0]])

		var wildCount = 0  // wild图标个数
		var index byte = 0 // 中奖图标
		var count = 0      // 中奖图标数

		for j := 0; j < len(value); j++ {
			if value[j] == WILD {
				// 百搭行首不计算百搭
				if j == 0 {
					break
				}
				count++
				wildCount++
				continue
			} else if value[j] == BONUS {
				break
			}

			// 未找到中奖图标Id时
			if index == 0 {
				index = value[j]
				count++
				continue
			} else {
				// 找到中奖图标Id后
				if index != value[j] {
					break
				}

				count++
			}
		}

		// 取到百搭符，但是没有取到其他图标，中奖图标为百搭
		if wildCount > 0 && index == 0 {
			index = WILD
		}

		var pow int64
		doubleIndex := count - 1
		switch index {
		case 1:
			pow = IMG_1_DOUBLE[doubleIndex]
		case 2:
			pow = IMG_2_DOUBLE[doubleIndex]
		case 3:
			pow = IMG_3_DOUBLE[doubleIndex]
		case 4:
			pow = IMG_4_DOUBLE[doubleIndex]
		case 5:
			pow = IMG_5_DOUBLE[doubleIndex]
		case 6:
			pow = IMG_6_DOUBLE[doubleIndex]
		case 7:
			pow = IMG_7_DOUBLE[doubleIndex]
		case 8:
			pow = IMG_8_DOUBLE[doubleIndex]
		}

		if pow == 0 {
			continue
		}

		// 添加需要显示的图标位置
		l := append([]byte{}, v[4]+1, v[3]+1, v[2]+1, v[1]+1, v[0]+1)
		for i := 0; i >= count; i++ {
			isNeedAdd := true
			for _, value := range isShow {
				if value == l[i] {
					isNeedAdd = false
					break
				}
			}
			if isNeedAdd {
				isShow = append(isShow, l[i])
			}
		}

		Pow = Pow + pow
		Lines[i+len(AllLine)] = append(Lines[i+len(AllLine)], l...)
		Lines[i+len(AllLine)] = append(Lines[i+len(AllLine)], byte(count))
	}

	return Lines, Pow, isShow
}

func (this *MgrImg) GetBonusCount(Scenes []byte) int64 {
	var count int64 = 0
	for _, v := range Scenes {
		if v == IMG_10 {
			count++
		}
	}

	return count
}

// 获取返回的图标  个数为imgCount
func (this *MgrImg) GetPlayResult(lineCount int64) ([]byte, [18][]byte, int64, []byte) {
	Bonus := append([]byte{}, this.Bonus...)
	Wild := append([]byte{}, this.Wild...)
	MaxWild := this.MaxWild

	var Scenes []byte
	for i := 0; i < 3; i++ {
		var Weight int32 = 0
		for j := 0; j < 5; j++ {
			if Bonus[j] > 0 {
				Weight = this.TWeight
			} else {
				Weight = this.BWeight
			}

			if Wild[j] == 0 || MaxWild == 0 {
				Weight = Weight - int32(WEIGHT_9)
			}

			rand := RandInt64(int64(Weight)) + 1
			for k, v := range this.Weight {
				rand = rand - int64(v)
				if rand <= 0 {
					if k == 8 {
						if Wild[j] > 0 && MaxWild > 0 {
							Wild[j] = Wild[j] - 1
							MaxWild = MaxWild - 1
						} else {
							k = k + 1
						}
					}

					if k == 9 {
						if Bonus[j] > 0 {
							Bonus[j] = Bonus[j] - 1
						}
					}
					Scenes = append(Scenes, byte(k)+1)
					break
				}
			}
		}
	}

	for i := len(Scenes); i < 15; i++ {
		k := RandInt64(int64(8)) + 1
		Scenes = append(Scenes, byte(k)+1)
	}

	Lines, Pow, isShow := this.CheckLine(Scenes, lineCount)
	return Scenes, Lines, Pow, isShow
}

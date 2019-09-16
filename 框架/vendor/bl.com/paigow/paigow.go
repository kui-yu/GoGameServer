package paigow

const (
	TYPE_ZHIZUN   = iota + 1 // 至尊1
	TYPE_TIAN                // 双天2
	TYPE_DI                  // 双地3
	TYPE_REN                 // 孖人4
	TYPE_HE                  // 孖和5
	TYPE_MEI                 // 孖梅6
	TYPE_CHANG               // 孖长7
	TYPE_BANDENG             // 孖板凳8
	TYPE_FU                  // 孖斧头9
	TYPE_HONG                // 孖红头10
	TYPE_GAO                 // 孖高脚11
	TYPE_ZERO                // 孖零霖12
	TYPE_NINE                // 杂九13
	TYPE_EIGHT               // 杂八14
	TYPE_SEVEN               // 杂七15
	TYPE_FIVE                // 杂五16
	TYPE_TIANKING            // 天王
	TYPE_DIKING              // 地王
	TYPE_TIANGANG            // 天杠
	TYPE_DIGANG              // 地杠
	TYPE_TIANNINE            // 天高九
	TYPE_DININE              // 地高九
	// 牌大小
	TYPE_NINE_1
	TYPE_NINE_2
	TYPE_EIGHT_1
	TYPE_EIGHT_2
	TYPE_SEVEN_1
	TYPE_SEVEN_2
	TYPE_FIVE_1
	TYPE_FIVE_2
	TYPE_SIX
	TYPE_THREE
)

const (
	// 文牌
	Card_WEN1_2  = 0x02
	Card_WEN1_4  = 0x04
	Card_WEN1_6  = 0x06
	Card_WEN1_7  = 0x07
	Card_WEN1_8  = 0x08
	Card_WEN1_10 = 0x0A
	Card_WEN1_11 = 0x0B
	Card_WEN1_12 = 0x0C

	Card_WEN2_4  = 0x14
	Card_WEN2_6  = 0x16
	Card_WEN2_10 = 0x1A
)

const (
	// 武牌
	Card_WU1_3 = 0x23
	Card_WU1_5 = 0x25
	Card_WU1_6 = 0x26
	Card_WU1_7 = 0x27
	Card_WU1_8 = 0x28
	Card_WU1_9 = 0x29
)

const (
	// 武牌
	Card_WU2_5 = 0x35
	Card_WU2_7 = 0x37
	Card_WU2_8 = 0x38
	Card_WU2_9 = 0x39
)

const (
	Color_Wen = 2 // 文牌 0、1
	Color_Wu  = 4 // 武牌 2、3
)

func GetInitNormalCards() []byte {
	var MVCard = []byte{}

	// 文牌
	cards := []byte{Card_WEN1_2, Card_WEN1_4, Card_WEN1_6, Card_WEN1_7,
		Card_WEN1_8, Card_WEN1_10, Card_WEN1_11,
		Card_WEN1_12, Card_WEN2_4, Card_WEN2_6, Card_WEN2_10}

	// 文牌一对
	for _, v := range cards {
		MVCard = append(MVCard, v)
		MVCard = append(MVCard, v)
	}

	// 武牌
	cards = []byte{Card_WU1_3, Card_WU1_5, Card_WU1_6,
		Card_WU1_7, Card_WU1_8, Card_WU1_9,
		Card_WU2_5, Card_WU2_7, Card_WU2_8, Card_WU2_9}

	// 武牌一张
	for _, v := range cards {
		MVCard = append(MVCard, v)
	}

	return MVCard[:]
}

// 获取牌型
func GetCardType(c int32) int32 {
	switch c {
	case Card_WEN1_12:
		return TYPE_TIAN
	case Card_WEN1_2:
		return TYPE_DI
	case Card_WEN1_8:
		return TYPE_REN
	case Card_WEN1_4:
		return TYPE_HE
	case Card_WEN1_10:
		return TYPE_MEI
	case Card_WEN1_6:
		return TYPE_CHANG
	case Card_WEN2_4:
		return TYPE_BANDENG
	case Card_WEN1_11:
		return TYPE_FU
	case Card_WEN2_10:
		return TYPE_HONG
	case Card_WEN1_7:
		return TYPE_GAO
	case Card_WEN2_6:
		return TYPE_ZERO
	// 武牌
	case Card_WU1_9:
		return TYPE_NINE_1
	case Card_WU2_9:
		return TYPE_NINE_2
	case Card_WU1_8:
		return TYPE_EIGHT_1
	case Card_WU2_8:
		return TYPE_EIGHT_2
	case Card_WU1_7:
		return TYPE_SEVEN_1
	case Card_WU2_7:
		return TYPE_SEVEN_2
	case Card_WU1_5:
		return TYPE_FIVE_1
	case Card_WU2_5:
		return TYPE_FIVE_2
	case Card_WU1_6:
		return TYPE_SIX
	case Card_WU1_3:
		return TYPE_THREE
	}

	return 0
}

// 获取牌型
func GetCardsType(c []int32) int32 {
	if c[0] == Card_WU1_6 && c[1] == Card_WU1_3 {
		return TYPE_ZHIZUN
	}
	if c[0] == c[1] {
		switch c[0] {
		case Card_WEN1_12:
			return TYPE_TIAN
		case Card_WEN1_2:
			return TYPE_DI
		case Card_WEN1_8:
			return TYPE_REN
		case Card_WEN1_4:
			return TYPE_HE
		case Card_WEN1_10:
			return TYPE_MEI
		case Card_WEN1_6:
			return TYPE_CHANG
		case Card_WEN2_4:
			return TYPE_BANDENG
		case Card_WEN1_11:
			return TYPE_FU
		case Card_WEN2_10:
			return TYPE_HONG
		case Card_WEN1_7:
			return TYPE_GAO
		case Card_WEN2_6:
			return TYPE_ZERO
		}
	}

	v1 := c[0] & 0x0F
	v2 := c[1] & 0x0F

	if (c[0] == Card_WU1_9 && c[1] == Card_WU2_9) || (c[1] == Card_WU1_9 && c[0] == Card_WU2_9) {
		return TYPE_NINE
	}
	if (c[0] == Card_WU1_8 && c[1] == Card_WU2_8) || (c[1] == Card_WU1_8 && c[0] == Card_WU2_8) {
		return TYPE_EIGHT
	}
	if (c[0] == Card_WU1_7 && c[1] == Card_WU2_7) || (c[1] == Card_WU1_7 && c[0] == Card_WU2_7) {
		return TYPE_SEVEN
	}
	if (c[0] == Card_WU1_5 && c[1] == Card_WU2_5) || (c[1] == Card_WU1_5 && c[0] == Card_WU2_5) {
		return TYPE_FIVE
	}

	if v1 == 0x0C {
		switch v2 {
		case 0x09:
			return TYPE_TIANKING
		case 0x08:
			return TYPE_TIANGANG
		case 0x07:
			return TYPE_TIANNINE
		}
	}

	if v1 == 0x02 {
		switch v2 {
		case 0x09:
			return TYPE_DIKING
		case 0x08:
			return TYPE_DIGANG
		case 0x07:
			return TYPE_DININE
		}
	}

	dian := (v1 + v2) % 10
	var dian2 int32 = 0
	if c[1] == Card_WU1_3 {
		dian2 = (v1 + 6) % 10
	} else if c[1] == Card_WU1_6 {
		dian2 = (v1 + 3) % 10
	}

	if dian2 > dian {
		dian = dian2
	}

	return 0 - dian
}

// 比牌
func CompareCard(b, p []int32) bool {
	bT := GetCardsType(b)
	pT := GetCardsType(p)

	if (bT > 0 && pT <= 0) || pT > 0 && bT <= 0 {
		if bT > pT {
			return true
		} else if bT == pT {
			for i := 0; i < 2; i++ {
				t1 := GetCardType(b[i])
				t2 := GetCardType(p[i])
				// 单牌类型越大，实际越小
				if t2 > t1 {
					return true
				} else if t2 < t1 {
					return false
				}
			}

			// 两张完全一致，庄赢
			return true
		}

		return false
	}

	// 都是大小牌，类型越小，大小越大
	// 都是点牌，bT、pT内容为负点数
	// 负点越大，大小越小
	// if bT > 0 || pT > 0 {
	if bT < pT {
		return true
	} else if bT == pT {
		for i := 0; i < 2; i++ {
			t1 := GetCardType(b[i])
			t2 := GetCardType(p[i])
			// 单牌类型越大，实际越小
			if t2 > t1 {
				return true
			} else if t2 < t1 {
				return false
			}
		}

		// 两张完全一致，庄赢
		return true
	}

	return false
}

// 文前面武后
func Sort(c []int32) []int32 {
	t1 := GetCardType(c[0])
	t2 := GetCardType(c[1])
	if t2 < t1 {
		c[0], c[1] = c[1], c[0]
	}

	return c
}

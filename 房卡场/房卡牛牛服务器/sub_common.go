package main

import (
	"math/rand"
	"time"
)

/*
* 为什么使用16进制做牌？
采用16进制，没16 向前进一位，所以我们很容易做出1到13的扑克牌，使用 二进制判断大小，15的二进制位 0000 1111，到时候我们使用15与牌相与，
如果>=16的数字我们不需要进行和计算，到头来还是拆分成 15以下，所以很容易就得出了扑克牌的值，判断花色使用 240，转换乘2进制位1111 0000，
我们也与牌相与并且>>4，当牌小于15的时候，值就是0， 当牌大于15 小于31 的时候，值就是1，当牌大于32，小于48的时候，值就是2，
当牌大于48，小于63的时候，值位3， 我们让这三个数字分别当成一个花色，这样就可以做到取值，也能取花色
*/
const (
	CARD_COLOR = 0xF0 //花色掩码  240
	CARD_VALUE = 0x0F //数值掩码	15
)

//方块
const (
	Card_Fang_1  = iota + 0x01 //方块 1   1
	Card_Fang_2                //2
	Card_Fang_3                //3
	Card_Fang_4                //4
	Card_Fang_5                //5
	Card_Fang_6                //6
	Card_Fang_7                //7
	Card_Fang_8                //8
	Card_Fang_9                //9
	Card_Fang_10               //10
	Card_Fang_J                //11
	Card_Fang_Q                //12
	Card_Fang_K                //13
)

//梅花
const (
	Card_Mei_1  = iota + 0x11 //梅花1  17
	Card_Mei_2                //18
	Card_Mei_3                //19
	Card_Mei_4                //20
	Card_Mei_5                //21
	Card_Mei_6                //22
	Card_Mei_7                //23
	Card_Mei_8                //24
	Card_Mei_9                //25
	Card_Mei_10               //26
	Card_Mei_J                //27
	Card_Mei_Q                //28
	Card_Mei_K                //29
)

//红桃
const (
	Card_Hong_1  = iota + 0x21 //红桃1  33
	Card_Hong_2                //34
	Card_Hong_3                //35
	Card_Hong_4                //36
	Card_Hong_5                //37
	Card_Hong_6                //38
	Card_Hong_7                //39
	Card_Hong_8                //40
	Card_Hong_9                //41
	Card_Hong_10               //42
	Card_Hong_J                //43
	Card_Hong_Q                //44
	Card_Hong_K                //45
)

//黑桃
const (
	Card_Hei_1  = iota + 0x31 //黑桃1  49
	Card_Hei_2                //50
	Card_Hei_3                //51
	Card_Hei_4                //52
	Card_Hei_5                //53
	Card_Hei_6                //54
	Card_Hei_7                //55
	Card_Hei_8                //56
	Card_Hei_9                //57
	Card_Hei_10               //58
	Card_Hei_J                //59
	Card_Hei_Q                //60
	Card_Hei_K                //61
)

//花色
const (
	CARD_COLOR_Fang = iota //方块 0
	CARD_COLOR_Mei         // 梅花 1
	CARD_COLOR_Hong        //红心 2
	CARD_COLOR_Hei         //黑桃 3
)

//牛类型
const (
	NIU_FIVE_SMALL = 14 //五小牛
	NIU_BOOM       = 13 //四炸
	NIU_FIVE_COLOR = 12 //五花牛
	NIU_FORE_COLOR = 11 //四花牛
	NIU_NIU        = 10 //牛牛
	NIU_NINE       = 9  //牛9
	NIU_EIGHT      = 8  //牛8
	NIU_SEVEN      = 7  //牛7
	NIU_SIX        = 6  //牛6
	NIU_FIVE       = 5  //牛5
	NIU_FORE       = 4  //牛4
	NIU_THREE      = 3  //牛3
	NIU_TWO        = 2  //牛2
	NIU_ONE        = 1  //牛1
	NIU_ZERO       = 0  //无牛
)

/////////////////////////////////////////////////////////
//卡牌管理器，负责做牌
type MgrCard struct {
	MVCard       []int //初始牌
	MVSourceCard []int //洗过的牌（也就是之后用来发的牌）
}

//初始化
func (this *MgrCard) InitCards() {
	this.MVCard = []int{}
	this.MVSourceCard = []int{}
}

//赋值
func (this *MgrCard) InitNormalCards() {
	begaincard := []int{Card_Fang_1, Card_Mei_1, Card_Hong_1, Card_Hei_1} //先创建一个切牌，初始化方块1，梅花1，红心1，黑桃1。
	for _, v := range begaincard {                                        //遍历原始集合，利用每一种牌型的最小值逐渐往上添加同一牌型的扑克牌，直到添加到K
		for j := 0; j < 13; j++ {
			this.MVCard = append(this.MVCard, v+j)
		}
	}
}

//洗牌
func (this *MgrCard) Shuffle() {
	this.MVSourceCard = append([]int{}, this.MVCard...)  //将初始的牌组全部放进This.MVSourceCard 切片中
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //利用time.Now()获取当前时间通过unixNano()函数获取unix时间，通过这个时间数字使用rand.NewSource获取一个随机种子，再通过rand.New()方法获取一个随机数
	perm := r.Perm(len(this.MVCard))                     //根据传入的 int类型的参数n 返回一串 0到n随机打乱的int切片，例如 传入5 返回可能是 [0,2,1,3,4]，我们可以利用这个方法只要将一个切片的长度传入进去，然后就可以获得一个乱下标切片，然后我们可以根据切片打乱打乱原来的切片，达到洗牌的目的
	for i, randIndex := range perm {                     //i:当前循环的次数，可以理解为 perm的下标      randIndex:perm的每一个下标值
		this.MVSourceCard[i] = this.MVCard[randIndex] //MVSourceCard[i]下放入 this.MVCard[randIndex]的值，也就是将一个切片的顺序，根据perm返回的下标切片打乱。
	}
}

//发手牌
//参数 num为需要发给手牌的数量
func (this *MgrCard) SendHandCard(num int) []int {
	list := append([]int{}, this.MVSourceCard[0:num]...) //从 所有可发牌中赋值 五张给list
	//删除第0个元素开始
	i := 0
	this.MVSourceCard = append(this.MVSourceCard[:i], this.MVSourceCard[i+num:]...) //因为 i=0,所以this.MVSourceCard[:i]表示空的切片加上this.MVSourceCard[num:] 也就是减去了发出去的手牌
	//返回
	return list //返回手牌切片
}

// 获取卡牌花色
func GetCardColor(card int) int {
	return (card & CARD_COLOR) >> 4 //CARD_COLOR 16进制 转换成10进制为 240  240转换为2进制为 11110000，如果传入一个梅花5，为21 ，转换为2进制为10101   21&240=10000 >>4=1  梅花:1 所以这样就可以取出牌的花色了。
}

// 获取卡牌值
func GetCardValue(card int) int {
	return (card & CARD_VALUE) //CARD_VALUE 16进制 转换为10进制为 15  相与得出的数一定是比较小的那个数还要小的。所以按照上面的方法很容易得出大小
}

//转化值
func GetLogicValue(card int) int { //获取每张牌的实际值（用来判断是否有牛）
	d := GetCardValue(card)
	if d > 9 {
		return 10
	}
	return d
}

//获取卡牌值
func GetSmallValue(cards []int) []int { //转化每一张卡牌的值，将其值存入数组中返回
	var smallCards []int
	for i := 0; i < len(cards); i++ {
		vi := GetCardValue(cards[i])
		smallCards = append(smallCards, vi)
	}
	return smallCards
}

//获取卡牌牛牛值
func GetSmallNiuValue(cards []int) []int { //转换每一张卡牌的值 （例如将 J(11),Q(12),K(13)转换成 10),用来进行 有无牛判断
	var smallCards []int
	for i := 0; i < len(cards); i++ {
		vi := GetLogicValue(cards[i])
		smallCards = append(smallCards, vi)
	}
	return smallCards //返回一个 经过转换最终 用来判断是否有牛的牌切片
}

//排序
func Sort(cards []int) []int {
	cs := append([]int{}, cards...)  //将牌组切片参数 添加到cs 变量中
	for i := 0; i < len(cs)-1; i++ { //将大的数往后排，小的数往前牌
		for j := i + 1; j < len(cs); j++ {
			vi := GetCardValue(cs[i])
			vj := GetCardValue(cs[j])
			if vi < vj || ((vi == vj) && (GetCardColor(cs[i]) > GetCardColor(cs[j]))) { //如果值
				vt := cs[i]
				cs[i] = cs[j]
				cs[j] = vt
			}
		}
	}
	return cs
}

//算五小牛
func MathFiveSmall(cards []int) int {
	var smallCards = GetSmallValue(cards)

	var sum int

	for i := 0; i < len(smallCards); i++ {
		if smallCards[i] > 4 {
			return NIU_ZERO
		}
		sum += smallCards[i]
	}
	if sum < 10 {
		return NIU_FIVE_SMALL
	}
	return NIU_ZERO
}

//算四炸
func MathBoom(cards []int) (int, []int) { //返回牛几，如果是四炸返回炸弹牌组切片

	niuCards := append([]int{}, cards...)

	smallCards := GetSmallValue(cards)

	var count1 int
	var count2 int

	for i := 0; i < len(smallCards); i++ {
		if smallCards[0] == smallCards[i] {
			count1++
		}
		if smallCards[1] == smallCards[i] {
			count2++
		}
	}
	if count1 == 4 {
		for i := 0; i < len(smallCards); i++ {
			if smallCards[0] != smallCards[i] {
				tempCard := niuCards[i]
				//移除
				niuCards = append(niuCards[:i], niuCards[i+1:]...)
				//后面追加
				niuCards = append(niuCards, tempCard)
				break
			}
		}
		return NIU_BOOM, niuCards
	}
	if count2 == 4 {
		for i := 0; i < len(smallCards); i++ {
			if smallCards[1] != smallCards[i] {
				tempCard := niuCards[i]
				//移除
				niuCards = append(niuCards[:i], niuCards[i+1:]...)
				//后面追加
				niuCards = append(niuCards, tempCard)
				break
			}
		}
		return NIU_BOOM, niuCards
	}

	return NIU_ZERO, niuCards
}

//算五花牛/四花牛
func MathColor(cards []int) (int, []int) {

	var smallCards = GetSmallValue(cards)

	var count int
	for i := 0; i < len(smallCards); i++ {
		if smallCards[i] > 10 {
			count++
		}
	}
	//五花牛
	if count == 5 {
		return NIU_FIVE_COLOR, cards
	}
	//四花牛
	// if count == 4 {
	// 	return NIU_FORE_COLOR, cards
	// }
	return NIU_ZERO, cards
}

//算牛牛
func MathTen(cards []int) (int, []int) {
	var smallCards = GetSmallNiuValue(cards)

	//牛点
	var niuPoint int

	for i := 0; i < len(smallCards); i++ {
		for j := i + 1; j < len(smallCards); j++ {
			for k := j + 1; k < len(smallCards); k++ {
				//有牛
				if (smallCards[i]+smallCards[j]+smallCards[k])%10 == 0 {
					var fourKey int
					for l := 0; l < len(smallCards); l++ {
						if l != i && l != j && l != k {
							fourKey = l
							break
						}
					}
					var fiveKey int
					for l := 0; l < len(smallCards); l++ {
						if l != i && l != j && l != k && l != fourKey {
							fiveKey = l
							break
						}
					}
					//123
					niuPoint = (smallCards[fourKey] + smallCards[fiveKey]) % 10
					if niuPoint == 0 {
						niuPoint = NIU_NIU
					}
					temp := []int{cards[i], cards[j], cards[k], cards[fourKey], cards[fiveKey]}
					return niuPoint, temp
				}
			}
		}
	}
	return niuPoint, cards
}

//获取结果
func GetResult(cards []int) (int, []int) {

	var sortCards []int = Sort(cards)
	// logs.Debug("收到手牌", sortCards)

	//五小牛判断
	// fiveSmall := MathFiveSmall(sortCards)
	// if fiveSmall > 0 {
	// 	return fiveSmall, sortCards
	// }

	//炸弹判断
	boom, boomCards := MathBoom(sortCards)
	if boom > 0 {
		return boom, boomCards
	}

	//花牛
	color, colorCards := MathColor(sortCards)
	if color > 0 {
		return color, colorCards
	}

	//判断牛牛
	return MathTen(sortCards)
}

//比牌
func SoloResult(player1 *ExtPlayer, player2 *ExtPlayer) int { //标记牛类型大小
	if player1.NiuPoint > player2.NiuPoint {
		return 1
	} else if player1.NiuPoint < player2.NiuPoint {
		return 2
	} else {
		//比大小
		if player1.NiuPoint == NIU_BOOM {
			//炸弹比较
			if player1.NiuCards[0] > player2.NiuCards[0] {
				return 1
			} else {
				return 2
			}
		} else {
			var sortCards1 []int = Sort(player1.HandCard)
			var sortCards2 []int = Sort(player2.HandCard)
			//比最大
			var card1 = GetCardValue(sortCards1[0])
			var card2 = GetCardValue(sortCards2[0])
			if card1 > card2 {
				return 1
			} else if card1 < card2 {
				return 2
			} else {
				//比花色
				if GetCardColor(sortCards1[0]) > GetCardColor(sortCards2[0]) {
					return 1
				} else {
					return 2
				}
			}
		}
	}
}

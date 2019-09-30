package main

import (
	"fmt"
	"time"

	// "fmt"
	"crypto/rand"
	"logs"
	"math/big"
)

type ExtDesk struct {
	Desk
	Qzhoushicount  int                   //区域走势限制局数
	Zzhoushicount  int                   //庄家走势限制局数
	DownBet        []int64               //区域总下注金币集合  0,青龙  1，白虎, 2，朱雀, 3，玄武
	DownBetZhenshi []int64               //真实玩家区域总下注集合
	ChipList       []int64               //筹码列表
	MaxBet         int64                 //区域限红
	ChairList      []PlayerInfoByChair   //座位玩家信息
	DownCards      []uint8               // 所有的牌
	StatusAndTimes StatuAndTimes         //状态及状态时间
	CardGroupArray map[int]CardGroupInfo //玩家和庄家的牌，庄家牌索引是最后一个
	NeedBro        bool                  //是否需要广播下注（即是否 有新的下注，需要广播）
	NeedUpdata     bool                  //是否需要更新其他玩家下注
	GameZs         ZouShiToClient
	WinArea        []int //胜利区域
	AreaWinDouble  []int //存区域与庄家比牌，胜者的倍数
	BalanceResult  []int64
}

//桌子初始化
func (this *ExtDesk) InitExtData() {
	logs.Debug("初始化ExtDesk......")
	//初始化下注区域集合
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.DownBet = append(this.DownBet, 0)
		this.DownBetZhenshi = append(this.DownBetZhenshi, 0)
	}
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.BalanceResult = append(this.BalanceResult, 0)
	}

	//初始化玩家座位
	for i := 0; i < gameConfig.LimitInfo.SeatCount; i++ {
		this.ChairList = append(this.ChairList, PlayerInfoByChair{
			Uid:   0,
			Nick:  "",
			Coins: 0,
		})
	}
	//初始化状态和状态时间
	this.StatusAndTimes = StatuAndTimes{
		ShuffleId:       gameConfig.GameStatusTimer.ShuffleId,
		ShuffleMs:       gameConfig.GameStatusTimer.ShuffleMs,
		FaCardId:        gameConfig.GameStatusTimer.FaCardId,
		FaCardMS:        gameConfig.GameStatusTimer.FaCardMS,
		StartdownbetsId: gameConfig.GameStatusTimer.StartdownbetsId,
		StartdownbetsMs: gameConfig.GameStatusTimer.StartdownbetsMs,
		StopdownbetsId:  gameConfig.GameStatusTimer.StopdownbetsId,
		StopdownbetsMs:  gameConfig.GameStatusTimer.StopdownbetsMs,
		DownBetsId:      gameConfig.GameStatusTimer.DownBetsId,
		DownBetsMS:      gameConfig.GameStatusTimer.DownBetsMS,
		OpenCardId:      gameConfig.GameStatusTimer.OpenCardId,
		OpenCardMS:      gameConfig.GameStatusTimer.OpenCardMS,
		BalanceId:       gameConfig.GameStatusTimer.BalanceId,
		BalanceMS:       gameConfig.GameStatusTimer.BalanceMS,
	}
	//初始化手牌
	this.CardGroupArray = make(map[int]CardGroupInfo)
	//绑定监听
	this.Handle[MSG_GAME_AUTO] = this.HandleAuto
	this.Handle[MSG_GAME_INFO_QDESKINFO] = this.HandleQDeskInfo
	this.Handle[MSG_GAME_INFO_DOWNBET] = this.HandleDownBet
	this.Handle[MSG_GAME_INFO_GETMOREPLAYER] = this.HandleGetMorePlayer
	this.Handle[MSG_GAME_DISCONNECT] = this.HandleDisconnect
	this.Handle[MSG_GAME_INFO_GET_RECORD] = this.HandleGetRecord
	//初始化needbro
	this.NeedBro = false
	//初始化走势
	zstc := ZouShiToClient{}
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount+1; i++ {
		zstc.Zoushi = append(zstc.Zoushi, []Pzshi{})
	}
	this.GameZs = zstc
	//初始化胜利区域
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.WinArea = append(this.WinArea, 0)
	}
	//进入洗牌状态
	this.AddTimer(0, 0, this.status_shuffle, nil)
}

//重置桌子方法
func (this *ExtDesk) Rest() {
	//初始化下注区域集合
	this.DownBet = []int64{}
	this.DownBetZhenshi = []int64{}
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.DownBet = append(this.DownBet, 0)
		this.DownBetZhenshi = append(this.DownBetZhenshi, 0)
	}
	//初始化手牌
	this.CardGroupArray = make(map[int]CardGroupInfo)
	this.NeedBro = false
	this.BalanceResult = []int64{}
	for i := 0; i < gameConfig.LimitInfo.BetAreaCount; i++ {
		this.BalanceResult = append(this.BalanceResult, 0)
	}
}

// 随机数生成器
func GetRandomNum(min, max int) (int, error) {
	maxBigInt := big.NewInt(int64(max - min))
	i, err := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < 0 {
		return 0, err
	}
	return int(i.Int64()) + min, err
}

//根据玩家金币判断还有哪些筹码能够下注
func (this *ExtDesk) CanUseChip(p *ExtPlayer) int {
	indexs := -1
	var allbet int64
	for _, v := range p.DownBet {
		allbet += v
	}

	// logs.Debug("````````````````````````````:", G_DbGetGameServerData.GameConfig.TenChips)
	for i, v := range G_DbGetGameServerData.GameConfig.TenChips {
		if v+allbet <= (p.Coins-v)/gameConfig.LimitInfo.Downbet_Double_Comp {
			indexs = i
		}
	}
	if !p.Robot {
		// logs.Debug("玩家总下注：``````````````````````", allbet)
		// logs.Debug("玩家金币:````````````````", p.Coins)
		// logs.Debug("indexs:`````````````````````", indexs)
	}
	return indexs
}

//广播状态
func (this *ExtDesk) BroStatusTime(times int) {
	status := GameStatuInfo{
		Id:         MSG_GAME_INFO_STATUSCHANGE, //协议号
		Status:     this.GameState,
		StatusTime: times,
	}
	this.BroadcastAll(MSG_GAME_INFO_STATUSCHANGE, status)
}

//获取当前进入方法名称 例如: 荣耀厅 ，王牌厅，战神厅
func (this *ExtDesk) GetRoomName() string {
	gradeId := GCONFIG.GradeType
	var roomName string
	if gradeId == 1 {
		roomName = "荣耀厅"
	} else if gradeId == 2 {
		roomName = "王牌厅"
	} else {
		roomName = "战神厅"
	}
	return roomName
}

//将名字改为****
func (this *ExtDesk) ChangeNick(nick string) string {
	if len(nick) > 3 {
		return "***" + nick[3:]
	} else {
		return "***" + nick
	}
}

//检测还没有在展示列表(座位)上的玩家,返回玩家id切片
func (this *ExtDesk) FindNoChairPlayer() []int64 {
	uidlist := []int64{}
	for _, v := range this.Players {
		if !v.IsOnChair {
			uidlist = append(uidlist, v.Uid)
		}
	}
	if len(uidlist) > 0 {
		return uidlist
	} else {
		return nil
	}
}

//玩家入座
func (this *ExtDesk) OnChair(p *ExtPlayer) {
	//判断是否存在空位
	var index = -1
	for i, v := range this.ChairList {
		if v.Uid == p.Uid {
			logs.Debug("头像已经存在")
			return
		}
		if v.Uid == 0 {
			index = i
			break
		}
	}
	//玩家入座
	if index != -1 {
		this.ChairList[index].Uid = p.Uid
		this.ChairList[index].Nick = p.Nick
		this.ChairList[index].Avatar = p.Head
		this.ChairList[index].Coins = p.Coins
		p.IsOnChair = true
	}
}

//玩家离座
func (this *ExtDesk) UpChair(p *ExtPlayer) {
	//判断玩家是否在座位上，如果有，则离开,
	for i, v := range this.ChairList {
		if v.Uid == p.Uid {
			this.ChairList[i].Uid = 0
			return
		}
	}
	//查找没在座位上的玩家，将其入座
	plist := this.FindNoChairPlayer()
	if plist != nil {
		var pl *ExtPlayer
		for _, v := range this.Players {
			if v.Uid == plist[0] {
				pl = v
				break
			}
		}
		this.OnChair(pl)
	}
}

//座位变更通知
func (this *ExtDesk) BroChairChange() {
	for _, v := range this.Players {
		v.SendNativeMsg(MSG_GAME_INFO_CHAIRCHANGE, &struct {
			Id        int
			ChairList []PlayerInfoByChair
		}{
			Id:        MSG_GAME_INFO_CHAIRCHANGE,
			ChairList: v.getChairList(),
		})
	}
}

//座位变更通知
func (this *ExtDesk) BroChairChangeNoto(p *ExtPlayer) {
	for _, v := range this.Players {
		if v.Uid != p.Uid {
			v.SendNativeMsg(MSG_GAME_INFO_CHAIRCHANGE, &struct {
				Id        int
				ChairList []PlayerInfoByChair
			}{
				Id:        MSG_GAME_INFO_CHAIRCHANGE,
				ChairList: v.getChairList(),
			})
		}
	}
}

//广播玩家下注信息
func (this *ExtDesk) BroDownBetInfo(d interface{}) {
	if this.NeedBro {
		for _, v := range this.Players {
			var downbet []int64
			for i, v1 := range v.OtherDownBet {
				res := v1 - v.OldtherDownBet[i]
				downbet = append(downbet, res)
			}
			v.SendNativeMsg(MSG_GAME_INFO_DOWNBET_BRO, DownBetBro{
				Id:           MSG_GAME_INFO_DOWNBET_BRO,
				DownBet:      this.DownBet,
				OtherDownBet: downbet,
				MyDownBet:    v.DownBet,
			})
		}
		this.NeedBro = false
		this.NeedUpdata = true
	}
	if this.GameState == gameConfig.GameStatusTimer.DownBetsId || this.GameState == this.StatusAndTimes.StartdownbetsId {
		this.AddTimer(999, 1, this.BroDownBetInfo, nil)
	}
}

//比大小
func GetResult(cInfo map[int]CardGroupInfo) ([]int, []int) { // 0 输 1赢 倍率
	seatCount := gameConfig.LimitInfo.BetAreaCount
	if len(cInfo) != seatCount+1 {
		logs.Debug("牌组长度不对")
	}

	var seatResult []int
	var multipe []int

	zCard := cInfo[seatCount]
	zCardLv := zCard.CardGroupType
	if zCardLv >= CardGroupType_None {
		zCardLv = 0
	}
	// fmt.Println(zCardLv, "庄等级")
	zMultiple := GetCardMultiple(zCardLv)

	for i := 0; i < seatCount; i++ {
		pCard := cInfo[i]
		pCardLv := pCard.CardGroupType
		// fmt.Println(pCardLv, "闲等级")
		if pCardLv >= CardGroupType_None {
			pCardLv = 0
		}
		pMultiple := GetCardMultiple(pCardLv)

		if zCardLv > pCardLv {
			seatResult = append(seatResult, 0)
			multipe = append(multipe, zMultiple)
		} else if pCardLv > zCardLv {
			seatResult = append(seatResult, 1)
			multipe = append(multipe, pMultiple)
		} else {
			if zCard.MaxCard&0x0F > pCard.MaxCard&0x0F {
				seatResult = append(seatResult, 0)
				multipe = append(multipe, zMultiple)
			} else if zCard.MaxCard&0x0F < pCard.MaxCard&0x0F {
				seatResult = append(seatResult, 1)
				multipe = append(multipe, pMultiple)
			} else {
				if GetCardColor(zCard.MaxCard) > GetCardColor(pCard.MaxCard) {
					logs.Debug("庄家赢")
					seatResult = append(seatResult, 0)
					multipe = append(multipe, zMultiple)
				} else if GetCardColor(zCard.MaxCard) < GetCardColor(pCard.MaxCard) {
					logs.Debug("庄家输")
					seatResult = append(seatResult, 1)
					multipe = append(multipe, pMultiple)
				} else {
					logs.Debug("其他")
					seatResult = append(seatResult, 0)
					multipe = append(multipe, zMultiple)
				}
			}
		}
	}
	return seatResult, multipe
}

//获取倍数
func GetCardMultiple(cardLv CardGroupType) int {
	if cardLv <= 6 {
		return 1
	} else if cardLv <= 7 {
		return 2
	} else if cardLv <= 9 {
		return 3
	} else if cardLv <= 10 {
		return 4
	} else {
		return 5
	}
}

//风控
func (this *ExtDesk) allotCard() {
	logs.Debug("开牌")
	//四个区域
	for area := 0; area < gameConfig.LimitInfo.BetAreaCount; area++ {
		cardGroupInfo := this.CardGroupArray[area]
		checkCards := this.SendCard(5)
		cardGroupInfo.CardGroupType, cardGroupInfo.MaxCard, cardGroupInfo.Cards = CalcCards(checkCards)
		this.CardGroupArray[area] = cardGroupInfo
	}
	//庄家区域
	var bankGroupInfo CardGroupInfo
	bankerPoker := this.SendCard(5)
	bankGroupInfo.CardGroupType, bankGroupInfo.MaxCard, bankGroupInfo.Cards = CalcCards(bankerPoker)
	this.CardGroupArray[4] = bankGroupInfo
	//以上的牌为纯随机，现在判断是否进入风控
	logs.Debug("现在没有进行风控转换", this.CardGroupArray)
	//获取真实玩家下注
	allbet := this.getUserBet()
	if allbet > 0 && CD-CalPkAll(StartControlTime, time.Now().Unix()) < 0 && GetCostType() != 2 {
		//进入风控
		this.ContollerWinOrLost(true)
		logs.Debug("进入风控", CD, CalPkAll(StartControlTime, time.Now().Unix()))
	} else if allbet <= 0 {
		//控制75%的胜率
		ra, _ := GetRandomNum(0, 100)
		fmt.Println("ra:", ra)
		if ra < 75 {
			this.ContollerWinOrLost(false)
		}
		logs.Debug("控制为胜率75")
	}
	logs.Debug("进行风控转换完成之后的牌:", this.CardGroupArray)
}

func (this *ExtDesk) ContollerWinOrLost(win bool) {
	if win {
		result := this.getWin()
		if len(result) != 0 {
			this.CardGroupArray = result
		}
	} else {
		result := this.getLose()
		if len(result) != 0 {
			this.CardGroupArray = result
		}
	}
}

func (this *ExtDesk) getWin() (winResult map[int]CardGroupInfo) {
	cardarry := []CardGroupInfo{}
	fmt.Println("this.cards:", this.CardGroupArray)
	for i := 0; i < len(this.CardGroupArray); i++ {
		cardarry = append(cardarry, this.CardGroupArray[i])
	}
	fmt.Println("cardarry:", cardarry)
	cardarrys := [][]CardGroupInfo{}
	for i := 0; i < len(cardarry); i++ {
		cardarry = append(cardarry[len(cardarry)-1:], cardarry[:len(cardarry)-1]...)
		cardarrys = append(cardarrys, cardarry)
	}
	fmt.Println("cardsarryssss:", cardarrys)
	result := []map[int]CardGroupInfo{}
	for i := 0; i < len(cardarrys); i++ {
		result_ := map[int]CardGroupInfo{}
		fmt.Println("len(cardssssss):", len(cardarrys[i]))
		for j := 0; j < len(cardarrys[i]); j++ {
			result_[j] = cardarrys[i][j]
		}
		result = append(result, result_)
	}
	//判断开奖结果庄家是不是赚的
	for _, v := range result {
		areaResult, multipleResult := GetResult(v)
		var bankerResult int64 = 0 // 庄家输赢结果
		for i, v := range areaResult {
			if v == 0 {
				//庄家赢
				bankerResult += this.DownBetZhenshi[i] * int64(multipleResult[i])
			} else {
				//庄家输
				bankerResult -= this.DownBetZhenshi[i] * int64(multipleResult[i])
			}
		}
		if bankerResult > 0 {
			winResult = v
			return
		}
	}
	return
}
func (this *ExtDesk) getLose() (loseResult map[int]CardGroupInfo) {
	cardarry := []CardGroupInfo{}
	fmt.Println("this.cards:", this.CardGroupArray)
	for i := 0; i < len(this.CardGroupArray); i++ {
		cardarry = append(cardarry, this.CardGroupArray[i])
	}
	fmt.Println("cardarry:", cardarry)
	cardarrys := [][]CardGroupInfo{}
	for i := 0; i < len(cardarry); i++ {
		cardarry = append(cardarry[len(cardarry)-1:], cardarry[:len(cardarry)-1]...)
		fmt.Println("lencardarrrryyyyy:", len(cardarry))
		cardarrys = append(cardarrys, cardarry)
	}
	fmt.Println("cardsarryssss:", cardarrys)
	result := []map[int]CardGroupInfo{}
	for i := 0; i < len(cardarrys); i++ {
		result_ := map[int]CardGroupInfo{}
		for j := 0; j < len(cardarrys[i]); j++ {
			result_[j] = cardarrys[i][j]
		}
		result = append(result, result_)
	}
	//判断开奖结果庄家是不是赚的
	for _, v := range result {
		fmt.Println("v::", v)
		areaResult, multipleResult := GetResult(v)
		var bankerResult int64 = 0 // 庄家输赢结果
		for i, v := range areaResult {
			if v == 0 {
				//庄家赢
				bankerResult += this.DownBetZhenshi[i] * int64(multipleResult[i])
			} else {
				//庄家输
				bankerResult -= this.DownBetZhenshi[i] * int64(multipleResult[i])
			}
		}
		if bankerResult < 0 {
			loseResult = v
			return
		}
	}
	return
}

//判断玩家真实下注
func (this *ExtDesk) getUserBet() (allbet int64) {
	for _, v := range this.Players {
		if !v.Robot {
			for _, v1 := range v.DownBet {
				allbet += v1
			}
		}
	}
	return
}

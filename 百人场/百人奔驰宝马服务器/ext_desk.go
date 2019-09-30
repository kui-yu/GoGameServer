package main

import (
	"fmt"
	"logs"
)

type ExtDesk struct {
	Desk

	Mark       string                // 房间标志
	Bets       map[int]GBetCountInfo // 下注区域信息
	GameResult int                   // 开奖结果
	Rate       float64               // 抽水率
	LogoRecord []int                 // 游戏开奖记录

	fsms       map[int]FSMBase // fsm状态机集合
	currentFSM FSMBase         // 当前状态机
	upFSM      FSMBase         // 上一个状态机
}

func (this *ExtDesk) ResetExtDesk() {
	this.GameResult = -1
	for _, v := range this.Bets {
		v.DownBetValue = 0
		v.UserBetValue = 0
	}
}

func (this *ExtDesk) InitExtData() {
	logs.Debug("初始化ext_dest")

	this.Mark = FormatDeskId(this.Id, GCONFIG.GradeType)

	InitMultiple()

	// 初始化下注区域
	this.Bets = make(map[int]GBetCountInfo)

	betCount := gameConfig.LimitInfo.BetCount
	// DebugLog("座位号：", betCount)
	for i := 0; i < betCount; i++ { //8个下注区
		this.Bets[i] = GBetCountInfo{
			Id:           i,
			DownBetValue: 0,
			UserBetValue: 0,
		}
	}

	this.fsms = make(map[int]FSMBase)

	//状态阶段
	this.addFSM(GAME_STATUS_DOWNBET, new(FSMDownBet))
	this.addFSM(GAME_STATUS_LOTTERY, new(FSMLottery))
	this.addFSM(GAME_STATUS_BALANCE, new(FSMSettle))

	this.addListener()

	DebugLog("开始运行")

	// 启动状态
	this.RunFSM(GAME_STATUS_DOWNBET)
}

func (this *ExtDesk) GetFSM(mark int) FSMBase {
	if mark != 0 {
		return this.fsms[mark]
	}
	return this.currentFSM
}

func (this *ExtDesk) RunFSM(mark int) {
	if this.currentFSM != nil {
		this.upFSM = this.currentFSM
		this.upFSM.Leave()
	}

	this.currentFSM = this.GetFSM(mark)
	this.currentFSM.Run()
}

// deskmgr -> 直接调用
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d2 *DkInMsg) {
	DebugLog("============接收到重新连接")
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:     MSG_GAME_RECONNECT_REPLY,
		Result: 0,
	})

	p.SendNativeMsg(MSG_GAME_INFO_RRECONNECT, &GReConnectReply{
		Id:     MSG_GAME_INFO_RRECONNECT,
		Result: 0,
	})
}

// 玩家离开 deskmgr -> 直接调用
func (this *ExtDesk) Leave(p *ExtPlayer) { // 400007
	DebugLog(">>>>>>>>>>>>>>>>>Leave 玩家离开")
	this.handleGameBack(p, nil)
}

func (this *ExtDesk) addListener() {
	this.Handle[MSG_GAME_AUTO] = this.handleGameAuto           // 玩家进入返回匹配结果
	this.Handle[MSG_GAME_INFO_QDESKINFO] = this.handleDeskInfo // 请求游戏桌子信息
	this.Handle[MSG_GAME_DISCONNECT] = this.handleDisConnect   // 用户掉线，处理与退出房间一致
}

func (this *ExtDesk) addFSM(mark int, fsm FSMBase) {
	fsm.InitFSM(mark, this)
	this.fsms[mark] = fsm
}

//========================大厅消息=================
// 返回匹配结果
func (this *ExtDesk) handleGameAuto(p *ExtPlayer, d *DkInMsg) {
	DebugLog("===========handleGameAuto")
	p.LiXian = false
	p.Online = false

	// 初始化
	p.Init()

	// fmt.Println("收到匹配请求,player head: ", p.Player.Head)
	// 发送匹配成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:     MSG_GAME_AUTO_REPLY,
		Result: 0,
	})
}

// 请求桌子信息
func (this *ExtDesk) handleDeskInfo(p *ExtPlayer, d *DkInMsg) {
	DebugLog("============接收到请求游戏桌子信息")

	// 发送桌子信息
	this.sendDeskInfo(p)

}

// 用户掉线，处理与退出房间一致
func (this *ExtDesk) handleDisConnect(p *ExtPlayer, d *DkInMsg) {
	DebugLog("==============handleDisConnect %d\n", p.Uid)

	if this.GetFSM(0).GetMark() == GAME_STATUS_BALANCE {
		this.LeaveByForce(p)
		return
	}

	if p.getDownBetTotal() == 0 {
		this.LeaveByForce(p)
		return
	}

	p.LiXian = true
}

// 用户请求返回
func (this *ExtDesk) handleGameBack(p *ExtPlayer, d *DkInMsg) {
	// 判断用户是否已下注或者不在座位上，允许玩家离开

	// TestLog("下注数量:%d", p.getDownBetTotal())

	var result int32 = 0
	err := ""

	if p.getDownBetTotal() == 0 {
		result = 0
		err = ""
	} else {
		result = 1
		err = "已在游戏中，退出失败。"
	}

	p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
		Id:     MSG_GAME_LEAVE_REPLY,
		Result: result,
		Err:    err,
		Uid:    p.Uid,
	})

	if result == 0 {
		this.LeaveByForce(p)
	}
}

//========================发送网络数据=================
// 发送桌子信息
func (this *ExtDesk) sendDeskInfo(p *ExtPlayer) {
	p.LiXian = false
	p.Online = true

	// fmt.Println("玩家金币：", p.Coins)

	betInfos := []GClientBetCountInfo{}

	for _, betInfo := range this.Bets { //

		sendBetInfo := GClientBetCountInfo{
			Id:           betInfo.Id,
			DownBetTotal: betInfo.DownBetValue,
		}

		betInfos = append(betInfos, sendBetInfo)
	}

	var restTime int64 = 0
	if this.GetFSM(0) != nil {
		restTime = this.GetFSM(0).getRestTime()
	}

	info := &GClientDeskInfo{
		Id:      MSG_GAME_INFO_RDESKINFO,
		Result:  0,
		FangHao: this.Mark,
		JuHao:   this.JuHao,
		Bets:    betInfos,

		BetLevels: gameConfig.LimitInfo.BetLevels,
		PlayerMassage: PlayerMsg{
			Uid:          p.Uid,
			MyUserAvatar: p.Head,
			MyUserName:   p.Account,
			MyUserCoin:   p.Coins,
		},

		GameStatus:         this.GetFSM(0).GetMark(),
		GameStatusDuration: restTime / 1000,

		AreaMaxCoin: gameConfig.LimitInfo.AreaMaxCoin,
	}

	info.PlayerMassage.MyDownBets = make([]int64, gameConfig.LimitInfo.BetCount)
	for id, v := range p.DownBets {
		info.PlayerMassage.MyDownBets[int(id)] = v
	}

	info.AreaCoin = make([]int64, gameConfig.LimitInfo.BetCount)
	for id, v := range this.Bets {
		info.AreaCoin[id] = v.DownBetValue
	}

	DebugLog("发送游戏桌子")
	DebugLog("", info)

	p.SendNativeMsg(MSG_GAME_INFO_RDESKINFO, info)
}

// 发送桌子状态
func (this *ExtDesk) SendGameState(gameStatus int, gameStatusDuration int64) {
	for _, v := range this.Players {
		if v.Online {
			v.SendNativeMsg(MSG_GAME_INFO_NSTATUS_CHANGE, &GSGameStatusInfo{
				Id:                 MSG_GAME_INFO_NSTATUS_CHANGE,
				Result:             0,
				GameStatus:         gameStatus,
				GameStatusDuration: gameStatusDuration / 1000,
			})
		}
	}
}

// 发送通知信息
func (this *ExtDesk) SendNotice(id int, info interface{}, filterRobot bool, hand func(p *ExtPlayer, hand interface{}) interface{}) {
	for _, v := range this.Players {
		if v.Online {
			if v.Robot == false || filterRobot == false {
				if hand != nil {
					v.SendNativeMsg(id, hand(v, info))
				} else {
					v.SendNativeMsg(id, info)
				}
			}
		}
	}
}

//========================对外方法=================
// 玩家下注, 并处理网络消息
func (this *ExtDesk) UserDownBet(p *ExtPlayer, betIdx int, coinIdx int) {
	if p == nil {
		return
	}
	//
	var result int32 = 0
	var err string = ""

	// 判断金币下注索引
	if coinIdx >= 0 && coinIdx >= gameConfig.LimitInfo.BetLevelCount {
		p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, &GSDownBet{
			Id:     MSG_GAME_INFO_RDOWNBET,
			Result: 2,
			Err:    "金额错误,下注失败",
		})
		return
	}

	betInfo, ok := this.Bets[betIdx]
	DebugLog("是否有该座位：", betIdx, ok)
	if ok == false {
		p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, &GSDownBet{
			Id:     MSG_GAME_INFO_RDOWNBET,
			Result: 2,
			Err:    "金额错误,下注失败",
		})
		return
	}
	var downCoin int64 = 0
	downCoin = gameConfig.LimitInfo.BetLevels[coinIdx]

	// 下注的金币是否足够
	if downCoin > p.Coins {
		p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, &GSDownBet{
			Id:     MSG_GAME_INFO_RDOWNBET,
			Result: 3,
			Err:    "金币不足,下注失败",
		})
		return
	}

	// 判断是否超过了单区域下注最大值
	value, ok := p.DownBets[betIdx]
	var areaValue int64 = 0
	if ok {
		areaValue = value + downCoin
	} else {
		areaValue = downCoin
	}

	isOutArea := false

	isOutArea = areaValue > int64(gameConfig.LimitInfo.AreaMaxCoin)

	if isOutArea {

		p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, &GSDownBet{
			Id:     MSG_GAME_INFO_RDOWNBET,
			Result: 4,
			Err:    "下注金币超过单个区域的最大值",
		})

		return
	}

	if ok {
		p.DownBets[betIdx] += downCoin
	} else {
		p.DownBets[betIdx] = downCoin
	}

	p.Coins -= downCoin
	betInfo.DownBetValue += downCoin //区域下注金额增加

	if !p.Robot {
		betInfo.UserBetValue += downCoin
	}

	this.Bets[betIdx] = betInfo //更新区域属性

	msg := GSDownBet{
		Id:     MSG_GAME_INFO_RDOWNBET,
		Result: result,
		Err:    err,
	}
	fmt.Println(msg)
	// 玩家请求下注回复
	p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, &msg)

	// 通知所有玩家有新的下注
	this.SendNotice(MSG_GAME_INFO_NDOWNBET, &GNDownBet{
		Id:             MSG_GAME_INFO_NDOWNBET,
		Uid:            p.Uid,
		BetIdx:         betIdx,
		Coin:           downCoin,
		CoinIdx:        coinIdx,
		AreaCoin:       this.Bets[betIdx].DownBetValue,
		PlayerAreaCoin: p.DownBets[betIdx],
	}, true, nil)
}

// 结算时调用，删除所有离线玩家
func (this *ExtDesk) RemoveAllOfflineAndExistSeat() {
	undoWarning := int32(gameConfig.GameBetDownUndo.Warning)
	undoExit := int32(gameConfig.GameBetDownUndo.Exit)

	plen := len(this.Players)
	for i := 0; i < plen; i++ {
		p := this.Players[i]

		isDelete := false
		if p.LiXian == true {
			isDelete = true
		} else if p.UnbetsCount == undoWarning {
			p.SendNativeMsg(MSG_GAME_INFO_NTIPS, &GSTips{
				Id:  MSG_GAME_INFO_NTIPS,
				Msg: fmt.Sprintf("%d局未下注，退出房间", undoExit-undoWarning),
			})
		} else if p.UnbetsCount == undoExit {
			isDelete = true
		}

		if isDelete {
			this.Desk.DelPlayer(p.Uid)
			this.Desk.DeskMgr.LeaveDo(p.Uid)
			plen -= 1
			i -= 1
		}
		p.ResetExtPlayer()
	}
}

// 统计桌面金币，并分配牌，实现输赢控制器
func (this *ExtDesk) allotCard() int {
	var blackNum int = 0
	var blackBets map[int]int64 = make(map[int]int64)

	for _, v := range this.Players {
		hierarchyRate := GetRateByHierarchyId(v.HierarchyId)
		DebugLog("层级概率：", hierarchyRate)
		if hierarchyRate == -1 {
			blackNum++
			for id, value := range v.DownBets {
				blackBets[id] += value
			}
		}
	}

	var resultAry []int = make([]int, 0)

	if blackNum > 0 { // 黑名单玩家处理
		var blackCoin int64 = 0
		var LoseCoin map[int]int64 = make(map[int]int64)
		for id, v := range blackBets {
			blackCoin += v
			LoseCoin[id] = (v * int64(CarTypeMultiple[id]*10)) / 10
		}

		for id, v := range LoseCoin {
			if v < blackCoin {
				resultAry = append(resultAry, id)
			}
		}
	}
	DebugLog("黑名单数组区域：", resultAry)

	//库存值判断
	var winResult []int
	intervalRate := GetRateByInterval()
	DebugLog("库存值：", intervalRate)
	rnum, _ := GetRandomNum(0, 100)
	if rnum < int(intervalRate*10000/100) { //进入库存值，玩家输
		var allCoin int64 = 0
		var LoseCoin map[int]int64 = make(map[int]int64)
		for id, v := range this.Bets {
			allCoin += v.UserBetValue
			LoseCoin[id] = (v.UserBetValue * int64(CarTypeMultiple[id]*10)) / 10
		}

		for id, _ := range LoseCoin {
			if allCoin > LoseCoin[id] {
				winResult = append(winResult, id)
			}
		}
	}

	var endResult []int
	if len(winResult) > 0 { // 获取玩家输集合
		for i := 0; i < len(resultAry); i++ {
			for j := 0; j < len(winResult); j++ {
				if resultAry[i] == winResult[j] {
					endResult = append(endResult, resultAry[i])
				}
			}
		}

		if len(endResult) == 0 && len(resultAry) > 0 { // 优先黑名单
			endResult = resultAry
		} else if len(endResult) == 0 && len(resultAry) == 0 { // 没有黑名单则选取胜利区
			endResult = winResult
		}
	} else { // 没进库存
		if len(resultAry) > 0 { // 处理黑名单
			endResult = resultAry
		} else { // 纯随机
			endResult = []int{0, 1, 2, 3, 4, 5, 6, 7}
		}
	}

	DebugLog("结果数组：", endResult)
	num, _ := GetRandomNum(0, len(endResult))

	return endResult[num]
}

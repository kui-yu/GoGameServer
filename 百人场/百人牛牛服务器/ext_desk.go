package main

import (
	"encoding/json"
	"fmt"
	"logs"
	// "math"
)

type ExtDesk struct {
	Desk
	Gchh GClientHallHistory
	Mark string // 房间标志
	// DownCardIdx    uint8                 // 牌发到的位置
	DownCards      []uint8               // 所有的牌
	CardGroupArray map[int]CardGroupInfo // 玩家和庄家的牌 庄家牌索引最后一个4

	Seats     map[int]SeatInfo      // 座位信息
	RankUsers []int64               // 展示用户
	ManyUsers []GClientManyUserInfo // 更多玩家

	fsms       map[int]FSMBase // fsm状态机集合
	currentFSM FSMBase         // 当前状态机
	upFSM      FSMBase         // 上一个状态机

}

func (this *ExtDesk) InitExtData() {
	logs.Debug("初始化ext_desk")
	this.Gchh.Trend = make([][]int, 4)
	this.Mark = FormatDeskId(this.Id, GCONFIG.GradeType)

	this.CardGroupArray = make(map[int]CardGroupInfo)
	this.Seats = make(map[int]SeatInfo)

	// 初始化下注区域
	seatCount := gameConfig.GameLimtInfo.SeatCount
	for i := 0; i < seatCount; i++ {
		this.Seats[i] = SeatInfo{
			Id:              i,
			UserId:          0,
			DownBetValue:    0,
			UserBetValue:    0,
			MinDownBetCount: 0,
		}
	}

	this.fsms = make(map[int]FSMBase)

	this.addFSM(GAME_STATUS_WAITSTART, new(FSMWaitStart))
	this.addFSM(GAME_STATUS_SEATBET, new(FSMSeatBet))
	this.addFSM(GAME_STATUS_FACARD, new(FSMFaCard))
	this.addFSM(GAME_STATUS_DOWNBTES, new(FSMDownBets))
	this.addFSM(GAME_STATUS_OPENCARD, new(FSMOpenCard))
	this.addFSM(GAME_STATUS_BALANCE, new(FSMBalance))
	this.addFSM(GAME_STATUS_SHUFFLECARD, new(FSMShuffleCard))
	this.addListener()

	// 启动状态
	this.RunFSM(GAME_STATUS_WAITSTART)
}

func (this *ExtDesk) GetFSM(mark int) FSMBase {
	if mark != 0 {
		return this.fsms[mark]
	}
	return this.currentFSM
}

func (this *ExtDesk) RunFSM(mark int) {
	var upMark int = 0
	if this.currentFSM != nil {
		this.upFSM = this.currentFSM
		this.upFSM.Leave()
		upMark = this.upFSM.GetMark()
	}

	this.currentFSM = this.GetFSM(mark)
	this.currentFSM.Run(upMark)
}

// deskmgr -> 直接调用
func (this *ExtDesk) HandleReconnect(p *ExtPlayer, d2 *DkInMsg) {
	DebugLog("============接收到重新连接")
	p.SendNativeMsg(MSG_GAME_RECONNECT_REPLY, &GReConnectReply{
		Id:       MSG_GAME_RECONNECT_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})

	// 重连消息
	p.SendNativeMsg(MSG_GAME_INFO_RECONNECT, &GReConnectReply{
		Id:     MSG_GAME_INFO_RECONNECT,
		Result: 0,
	})
	this.sendDeskInfo(p)
}

// 玩家离开
func (this *ExtDesk) Leave(p *ExtPlayer) {
	DebugLog(">>>>>>>>>>>>>>>>>Leave 玩家离开")
	this.handleGameBack(p, nil)
}

func (this *ExtDesk) addListener() {
	this.Handle[MSG_GAME_AUTO] = this.handleGameAuto // 玩家进入返回匹配结果
	//this.Handle[MSG_GAME_QDESKINFO] = this.handleDeskInfo   // 请求游戏桌子信息
	this.Handle[MSG_GAME_DISCONNECT] = this.handleDisConnect  // 用户掉线，处理与退出房间一致
	this.Handle[MSG_GAME_QHISTORY] = this.handleGameHistory   // 请求走势
	this.Handle[MSG_GAME_QMANYUSER] = this.handleGameManyUser // 请求更多玩家信息
	this.Handle[MSG_GAME_QBACK] = this.handleGameBack         // 请求返回大厅
	this.Handle[MSG_GAME_QSEATINFO] = this.handleSeatInfo     // 请求座位信息
	this.Handle[MSG_GAME_INFO_GET_RECORD] = this.HandleGetRecord
}

func (this *ExtDesk) addFSM(mark int, fsm FSMBase) {
	fsm.InitFSM(mark, this)
	this.fsms[mark] = fsm
}

//========================大厅消息=================
// 返回匹配结果
func (this *ExtDesk) handleGameAuto(p *ExtPlayer, d *DkInMsg) {
	DebugLog("===========handleGameAuto")
	// p.Online = false // 设置玩家离线
	// p.LiXian = false

	// 初始化
	p.Init()
	// 发送匹配成功
	p.SendNativeMsg(MSG_GAME_AUTO_REPLY, &GAutoGameReply{
		Id:       MSG_GAME_AUTO_REPLY,
		CostType: GetCostType(),
		Result:   0,
	})
	// 发送桌子信息
	this.sendDeskInfo(p)

	this.GetFSM(0).onUserOnline(p)
}

// 请求桌子信息
func (this *ExtDesk) handleDeskInfo(p *ExtPlayer, d *DkInMsg) {
	DebugLog("============接收到请求游戏桌子信息")

	// 发送桌子信息
	this.sendDeskInfo(p)

	this.GetFSM(0).onUserOnline(p)
}

// 请求走势图
func (this *ExtDesk) handleGameHistory(p *ExtPlayer, d *DkInMsg) {
	var trendHistory map[int][]CardGroupType = make(map[int][]CardGroupType)

	for _, seat := range this.Seats {
		trendHistory[seat.Id] = seat.TrendHistory
	}

	p.SendNativeMsg(MSG_GAME_RHISTORY, &GClientRGameHistory{
		Id:       MSG_GAME_RHISTORY,
		Result:   0,
		Historys: trendHistory,
	})
}

// 请求更多玩家信息
func (this *ExtDesk) handleGameManyUser(p *ExtPlayer, d *DkInMsg) {
	ManyUsers := []GClientManyUserInfo{}
	for _, user := range this.ManyUsers {
		name := user.NickName
		if user.Uid != p.Uid {
			name = MarkNickName(name)
		}
		ManyUsers = append(ManyUsers, GClientManyUserInfo{
			Uid:       user.Uid,
			NickName:  name,
			Avatar:    user.Avatar,
			Coin:      user.Coin,
			GameCount: user.GameCount,
			Victory:   user.Victory,
			DownBet:   user.DownBet,
		})
	}

	p.SendNativeMsg(MSG_GAME_RMANYUSER, &GClientRManyUser{
		Id:        MSG_GAME_RMANYUSER,
		Result:    0,
		ManyUsers: ManyUsers,
	})
}

// 用户掉线，处理与退出房间一致
func (this *ExtDesk) handleDisConnect(p *ExtPlayer, d *DkInMsg) {
	DebugLog("==============handleDisConnect %d\n", p.Uid)
	this.GetFSM(0).onUserOffline(p)

	// 座位玩家离开
	if _, ok := this.findUserSeatDown(p.Uid); ok {
		if this.GetFSM(0).GetMark() == GAME_STATUS_BALANCE ||
			this.GetFSM(0).GetMark() == GAME_STATUS_SHUFFLECARD ||
			this.GetFSM(0).GetMark() == GAME_STATUS_WAITSTART {
			if this.UserSeatUp(p, false, false) {
				this.LeaveByForce(p)
			} else {
				p.LiXian = true
			}

		} else {
			p.LiXian = true
		}
		return
	}

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

	TestLog("下注数量:%d", p.getDownBetTotal())

	var result int32 = 0
	err := ""

	if p.getDownBetTotal() == 0 {
		if _, ok := this.findUserSeatDown(p.Uid); ok {
			if this.UserSeatUp(p, false, false) == true {
				result = 0
				err = ""
			} else {
				result = 2
				err = "已入座，退出失败。"
			}
		} else {
			result = 0
			err = ""
		}

	} else {
		result = 1
		err = "已在游戏中，退出失败。"
	}

	p.SendNativeMsg(MSG_GAME_RBACK, &GClientRGameBack{
		Id:     MSG_GAME_RBACK,
		Result: result,
		Err:    err,
	})

	if result == 0 {
		this.LeaveByForce(p)

		this.RefreshManyUsers()
		this.RefreshRank()
	} else {
		p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Result: result,
			Err:    err,
		})
	}
}

func (this *ExtDesk) handleSeatInfo(p *ExtPlayer, d *DkInMsg) {
	seatInfos := []GClientSeatInfo{}

	for _, seatInfo := range this.Seats {

		sendSeatInfo := GClientSeatInfo{
			Id:            seatInfo.Id,
			DownBetTotal:  seatInfo.DownBetValue,
			SeatDownCount: 0,
		}

		if seatInfo.UserId != 0 {
			extPlayer := this.GetPlayer(seatInfo.UserId)
			if extPlayer != nil {
				sendSeatInfo.UserId = seatInfo.UserId
				sendSeatInfo.Name = extPlayer.Account
				sendSeatInfo.Avatar = extPlayer.Head
				sendSeatInfo.SeatDownCount = seatInfo.SeatDownCount + 1

				if seatInfo.UserId != p.Uid {
					sendSeatInfo.Name = MarkNickName(sendSeatInfo.Name)
				}
			}
		}
		seatInfos = append(seatInfos, sendSeatInfo)
	}

	p.SendNativeMsg(MSG_GAME_RSEATINFO, struct {
		Id   int
		Data []GClientSeatInfo
	}{
		Id:   MSG_GAME_RSEATINFO,
		Data: seatInfos,
	})
}

//========================发送网络数据=================
// 发送桌子信息
func (this *ExtDesk) sendDeskInfo(p *ExtPlayer) {
	p.LiXian = false

	seatInfos := []GClientSeatInfo{}

	for _, seatInfo := range this.Seats {

		sendSeatInfo := GClientSeatInfo{
			Id:            seatInfo.Id,
			DownBetTotal:  seatInfo.DownBetValue,
			UserId:        0,
			SeatDownCount: 0,
		}

		if seatInfo.UserId != 0 {
			extPlayer := this.GetPlayer(seatInfo.UserId)
			if extPlayer != nil {
				sendSeatInfo.UserId = seatInfo.UserId
				sendSeatInfo.Name = extPlayer.Account
				sendSeatInfo.Avatar = extPlayer.Head
				sendSeatInfo.SeatDownCount = seatInfo.SeatDownCount + 1

				if seatInfo.UserId != p.Uid {
					sendSeatInfo.Name = MarkNickName(sendSeatInfo.Name)
				}
			}
		}
		seatInfos = append(seatInfos, sendSeatInfo)
	}

	rankUsers := []GClientRankInfo{}
	for _, userId := range this.RankUsers {
		extPlayer := this.GetPlayer(userId)

		if extPlayer != nil {
			rankInfo := GClientRankInfo{
				UserId: userId,
				Avatar: extPlayer.Head,
			}

			rankUsers = append(rankUsers, rankInfo)
		}

	}

	var restTime int64 = 0
	if this.GetFSM(0) != nil {
		restTime = this.GetFSM(0).getRestTime()
	}
	info := &GClientDeskInfo{
		Id:           MSG_GAME_RDESKINFO,
		Result:       0,
		FangHao:      this.Mark,
		JuHao:        this.JuHao,
		Seats:        seatInfos,
		RankUsers:    rankUsers,
		BetLevels:    gameConfig.GameLimtInfo.BetLevels,
		MyUserAvatar: p.Head,
		MyUserName:   p.Account,
		MyUserCoin:   p.Coins,
		MyDownBets:   p.DownBets,

		GameStatus:          this.GetFSM(0).GetMark(),
		GameStatusDuration:  restTime,
		CardGroupArray:      this.CardGroupArray,
		SeatDownMinCoinCond: gameConfig.GameLimtInfo.SeatDownCond,
		SeatDownMinBetCond:  gameConfig.GameLimtInfo.SeatDownMinBet,
		AreaMaxCoin:         gameConfig.GameLimtInfo.AreaMaxCoin,
		AreaMaxCoinDownSeat: gameConfig.GameLimtInfo.AreaMaxCoinDownSeat,
		SeatUpTotalCount:    gameConfig.GameLimtInfo.SeatDownNum,
	}

	DebugLog("发送游戏桌子")
	DebugLog("", info)

	p.SendNativeMsg(MSG_GAME_RDESKINFO, info)

	this.RefreshManyUsers()
	this.RefreshRank()
}

// 发送桌子状态
func (this *ExtDesk) SendGameState(gameStatus int, gameStatusDuration int64) {
	for _, v := range this.Players {
		if v.LiXian == false {
			v.SendNativeMsg(MSG_GAME_NSTATUS_CHANGE, &GClientGameStatusInfo{
				Id:                 MSG_GAME_NSTATUS_CHANGE,
				Result:             0,
				GameStatus:         gameStatus,
				GameStatusDuration: gameStatusDuration,
			})
		}
	}
}

// 发送通知信息
func (this *ExtDesk) SendNotice(id int, info interface{}, filterRobot bool, hand func(p *ExtPlayer, hand interface{}) interface{}) {
	for _, v := range this.Players {
		if v.LiXian == false {
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
// 查询玩家是否在座位
func (this *ExtDesk) findUserSeatDown(uid int64) (SeatInfo, bool) {
	for _, v := range this.Seats {
		if v.UserId == uid {
			return v, true
		}
	}

	return SeatInfo{}, false
}

// 玩家坐下, 并处理网络消息
func (this *ExtDesk) UserSeatDown(p *ExtPlayer, seatId int) bool {
	DebugLog("=============调用玩家座位坐下 %d\n", p.Uid)
	var result int32 = 0
	err := ""
	seatDownCond := gameConfig.GameLimtInfo.SeatDownCond
	seatCount := gameConfig.GameLimtInfo.SeatCount
	if p.Coins < int64(seatDownCond) {
		result = 1
		err = fmt.Sprintf("携带金币不足%d", seatDownCond)
	}

	if seatId >= seatCount {
		result = 2
		err = "座位不存在"
	}

	sInfo, ok := this.Seats[seatId]
	if ok == false {
		result = 4
		err = "座位不存在"
	}
	if sInfo.UserId != 0 {
		result = 3
		err = "已有玩家坐下"
	}

	if _, ok := this.findUserSeatDown(p.Uid); ok {
		result = 5
		err = "你已在座位上"
	}

	p.SendNativeMsg(MSG_GAME_RSEATDOWN, &GClientRSeatDown{
		Id:     MSG_GAME_RSEATDOWN,
		Result: result,
		Err:    err,
	})

	// 玩家坐下成功，更新数据，并通知所有玩家有玩家坐下
	if result == 0 {
		sInfo.UserId = p.Uid
		sInfo.SeatDownCount = 0
		sInfo.MinDownBetCount = 0

		this.Seats[seatId] = sInfo
		DebugLog("=============通知有玩家座位坐下 %d\n", p.Uid)

		name := p.Account
		this.SendNotice(MSG_GAME_NSEATDOWN, &GClientSeatDownChange{
			Id:            MSG_GAME_NSEATDOWN,
			Type:          0,
			SeatId:        seatId,
			GameStatus:    this.GetFSM(0).GetMark(),
			OldUserId:     0,
			NewUserId:     p.Uid,
			NewUserAvatar: p.Head,
			NewUserName:   name,
		}, false, func(p *ExtPlayer, temp interface{}) interface{} {
			info := temp.(*GClientSeatDownChange)
			if info.NewUserId != p.Uid {
				info.NewUserName = MarkNickName(name)
			} else {
				info.NewUserName = name
			}
			return info
		})

		return true
	}
	return false
}

// 玩家站立， 并处理消息
func (this *ExtDesk) UserSeatUp(p *ExtPlayer, isSendRetMsg bool, isForceUp bool) bool {
	DebugLog("=============调用座位玩家离开", p.Uid)
	seatIdx := -1
	for k, seatInfo := range this.Seats {
		if seatInfo.UserId == p.Uid {
			seatIdx = k
			break
		}
	}

	var result int32 = 0 //0成功，其他失败
	var err string = ""  //
	if seatIdx == -1 {
		result = 1
		err = "玩家不在座位上"
	}

	if (result == 0 && isForceUp == false && this.Seats[seatIdx].SeatDownCount < gameConfig.GameLimtInfo.SeatDownNum) &&
		(p.Coins >= int64(gameConfig.GameLimtInfo.SeatDownCond)) {
		result = 2
		err = fmt.Sprintf("你还没玩够%d/%d局呢", this.Seats[seatIdx].SeatDownCount, gameConfig.GameLimtInfo.SeatDownNum)

		DebugLog("=============座位玩家再玩 局离开", gameConfig.GameLimtInfo.SeatDownNum-this.Seats[seatIdx].SeatDownCount)
	}

	if result == 0 {
		seatInfo, ok := this.Seats[seatIdx]
		if ok == false {
			return false
		}
		seatInfo.UserId = 0
		seatInfo.SeatDownCount = 0
		seatInfo.MinDownBetCount = 0
		this.Seats[seatIdx] = seatInfo

		DebugLog("=============调用座位玩家离开成功 %d\n", p.Uid)
		this.SendNotice(MSG_GAME_NSEATDOWN, &GClientSeatDownChange{
			Id:         MSG_GAME_NSEATDOWN,
			Type:       2,       // 0添加 1修改 2删除
			SeatId:     seatIdx, // 座位号
			GameStatus: this.GetFSM(0).GetMark(),
			OldUserId:  p.Uid, // 玩家Id
			NewUserId:  0,     // 玩家Id
		}, false, nil)
	}

	if isSendRetMsg {
		p.SendNativeMsg(MSG_GAME_RSEATUP, &GClientRSeatUp{
			Id:     MSG_GAME_RSEATUP,
			Result: result,
			Err:    err,
		})
	}

	return result == 0
}

// 玩家下注, 并处理网络消息
func (this *ExtDesk) UserDownBet(p *ExtPlayer, seatId int, betIdx int, isSendRetMsg bool) {
	if p == nil {
		return
	}
	// 判断是否是抢座状态
	var result int32 = 0
	var err string = ""
	if this.GetFSM(0).GetMark() == GAME_STATUS_SEATBET {
		_, ok := this.findUserSeatDown(p.Uid)
		if !ok {
			if isSendRetMsg {
				p.SendNativeMsg(MSG_GAME_RDOWNBET, &GClientRDownBet{
					Id:     MSG_GAME_RDOWNBET,
					Result: 1,
					Err:    "玩家不在座位,下注失败",
				})
			}
			return
		}
	}
	if this.GetFSM(0).GetMark() != GAME_STATUS_SEATBET &&
		this.GetFSM(0).GetMark() != GAME_STATUS_DOWNBTES {
		p.SendNativeMsg(MSG_GAME_RDOWNBET, &GClientRDownBet{
			Id:     MSG_GAME_RDOWNBET,
			Result: 1,
			Err:    "不是下注阶段,下注失败",
		})
	}
	// 判断金币下注索引
	if betIdx >= 0 && betIdx >= gameConfig.GameLimtInfo.BetLevelCount {
		if isSendRetMsg {
			p.SendNativeMsg(MSG_GAME_RDOWNBET, &GClientRDownBet{
				Id:     MSG_GAME_RDOWNBET,
				Result: 2,
				Err:    "金额错误,下注失败",
			})
		}
		return
	}

	seatInfo, ok := this.Seats[seatId]
	if ok == false {
		if isSendRetMsg {
			p.SendNativeMsg(MSG_GAME_RDOWNBET, &GClientRDownBet{
				Id:     MSG_GAME_RDOWNBET,
				Result: 2,
				Err:    "金额错误,下注失败",
			})
		}
		return
	}
	var downCoin int64 = 0
	if betIdx != -1 {
		downCoin = gameConfig.GameLimtInfo.BetLevels[betIdx]
	} else {
		downCoin = int64(gameConfig.GameLimtInfo.SeatDownMinBet)
	}

	// 下注的金币是否足够
	// if GetCostType() == 1 {
	if downCoin > p.Coins {
		if isSendRetMsg {
			p.SendNativeMsg(MSG_GAME_RDOWNBET, &GClientRDownBet{
				Id:     MSG_GAME_RDOWNBET,
				Result: 3,
				Err:    "金币不足,下注失败",
			})
		}
		return
	}
	// }

	// 下注的金币*赔的倍数 <自身携带的金币  体验场
	// if GetCostType() == 1 {
	if (int64(p.getDownBetTotal())+downCoin)*int64(gameConfig.GameLimtInfo.CompDownbetDouble-1) > (p.Coins - downCoin) {
		if isSendRetMsg {
			p.SendNativeMsg(MSG_GAME_RDOWNBET, &GClientRDownBet{
				Id:     MSG_GAME_RDOWNBET,
				Result: 3,
				Err:    fmt.Sprintf("自身金币不足下注金币的%d倍,下注失败", gameConfig.GameLimtInfo.CompDownbetDouble),
			})
		}
		return
	}
	// }

	// 判断是否超过了单区域下注最大值
	value, ok := p.DownBets[uint8(seatId)]
	var areaValue int64 = 0
	if ok {
		areaValue = value + downCoin
	} else {
		areaValue = downCoin
	}

	isOutArea := false
	if _, ok := this.findUserSeatDown(p.Uid); ok {
		isOutArea = areaValue > int64(gameConfig.GameLimtInfo.AreaMaxCoinDownSeat)
	} else {
		isOutArea = areaValue > int64(gameConfig.GameLimtInfo.AreaMaxCoin)
	}
	// if GetCostType() == 1 {
	if isOutArea {
		if isSendRetMsg {
			p.SendNativeMsg(MSG_GAME_RDOWNBET, &GClientRDownBet{
				Id:     MSG_GAME_RDOWNBET,
				Result: 4,
				Err:    "下注金币超过单个区域的最大值",
			})
		}
		return
	}
	// }

	if ok {
		p.DownBets[uint8(seatId)] += downCoin
	} else {
		p.DownBets[uint8(seatId)] = downCoin
	}

	p.Coins -= downCoin
	seatInfo.DownBetValue += downCoin

	if !p.Robot {
		seatInfo.UserBetValue += downCoin
	}

	this.Seats[seatId] = seatInfo

	// 玩家请求下注回复
	if !p.Robot {
		if isSendRetMsg {
			p.SendNativeMsg(MSG_GAME_RDOWNBET, &GClientRDownBet{
				Id:     MSG_GAME_RDOWNBET,
				Coins:  p.Coins,
				Result: result,
				Err:    err,
			})
		}
	}

	// 通知所有玩家有新的下注
	this.SendNotice(MSG_GAME_NDOWNBET, &GClientNDownBet{
		Id:      MSG_GAME_NDOWNBET,
		Uid:     p.Uid,
		SeatIdx: seatId,
		Coin:    uint32(downCoin),
		CoinIdx: betIdx,
	}, true, nil)
}

// 结算时调用，删除所有离线玩家,以及座位金币不足站立和3局最低下注
func (this *ExtDesk) RemoveAllOfflineAndExistSeat() {
	undoWarning := int32(gameConfig.GameBetDownUndo.Warning)
	undoExit := int32(gameConfig.GameBetDownUndo.Exit)

	plen := len(this.Players)
	for i := 0; i < plen; i++ {
		p := this.Players[i]
		seatUser, seatok := this.findUserSeatDown(p.Uid)

		isDelete := false
		if p.LiXian == true {
			if seatok &&
				p.Coins >= int64(gameConfig.GameLimtInfo.SeatDownCond) &&
				seatUser.SeatDownCount < gameConfig.GameLimtInfo.SeatDownNum {
				isDelete = false
			} else {
				isDelete = true
			}
		} else if p.UnbetsCount == undoWarning {
			p.SendNativeMsg(MSG_GAME_NTIPS, &GClientTips{
				Id:   MSG_GAME_NTIPS,
				Code: GCLIENT_TIPS_NOTBET,
				Msg:  fmt.Sprintf("%d局未下注，退出房间", undoExit-undoWarning),
			})
		} else if p.UnbetsCount == undoExit {
			isDelete = true
		}

		// 检查是否在座位上
		if isDelete == false {
			if seatok {
				if p.Coins < int64(gameConfig.GameLimtInfo.SeatDownCond) {
					this.UserSeatUp(p, true, false)
				} else if p.getDownBetTotal() == int64(gameConfig.GameLimtInfo.SeatDownMinBet) {
					seatUser.MinDownBetCount++
					this.Seats[seatUser.Id] = seatUser

					if seatUser.MinDownBetCount >= gameConfig.GameLimtInfo.SeatDownAutoUp {
						this.UserSeatUp(p, true, true)
					}
				} else if seatUser.MinDownBetCount != 0 {
					seatUser.MinDownBetCount = 0
					this.Seats[seatUser.Id] = seatUser
				}
			}
		}

		if isDelete {
			if seatok {
				this.UserSeatUp(p, false, false)
			}
			this.LeaveByForce(p)
			plen -= 1
			i -= 1
		}
	}
}

// 统计桌面金币，并分配牌，实现输赢控制器
func (this *ExtDesk) allotCard() {

	//1.通杀概率 改黑名单概率
	killRate := G_DbGetGameServerData.GameConfig.MiniRate
	logs.Debug("通杀概率", killRate)
	kill, _ := GetRandomNum(0, 100)
	if kill < int(killRate*10000/100) && GCONFIG.GradeType != 6 {
		this.InAllWinAllot()
		return
	}

	//2.点杀层
	seatCount := gameConfig.GameLimtInfo.SeatCount

	var noPlayerArea []int
	for i := 0; i < seatCount; i++ {
		if this.Seats[i].UserBetValue == 0 {
			noPlayerArea = append(noPlayerArea, i)
		}
	}

	//如果有玩家下注，进入库存点杀判断   体验场不点杀
	if len(noPlayerArea) < seatCount && GetCostType() == 1 {
		//黑名单位置
		var blackArea []int

		//初始化点杀玩家区域金额
		var blackAreaBets []int64
		for i := 0; i < seatCount; i++ {
			blackAreaBets = append(blackAreaBets, 0)
		}

		var blackNum int
		for _, v := range this.Players {
			hierarchyRate := GetRateByHierarchyId(v.HierarchyId)
			if hierarchyRate == -1 {
				for area, bet := range v.DownBets {
					if bet > 0 {
						blackAreaBets[area] += bet
						blackNum++
					}
				}
			}
		}

		if blackNum > 0 {
			var moreAreas [][]int
			//规则
			areas := [][]int{
				{0, 0, 0, 0}, {0, 0, 0, 1}, {0, 0, 1, 0}, {0, 0, 1, 1},
				{0, 1, 0, 0}, {0, 1, 0, 1}, {0, 1, 1, 0}, {0, 1, 1, 1},
				{1, 0, 0, 0}, {1, 0, 0, 1}, {1, 0, 1, 0}, {1, 0, 1, 1},
				{1, 1, 0, 0}, {1, 1, 0, 1}, {1, 1, 1, 0}, {1, 1, 1, 1},
			}
			//按规则筛选
			var winAreas [][]int
			for i := 0; i < len(areas); i++ {
				area := areas[i]
				var betAreaTotal int64
				for j := 0; j < len(area); j++ {
					if area[j] > 0 {
						betAreaTotal += blackAreaBets[j]
					} else {
						betAreaTotal -= blackAreaBets[j]
					}
				}
				if betAreaTotal >= 0 {
					winAreas = append(winAreas, area)
				}
			}
			logs.Debug("区域组合", winAreas)

			//转化庄赢数组
			for _, area := range winAreas {
				var index []int
				for j, value := range area {
					if value > 0 {
						index = append(index, j)
					}
				}
				//最多杀两个区域
				if len(index) > 0 && len(index) <= 2 {
					moreAreas = append(moreAreas, index)
				}
			}

			logs.Debug("点杀区域组合", moreAreas)
			if len(moreAreas) > 0 {
				r, _ := GetRandomNum(0, len(moreAreas))
				blackArea = append(blackArea, moreAreas[r]...)
				logs.Debug("黑名单区域", blackArea)
			}
		}

		//----------------------区域赋值-------------------------------
		noPlayerArea = []int{}
		for i := 0; i < seatCount; i++ {
			var flag bool = false
			//黑名单区域
			for _, areaId := range blackArea {
				if i == areaId {
					flag = true
				}
			}

			if !flag && this.Seats[i].UserBetValue == 0 {
				noPlayerArea = append(noPlayerArea, i)
			}
		}

		//其他正常区域
		var otherArea []int
		for i := 0; i < seatCount; i++ {
			var flag bool = false
			//黑名单区域
			for _, areaId := range blackArea {
				if i == areaId {
					flag = true
				}
			}
			//没有玩家区域
			for _, areaId := range noPlayerArea {
				if i == areaId {
					flag = true
				}
			}
			if !flag {
				otherArea = append(otherArea, i)
			}
		}

		//记录最大值
		var getMaxGroupType CardGroupType
		var getMaxCard uint8

		//---------------点杀区域处理-----------------
		for i := 0; i < len(blackArea); i++ {
			area := blackArea[i]
			logs.Debug("进点杀,区域", area)
			// rnum, _ := GetRandomNum(0, 100)
			// if rnum < int(killRate*10000/100) {
			var cardGroupInfo CardGroupInfo = this.CardGroupArray[area]
			for i := 0; i < 1000; i++ {
				sendPoker := this.SendCard(3)
				checkCards := append(cardGroupInfo.Cards, sendPoker...)
				groupType, maxCard, handCards := CalcCards(checkCards)
				//小等于牛6
				if groupType <= CardGroupType_Cattle_6 {
					cardGroupInfo.CardGroupType = groupType
					cardGroupInfo.Cards = handCards
					cardGroupInfo.MaxCard = maxCard
					break
				}
				//回收
				this.RecoverCard(sendPoker)
				//重新打乱牌组
				this.UpsetCard()
			}
			this.CardGroupArray[area] = cardGroupInfo
			//判断最大牌
			if cardGroupInfo.CardGroupType > getMaxGroupType {
				getMaxGroupType = cardGroupInfo.CardGroupType
				getMaxCard = cardGroupInfo.MaxCard
			}
			// }
		}

		//------------其他真实下注区域按库存开启-------------------
		for i := 0; i < len(otherArea); i++ {
			area := otherArea[i]
			logs.Debug("进库存判断,区域", area)
			//2.库存概率
			intervalRate := GetRateByInterval()
			rnum, _ := GetRandomNum(0, 100)
			if rnum < int(intervalRate*10000/100) {
				var cardGroupInfo CardGroupInfo = this.CardGroupArray[area]
				for i := 0; i < 1000; i++ {
					sendPoker := this.SendCard(3)
					checkCards := append(cardGroupInfo.Cards, sendPoker...)
					groupType, maxCard, handCards := CalcCards(checkCards)
					//小等于牛6
					if groupType <= CardGroupType_Cattle_6 {
						cardGroupInfo.CardGroupType = groupType
						cardGroupInfo.Cards = handCards
						cardGroupInfo.MaxCard = maxCard
						break
					}
					//回收
					this.RecoverCard(sendPoker)
					//重新打乱牌组
					this.UpsetCard()
				}
				this.CardGroupArray[area] = cardGroupInfo
				//判断最大牌
				if cardGroupInfo.CardGroupType > getMaxGroupType {
					getMaxGroupType = cardGroupInfo.CardGroupType
					getMaxCard = cardGroupInfo.MaxCard
				}
			} else {
				rnum, _ := GetRandomNum(0, 100)
				if rnum < int(intervalRate*10000/100) {
					//按库存概率随机开牛7-9
					var cardGroupInfo CardGroupInfo = this.CardGroupArray[area]
					for i := 0; i < 1000; i++ {
						sendPoker := this.SendCard(3)
						checkCards := append(cardGroupInfo.Cards, sendPoker...)
						groupType, maxCard, handCards := CalcCards(checkCards)
						if groupType >= CardGroupType_Cattle_6 && groupType <= CardGroupType_Cattle_9 {
							cardGroupInfo.CardGroupType = groupType
							cardGroupInfo.Cards = handCards
							cardGroupInfo.MaxCard = maxCard
							break
						}
						//回收
						this.RecoverCard(sendPoker)
						//重新打乱牌组
						this.UpsetCard()
					}
					this.CardGroupArray[area] = cardGroupInfo
				} else {
					//随机出牌
					cardGroupInfo := this.CardGroupArray[area]
					sendPoker := this.SendCard(3)
					checkCards := append(cardGroupInfo.Cards, sendPoker...)
					cardGroupInfo.CardGroupType, cardGroupInfo.MaxCard, cardGroupInfo.Cards = CalcCards(checkCards)
					this.CardGroupArray[area] = cardGroupInfo
				}
			}
		}

		if getMaxGroupType > 0 && getMaxCard > 0 {
			//庄家牌
			var bankGroupInfo CardGroupInfo
			for i := 0; i < 1000; i++ {
				bankerPoker := this.SendCard(5)
				groupType, maxCard, handCards := CalcCards(bankerPoker)
				if (groupType > getMaxGroupType ||
					(groupType == getMaxGroupType && maxCard&0xF > getMaxCard&0xF) ||
					(groupType == getMaxGroupType && maxCard&0xF == getMaxCard&0xF && maxCard>>4 > getMaxCard>>4)) &&
					groupType <= CardGroupType_Cattle_WUHUA {

					bankGroupInfo.CardGroupType = groupType
					bankGroupInfo.Cards = handCards
					bankGroupInfo.MaxCard = maxCard
					break
				}
				this.RecoverCard(bankerPoker)
				//重新打乱牌组
				this.UpsetCard()
			}
			this.CardGroupArray[4] = bankGroupInfo
		} else {
			var bankGroupInfo CardGroupInfo
			bankerPoker := this.SendCard(5)
			groupType, maxCard, handCards := CalcCards(bankerPoker)
			bankGroupInfo.CardGroupType = groupType
			bankGroupInfo.Cards = handCards
			bankGroupInfo.MaxCard = maxCard

			this.CardGroupArray[4] = bankGroupInfo
		}

		//------------------no player area--------------------------
		logs.Debug("没有下注区域", noPlayerArea)
		// logs.Debug("剩余牌组", this.DownCards)
		for _, area := range noPlayerArea {

			var cardGroupInfo CardGroupInfo = this.CardGroupArray[area]
			// logs.Debug("区域牌组", cardGroupInfo.Cards)
			for i := 0; i < 1000; i++ {
				var sendPoker []int
				if len(cardGroupInfo.Cards) == 0 {
					sendPoker = this.SendCard(5)
				} else {
					sendPoker = this.SendCard(3)
				}
				checkCards := append(cardGroupInfo.Cards, sendPoker...)
				groupType, maxCard, handCards := CalcCards(checkCards)
				//大于庄家
				if groupType > CardGroupType_Cattle_7 && groupType != CardGroupType_NotCattle {

					cardGroupInfo.CardGroupType = groupType
					cardGroupInfo.Cards = handCards
					cardGroupInfo.MaxCard = maxCard
					break
				}
				//回收
				this.RecoverCard(sendPoker)
				//重新打乱牌组
				this.UpsetCard()
			}
			this.CardGroupArray[area] = cardGroupInfo
		}
	} else {
		//无人下注，纯随机
		logs.Debug("无人下注，纯随机")
		//四个区域
		for area := 0; area < seatCount; area++ {
			cardGroupInfo := this.CardGroupArray[area]
			sendPoker := this.SendCard(3)
			checkCards := append(cardGroupInfo.Cards, sendPoker...)
			cardGroupInfo.CardGroupType, cardGroupInfo.MaxCard, cardGroupInfo.Cards = CalcCards(checkCards)
			this.CardGroupArray[area] = cardGroupInfo
		}

		//庄家区域
		var bankGroupInfo CardGroupInfo
		bankerPoker := this.SendCard(5)
		bankGroupInfo.CardGroupType, bankGroupInfo.MaxCard, bankGroupInfo.Cards = CalcCards(bankerPoker)
		this.CardGroupArray[4] = bankGroupInfo
	}

	// for i := 0; i < len(this.CardGroupArray); i++ {
	// 	logs.Debug("第", this.CardGroupArray[i])
	// }
}

// 刷新排行
func (this *ExtDesk) RefreshRank() {
	this.RankUsers = []int64{}

	rankCount := gameConfig.GameLimtInfo.RankCount

	s := len(this.Players)
	if s <= rankCount {
		for i := 0; i < s; i++ {
			this.RankUsers = append(this.RankUsers, this.Players[i].Uid)
		}
	} else {
		var rankIdxs []int
		for i := 0; i < rankCount; i++ {
			idx, _ := GetRandomNum(0, s-1)
			for k := 0; k < 100; k++ {
				idx, _ = GetRandomNum(0, s-1)
				isExists := false
				for _, oldKey := range rankIdxs {
					if oldKey == idx {
						isExists = true
						break
					}
				}

				if isExists == false {
					break
				}
			}
			rankIdxs = append(rankIdxs, idx)
			this.RankUsers = append(this.RankUsers, this.Players[idx].Uid)
		}
	}

	// 随机排行榜
	rankUsers := []GClientRankInfo{}
	for _, userId := range this.RankUsers {
		extPlayer := this.GetPlayer(userId)

		rankUsers = append(rankUsers, GClientRankInfo{
			UserId: userId,
			Avatar: extPlayer.Head,
		})
	}

	senddata := struct {
		Id   int
		Data []GClientRankInfo
	}{
		Id:   MSG_GAME_NRANKLIST,
		Data: rankUsers,
	}

	this.SendNotice(MSG_GAME_NRANKLIST, &senddata, true, nil)
}

// 刷新更多在线玩家
func (this *ExtDesk) RefreshManyUsers() {
	var manyUsers []GClientManyUserInfo
	userListRecordCount := gameConfig.GameLimtInfo.UserListRecordCount
	userListCount := gameConfig.GameLimtInfo.UserListCount

	orderPlayers := this.Players
	ulen := len(orderPlayers)

	for i := 0; i < ulen; i++ {
		for j := i + 1; j < ulen; j++ {
			if orderPlayers[i].Coins < orderPlayers[j].Coins {
				orderPlayers[i], orderPlayers[j] = orderPlayers[j], orderPlayers[i]
			}
		}
	}

	olen := userListCount
	if ulen < olen {
		olen = ulen
	}

	for i := 0; i < olen; i++ {
		p := orderPlayers[i]
		manyUsers = append(manyUsers, GClientManyUserInfo{
			Uid:       p.Uid,
			NickName:  p.Account,
			Avatar:    p.Head,
			Coin:      int64(p.Coins),
			GameCount: int32(userListRecordCount),
			Victory:   p.HVictoryCount,
			DownBet:   p.HDownBetTotal,
		})
	}

	this.ManyUsers = manyUsers
}

//通杀
func (this *ExtDesk) InAllWinAllot() {
	logs.Debug("庄家通杀")

	//四个区域位置
	seatCount := gameConfig.GameLimtInfo.SeatCount

	//四个区域通杀
	var blackArea []int
	for i := 0; i < seatCount; i++ {
		blackArea = append(blackArea, i)
	}

	//庄家牌组 大于等牛6
	var bankGroupInfo CardGroupInfo
	for i := 0; i < 1000; i++ {
		bankerPoker := this.SendCard(5)
		groupType, maxCard, handCards := CalcCards(bankerPoker)
		if groupType >= CardGroupType_Cattle_6 && groupType <= CardGroupType_Cattle_WUHUA {
			bankGroupInfo.CardGroupType = groupType
			bankGroupInfo.Cards = handCards
			bankGroupInfo.MaxCard = maxCard
			break
		}

		this.RecoverCard(bankerPoker)
		//重新打乱牌组
		this.UpsetCard()
	}
	this.CardGroupArray[4] = bankGroupInfo

	//闲家牌组 小于庄家牌组
	for i := 0; i < seatCount; i++ {
		var cardGroupInfo CardGroupInfo = this.CardGroupArray[i]
		for i := 0; i < 1000; i++ {
			sendPoker := this.SendCard(3)
			checkCards := append(cardGroupInfo.Cards, sendPoker...)
			groupType, maxCard, handCards := CalcCards(checkCards)

			if (groupType == bankGroupInfo.CardGroupType && maxCard&0xF < bankGroupInfo.MaxCard&0xF) ||
				groupType < bankGroupInfo.CardGroupType {

				cardGroupInfo.CardGroupType = groupType
				cardGroupInfo.Cards = handCards
				cardGroupInfo.MaxCard = maxCard
				break
			}

			this.RecoverCard(sendPoker)
			//重新打乱牌组
			this.UpsetCard()
		}
		this.CardGroupArray[i] = cardGroupInfo
	}

	logs.Debug("通杀")
}

//纯随机发牌
func (this *ExtDesk) AllotCardRand() {
	logs.Debug("纯随机发牌")

	seatCount := gameConfig.GameLimtInfo.SeatCount
	//四个区域
	for area := 0; area < seatCount; area++ {
		cardGroupInfo := this.CardGroupArray[area]
		sendPoker := this.SendCard(3)
		checkCards := append(cardGroupInfo.Cards, sendPoker...)
		cardGroupInfo.CardGroupType, cardGroupInfo.MaxCard, cardGroupInfo.Cards = CalcCards(checkCards)
		this.CardGroupArray[area] = cardGroupInfo
	}

	//庄家区域
	var bankGroupInfo CardGroupInfo
	bankerPoker := this.SendCard(5)
	bankGroupInfo.CardGroupType, bankGroupInfo.MaxCard, bankGroupInfo.Cards = CalcCards(bankerPoker)
	this.CardGroupArray[4] = bankGroupInfo
}

//庄家输牌
func (this *ExtDesk) AllotCardBankerLose() {
	logs.Debug("庄家输牌")

	//四个区域位置
	seatCount := gameConfig.GameLimtInfo.SeatCount

	//四个区域通杀
	var blackArea []int
	for i := 0; i < seatCount; i++ {
		blackArea = append(blackArea, i)
	}

	//庄家牌组 小于牛6
	var bankGroupInfo CardGroupInfo
	for i := 0; i < 1000; i++ {
		bankerPoker := this.SendCard(5)
		groupType, maxCard, handCards := CalcCards(bankerPoker)
		if groupType < CardGroupType_Cattle_6 {
			bankGroupInfo.CardGroupType = groupType
			bankGroupInfo.Cards = handCards
			bankGroupInfo.MaxCard = maxCard
			break
		}

		this.RecoverCard(bankerPoker)
		//重新打乱牌组
		this.UpsetCard()
	}
	this.CardGroupArray[4] = bankGroupInfo

	//闲家牌组 大于庄家牌组
	for i := 0; i < seatCount; i++ {
		var cardGroupInfo CardGroupInfo = this.CardGroupArray[i]
		for i := 0; i < 1000; i++ {
			sendPoker := this.SendCard(3)
			checkCards := append(cardGroupInfo.Cards, sendPoker...)
			groupType, maxCard, handCards := CalcCards(checkCards)

			if (groupType == bankGroupInfo.CardGroupType && maxCard&0xF > bankGroupInfo.MaxCard&0xF) ||
				groupType > bankGroupInfo.CardGroupType {

				cardGroupInfo.CardGroupType = groupType
				cardGroupInfo.Cards = handCards
				cardGroupInfo.MaxCard = maxCard
				break
			}

			this.RecoverCard(sendPoker)
			//重新打乱牌组
			this.UpsetCard()
		}
		this.CardGroupArray[i] = cardGroupInfo
	}
}

func (this *ExtDesk) GetZouShi(fromId int32) {
	rsp := GGetZouShiReply{
		Id: MSG_GAME_GETZOUSHI_REPLY,
	}

	rsp.Data.SerId = this.GetServerId()
	rsp.Data.GameInfo = GameTypeDetail{GameType: int32(GCONFIG.GameType),
		RoomType:  int32(GCONFIG.RoomType),
		GradeType: int32(GCONFIG.GradeType),
	}
	currentStageTime := this.GetFSM(0).GetMark()
	switch currentStageTime {
	case GAME_STATUS_FACARD:
		this.Gchh.GameState = "发牌中"
	case GAME_STATUS_DOWNBTES:
		this.Gchh.GameState = "下注中"
	case GAME_STATUS_BALANCE:
		this.Gchh.GameState = "结算中"
	}
	this.Gchh.LimitRed = gameConfig.GameLimtInfo.AreaMaxCoin
	for i, _ := range this.Gchh.Trend {
		if len(this.Gchh.Trend[i]) == 0 {
			this.Gchh.Trend = make([][]int, 0)
			this.Gchh.WinCount = []int{0, 0, 0, 0}
			break
		}
	}
	stageTimeCounts, RemainTime := this.AllTimes()
	this.Gchh.StageTimeCount = stageTimeCounts
	// logs.Debug("游戏总时间:%v", this.Gchh.StageTimeCount)
	this.Gchh.RemainTime = RemainTime
	zs, err := json.Marshal(this.Gchh)
	if err != nil {
		logs.Error("更新走势时json转换出错")
	}
	rsp.Data.Data.ZouShi = string(zs)
	rsp.Data.Data.PlayerNum = int32(len(this.DeskMgr.MapPlayers))
	rsp.Data.Data.GameStatus = int32(currentStageTime)
	stageTime, fid := this.GetStagetime(currentStageTime)
	rsp.Data.Data.StageTime = stageTime
	rsp.Data.Data.MaxBet = int64(this.Gchh.LimitRed)
	rsp.Data.UpdateT = int64(this.GetTimerNum(fid))
	rsp.Data.GradeNumber = int32(GCONFIG.GradeNumber)
	//logs.Debug("*****************大厅走势总结构体:%v", rsp)
	this.DeskMgr.SendNativeMsgNoPlayer(MSG_GAME_GETZOUSHI_REPLY, 0, fromId, &rsp)
}

//获取状态时间及状态ID
func (this *ExtDesk) GetStagetime(current int) (int, int) {
	switch current {
	case GAME_STATUS_WAITSTART:
		return gameConfig.GameStatusTimer.WaitstartMS, gameConfig.GameStatusTimer.WaitstartId
	case GAME_STATUS_SEATBET:
		return gameConfig.GameStatusTimer.RobSeatMS, gameConfig.GameStatusTimer.RobSeatId
	case GAME_STATUS_FACARD:
		return gameConfig.GameStatusTimer.FaCardMS, gameConfig.GameStatusTimer.FaCardId
	case GAME_STATUS_DOWNBTES:
		return gameConfig.GameStatusTimer.DownBetsMS, gameConfig.GameStatusTimer.DownBetsId
	case GAME_STATUS_OPENCARD:
		return gameConfig.GameStatusTimer.OpenCardMS, gameConfig.GameStatusTimer.OpenCardId
	case GAME_STATUS_BALANCE:
		return gameConfig.GameStatusTimer.BalanceMS, gameConfig.GameStatusTimer.BalanceId
	case GAME_STATUS_SHUFFLECARD:
		return gameConfig.GameStatusTimer.ShufflecardMS, gameConfig.GameStatusTimer.ShufflecardId
	default:
		return 0, 0
	}

}

func (this *ExtDesk) AllTimes() (int, int) {
	stageTimes := gameConfig.GameStatusTimer
	allStageTime := []int{stageTimes.WaitstartMS, stageTimes.RobSeatMS, stageTimes.FaCardMS, stageTimes.DownBetsMS, stageTimes.OpenCardMS, stageTimes.BalanceMS, stageTimes.ShufflecardMS}
	// logs.Debug("全部时间:%v", allStageTime)
	var stageTimeCount int = 0 //游戏总时间

	for i, _ := range allStageTime {
		stageTimeCount += allStageTime[i]
	}
	//获取当前时间
	currentStageTime := this.GetFSM(0).GetMark()
	_, fid := this.GetStagetime(currentStageTime)
	// logs.Debug("-------------------------当前阶段:%v;当前阶段时间:%v", currentStageTime, stageMS)
	//剩余时间！！！！！
	var RemainTime int
	switch currentStageTime {
	case GAME_STATUS_WAITSTART:
		RemainTime = stageTimeCount - this.GetTimerNum(fid)
	case GAME_STATUS_SEATBET:
		RemainTime = stageTimeCount - (allStageTime[0] + this.GetTimerNum(fid))
	case GAME_STATUS_FACARD:
		RemainTime = stageTimeCount - (timeCount(allStageTime[:2]) + this.GetTimerNum(fid))
	case GAME_STATUS_DOWNBTES:
		RemainTime = stageTimeCount - (timeCount(allStageTime[:3]) + this.GetTimerNum(fid))
	case GAME_STATUS_OPENCARD:
		RemainTime = stageTimeCount - (timeCount(allStageTime[:4]) + this.GetTimerNum(fid))
	case GAME_STATUS_BALANCE:
		RemainTime = stageTimeCount - (timeCount(allStageTime[:5]) + this.GetTimerNum(fid))
	case GAME_STATUS_SHUFFLECARD:
		RemainTime = this.GetTimerNum(fid)
	default:
		RemainTime = 0
	}
	return stageTimeCount, RemainTime
}

func timeCount(times []int) int {
	var count int
	for i, _ := range times {
		count += times[i]
	}
	return count
}

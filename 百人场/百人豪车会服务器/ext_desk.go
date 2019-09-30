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

	fsms        map[int]FSMBase   // fsm状态机集合
	currentFSM  FSMBase           // 当前状态机
	upFSM       FSMBase           // 上一个状态机
	AllBetCount int               //总下注金额
	BetAgain    map[int64][]int64 //重复下注
	Car         int               //上一把开奖结果
	ChangeCar   int
	Chair       map[int]*GUserInfoByChair //座位
	// ChipList    []int64                   //筹码列表
	// MaxBet      int64                     //限红

}

func (this *ExtDesk) ResetExtDesk() {
	this.GameResult = -1
	betCount := gameConfig.LimitInfo.BetCount
	this.Bets = make(map[int]GBetCountInfo)
	for i := 0; i < betCount; i++ { //8个下注区
		this.Bets[i] = GBetCountInfo{
			Id:           i,
			DownBetValue: 0,
			UserBetValue: 0,
		}
	}
}

func (this *ExtDesk) InitExtData() {
	logs.Debug("初始化ext_dest")
	this.Mark = FormatDeskId(this.Id, GCONFIG.GradeType)
	this.BetAgain = make(map[int64][]int64)
	this.Chair = make(map[int]*GUserInfoByChair)
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
	for i := 0; i < 6; i++ {
		this.Chair[i] = &GUserInfoByChair{
			Uid:   0,
			Nick:  "",
			Head:  "",
			Coins: 0,
		}
	}
	this.fsms = make(map[int]FSMBase)

	//状态阶段
	this.addFSM(GAME_STATUS_DOWNBET, new(FSMDownBet))
	this.addFSM(GAME_STATUS_LOTTERY, new(FSMLottery))
	this.addFSM(GAME_STATUS_BALANCE, new(FSMSettle))
	this.addFSM(GAME_STATUS_READY, new(FSMReady))
	this.addListener()

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
	this.recDestInfo(p)
}

// 玩家离开 deskmgr -> 直接调用
func (this *ExtDesk) Leave(p *ExtPlayer) { // 400007
	DebugLog(">>>>>>>>>>>>>>>>>Leave 玩家离开")
	logs.Debug("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!玩家请求离开游戏")
	this.handleGameBack(p, nil)
}

func (this *ExtDesk) addListener() {
	this.Handle[MSG_GAME_AUTO] = this.handleGameAuto              // 玩家进入返回匹配结果
	this.Handle[MSG_GAME_INFO_QDESKINFO] = this.handleDeskInfo    // 请求游戏桌子信息
	this.Handle[MSG_GAME_DISCONNECT] = this.handleDisConnect      // 用户掉线，处理与退出房间一致
	this.Handle[MSG_GAME_INFO_MOREPLAYER] = this.handleMorePlayer //处理更多玩家请求
	this.Handle[MSG_GAME_INFO_BETAGAIN] = this.handleBetAgain     //重复下注
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

//玩家重复下注
func (this *ExtDesk) handleBetAgain(p *ExtPlayer, d *DkInMsg) {
	if this.GameState != GAME_STATUS_DOWNBET {
		p.SendNativeMsg(MSG_GAME_INFO_BETAGAIN_REPLAY, BetAgainReply{
			Id:     MSG_GAME_INFO_BETAGAIN_REPLAY,
			Result: 1,
			Err:    "不是下注阶段",
		})
		return
	}
	var pd1 = true

	for k, v := range this.BetAgain {
		if p.Uid == int64(k) {
			pd1 = false
			var pd bool = true
			for _, v1 := range v {
				if v1 != 0 {
					pd = false
				}
			}
			if pd {
				fmt.Println("发现玩家:", p.Nick, "上局为没有下注")
				p.SendNativeMsg(MSG_GAME_INFO_BETAGAIN_REPLAY, BetAgainReply{
					Id:      MSG_GAME_INFO_BETAGAIN_REPLAY,
					Result:  1,
					BetArea: v,
					Err:     "您上局未下注",
				})
				return
			} else {

				for i, v := range p.PAreaCoins {
					p.DownBets[i] = v
				}
				var arr []int64 = []int64{0, 0, 0, 0, 0, 0, 0, 0}
				for ix, v := range v {
					arr[ix] = v + p.PAreaCoins[ix]
				}
				var allCoins int64
				for _, v4 := range v {
					allCoins += v4
				}
				if p.Coins-allCoins < 0 {
					p.SendNativeMsg(MSG_GAME_INFO_BETAGAIN_REPLAY, BetAgainReply{
						Id:     MSG_GAME_INFO_BETAGAIN_REPLAY,
						Result: 2,
						Err:    "重复下注金币不足，请您手动下注",
					})
					return
				}
				var newParea = []int64{0, 0, 0, 0, 0, 0, 0, 0}
				for i1, v4 := range v {
					newParea[i1] = p.PAreaCoins[i1] + v4
				}
				//判断是否到达区域限红
				for _, v := range newParea {
					if v > G_DbGetGameServerData.GameConfig.LimitRedMax {
						p.SendNativeMsg(MSG_GAME_INFO_BETAGAIN_REPLAY, BetAgainReply{
							Id:     MSG_GAME_INFO_BETAGAIN_REPLAY,
							Result: 3,
							Err:    "区域到达限红",
						})
						return
					}
				}
				for i1, v2 := range v {
					info := GBetCountInfo{
						Id:           i1,
						DownBetValue: this.Bets[i1].DownBetValue + v2,
					}
					if !p.Robot {
						info.UserBetValue = this.Bets[i1].DownBetValue + v2
					}
					this.Bets[i1] = info
				}
				for i1, v4 := range v {
					p.Coins -= v4
					p.TotalBet += v4
					p.DownBets[i1] += v4
					p.PAreaCoins[i1] += v4
				}
				replay := BetAgainReply{
					Id:          MSG_GAME_INFO_BETAGAIN_REPLAY,
					Result:      0,
					BetArea:     arr,
					Coins:       p.Coins,
					CanUserChip: this.CanUseChip(p),
				}
				p.SendNativeMsg(MSG_GAME_INFO_BETAGAIN_REPLAY, replay)
				chongfuxiazhu = true
				break
			}
		}
	}
	if pd1 {
		p.SendNativeMsg(MSG_GAME_INFO_BETAGAIN_REPLAY, BetAgainReply{
			Id:     MSG_GAME_INFO_BETAGAIN_REPLAY,
			Result: 1,
			Err:    "您上局未下注",
		})
	}
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
	// if this.GetFSM(0).GetMark() == GAME_STATUS_BALANCE {
	// 	logs.Debug("用户被踢出")
	// 	this.LeaveByForce(p)
	// 	return
	// }
	total := p.getDownBetTotal()
	DebugLog("total:", total)
	if total == 0 {
		logs.Debug("该玩家未下注，所以直接踢出")
		key := this.isInChair(p)
		if key != -1 {
			e := this.getElsePlayer(p)
			if len(e) == 0 {
				this.resetOneCharid(key)
			} else {
				xx := e[len(e)-1]
				this.Chair[key].Uid = xx.Uid
				this.Chair[key].Nick = xx.Nick
				this.Chair[key].Coins = xx.Coins
				this.Chair[key].Head = xx.Head
			}
		}
		//座位变更通知
		for _, vp := range this.Players {
			cu := ChairUpdate{
				Id: MSG_GAME_INFO_CHAIR_UPDATE,
			}
			ccc := make(map[int]GUserInfoByChair)
			for i, v := range this.Chair {
				ccc[i] = *v
			}

			for i, v := range ccc {
				if v.Uid != vp.Uid && v.Uid != 0 {
					// fmt.Println("开始换啦！！！！!!!!!!!!!!!!!!!!!")
					// fmt.Println("发送给:", vp.Uid, "名称为:", vp.Nick)
					// fmt.Println("将", v.Uid, "名称为:", v.Nick, "换掉")
					ccc[i] = GUserInfoByChair{
						Uid:   v.Uid,
						Nick:  "***" + v.Nick[len(v.Nick)-4:],
						Head:  v.Head,
						Coins: v.Coins,
					}
				}
			}
			// for _, v := range ccc {
			// 	fmt.Println(v)
			// }

			// for _, v := range this.Chair {
			// 	fmt.Println(v)
			// }
			cu.Chair = ccc
			vp.SendNativeMsg(MSG_GAME_INFO_CHAIR_UPDATE, cu)
		}
		delete(this.BetAgain, p.Uid)
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

	if result == 0 {
		key := this.isInChair(p)
		if key != -1 {
			if len(this.getElsePlayer(p)) == 0 {
				this.resetOneCharid(key)
			} else {
				e := this.getElsePlayer(p)
				xx := e[len(e)-1]
				this.Chair[key].Uid = xx.Uid
				this.Chair[key].Nick = xx.Nick
				this.Chair[key].Coins = xx.Coins
				this.Chair[key].Head = xx.Head
			}
		}
		//座位变更通知
		for _, vp := range this.Players {
			cu := ChairUpdate{
				Id: MSG_GAME_INFO_CHAIR_UPDATE,
			}
			ccc := make(map[int]GUserInfoByChair)
			for i, v := range this.Chair {
				ccc[i] = *v
			}

			for i, v := range ccc {
				if v.Uid != vp.Uid && v.Uid != 0 {
					// fmt.Println("开始换啦！！！！!!!!!!!!!!!!!!!!!")
					// fmt.Println("发送给:", vp.Uid, "名称为:", vp.Nick)
					// fmt.Println("将", v.Uid, "名称为:", v.Nick, "换掉")
					ccc[i] = GUserInfoByChair{
						Uid:   v.Uid,
						Nick:  "***" + v.Nick[len(v.Nick)-4:],
						Head:  v.Head,
						Coins: v.Coins,
					}
				}
			}
			// for _, v := range ccc {
			// 	fmt.Println(v)
			// }

			// for _, v := range this.Chair {
			// 	fmt.Println(v)
			// }
			cu.Chair = ccc
			vp.SendNativeMsg(MSG_GAME_INFO_CHAIR_UPDATE, cu)
		}
		delete(this.BetAgain, p.Uid)
		this.LeaveByForce(p)
	} else {
		gl := &GLeaveReply{
			Id:     MSG_GAME_LEAVE_REPLY,
			Result: result,
			Err:    err,
			Uid:    p.Uid,
		}
		p.SendNativeMsg(MSG_GAME_LEAVE_REPLY, gl)
	}
}

//判断玩家是否为 座位上的玩家，如果有，则返回座位id，如果没有，则返回-1
func (this *ExtDesk) isInChair(p *ExtPlayer) int {
	for i, v := range this.Chair {
		if v.Uid == p.Uid {
			return i
		}
	}
	return -1
}

func (this *ExtDesk) resetOneCharid(key int) {
	this.Chair[key] = &GUserInfoByChair{
		Uid:   0,
		Nick:  "",
		Head:  "",
		Coins: 0,
	}
	//现在将该座位重置了，需要传正确位置的数组
	count := 0
	cha := make(map[int]*GUserInfoByChair)
	for _, v := range this.Chair {
		if v.Uid != 0 {
			cha[count] = v
			count += 1
		}
	}
	this.Chair = cha
	oldlen := len(this.Chair) - 1
	c := 6 - oldlen

	if c > 0 {
		for i := 1; i <= c; i++ {
			this.Chair[oldlen+i] = &GUserInfoByChair{}
		}
	}
}

//玩家入座
func (this *ExtDesk) onChair(p *ExtPlayer, key int) {
	this.Chair[key] = &GUserInfoByChair{
		Uid:   p.Uid,
		Coins: p.Coins,
		Head:  p.Head,
		Nick:  p.Nick,
	}
}

//获取出自己外的座位玩家信息
func (this *ExtDesk) getChairPlayer(p *ExtPlayer) map[int]*ExtPlayer {
	var arr map[int]*ExtPlayer = make(map[int]*ExtPlayer)
	for _, v := range this.Players {
		key := this.isInChair(v)
		if key != -1 {
			// fmt.Println("执行次数！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！")
			if v.Uid != p.Uid {
				arr[key] = v
			}
		}
	}
	return arr
}

//获取除自己外的其他玩家信息
func (this *ExtDesk) getElsePlayer(p *ExtPlayer) []*ExtPlayer {
	var arr []*ExtPlayer
	for _, v := range this.Players {
		if this.isInChair(v) == -1 {
			if v.Uid != p.Uid {
				arr = append(arr, v)
			}
		}
	}
	return arr
}

//========================发送网络数据=================
// 发送桌子信息
func (this *ExtDesk) sendDeskInfo(p *ExtPlayer) {
	p.LiXian = false
	p.Online = true
	// fmt.Println("玩家金币：", p.Coins)
	//判断该玩家进入桌子之后，是否还有座位给他座
	if this.isInChair(p) == -1 {
		for i := 0; i < 6; i++ {
			v := this.Chair[i]
			if v != nil {
				if v.Uid == 0 {
					logs.Debug("将玩家入座", p.Nick)
					v.Uid = p.Uid
					v.Coins = p.Coins
					v.Head = p.Head
					v.Nick = p.Nick
					break
				}
			}
		}
	}
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
	var allbetcount int64 = 0
	for _, v := range betInfos {
		allbetcount += v.DownBetTotal
	}
	mult := []float32{0, 0, 0, 0, 0, 0, 0, 0}
	//fmt.Println(CarTypeMultiple)
	for i, v := range CarTypeMultiple {
		mult[i] = v
	}

	var allbet int64
	for _, v := range p.DownBets {
		allbet += v
	}
	info := &GClientDeskInfo{
		Id:          MSG_GAME_INFO_RDESKINFO,
		Result:      0,
		FangHao:     this.Mark,
		JuHao:       this.JuHao,
		Bets:        betInfos,
		AllBetCount: allbetcount,
		BetLevels:   G_DbGetGameServerData.GameConfig.TenChips,
		PlayerMassage: PlayerMsg{
			Uid:          p.Uid,
			MyUserAvatar: p.Head,
			MyUserName:   p.Account,
			MyUserCoin:   p.Coins,
		},
		Log:                this.LogoRecord,
		GameStatus:         this.GetFSM(0).GetMark(),
		GameStatusDuration: restTime / 1000,
		Multiple:           mult,
		AreaMaxCoin:        G_DbGetGameServerData.GameConfig.LimitRedMax,
		Car:                this.Car,
		ChangeCard:         this.ChangeCar,
		Chair:              this.Chair,
		CanUserChip:        this.CanUseChip(p),
		// OldOther:           p.OtherBet,
		Head:      p.Head,
		BetAll:    allbet,
		Nick:      p.Nick,
		CarName:   numBecomString([]int{this.GameResult})[0],
		Index:     this.GameResult,
		Double:    CarTypeMultiple[this.GameResult],
		ShengTime: this.GetTimerNum(this.GameState),
	}
	// 广播  座位变更通知
	for _, vp := range this.Players {
		cu := ChairUpdate{
			Id: MSG_GAME_INFO_CHAIR_UPDATE,
		}
		ccc := make(map[int]GUserInfoByChair)
		for i, v := range this.Chair {
			ccc[i] = *v
		}

		for i, v := range ccc {
			if v.Uid != vp.Uid && v.Uid != 0 {
				// fmt.Println("开始换啦！！！！!!!!!!!!!!!!!!!!!")
				// fmt.Println("发送给:", vp.Uid, "名称为:", vp.Nick)
				// fmt.Println("将", v.Uid, "名称为:", v.Nick, "换掉")
				ccc[i] = GUserInfoByChair{
					Uid:   v.Uid,
					Nick:  "***" + v.Nick[len(v.Nick)-4:],
					Head:  v.Head,
					Coins: v.Coins,
				}
			}
		}
		// for _, v := range ccc {
		// 	fmt.Println(v)
		// }

		// for _, v := range this.Chair {
		// 	fmt.Println(v)
		// }
		cu.Chair = ccc
		vp.SendNativeMsg(MSG_GAME_INFO_CHAIR_UPDATE, cu)
	}
	info.PlayerMassage.MyDownBets = make([]int64, gameConfig.LimitInfo.BetCount)
	for id, v := range p.DownBets {
		info.PlayerMassage.MyDownBets[int(id)] = v
	}

	info.AreaCoin = make([]int64, gameConfig.LimitInfo.BetCount)
	for id, v := range this.Bets {
		info.AreaCoin[id] = v.DownBetValue
	}

	p.SendNativeMsg(MSG_GAME_INFO_RDESKINFO, info)
}
func (this *ExtDesk) recDestInfo(p *ExtPlayer) {
	p.LiXian = false
	p.Online = true
	betInfos := []GClientBetCountInfo{}
	if this.isInChair(p) == -1 {
		for i := 0; i < 6; i++ {
			v := this.Chair[i]
			if v != nil {
				if v.Uid == 0 {
					v.Uid = p.Uid
					v.Coins = p.Coins
					v.Head = p.Head
					v.Nick = p.Nick
					break
				}
			}
		}

	}

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
	var allbetcount int64 = 0
	for _, v := range betInfos {
		allbetcount += v.DownBetTotal
	}
	mult := []float32{0, 0, 0, 0, 0, 0, 0, 0}
	//fmt.Println(CarTypeMultiple)
	for i, v := range CarTypeMultiple {
		mult[i] = v
	}
	var allbet int64
	for _, v := range p.DownBets {
		allbet += v
	}
	info := &GClientDeskInfo{
		Id:          MSG_GAME_INFO_RRECONNECT_REPLAY,
		Result:      0,
		FangHao:     this.Mark,
		JuHao:       this.JuHao,
		Bets:        betInfos,
		AllBetCount: allbetcount,
		BetLevels:   G_DbGetGameServerData.GameConfig.TenChips,
		PlayerMassage: PlayerMsg{
			Uid:          p.Uid,
			MyUserAvatar: p.Head,
			MyUserName:   p.Account,
			MyUserCoin:   p.Coins,
		},
		Log:                this.LogoRecord,
		GameStatus:         this.GetFSM(0).GetMark(),
		GameStatusDuration: restTime / 1000,
		Multiple:           mult,
		AreaMaxCoin:        G_DbGetGameServerData.GameConfig.LimitRedMax,
		Car:                this.Car,
		ChangeCard:         this.ChangeCar,
		Chair:              this.Chair,
		CanUserChip:        this.CanUseChip(p),
		// OldOther:           p.OtherBet,
		Head:      p.Head,
		BetAll:    allbet,
		Nick:      p.Nick,
		CarName:   numBecomString([]int{this.GameResult})[0],
		Index:     this.GameResult,
		Double:    CarTypeMultiple[this.GameResult],
		ShengTime: this.GetTimerNum(this.GameState),
	}
	for _, vp := range this.Players {
		cu := ChairUpdate{
			Id: MSG_GAME_INFO_CHAIR_UPDATE,
		}
		ccc := make(map[int]GUserInfoByChair)
		for i, v := range this.Chair {
			ccc[i] = *v
		}

		for i, v := range ccc {
			if v.Uid != vp.Uid && v.Uid != 0 {
				// fmt.Println("开始换啦！！！！!!!!!!!!!!!!!!!!!")
				// fmt.Println("发送给:", vp.Uid, "名称为:", vp.Nick)
				// fmt.Println("将", v.Uid, "名称为:", v.Nick, "换掉")
				ccc[i] = GUserInfoByChair{
					Uid:   v.Uid,
					Nick:  "***" + v.Nick[len(v.Nick)-4:],
					Head:  v.Head,
					Coins: v.Coins,
				}
			}
		}
		cu.Chair = ccc
		vp.SendNativeMsg(MSG_GAME_INFO_CHAIR_UPDATE, cu)
	}

	info.PlayerMassage.MyDownBets = make([]int64, gameConfig.LimitInfo.BetCount)
	for id, v := range p.DownBets {
		info.PlayerMassage.MyDownBets[int(id)] = v
	}

	info.AreaCoin = make([]int64, gameConfig.LimitInfo.BetCount)
	for id, v := range this.Bets {
		info.AreaCoin[id] = v.DownBetValue
	}

	p.SendNativeMsg(MSG_GAME_INFO_RRECONNECT_REPLAY, info)
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
	if this.GameState != GAME_STATUS_DOWNBET {
		msg := GSDownBet{
			Id:     MSG_GAME_INFO_RDOWNBET,
			Result: 1,
			Err:    "不是下注状态，下注失败！",
		}
		p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, msg)
		return
	}
	if p == nil {
		return
	}
	//
	var result int32 = 0
	var err string = ""

	// 判断金币下注索引
	if coinIdx >= 0 && coinIdx >= len(G_DbGetGameServerData.GameConfig.TenChips) {
		p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, &GSDownBet{
			Id:     MSG_GAME_INFO_RDOWNBET,
			Result: 1,
		})
		return
	}
	betInfo, ok := this.Bets[betIdx]
	if ok == false {
		p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, &GSDownBet{
			Id:     MSG_GAME_INFO_RDOWNBET,
			Result: 1,
			Err:    "金额错误,下注失败",
		})
		return
	}

	var downCoin int64 = 0
	downCoin = G_DbGetGameServerData.GameConfig.TenChips[coinIdx]
	// 下注的金币是否足够
	if downCoin > p.Coins {
		msg := GSDownBet{
			Id:     MSG_GAME_INFO_RDOWNBET,
			Result: 2,
			Err:    "金币不足，请充值",
		}
		p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, msg)
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

	isOutArea = areaValue > G_DbGetGameServerData.GameConfig.LimitRedMax

	if isOutArea {

		msg := GSDownBet{
			Id:     MSG_GAME_INFO_RDOWNBET,
			Result: 3,
			Err:    "下注区域超过最大值",
		}
		p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, msg)
		return
	}
	if ok {
		p.DownBets[betIdx] += downCoin
	} else {
		p.DownBets[betIdx] = downCoin
	}
	p.TotalBet += downCoin
	p.Coins -= downCoin
	betInfo.DownBetValue += downCoin //区域下注金额增加
	p.PAreaCoins[betIdx] += downCoin
	if !p.Robot {
		betInfo.UserBetValue += downCoin
	}
	this.Bets[betIdx] = betInfo //更新区域属性

	msg := GSDownBet{
		Id:          MSG_GAME_INFO_RDOWNBET,
		Result:      result,
		Err:         err,
		PAreaCoins:  p.PAreaCoins,
		Coins:       p.Coins,
		AreaId:      betIdx,
		CoinId:      coinIdx,
		CanUserChip: this.CanUseChip(p),
		DownCoins:   downCoin,
	}

	// 玩家请求下注回复
	p.SendNativeMsg(MSG_GAME_INFO_RDOWNBET, &msg)
	xaizhu = true
}

var chongfuxiazhu bool = false
var broCount int = 1 //记录广播次数
var oldbet [8]int64  //记录旧的 各区域下注
var oldOther map[int64][]int64 = make(map[int64][]int64)
var xaizhu bool = false

func (this *ExtDesk) BroBetMsg(d interface{}) {
	// fmt.Println("广播:", broCount)
	var bet = [8]int64{0, 0, 0, 0, 0, 0, 0, 0}
	for key, v := range this.Bets {
		bet[key] += v.DownBetValue
	}
	if broCount == 1 {
		oldbet = bet
	}
	if oldbet == bet && !chongfuxiazhu && !xaizhu {
		// fmt.Println("没有玩家下注，不广播")
		broCount++
		if broCount <= gameConfig.StateInfo.DownBetTime/1000 {
			this.AddTimer(gameConfig.StateInfo.BroMsg, gameConfig.StateInfo.BroMsgTime, this.BroBetMsg, nil)
		}
		return
	}
	var allbet int64
	for _, v := range bet {
		allbet += v
	}

	// fmt.Println("有玩家下注，广播")

	//广播下注信息
	for _, pl := range this.Players {
		msg := GNDownBet{
			Id:      MSG_GAME_INFO_NDOWNBET,
			Bets:    bet,
			AllBets: allbet,
		}
		// var seat [][]int64
		// for i := 0; i < 6; i++ {
		// 	seat = append(seat, []int64{0, 0, 0, 0, 0, 0, 0, 0})
		// }
		// msg.SeatBetList = seat
		// seatPlayer := this.getChairPlayer(pl)
		// fmt.Println("seatPlayer:", seatPlayer)
		// otherPlayer := this.getElsePlayer(pl)
		//获取座位玩家下注信息
		// for _, v1 := range seatPlayer {
		// 	fmt.Println("v1.PAreaCoins:", v1.PAreaCoins)
		// 	for i, v2 := range v1.PAreaCoins {
		// 		// msg.SeatBetList[i1][i] = v2
		// 	}
		// }
		// fmt.Println("msg.SeatBetList:", msg.SeatBetList)
		//获取其他玩家（未在座位上的玩家）下注信息
		var other []int64 = []int64{0, 0, 0, 0, 0, 0, 0, 0}
		var ok = true
		for _, v1 := range this.Players {
			// fmt.Println("其他玩家有值？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？")
			if pl.Uid != v1.Uid {
				for i, v2 := range v1.PAreaCoins {
					other[i] += v2
				}
			}
		}
		msg.OldOtherBetList = other
		// fmt.Println("第一个oldther", msg.OldOtherBetList)
		_, ok = oldOther[pl.Uid]
		if !ok {
			// logs.Debug("走了第一个")
			oldOther[pl.Uid] = other
			// fmt.Println("other:", other)
			// fmt.Println("oldOther:1111111111111111111111111111111", oldOther[pl.Uid])
		}
		if ok {
			// fmt.Println("oldOther:222222222222222222222222222222222", oldOther[pl.Uid])
			for i, v := range msg.OldOtherBetList {
				//  fmt.Println("下注广播", i+1, "次！！！！！！！！！！！！！！！！！！！！！！！！！！")
				msg.OtherBetList = append(msg.OtherBetList, v-oldOther[pl.Uid][i])
			}
		} else {
			msg.OtherBetList = msg.OldOtherBetList
		}
		for i, v := range msg.OtherBetList {
			if v < 0 {
				msg.OtherBetList[i] = 0
			}
		}
		// fmt.Println("msg.OtherBetList:", msg.OtherBetList)
		msg.PAreaCoins = pl.PAreaCoins
		// fmt.Println("other:", msg.OldOtherBetList)
		// fmt.Println("newother:", msg.OtherBetList)
		pl.SendNativeMsg(MSG_GAME_INFO_NDOWNBET, msg)
		// pl.OtherBet = msg.OldOtherBetList
		// fmt.Println("其他玩家下注", pl.OtherBet)
		oldOther[pl.Uid] = msg.OldOtherBetList
	}
	// // 通知所有玩家有新的下注
	// this.BroadcastAll(MSG_GAME_INFO_NDOWNBET, &GNDownBet{
	// 	Id:      MSG_GAME_INFO_NDOWNBET,
	// 	Bets:    bet,
	// 	AllBets: allbet,
	// })

	if broCount <= gameConfig.StateInfo.DownBetTime/1000 {
		this.AddTimer(gameConfig.StateInfo.BroMsg, gameConfig.StateInfo.BroMsgTime, this.BroBetMsg, nil)
	}
	broCount++
	chongfuxiazhu = false
	oldbet = bet
	xaizhu = false
}
func (this *ExtDesk) getPlayer() []*ExtPlayer {
	players := this.Players
	for i := 0; i < len(players); i++ {
		for j := 1; j < len(players)-i; j++ {
			if players[j].Coins > players[j-1].Coins {
				//交换
				players[j], players[j-1] = players[j-1], players[j]
			}
		}
	}
	if len(players) > gameConfig.LimitInfo.PlayerList {
		players = append(players, players[:20]...)
	}

	return players
}
func (this *ExtDesk) handleMorePlayer(p *ExtPlayer, d *DkInMsg) {
	fmt.Println("监听到玩家请求更多玩家 ")
	userInfo := []GUserInfo{}
	newuserInfo := []GUserInfo{}
	var players []*ExtPlayer

	players = this.getPlayer()
	for _, v := range players {
		fmt.Print(v.Nick)
		fmt.Print(v.Uid)
	}

	for i, v := range players {
		if i > gameConfig.LimitInfo.PlayerList-1 {
			break
		}
		info := GUserInfo{
			Uid:   v.Uid,
			Head:  v.Head,
			Coins: v.Coins,
		}
		if p.Uid != v.Uid { //不是自己的玩家 名字隐藏
			info.Nick = "***" + v.Nick[4:]
		} else {
			info.Nick = p.Nick
		}
		info.TotBet = v.TotalBet
		if len(v.BetHistorys) <= 20 {
			for _, v := range v.BetHistorys {
				if v == 0 {
					info.WinCount += 1
				}
			}
			// v.Match = len(v.BetHistorys)
			v.Match = 20
			fmt.Println("MATCH:", v.Match)
		} else {
			his := v.BetHistorys[len(v.BetHistorys)-20:] //截取20个
			for _, v := range his {
				if v == 0 {
					info.WinCount += 1
				}
			}
			v.Match = 20
		}
		info.Index = i + 1
		info.Match = v.Match
		if info.TotBet == 0 {
			newuserInfo = append(newuserInfo, info)
		} else {
			userInfo = append(userInfo, info)

		}
	}
	userInfo = append(userInfo, newuserInfo...)
	msg := GUserInfoReply{
		Id:       MSG_GAME_INFO_MOREPLAYER_REPLY,
		UserInfo: userInfo,
	}
	p.SendNativeMsg(MSG_GAME_INFO_MOREPLAYER_REPLY, msg)
	// fmt.Println(msg)
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
			this.LeaveByForce(p)
			delete(this.BetAgain, p.Uid)
			plen -= 1
			i -= 1
			key := this.isInChair(p)
			if key != -1 {
				if len(this.getElsePlayer(p)) == 0 {
					this.resetOneCharid(key)
				} else {
					e := this.getElsePlayer(p)
					xx := e[len(e)-1]
					this.Chair[key].Uid = xx.Uid
					this.Chair[key].Nick = xx.Nick
					this.Chair[key].Coins = xx.Coins
					this.Chair[key].Head = xx.Head
				}
			}
			//座位变更通知
			for _, vp := range this.Players {
				cu := ChairUpdate{
					Id: MSG_GAME_INFO_CHAIR_UPDATE,
				}
				ccc := make(map[int]GUserInfoByChair)
				for i, v := range this.Chair {
					ccc[i] = *v
				}

				for i, v := range ccc {
					if v.Uid != vp.Uid && v.Uid != 0 {
						// fmt.Println("开始换啦！！！！!!!!!!!!!!!!!!!!!")
						// fmt.Println("发送给:", vp.Uid, "名称为:", vp.Nick)
						// fmt.Println("将", v.Uid, "名称为:", v.Nick, "换掉")
						ccc[i] = GUserInfoByChair{
							Uid:   v.Uid,
							Nick:  "***" + v.Nick[len(v.Nick)-4:],
							Head:  v.Head,
							Coins: v.Coins,
						}
					}
				}
				// for _, v := range ccc {
				// 	fmt.Println(v)
				// }

				// for _, v := range this.Chair {
				// 	fmt.Println(v)
				// }
				cu.Chair = ccc
				vp.SendNativeMsg(MSG_GAME_INFO_CHAIR_UPDATE, cu)
			}
		}
		p.ResetExtPlayer()
	}
}

//小郑
// 统计桌面金币，并分配牌，实现输赢控制器
func (this *ExtDesk) allotCard(wantwin bool) (result int) {
	logs.Debug("进入风控，控制：", wantwin)
	var allcoin int64
	allAreacoin := make([]int64, 8, 8)
	for i, b := range this.Bets {
		allcoin += b.UserBetValue
		allAreacoin[i] += b.UserBetValue
	}
	logs.Debug("玩家总下注：", allcoin, "玩家区域下注：", allAreacoin)
	//赢的区域数组
	var winArea []int
	//输的区域数组
	var failArea []int
	//查找赢的区域
	for i, _ := range allAreacoin {
		if float32(allAreacoin[i])*CarTypeMultiple[i] < float32(allcoin) {
			winArea = append(winArea, i)
		} else {
			failArea = append(failArea, i)
		}
	}
	logs.Debug("赢的区域：", winArea, "输的区域：", failArea)
	if wantwin {
		if len(winArea) == 0 {
			fmt.Println("没有赢的区域")
			result = int(RandInt64(8))
			return
		}
		result = winArea[RandInt64(int64(len(winArea)))]
		return
	}
	if len(failArea) == 0 {
		fmt.Println("没有失败的区域")
		result = int(RandInt64(8))
		return
	}
	result = failArea[RandInt64(int64(len(failArea)))]
	return
}

//判断用户还能使用的金币筹码index

//根据玩家金币判断还有哪些筹码能够下注
func (this *ExtDesk) CanUseChip(p *ExtPlayer) int {
	indexs := -1
	var allbet int64
	for _, v := range p.DownBets {
		allbet += v
	}
	for i, v := range G_DbGetGameServerData.GameConfig.TenChips {
		if v <= p.Coins {
			indexs = i
		}
	}
	return indexs
}

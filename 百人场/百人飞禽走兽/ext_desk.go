package main

import (
	// "fmt"
	"crypto/rand"
	"logs"
	"math/big"
)

type ExtDesk struct {
	Desk
	DownBet             []int64             //区域总下注金币集合
	ChairList           []PlayerInfoByChair //座位玩家信息
	NeedBro             bool                //是否需要广播下注（即是否 有新的下注，需要广播）
	NeedUpdata          bool                //是否需要更新其他玩家下注
	GameResult          int                 //开奖结果
	BalanceResult       []int64             //结算集合
	GameResultHistory   []int               //开奖记录
	Zhuang              PlayerMsg           //是否为玩家上庄
	WaitZhuang          []PlayerMsg         //等待上装玩家集合
	ZhuangBalanceResult []int64             //庄家结算
	ZhenShiDownBet      []int64             //正式玩家押注区域
	ZhenShisGetcoins    int64               //玩家真实输赢
	ZhuangMatchCount    int                 //庄家局数
}

//桌子初始化
func (this *ExtDesk) InitExtData() {
	for i := 0; i < gameConfig.LimitInfo.BetCount; i++ {
		this.DownBet = append(this.DownBet, 0)
	}
	for i := 0; i < gameConfig.LimitInfo.ChairNum; i++ {
		this.ChairList = append(this.ChairList, PlayerInfoByChair{})
	}
	//初始化倍数
	InitMultiple()
	//直接进入准备状态
	this.status_ready(nil)
}

func (this *ExtDesk) AddListen() {
	this.Handle[MSG_GAME_INFO_QDESKINFO] = this.HandleQDeskInfo
	this.Handle[MSG_GAME_INFO_DOWNBET] = this.HandleDownBet
	this.Handle[MSG_GAME_INFO_GETMOREPLAYER] = this.HandleGetMorePlayer
	this.Handle[MSG_GAME_INFO_UPZHUANG] = this.HandleUpZhuang
	this.Handle[MSG_GAME_INFO_DOWNZHUANG] = this.HandleDownZhuang
	this.Handle[MSG_GAME_INFO_GET_RECORD] = this.HandleGetRecord
}

//重置桌子方法
func (this *ExtDesk) Rest() {
	for i := 0; i < gameConfig.LimitInfo.BetCount; i++ {
		this.DownBet = append(this.DownBet, 0)
	}
	for i := 0; i < gameConfig.LimitInfo.ChairNum; i++ {
		this.ChairList = append(this.ChairList, PlayerInfoByChair{})
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
	for i, v := range G_DbGetGameServerData.GameConfig.TenChips {
		if v < p.Coins {
			indexs = i
		}
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
			p.IsOnChair = false
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
		logs.Debug("发现玩家下注，广播")
		for _, v := range this.Players {
			if v.Uid != this.Zhuang.Uid {
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
			} else {
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
					IsZhuang:     true,
				})
			}
		}
		this.NeedBro = false
		this.NeedUpdata = true
	} else {
		logs.Debug("发现玩家未下注 广播")
	}
	if this.GameState == GAME_STATUS_STARTBET || this.GameState == GAME_STATUS_DOWNBET || this.GameState == GAME_STATUS_ENDBET {
		this.AddTimer(999, 1, this.BroDownBetInfo, nil) //开始每隔一秒更新一次下注情况
	}
}

//输赢控制器
func (this *ExtDesk) ControlWinOrLose(win bool) (result int) {
	/*
		思路： 获取所有玩家下注总额，再获取所有玩家再每个区域的下注金额，以此来获取如果开每个区域，每个区域将会亏损多少钱，将亏损的钱与玩家总下注的钱进行对比，
		如果亏损钱小于赚的钱，那么证明该区域会使庄家有利，将该区域存在可开奖区域中
	*/
	var allBet int64                 //玩家下注总额
	areaBet := make(map[int]int64)   //玩家每个区域下注总额
	loseCoins := make(map[int]int64) //开奖区域输钱集合
	var openresult []int             //开奖区域集合
	var openresulttolose []int       //玩家赢
	for i, v := range this.DownBet {
		allBet += v
		areaBet[i] += v
	}
	for i, v := range areaBet {
		loseCoins[i] = (v * int64(LotteryDouble[i]*10)) / 10
	}
	for i, v := range loseCoins {
		if v < allBet {
			openresult = append(openresult, i)
		}
		if v >= allBet {
			openresulttolose = append(openresulttolose, i)
		}
	}
	if win {
		result = openresult[0]
		logs.Debug("庄赢")
	} else {
		logs.Debug("庄尽量输")
		if len(openresulttolose) > 0 {
			result = openresulttolose[0]
		} else {
			result = this.randresult()
		}
	}
	return
}

//纯随机开奖
func (this *ExtDesk) randresult() (result int) {
	logs.Debug("纯随机开奖")
	result, _ = GetRandomNum(0, gameConfig.LimitInfo.BetCount)
	return
}

//玩家上庄
func (this *ExtDesk) UpZhuang(p *ExtPlayer) {
	playermsg := PlayerMsg{
		Uid:          p.Uid,
		MyUserAvatar: p.Head,
		MyUserCoin:   p.Coins,
		MyUserName:   p.Nick,
	}
	if this.Zhuang.Uid == 0 {
		this.Zhuang = playermsg
		p.SendNativeMsg(MSG_GAME_INFO_UPZHUANG_REPLY, ChangZhuangReply{
			Id:     MSG_GAME_INFO_UPZHUANG_REPLY,
			Result: 0,
		})
		this.ZhuangMatchCount = 0
		this.ChangeZhuang()
	} else {
		if len(this.WaitZhuang) < gameConfig.LimitInfo.WaitZhuangCount {
			p.SendNativeMsg(MSG_GAME_INFO_UPZHUANG_REPLY, ChangZhuangReply{
				Id:     MSG_GAME_INFO_UPZHUANG_REPLY,
				Result: 1,
				Err:    "庄家已经存在，已经将您添加到等待队列，轮到您的时候自动为您上庄",
			})
			this.WaitZhuang = append(this.WaitZhuang, playermsg)
		} else {
			p.SendNativeMsg(MSG_GAME_INFO_UPZHUANG_REPLY, ChangZhuangReply{
				Id:     MSG_GAME_INFO_UPZHUANG_REPLY,
				Result: 2,
				Err:    "玩家已经满了,请晚些时候请求坐庄把！",
			})
		}
	}
}

//玩家下庄
func (this *ExtDesk) DownZhuang(p *ExtPlayer) {
	if this.Zhuang.Uid != p.Uid {
		logs.Debug("该玩家不是庄家，无法执行下庄操作")
		p.SendNativeMsg(MSG_GAME_INFO_DOWNZHUANG_REPLY, ChangZhuangReply{
			Id:     MSG_GAME_INFO_DOWNZHUANG_REPLY,
			Result: 1,
			Err:    "您不是庄家，无法下庄 ",
		})
		return
	}
	p.SendNativeMsg(MSG_GAME_INFO_DOWNZHUANG_REPLY, ChangZhuangReply{
		Id:     MSG_GAME_INFO_DOWNZHUANG_REPLY,
		Result: 0,
	})
	this.Zhuang.Uid = 0
	if len(this.WaitZhuang) > 0 {
		this.Zhuang = this.WaitZhuang[0]
		this.WaitZhuang = this.WaitZhuang[1:]
		this.ZhuangMatchCount = 0
	}
	this.ChangeZhuang()
}

//庄家改变通知
func (this *ExtDesk) ChangeZhuang() {
	zinfo := ZhuangInfo{
		Id:       MSG_GAME_INFO_CHANGEZHUANG,
		Info:     this.Zhuang,
		WaitList: this.WaitZhuang,
	}
	for _, v := range this.Players {
		if v.Uid != this.Zhuang.Uid {
			zinfo.Info.MyUserName = this.ChangeNick(zinfo.Info.MyUserName)
		}
		for i, v1 := range zinfo.WaitList {
			if v.Uid != v1.Uid {
				zinfo.WaitList[i].MyUserName = this.ChangeNick(zinfo.Info.MyUserName)
			}
		}
		v.SendNativeMsg(MSG_GAME_INFO_CHANGEZHUANG, zinfo)
	}
}

//庄家改变通知（只通知一个人)
func (this *ExtDesk) ChangeZhuangToOne(p *ExtPlayer) {
	zinfo := ZhuangInfo{
		Id:       MSG_GAME_INFO_CHANGEZHUANG,
		Info:     this.Zhuang,
		WaitList: this.WaitZhuang,
	}
	for _, v := range this.Players {
		if p.Uid == v.Uid {
			if v.Uid != this.Zhuang.Uid {
				zinfo.Info.MyUserName = this.ChangeNick(zinfo.Info.MyUserName)
			}
			for i, v1 := range zinfo.WaitList {
				if v.Uid != v1.Uid {
					zinfo.WaitList[i].MyUserName = this.ChangeNick(zinfo.Info.MyUserName)
				}
			}
			v.SendNativeMsg(MSG_GAME_INFO_CHANGEZHUANG, zinfo)
		}
	}
}

//将玩家移除等待队列
func (this *ExtDesk) RemoveWait(p *ExtPlayer) {
	for i, v := range this.WaitZhuang {
		if v.Uid == p.Uid {
			this.WaitZhuang = append(this.WaitZhuang[i:0], this.WaitZhuang[i+1:]...)
		}
	}
}

//获取庄家信息
func (this *ExtDesk) getZhuangInfo(p *ExtPlayer) {
	plmsg := this.Zhuang
	if p.Uid != plmsg.Uid {
		plmsg.MyUserName = this.ChangeNick(plmsg.MyUserName)
	}
	return plmsg
}

//获取玩家信息
func (this *ExtDesk) getWaitZhuangListInfo(p *ExtPlayer) {
	wai := this.WaitZhuang
	for i, v := range wai {
		if v.Uid != p.Uid {
			wai[i].MyUserName = this.ChangeNick(wai[i].MyUserName)
		}
	}
	return wai
}

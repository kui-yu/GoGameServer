package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/tidwall/gjson"
)

type ExtRobotClient struct {
	RobotClient

	DeskInfo GClientDeskInfo
}

//=======================对外方法=====================
// 机器人开始
func (this *ExtRobotClient) Start() {
	this.ChildPtr = this
	this.IsStop = false
	this.GameIn = false
	this.IsConnect = false
	this.Coin = int64(gameConfig.getGameConfigCoin("initcoin"))

	this.SendMsg = make(chan *SendMsg, 1000)
	this.RecvMsg = make(chan *RecvMsg, 1000)
	this.Handle = make(map[uint32]func(string))
	this.EventMsg = make(chan *EventMsg, 1000)
	this.EventHandle = make(map[int32]func(interface{}))
	this.TimeTicker = new(TimeTicker)
	this.SysTicker = nil

	DebugLog("机器人开始了")
	this.Handle[MSG_GAME_AUTO_REPLY] = this.Login
	//投注阶段
	this.Handle[MSG_GAME_INFO_STAGE_BET] = this.RobotBet
	//游戏结果
	this.Handle[MSG_GAME_INFO_STAGE_GAME_RESULT] = this.GameResult
	//结算信息
	this.Handle[MSG_GAME_INFO_STAGE_SETTLE] = this.GameSettle
	this.addEventHandler()
	this.CheckOnLineId = this.TimeTicker.AddTimer(30, func(int, interface{}) {
		this.RestChekIsNo()
	}, nil)
	go func() {
		this.AddHandler()
		this.requstAccountInfo(this.Coin)
		this.OnMessageHandle()
	}()
}
func (this *ExtRobotClient) RestChekIsNo() {
	time.Sleep(time.Second * 30)
	controller.sendEvent(EVENT_CONT_ROBOTSHIFT, this)
}
func (this *ExtRobotClient) addEventHandler() {
	this.EventHandle[EVENT_CONNECT_SUCCESS] = this.onConnected
}

func (this *ExtRobotClient) AddHandler() {
	this.Handle[MSG_HALL_ROBOT_LOGIN_REPLY] = this.onLogin
	this.Handle[MSG_HALL_JOIN_GAME_REPLY] = this.onJoinGame

}

// 添加网络消息监听
func (this *ExtRobotClient) AddGameHandler() {

	this.Handle[MSG_GAME_NDOWNBET] = this.onDownBetN
	this.Handle[MSG_GAME_NSEATDOWN] = this.onNSeatDown // 座位信息改变通知
	this.Handle[MSG_GAME_RMANYUSER] = this.onManyUser
}

func (this *ExtRobotClient) RemoveAllHandler() {
	delete(this.Handle, MSG_HALL_ROBOT_LOGIN_REPLY)
	delete(this.Handle, MSG_HALL_JOIN_GAME_REPLY)
	delete(this.Handle, MSG_GAME_AUTO_REPLY)

	delete(this.Handle, MSG_GAME_NSTATUS_CHANGE)
	delete(this.Handle, MSG_GAME_NDOWNBET)
	delete(this.Handle, MSG_GAME_NSEATDOWN)
	delete(this.Handle, MSG_GAME_RMANYUSER)
}

// 下注通知
func (this *ExtRobotClient) onDownBetN(str string) {
	DebugLog("接收到下注", str)

	uid := gjson.Get(str, "Uid").Int()
	if this.UserInfo.Uid == uid {
		coin := gjson.Get(str, "Coin").Int()
		DebugLog("用户自己下注扣除金币", coin)
	}
}

func (this *ExtRobotClient) onManyUser(json string) {
	TestLog("请求更多玩家返回", json)
}

func (this *ExtRobotClient) onAllocSeat() {
	// 判断是否需要机器人开始抢座
	seatLen := 0
	for _, seat := range this.DeskInfo.Seats {
		if seat.UserId != 0 {
			seatLen++
		}
	}

	emptySeatLen := 4 - seatLen
	minemptyseat := gameConfig.getGameConfigInt("minemptyseat")

	if emptySeatLen <= minemptyseat {
		return
	}

	num := emptySeatLen - minemptyseat

	clients := controller.getRobotClients()
	// 删除坐下的玩家
	var useclients []*ExtRobotClient
	for _, v := range clients {
		var exists = false
		for _, s := range v.DeskInfo.Seats {
			if v.UserInfo.Uid == s.UserId {
				exists = true
				break
			}
		}

		if exists == false {
			useclients = append(useclients, v)
		}
	}

	clen := len(useclients)

	if clen <= num {
		for i := 0; i < clen; i++ {
			useclients[i].AddEventNative(EVENT_ROBOT_SEATDOWN, this.DeskInfo.GameStatusDuration)
		}
	} else {
		var oldIdxs map[int]bool = make(map[int]bool)

		for {

			idx, _ := GetRandomNum(0, clen)
			if _, ok := oldIdxs[idx]; ok == false {
				oldIdxs[idx] = true
				useclients[idx].AddEventNative(EVENT_ROBOT_SEATDOWN, this.DeskInfo.GameStatusDuration)
				num--
				if num == 0 {
					break
				}
			}
		}
	}
}

// 坐下通知
func (this *ExtRobotClient) onNSeatDown(str string) {
	DebugLog("座位信息改变通知", str)

	idx := gjson.Get(str, "SeatId").Int()
	stype := gjson.Get(str, "Type").Int() //0添加 1修改 2删除

	if stype == 0 {
		this.DeskInfo.Seats = append(this.DeskInfo.Seats, GClientSeatInfo{
			Id:            int(idx),
			UserId:        gjson.Get(str, "NewUserId").Int(),
			Name:          gjson.Get(str, "NewUserName").String(),
			Avatar:        gjson.Get(str, "NewUserAvatar").String(),
			SeatDownCount: 1,
			DownBetTotal:  0,
		})
	} else if stype == 1 {
		for i, seat := range this.DeskInfo.Seats {
			if seat.Id == int(idx) {
				seat.UserId = gjson.Get(str, "NewUserId").Int()
				seat.Name = gjson.Get(str, "NewUserName").String()
				seat.Avatar = gjson.Get(str, "NewUserAvatar").String()
				this.DeskInfo.Seats[i] = seat
				break
			}
		}
	} else if stype == 2 {
		for i, seat := range this.DeskInfo.Seats {
			if seat.Id == int(idx) {
				this.DeskInfo.Seats = append(this.DeskInfo.Seats[:i], this.DeskInfo.Seats[i+1:]...)
				break
			}
		}
	}
}

// 判断是否在座位上
func (this *ExtRobotClient) IsSeatDown() bool {
	seats := this.DeskInfo.Seats

	isSeatDown := false
	for _, v := range seats {
		if v.UserId == this.UserInfo.Uid {
			isSeatDown = true
			break
		}
	}

	return isSeatDown
}

// 判断是否在座位上
func (this *ExtRobotClient) GetMySeatDown() *GClientSeatInfo {
	seats := this.DeskInfo.Seats

	for _, v := range seats {
		if v.UserId == this.UserInfo.Uid {
			return &v
		}
	}

	return nil
}

//骰宝添加的
func (this *ExtRobotClient) Login(str string) {
	this.AddMsgNative(MSG_GAME_INFO_PLAYER_IN, struct{ Id int }{Id: MSG_GAME_INFO_PLAYER_IN})
}

//机器人下注
func (this *ExtRobotClient) RobotBet(str string) {
	obj := new(GABetInfo)
	obj.Id = MSG_GAME_INFO_PLAYER_BET
	arr := []int64{1, 5, 10, 50, 100}
	rand.Seed(time.Now().UnixNano())
	r := make([]int64, 4)
	for i := 1; i <= rand.Intn(ROBOT_MAX_BET_COUNT+1); i++ {
		r[rand.Intn(4)] = arr[rand.Intn(len(arr))]
		obj.Big = r[0]
		obj.Small = r[1]
		obj.Odd = r[2]
		obj.Even = r[3]
		this.AddMsgNative(MSG_GAME_INFO_PLAYER_BET, obj)
		r = make([]int64, 4)
	}
}

//游戏结果
func (this *ExtRobotClient) GameResult(str string) {
	obj := new(GSGameHistory)
	json.Unmarshal([]byte(str), obj)
	fmt.Println("================")
	fmt.Println("骰子1: ", obj.Info.NumberOne)
	fmt.Println("骰子2: ", obj.Info.NumberTwo)
	fmt.Println("骰子3: ", obj.Info.NumberThree)
	fmt.Println("开大: ", obj.Info.Big)
	fmt.Println("开小: ", obj.Info.Small)
	fmt.Println("开单: ", obj.Info.Odd)
	fmt.Println("开双: ", obj.Info.Even)
	fmt.Println("================")
}

//游戏结算结果
func (this *ExtRobotClient) GameSettle(str string) {
	obj := new(GSSettleInfo)
	json.Unmarshal([]byte(str), obj)
	fmt.Println("======================")
	fmt.Println("Id: ", obj.Id)
	fmt.Println("Uid: ", obj.Uid)
	fmt.Println("输赢的钱: ", obj.Count)
	fmt.Println("扣掉的费率: ", obj.CountRate)
	fmt.Println("玩家总金额: ", obj.Coin)
	fmt.Println("赢时押的位置: ", obj.Place)
	fmt.Println("当前可下注范围: ", obj.MaxBet)
	fmt.Println("======================")
}
func (this *ExtRobotClient) Stop() bool {
	if this.IsRobotStop() {
		return false
	}

	this.RWMutex.Lock()
	this.IsStop = true
	defer this.RWMutex.Unlock()

	this.BaseStop()

	this.RemoveAllHandler()

	return true
}

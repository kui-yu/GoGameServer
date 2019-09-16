package seatlist

import (
	"bl.com/util"
)

type IPlayer interface {
	GetUid() int64          // 获取用户ID
	GetCoins() int64        // 获取金币
	GetNewBetList() []int64 // 获取新下注列表
	GetBetCoins() int64     // 获取下注金币
	GetWinCoins() int64     // 获取赢取金币
}

type MgrSeat struct {
	Player []interface{}

	Seat []interface{}

	SeatNum int
}

func (this *MgrSeat) SetSeatNum(num int) bool {
	if this.SeatNum == 0 {
		this.SeatNum = num
		return true
	}

	return false
}

// 添加玩家
func (this *MgrSeat) AddPlayer(p interface{}) {
	var isAppended = false
	var p1 = p.(IPlayer)
	for _, v := range this.Player {
		p2 := v.(IPlayer)
		if p1.GetUid() == p2.GetUid() {
			isAppended = true
			break
		}
	}

	if isAppended {
		return
	}

	this.Player = append(this.Player, p)

	isAppended = false
	for _, v := range this.Seat {
		p2 := v.(IPlayer)
		if p1.GetUid() == p2.GetUid() {
			isAppended = true
			break
		}
	}

	if isAppended {
		return
	}

	// 座位玩家不足 补充座位玩家
	if len(this.Player) < this.SeatNum {
		this.Seat = append(this.Seat, p)
	}
}

func (this *MgrSeat) DelPlayer(p interface{}) bool {
	var p1 = p.(IPlayer)
	for i, v := range this.Player {
		p2 := v.(IPlayer)
		if p1.GetUid() == p2.GetUid() {
			this.Player = append(this.Player[:i], this.Player[i+1:]...)
			return true
		}
	}

	return false
}

// 更新座位列表
func (this *MgrSeat) UpdateSeatList() {
	var player []interface{}
	for _, v := range this.Player {
		player = append(player, v)
	}

	length := len(player)
	num := length
	if num > this.SeatNum {
		num = this.SeatNum
	}

	this.Seat = []interface{}{}
	for i := 0; i < num; i++ {
		rand, _ := util.GetRandomNum(0, length-i)
		this.Seat = append(this.Seat, player[rand])
		player = append(player[:rand], player[rand+1:]...)
	}
}

// 获取座位列表
func (this *MgrSeat) GetSeatList() []interface{} {
	length := len(this.Seat)
	lengthP := len(this.Player)

	if length < this.SeatNum && length < lengthP {
		count := lengthP
		if count > this.SeatNum {
			count = this.SeatNum
		}

		for i := length; i < count; i++ {
			this.Seat = append(this.Seat, this.Player[i])
		}
	}

	return this.Seat
}

// 获取座位玩家下注信息
func (this *MgrSeat) GetSeatNewBetList() [][]int64 {
	var ret [][]int64

	for _, seat := range this.Seat {
		ret = append(ret, seat.(IPlayer).GetNewBetList())
	}

	return ret
}

// 判断用户是否在座位上
func (this *MgrSeat) IsOnSeat(p interface{}) bool {
	var p1 = p.(IPlayer)
	for _, v := range this.Player {
		p2 := v.(IPlayer)
		if p1.GetUid() == p2.GetUid() {
			return true
		}
	}

	return false
}

//区别GetOtherNewBetList
func (this *MgrSeat) GetOtherNewBetList2(selfUid int64) []int64 {
	// var isSeat = false
	var ret []int64

	if len(this.Player) > 0 {
		list := this.Player[0].(IPlayer).GetNewBetList()
		for range list {
			ret = append(ret, 0)
		}
	}

	for _, p := range this.Player {
		if p.(IPlayer).GetUid() == selfUid {
			continue
		}

		list := p.(IPlayer).GetNewBetList()
		for i, v := range list {
			ret[i] += v
		}
	}

	return ret
}

// 获取其他玩家总下注信息
func (this *MgrSeat) GetOtherNewBetList() []int64 {
	var isSeat = false
	var ret []int64

	if len(this.Player) > 0 {
		list := this.Player[0].(IPlayer).GetNewBetList()
		for range list {
			ret = append(ret, 0)
		}
	}

	for _, p := range this.Player {
		isSeat = false
		for _, seat := range this.Seat {
			if p.(IPlayer).GetUid() == seat.(IPlayer).GetUid() {
				isSeat = true
				break
			}
		}

		if isSeat {
			continue
		}

		list := p.(IPlayer).GetNewBetList()
		for i, v := range list {
			ret[i] += v
		}
	}

	return ret
}

// 获取玩家列表
func (this *MgrSeat) GetUserList(count int) []interface{} {
	var ret []interface{}

	lengthP := len(this.Player)
	if count > lengthP {
		count = lengthP
	}

	for i := 0; i < count; i++ {
		ret = append(ret, this.Player[i])
	}

	return ret
}

///////////////////////////////////////////////////
// 玩家列表排序
// 按身上金币
func (this *MgrSeat) OrderByCoins() {
	util.SortBody(this.Player, func(p, q *interface{}) bool {
		p1 := (*p).(IPlayer)
		p2 := (*q).(IPlayer)

		i := p1.GetCoins()
		j := p2.GetCoins()

		return i < j
	})
}

// 按下注总金币
func (this *MgrSeat) OrderByBetCoins() {
	util.SortBody(this.Player, func(p, q *interface{}) bool {
		p1 := (*p).(IPlayer)
		p2 := (*q).(IPlayer)

		i := p1.GetBetCoins()
		j := p2.GetBetCoins()

		return i > j
	})
}

// 按赢取总金币
func (this *MgrSeat) OrderByWinCoins() {
	util.SortBody(this.Player, func(p, q *interface{}) bool {
		p1 := (*p).(IPlayer)
		p2 := (*q).(IPlayer)

		i := p1.GetWinCoins()
		j := p2.GetWinCoins()

		return i < j
	})
}

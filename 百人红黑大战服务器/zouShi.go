package main

import (
	// "encoding/json"
	"logs"
	"sync"
	"time"
)

var GMgrGameZouShi MgrGameZouShi

func init() {
	GMgrGameZouShi.Zs = make(map[GameTypeDetail]*[]*GameZouShi)
	go GMgrGameZouShi.Run()
}

type MgrGameZouShi struct {
	Zs map[GameTypeDetail]*[]*GameZouShi
	Lk sync.RWMutex
}

func (this *MgrGameZouShi) GetZouShi(gt *GameTypeDetail) []GameZouShi {
	this.Lk.RLock()
	defer this.Lk.RUnlock()
	//
	result := []GameZouShi{}
	v, ok := this.Zs[*gt]
	if !ok {
		return result
	}
	//
	t := time.Now().Unix()
	for _, g := range *v {
		ad := *g
		if ad.GradeNumber != 0 {
			ad.UpdateT = g.UpdateT - t
			if ad.UpdateT < 0 {
				ad.UpdateT = 2
			}
			result = append(result, ad)
		}
	}
	return result
}

func (this *MgrGameZouShi) GetZouShiSingle(gt *GameTypeDetail, serId int32) *GameZouShi {
	this.Lk.RLock()
	defer this.Lk.RUnlock()
	//
	result := GameZouShi{}
	v, ok := this.Zs[*gt]
	if !ok {
		return nil
	}
	//
	for _, g := range *v {
		if g.GradeNumber != 0 {
			result = *g
			result.UpdateT = g.UpdateT - time.Now().Unix()
			if result.UpdateT < 0 {
				result.UpdateT = 2
			}
			return &result
		}
	}

	return nil
}

func (this *MgrGameZouShi) AddZouShi(zs *GameZouShi) {
	this.Lk.Lock()
	defer this.Lk.Unlock()
	//
	v, ok := this.Zs[zs.GameInfo]
	logs.Debug("。。。。。。。。。。更新新走势1", ok, *zs)
	if !ok {
		return
	}
	//

	for _, z := range *v {
		logs.Debug("。。。。。。。。。。更新新走势2")
		z.UpdateT = time.Now().Unix() + zs.UpdateT
		z.GradeNumber = zs.GradeNumber
		z.Data = zs.Data
		break
	}
}

func (this *MgrGameZouShi) TimerSendUpdate() {
	this.Lk.Lock()
	defer this.Lk.Unlock()
	if len(this.Zs) == 0 {
		this.Zs[GameTypeDetail{
			GameType:  int32(GCONFIG.GameType),
			RoomType:  int32(GCONFIG.RoomType),
			GradeType: int32(GCONFIG.GradeType),
		}] = &([]*GameZouShi{
			&GameZouShi{
				SerId: 1,
				GameInfo: GameTypeDetail{
					GameType:  int32(GCONFIG.GameType),
					RoomType:  int32(GCONFIG.RoomType),
					GradeType: int32(GCONFIG.GradeType),
				},
				UpdateT:     time.Now().Unix() + 1,
				GradeNumber: 0,
			},
		})
	}

	t := time.Now().Unix()
	for _, v := range this.Zs {
		for _, zv := range *v {
			jg := t - zv.UpdateT
			if jg > 0 {
				GDeskMgr.AddNativeMsg(MSG_GAME_GETZOUSHI, 0, &GGetZouShi{
					Id: MSG_GAME_GETZOUSHI,
				})
			}
			if jg > 10 {
				zv.GradeNumber = 0
			}
		}
	}
}

func (this *MgrGameZouShi) Run() {
	t1 := time.NewTicker(time.Second)
	for {
		select {
		case <-t1.C:
			this.TimerSendUpdate()
		}
	}
}

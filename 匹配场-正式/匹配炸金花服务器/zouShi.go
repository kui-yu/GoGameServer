package main

import (
	// "encoding/json"
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
	for _, g := range *v {
		result = append(result, *g)
	}
	return result
}

//刷新，删除关掉游戏的走势
func (this *MgrGameZouShi) Update() {
	this.Lk.Lock()
	defer this.Lk.Unlock()
	nowt := time.Now().Unix()
	//
	for _, v := range this.Zs {
		newv := []*GameZouShi{}
		for _, zv := range *v {
			if nowt-zv.UpdateT < 10 {
				newv = append(newv, zv)
			}
		}
		*v = newv
	}
}

func (this *MgrGameZouShi) AddZouShi(zs *GameZouShi) {
	this.Lk.Lock()
	defer this.Lk.Unlock()
	//
	v, ok := this.Zs[zs.GameInfo]
	if !ok {
		this.Zs[zs.GameInfo] = &([]*GameZouShi{zs})
		return
	}
	//
	exist := false
	for _, z := range *v {
		if z.SerId == zs.SerId {
			z.SerId = zs.SerId
			z.UpdateT = time.Now().Unix()
			z.ZouShi = zs.ZouShi
			z.PlayerNum = zs.PlayerNum
			exist = true
			break
		}
	}
	if !exist {
		*v = append(*v, zs)
	}
}

func (this *MgrGameZouShi) TimerSendUpdate() {
	GDeskMgr.AddNativeMsg(MSG_GAME_GETZOUSHI, 0, &GGetZouShi{
		Id: MSG_GAME_GETZOUSHI,
	})
}

func (this *MgrGameZouShi) Run() {
	t1 := time.NewTicker(time.Second * 5)
	t2 := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t1.C:
			this.TimerSendUpdate()
		case <-t2.C:
			this.Update()
		}
	}
}

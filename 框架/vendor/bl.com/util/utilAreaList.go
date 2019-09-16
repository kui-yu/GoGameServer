package util

import (
	"sync"
)

type AreaList struct {
	sync.RWMutex
	areaList []Area
}

func (this *AreaList) Init(count int) {
	this.areaList = []Area{}
	for i := 0; i < count; i++ {
		area := Area{}
		this.areaList = append(this.areaList, area)
	}
}

func (this *AreaList) GetLength() int {
	return len(this.areaList)
}

func (this *AreaList) AddValue(area int, value int64) bool {
	this.RLock()
	defer this.RUnlock()

	if len(this.areaList) <= int(area) {
		return false
	}

	this.areaList[area].AddValue(value)

	return true
}

func (this *AreaList) GetValue(area int) (value int64) {
	this.Lock()
	defer this.Unlock()

	if len(this.areaList) <= area {
		return
	}

	value = this.areaList[area].GetValue()

	return
}

func (this *AreaList) SetValue(area int, value int64) {
	this.Lock()
	defer this.Unlock()

	this.areaList[area].SetValue(value)
}

func (this *AreaList) GetTotValue() (value int64) {
	this.Lock()
	defer this.Unlock()

	for _, v := range this.areaList {
		value += v.GetValue()
	}

	return
}

func (this *AreaList) GetValueList() (list []int64) {
	this.Lock()
	defer this.Unlock()

	for _, v := range this.areaList {
		list = append(list, v.GetValue())
	}

	return
}

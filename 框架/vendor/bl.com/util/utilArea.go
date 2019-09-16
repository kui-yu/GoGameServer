// 线程安全的添加、获取数据
package util

import (
	"sync"
)

type Area struct {
	sync.RWMutex
	value int64
}

func (this *Area) AddValue(v int64) {
	this.Lock()
	defer this.Unlock()

	this.value += v
}

func (this *Area) GetValue() int64 {
	this.RLock()
	defer this.RUnlock()

	return this.value
}

func (this *Area) SetValue(v int64) {
	this.Lock()
	defer this.Unlock()

	this.value = v
}

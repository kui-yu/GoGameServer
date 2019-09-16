package util

import (
	"errors"
	"sync"
)

type List struct {
	sync.RWMutex

	v      []interface{}
	length int
}

func (this *List) SetMaxListLength(length int) {
	this.Lock()
	defer this.Unlock()

	this.length = length
	if len(this.v) > this.length {
		this.v = this.v[:this.length]
	}
}

func (this *List) AddValue(v interface{}) {
	this.Lock()
	defer this.Unlock()

	this.v = append(this.v, v)
	if len(this.v) > this.length {
		this.v = this.v[1:]
	}
}

func (this *List) GetValue(index int) (interface{}, error) {
	this.RLock()
	defer this.RUnlock()

	if index > len(this.v) {
		return 0, errors.New("长度错误")
	}

	return this.v[index], nil
}

func (this *List) GetList() []interface{} {
	this.RLock()
	defer this.RUnlock()

	return this.v
}

func (this *List) GetBoolList() []bool {
	this.RLock()
	defer this.RUnlock()

	var ret []bool
	for _, v := range this.v {
		ret = append(ret, v.(bool))
	}

	return ret
}

func (this *List) GetByteList() []byte {
	this.RLock()
	defer this.RUnlock()

	var ret []byte
	for _, v := range this.v {
		ret = append(ret, v.(byte))
	}

	return ret
}

func (this *List) GetInt32List() []int32 {
	this.RLock()
	defer this.RUnlock()

	var ret []int32
	for _, v := range this.v {
		ret = append(ret, v.(int32))
	}

	return ret
}

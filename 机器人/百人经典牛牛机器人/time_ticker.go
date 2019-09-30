package main

//此文件不用改动
type Timer struct {
	Id int
	H  func(int, interface{})
	T  int //定时时间
	D  interface{}
}

type TimeTicker struct {
	TList []*Timer // 定时器列表
}

//定时器
func (this *TimeTicker) DoTimer() {
	if len(this.TList) == 0 {
		return
	}
	nlist := []*Timer{}
	olist := []*Timer{}
	for _, v := range this.TList {
		v.T--
		if v.T <= 0 {
			olist = append(olist, v)
		} else {
			nlist = append(nlist, v)
		}
	}
	this.TList = nlist
	for _, v := range olist {
		v.H(v.Id, v.D)
	}
}

func (this *TimeTicker) AddTimer(t int, h func(int, interface{}), d interface{}) int {
	var id int = 1
	for {
		var exists bool = false
		for i := 0; i < len(this.TList); i++ {
			if this.TList[i].Id == id {
				exists = true
				id++
				break
			}
		}
		if exists == false {
			break
		}
	}

	this.TList = append(this.TList, &Timer{
		Id: id,
		H:  h,
		T:  t,
		D:  d,
	})

	return id
}

func (this *TimeTicker) DelTimer(id int) {
	for i, v := range this.TList {
		if v.Id == id {
			this.TList = append(this.TList[:i], this.TList[i+1:]...)
			break
		}
	}
}

func (this *TimeTicker) GetTimerNum(id int) int {
	for _, v := range this.TList {
		if v.Id == id {
			return v.T
		}
	}
	return 0
}

//清空定时器
func (this *TimeTicker) ClearTimer() {
	this.TList = []*Timer{}
}

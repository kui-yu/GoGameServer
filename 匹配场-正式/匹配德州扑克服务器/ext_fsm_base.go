/**
* 状态机
**/
package main

type FSMBase interface {
	InitFsm(mark int, extDesk *ExtDesk)

	GetMark() int                        //获得当前状态机标志
	Run(upMark int, args ...interface{}) //运行当前状态
	Leave()                              //离开状态
	Reset()                              //重置状态

	GetRestTime() int64 //剩余时间

	OnUserOnline(p *ExtPlayer)
	OnUserOffline(p *ExtPlayer)
}

/**
* 状态机接口
**/
package main

type FSMBase interface {
	InitFSM(mark int, extDest *ExtDesk)

	GetMark() int
	Run()
	Leave()

	getRestTime() int64 // 剩余时间
}

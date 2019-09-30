/**
* 状态机接口
**/
package main

type FSMBase interface {
	InitFSM(mark int, extDest *ExtDesk)

	GetMark() int
	Run(upMark int)
	Leave()

	getRestTime() int64 // 剩余时间

	onUserOnline(p *ExtPlayer)
	onUserOffline(p *ExtPlayer)
}

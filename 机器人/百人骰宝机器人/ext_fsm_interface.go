/**
* 状态机接口
**/
package main

type FSMBase interface {
	InitFSM(mark int, client *ExtRobotClient)

	GetMark() int
	Run(upMark int)
	Leave()

	onEvent(interface{})
}

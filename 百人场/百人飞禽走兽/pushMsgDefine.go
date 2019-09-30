package main

const (
	MSG_HALL_PUSH_START = 350000 + iota
	MSG_HALL_PUSH_CHANGECOIN
	MSG_HALL_PUSH_CLIENT //350001,消息推送
)

///////////////////////////////////////////////////////////////
//推送后统一应答，如不需要返回自己定义结构体，如有多余的继承这个结构，如下
// type PRspChangeCoin struct {
// 	PBaseMsg
// 	//
// 	Detail string
// }
type PBaseMsg struct {
	Code int
	Err  string
}

////////////////////////////////////////////////////////////////
//推送消息定义,统一以PMsg开头
type PMsgToClientChangeCoin struct {
	Id   int32 //协议号
	Coin int64 //改变的金币
}

type MsgServer struct {
	Uid  int64 //用户uid
	Coin int32 //改变的金币
}

type PMsgToHallWebMsg struct {
	Server map[string]interface{} // 自己需要解析的数据
	Client interface{}            // 直接传递给客户端的数据
}

type PMsgToClientWebMsg struct {
	Id  int32 // 协议号
	Msg interface{}
}

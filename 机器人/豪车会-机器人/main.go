package main

import (
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"time"
)

//定义机器人管理器，包含当前所有机器人的引用,全局变量
var controller ExtController = ExtController{}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	fmt.Println("机器人启动")

	//连接机器人管理中心
	StartRobotServer()
	//机器人管理的初始化
	controller.InitBase()
	controller.Init()
	InitHttpHandle()

	err := http.ListenAndServe(":"+strconv.Itoa(gameConfig.OpenWebInterface.Port), nil)
	if err != nil {
		fmt.Println("web接口启动失败")
	}
	for {
		time.Sleep(time.Second)
	}
}

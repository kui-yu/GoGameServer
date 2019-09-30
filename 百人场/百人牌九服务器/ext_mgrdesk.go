package main

import (
	"logs"
	"os"
)

//此处定义配置的全局变量

//init此处读取配置，初始化全局变量
func init() {
	err := GetGameConfig()
	if err != nil {
		logs.Error("get game config faild! err: #%v", err)
		os.Exit(1)
	}

	// 倍率后面添加一个倍率，用于错误区域倍率
	gameConfig.Double = append(gameConfig.Double, 0)
}

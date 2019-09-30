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
}

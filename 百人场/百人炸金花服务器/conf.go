package main

import (
	"logs"
	"os"
	"strconv"

	//	"strings"

	//	"net/http"

	. "github.com/goconfig"
)

var GCONFIG Config

type Config struct {
	Id          int
	Ip          string
	GameType    int
	RoomType    int
	GradeType   int
	GradeNumber int
	Port        int
	DeskNum     int
	PlayerNum   int
	//
	WebDbIp    string
	WebGameIp  string
	WebRobotIp string
}

func init() {
	InitConfig()
	logs.Debug("网关配置文件初始化结束...", GCONFIG)
}

func InitConfig() {
	path := "conf/conf.ini"
	//
	c, err := LoadConfigFile(path)
	if err != nil {
		logs.Debug("LoadConfigFile failed, path:", path)
		os.Exit(1)
	}
	//----------------------------------
	GCONFIG.Port, err = GetConfigInt(c, "Server", "port")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.Ip, err = GetConfigString(c, "Server", "ip")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.GameType, err = GetConfigInt(c, "Server", "gametype")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.RoomType, err = GetConfigInt(c, "Server", "roomtype")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.GradeType, err = GetConfigInt(c, "Server", "gradetype")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.GradeNumber, err = GetConfigInt(c, "Server", "gradenumber")
	if err != nil {
		GCONFIG.GradeNumber = 0
	}
	GCONFIG.DeskNum, err = GetConfigInt(c, "Server", "desknum")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.PlayerNum, err = GetConfigInt(c, "Server", "playernum")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.WebDbIp, err = GetConfigString(c, "Server", "webdbip")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.WebDbIp = "http://" + GCONFIG.WebDbIp
	//
	GCONFIG.WebGameIp, err = GetConfigString(c, "Server", "webgameip")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.WebGameIp = "http://" + GCONFIG.WebGameIp
	//
	GCONFIG.WebRobotIp, err = GetConfigString(c, "Server", "webrobot")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.WebRobotIp = "http://" + GCONFIG.WebRobotIp
}

func GetConfigInt(c *ConfigFile, region string, key string) (int, error) {
	v, err := c.GetValue(region, key)
	if err != nil {
		logs.Debug("GetValue:", region, key, err)
		return 0, err
	}
	vint, err := strconv.Atoi(v)
	if err != nil {
		logs.Debug("Atoi:", vint, err)
		return 0, err
	}
	//
	return vint, nil
}

func GetConfigString(c *ConfigFile, region string, key string) (string, error) {
	v, err := c.GetValue(region, key)
	if err != nil {
		logs.Debug("GetValue:", region, key, err)
		return "", err
	}
	return v, nil
}

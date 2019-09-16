package main

import (
	"logs"
	"os"
	"strconv"

	. "github.com/goconfig"
)

var GCONFIG Config

type Config struct {
	Ip   string
	Port int

	BgRobotGetUrl     string // 请求获得机器人url
	BgRobotRestUrl    string // 请求归还机器人url
	BgRobotRestAllUrl string // 归还所有机器人url
	BgRobotAddCoinUrl string // 请求添加金币接口
	BgRobotToken      string // 秘钥
	BgRobotHallId     int    // 大厅

	GetGameListUrlLastInterface string // 对外获得游戏列表接口
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

	GCONFIG.BgRobotGetUrl, err = GetConfigString(c, "Server", "bg_robot_get_url")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.BgRobotRestUrl, err = GetConfigString(c, "Server", "bg_robot_reset_url")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.BgRobotRestAllUrl, err = GetConfigString(c, "Server", "bg_robot_resetall_url")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.BgRobotAddCoinUrl, err = GetConfigString(c, "Server", "bg_robot_addcoin_url")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.BgRobotToken, err = GetConfigString(c, "Server", "bg_robot_token")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.BgRobotHallId, err = GetConfigInt(c, "Server", "bg_robot_hallId")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.GetGameListUrlLastInterface, err = GetConfigString(c, "Server", "get_gamelist_urllast")
	if err != nil {
		os.Exit(1)
	}
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

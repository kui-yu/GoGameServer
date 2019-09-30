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
	Ip      string
	Port    int
	Num     int
	ShowUid int64
}

func init() {
	InitConfig()
	logs.Debug("机器人初始化配置结束...", GCONFIG)
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
	GCONFIG.Port, err = GetConfigInt(c, "Robot", "port")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.Ip, err = GetConfigString(c, "Robot", "ip")
	if err != nil {
		os.Exit(1)
	}
	GCONFIG.Num, err = GetConfigInt(c, "Robot", "num")
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

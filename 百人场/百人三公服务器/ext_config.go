package main

import (
	"io/ioutil"
	"logs"

	"github.com/go-yaml/yaml-2"
)

type GameConfigInfo struct {
	Double_Timer    int `yaml:"double_timer"`
	Start_Timer     int `yaml:"start_timer"`
	Bet_Timer       int `yaml:"bet_timer"`
	Stop_Bet_Timer  int `yaml:"stop_bet_timer"`
	Send_Card_Timer int `yaml:"send_card_timer"`
	Open_Timer      int `yaml:"open_timer"`
	Settle_Timer    int `yaml:"settle_timer"`
	Shuffle_Timer   int `yaml:"shuffle_timer"`
}

var gameConfigInfo GameConfigInfo
var maxBet int64

func init() {
	gameConfig, err := ioutil.ReadFile("./ext_conf/gameConf.yaml")
	if err != nil {
		logs.Debug("读取gameConf.yaml失败:%v", err)
	}
	err = yaml.Unmarshal(gameConfig, &gameConfigInfo)
	if err != nil {
		logs.Debug("解析gameConf.yaml失败:%v", err)
	}
	maxBet = GetBetInfo()
}

//获取对应场次的下注限红
func GetBetInfo() int64 {
	if GCONFIG.GradeType == 1 || GCONFIG.GradeType == 6 { //体验,荣耀
		return 500000
	} else if GCONFIG.GradeType == 2 { //王者
		return 800000
	} else if GCONFIG.GradeType == 3 { //战神
		return 2500000
	} else {
		return 0
	}
}

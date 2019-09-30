package main

import (
	"io/ioutil"
	"logs"

	"github.com/go-yaml/yaml-2"
)

var gameConfig GameConfig

type GameTimer struct {
	RandBank    int `yaml:"randbank_ms"`    //随机庄家
	HoleCards   int `yaml:"holecards_ms"`   //发给玩家的两张牌
	PublicCards int `yaml:"publiccard_ms"`  //发公共牌
	UserOperate int `yaml:"useroperate_ms"` //玩家操作
	GameOver    int `yaml:"gameover_ms"`    //等待游戏结束
}

type ALimitInfo struct {
	Id              int     `yaml:"id"`
	RobotNum        int     `yaml:"robotNum"`
	DefSettCoin     int64   `yaml:"defSettCoin"`
	UserSettCoinMin int64   `yaml:"userSettCoinMin"`
	UserSettCoinMax int64   `yaml:"userSettCoinMax"`
	SmallBlind      int64   `yaml:"smallBlind"`
	BigBlind        int64   `yaml:"bigBlind"`
	StageDownBet    []int64 `yaml:"stageDownBet"`
}

type AnalysisConfig struct {
	Timer      GameTimer    `yaml:"timer"`
	LimitInfos []ALimitInfo `yaml:"limitInfo"`
}

type GameConfig struct {
	GameTimer         GameTimer //游戏时间
	RobotNum          int       //机器人数量
	SmallBlindCoin    int64     //小盲
	BigBlindCoin      int64     //大盲
	StageDownBetLimit []int64   //每个阶段的最高下注
	DefSettCoin       int64     //默认设置筹码
	UserSettCoinMin   int64     //玩家最小携带金币
	UserSettCoinMax   int64     //玩家最大携带金币
}

func init() {
	analysisConfig := AnalysisConfig{}
	configFile, err := ioutil.ReadFile("ext_conf/gameConf.yaml")
	err = yaml.Unmarshal(configFile, &analysisConfig)

	DebugLog("gameConfig", analysisConfig)

	if err != nil {
		logs.Error("yamlFile.Get err %v", err)
	}

	gameConfig.GameTimer = analysisConfig.Timer
	for _, limitInfo := range analysisConfig.LimitInfos {
		if GCONFIG.GradeType == limitInfo.Id {
			gameConfig.SmallBlindCoin = limitInfo.SmallBlind
			gameConfig.BigBlindCoin = limitInfo.BigBlind
			gameConfig.StageDownBetLimit = limitInfo.StageDownBet
			gameConfig.RobotNum = limitInfo.RobotNum
			gameConfig.DefSettCoin = limitInfo.DefSettCoin
			gameConfig.UserSettCoinMin = limitInfo.UserSettCoinMin
			gameConfig.UserSettCoinMax = limitInfo.UserSettCoinMax
			break
		}
	}
}

package main

import (
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml-2"
)

var gameConfig GameConfig

var configFile []byte

type GameConfig struct {
	Double    []float64   `yaml:"double"`
	LimitInfo LimitInfo   `yaml:"limitInfo"`
	Undo      Undo        `yaml:"undo"`
	DeskInfo  DeskInfo    `yaml:"deskInfo"`
	Timer     TimerConfig `yaml:"timer"`
}

type LimitInfo struct {
	Limit    []Limit    `yaml:"limit"`
	BetCoins []BetCoins `yaml:"betCoins"`
}

type Limit struct {
	Low  int64 `yaml:"low"`
	High int64 `yaml:"high"`
}

type BetCoins struct {
	Bet []int64 `yaml:"bet"`
}

type Undo struct {
	Warning int32 `yaml:"warning"`
	Exit    int32 `yaml:"exit"`
}

type DeskInfo struct {
	SeatCount     int `yaml:"seatCount"`
	ListCount     int `yaml:"listCount"`
	RunChartCount int `yaml:"runChartCount"`
	BetLimit      int `yaml:"betLimit"`
	Win           int `yaml:"win"`
}

type TimerConfig struct {
	Shuffle     int `yaml:"shuffle"`
	ShuffleNum  int `yaml:"shuffleNum"`
	Ready       int `yaml:"ready"`
	ReadyNum    int `yaml:"readyNum"`
	SendCard    int `yaml:"sendCard"`
	SendCardNum int `yaml:"sendCardNum"`
	Bet         int `yaml:"bet"`
	BetNum      int `yaml:"betNum"`
	StopBet     int `yaml:"stopBet"`
	StopBetNum  int `yaml:"stopBetNum"`
	Open        int `yaml:"open"`
	OpenNum     int `yaml:"openNum"`
	AddNum      int `yaml:"addNum"`
	Award       int `yaml:"award"`
	AwardNum    int `yaml:"awardNum"`
	Over        int `yaml:"over"`
	OverNum     int `yaml:"overNum"`
	NewBet      int `yaml:"newBet"`
	NewBetNum   int `yaml:"newBetNum"`
}

func GetGameConfig() (err error) {
	err = yaml.Unmarshal(configFile, &gameConfig)
	return err
}

func init() {
	var err error
	configFile, err = ioutil.ReadFile("ext_conf/gameConf.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
}

// 下注区域
const (
	// 龙、虎、和
	INDEX_DRAGON int = 1 + iota
	INDEX_TIGER
	INDEX_DRAW

	// 龙   方、梅、红、黑
	INDEX_DRAGONSPADE
	INDEX_DRAGONPLUM
	INDEX_DRAGONRED
	INDEX_DRAGONBLOCK

	// 虎   方、梅、红、黑
	INDEX_TIGERSPADE
	INDEX_TIGERPLUM
	INDEX_TIGERRED
	INDEX_TIGERBLOCK

	// 庄赢、庄输
	INDEX_BANKERWIN
	INDEX_BANKERLOSE

	// 错误下标
	INDEX_ERROR
)

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
	// 天、地、人
	INDEX_TIAN_WIN int = 1 + iota
	INDEX_TIAN_LOSS
	INDEX_DI_WIN
	INDEX_DI_LOSS
	INDEX_REN_WIN
	INDEX_REN_LOSS

	// 天王、至尊
	INDEX_BANKER_TIAN
	INDEX_BANKER_ZHIZUN
	// 错误下标
	INDEX_ERROR
)


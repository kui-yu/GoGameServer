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
	SeatCount          int `yaml:"seatCount"`
	ListCount          int `yaml:"listCount"`
	RunChartCount      int `yaml:"runChartCount"`
	CardTypeChartCount int `yaml:"cardTypeChartCount"`
	BetLimit           int `yaml:"betLimit"`
	Win                int `yaml:"win"`
}

//定时器状态和时间
type TimerConfig struct {
	Shuffle     int `yaml:"shuffle"`
	ShuffleNum  int `yaml:"shuffleNum"` //洗牌时间
	Ready       int `yaml:"ready"`
	ReadyNum    int `yaml:"readyNum"` //准备时间
	SendCard    int `yaml:"sendCard"`
	SendCardNum int `yaml:"sendCardNum"` //发牌时间
	Bet         int `yaml:"bet"`
	BetNum      int `yaml:"betNum"` //下注时间
	StopBet     int `yaml:"stopBet"`
	StopBetNum  int `yaml:"stopBetNum"` //停止下注时间
	Open        int `yaml:"open"`
	OpenNum     int `yaml:"openNum"` //开牌时间
	AddNum      int `yaml:"addNum"`
	Award       int `yaml:"award"`
	AwardNum    int `yaml:"awardNum"` //派奖时间
	Over        int `yaml:"over"`
	OverNum     int `yaml:"overNum"` //重新回到洗牌时间
	NewBet      int `yaml:"newBet"`
	NewBetNum   int `yaml:"newBetNum"` //新下注广播
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
	//红、黑
	INDEX_RED = 1 + iota
	INDEX_BLACK

	//幸运一击
	INDEX_LUCKYBLOW

	//错误下标
	INDEX_ERROR
)

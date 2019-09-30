package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml-2"
)

var gameConfig GameConfig

type GameConfig struct {
	StateInfo       StateInfo       `yaml:"stateInfo"`
	LimitInfo       LimitInfo       `yaml:"limitInfo"`
	GameBetDownUndo GameBetDownUndo `yaml:"undo"`
	Multiple        Multiple        `yaml:"multiple"`
}

type Multiple struct {
	Porsche      float32 `yaml:"porsche"`
	Benz         float32 `yaml:"benz"`
	Ferrari      float32 `yaml:"ferrari"`
	Maserati     float32 `yaml:"maserati"`
	Min_Porsche  float32 `yaml:"min_Porsche"`
	Min_Benz     float32 `yaml:"min_Benz"`
	Min_Ferrari  float32 `yaml:"min_Ferrari"`
	Min_Maserati float32 `yaml:"min_Maserati"`
}
type StateInfo struct {
	DownBet        int `yaml:"downbet"`
	DownBetTime    int `yaml:"downbet_time"`
	RunCarLogo     int `yaml:"runcarlogo"`
	RunCarLogoTime int `yaml:"runcarlogo_time"`
	Balance        int `yaml:"balance"`
	BalanceTime    int `yaml:"balance_time"`
	BroMsg         int `yaml:"bromsg"`
	BroMsgTime     int `yaml:"bromsg_time"`
	Ready          int `yaml:"ready"`
	ReadyTime      int `yaml:"ready_time"`
}
type Bets struct {
	Bet []int64 `yaml:"bet"`
}
type Maxs struct {
	Max int `yaml:"max"`
}
type LimitInfo struct {
	BetLevelCount int    `yaml:"bet_level_count"`
	BetLevels     []Bets `yaml:"bet_levels"`
	BetCount      int    `yaml:"bet_count"`
	AreaMaxCoin   []Maxs `yaml:"areamaxcoin"`
	LogoLimit     int    `yaml:"logolimit"`
	PlayerList    int    `yaml:"playerlist"`
}

type GameBetDownUndo struct {
	Warning int `yaml:"warning"`
	Exit    int `yaml:"exit"`
}

func init() {
	filebuffer, err := ioutil.ReadFile("ext_conf/conf.yaml")
	err = yaml.Unmarshal(filebuffer, &gameConfig)

	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
	fmt.Println("===================gameconfig")
	fmt.Println(gameConfig)
}

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
}

type StateInfo struct {
	DownBet        int `yaml:"downbet"`
	DownBetTime    int `yaml:"downbet_time"`
	RunCarLogo     int `yaml:"runcarlogo"`
	RunCarLogoTime int `yaml:"runcarlogo_time"`
	Balance        int `yaml:"balance"`
	BalanceTime    int `yaml:"balance_time"`
}

type LimitInfo struct {
	BetLevelCount int     `yaml:"bet_level_count"`
	BetLevels     []int64 `yaml:"bet_levels"`
	BetCount      int     `yaml:"bet_count"`
	AreaMaxCoin   int     `yaml:"areamaxcoin"`
	LogoLimit     int     `yaml:"logolimit"`
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

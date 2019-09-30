package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml-2"
)

var gameConfig GameConfig

type GameStatusTimer struct {
	WaitstartId   int `yaml:"waitstart_id"`
	WaitstartMS   int `yaml:"waitstart_ms"`
	RobSeatId     int `yaml:"robseat_id"`
	RobSeatMS     int `yaml:"robseat_ms"`
	FaCardId      int `yaml:"facard_id"`
	FaCardMS      int `yaml:"facard_ms"`
	ShufflecardId int `yaml:"shufflecard_id"`
	ShufflecardMS int `yaml:"shufflecard_ms"`
	DownBetsId    int `yaml:"downbets_id"`
	DownBetsMS    int `yaml:"downbets_ms"`
	OpenCardId    int `yaml:"opencard_id"`
	OpenCardMS    int `yaml:"opencard_ms"`
	BalanceId     int `yaml:"balance_id"`
	BalanceMS     int `yaml:"balance_ms"`
}

type GameBetDownUndo struct {
	Warning int `yaml:"warning"`
	Exit    int `yaml:"exit"`
}

type GameLimtInfo struct {
	BetLevelCount  int     `yaml:"bet_level_count"`
	BetLevels      []int64 `yaml:"bet_levels"`
	SeatCount      int     `yaml:"seat_count"`
	SeatDownCond   int     `yaml:"seatdown_cond"`
	SeatDownMinBet int     `yaml:"seatdown_minbet"`
	SeatDownNum    int     `yaml:"seatdown_num"`
	SeatDownAutoUp int     `yaml:"seatdown_autoup_num"`

	RankCount           int `yaml:"rank_count"`
	UserListCount       int `yaml:"userlist_count"`
	UserListRecordCount int `yaml:"userlist_record_count"`
	RunchartCount       int `yaml:"runchart_count"`

	AreaMaxCoin         int  `yaml:"area_maxcoin"`
	AreaMaxCoinDownSeat int  `yaml:"area_maxcoin_downseat"`
	CompDownbetDouble   int  `yaml:"comp_downbet_double"`
	ExistsMaxminking    bool `yaml:"exists_maxminking"`
}

type GameCtrlInfo struct {
	BankerWinProb  int `yaml:"banker_win_probability"`
	BankerLoseProb int `yaml:"banker_lose_probability"`
}

type GameConfig struct {
	GameStatusTimer GameStatusTimer `yaml:"timer"`
	GameBetDownUndo GameBetDownUndo `yaml:"undo"`
	GameLimtInfo    GameLimtInfo    `yaml:"limitInfo"`
	GameCtrlInfo    GameCtrlInfo    `yaml:"ctrlInfo"`
}

func init() {
	configFile, err := ioutil.ReadFile("ext_conf/gameConf.yaml")
	err = yaml.Unmarshal(configFile, &gameConfig)

	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}

	fmt.Println("===================gameconfig")
	fmt.Println(gameConfig)
}

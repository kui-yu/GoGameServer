package main

import (
	"io/ioutil"
	"logs"

	"github.com/go-yaml/yaml-2"
)

var gameConfig GameConfig

type GameConfig struct {
	LimitInfo       LimitInfo       `yaml:"limitInfo"` //限制信息
	GameStatusTimer GameStatusTimer `yaml:"timer"`
}
type GameStatusTimer struct {
	ShuffleId       int `yaml:"shuffle_id"`
	ShuffleMs       int `yaml:"shuffle_ms"`
	StartdownbetsId int `yaml:"startdownbets_id"`
	StartdownbetsMs int `yaml:"startdownbets_ms"`
	DownBetsId      int `yaml:"downbets_id"`
	DownBetsMS      int `yaml:"downbets_ms"`
	StopdownbetsId  int `yaml:"stopdownbets_id"`
	StopdownbetsMs  int `yaml:"stopdownbets_ms"`
	FaCardId        int `yaml:"facard_id"`
	FaCardMS        int `yaml:"facard_ms"`
	OpenCardId      int `yaml:"opencard_id"`
	OpenCardMS      int `yaml:"opencard_ms"`
	// ThenCardId    int `yaml:"thencard_id"`
	// ThenCardMs    int `yaml:"thencard_ms"`
	BalanceId int `yaml:"balance_id"`
	BalanceMS int `yaml:"balance_ms"`
}
type LimitInfo struct {
	BetLevelCount         int    `yaml:"bet_level_count"` //筹码个数
	BetLevels             []Bets `yaml:"bet_levels"`      //筹码级别
	BetAreaCount          int    `yaml:"bet_count"`       //可下注区域数
	AreaMaxCoin           []Maxs `yaml:"areamaxcoin"`     //区域限红
	SeatCount             int    `yaml:"seatcount"`       //座位个数
	Qzhoushicount         int    `yaml:"qzoushicount"`    //区域走势限制
	Zzhoushicount         int    `yaml:"zzoushicount"`    //庄家走势限制
	Userlist_count        int    `yaml:"userlist_count"`
	Userlist_record_count int    `yaml:"userlist_record_count"` //玩家记录输赢局数限制
	NobetWarning          int    `yaml:"nobetwarning"`          //几局未下注提示
	NobetRemove           int    `yaml:"nobetremove"`           //几局未下注退出
	Downbet_Double_Comp   int64  `yaml:"downbet_coins_comp"`    //携带金币比押注倍数比例
}
type Bets struct {
	Bet []int64 `yaml:"bet"` //筹码值
}
type Maxs struct {
	Max int64 `yaml:"max"` //限红
}

func init() {
	filebuffer, err := ioutil.ReadFile("ext_conf/gameConf.yaml")
	logs.Debug("filebuffer:", filebuffer)
	err = yaml.Unmarshal(filebuffer, &gameConfig)
	if err != nil {
		logs.Error("yaml配置文件出错！！ err is:", err)
	}
	logs.Debug("==================gameConfig:", gameConfig.GameStatusTimer.ShuffleId)
	logs.Debug("==================gameConfig:", gameConfig.GameStatusTimer.ShuffleMs)
}

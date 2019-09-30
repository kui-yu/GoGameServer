package main

import (
	"fmt"
	"io/ioutil"
	"logs"

	"github.com/go-yaml/yaml-2"
)

type GameConfig struct {
	Idle_Multiple               []int64 `yaml:"idle_multiple"`
	Max_Multiple                int64   `yaml:"max_multiple"`
	Rate                        float64 `yaml:"rate"`
	Game_Auto_Timer             int     `yaml:"game_auto_timer"`
	Stage_Start_Timer           int     `yaml:"stage_start_timer"`
	Stage_Shuffle_Cards_Timer   int     `yaml:"stage_shuffle_cards_timer"`
	Stage_Send_Cards_Timer      int     `yaml:"stage_send_cards_timer"`
	Stage_Banker_Multiple_Timer int     `yaml:"stage_banker_multiple_timer"`
	Stage_Idle_Multiple_Timer   int     `yaml:"stage_idle_multiple_timer"`
	Stage_Open_Cards_Timer      int     `yaml:"stage_open_cards_timer"`
	Stage_Game_Settle_Timer     int     `yaml:"stage_game_settle_timer"`
}

var gameConfig GameConfig

func init() {
	f, err := ioutil.ReadFile("./ext_conf/config.yaml")
	if err != nil {
		logs.Debug("读取配置文件config.yaml失败:%v", err)
	}
	err = yaml.Unmarshal(f, &gameConfig)
	if err != nil {
		logs.Debug("config.yaml解析成结构体失败:%v", err)
	}
	fmt.Println("配置文件读取成功:", gameConfig)
	//switch GCONFIG.GradeType {
	//case 1:
	//	gameConfig.Max_Bet = gameConfig.Max_BetArr[1]
	//case 2:
	//	gameConfig.Max_Bet = gameConfig.Max_BetArr[2]
	//case 3:
	//	gameConfig.Max_Bet = gameConfig.Max_BetArr[3]
	//case 6:
	//	gameConfig.Max_Bet = gameConfig.Max_BetArr[0]
	//default:
	//	logs.Debug("GradeType配置错误")
	//}
}

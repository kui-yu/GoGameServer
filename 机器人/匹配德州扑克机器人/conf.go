package main

//此文件不可修改
import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/go-yaml/yaml-2"
)

var gameConfig GameConfig

func (this *GameConfig) getGameConfigString(key string) string {
	for _, v := range gameConfig.RobotConfigItems {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}
func (this *GameConfig) getGameConfigInt(key string) int {
	for _, v := range gameConfig.RobotConfigItems {
		if v.Name == key {
			value, _ := strconv.Atoi(v.Value)
			return value
		}
	}
	return 0
}

// 获取金币配置接口，如果游戏有用到金币相关的，需要调用该方法
func (this *GameConfig) getGameConfigCoin(key string) int64 {
	return int64(this.getGameConfigInt(key)) * 100
}

type GCRobotManager struct {
	Ip   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type GCGameServer struct {
	BackstageUrl        string `yaml:"backstageUrl"`
	WebsocketUrl        string `yaml:"websocketUrl"`
	UrlPathLogin        string `yaml:"urlPathLogin"`
	IsCustomHallConnect bool   `yaml:"isCustomHallConnect"`
	HallIp              string `yaml:"hallIp"`
	HallPort            int    `yaml:"hallPort"`
	LogPrint            bool   `yaml:"logPrint"`
	ErrorPrint          bool   `yaml:"errorPrint"`
}

type GCEnterGame struct {
	Name      string `yaml:"name"`
	Gametype  int    `yaml:"gametype"`
	Roomtype  int    `yaml:"roomtype"`
	Gradetype int    `yaml:"gradetype"`
}

type GCRabbitClient struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
	Pass      string `yaml:"pass"`
	QueueName string `yaml:"queueName"`
}

type RobotConfigItem struct {
	Alias string `yaml:"alias"json:"alias"`
	Name  string `yaml:"name"json:"name"`
	Value string `yaml:"value"json:"value"`
	Type  string `yaml:"type"json:"type"`
}

type OpenWebInterface struct {
	Ip                string `yaml:"Ip"`
	Port              int    `yaml:"Port"`
	GetRobotConfigUrl string `yaml:"getRobotConfigUrl"`
	PutRobotConfigUrl string `yaml:"putRobotConfigUrl"`
	CheckRobotUrl     string `yaml:"checkRobotUrl"`
	OfflineRobotUrl   string `yaml:"offlineRobotUrl"`
	Forceoffroboturl  string `yaml:"forceoffroboturl"`
	Forceonroboturl   string `yaml:"forceonroboturl"`
}

type GameConfig struct {
	GCRobotManager   GCRobotManager    `yaml:"robotManager"`
	GCRabbitClient   GCRabbitClient    `yaml:"rabbitClient"`
	GCGameServer     GCGameServer      `yaml:"gameServer"`
	GCEnterGame      GCEnterGame       `yaml:"enterGame"`
	OpenWebInterface OpenWebInterface  `yaml:"openWebInterface"`
	RobotConfigItems []RobotConfigItem `yaml:"robotConfig"`
}

func init() {
	configFile, err := ioutil.ReadFile("conf/conf.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
	err = yaml.Unmarshal(configFile, &gameConfig)

	fmt.Println(gameConfig)
}

// 保存配置文件
func SaveConf() {
	bytes, err := yaml.Marshal(gameConfig)

	if err != nil {
		ErrorLog("保存配置文件失败", err)
		return
	}

	err = ioutil.WriteFile("conf/conf.yaml", bytes, 0666)
	if err != nil {
		ErrorLog("保存配置文件失败", err)
		return
	}
}

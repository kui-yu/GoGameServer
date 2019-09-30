package main

import (
	"fmt"
	"io/ioutil"
	"logs"

	"github.com/go-yaml/yaml-2"
)

var gameConfig GameConfig

type GameConfig struct {
	LimitInfo  LimitInfo  `yaml:"limitinfo"`  //游戏限制信息
	AreaDouble AreaDouble `yaml:"areadouble"` //游戏倍数
}
type LimitInfo struct {
	BetCount          int   `yaml:"betcount"`
	ChairNum          int   `yaml:"chairnum"`          //座位数量 即展示玩家数量
	MorePlayerNum     int   `yaml:"moreplayernum"`     //更多玩家数量
	HistoryNum        int   `yaml:"historynum"`        //更多玩家显示数量
	NobetWarning      int   `yaml:"nobetwarning"`      //几局未下注提示
	NobetRemove       int   `yaml:"nobetremove"`       //几句未下注踢出
	UpZhuangNeed      int64 `yaml:"upzhuangneed"`      //上庄金额
	ChangeZhuangCount int   `yaml:"changezhuangcount"` //几局之后自动换庄
	WaitZhuangCount   int   `yaml:"waitzhuangcount"`
}
type AreaDouble struct {
	Feiqing  float32 `yaml:"feiqing"`
	Zoushou  float32 `yaml:"zoushou"`
	Tuzi     float32 `yaml:"tuzi"`
	Yanzi    float32 `yaml:"yanzi"`
	Gezi     float32 `yaml:"gezi"`
	Houzi    float32 `yaml:"houzi"`
	Xiongmao float32 `yaml:"xiongmao"`
	Kongque  float32 `yaml:"kongque"`
	Shizi    float32 `yaml:"shizi"`
	Laoying  float32 `yaml:"laoying"`
	Shayu    float32 `yaml:"shayu"`
	Jinshayu float32 `yaml:"jinshayu"`
}

func init() {
	filebuffer, err := ioutil.ReadFile("ext_conf/gameConf.yaml")
	err = yaml.Unmarshal(filebuffer, &gameConfig)
	if err != nil {
		logs.Debug("yaml文件读取失败", err)
	}
	fmt.Println("---------------------------gameconfig")
	fmt.Println(gameConfig)
}

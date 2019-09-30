package main

import (
	"io/ioutil"
	"log"

	yaml "github.com/go-yaml/yaml-2"
)

var gameConfig GameConfig

var configFile []byte

type GameConfig struct {
	DeskInfo DeskInfo `yaml:"deskInfo"`
}

type DeskInfo struct {
	OpenDimension    int             `yaml:"openDimension"`
	DimensionNum     int             `yaml:"dimensionNum"`
	DimensionDoor    []DimensionDoor `yaml:"dimensionDoor"`
	DimensionWeight  []int           `yaml:"dimensionWeight"`
	GameMinBroadcast int             `yaml:"gameMinBroadcast"`
	Win              int             `yaml:"win"`
}

type DimensionDoor struct {
	Coins []int64 `yaml:"coins"`
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

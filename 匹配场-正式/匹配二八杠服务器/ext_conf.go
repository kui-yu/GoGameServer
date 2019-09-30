package main

//此文件不可修改
import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml-2"
)

type ExtRobotConfig struct {
	MaxRobot   int   `yaml:"maxrobot"`
	RobotRate  []int `yaml:"robotrate"`
	PlayerRate int   `yaml:"playerrate"`
}

var GExtRobot ExtRobotConfig

func init() {
	configFile, err := ioutil.ReadFile("ext_conf/robot.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
	err = yaml.Unmarshal(configFile, &GExtRobot)

	fmt.Println(GExtRobot)
}

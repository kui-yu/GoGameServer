package main

//此文件不可修改
import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml-2"
)

var GExtConfig GExtGameConfig

type GExtGameConfig struct {
	PlayTime int `yaml:"playTime"`
}

func init() {
	configFile, err := ioutil.ReadFile("ext_conf/ext_conf.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
	err = yaml.Unmarshal(configFile, &GExtConfig)

	fmt.Println(GExtConfig)
}

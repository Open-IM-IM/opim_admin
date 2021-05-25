package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var Config config

type config struct {
	Secret string `yaml:"secret"`

	Api string `yaml:"api"`

	Mysql struct {
		DBAddress      []string `yaml:"dbAddress"`
		DBUserName     string   `yaml:"dbUserName"`
		DBPassword     string   `yaml:"dbPassword"`
		DBChatName     string   `yaml:"dbChatName"` // 默认使用DBAddress[0]
		DBMsgName      string   `yaml:"dbMsgName"`
		DBMsgTableNum  int      `yaml:"dbMsgTableNum"`
		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`
	}
}

func init() {
	bytes, err := ioutil.ReadFile("../config/config.yaml")
	if err != nil {
		fmt.Errorf("read config fail! err = %v", err)
		return
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		fmt.Errorf("unmarshal config fail! err = %v", err)
		return
	}
	fmt.Println("init config success...")
}

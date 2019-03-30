package conf

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config 配置文件
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	}
	Project struct {
		Name string `yaml:"name"`
	}
	Respository struct {
		Root  string `yaml:"root-path"`
		Store string `yaml:"store-root-path"`
	}
	Log struct {
		File string `yaml:"file"`
	}
	Session struct {
		Name    string `yaml:"name"`
		Timeout int64  `yaml:"timeout"` // 秒
	}
}

// C 配置文件
var C Config

// Init 初始化
func Init() (err error) {

	// 初始化配置文件
	data, err := ioutil.ReadFile("etc/conf")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &C)
	if err != nil {
		return err
	}

	// 打印banner
	printBanner()

	return nil
}

func printBanner() {
	data, err := ioutil.ReadFile("etc/banner.txt")
	if err == nil {
		fmt.Println(string(data))
	}

}

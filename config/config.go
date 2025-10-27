package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App []App `yaml:"apps"` // App Client들을 배열로 관리
}

type App struct {
	App struct {
		Port    string `yaml:"port"`
		Version string `yaml:"version"`
		Name    string `yaml:"name"`
	} `yaml:"app"`

	Http     HttpCfg  `yaml:"http"`
	Producer Producer `yaml:"kafka"`
	// Producer - 1 Producer
	// HTTP     - N Client
}

func NewCfg(path string) Config {
	file, err := os.ReadFile(path)

	if err != nil {
		panic(err.Error())
	}

	var c Config

	err = yaml.Unmarshal(file, &c)

	if err != nil {
		panic(err.Error())
	}

	return c
}

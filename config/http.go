package config

import "api-gateway-module/types/http"

type HttpCfg struct { // 하나의 플랫폼에 하나의 http client와 연결
	Router  []Router `yaml:"router"`
	BaseUrl string   `yaml:"base_url"`
}

type Router struct {
	Method   http.HttpMethod `yaml:"method"`
	GetType  http.GetType    `yaml:"get_type"`
	Variable []string        `yaml:"variable"`
	Path     string          `yaml:"path"`

	Auth   *Auth             `yaml:"auth"`
	Header map[string]string `yaml:"header"`
}

type Auth struct {
	Key   string `yaml:"key"`
	Token string `yaml:"token"`
}

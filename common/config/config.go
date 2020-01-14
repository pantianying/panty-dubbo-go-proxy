package config

import (
	"github.com/pantianying/dubbo-go-proxy/common/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	confConFile = "./conf/proxy.yml"
	Config      *ProxyConfig
)

type ProxyConfig struct {
	HttpListenAddr string `yaml:"http.addr" default:"5s"`

	Retries               string `yaml:"dubbo.retries" default:"0"`
	ResultFiledHumpToLine bool   `yaml:"dubbo.resultFiledHumpToLine" default:true`

	RedisAddr     string `yaml:"redis.addr"`
	RedisPassword string `yaml:"redis.password"`
}

func init() {
	confFileStream, err := ioutil.ReadFile(confConFile)
	if err != nil {
		logger.Error("get config err", err)
	}
	Config = &ProxyConfig{}
	err = yaml.Unmarshal(confFileStream, Config)
	if err != nil {
		logger.Error("get config err", err)
	}
	logger.Debugf("read config :%+v", Config)
}

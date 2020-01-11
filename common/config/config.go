package config

import (
	"github.com/pantianying/dubbo-go-proxy/common/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	confConFile = "./conf/http-proxy.yml"
	Config      *BaseConfig
)

type BaseConfig struct {
	HttpListenAddr string `yaml:"httpListenAddr" default:"5s"`
}

func init() {
	confFileStream, err := ioutil.ReadFile(confConFile)
	if err != nil {
		logger.Error("get config err", err)
	}
	Config = &BaseConfig{}
	err = yaml.Unmarshal(confFileStream, Config)
	if err != nil {
		logger.Error("get config err", err)
	}
}

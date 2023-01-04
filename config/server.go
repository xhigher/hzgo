package config

import (
	"flag"
	"fmt"
	"github.com/xhigher/hzgo/env"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Env string `yaml:"env"`
	Name        string             `yaml:"name"`
	Addr        string             `yaml:"addr"`
	Logger      *LoggerConfig         `yaml:"logger"`
	Mysql       []*MysqlConfig     `yaml:"mysql"`
	Redis       []*RedisConfig     `yaml:"redis"`
	Pay         *PayConfig         `yaml:"pay"`
}

var (
	serverConfig = ServerConfig{}
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "deploy/config.yml", "config file")
}

func GetConfig() *ServerConfig {
	return &serverConfig
}

func Init() {
	flag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error read config %s", configFile))
	}
	err = viper.Unmarshal(&serverConfig)
	if err != nil {
		panic(fmt.Errorf("errpr viper.Unnarshal %v", err))
	}
	env.Init(serverConfig.Env)
}




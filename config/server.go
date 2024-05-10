package config

import (
	"flag"
	"fmt"
	"github.com/xhigher/hzgo/env"
	"github.com/xhigher/hzgo/srd"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Env        string         `mapstructure:"env" json:"env" yaml:"env"`
	Name       string         `mapstructure:"name" json:"name" yaml:"name"`
	Addr       string         `mapstructure:"addr" json:"addr" yaml:"addr"`
	OuterAddr  string         `mapstructure:"outer-addr" json:"outer-addr" yaml:"outer-addr"`
	InnerAddr  string         `mapstructure:"inner-addr" json:"inner-addr" yaml:"inner-addr"`
	MaxReqSize int            `mapstructure:"max-req-size" json:"max-req-size" yaml:"max-req-size"`
	Logger     *LoggerConfig  `mapstructure:"logger" json:"logger" yaml:"logger"`
	JWT        *JWTConfig     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Sec        *SecConfig     `mapstructure:"sec" json:"sec" yaml:"sec"`
	Mysql      []*MysqlConfig `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis      []*RedisConfig `mapstructure:"redis" json:"redis" yaml:"redis"`
	Srd        *srd.Config    `mapstructure:"srd" json:"srd" yaml:"srd"`
}

var (
	serverConfig = ServerConfig{}
	configFile   string
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

	if serverConfig.MaxReqSize == 0 {
		serverConfig.MaxReqSize = 20 << 20
	}
	if len(serverConfig.Addr) == 0 {
		serverConfig.Addr = "0.0.0.0:8888"
	}

	env.Init(serverConfig.Env)
}

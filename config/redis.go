package config

type RedisConfig struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Db       int    `mapstructure:"db" json:"db" yaml:"db"`
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

package config

type RedisConfig struct {
	Name     string `yaml:"name"`
	Db       int    `yaml:"db"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

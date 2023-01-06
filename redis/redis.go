package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/logger"
	"strings"
)
//If you are using Redis 6, install go-redis/v8
//If you are using Redis 7, install go-redis/v9

var (
	clients map[string]*redis.Client
)

func Client(name string) *redis.Client {
	return clients[name]
}

func InitRedis(configs []*config.RedisConfig) {
	if len(configs) == 0 {
		logger.Warnf("redis config nil")
		return
	}
	clients = make(map[string]*redis.Client)
	for _, conf := range configs {
		if len(conf.Name) == 0 {
			conf.Name = "default"
		}
		addr := conf.Addr
		if len(strings.Split(conf.Addr, ":")) == 1 {
			addr = addr + ":6379"
		}
		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: conf.Password,
			DB:       conf.Db,
			PoolSize: 100,
		})
		if client != nil {
			clients[conf.Name] = client
			logger.Infof("Redis init done, name: %s, addr: %s", conf.Name, conf.Addr)
		} else {
			logger.Errorf("Redis init failed, name: %s, addr: %s", conf.Name, conf.Addr)
			return
		}
	}
}

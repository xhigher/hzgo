package cache

import (
	"github.com/go-redis/redis/v8"
	redispool "github.com/xhigher/hzgo/redis"
)

func client() *redis.Client {
	return redispool.Client("misc")
}

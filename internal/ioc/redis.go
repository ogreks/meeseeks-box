package ioc

import (
	"github.com/go-redis/redis"
	"github.com/ogreks/meeseeks-box/configs"
)

// NewRedisClient redis client Cache
func NewRedisClient(cfg configs.Config) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RCache.Driver,
		Password: cfg.RCache.Password,
	})
}

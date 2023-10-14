package core

import (
	"context"
	"hios/config"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

var Cache = cache.New(5*time.Hour, 10*time.Hour)

var Redis *redis.Client

// InitRedis 初始化redis客户端
func InitRedis() error {
	opt, err := redis.ParseURL(config.CONF.Redis.RedisUrl)
	if err != nil {
		return err
	}
	opt.PoolSize = 100
	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	Redis = client
	return nil
}

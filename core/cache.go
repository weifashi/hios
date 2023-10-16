package core

import (
	"context"
	"hios/config"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

var Cache = cache.New(5*time.Hour, 10*time.Hour)

// // 设置缓存
// core.Cache.Set("key", "value", 3*time.Second)

// // 获取缓存
// if value, found := core.Cache.Get("key"); found {
//     // 缓存命中，使用 value
// } else {
//     // 缓存未命中，进行其他逻辑处理
// }

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

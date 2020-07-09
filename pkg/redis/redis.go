package redis

import (
	"time"
	"word/pkg/app"

	"github.com/go-redis/redis"
)

type config struct {
	Addr         string `mapstructure:"addr" toml:"addr" env:"REDIS_URL"`
	Password     string `mapstructure:"password" toml:"password" env:"REDIS_PASS"`
	Db           int    `mapstructure:"db" toml:"db" env:"REDIS_DB"`
	PoolSize     int    `mapstructure:"pool_size" toml:"pool_size" env:"REDIS_POOL_SIZE"`
	MinIdleConns int    `mapstructure:"min_idle_conns" toml:"min_idle_conns" env:"REDIS_MIN_CONN"`
}

var (
	// Client redis连接资源
	Client *redis.Client
	conf   config
)

// Start 启动redis
func Start() {
	err := app.InitConfig("redis", &conf)
	if err != nil {
		return
	}
	Client = redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           conf.Db,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
	})
}

// CacheGet 获取指定key的值,如果值不存在,就执行f方法将返回值存入redis
func CacheGet(key string, expiration time.Duration, f func() string) string {
	cmd := Client.Get(key)
	var val string
	result, _ := cmd.Result()
	if len(result) == 0 {
		val = f()
		Client.Set(key, val, expiration)
		return val
	}
	return result
}

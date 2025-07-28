package redisx

import (
	"Tiktok/configs"
	"Tiktok/log/zlog"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const (
	redisAddr = "%s:%d"
)

func GetRedisClient(config configs.Config) (*redis.Client, error) {
	if !config.Redis.Enable {
		zlog.Warnf("不使用Redis模式")
		return nil, nil
	}
	client := redis.NewClient(&redis.Options{
		Network:         "",
		Addr:            fmt.Sprintf(redisAddr, configs.Conf.Redis.Host, configs.Conf.Redis.Port),
		Dialer:          nil,
		OnConnect:       nil,
		Username:        "",
		Password:        configs.Conf.Redis.Password,
		DB:              configs.Conf.Redis.DB,
		MaxRetries:      0,
		MinRetryBackoff: 0,
		MaxRetryBackoff: 0,
		DialTimeout:     0,
		ReadTimeout:     0,
		WriteTimeout:    0,
		PoolFIFO:        false,
		PoolSize:        1000,
		MinIdleConns:    1,
		PoolTimeout:     0,
		TLSConfig:       nil,
		Limiter:         nil,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		zlog.Fatalf("redis无法链接 %v", err)
		return nil, err
	}
	return client, nil
}

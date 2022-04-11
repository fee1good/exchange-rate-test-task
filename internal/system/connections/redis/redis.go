package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/fee1good/exchange-rate-test-task/internal/config"
)

const (
	readRedisTimeOut  = 200 * time.Millisecond
	writeRedisTimeOut = 200 * time.Millisecond
)

func New(cfg *config.Redis, ctx context.Context) (*redis.Client, error) {
	redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: []string{cfg.URL},
		ReadTimeout:   readRedisTimeOut,
		WriteTimeout:  writeRedisTimeOut,
		Password:      cfg.Password,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return redisClient, nil
}

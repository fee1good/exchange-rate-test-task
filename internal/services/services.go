package services

import (
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"

	"github.com/fee1good/exchange-rate-test-task/internal/config"
	"github.com/fee1good/exchange-rate-test-task/internal/repositories"
	"github.com/fee1good/exchange-rate-test-task/internal/services/rates"
	"github.com/fee1good/exchange-rate-test-task/internal/services/redis/locker"
)

type Container struct {
	Locker locker.Locker
	Rates  *rates.Service
}

func NewContainer(cfg *config.Config, logger *zerolog.Logger, redisClient *redis.Client, repositories *repositories.Container) *Container {
	return &Container{
		Locker: locker.New(redisClient),
		Rates:  rates.NewService(cfg, logger, repositories.Rates),
	}
}

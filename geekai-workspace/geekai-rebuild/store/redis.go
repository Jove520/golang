package store

import (
	"context"
	"time"

	"geekai-rebuild/core/types"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedisClient(appConfig *types.AppConfig, log *zap.SugaredLogger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        appConfig.Redis.Addr(),
		Password:    appConfig.Redis.Password,
		DB:          appConfig.Redis.DB,
		PoolSize:    20,
		PoolTimeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Infow("redis connected", "addr", appConfig.Redis.Addr(), "db", appConfig.Redis.DB)
	return client, nil
}

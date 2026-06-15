package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func CheckRedis(redisClient *redis.Client, log *zap.SugaredLogger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	log.Info("redis health check passed")
	return nil
}

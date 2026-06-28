package cmd

import (
	"context"
	"fmt"
	"grubzo/internal/config"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func getRedisClient(cfg *config.Config) (*redis.Client, error) {
	if cfg.SessionStorage == "local" {
		mr, err := miniredis.Run()
		if err != nil {
			return nil, err
		}

		return redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		}), nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, err
	}

	return rdb, nil
}

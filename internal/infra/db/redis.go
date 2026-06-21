package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mohamed/go-clean-architecture/config"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and validates a connection to Redis
func NewRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	slog.Info("connecting to Redis...")

	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse redis url: %w", err)
	}

	rdb := redis.NewClient(opt)

	// Validate connection
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("unable to ping redis: %w", err)
	}

	slog.Info("successfully connected to Redis")
	return rdb, nil
}

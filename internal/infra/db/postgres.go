package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohamed/go-clean-architecture/config"
)

// NewPostgresPool creates a new high-performance connection pool for PostgreSQL
func NewPostgresPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	slog.Info("connecting to PostgreSQL...")

	poolCfg, err := pgxpool.ParseConfig(cfg.DB.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database url: %w", err)
	}

	// Apply connection pool settings
	poolCfg.MaxConns = cfg.DB.MaxConns
	poolCfg.MinConns = cfg.DB.MinConns
	poolCfg.MaxConnIdleTime = cfg.DB.MaxConnIdleTime
	poolCfg.MaxConnLifetime = cfg.DB.MaxConnLifeTime

	// Set connection timeouts
	poolCfg.ConnConfig.ConnectTimeout = 5 * time.Second

	dbPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := dbPool.Ping(pingCtx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	slog.Info("successfully connected to PostgreSQL connection pool",
		slog.Int("max_conns", int(cfg.DB.MaxConns)),
		slog.Int("min_conns", int(cfg.DB.MinConns)),
	)

	return dbPool, nil
}

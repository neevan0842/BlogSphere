package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neevan0842/BlogSphere/backend/config"
)

func NewDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	// Database connection pool
	poolConfig, err := pgxpool.ParseConfig(config.Envs.DATABASE_URL)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	poolConfig.MaxConns = config.Envs.DB_MAX_OPEN_CONNS
	poolConfig.MinConns = config.Envs.DB_MAX_IDLE_CONNS
	maxIdleTime, err := time.ParseDuration(config.Envs.DB_MAX_IDLE_TIME)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConnIdleTime = maxIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

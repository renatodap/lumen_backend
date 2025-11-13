package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lumen/backend/pkg/logger"
	"go.uber.org/zap"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(dsn string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established",
		zap.Int32("max_conns", config.MaxConns),
		zap.Int32("min_conns", config.MinConns),
	)

	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		logger.Info("Database connection closed")
	}
}

func (db *Database) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return db.Pool.Ping(ctx)
}

func (db *Database) Stats() map[string]interface{} {
	stats := db.Pool.Stat()
	return map[string]interface{}{
		"acquired_conns":     stats.AcquiredConns(),
		"idle_conns":         stats.IdleConns(),
		"total_conns":        stats.TotalConns(),
		"max_conns":          stats.MaxConns(),
		"acquire_count":      stats.AcquireCount(),
		"acquire_duration":   stats.AcquireDuration().String(),
		"canceled_count":     stats.CanceledAcquireCount(),
		"constructing_conns": stats.ConstructingConns(),
		"empty_acquire":      stats.EmptyAcquireCount(),
	}
}

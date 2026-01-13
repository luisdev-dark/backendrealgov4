package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
	poolErr  error
)

// GetPool returns a singleton connection pool to the database.
// It initializes the pool on first call using sync.Once for thread safety.
// This is serverless-friendly as it reuses the pool across invocations.
func GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	poolOnce.Do(func() {
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			poolErr = fmt.Errorf("DATABASE_URL environment variable is not set")
			return
		}

		// Parse and configure the pool
		config, err := pgxpool.ParseConfig(databaseURL)
		if err != nil {
			poolErr = fmt.Errorf("failed to parse DATABASE_URL: %w", err)
			return
		}

		// Ensure SSL mode is required for Neon (if not already in URL)
		// Note: Most Neon URLs already include sslmode=require
		// This is just a safety check
		if config.ConnConfig.TLSConfig == nil {
			poolErr = fmt.Errorf("TLS/SSL is required for database connection. Ensure DATABASE_URL includes sslmode=require")
			return
		}

		// Configure pool settings for serverless environment
		// Keep connections minimal to avoid exhausting database limits
		config.MaxConns = 10
		config.MinConns = 0

		// Create the pool
		pool, poolErr = pgxpool.NewWithConfig(ctx, config)
		if poolErr != nil {
			poolErr = fmt.Errorf("failed to create connection pool: %w", poolErr)
			return
		}

		// Test the connection
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			pool = nil
			poolErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}
	})

	if poolErr != nil {
		return nil, poolErr
	}

	return pool, nil
}

// Close closes the database connection pool.
// This should typically be called during application shutdown.
func Close() {
	if pool != nil {
		pool.Close()
	}
}

package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"restobill/internal/logger"
)

// Client holds the database connection pool
type Client struct {
	Pool *pgxpool.Pool
}

// NewClient creates a new database client and verifies the connection
func NewClient(databaseURL string) (*Client, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	logger.Info("database connection established")
	return &Client{Pool: pool}, nil
}

// Close closes the database connection pool
func (c *Client) Close() {
	c.Pool.Close()
	logger.Info("database connection closed")
}

// Ping checks if the database is reachable
func (c *Client) Ping(ctx context.Context) error {
	return c.Pool.Ping(ctx)
}

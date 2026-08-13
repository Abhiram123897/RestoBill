package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

// --- PostgreSQL Error Helpers (Tharsis pattern) ---

// AsPgError extracts a PostgreSQL error from a generic error
func AsPgError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}

// IsUniqueViolation checks if the error is a unique constraint violation
func IsUniqueViolation(pgErr *pgconn.PgError) bool {
	return pgErr.Code == pgerrcode.UniqueViolation
}

// IsForeignKeyViolation checks if the error is a foreign key constraint violation
func IsForeignKeyViolation(pgErr *pgconn.PgError) bool {
	return pgErr.Code == pgerrcode.ForeignKeyViolation
}

// IsInvalidTextRepresentation checks if the error is an invalid UUID/input format
func IsInvalidTextRepresentation(pgErr *pgconn.PgError) bool {
	return pgErr.Code == pgerrcode.InvalidTextRepresentation
}

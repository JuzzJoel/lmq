package database

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/001_init.sql
var initSQL string

//go:embed migrations/002_add_passwords.sql
var addPasswordsSQL string

// NewPostgresPool creates and returns a new pgxpool connected to the database.
func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}
	
	if config.ConnConfig.TLSConfig != nil {
		config.ConnConfig.TLSConfig.InsecureSkipVerify = true
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}

// RunMigrations executes the embedded SQL migrations.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	log.Printf("[Database]: Running migrations...")
	_, err := pool.Exec(ctx, initSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations (001): %w", err)
	}
	_, err = pool.Exec(ctx, addPasswordsSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations (002): %w", err)
	}
	log.Printf("[Database]: Migrations completed successfully.")
	return nil
}

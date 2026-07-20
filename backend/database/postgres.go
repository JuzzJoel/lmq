package database

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/001_init.sql
var initSQL string

//go:embed migrations/002_add_passwords.sql
var addPasswordsSQL string

//go:embed migrations/003_add_routes.sql
var addRoutesSQL string

//go:embed migrations/004_add_burn_after_reading.sql
var addBurnAfterReadingSQL string

//go:embed migrations/005_add_tags.sql
var addTagsSQL string

// NewPostgresPool creates and returns a new pgxpool connected to the database.
func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}
	
	if config.ConnConfig.TLSConfig != nil {
		config.ConnConfig.TLSConfig.InsecureSkipVerify = true
	}

	config.MaxConns = 3
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	// Implement robust connection retry for serverless DB transient errors (e.g. Supabase cold starts)
	var pingErr error
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		if pingErr = pool.Ping(ctx); pingErr == nil {
			return pool, nil // Success
		}
		log.Printf("[Database]: Ping failed (attempt %d/%d): %v. Retrying in 2 seconds...", i, maxRetries, pingErr)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to ping postgres after %d attempts: %w", maxRetries, pingErr)
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
	_, err = pool.Exec(ctx, addRoutesSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations (003): %w", err)
	}
	_, err = pool.Exec(ctx, addBurnAfterReadingSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations (004): %w", err)
	}
	_, err = pool.Exec(ctx, addTagsSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations (005): %w", err)
	}
	log.Printf("[Database]: Migrations completed successfully.")
	return nil
}

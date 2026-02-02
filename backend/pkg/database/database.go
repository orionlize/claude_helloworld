package database

import (
	"context"
	"fmt"
	"apihub/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	var connString string

	// Use DATABASE_URL if provided (Supabase), otherwise build from components
	if cfg.URL != "" {
		connString = cfg.URL
	} else {
		connString = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
			cfg.SSLMode,
		)
	}

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pool, nil
}

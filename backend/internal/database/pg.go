package database

import (
	"context"
	"fmt"
	"time"

	"github.com/crimsonn/zm_server/internal/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres() (*sqlx.DB, error) {
	config := struct {
		host     string
		port     int
		user     string
		password string
		name     string
		sslMode  string
	}{
		host:     env.GetEnvString("DB_HOST", "localhost"),
		port:     env.GetEnvNumber("DB_PORT", 5432),
		user:     env.GetEnvString("DB_USERNAME", ""),
		password: env.GetEnvString("DB_PASSWORD", ""),
		name:     env.GetEnvString("DB_DATABASE", ""),
		sslMode:  env.GetEnvString("DB_SSLMODE", "disable"),
	}

	if config.user == "" || config.name == "" {
		return nil, fmt.Errorf("DB_USERNAME and DB_DATABASE must be set")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.host,
		config.port,
		config.user,
		config.password,
		config.name,
		config.sslMode,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres connection: %w", err)
	}

	return db, nil
}

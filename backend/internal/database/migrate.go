package database

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

func RunMigrations(db *sqlx.DB, migrationsDir string) error {
	m, err := newMigrator(db, migrationsDir)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply up migrations: %w", err)
	}

	return nil
}

func RollbackLastMigration(db *sqlx.DB, migrationsDir string) error {
	m, err := newMigrator(db, migrationsDir)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply down migration step: %w", err)
	}

	return nil
}

func newMigrator(db *sqlx.DB, migrationsDir string) (*migrate.Migrate, error) {
	if migrationsDir == "" {
		migrationsDir = "./migrations"
	}

	absMigrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("resolve migrations directory %q: %w", migrationsDir, err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("init postgres migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+absMigrationsDir, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return m, nil
}

func closeMigrator(m *migrate.Migrate) {
	sourceErr, dbErr := m.Close()
	if sourceErr != nil || dbErr != nil {
		return
	}
}

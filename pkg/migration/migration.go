package migration

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/minhducta/agent-core-service/pkg/config"
)

// Run applies all pending up-migrations from the given directory.
func Run(dbCfg config.DatabaseConfig, migrationsPath string) (int, error) {
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName, dbCfg.SSLMode,
	)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return 0, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	versionBefore, _, _ := m.Version()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return 0, fmt.Errorf("failed to run migrations: %w", err)
	}

	versionAfter, _, _ := m.Version()
	applied := int(versionAfter) - int(versionBefore)
	if applied < 0 {
		applied = 0
	}

	return applied, nil
}

// Version returns the current migration version and dirty flag.
func Version(dbCfg config.DatabaseConfig, migrationsPath string) (uint, bool, error) {
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName, dbCfg.SSLMode,
	)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	return m.Version()
}

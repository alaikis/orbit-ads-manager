package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"orbit/apps/api/config"
)

type Migration struct {
	Version string
	Content string
}

func RunMigrations() error {
	cfg := config.AppCfg
	if cfg.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	migrationsDir := "migrations"
	if envDir := os.Getenv("MIGRATIONS_DIR"); envDir != "" {
		migrationsDir = envDir
	}

	migrations, err := loadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if len(migrations) == 0 {
		log.Println("no migrations to apply")
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := createMigrationsTable(sqlDB); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := getAppliedVersions(sqlDB)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for _, m := range migrations {
		if _, ok := applied[m.Version]; ok {
			log.Printf("migration %s already applied, skipping", m.Version)
			continue
		}

		log.Printf("applying migration %s", m.Version)
		if _, err := sqlDB.Exec(m.Content); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", m.Version, err)
		}

		if _, err := sqlDB.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", m.Version); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", m.Version, err)
		}

		log.Printf("migration %s applied successfully", m.Version)
	}

	log.Printf("all migrations applied successfully (%d total)", len(migrations))
	return nil
}

func loadMigrations(dir string) ([]Migration, error) {
	var migrations []Migration

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", entry.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			Content: string(content),
		})
	}

	return migrations, nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func getAppliedVersions(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

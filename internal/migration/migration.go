package migration

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Run executes all SQL migrations from the migrations directory
func Run(ctx context.Context, db *pgxpool.Pool, migrationsDir string) error {
	// Create migrations tracking table if not exists
	if err := createMigrationsTable(ctx, db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	files, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Run each migration
	for _, file := range files {
		filename := filepath.Base(file)

		// Check if migration was already applied
		applied, err := isMigrationApplied(ctx, db, filename)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if applied {
			log.Printf("Migration %s already applied, skipping", filename)
			continue
		}

		// Read and execute migration
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		log.Printf("Running migration: %s", filename)
		if _, err := db.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record migration as applied
		if err := recordMigration(ctx, db, filename); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		log.Printf("Migration %s completed successfully", filename)
	}

	return nil
}

// createMigrationsTable creates the schema_migrations table to track applied migrations
func createMigrationsTable(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`
	_, err := db.Exec(ctx, query)
	return err
}

// getMigrationFiles returns sorted list of .sql files from the migrations directory
func getMigrationFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort files to ensure consistent order
	sort.Strings(files)
	return files, nil
}

// isMigrationApplied checks if a migration has already been applied
func isMigrationApplied(ctx context.Context, db *pgxpool.Pool, filename string) (bool, error) {
	var count int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE filename = $1", filename).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// recordMigration records a migration as applied
func recordMigration(ctx context.Context, db *pgxpool.Pool, filename string) error {
	_, err := db.Exec(ctx, "INSERT INTO schema_migrations (filename) VALUES ($1)", filename)
	return err
}

package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"github.com/OderoCeasar/system/config"

	_ "github.com/lib/pq"
)

type Migration struct {
	Version int
	Name    string
	SQL     string
}

// RunMigrations executes all pending database migrations
func RunMigrations(cfg *config.Config) error {
	// Connect to database using database/sql for migrations
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all up migration files
	migrations, err := loadMigrations("db/migrations", "up")
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedVersions, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range migrations {
		if contains(appliedVersions, migration.Version) {
			log.Printf("Migration %d (%s) already applied, skipping", migration.Version, migration.Name)
			continue
		}

		log.Printf("Running migration %d: %s", migration.Version, migration.Name)
		if err := runMigration(db, migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", migration.Version, err)
		}

		log.Printf("Migration %d completed successfully", migration.Version)
	}

	log.Println("All migrations completed successfully")
	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	return err
}

func loadMigrations(dir, direction string) ([]Migration, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	pattern := regexp.MustCompile(`^(\d+)_.*\.(up|down)\.sql$`)
	var migrations []Migration
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		matches := pattern.FindStringSubmatch(file.Name())
		if len(matches) != 3 {
			continue
		}
		if matches[2] != direction {
			continue
		}

		version, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Printf("Warning: Could not parse version from %s, skipping", file.Name())
			continue
		}

		// Read SQL content
		content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", file.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    file.Name(),
			SQL:     string(content),
		})
	}

	// Sort migrations by version (ascending for up, descending for down)
	sort.Slice(migrations, func(i, j int) bool {
		if direction == "down" {
			return migrations[i].Version > migrations[j].Version
		}
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) ([]int, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, rows.Err()
}

func runMigration(db *sql.DB, migration Migration) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(migration.SQL); err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	// Record migration
	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
		migration.Version, migration.Name,
	); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// RollbackMigration rolls back a specific migration (advanced feature)
func RollbackMigration(cfg *config.Config, version int) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	appliedVersions, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	if len(appliedVersions) == 0 {
		return fmt.Errorf("no applied migrations found")
	}

	targetVersion := version
	if targetVersion <= 0 {
		targetVersion = appliedVersions[len(appliedVersions)-1]
	}
	if !contains(appliedVersions, targetVersion) {
		return fmt.Errorf("migration version %d is not applied", targetVersion)
	}

	migrations, err := loadMigrations("db/migrations", "down")
	if err != nil {
		return fmt.Errorf("failed to load down migrations: %w", err)
	}

	var target *Migration
	for i := range migrations {
		if migrations[i].Version == targetVersion {
			target = &migrations[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("down migration for version %d not found", targetVersion)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(target.SQL); err != nil {
		return fmt.Errorf("failed to execute down migration SQL: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", targetVersion); err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Rollback migration %d completed successfully", targetVersion)
	return nil
}

package migrate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/techtonic-org/rf-migrate/pkg/db"
)

// Migrator handles database migrations
type Migrator struct {
	DB            db.DB
	MigrationDir  string
	CurrentSQL    string
	MigrationsDir string
}

// NewMigrator creates a new migrator
func NewMigrator(database db.DB, migrationDir string) (*Migrator, error) {
	// Ensure migrations table exists
	if err := database.EnsureMigrationsTable(); err != nil {
		return nil, err
	}

	return &Migrator{
		DB:            database,
		MigrationDir:  migrationDir,
		CurrentSQL:    filepath.Join(migrationDir, "current.sql"),
		MigrationsDir: filepath.Join(migrationDir, "migrations"),
	}, nil
}

// Apply applies the current SQL migration file
func (m *Migrator) Apply() error {
	content, err := os.ReadFile(m.CurrentSQL)
	if err != nil {
		return fmt.Errorf("failed to read current.sql: %w", err)
	}

	if len(content) == 0 {
		return nil // Nothing to apply
	}

	return m.DB.Execute(string(content))
}

// Watch watches the current SQL file and reapplies it on changes
func (m *Migrator) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Apply initially
	if err := m.Apply(); err != nil {
		return err
	}

	// Watch for changes
	if err := watcher.Add(m.CurrentSQL); err != nil {
		return fmt.Errorf("failed to watch current.sql: %w", err)
	}

	fmt.Println("Watching for changes to current.sql...")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("current.sql modified, reapplying...")
				if err := m.Apply(); err != nil {
					fmt.Printf("Error reapplying: %v\n", err)
				} else {
					fmt.Println("Applied successfully")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}

// Commit commits the current SQL file to a new migration
func (m *Migrator) Commit(name string) error {
	// Read current content
	content, err := os.ReadFile(m.CurrentSQL)
	if err != nil {
		return fmt.Errorf("failed to read current.sql: %w", err)
	}

	if len(content) == 0 {
		return fmt.Errorf("nothing to commit: current.sql is empty")
	}

	// Generate timestamp and filename
	timestamp := time.Now().UTC().Format("20060102150405")
	sanitizedName := strings.ReplaceAll(name, " ", "_")
	fileName := fmt.Sprintf("%s_%s.sql", timestamp, sanitizedName)
	fullPath := filepath.Join(m.MigrationsDir, fileName)

	// Calculate hash
	hash := computeHash(content)

	// Get previous hash
	var previousHash string
	migrations, err := m.DB.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	if len(migrations) > 0 {
		previousHash = migrations[len(migrations)-1].Hash
	}

	// Write migration file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	// Apply the migration and record it
	if err := m.DB.Execute(string(content)); err != nil {
		return fmt.Errorf("failed to apply migration: %w", err)
	}

	// Record the migration
	if err := m.DB.ApplyMigration(fileName, hash, previousHash); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Clear current.sql
	if err := os.WriteFile(m.CurrentSQL, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to clear current.sql: %w", err)
	}

	fmt.Printf("Committed migration: %s\n", fileName)
	return nil
}

// Migrate applies all unapplied migrations
func (m *Migrator) Migrate() error {
	// Get applied migrations
	appliedMigrations, err := m.DB.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get all migration files
	files, err := getFiles(m.MigrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Create a map of applied migration filenames
	appliedFiles := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedFiles[migration.FileName] = true
	}

	// Find and apply unapplied migrations
	appliedCount := 0
	var lastHash string
	if len(appliedMigrations) > 0 {
		lastHash = appliedMigrations[len(appliedMigrations)-1].Hash
	}

	for _, file := range files {
		if !appliedFiles[file] {
			// Read migration file
			fullPath := filepath.Join(m.MigrationsDir, file)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", file, err)
			}

			// Apply migration
			if err := m.DB.Execute(string(content)); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", file, err)
			}

			// Calculate hash
			hash := computeHash(content)

			// Record migration
			if err := m.DB.ApplyMigration(file, hash, lastHash); err != nil {
				return fmt.Errorf("failed to record migration %s: %w", file, err)
			}

			lastHash = hash
			appliedCount++
			fmt.Printf("Applied migration: %s\n", file)
		}
	}

	fmt.Printf("Applied %d migrations\n", appliedCount)
	return nil
}

// Uncommit removes the last migration
func (m *Migrator) Uncommit() error {
	// Remove last migration
	migration, err := m.DB.RemoveLastMigration()
	if err != nil {
		return fmt.Errorf("failed to remove last migration: %w", err)
	}

	// Read migration file
	migrationPath := filepath.Join(m.MigrationsDir, migration.FileName)
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Append to current.sql
	currentContent, err := os.ReadFile(m.CurrentSQL)
	if err != nil {
		return fmt.Errorf("failed to read current.sql: %w", err)
	}

	combinedContent := append(content, currentContent...)
	if err := os.WriteFile(m.CurrentSQL, combinedContent, 0644); err != nil {
		return fmt.Errorf("failed to update current.sql: %w", err)
	}

	// Delete migration file
	if err := os.Remove(migrationPath); err != nil {
		return fmt.Errorf("failed to delete migration file: %w", err)
	}

	fmt.Printf("Uncommitted migration: %s\n", migration.FileName)
	return nil
}

// computeHash calculates a SHA-256 hash of the content
func computeHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// getFiles returns all .sql files in the directory, sorted by name
func getFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

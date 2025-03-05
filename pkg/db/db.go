package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB represents a database connection
type DB interface {
	// Execute runs a SQL query with no rows returned
	Execute(query string) error

	// Query runs a SQL query with rows returned
	Query(query string) (*sql.Rows, error)

	// Close closes the database connection
	Close() error

	// EnsureMigrationsTable ensures that the migrations table exists
	EnsureMigrationsTable() error

	// ApplyMigration applies a migration and records it
	ApplyMigration(fileName string, hash string, previousHash string) error

	// GetAppliedMigrations returns all applied migrations
	GetAppliedMigrations() ([]Migration, error)

	// RemoveLastMigration removes the last migration from the migrations table
	RemoveLastMigration() (Migration, error)
}

// Migration represents a migration record
type Migration struct {
	Hash         string
	PreviousHash string
	FileName     string
	Date         time.Time
}

// PostgresDB is a PostgreSQL implementation of DB
type PostgresDB struct {
	db *sql.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(url string) (DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{db: db}, nil
}

// Execute runs a SQL query with no rows returned
func (pdb *PostgresDB) Execute(query string) error {
	_, err := pdb.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// Query runs a SQL query with rows returned
func (pdb *PostgresDB) Query(query string) (*sql.Rows, error) {
	rows, err := pdb.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	return rows, nil
}

// Close closes the database connection
func (pdb *PostgresDB) Close() error {
	return pdb.db.Close()
}

// EnsureMigrationsTable ensures that the migrations table exists
func (pdb *PostgresDB) EnsureMigrationsTable() error {
	// Create schema if not exists
	schemaQuery := `create schema if not exists rf_migrate;`
	if _, err := pdb.db.Exec(schemaQuery); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Create table if not exists
	query := `
	create table if not exists rf_migrate.migrations (
		hash text primary key,
		previous_hash text,
		file_name text not null,
		date timestamp not null default now()
	);`

	if _, err := pdb.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

// ApplyMigration applies a migration and records it
func (pdb *PostgresDB) ApplyMigration(fileName string, hash string, previousHash string) error {
	tx, err := pdb.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Insert migration record
	query := `
	insert into rf_migrate.migrations (hash, previous_hash, file_name, date)
	values ($1, $2, $3, now());`

	_, err = tx.Exec(query, hash, previousHash, fileName)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return fmt.Errorf("failed to insert migration record: %w", err)
	}

	return tx.Commit()
}

// GetAppliedMigrations returns all applied migrations
func (pdb *PostgresDB) GetAppliedMigrations() ([]Migration, error) {
	query := `
	select hash, previous_hash, file_name, date
	from rf_migrate.migrations
	order by date asc;`

	rows, err := pdb.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		if err := rows.Scan(&m.Hash, &m.PreviousHash, &m.FileName, &m.Date); err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		migrations = append(migrations, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate migration rows: %w", err)
	}

	return migrations, nil
}

// RemoveLastMigration removes the last migration from the migrations table
func (pdb *PostgresDB) RemoveLastMigration() (Migration, error) {
	tx, err := pdb.db.Begin()
	if err != nil {
		return Migration{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Get the latest migration
	query := `
	select hash, previous_hash, file_name, date
	from rf_migrate.migrations
	order by date desc
	limit 1;`

	var m Migration
	err = tx.QueryRow(query).Scan(&m.Hash, &m.PreviousHash, &m.FileName, &m.Date)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return Migration{}, fmt.Errorf("failed to get last migration: %w", err)
	}

	// Delete the migration
	deleteQuery := `delete from rf_migrate.migrations where hash = $1;`
	_, err = tx.Exec(deleteQuery, m.Hash)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return Migration{}, fmt.Errorf("failed to delete migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return Migration{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return m, nil
}

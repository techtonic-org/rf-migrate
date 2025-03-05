# RF-Migrate

A simple, powerful database migration tool for PostgreSQL.

## Features

- **Development-friendly**: Edit and apply your SQL migrations continuously during development
- **Idempotent migrations**: Migrations are designed to be reapplied without errors
- **Version tracking**: Migrations are tracked in the database with hash-based integrity
- **Simple CLI**: Easy to use commands for all migration operations

## Installation

### Using Go Install (Recommended)

```bash
go install github.com/techtonic-org/rf-migrate@latest
```

### Binary Downloads

Visit the [releases page](https://github.com/techtonic-org/rf-migrate/releases) to download the pre-compiled binary for your platform.

### Using as a Tool in Another Go Project

To add rf-migrate as a tool dependency in your Go project:

1. Create a `tools/tools.go` file:
```go
//go:build tools
// +build tools

package tools

import (
	_ "github.com/techtonic-org/rf-migrate" // Import for go mod dependency
)
```

2. Add to your Makefile:
```makefile
# Default database connection URL and migrations directory
DB_URL ?= postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable
MIGRATIONS_DIR ?= ./db/migrations

migrate:
	@echo "Applying migrations..."
	@DATABASE_URL=$(DB_URL) RF_MIGRATION_DIR=$(MIGRATIONS_DIR) go run github.com/techtonic-org/rf-migrate migrate

migrate-watch:
	@echo "Watching migrations..."
	@DATABASE_URL=$(DB_URL) RF_MIGRATION_DIR=$(MIGRATIONS_DIR) go run github.com/techtonic-org/rf-migrate watch

migrate-commit:
	@if [ -z "$(name)" ]; then echo "Error: migration name required. Use 'make migrate-commit name=migration_name'"; exit 1; fi
	@echo "Committing migration: $(name)"
	@DATABASE_URL=$(DB_URL) RF_MIGRATION_DIR=$(MIGRATIONS_DIR) go run github.com/techtonic-org/rf-migrate commit --name "$(name)"
```

## Usage

### Configuration

RF-Migrate can be configured in multiple ways (in order of precedence):

1. **Command-line flags** (highest precedence):
   ```bash
   rf-migrate --database-url "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable" --migration-dir "./migrations"
   ```

2. **Environment variables**:
   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"
   export RF_MIGRATION_DIR="./migrations"
   rf-migrate migrate
   ```

3. **Configuration file** (YAML or JSON):
   ```yaml
   # rfmigrate.yaml example
   databaseUrl: "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"
   migrationDir: "./migrations"
   ```

   The config file location can be specified with `-c` or `--config` flag:
   ```bash
   rf-migrate -c custom-config.yaml migrate
   ```
   
   If not specified, it defaults to `./rfmigrate.yaml` or `./rfmigrate.json`.

### Commands

#### Initialize

First, set up your migration directory structure:

```bash
mkdir -p migrations
touch migrations/current.sql
```

#### Development Workflow

1. **Edit** your `current.sql` file with your database changes
2. **Watch** mode during development:
   ```bash
   rf-migrate watch
   ```
   This continuously applies changes from `current.sql` to your database

3. **Commit** your changes when satisfied:
   ```bash
   rf-migrate commit --name "add_users_table"
   ```
   This creates a timestamped migration file like `20231010123045_add_users_table.sql`

4. **Uncommit** if needed:
   ```bash
   rf-migrate uncommit
   ```
   This restores the last migration to `current.sql`

#### Deployment

Apply all migrations:

```bash
rf-migrate migrate
```

This applies all migration files that haven't been applied yet.

## Migration Format

Migrations should be idempotent, typically using `IF EXISTS` and `IF NOT EXISTS` clauses:

```sql
-- Drop if exists
DROP TABLE IF EXISTS users;

-- Create new
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## License

MIT 
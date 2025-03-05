# RF-Migrate

A simple, powerful roll-forward database migration tool for PostgreSQL written in Go. Inspired by [graphile-migrate](https://github.com/graphile/migrate).

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

#### For Go 1.24+ (Recommended)

Go 1.24 introduces the new `tool` directive that simplifies managing tool dependencies:

```bash
# Add rf-migrate as a tool dependency
go get -tool github.com/techtonic-org/rf-migrate@latest

# Use it directly with go tool
go tool rf-migrate migrate
```

This will add a tool directive to your `go.mod` file:

```
tool github.com/techtonic-org/rf-migrate v1.2.3
```

You can now run rf-migrate directly using `go tool rf-migrate` in your Makefile:

```makefile
# Default database connection URL and migrations directory
DB_URL ?= postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable
MIGRATIONS_DIR ?= ./db/migrations

migrate:
	@echo "Applying migrations..."
	@DATABASE_URL=$(DB_URL) RF_MIGRATION_DIR=$(MIGRATIONS_DIR) go tool rf-migrate migrate

migrate-watch:
	@echo "Watching migrations..."
	@DATABASE_URL=$(DB_URL) RF_MIGRATION_DIR=$(MIGRATIONS_DIR) go tool rf-migrate watch

migrate-commit:
	@if [ -z "$(name)" ]; then echo "Error: migration name required. Use 'make migrate-commit name=migration_name'"; exit 1; fi
	@echo "Committing migration: $(name)"
	@DATABASE_URL=$(DB_URL) RF_MIGRATION_DIR=$(MIGRATIONS_DIR) go tool rf-migrate commit --name "$(name)"
```

#### For Go <1.24 (Legacy Approach)

For older Go versions, use the tools.go approach:

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

## Development

### Semantic Versioning and Releases

This project uses [go-semantic-release](https://github.com/go-semantic-release/semantic-release) to automatically determine the version number and publish releases. New versions are automatically released when changes are pushed to the main branch.

To trigger specific types of releases, use the following commit message format:

- `feat: add new feature` - Triggers a MINOR version bump (e.g., 1.1.0 → 1.2.0)
- `fix: fix a bug` - Triggers a PATCH version bump (e.g., 1.1.0 → 1.1.1)
- `feat!: add feature with breaking change` - Triggers a MAJOR version bump (e.g., 1.1.0 → 2.0.0)

Other commit types (docs, style, refactor, perf, test, chore) won't trigger a new release.

Run `make help` to see a full list of supported commit types and formats.

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

The tool will automatically create the migration directory and an empty `current.sql` file if they don't exist. 

However, you can also manually set up your migration directory structure:

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
   This creates a timestamped migration file like `20231010123045_add_users_table.sql` in the same directory as `current.sql`

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

This applies all timestamped migration files in the specified migration directory that haven't been applied yet.

## Migration Format

Migrations should be idempotent, typically using `IF EXISTS` and `IF NOT EXISTS` clauses:

```sql
-- Drop if exists
drop table if exists users;

-- Create new
create table if not exists users (
  id serial primary key,
  username varchar(255) not null unique,
  created_at timestamp not null default now()
);
```

## License

MIT 
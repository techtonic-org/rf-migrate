.PHONY: build test clean help

VERSION := $(shell cat VERSION 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

# Build binary
build:
	@echo "Building rf-migrate..."
	@go build $(LDFLAGS) -o rf-migrate

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f rf-migrate
	@rm -rf dist

# Install from source
install:
	@echo "Installing rf-migrate..."
	@go install $(LDFLAGS)

# Sample test command with PostgreSQL
test-postgres:
	@echo "Running with local PostgreSQL..."
	@export DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" && \
	export RF_MIGRATION_DIR="./test/migrations" && \
	go run main.go config

# Show help for using semantic commits
help:
	@echo "RF-Migrate uses semantic-release for versioning"
	@echo ""
	@echo "To trigger releases, use semantic commit messages:"
	@echo "  feat: add new feature (triggers MINOR version bump)"
	@echo "  fix: fix a bug (triggers PATCH version bump)"
	@echo "  docs: documentation changes only"
	@echo "  style: formatting, missing semi colons, etc"
	@echo "  refactor: code change that neither fixes a bug nor adds a feature"
	@echo "  perf: code change that improves performance"
	@echo "  test: adding missing tests"
	@echo "  chore: updating build tasks, package manager configs, etc"
	@echo ""
	@echo "Breaking changes (trigger MAJOR version bump):"
	@echo "  feat!: add feature with breaking change"
	@echo "  or include 'BREAKING CHANGE:' in commit body" 
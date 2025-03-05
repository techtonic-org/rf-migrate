.PHONY: build test clean release

VERSION := $(shell cat VERSION)
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

# Create a new release tag
release:
	@echo "Current version: $(VERSION)"
	@read -p "Enter new version: " VERSION; \
	echo $$VERSION > VERSION; \
	git add VERSION; \
	git commit -m "Bump version to $$VERSION"; \
	git tag -a v$$VERSION -m "Version $$VERSION"; \
	echo "Run 'git push && git push --tags' to publish"

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
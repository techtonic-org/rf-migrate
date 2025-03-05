package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	DatabaseURL  string `mapstructure:"databaseUrl"`
	MigrationDir string `mapstructure:"migrationDir"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("databaseUrl", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	v.SetDefault("migrationDir", "./migrations")

	// If config file is provided
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Look for default config files
		v.SetConfigName("rfmigrate")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, we'll use defaults and CLI flags
	}

	// Bind environment variables
	if err := v.BindEnv("databaseUrl", "DATABASE_URL"); err != nil {
		return nil, fmt.Errorf("failed to bind environment variable: %w", err)
	}
	if err := v.BindEnv("migrationDir", "RF_MIGRATION_DIR"); err != nil {
		return nil, fmt.Errorf("failed to bind environment variable: %w", err)
	}

	// Read environment variables
	v.AutomaticEnv()

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Ensure migration directory exists
	if err := ensureMigrationDir(config.MigrationDir); err != nil {
		return nil, err
	}

	return &config, nil
}

// ensureMigrationDir makes sure the migration directory and its structure exist
func ensureMigrationDir(dir string) error {
	// Create main migration directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create migration directory: %w", err)
	}

	// Create current.sql if it doesn't exist
	currentSqlPath := filepath.Join(dir, "current.sql")
	if _, err := os.Stat(currentSqlPath); os.IsNotExist(err) {
		emptyFile, err := os.Create(currentSqlPath)
		if err != nil {
			return fmt.Errorf("failed to create current.sql: %w", err)
		}
		emptyFile.Close()
	}

	return nil
}

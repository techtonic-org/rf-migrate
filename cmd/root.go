package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/techtonic-org/rf-migrate/pkg/config"
	"github.com/techtonic-org/rf-migrate/pkg/db"
	"github.com/techtonic-org/rf-migrate/pkg/migrate"
)

var (
	// Used for flags
	cfgFile      string
	databaseURL  string
	migrationDir string
	showVersion  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rf-migrate",
	Short: "A simple database migration tool",
	Long: `rf-migrate is a CLI tool for managing database migrations.
It allows you to develop, apply, and track database schema changes.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./rfmigrate.yaml)")
	rootCmd.PersistentFlags().StringVar(&databaseURL, "database-url", "", "Database connection URL")
	rootCmd.PersistentFlags().StringVar(&migrationDir, "migration-dir", "", "Directory for migration files")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Show version information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Config already provided via flag
}

// loadConfig loads configuration and applies CLI flag overrides
func loadConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	// Override with CLI flags if provided
	if databaseURL != "" {
		cfg.DatabaseURL = databaseURL
	}
	if migrationDir != "" {
		cfg.MigrationDir = migrationDir
	}

	return cfg, nil
}

// createMigrator creates a new migrator instance
func createMigrator() (*migrate.Migrator, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	// Connect to database
	database, err := db.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Create migrator
	return migrate.NewMigrator(database, cfg.MigrationDir)
}

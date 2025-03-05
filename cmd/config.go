package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd shows the loaded configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the loaded configuration",
	Long:  `Displays the current configuration after resolving from all sources: defaults, config file, environment variables, and CLI flags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		fmt.Println("Current Configuration:")
		fmt.Println("----------------------")
		fmt.Printf("Database URL: %s\n", cfg.DatabaseURL)
		fmt.Printf("Migration Directory: %s\n", cfg.MigrationDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

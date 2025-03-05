package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply all unapplied migrations",
	Long: `Applies all migration files from the migrations directory that have not yet been applied.
This command is typically used in production or staging environments to bring
the database schema up to date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator, err := createMigrator()
		if err != nil {
			return err
		}

		fmt.Println("Applying migrations...")
		if err := migrator.Migrate(); err != nil {
			return err
		}

		fmt.Println("All migrations applied successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

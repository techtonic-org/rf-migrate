package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the current.sql migration",
	Long: `Applies the current.sql migration file to the database.
This is a one-time application, unlike the watch command which
continuously applies the migration on changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator, err := createMigrator()
		if err != nil {
			return err
		}

		fmt.Println("Applying current.sql...")
		if err := migrator.Apply(); err != nil {
			return err
		}

		fmt.Println("Successfully applied current.sql")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

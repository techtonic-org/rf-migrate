package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var commitName string

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit current.sql to a migration file",
	Long: `Commits the current.sql file to a timestamped migration file in the migrations directory.
The migration is applied and recorded in the migrations table.
The current.sql file is cleared after a successful commit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if commitName == "" {
			return errors.New("migration name is required")
		}

		migrator, err := createMigrator()
		if err != nil {
			return err
		}

		fmt.Printf("Committing migration '%s'...\n", commitName)
		if err := migrator.Commit(commitName); err != nil {
			return err
		}

		fmt.Println("Migration committed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&commitName, "name", "n", "", "Name for the migration (required)")
	if err := commitCmd.MarkFlagRequired("name"); err != nil {
		fmt.Printf("Error marking flag as required: %v\n", err)
	}
}

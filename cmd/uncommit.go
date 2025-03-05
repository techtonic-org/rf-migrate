package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// uncommitCmd represents the uncommit command
var uncommitCmd = &cobra.Command{
	Use:   "uncommit",
	Short: "Uncommit the last migration",
	Long: `Removes the last migration from the migrations table,
deletes the migration file, and puts its content back into current.sql.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator, err := createMigrator()
		if err != nil {
			return err
		}

		fmt.Println("Uncommitting the last migration...")
		if err := migrator.Uncommit(); err != nil {
			return err
		}

		fmt.Println("Last migration uncommitted successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uncommitCmd)
}

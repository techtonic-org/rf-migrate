package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch current.sql and reapply on changes",
	Long: `Watches the current.sql file and reapplies it to the database whenever it changes.
This is useful during development to test your migrations without having to manually reapply them.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator, err := createMigrator()
		if err != nil {
			return err
		}

		fmt.Println("Watching current.sql for changes...")
		return migrator.Watch()
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

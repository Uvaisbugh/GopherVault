package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm"},
	Short:   "Delete an entry from the vault",
	Long: `Remove an entry by its ID or ID prefix.

Example:
  gopher delete a1b2c3d4
  gopher rm a1b2c3d4`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		if err := vaultService.DeleteEntry(id); err != nil {
			return err
		}

		fmt.Printf("🗑️ Entry '%s' deleted successfully.\n", id)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import <filename>",
	Short: "Import entries from a JSON backup file",
	Long: `Import and restore knowledge entries from an exported JSON file into GopherVault.

Example:
  gopher import backup.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]

		count, err := vaultService.Import(source)
		if err != nil {
			return err
		}

		fmt.Printf("📥 Successfully imported %d %s from '%s'\n", count, pluralize(count, "entry", "entries"), source)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [filename]",
	Short: "Export all entries to a JSON backup file",
	Long: `Export all knowledge entries stored in GopherVault to a JSON file.

Example:
  gopher export backup.json
  gopher export my-vault-2026.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dest := "vault-backup.json"
		if len(args) > 0 && args[0] != "" {
			dest = args[0]
		}

		if err := vaultService.Export(dest); err != nil {
			return err
		}

		fmt.Printf("📦 Successfully exported all vault entries to '%s'\n", dest)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}

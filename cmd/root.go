// Package cmd defines all CLI commands using the Cobra library.
//
// 📚 Go concept: In Go, Cobra is the industry-standard library for building
// production-grade CLI tools (used by Kubernetes kubectl, Hugo, GitHub CLI gh, etc.).
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Uvaisbugh/gophervault/internal/service"
	"github.com/Uvaisbugh/gophervault/internal/storage"
)

var (
	// customDBPath is a global flag to allow specifying a custom JSON storage file.
	customDBPath string

	// vaultService holds the shared business logic instance initialized in PersistentPreRun.
	vaultService *service.VaultService
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "gopher",
	Short: "GopherVault is a local-first CLI knowledge manager for developers",
	Long: `🐹 GopherVault
A local-first command-line knowledge manager built in Go.

Save code snippets, terminal commands, debug notes, and configuration examples
so you never have to search for solutions you've already solved.

Get started:
  gopher add --title "Docker prune" --tag devops --body "docker system prune -a"
  gopher search "docker"
  gopher list
  gopher stats`,
	// PersistentPreRunE runs before any child command executes.
	// It initializes the storage repository and service.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewJSONStore(customDBPath)
		if err != nil {
			return fmt.Errorf("failed to initialize vault store: %w", err)
		}
		vaultService = service.NewVaultService(store)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main().
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// 📚 Go concept: "init()" runs automatically before main() when the package loads.
	// We bind persistent flags here (flags available to root and all subcommands).
	RootCmd.PersistentFlags().StringVar(&customDBPath, "db", "", "custom path to vault.json storage file")
}

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Uvaisbugh/gophervault/internal/models"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search saved knowledge by keyword",
	Long: `Search your stored knowledge entries.

Searches across titles, snippet bodies, and tags.

Example:
  gopher search "docker"
  gopher search prune`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		results, err := vaultService.SearchEntries(query)
		if err != nil {
			return err
		}

		if len(results) == 0 {
			fmt.Printf("🔍 No entries found matching '%s'.\n", query)
			return nil
		}

		fmt.Printf("🔍 Found %d %s matching '%s':\n\n", len(results), pluralize(len(results), "entry", "entries"), query)
		for _, e := range results {
			printEntryResult(e)
		}

		return nil
	},
}

func printEntryResult(e *models.Entry) {
	shortID := e.ID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}

	tags := strings.Join(e.Tags, ", ")
	if tags == "" {
		tags = "-"
	}

	fmt.Printf("  • [%s] %s\n", shortID, e.Title)
	fmt.Printf("    Category: %s | Tags: [%s]\n", e.Category, tags)

	preview := e.Body
	if len(preview) > 75 {
		preview = preview[:72] + "..."
	}
	fmt.Printf("    Snippet:  %s\n\n", preview)
}

func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

func init() {
	RootCmd.AddCommand(searchCmd)
}

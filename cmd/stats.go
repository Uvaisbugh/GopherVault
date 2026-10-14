package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics about saved entries and popular tags",
	Long: `Display high-level metrics about your knowledge base:
total entries, unique tags, categories, and most frequently used tags.

Example:
  gopher stats`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := vaultService.GetStats()
		if err != nil {
			return err
		}

		fmt.Println("📊 GopherVault Statistics")
		fmt.Println("────────────────────────────────────────────")
		fmt.Printf("  Total Entries:  %d\n", stats.TotalEntries)
		fmt.Printf("  Unique Tags:    %d\n", stats.TotalTags)
		fmt.Printf("  Categories:     %d\n", len(stats.Categories))
		fmt.Println("────────────────────────────────────────────")

		if len(stats.Categories) > 0 {
			fmt.Println("📁 Entries by Category:")
			for cat, count := range stats.Categories {
				fmt.Printf("  • %-16s : %d\n", cat, count)
			}
			fmt.Println()
		}

		if len(stats.PopularTags) > 0 {
			fmt.Println("🏷️  Popular Tags:")
			maxDisplay := 10
			if len(stats.PopularTags) < maxDisplay {
				maxDisplay = len(stats.PopularTags)
			}
			for i := 0; i < maxDisplay; i++ {
				item := stats.PopularTags[i]
				bar := strings.Repeat("■", item.Count)
				fmt.Printf("  %-16s %-15s (%d)\n", item.Tag, bar, item.Count)
			}
			fmt.Println()
		}

		if stats.TotalEntries == 0 {
			fmt.Println("Your vault is currently empty! Try adding an entry:")
			fmt.Println("  gopher add --title \"My First Note\" --body \"hello world\" --tag initial")
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
}

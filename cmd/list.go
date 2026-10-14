package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	listTag      string
	listCategory string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List entries filtered by tag or category",
	Long: `List entries saved in GopherVault.

Optionally filter by tag or category.

Examples:
  gopher list
  gopher list --tag devops
  gopher list --category git
  gopher list --tag devops --category docker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := vaultService.ListEntries(listCategory, listTag)
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			filterDesc := ""
			if listTag != "" {
				filterDesc += fmt.Sprintf(" with tag '%s'", listTag)
			}
			if listCategory != "" {
				filterDesc += fmt.Sprintf(" in category '%s'", listCategory)
			}
			fmt.Printf("📂 No entries found%s.\n", filterDesc)
			return nil
		}

		fmt.Printf("📂 Saved Knowledge (%d %s):\n\n", len(entries), pluralize(len(entries), "entry", "entries"))

		// Print header
		fmt.Printf("  %-10s  %-24s  %-12s  %s\n", "ID", "TITLE", "CATEGORY", "TAGS")
		fmt.Printf("  %-10s  %-24s  %-12s  %s\n", "----------", "------------------------", "------------", "----")

		for _, e := range entries {
			shortID := e.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}

			title := e.Title
			if len(title) > 24 {
				title = title[:21] + "..."
			}

			category := e.Category
			if len(category) > 12 {
				category = category[:10] + ".."
			}

			tags := strings.Join(e.Tags, ", ")
			if tags == "" {
				tags = "-"
			}

			fmt.Printf("  %-10s  %-24s  %-12s  %s\n", shortID, title, category, tags)
		}

		fmt.Println("\nTip: Use 'gopher show <id>' to inspect full note and snippet.")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listTag, "tag", "", "filter entries by tag")
	listCmd.Flags().StringVarP(&listCategory, "category", "c", "", "filter entries by category")
}

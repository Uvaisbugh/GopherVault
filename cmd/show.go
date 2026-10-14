package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Display the full details of an entry",
	Long: `Display the full details, timestamps, tags, and body snippet of an entry by ID.

You can provide either the full ID or a prefix (e.g. first 6-8 characters).

Example:
  gopher show a1b2c3d4`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		entry, err := vaultService.GetEntry(id)
		if err != nil {
			return err
		}

		tags := strings.Join(entry.Tags, ", ")
		if tags == "" {
			tags = "none"
		}

		fmt.Println("────────────────────────────────────────────────────────────")
		fmt.Printf("📌 %s\n", entry.Title)
		fmt.Println("────────────────────────────────────────────────────────────")
		fmt.Printf("ID:         %s\n", entry.ID)
		fmt.Printf("Category:   %s\n", entry.Category)
		fmt.Printf("Tags:       %s\n", tags)
		fmt.Printf("Created:    %s\n", entry.CreatedAt.Format("2006-01-02 15:04:05 MST"))
		fmt.Printf("Updated:    %s\n", entry.UpdatedAt.Format("2006-01-02 15:04:05 MST"))
		fmt.Println("────────────────────────────────────────────────────────────")
		fmt.Println("CONTENT / SNIPPET:")
		fmt.Println(entry.Body)
		fmt.Println("────────────────────────────────────────────────────────────")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(showCmd)
}

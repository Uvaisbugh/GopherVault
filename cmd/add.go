package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	addTitle    string
	addBody     string
	addCategory string
	addTags     []string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Save a new note, command, or code snippet",
	Long: `Save a new note or code snippet to GopherVault.

You can supply values via command-line flags or run interactively if flags are omitted.

Examples:
  gopher add --title "Docker cleanup" --tag devops --body "docker system prune -a"
  gopher add -t "Git stash pop" -b "git stash pop" -c git --tag git,vcs
  gopher add (interactive mode)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		// Interactive fallback if title is not passed
		if strings.TrimSpace(addTitle) == "" {
			fmt.Print("Title: ")
			line, _ := reader.ReadString('\n')
			addTitle = strings.TrimSpace(line)
		}

		// Interactive fallback if body is not passed
		if strings.TrimSpace(addBody) == "" {
			fmt.Print("Snippet / Body: ")
			line, _ := reader.ReadString('\n')
			addBody = strings.TrimSpace(line)
		}

		// If category not specified, prompt or default
		if strings.TrimSpace(addCategory) == "" && len(args) == 0 && addTitle != "" {
			// keep default "general" if blank
			if addCategory == "" {
				addCategory = "general"
			}
		}

		// Parse comma-separated tags if supplied like --tag "go,cli,dev"
		var parsedTags []string
		for _, tagGroup := range addTags {
			for _, t := range strings.Split(tagGroup, ",") {
				t = strings.TrimSpace(t)
				if t != "" {
					parsedTags = append(parsedTags, t)
				}
			}
		}

		entry, err := vaultService.AddEntry(addTitle, addBody, addCategory, parsedTags)
		if err != nil {
			return err
		}

		fmt.Println("\n✨ Entry saved successfully!")
		fmt.Printf("   ID:       %s\n", entry.ID)
		fmt.Printf("   Title:    %s\n", entry.Title)
		fmt.Printf("   Category: %s\n", entry.Category)
		if len(entry.Tags) > 0 {
			fmt.Printf("   Tags:     [%s]\n", strings.Join(entry.Tags, ", "))
		}
		fmt.Printf("   Created:  %s\n\n", entry.CreatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addTitle, "title", "t", "", "title of the snippet or note")
	addCmd.Flags().StringVarP(&addBody, "body", "b", "", "content or code snippet")
	addCmd.Flags().StringVarP(&addCategory, "category", "c", "general", "category (e.g. devops, git, go, database)")
	addCmd.Flags().StringSliceVar(&addTags, "tag", []string{}, "tags for organizing knowledge (can be repeated or comma-separated)")
}

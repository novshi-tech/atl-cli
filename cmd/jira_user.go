package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage Jira users",
}

var userSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for users by name or email",
	RunE:  runUserSearch,
}

func init() {
	userSearchCmd.Flags().StringP("query", "q", "", "Search string (display name or email)")
	userSearchCmd.MarkFlagRequired("query")
	userSearchCmd.Flags().Int("max", 50, "Maximum number of results")

	userCmd.AddCommand(userSearchCmd)
	jiraCmd.AddCommand(userCmd)
}

func runUserSearch(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	query, _ := cmd.Flags().GetString("query")
	max, _ := cmd.Flags().GetInt("max")

	users, err := client.SearchUsers(query, max)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		items := make([]JSONUserItem, len(users))
		for i, u := range users {
			items[i] = JSONUserItem{
				AccountID:    u.AccountID,
				DisplayName:  u.DisplayName,
				EmailAddress: u.EmailAddress,
				Active:       u.Active,
			}
		}
		return printJSON(items)
	}

	if len(users) == 0 {
		fmt.Println("No users found.")
		return nil
	}

	fmt.Printf("Found %d user(s):\n\n", len(users))
	for _, u := range users {
		email := ""
		if u.EmailAddress != "" {
			email = u.EmailAddress
		}
		active := "active"
		if !u.Active {
			active = "inactive"
		}
		fmt.Printf("%-40s  %-25s  %-30s  %s\n",
			u.AccountID,
			u.DisplayName,
			email,
			active,
		)
	}
	return nil
}

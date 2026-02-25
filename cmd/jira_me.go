package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var jiraMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show the currently authenticated Jira user",
	RunE:  runJiraMe,
}

func init() {
	jiraCmd.AddCommand(jiraMeCmd)
}

func runJiraMe(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	user, err := client.GetMyself()
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		return printJSON(JSONUserItem{
			AccountID:    user.AccountID,
			DisplayName:  user.DisplayName,
			EmailAddress: user.EmailAddress,
			Active:       user.Active,
		})
	}

	active := "active"
	if !user.Active {
		active = "inactive"
	}
	fmt.Printf("Account ID:  %s\n", user.AccountID)
	fmt.Printf("Name:        %s\n", user.DisplayName)
	fmt.Printf("Email:       %s\n", user.EmailAddress)
	fmt.Printf("Active:      %s\n", active)
	return nil
}

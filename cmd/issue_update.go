package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing Jira issue",
	RunE:  runIssueUpdate,
}

func init() {
	issueUpdateCmd.Flags().StringP("key", "k", "", "Issue key (required)")
	issueUpdateCmd.MarkFlagRequired("key")
	issueUpdateCmd.Flags().StringP("summary", "s", "", "New summary")
	issueUpdateCmd.Flags().StringP("description", "d", "", "New description")
	issueUpdateCmd.Flags().String("status", "", "Transition to this status")
	issueCmd.AddCommand(issueUpdateCmd)
}

func runIssueUpdate(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	key, _ := cmd.Flags().GetString("key")
	summary, _ := cmd.Flags().GetString("summary")
	description, _ := cmd.Flags().GetString("description")
	status, _ := cmd.Flags().GetString("status")

	if summary != "" || description != "" {
		if err := client.UpdateIssue(key, summary, description); err != nil {
			return err
		}
		fmt.Printf("Updated issue: %s\n", key)
	}

	if status != "" {
		if err := client.TransitionIssue(key, status); err != nil {
			return err
		}
		fmt.Printf("Transitioned %s to %q\n", key, status)
	}

	if summary == "" && description == "" && status == "" {
		fmt.Println("Nothing to update. Specify --summary, --description, or --status.")
	}

	fmt.Printf("URL: %s/browse/%s\n", client.BaseURL(), key)
	return nil
}

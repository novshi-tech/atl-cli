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
	issueUpdateCmd.Flags().String("assignee", "", "Assignee account ID (use \"none\" to unassign)")
	issueUpdateCmd.Flags().String("due", "", "Due date (YYYY-MM-DD)")
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
	assignee, _ := cmd.Flags().GetString("assignee")
	assigneeChanged := cmd.Flags().Changed("assignee")
	due, _ := cmd.Flags().GetString("due")

	if summary == "" && description == "" && status == "" && !assigneeChanged && due == "" {
		if jsonMode(cmd) {
			return printJSON(JSONMutationResult{Key: key, URL: fmt.Sprintf("%s/browse/%s", client.BaseURL(), key)})
		}
		fmt.Println("Nothing to update. Specify --summary, --description, --status, --assignee, or --due.")
		fmt.Printf("URL: %s/browse/%s\n", client.BaseURL(), key)
		return nil
	}

	if summary != "" || description != "" || due != "" {
		if err := client.UpdateIssue(key, summary, description, due); err != nil {
			return err
		}
		if !jsonMode(cmd) {
			fmt.Printf("Updated issue: %s\n", key)
		}
	}

	if assigneeChanged {
		if assignee == "none" || assignee == "" {
			if err := client.AssignIssue(key, nil); err != nil {
				return err
			}
			if !jsonMode(cmd) {
				fmt.Printf("Unassigned %s\n", key)
			}
		} else {
			if err := client.AssignIssue(key, &assignee); err != nil {
				return err
			}
			if !jsonMode(cmd) {
				fmt.Printf("Assigned %s to %s\n", key, assignee)
			}
		}
	}

	if status != "" {
		if err := client.TransitionIssue(key, status); err != nil {
			return err
		}
		if !jsonMode(cmd) {
			fmt.Printf("Transitioned %s to %q\n", key, status)
		}
	}

	if jsonMode(cmd) {
		return printJSON(JSONMutationResult{
			Key: key,
			URL: fmt.Sprintf("%s/browse/%s", client.BaseURL(), key),
		})
	}

	fmt.Printf("URL: %s/browse/%s\n", client.BaseURL(), key)
	return nil
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Jira issue",
	RunE:  runIssueCreate,
}

func init() {
	issueCreateCmd.Flags().StringP("project", "p", "", "Project key (required)")
	issueCreateCmd.MarkFlagRequired("project")
	issueCreateCmd.Flags().StringP("summary", "s", "", "Issue summary (required)")
	issueCreateCmd.MarkFlagRequired("summary")
	issueCreateCmd.Flags().StringP("type", "t", "Task", "Issue type")
	issueCreateCmd.Flags().StringP("description", "d", "", "Issue description")
	issueCreateCmd.Flags().String("due", "", "Due date (YYYY-MM-DD)")
	issueCmd.AddCommand(issueCreateCmd)
}

func runIssueCreate(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	project, _ := cmd.Flags().GetString("project")
	summary, _ := cmd.Flags().GetString("summary")
	issueType, _ := cmd.Flags().GetString("type")
	description, _ := cmd.Flags().GetString("description")
	due, _ := cmd.Flags().GetString("due")

	resp, err := client.CreateIssue(project, issueType, summary, description, due)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		return printJSON(JSONMutationResult{
			Key: resp.Key,
			URL: fmt.Sprintf("%s/browse/%s", client.BaseURL(), resp.Key),
		})
	}

	fmt.Printf("Created issue: %s\n", resp.Key)
	fmt.Printf("URL: %s/browse/%s\n", client.BaseURL(), resp.Key)
	return nil
}

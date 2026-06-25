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
	issueCreateCmd.Flags().StringP("type", "t", "", "Issue type name (required; run 'atl jira issuetype list --project <key>' to discover)")
	issueCreateCmd.MarkFlagRequired("type")
	issueCreateCmd.Flags().StringP("description", "d", "", "Issue description")
	issueCreateCmd.Flags().String("due", "", "Due date (YYYY-MM-DD)")
	issueCreateCmd.Flags().String("epic", "", "Epic key to link this issue to")
	issueCreateCmd.Flags().String("parent", "", "Parent issue key (e.g. parent task for a sub-task)")
	issueCreateCmd.MarkFlagsMutuallyExclusive("epic", "parent")
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
	epic, _ := cmd.Flags().GetString("epic")
	parent, _ := cmd.Flags().GetString("parent")

	parentKey := epic
	if parent != "" {
		parentKey = parent
	}

	resp, err := client.CreateIssue(project, issueType, summary, description, due, parentKey)
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

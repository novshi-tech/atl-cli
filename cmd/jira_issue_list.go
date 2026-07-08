package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "Search for issues using JQL",
	RunE:  runIssueList,
}

func init() {
	issueListCmd.Flags().String("jql", "", "JQL query string (required)")
	issueListCmd.Flags().Int("max", 50, "Maximum number of results")
	issueCmd.AddCommand(issueListCmd)
}

func runIssueList(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	jql, _ := cmd.Flags().GetString("jql")
	max, _ := cmd.Flags().GetInt("max")

	if jql == "" {
		return fmt.Errorf("--jql is required")
	}

	resp, err := client.SearchIssues(jql, max)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		items := make([]JSONIssueItem, len(resp.Issues))
		for i, issue := range resp.Issues {
			assignee := ""
			if issue.Fields.Assignee != nil {
				assignee = issue.Fields.Assignee.DisplayName
			}
			items[i] = JSONIssueItem{
				Key:         issue.Key,
				Summary:     issue.Fields.Summary,
				Status:      issue.Fields.Status.Name,
				Type:        issue.Fields.IssueType.Name,
				Assignee:    assignee,
				StoryPoints: issue.Fields.StoryPoints,
			}
		}
		return printJSON(items)
	}

	if len(resp.Issues) == 0 {
		fmt.Println("No issues found.")
		return nil
	}

	fmt.Printf("Found %d issue(s):\n\n", len(resp.Issues))
	for _, issue := range resp.Issues {
		assigneeName := "Unassigned"
		if issue.Fields.Assignee != nil {
			assigneeName = issue.Fields.Assignee.DisplayName
		}
		fmt.Printf("%-12s  %-15s  %-6s  %-20s  %s\n",
			issue.Key,
			issue.Fields.Status.Name,
			formatStoryPoints(issue.Fields.StoryPoints),
			assigneeName,
			issue.Fields.Summary,
		)
	}
	return nil
}

// formatStoryPoints renders a Story Points value for text output, trimming
// trailing zeros (e.g. "3" rather than "3.0") since half-point estimates
// (e.g. "2.5") are also valid.
func formatStoryPoints(sp *float64) string {
	if sp == nil {
		return "-"
	}
	return strconv.FormatFloat(*sp, 'f', -1, 64)
}

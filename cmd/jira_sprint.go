package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sprintCmd = &cobra.Command{
	Use:   "sprint",
	Short: "Manage Jira sprints",
}

var sprintListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sprints for a board",
	RunE:  runSprintList,
}

var sprintIssuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "List issues in a sprint",
	RunE:  runSprintIssues,
}

func init() {
	sprintListCmd.Flags().Int("board", 0, "Board ID (required)")
	sprintListCmd.MarkFlagRequired("board")
	sprintListCmd.Flags().String("state", "", "Filter by state (active, closed, future)")

	sprintIssuesCmd.Flags().Int("sprint", 0, "Sprint ID (required)")
	sprintIssuesCmd.MarkFlagRequired("sprint")

	sprintCmd.AddCommand(sprintListCmd)
	sprintCmd.AddCommand(sprintIssuesCmd)
	jiraCmd.AddCommand(sprintCmd)
}

func runSprintList(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	boardID, _ := cmd.Flags().GetInt("board")
	state, _ := cmd.Flags().GetString("state")

	resp, err := client.ListSprints(boardID, state)
	if err != nil {
		return err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No sprints found.")
		return nil
	}

	for _, s := range resp.Values {
		goal := ""
		if s.Goal != "" {
			goal = fmt.Sprintf("  (%s)", s.Goal)
		}
		fmt.Printf("%-6d  %-10s  %s%s\n", s.ID, s.State, s.Name, goal)
	}
	return nil
}

func runSprintIssues(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	sprintID, _ := cmd.Flags().GetInt("sprint")

	resp, err := client.GetSprintIssues(sprintID)
	if err != nil {
		return err
	}

	if len(resp.Issues) == 0 {
		fmt.Println("No issues found in this sprint.")
		return nil
	}

	fmt.Printf("Found %d issue(s):\n\n", resp.Total)
	for _, issue := range resp.Issues {
		assignee := "Unassigned"
		if issue.Fields.Assignee != nil {
			assignee = issue.Fields.Assignee.DisplayName
		}
		fmt.Printf("%-12s  %-15s  %-20s  %s\n",
			issue.Key,
			issue.Fields.Status.Name,
			assignee,
			issue.Fields.Summary,
		)
	}
	return nil
}

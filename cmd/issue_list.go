package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "Search for issues using JQL",
	RunE:  runIssueList,
}

func init() {
	issueListCmd.Flags().String("jql", "", "JQL query string")
	issueListCmd.Flags().StringP("project", "p", "", "Filter by project key")
	issueListCmd.Flags().String("status", "", "Filter by status")
	issueListCmd.Flags().String("assignee", "", "Filter by assignee (use 'me' for current user)")
	issueListCmd.Flags().Int("max", 50, "Maximum number of results")
	issueCmd.AddCommand(issueListCmd)
}

func runIssueList(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	jql, _ := cmd.Flags().GetString("jql")
	project, _ := cmd.Flags().GetString("project")
	status, _ := cmd.Flags().GetString("status")
	assignee, _ := cmd.Flags().GetString("assignee")
	max, _ := cmd.Flags().GetInt("max")

	if jql == "" {
		var clauses []string
		if project != "" {
			clauses = append(clauses, fmt.Sprintf("project = %s", project))
		}
		if status != "" {
			clauses = append(clauses, fmt.Sprintf("status = \"%s\"", status))
		}
		if assignee != "" {
			if assignee == "me" {
				clauses = append(clauses, "assignee = currentUser()")
			} else {
				clauses = append(clauses, fmt.Sprintf("assignee = \"%s\"", assignee))
			}
		}
		if len(clauses) == 0 {
			return fmt.Errorf("specify --jql or at least one of --project, --status, --assignee")
		}
		jql = strings.Join(clauses, " AND ") + " ORDER BY updated DESC"
	}

	resp, err := client.SearchIssues(jql, max)
	if err != nil {
		return err
	}

	if len(resp.Issues) == 0 {
		fmt.Println("No issues found.")
		return nil
	}

	fmt.Printf("Found %d issue(s):\n\n", resp.Total)
	for _, issue := range resp.Issues {
		assigneeName := "Unassigned"
		if issue.Fields.Assignee != nil {
			assigneeName = issue.Fields.Assignee.DisplayName
		}
		fmt.Printf("%-12s  %-15s  %-20s  %s\n",
			issue.Key,
			issue.Fields.Status.Name,
			assigneeName,
			issue.Fields.Summary,
		)
	}
	return nil
}

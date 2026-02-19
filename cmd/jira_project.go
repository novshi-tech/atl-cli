package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage Jira projects",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Jira projects",
	RunE:  runProjectList,
}

func init() {
	projectListCmd.Flags().StringP("query", "q", "", "Filter projects by name")
	projectListCmd.Flags().Int("max", 50, "Maximum number of results")

	projectCmd.AddCommand(projectListCmd)
	jiraCmd.AddCommand(projectCmd)
}

func runProjectList(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	query, _ := cmd.Flags().GetString("query")
	max, _ := cmd.Flags().GetInt("max")

	resp, err := client.ListProjects(query, max)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		items := make([]JSONProjectItem, len(resp.Values))
		for i, p := range resp.Values {
			items[i] = JSONProjectItem{Key: p.Key, Name: p.Name, Type: p.ProjectTypeKey}
		}
		return printJSON(items)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No projects found.")
		return nil
	}

	fmt.Printf("Found %d project(s):\n\n", len(resp.Values))
	for _, p := range resp.Values {
		fmt.Printf("%-12s  %-40s  %s\n", p.Key, p.Name, p.ProjectTypeKey)
	}
	return nil
}

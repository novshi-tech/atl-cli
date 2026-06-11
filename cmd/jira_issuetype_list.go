package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueTypeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available issue types",
	RunE:  runIssueTypeList,
}

func init() {
	issueTypeListCmd.Flags().StringP("project", "p", "", "Filter by project ID")
	issueTypeCmd.AddCommand(issueTypeListCmd)
}

func runIssueTypeList(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	project, _ := cmd.Flags().GetString("project")

	types, err := client.GetIssueTypes(project)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		items := make([]JSONIssueTypeItem, len(types))
		for i, t := range types {
			items[i] = JSONIssueTypeItem{
				ID:          t.ID,
				Name:        t.Name,
				Description: t.Description,
				Subtask:     t.Subtask,
			}
		}
		return printJSON(items)
	}

	if len(types) == 0 {
		fmt.Println("No issue types found.")
		return nil
	}

	for _, t := range types {
		subtask := ""
		if t.Subtask {
			subtask = " [subtask]"
		}
		fmt.Printf("%-20s  %s%s\n", t.Name, t.Description, subtask)
	}
	return nil
}

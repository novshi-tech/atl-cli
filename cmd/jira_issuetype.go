package cmd

import "github.com/spf13/cobra"

var issueTypeCmd = &cobra.Command{
	Use:   "issuetype",
	Short: "Manage Jira issue types",
}

func init() {
	jiraCmd.AddCommand(issueTypeCmd)
}

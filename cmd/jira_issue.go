package cmd

import "github.com/spf13/cobra"

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage Jira issues",
}

func init() {
	jiraCmd.AddCommand(issueCmd)
}

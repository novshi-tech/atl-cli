package cmd

import "github.com/spf13/cobra"

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Manage Jira Cloud resources",
}

func init() {
	jiraCmd.PersistentFlags().String("site", "", "Site alias to use (defaults to the default site)")
	rootCmd.AddCommand(jiraCmd)
}

package cmd

import "github.com/spf13/cobra"

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage Jira issues",
}

func init() {
	issueCmd.PersistentFlags().String("site", "", "Site alias to use (defaults to the default site)")
	rootCmd.AddCommand(issueCmd)
}

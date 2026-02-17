package cmd

import "github.com/spf13/cobra"

var bbPRCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage Bitbucket pull requests",
}

func init() {
	bitbucketCmd.AddCommand(bbPRCmd)
}

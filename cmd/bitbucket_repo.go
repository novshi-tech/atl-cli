package cmd

import "github.com/spf13/cobra"

var bbRepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage Bitbucket repositories",
}

func init() {
	bitbucketCmd.AddCommand(bbRepoCmd)
}

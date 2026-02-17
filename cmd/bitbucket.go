package cmd

import "github.com/spf13/cobra"

var bitbucketCmd = &cobra.Command{
	Use:   "bitbucket",
	Short: "Manage Bitbucket Cloud resources",
}

func init() {
	bitbucketCmd.PersistentFlags().String("site", "", "Site alias to use (defaults to the default site)")
	rootCmd.AddCommand(bitbucketCmd)
}

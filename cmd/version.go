package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version of atl",
	Run: func(cmd *cobra.Command, args []string) {
		v := currentBuildVersion()
		if v == "" || v == "(devel)" {
			v = "dev"
		}
		fmt.Println(v)
	},
}

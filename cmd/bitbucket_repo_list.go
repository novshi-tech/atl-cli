package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bbRepoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories in a workspace",
	RunE:  runBBRepoList,
}

func init() {
	bbRepoListCmd.Flags().String("workspace", "", "Workspace slug (required)")
	bbRepoListCmd.MarkFlagRequired("workspace")
	bbRepoListCmd.Flags().Int("max", 25, "Maximum number of results")
	bbRepoCmd.AddCommand(bbRepoListCmd)
}

func runBBRepoList(cmd *cobra.Command, args []string) error {
	client, err := newBitbucketClient(cmd)
	if err != nil {
		return err
	}

	workspace, _ := cmd.Flags().GetString("workspace")
	max, _ := cmd.Flags().GetInt("max")

	resp, err := client.ListRepositories(workspace, 1, max)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		items := make([]JSONRepoItem, len(resp.Values))
		for i, r := range resp.Values {
			items[i] = JSONRepoItem{
				Slug:      r.Slug,
				Name:      r.Name,
				Language:  r.Language,
				IsPrivate: r.IsPrivate,
			}
		}
		return printJSON(items)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No repositories found.")
		return nil
	}

	fmt.Printf("Found %d repositor(ies):\n\n", len(resp.Values))
	for _, r := range resp.Values {
		private := "public"
		if r.IsPrivate {
			private = "private"
		}
		fmt.Printf("%-30s  %-20s  %-10s  %s\n", r.Slug, r.Name, r.Language, private)
	}
	return nil
}

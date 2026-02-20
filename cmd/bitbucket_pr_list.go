package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bbPRListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests in a repository",
	RunE:  runBBPRList,
}

func init() {
	bbPRListCmd.Flags().String("workspace", "", "Workspace slug (required)")
	bbPRListCmd.MarkFlagRequired("workspace")
	bbPRListCmd.Flags().String("repo", "", "Repository slug (required)")
	bbPRListCmd.MarkFlagRequired("repo")
	bbPRListCmd.Flags().String("state", "", "Filter by state: OPEN, MERGED, DECLINED, SUPERSEDED (default: OPEN)")
	bbPRListCmd.Flags().Int("max", 25, "Maximum number of results")
	bbPRCmd.AddCommand(bbPRListCmd)
}

func runBBPRList(cmd *cobra.Command, args []string) error {
	client, err := newBitbucketClient(cmd)
	if err != nil {
		return err
	}

	workspace, _ := cmd.Flags().GetString("workspace")
	repo, _ := cmd.Flags().GetString("repo")
	state, _ := cmd.Flags().GetString("state")
	max, _ := cmd.Flags().GetInt("max")

	resp, err := client.ListPullRequests(workspace, repo, state, 1, max)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		items := make([]JSONPRItem, len(resp.Values))
		for i, pr := range resp.Values {
			items[i] = JSONPRItem{
				ID:     pr.ID,
				Title:  pr.Title,
				State:  pr.State,
				Author: pr.Author.DisplayName,
				Source: pr.Source.Branch.Name,
				Dest:   pr.Destination.Branch.Name,
			}
		}
		return printJSON(items)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No pull requests found.")
		return nil
	}

	fmt.Printf("Found %d pull request(s):\n\n", len(resp.Values))
	for _, pr := range resp.Values {
		fmt.Printf("#%-6d  %-10s  %-20s  %s→%s  %s\n",
			pr.ID, pr.State, pr.Author.DisplayName,
			pr.Source.Branch.Name, pr.Destination.Branch.Name, pr.Title)
	}
	return nil
}

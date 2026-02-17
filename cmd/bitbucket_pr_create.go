package cmd

import (
	"fmt"

	"github.com/novshi-tech/atl-cli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbPRCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new pull request",
	RunE:  runBBPRCreate,
}

func init() {
	bbPRCreateCmd.Flags().String("workspace", "", "Workspace slug (required)")
	bbPRCreateCmd.MarkFlagRequired("workspace")
	bbPRCreateCmd.Flags().String("repo", "", "Repository slug (required)")
	bbPRCreateCmd.MarkFlagRequired("repo")
	bbPRCreateCmd.Flags().String("title", "", "Pull request title (required)")
	bbPRCreateCmd.MarkFlagRequired("title")
	bbPRCreateCmd.Flags().String("source", "", "Source branch (required)")
	bbPRCreateCmd.MarkFlagRequired("source")
	bbPRCreateCmd.Flags().String("dest", "", "Destination branch (defaults to repo main branch)")
	bbPRCreateCmd.Flags().StringP("description", "d", "", "Pull request description")
	bbPRCmd.AddCommand(bbPRCreateCmd)
}

func runBBPRCreate(cmd *cobra.Command, args []string) error {
	client, err := newBitbucketClient(cmd)
	if err != nil {
		return err
	}

	workspace, _ := cmd.Flags().GetString("workspace")
	repo, _ := cmd.Flags().GetString("repo")
	title, _ := cmd.Flags().GetString("title")
	source, _ := cmd.Flags().GetString("source")
	dest, _ := cmd.Flags().GetString("dest")
	description, _ := cmd.Flags().GetString("description")

	req := bitbucket.CreatePRRequest{
		Title:  title,
		Source: bitbucket.CreatePRRef{Branch: bitbucket.CreatePRBranch{Name: source}},
	}
	if dest != "" {
		req.Destination = &bitbucket.CreatePRRef{Branch: bitbucket.CreatePRBranch{Name: dest}}
	}
	if description != "" {
		req.Description = description
	}

	pr, err := client.CreatePullRequest(workspace, repo, req)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		return printJSON(JSONMutationResult{
			Key: fmt.Sprintf("%d", pr.ID),
			URL: pr.Links.HTML.Href,
		})
	}

	fmt.Printf("Created pull request: #%d\n", pr.ID)
	fmt.Printf("URL: %s\n", pr.Links.HTML.Href)
	return nil
}

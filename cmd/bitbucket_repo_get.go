package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bbRepoGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a repository",
	RunE:  runBBRepoGet,
}

func init() {
	bbRepoGetCmd.Flags().String("workspace", "", "Workspace slug (required)")
	bbRepoGetCmd.MarkFlagRequired("workspace")
	bbRepoGetCmd.Flags().String("repo", "", "Repository slug (required)")
	bbRepoGetCmd.MarkFlagRequired("repo")
	bbRepoCmd.AddCommand(bbRepoGetCmd)
}

func runBBRepoGet(cmd *cobra.Command, args []string) error {
	client, err := newBitbucketClient(cmd)
	if err != nil {
		return err
	}

	workspace, _ := cmd.Flags().GetString("workspace")
	repo, _ := cmd.Flags().GetString("repo")

	r, err := client.GetRepository(workspace, repo)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		mainbranch := ""
		if r.MainBranch != nil {
			mainbranch = r.MainBranch.Name
		}
		return printJSON(JSONRepoDetail{
			Slug:        r.Slug,
			Name:        r.Name,
			FullName:    r.FullName,
			Description: r.Description,
			Language:    r.Language,
			IsPrivate:   r.IsPrivate,
			MainBranch:  mainbranch,
			UpdatedOn:   r.UpdatedOn,
		})
	}

	fmt.Printf("Slug:         %s\n", r.Slug)
	fmt.Printf("Name:         %s\n", r.Name)
	fmt.Printf("Full Name:    %s\n", r.FullName)
	fmt.Printf("Description:  %s\n", r.Description)
	fmt.Printf("Language:     %s\n", r.Language)
	private := "No"
	if r.IsPrivate {
		private = "Yes"
	}
	fmt.Printf("Private:      %s\n", private)
	if r.MainBranch != nil {
		fmt.Printf("Main Branch:  %s\n", r.MainBranch.Name)
	}
	fmt.Printf("Updated:      %s\n", r.UpdatedOn)
	return nil
}

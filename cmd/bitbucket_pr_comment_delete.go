package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bbPRCommentDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a comment on a pull request",
	RunE:  runBBPRCommentDelete,
}

func init() {
	bbPRCommentDeleteCmd.Flags().String("workspace", "", "Workspace slug")
	bbPRCommentDeleteCmd.Flags().String("repo", "", "Repository slug (required)")
	bbPRCommentDeleteCmd.MarkFlagRequired("repo")
	bbPRCommentDeleteCmd.Flags().Int("pr", 0, "Pull request ID (required)")
	bbPRCommentDeleteCmd.MarkFlagRequired("pr")
	bbPRCommentDeleteCmd.Flags().Int("id", 0, "Comment ID (required)")
	bbPRCommentDeleteCmd.MarkFlagRequired("id")
	bbPRCommentCmd.AddCommand(bbPRCommentDeleteCmd)
}

func runBBPRCommentDelete(cmd *cobra.Command, args []string) error {
	client, err := newBitbucketClient(cmd)
	if err != nil {
		return err
	}

	workspace, err := resolveBBWorkspace(cmd)
	if err != nil {
		return err
	}
	repo, _ := cmd.Flags().GetString("repo")
	prID, _ := cmd.Flags().GetInt("pr")
	commentID, _ := cmd.Flags().GetInt("id")

	if err := client.DeletePRComment(workspace, repo, prID, commentID); err != nil {
		return err
	}

	prURL := fmt.Sprintf("https://bitbucket.org/%s/%s/pull-requests/%d", workspace, repo, prID)

	if jsonMode(cmd) {
		return printJSON(JSONMutationResult{
			Key: fmt.Sprintf("%d", commentID),
			URL: prURL,
		})
	}

	fmt.Printf("Comment #%d deleted from pull request #%d\n", commentID, prID)
	fmt.Printf("URL: %s\n", prURL)
	return nil
}

package cmd

import (
	"fmt"

	"github.com/novshi-tech/atl-cli/internal/bitbucket"
	"github.com/spf13/cobra"
)

var bbPRCommentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a comment on a pull request",
	RunE:  runBBPRCommentCreate,
}

func init() {
	bbPRCommentCreateCmd.Flags().String("workspace", "", "Workspace slug")
	bbPRCommentCreateCmd.Flags().String("repo", "", "Repository slug (required)")
	bbPRCommentCreateCmd.MarkFlagRequired("repo")
	bbPRCommentCreateCmd.Flags().Int("pr", 0, "Pull request ID (required)")
	bbPRCommentCreateCmd.MarkFlagRequired("pr")
	bbPRCommentCreateCmd.Flags().StringP("body", "b", "", "Comment body (required)")
	bbPRCommentCreateCmd.MarkFlagRequired("body")
	bbPRCommentCreateCmd.Flags().String("path", "", "File path for inline comment (ignored when --parent is set; inherited from the parent comment)")
	bbPRCommentCreateCmd.Flags().Int("line", 0, "Line number for inline comment (ignored when --parent is set; inherited from the parent comment)")
	bbPRCommentCreateCmd.Flags().Int("parent", 0, "Parent comment ID (for replies); inline position, if any, is inherited from the parent")
	bbPRCommentCmd.AddCommand(bbPRCommentCreateCmd)
}

func runBBPRCommentCreate(cmd *cobra.Command, args []string) error {
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
	body, _ := cmd.Flags().GetString("body")
	path, _ := cmd.Flags().GetString("path")
	line, _ := cmd.Flags().GetInt("line")
	parentID, _ := cmd.Flags().GetInt("parent")

	req := bitbucket.CreatePRCommentRequest{
		Content: bitbucket.PRCommentContent{Raw: body},
	}

	if parentID > 0 {
		parent, err := client.GetPRComment(workspace, repo, prID, parentID)
		if err != nil {
			return fmt.Errorf("resolving parent comment #%d: %w", parentID, err)
		}
		req.Parent = &bitbucket.PRCommentParent{ID: parentID}
		// Replies must anchor to the exact from/to the parent used; reconstructing
		// it from --path/--line can mismatch (e.g. a parent spanning from+to) and
		// gets rejected by the API, so inherit it verbatim instead.
		req.Inline = parent.Inline
	} else if path != "" {
		inline := &bitbucket.PRInline{Path: path}
		if line > 0 {
			inline.To = &line
		}
		req.Inline = inline
	}

	comment, err := client.CreatePRComment(workspace, repo, prID, req)
	if err != nil {
		return err
	}

	prURL := fmt.Sprintf("https://bitbucket.org/%s/%s/pull-requests/%d", workspace, repo, prID)

	if jsonMode(cmd) {
		return printJSON(JSONMutationResult{
			Key: fmt.Sprintf("%d", comment.ID),
			URL: prURL,
		})
	}

	fmt.Printf("Comment added to pull request #%d\n", prID)
	fmt.Printf("URL: %s\n", prURL)
	return nil
}

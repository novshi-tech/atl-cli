package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bbPRCommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "List comments on a pull request",
	RunE:  runBBPRComment,
}

func init() {
	bbPRCommentCmd.Flags().String("workspace", "", "Workspace slug (required)")
	bbPRCommentCmd.MarkFlagRequired("workspace")
	bbPRCommentCmd.Flags().String("repo", "", "Repository slug (required)")
	bbPRCommentCmd.MarkFlagRequired("repo")
	bbPRCommentCmd.Flags().Int("pr", 0, "Pull request ID (required)")
	bbPRCommentCmd.MarkFlagRequired("pr")
	bbPRCmd.AddCommand(bbPRCommentCmd)
}

func runBBPRComment(cmd *cobra.Command, args []string) error {
	client, err := newBitbucketClient(cmd)
	if err != nil {
		return err
	}

	workspace, _ := cmd.Flags().GetString("workspace")
	repo, _ := cmd.Flags().GetString("repo")
	prID, _ := cmd.Flags().GetInt("pr")

	resp, err := client.ListPRComments(workspace, repo, prID)
	if err != nil {
		return err
	}

	// Filter out inline comments (code review comments)
	var comments []struct {
		Author  string
		Created string
		Body    string
	}
	for _, c := range resp.Values {
		if c.Inline != nil {
			continue
		}
		comments = append(comments, struct {
			Author  string
			Created string
			Body    string
		}{
			Author:  c.User.DisplayName,
			Created: c.CreatedOn,
			Body:    c.Content.Raw,
		})
	}

	if jsonMode(cmd) {
		items := make([]JSONCommentItem, len(comments))
		for i, c := range comments {
			items[i] = JSONCommentItem{
				Author:  c.Author,
				Created: c.Created,
				Body:    c.Body,
			}
		}
		return printJSON(items)
	}

	if len(comments) == 0 {
		fmt.Println("No comments found.")
		return nil
	}

	fmt.Printf("Found %d comment(s):\n\n", len(comments))
	for _, c := range comments {
		fmt.Printf("[%s] %s:\n%s\n\n", c.Created, c.Author, c.Body)
	}
	return nil
}

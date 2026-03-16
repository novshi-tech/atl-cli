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
	bbPRCommentCmd.Flags().String("workspace", "", "Workspace slug")
	bbPRCommentCmd.Flags().String("repo", "", "Repository slug (required)")
	bbPRCommentCmd.MarkFlagRequired("repo")
	bbPRCommentCmd.Flags().Int("pr", 0, "Pull request ID (required)")
	bbPRCommentCmd.MarkFlagRequired("pr")
	bbPRCommentCmd.Flags().Bool("inline", false, "Include inline code review comments")
	bbPRCmd.AddCommand(bbPRCommentCmd)
}

func runBBPRComment(cmd *cobra.Command, args []string) error {
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
	showInline, _ := cmd.Flags().GetBool("inline")

	resp, err := client.ListPRComments(workspace, repo, prID)
	if err != nil {
		return err
	}

	var comments []JSONCommentItem
	var inlineComments []JSONInlineCommentItem
	for _, c := range resp.Values {
		if c.Inline != nil {
			if showInline {
				inlineComments = append(inlineComments, JSONInlineCommentItem{
					Author:  c.User.DisplayName,
					Created: c.CreatedOn,
					Path:    c.Inline.Path,
					From:    c.Inline.From,
					To:      c.Inline.To,
					Body:    c.Content.Raw,
				})
			}
			continue
		}
		comments = append(comments, JSONCommentItem{
			Author:  c.User.DisplayName,
			Created: c.CreatedOn,
			Body:    c.Content.Raw,
		})
	}

	if jsonMode(cmd) {
		if showInline {
			return printJSON(struct {
				Comments       []JSONCommentItem       `json:"comments"`
				InlineComments []JSONInlineCommentItem `json:"inline_comments"`
			}{comments, inlineComments})
		}
		return printJSON(comments)
	}

	if len(comments) == 0 && (!showInline || len(inlineComments) == 0) {
		fmt.Println("No comments found.")
		return nil
	}

	if len(comments) > 0 {
		fmt.Printf("Found %d comment(s):\n\n", len(comments))
		for _, c := range comments {
			fmt.Printf("[%s] %s:\n%s\n\n", c.Created, c.Author, c.Body)
		}
	}

	if showInline && len(inlineComments) > 0 {
		fmt.Printf("Found %d inline comment(s):\n\n", len(inlineComments))
		for _, c := range inlineComments {
			lineRef := prInlineLineRef(c.From, c.To)
			fmt.Printf("[%s] %s on %s%s:\n%s\n\n", c.Created, c.Author, c.Path, lineRef, c.Body)
		}
	}

	return nil
}

func prInlineLineRef(from, to *int) string {
	if from != nil && to != nil {
		return fmt.Sprintf(" (lines %d-%d)", *from, *to)
	}
	if to != nil {
		return fmt.Sprintf(" (line %d)", *to)
	}
	return ""
}

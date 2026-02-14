package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueCommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Add a comment to a Jira issue",
	RunE:  runIssueComment,
}

func init() {
	issueCommentCmd.Flags().StringP("key", "k", "", "Issue key (required)")
	issueCommentCmd.MarkFlagRequired("key")
	issueCommentCmd.Flags().StringP("body", "b", "", "Comment body (required)")
	issueCommentCmd.MarkFlagRequired("body")
	issueCmd.AddCommand(issueCommentCmd)
}

func runIssueComment(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	key, _ := cmd.Flags().GetString("key")
	body, _ := cmd.Flags().GetString("body")

	if err := client.AddComment(key, body); err != nil {
		return err
	}

	if jsonMode(cmd) {
		return printJSON(JSONMutationResult{
			Key: key,
			URL: fmt.Sprintf("%s/browse/%s", client.BaseURL(), key),
		})
	}

	fmt.Printf("Comment added to %s\n", key)
	fmt.Printf("URL: %s/browse/%s\n", client.BaseURL(), key)
	return nil
}

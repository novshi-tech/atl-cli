package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/novshi-tech/atl-cli/internal/adf"
)

var issueViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View details of a Jira issue",
	RunE:  runIssueView,
}

func init() {
	issueViewCmd.Flags().StringP("key", "k", "", "Issue key (required)")
	issueViewCmd.MarkFlagRequired("key")
	issueCmd.AddCommand(issueViewCmd)
}

func runIssueView(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	key, _ := cmd.Flags().GetString("key")

	issue, err := client.GetIssue(key)
	if err != nil {
		return err
	}

	assignee := "Unassigned"
	if issue.Fields.Assignee != nil {
		assignee = issue.Fields.Assignee.DisplayName
	}

	if jsonMode(cmd) {
		detail := JSONIssueDetail{
			Key:      issue.Key,
			Summary:  issue.Fields.Summary,
			Status:   issue.Fields.Status.Name,
			Type:     issue.Fields.IssueType.Name,
			Assignee: assignee,
			URL:      fmt.Sprintf("%s/browse/%s", client.BaseURL(), issue.Key),
		}
		if issue.Fields.Description != nil {
			detail.Description = adfToText(issue.Fields.Description)
		}
		if issue.Fields.Comment != nil {
			for _, c := range issue.Fields.Comment.Comments {
				detail.Comments = append(detail.Comments, JSONCommentItem{
					Author:  c.Author.DisplayName,
					Created: c.Created,
					Body:    adfToText(&c.Body),
				})
			}
		}
		return printJSON(detail)
	}

	fmt.Printf("Key:       %s\n", issue.Key)
	fmt.Printf("Summary:   %s\n", issue.Fields.Summary)
	fmt.Printf("Status:    %s\n", issue.Fields.Status.Name)
	fmt.Printf("Type:      %s\n", issue.Fields.IssueType.Name)
	fmt.Printf("Assignee:  %s\n", assignee)
	fmt.Printf("URL:       %s/browse/%s\n", client.BaseURL(), issue.Key)

	if issue.Fields.Description != nil {
		fmt.Printf("\n--- Description ---\n%s\n", adfToText(issue.Fields.Description))
	}

	if issue.Fields.Comment != nil && len(issue.Fields.Comment.Comments) > 0 {
		fmt.Printf("\n--- Comments (%d) ---\n", len(issue.Fields.Comment.Comments))
		for _, c := range issue.Fields.Comment.Comments {
			fmt.Printf("\n[%s] %s:\n%s\n", c.Created, c.Author.DisplayName, adfToText(&c.Body))
		}
	}

	return nil
}

// adfToText extracts plain text from an ADF node.
func adfToText(node *adf.Node) string {
	if node == nil {
		return ""
	}
	if node.Type == "text" {
		return node.Text
	}
	var parts []string
	for i := range node.Content {
		t := adfToText(&node.Content[i])
		if t != "" {
			parts = append(parts, t)
		}
	}
	sep := ""
	if node.Type == "doc" || node.Type == "paragraph" || node.Type == "bulletList" || node.Type == "orderedList" {
		sep = "\n"
	}
	return strings.Join(parts, sep)
}

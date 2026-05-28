package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueAttachmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List attachments on a Jira issue",
	RunE:  runIssueAttachmentList,
}

func init() {
	issueAttachmentListCmd.Flags().StringP("key", "k", "", "Issue key (required)")
	issueAttachmentListCmd.MarkFlagRequired("key")
	issueAttachmentCmd.AddCommand(issueAttachmentListCmd)
}

func runIssueAttachmentList(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	key, _ := cmd.Flags().GetString("key")

	attachments, err := client.GetAttachments(key)
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		items := make([]JSONAttachmentItem, len(attachments))
		for i, a := range attachments {
			author := ""
			if a.Author != nil {
				author = a.Author.DisplayName
			}
			items[i] = JSONAttachmentItem{
				ID:       a.ID,
				Filename: a.Filename,
				Size:     a.Size,
				MimeType: a.MimeType,
				Author:   author,
				Created:  a.Created,
				Content:  a.Content,
			}
		}
		return printJSON(items)
	}

	if len(attachments) == 0 {
		fmt.Printf("No attachments on %s.\n", key)
		return nil
	}

	fmt.Printf("Found %d attachment(s) on %s:\n\n", len(attachments), key)
	for _, a := range attachments {
		author := "Unknown"
		if a.Author != nil {
			author = a.Author.DisplayName
		}
		fmt.Printf("%-12s  %-10s  %-20s  %s\n", a.ID, formatSize(a.Size), author, a.Filename)
	}
	return nil
}

func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case size >= GB:
		return fmt.Sprintf("%.1fGB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.1fMB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.1fKB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%dB", size)
	}
}

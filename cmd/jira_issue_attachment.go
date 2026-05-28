package cmd

import "github.com/spf13/cobra"

var issueAttachmentCmd = &cobra.Command{
	Use:   "attachment",
	Short: "Manage Jira issue attachments",
}

func init() {
	issueCmd.AddCommand(issueAttachmentCmd)
}

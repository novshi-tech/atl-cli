package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var issueAttachmentDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a Jira issue attachment",
	RunE:  runIssueAttachmentDownload,
}

func init() {
	issueAttachmentDownloadCmd.Flags().String("id", "", "Attachment ID (required unless --key/--filename are used)")
	issueAttachmentDownloadCmd.Flags().StringP("key", "k", "", "Issue key (used with --filename to look up the attachment)")
	issueAttachmentDownloadCmd.Flags().String("filename", "", "Filename to match on the issue (used with --key)")
	issueAttachmentDownloadCmd.Flags().StringP("output", "o", "", "Output file path (default: server-provided filename in current directory; use '-' for stdout)")
	issueAttachmentCmd.AddCommand(issueAttachmentDownloadCmd)
}

func runIssueAttachmentDownload(cmd *cobra.Command, args []string) error {
	client, err := newJiraClient(cmd)
	if err != nil {
		return err
	}

	id, _ := cmd.Flags().GetString("id")
	key, _ := cmd.Flags().GetString("key")
	filename, _ := cmd.Flags().GetString("filename")
	output, _ := cmd.Flags().GetString("output")

	resolvedFilename := ""
	if id == "" {
		if key == "" || filename == "" {
			return fmt.Errorf("either --id or both --key and --filename must be specified")
		}
		attachments, err := client.GetAttachments(key)
		if err != nil {
			return err
		}
		for _, a := range attachments {
			if a.Filename == filename {
				id = a.ID
				resolvedFilename = a.Filename
				break
			}
		}
		if id == "" {
			return fmt.Errorf("no attachment named %q found on %s", filename, key)
		}
	}

	var buf bytes.Buffer
	serverFilename, err := client.DownloadAttachment(id, &buf)
	if err != nil {
		return err
	}
	if resolvedFilename == "" {
		resolvedFilename = serverFilename
	}

	if output == "-" {
		if _, err := buf.WriteTo(os.Stdout); err != nil {
			return fmt.Errorf("writing to stdout: %w", err)
		}
		return nil
	}

	outPath := output
	if outPath == "" {
		if resolvedFilename == "" {
			resolvedFilename = "attachment-" + id
		}
		outPath = resolvedFilename
	} else if info, err := os.Stat(outPath); err == nil && info.IsDir() {
		name := resolvedFilename
		if name == "" {
			name = "attachment-" + id
		}
		outPath = filepath.Join(outPath, name)
	}

	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	if jsonMode(cmd) {
		return printJSON(JSONAttachmentDownload{
			ID:       id,
			Filename: resolvedFilename,
			Path:     outPath,
			Size:     int64(buf.Len()),
		})
	}

	fmt.Printf("Downloaded attachment %s (%s) to %s\n", id, formatSize(int64(buf.Len())), outPath)
	return nil
}

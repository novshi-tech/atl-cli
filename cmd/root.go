package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"novshi-tech.com/jira-cli/internal/auth"
	"novshi-tech.com/jira-cli/internal/jira"
)

var rootCmd = &cobra.Command{
	Use:   "jira-cli",
	Short: "A CLI tool for Jira Cloud",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// newJiraClient resolves the site alias from the --site flag (or default) and returns a Jira client.
func newJiraClient(cmd *cobra.Command) (*jira.Client, error) {
	store, err := auth.NewStore()
	if err != nil {
		return nil, err
	}

	site, _ := cmd.Flags().GetString("site")
	if site == "" {
		site, err = auth.GetDefaultSite(store)
		if err != nil {
			return nil, fmt.Errorf("no --site specified and no default site configured; run 'jira-cli configure --site <name>' first")
		}
	}

	return jira.NewClientFromStore(store, site)
}

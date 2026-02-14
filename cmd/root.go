package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"novshi-tech.com/atl/internal/auth"
	"novshi-tech.com/atl/internal/jira"
)

var rootCmd = &cobra.Command{
	Use:   "atl",
	Short: "A CLI for Atlassian Cloud products",
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func jsonMode(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("json")
	return v
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
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
			return nil, fmt.Errorf("no --site specified and no default site configured; run 'atl configure --site <name>' first")
		}
	}

	return jira.NewClientFromStore(store, site)
}

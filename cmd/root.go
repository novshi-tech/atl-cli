package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/novshi-tech/atl-cli/internal/auth"
	"github.com/novshi-tech/atl-cli/internal/bitbucket"
	"github.com/novshi-tech/atl-cli/internal/jira"
)

var rootCmd = &cobra.Command{
	Use:   "atl",
	Short: "A CLI for Atlassian Cloud products",
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
}

func Execute() {
	checkAndAutoUpdate()
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
		site = os.Getenv("ATL_SITE")
	}
	if site == "" {
		site, err = auth.GetDefaultSite(store)
		if err != nil {
			return nil, fmt.Errorf("no --site specified and no default site configured; run 'atl configure --site <name>' first")
		}
	}

	return jira.NewClientFromStore(store, site)
}

// newBitbucketClient resolves the site alias from the --site flag (or default) and returns a Bitbucket client.
func newBitbucketClient(cmd *cobra.Command) (*bitbucket.Client, error) {
	store, err := auth.NewStore()
	if err != nil {
		return nil, err
	}

	site, _ := cmd.Flags().GetString("site")
	if site == "" {
		site = os.Getenv("ATL_SITE")
	}
	if site == "" {
		site, err = auth.GetDefaultSite(store)
		if err != nil {
			return nil, fmt.Errorf("no --site specified and no default site configured; run 'atl configure --site <name>' first")
		}
	}

	return bitbucket.NewClientFromStore(store, site)
}

// resolveBBWorkspace resolves the workspace from the --workspace flag and site configuration.
// If both are set, they must match. If neither is set, an error is returned.
func resolveBBWorkspace(cmd *cobra.Command) (string, error) {
	flagWS, _ := cmd.Flags().GetString("workspace")

	store, err := auth.NewStore()
	if err != nil {
		return "", err
	}

	site, _ := cmd.Flags().GetString("site")
	if site == "" {
		site = os.Getenv("ATL_SITE")
	}
	if site == "" {
		site, err = auth.GetDefaultSite(store)
		if err != nil {
			return "", fmt.Errorf("no --site specified and no default site configured; run 'atl configure --site <name>' first")
		}
	}

	creds, err := auth.LoadSite(store, site)
	if err != nil {
		return "", err
	}
	savedWS := creds.BBWorkspace

	switch {
	case flagWS != "" && savedWS != "":
		if flagWS != savedWS {
			return "", fmt.Errorf("workspace mismatch: --workspace %q does not match configured workspace %q for site %q", flagWS, savedWS, site)
		}
		return savedWS, nil
	case flagWS != "":
		return flagWS, nil
	case savedWS != "":
		return savedWS, nil
	default:
		return "", fmt.Errorf("no workspace specified; use --workspace flag or configure it with 'atl configure --site <name>'")
	}
}

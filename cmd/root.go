package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/novshi-tech/atl-cli/internal/auth"
	"github.com/novshi-tech/atl-cli/internal/bitbucket"
	"github.com/novshi-tech/atl-cli/internal/jira"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "atl",
	Short: "A CLI for Atlassian Cloud products",
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
}

func Execute() {
	checkForUpdate()
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

// loadSiteCredentials returns the credentials a command should use.
//
// Credentials supplied through the environment (ATL_BASE_URL / ATL_EMAIL /
// ATL_API_TOKEN, see auth.LoadFromEnv) win outright and never touch the
// credential store; otherwise the site alias is resolved from --site, then
// ATL_SITE, then the stored default site, and loaded from the store.
func loadSiteCredentials(cmd *cobra.Command) (auth.SiteCredentials, error) {
	if creds, ok, err := auth.LoadFromEnv(); ok || err != nil {
		return creds, err
	}

	store, err := auth.NewStore()
	if err != nil {
		return auth.SiteCredentials{}, err
	}

	site, err := resolveSiteAlias(cmd, store)
	if err != nil {
		return auth.SiteCredentials{}, err
	}

	return auth.LoadSite(store, site)
}

// resolveSiteAlias picks the site alias from --site, ATL_SITE, or the stored default.
func resolveSiteAlias(cmd *cobra.Command, store auth.CredentialStore) (string, error) {
	site, _ := cmd.Flags().GetString("site")
	if site == "" {
		site = os.Getenv("ATL_SITE")
	}
	if site == "" {
		var err error
		site, err = auth.GetDefaultSite(store)
		if err != nil {
			return "", fmt.Errorf("no --site specified and no default site configured; run 'atl configure --site <name>' first")
		}
	}
	return site, nil
}

// newJiraClient returns a Jira client for the credentials selected by loadSiteCredentials.
func newJiraClient(cmd *cobra.Command) (*jira.Client, error) {
	creds, err := loadSiteCredentials(cmd)
	if err != nil {
		return nil, err
	}
	return jira.NewClient(creds), nil
}

// newBitbucketClient returns a Bitbucket client for the credentials selected by loadSiteCredentials.
func newBitbucketClient(cmd *cobra.Command) (*bitbucket.Client, error) {
	creds, err := loadSiteCredentials(cmd)
	if err != nil {
		return nil, err
	}
	return bitbucket.NewClient(creds), nil
}

// resolveBBWorkspace resolves the workspace from the --workspace flag and site configuration.
// If both are set, they must match. If neither is set, an error is returned.
func resolveBBWorkspace(cmd *cobra.Command) (string, error) {
	flagWS, _ := cmd.Flags().GetString("workspace")

	creds, err := loadSiteCredentials(cmd)
	if err != nil {
		return "", err
	}
	savedWS := creds.BBWorkspace

	switch {
	case flagWS != "" && savedWS != "":
		if flagWS != savedWS {
			return "", fmt.Errorf("workspace mismatch: --workspace %q does not match configured workspace %q", flagWS, savedWS)
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

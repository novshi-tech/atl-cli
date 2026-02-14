package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"novshi-tech.com/atl/internal/auth"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure authentication for an Atlassian site",
	RunE:  runConfigure,
}

func init() {
	configureCmd.Flags().StringP("site", "s", "", "Site alias name (required)")
	configureCmd.MarkFlagRequired("site")
	configureCmd.Flags().Bool("default", false, "Set this site as the default")
	rootCmd.AddCommand(configureCmd)
}

func runConfigure(cmd *cobra.Command, args []string) error {
	alias, _ := cmd.Flags().GetString("site")
	setDefault, _ := cmd.Flags().GetBool("default")

	store, err := auth.NewStore()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Configuring site %q\n", alias)

	fmt.Print("Site URL (e.g., https://yourcompany.atlassian.net): ")
	siteURL, _ := reader.ReadString('\n')
	siteURL = strings.TrimSpace(siteURL)

	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("API Token: ")
	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("reading API token: %w", err)
	}
	fmt.Println()
	apiToken := string(tokenBytes)

	creds := auth.SiteCredentials{
		BaseURL:  siteURL,
		Email:    email,
		APIToken: apiToken,
	}

	if err := auth.SaveSite(store, alias, creds); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}

	if err := auth.AddSiteToList(store, alias); err != nil {
		return fmt.Errorf("updating site list: %w", err)
	}

	// Auto-set as default if this is the first site or --default is specified
	if setDefault {
		if err := auth.SetDefaultSite(store, alias); err != nil {
			return fmt.Errorf("setting default site: %w", err)
		}
		fmt.Printf("Site %q set as default.\n", alias)
	} else {
		sites, _ := auth.ListSites(store)
		if len(sites) == 1 {
			if err := auth.SetDefaultSite(store, alias); err != nil {
				return fmt.Errorf("setting default site: %w", err)
			}
			fmt.Printf("Site %q set as default (first site).\n", alias)
		}
	}

	fmt.Printf("Site %q configured successfully.\n", alias)
	return nil
}

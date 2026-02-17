package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/novshi-tech/atl-cli/internal/auth"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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

	// Load existing credentials so Enter preserves current values
	existing, _ := auth.LoadSite(store, alias)

	fmt.Printf("Configuring site %q\n", alias)

	siteURL := promptText(reader, "Site URL", existing.BaseURL)
	email := promptText(reader, "Email", existing.Email)
	apiToken := promptSecret("API Token", existing.APIToken)
	bbAppPassword := promptSecret("Bitbucket App Password", existing.BBAppPassword)

	creds := auth.SiteCredentials{
		BaseURL:       siteURL,
		Email:         email,
		APIToken:      apiToken,
		BBAppPassword: bbAppPassword,
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

// promptText prompts for a text value. If the user presses Enter, the existing value is kept.
func promptText(reader *bufio.Reader, label, current string) string {
	if current != "" {
		fmt.Printf("%s [%s]: ", label, current)
	} else {
		fmt.Printf("%s: ", label)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return current
	}
	return input
}

// promptSecret prompts for a secret value (input hidden). If the user presses Enter, the existing value is kept.
func promptSecret(label, current string) string {
	if current != "" {
		fmt.Printf("%s [****]: ", label)
	} else {
		fmt.Printf("%s (Enter to skip): ", label)
	}
	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil || len(b) == 0 {
		return current
	}
	return string(b)
}

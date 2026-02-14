package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"novshi-tech.com/jira-cli/internal/auth"
)

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "List configured Jira sites",
	RunE:  runSites,
}

func init() {
	rootCmd.AddCommand(sitesCmd)
}

func runSites(cmd *cobra.Command, args []string) error {
	store, err := auth.NewStore()
	if err != nil {
		return err
	}

	sites, err := auth.ListSites(store)
	if err != nil {
		return err
	}

	if len(sites) == 0 {
		fmt.Println("No sites configured. Run 'jira-cli configure --site <name>' to add one.")
		return nil
	}

	defaultSite, _ := auth.GetDefaultSite(store)

	for _, s := range sites {
		if s == defaultSite {
			fmt.Printf("  * %s (default)\n", s)
		} else {
			fmt.Printf("    %s\n", s)
		}
	}
	return nil
}

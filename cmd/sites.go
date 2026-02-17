package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/novshi-tech/atl-cli/internal/auth"
)

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "List configured Atlassian sites",
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

	defaultSite, _ := auth.GetDefaultSite(store)

	if jsonMode(cmd) {
		items := make([]JSONSiteItem, len(sites))
		for i, s := range sites {
			items[i] = JSONSiteItem{
				Name:    s,
				Default: s == defaultSite,
			}
		}
		return printJSON(items)
	}

	if len(sites) == 0 {
		fmt.Println("No sites configured. Run 'atl configure --site <name>' to add one.")
		return nil
	}

	for _, s := range sites {
		if s == defaultSite {
			fmt.Printf("  * %s (default)\n", s)
		} else {
			fmt.Printf("    %s\n", s)
		}
	}
	return nil
}
